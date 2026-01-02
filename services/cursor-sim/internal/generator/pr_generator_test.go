package generator

import (
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPRStore implements PRStore for testing
type mockPRStore struct {
	commits        []models.Commit
	prs            []models.PullRequest
	nextPRNumbers  map[string]int
	reviewComments []models.ReviewComment
}

func newMockPRStore() *mockPRStore {
	return &mockPRStore{
		commits:        []models.Commit{},
		prs:            []models.PullRequest{},
		nextPRNumbers:  make(map[string]int),
		reviewComments: []models.ReviewComment{},
	}
}

func (m *mockPRStore) GetCommitsByTimeRange(from, to time.Time) []models.Commit {
	result := []models.Commit{}
	for _, c := range m.commits {
		if !c.CommitTs.Before(from) && !c.CommitTs.After(to) {
			result = append(result, c)
		}
	}
	return result
}

func (m *mockPRStore) GetDeveloper(userID string) (*seed.Developer, error) {
	return &seed.Developer{
		UserID: userID,
		Email:  userID + "@example.com",
		Name:   "Test User",
		Team:   "Platform",
	}, nil
}

func (m *mockPRStore) AddPR(pr models.PullRequest) error {
	m.prs = append(m.prs, pr)
	return nil
}

func (m *mockPRStore) GetNextPRNumber(repoName string) int {
	num, ok := m.nextPRNumbers[repoName]
	if !ok {
		num = 1
	}
	m.nextPRNumbers[repoName] = num + 1
	return num
}

func (m *mockPRStore) AddReviewComment(comment models.ReviewComment) error {
	m.reviewComments = append(m.reviewComments, comment)
	return nil
}

func TestPRGenerator_ClusterCommits(t *testing.T) {
	store := newMockPRStore()
	now := time.Now()

	// Add commits that should cluster into 2 PRs:
	// - PR 1: user_001 on feature/auth with 3 commits
	// - PR 2: user_002 on feature/api with 2 commits
	store.commits = []models.Commit{
		{CommitHash: "a1", UserID: "user_001", BranchName: "feature/auth", RepoName: "acme/api", CommitTs: now.Add(-2 * time.Hour)},
		{CommitHash: "a2", UserID: "user_001", BranchName: "feature/auth", RepoName: "acme/api", CommitTs: now.Add(-1 * time.Hour)},
		{CommitHash: "a3", UserID: "user_001", BranchName: "feature/auth", RepoName: "acme/api", CommitTs: now},
		{CommitHash: "b1", UserID: "user_002", BranchName: "feature/api", RepoName: "acme/api", CommitTs: now.Add(-30 * time.Minute)},
		{CommitHash: "b2", UserID: "user_002", BranchName: "feature/api", RepoName: "acme/api", CommitTs: now},
	}

	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			CycleTimes: seed.CycleTimes{
				ReviewLeadTime: seed.TimeDistribution{
					Params: seed.StatParams{Mean: 8, Std: 4},
				},
			},
		},
		TextTemplates: seed.TextTemplates{
			PRTitles: []string{"feat: add new feature"},
		},
	}

	gen := NewPRGenerator(seedData, store)
	clusters := gen.ClusterCommits(store.commits)

	assert.Len(t, clusters, 2, "should create 2 clusters")

	// Verify cluster sizes
	clusterSizes := map[string]int{}
	for _, cluster := range clusters {
		key := cluster.AuthorID + ":" + cluster.Branch
		clusterSizes[key] = len(cluster.Commits)
	}

	assert.Equal(t, 3, clusterSizes["user_001:feature/auth"])
	assert.Equal(t, 2, clusterSizes["user_002:feature/api"])
}

func TestPRGenerator_GeneratePRsFromCommits(t *testing.T) {
	store := newMockPRStore()
	now := time.Now()

	store.commits = []models.Commit{
		{
			CommitHash:         "a1",
			UserID:             "user_001",
			UserEmail:          "user@example.com",
			UserName:           "Test User",
			BranchName:         "feature/auth",
			RepoName:           "acme/api",
			TotalLinesAdded:    100,
			TotalLinesDeleted:  20,
			TabLinesAdded:      60,
			ComposerLinesAdded: 20,
			CommitTs:           now.Add(-2 * time.Hour),
		},
		{
			CommitHash:         "a2",
			UserID:             "user_001",
			UserEmail:          "user@example.com",
			UserName:           "Test User",
			BranchName:         "feature/auth",
			RepoName:           "acme/api",
			TotalLinesAdded:    50,
			TotalLinesDeleted:  10,
			TabLinesAdded:      30,
			ComposerLinesAdded: 10,
			CommitTs:           now.Add(-1 * time.Hour),
		},
	}

	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			CycleTimes: seed.CycleTimes{
				ReviewLeadTime: seed.TimeDistribution{
					Params: seed.StatParams{Mean: 8, Std: 4},
				},
			},
		},
		TextTemplates: seed.TextTemplates{
			PRTitles: []string{"feat: {description}"},
		},
	}

	gen := NewPRGenerator(seedData, store)
	err := gen.GeneratePRsFromCommits(now.Add(-24*time.Hour), now)
	require.NoError(t, err)

	require.Len(t, store.prs, 1, "should create 1 PR")

	pr := store.prs[0]
	assert.Equal(t, "acme/api", pr.RepoName)
	assert.Equal(t, "user_001", pr.AuthorID)
	assert.Equal(t, "feature/auth", pr.HeadBranch)
	assert.Equal(t, 2, pr.CommitCount)

	// Verify aggregated metrics
	assert.Equal(t, 150, pr.Additions)
	assert.Equal(t, 30, pr.Deletions)
	assert.Equal(t, 90, pr.TabLines)
	assert.Equal(t, 30, pr.ComposerLines)

	// AI ratio should be (90+30)/150 = 0.8
	assert.InDelta(t, 0.8, pr.AIRatio, 0.01)
}

