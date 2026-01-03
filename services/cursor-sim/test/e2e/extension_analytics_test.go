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

// TestE2E_TeamFileExtensions verifies that /analytics/team/top-file-extensions returns extension data.
func TestE2E_TeamFileExtensions(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate commits first (extensions are derived from commit file changes)
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 7) // 1 week of data
	require.NoError(t, err, "Failed to generate commits")

	// Generate file extension events
	extensionGen := generator.NewExtensionGenerator(seedData, store, "medium")
	err = extensionGen.GenerateFileExtensions(ctx, 7)
	require.NoError(t, err, "Failed to generate file extensions")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/team/top-file-extensions", nil)
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

	// Verify extensions are realistic (should match repo languages)
	extensions := extractExtensionNames(response.Data)
	assert.NotEmpty(t, extensions, "Expected at least one file extension in response")

	// Extensions should match repo languages (Go → go, TypeScript → ts/tsx, etc.)
	// Common extensions we expect: go, ts, tsx, js, json, md (without leading dot)
	validExtensions := []string{"go", "ts", "tsx", "js", "jsx", "json", "md", "yaml", "yml", "py", "java", "rb"}
	foundValid := false
	for _, ext := range extensions {
		for _, valid := range validExtensions {
			if ext == valid {
				foundValid = true
				break
			}
		}
		if foundValid {
			break
		}
	}
	assert.True(t, foundValid, "Expected at least one valid file extension")
}

// TestE2E_ByUserFileExtensions verifies that /analytics/by-user/top-file-extensions returns per-user extension data.
func TestE2E_ByUserFileExtensions(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate commits first
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 7)
	require.NoError(t, err, "Failed to generate commits")

	// Generate file extension events
	extensionGen := generator.NewExtensionGenerator(seedData, store, "medium")
	err = extensionGen.GenerateFileExtensions(ctx, 7)
	require.NoError(t, err, "Failed to generate file extensions")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/by-user/top-file-extensions", nil)
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

// TestE2E_FileExtensionsEmpty_WithoutGenerator documents bug: returns empty without extension generation.
func TestE2E_FileExtensionsEmpty_WithoutGenerator(t *testing.T) {
	// This test shows the bug: without ExtensionGenerator, endpoints return empty

	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// NOTE: Intentionally NOT calling ExtensionGenerator to show the bug

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Test team file-extensions endpoint
	req := httptest.NewRequest("GET", "/analytics/team/top-file-extensions", nil)
	req.SetBasicAuth("test-api-key", "")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected 200 OK")

	var response TeamAnalyticsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")

	// BUG: Returns empty array without extension generation
	assert.Empty(t, response.Data, "BUG: Without ExtensionGenerator, returns empty data")
}

// extractExtensionNames extracts file extension strings from analytics data.
// Response format: [{"event_date": "...", "file_extension": "go", ...}, ...]
func extractExtensionNames(data []map[string]interface{}) []string {
	extensions := make([]string, 0)
	seen := make(map[string]bool)

	for _, item := range data {
		// Look for "file_extension" field
		if ext, ok := item["file_extension"].(string); ok {
			if ext != "" && !seen[ext] {
				extensions = append(extensions, ext)
				seen[ext] = true
			}
		}
	}

	return extensions
}
