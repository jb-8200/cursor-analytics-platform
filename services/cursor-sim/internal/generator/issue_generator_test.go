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

func TestNewIssueGenerator(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{UserID: "user-1", Email: "dev1@example.com", Name: "Dev One"},
			{UserID: "user-2", Email: "dev2@example.com", Name: "Dev Two"},
		},
		TextTemplates: seed.TextTemplates{
			PRTitles: []string{"Fix bug", "Add feature"},
		},
	}

	rng := rand.New(rand.NewSource(12345))

	gen := NewIssueGenerator(seedData, rng)
	assert.NotNil(t, gen)
	assert.Equal(t, seedData, gen.seed)
	assert.NotNil(t, gen.rng)
}

func TestGenerateIssuesForPRs_40PercentOfMergedPRs(t *testing.T) {
	seedData := createTestSeedData()
	rng := rand.New(rand.NewSource(12345))
	gen := NewIssueGenerator(seedData, rng)

	mergedPRs := make([]models.PullRequest, 100)
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 100; i++ {
		mergedAt := baseTime.Add(time.Duration(i) * time.Hour)
		mergedPRs[i] = models.PullRequest{
			Number:    i + 1,
			Title:     "Fix bug",
			State:     models.PRStateMerged,
			AuthorID:  "dev1@example.com",
			RepoName:  "test-repo",
			CreatedAt: baseTime.Add(time.Duration(i) * time.Hour),
			MergedAt:  &mergedAt,
		}
	}

	issues := gen.GenerateIssuesForPRs(mergedPRs, "test-repo")

	assert.GreaterOrEqual(t, len(issues), 30, "Should generate at least 30 issues")
	assert.LessOrEqual(t, len(issues), 50, "Should generate at most 50 issues")
}

func TestGenerateIssuesForPRs_10PercentRemainOpen(t *testing.T) {
	seedData := createTestSeedData()
	rng := rand.New(rand.NewSource(12345))
	gen := NewIssueGenerator(seedData, rng)

	mergedPRs := make([]models.PullRequest, 100)
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 100; i++ {
		mergedAt := baseTime.Add(time.Duration(i) * time.Hour)
		mergedPRs[i] = models.PullRequest{
			Number:    i + 1,
			Title:     "Fix bug",
			State:     models.PRStateMerged,
			AuthorID:  "dev1@example.com",
			RepoName:  "test-repo",
			CreatedAt: baseTime.Add(time.Duration(i) * time.Hour),
			MergedAt:  &mergedAt,
		}
	}

	issues := gen.GenerateIssuesForPRs(mergedPRs, "test-repo")

	openCount := 0
	for _, issue := range issues {
		if issue.State == models.IssueStateOpen {
			openCount++
		}
	}

	totalIssues := len(issues)
	if totalIssues > 0 {
		openPercentage := float64(openCount) / float64(totalIssues)
		assert.GreaterOrEqual(t, openPercentage, 0.01, "At least 1% should be open")
		assert.LessOrEqual(t, openPercentage, 0.25, "At most 25% should be open")
	}
}

func TestGenerateIssuesForPRs_IssueCreatedBeforePR(t *testing.T) {
	seedData := createTestSeedData()
	rng := rand.New(rand.NewSource(12345))
	gen := NewIssueGenerator(seedData, rng)

	baseTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	mergedAt := baseTime.Add(24 * time.Hour)
	prs := []models.PullRequest{
		{
			Number:    1,
			Title:     "Fix authentication bug",
			State:     models.PRStateMerged,
			AuthorID:  "dev1@example.com",
			RepoName:  "test-repo",
			CreatedAt: baseTime,
			MergedAt:  &mergedAt,
		},
	}

	issues := gen.GenerateIssuesForPRs(prs, "test-repo")

	for _, issue := range issues {
		if issue.ClosedByPRID != nil {
			for _, pr := range prs {
				if pr.Number == *issue.ClosedByPRID {
					assert.True(t, issue.CreatedAt.Before(pr.CreatedAt),
						"Issue created_at (%v) should be before PR created_at (%v)",
						issue.CreatedAt, pr.CreatedAt)
				}
			}
		}
	}
}

