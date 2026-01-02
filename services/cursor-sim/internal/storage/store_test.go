package storage

import (
	"sync"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore_LoadDevelopers(t *testing.T) {
	store := NewMemoryStore()

	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob"},
	}

	err := store.LoadDevelopers(developers)
	require.NoError(t, err)

	// Verify retrieval by user ID
	dev, err := store.GetDeveloper("user_001")
	require.NoError(t, err)
	assert.Equal(t, "alice@example.com", dev.Email)

	// Verify retrieval by email
	dev, err = store.GetDeveloperByEmail("bob@example.com")
	require.NoError(t, err)
	assert.Equal(t, "user_002", dev.UserID)

	// Verify list all
	all := store.ListDevelopers()
	assert.Len(t, all, 2)
}

func TestMemoryStore_GetDeveloper_NotFound(t *testing.T) {
	store := NewMemoryStore()

	dev, err := store.GetDeveloper("nonexistent")
	require.Error(t, err)
	assert.Nil(t, dev)
	assert.Contains(t, err.Error(), "developer not found")
}

func TestMemoryStore_AddCommit(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	commit := models.Commit{
		CommitHash:      "abc123",
		UserID:          "user_001",
		UserEmail:       "test@example.com",
		RepoName:        "test/repo",
		TotalLinesAdded: 100,
		CommitTs:        now,
		CreatedAt:       now,
	}

	err := store.AddCommit(commit)
	require.NoError(t, err)

	// Verify retrieval by hash
	retrieved, err := store.GetCommitByHash("abc123")
	require.NoError(t, err)
	assert.Equal(t, "user_001", retrieved.UserID)
	assert.Equal(t, 100, retrieved.TotalLinesAdded)
}

func TestMemoryStore_GetCommitByHash_NotFound(t *testing.T) {
	store := NewMemoryStore()

	commit, err := store.GetCommitByHash("nonexistent")
	require.Error(t, err)
	assert.Nil(t, commit)
}

func TestMemoryStore_GetCommitsByTimeRange(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	// Add commits at different times
	commits := []models.Commit{
		{CommitHash: "c1", UserID: "user_001", CommitTs: now.Add(-2 * time.Hour), CreatedAt: now},
		{CommitHash: "c2", UserID: "user_001", CommitTs: now.Add(-1 * time.Hour), CreatedAt: now},
		{CommitHash: "c3", UserID: "user_001", CommitTs: now, CreatedAt: now},
		{CommitHash: "c4", UserID: "user_001", CommitTs: now.Add(1 * time.Hour), CreatedAt: now},
	}

	for _, c := range commits {
		err := store.AddCommit(c)
		require.NoError(t, err)
	}

	// Query range: -90 minutes to +30 minutes
	from := now.Add(-90 * time.Minute)
	to := now.Add(30 * time.Minute)

	results := store.GetCommitsByTimeRange(from, to)

	// Should return c2 and c3
	assert.Len(t, results, 2)
	hashes := []string{results[0].CommitHash, results[1].CommitHash}
	assert.Contains(t, hashes, "c2")
	assert.Contains(t, hashes, "c3")
}

func TestMemoryStore_GetCommitsByUser(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	commits := []models.Commit{
		{CommitHash: "c1", UserID: "user_001", CommitTs: now, CreatedAt: now},
		{CommitHash: "c2", UserID: "user_002", CommitTs: now, CreatedAt: now},
		{CommitHash: "c3", UserID: "user_001", CommitTs: now, CreatedAt: now},
	}

	for _, c := range commits {
		err := store.AddCommit(c)
		require.NoError(t, err)
	}

	// Get commits for user_001
	from := now.Add(-1 * time.Hour)
	to := now.Add(1 * time.Hour)
	results := store.GetCommitsByUser("user_001", from, to)

	assert.Len(t, results, 2)
	for _, c := range results {
		assert.Equal(t, "user_001", c.UserID)
	}
}

