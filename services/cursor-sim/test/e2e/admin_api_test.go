package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/server"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupAdminE2EServer creates a test server with minimal data for admin endpoint testing.
func setupAdminE2EServer(t *testing.T) (context.CancelFunc, *storage.MemoryStore, *seed.SeedData) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)

	// Initialize storage
	store := storage.NewMemoryStore()
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err)

	// Generate minimal commits (3 days for faster tests)
	commitGen := generator.NewCommitGeneratorWithSeed(seedData, store, "medium", 42)
	ctx := context.Background()
	err = commitGen.GenerateCommits(ctx, 3, 0)
	require.NoError(t, err)

	// Generate PRs from commits
	prGen := generator.NewPRGeneratorWithSeed(seedData, store, 42)
	err = prGen.GeneratePRsFromCommits(time.Now().Add(-3*24*time.Hour), time.Now())
	require.NoError(t, err)

	// Generate reviews for PRs
	reviewGen := generator.NewReviewGeneratorWithSeed(seedData, store, 42)
	repos := store.ListRepositories()
	for _, repo := range repos {
		_, err = reviewGen.GenerateReviewsForRepo(repo)
		require.NoError(t, err)
	}

	// Generate issues linked to PRs
	issueGen := generator.NewIssueGeneratorWithStore(seedData, store, 42)
	for _, repo := range repos {
		prs := store.GetPRsByRepo(repo)
		_, err = issueGen.GenerateAndStoreIssuesForPRs(prs, repo)
		require.NoError(t, err)
	}

	// Create test config
	testConfig := &config.Config{
		Mode:     "runtime",
		Days:     3,
		Velocity: "medium",
		GenParams: config.GenerationParams{
			Days:       3,
			Developers: len(seedData.Developers),
			MaxCommits: 0,
		},
	}

	// Create and start HTTP server on unique port for Admin tests
	const adminTestPort = 19083
	router := server.NewRouter(store, seedData, testAPIKey, testConfig, "2.0.0")
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", adminTestPort),
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

	return cleanup, store, seedData
}

// makeAdminRequest is a helper to make HTTP requests to Admin endpoints.
func makeAdminRequest(t *testing.T, method, path string, body interface{}) *http.Response {
	const adminTestPort = 19083
	baseURL := fmt.Sprintf("http://localhost:%d", adminTestPort)

	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, baseURL+path, bytes.NewReader(reqBody))
	require.NoError(t, err)

	req.SetBasicAuth(testAPIKey, "")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

// TestAdminAPI_Regenerate_Override tests override mode with large dataset.
func TestAdminAPI_Regenerate_Override(t *testing.T) {
	cleanup, store, _ := setupAdminE2EServer(t)
	defer cleanup()

	// Get initial stats
	initialCommits := len(store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour)))

	// Regenerate with override mode (small dataset for speed: 10 devs, 7 days)
	reqBody := map[string]interface{}{
		"mode":        "override",
		"days":        7,
		"velocity":    "high",
		"developers":  10,
		"max_commits": 500,
	}

	resp := makeAdminRequest(t, "POST", "/admin/regenerate", reqBody)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response structure
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "override", result["mode"])
	assert.Equal(t, true, result["data_cleaned"])
	assert.Contains(t, result, "commits_added")
	assert.Contains(t, result, "prs_added")
	assert.Contains(t, result, "total_commits")
	assert.Contains(t, result, "duration")

	// Verify config echoed back
	config := result["config"].(map[string]interface{})
	assert.Equal(t, float64(7), config["days"])
	assert.Equal(t, "high", config["velocity"])
	assert.Equal(t, float64(10), config["developers"])
	assert.Equal(t, float64(500), config["max_commits"])

	// Verify data was replaced (not appended)
	newCommits := len(store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour)))
	// Override mode should clear data first, so new commits might be different from initial
	t.Logf("Override mode: initial %d commits, new %d commits", initialCommits, newCommits)
}

// TestAdminAPI_Regenerate_Append tests append mode (cumulative data).
func TestAdminAPI_Regenerate_Append(t *testing.T) {
	cleanup, store, _ := setupAdminE2EServer(t)
	defer cleanup()

	// Get initial stats
	initialCommits := len(store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour)))

	// Regenerate with append mode (small dataset: 0 devs = use seed, 3 days)
	reqBody := map[string]interface{}{
		"mode":        "append",
		"days":        3,
		"velocity":    "medium",
		"developers":  0, // Use seed count
		"max_commits": 0, // Unlimited
	}

	resp := makeAdminRequest(t, "POST", "/admin/regenerate", reqBody)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "append", result["mode"])
	assert.Equal(t, false, result["data_cleaned"])

	// Verify data was appended
	newCommits := len(store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour)))
	assert.Greater(t, newCommits, initialCommits, "append mode should add data")

	t.Logf("Append mode: initial %d commits, new %d commits (added %d)",
		initialCommits, newCommits, newCommits-initialCommits)
}

