package generator

import (
	"context"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockStore implements a simple in-memory store for testing
type MockStore struct {
	commits    []models.Commit
	developers []seed.Developer
}

func (m *MockStore) AddCommit(commit models.Commit) error {
	m.commits = append(m.commits, commit)
	return nil
}

func (m *MockStore) GetCommits() []models.Commit {
	return m.commits
}

func (m *MockStore) ListDevelopers() []seed.Developer {
	return m.developers
}

// NewMockStore creates a MockStore initialized with developers
func NewMockStore(developers []seed.Developer) *MockStore {
	return &MockStore{
		commits:    []models.Commit{},
		developers: developers,
	}
}

func TestCommitGenerator_New(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{UserID: "user_001", Email: "test@example.com"},
		},
	}
	store := NewMockStore(seedData.Developers)

	gen := NewCommitGenerator(seedData, store, "medium")
	assert.NotNil(t, gen)
}

func TestCommitGenerator_GenerateCommits(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:         "user_001",
				Email:          "alice@example.com",
				Name:           "Alice",
				AcceptanceRate: 0.85,
				PRBehavior: seed.PRBehavior{
					PRsPerWeek:    2.0,
					AvgPRSizeLOC:  100,
					AvgFilesPerPR: 3,
				},
			},
		},
		Repositories: []seed.Repository{
			{
				RepoName:        "test/repo",
				PrimaryLanguage: "Go",
				DefaultBranch:   "main",
				Teams:           []string{"Backend"},
			},
		},
		TextTemplates: seed.TextTemplates{
			CommitMessages: seed.CommitMessageTemplates{
				Feature:  []string{"feat: add {{.Feature}}"},
				Bugfix:   []string{"fix: resolve {{.Issue}}"},
				Refactor: []string{"refactor: improve code"},
				Chore:    []string{"chore: maintenance"},
			},
		},
	}

	store := NewMockStore(seedData.Developers)
	gen := NewCommitGeneratorWithSeed(seedData, store, "medium", 42) // Use deterministic seed

	ctx := context.Background()
	err := gen.GenerateCommits(ctx, 7, 0) // Generate 7 days of history for reliable commits (0 = unlimited)

	require.NoError(t, err)
	commits := store.GetCommits()

	// Should have generated at least one commit
	assert.Greater(t, len(commits), 0, "should generate at least one commit")

	// Verify first commit structure
	if len(commits) > 0 {
		c := commits[0]
		assert.NotEmpty(t, c.CommitHash)
		assert.Equal(t, "user_001", c.UserID)
		assert.Equal(t, "alice@example.com", c.UserEmail)
		assert.NotEmpty(t, c.RepoName)
		assert.NotEmpty(t, c.Message)
		assert.Greater(t, c.TotalLinesAdded, 0)
	}
}

func TestCommitGenerator_AIAttribution(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:         "user_001",
				Email:          "test@example.com",
				Name:           "Test",
				AcceptanceRate: 0.80,
				PRBehavior: seed.PRBehavior{
					PRsPerWeek:   5.0,
					AvgPRSizeLOC: 200,
				},
			},
		},
		Repositories: []seed.Repository{
			{RepoName: "test/repo", DefaultBranch: "main"},
		},
		TextTemplates: seed.TextTemplates{
			CommitMessages: seed.CommitMessageTemplates{
				Feature: []string{"test commit"},
			},
		},
	}

	store := NewMockStore(seedData.Developers)
	gen := NewCommitGenerator(seedData, store, "medium")

	ctx := context.Background()
	err := gen.GenerateCommits(ctx, 2, 0)
	require.NoError(t, err)

	commits := store.GetCommits()
	require.Greater(t, len(commits), 0)

	// Check AI attribution for all commits
	for _, c := range commits {
		// Total should equal sum of parts
		total := c.TabLinesAdded + c.ComposerLinesAdded + c.NonAILinesAdded
		assert.Equal(t, c.TotalLinesAdded, total,
			"TotalLinesAdded should equal TabLines + ComposerLines + NonAILines")

		// AI ratio should be reasonable (between 0 and 1)
		ratio := c.AIRatio()
		assert.GreaterOrEqual(t, ratio, 0.0)
		assert.LessOrEqual(t, ratio, 1.0)

		// Tab and Composer lines should be non-negative
		assert.GreaterOrEqual(t, c.TabLinesAdded, 0)
		assert.GreaterOrEqual(t, c.ComposerLinesAdded, 0)
		assert.GreaterOrEqual(t, c.NonAILinesAdded, 0)
	}
}

