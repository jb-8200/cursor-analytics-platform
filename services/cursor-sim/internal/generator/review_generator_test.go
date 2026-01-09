package generator

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockReviewStore implements ReviewStore for testing
type mockReviewStore struct {
	prs            map[string]map[int]*models.PullRequest
	reviewComments map[string]map[int][]models.ReviewComment
	developers     map[string]*seed.Developer
}

func newMockReviewStore() *mockReviewStore {
	return &mockReviewStore{
		prs:            make(map[string]map[int]*models.PullRequest),
		reviewComments: make(map[string]map[int][]models.ReviewComment),
		developers:     make(map[string]*seed.Developer),
	}
}

func (m *mockReviewStore) GetPR(repoName string, number int) (*models.PullRequest, error) {
	if repoPRs, ok := m.prs[repoName]; ok {
		if pr, ok := repoPRs[number]; ok {
			return pr, nil
		}
	}
	return nil, nil
}

func (m *mockReviewStore) GetPRsByRepo(repoName string) []models.PullRequest {
	result := []models.PullRequest{}
	if repoPRs, ok := m.prs[repoName]; ok {
		for _, pr := range repoPRs {
			result = append(result, *pr)
		}
	}
	return result
}

func (m *mockReviewStore) UpdatePR(pr models.PullRequest) error {
	if m.prs[pr.RepoName] == nil {
		m.prs[pr.RepoName] = make(map[int]*models.PullRequest)
	}
	m.prs[pr.RepoName][pr.Number] = &pr
	return nil
}

func (m *mockReviewStore) AddReviewComment(comment models.ReviewComment) error {
	if m.reviewComments[comment.RepoName] == nil {
		m.reviewComments[comment.RepoName] = make(map[int][]models.ReviewComment)
	}
	m.reviewComments[comment.RepoName][comment.PRNumber] = append(
		m.reviewComments[comment.RepoName][comment.PRNumber],
		comment,
	)
	return nil
}

func (m *mockReviewStore) GetReviewComments(repoName string, prNumber int) []models.ReviewComment {
	if repoComments, ok := m.reviewComments[repoName]; ok {
		return repoComments[prNumber]
	}
	return []models.ReviewComment{}
}

func (m *mockReviewStore) ListDevelopers() []seed.Developer {
	result := make([]seed.Developer, 0, len(m.developers))
	for _, dev := range m.developers {
		result = append(result, *dev)
	}
	return result
}

func (m *mockReviewStore) GetDeveloper(userID string) (*seed.Developer, error) {
	if dev, ok := m.developers[userID]; ok {
		return dev, nil
	}
	return nil, nil
}

func (m *mockReviewStore) addPR(pr models.PullRequest) {
	if m.prs[pr.RepoName] == nil {
		m.prs[pr.RepoName] = make(map[int]*models.PullRequest)
	}
	m.prs[pr.RepoName][pr.Number] = &pr
}

func (m *mockReviewStore) addDeveloper(dev seed.Developer) {
	m.developers[dev.UserID] = &dev
}

func TestReviewGenerator_SelectReviewers(t *testing.T) {
	store := newMockReviewStore()

	// Add developers on the same team
	store.addDeveloper(seed.Developer{UserID: "user_001", Email: "alice@example.com", Name: "Alice", Team: "Platform"})
	store.addDeveloper(seed.Developer{UserID: "user_002", Email: "bob@example.com", Name: "Bob", Team: "Platform"})
	store.addDeveloper(seed.Developer{UserID: "user_003", Email: "charlie@example.com", Name: "Charlie", Team: "Platform"})
	store.addDeveloper(seed.Developer{UserID: "user_004", Email: "diana@example.com", Name: "Diana", Team: "Backend"})

	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			ReviewPatterns: seed.ReviewPatterns{
				ReviewerCount: seed.ReviewerCountParams{
					Base: 2,
				},
			},
		},
	}

	gen := NewReviewGeneratorWithSeed(seedData, store, 42)

	// Author is user_001 on Platform team
	reviewers := gen.SelectReviewers("user_001", "Platform")

	// Should select reviewers from same team but not the author
	assert.NotEmpty(t, reviewers, "should select at least one reviewer")
	assert.LessOrEqual(t, len(reviewers), 2, "should not exceed max reviewers")

	for _, r := range reviewers {
		assert.NotEqual(t, "user_001", r, "author should not be a reviewer")
	}
}

