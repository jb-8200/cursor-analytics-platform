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

func TestReviewQualityAnalytics_BasicMetrics(t *testing.T) {
	store := storage.NewMemoryStore()

	// Setup: Create PRs with reviews
	// PR 1: 2 approvals, 1 changes requested
	pr1 := models.PullRequest{
		ID:        1,
		Number:    1,
		RepoName:  "test/repo",
		State:     models.PRStateMerged,
		CreatedAt: time.Date(2026, 1, 5, 10, 0, 0, 0, time.UTC),
		MergedAt:  timePtr(time.Date(2026, 1, 6, 10, 0, 0, 0, time.UTC)),
	}
	require.NoError(t, store.StorePR(pr1))

	review1 := models.Review{
		ID:          1,
		PRID:        1,
		Reviewer:    "reviewer1@example.com",
		State:       models.ReviewStateApproved,
		SubmittedAt: time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC),
		Comments:    []models.ReviewComment{{Body: "LGTM"}, {Body: "Nice work"}},
	}
	require.NoError(t, store.StoreReview(review1))

	review2 := models.Review{
		ID:          2,
		PRID:        1,
		Reviewer:    "reviewer2@example.com",
		State:       models.ReviewStateApproved,
		SubmittedAt: time.Date(2026, 1, 5, 14, 0, 0, 0, time.UTC),
		Comments:    []models.ReviewComment{{Body: "Good"}},
	}
	require.NoError(t, store.StoreReview(review2))

	review3 := models.Review{
		ID:          3,
		PRID:        1,
		Reviewer:    "reviewer3@example.com",
		State:       models.ReviewStateChangesRequested,
		SubmittedAt: time.Date(2026, 1, 5, 15, 0, 0, 0, time.UTC),
		Comments:    []models.ReviewComment{{Body: "Fix this"}, {Body: "Change that"}, {Body: "One more thing"}},
	}
	require.NoError(t, store.StoreReview(review3))

	// PR 2: 1 approval, 1 pending
	pr2 := models.PullRequest{
		ID:        2,
		Number:    2,
		RepoName:  "test/repo",
		State:     models.PRStateMerged,
		CreatedAt: time.Date(2026, 1, 6, 10, 0, 0, 0, time.UTC),
		MergedAt:  timePtr(time.Date(2026, 1, 7, 10, 0, 0, 0, time.UTC)),
	}
	require.NoError(t, store.StorePR(pr2))

	review4 := models.Review{
		ID:          4,
		PRID:        2,
		Reviewer:    "reviewer1@example.com",
		State:       models.ReviewStateApproved,
		SubmittedAt: time.Date(2026, 1, 6, 12, 0, 0, 0, time.UTC),
		Comments:    []models.ReviewComment{},
	}
	require.NoError(t, store.StoreReview(review4))

	review5 := models.Review{
		ID:          5,
		PRID:        2,
		Reviewer:    "reviewer4@example.com",
		State:       models.ReviewStatePending,
		SubmittedAt: time.Date(2026, 1, 6, 13, 0, 0, 0, time.UTC),
		Comments:    []models.ReviewComment{{Body: "Checking"}},
	}
	require.NoError(t, store.StoreReview(review5))

	// Make request
	handler := ReviewQualityAnalytics(store)
	req := httptest.NewRequest("GET", "/analytics/github/review-quality", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Verify response
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp ReviewQualityResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

	// Verify metrics
	// Total reviews: 5
	// Approvals: 3 → approval rate = 3/5 = 0.6
	// Changes requested: 1 → changes requested rate = 1/5 = 0.2
	// Pending: 1 → pending rate = 1/5 = 0.2
	assert.Equal(t, 5, resp.Data.TotalReviews)
	assert.InDelta(t, 0.6, resp.Data.ApprovalRate, 0.01)
	assert.InDelta(t, 0.2, resp.Data.ChangesRequestedRate, 0.01)
	assert.InDelta(t, 0.2, resp.Data.PendingRate, 0.01)

	// Total comments: 2 + 1 + 3 + 0 + 1 = 7
	// Avg comments per review: 7/5 = 1.4
	assert.InDelta(t, 1.4, resp.Data.AvgCommentsPerReview, 0.01)

	// Total PRs reviewed: 2
	assert.Equal(t, 2, resp.Data.TotalPRsReviewed)

	// Unique reviewers on merged PRs
	// PR 1 has 3 reviewers, PR 2 has 2 reviewers
	// Avg reviewers per PR = (3 + 2) / 2 = 2.5
	assert.InDelta(t, 2.5, resp.Data.AvgReviewersPerPR, 0.01)
}

