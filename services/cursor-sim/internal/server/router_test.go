package server

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

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

func TestNewRouter(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	apiKey := "test-key"

	router := NewRouter(store, seedData, apiKey)

	assert.NotNil(t, router)
}

func TestRouter_HealthEndpoint(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	router := NewRouter(store, seedData, "test-key")

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestRouter_TeamsMembers(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	router := NewRouter(store, seedData, "test-key")

	req := httptest.NewRequest("GET", "/teams/members", nil)
	req.SetBasicAuth("test-key", "")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestRouter_AICodeCommits(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	router := NewRouter(store, seedData, "test-key")

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits", nil)
	req.SetBasicAuth("test-key", "")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestRouter_WithoutAuth(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()
	router := NewRouter(store, seedData, "test-key")

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
	router := NewRouter(store, seedData, "test-key")

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

	router := NewRouter(store, seedData, "test-key")

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

	router := NewRouter(store, seedData, "test-key")

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

	router := NewRouter(store, seedData, "test-key")

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

	router := NewRouter(store, seedData, "test-key")

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
	router := NewRouter(store, seedData, apiKey)

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
	router := NewRouter(store, seedData, apiKey)

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
	router := NewRouter(store, seedData, apiKey)

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
	router := NewRouter(store, seedData, apiKey)

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
	router := NewRouter(store, seedData, apiKey)

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
