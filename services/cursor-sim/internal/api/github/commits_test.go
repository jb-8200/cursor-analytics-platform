package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

func TestListCommits(t *testing.T) {
	// Setup
	store := setupTestStore()
	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	// Add test commits
	commits := []models.Commit{
		{
			CommitHash:      "abc123",
			UserID:          "user_001",
			UserEmail:       "alice@example.com",
			UserName:        "Alice",
			RepoName:        "acme/platform",
			BranchName:      "main",
			TotalLinesAdded: 50,
			TabLinesAdded:   20,
			CommitTs:        baseTime,
		},
		{
			CommitHash:      "def456",
			UserID:          "user_002",
			UserEmail:       "bob@example.com",
			UserName:        "Bob",
			RepoName:        "acme/platform",
			BranchName:      "main",
			TotalLinesAdded: 30,
			TabLinesAdded:   10,
			CommitTs:        baseTime.Add(1 * time.Hour),
		},
	}

	for _, commit := range commits {
		if err := store.AddCommit(commit); err != nil {
			t.Fatalf("Failed to add commit: %v", err)
		}
	}

	// Test
	req := httptest.NewRequest("GET", "/repos/acme/platform/commits", nil)
	w := httptest.NewRecorder()

	handler := ListCommits(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []models.Commit
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 commits, got %d", len(response))
	}
}

func TestGetCommit(t *testing.T) {
	// Setup
	store := setupTestStore()
	baseTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	commit := models.Commit{
		CommitHash:      "abc123",
		UserID:          "user_001",
		UserEmail:       "alice@example.com",
		UserName:        "Alice",
		RepoName:        "acme/platform",
		BranchName:      "main",
		TotalLinesAdded: 50,
		TabLinesAdded:   20,
		CommitTs:        baseTime,
	}

	if err := store.AddCommit(commit); err != nil {
		t.Fatalf("Failed to add commit: %v", err)
	}

	// Test
	req := httptest.NewRequest("GET", "/repos/acme/platform/commits/abc123", nil)
	w := httptest.NewRecorder()

	handler := GetCommit(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.Commit
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.CommitHash != "abc123" {
		t.Errorf("Expected commit hash abc123, got %s", response.CommitHash)
	}
}

func TestGetCommitNotFound(t *testing.T) {
	// Setup
	store := setupTestStore()

	// Test
	req := httptest.NewRequest("GET", "/repos/acme/platform/commits/notfound", nil)
	w := httptest.NewRecorder()

	handler := GetCommit(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestListPullCommits(t *testing.T) {
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
		Deletions:    15,
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
	req := httptest.NewRequest("GET", "/repos/acme/platform/pulls/1/commits", nil)
	w := httptest.NewRecorder()

	handler := ListPullCommits(store)
	handler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var response []models.Commit
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 commits, got %d", len(response))
	}
}

