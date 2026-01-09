package server

import (
	"net/http/httptest"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
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
