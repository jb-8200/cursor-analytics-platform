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

func TestPRCycleTimeAnalytics_BasicMetrics(t *testing.T) {
	store := storage.NewMemoryStore()
	now := time.Now()

	// Create test PRs with known cycle times
	prs := []models.PullRequest{
		{
			ID:        1,
			Number:    1,
			State:     models.PRStateMerged,
			AuthorID:  "dev1",
			RepoName:  "test-repo",
			CreatedAt: now.Add(-10 * 24 * time.Hour),
			MergedAt:  timePtr(now.Add(-3 * 24 * time.Hour)),
			FirstReviewAt: timePtr(now.Add(-9 * 24 * time.Hour)), // 1 day to first review
		},
		{
			ID:        2,
			Number:    2,
			State:     models.PRStateMerged,
			AuthorID:  "dev2",
			RepoName:  "test-repo",
			CreatedAt: now.Add(-8 * 24 * time.Hour),
			MergedAt:  timePtr(now.Add(-2 * 24 * time.Hour)),
			FirstReviewAt: timePtr(now.Add(-6 * 24 * time.Hour)), // 2 days to first review
		},
		{
			ID:        3,
			Number:    3,
			State:     models.PRStateMerged,
			AuthorID:  "dev3",
			RepoName:  "test-repo",
			CreatedAt: now.Add(-6 * 24 * time.Hour),
			MergedAt:  timePtr(now.Add(-1 * 24 * time.Hour)),
			FirstReviewAt: timePtr(now.Add(-3 * 24 * time.Hour)), // 3 days to first review
		},
	}

	for _, pr := range prs {
		require.NoError(t, store.StorePR(pr))
	}

	handler := PRCycleTimeAnalytics(store)
	req := httptest.NewRequest(http.MethodGet, "/analytics/github/pr-cycle-time", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PRCycleTimeResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Avg time to first review: (1 + 2 + 3) / 3 = 2 days = 172800 seconds
	assert.InDelta(t, 172800.0, response.Data.AvgTimeToFirstReview, 3600.0) // Allow 1 hour delta

	// Avg time to merge: (7 + 6 + 5) / 3 = 6 days = 518400 seconds
	assert.InDelta(t, 518400.0, response.Data.AvgTimeToMerge, 3600.0)

	assert.Equal(t, 3, response.Data.TotalPRsAnalyzed)
	assert.True(t, response.Data.MedianTimeToMerge > 0)
	assert.True(t, response.Data.P50TimeToMerge > 0)
	assert.True(t, response.Data.P75TimeToMerge > 0)
	assert.True(t, response.Data.P90TimeToMerge > 0)
}

func TestPRCycleTimeAnalytics_DateRangeFiltering(t *testing.T) {
	store := storage.NewMemoryStore()
	now := time.Now()

	prs := []models.PullRequest{
		{
			ID:        1,
			Number:    1,
			State:     models.PRStateMerged,
			AuthorID:  "dev1",
			RepoName:  "test-repo",
			CreatedAt: now.Add(-30 * 24 * time.Hour),
			MergedAt:  timePtr(now.Add(-25 * 24 * time.Hour)),
			FirstReviewAt: timePtr(now.Add(-29 * 24 * time.Hour)),
		},
		{
			ID:        2,
			Number:    2,
			State:     models.PRStateMerged,
			AuthorID:  "dev2",
			RepoName:  "test-repo",
			CreatedAt: now.Add(-10 * 24 * time.Hour),
			MergedAt:  timePtr(now.Add(-5 * 24 * time.Hour)),
			FirstReviewAt: timePtr(now.Add(-9 * 24 * time.Hour)),
		},
	}

	for _, pr := range prs {
		require.NoError(t, store.StorePR(pr))
	}

	// Filter to only include recent PR
	handler := PRCycleTimeAnalytics(store)
	req := httptest.NewRequest(http.MethodGet, "/analytics/github/pr-cycle-time?from="+now.Add(-15*24*time.Hour).Format("2006-01-02")+"&to="+now.Format("2006-01-02"), nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PRCycleTimeResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Should only include PR #2
	assert.Equal(t, 1, response.Data.TotalPRsAnalyzed)
}

func TestPRCycleTimeAnalytics_EmptyData(t *testing.T) {
	store := storage.NewMemoryStore()

	handler := PRCycleTimeAnalytics(store)
	req := httptest.NewRequest(http.MethodGet, "/analytics/github/pr-cycle-time", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PRCycleTimeResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, 0.0, response.Data.AvgTimeToFirstReview)
	assert.Equal(t, 0.0, response.Data.AvgTimeToMerge)
	assert.Equal(t, 0.0, response.Data.MedianTimeToMerge)
	assert.Equal(t, 0, response.Data.TotalPRsAnalyzed)
}

func TestPRCycleTimeAnalytics_NoReviews(t *testing.T) {
	store := storage.NewMemoryStore()
	now := time.Now()

	// PR without FirstReviewAt
	pr := models.PullRequest{
		ID:        1,
		Number:    1,
		State:     models.PRStateMerged,
		AuthorID:  "dev1",
		RepoName:  "test-repo",
		CreatedAt: now.Add(-10 * 24 * time.Hour),
		MergedAt:  timePtr(now.Add(-3 * 24 * time.Hour)),
		FirstReviewAt: nil, // No review
	}
	require.NoError(t, store.StorePR(pr))

	handler := PRCycleTimeAnalytics(store)
	req := httptest.NewRequest(http.MethodGet, "/analytics/github/pr-cycle-time", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PRCycleTimeResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Avg time to first review should be 0 (no reviews)
	assert.Equal(t, 0.0, response.Data.AvgTimeToFirstReview)

	// But should still have merge time
	assert.True(t, response.Data.AvgTimeToMerge > 0)
	assert.Equal(t, 1, response.Data.TotalPRsAnalyzed)
}

func TestPRCycleTimeAnalytics_NoMergedPRs(t *testing.T) {
	store := storage.NewMemoryStore()
	now := time.Now()

	// Only open or closed PRs, no merged
	prs := []models.PullRequest{
		{
			ID:        1,
			Number:    1,
			State:     models.PRStateOpen,
			AuthorID:  "dev1",
			RepoName:  "test-repo",
			CreatedAt: now.Add(-10 * 24 * time.Hour),
			FirstReviewAt: timePtr(now.Add(-9 * 24 * time.Hour)),
		},
		{
			ID:        2,
			Number:    2,
			State:     models.PRStateClosed,
			AuthorID:  "dev2",
			RepoName:  "test-repo",
			CreatedAt: now.Add(-8 * 24 * time.Hour),
			ClosedAt:  timePtr(now.Add(-2 * 24 * time.Hour)),
		},
	}

	for _, pr := range prs {
		require.NoError(t, store.StorePR(pr))
	}

	handler := PRCycleTimeAnalytics(store)
	req := httptest.NewRequest(http.MethodGet, "/analytics/github/pr-cycle-time", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PRCycleTimeResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// No merged PRs
	assert.Equal(t, 0.0, response.Data.AvgTimeToMerge)
	assert.Equal(t, 0.0, response.Data.MedianTimeToMerge)
	assert.Equal(t, 0, response.Data.TotalPRsAnalyzed)
}

func TestPRCycleTimeAnalytics_PercentilesCorrect(t *testing.T) {
	store := storage.NewMemoryStore()
	now := time.Now()

	// Create PRs with known merge times
	mergeTimes := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} // days
	for i, days := range mergeTimes {
		pr := models.PullRequest{
			ID:        i + 1,
			Number:    i + 1,
			State:     models.PRStateMerged,
			AuthorID:  "dev",
			RepoName:  "test-repo",
			CreatedAt: now.Add(time.Duration(-days-2) * 24 * time.Hour),
			MergedAt:  timePtr(now.Add(time.Duration(-2) * 24 * time.Hour)),
			FirstReviewAt: timePtr(now.Add(time.Duration(-days-1) * 24 * time.Hour)),
		}
		require.NoError(t, store.StorePR(pr))
	}

	handler := PRCycleTimeAnalytics(store)
	req := httptest.NewRequest(http.MethodGet, "/analytics/github/pr-cycle-time", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PRCycleTimeResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// P50 should be around 5 days
	assert.InDelta(t, 5*24*3600.0, response.Data.P50TimeToMerge, 2*24*3600.0)

	// P75 should be around 7-8 days
	assert.InDelta(t, 7.5*24*3600.0, response.Data.P75TimeToMerge, 2*24*3600.0)

	// P90 should be around 9 days
	assert.InDelta(t, 9*24*3600.0, response.Data.P90TimeToMerge, 2*24*3600.0)

	// P90 >= P75 >= P50
	assert.True(t, response.Data.P90TimeToMerge >= response.Data.P75TimeToMerge)
	assert.True(t, response.Data.P75TimeToMerge >= response.Data.P50TimeToMerge)
}

func TestPRCycleTimeAnalytics_InvalidDateFormat(t *testing.T) {
	store := storage.NewMemoryStore()

	handler := PRCycleTimeAnalytics(store)
	req := httptest.NewRequest(http.MethodGet, "/analytics/github/pr-cycle-time?from=invalid-date", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPRCycleTimeAnalytics_SinglePR(t *testing.T) {
	store := storage.NewMemoryStore()
	now := time.Now()

	pr := models.PullRequest{
		ID:        1,
		Number:    1,
		State:     models.PRStateMerged,
		AuthorID:  "dev1",
		RepoName:  "test-repo",
		CreatedAt: now.Add(-5 * 24 * time.Hour),
		MergedAt:  timePtr(now.Add(-1 * 24 * time.Hour)),
		FirstReviewAt: timePtr(now.Add(-4 * 24 * time.Hour)),
	}
	require.NoError(t, store.StorePR(pr))

	handler := PRCycleTimeAnalytics(store)
	req := httptest.NewRequest(http.MethodGet, "/analytics/github/pr-cycle-time", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PRCycleTimeResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// With single PR, all percentiles should be the same
	assert.Equal(t, response.Data.MedianTimeToMerge, response.Data.P50TimeToMerge)
	assert.Equal(t, response.Data.P50TimeToMerge, response.Data.P75TimeToMerge)
	assert.Equal(t, response.Data.P75TimeToMerge, response.Data.P90TimeToMerge)
	assert.Equal(t, 1, response.Data.TotalPRsAnalyzed)
}

// Helper function to create time pointers
func timePtr(t time.Time) *time.Time {
	return &t
}