func TestReviewGenerator_SelectReviewers_ExcludesAuthor(t *testing.T) {
	store := newMockReviewStore()

	// Only 2 developers on the team
	store.addDeveloper(seed.Developer{UserID: "author", Email: "author@example.com", Name: "Author", Team: "Platform"})
	store.addDeveloper(seed.Developer{UserID: "reviewer", Email: "reviewer@example.com", Name: "Reviewer", Team: "Platform"})

	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			ReviewPatterns: seed.ReviewPatterns{
				ReviewerCount: seed.ReviewerCountParams{
					Base: 2,
				},
			},
		},
	}

	gen := NewReviewGeneratorWithSeed(seedData, store, 42)
	reviewers := gen.SelectReviewers("author", "Platform")

	// Should only select "reviewer", not "author"
	assert.Len(t, reviewers, 1)
	assert.Equal(t, "reviewer", reviewers[0])
}

func TestReviewGenerator_GenerateReviewComments(t *testing.T) {
	store := newMockReviewStore()
	now := time.Now()

	// Add PR
	pr := models.PullRequest{
		Number:    1,
		RepoName:  "acme/api",
		AuthorID:  "user_001",
		State:     models.PRStateOpen,
		Additions: 100,
		Deletions: 20,
		CreatedAt: now.Add(-2 * time.Hour),
	}
	store.addPR(pr)

	store.addDeveloper(seed.Developer{UserID: "user_001", Team: "Platform"})
	store.addDeveloper(seed.Developer{UserID: "user_002", Team: "Platform"})

	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			ReviewPatterns: seed.ReviewPatterns{
				CommentsPer100LOC: seed.CommentDensity{
					Base: 3.0,
				},
			},
		},
		TextTemplates: seed.TextTemplates{
			ReviewComments: seed.ReviewCommentTemplates{
				Approval:   []string{"LGTM"},
				Suggestion: []string{"Consider refactoring this"},
				Style:      []string{"Nice improvement!"},
				Logic:      []string{"This might cause issues"},
			},
		},
	}

	gen := NewReviewGeneratorWithSeed(seedData, store, 42)

	// Generate comments for the PR
	err := gen.GenerateReviewForPR("acme/api", 1, "user_002")
	require.NoError(t, err)

	comments := store.GetReviewComments("acme/api", 1)
	assert.NotEmpty(t, comments, "should generate at least one comment")

	// All comments should be from the reviewer
	for _, c := range comments {
		assert.Equal(t, "user_002", c.AuthorID)
		assert.Equal(t, 1, c.PRNumber)
		assert.Equal(t, "acme/api", c.RepoName)
		assert.NotEmpty(t, c.Body)
	}
}

func TestReviewGenerator_GenerateApprovalDecision(t *testing.T) {
	tests := []struct {
		name      string
		aiRatio   float64
		additions int
	}{
		{"low AI ratio", 0.2, 50},
		{"high AI ratio", 0.9, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMockReviewStore()

			pr := models.PullRequest{
				Number:    1,
				RepoName:  "acme/api",
				AuthorID:  "user_001",
				State:     models.PRStateOpen,
				AIRatio:   tt.aiRatio,
				Additions: tt.additions,
			}
			store.addPR(pr)

			seedData := &seed.SeedData{
				Version: "1.0.0",
				PRLifecycle: seed.PRLifecycle{
					QualityOutcomes: seed.QualityOutcomes{
						RevertProbability: seed.OutcomeParams{
							Base: 0.1,
						},
					},
				},
			}

			gen := NewReviewGeneratorWithSeed(seedData, store, 42)
			decision := gen.GenerateApprovalDecision(&pr)

			// Decision should be valid
			assert.Contains(t, []models.ReviewState{
				models.ReviewStateApproved,
				models.ReviewStateChangesRequested,
				models.ReviewStatePending,
			}, decision)
		})
	}
}