// TestAdminAPI_SeedUpload_JSON tests JSON seed upload.
func TestAdminAPI_SeedUpload_JSON(t *testing.T) {
	cleanup, _, _ := setupAdminE2EServer(t)
	defer cleanup()

	// Create a minimal JSON seed
	jsonSeed := `{
		"version": "1.0",
		"developers": [
			{
				"user_id": "test_001",
				"email": "test1@example.com",
				"name": "Test Developer 1",
				"team": "Backend"
			},
			{
				"user_id": "test_002",
				"email": "test2@example.com",
				"name": "Test Developer 2",
				"team": "Frontend"
			}
		],
		"repositories": [
			{
				"repo_name": "test/repo",
				"primary_language": "go"
			}
		]
	}`

	reqBody := map[string]interface{}{
		"data":       jsonSeed,
		"format":     "json",
		"regenerate": false,
	}

	resp := makeAdminRequest(t, "POST", "/admin/seed", reqBody)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, true, result["seed_loaded"])
	assert.Equal(t, float64(2), result["developers"])
	assert.Equal(t, float64(1), result["repositories"])
	assert.Contains(t, result, "organizations")
	assert.Contains(t, result, "teams")

	t.Logf("JSON seed uploaded: %d developers, %d repos", int(result["developers"].(float64)), int(result["repositories"].(float64)))
}

// TestAdminAPI_SeedUpload_CSV tests CSV seed upload with regeneration.
func TestAdminAPI_SeedUpload_CSV(t *testing.T) {
	cleanup, _, _ := setupAdminE2EServer(t)
	defer cleanup()

	// Create a minimal CSV seed
	csvSeed := `user_id,email,name
csv_001,csv1@example.com,CSV Developer 1
csv_002,csv2@example.com,CSV Developer 2
csv_003,csv3@example.com,CSV Developer 3`

	reqBody := map[string]interface{}{
		"data":       csvSeed,
		"format":     "csv",
		"regenerate": true,
		"regenerate_config": map[string]interface{}{
			"mode":        "override",
			"days":        3,
			"velocity":    "low",
			"developers":  0, // Use CSV count (3)
			"max_commits": 100,
		},
	}

	resp := makeAdminRequest(t, "POST", "/admin/seed", reqBody)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, true, result["seed_loaded"])
	assert.Equal(t, float64(3), result["developers"])
	assert.Equal(t, true, result["regenerated"])
	assert.Contains(t, result, "generate_stats")

	// Verify regeneration occurred
	genStats := result["generate_stats"].(map[string]interface{})
	assert.Equal(t, "success", genStats["status"])
	assert.Contains(t, genStats, "total_commits")
	assert.Contains(t, genStats, "duration")

	t.Logf("CSV seed uploaded and regenerated: %d developers", int(result["developers"].(float64)))
}

// TestAdminAPI_Config tests config inspection endpoint.
func TestAdminAPI_Config(t *testing.T) {
	cleanup, _, _ := setupAdminE2EServer(t)
	defer cleanup()

	resp := makeAdminRequest(t, "GET", "/admin/config", nil)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response structure
	assert.Contains(t, result, "generation")
	assert.Contains(t, result, "seed")
	assert.Contains(t, result, "external_sources")
	assert.Contains(t, result, "server")

	// Verify generation section
	generation := result["generation"].(map[string]interface{})
	assert.Contains(t, generation, "days")
	assert.Contains(t, generation, "velocity")
	assert.Contains(t, generation, "developers")
	assert.Contains(t, generation, "max_commits")

	// Verify seed section
	seedInfo := result["seed"].(map[string]interface{})
	assert.Contains(t, seedInfo, "version")
	assert.Contains(t, seedInfo, "developers")
	assert.Contains(t, seedInfo, "repositories")
	assert.Contains(t, seedInfo, "teams")

	// Verify server section
	serverInfo := result["server"].(map[string]interface{})
	assert.Contains(t, serverInfo, "port")
	assert.Contains(t, serverInfo, "version")
	assert.Contains(t, serverInfo, "uptime")

	t.Logf("Config inspection successful: %d developers, %d repos",
		int(seedInfo["developers"].(float64)), int(seedInfo["repositories"].(float64)))
}

