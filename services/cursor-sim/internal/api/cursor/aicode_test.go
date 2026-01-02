package cursor

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestStore() *storage.MemoryStore {
	store := storage.NewMemoryStore()

	// Load developers
	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob"},
	}
	store.LoadDevelopers(developers)

	// Add commits
	now := time.Now()
	commits := []models.Commit{
		{
			CommitHash:         "commit1",
			UserID:             "user_001",
			UserEmail:          "alice@example.com",
			RepoName:           "acme/api",
			TotalLinesAdded:    100,
			TabLinesAdded:      60,
			ComposerLinesAdded: 20,
			NonAILinesAdded:    20,
			CommitTs:           now.Add(-2 * time.Hour),
			CreatedAt:          now,
		},
		{
			CommitHash:         "commit2",
			UserID:             "user_002",
			UserEmail:          "bob@example.com",
			RepoName:           "acme/web",
			TotalLinesAdded:    50,
			TabLinesAdded:      30,
			ComposerLinesAdded: 10,
			NonAILinesAdded:    10,
			CommitTs:           now.Add(-1 * time.Hour),
			CreatedAt:          now,
		},
		{
			CommitHash:         "commit3",
			UserID:             "user_001",
			UserEmail:          "alice@example.com",
			RepoName:           "acme/api",
			TotalLinesAdded:    200,
			TabLinesAdded:      120,
			ComposerLinesAdded: 40,
			NonAILinesAdded:    40,
			CommitTs:           now.Add(-10 * time.Hour), // Older
			CreatedAt:          now,
		},
	}

	for _, c := range commits {
		store.AddCommit(c)
	}

	return store
}

func TestAICodeCommits_Success(t *testing.T) {
	store := setupTestStore()
	handler := AICodeCommits(store)

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response models.CommitsResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should return all commits (default 30 days range)
	assert.GreaterOrEqual(t, len(response.Items), 2)
}

func TestAICodeCommits_WithDateRange(t *testing.T) {
	store := setupTestStore()
	handler := AICodeCommits(store)

	// Request with startDate/endDate parameters (Cursor API spec)
	startDate := time.Now().Add(-3 * time.Hour).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits?startDate="+startDate+"&endDate="+endDate, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	var response models.CommitsResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should return commits within range
	assert.NotNil(t, response.Items)
	assert.GreaterOrEqual(t, response.TotalCount, 0)
}

func TestAICodeCommits_WithUserFilter(t *testing.T) {
	store := setupTestStore()
	handler := AICodeCommits(store)

	// Use 'user' parameter (Cursor API spec)
	req := httptest.NewRequest("GET", "/analytics/ai-code/commits?user=user_001", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	var response models.CommitsResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// All commits should be from user_001
	for _, c := range response.Items {
		assert.Equal(t, "user_001", c.UserID)
	}
}

func TestAICodeCommits_WithPagination(t *testing.T) {
	store := setupTestStore()
	handler := AICodeCommits(store)

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits?page=1&pageSize=2", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	var response models.CommitsResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 2, response.PageSize)
	assert.LessOrEqual(t, len(response.Items), 2)
}

func TestAICodeCommits_EmptyStore(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := AICodeCommits(store)

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	var response models.CommitsResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should return empty array
	assert.Len(t, response.Items, 0)
	assert.Equal(t, 0, response.TotalCount)
}

func TestAICodeCommits_InvalidDateFormat(t *testing.T) {
	store := setupTestStore()
	handler := AICodeCommits(store)

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits?startDate=invalid-date", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should return error
	assert.Equal(t, 400, rec.Code)

	var errorResp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
	require.NoError(t, err)
	assert.Contains(t, errorResp["error"], "startDate")
}

func TestAICodeCommits_InvalidPagination(t *testing.T) {
	store := setupTestStore()
	handler := AICodeCommits(store)

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits?page=0", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should return error
	assert.Equal(t, 400, rec.Code)
}

func TestAICodeCommits_ResponseStructure(t *testing.T) {
	store := setupTestStore()
	handler := AICodeCommits(store)

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	// Verify CommitsResponse structure matches OpenAPI spec
	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify required fields per OpenAPI spec
	assert.Contains(t, response, "items")
	assert.Contains(t, response, "totalCount")
	assert.Contains(t, response, "page")
	assert.Contains(t, response, "pageSize")

	// Verify field types
	_, ok := response["items"].([]interface{})
	assert.True(t, ok, "items should be an array")

	page := int(response["page"].(float64))
	assert.Greater(t, page, 0)

	pageSize := int(response["pageSize"].(float64))
	assert.Greater(t, pageSize, 0)
}

func TestAICodeCommits_CommitFields(t *testing.T) {
	store := setupTestStore()
	handler := AICodeCommits(store)

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	var response models.CommitsResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	if len(response.Items) > 0 {
		commit := response.Items[0]
		// Verify required fields per OpenAPI CommitRecord schema
		assert.NotEmpty(t, commit.CommitHash)
		assert.NotEmpty(t, commit.UserID)
		assert.NotEmpty(t, commit.UserEmail)
		assert.NotEmpty(t, commit.RepoName)
		assert.GreaterOrEqual(t, commit.TotalLinesAdded, 0)
		assert.GreaterOrEqual(t, commit.TabLinesAdded, 0)
		assert.GreaterOrEqual(t, commit.ComposerLinesAdded, 0)
		assert.GreaterOrEqual(t, commit.NonAILinesAdded, 0)
		assert.False(t, commit.CommitTs.IsZero())
	}
}

func TestAICodeCommitsCSV_Success(t *testing.T) {
	store := setupTestStore()
	handler := AICodeCommitsCSV(store)

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits.csv", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "text/csv", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Header().Get("Content-Disposition"), "attachment")

	csv := rec.Body.String()
	assert.Contains(t, csv, "commitHash")
	assert.Contains(t, csv, "userId")
	assert.Contains(t, csv, "userEmail")
	assert.Contains(t, csv, "totalLinesAdded")
}

func TestAICodeCommitsCSV_WithFilters(t *testing.T) {
	store := setupTestStore()
	handler := AICodeCommitsCSV(store)

	// Use 'user' parameter (Cursor API spec)
	req := httptest.NewRequest("GET", "/analytics/ai-code/commits.csv?user=user_001", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	csv := rec.Body.String()
	// Should contain commits from user_001
	assert.Contains(t, csv, "user_001")
	assert.Contains(t, csv, "alice@example.com")
}

func TestAICodeCommitsCSV_EmptyStore(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := AICodeCommitsCSV(store)

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits.csv", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	csv := rec.Body.String()
	// Should have header only
	assert.Contains(t, csv, "commitHash")
	lines := 0
	for _, c := range csv {
		if c == '\n' {
			lines++
		}
	}
	assert.Equal(t, 1, lines, "should have only header row")
}

func TestAICodeCommitsCSV_InvalidParams(t *testing.T) {
	store := setupTestStore()
	handler := AICodeCommitsCSV(store)

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits.csv?startDate=invalid", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should return error
	assert.Equal(t, 400, rec.Code)
}
