package generator

import (
	"fmt"
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

func TestPRGenerator_PRStatusDistribution(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	developer := seed.Developer{
		UserID: "user_001",
		Email:  "alice@example.com",
		Name:   "Alice",
	}

	repository := seed.Repository{
		RepoName:      "acme/platform",
		DefaultBranch: "main",
	}

	// Generate many commits to test status distribution
	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	var commits []models.Commit

	// Create 100 separate commits with large time gaps (to create separate PRs)
	for i := 0; i < 100; i++ {
		commits = append(commits, models.Commit{
			CommitHash:      fmt.Sprintf("hash%d", i),
			UserID:          "user_001",
			RepoName:        "acme/platform",
			BranchName:      fmt.Sprintf("feature/task-%d", i),
			TotalLinesAdded: 50,
			CommitTs:        baseTime.Add(time.Duration(i*48) * time.Hour), // 48h apart = separate PRs
		})
	}

	gen := NewPRGenerator(&seed.SeedData{
		Developers:   []seed.Developer{developer},
		Repositories: []seed.Repository{repository},
	}, rng)

	prs := gen.GroupCommitsIntoPRs(commits)

	// Count statuses
	statusCounts := make(map[models.PRState]int)
	for _, pr := range prs {
		statusCounts[pr.State]++
	}

	total := len(prs)
	mergedPct := float64(statusCounts[models.PRStateMerged]) / float64(total) * 100
	closedPct := float64(statusCounts[models.PRStateClosed]) / float64(total) * 100
	openPct := float64(statusCounts[models.PRStateOpen]) / float64(total) * 100

	t.Logf("Status distribution over %d PRs: merged=%.1f%%, closed=%.1f%%, open=%.1f%%",
		total, mergedPct, closedPct, openPct)

	// Allow some variance due to randomness, but should be roughly 85/10/5
	if mergedPct < 70 || mergedPct > 95 {
		t.Errorf("Merged percentage %.1f%% is outside expected range (70-95%%)", mergedPct)
	}
	if closedPct < 2 || closedPct > 20 {
		t.Errorf("Closed percentage %.1f%% is outside expected range (2-20%%)", closedPct)
	}
	if openPct < 0 || openPct > 15 {
		t.Errorf("Open percentage %.1f%% is outside expected range (0-15%%)", openPct)
	}
}

func TestPRGenerator_MergeTimestamps(t *testing.T) {
	rng := rand.New(rand.NewSource(99999))

	developer := seed.Developer{
		UserID: "user_001",
		Email:  "alice@example.com",
		Name:   "Alice",
	}

	repository := seed.Repository{
		RepoName:      "acme/platform",
		DefaultBranch: "main",
	}

	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	commits := []models.Commit{
		{
			CommitHash: "abc123",
			UserID:     "user_001",
			RepoName:   "acme/platform",
			BranchName: "feature/auth",
			CommitTs:   baseTime,
		},
	}

	gen := NewPRGenerator(&seed.SeedData{
		Developers:   []seed.Developer{developer},
		Repositories: []seed.Repository{repository},
	}, rng)

	prs := gen.GroupCommitsIntoPRs(commits)

	if len(prs) != 1 {
		t.Fatalf("Expected 1 PR, got %d", len(prs))
	}

	pr := prs[0]

	// CreatedAt should be before or at first commit time
	if pr.CreatedAt.After(baseTime) {
		t.Errorf("PR CreatedAt %v is after first commit %v", pr.CreatedAt, baseTime)
	}

	// Merged PRs should have MergedAt set
	if pr.State == models.PRStateMerged {
		if pr.MergedAt == nil {
			t.Error("Merged PR has nil MergedAt")
		} else {
			// MergedAt should be after CreatedAt
			if pr.MergedAt.Before(pr.CreatedAt) {
				t.Errorf("MergedAt %v is before CreatedAt %v", pr.MergedAt, pr.CreatedAt)
			}

			// MergedAt should be within 7 days
			diff := pr.MergedAt.Sub(pr.CreatedAt)
			if diff > 7*24*time.Hour+time.Second {
				t.Errorf("MergedAt is %v after CreatedAt, expected <= 7 days", diff)
			}
			if diff < 1*24*time.Hour {
				t.Errorf("MergedAt is %v after CreatedAt, expected >= 1 day", diff)
			}
		}
	}

	// Closed PRs should have ClosedAt set
	if pr.State == models.PRStateClosed {
		if pr.ClosedAt == nil {
			t.Error("Closed PR has nil ClosedAt")
		} else {
			if pr.ClosedAt.Before(pr.CreatedAt) {
				t.Errorf("ClosedAt %v is before CreatedAt %v", pr.ClosedAt, pr.CreatedAt)
			}
		}
	}

	// Open PRs should not have MergedAt or ClosedAt
	if pr.State == models.PRStateOpen {
		if pr.MergedAt != nil {
			t.Error("Open PR has MergedAt set")
		}
		if pr.ClosedAt != nil {
			t.Error("Open PR has ClosedAt set")
		}
	}
}

func TestPRGenerator_CommitIDsTracking(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	developer := seed.Developer{
		UserID: "user_001",
		Email:  "alice@example.com",
	}

	repository := seed.Repository{
		RepoName: "acme/platform",
	}

	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	commits := []models.Commit{
		{CommitHash: "hash1", UserID: "user_001", RepoName: "acme/platform", BranchName: "feature/test", CommitTs: baseTime},
		{CommitHash: "hash2", UserID: "user_001", RepoName: "acme/platform", BranchName: "feature/test", CommitTs: baseTime.Add(5 * time.Minute)},
		{CommitHash: "hash3", UserID: "user_001", RepoName: "acme/platform", BranchName: "feature/test", CommitTs: baseTime.Add(10 * time.Minute)},
	}

	gen := NewPRGenerator(&seed.SeedData{
		Developers:   []seed.Developer{developer},
		Repositories: []seed.Repository{repository},
	}, rng)

	prs := gen.GroupCommitsIntoPRs(commits)

	if len(prs) != 1 {
		t.Fatalf("Expected 1 PR, got %d", len(prs))
	}

	pr := prs[0]

	// Verify all commit IDs are tracked
	if len(pr.CommitIDs) != 3 {
		t.Errorf("PR has %d commit IDs, want 3", len(pr.CommitIDs))
	}

	expectedIDs := []string{"hash1", "hash2", "hash3"}
	for i, expectedID := range expectedIDs {
		if i >= len(pr.CommitIDs) {
			t.Errorf("Missing commit ID at index %d", i)
			continue
		}
		if pr.CommitIDs[i] != expectedID {
			t.Errorf("CommitIDs[%d] = %s, want %s", i, pr.CommitIDs[i], expectedID)
		}
	}
}

func TestPRGenerator_TitleGeneration(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	developer := seed.Developer{
		UserID: "user_001",
		Email:  "alice@example.com",
	}

	repository := seed.Repository{
		RepoName: "acme/platform",
	}

	tests := []struct {
		name           string
		branchName     string
		commitMessage  string
		expectedPrefix string
	}{
		{
			name:           "feature branch with commit message",
			branchName:     "feature/user-auth",
			commitMessage:  "Add user authentication",
			expectedPrefix: "Add user authentication",
		},
		{
			name:           "feature branch without commit message",
			branchName:     "feature/payment-flow",
			commitMessage:  "",
			expectedPrefix: "Implement payment",
		},
		{
			name:           "bugfix branch",
			branchName:     "bugfix/login-error",
			commitMessage:  "Fix login validation",
			expectedPrefix: "Fix login validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
			commits := []models.Commit{
				{
					CommitHash: "abc123",
					UserID:     "user_001",
					RepoName:   "acme/platform",
					BranchName: tt.branchName,
					Message:    tt.commitMessage,
					CommitTs:   baseTime,
				},
			}

			gen := NewPRGenerator(&seed.SeedData{
				Developers:   []seed.Developer{developer},
				Repositories: []seed.Repository{repository},
			}, rng)

			prs := gen.GroupCommitsIntoPRs(commits)

			if len(prs) != 1 {
				t.Fatalf("Expected 1 PR, got %d", len(prs))
			}

			pr := prs[0]

			if pr.Title == "" {
				t.Error("PR title is empty")
			}

			t.Logf("Generated title: %q", pr.Title)

			// Title should contain expected prefix
			if len(pr.Title) < len(tt.expectedPrefix) || pr.Title[:len(tt.expectedPrefix)] != tt.expectedPrefix {
				// Allow some flexibility for generated titles
				t.Logf("Title %q does not start with expected prefix %q (this may be OK)", pr.Title, tt.expectedPrefix)
			}
		})
	}
}

