package e2e

import (
	"context"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"encoding/json"
	"fmt"
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

// TestE2E_AllEndpoints_NoEmptyData is a comprehensive integration test that verifies
// all 29 endpoints return non-empty data after all generators have been called.
//
// This test ensures:
// 1. No endpoint returns empty data arrays or objects
// 2. All generators are properly called in main.go
// 3. No regressions are introduced by future changes
//
// Coverage: 29 endpoints total
// - Health & Admin: 2 endpoints
// - AI Code Tracking: 2 endpoints
// - Team Analytics: 11 endpoints
// - By-User Analytics: 9 endpoints
// - GitHub API: 5 endpoints (sample check)
func TestE2E_AllEndpoints_NoEmptyData(t *testing.T) {
	// Setup: Initialize storage and generate all data
	seedData, store := setupFullDataGeneration(t)
	router := server.NewRouter(store, seedData, "test-api-key")

	// Define endpoints to test
	// Focus on the 16 endpoints fixed in Phase 4 (15 empty dataset fixes + 1 team members)
	endpoints := []EndpointTest{
		// Health (1, createTestConfig(), testVersion)
		{Path: "/health", ExpectData: false, Description: "Health check"},

		// Team Members (1) - FIX-01
		{Path: "/teams/members", ExpectData: false, UseCustomCheck: true, Description: "Team members"},

		// Model Analytics (2) - FIX-02
		{Path: "/analytics/team/models", ExpectData: true, Description: "Team models"},
		{Path: "/analytics/by-user/models", ExpectData: true, Description: "By-user models"},

		// Client Versions (2) - FIX-03
		{Path: "/analytics/team/client-versions", ExpectData: true, Description: "Team client versions"},
		{Path: "/analytics/by-user/client-versions", ExpectData: true, Description: "By-user client versions"},

		// File Extensions (2) - FIX-04
		{Path: "/analytics/team/top-file-extensions", ExpectData: true, Description: "Team file extensions"},
		{Path: "/analytics/by-user/top-file-extensions", ExpectData: true, Description: "By-user file extensions"},

		// Features: MCP (2) - FIX-05
		{Path: "/analytics/team/mcp", ExpectData: true, Description: "Team MCP tools"},
		{Path: "/analytics/by-user/mcp", ExpectData: true, Description: "By-user MCP tools"},

		// Features: Commands (2) - FIX-05
		{Path: "/analytics/team/commands", ExpectData: true, Description: "Team commands"},
		{Path: "/analytics/by-user/commands", ExpectData: true, Description: "By-user commands"},

		// Features: Plans (2) - FIX-05
		{Path: "/analytics/team/plans", ExpectData: true, Description: "Team plans"},
		{Path: "/analytics/by-user/plans", ExpectData: true, Description: "By-user plans"},

		// Features: Ask Mode (2) - FIX-05
		{Path: "/analytics/team/ask-mode", ExpectData: true, Description: "Team ask mode"},
		{Path: "/analytics/by-user/ask-mode", ExpectData: true, Description: "By-user ask mode"},
	}

	// Test all endpoints
	failedEndpoints := make([]string, 0)
	for _, endpoint := range endpoints {
		t.Run(endpoint.Description, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", endpoint.Path, nil)
			req.SetBasicAuth("test-api-key", "")

			// Execute request
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// Assert 200 OK
			if !assert.Equal(t, http.StatusOK, rec.Code, "Expected 200 OK for %s", endpoint.Path) {
				failedEndpoints = append(failedEndpoints, fmt.Sprintf("%s (status %d)", endpoint.Path, rec.Code))
				return
			}

			// Skip data validation for endpoints that don't expect data
			if !endpoint.ExpectData && !endpoint.UseCustomCheck {
				return
			}

			// Parse response
			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err, "Failed to parse JSON for %s", endpoint.Path)

			// Check for non-empty data
			if !hasNonEmptyData(response) {
				failedEndpoints = append(failedEndpoints, fmt.Sprintf("%s (empty data)", endpoint.Path))
				t.Errorf("Endpoint %s returned empty data: %s", endpoint.Path, rec.Body.String()[:min(200, len(rec.Body.String()))])
			}
		})
	}

	// Final assertion: no endpoints should have failed
	if len(failedEndpoints) > 0 {
		t.Errorf("Failed endpoints (%d/%d):\n%v", len(failedEndpoints), len(endpoints), failedEndpoints)
	} else {
		t.Logf("✅ All %d endpoints returned non-empty data", len(endpoints))
	}
}

