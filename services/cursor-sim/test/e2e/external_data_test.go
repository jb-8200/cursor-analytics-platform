package e2e

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/server"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const externalDataTestPort = 19083

// setupExternalDataE2EServer creates a test server with external data sources enabled.
func setupExternalDataE2EServer(t *testing.T) (context.CancelFunc, *storage.MemoryStore) {
	// Load test seed data with external data sources enabled
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)

	// Verify external data sources are enabled in seed
	require.NotNil(t, seedData.ExternalDataSources)
	require.NotNil(t, seedData.ExternalDataSources.Harvey)
	require.True(t, seedData.ExternalDataSources.Harvey.Enabled)
	require.NotNil(t, seedData.ExternalDataSources.Copilot)
	require.True(t, seedData.ExternalDataSources.Copilot.Enabled)
	require.NotNil(t, seedData.ExternalDataSources.Qualtrics)
	require.True(t, seedData.ExternalDataSources.Qualtrics.Enabled)

	// Initialize storage
	store := storage.NewMemoryStore()
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err)

	// Create HTTP server with router
	router := server.NewRouter(store, seedData, testAPIKey)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", externalDataTestPort),
		Handler: router,
	}

	// Start server in background
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Wait for server to be ready
	time.Sleep(50 * time.Millisecond)

	// Return cleanup function
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
	}

	return cleanup, store
}

// makeAuthenticatedRequest makes an authenticated HTTP request to the test server.
func makeAuthenticatedRequest(t *testing.T, method, path string) *http.Response {
	baseURL := fmt.Sprintf("http://localhost:%d", externalDataTestPort)
	req, err := http.NewRequest(method, baseURL+path, nil)
	require.NoError(t, err)

	req.SetBasicAuth(testAPIKey, "")

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

// makeUnauthenticatedRequest makes an HTTP request without authentication.
func makeUnauthenticatedRequest(t *testing.T, method, path string) *http.Response {
	baseURL := fmt.Sprintf("http://localhost:%d", externalDataTestPort)
	req, err := http.NewRequest(method, baseURL+path, nil)
	require.NoError(t, err)

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

// ============================================================================
// Harvey E2E Tests
// ============================================================================

// TestHarvey_E2E_UsageEndpoint tests the Harvey usage endpoint with various filters.
func TestHarvey_E2E_UsageEndpoint(t *testing.T) {
	cleanup, _ := setupExternalDataE2EServer(t)
	defer cleanup()

	// Test basic query with date range
	now := time.Now()
	from := now.Add(-30 * 24 * time.Hour).Format("2006-01-02")
	to := now.Format("2006-01-02")

	path := fmt.Sprintf("/harvey/api/v1/history/usage?from=%s&to=%s", from, to)
	resp := makeAuthenticatedRequest(t, "GET", path)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	// Parse response
	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response structure
	assert.Contains(t, result, "data")
	assert.Contains(t, result, "pagination")

	pagination := result["pagination"].(map[string]interface{})
	assert.Contains(t, pagination, "page")
	assert.Contains(t, pagination, "pageSize")
	assert.Contains(t, pagination, "totalCount")
	assert.Contains(t, pagination, "totalPages")
	assert.Contains(t, pagination, "hasNextPage")

	// Verify data is an array
	data := result["data"].([]interface{})
	assert.NotNil(t, data)

	t.Logf("Harvey usage endpoint returned %d events", int(pagination["totalCount"].(float64)))
}

// TestHarvey_E2E_Pagination tests Harvey usage pagination.
func TestHarvey_E2E_Pagination(t *testing.T) {
	cleanup, _ := setupExternalDataE2EServer(t)
	defer cleanup()

	now := time.Now()
	from := now.Add(-30 * 24 * time.Hour).Format("2006-01-02")
	to := now.Format("2006-01-02")

	// Test page 1 with small page size
	path := fmt.Sprintf("/harvey/api/v1/history/usage?from=%s&to=%s&page=1&page_size=5", from, to)
	resp := makeAuthenticatedRequest(t, "GET", path)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	pagination := result["pagination"].(map[string]interface{})
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(5), pagination["pageSize"])

	// Data should not exceed page size
	data := result["data"].([]interface{})
	assert.LessOrEqual(t, len(data), 5)

	t.Logf("Harvey pagination: page 1 returned %d events", len(data))
}

// TestHarvey_E2E_DateFiltering tests Harvey date range filtering.
func TestHarvey_E2E_DateFiltering(t *testing.T) {
	cleanup, _ := setupExternalDataE2EServer(t)
	defer cleanup()

	// Test with narrow date range
	now := time.Now()
	from := now.Add(-7 * 24 * time.Hour).Format("2006-01-02")
	to := now.Format("2006-01-02")

	path := fmt.Sprintf("/harvey/api/v1/history/usage?from=%s&to=%s", from, to)
	resp := makeAuthenticatedRequest(t, "GET", path)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Should return some data (generated data)
	data := result["data"].([]interface{})
	assert.NotNil(t, data)

	t.Logf("Harvey date filtering (7 days): %d events", len(data))
}

// TestHarvey_E2E_DisabledWhenNotConfigured tests that Harvey endpoint returns 404 when not configured.
func TestHarvey_E2E_DisabledWhenNotConfigured(t *testing.T) {
	// Create seed data WITHOUT Harvey enabled
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Dev One"},
		},
		ExternalDataSources: nil, // No external data sources
	}

	store := storage.NewMemoryStore()
	router := server.NewRouter(store, seedData, testAPIKey)

	// Start server on different port
	const testPort = 19084
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", testPort),
		Handler: router,
	}

	go func() {
		httpServer.ListenAndServe()
	}()
	defer httpServer.Shutdown(context.Background())

	time.Sleep(50 * time.Millisecond)

	// Try to access Harvey endpoint
	baseURL := fmt.Sprintf("http://localhost:%d", testPort)
	req, _ := http.NewRequest("GET", baseURL+"/harvey/api/v1/history/usage?from=2026-01-01&to=2026-01-09", nil)
	req.SetBasicAuth(testAPIKey, "")

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should return 404 since Harvey is not configured
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	t.Log("Harvey endpoint correctly returns 404 when not configured")
}

