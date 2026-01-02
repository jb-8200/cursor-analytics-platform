package generator

import (
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockResearchStore implements the ResearchStore interface for testing.
type MockResearchStore struct {
	commits    []models.Commit
	prs        []models.PullRequest
	reviews    map[string][]models.ReviewComment // key: "repoName:prNumber"
	developers []seed.Developer
}

func NewMockResearchStore() *MockResearchStore {
	return &MockResearchStore{
		reviews: make(map[string][]models.ReviewComment),
	}
}

func (m *MockResearchStore) GetCommitsByTimeRange(from, to time.Time) []models.Commit {
	var result []models.Commit
	for _, c := range m.commits {
		if (c.CommitTs.Equal(from) || c.CommitTs.After(from)) &&
			(c.CommitTs.Equal(to) || c.CommitTs.Before(to)) {
			result = append(result, c)
		}
	}
	return result
}

func (m *MockResearchStore) GetPRsByRepo(repoName string) []models.PullRequest {
	var result []models.PullRequest
	for _, pr := range m.prs {
		if pr.RepoName == repoName {
			result = append(result, pr)
		}
	}
	return result
}

func (m *MockResearchStore) GetPRsByRepoAndState(repoName string, state models.PRState) []models.PullRequest {
	var result []models.PullRequest
	for _, pr := range m.prs {
		if pr.RepoName == repoName && pr.State == state {
			result = append(result, pr)
		}
	}
	return result
}

func (m *MockResearchStore) GetReviewComments(repoName string, prNumber int) []models.ReviewComment {
	key := repoName + ":" + string(rune(prNumber))
	return m.reviews[key]
}

func (m *MockResearchStore) GetDeveloper(userID string) (*seed.Developer, error) {
	for _, d := range m.developers {
		if d.UserID == userID {
			return &d, nil
		}
	}
	return nil, nil
}

func (m *MockResearchStore) ListRepositories() []string {
	repoSet := make(map[string]bool)
	for _, pr := range m.prs {
		repoSet[pr.RepoName] = true
	}
	var repos []string
	for repo := range repoSet {
		repos = append(repos, repo)
	}
	return repos
}

func (m *MockResearchStore) AddCommit(c models.Commit) {
	m.commits = append(m.commits, c)
}

func (m *MockResearchStore) AddPR(pr models.PullRequest) {
	m.prs = append(m.prs, pr)
}

func (m *MockResearchStore) AddReviewComment(comment models.ReviewComment) {
	key := comment.RepoName + ":" + string(rune(comment.PRNumber))
	m.reviews[key] = append(m.reviews[key], comment)
}

func (m *MockResearchStore) AddDeveloper(d seed.Developer) {
	m.developers = append(m.developers, d)
}

func TestResearchGenerator_New(t *testing.T) {
	store := NewMockResearchStore()
	seedData := &seed.SeedData{}

	gen := NewResearchGenerator(seedData, store)
	assert.NotNil(t, gen)
}

func TestResearchGenerator_JoinCommitPRData(t *testing.T) {
	store := NewMockResearchStore()

	// Add test developer
	store.AddDeveloper(seed.Developer{
		UserID:    "user_001",
		Email:     "alice@example.com",
		Name:      "Alice",
		Seniority: "senior",
	})

	// Add test commits
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	commit1 := models.Commit{
		CommitHash:         "abc123",
		UserID:             "user_001",
		UserEmail:          "alice@example.com",
		RepoName:           "test/repo",
		BranchName:         "feature-1",
		TotalLinesAdded:    100,
		TabLinesAdded:      30,
		ComposerLinesAdded: 20,
		NonAILinesAdded:    50,
		CommitTs:           baseTime,
	}
	store.AddCommit(commit1)

	// Add test PR
	mergeTime := baseTime.Add(24 * time.Hour)
	pr := models.PullRequest{
		Number:        1,
		RepoName:      "test/repo",
		AuthorID:      "user_001",
		State:         models.PRStateMerged,
		HeadBranch:    "feature-1",
		Additions:     100,
		Deletions:     20,
		ChangedFiles:  3,
		AIRatio:       0.5,
		TabLines:      30,
		ComposerLines: 20,
		CreatedAt:     baseTime.Add(-1 * time.Hour),
		MergedAt:      &mergeTime,
	}
	store.AddPR(pr)

	seedData := &seed.SeedData{
		Repositories: []seed.Repository{
			{RepoName: "test/repo", Maturity: seed.Maturity{AgeDays: 365}}, // mature
		},
	}

	gen := NewResearchGenerator(seedData, store)

	dataPoints, err := gen.JoinCommitPRData(baseTime.Add(-1*time.Hour), baseTime.Add(2*time.Hour))
	require.NoError(t, err)
	require.Len(t, dataPoints, 1)

	dp := dataPoints[0]
	assert.Equal(t, "abc123", dp.CommitHash)
	assert.Equal(t, 1, dp.PRNumber)
	assert.Equal(t, "user_001", dp.AuthorID)
	assert.Equal(t, 0.5, dp.AIRatio)
	assert.Equal(t, "senior", dp.AuthorSeniority)
}

func TestResearchGenerator_CalculateCycleTime(t *testing.T) {
	store := NewMockResearchStore()
	seedData := &seed.SeedData{}
	gen := NewResearchGenerator(seedData, store)

	// Test coding lead time (first commit to last commit in PR)
	commits := []models.Commit{
		{CommitTs: time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)},
		{CommitTs: time.Date(2026, 1, 15, 14, 0, 0, 0, time.UTC)},
		{CommitTs: time.Date(2026, 1, 15, 18, 0, 0, 0, time.UTC)},
	}

	codingTime := gen.CalculateCodingLeadTime(commits)
	assert.Equal(t, 8.0, codingTime) // 8 hours from first to last commit

	// Test review lead time (PR created to approved/merged)
	createdAt := time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC)
	mergedAt := time.Date(2026, 1, 16, 17, 0, 0, 0, time.UTC)
	pr := models.PullRequest{
		CreatedAt: createdAt,
		MergedAt:  &mergedAt,
	}

	reviewTime := gen.CalculateReviewLeadTime(pr)
	assert.Equal(t, 32.0, reviewTime) // 32 hours from created to merged

	// Test merge lead time (PR approved to merged)
	approvedAt := time.Date(2026, 1, 16, 14, 0, 0, 0, time.UTC)
	mergeTime := gen.CalculateMergeLeadTime(approvedAt, mergedAt)
	assert.Equal(t, 3.0, mergeTime) // 3 hours from approved to merged
}

