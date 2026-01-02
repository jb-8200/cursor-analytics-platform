package generator

import (
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