func TestReviewGenerator_SimulateReviewIterations(t *testing.T) {
	store := newMockReviewStore()
	now := time.Now()

	// Create PR with high AI ratio (more likely to need iterations)
	pr := models.PullRequest{
		Number:    1,
		RepoName:  "acme/api",
		AuthorID:  "user_001",
		State:     models.PRStateOpen,
		AIRatio:   0.8,
		Additions: 200,
		CreatedAt: now.Add(-24 * time.Hour),
	}
	store.addPR(pr)

	store.addDeveloper(seed.Developer{UserID: "user_001", Team: "Platform"})
	store.addDeveloper(seed.Developer{UserID: "user_002", Team: "Platform"})

	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			ReviewPatterns: seed.ReviewPatterns{
				ReviewerCount: seed.ReviewerCountParams{
					Base: 1,
				},
				Iterations: seed.IterationParams{
					BaseDistribution: "poisson",
					Params:           map[string]float64{"lambda": 1.5},
				},
				CommentsPer100LOC: seed.CommentDensity{
					Base: 2.0,
				},
			},
		},
		TextTemplates: seed.TextTemplates{
			ReviewComments: seed.ReviewCommentTemplates{
				Approval:   []string{"Good"},
				Suggestion: []string{"Needs work"},
				Style:      []string{"LGTM"},
				Logic:      []string{"Check this"},
			},
		},
	}

	gen := NewReviewGeneratorWithSeed(seedData, store, 42)

	iterations, err := gen.SimulateReviewIterations("acme/api", 1)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, iterations, 1, "should have at least 1 iteration")
	assert.LessOrEqual(t, iterations, 5, "should not have excessive iterations")
}

func TestReviewGenerator_Reproducibility(t *testing.T) {
	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			ReviewPatterns: seed.ReviewPatterns{
				ReviewerCount: seed.ReviewerCountParams{
					Base: 2,
				},
			},
		},
	}

	// Create two identical stores
	store1 := newMockReviewStore()
	store2 := newMockReviewStore()

	for _, store := range []*mockReviewStore{store1, store2} {
		store.addDeveloper(seed.Developer{UserID: "user_001", Team: "Platform"})
		store.addDeveloper(seed.Developer{UserID: "user_002", Team: "Platform"})
		store.addDeveloper(seed.Developer{UserID: "user_003", Team: "Platform"})
	}

	gen1 := NewReviewGeneratorWithSeed(seedData, store1, 12345)
	gen2 := NewReviewGeneratorWithSeed(seedData, store2, 12345)

	// Same seed should produce same reviewers
	reviewers1 := gen1.SelectReviewers("user_001", "Platform")
	reviewers2 := gen2.SelectReviewers("user_001", "Platform")

	assert.Equal(t, reviewers1, reviewers2, "same seed should produce same reviewers")
}

func TestReviewGenerator_GenerateReviewsForRepo(t *testing.T) {
	store := newMockReviewStore()
	now := time.Now()

	// Add multiple PRs
	prs := []models.PullRequest{
		{Number: 1, RepoName: "acme/api", AuthorID: "user_001", State: models.PRStateOpen, CreatedAt: now, Additions: 50},
		{Number: 2, RepoName: "acme/api", AuthorID: "user_002", State: models.PRStateOpen, CreatedAt: now, Additions: 50},
		{Number: 3, RepoName: "acme/api", AuthorID: "user_001", State: models.PRStateMerged, CreatedAt: now, Additions: 50}, // Already merged
	}
	for _, pr := range prs {
		store.addPR(pr)
	}

	store.addDeveloper(seed.Developer{UserID: "user_001", Team: "Platform"})
	store.addDeveloper(seed.Developer{UserID: "user_002", Team: "Platform"})
	store.addDeveloper(seed.Developer{UserID: "user_003", Team: "Platform"})

	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			ReviewPatterns: seed.ReviewPatterns{
				ReviewerCount: seed.ReviewerCountParams{
					Base: 1,
				},
				CommentsPer100LOC: seed.CommentDensity{
					Base: 1.0,
				},
			},
		},
		TextTemplates: seed.TextTemplates{
			ReviewComments: seed.ReviewCommentTemplates{
				Approval: []string{"LGTM"},
			},
		},
	}

	gen := NewReviewGeneratorWithSeed(seedData, store, 42)

	count, err := gen.GenerateReviewsForRepo("acme/api")
	require.NoError(t, err)

	// Should generate reviews for open PRs (1 and 2), not for merged PR (3)
	assert.GreaterOrEqual(t, count, 2, "should generate reviews for open PRs")
}

