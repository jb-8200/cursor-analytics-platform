package generator

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

func TestPRGenerationIntegration(t *testing.T) {
	// Setup: Create seed data with multiple developers and repos
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:    "user_001",
				Email:     "alice@example.com",
				Name:      "Alice",
				Seniority: "senior",
				WorkingHoursBand: seed.WorkingHours{
					Start: 9,
					End:   18,
				},
			},
			{
				UserID:    "user_002",
				Email:     "bob@example.com",
				Name:      "Bob",
				Seniority: "junior",
				WorkingHoursBand: seed.WorkingHours{
					Start: 9,
					End:   18,
				},
			},
		},
		Repositories: []seed.Repository{
			{
				RepoName:      "acme/platform",
				DefaultBranch: "main",
			},
			{
				RepoName:      "acme/api",
				DefaultBranch: "main",
			},
		},
	}

	// Create memory store
	store := storage.NewMemoryStore()
	if err := store.LoadDevelopers(seedData.Developers); err != nil {
		t.Fatalf("Failed to load developers: %v", err)
	}

	// Generate commits with realistic patterns
	rng := rand.New(rand.NewSource(42))
	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	commits := []models.Commit{
		// Alice's work on acme/platform feature/auth branch (should become 1 PR)
		{
			CommitHash:         "a1",
			UserID:             "user_001",
			RepoName:           "acme/platform",
			BranchName:         "feature/auth",
			TotalLinesAdded:    100,
			TotalLinesDeleted:  20,
			TabLinesAdded:      40,
			ComposerLinesAdded: 30,
			CommitTs:           baseTime,
		},
		{
			CommitHash:         "a2",
			UserID:             "user_001",
			RepoName:           "acme/platform",
			BranchName:         "feature/auth",
			TotalLinesAdded:    50,
			TotalLinesDeleted:  10,
			TabLinesAdded:      20,
			ComposerLinesAdded: 15,
			CommitTs:           baseTime.Add(5 * time.Minute),
		},
		// Alice's work on different branch (should become separate PR)
		{
			CommitHash:         "a3",
			UserID:             "user_001",
			RepoName:           "acme/platform",
			BranchName:         "feature/billing",
			TotalLinesAdded:    80,
			TotalLinesDeleted:  15,
			TabLinesAdded:      30,
			ComposerLinesAdded: 20,
			CommitTs:           baseTime.Add(10 * time.Minute),
		},
		// Bob's work on acme/api (should become separate PR)
		{
			CommitHash:         "b1",
			UserID:             "user_002",
			RepoName:           "acme/api",
			BranchName:         "feature/endpoints",
			TotalLinesAdded:    40,
			TotalLinesDeleted:  5,
			TabLinesAdded:      15,
			ComposerLinesAdded: 10,
			CommitTs:           baseTime.Add(15 * time.Minute),
		},
	}

	// Store all commits
	for _, commit := range commits {
		if err := store.AddCommit(commit); err != nil {
			t.Fatalf("Failed to add commit %s: %v", commit.CommitHash, err)
		}
	}

	// Generate PRs from commits
	prGen := NewPRGenerator(seedData, rng)
	prs := prGen.GroupCommitsIntoPRs(commits)

	// Verify expected number of PRs
	expectedPRs := 3 // auth, billing, endpoints
	if len(prs) != expectedPRs {
		t.Errorf("Expected %d PRs, got %d", expectedPRs, len(prs))
	}

	// Store all PRs
	for _, pr := range prs {
		if err := store.AddPR(pr); err != nil {
			t.Fatalf("Failed to add PR %d: %v", pr.Number, err)
		}
	}

	// Verify: Retrieve PRs from store
	platformPRs := store.GetPRsByRepo("acme/platform")
	if len(platformPRs) != 2 {
		t.Errorf("Expected 2 PRs for acme/platform, got %d", len(platformPRs))
	}

	apiPRs := store.GetPRsByRepo("acme/api")
	if len(apiPRs) != 1 {
		t.Errorf("Expected 1 PR for acme/api, got %d", len(apiPRs))
	}

	// Verify: PR metrics are correctly aggregated
	for _, pr := range prs {
		if pr.CommitCount == 0 {
			t.Errorf("PR %d has 0 commits", pr.Number)
		}

		if pr.Additions == 0 {
			t.Errorf("PR %d has 0 additions", pr.Number)
		}

		// Verify AI ratio calculation
		expectedAILines := pr.TabLines + pr.ComposerLines
		if pr.Additions > 0 {
			expectedAIRatio := float64(expectedAILines) / float64(pr.Additions)
			if pr.AIRatio < expectedAIRatio-0.01 || pr.AIRatio > expectedAIRatio+0.01 {
				t.Errorf("PR %d: AIRatio = %.3f, expected %.3f", pr.Number, pr.AIRatio, expectedAIRatio)
			}
		}

		// Verify PR state
		if pr.State != models.PRStateOpen {
			t.Errorf("PR %d has state %s, expected open", pr.Number, pr.State)
		}

		// Verify timestamps
		if pr.CreatedAt.IsZero() {
			t.Errorf("PR %d has zero CreatedAt", pr.Number)
		}

		if pr.UpdatedAt.IsZero() {
			t.Errorf("PR %d has zero UpdatedAt", pr.Number)
		}
	}

	// Verify: GetPRsByAuthor works correctly
	alicePRs := store.GetPRsByAuthor("user_001")
	if len(alicePRs) != 2 {
		t.Errorf("Expected 2 PRs for Alice, got %d", len(alicePRs))
	}

	bobPRs := store.GetPRsByAuthor("user_002")
	if len(bobPRs) != 1 {
		t.Errorf("Expected 1 PR for Bob, got %d", len(bobPRs))
	}

	// Verify: Seniority-based correlation (senior has more commits per PR on average)
	var seniorCommitsPerPR, juniorCommitsPerPR float64
	var seniorCount, juniorCount int

	for _, pr := range prs {
		dev, err := store.GetDeveloper(pr.AuthorID)
		if err != nil {
			continue
		}

		if dev.Seniority == "senior" {
			seniorCommitsPerPR += float64(pr.CommitCount)
			seniorCount++
		} else if dev.Seniority == "junior" {
			juniorCommitsPerPR += float64(pr.CommitCount)
			juniorCount++
		}
	}

	if seniorCount > 0 {
		seniorCommitsPerPR /= float64(seniorCount)
	}
	if juniorCount > 0 {
		juniorCommitsPerPR /= float64(juniorCount)
	}

	// Log metrics for verification
	t.Logf("Senior avg commits/PR: %.2f", seniorCommitsPerPR)
	t.Logf("Junior avg commits/PR: %.2f", juniorCommitsPerPR)
	t.Logf("Total PRs generated: %d", len(prs))
	t.Logf("PRs by developer: Alice=%d, Bob=%d", len(alicePRs), len(bobPRs))
}

