package storage

import (
	"sync"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMemoryStore_StorePRByID tests storing and retrieving PRs by ID.
func TestMemoryStore_StorePRByID(t *testing.T) {
	store := NewMemoryStore()
	pr := models.PullRequest{
		ID:          1,
		Number:      42,
		RepoName:    "acme/platform",
		AuthorEmail: "dev@company.com",
		AuthorID:    "user-1",
		State:       models.PRStateMerged,
		Title:       "Add feature",
		BaseBranch:  "main",
		HeadBranch:  "feature-branch",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := store.StorePR(pr)
	require.NoError(t, err)

	retrieved, err := store.GetPRByID(1)
	require.NoError(t, err)
	assert.Equal(t, pr.Number, retrieved.Number)
	assert.Equal(t, pr.RepoName, retrieved.RepoName)
	assert.Equal(t, pr.AuthorEmail, retrieved.AuthorEmail)
	assert.Equal(t, pr.State, retrieved.State)
}

// TestMemoryStore_GetPRByID_NotFound tests retrieving non-existent PR.
func TestMemoryStore_GetPRByID_NotFound(t *testing.T) {
	store := NewMemoryStore()

	_, err := store.GetPRByID(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PR not found")
}

// TestMemoryStore_GetPRsByStatus tests filtering PRs by state.
func TestMemoryStore_GetPRsByStatus(t *testing.T) {
	store := NewMemoryStore()

	// Store PRs with different states
	prs := []models.PullRequest{
		{
			ID:          1,
			Number:      1,
			RepoName:    "acme/platform",
			AuthorEmail: "dev1@company.com",
			AuthorID:    "user-1",
			State:       models.PRStateOpen,
			Title:       "PR 1",
			BaseBranch:  "main",
			HeadBranch:  "feat-1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Number:      2,
			RepoName:    "acme/platform",
			AuthorEmail: "dev2@company.com",
			AuthorID:    "user-2",
			State:       models.PRStateMerged,
			Title:       "PR 2",
			BaseBranch:  "main",
			HeadBranch:  "feat-2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          3,
			Number:      3,
			RepoName:    "acme/platform",
			AuthorEmail: "dev3@company.com",
			AuthorID:    "user-3",
			State:       models.PRStateOpen,
			Title:       "PR 3",
			BaseBranch:  "main",
			HeadBranch:  "feat-3",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, pr := range prs {
		err := store.StorePR(pr)
		require.NoError(t, err)
	}

	// Get open PRs
	openPRs, err := store.GetPRsByStatus(models.PRStateOpen)
	require.NoError(t, err)
	assert.Len(t, openPRs, 2)

	// Get merged PRs
	mergedPRs, err := store.GetPRsByStatus(models.PRStateMerged)
	require.NoError(t, err)
	assert.Len(t, mergedPRs, 1)
	assert.Equal(t, 2, mergedPRs[0].Number)
}

// TestMemoryStore_GetPRsByAuthorEmail tests retrieving PRs by author email.
func TestMemoryStore_GetPRsByAuthorEmail(t *testing.T) {
	store := NewMemoryStore()

	authorEmail := "dev@company.com"

	prs := []models.PullRequest{
		{
			ID:          1,
			Number:      1,
			RepoName:    "acme/platform",
			AuthorEmail: authorEmail,
			AuthorID:    "user-1",
			State:       models.PRStateOpen,
			Title:       "PR 1",
			BaseBranch:  "main",
			HeadBranch:  "feat-1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Number:      2,
			RepoName:    "acme/platform",
			AuthorEmail: "other@company.com",
			AuthorID:    "user-2",
			State:       models.PRStateMerged,
			Title:       "PR 2",
			BaseBranch:  "main",
			HeadBranch:  "feat-2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          3,
			Number:      3,
			RepoName:    "acme/api",
			AuthorEmail: authorEmail,
			AuthorID:    "user-1",
			State:       models.PRStateMerged,
			Title:       "PR 3",
			BaseBranch:  "main",
			HeadBranch:  "feat-3",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, pr := range prs {
		err := store.StorePR(pr)
		require.NoError(t, err)
	}

	authorPRs, err := store.GetPRsByAuthorEmail(authorEmail)
	require.NoError(t, err)
	assert.Len(t, authorPRs, 2)

	// Verify both PRs belong to the author
	for _, pr := range authorPRs {
		assert.Equal(t, authorEmail, pr.AuthorEmail)
	}
}

// TestMemoryStore_GetPRsByRepoWithPagination tests paginated PR retrieval.
func TestMemoryStore_GetPRsByRepoWithPagination(t *testing.T) {
	store := NewMemoryStore()

	repoName := "acme/platform"

	// Create 25 PRs
	for i := 1; i <= 25; i++ {
		pr := models.PullRequest{
			ID:          i,
			Number:      i,
			RepoName:    repoName,
			AuthorEmail: "dev@company.com",
			AuthorID:    "user-1",
			State:       models.PRStateOpen,
			Title:       "PR",
			BaseBranch:  "main",
			HeadBranch:  "feat",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := store.StorePR(pr)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		state         string
		page          int
		pageSize      int
		expectedLen   int
		expectedTotal int
	}{
		{"first page", "", 1, 10, 10, 25},
		{"second page", "", 2, 10, 10, 25},
		{"third page", "", 3, 10, 5, 25},
		{"filter by open", "open", 1, 10, 10, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prs, total, err := store.GetPRsByRepoWithPagination(repoName, tt.state, tt.page, tt.pageSize)
			require.NoError(t, err)
			assert.Len(t, prs, tt.expectedLen)
			assert.Equal(t, tt.expectedTotal, total)
		})
	}
}

// TestMemoryStore_StoreReview tests storing and retrieving reviews.
func TestMemoryStore_StoreReview(t *testing.T) {
	store := NewMemoryStore()

	review := models.Review{
		ID:          1,
		PRID:        42,
		Reviewer:    "reviewer@company.com",
		State:       models.ReviewStateApproved,
		SubmittedAt: time.Now(),
		Body:        "LGTM",
	}

	err := store.StoreReview(review)
	require.NoError(t, err)

	// Verify storage by retrieving via PR ID
	reviews, err := store.GetReviewsByPRID(42)
	require.NoError(t, err)
	assert.Len(t, reviews, 1)
	assert.Equal(t, review.ID, reviews[0].ID)
	assert.Equal(t, review.Reviewer, reviews[0].Reviewer)
}

// TestMemoryStore_GetReviewsByPRID tests retrieving reviews by PR ID.
func TestMemoryStore_GetReviewsByPRID(t *testing.T) {
	store := NewMemoryStore()

	prID := int64(42)

	reviews := []models.Review{
		{
			ID:          1,
			PRID:        int(prID),
			Reviewer:    "reviewer1@company.com",
			State:       models.ReviewStateApproved,
			SubmittedAt: time.Now(),
		},
		{
			ID:          2,
			PRID:        int(prID),
			Reviewer:    "reviewer2@company.com",
			State:       models.ReviewStateChangesRequested,
			SubmittedAt: time.Now(),
		},
		{
			ID:          3,
			PRID:        99, // Different PR
			Reviewer:    "reviewer3@company.com",
			State:       models.ReviewStateApproved,
			SubmittedAt: time.Now(),
		},
	}

	for _, review := range reviews {
		err := store.StoreReview(review)
		require.NoError(t, err)
	}

	prReviews, err := store.GetReviewsByPRID(prID)
	require.NoError(t, err)
	assert.Len(t, prReviews, 2)
}

// TestMemoryStore_GetReviewsByReviewer tests retrieving reviews by reviewer email.
func TestMemoryStore_GetReviewsByReviewer(t *testing.T) {
	store := NewMemoryStore()

	reviewerEmail := "reviewer@company.com"

	reviews := []models.Review{
		{
			ID:          1,
			PRID:        42,
			Reviewer:    reviewerEmail,
			State:       models.ReviewStateApproved,
			SubmittedAt: time.Now(),
		},
		{
			ID:          2,
			PRID:        43,
			Reviewer:    "other@company.com",
			State:       models.ReviewStateChangesRequested,
			SubmittedAt: time.Now(),
		},
		{
			ID:          3,
			PRID:        44,
			Reviewer:    reviewerEmail,
			State:       models.ReviewStateApproved,
			SubmittedAt: time.Now(),
		},
	}

	for _, review := range reviews {
		err := store.StoreReview(review)
		require.NoError(t, err)
	}

	reviewerReviews, err := store.GetReviewsByReviewer(reviewerEmail)
	require.NoError(t, err)
	assert.Len(t, reviewerReviews, 2)

	// Verify all reviews belong to the reviewer
	for _, review := range reviewerReviews {
		assert.Equal(t, reviewerEmail, review.Reviewer)
	}
}

// TestMemoryStore_GetReviewsByRepoPR tests retrieving reviews by repo and PR number.
func TestMemoryStore_GetReviewsByRepoPR(t *testing.T) {
	store := NewMemoryStore()

	// First, store PRs
	pr1 := models.PullRequest{
		ID:          1,
		Number:      42,
		RepoName:    "acme/platform",
		AuthorEmail: "dev@company.com",
		AuthorID:    "user-1",
		State:       models.PRStateOpen,
		Title:       "PR 42",
		BaseBranch:  "main",
		HeadBranch:  "feat-1",
		CreatedAt:   time.Now(),
	}
	pr2 := models.PullRequest{
		ID:          2,
		Number:      43,
		RepoName:    "acme/platform",
		AuthorEmail: "dev@company.com",
		AuthorID:    "user-1",
		State:       models.PRStateOpen,
		Title:       "PR 43",
		BaseBranch:  "main",
		HeadBranch:  "feat-2",
		CreatedAt:   time.Now(),
	}

	require.NoError(t, store.StorePR(pr1))
	require.NoError(t, store.StorePR(pr2))

	// Store reviews
	reviews := []models.Review{
		{
			ID:          1,
			PRID:        1, // PR ID, not number
			Reviewer:    "reviewer1@company.com",
			State:       models.ReviewStateApproved,
			SubmittedAt: time.Now(),
		},
		{
			ID:          2,
			PRID:        1,
			Reviewer:    "reviewer2@company.com",
			State:       models.ReviewStateChangesRequested,
			SubmittedAt: time.Now(),
		},
		{
			ID:          3,
			PRID:        2, // Different PR
			Reviewer:    "reviewer3@company.com",
			State:       models.ReviewStateApproved,
			SubmittedAt: time.Now(),
		},
	}

	for _, review := range reviews {
		err := store.StoreReview(review)
		require.NoError(t, err)
	}

	// Get reviews for acme/platform#42
	prReviews, err := store.GetReviewsByRepoPR("acme/platform", 42)
	require.NoError(t, err)
	assert.Len(t, prReviews, 2)
}

// TestMemoryStore_StoreIssue tests storing and retrieving issues.
func TestMemoryStore_StoreIssue(t *testing.T) {
	store := NewMemoryStore()

	issue := models.Issue{
		Number:    1,
		Title:     "Bug in feature",
		Body:      "Description of bug",
		State:     models.IssueStateOpen,
		AuthorID:  "user-1",
		RepoName:  "acme/platform",
		CreatedAt: time.Now(),
	}

	err := store.StoreIssue(issue)
	require.NoError(t, err)

	retrieved, err := store.GetIssueByNumber("acme/platform", 1)
	require.NoError(t, err)
	assert.Equal(t, issue.Title, retrieved.Title)
	assert.Equal(t, issue.State, retrieved.State)
}

// TestMemoryStore_GetIssueByNumber_NotFound tests retrieving non-existent issue.
func TestMemoryStore_GetIssueByNumber_NotFound(t *testing.T) {
	store := NewMemoryStore()

	_, err := store.GetIssueByNumber("acme/platform", 999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "issue not found")
}

// TestMemoryStore_GetIssuesByState tests filtering issues by state.
func TestMemoryStore_GetIssuesByState(t *testing.T) {
	store := NewMemoryStore()

	repoName := "acme/platform"

	issues := []models.Issue{
		{
			Number:    1,
			Title:     "Issue 1",
			State:     models.IssueStateOpen,
			AuthorID:  "user-1",
			RepoName:  repoName,
			CreatedAt: time.Now(),
		},
		{
			Number:    2,
			Title:     "Issue 2",
			State:     models.IssueStateClosed,
			AuthorID:  "user-2",
			RepoName:  repoName,
			CreatedAt: time.Now(),
		},
		{
			Number:    3,
			Title:     "Issue 3",
			State:     models.IssueStateOpen,
			AuthorID:  "user-3",
			RepoName:  repoName,
			CreatedAt: time.Now(),
		},
	}

	for _, issue := range issues {
		err := store.StoreIssue(issue)
		require.NoError(t, err)
	}

	// Get open issues
	openIssues, err := store.GetIssuesByState(repoName, models.IssueStateOpen)
	require.NoError(t, err)
	assert.Len(t, openIssues, 2)

	// Get closed issues
	closedIssues, err := store.GetIssuesByState(repoName, models.IssueStateClosed)
	require.NoError(t, err)
	assert.Len(t, closedIssues, 1)
	assert.Equal(t, 2, closedIssues[0].Number)
}

// TestMemoryStore_GetIssuesByRepo tests retrieving all issues for a repo.
func TestMemoryStore_GetIssuesByRepo(t *testing.T) {
	store := NewMemoryStore()

	repo1 := "acme/platform"
	repo2 := "acme/api"

	issues := []models.Issue{
		{
			Number:    1,
			Title:     "Issue 1",
			State:     models.IssueStateOpen,
			AuthorID:  "user-1",
			RepoName:  repo1,
			CreatedAt: time.Now(),
		},
		{
			Number:    2,
			Title:     "Issue 2",
			State:     models.IssueStateClosed,
			AuthorID:  "user-2",
			RepoName:  repo1,
			CreatedAt: time.Now(),
		},
		{
			Number:    1,
			Title:     "Issue 1",
			State:     models.IssueStateOpen,
			AuthorID:  "user-3",
			RepoName:  repo2,
			CreatedAt: time.Now(),
		},
	}

	for _, issue := range issues {
		err := store.StoreIssue(issue)
		require.NoError(t, err)
	}

	repo1Issues, err := store.GetIssuesByRepo(repo1)
	require.NoError(t, err)
	assert.Len(t, repo1Issues, 2)

	repo2Issues, err := store.GetIssuesByRepo(repo2)
	require.NoError(t, err)
	assert.Len(t, repo2Issues, 1)
}

// TestMemoryStore_ThreadSafety tests concurrent access to storage methods.
func TestMemoryStore_ThreadSafety(t *testing.T) {
	store := NewMemoryStore()

	var wg sync.WaitGroup
	concurrency := 10

	// Concurrent PR writes
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			pr := models.PullRequest{
				ID:          id,
				Number:      id,
				RepoName:    "acme/platform",
				AuthorEmail: "dev@company.com",
				AuthorID:    "user-1",
				State:       models.PRStateOpen,
				Title:       "PR",
				BaseBranch:  "main",
				HeadBranch:  "feat",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			_ = store.StorePR(pr)
		}(i)
	}
	wg.Wait()

	// Concurrent reads
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			_, _ = store.GetPRByID(id)
		}(i)
	}
	wg.Wait()

	// Verify all PRs were stored
	prs, err := store.GetPRsByStatus(models.PRStateOpen)
	require.NoError(t, err)
	assert.Equal(t, concurrency, len(prs))
}
