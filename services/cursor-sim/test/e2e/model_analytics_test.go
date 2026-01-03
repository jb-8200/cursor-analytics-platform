package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/server"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TeamAnalyticsResponse represents the generic team analytics response format.
type TeamAnalyticsResponse struct {
	Data   []map[string]interface{} `json:"data"`
	Params map[string]interface{}   `json:"params"`
}

// ByUserAnalyticsResponse represents the by-user analytics response format.
type ByUserAnalyticsResponse struct {
	Data       map[string]interface{} `json:"data"`
	Pagination map[string]interface{} `json:"pagination"`
	Params     map[string]interface{} `json:"params"`
}

// TestE2E_TeamModels verifies that /analytics/team/models returns model usage data.
func TestE2E_TeamModels(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate model usage events
	ctx := context.Background()
	modelGen := generator.NewModelGenerator(seedData, store, "medium")
	err = modelGen.GenerateModelUsage(ctx, 7) // 1 week of data
	require.NoError(t, err, "Failed to generate model usage")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/team/models", nil)
	req.SetBasicAuth("test-api-key", "")

	// Execute request
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code, "Expected 200 OK")

	// Parse response
	var response TeamAnalyticsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")

	// Verify data is present
	assert.NotEmpty(t, response.Data, "Expected non-empty data array")

	// Verify models match seed file preferences
	// Seed file has: ["gpt-4-turbo", "claude-3-sonnet"] for Alice
	// and ["gpt-4-turbo"] for Bob
	modelNames := extractModelNames(response.Data)

	// Debug: print actual models if assertion will fail
	if len(modelNames) == 0 {
		t.Logf("Response data: %+v", response.Data)
	}

	assert.NotEmpty(t, modelNames, "Expected at least one model in response")

	// Verify expected models are present
	// Seed file has: ["gpt-4-turbo", "claude-3-sonnet"] for Alice and ["gpt-4-turbo"] for Bob
	assert.Contains(t, modelNames, "gpt-4-turbo", "Expected gpt-4-turbo in models")
	assert.Contains(t, modelNames, "claude-3-sonnet", "Expected claude-3-sonnet in models")
}

// TestE2E_ByUserModels verifies that /analytics/by-user/models returns per-user model data.
func TestE2E_ByUserModels(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate model usage events
	ctx := context.Background()
	modelGen := generator.NewModelGenerator(seedData, store, "medium")
	err = modelGen.GenerateModelUsage(ctx, 7) // 1 week of data
	require.NoError(t, err, "Failed to generate model usage")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/by-user/models", nil)
	req.SetBasicAuth("test-api-key", "")

	// Execute request
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code, "Expected 200 OK")

	// Parse response
	var response ByUserAnalyticsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")

	// Verify data is present (should be map with user emails as keys)
	assert.NotEmpty(t, response.Data, "Expected non-empty data object")

	// Verify both users have data
	dataJSON, _ := json.Marshal(response.Data)
	dataStr := string(dataJSON)
	assert.Contains(t, dataStr, "alice@example.com", "Expected alice@example.com in response")
	assert.Contains(t, dataStr, "bob@example.com", "Expected bob@example.com in response")
}

// TestE2E_ModelsEmpty_WithoutGenerator documents bug: returns empty without model generation.
func TestE2E_ModelsEmpty_WithoutGenerator(t *testing.T) {
	// This test shows the bug: without ModelGenerator, endpoints return empty

	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// NOTE: Intentionally NOT calling ModelGenerator to show the bug

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Test team models endpoint
	req := httptest.NewRequest("GET", "/analytics/team/models", nil)
	req.SetBasicAuth("test-api-key", "")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected 200 OK")

	var response TeamAnalyticsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")

	// BUG: Returns empty array without model generation
	assert.Empty(t, response.Data, "BUG: Without ModelGenerator, returns empty data")
}

// extractModelNames extracts model names from analytics data.
// The response format is: [{date: "...", model_breakdown: {model1: {...}, model2: {...}}}]
func extractModelNames(data []map[string]interface{}) []string {
	names := make([]string, 0)
	seen := make(map[string]bool)

	for _, item := range data {
		if breakdown, ok := item["model_breakdown"].(map[string]interface{}); ok {
			for modelName := range breakdown {
				if !seen[modelName] {
					names = append(names, modelName)
					seen[modelName] = true
				}
			}
		}
	}
	return names
}