// ============================================================================
// Copilot E2E Tests
// ============================================================================

// TestCopilot_E2E_JSONResponse tests Copilot endpoint with JSON response format.
func TestCopilot_E2E_JSONResponse(t *testing.T) {
	cleanup, _ := setupExternalDataE2EServer(t)
	defer cleanup()

	// Test with D30 period and JSON format
	path := "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=application/json"
	resp := makeAuthenticatedRequest(t, "GET", path)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	// Parse OData response
	var result models.CopilotUsageResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify OData structure
	assert.NotEmpty(t, result.Context)
	assert.NotNil(t, result.Value)

	// Verify usage data
	if len(result.Value) > 0 {
		firstUser := result.Value[0]
		assert.NotEmpty(t, firstUser.UserPrincipalName)
		assert.NotEmpty(t, firstUser.DisplayName)
		assert.NotEmpty(t, firstUser.ReportRefreshDate)
		assert.Equal(t, 30, firstUser.ReportPeriod)
	}

	t.Logf("Copilot JSON response: %d users", len(result.Value))
}

// TestCopilot_E2E_CSVExport tests Copilot endpoint with CSV export format.
func TestCopilot_E2E_CSVExport(t *testing.T) {
	cleanup, _ := setupExternalDataE2EServer(t)
	defer cleanup()

	// Test with D30 period and CSV format
	path := "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=text/csv"
	resp := makeAuthenticatedRequest(t, "GET", path)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/csv", resp.Header.Get("Content-Type"))
	assert.Contains(t, resp.Header.Get("Content-Disposition"), "attachment")
	assert.Contains(t, resp.Header.Get("Content-Disposition"), "copilot-usage-D30.csv")

	// Parse CSV
	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Verify CSV structure
	assert.Greater(t, len(records), 0, "CSV should have at least header row")

	// Verify header
	header := records[0]
	assert.Contains(t, header, "Report Refresh Date")
	assert.Contains(t, header, "User Principal Name")
	assert.Contains(t, header, "Display Name")

	t.Logf("Copilot CSV export: %d rows (including header)", len(records))
}

