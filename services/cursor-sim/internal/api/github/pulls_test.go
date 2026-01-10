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

func setupTestStore() *storage.MemoryStore {
	store := storage.NewMemoryStore()
	now := time.Now()

	// Add some test PRs
	_ = store.AddPR(models.PullRequest{
		Number:      1,
		Title:       "feat: add new feature",
		Body:        "This is a new feature",
		State:       models.PRStateOpen,
		AuthorID:    "user_001",
		AuthorEmail: "alice@example.com",
		AuthorName:  "Alice",
		RepoName:    "acme/api",
		BaseBranch:  "main",
		HeadBranch:  "feature/new",
		Additions:   100,
		Deletions:   20,
		AIRatio:     0.6,
		TabLines:    60,
		CreatedAt:   now.Add(-24 * time.Hour),
		UpdatedAt:   now,
	})

	_ = store.AddPR(models.PullRequest{
		Number:      2,
		Title:       "fix: resolve bug",
		State:       models.PRStateMerged,
		AuthorID:    "user_002",
		RepoName:    "acme/api",
		BaseBranch:  "main",
		HeadBranch:  "fix/bug",
		Additions:   50,
		Deletions:   10,
		CreatedAt:   now.Add(-48 * time.Hour),
		MergedAt:    &now,
	})

	_ = store.AddPR(models.PullRequest{
		Number:    1,
		Title:     "chore: update deps",
		State:     models.PRStateOpen,
		AuthorID:  "user_001",
		RepoName:  "acme/frontend",
		CreatedAt: now,
	})

	return store
}

func TestListPulls(t *testing.T) {
	store := setupTestStore()

	t.Run("list all PRs for repo", func(t *testing.T) {
		handler := ListPulls(store)

		req := httptest.NewRequest(http.MethodGet, "/repos/acme/api/pulls", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []models.PullRequest
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response, 2, "should return 2 PRs for acme/api")
	})

	t.Run("filter by state", func(t *testing.T) {
		handler := ListPulls(store)

		req := httptest.NewRequest(http.MethodGet, "/repos/acme/api/pulls?state=open", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []models.PullRequest
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response, 1, "should return 1 open PR")
		assert.Equal(t, models.PRStateOpen, response[0].State)
	})

	t.Run("pagination", func(t *testing.T) {
		handler := ListPulls(store)

		req := httptest.NewRequest(http.MethodGet, "/repos/acme/api/pulls?per_page=1&page=1", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []models.PullRequest
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response, 1, "should return 1 PR per page")
	})
}

func TestGetPull(t *testing.T) {
	store := setupTestStore()

	t.Run("get existing PR", func(t *testing.T) {
		handler := GetPull(store)

		req := httptest.NewRequest(http.MethodGet, "/repos/acme/api/pulls/1", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.PullRequest
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.Number)
		assert.Equal(t, "feat: add new feature", response.Title)
	})

	t.Run("PR not found", func(t *testing.T) {
		handler := GetPull(store)

		req := httptest.NewRequest(http.MethodGet, "/repos/acme/api/pulls/999", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestListPullReviews(t *testing.T) {
	store := setupTestStore()

	// Get the PR to find its ID
	pr, err := store.GetPR("acme/api", 1)
	require.NoError(t, err)
	require.NotNil(t, pr)

	// Add some reviews
	now := time.Now()
	_ = store.StoreReview(models.Review{
		ID:          1,
		PRID:        int(pr.ID),
		Reviewer:    "user_002",
		Body:        "LGTM",
		State:       models.ReviewStateApproved,
		SubmittedAt: now,
	})

	_ = store.StoreReview(models.Review{
		ID:          2,
		PRID:        int(pr.ID),
		Reviewer:    "user_003",
		Body:        "Needs refactoring",
		State:       models.ReviewStateChangesRequested,
		SubmittedAt: now,
	})

	t.Run("list reviews for PR", func(t *testing.T) {
		handler := ListPullReviews(store)

		req := httptest.NewRequest(http.MethodGet, "/repos/acme/api/pulls/1/reviews", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []models.Review
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response, 2, "should return 2 reviews")
	})

	t.Run("no reviews returns empty array", func(t *testing.T) {
		handler := ListPullReviews(store)

		req := httptest.NewRequest(http.MethodGet, "/repos/acme/api/pulls/2/reviews", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []models.Review
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response, 0, "should return empty array")
	})
}