func TestReviewQualityAnalytics_DateFiltering(t *testing.T) {
	store := storage.NewMemoryStore()

	// Create PRs with different merged dates
	pr1 := models.PullRequest{
		ID:        1,
		Number:    1,
		RepoName:  "test/repo",
		State:     models.PRStateMerged,
		CreatedAt: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
		MergedAt:  timePtr(time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)),
	}
	require.NoError(t, store.StorePR(pr1))

	review1 := models.Review{
		ID:          1,
		PRID:        1,
		Reviewer:    "reviewer1@example.com",
		State:       models.ReviewStateApproved,
		SubmittedAt: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
		Comments:    []models.ReviewComment{},
	}
	require.NoError(t, store.StoreReview(review1))

	pr2 := models.PullRequest{
		ID:        2,
		Number:    2,
		RepoName:  "test/repo",
		State:     models.PRStateMerged,
		CreatedAt: time.Date(2026, 1, 5, 10, 0, 0, 0, time.UTC),
		MergedAt:  timePtr(time.Date(2026, 1, 6, 10, 0, 0, 0, time.UTC)),
	}
	require.NoError(t, store.StorePR(pr2))

	review2 := models.Review{
		ID:          2,
		PRID:        2,
		Reviewer:    "reviewer2@example.com",
		State:       models.ReviewStateChangesRequested,
		SubmittedAt: time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC),
		Comments:    []models.ReviewComment{},
	}
	require.NoError(t, store.StoreReview(review2))

	// Filter to only include PR 2 (merged on 2026-01-06)
	handler := ReviewQualityAnalytics(store)
	req := httptest.NewRequest("GET", "/analytics/github/review-quality?from=2026-01-05&to=2026-01-07", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp ReviewQualityResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

	// Should only include review2 (from PR2)
	assert.Equal(t, 1, resp.Data.TotalReviews)
	assert.Equal(t, 1, resp.Data.TotalPRsReviewed)
	assert.InDelta(t, 0.0, resp.Data.ApprovalRate, 0.01)       // No approvals
	assert.InDelta(t, 1.0, resp.Data.ChangesRequestedRate, 0.01) // 1 changes requested

	// Verify params
	assert.Equal(t, "2026-01-05", resp.Params.From)
	assert.Equal(t, "2026-01-07", resp.Params.To)
}

func TestReviewQualityAnalytics_NoData(t *testing.T) {
	store := storage.NewMemoryStore()

	// Make request with no data
	handler := ReviewQualityAnalytics(store)
	req := httptest.NewRequest("GET", "/analytics/github/review-quality", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp ReviewQualityResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

	// Should return zero values
	assert.Equal(t, 0, resp.Data.TotalReviews)
	assert.Equal(t, 0, resp.Data.TotalPRsReviewed)
	assert.Equal(t, 0.0, resp.Data.ApprovalRate)
	assert.Equal(t, 0.0, resp.Data.AvgReviewersPerPR)
	assert.Equal(t, 0.0, resp.Data.AvgCommentsPerReview)
}

func TestReviewQualityAnalytics_InvalidDateFormat(t *testing.T) {
	store := storage.NewMemoryStore()

	handler := ReviewQualityAnalytics(store)

	tests := []struct {
		name  string
		query string
	}{
		{"invalid from date", "?from=invalid"},
		{"invalid to date", "?to=2026-13-45"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/analytics/github/review-quality"+tt.query, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestReviewQualityAnalytics_OnlyMergedPRs(t *testing.T) {
	store := storage.NewMemoryStore()

	// Create open PR with reviews
	prOpen := models.PullRequest{
		ID:        1,
		Number:    1,
		RepoName:  "test/repo",
		State:     models.PRStateOpen,
		CreatedAt: time.Date(2026, 1, 5, 10, 0, 0, 0, time.UTC),
	}
	require.NoError(t, store.StorePR(prOpen))

	reviewOpen := models.Review{
		ID:          1,
		PRID:        1,
		Reviewer:    "reviewer1@example.com",
		State:       models.ReviewStateApproved,
		SubmittedAt: time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC),
		Comments:    []models.ReviewComment{},
	}
	require.NoError(t, store.StoreReview(reviewOpen))

	// Create merged PR with reviews
	prMerged := models.PullRequest{
		ID:        2,
		Number:    2,
		RepoName:  "test/repo",
		State:     models.PRStateMerged,
		CreatedAt: time.Date(2026, 1, 6, 10, 0, 0, 0, time.UTC),
		MergedAt:  timePtr(time.Date(2026, 1, 7, 10, 0, 0, 0, time.UTC)),
	}
	require.NoError(t, store.StorePR(prMerged))

	reviewMerged := models.Review{
		ID:          2,
		PRID:        2,
		Reviewer:    "reviewer2@example.com",
		State:       models.ReviewStateApproved,
		SubmittedAt: time.Date(2026, 1, 6, 12, 0, 0, 0, time.UTC),
		Comments:    []models.ReviewComment{},
	}
	require.NoError(t, store.StoreReview(reviewMerged))

	handler := ReviewQualityAnalytics(store)
	req := httptest.NewRequest("GET", "/analytics/github/review-quality", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp ReviewQualityResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

	// Should only include merged PR
	assert.Equal(t, 1, resp.Data.TotalReviews)
	assert.Equal(t, 1, resp.Data.TotalPRsReviewed)
}