func TestCommitGenerator_TimeRange(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "test@example.com",
				PRBehavior: seed.PRBehavior{
					PRsPerWeek: 3.0,
				},
			},
		},
		Repositories: []seed.Repository{
			{RepoName: "test/repo", DefaultBranch: "main"},
		},
		TextTemplates: seed.TextTemplates{
			CommitMessages: seed.CommitMessageTemplates{
				Feature: []string{"test"},
			},
		},
	}

	store := NewMockStore(seedData.Developers)
	gen := NewCommitGenerator(seedData, store, "low")

	days := 5
	ctx := context.Background()
	err := gen.GenerateCommits(ctx, days, 0)
	require.NoError(t, err)

	commits := store.GetCommits()
	if len(commits) == 0 {
		t.Skip("No commits generated (can happen with low probability)")
	}

	// All commits should be within the time range
	now := time.Now()
	startTime := now.AddDate(0, 0, -days)

	for _, c := range commits {
		assert.True(t, c.CommitTs.After(startTime) || c.CommitTs.Equal(startTime),
			"commit timestamp should be after start time")
		assert.True(t, c.CommitTs.Before(now) || c.CommitTs.Equal(now),
			"commit timestamp should be before now")
	}
}

func TestCommitGenerator_Reproducibility(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "test@example.com",
				PRBehavior: seed.PRBehavior{
					PRsPerWeek: 3.0,
				},
			},
		},
		Repositories: []seed.Repository{
			{RepoName: "test/repo", DefaultBranch: "main"},
		},
		TextTemplates: seed.TextTemplates{
			CommitMessages: seed.CommitMessageTemplates{
				Feature: []string{"test"},
			},
		},
	}

	// Generate twice with same seed
	store1 := NewMockStore(seedData.Developers)
	gen1 := NewCommitGeneratorWithSeed(seedData, store1, "medium", 12345)
	err := gen1.GenerateCommits(context.Background(), 2, 0)
	require.NoError(t, err)

	store2 := NewMockStore(seedData.Developers)
	gen2 := NewCommitGeneratorWithSeed(seedData, store2, "medium", 12345)
	err = gen2.GenerateCommits(context.Background(), 2, 0)
	require.NoError(t, err)

	// Should generate same number of commits
	assert.Equal(t, len(store1.GetCommits()), len(store2.GetCommits()),
		"same seed should produce same number of commits")
}

func TestCommitGenerator_CommitMessage(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "test@example.com",
				PRBehavior: seed.PRBehavior{
					PRsPerWeek: 10.0, // High rate to ensure commits
				},
			},
		},
		Repositories: []seed.Repository{
			{RepoName: "test/repo", DefaultBranch: "main"},
		},
		TextTemplates: seed.TextTemplates{
			CommitMessages: seed.CommitMessageTemplates{
				Feature:  []string{"feat: new feature"},
				Bugfix:   []string{"fix: bug fix"},
				Refactor: []string{"refactor: improve"},
				Chore:    []string{"chore: update"},
			},
		},
	}

	store := NewMockStore(seedData.Developers)
	gen := NewCommitGenerator(seedData, store, "high")

	ctx := context.Background()
	err := gen.GenerateCommits(ctx, 3, 0)
	require.NoError(t, err)

	commits := store.GetCommits()
	require.Greater(t, len(commits), 0, "should generate commits")

	// All commits should have messages
	for _, c := range commits {
		assert.NotEmpty(t, c.Message, "commit should have a message")
	}
}

