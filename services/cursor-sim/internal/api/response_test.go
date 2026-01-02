package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRespondJSON_Success(t *testing.T) {
	rec := httptest.NewRecorder()

	data := map[string]string{"message": "hello"}
	err := RespondJSON(rec, 200, data)
	require.NoError(t, err)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var result map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "hello", result["message"])
}

func TestRespondJSON_EmptyData(t *testing.T) {
	rec := httptest.NewRecorder()

	err := RespondJSON(rec, 200, []string{})
	require.NoError(t, err)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "[]\n", rec.Body.String())
}

func TestRespondJSON_NilData(t *testing.T) {
	rec := httptest.NewRecorder()

	err := RespondJSON(rec, 200, nil)
	require.NoError(t, err)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "null\n", rec.Body.String())
}

func TestRespondError_BasicError(t *testing.T) {
	rec := httptest.NewRecorder()

	RespondError(rec, 400, "bad request")

	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var result map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "bad request", result["error"])
}

func TestRespondError_ServerError(t *testing.T) {
	rec := httptest.NewRecorder()

	RespondError(rec, 500, "internal error")

	assert.Equal(t, 500, rec.Code)

	var result map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "internal error", result["error"])
}

func TestBuildPaginatedResponse_FirstPage(t *testing.T) {
	data := []string{"item1", "item2", "item3"}
	params := models.Params{
		From:     "2024-01-01",
		To:       "2024-01-31",
		Page:     1,
		PageSize: 10,
	}

	resp := BuildPaginatedResponse(data, params, 25)

	assert.Equal(t, data, resp.Data)
	assert.Equal(t, 1, resp.Pagination.Page)
	assert.Equal(t, 10, resp.Pagination.PageSize)
	assert.Equal(t, 3, resp.Pagination.TotalPages) // 25 items / 10 per page = 3 pages
	assert.True(t, resp.Pagination.HasNextPage)
	assert.False(t, resp.Pagination.HasPreviousPage)
}

func TestBuildPaginatedResponse_MiddlePage(t *testing.T) {
	data := []string{"item11", "item12"}
	params := models.Params{
		Page:     2,
		PageSize: 10,
	}

	resp := BuildPaginatedResponse(data, params, 25)

	assert.Equal(t, 2, resp.Pagination.Page)
	assert.True(t, resp.Pagination.HasNextPage)
	assert.True(t, resp.Pagination.HasPreviousPage)
}

func TestBuildPaginatedResponse_LastPage(t *testing.T) {
	data := []string{"item21", "item22", "item23", "item24", "item25"}
	params := models.Params{
		Page:     3,
		PageSize: 10,
	}

	resp := BuildPaginatedResponse(data, params, 25)

	assert.Equal(t, 3, resp.Pagination.Page)
	assert.False(t, resp.Pagination.HasNextPage)
	assert.True(t, resp.Pagination.HasPreviousPage)
}

func TestBuildPaginatedResponse_EmptyData(t *testing.T) {
	data := []models.Commit{}
	params := models.Params{
		Page:     1,
		PageSize: 10,
	}

	resp := BuildPaginatedResponse(data, params, 0)

	assert.Equal(t, 0, len(resp.Data.([]models.Commit)))
	assert.Equal(t, 0, resp.Pagination.TotalPages)
	assert.False(t, resp.Pagination.HasNextPage)
	assert.False(t, resp.Pagination.HasPreviousPage)
}

func TestBuildPaginatedResponse_SinglePage(t *testing.T) {
	data := []string{"item1", "item2"}
	params := models.Params{
		Page:     1,
		PageSize: 10,
	}

	resp := BuildPaginatedResponse(data, params, 2)

	assert.Equal(t, 1, resp.Pagination.TotalPages)
	assert.False(t, resp.Pagination.HasNextPage)
	assert.False(t, resp.Pagination.HasPreviousPage)
}