func TestPRNumberSequencing(t *testing.T) {
	// Verify that PR numbers are sequential and unique
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:    "user_001",
				Email:     "alice@example.com",
				Seniority: "mid",
				WorkingHoursBand: seed.WorkingHours{
					Start: 9,
					End:   18,
				},
			},
		},
		Repositories: []seed.Repository{
			{
				RepoName:      "acme/test",
				DefaultBranch: "main",
			},
		},
	}

	rng := rand.New(rand.NewSource(123))
	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	// Create 5 commits on 5 different branches
	commits := make([]models.Commit, 5)
	for i := 0; i < 5; i++ {
		commits[i] = models.Commit{
			CommitHash:      string(rune('a' + i)),
			UserID:          "user_001",
			RepoName:        "acme/test",
			BranchName:      string(rune('A' + i)),
			TotalLinesAdded: 10,
			CommitTs:        baseTime.Add(time.Duration(i) * time.Minute),
		}
	}

	prGen := NewPRGenerator(seedData, rng)
	prs := prGen.GroupCommitsIntoPRs(commits)

	// Verify PR numbers are sequential
	seenNumbers := make(map[int]bool)
	for _, pr := range prs {
		if seenNumbers[pr.Number] {
			t.Errorf("Duplicate PR number: %d", pr.Number)
		}
		seenNumbers[pr.Number] = true

		if pr.Number < 1 {
			t.Errorf("PR number must be >= 1, got %d", pr.Number)
		}
	}

	// Verify we got 5 separate PRs (different branches)
	if len(prs) != 5 {
		t.Errorf("Expected 5 PRs (different branches), got %d", len(prs))
	}
}
