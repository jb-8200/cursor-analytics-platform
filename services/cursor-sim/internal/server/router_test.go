package server

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestSeedData() *seed.SeedData {
	return &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:    "user_001",
				Email:     "test@example.com",
				Name:      "Test Developer",
				Seniority: "mid",
				WorkingHoursBand: seed.WorkingHours{
					Start: 9,
					End:   18,
				},
			},
		},
		Repositories: []seed.Repository{
			{
				RepoName:      "test/repo",
				DefaultBranch: "main",
			},
		},
	}
}

func createTestConfig() *config.Config {
	return &config.Config{
		Mode:       "runtime",
		Port:       8080,
		Days:       1,
		Velocity:   "medium",
		SeedPath:   "",
		CorpusPath: "",
	}
}

const testVersion = "1.0.0"

func TestNewRouter(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	apiKey := "test-key"
	cfg := createTestConfig()

	router := NewRouter(store, seedData, apiKey, cfg, testVersion)

	assert.NotNil(t, router)
}

func TestRouter_HealthEndpoint(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	router := NewRouter(store, seedData, "test-key", createTestConfig(), testVersion)

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestRouter_TeamsMembers(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	router := NewRouter(store, seedData, "test-key", createTestConfig(), testVersion)

	req := httptest.NewRequest("GET", "/teams/members", nil)
	req.SetBasicAuth("test-key", "")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestRouter_AICodeCommits(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	router := NewRouter(store, seedData, "test-key", createTestConfig(), testVersion)

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits", nil)
	req.SetBasicAuth("test-key", "")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestRouter_WithoutAuth(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	router := NewRouter(store, seedData, "test-key", createTestConfig(), testVersion)

	req := httptest.NewRequest("GET", "/teams/members", nil)
	// No auth header
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should return 401 Unauthorized
	assert.Equal(t, 401, rec.Code)
}

func TestRouter_NotFound(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	router := NewRouter(store, seedData, "test-key", createTestConfig(), testVersion)

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	req.SetBasicAuth("test-key", "")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should return 404 Not Found
	assert.Equal(t, 404, rec.Code)
}

// TestRouter_HarveyRoutes_Enabled verifies Harvey routes are registered when seed data contains Harvey configuration.
func TestRouter_HarveyRoutes_Enabled(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "test@example.com",
				Name:   "Test Developer",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Harvey: &seed.HarveySeedConfig{
				Enabled:    true,
				TotalUsage: seed.UsageRange{Min: 100, Max: 500},
				ModelsUsed: []string{"gpt-4"},
			},
		},
	}

	router := NewRouter(store, seedData, "test-key", createTestConfig(), testVersion)

	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=2025-01-01&to=2025-01-31", nil)
	req.SetBasicAuth("test-key", "")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Harvey routes should be registered and return 200 (not 404)
	assert.NotEqual(t, 404, rec.Code, "Harvey route should be registered when Harvey config exists")
}

// TestRouter_HarveyRoutes_DisabledNilConfig verifies Harvey routes are NOT registered when ExternalDataSources is nil.
func TestRouter_HarveyRoutes_DisabledNilConfig(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "test@example.com",
			},
		},
		ExternalDataSources: nil, // No external data sources
	}

	router := NewRouter(store, seedData, "test-key", createTestConfig(), testVersion)

	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=2025-01-01&to=2025-01-31", nil)
	req.SetBasicAuth("test-key", "")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Harvey routes should NOT be registered, should return 404
	assert.Equal(t, 404, rec.Code, "Harvey route should not be registered when ExternalDataSources is nil")
}

// TestRouter_HarveyRoutes_DisabledNilHarvey verifies Harvey routes are NOT registered when Harvey config is nil.
func TestRouter_HarveyRoutes_DisabledNilHarvey(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "test@example.com",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Harvey: nil, // Harvey config is nil
		},
	}

	router := NewRouter(store, seedData, "test-key", createTestConfig(), testVersion)

	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=2025-01-01&to=2025-01-31", nil)
	req.SetBasicAuth("test-key", "")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Harvey routes should NOT be registered, should return 404
	assert.Equal(t, 404, rec.Code, "Harvey route should not be registered when Harvey config is nil")
}

// TestRouter_HarveyRoutes_DisabledFalse verifies Harvey routes are NOT registered when Harvey.Enabled is false.
func TestRouter_HarveyRoutes_DisabledFalse(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "test@example.com",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Harvey: &seed.HarveySeedConfig{
				Enabled: false, // Harvey is disabled
			},
		},
	}

	router := NewRouter(store, seedData, "test-key", createTestConfig(), testVersion)

	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=2025-01-01&to=2025-01-31", nil)
	req.SetBasicAuth("test-key", "")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Harvey routes should NOT be registered, should return 404
	assert.Equal(t, 404, rec.Code, "Harvey route should not be registered when Harvey.Enabled is false")
}