func TestMemoryStore_GetCommitsByRepo(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	commits := []models.Commit{
		{CommitHash: "c1", RepoName: "acme/api", UserID: "user_001", CommitTs: now, CreatedAt: now},
		{CommitHash: "c2", RepoName: "acme/web", UserID: "user_001", CommitTs: now, CreatedAt: now},
		{CommitHash: "c3", RepoName: "acme/api", UserID: "user_002", CommitTs: now, CreatedAt: now},
	}

	for _, c := range commits {
		err := store.AddCommit(c)
		require.NoError(t, err)
	}

	from := now.Add(-1 * time.Hour)
	to := now.Add(1 * time.Hour)
	results := store.GetCommitsByRepo("acme/api", from, to)

	assert.Len(t, results, 2)
	for _, c := range results {
		assert.Equal(t, "acme/api", c.RepoName)
	}
}

func TestMemoryStore_ConcurrentAccess(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	// Concurrent writes
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			commit := models.Commit{
				CommitHash: string(rune('a' + idx)),
				UserID:     "user_001",
				CommitTs:   now,
				CreatedAt:  now,
			}
			if err := store.AddCommit(commit); err != nil {
				errChan <- err
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			from := now.Add(-1 * time.Hour)
			to := now.Add(1 * time.Hour)
			_ = store.GetCommitsByTimeRange(from, to)
		}()
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		t.Errorf("concurrent access error: %v", err)
	}
}

func TestMemoryStore_TimeRangePerformance(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	// Add 1000 commits over 30 days
	for i := 0; i < 1000; i++ {
		commit := models.Commit{
			CommitHash: string(rune(i)),
			UserID:     "user_001",
			CommitTs:   now.Add(-time.Duration(i) * time.Hour),
			CreatedAt:  now,
		}
		err := store.AddCommit(commit)
		require.NoError(t, err)
	}

	// Query a 7-day range (should be fast)
	start := time.Now()
	from := now.Add(-7 * 24 * time.Hour)
	to := now
	results := store.GetCommitsByTimeRange(from, to)
	elapsed := time.Since(start)

	// Should complete in < 10ms
	assert.Less(t, elapsed, 10*time.Millisecond, "query should be fast")
	assert.Greater(t, len(results), 0, "should find commits in range")
}

func TestMemoryStore_EmptyStore(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	// All queries should return empty results
	assert.Len(t, store.ListDevelopers(), 0)
	assert.Len(t, store.GetCommitsByTimeRange(now, now.Add(1*time.Hour)), 0)
	assert.Len(t, store.GetCommitsByUser("user_001", now, now.Add(1*time.Hour)), 0)
	assert.Len(t, store.GetCommitsByRepo("test/repo", now, now.Add(1*time.Hour)), 0)
}

func TestMemoryStore_DuplicateCommitHash(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	commit1 := models.Commit{
		CommitHash: "abc123",
		UserID:     "user_001",
		CommitTs:   now,
		CreatedAt:  now,
	}

	commit2 := models.Commit{
		CommitHash: "abc123", // Same hash
		UserID:     "user_002",
		CommitTs:   now,
		CreatedAt:  now,
	}

	err := store.AddCommit(commit1)
	require.NoError(t, err)

	// Second commit with same hash should be rejected or overwrite
	err = store.AddCommit(commit2)
	require.NoError(t, err) // We allow overwrite

	// Should get the latest commit
	retrieved, err := store.GetCommitByHash("abc123")
	require.NoError(t, err)
	assert.Equal(t, "user_002", retrieved.UserID)
}

// PR Storage Tests

