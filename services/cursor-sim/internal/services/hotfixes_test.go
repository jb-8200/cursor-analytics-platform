package services

import (
	"strings"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

// MockHotfixStore implements HotfixPRStore for testing.
type MockHotfixStore struct {
	prs []models.PullRequest
}

func (m *MockHotfixStore) GetPRsByRepo(repoName string) []models.PullRequest {
	var result []models.PullRequest
	for _, pr := range m.prs {
		if pr.RepoName == repoName {
			result = append(result, pr)
		}
	}
	return result
}

func (m *MockHotfixStore) GetPRsByRepoAndState(repoName string, state models.PRState) []models.PullRequest {
	var result []models.PullRequest
	for _, pr := range m.prs {
		if pr.RepoName == repoName && pr.State == state {
			result = append(result, pr)
		}
	}
	return result
}

func TestIsHotfixPR(t *testing.T) {
	tests := []struct {
		title    string
		body     string
		expected bool
	}{
		{"Fix critical bug", "", true},
		{"Hotfix authentication", "", true},
		{"Urgent patch for production", "", true},
		{"Patch database issue", "", true},
		{"Add new feature", "Contains fix keyword", true},
		{"Update documentation", "No hotfix here", true},
		{"Refactor utils", "Add new method", false},
	}

	for _, tt := range tests {
		pr := models.PullRequest{Title: tt.title, Body: tt.body}
		result := IsHotfixPR(pr)
		if result != tt.expected {
			t.Errorf("IsHotfixPR(title=%q, body=%q) = %v, expected %v",
				tt.title, tt.body, result, tt.expected)
		}
	}
}

func TestFilesOverlap(t *testing.T) {
	pr1 := models.PullRequest{
		RepoName:     "acme/platform",
		ChangedFiles: 5,
	}
	pr2 := models.PullRequest{
		RepoName:     "acme/platform",
		ChangedFiles: 4,
	}

	overlap := FilesOverlap(pr1, pr2)

	// Should have some overlap
	if len(overlap) == 0 {
		t.Error("Expected some file overlap, got empty")
	}

	// All files should have file extensions
	validExtensions := []string{".ts", ".tsx", ".go", ".py", ".js", ".jsx"}
	for _, f := range overlap {
		hasValidExtension := false
		for _, ext := range validExtensions {
			if strings.HasSuffix(f, ext) {
				hasValidExtension = true
				break
			}
		}
		if !hasValidExtension {
			t.Errorf("Expected file with valid extension, got %q", f)
		}
	}
}

func TestFilesOverlap_DifferentRepos(t *testing.T) {
	pr1 := models.PullRequest{RepoName: "repo1", ChangedFiles: 5}
	pr2 := models.PullRequest{RepoName: "repo2", ChangedFiles: 5}

	overlap := FilesOverlap(pr1, pr2)

	if len(overlap) != 0 {
		t.Error("Expected no overlap for different repos")
	}
}

func TestGetHotfixes_EmptyRepo(t *testing.T) {
	store := &MockHotfixStore{prs: []models.PullRequest{}}
	svc := NewHotfixService(store)

	since := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	analysis := svc.GetHotfixes("acme/platform", 48, since, until)

	if analysis.TotalPRsMerged != 0 {
		t.Errorf("Expected 0 PRs merged, got %d", analysis.TotalPRsMerged)
	}

	if analysis.PRsWithHotfix != 0 {
		t.Errorf("Expected 0 hotfixes, got %d", analysis.PRsWithHotfix)
	}

	if analysis.HotfixRate != 0.0 {
		t.Errorf("Expected 0.0 hotfix rate, got %.2f", analysis.HotfixRate)
	}
}

func TestGetHotfixes_WithMergedPRs(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	// First PR
	mergedAt1 := baseTime
	// Hotfix PR 5 hours after
	mergedAt2 := baseTime.Add(5 * time.Hour)

	store := &MockHotfixStore{
		prs: []models.PullRequest{
			{
				Number:       1,
				State:        models.PRStateMerged,
				RepoName:     "acme/platform",
				Title:        "Add auth feature",
				Body:         "Adds new authentication",
				MergedAt:     &mergedAt1,
				ChangedFiles: 5,
			},
			{
				Number:       2,
				State:        models.PRStateMerged,
				RepoName:     "acme/platform",
				Title:        "Fix auth bug",
				Body:         "Fixes security issue",
				MergedAt:     &mergedAt2,
				ChangedFiles: 3,
			},
		},
	}

	svc := NewHotfixService(store)

	since := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	analysis := svc.GetHotfixes("acme/platform", 48, since, until)

	if analysis.TotalPRsMerged != 2 {
		t.Errorf("Expected 2 PRs merged, got %d", analysis.TotalPRsMerged)
	}

	if analysis.WindowHours != 48 {
		t.Errorf("Expected window_hours=48, got %d", analysis.WindowHours)
	}

	// Hotfix rate should be <= 1.0
	if analysis.HotfixRate < 0 || analysis.HotfixRate > 1 {
		t.Errorf("Expected hotfix rate between 0 and 1, got %.2f", analysis.HotfixRate)
	}
}

func TestGetHotfixes_OutsideWindow(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	// First PR
	mergedAt1 := baseTime
	// Potential hotfix 60 hours later (outside 48h window)
	mergedAt2 := baseTime.Add(60 * time.Hour)

	store := &MockHotfixStore{
		prs: []models.PullRequest{
			{
				Number:       1,
				State:        models.PRStateMerged,
				RepoName:     "acme/platform",
				Title:        "Add auth feature",
				MergedAt:     &mergedAt1,
				ChangedFiles: 5,
			},
			{
				Number:       2,
				State:        models.PRStateMerged,
				RepoName:     "acme/platform",
				Title:        "Fix auth bug",
				MergedAt:     &mergedAt2,
				ChangedFiles: 3,
			},
		},
	}

	svc := NewHotfixService(store)

	since := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	analysis := svc.GetHotfixes("acme/platform", 48, since, until)

	// Should have no hotfixes (outside window)
	if analysis.PRsWithHotfix != 0 {
		t.Errorf("Expected 0 hotfixes (outside window), got %d", analysis.PRsWithHotfix)
	}
}
