package generator

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

func TestPRGenerator_GroupCommitsIntoPRs(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	developer := seed.Developer{
		UserID:    "user_001",
		Email:     "alice@example.com",
		Name:      "Alice",
		Seniority: "mid",
		WorkingHoursBand: seed.WorkingHours{
			Start: 9,
			End:   18,
		},
	}

	repository := seed.Repository{
		RepoName: "acme/platform",
	}

	// Create commits on same branch
	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	commits := []models.Commit{
		{
			CommitHash:         "abc123",
			UserID:             "user_001",
			RepoName:           "acme/platform",
			BranchName:         "feature/auth",
			TotalLinesAdded:    50,
			TotalLinesDeleted:  10,
			TabLinesAdded:      20,
			ComposerLinesAdded: 10,
			CommitTs:           baseTime,
		},
		{
			CommitHash:         "def456",
			UserID:             "user_001",
			RepoName:           "acme/platform",
			BranchName:         "feature/auth",
			TotalLinesAdded:    30,
			TotalLinesDeleted:  5,
			TabLinesAdded:      15,
			ComposerLinesAdded: 5,
			CommitTs:           baseTime.Add(10 * time.Minute), // 10 min gap (< min inactivity 15min)
		},
		{
			CommitHash:         "ghi789",
			UserID:             "user_001",
			RepoName:           "acme/platform",
			BranchName:         "feature/auth",
			TotalLinesAdded:    40,
			TotalLinesDeleted:  8,
			TabLinesAdded:      18,
			ComposerLinesAdded: 8,
			CommitTs:           baseTime.Add(20 * time.Minute), // 10 min gap (< min inactivity 15min)
		},
	}

	gen := NewPRGenerator(&seed.SeedData{
		Developers:   []seed.Developer{developer},
		Repositories: []seed.Repository{repository},
	}, rng)

	prs := gen.GroupCommitsIntoPRs(commits)

	// Should create 1 PR from these 3 commits (all within inactivity gap)
	if len(prs) != 1 {
		t.Fatalf("Expected 1 PR, got %d", len(prs))
	}

	pr := prs[0]

	// Verify PR metadata
	if pr.RepoName != "acme/platform" {
		t.Errorf("PR.RepoName = %s, want acme/platform", pr.RepoName)
	}

	if pr.HeadBranch != "feature/auth" {
		t.Errorf("PR.HeadBranch = %s, want feature/auth", pr.HeadBranch)
	}

	if pr.AuthorID != "user_001" {
		t.Errorf("PR.AuthorID = %s, want user_001", pr.AuthorID)
	}

	// Verify PR aggregates
	expectedAdditions := 50 + 30 + 40 // 120
	if pr.Additions != expectedAdditions {
		t.Errorf("PR.Additions = %d, want %d", pr.Additions, expectedAdditions)
	}

	expectedDeletions := 10 + 5 + 8 // 23
	if pr.Deletions != expectedDeletions {
		t.Errorf("PR.Deletions = %d, want %d", pr.Deletions, expectedDeletions)
	}

	expectedCommits := 3
	if pr.CommitCount != expectedCommits {
		t.Errorf("PR.CommitCount = %d, want %d", pr.CommitCount, expectedCommits)
	}

	// Verify AI metrics
	totalTabLines := 20 + 15 + 18   // 53
	totalComposerLines := 10 + 5 + 8 // 23
	totalAILines := totalTabLines + totalComposerLines // 76
	expectedAIRatio := float64(totalAILines) / float64(expectedAdditions) // 76/120 = 0.633

	if pr.TabLines != totalTabLines {
		t.Errorf("PR.TabLines = %d, want %d", pr.TabLines, totalTabLines)
	}

	if pr.ComposerLines != totalComposerLines {
		t.Errorf("PR.ComposerLines = %d, want %d", pr.ComposerLines, totalComposerLines)
	}

	// Check AI ratio (allow small floating point error)
	if pr.AIRatio < expectedAIRatio-0.01 || pr.AIRatio > expectedAIRatio+0.01 {
		t.Errorf("PR.AIRatio = %.3f, want %.3f", pr.AIRatio, expectedAIRatio)
	}
}