// =============================================================================
// Tests for the new GenerateReviewsForPR(pr PullRequest) []Review method
// =============================================================================

func TestGenerateReviewsForPR_Basic(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice", Team: "platform"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob", Team: "platform"},
		{UserID: "user_003", Email: "charlie@example.com", Name: "Charlie", Team: "platform"},
		{UserID: "user_004", Email: "diana@example.com", Name: "Diana", Team: "platform"},
	}

	seedData := &seed.SeedData{
		Developers: developers,
		TextTemplates: seed.TextTemplates{
			ReviewComments: seed.ReviewCommentTemplates{
				Style:      []string{"Consider adding a comment here"},
				Logic:      []string{"This logic could be simplified"},
				Suggestion: []string{"Maybe try a different approach"},
				Approval:   []string{"LGTM!", "Looks great!"},
			},
		},
	}

	gen := NewReviewGenerator(seedData, rng)

	// Create a sample merged PR
	prCreatedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	mergedAt := prCreatedAt.Add(48 * time.Hour)

	pr := models.PullRequest{
		Number:      1,
		Title:       "Add authentication flow",
		AuthorID:    "user_001",
		AuthorEmail: "alice@example.com",
		RepoName:    "acme/platform",
		State:       models.PRStateMerged,
		CreatedAt:   prCreatedAt,
		MergedAt:    &mergedAt,
	}

	reviews := gen.GenerateReviewsForPR(pr)

	// Should generate 1-3 reviews
	assert.GreaterOrEqual(t, len(reviews), 1, "should generate at least 1 review")
	assert.LessOrEqual(t, len(reviews), 3, "should generate at most 3 reviews")

	for _, review := range reviews {
		// Review should reference the PR
		assert.Equal(t, pr.Number, review.PRID, "review PRID should match PR number")

		// Reviewer should NOT be the PR author
		assert.NotEqual(t, pr.AuthorEmail, review.Reviewer, "reviewer should not be PR author")

		// Review timestamp should be between PR creation and merge
		assert.False(t, review.SubmittedAt.Before(pr.CreatedAt), "review should not be before PR creation")
		assert.False(t, review.SubmittedAt.After(*pr.MergedAt), "review should not be after PR merge")

		// Review should be valid
		require.NoError(t, review.Validate(), "review should be valid")
	}
}

func TestGenerateReviewsForPR_ReviewerNotPRAuthor(t *testing.T) {
	rng := rand.New(rand.NewSource(99999))

	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice", Team: "platform"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob", Team: "platform"},
	}

	seedData := &seed.SeedData{Developers: developers}
	gen := NewReviewGenerator(seedData, rng)

	prCreatedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	mergedAt := prCreatedAt.Add(24 * time.Hour)

	pr := models.PullRequest{
		Number:      1,
		AuthorID:    "user_001",
		AuthorEmail: "alice@example.com",
		State:       models.PRStateMerged,
		CreatedAt:   prCreatedAt,
		MergedAt:    &mergedAt,
	}

	// Generate many reviews to ensure no reviewer is the author
	for i := 0; i < 100; i++ {
		reviews := gen.GenerateReviewsForPR(pr)
		for _, review := range reviews {
			assert.NotEqual(t, pr.AuthorEmail, review.Reviewer,
				"reviewer %s should not be PR author on iteration %d", review.Reviewer, i)
		}
	}
}