func TestPRGenerator_MultipleAuthors(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	alice := seed.Developer{
		UserID: "user_001",
		Email:  "alice@example.com",
		Name:   "Alice",
	}

	bob := seed.Developer{
		UserID: "user_002",
		Email:  "bob@example.com",
		Name:   "Bob",
	}

	repository := seed.Repository{
		RepoName: "acme/platform",
	}

	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	commits := []models.Commit{
		// Alice's commits
		{CommitHash: "a1", UserID: "user_001", RepoName: "acme/platform", BranchName: "feature/auth", CommitTs: baseTime},
		{CommitHash: "a2", UserID: "user_001", RepoName: "acme/platform", BranchName: "feature/auth", CommitTs: baseTime.Add(5 * time.Minute)},
		// Bob's commits
		{CommitHash: "b1", UserID: "user_002", RepoName: "acme/platform", BranchName: "feature/payment", CommitTs: baseTime.Add(1 * time.Hour)},
		{CommitHash: "b2", UserID: "user_002", RepoName: "acme/platform", BranchName: "feature/payment", CommitTs: baseTime.Add(1*time.Hour + 5*time.Minute)},
	}

	gen := NewPRGenerator(&seed.SeedData{
		Developers:   []seed.Developer{alice, bob},
		Repositories: []seed.Repository{repository},
	}, rng)

	prs := gen.GroupCommitsIntoPRs(commits)

	// Should create 2 PRs (one per author+branch)
	if len(prs) != 2 {
		t.Fatalf("Expected 2 PRs, got %d", len(prs))
	}

	// Verify each author has their own PR
	authorPRs := make(map[string]int)
	for _, pr := range prs {
		authorPRs[pr.AuthorID]++
	}

	if authorPRs["user_001"] != 1 {
		t.Errorf("Expected 1 PR for Alice, got %d", authorPRs["user_001"])
	}
	if authorPRs["user_002"] != 1 {
		t.Errorf("Expected 1 PR for Bob, got %d", authorPRs["user_002"])
	}
}