func TestMemoryStore_AddPR(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	pr := models.PullRequest{
		Number:      1,
		Title:       "feat: add authentication",
		State:       models.PRStateOpen,
		AuthorID:    "user_001",
		RepoName:    "acme/platform",
		BaseBranch:  "main",
		HeadBranch:  "feature/auth",
		Reviewers:   []string{"user_002"},
		Labels:      []string{"enhancement"},
		Additions:   150,
		Deletions:   20,
		CommitCount: 3,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := store.AddPR(pr)
	require.NoError(t, err)

	// Verify retrieval
	retrieved, err := store.GetPR("acme/platform", 1)
	require.NoError(t, err)
	assert.Equal(t, "feat: add authentication", retrieved.Title)
	assert.Equal(t, "user_001", retrieved.AuthorID)
}

func TestMemoryStore_GetPR_NotFound(t *testing.T) {
	store := NewMemoryStore()

	pr, err := store.GetPR("nonexistent/repo", 999)
	require.Error(t, err)
	assert.Nil(t, pr)
	assert.Contains(t, err.Error(), "not found")
}

func TestMemoryStore_GetPRsByRepo(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	// Add PRs to different repos
	prs := []models.PullRequest{
		{Number: 1, RepoName: "acme/api", State: models.PRStateOpen, AuthorID: "user_001", CreatedAt: now, UpdatedAt: now},
		{Number: 2, RepoName: "acme/api", State: models.PRStateMerged, AuthorID: "user_002", CreatedAt: now, UpdatedAt: now},
		{Number: 1, RepoName: "acme/web", State: models.PRStateOpen, AuthorID: "user_001", CreatedAt: now, UpdatedAt: now},
	}

	for _, pr := range prs {
		err := store.AddPR(pr)
		require.NoError(t, err)
	}

	// Get PRs for acme/api
	results := store.GetPRsByRepo("acme/api")
	assert.Len(t, results, 2)

	// Get PRs for acme/web
	results = store.GetPRsByRepo("acme/web")
	assert.Len(t, results, 1)

	// Get PRs for nonexistent repo
	results = store.GetPRsByRepo("nonexistent")
	assert.Len(t, results, 0)
}

func TestMemoryStore_GetPRsByState(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	prs := []models.PullRequest{
		{Number: 1, RepoName: "acme/api", State: models.PRStateOpen, AuthorID: "user_001", CreatedAt: now, UpdatedAt: now},
		{Number: 2, RepoName: "acme/api", State: models.PRStateMerged, AuthorID: "user_001", CreatedAt: now, UpdatedAt: now},
		{Number: 3, RepoName: "acme/api", State: models.PRStateOpen, AuthorID: "user_002", CreatedAt: now, UpdatedAt: now},
		{Number: 4, RepoName: "acme/api", State: models.PRStateClosed, AuthorID: "user_001", CreatedAt: now, UpdatedAt: now},
	}

	for _, pr := range prs {
		err := store.AddPR(pr)
		require.NoError(t, err)
	}

	// Get open PRs
	results := store.GetPRsByRepoAndState("acme/api", models.PRStateOpen)
	assert.Len(t, results, 2)
	for _, pr := range results {
		assert.Equal(t, models.PRStateOpen, pr.State)
	}

	// Get merged PRs
	results = store.GetPRsByRepoAndState("acme/api", models.PRStateMerged)
	assert.Len(t, results, 1)
}

func TestMemoryStore_GetPRsByAuthor(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	prs := []models.PullRequest{
		{Number: 1, RepoName: "acme/api", State: models.PRStateOpen, AuthorID: "user_001", CreatedAt: now, UpdatedAt: now},
		{Number: 2, RepoName: "acme/api", State: models.PRStateMerged, AuthorID: "user_002", CreatedAt: now, UpdatedAt: now},
		{Number: 3, RepoName: "acme/web", State: models.PRStateOpen, AuthorID: "user_001", CreatedAt: now, UpdatedAt: now},
	}

	for _, pr := range prs {
		err := store.AddPR(pr)
		require.NoError(t, err)
	}

	// Get PRs by user_001
	results := store.GetPRsByAuthor("user_001")
	assert.Len(t, results, 2)
	for _, pr := range results {
		assert.Equal(t, "user_001", pr.AuthorID)
	}
}

func TestMemoryStore_UpdatePR(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	pr := models.PullRequest{
		Number:    1,
		RepoName:  "acme/api",
		State:     models.PRStateOpen,
		AuthorID:  "user_001",
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := store.AddPR(pr)
	require.NoError(t, err)

	// Update the PR
	mergedAt := now.Add(24 * time.Hour)
	pr.State = models.PRStateMerged
	pr.MergedAt = &mergedAt
	pr.UpdatedAt = mergedAt

	err = store.UpdatePR(pr)
	require.NoError(t, err)

	// Verify update
	retrieved, err := store.GetPR("acme/api", 1)
	require.NoError(t, err)
	assert.Equal(t, models.PRStateMerged, retrieved.State)
	assert.NotNil(t, retrieved.MergedAt)
}

func TestMemoryStore_ListRepositories(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	// Add PRs to create repos implicitly
	prs := []models.PullRequest{
		{Number: 1, RepoName: "acme/api", State: models.PRStateOpen, AuthorID: "user_001", CreatedAt: now, UpdatedAt: now},
		{Number: 1, RepoName: "acme/web", State: models.PRStateOpen, AuthorID: "user_001", CreatedAt: now, UpdatedAt: now},
		{Number: 2, RepoName: "acme/api", State: models.PRStateOpen, AuthorID: "user_001", CreatedAt: now, UpdatedAt: now},
	}

	for _, pr := range prs {
		_ = store.AddPR(pr)
	}

	repos := store.ListRepositories()
	assert.Len(t, repos, 2)
	assert.Contains(t, repos, "acme/api")
	assert.Contains(t, repos, "acme/web")
}

// ReviewComment Storage Tests

func TestMemoryStore_AddReviewComment(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	comment := models.ReviewComment{
		ID:        1,
		PRNumber:  42,
		RepoName:  "acme/platform",
		AuthorID:  "user_002",
		Body:      "LGTM!",
		State:     models.ReviewStateApproved,
		CreatedAt: now,
	}

	err := store.AddReviewComment(comment)
	require.NoError(t, err)

	// Verify retrieval
	comments := store.GetReviewComments("acme/platform", 42)
	assert.Len(t, comments, 1)
	assert.Equal(t, "LGTM!", comments[0].Body)
}

func TestMemoryStore_GetReviewComments_MultipleComments(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	comments := []models.ReviewComment{
		{ID: 1, PRNumber: 42, RepoName: "acme/api", AuthorID: "user_001", Body: "Needs changes", State: models.ReviewStateChangesRequested, CreatedAt: now},
		{ID: 2, PRNumber: 42, RepoName: "acme/api", AuthorID: "user_002", Body: "Good idea", State: models.ReviewStatePending, CreatedAt: now.Add(time.Hour)},
		{ID: 3, PRNumber: 42, RepoName: "acme/api", AuthorID: "user_001", Body: "LGTM now", State: models.ReviewStateApproved, CreatedAt: now.Add(2 * time.Hour)},
		{ID: 4, PRNumber: 99, RepoName: "acme/api", AuthorID: "user_001", Body: "Different PR", State: models.ReviewStateApproved, CreatedAt: now},
	}

	for _, c := range comments {
		err := store.AddReviewComment(c)
		require.NoError(t, err)
	}

	// Get comments for PR 42
	results := store.GetReviewComments("acme/api", 42)
	assert.Len(t, results, 3)

	// Get comments for PR 99
	results = store.GetReviewComments("acme/api", 99)
	assert.Len(t, results, 1)

	// Get comments for nonexistent PR
	results = store.GetReviewComments("acme/api", 0)
	assert.Len(t, results, 0)
}

func TestMemoryStore_GetNextPRNumber(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	// First PR should get number 1
	num := store.GetNextPRNumber("acme/api")
	assert.Equal(t, 1, num)

	// Add a PR with that number
	pr := models.PullRequest{Number: 1, RepoName: "acme/api", State: models.PRStateOpen, AuthorID: "user_001", CreatedAt: now, UpdatedAt: now}
	_ = store.AddPR(pr)

	// Next PR should get number 2
	num = store.GetNextPRNumber("acme/api")
	assert.Equal(t, 2, num)

	// Add PRs 2 and 5
	pr.Number = 2
	_ = store.AddPR(pr)
	pr.Number = 5
	_ = store.AddPR(pr)

	// Next should be 6
	num = store.GetNextPRNumber("acme/api")
	assert.Equal(t, 6, num)

	// Different repo should start at 1
	num = store.GetNextPRNumber("acme/web")
	assert.Equal(t, 1, num)
}

func TestMemoryStore_PRConcurrentAccess(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	// Concurrent PR writes
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			pr := models.PullRequest{
				Number:    idx + 1,
				RepoName:  "acme/api",
				State:     models.PRStateOpen,
				AuthorID:  "user_001",
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := store.AddPR(pr); err != nil {
				errChan <- err
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = store.GetPRsByRepo("acme/api")
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Errorf("concurrent PR access error: %v", err)
	}
}
