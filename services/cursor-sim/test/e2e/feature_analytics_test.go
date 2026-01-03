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

// TestE2E_TeamMCP verifies that /analytics/team/mcp returns MCP tool usage data.
func TestE2E_TeamMCP(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate commits first (features derive from commits)
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 7) // 1 week of data
	require.NoError(t, err, "Failed to generate commits")

	// Generate feature events (MCP, Commands, Plans, AskMode)
	featureGen := generator.NewFeatureGenerator(seedData, store, "medium")
	err = featureGen.GenerateFeatures(ctx, 7)
	require.NoError(t, err, "Failed to generate features")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/team/mcp", nil)
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

	// Verify MCP tools are realistic
	mcpTools := extractMCPTools(response.Data)
	assert.NotEmpty(t, mcpTools, "Expected at least one MCP tool in response")

	// Common MCP tools: read_file, write_file, list_directory, search, execute_command
	validTools := []string{"read_file", "write_file", "list_directory", "search", "execute_command", "git_status", "grep"}
	foundValid := false
	for _, tool := range mcpTools {
		for _, valid := range validTools {
			if tool == valid {
				foundValid = true
				break
			}
		}
		if foundValid {
			break
		}
	}
	assert.True(t, foundValid, "Expected at least one valid MCP tool")
}

// TestE2E_TeamCommands verifies that /analytics/team/commands returns command usage data.
func TestE2E_TeamCommands(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate commits
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 7)
	require.NoError(t, err, "Failed to generate commits")

	// Generate feature events
	featureGen := generator.NewFeatureGenerator(seedData, store, "medium")
	err = featureGen.GenerateFeatures(ctx, 7)
	require.NoError(t, err, "Failed to generate features")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/team/commands", nil)
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

	// Verify commands are realistic
	commands := extractCommandNames(response.Data)
	assert.NotEmpty(t, commands, "Expected at least one command in response")

	// Common commands: explain, refactor, fix, optimize, test, document
	validCommands := []string{"explain", "refactor", "fix", "optimize", "test", "document", "review"}
	foundValid := false
	for _, cmd := range commands {
		for _, valid := range validCommands {
			if cmd == valid {
				foundValid = true
				break
			}
		}
		if foundValid {
			break
		}
	}
	assert.True(t, foundValid, "Expected at least one valid command")
}

// TestE2E_TeamPlans verifies that /analytics/team/plans returns plan usage data.
func TestE2E_TeamPlans(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate commits
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 7)
	require.NoError(t, err, "Failed to generate commits")

	// Generate feature events
	featureGen := generator.NewFeatureGenerator(seedData, store, "medium")
	err = featureGen.GenerateFeatures(ctx, 7)
	require.NoError(t, err, "Failed to generate features")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/team/plans", nil)
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

	// Verify models are realistic (plans use models)
	models := extractModelNamesFromFeatures(response.Data)
	assert.NotEmpty(t, models, "Expected at least one model in response")
}

// TestE2E_TeamAskMode verifies that /analytics/team/ask-mode returns ask mode usage data.
func TestE2E_TeamAskMode(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate commits
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 7)
	require.NoError(t, err, "Failed to generate commits")

	// Generate feature events
	featureGen := generator.NewFeatureGenerator(seedData, store, "medium")
	err = featureGen.GenerateFeatures(ctx, 7)
	require.NoError(t, err, "Failed to generate features")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/team/ask-mode", nil)
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

	// Verify models are realistic
	models := extractModelNamesFromFeatures(response.Data)
	assert.NotEmpty(t, models, "Expected at least one model in response")
}

// TestE2E_ByUserMCP verifies that /analytics/by-user/mcp returns per-user MCP data.
func TestE2E_ByUserMCP(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate commits
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 7)
	require.NoError(t, err, "Failed to generate commits")

	// Generate feature events
	featureGen := generator.NewFeatureGenerator(seedData, store, "medium")
	err = featureGen.GenerateFeatures(ctx, 7)
	require.NoError(t, err, "Failed to generate features")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/by-user/mcp", nil)
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

// TestE2E_ByUserCommands verifies that /analytics/by-user/commands returns per-user command data.
func TestE2E_ByUserCommands(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate commits
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 7)
	require.NoError(t, err, "Failed to generate commits")

	// Generate feature events
	featureGen := generator.NewFeatureGenerator(seedData, store, "medium")
	err = featureGen.GenerateFeatures(ctx, 7)
	require.NoError(t, err, "Failed to generate features")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/by-user/commands", nil)
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

	// Verify data is present
	assert.NotEmpty(t, response.Data, "Expected non-empty data object")

	// Verify both users have data
	dataJSON, _ := json.Marshal(response.Data)
	dataStr := string(dataJSON)
	assert.Contains(t, dataStr, "alice@example.com", "Expected alice@example.com in response")
	assert.Contains(t, dataStr, "bob@example.com", "Expected bob@example.com in response")
}