func TestGenerateReviewsForPR_StateDistribution(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob"},
		{UserID: "user_003", Email: "charlie@example.com", Name: "Charlie"},
	}

	seedData := &seed.SeedData{Developers: developers}
	gen := NewReviewGenerator(seedData, rng)

	prCreatedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	// Generate many reviews across multiple PRs
	stateCounts := make(map[models.ReviewState]int)
	totalReviews := 0

	for i := 0; i < 200; i++ {
		mergedAt := prCreatedAt.Add(time.Duration(i*24+48) * time.Hour)

		pr := models.PullRequest{
			Number:      i + 1,
			AuthorID:    "user_001",
			AuthorEmail: "alice@example.com",
			State:       models.PRStateMerged,
			CreatedAt:   prCreatedAt.Add(time.Duration(i*24) * time.Hour),
			MergedAt:    &mergedAt,
		}

		reviews := gen.GenerateReviewsForPR(pr)
		for _, review := range reviews {
			stateCounts[review.State]++
			totalReviews++
		}
	}

	// Calculate percentages
	approvedPct := float64(stateCounts[models.ReviewStateApproved]) / float64(totalReviews) * 100
	changesPct := float64(stateCounts[models.ReviewStateChangesRequested]) / float64(totalReviews) * 100
	pendingPct := float64(stateCounts[models.ReviewStatePending]) / float64(totalReviews) * 100

	t.Logf("State distribution over %d reviews: approved=%.1f%%, changes_requested=%.1f%%, pending=%.1f%%",
		totalReviews, approvedPct, changesPct, pendingPct)

	// Expected: 70% approved, 20% changes_requested, 10% pending
	// Allow variance due to randomness
	assert.Greater(t, approvedPct, 55.0, "approved should be > 55%%")
	assert.Less(t, approvedPct, 85.0, "approved should be < 85%%")
	assert.Greater(t, changesPct, 10.0, "changes_requested should be > 10%%")
	assert.Less(t, changesPct, 35.0, "changes_requested should be < 35%%")
	assert.Greater(t, pendingPct, 2.0, "pending should be > 2%%")
	assert.Less(t, pendingPct, 20.0, "pending should be < 20%%")
}

func TestGenerateReviewsForPR_TimingWithinBounds(t *testing.T) {
	rng := rand.New(rand.NewSource(54321))

	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob"},
	}

	seedData := &seed.SeedData{Developers: developers}
	gen := NewReviewGenerator(seedData, rng)

	prCreatedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	mergedAt := prCreatedAt.Add(72 * time.Hour) // 3 days later

	pr := models.PullRequest{
		Number:      1,
		AuthorID:    "user_001",
		AuthorEmail: "alice@example.com",
		State:       models.PRStateMerged,
		CreatedAt:   prCreatedAt,
		MergedAt:    &mergedAt,
	}

	// Generate many reviews to test timing consistency
	for i := 0; i < 100; i++ {
		reviews := gen.GenerateReviewsForPR(pr)
		for _, review := range reviews {
			assert.False(t, review.SubmittedAt.Before(prCreatedAt),
				"review should not be before PR creation")
			assert.False(t, review.SubmittedAt.After(mergedAt),
				"review should not be after PR merge")
		}
	}
}

func TestGenerateReviewsForPR_OpenPR(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob"},
	}

	seedData := &seed.SeedData{Developers: developers}
	gen := NewReviewGenerator(seedData, rng)

	prCreatedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	pr := models.PullRequest{
		Number:      1,
		AuthorID:    "user_001",
		AuthorEmail: "alice@example.com",
		State:       models.PRStateOpen,
		CreatedAt:   prCreatedAt,
		// No MergedAt - PR is open
	}

	reviews := gen.GenerateReviewsForPR(pr)

	// Open PRs can still have reviews
	maxReviewTime := prCreatedAt.Add(7 * 24 * time.Hour) // 7 days max
	for _, review := range reviews {
		assert.False(t, review.SubmittedAt.Before(prCreatedAt),
			"review should not be before PR creation")
		assert.False(t, review.SubmittedAt.After(maxReviewTime),
			"review should not be more than 7 days after PR creation for open PR")
	}
}