func TestGenerateIssuesForPRs_LabelsAssigned(t *testing.T) {
	seedData := createTestSeedData()
	rng := rand.New(rand.NewSource(12345))
	gen := NewIssueGenerator(seedData, rng)

	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	mergedAt := baseTime.Add(24 * time.Hour)
	prs := []models.PullRequest{
		{
			Number:    1,
			Title:     "Fix bug",
			State:     models.PRStateMerged,
			AuthorID:  "dev1@example.com",
			RepoName:  "test-repo",
			CreatedAt: baseTime,
			MergedAt:  &mergedAt,
		},
	}

	issues := gen.GenerateIssuesForPRs(prs, "test-repo")

	validLabels := map[string]bool{
		"bug":         true,
		"feature":     true,
		"enhancement": true,
	}

	for _, issue := range issues {
		assert.NotEmpty(t, issue.Labels, "Issue should have at least one label")

		for _, label := range issue.Labels {
			assert.True(t, validLabels[label],
				"Label '%s' should be one of: bug, feature, enhancement", label)
		}
	}
}

func TestGenerateIssuesForPRs_OnlyMergedPRs(t *testing.T) {
	seedData := createTestSeedData()
	rng := rand.New(rand.NewSource(12345))
	gen := NewIssueGenerator(seedData, rng)

	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	mergedAt := baseTime.Add(24 * time.Hour)
	closedAt := baseTime.Add(12 * time.Hour)

	prs := []models.PullRequest{
		{
			Number:    1,
			Title:     "Fix bug",
			State:     models.PRStateMerged,
			AuthorID:  "dev1@example.com",
			RepoName:  "test-repo",
			CreatedAt: baseTime,
			MergedAt:  &mergedAt,
		},
		{
			Number:    2,
			Title:     "Add feature",
			State:     models.PRStateOpen,
			AuthorID:  "dev2@example.com",
			RepoName:  "test-repo",
			CreatedAt: baseTime,
		},
		{
			Number:    3,
			Title:     "Refactor code",
			State:     models.PRStateClosed,
			AuthorID:  "dev1@example.com",
			RepoName:  "test-repo",
			CreatedAt: baseTime,
			ClosedAt:  &closedAt,
		},
	}

	issues := gen.GenerateIssuesForPRs(prs, "test-repo")

	for _, issue := range issues {
		if issue.ClosedByPRID != nil {
			assert.Equal(t, 1, *issue.ClosedByPRID,
				"Issues should only reference merged PRs")
		}
	}
}

func TestGenerateIssuesForPRs_EmptyPRList(t *testing.T) {
	seedData := createTestSeedData()
	rng := rand.New(rand.NewSource(12345))
	gen := NewIssueGenerator(seedData, rng)

	issues := gen.GenerateIssuesForPRs([]models.PullRequest{}, "test-repo")
	assert.Empty(t, issues, "Should return empty slice for empty PR list")
}

func TestGenerateIssuesForPRs_UniqueIssueNumbers(t *testing.T) {
	seedData := createTestSeedData()
	rng := rand.New(rand.NewSource(12345))
	gen := NewIssueGenerator(seedData, rng)

	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	mergedPRs := make([]models.PullRequest, 50)
	for i := 0; i < 50; i++ {
		mergedAt := baseTime.Add(time.Duration(i) * time.Hour)
		mergedPRs[i] = models.PullRequest{
			Number:    i + 1,
			Title:     "Fix bug",
			State:     models.PRStateMerged,
			AuthorID:  "dev1@example.com",
			RepoName:  "test-repo",
			CreatedAt: baseTime.Add(time.Duration(i) * time.Hour),
			MergedAt:  &mergedAt,
		}
	}

	issues := gen.GenerateIssuesForPRs(mergedPRs, "test-repo")

	seenNumbers := make(map[int]bool)
	for _, issue := range issues {
		assert.False(t, seenNumbers[issue.Number],
			"Issue number %d is duplicated", issue.Number)
		seenNumbers[issue.Number] = true
	}
}

func TestGenerateIssuesForPRs_ReproducibilityWithSeed(t *testing.T) {
	seedData := createTestSeedData()

	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	mergedPRs := make([]models.PullRequest, 20)
	for i := 0; i < 20; i++ {
		mergedAt := baseTime.Add(time.Duration(i) * time.Hour)
		mergedPRs[i] = models.PullRequest{
			Number:    i + 1,
			Title:     "Fix bug",
			State:     models.PRStateMerged,
			AuthorID:  "dev1@example.com",
			RepoName:  "test-repo",
			CreatedAt: baseTime.Add(time.Duration(i) * time.Hour),
			MergedAt:  &mergedAt,
		}
	}

	gen1 := NewIssueGenerator(seedData, rand.New(rand.NewSource(12345)))
	issues1 := gen1.GenerateIssuesForPRs(mergedPRs, "test-repo")

	gen2 := NewIssueGenerator(seedData, rand.New(rand.NewSource(12345)))
	issues2 := gen2.GenerateIssuesForPRs(mergedPRs, "test-repo")

	require.Equal(t, len(issues1), len(issues2), "Should generate same number of issues")

	for i := range issues1 {
		assert.Equal(t, issues1[i].Number, issues2[i].Number)
		assert.Equal(t, issues1[i].Title, issues2[i].Title)
		assert.Equal(t, issues1[i].State, issues2[i].State)
		assert.Equal(t, issues1[i].Labels, issues2[i].Labels)
	}
}

