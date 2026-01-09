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

// setupPRAnalyticsTestStore creates a test store with sample PR data.
func setupPRAnalyticsTestStore() *storage.MemoryStore {
	store := storage.NewMemoryStore()
	now := time.Now()

	// Add test PRs with various states and authors
	prs := []models.PullRequest{
		{
			ID:          1,
			Number:      1,
			Title:       "feat: add authentication",
			State:       models.PRStateOpen,
			AuthorEmail: "alice@example.com",
			AuthorName:  "Alice",
			RepoName:    "acme/api",
			BaseBranch:  "main",
			HeadBranch:  "feature/auth",
			Additions:   150,
			Deletions:   30,
			AIRatio:     0.7,
			TabLines:    105,
			CreatedAt:   now.Add(-72 * time.Hour),
		},
		{
			ID:          2,
			Number:      2,
			Title:       "fix: resolve login bug",
			State:       models.PRStateMerged,
			AuthorEmail: "bob@example.com",
			AuthorName:  "Bob",
			RepoName:    "acme/api",
			BaseBranch:  "main",
			HeadBranch:  "fix/login",
			Additions:   50,
			Deletions:   20,
			AIRatio:     0.5,
			TabLines:    25,
			CreatedAt:   now.Add(-48 * time.Hour),
			MergedAt:    ptrTime(now.Add(-24 * time.Hour)),
		},
		{
			ID:          3,
			Number:      3,
			Title:       "chore: update dependencies",
			State:       models.PRStateClosed,
			AuthorEmail: "alice@example.com",
			AuthorName:  "Alice",
			RepoName:    "acme/api",
			BaseBranch:  "main",
			HeadBranch:  "chore/deps",
			Additions:   10,
			Deletions:   5,
			CreatedAt:   now.Add(-96 * time.Hour),
			ClosedAt:    ptrTime(now.Add(-72 * time.Hour)),
		},
		{
			ID:          4,
			Number:      4,
			Title:       "feat: add new dashboard",
			State:       models.PRStateMerged,
			AuthorEmail: "alice@example.com",
			AuthorName:  "Alice",
			RepoName:    "acme/frontend",
			BaseBranch:  "main",
			HeadBranch:  "feature/dashboard",
			Additions:   200,
			Deletions:   50,
			CreatedAt:   now.Add(-120 * time.Hour),
			MergedAt:    ptrTime(now.Add(-24 * time.Hour)),
		},
	}

	for _, pr := range prs {
		_ = store.StorePR(pr)
	}

	return store
}

// ptrTime returns a pointer to a time.Time value.
func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestListPRsAnalytics(t *testing.T) {
	store := setupPRAnalyticsTestStore()

	t.Run("list all PRs without filters", func(t *testing.T) {
		handler := ListPRsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/prs", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PRsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 4, "should return all 4 PRs")
		assert.Equal(t, 4, response.Pagination.Total)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 20, response.Pagination.PageSize)
	})

	t.Run("filter by status", func(t *testing.T) {
		handler := ListPRsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/prs?status=merged", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PRsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 2, "should return 2 merged PRs")
		for _, pr := range response.Data {
			assert.Equal(t, models.PRStateMerged, pr.State)
		}
	})

	t.Run("filter by author", func(t *testing.T) {
		handler := ListPRsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/prs?author=alice@example.com", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PRsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 3, "should return 3 PRs by alice@example.com")
		for _, pr := range response.Data {
			assert.Equal(t, "alice@example.com", pr.AuthorEmail)
		}
	})

	t.Run("filter by date range", func(t *testing.T) {
		handler := ListPRsAnalytics(store)

		now := time.Now()
		// Use a wide range to ensure we get PRs
		startDate := now.Add(-150 * time.Hour).Format("2006-01-02")
		endDate := now.Add(24 * time.Hour).Format("2006-01-02")

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/prs?start_date="+startDate+"&end_date="+endDate, nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PRsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should include all PRs within the range (all 4 created in the last 150 hours)
		assert.GreaterOrEqual(t, len(response.Data), 1, "should return PRs within date range")
		// Verify the date filters are being applied
		assert.NotEmpty(t, response.Params.StartDate)
		assert.NotEmpty(t, response.Params.EndDate)
	})

	t.Run("pagination - page 1", func(t *testing.T) {
		handler := ListPRsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/prs?page=1&page_size=2", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PRsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 2, "should return 2 PRs on page 1")
		assert.Equal(t, 4, response.Pagination.Total)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 2, response.Pagination.PageSize)
	})

	t.Run("pagination - page 2", func(t *testing.T) {
		handler := ListPRsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/prs?page=2&page_size=2", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PRsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 2, "should return 2 PRs on page 2")
		assert.Equal(t, 4, response.Pagination.Total)
		assert.Equal(t, 2, response.Pagination.Page)
	})

	t.Run("combined filters", func(t *testing.T) {
		handler := ListPRsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/prs?status=merged&author=alice@example.com", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PRsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 1, "should return 1 merged PR by alice@example.com")
		assert.Equal(t, models.PRStateMerged, response.Data[0].State)
		assert.Equal(t, "alice@example.com", response.Data[0].AuthorEmail)
	})

	t.Run("invalid page parameter", func(t *testing.T) {
		handler := ListPRsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/prs?page=-1", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid page_size parameter", func(t *testing.T) {
		handler := ListPRsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/prs?page_size=0", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid date format", func(t *testing.T) {
		handler := ListPRsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/prs?start_date=invalid", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("empty result set", func(t *testing.T) {
		handler := ListPRsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/prs?author=nonexistent@example.com", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PRsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 0, "should return empty array")
		assert.Equal(t, 0, response.Pagination.Total)
	})

	t.Run("response includes metrics", func(t *testing.T) {
		handler := ListPRsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/prs", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PRsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify first PR has all expected fields
		pr := response.Data[0]
		assert.NotEmpty(t, pr.Title)
		assert.NotEmpty(t, pr.AuthorEmail)
		assert.NotEmpty(t, pr.State)
		assert.NotZero(t, pr.CreatedAt)
	})
}