// TestRouter_CopilotRoutes_Enabled verifies Copilot routes are registered when seed data contains Copilot configuration.
func TestRouter_CopilotRoutes_Enabled(t *testing.T) {
	// Arrange: Create seed data WITH Copilot enabled
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID: "dev1",
				Email:  "dev1@example.com",
				Name:   "Developer One",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled:       true,
				TotalLicenses: 10,
				ActiveUsers:   8,
			},
		},
	}

	store := storage.NewMemoryStore()
	apiKey := "test-api-key"
	router := NewRouter(store, seedData, apiKey, createTestConfig(), testVersion)

	// Act: Request the Copilot endpoint
	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')", nil)
	req.SetBasicAuth(apiKey, "")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert: Should return 200 OK (route exists and handler works)
	assert.Equal(t, 200, rec.Code, "Expected Copilot route to exist when enabled")

	// Verify JSON response structure
	var response map[string]interface{}
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err, "Response should be valid JSON")
	assert.Contains(t, response, "@odata.context", "Response should have OData context")
	assert.Contains(t, response, "value", "Response should have value array")
}

// TestRouter_CopilotRoutes_Disabled verifies Copilot routes are NOT registered when ExternalDataSources is nil.
func TestRouter_CopilotRoutes_Disabled(t *testing.T) {
	// Arrange: Create seed data WITHOUT Copilot enabled
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID: "dev1",
				Email:  "dev1@example.com",
				Name:   "Developer One",
			},
		},
		ExternalDataSources: nil, // No external data sources
	}

	store := storage.NewMemoryStore()
	apiKey := "test-api-key"
	router := NewRouter(store, seedData, apiKey, createTestConfig(), testVersion)

	// Act: Request the Copilot endpoint
	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')", nil)
	req.SetBasicAuth(apiKey, "")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert: Should return 404 (route does not exist)
	assert.Equal(t, 404, rec.Code, "Expected 404 when Copilot is not enabled")
}

// TestRouter_CopilotRoutes_EnabledButFalse verifies Copilot routes are NOT registered when Copilot.Enabled is false.
func TestRouter_CopilotRoutes_EnabledButFalse(t *testing.T) {
	// Arrange: Create seed data WITH Copilot config but Enabled=false
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID: "dev1",
				Email:  "dev1@example.com",
				Name:   "Developer One",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled: false, // Explicitly disabled
			},
		},
	}

	store := storage.NewMemoryStore()
	apiKey := "test-api-key"
	router := NewRouter(store, seedData, apiKey, createTestConfig(), testVersion)

	// Act: Request the Copilot endpoint
	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')", nil)
	req.SetBasicAuth(apiKey, "")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert: Should return 404 (route does not exist)
	assert.Equal(t, 404, rec.Code, "Expected 404 when Copilot is disabled")
}

// TestRouter_CopilotRoutes_RequiresAuth verifies Copilot routes require authentication.
func TestRouter_CopilotRoutes_RequiresAuth(t *testing.T) {
	// Arrange: Create seed data WITH Copilot enabled
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID: "dev1",
				Email:  "dev1@example.com",
				Name:   "Developer One",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled:       true,
				TotalLicenses: 10,
			},
		},
	}

	store := storage.NewMemoryStore()
	apiKey := "test-api-key"
	router := NewRouter(store, seedData, apiKey, createTestConfig(), testVersion)

	// Act: Request WITHOUT authentication
	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert: Should return 401 Unauthorized
	assert.Equal(t, 401, rec.Code, "Expected 401 when authentication is missing")
}

// TestRouter_CopilotRoutes_CSVFormat verifies Copilot CSV export works.
func TestRouter_CopilotRoutes_CSVFormat(t *testing.T) {
	// Arrange: Create seed data WITH Copilot enabled
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID: "dev1",
				Email:  "dev1@example.com",
				Name:   "Developer One",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled:       true,
				TotalLicenses: 10,
				ActiveUsers:   8,
			},
		},
	}

	store := storage.NewMemoryStore()
	apiKey := "test-api-key"
	router := NewRouter(store, seedData, apiKey, createTestConfig(), testVersion)

	// Act: Request CSV format
	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=text/csv", nil)
	req.SetBasicAuth(apiKey, "")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert: Should return CSV
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "text/csv", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, rec.Header().Get("Content-Disposition"), "copilot-usage-D30.csv")
}