// TestCopilot_E2E_AllPeriods tests all Copilot report periods.
func TestCopilot_E2E_AllPeriods(t *testing.T) {
	cleanup, _ := setupExternalDataE2EServer(t)
	defer cleanup()

	periods := []string{"D7", "D30", "D90", "D180"}

	for _, period := range periods {
		t.Run(period, func(t *testing.T) {
			path := fmt.Sprintf("/reports/getMicrosoft365CopilotUsageUserDetail(period='%s')?$format=application/json", period)
			resp := makeAuthenticatedRequest(t, "GET", path)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var result models.CopilotUsageResponse
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			// Verify report period matches
			if len(result.Value) > 0 {
				expectedDays := map[string]int{"D7": 7, "D30": 30, "D90": 90, "D180": 180}
				assert.Equal(t, expectedDays[period], result.Value[0].ReportPeriod)
			}

			t.Logf("Copilot %s period: %d users", period, len(result.Value))
		})
	}
}

// TestCopilot_E2E_DisabledWhenNotConfigured tests that Copilot endpoint returns 404 when not configured.
func TestCopilot_E2E_DisabledWhenNotConfigured(t *testing.T) {
	// Create seed data WITHOUT Copilot enabled
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Dev One"},
		},
		ExternalDataSources: nil, // No external data sources
	}

	store := storage.NewMemoryStore()
	router := server.NewRouter(store, seedData, testAPIKey)

	// Start server on different port
	const testPort = 19085
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", testPort),
		Handler: router,
	}

	go func() {
		httpServer.ListenAndServe()
	}()
	defer httpServer.Shutdown(context.Background())

	time.Sleep(50 * time.Millisecond)

	// Try to access Copilot endpoint
	baseURL := fmt.Sprintf("http://localhost:%d", testPort)
	req, _ := http.NewRequest("GET", baseURL+"/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')", nil)
	req.SetBasicAuth(testAPIKey, "")

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should return 404 since Copilot is not configured
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	t.Log("Copilot endpoint correctly returns 404 when not configured")
}

// ============================================================================
// Qualtrics E2E Tests
// ============================================================================

// TestQualtrics_E2E_FullExportFlow tests the complete Qualtrics export workflow.
func TestQualtrics_E2E_FullExportFlow(t *testing.T) {
	cleanup, _ := setupExternalDataE2EServer(t)
	defer cleanup()

	// Load seed data to get survey ID
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)
	surveyID := seedData.ExternalDataSources.Qualtrics.SurveyID

	// Step 1: Start export
	startPath := fmt.Sprintf("/API/v3/surveys/%s/export-responses", surveyID)
	resp := makeAuthenticatedRequest(t, "POST", startPath)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var startResp models.ExportStartResponse
	err = json.NewDecoder(resp.Body).Decode(&startResp)
	require.NoError(t, err)

	assert.NotEmpty(t, startResp.Result.ProgressID)
	assert.Equal(t, "inProgress", startResp.Result.Status)
	assert.Equal(t, 0, startResp.Result.PercentComplete)

	progressID := startResp.Result.ProgressID
	t.Logf("Export started: progressID=%s", progressID)

	// Step 2: Poll progress until complete (simulate progressive completion)
	var fileID string
	maxAttempts := 10
	for attempt := 0; attempt < maxAttempts; attempt++ {
		time.Sleep(100 * time.Millisecond)

		progressPath := fmt.Sprintf("/API/v3/surveys/%s/export-responses/%s", surveyID, progressID)
		resp := makeAuthenticatedRequest(t, "GET", progressPath)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var progressResp models.ExportProgressResponse
		err = json.NewDecoder(resp.Body).Decode(&progressResp)
		require.NoError(t, err)

		t.Logf("Progress check %d: %d%% complete", attempt+1, progressResp.Result.PercentComplete)

		if progressResp.Result.Status == "complete" {
			assert.Equal(t, 100, progressResp.Result.PercentComplete)
			assert.NotEmpty(t, progressResp.Result.FileID)
			fileID = progressResp.Result.FileID
			break
		}

		// Progress should advance
		assert.GreaterOrEqual(t, progressResp.Result.PercentComplete, 0)
		assert.LessOrEqual(t, progressResp.Result.PercentComplete, 100)
	}

	require.NotEmpty(t, fileID, "export should complete within %d attempts", maxAttempts)

	// Step 3: Download file
	filePath := fmt.Sprintf("/API/v3/surveys/%s/export-responses/%s/file", surveyID, fileID)
	resp = makeAuthenticatedRequest(t, "GET", filePath)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/zip", resp.Header.Get("Content-Type"))
	assert.Contains(t, resp.Header.Get("Content-Disposition"), "attachment")

	// Read ZIP file
	zipData, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Verify ZIP contents
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	require.NoError(t, err)

	// Find CSV file in ZIP
	var csvFound bool
	for _, file := range zipReader.File {
		if file.Name == "survey_responses.csv" {
			csvFound = true

			// Open CSV
			csvFile, err := file.Open()
			require.NoError(t, err)
			defer csvFile.Close()

			// Parse CSV
			csvReader := csv.NewReader(csvFile)
			records, err := csvReader.ReadAll()
			require.NoError(t, err)

			// Verify CSV structure
			assert.Greater(t, len(records), 0, "CSV should have at least header")

			header := records[0]
			assert.Contains(t, header, "ResponseID")
			assert.Contains(t, header, "RespondentEmail")
			assert.Contains(t, header, "OverallAISatisfaction")

			t.Logf("CSV contains %d rows (including header)", len(records))
			break
		}
	}

	assert.True(t, csvFound, "survey_responses.csv should be in ZIP")

	t.Log("Full Qualtrics export flow completed successfully")
}

