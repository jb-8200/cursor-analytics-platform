package qualtrics

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestSeedData creates a minimal seed data for testing.
func createTestSeedData() *seed.SeedData {
	return &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev1@company.com"},
			{Email: "dev2@company.com"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				ResponseCount: 10,
			},
		},
	}
}

func TestQualtricsStartExport_Success(t *testing.T) {
	seedData := createTestSeedData()
	gen := generator.NewSurveyGenerator(seedData)
	manager := services.NewExportJobManager(gen)
	handlers := NewExportHandlers(manager)

	req := httptest.NewRequest("POST", "/API/v3/surveys/SV_abc123/export-responses", nil)
	rr := httptest.NewRecorder()

	handlers.StartExportHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var resp models.ExportStartResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	assert.NotEmpty(t, resp.Result.ProgressID)
	assert.Equal(t, "inProgress", resp.Result.Status)
	assert.Equal(t, 0, resp.Result.PercentComplete)
	assert.Equal(t, "200 - OK", resp.Meta.HTTPStatus)
}

func TestQualtricsProgress_InProgress(t *testing.T) {
	seedData := createTestSeedData()
	gen := generator.NewSurveyGenerator(seedData)
	manager := services.NewExportJobManager(gen)
	handlers := NewExportHandlers(manager)

	// Start export first
	job, err := manager.StartExport("SV_abc123")
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/"+job.ProgressID, nil)
	rr := httptest.NewRecorder()

	handlers.ProgressHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp models.ExportProgressResponse
	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	// Progress should have advanced from 0
	assert.True(t, resp.Result.PercentComplete > 0)
	assert.Equal(t, "inProgress", resp.Result.Status)
}

func TestQualtricsProgress_Complete(t *testing.T) {
	seedData := createTestSeedData()
	gen := generator.NewSurveyGenerator(seedData)
	manager := services.NewExportJobManager(gen)
	handlers := NewExportHandlers(manager)

	job, _ := manager.StartExport("SV_abc123")

	// Poll until complete
	var resp models.ExportProgressResponse
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/"+job.ProgressID, nil)
		rr := httptest.NewRecorder()

		handlers.ProgressHandler().ServeHTTP(rr, req)

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)

		if resp.Result.Status == "complete" {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	assert.Equal(t, "complete", resp.Result.Status)
	assert.Equal(t, 100, resp.Result.PercentComplete)
	assert.NotEmpty(t, resp.Result.FileID)
}

func TestQualtricsFileDownload_Success(t *testing.T) {
	seedData := createTestSeedData()
	gen := generator.NewSurveyGenerator(seedData)
	manager := services.NewExportJobManager(gen)
	handlers := NewExportHandlers(manager)

	job, _ := manager.StartExport("SV_abc123")

	// Poll until complete
	for job.Status == models.ExportStatusInProgress {
		job, _ = manager.GetProgress(job.ProgressID)
	}

	req := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/"+job.FileID+"/file", nil)
	rr := httptest.NewRecorder()

	handlers.FileDownloadHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/zip", rr.Header().Get("Content-Type"))
	assert.Contains(t, rr.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, rr.Header().Get("Content-Disposition"), "survey_responses.zip")

	// Verify it's a valid ZIP
	body := rr.Body.Bytes()
	assert.True(t, len(body) > 0)

	reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	require.NoError(t, err)
	assert.Len(t, reader.File, 1)
	assert.Equal(t, "survey_responses.csv", reader.File[0].Name)
}

