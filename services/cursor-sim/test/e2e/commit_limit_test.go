package e2e

import (
	"context"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"encoding/json"
	"fmt"
	"net/http"
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
	commitLimitTestPort   = 18090
	commitLimitTestAPIKey = "test-commit-limit-key"
)

// TestE2E_CommitLimit verifies that maxCommits parameter correctly limits
// the total number of commits generated and returned by the API.
func TestE2E_CommitLimit(t *testing.T) {
	maxCommits := 50

	// Setup test server with commit limit
	cleanup := setupTestServerWithCommitLimit(t, maxCommits)
	defer cleanup()

	// Query the commits endpoint with pagination
	allCommits := fetchAllCommits(t)

	// Assert exact count matches the limit
	assert.Equal(t, maxCommits, len(allCommits),
		"Expected exactly %d commits, got %d", maxCommits, len(allCommits))
}

// TestE2E_CommitLimitPagination verifies that pagination works correctly
// when commit limit is enforced.
func TestE2E_CommitLimitPagination(t *testing.T) {
	maxCommits := 25
	pageSize := 10

	cleanup := setupTestServerWithCommitLimit(t, maxCommits)
	defer cleanup()

	// Fetch page 1
	page1 := fetchCommitsPage(t, 1, pageSize)
	assert.Equal(t, pageSize, len(page1["items"].([]interface{})),
		"Page 1 should have %d items", pageSize)

	// Fetch page 2
	page2 := fetchCommitsPage(t, 2, pageSize)
	assert.Equal(t, pageSize, len(page2["items"].([]interface{})),
		"Page 2 should have %d items", pageSize)

	// Fetch page 3 (should have remaining 5)
	page3 := fetchCommitsPage(t, 3, pageSize)
	assert.Equal(t, 5, len(page3["items"].([]interface{})),
		"Page 3 should have 5 items (25 total - 20 from first 2 pages)")

	// Verify totalCount is correct
	totalCount := int(page1["totalCount"].(float64))
	assert.Equal(t, maxCommits, totalCount,
		"Total count should match max commits")
}

// setupTestServerWithCommitLimit creates a test server with a specific commit limit
func setupTestServerWithCommitLimit(t *testing.T, maxCommits int) context.CancelFunc {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err)

	// Generate commits with limit
	// Use 30 days to ensure we hit the limit, commits fall within default 30-day query range
	gen := generator.NewCommitGeneratorWithSeed(seedData, store, "high", 42)
	ctx := context.Background()
	err = gen.GenerateCommits(ctx, 30, maxCommits) // 30 days, limited commits
	require.NoError(t, err)

	// Verify we generated exactly maxCommits
	allCommits := store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour))
	require.Equal(t, maxCommits, len(allCommits),
		"Store should contain exactly %d commits", maxCommits)

	// Create and start HTTP server
	router := server.NewRouter(store, seedData, commitLimitTestAPIKey, createTestConfig(), testVersion)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", commitLimitTestPort),
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
	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			t.Logf("Server shutdown error: %v", err)
		}
	}
}

// fetchAllCommits fetches all commits across all pages
func fetchAllCommits(t *testing.T) []interface{} {
	var allCommits []interface{}
	page := 1
	pageSize := 20

	for {
		result := fetchCommitsPage(t, page, pageSize)

		items := result["items"].([]interface{})
		allCommits = append(allCommits, items...)

		totalCount := int(result["totalCount"].(float64))
		currentPage := int(result["page"].(float64))
		currentPageSize := int(result["pageSize"].(float64))

		// Check if we've fetched all items
		if len(allCommits) >= totalCount {
			break
		}

		// Check if this was the last page
		if currentPage*currentPageSize >= totalCount {
			break
		}

		page++
	}

	return allCommits
}

// fetchCommitsPage fetches a single page of commits
func fetchCommitsPage(t *testing.T, page, pageSize int) map[string]interface{} {
	url := fmt.Sprintf("http://localhost:%d/analytics/ai-code/commits?page=%d&pageSize=%d",
		commitLimitTestPort, page, pageSize)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	req.SetBasicAuth(commitLimitTestAPIKey, "")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response structure
	assert.Contains(t, result, "items")
	assert.Contains(t, result, "totalCount")
	assert.Contains(t, result, "page")
	assert.Contains(t, result, "pageSize")

	return result
}