func TestGenerateReviewsForPR_CommentsForNonApproved(t *testing.T) {
	rng := rand.New(rand.NewSource(77777))

	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob"},
		{UserID: "user_003", Email: "charlie@example.com", Name: "Charlie"},
	}

	seedData := &seed.SeedData{
		Developers: developers,
		TextTemplates: seed.TextTemplates{
			ReviewComments: seed.ReviewCommentTemplates{
				Style:      []string{"Fix formatting", "Add whitespace"},
				Logic:      []string{"Simplify logic", "Consider edge cases"},
				Suggestion: []string{"Try this instead", "Consider refactoring"},
				Approval:   []string{"LGTM!", "Ship it!"},
			},
		},
	}

	gen := NewReviewGenerator(seedData, rng)

	prCreatedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	// Generate many reviews to find changes_requested ones
	changesRequestedCount := 0
	changesRequestedWithComments := 0

	for i := 0; i < 500; i++ {
		mergedAt := prCreatedAt.Add(time.Duration(48+i) * time.Hour)

		pr := models.PullRequest{
			Number:      i + 1,
			AuthorID:    "user_001",
			AuthorEmail: "alice@example.com",
			RepoName:    "acme/platform",
			State:       models.PRStateMerged,
			CreatedAt:   prCreatedAt.Add(time.Duration(i) * time.Hour),
			MergedAt:    &mergedAt,
		}

		reviews := gen.GenerateReviewsForPR(pr)
		for _, review := range reviews {
			if review.State == models.ReviewStateChangesRequested {
				changesRequestedCount++
				if len(review.Comments) > 0 || review.Body != "" {
					changesRequestedWithComments++
				}
			}
		}
	}

	t.Logf("Found %d changes_requested reviews, %d with comments/body",
		changesRequestedCount, changesRequestedWithComments)

	// Most changes_requested reviews should have comments or body
	if changesRequestedCount > 0 {
		commentRate := float64(changesRequestedWithComments) / float64(changesRequestedCount)
		assert.Greater(t, commentRate, 0.5,
			"at least 50%% of changes_requested reviews should have comments")
	}
}

func TestGenerateReviewsForPR_NoAvailableReviewers(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	// Only one developer (the PR author)
	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
	}

	seedData := &seed.SeedData{Developers: developers}
	gen := NewReviewGenerator(seedData, rng)

	prCreatedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	mergedAt := prCreatedAt.Add(24 * time.Hour)

	pr := models.PullRequest{
		Number:      1,
		AuthorID:    "user_001",
		AuthorEmail: "alice@example.com",
		State:       models.PRStateMerged,
		CreatedAt:   prCreatedAt,
		MergedAt:    &mergedAt,
	}

	// Should not panic when no reviewers available
	reviews := gen.GenerateReviewsForPR(pr)

	// Should return empty slice or handle gracefully
	assert.Empty(t, reviews, "should return empty reviews when no other developers available")
}

func TestGenerateReviewsForPR_UniqueReviewers(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob"},
		{UserID: "user_003", Email: "charlie@example.com", Name: "Charlie"},
		{UserID: "user_004", Email: "diana@example.com", Name: "Diana"},
	}

	seedData := &seed.SeedData{Developers: developers}
	gen := NewReviewGenerator(seedData, rng)

	prCreatedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	mergedAt := prCreatedAt.Add(24 * time.Hour)

	pr := models.PullRequest{
		Number:      1,
		AuthorID:    "user_001",
		AuthorEmail: "alice@example.com",
		State:       models.PRStateMerged,
		CreatedAt:   prCreatedAt,
		MergedAt:    &mergedAt,
	}

	// Check across multiple iterations
	for i := 0; i < 100; i++ {
		reviews := gen.GenerateReviewsForPR(pr)

		// Each review should have a unique reviewer
		reviewers := make(map[string]bool)
		for _, review := range reviews {
			assert.False(t, reviewers[review.Reviewer],
				"reviewer %s should not be duplicated on iteration %d", review.Reviewer, i)
			reviewers[review.Reviewer] = true
		}
	}
}