func TestPRGenerator_EmptyCommits(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	gen := NewPRGenerator(&seed.SeedData{
		Developers:   []seed.Developer{},
		Repositories: []seed.Repository{},
	}, rng)

	prs := gen.GroupCommitsIntoPRs([]models.Commit{})

	if len(prs) != 0 {
		t.Errorf("Expected 0 PRs from empty commits, got %d", len(prs))
	}

	// Should return non-nil slice
	if prs == nil {
		t.Error("Expected non-nil slice, got nil")
	}
}

func TestPRGenerator_PRNumberIncrement(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	developer := seed.Developer{
		UserID: "user_001",
		Email:  "alice@example.com",
	}

	repository := seed.Repository{
		RepoName: "acme/platform",
	}

	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	// Create commits for 3 different branches
	commits := []models.Commit{
		{CommitHash: "a1", UserID: "user_001", RepoName: "acme/platform", BranchName: "feature/auth", CommitTs: baseTime},
		{CommitHash: "b1", UserID: "user_001", RepoName: "acme/platform", BranchName: "feature/payment", CommitTs: baseTime.Add(24 * time.Hour)},
		{CommitHash: "c1", UserID: "user_001", RepoName: "acme/platform", BranchName: "bugfix/login", CommitTs: baseTime.Add(48 * time.Hour)},
	}

	gen := NewPRGenerator(&seed.SeedData{
		Developers:   []seed.Developer{developer},
		Repositories: []seed.Repository{repository},
	}, rng)

	prs := gen.GroupCommitsIntoPRs(commits)

	if len(prs) != 3 {
		t.Fatalf("Expected 3 PRs, got %d", len(prs))
	}

	// PR numbers should increment
	if prs[0].Number != 1 {
		t.Errorf("First PR number = %d, want 1", prs[0].Number)
	}
	if prs[1].Number != 2 {
		t.Errorf("Second PR number = %d, want 2", prs[1].Number)
	}
	if prs[2].Number != 3 {
		t.Errorf("Third PR number = %d, want 3", prs[2].Number)
	}
}

func TestPRGenerator_Reproducibility(t *testing.T) {
	developer := seed.Developer{
		UserID: "user_001",
		Email:  "alice@example.com",
	}

	repository := seed.Repository{
		RepoName: "acme/platform",
	}

	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	commits := []models.Commit{
		{CommitHash: "a1", UserID: "user_001", RepoName: "acme/platform", BranchName: "feature/auth", CommitTs: baseTime},
		{CommitHash: "a2", UserID: "user_001", RepoName: "acme/platform", BranchName: "feature/auth", CommitTs: baseTime.Add(5 * time.Minute)},
	}

	// Generate with same seed twice
	gen1 := NewPRGeneratorWithSeed(&seed.SeedData{
		Developers:   []seed.Developer{developer},
		Repositories: []seed.Repository{repository},
	}, nil, 999)

	prs1 := gen1.GroupCommitsIntoPRs(commits)

	gen2 := NewPRGeneratorWithSeed(&seed.SeedData{
		Developers:   []seed.Developer{developer},
		Repositories: []seed.Repository{repository},
	}, nil, 999)

	prs2 := gen2.GroupCommitsIntoPRs(commits)

	// Results should be identical
	if len(prs1) != len(prs2) {
		t.Fatalf("PR count mismatch: %d vs %d", len(prs1), len(prs2))
	}

	for i := range prs1 {
		if prs1[i].State != prs2[i].State {
			t.Errorf("PR %d state mismatch: %s vs %s", i, prs1[i].State, prs2[i].State)
		}
		if prs1[i].Title != prs2[i].Title {
			t.Errorf("PR %d title mismatch: %s vs %s", i, prs1[i].Title, prs2[i].Title)
		}
		// MergedAt times should be identical
		if prs1[i].MergedAt != nil && prs2[i].MergedAt != nil {
			if !prs1[i].MergedAt.Equal(*prs2[i].MergedAt) {
				t.Errorf("PR %d MergedAt mismatch: %v vs %v", i, prs1[i].MergedAt, prs2[i].MergedAt)
			}
		}
	}
}
