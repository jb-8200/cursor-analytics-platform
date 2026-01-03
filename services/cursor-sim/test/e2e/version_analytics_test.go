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

// TestE2E_TeamClientVersions verifies that /analytics/team/client-versions returns version data.
func TestE2E_TeamClientVersions(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate client version events
	ctx := context.Background()
	versionGen := generator.NewVersionGenerator(seedData, store, "medium")
	err = versionGen.GenerateClientVersions(ctx, 7) // 1 week of data
	require.NoError(t, err, "Failed to generate client versions")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/team/client-versions", nil)
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

	// Verify versions are realistic (should look like semantic versions)
	versions := extractVersionNumbers(response.Data)
	assert.NotEmpty(t, versions, "Expected at least one version in response")

	// Versions should follow pattern like "0.41.0", "0.42.1", etc.
	for _, version := range versions {
		assert.Regexp(t, `^\d+\.\d+\.\d+$`, version, "Version should be in semantic version format")
	}
}

// TestE2E_ByUserClientVersions verifies that /analytics/by-user/client-versions returns per-user version data.
func TestE2E_ByUserClientVersions(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate client version events
	ctx := context.Background()
	versionGen := generator.NewVersionGenerator(seedData, store, "medium")
	err = versionGen.GenerateClientVersions(ctx, 7) // 1 week of data
	require.NoError(t, err, "Failed to generate client versions")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/by-user/client-versions", nil)
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

// TestE2E_ClientVersionsEmpty_WithoutGenerator documents bug: returns empty without version generation.
func TestE2E_ClientVersionsEmpty_WithoutGenerator(t *testing.T) {
	// This test shows the bug: without VersionGenerator, endpoints return empty

	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// NOTE: Intentionally NOT calling VersionGenerator to show the bug

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Test team client-versions endpoint
	req := httptest.NewRequest("GET", "/analytics/team/client-versions", nil)
	req.SetBasicAuth("test-api-key", "")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected 200 OK")

	var response TeamAnalyticsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")

	// BUG: Returns empty array without version generation
	assert.Empty(t, response.Data, "BUG: Without VersionGenerator, returns empty data")
}

// extractVersionNumbers extracts version strings from analytics data.
// Response format may vary, so we check common fields: "version", "client_version", etc.
func extractVersionNumbers(data []map[string]interface{}) []string {
	versions := make([]string, 0)
	seen := make(map[string]bool)

	for _, item := range data {
		// Try different possible field names for version data
		var versionStr string

		// Check for "version" field
		if v, ok := item["version"].(string); ok {
			versionStr = v
		} else if v, ok := item["client_version"].(string); ok {
			versionStr = v
		} else if breakdown, ok := item["version_breakdown"].(map[string]interface{}); ok {
			// If there's a breakdown, extract version names
			for ver := range breakdown {
				if !seen[ver] {
					versions = append(versions, ver)
					seen[ver] = true
				}
			}
			continue
		}

		if versionStr != "" && !seen[versionStr] {
			versions = append(versions, versionStr)
			seen[versionStr] = true
		}
	}

	return versions
}