func TestGenerateReviewsForPR_Reproducibility(t *testing.T) {
	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob"},
	}

	seedData := &seed.SeedData{Developers: developers}

	prCreatedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	mergedAt := prCreatedAt.Add(24 * time.Hour)

	pr := models.PullRequest{
		Number:      1,
		AuthorID:    "user_001",
		AuthorEmail: "alice@example.com",
		State:       models.PRStateMerged,
		CreatedAt:   prCreatedAt,
		MergedAt:    &mergedAt,
	}

	// Generate with same seed twice
	gen1 := NewReviewGenerator(seedData, rand.New(rand.NewSource(999)))
	reviews1 := gen1.GenerateReviewsForPR(pr)

	gen2 := NewReviewGenerator(seedData, rand.New(rand.NewSource(999)))
	reviews2 := gen2.GenerateReviewsForPR(pr)

	// Results should be identical
	require.Equal(t, len(reviews1), len(reviews2), "review count should match")

	for i := range reviews1 {
		assert.Equal(t, reviews1[i].Reviewer, reviews2[i].Reviewer,
			"reviewer should match for review %d", i)
		assert.Equal(t, reviews1[i].State, reviews2[i].State,
			"state should match for review %d", i)
		assert.True(t, reviews1[i].SubmittedAt.Equal(reviews2[i].SubmittedAt),
			"submitted_at should match for review %d", i)
	}
}

func TestGenerateReviewsForPR_ReviewIDsIncrement(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob"},
		{UserID: "user_003", Email: "charlie@example.com", Name: "Charlie"},
	}

	seedData := &seed.SeedData{Developers: developers}
	gen := NewReviewGenerator(seedData, rng)

	prCreatedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	var allReviewIDs []int

	// Generate reviews for multiple PRs
	for prNum := 1; prNum <= 5; prNum++ {
		mergedAt := prCreatedAt.Add(time.Duration(prNum*24) * time.Hour)

		pr := models.PullRequest{
			Number:      prNum,
			AuthorID:    "user_001",
			AuthorEmail: "alice@example.com",
			State:       models.PRStateMerged,
			CreatedAt:   prCreatedAt.Add(time.Duration(prNum-1) * 24 * time.Hour),
			MergedAt:    &mergedAt,
		}

		reviews := gen.GenerateReviewsForPR(pr)
		for _, review := range reviews {
			allReviewIDs = append(allReviewIDs, review.ID)
		}
	}

	// All IDs should be unique and positive
	seen := make(map[int]bool)
	for i, id := range allReviewIDs {
		assert.Greater(t, id, 0, "review ID at index %d should be positive", i)
		assert.False(t, seen[id], "review ID %d at index %d should not be duplicated", id, i)
		seen[id] = true
	}
}

func TestGenerateReviewsForPR_ClosedPR(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob"},
	}

	seedData := &seed.SeedData{Developers: developers}
	gen := NewReviewGenerator(seedData, rng)

	prCreatedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	closedAt := prCreatedAt.Add(48 * time.Hour)

	pr := models.PullRequest{
		Number:      1,
		AuthorID:    "user_001",
		AuthorEmail: "alice@example.com",
		State:       models.PRStateClosed,
		CreatedAt:   prCreatedAt,
		ClosedAt:    &closedAt,
		// No MergedAt - PR was closed without merging
	}

	reviews := gen.GenerateReviewsForPR(pr)

	// Closed PRs can still have reviews, timing should be based on ClosedAt
	for _, review := range reviews {
		assert.False(t, review.SubmittedAt.Before(prCreatedAt),
			"review should not be before PR creation")
		assert.False(t, review.SubmittedAt.After(closedAt),
			"review should not be after PR closed")
	}
}
