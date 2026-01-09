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

// setupReviewAnalyticsTestStore creates a test store with sample review data.
func setupReviewAnalyticsTestStore() *storage.MemoryStore {
	store := storage.NewMemoryStore()
	now := time.Now()

	// Add test PRs first
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
			CreatedAt:   now.Add(-48 * time.Hour),
			MergedAt:    ptrTime(now.Add(-24 * time.Hour)),
		},
	}

	for _, pr := range prs {
		_ = store.StorePR(pr)
	}

	// Add test reviews with various states and reviewers
	reviews := []models.Review{
		{
			ID:          1,
			PRID:        1,
			Reviewer:    "bob@example.com",
			State:       models.ReviewStateApproved,
			SubmittedAt: now.Add(-70 * time.Hour),
			Body:        "Looks good to me!",
		},
		{
			ID:          2,
			PRID:        1,
			Reviewer:    "charlie@example.com",
			State:       models.ReviewStateChangesRequested,
			SubmittedAt: now.Add(-68 * time.Hour),
			Body:        "Please add more tests",
		},
		{
			ID:          3,
			PRID:        2,
			Reviewer:    "alice@example.com",
			State:       models.ReviewStateApproved,
			SubmittedAt: now.Add(-46 * time.Hour),
			Body:        "LGTM",
		},
		{
			ID:          4,
			PRID:        2,
			Reviewer:    "bob@example.com",
			State:       models.ReviewStateApproved,
			SubmittedAt: now.Add(-45 * time.Hour),
			Body:        "Approved",
		},
	}

	for _, review := range reviews {
		_ = store.StoreReview(review)
	}

	return store
}

func TestListReviewsAnalytics(t *testing.T) {
	store := setupReviewAnalyticsTestStore()

	t.Run("list all reviews without filters", func(t *testing.T) {
		handler := ListReviewsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/reviews", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ReviewsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 4, "should return all 4 reviews")
		assert.Equal(t, 4, response.Pagination.Total)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 20, response.Pagination.PageSize)
	})

	t.Run("filter by pr_id", func(t *testing.T) {
		handler := ListReviewsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/reviews?pr_id=1", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ReviewsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 2, "should return 2 reviews for PR #1")
		for _, review := range response.Data {
			assert.Equal(t, 1, review.PRID)
		}
	})

	t.Run("filter by reviewer", func(t *testing.T) {
		handler := ListReviewsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/reviews?reviewer=bob@example.com", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ReviewsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 2, "should return 2 reviews by bob@example.com")
		for _, review := range response.Data {
			assert.Equal(t, "bob@example.com", review.Reviewer)
		}
	})

	t.Run("combined filters - pr_id and reviewer", func(t *testing.T) {
		handler := ListReviewsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/reviews?pr_id=1&reviewer=bob@example.com", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ReviewsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 1, "should return 1 review matching both filters")
		assert.Equal(t, 1, response.Data[0].PRID)
		assert.Equal(t, "bob@example.com", response.Data[0].Reviewer)
	})

	t.Run("pagination - page 1", func(t *testing.T) {
		handler := ListReviewsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/reviews?page=1&page_size=2", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ReviewsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 2, "should return 2 reviews on page 1")
		assert.Equal(t, 4, response.Pagination.Total)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 2, response.Pagination.PageSize)
	})

	t.Run("pagination - page 2", func(t *testing.T) {
		handler := ListReviewsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/reviews?page=2&page_size=2", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ReviewsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 2, "should return 2 reviews on page 2")
		assert.Equal(t, 4, response.Pagination.Total)
		assert.Equal(t, 2, response.Pagination.Page)
	})

	t.Run("empty result set", func(t *testing.T) {
		handler := ListReviewsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/reviews?reviewer=nonexistent@example.com", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ReviewsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 0, "should return empty array")
		assert.Equal(t, 0, response.Pagination.Total)
	})

	t.Run("invalid pr_id parameter", func(t *testing.T) {
		handler := ListReviewsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/reviews?pr_id=invalid", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid page parameter", func(t *testing.T) {
		handler := ListReviewsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/reviews?page=-1", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid page_size parameter", func(t *testing.T) {
		handler := ListReviewsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/reviews?page_size=0", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("response includes all fields", func(t *testing.T) {
		handler := ListReviewsAnalytics(store)

		req := httptest.NewRequest(http.MethodGet, "/analytics/github/reviews", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ReviewsAnalyticsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify first review has all expected fields
		review := response.Data[0]
		assert.NotZero(t, review.ID)
		assert.NotZero(t, review.PRID)
		assert.NotEmpty(t, review.Reviewer)
		assert.NotEmpty(t, review.State)
		assert.NotZero(t, review.SubmittedAt)
	})
}