// TestRouter_QualtricsRoutes_Enabled verifies Qualtrics routes are registered when seed data contains Qualtrics configuration.
func TestRouter_QualtricsRoutes_Enabled(t *testing.T) {
	// Arrange: Create seed data WITH Qualtrics enabled
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID: "dev1",
				Email:  "dev1@example.com",
				Name:   "Developer One",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				SurveyName:    "Test Survey",
				ResponseCount: 50,
			},
		},
	}

	store := storage.NewMemoryStore()
	apiKey := "test-api-key"
	router := NewRouter(store, seedData, apiKey, createTestConfig(), testVersion)

	// Act: Request the Qualtrics start export endpoint
	req := httptest.NewRequest("POST", "/API/v3/surveys/SV_abc123/export-responses", nil)
	req.SetBasicAuth(apiKey, "")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert: Should return 200 OK (route exists and handler works)
	assert.Equal(t, 200, rec.Code, "Expected Qualtrics route to exist when enabled")

	// Verify JSON response structure
	var response map[string]interface{}
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err, "Response should be valid JSON")
	assert.Contains(t, response, "result", "Response should have result object")
	assert.Contains(t, response, "meta", "Response should have meta object")
}

// TestRouter_QualtricsRoutes_Disabled verifies Qualtrics routes are NOT registered when ExternalDataSources is nil.
func TestRouter_QualtricsRoutes_Disabled(t *testing.T) {
	// Arrange: Create seed data WITHOUT Qualtrics enabled
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID: "dev1",
				Email:  "dev1@example.com",
				Name:   "Developer One",
			},
		},
		ExternalDataSources: nil, // No external data sources
	}

	store := storage.NewMemoryStore()
	apiKey := "test-api-key"
	router := NewRouter(store, seedData, apiKey, createTestConfig(), testVersion)

	// Act: Request the Qualtrics endpoint
	req := httptest.NewRequest("POST", "/API/v3/surveys/SV_abc123/export-responses", nil)
	req.SetBasicAuth(apiKey, "")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert: Should return 404 (route does not exist)
	assert.Equal(t, 404, rec.Code, "Expected 404 when Qualtrics is not enabled")
}

// TestRouter_QualtricsRoutes_EnabledButFalse verifies Qualtrics routes are NOT registered when Qualtrics.Enabled is false.
func TestRouter_QualtricsRoutes_EnabledButFalse(t *testing.T) {
	// Arrange: Create seed data WITH Qualtrics config but Enabled=false
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID: "dev1",
				Email:  "dev1@example.com",
				Name:   "Developer One",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled: false, // Explicitly disabled
			},
		},
	}

	store := storage.NewMemoryStore()
	apiKey := "test-api-key"
	router := NewRouter(store, seedData, apiKey, createTestConfig(), testVersion)

	// Act: Request the Qualtrics endpoint
	req := httptest.NewRequest("POST", "/API/v3/surveys/SV_abc123/export-responses", nil)
	req.SetBasicAuth(apiKey, "")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert: Should return 404 (route does not exist)
	assert.Equal(t, 404, rec.Code, "Expected 404 when Qualtrics is disabled")
}

// TestRouter_QualtricsRoutes_RequiresAuth verifies Qualtrics routes require authentication.
func TestRouter_QualtricsRoutes_RequiresAuth(t *testing.T) {
	// Arrange: Create seed data WITH Qualtrics enabled
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID: "dev1",
				Email:  "dev1@example.com",
				Name:   "Developer One",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				ResponseCount: 50,
			},
		},
	}

	store := storage.NewMemoryStore()
	apiKey := "test-api-key"
	router := NewRouter(store, seedData, apiKey, createTestConfig(), testVersion)

	// Act: Request WITHOUT authentication
	req := httptest.NewRequest("POST", "/API/v3/surveys/SV_abc123/export-responses", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert: Should return 401 Unauthorized
	assert.Equal(t, 401, rec.Code, "Expected 401 when authentication is missing")
}

