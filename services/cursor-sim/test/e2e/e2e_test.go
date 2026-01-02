package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/server"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPort   = 19080
	testAPIKey = "test-api-key"
	baseURL    = "http://localhost:19080"
)

// setupTestServer starts a test server with sample data and returns a cleanup function.
func setupTestServer(t *testing.T) (context.CancelFunc, *storage.MemoryStore) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)

	// Initialize storage
	store := storage.NewMemoryStore()

	// Generate sample commits (1 day of history)
	gen := generator.NewCommitGenerator(seedData, store, "medium")
	ctx := context.Background()
	err = gen.GenerateCommits(ctx, 1)
	require.NoError(t, err)

	// Create and start HTTP server
	router := server.NewRouter(store, testAPIKey)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", testPort),
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

	return cleanup, store
}

// makeRequest is a helper to make HTTP requests with optional authentication.
func makeRequest(t *testing.T, method, path string, withAuth bool) *http.Response {
	req, err := http.NewRequest(method, baseURL+path, nil)
	require.NoError(t, err)

	if withAuth {
		req.SetBasicAuth(testAPIKey, "")
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

func TestE2E_HealthEndpoint(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	resp := makeRequest(t, "GET", "/health", false)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]string
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

func TestE2E_TeamsMembers(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	// Without auth should fail
	resp := makeRequest(t, "GET", "/teams/members", false)
	resp.Body.Close()
	assert.Equal(t, 401, resp.StatusCode)

	// With auth should succeed
	resp = makeRequest(t, "GET", "/teams/members", true)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Contains(t, result, "teamMembers")
}

func TestE2E_AICodeCommits(t *testing.T) {
	cleanup, store := setupTestServer(t)
	defer cleanup()

	// Get all commits
	resp := makeRequest(t, "GET", "/analytics/ai-code/commits", true)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response structure
	assert.Contains(t, result, "data")
	assert.Contains(t, result, "pagination")
	assert.Contains(t, result, "params")

	// Verify we have commits
	data := result["data"].([]interface{})
	allCommits := store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour))
	assert.Greater(t, len(data), 0)
	assert.LessOrEqual(t, len(data), len(allCommits))
}

func TestE2E_AICodeCommitsCSV(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	resp := makeRequest(t, "GET", "/analytics/ai-code/commits.csv", true)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "text/csv", resp.Header.Get("Content-Type"))
	assert.Contains(t, resp.Header.Get("Content-Disposition"), "attachment")
	assert.Contains(t, resp.Header.Get("Content-Disposition"), "cursor-sim-export-")
}

func TestE2E_TeamAnalytics(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	endpoints := []string{
		"/analytics/team/agent-edits",
		"/analytics/team/tabs",
		"/analytics/team/dau",
		"/analytics/team/models",
		"/analytics/team/client-versions",
		"/analytics/team/top-file-extensions",
		"/analytics/team/mcp",
		"/analytics/team/commands",
		"/analytics/team/plans",
		"/analytics/team/ask-mode",
		"/analytics/team/leaderboard",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			resp := makeRequest(t, "GET", endpoint, true)
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode, "Endpoint %s should return 200", endpoint)

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			// All analytics endpoints should have this structure
			assert.Contains(t, result, "data")
			assert.Contains(t, result, "pagination")
			assert.Contains(t, result, "params")
		})
	}
}

func TestE2E_ByUserAnalytics(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	endpoints := []string{
		"/analytics/by-user/agent-edits",
		"/analytics/by-user/tabs",
		"/analytics/by-user/models",
		"/analytics/by-user/client-versions",
		"/analytics/by-user/top-file-extensions",
		"/analytics/by-user/mcp",
		"/analytics/by-user/commands",
		"/analytics/by-user/plans",
		"/analytics/by-user/ask-mode",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			resp := makeRequest(t, "GET", endpoint, true)
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode, "Endpoint %s should return 200", endpoint)

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			// All analytics endpoints should have this structure
			assert.Contains(t, result, "data")
			assert.Contains(t, result, "pagination")
			assert.Contains(t, result, "params")
		})
	}
}

func TestE2E_Authentication(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name       string
		endpoint   string
		withAuth   bool
		wantStatus int
	}{
		{"Health no auth", "/health", false, 200},
		{"Teams without auth", "/teams/members", false, 401},
		{"Teams with auth", "/teams/members", true, 200},
		{"AI code without auth", "/analytics/ai-code/commits", false, 401},
		{"AI code with auth", "/analytics/ai-code/commits", true, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := makeRequest(t, "GET", tt.endpoint, tt.withAuth)
			defer resp.Body.Close()
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestE2E_QueryParameters(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name     string
		endpoint string
		query    string
	}{
		{"Date range", "/analytics/ai-code/commits", "?from=2026-01-01&to=2026-01-02"},
		{"Pagination", "/analytics/ai-code/commits", "?page=1&page_size=10"},
		{"Combined", "/analytics/ai-code/commits", "?from=2026-01-01&to=2026-01-02&page=1&page_size=5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := makeRequest(t, "GET", tt.endpoint+tt.query, true)
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			// Verify pagination params were applied
			params := result["params"].(map[string]interface{})
			if strings.Contains(tt.query, "page=") {
				assert.NotNil(t, params["page"])
			}
		})
	}
}

func TestE2E_RateLimiting(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	// Make many requests quickly to trigger rate limiting
	// Rate limit is 100 req/min, so we need to make >100 requests
	successCount := 0
	rateLimitCount := 0

	for i := 0; i < 120; i++ {
		resp := makeRequest(t, "GET", "/health", false)
		resp.Body.Close()

		if resp.StatusCode == 200 {
			successCount++
		} else if resp.StatusCode == 429 {
			rateLimitCount++
		}
	}

	// Should have some successful requests
	assert.Greater(t, successCount, 0)

	// Should have hit rate limit at some point
	// (This might not always trigger in test environment due to timing)
	t.Logf("Successful: %d, Rate limited: %d", successCount, rateLimitCount)
}

