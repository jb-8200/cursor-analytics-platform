package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

func TestListPullFiles(t *testing.T) {
	// Setup
	store := setupTestStore()
	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	// Create PR
	pr := models.PullRequest{
		Number:       1,
		RepoName:     "acme/platform",
		HeadBranch:   "feature/auth",
		CommitCount:  2,
		Additions:    80,
		ChangedFiles: 2,
	}
	if err := store.AddPR(pr); err != nil {
		t.Fatalf("Failed to add PR: %v", err)
	}

	// Add commits on the PR's branch
	commits := []models.Commit{
		{
			CommitHash:      "abc123",
			UserID:          "user_001",
			RepoName:        "acme/platform",
			BranchName:      "feature/auth",
			TotalLinesAdded: 50,
			CommitTs:        baseTime,
		},
		{
			CommitHash:      "def456",
			UserID:          "user_001",
			RepoName:        "acme/platform",
			BranchName:      "feature/auth",
			TotalLinesAdded: 30,
			CommitTs:        baseTime.Add(1 * time.Hour),
		},
	}

	for _, commit := range commits {
		if err := store.AddCommit(commit); err != nil {
			t.Fatalf("Failed to add commit: %v", err)
		}
	}

	// Test
	req := httptest.NewRequest("GET", "/repos/acme/platform/pulls/1/files", nil)
	w := httptest.NewRecorder()

	handler := ListPullFiles(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var response []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 files, got %d", len(response))
	}

	// Verify each file has required fields
	for _, file := range response {
		if _, ok := file["filename"]; !ok {
			t.Error("File missing 'filename' field")
		}
		if _, ok := file["additions"]; !ok {
			t.Error("File missing 'additions' field")
		}
		if _, ok := file["greenfield_index"]; !ok {
			t.Error("File missing 'greenfield_index' field")
		}
	}
}

func TestListPullFilesWithGreenfieldIndex(t *testing.T) {
	// Setup
	store := setupTestStore()
	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	// Create PR (first commit on this branch = greenfield)
	pr := models.PullRequest{
		Number:       1,
		RepoName:     "acme/platform",
		HeadBranch:   "feature/auth",
		CommitCount:  1,
		Additions:    50,
		ChangedFiles: 1,
		CreatedAt:    baseTime, // Greenfield since PR created at baseTime
	}
	if err := store.AddPR(pr); err != nil {
		t.Fatalf("Failed to add PR: %v", err)
	}

	// Add commit (matches PR creation time = 100% greenfield)
	commit := models.Commit{
		CommitHash:      "abc123",
		UserID:          "user_001",
		RepoName:        "acme/platform",
		BranchName:      "feature/auth",
		TotalLinesAdded: 50,
		CommitTs:        baseTime,
	}
	if err := store.AddCommit(commit); err != nil {
		t.Fatalf("Failed to add commit: %v", err)
	}

	// Test
	req := httptest.NewRequest("GET", "/repos/acme/platform/pulls/1/files", nil)
	w := httptest.NewRecorder()

	handler := ListPullFiles(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) > 0 {
		// Files created < 30 days ago should be greenfield
		for _, file := range response {
			if greenfield, ok := file["is_greenfield"].(bool); ok && !greenfield {
				t.Error("Expected is_greenfield to be true for new files")
			}
		}
	}
}
