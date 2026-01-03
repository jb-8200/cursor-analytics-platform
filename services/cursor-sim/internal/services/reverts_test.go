package services

import (
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// MockPRStore implements PRStore for testing.
type MockPRStore struct {
	prs       []models.PullRequest
	developers map[string]*seed.Developer
}

func (m *MockPRStore) GetPRsByRepo(repoName string) []models.PullRequest {
	var result []models.PullRequest
	for _, pr := range m.prs {
		if pr.RepoName == repoName {
			result = append(result, pr)
		}
	}
	return result
}

func (m *MockPRStore) GetPRsByRepoAndState(repoName string, state models.PRState) []models.PullRequest {
	var result []models.PullRequest
	for _, pr := range m.prs {
		if pr.RepoName == repoName && pr.State == state {
			result = append(result, pr)
		}
	}
	return result
}

func (m *MockPRStore) GetDeveloper(userID string) (*seed.Developer, error) {
	dev, exists := m.developers[userID]
	if !exists {
		return nil, nil
	}
	return dev, nil
}

func TestIsRevertMessage(t *testing.T) {
	tests := []struct {
		message  string
		expected bool
	}{
		{"Revert \"Add feature X\"", true},
		{"revert changes from PR #123", true},
		{"Rollback deployment", true},
		{"ROLLBACK broken changes", true},
		{"Reverting bad commit", true},
		{"Backing out changes", true},
		{"Add new feature", false},
		{"Fix bug in authentication", false},
		{"Normal commit message", false},
	}

	for _, tt := range tests {
		result := IsRevertMessage(tt.message)
		if result != tt.expected {
			t.Errorf("IsRevertMessage(%q) = %v, expected %v", tt.message, result, tt.expected)
		}
	}
}

func TestExtractPRNumber(t *testing.T) {
	tests := []struct {
		message  string
		expected int
		found    bool
	}{
		{"Revert PR #123", 123, true},
		{"Fixes #456", 456, true},
		{"Related to PR 789", 789, true},
		{"Closes pull request #234", 234, true},
		{"No PR number here", 0, false},
		{"PR without number", 0, false},
	}

	for _, tt := range tests {
		prNum, found := ExtractPRNumber(tt.message)
		if found != tt.found {
			t.Errorf("ExtractPRNumber(%q) found = %v, expected %v", tt.message, found, tt.found)
		}
		if found && prNum != tt.expected {
			t.Errorf("ExtractPRNumber(%q) = %d, expected %d", tt.message, prNum, tt.expected)
		}
	}
}

func TestCalculateRevertRisk(t *testing.T) {
	tests := []struct {
		name          string
		pr            models.PullRequest
		developer     *seed.Developer
		expectedRange []float64 // [min, max] range
	}{
		{
			name: "High AI ratio with junior developer",
			pr: models.PullRequest{
				AIRatio: 0.9,
			},
			developer: &seed.Developer{
				Seniority:     "junior",
				ActivityLevel: "high",
			},
			expectedRange: []float64{0.08, 0.15}, // High risk
		},
		{
			name: "Low AI ratio with senior developer",
			pr: models.PullRequest{
				AIRatio: 0.1,
			},
			developer: &seed.Developer{
				Seniority:     "senior",
				ActivityLevel: "low",
			},
			expectedRange: []float64{0.0, 0.08}, // Low risk (adjusted)
		},
		{
			name: "Medium AI ratio with mid developer",
			pr: models.PullRequest{
				AIRatio: 0.5,
			},
			developer: &seed.Developer{
				Seniority:     "mid",
				ActivityLevel: "medium",
			},
			expectedRange: []float64{0.05, 0.13}, // Medium risk (adjusted)
		},
		{
			name: "No developer info",
			pr: models.PullRequest{
				AIRatio: 0.7,
			},
			developer:     nil,
			expectedRange: []float64{0.0, 0.15},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			risk := CalculateRevertRisk(tt.pr, tt.developer)

			if risk < tt.expectedRange[0] || risk > tt.expectedRange[1] {
				t.Errorf("CalculateRevertRisk() = %.3f, expected range [%.3f, %.3f]",
					risk, tt.expectedRange[0], tt.expectedRange[1])
			}

			// Verify risk is capped at 0.15
			if risk > 0.15 {
				t.Errorf("CalculateRevertRisk() = %.3f, expected max 0.15", risk)
			}
		})
	}
}

func TestCalculateRevertRisk_AIRatioCorrelation(t *testing.T) {
	// Verify that higher AI ratio leads to higher risk
	dev := &seed.Developer{
		Seniority:     "mid",
		ActivityLevel: "medium",
	}

	pr1 := models.PullRequest{AIRatio: 0.1}
	pr2 := models.PullRequest{AIRatio: 0.5}
	pr3 := models.PullRequest{AIRatio: 0.9}

	risk1 := CalculateRevertRisk(pr1, dev)
	risk2 := CalculateRevertRisk(pr2, dev)
	risk3 := CalculateRevertRisk(pr3, dev)

	if !(risk1 < risk2 && risk2 < risk3) {
		t.Errorf("Expected increasing risk with AI ratio: %.3f < %.3f < %.3f",
			risk1, risk2, risk3)
	}
}

