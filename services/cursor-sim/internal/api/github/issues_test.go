package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupIssuesTestStore creates a test store with sample issues.
func setupIssuesTestStore() storage.Store {
	store := storage.NewMemoryStore()

	// Create sample issues across multiple repos
	issues := []models.Issue{
		{
			Number:    1,
			Title:     "Bug: Login fails",
			State:     models.IssueStateOpen,
			AuthorID:  "user1",
			RepoName:  "org/repo1",
			Labels:    []string{"bug", "critical"},
			CreatedAt: time.Date(2025, 1, 5, 10, 0, 0, 0, time.UTC),
		},
		{
			Number:    2,
			Title:     "Feature: Add dark mode",
			State:     models.IssueStateOpen,
			AuthorID:  "user2",
			RepoName:  "org/repo1",
			Labels:    []string{"feature", "enhancement"},
			CreatedAt: time.Date(2025, 1, 6, 11, 0, 0, 0, time.UTC),
		},
		{
			Number:    3,
			Title:     "Fix typo in README",
			State:     models.IssueStateClosed,
			AuthorID:  "user3",
			RepoName:  "org/repo1",
			Labels:    []string{"documentation"},
			CreatedAt: time.Date(2025, 1, 4, 9, 0, 0, 0, time.UTC),
		},
		{
			Number:    1,
			Title:     "Performance issue",
			State:     models.IssueStateOpen,
			AuthorID:  "user1",
			RepoName:  "org/repo2",
			Labels:    []string{"bug", "performance"},
			CreatedAt: time.Date(2025, 1, 7, 12, 0, 0, 0, time.UTC),
		},
		{
			Number:    2,
			Title:     "Refactor database layer",
			State:     models.IssueStateClosed,
			AuthorID:  "user2",
			RepoName:  "org/repo2",
			Labels:    []string{"refactor"},
			CreatedAt: time.Date(2025, 1, 3, 8, 0, 0, 0, time.UTC),
		},
	}

	for _, issue := range issues {
		if err := store.StoreIssue(issue); err != nil {
			panic(err)
		}
	}

	return store
}

func TestListIssuesAnalytics_NoFilters(t *testing.T) {
	store := setupIssuesTestStore()
	handler := ListIssuesAnalytics(store)

	req := httptest.NewRequest(http.MethodGet, "/analytics/github/issues", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response IssuesAnalyticsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Should return all issues with default pagination
	assert.Equal(t, 5, response.Pagination.Total)
	assert.Equal(t, 5, len(response.Data))
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 20, response.Pagination.PageSize)
}

