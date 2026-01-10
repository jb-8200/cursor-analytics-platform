package cursor

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	// Create test config
	cfg := &config.Config{
		Days:     90,
		Velocity: "medium",
		Port:     8080,
		GenParams: config.GenerationParams{
			Developers: 2,
			Days:       90,
			MaxCommits: 1000,
		},
	}

	// Create test seed data
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID:    "user_001",
				Name:      "Alice Chen",
				Email:     "alice@example.com",
				Org:       "acme-corp",
				Division:  "Engineering",
				Team:      "Backend",
				Region:    "US",
				Seniority: "senior",
			},
			{
				UserID:    "user_002",
				Name:      "Bob Smith",
				Email:     "bob@example.com",
				Org:       "acme-corp",
				Division:  "Engineering",
				Team:      "Frontend",
				Region:    "EU",
				Seniority: "mid",
			},
		},
		Repositories: []seed.Repository{
			{RepoName: "acme-corp/payment-service"},
			{RepoName: "acme-corp/web-app"},
		},
	}

	// Create handler
	handler := GetConfig(cfg, seedData, "2.0.0")

	// Make request
	req := httptest.NewRequest("GET", "/admin/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Check status code
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	// Parse response
	var response models.ConfigResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify generation config
	assert.Equal(t, 90, response.Generation.Days)
	assert.Equal(t, "medium", response.Generation.Velocity)
	assert.Equal(t, 2, response.Generation.Developers)
	assert.Equal(t, 1000, response.Generation.MaxCommits)

	// Verify seed config
	assert.Equal(t, "1.0", response.Seed.Version)
	assert.Equal(t, 2, response.Seed.Developers)
	assert.Equal(t, 2, response.Seed.Repositories)
	assert.ElementsMatch(t, []string{"acme-corp"}, response.Seed.Organizations)
	assert.ElementsMatch(t, []string{"Engineering"}, response.Seed.Divisions)
	assert.ElementsMatch(t, []string{"Backend", "Frontend"}, response.Seed.Teams)
	assert.ElementsMatch(t, []string{"US", "EU"}, response.Seed.Regions)

	// Verify developer breakdowns
	assert.Equal(t, 1, response.Seed.BySeniority["senior"])
	assert.Equal(t, 1, response.Seed.BySeniority["mid"])
	assert.Equal(t, 1, response.Seed.ByRegion["US"])
	assert.Equal(t, 1, response.Seed.ByRegion["EU"])
	assert.Equal(t, 1, response.Seed.ByTeam["Backend"])
	assert.Equal(t, 1, response.Seed.ByTeam["Frontend"])

	// Verify external sources (disabled by default)
	assert.False(t, response.ExternalSources.Harvey.Enabled)
	assert.False(t, response.ExternalSources.Copilot.Enabled)
	assert.False(t, response.ExternalSources.Qualtrics.Enabled)

	// Verify server config
	assert.Equal(t, 8080, response.Server.Port)
	assert.Equal(t, "2.0.0", response.Server.Version)
	assert.NotEmpty(t, response.Server.Uptime)
}

func TestGetConfig_ExternalSources(t *testing.T) {
	// Create test config
	cfg := &config.Config{
		Days:     180,
		Velocity: "high",
		Port:     8080,
		GenParams: config.GenerationParams{
			Developers: 100,
			Days:       180,
			MaxCommits: 500,
		},
	}

	// Create test seed data with external sources
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID:    "user_001",
				Name:      "Developer One",
				Email:     "dev1@example.com",
				Org:       "example-corp",
				Division:  "Engineering",
				Team:      "Platform",
				Region:    "US",
				Seniority: "senior",
			},
		},
		Repositories: []seed.Repository{
			{RepoName: "example-corp/api"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Harvey: &seed.HarveySeedConfig{
				Enabled:       true,
				ModelsUsed:    []string{"gpt-4", "claude-3-sonnet"},
				PracticeAreas: []string{"Contract Review", "Legal Research"},
			},
			Copilot: &seed.CopilotSeedConfig{
				Enabled:       true,
				TotalLicenses: 50,
				ActiveUsers:   35,
			},
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_aitools_q1_2026",
				SurveyName:    "AI Tools Survey Q1 2026",
				ResponseCount: 150,
			},
		},
	}

	// Create handler
	handler := GetConfig(cfg, seedData, "2.0.0")

	// Make request
	req := httptest.NewRequest("GET", "/admin/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Check status code
	assert.Equal(t, 200, rec.Code)

	// Parse response
	var response models.ConfigResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify Harvey config
	assert.True(t, response.ExternalSources.Harvey.Enabled)
	assert.ElementsMatch(t, []string{"gpt-4", "claude-3-sonnet"}, response.ExternalSources.Harvey.Models)

	// Verify Copilot config
	assert.True(t, response.ExternalSources.Copilot.Enabled)
	assert.Equal(t, 50, response.ExternalSources.Copilot.TotalLicenses)
	assert.Equal(t, 35, response.ExternalSources.Copilot.ActiveUsers)

	// Verify Qualtrics config
	assert.True(t, response.ExternalSources.Qualtrics.Enabled)
	assert.Equal(t, "SV_aitools_q1_2026", response.ExternalSources.Qualtrics.SurveyID)
	assert.Equal(t, 150, response.ExternalSources.Qualtrics.ResponseCount)
}

func TestGetConfig_MethodNotAllowed(t *testing.T) {
	// Create minimal test config
	cfg := &config.Config{
		Days:     90,
		Velocity: "medium",
		Port:     8080,
		GenParams: config.GenerationParams{
			Developers: 1,
			Days:       90,
			MaxCommits: 0,
		},
	}

	seedData := &seed.SeedData{
		Version:      "1.0",
		Developers:   []seed.Developer{{UserID: "test"}},
		Repositories: []seed.Repository{{RepoName: "test/repo"}},
	}

	handler := GetConfig(cfg, seedData, "2.0.0")

	// Test POST method (should fail)
	req := httptest.NewRequest("POST", "/admin/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 405, rec.Code)
}