func TestQualtricsStartExport_NoAuth(t *testing.T) {
	seedData := createTestSeedData()
	gen := generator.NewSurveyGenerator(seedData)
	manager := services.NewExportJobManager(gen)
	handlers := NewExportHandlers(manager)

	req := httptest.NewRequest("POST", "/API/v3/surveys/SV_abc123/export-responses", nil)
	// Don't set authentication

	rr := httptest.NewRecorder()
	handlers.StartExportHandler().ServeHTTP(rr, req)

	// Note: Authentication is handled by middleware in router, not handlers
	// This test verifies handler works without auth (middleware will enforce)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestQualtricsProgress_NotFound(t *testing.T) {
	seedData := createTestSeedData()
	gen := generator.NewSurveyGenerator(seedData)
	manager := services.NewExportJobManager(gen)
	handlers := NewExportHandlers(manager)

	req := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/ES_nonexistent", nil)
	rr := httptest.NewRecorder()

	handlers.ProgressHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestQualtricsFileDownload_NotFound(t *testing.T) {
	seedData := createTestSeedData()
	gen := generator.NewSurveyGenerator(seedData)
	manager := services.NewExportJobManager(gen)
	handlers := NewExportHandlers(manager)

	req := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/FILE_nonexistent/file", nil)
	rr := httptest.NewRecorder()

	handlers.FileDownloadHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestQualtricsStartExport_MultipleSurveys(t *testing.T) {
	seedData := createTestSeedData()
	gen := generator.NewSurveyGenerator(seedData)
	manager := services.NewExportJobManager(gen)
	handlers := NewExportHandlers(manager)

	// Start exports for different survey IDs
	surveys := []string{"SV_survey1", "SV_survey2", "SV_survey3"}
	progressIDs := make([]string, len(surveys))

	for i, surveyID := range surveys {
		req := httptest.NewRequest("POST", "/API/v3/surveys/"+surveyID+"/export-responses", nil)
		rr := httptest.NewRecorder()

		handlers.StartExportHandler().ServeHTTP(rr, req)

		var resp models.ExportStartResponse
		json.NewDecoder(rr.Body).Decode(&resp)
		progressIDs[i] = resp.Result.ProgressID
	}

	// Verify all progress IDs are unique
	seen := make(map[string]bool)
	for _, id := range progressIDs {
		assert.False(t, seen[id], "Progress ID should be unique")
		seen[id] = true
	}
}

func TestQualtricsProgress_ResponseFormat(t *testing.T) {
	seedData := createTestSeedData()
	gen := generator.NewSurveyGenerator(seedData)
	manager := services.NewExportJobManager(gen)
	handlers := NewExportHandlers(manager)

	job, _ := manager.StartExport("SV_abc123")

	req := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/"+job.ProgressID, nil)
	rr := httptest.NewRecorder()

	handlers.ProgressHandler().ServeHTTP(rr, req)

	// Verify response structure matches Qualtrics API spec
	var resp models.ExportProgressResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	// Verify meta field exists
	assert.NotEmpty(t, resp.Meta.HTTPStatus)
	assert.Equal(t, "200 - OK", resp.Meta.HTTPStatus)

	// Verify result field structure
	assert.NotEmpty(t, resp.Result.Status)
	assert.True(t, resp.Result.PercentComplete >= 0 && resp.Result.PercentComplete <= 100)
}

func TestQualtricsStartExport_InvalidPath(t *testing.T) {
	seedData := createTestSeedData()
	gen := generator.NewSurveyGenerator(seedData)
	manager := services.NewExportJobManager(gen)
	handlers := NewExportHandlers(manager)

	// Invalid path without survey ID
	req := httptest.NewRequest("POST", "/API/v3/surveys//export-responses", nil)
	rr := httptest.NewRecorder()

	handlers.StartExportHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestQualtricsProgress_InvalidPath(t *testing.T) {
	seedData := createTestSeedData()
	gen := generator.NewSurveyGenerator(seedData)
	manager := services.NewExportJobManager(gen)
	handlers := NewExportHandlers(manager)

	// Invalid path without progress ID
	req := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/", nil)
	rr := httptest.NewRecorder()

	handlers.ProgressHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestQualtricsFileDownload_InvalidPath(t *testing.T) {
	seedData := createTestSeedData()
	gen := generator.NewSurveyGenerator(seedData)
	manager := services.NewExportJobManager(gen)
	handlers := NewExportHandlers(manager)

	// Invalid path without file ID
	req := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses//file", nil)
	rr := httptest.NewRecorder()

	handlers.FileDownloadHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