func TestListIssuesAnalytics_FilterByState_Open(t *testing.T) {
	store := setupIssuesTestStore()
	handler := ListIssuesAnalytics(store)

	req := httptest.NewRequest(http.MethodGet, "/analytics/github/issues?state=open", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response IssuesAnalyticsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Should return only open issues
	assert.Equal(t, 3, response.Pagination.Total)
	assert.Equal(t, 3, len(response.Data))
	assert.Equal(t, "open", response.Params.State)

	// Verify all issues are open
	for _, issue := range response.Data {
		assert.Equal(t, models.IssueStateOpen, issue.State)
	}
}

func TestListIssuesAnalytics_FilterByState_Closed(t *testing.T) {
	store := setupIssuesTestStore()
	handler := ListIssuesAnalytics(store)

	req := httptest.NewRequest(http.MethodGet, "/analytics/github/issues?state=closed", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response IssuesAnalyticsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Should return only closed issues
	assert.Equal(t, 2, response.Pagination.Total)
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, "closed", response.Params.State)

	// Verify all issues are closed
	for _, issue := range response.Data {
		assert.Equal(t, models.IssueStateClosed, issue.State)
	}
}

func TestListIssuesAnalytics_FilterByLabels_Single(t *testing.T) {
	store := setupIssuesTestStore()
	handler := ListIssuesAnalytics(store)

	req := httptest.NewRequest(http.MethodGet, "/analytics/github/issues?labels=bug", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response IssuesAnalyticsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Should return only issues with "bug" label
	assert.Equal(t, 2, response.Pagination.Total)
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, "bug", response.Params.Labels)

	// Verify all issues have "bug" label
	for _, issue := range response.Data {
		assert.Contains(t, issue.Labels, "bug")
	}
}

func TestListIssuesAnalytics_FilterByLabels_Multiple(t *testing.T) {
	store := setupIssuesTestStore()
	handler := ListIssuesAnalytics(store)

	req := httptest.NewRequest(http.MethodGet, "/analytics/github/issues?labels=bug,critical", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response IssuesAnalyticsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Should return only issues with both "bug" AND "critical" labels
	assert.Equal(t, 1, response.Pagination.Total)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, "bug,critical", response.Params.Labels)

	// Verify issue has both labels
	issue := response.Data[0]
	assert.Contains(t, issue.Labels, "bug")
	assert.Contains(t, issue.Labels, "critical")
}

func TestListIssuesAnalytics_CombinedFilters(t *testing.T) {
	store := setupIssuesTestStore()
	handler := ListIssuesAnalytics(store)

	req := httptest.NewRequest(http.MethodGet, "/analytics/github/issues?state=open&labels=bug", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response IssuesAnalyticsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Should return only open issues with "bug" label
	assert.Equal(t, 2, response.Pagination.Total)
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, "open", response.Params.State)
	assert.Equal(t, "bug", response.Params.Labels)

	// Verify all issues match criteria
	for _, issue := range response.Data {
		assert.Equal(t, models.IssueStateOpen, issue.State)
		assert.Contains(t, issue.Labels, "bug")
	}
}

func TestListIssuesAnalytics_Pagination(t *testing.T) {
	store := setupIssuesTestStore()
	handler := ListIssuesAnalytics(store)

	// Request first page with 2 items per page
	req := httptest.NewRequest(http.MethodGet, "/analytics/github/issues?page=1&page_size=2", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response IssuesAnalyticsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, 5, response.Pagination.Total)
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 2, response.Pagination.PageSize)

	// Request second page
	req = httptest.NewRequest(http.MethodGet, "/analytics/github/issues?page=2&page_size=2", nil)
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, 5, response.Pagination.Total)
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, 2, response.Pagination.Page)
	assert.Equal(t, 2, response.Pagination.PageSize)
}

func TestListIssuesAnalytics_InvalidState(t *testing.T) {
	store := setupIssuesTestStore()
	handler := ListIssuesAnalytics(store)

	req := httptest.NewRequest(http.MethodGet, "/analytics/github/issues?state=invalid", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListIssuesAnalytics_InvalidPage(t *testing.T) {
	store := setupIssuesTestStore()
	handler := ListIssuesAnalytics(store)

	req := httptest.NewRequest(http.MethodGet, "/analytics/github/issues?page=0", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListIssuesAnalytics_InvalidPageSize(t *testing.T) {
	store := setupIssuesTestStore()
	handler := ListIssuesAnalytics(store)

	req := httptest.NewRequest(http.MethodGet, "/analytics/github/issues?page_size=-1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListIssuesAnalytics_EmptyResult(t *testing.T) {
	store := setupIssuesTestStore()
	handler := ListIssuesAnalytics(store)

	// Query for non-existent labels
	req := httptest.NewRequest(http.MethodGet, "/analytics/github/issues?labels=nonexistent", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response IssuesAnalyticsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Should return empty array, not nil
	assert.Equal(t, 0, response.Pagination.Total)
	assert.NotNil(t, response.Data)
	assert.Equal(t, 0, len(response.Data))
}

func TestListIssuesAnalytics_PageSizeLimit(t *testing.T) {
	store := setupIssuesTestStore()
	handler := ListIssuesAnalytics(store)

	// Request page_size > 100 (should cap at 100)
	req := httptest.NewRequest(http.MethodGet, "/analytics/github/issues?page_size=500", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response IssuesAnalyticsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Should be capped at 100
	assert.Equal(t, 100, response.Pagination.PageSize)
}