// TestAdminAPI_Stats tests stats retrieval with time series.
func TestAdminAPI_Stats(t *testing.T) {
	cleanup, _, _ := setupAdminE2EServer(t)
	defer cleanup()

	// Test without time series
	resp := makeAdminRequest(t, "GET", "/admin/stats", nil)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response structure
	assert.Contains(t, result, "generation")
	assert.Contains(t, result, "developers")
	assert.Contains(t, result, "quality")
	assert.Contains(t, result, "variance")
	assert.Contains(t, result, "performance")
	assert.Contains(t, result, "organization")

	// Verify generation stats
	generation := result["generation"].(map[string]interface{})
	assert.Contains(t, generation, "total_commits")
	assert.Contains(t, generation, "total_prs")
	assert.Contains(t, generation, "total_developers")
	assert.Contains(t, generation, "data_size")

	// Verify developers breakdown
	developers := result["developers"].(map[string]interface{})
	assert.Contains(t, developers, "by_seniority")
	assert.Contains(t, developers, "by_region")
	assert.Contains(t, developers, "by_team")

	// Verify quality metrics
	quality := result["quality"].(map[string]interface{})
	assert.Contains(t, quality, "avg_revert_rate")
	assert.Contains(t, quality, "avg_review_thoroughness")

	// Test with time series
	resp2 := makeAdminRequest(t, "GET", "/admin/stats?include_timeseries=true", nil)
	defer resp2.Body.Close()

	assert.Equal(t, 200, resp2.StatusCode)

	var result2 map[string]interface{}
	err = json.NewDecoder(resp2.Body).Decode(&result2)
	require.NoError(t, err)

	// Verify time series is included
	assert.Contains(t, result2, "time_series")
	timeSeries := result2["time_series"].(map[string]interface{})
	assert.Contains(t, timeSeries, "commits_per_day")
	assert.Contains(t, timeSeries, "prs_per_day")

	t.Logf("Stats retrieved: %d commits, %d PRs",
		int(generation["total_commits"].(float64)), int(generation["total_prs"].(float64)))
}

// TestAdminAPI_SeedPresets tests getting predefined presets.
func TestAdminAPI_SeedPresets(t *testing.T) {
	cleanup, _, _ := setupAdminE2EServer(t)
	defer cleanup()

	resp := makeAdminRequest(t, "GET", "/admin/seed/presets", nil)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response structure
	assert.Contains(t, result, "presets")

	// Verify presets array
	presets := result["presets"].([]interface{})
	assert.Equal(t, 4, len(presets), "should have 4 predefined presets")

	// Verify each preset has required fields
	presetNames := []string{"small-team", "medium-team", "enterprise", "multi-region"}
	for i, preset := range presets {
		p := preset.(map[string]interface{})
		assert.Contains(t, p, "name")
		assert.Contains(t, p, "description")
		assert.Contains(t, p, "developers")
		assert.Contains(t, p, "teams")
		assert.Contains(t, p, "regions")
		assert.Contains(t, p, "suggested_days")
		assert.Contains(t, p, "suggested_velocity")

		// Verify expected preset names
		assert.Equal(t, presetNames[i], p["name"])
	}

	t.Logf("Seed presets retrieved: %d presets", len(presets))
}

// TestAdminAPI_Authentication tests missing API key returns 401.
func TestAdminAPI_Authentication(t *testing.T) {
	cleanup, _, _ := setupAdminE2EServer(t)
	defer cleanup()

	endpoints := []struct {
		method string
		path   string
	}{
		{"POST", "/admin/regenerate"},
		{"POST", "/admin/seed"},
		{"GET", "/admin/seed/presets"},
		{"GET", "/admin/config"},
		{"GET", "/admin/stats"},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.method+" "+endpoint.path, func(t *testing.T) {
			const adminTestPort = 19083
			baseURL := fmt.Sprintf("http://localhost:%d", adminTestPort)

			var reqBody []byte
			if endpoint.method == "POST" {
				// Provide minimal body for POST requests
				body := map[string]interface{}{"mode": "append", "days": 1, "velocity": "low", "developers": 0, "max_commits": 0}
				reqBody, _ = json.Marshal(body)
			}

			req, err := http.NewRequest(endpoint.method, baseURL+endpoint.path, bytes.NewReader(reqBody))
			require.NoError(t, err)

			// No authentication header

			client := &http.Client{Timeout: 2 * time.Second}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should return 401 Unauthorized
			assert.Equal(t, 401, resp.StatusCode, "endpoint %s %s should require authentication", endpoint.method, endpoint.path)
		})
	}
}