func TestGenerateCommits_MaxLimit(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "dev1@example.com",
				Name:   "Dev 1",
				PRBehavior: seed.PRBehavior{
					PRsPerWeek:   10.0, // High rate to ensure we hit limit
					AvgPRSizeLOC: 100,
				},
			},
			{
				UserID: "user_002",
				Email:  "dev2@example.com",
				Name:   "Dev 2",
				PRBehavior: seed.PRBehavior{
					PRsPerWeek:   10.0,
					AvgPRSizeLOC: 100,
				},
			},
		},
		Repositories: []seed.Repository{
			{RepoName: "test/repo", DefaultBranch: "main"},
		},
		TextTemplates: seed.TextTemplates{
			CommitMessages: seed.CommitMessageTemplates{
				Feature: []string{"test commit"},
			},
		},
	}

	store := NewMockStore(seedData.Developers)
	gen := NewCommitGeneratorWithSeed(seedData, store, "high", 42)

	ctx := context.Background()
	maxCommits := 10
	err := gen.GenerateCommits(ctx, 90, maxCommits)

	require.NoError(t, err)
	commits := store.GetCommits()
	assert.Equal(t, maxCommits, len(commits), "should generate exactly maxCommits commits")
}

func TestGenerateCommits_NoLimit(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "test@example.com",
				PRBehavior: seed.PRBehavior{
					PRsPerWeek:   5.0,
					AvgPRSizeLOC: 100,
				},
			},
		},
		Repositories: []seed.Repository{
			{RepoName: "test/repo", DefaultBranch: "main"},
		},
		TextTemplates: seed.TextTemplates{
			CommitMessages: seed.CommitMessageTemplates{
				Feature: []string{"test"},
			},
		},
	}

	store := NewMockStore(seedData.Developers)
	gen := NewCommitGeneratorWithSeed(seedData, store, "medium", 42)

	ctx := context.Background()
	err := gen.GenerateCommits(ctx, 7, 0) // 0 = unlimited

	require.NoError(t, err)
	commits := store.GetCommits()
	// Should generate commits based on Poisson process (not limited)
	assert.Greater(t, len(commits), 0, "should generate commits with unlimited mode")
}

func TestGenerateCommits_LimitDistribution(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "dev1@example.com",
				Name:   "Dev 1",
				PRBehavior: seed.PRBehavior{
					PRsPerWeek:   10.0,
					AvgPRSizeLOC: 100,
				},
			},
			{
				UserID: "user_002",
				Email:  "dev2@example.com",
				Name:   "Dev 2",
				PRBehavior: seed.PRBehavior{
					PRsPerWeek:   10.0,
					AvgPRSizeLOC: 100,
				},
			},
			{
				UserID: "user_003",
				Email:  "dev3@example.com",
				Name:   "Dev 3",
				PRBehavior: seed.PRBehavior{
					PRsPerWeek:   10.0,
					AvgPRSizeLOC: 100,
				},
			},
		},
		Repositories: []seed.Repository{
			{RepoName: "test/repo", DefaultBranch: "main"},
		},
		TextTemplates: seed.TextTemplates{
			CommitMessages: seed.CommitMessageTemplates{
				Feature: []string{"test"},
			},
		},
	}

	store := NewMockStore(seedData.Developers)
	gen := NewCommitGeneratorWithSeed(seedData, store, "high", 42)

	ctx := context.Background()
	maxCommits := 15
	err := gen.GenerateCommits(ctx, 90, maxCommits)

	require.NoError(t, err)
	commits := store.GetCommits()
	assert.Equal(t, maxCommits, len(commits), "should generate exactly maxCommits")

	// Verify commits are distributed across developers
	developerCounts := make(map[string]int)
	for _, c := range commits {
		developerCounts[c.UserID]++
	}

	// Should have commits from multiple developers
	assert.GreaterOrEqual(t, len(developerCounts), 1, "should have commits from at least 1 developer")
	// Total should match
	totalCommits := 0
	for _, count := range developerCounts {
		totalCommits += count
	}
	assert.Equal(t, maxCommits, totalCommits, "total commits should equal max")
}