func TestPRGenerator_PRState(t *testing.T) {
	store := newMockPRStore()
	now := time.Now()

	// Old commits should result in merged PR
	store.commits = []models.Commit{
		{
			CommitHash: "old1",
			UserID:     "user_001",
			BranchName: "feature/old",
			RepoName:   "acme/api",
			CommitTs:   now.Add(-72 * time.Hour), // 3 days old
		},
	}

	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			CycleTimes: seed.CycleTimes{
				ReviewLeadTime: seed.TimeDistribution{
					Params: seed.StatParams{Mean: 8, Std: 4},
				},
			},
		},
	}

	gen := NewPRGeneratorWithSeed(seedData, store, 42)
	err := gen.GeneratePRsFromCommits(now.Add(-7*24*time.Hour), now)
	require.NoError(t, err)

	require.Len(t, store.prs, 1)
	// Older PRs should be merged
	assert.Equal(t, models.PRStateMerged, store.prs[0].State)
	assert.NotNil(t, store.prs[0].MergedAt)
}

func TestPRGenerator_AIRatioCalculation(t *testing.T) {
	tests := []struct {
		name            string
		tabLines        int
		composerLines   int
		totalLines      int
		expectedAIRatio float64
	}{
		{"100% AI", 70, 30, 100, 1.0},
		{"50% AI", 30, 20, 100, 0.5},
		{"0% AI", 0, 0, 100, 0.0},
		{"No lines", 0, 0, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commits := []models.Commit{
				{
					TotalLinesAdded:    tt.totalLines,
					TabLinesAdded:      tt.tabLines,
					ComposerLinesAdded: tt.composerLines,
				},
			}

			store := newMockPRStore()
			gen := NewPRGeneratorWithSeed(nil, store, 42)
			pr := gen.createPRFromCluster(CommitCluster{
				Commits:  commits,
				AuthorID: "user_001",
				Branch:   "feature/test",
				RepoName: "test/repo",
			}, 1)

			assert.InDelta(t, tt.expectedAIRatio, pr.AIRatio, 0.01)
		})
	}
}

func TestPRGenerator_Reproducibility(t *testing.T) {
	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			CycleTimes: seed.CycleTimes{
				ReviewLeadTime: seed.TimeDistribution{
					Params: seed.StatParams{Mean: 8, Std: 4},
				},
			},
		},
		TextTemplates: seed.TextTemplates{
			PRTitles: []string{"feat: title 1", "fix: title 2", "chore: title 3"},
		},
	}

	// Same seed should produce same results
	store1 := newMockPRStore()
	store2 := newMockPRStore()

	now := time.Now()
	commits := []models.Commit{
		{CommitHash: "a1", UserID: "user_001", BranchName: "feature/auth", RepoName: "acme/api", CommitTs: now},
	}
	store1.commits = commits
	store2.commits = commits

	gen1 := NewPRGeneratorWithSeed(seedData, store1, 12345)
	gen2 := NewPRGeneratorWithSeed(seedData, store2, 12345)

	_ = gen1.GeneratePRsFromCommits(now.Add(-time.Hour), now)
	_ = gen2.GeneratePRsFromCommits(now.Add(-time.Hour), now)

	require.Len(t, store1.prs, 1)
	require.Len(t, store2.prs, 1)

	// Same seed should produce same PR title
	assert.Equal(t, store1.prs[0].Title, store2.prs[0].Title)
}

func TestPRGenerator_TimeWindowClustering(t *testing.T) {
	store := newMockPRStore()
	now := time.Now()

	// Commits far apart should be in separate PRs
	store.commits = []models.Commit{
		{CommitHash: "a1", UserID: "user_001", BranchName: "feature/auth", RepoName: "acme/api", CommitTs: now.Add(-48 * time.Hour)},
		{CommitHash: "a2", UserID: "user_001", BranchName: "feature/auth", RepoName: "acme/api", CommitTs: now}, // 48 hours later
	}

	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			CycleTimes: seed.CycleTimes{
				CodingLeadTime: seed.TimeDistribution{
					Params: seed.StatParams{Mean: 4, Std: 2},
				},
			},
		},
	}

	gen := NewPRGenerator(seedData, store)
	clusters := gen.ClusterCommits(store.commits)

	// Commits > 24 hours apart on same branch should be different PRs
	assert.GreaterOrEqual(t, len(clusters), 1, "should create at least 1 cluster")
}