// TestRouter_QualtricsRoutes_ProgressEndpoint verifies the progress endpoint works.
func TestRouter_QualtricsRoutes_ProgressEndpoint(t *testing.T) {
	// Arrange: Create seed data WITH Qualtrics enabled
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID: "dev1",
				Email:  "dev1@example.com",
				Name:   "Developer One",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				ResponseCount: 50,
			},
		},
	}

	store := storage.NewMemoryStore()
	apiKey := "test-api-key"
	router := NewRouter(store, seedData, apiKey, createTestConfig(), testVersion)

	// First, start an export to get a progressID
	startReq := httptest.NewRequest("POST", "/API/v3/surveys/SV_abc123/export-responses", nil)
	startReq.SetBasicAuth(apiKey, "")
	startRec := httptest.NewRecorder()
	router.ServeHTTP(startRec, startReq)

	require.Equal(t, 200, startRec.Code, "Start export should succeed")

	var startResp map[string]interface{}
	err := json.NewDecoder(startRec.Body).Decode(&startResp)
	require.NoError(t, err)

	result := startResp["result"].(map[string]interface{})
	progressID := result["progressId"].(string)

	// Act: Check progress
	progressReq := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/"+progressID, nil)
	progressReq.SetBasicAuth(apiKey, "")
	progressRec := httptest.NewRecorder()
	router.ServeHTTP(progressRec, progressReq)

	// Assert: Should return 200 OK
	assert.Equal(t, 200, progressRec.Code, "Progress endpoint should work")

	var progressResp map[string]interface{}
	err = json.NewDecoder(progressRec.Body).Decode(&progressResp)
	require.NoError(t, err, "Response should be valid JSON")
	assert.Contains(t, progressResp, "result", "Response should have result object")
}

// TestRouter_DocsLoginPage verifies /docs/login endpoint is accessible without auth.
func TestRouter_DocsLoginPage(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	apiKey := "test-key"
	cfg := createTestConfig()
	router := NewRouter(store, seedData, apiKey, cfg, testVersion)

	req := httptest.NewRequest("GET", "/docs/login", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should return 200 and HTML content
	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
}

// TestRouter_DocsWithoutSession verifies /docs redirects to login without session.
func TestRouter_DocsWithoutSession(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	apiKey := "test-key"
	cfg := createTestConfig()
	router := NewRouter(store, seedData, apiKey, cfg, testVersion)

	req := httptest.NewRequest("GET", "/docs", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should redirect to login
	assert.Equal(t, 302, rec.Code)
	assert.Contains(t, rec.Header().Get("Location"), "/docs/login")
}

// TestRouter_DocsOpenAPISpec verifies /docs/openapi/ endpoint serves YAML specs.
func TestRouter_DocsOpenAPISpec(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	apiKey := "test-key"
	cfg := createTestConfig()
	router := NewRouter(store, seedData, apiKey, cfg, testVersion)

	req := httptest.NewRequest("GET", "/docs/openapi/cursor-api.yaml", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should return 200 and YAML content type
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "application/yaml", rec.Header().Get("Content-Type"))
	assert.NotEmpty(t, rec.Body.String())
}

// TestRouter_QualtricsRoutes_FileDownloadEndpoint verifies the file download endpoint works.
func TestRouter_QualtricsRoutes_FileDownloadEndpoint(t *testing.T) {
	// Arrange: Create seed data WITH Qualtrics enabled
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID: "dev1",
				Email:  "dev1@example.com",
				Name:   "Developer One",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				ResponseCount: 50,
			},
		},
	}

	store := storage.NewMemoryStore()
	apiKey := "test-api-key"
	router := NewRouter(store, seedData, apiKey, createTestConfig(), testVersion)

	// First, start an export
	startReq := httptest.NewRequest("POST", "/API/v3/surveys/SV_abc123/export-responses", nil)
	startReq.SetBasicAuth(apiKey, "")
	startRec := httptest.NewRecorder()
	router.ServeHTTP(startRec, startReq)

	var startResp map[string]interface{}
	json.NewDecoder(startRec.Body).Decode(&startResp)
	result := startResp["result"].(map[string]interface{})
	progressID := result["progressId"].(string)

	// Poll progress until complete
	var fileID string
	for i := 0; i < 10; i++ {
		progressReq := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/"+progressID, nil)
		progressReq.SetBasicAuth(apiKey, "")
		progressRec := httptest.NewRecorder()
		router.ServeHTTP(progressRec, progressReq)

		var progressResp map[string]interface{}
		json.NewDecoder(progressRec.Body).Decode(&progressResp)
		progressResult := progressResp["result"].(map[string]interface{})

		if progressResult["status"] == "complete" {
			fileID = progressResult["fileId"].(string)
			break
		}
	}

	require.NotEmpty(t, fileID, "Should have a file ID after completion")

	// Act: Download the file
	fileReq := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/"+fileID+"/file", nil)
	fileReq.SetBasicAuth(apiKey, "")
	fileRec := httptest.NewRecorder()
	router.ServeHTTP(fileRec, fileReq)

	// Assert: Should return 200 OK with ZIP content
	assert.Equal(t, 200, fileRec.Code, "File download should work")
	assert.Equal(t, "application/zip", fileRec.Header().Get("Content-Type"), "Should return ZIP file")
	assert.Contains(t, fileRec.Header().Get("Content-Disposition"), "attachment", "Should be an attachment")
}