func TestCalculateRevertRisk_SeniorityCorrelation(t *testing.T) {
	// Verify that senior developers have lower risk than juniors
	pr := models.PullRequest{AIRatio: 0.7}

	junior := &seed.Developer{Seniority: "junior", ActivityLevel: "medium"}
	mid := &seed.Developer{Seniority: "mid", ActivityLevel: "medium"}
	senior := &seed.Developer{Seniority: "senior", ActivityLevel: "medium"}

	riskJunior := CalculateRevertRisk(pr, junior)
	riskMid := CalculateRevertRisk(pr, mid)
	riskSenior := CalculateRevertRisk(pr, senior)

	// All risks should be capped at 0.15
	if riskJunior > 0.15 || riskMid > 0.15 || riskSenior > 0.15 {
		t.Errorf("Risk should be capped at 0.15: junior=%.3f, mid=%.3f, senior=%.3f",
			riskJunior, riskMid, riskSenior)
	}

	// Junior should have highest risk, senior should have lowest
	if !(riskSenior < riskJunior) {
		t.Errorf("Expected senior risk < junior risk: senior(%.3f) < junior(%.3f)",
			riskSenior, riskJunior)
	}
}

func TestGetReverts_EmptyRepo(t *testing.T) {
	store := &MockPRStore{
		prs:        []models.PullRequest{},
		developers: make(map[string]*seed.Developer),
	}

	svc := NewRevertServiceWithSeed(store, 42)

	since := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	analysis, err := svc.GetReverts("acme/platform", 7, since, until)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if analysis.TotalPRsMerged != 0 {
		t.Errorf("Expected 0 PRs merged, got %d", analysis.TotalPRsMerged)
	}

	if analysis.TotalPRsReverted != 0 {
		t.Errorf("Expected 0 PRs reverted, got %d", analysis.TotalPRsReverted)
	}

	if analysis.RevertRate != 0.0 {
		t.Errorf("Expected 0.0 revert rate, got %.2f", analysis.RevertRate)
	}
}

func TestGetReverts_WithMergedPRs(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	mergedAt1 := baseTime
	mergedAt2 := baseTime.Add(24 * time.Hour)

	store := &MockPRStore{
		prs: []models.PullRequest{
			{
				Number:      1,
				State:       models.PRStateMerged,
				RepoName:    "acme/platform",
				AuthorID:    "user_001",
				AIRatio:     0.8, // High AI ratio
				MergedAt:    &mergedAt1,
			},
			{
				Number:      2,
				State:       models.PRStateMerged,
				RepoName:    "acme/platform",
				AuthorID:    "user_002",
				AIRatio:     0.2, // Low AI ratio
				MergedAt:    &mergedAt2,
			},
		},
		developers: map[string]*seed.Developer{
			"user_001": {UserID: "user_001", Seniority: "junior", ActivityLevel: "high"},
			"user_002": {UserID: "user_002", Seniority: "senior", ActivityLevel: "low"},
		},
	}

	svc := NewRevertServiceWithSeed(store, 42)

	since := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	analysis, err := svc.GetReverts("acme/platform", 7, since, until)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if analysis.TotalPRsMerged != 2 {
		t.Errorf("Expected 2 PRs merged, got %d", analysis.TotalPRsMerged)
	}

	if analysis.WindowDays != 7 {
		t.Errorf("Expected window_days=7, got %d", analysis.WindowDays)
	}

	// Revert rate should be between 0 and 1
	if analysis.RevertRate < 0 || analysis.RevertRate > 1 {
		t.Errorf("Expected revert rate between 0 and 1, got %.2f", analysis.RevertRate)
	}

	// Number of reverted PRs should match revert rate calculation
	expectedReverted := int(float64(analysis.TotalPRsMerged) * analysis.RevertRate)
	if analysis.TotalPRsReverted != expectedReverted {
		t.Errorf("Inconsistent revert count: got %d, expected %d", analysis.TotalPRsReverted, expectedReverted)
	}
}

func TestGetReverts_Reproducibility(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	mergedAt := baseTime

	store := &MockPRStore{
		prs: []models.PullRequest{
			{
				Number:   1,
				State:    models.PRStateMerged,
				RepoName: "acme/platform",
				AuthorID: "user_001",
				AIRatio:  0.7,
				MergedAt: &mergedAt,
			},
		},
		developers: map[string]*seed.Developer{
			"user_001": {UserID: "user_001", Seniority: "mid", ActivityLevel: "medium"},
		},
	}

	since := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	// Run twice with same seed
	svc1 := NewRevertServiceWithSeed(store, 42)
	analysis1, _ := svc1.GetReverts("acme/platform", 7, since, until)

	svc2 := NewRevertServiceWithSeed(store, 42)
	analysis2, _ := svc2.GetReverts("acme/platform", 7, since, until)

	// Results should be identical
	if analysis1.TotalPRsReverted != analysis2.TotalPRsReverted {
		t.Errorf("Expected reproducible reverted count: %d vs %d",
			analysis1.TotalPRsReverted, analysis2.TotalPRsReverted)
	}

	if analysis1.RevertRate != analysis2.RevertRate {
		t.Errorf("Expected reproducible revert rate: %.3f vs %.3f",
			analysis1.RevertRate, analysis2.RevertRate)
	}
}

func TestGenerateRevertMessage(t *testing.T) {
	message := GenerateRevertMessage(123, "Add authentication feature")

	// Should contain PR number
	if !contains(message, "123") {
		t.Errorf("Expected message to contain PR number 123, got: %s", message)
	}

	// Should indicate it's a revert
	if !IsRevertMessage(message) {
		t.Errorf("Expected message to be detected as revert, got: %s", message)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