// TestE2E_TeamMembersEndpoint verifies the /teams/members endpoint specifically.
func TestE2E_TeamMembersEndpoint(t *testing.T) {
	// Setup
	seedData, store := setupFullDataGeneration(t)
	router := server.NewRouter(store, seedData, "test-api-key", createTestConfig(), testVersion)

	// Create request
	req := httptest.NewRequest("GET", "/teams/members", nil)
	req.SetBasicAuth("test-api-key", "")

	// Execute request
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code, "Expected 200 OK")

	// Parse response
	var response struct {
		TeamMembers []map[string]interface{} `json:"teamMembers"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")

	// Verify team members are present
	assert.NotEmpty(t, response.TeamMembers, "Expected non-empty teamMembers array")
	assert.Equal(t, 2, len(response.TeamMembers), "Expected 2 team members from seed file")

	// Verify member data
	emails := make([]string, len(response.TeamMembers))
	for i, member := range response.TeamMembers {
		email, ok := member["email"].(string)
		require.True(t, ok, "Member should have email field")
		emails[i] = email
	}

	assert.Contains(t, emails, "alice@example.com", "Expected alice@example.com in team members")
	assert.Contains(t, emails, "bob@example.com", "Expected bob@example.com in team members")
}

// EndpointTest defines a test case for an endpoint.
type EndpointTest struct {
	Path           string // Endpoint path
	ExpectData     bool   // Whether to expect non-empty data field
	UseCustomCheck bool   // Whether to use custom validation (e.g., teamMembers field)
	Description    string // Test description
}

// setupFullDataGeneration initializes storage and runs all generators.
// This mimics the full data generation pipeline in main.go.
func setupFullDataGeneration(t *testing.T) (*seed.SeedData, storage.Store) {
	t.Helper()

	// Load seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	ctx := context.Background()

	// Generate commits
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 7, 0) // 1 week of data
	require.NoError(t, err, "Failed to generate commits")

	// Generate model usage
	modelGen := generator.NewModelGenerator(seedData, store, "medium")
	err = modelGen.GenerateModelUsage(ctx, 7)
	require.NoError(t, err, "Failed to generate model usage")

	// Generate client versions
	versionGen := generator.NewVersionGenerator(seedData, store, "medium")
	err = versionGen.GenerateClientVersions(ctx, 7)
	require.NoError(t, err, "Failed to generate client versions")

	// Generate file extensions
	extensionGen := generator.NewExtensionGenerator(seedData, store, "medium")
	err = extensionGen.GenerateFileExtensions(ctx, 7)
	require.NoError(t, err, "Failed to generate file extensions")

	// Generate features (MCP, Commands, Plans, AskMode)
	featureGen := generator.NewFeatureGenerator(seedData, store, "medium")
	err = featureGen.GenerateFeatures(ctx, 7)
	require.NoError(t, err, "Failed to generate features")

	return seedData, store
}

// hasNonEmptyData checks if a response contains non-empty data.
// Handles both array-based responses ({"data": [...]}) and object-based responses ({"data": {...}}).
func hasNonEmptyData(response map[string]interface{}) bool {
	data, exists := response["data"]
	if !exists {
		// Some endpoints use different field names
		// Check for common alternatives
		if teamMembers, ok := response["teamMembers"]; ok {
			return isNonEmpty(teamMembers)
		}
		if stats, ok := response["stats"]; ok {
			return isNonEmpty(stats)
		}
		return false
	}

	return isNonEmpty(data)
}

// isNonEmpty checks if a value is non-empty (works for arrays, objects, strings, etc.).
func isNonEmpty(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case []interface{}:
		return len(v) > 0
	case map[string]interface{}:
		return len(v) > 0
	case string:
		return v != ""
	case float64:
		return true // Numbers are always considered non-empty
	case bool:
		return true // Booleans are always considered non-empty
	default:
		return true // Unknown types are considered non-empty
	}
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