func TestGenerateIssuesForPRs_ValidationPasses(t *testing.T) {
	seedData := createTestSeedData()
	rng := rand.New(rand.NewSource(12345))
	gen := NewIssueGenerator(seedData, rng)

	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	mergedAt := baseTime.Add(24 * time.Hour)
	prs := []models.PullRequest{
		{
			Number:    1,
			Title:     "Fix bug",
			State:     models.PRStateMerged,
			AuthorID:  "dev1@example.com",
			RepoName:  "test-repo",
			CreatedAt: baseTime,
			MergedAt:  &mergedAt,
		},
	}

	issues := gen.GenerateIssuesForPRs(prs, "test-repo")

	for _, issue := range issues {
		err := issue.Validate()
		assert.NoError(t, err, "Issue should pass validation: %+v", issue)
	}
}

func TestGenerateIssuesForPRs_ClosedAtSetForClosedIssues(t *testing.T) {
	seedData := createTestSeedData()
	rng := rand.New(rand.NewSource(54321)) // Different seed to get closed issues
	gen := NewIssueGenerator(seedData, rng)

	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Generate many PRs to get some closed issues
	mergedPRs := make([]models.PullRequest, 50)
	for i := 0; i < 50; i++ {
		mergedAt := baseTime.Add(time.Duration(i+1) * 24 * time.Hour)
		mergedPRs[i] = models.PullRequest{
			Number:    i + 1,
			Title:     "Fix bug",
			State:     models.PRStateMerged,
			AuthorID:  "dev1@example.com",
			RepoName:  "test-repo",
			CreatedAt: baseTime.Add(time.Duration(i) * time.Hour),
			MergedAt:  &mergedAt,
		}
	}

	issues := gen.GenerateIssuesForPRs(mergedPRs, "test-repo")

	for _, issue := range issues {
		if issue.State == models.IssueStateClosed {
			assert.NotNil(t, issue.ClosedAt, "Closed issue should have ClosedAt timestamp")
			assert.NotNil(t, issue.ClosedByPRID, "Closed issue should have ClosedByPRID")

			assert.True(t, issue.ClosedAt.After(issue.CreatedAt),
				"ClosedAt (%v) should be after CreatedAt (%v)",
				issue.ClosedAt, issue.CreatedAt)
		} else if issue.State == models.IssueStateOpen {
			assert.Nil(t, issue.ClosedAt, "Open issue should not have ClosedAt timestamp")
		}
	}
}

func TestGenerateIssuesForPRs_AuthorFromDevelopers(t *testing.T) {
	seedData := createTestSeedData()
	rng := rand.New(rand.NewSource(12345))
	gen := NewIssueGenerator(seedData, rng)

	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	mergedAt := baseTime.Add(24 * time.Hour)
	prs := []models.PullRequest{
		{
			Number:    1,
			Title:     "Fix bug",
			State:     models.PRStateMerged,
			AuthorID:  "dev1@example.com",
			RepoName:  "test-repo",
			CreatedAt: baseTime,
			MergedAt:  &mergedAt,
		},
	}

	issues := gen.GenerateIssuesForPRs(prs, "test-repo")

	validAuthors := make(map[string]bool)
	for _, dev := range seedData.Developers {
		validAuthors[dev.Email] = true
	}

	for _, issue := range issues {
		assert.True(t, validAuthors[issue.AuthorID],
			"Issue author '%s' should be from seed developers", issue.AuthorID)
	}
}

func createTestSeedData() *seed.SeedData {
	return &seed.SeedData{
		Developers: []seed.Developer{
			{UserID: "user-1", Email: "dev1@example.com", Name: "Dev One"},
			{UserID: "user-2", Email: "dev2@example.com", Name: "Dev Two"},
			{UserID: "user-3", Email: "dev3@example.com", Name: "Dev Three"},
		},
		TextTemplates: seed.TextTemplates{
			PRTitles: []string{
				"Fix authentication bug",
				"Add user profile feature",
				"Refactor database queries",
				"Update API endpoints",
			},
		},
	}
}