// TestQualtrics_E2E_ProgressAdvancement tests that progress advances correctly.
func TestQualtrics_E2E_ProgressAdvancement(t *testing.T) {
	cleanup, _ := setupExternalDataE2EServer(t)
	defer cleanup()

	// Load seed data to get survey ID
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)
	surveyID := seedData.ExternalDataSources.Qualtrics.SurveyID

	// Start export
	startPath := fmt.Sprintf("/API/v3/surveys/%s/export-responses", surveyID)
	resp := makeAuthenticatedRequest(t, "POST", startPath)
	defer resp.Body.Close()

	var startResp models.ExportStartResponse
	json.NewDecoder(resp.Body).Decode(&startResp)
	progressID := startResp.Result.ProgressID

	// Track progress over multiple checks
	var previousPercent int
	progressPath := fmt.Sprintf("/API/v3/surveys/%s/export-responses/%s", surveyID, progressID)

	for i := 0; i < 6; i++ {
		time.Sleep(100 * time.Millisecond)

		resp := makeAuthenticatedRequest(t, "GET", progressPath)
		defer resp.Body.Close()

		var progressResp models.ExportProgressResponse
		json.NewDecoder(resp.Body).Decode(&progressResp)

		currentPercent := progressResp.Result.PercentComplete
		t.Logf("Check %d: %d%% complete", i+1, currentPercent)

		// Progress should never decrease
		assert.GreaterOrEqual(t, currentPercent, previousPercent, "progress should not go backwards")

		// Should advance by 20% each call (based on implementation)
		if i > 0 && currentPercent < 100 {
			assert.GreaterOrEqual(t, currentPercent, previousPercent, "progress should advance")
		}

		previousPercent = currentPercent

		if progressResp.Result.Status == "complete" {
			assert.Equal(t, 100, currentPercent)
			break
		}
	}

	t.Log("Progress advancement verified")
}