// TestE2E_ByUserPlans verifies that /analytics/by-user/plans returns per-user plan data.
func TestE2E_ByUserPlans(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate commits
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 7)
	require.NoError(t, err, "Failed to generate commits")

	// Generate feature events
	featureGen := generator.NewFeatureGenerator(seedData, store, "medium")
	err = featureGen.GenerateFeatures(ctx, 7)
	require.NoError(t, err, "Failed to generate features")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/by-user/plans", nil)
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

	// Verify data is present
	assert.NotEmpty(t, response.Data, "Expected non-empty data object")

	// Verify both users have data
	dataJSON, _ := json.Marshal(response.Data)
	dataStr := string(dataJSON)
	assert.Contains(t, dataStr, "alice@example.com", "Expected alice@example.com in response")
	assert.Contains(t, dataStr, "bob@example.com", "Expected bob@example.com in response")
}

// TestE2E_ByUserAskMode verifies that /analytics/by-user/ask-mode returns per-user ask mode data.
func TestE2E_ByUserAskMode(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate commits
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 7)
	require.NoError(t, err, "Failed to generate commits")

	// Generate feature events
	featureGen := generator.NewFeatureGenerator(seedData, store, "medium")
	err = featureGen.GenerateFeatures(ctx, 7)
	require.NoError(t, err, "Failed to generate features")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/analytics/by-user/ask-mode", nil)
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

	// Verify data is present
	assert.NotEmpty(t, response.Data, "Expected non-empty data object")

	// Verify both users have data
	dataJSON, _ := json.Marshal(response.Data)
	dataStr := string(dataJSON)
	assert.Contains(t, dataStr, "alice@example.com", "Expected alice@example.com in response")
	assert.Contains(t, dataStr, "bob@example.com", "Expected bob@example.com in response")
}

// TestE2E_FeaturesEmpty_WithoutGenerator documents bug: returns empty without feature generation.
func TestE2E_FeaturesEmpty_WithoutGenerator(t *testing.T) {
	// This test shows the bug: without FeatureGenerator, endpoints return empty

	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// NOTE: Intentionally NOT calling FeatureGenerator to show the bug

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Test team MCP endpoint
	req := httptest.NewRequest("GET", "/analytics/team/mcp", nil)
	req.SetBasicAuth("test-api-key", "")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected 200 OK")

	var response TeamAnalyticsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")

	// BUG: Returns empty array without feature generation
	assert.Empty(t, response.Data, "BUG: Without FeatureGenerator, returns empty data")
}

// extractMCPTools extracts MCP tool names from analytics data.
// Response format: [{"event_date": "...", "tool_name": "read_file", "mcp_server_name": "filesystem", "usage": 10}, ...]
func extractMCPTools(data []map[string]interface{}) []string {
	tools := make([]string, 0)
	seen := make(map[string]bool)

	for _, item := range data {
		if toolName, ok := item["tool_name"].(string); ok {
			if toolName != "" && !seen[toolName] {
				tools = append(tools, toolName)
				seen[toolName] = true
			}
		}
	}

	return tools
}

// extractCommandNames extracts command names from analytics data.
// Response format: [{"event_date": "...", "command_name": "explain", "usage": 5}, ...]
func extractCommandNames(data []map[string]interface{}) []string {
	commands := make([]string, 0)
	seen := make(map[string]bool)

	for _, item := range data {
		if cmdName, ok := item["command_name"].(string); ok {
			if cmdName != "" && !seen[cmdName] {
				commands = append(commands, cmdName)
				seen[cmdName] = true
			}
		}
	}

	return commands
}

// extractModelNamesFromFeatures extracts model names from feature analytics data (plans, ask-mode).
// Response format: [{"event_date": "...", "model": "gpt-4-turbo", "usage": 3}, ...]
func extractModelNamesFromFeatures(data []map[string]interface{}) []string {
	models := make([]string, 0)
	seen := make(map[string]bool)

	for _, item := range data {
		if modelName, ok := item["model"].(string); ok {
			if modelName != "" && !seen[modelName] {
				models = append(models, modelName)
				seen[modelName] = true
			}
		}
	}

	return models
}