func TestRespondCSV_Commits(t *testing.T) {
	rec := httptest.NewRecorder()

	commits := []models.Commit{
		{
			CommitHash:         "abc123",
			UserID:             "user_001",
			UserEmail:          "test@example.com",
			RepoName:           "acme/api",
			TotalLinesAdded:    100,
			TabLinesAdded:      60,
			ComposerLinesAdded: 20,
			NonAILinesAdded:    20,
			CommitTs:           time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			CreatedAt:          time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			CommitHash:         "def456",
			UserID:             "user_002",
			UserEmail:          "alice@example.com",
			RepoName:           "acme/web",
			TotalLinesAdded:    50,
			TabLinesAdded:      30,
			ComposerLinesAdded: 10,
			NonAILinesAdded:    10,
			CommitTs:           time.Date(2024, 1, 2, 14, 0, 0, 0, time.UTC),
			CreatedAt:          time.Date(2024, 1, 2, 14, 0, 0, 0, time.UTC),
		},
	}

	err := RespondCSV(rec, commits)
	require.NoError(t, err)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "text/csv", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, rec.Header().Get("Content-Disposition"), "filename=")

	csv := rec.Body.String()
	assert.Contains(t, csv, "commitHash")
	assert.Contains(t, csv, "userId")
	assert.Contains(t, csv, "abc123")
	assert.Contains(t, csv, "user_001")
	assert.Contains(t, csv, "def456")
	assert.Contains(t, csv, "user_002")
}

func TestRespondCSV_EmptyData(t *testing.T) {
	rec := httptest.NewRecorder()

	err := RespondCSV(rec, []models.Commit{})
	require.NoError(t, err)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "text/csv", rec.Header().Get("Content-Type"))

	csv := rec.Body.String()
	// Should have header only
	assert.Contains(t, csv, "commitHash")
	// Count newlines (header only = 1 newline)
	lines := 0
	for _, c := range csv {
		if c == '\n' {
			lines++
		}
	}
	assert.Equal(t, 1, lines, "should have only header row")
}

func TestRespondCSV_FilenameTimestamp(t *testing.T) {
	rec := httptest.NewRecorder()

	commits := []models.Commit{
		{CommitHash: "abc", UserID: "user_001", CommitTs: time.Now(), CreatedAt: time.Now()},
	}

	err := RespondCSV(rec, commits)
	require.NoError(t, err)

	disposition := rec.Header().Get("Content-Disposition")
	assert.Contains(t, disposition, ".csv")
	// Filename should contain timestamp
	assert.Regexp(t, `filename=.*\d{8}.*\.csv`, disposition)
}

func TestParseQueryParams_ValidParams(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?from=2024-01-01&to=2024-01-31&page=2&pageSize=50", nil)

	params, err := ParseQueryParams(req)
	require.NoError(t, err)

	assert.Equal(t, "2024-01-01", params.From)
	assert.Equal(t, "2024-01-31", params.To)
	assert.Equal(t, 2, params.Page)
	assert.Equal(t, 50, params.PageSize)
}

func TestParseQueryParams_DefaultValues(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	params, err := ParseQueryParams(req)
	require.NoError(t, err)

	assert.Equal(t, 1, params.Page, "default page should be 1")
	assert.Equal(t, 100, params.PageSize, "default pageSize should be 100")
	assert.NotEmpty(t, params.From, "default from should be set")
	assert.NotEmpty(t, params.To, "default to should be set")
}

func TestParseQueryParams_InvalidPage(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?page=0", nil)

	_, err := ParseQueryParams(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "page")
}

func TestParseQueryParams_InvalidPageSize(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?pageSize=0", nil)

	_, err := ParseQueryParams(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pageSize")
}

func TestParseQueryParams_PageSizeTooLarge(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?pageSize=1001", nil)

	_, err := ParseQueryParams(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pageSize")
}

func TestParseQueryParams_InvalidDateFormat(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?from=2024/01/01", nil)

	_, err := ParseQueryParams(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "from")
}

func TestParseQueryParams_UserIDFilter(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?userId=user_001", nil)

	params, err := ParseQueryParams(req)
	require.NoError(t, err)

	assert.Equal(t, "user_001", params.UserID)
}

func TestParseQueryParams_RepoNameFilter(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?repoName=acme/api", nil)

	params, err := ParseQueryParams(req)
	require.NoError(t, err)

	assert.Equal(t, "acme/api", params.RepoName)
}
