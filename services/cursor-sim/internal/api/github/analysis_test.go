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