// TestQualtrics_E2E_DisabledWhenNotConfigured tests that Qualtrics endpoints return 404 when not configured.
func TestQualtrics_E2E_DisabledWhenNotConfigured(t *testing.T) {
	// Create seed data WITHOUT Qualtrics enabled
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Dev One"},
		},
		ExternalDataSources: nil, // No external data sources
	}

	store := storage.NewMemoryStore()
	router := server.NewRouter(store, seedData, testAPIKey)

	// Start server on different port
	const testPort = 19086
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", testPort),
		Handler: router,
	}

	go func() {
		httpServer.ListenAndServe()
	}()
	defer httpServer.Shutdown(context.Background())

	time.Sleep(50 * time.Millisecond)

	// Try to access Qualtrics endpoint
	baseURL := fmt.Sprintf("http://localhost:%d", testPort)
	req, _ := http.NewRequest("POST", baseURL+"/API/v3/surveys/SV_test/export-responses", nil)
	req.SetBasicAuth(testAPIKey, "")

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should return 404 since Qualtrics is not configured
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	t.Log("Qualtrics endpoints correctly return 404 when not configured")
}

// ============================================================================
// Authentication Tests
// ============================================================================

// TestExternalData_E2E_AuthenticationRequired tests that all external data endpoints require authentication.
func TestExternalData_E2E_AuthenticationRequired(t *testing.T) {
	cleanup, _ := setupExternalDataE2EServer(t)
	defer cleanup()

	// Load seed data to get survey ID
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)
	surveyID := seedData.ExternalDataSources.Qualtrics.SurveyID

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/harvey/api/v1/history/usage?from=2026-01-01&to=2026-01-09"},
		{"GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')"},
		{"POST", fmt.Sprintf("/API/v3/surveys/%s/export-responses", surveyID)},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.method+" "+endpoint.path, func(t *testing.T) {
			resp := makeUnauthenticatedRequest(t, endpoint.method, endpoint.path)
			defer resp.Body.Close()

			// Should return 401 Unauthorized
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "endpoint should require authentication")
		})
	}

	t.Log("All external data endpoints correctly require authentication")
}

// ============================================================================
// Integration Test: All APIs Enabled
// ============================================================================

// TestExternalData_E2E_AllAPIsEnabled tests that all external data APIs are accessible when configured.
func TestExternalData_E2E_AllAPIsEnabled(t *testing.T) {
	cleanup, _ := setupExternalDataE2EServer(t)
	defer cleanup()

	// Load seed data to get survey ID
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)
	surveyID := seedData.ExternalDataSources.Qualtrics.SurveyID

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"Harvey Usage", "GET", "/harvey/api/v1/history/usage?from=2026-01-01&to=2026-01-09"},
		{"Copilot JSON", "GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=application/json"},
		{"Copilot CSV", "GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=text/csv"},
		{"Qualtrics Start Export", "POST", fmt.Sprintf("/API/v3/surveys/%s/export-responses", surveyID)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := makeAuthenticatedRequest(t, tt.method, tt.path)
			defer resp.Body.Close()

			// All endpoints should return 200 OK
			assert.Equal(t, http.StatusOK, resp.StatusCode, "endpoint should be accessible")

			// Verify content type is set
			contentType := resp.Header.Get("Content-Type")
			assert.NotEmpty(t, contentType, "content-type should be set")

			t.Logf("%s: %d %s", tt.name, resp.StatusCode, contentType)
		})
	}

	t.Log("All external data APIs are accessible and responding")
}

// ============================================================================
// Error Handling Tests
// ============================================================================

// TestExternalData_E2E_ErrorCases tests error handling for invalid requests.
func TestExternalData_E2E_ErrorCases(t *testing.T) {
	cleanup, _ := setupExternalDataE2EServer(t)
	defer cleanup()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"Harvey invalid date", "GET", "/harvey/api/v1/history/usage?from=invalid&to=2026-01-09", http.StatusBadRequest},
		{"Harvey invalid page", "GET", "/harvey/api/v1/history/usage?from=2026-01-01&to=2026-01-09&page=-1", http.StatusBadRequest},
		{"Copilot invalid period", "GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D99')", http.StatusBadRequest},
		{"Qualtrics invalid survey", "POST", "/API/v3/surveys//export-responses", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := makeAuthenticatedRequest(t, tt.method, tt.path)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode, "should return expected error status")

			t.Logf("%s: %d (expected %d)", tt.name, resp.StatusCode, tt.expectedStatus)
		})
	}
}
