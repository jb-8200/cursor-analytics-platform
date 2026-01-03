package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

func TestSurvivalAnalysisHandler_Success(t *testing.T) {
	// Setup
	store := setupTestStore()
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	// Add test commits
	commits := []models.Commit{
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
	}

	for _, commit := range commits {
		if err := store.AddCommit(commit); err != nil {
			t.Fatalf("Failed to add commit: %v", err)
		}
	}

	// Test
	req := httptest.NewRequest("GET", "/repos/acme/platform/analysis/survival?cohort_start=2026-01-01&cohort_end=2026-01-31&observation_date=2026-02-28", nil)
	w := httptest.NewRecorder()

	handler := SurvivalAnalysisHandler(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var response models.SurvivalAnalysis
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify response structure
	if response.CohortStart != "2026-01-01" {
		t.Errorf("Expected cohort_start '2026-01-01', got '%s'", response.CohortStart)
	}

	if response.CohortEnd != "2026-01-31" {
		t.Errorf("Expected cohort_end '2026-01-31', got '%s'", response.CohortEnd)
	}

	if response.TotalLinesAdded == 0 {
		t.Error("Expected non-zero total lines added")
	}

	if response.SurvivalRate < 0 || response.SurvivalRate > 1 {
		t.Errorf("Expected survival rate between 0 and 1, got %.2f", response.SurvivalRate)
	}

	// Verify developer breakdown
	if len(response.ByDeveloper) == 0 {
		t.Error("Expected developer breakdown, got empty list")
	}

	for _, dev := range response.ByDeveloper {
		if dev.Email == "" {
			t.Error("Expected developer email")
		}

		if dev.SurvivalRate < 0 || dev.SurvivalRate > 1 {
			t.Errorf("Expected survival rate between 0 and 1 for %s, got %.2f", dev.Email, dev.SurvivalRate)
		}
	}
}

func TestSurvivalAnalysisHandler_DefaultDates(t *testing.T) {
	// Setup
	store := setupTestStore()

	// Test with no query parameters (should use defaults)
	req := httptest.NewRequest("GET", "/repos/acme/platform/analysis/survival", nil)
	w := httptest.NewRecorder()

	handler := SurvivalAnalysisHandler(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.SurvivalAnalysis
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify dates are set (default behavior)
	if response.CohortStart == "" {
		t.Error("Expected cohort_start to be set")
	}

	if response.CohortEnd == "" {
		t.Error("Expected cohort_end to be set")
	}

	if response.ObservationDate == "" {
		t.Error("Expected observation_date to be set")
	}
}

func TestSurvivalAnalysisHandler_InvalidRepo(t *testing.T) {
	// Setup
	store := setupTestStore()

	// Test with invalid repo path (missing owner/repo segments)
	req := httptest.NewRequest("GET", "/analysis/survival", nil)
	w := httptest.NewRecorder()

	handler := SurvivalAnalysisHandler(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestSurvivalAnalysisHandler_EmptyRepo(t *testing.T) {
	// Setup
	store := setupTestStore()

	// Test with repo that has no commits
	req := httptest.NewRequest("GET", "/repos/empty/repo/analysis/survival?cohort_start=2026-01-01&cohort_end=2026-01-31", nil)
	w := httptest.NewRecorder()

	handler := SurvivalAnalysisHandler(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.SurvivalAnalysis
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify empty response
	if response.TotalLinesAdded != 0 {
		t.Errorf("Expected 0 lines added for empty repo, got %d", response.TotalLinesAdded)
	}

	if response.SurvivalRate != 0.0 {
		t.Errorf("Expected 0.0 survival rate for empty repo, got %.2f", response.SurvivalRate)
	}
}

func TestRevertAnalysisHandler_Success(t *testing.T) {
	// Setup
	store := setupTestStore()
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	// Add test PRs
	mergedAt1 := baseTime
	mergedAt2 := baseTime.Add(24 * time.Hour)

	pr1 := models.PullRequest{
		Number:      1,
		State:       models.PRStateMerged,
		RepoName:    "acme/platform",
		AuthorID:    "user_001",
		AuthorEmail: "alice@example.com",
		AIRatio:     0.8,
		MergedAt:    &mergedAt1,
	}

	pr2 := models.PullRequest{
		Number:      2,
		State:       models.PRStateMerged,
		RepoName:    "acme/platform",
		AuthorID:    "user_002",
		AuthorEmail: "bob@example.com",
		AIRatio:     0.2,
		MergedAt:    &mergedAt2,
	}

	if err := store.AddPR(pr1); err != nil {
		t.Fatalf("Failed to add PR: %v", err)
	}
	if err := store.AddPR(pr2); err != nil {
		t.Fatalf("Failed to add PR: %v", err)
	}

	// Test
	req := httptest.NewRequest("GET", "/repos/acme/platform/analysis/reverts?window_days=7&since=2026-01-01&until=2026-01-31", nil)
	w := httptest.NewRecorder()

	handler := RevertAnalysisHandler(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var response models.RevertAnalysis
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify response structure
	if response.WindowDays != 7 {
		t.Errorf("Expected window_days 7, got %d", response.WindowDays)
	}

	if response.TotalPRsMerged != 2 {
		t.Errorf("Expected 2 PRs merged, got %d", response.TotalPRsMerged)
	}

	if response.RevertRate < 0 || response.RevertRate > 1 {
		t.Errorf("Expected revert rate between 0 and 1, got %.2f", response.RevertRate)
	}

	// Verify reverted PRs structure
	for _, reverted := range response.RevertedPRs {
		if reverted.PRNumber <= 0 {
			t.Error("Expected valid PR number")
		}

		if reverted.DaysToRevert < 0 || reverted.DaysToRevert > 7 {
			t.Errorf("Expected days_to_revert between 0 and 7, got %.2f", reverted.DaysToRevert)
		}

		if reverted.MergedAt == "" || reverted.RevertedAt == "" {
			t.Error("Expected non-empty timestamps")
		}
	}
}

func TestRevertAnalysisHandler_DefaultParameters(t *testing.T) {
	// Setup
	store := setupTestStore()

	// Test with no query parameters (should use defaults)
	req := httptest.NewRequest("GET", "/repos/acme/platform/analysis/reverts", nil)
	w := httptest.NewRecorder()

	handler := RevertAnalysisHandler(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.RevertAnalysis
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify defaults
	if response.WindowDays != 7 {
		t.Errorf("Expected default window_days=7, got %d", response.WindowDays)
	}
}

func TestRevertAnalysisHandler_InvalidRepo(t *testing.T) {
	// Setup
	store := setupTestStore()

	// Test with invalid repo path
	req := httptest.NewRequest("GET", "/analysis/reverts", nil)
	w := httptest.NewRecorder()

	handler := RevertAnalysisHandler(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRevertAnalysisHandler_EmptyRepo(t *testing.T) {
	// Setup
	store := setupTestStore()

	// Test with repo that has no PRs
	req := httptest.NewRequest("GET", "/repos/empty/repo/analysis/reverts?window_days=7", nil)
	w := httptest.NewRecorder()

	handler := RevertAnalysisHandler(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.RevertAnalysis
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify empty response
	if response.TotalPRsMerged != 0 {
		t.Errorf("Expected 0 PRs merged for empty repo, got %d", response.TotalPRsMerged)
	}

	if response.TotalPRsReverted != 0 {
		t.Errorf("Expected 0 PRs reverted for empty repo, got %d", response.TotalPRsReverted)
	}

	if response.RevertRate != 0.0 {
		t.Errorf("Expected 0.0 revert rate for empty repo, got %.2f", response.RevertRate)
	}
}