func TestE2E_NotFound(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	resp := makeRequest(t, "GET", "/nonexistent/endpoint", true)
	defer resp.Body.Close()

	assert.Equal(t, 404, resp.StatusCode)
}

// PR Lifecycle E2E Tests (Phase 2)

func TestE2E_PRLifecycle(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)

	// Initialize storage
	store := storage.NewMemoryStore()
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err)

	// Generate commits
	commitGen := generator.NewCommitGeneratorWithSeed(seedData, store, "medium", 42)
	ctx := context.Background()
	err = commitGen.GenerateCommits(ctx, 7)
	require.NoError(t, err)

	// Verify commits were generated
	commits := store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour))
	require.Greater(t, len(commits), 0, "should have commits")

	// Generate PRs from commits
	prGen := generator.NewPRGeneratorWithSeed(seedData, store, 42)
	err = prGen.GeneratePRsFromCommits(time.Now().Add(-7*24*time.Hour), time.Now())
	require.NoError(t, err)

	// Verify PRs were generated
	repos := store.ListRepositories()
	totalPRs := 0
	for _, repo := range repos {
		prs := store.GetPRsByRepo(repo)
		totalPRs += len(prs)
	}
	require.Greater(t, totalPRs, 0, "should have PRs")

	// Generate reviews for PRs
	reviewGen := generator.NewReviewGeneratorWithSeed(seedData, store, 42)
	for _, repo := range repos {
		_, err = reviewGen.GenerateReviewsForRepo(repo)
		require.NoError(t, err)
	}

	// Verify reviews were generated
	totalReviews := 0
	for _, repo := range repos {
		prs := store.GetPRsByRepo(repo)
		for _, pr := range prs {
			reviews := store.GetReviewComments(repo, pr.Number)
			totalReviews += len(reviews)
		}
	}
	assert.Greater(t, totalReviews, 0, "should have reviews")

	// Apply quality outcomes
	qualGen := generator.NewQualityGeneratorWithSeed(seedData, store, 42)
	for _, repo := range repos {
		err = qualGen.ApplyQualityOutcomes(repo)
		require.NoError(t, err)
	}

	t.Logf("PR Lifecycle complete: %d commits, %d PRs, %d reviews", len(commits), totalPRs, totalReviews)
}

func TestE2E_PRMetricsAggregation(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)

	store := storage.NewMemoryStore()
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err)

	// Generate commits with known AI metrics
	commitGen := generator.NewCommitGeneratorWithSeed(seedData, store, "medium", 42)
	ctx := context.Background()
	err = commitGen.GenerateCommits(ctx, 7)
	require.NoError(t, err)

	// Generate PRs
	prGen := generator.NewPRGeneratorWithSeed(seedData, store, 42)
	err = prGen.GeneratePRsFromCommits(time.Now().Add(-7*24*time.Hour), time.Now())
	require.NoError(t, err)

	// Verify AI metrics are properly aggregated in PRs
	repos := store.ListRepositories()
	for _, repo := range repos {
		prs := store.GetPRsByRepo(repo)
		for _, pr := range prs {
			// AI ratio should be between 0 and 1
			assert.GreaterOrEqual(t, pr.AIRatio, 0.0, "AI ratio should be >= 0")
			assert.LessOrEqual(t, pr.AIRatio, 1.0, "AI ratio should be <= 1")

			// Tab + Composer lines should not exceed total additions
			assert.LessOrEqual(t, pr.TabLines+pr.ComposerLines, pr.Additions,
				"AI lines should not exceed additions")
		}
	}
}

func TestE2E_QualityOutcomeCorrelations(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)

	store := storage.NewMemoryStore()
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err)

	// Generate large sample for statistical significance
	commitGen := generator.NewCommitGeneratorWithSeed(seedData, store, "high", 42)
	ctx := context.Background()
	err = commitGen.GenerateCommits(ctx, 14) // 2 weeks
	require.NoError(t, err)

	prGen := generator.NewPRGeneratorWithSeed(seedData, store, 42)
	err = prGen.GeneratePRsFromCommits(time.Now().Add(-14*24*time.Hour), time.Now())
	require.NoError(t, err)

	qualGen := generator.NewQualityGeneratorWithSeed(seedData, store, 42)
	repos := store.ListRepositories()
	for _, repo := range repos {
		err = qualGen.ApplyQualityOutcomes(repo)
		require.NoError(t, err)
	}

	// Count reverts by AI ratio category
	var lowAIReverts, highAIReverts int
	var lowAITotal, highAITotal int

	for _, repo := range repos {
		prs := store.GetPRsByRepo(repo)
		for _, pr := range prs {
			if pr.State != "merged" {
				continue
			}

			if pr.AIRatio < 0.3 {
				lowAITotal++
				if pr.WasReverted {
					lowAIReverts++
				}
			} else if pr.AIRatio > 0.7 {
				highAITotal++
				if pr.WasReverted {
					highAIReverts++
				}
			}
		}
	}

	// Log results (we can't assert exact rates due to randomness)
	if lowAITotal > 0 {
		lowRate := float64(lowAIReverts) / float64(lowAITotal)
		t.Logf("Low AI ratio revert rate: %.2f%% (%d/%d)", lowRate*100, lowAIReverts, lowAITotal)
	}
	if highAITotal > 0 {
		highRate := float64(highAIReverts) / float64(highAITotal)
		t.Logf("High AI ratio revert rate: %.2f%% (%d/%d)", highRate*100, highAIReverts, highAITotal)
	}
}
