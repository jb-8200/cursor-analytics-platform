package services

import (
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

// MockCommitStore implements CommitStore for testing.
type MockCommitStore struct {
	commits []models.Commit
}

func (m *MockCommitStore) GetCommitsByRepo(repoName string, from, to time.Time) []models.Commit {
	var result []models.Commit
	for _, commit := range m.commits {
		if commit.RepoName == repoName && !commit.CommitTs.Before(from) && commit.CommitTs.Before(to) {
			result = append(result, commit)
		}
	}
	return result
}

func TestCalculateSurvival_EmptyRepo(t *testing.T) {
	store := &MockCommitStore{commits: []models.Commit{}}
	svc := NewSurvivalServiceWithSeed(store, 42)

	cohortStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cohortEnd := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	observationDate := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	analysis, err := svc.CalculateSurvival("acme/platform", cohortStart, cohortEnd, observationDate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if analysis.TotalLinesAdded != 0 {
		t.Errorf("Expected 0 lines added, got %d", analysis.TotalLinesAdded)
	}

	if analysis.SurvivalRate != 0.0 {
		t.Errorf("Expected 0.0 survival rate, got %.2f", analysis.SurvivalRate)
	}
}

func TestCalculateSurvival_WithCommits(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	store := &MockCommitStore{
		commits: []models.Commit{
			{
				CommitHash:         "abc123",
				UserID:             "user_001",
				UserEmail:          "alice@example.com",
				RepoName:           "acme/platform",
				TotalLinesAdded:    100,
				TabLinesAdded:      40,
				ComposerLinesAdded: 30,
				CommitTs:           baseTime,
			},
			{
				CommitHash:         "def456",
				UserID:             "user_002",
				UserEmail:          "bob@example.com",
				RepoName:           "acme/platform",
				TotalLinesAdded:    50,
				TabLinesAdded:      20,
				ComposerLinesAdded: 10,
				CommitTs:           baseTime.Add(1 * time.Hour),
			},
		},
	}

	svc := NewSurvivalServiceWithSeed(store, 42)

	cohortStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cohortEnd := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	observationDate := time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC) // 28 days after cohort

	analysis, err := svc.CalculateSurvival("acme/platform", cohortStart, cohortEnd, observationDate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify basic metrics
	if analysis.TotalLinesAdded == 0 {
		t.Errorf("Expected non-zero lines added, got %d", analysis.TotalLinesAdded)
	}

	if analysis.SurvivalRate < 0 || analysis.SurvivalRate > 1 {
		t.Errorf("Expected survival rate between 0 and 1, got %.2f", analysis.SurvivalRate)
	}

	// Verify developer breakdown exists
	if len(analysis.ByDeveloper) == 0 {
		t.Error("Expected developer breakdown, got empty list")
	}

	// Verify format of dates
	if analysis.CohortStart != "2026-01-01" {
		t.Errorf("Expected cohort_start '2026-01-01', got '%s'", analysis.CohortStart)
	}

	if analysis.CohortEnd != "2026-01-31" {
		t.Errorf("Expected cohort_end '2026-01-31', got '%s'", analysis.CohortEnd)
	}
}

func TestCalculateSurvival_DeveloperBreakdown(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	store := &MockCommitStore{
		commits: []models.Commit{
			{
				CommitHash:         "abc123",
				UserEmail:          "alice@example.com",
				RepoName:           "acme/platform",
				TotalLinesAdded:    100,
				TabLinesAdded:      50,
				ComposerLinesAdded: 30,
				CommitTs:           baseTime,
			},
			{
				CommitHash:         "def456",
				UserEmail:          "bob@example.com",
				RepoName:           "acme/platform",
				TotalLinesAdded:    200,
				TabLinesAdded:      100,
				ComposerLinesAdded: 50,
				CommitTs:           baseTime.Add(1 * time.Hour),
			},
		},
	}

	svc := NewSurvivalServiceWithSeed(store, 42)

	cohortStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cohortEnd := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	observationDate := time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)

	analysis, err := svc.CalculateSurvival("acme/platform", cohortStart, cohortEnd, observationDate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify we have developers in breakdown
	if len(analysis.ByDeveloper) == 0 {
		t.Fatal("Expected developers in breakdown, got none")
	}

	// Verify each developer has valid metrics
	for _, dev := range analysis.ByDeveloper {
		if dev.Email == "" {
			t.Error("Expected developer email, got empty")
		}

		if dev.LinesAdded < 0 {
			t.Errorf("Expected non-negative lines added for %s, got %d", dev.Email, dev.LinesAdded)
		}

		if dev.SurvivalRate < 0 || dev.SurvivalRate > 1 {
			t.Errorf("Expected survival rate between 0 and 1 for %s, got %.2f", dev.Email, dev.SurvivalRate)
		}
	}
}

func TestCalculateSurvival_Reproducibility(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	store := &MockCommitStore{
		commits: []models.Commit{
			{
				CommitHash:      "abc123",
				UserEmail:       "alice@example.com",
				RepoName:        "acme/platform",
				TotalLinesAdded: 100,
				CommitTs:        baseTime,
			},
		},
	}

	cohortStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cohortEnd := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	observationDate := time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)

	// Run twice with same seed
	svc1 := NewSurvivalServiceWithSeed(store, 42)
	analysis1, _ := svc1.CalculateSurvival("acme/platform", cohortStart, cohortEnd, observationDate)

	svc2 := NewSurvivalServiceWithSeed(store, 42)
	analysis2, _ := svc2.CalculateSurvival("acme/platform", cohortStart, cohortEnd, observationDate)

	// Results should be identical
	if analysis1.TotalLinesAdded != analysis2.TotalLinesAdded {
		t.Errorf("Expected reproducible total lines: %d vs %d", analysis1.TotalLinesAdded, analysis2.TotalLinesAdded)
	}

	if analysis1.LinesSurviving != analysis2.LinesSurviving {
		t.Errorf("Expected reproducible surviving lines: %d vs %d", analysis1.LinesSurviving, analysis2.LinesSurviving)
	}

	if analysis1.SurvivalRate != analysis2.SurvivalRate {
		t.Errorf("Expected reproducible survival rate: %.3f vs %.3f", analysis1.SurvivalRate, analysis2.SurvivalRate)
	}
}