func TestResearchGenerator_DetermineGreenfield(t *testing.T) {
	store := NewMockResearchStore()
	seedData := &seed.SeedData{}
	gen := NewResearchGenerator(seedData, store)

	tests := []struct {
		name     string
		pr       models.PullRequest
		expected bool
	}{
		{
			name: "High additions, no deletions - greenfield",
			pr: models.PullRequest{
				Additions: 500,
				Deletions: 0,
			},
			expected: true,
		},
		{
			name: "Additions equal deletions - not greenfield",
			pr: models.PullRequest{
				Additions: 100,
				Deletions: 100,
			},
			expected: false,
		},
		{
			name: "Low additions ratio - not greenfield",
			pr: models.PullRequest{
				Additions: 100,
				Deletions: 200,
			},
			expected: false,
		},
		{
			name: "High additions ratio (>80%) - greenfield",
			pr: models.PullRequest{
				Additions: 400,
				Deletions: 50,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.IsGreenfield(tt.pr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResearchGenerator_ApplyControlVariables(t *testing.T) {
	store := NewMockResearchStore()
	store.AddDeveloper(seed.Developer{
		UserID:    "user_001",
		Seniority: "senior",
	})
	store.AddDeveloper(seed.Developer{
		UserID:    "user_002",
		Seniority: "junior",
	})

	seedData := &seed.SeedData{
		Repositories: []seed.Repository{
			{RepoName: "test/repo", Maturity: seed.Maturity{AgeDays: 365}},    // mature
			{RepoName: "test/new-repo", Maturity: seed.Maturity{AgeDays: 30}}, // greenfield
		},
	}

	gen := NewResearchGenerator(seedData, store)

	// Test senior developer in mature repo
	dp := &models.ResearchDataPoint{
		AuthorID: "user_001",
		RepoName: "test/repo",
	}
	gen.ApplyControlVariables(dp)
	assert.Equal(t, "senior", dp.AuthorSeniority)
	assert.Equal(t, "mature", dp.RepoMaturity) // AgeDays >= 180 = mature

	// Test junior developer in greenfield repo
	dp2 := &models.ResearchDataPoint{
		AuthorID: "user_002",
		RepoName: "test/new-repo",
	}
	gen.ApplyControlVariables(dp2)
	assert.Equal(t, "junior", dp2.AuthorSeniority)
	assert.Equal(t, "greenfield", dp2.RepoMaturity) // AgeDays < 90 = greenfield
}

func TestResearchGenerator_CountReviewIterations(t *testing.T) {
	store := NewMockResearchStore()
	seedData := &seed.SeedData{}
	gen := NewResearchGenerator(seedData, store)

	reviews := []models.ReviewComment{
		{State: models.ReviewStateChangesRequested},
		{State: models.ReviewStatePending},
		{State: models.ReviewStateChangesRequested},
		{State: models.ReviewStateApproved},
	}

	iterations := gen.CountReviewIterations(reviews)
	assert.Equal(t, 2, iterations) // 2 changes requested = 2 iterations
}

func TestResearchGenerator_GenerateDataset(t *testing.T) {
	store := NewMockResearchStore()

	// Set up test data
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	store.AddDeveloper(seed.Developer{
		UserID:    "user_001",
		Seniority: "senior",
	})

	// Add commits
	for i := 0; i < 5; i++ {
		store.AddCommit(models.Commit{
			CommitHash:         "hash_" + string(rune('a'+i)),
			UserID:             "user_001",
			RepoName:           "test/repo",
			BranchName:         "feature-1",
			TotalLinesAdded:    100,
			TabLinesAdded:      30,
			ComposerLinesAdded: 20,
			CommitTs:           baseTime.Add(time.Duration(i) * time.Hour),
		})
	}

	// Add PR
	mergeTime := baseTime.Add(48 * time.Hour)
	store.AddPR(models.PullRequest{
		Number:     1,
		RepoName:   "test/repo",
		AuthorID:   "user_001",
		State:      models.PRStateMerged,
		HeadBranch: "feature-1",
		Additions:  500,
		Deletions:  50,
		AIRatio:    0.5,
		CreatedAt:  baseTime.Add(-1 * time.Hour),
		MergedAt:   &mergeTime,
	})

	seedData := &seed.SeedData{
		Repositories: []seed.Repository{
			{RepoName: "test/repo", Maturity: seed.Maturity{AgeDays: 365}}, // mature
		},
	}

	gen := NewResearchGenerator(seedData, store)

	dataPoints, err := gen.GenerateDataset(baseTime.Add(-2*time.Hour), baseTime.Add(10*time.Hour))
	require.NoError(t, err)
	assert.Greater(t, len(dataPoints), 0)

	// Verify data point structure
	for _, dp := range dataPoints {
		assert.NotEmpty(t, dp.CommitHash)
		assert.NotEmpty(t, dp.AuthorID)
		assert.GreaterOrEqual(t, dp.AIRatio, 0.0)
		assert.LessOrEqual(t, dp.AIRatio, 1.0)
	}
}

func TestResearchGenerator_Reproducibility(t *testing.T) {
	setupStore := func() *MockResearchStore {
		store := NewMockResearchStore()
		baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
		store.AddDeveloper(seed.Developer{UserID: "user_001", Seniority: "senior"})
		store.AddCommit(models.Commit{
			CommitHash:      "abc123",
			UserID:          "user_001",
			RepoName:        "test/repo",
			BranchName:      "feature-1",
			TotalLinesAdded: 100,
			CommitTs:        baseTime,
		})
		mergeTime := baseTime.Add(24 * time.Hour)
		store.AddPR(models.PullRequest{
			Number:     1,
			RepoName:   "test/repo",
			AuthorID:   "user_001",
			State:      models.PRStateMerged,
			HeadBranch: "feature-1",
			CreatedAt:  baseTime.Add(-1 * time.Hour),
			MergedAt:   &mergeTime,
		})
		return store
	}

	seedData := &seed.SeedData{
		Repositories: []seed.Repository{{RepoName: "test/repo"}},
	}

	// Generate dataset twice with same seed
	store1 := setupStore()
	gen1 := NewResearchGeneratorWithSeed(seedData, store1, 42)
	from := time.Date(2026, 1, 15, 8, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	data1, err := gen1.GenerateDataset(from, to)
	require.NoError(t, err)

	store2 := setupStore()
	gen2 := NewResearchGeneratorWithSeed(seedData, store2, 42)
	data2, err := gen2.GenerateDataset(from, to)
	require.NoError(t, err)

	assert.Equal(t, len(data1), len(data2), "same seed should produce same number of data points")
}
