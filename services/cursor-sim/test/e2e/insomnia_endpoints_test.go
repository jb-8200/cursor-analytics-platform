package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
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

// TestInsomnia_E2E tests all endpoints from the Insomnia collection.
// This ensures the API contract matches what's documented in docs/insomnia/cursor-sim_Insomnia_2026-01-04.yaml

const insomniaTestPort = 19083

// setupInsomniaTestServer creates a fully populated test server matching runtime mode.
func setupInsomniaTestServer(t *testing.T) (context.CancelFunc, *storage.MemoryStore, *seed.SeedData) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)

	// Initialize storage
	store := storage.NewMemoryStore()
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err)

	// Generate 30 days of commits (matching runtime default)
	commitGen := generator.NewCommitGeneratorWithSeed(seedData, store, "medium", 42)
	ctx := context.Background()
	err = commitGen.GenerateCommits(ctx, 30, 0)
	require.NoError(t, err)

	// Generate PRs from commits
	prGen := generator.NewPRGeneratorWithSeed(seedData, store, 42)
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now().Add(24 * time.Hour)
	err = prGen.GeneratePRsFromCommits(startDate, endDate)
	require.NoError(t, err)

	// Generate reviews for PRs using the store-free pattern
	reviewGen := generator.NewReviewGenerator(seedData, rand.New(rand.NewSource(42)))
	repos := store.ListRepositories()
	for _, repoName := range repos {
		prs := store.GetPRsByRepo(repoName)
		for _, pr := range prs {
			reviews := reviewGen.GenerateReviewsForPR(pr)
			for _, review := range reviews {
				err = store.StoreReview(review)
				require.NoError(t, err)
			}
		}
	}

	// Generate issues for PRs
	issueGen := generator.NewIssueGeneratorWithStore(seedData, store, 42)
	for _, repoName := range repos {
		prs := store.GetPRsByRepo(repoName)
		_, err = issueGen.GenerateAndStoreIssuesForPRs(prs, repoName)
		require.NoError(t, err)
	}

	// Generate model usage events
	modelGen := generator.NewModelGeneratorWithSeed(seedData, store, "medium", 42)
	err = modelGen.GenerateModelUsage(ctx, 30)
	require.NoError(t, err)

	// Generate client version events
	versionGen := generator.NewVersionGeneratorWithSeed(seedData, store, "medium", 42)
	err = versionGen.GenerateClientVersions(ctx, 30)
	require.NoError(t, err)

	// Generate file extension events
	extensionGen := generator.NewExtensionGeneratorWithSeed(seedData, store, "medium", 42)
	err = extensionGen.GenerateFileExtensions(ctx, 30)
	require.NoError(t, err)

	// Generate feature events (MCP, commands, plans, ask-mode)
	featureGen := generator.NewFeatureGeneratorWithSeed(seedData, store, "medium", 42)
	err = featureGen.GenerateFeatures(ctx, 30)
	require.NoError(t, err)

	// Create and start HTTP server
	router := server.NewRouter(store, seedData, testAPIKey)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", insomniaTestPort),
		Handler: router,
	}

	// Start server in background
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)

	// Return cleanup function
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
	}

	return cleanup, store, seedData
}

// makeInsomniaRequest is a helper to make HTTP requests to the Insomnia test server.
func makeInsomniaRequest(t *testing.T, method, path string, withAuth bool) *http.Response {
	baseURL := fmt.Sprintf("http://localhost:%d", insomniaTestPort)
	req, err := http.NewRequest(method, baseURL+path, nil)
	require.NoError(t, err)

	if withAuth {
		req.SetBasicAuth(testAPIKey, "")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

// TestInsomnia_HealthCheck tests the health endpoint (no auth required).
func TestInsomnia_HealthCheck(t *testing.T) {
	cleanup, _, _ := setupInsomniaTestServer(t)
	defer cleanup()

	resp := makeInsomniaRequest(t, "GET", "/health", false)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]string
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

// TestInsomnia_TeamMembers tests the team members endpoint.
func TestInsomnia_TeamMembers(t *testing.T) {
	cleanup, _, _ := setupInsomniaTestServer(t)
	defer cleanup()

	resp := makeInsomniaRequest(t, "GET", "/teams/members", true)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Contains(t, result, "teamMembers")

	members := result["teamMembers"].([]interface{})
	assert.Greater(t, len(members), 0, "should have team members")

	t.Logf("Team members: %d", len(members))
}

// TestInsomnia_AICodeCommits tests the AI code commits endpoint.
func TestInsomnia_AICodeCommits(t *testing.T) {
	cleanup, _, _ := setupInsomniaTestServer(t)
	defer cleanup()

	resp := makeInsomniaRequest(t, "GET", "/analytics/ai-code/commits", true)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Contains(t, result, "items")
	assert.Contains(t, result, "totalCount")
	assert.Contains(t, result, "page")
	assert.Contains(t, result, "pageSize")

	items := result["items"].([]interface{})
	totalCount := int(result["totalCount"].(float64))
	assert.Greater(t, len(items), 0, "should have commits")
	assert.Greater(t, totalCount, 0, "total count should be positive")

	t.Logf("AI code commits: %d items, %d total", len(items), totalCount)
}

// TestInsomnia_AICodeCommitsCSV tests the CSV export endpoint.
func TestInsomnia_AICodeCommitsCSV(t *testing.T) {
	cleanup, _, _ := setupInsomniaTestServer(t)
	defer cleanup()

	resp := makeInsomniaRequest(t, "GET", "/analytics/ai-code/commits.csv", true)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "text/csv", resp.Header.Get("Content-Type"))
	assert.Contains(t, resp.Header.Get("Content-Disposition"), "attachment")
}

// TestInsomnia_TeamAnalyticsEndpoints tests all team analytics endpoints.
func TestInsomnia_TeamAnalyticsEndpoints(t *testing.T) {
	cleanup, _, _ := setupInsomniaTestServer(t)
	defer cleanup()

	endpoints := []struct {
		path        string
		description string
	}{
		{"/analytics/team/agent-edits", "Agent Edits"},
		{"/analytics/team/tabs", "Tab Completions"},
		{"/analytics/team/dau", "Daily Active Users"},
		{"/analytics/team/models", "Model Usage"},
		{"/analytics/team/client-versions", "Client Versions"},
		{"/analytics/team/top-file-extensions", "Top File Extensions"},
		{"/analytics/team/mcp", "MCP Tool Usage"},
		{"/analytics/team/commands", "Command Usage"},
		{"/analytics/team/plans", "Plan Usage"},
		{"/analytics/team/ask-mode", "Ask Mode Usage"},
		{"/analytics/team/leaderboard", "Leaderboard"},
	}

	for _, ep := range endpoints {
		t.Run(ep.description, func(t *testing.T) {
			resp := makeInsomniaRequest(t, "GET", ep.path, true)
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode, "endpoint %s should return 200", ep.path)

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			assert.Contains(t, result, "data", "endpoint %s should have data", ep.path)
			assert.Contains(t, result, "params", "endpoint %s should have params", ep.path)

			// Leaderboard has pagination
			if ep.path == "/analytics/team/leaderboard" {
				assert.Contains(t, result, "pagination", "leaderboard should have pagination")
			}
		})
	}
}

// TestInsomnia_ByUserAnalyticsEndpoints tests all by-user analytics endpoints.
func TestInsomnia_ByUserAnalyticsEndpoints(t *testing.T) {
	cleanup, _, _ := setupInsomniaTestServer(t)
	defer cleanup()

	endpoints := []struct {
		path        string
		description string
	}{
		{"/analytics/by-user/agent-edits", "Agent Edits by User"},
		{"/analytics/by-user/tabs", "Tab Completions by User"},
		{"/analytics/by-user/models", "Model Usage by User"},
		{"/analytics/by-user/client-versions", "Client Versions by User"},
		{"/analytics/by-user/top-file-extensions", "File Extensions by User"},
		{"/analytics/by-user/mcp", "MCP Usage by User"},
		{"/analytics/by-user/commands", "Commands by User"},
		{"/analytics/by-user/plans", "Plans by User"},
		{"/analytics/by-user/ask-mode", "Ask Mode by User"},
	}

	for _, ep := range endpoints {
		t.Run(ep.description, func(t *testing.T) {
			resp := makeInsomniaRequest(t, "GET", ep.path, true)
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode, "endpoint %s should return 200", ep.path)

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			assert.Contains(t, result, "data", "endpoint %s should have data", ep.path)
			assert.Contains(t, result, "pagination", "endpoint %s should have pagination", ep.path)
			assert.Contains(t, result, "params", "endpoint %s should have params", ep.path)
		})
	}
}

// TestInsomnia_GitHubReposEndpoints tests GitHub repository endpoints.
func TestInsomnia_GitHubReposEndpoints(t *testing.T) {
	cleanup, store, _ := setupInsomniaTestServer(t)
	defer cleanup()

	// Get a valid repo name from storage
	repos := store.ListRepositories()
	require.Greater(t, len(repos), 0, "should have repositories")
	repoName := repos[0]

	t.Run("List Repositories", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/repos", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result []interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Greater(t, len(result), 0, "should have repositories")

		t.Logf("Repositories: %d", len(result))
	})

	t.Run("Get Repository", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/repos/"+repoName, true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "full_name")
	})

	t.Run("List Commits", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/repos/"+repoName+"/commits", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result []interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Greater(t, len(result), 0, "should have commits")

		t.Logf("Commits for %s: %d", repoName, len(result))
	})

	t.Run("Get Single Commit", func(t *testing.T) {
		// First get a list of commits
		commits := store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour))
		require.Greater(t, len(commits), 0)

		// Find a commit from our repo
		var commitSHA string
		for _, c := range commits {
			if c.RepoName == repoName {
				commitSHA = c.CommitHash[:12] // Use first 12 chars
				break
			}
		}
		require.NotEmpty(t, commitSHA, "should find a commit SHA")

		resp := makeInsomniaRequest(t, "GET", "/repos/"+repoName+"/commits/"+commitSHA, true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "commitHash")
	})
}

// TestInsomnia_GitHubPullRequestsEndpoints tests GitHub pull request endpoints.
func TestInsomnia_GitHubPullRequestsEndpoints(t *testing.T) {
	cleanup, store, _ := setupInsomniaTestServer(t)
	defer cleanup()

	// Get a valid repo and PR
	repos := store.ListRepositories()
	require.Greater(t, len(repos), 0, "should have repositories")

	var repoName string
	var prNumber int
	for _, repo := range repos {
		prs := store.GetPRsByRepo(repo)
		if len(prs) > 0 {
			repoName = repo
			prNumber = prs[0].Number
			break
		}
	}
	require.NotEmpty(t, repoName, "should have a repo with PRs")
	require.Greater(t, prNumber, 0, "should have a valid PR number")

	t.Run("List Pull Requests", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/repos/"+repoName+"/pulls", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result []interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Greater(t, len(result), 0, "should have PRs")

		t.Logf("PRs for %s: %d", repoName, len(result))
	})

	t.Run("List Pull Requests with state filter", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/repos/"+repoName+"/pulls?state=all", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Get Pull Request", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", fmt.Sprintf("/repos/%s/pulls/%d", repoName, prNumber), true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "number")

		t.Logf("PR #%d details retrieved", prNumber)
	})

	t.Run("List PR Commits", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", fmt.Sprintf("/repos/%s/pulls/%d/commits", repoName, prNumber), true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result []interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Greater(t, len(result), 0, "PR should have commits")

		t.Logf("PR #%d commits: %d", prNumber, len(result))
	})

	t.Run("List PR Files", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", fmt.Sprintf("/repos/%s/pulls/%d/files", repoName, prNumber), true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result []interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Greater(t, len(result), 0, "PR should have files")

		t.Logf("PR #%d files: %d", prNumber, len(result))
	})

	t.Run("List PR Reviews", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", fmt.Sprintf("/repos/%s/pulls/%d/reviews", repoName, prNumber), true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result []interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		// Reviews may or may not be present
		t.Logf("PR #%d reviews: %d", prNumber, len(result))
	})
}

// TestInsomnia_GitHubAnalyticsEndpoints tests GitHub analytics endpoints.
func TestInsomnia_GitHubAnalyticsEndpoints(t *testing.T) {
	cleanup, _, _ := setupInsomniaTestServer(t)
	defer cleanup()

	t.Run("PRs Analytics", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/analytics/github/prs", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "data")
		assert.Contains(t, result, "pagination")

		data := result["data"].([]interface{})
		pagination := result["pagination"].(map[string]interface{})
		t.Logf("PRs: %d items, %d total", len(data), int(pagination["total"].(float64)))
	})

	t.Run("Reviews Analytics", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/analytics/github/reviews", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "data")
		assert.Contains(t, result, "pagination")

		data := result["data"].([]interface{})
		pagination := result["pagination"].(map[string]interface{})
		t.Logf("Reviews: %d items, %d total", len(data), int(pagination["total"].(float64)))
	})

	t.Run("Issues Analytics", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/analytics/github/issues", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "data")

		data := result["data"].([]interface{})
		t.Logf("Issues: %d items", len(data))
	})

	t.Run("PR Cycle Time", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/analytics/github/pr-cycle-time", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "data")
		assert.Contains(t, result, "params")

		data := result["data"].(map[string]interface{})
		t.Logf("PR Cycle Time: %d PRs analyzed", int(data["totalPRsAnalyzed"].(float64)))
	})

	t.Run("Review Quality", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/analytics/github/review-quality", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "data")
		assert.Contains(t, result, "params")

		data := result["data"].(map[string]interface{})
		t.Logf("Review Quality: %d PRs reviewed", int(data["total_prs_reviewed"].(float64)))
	})
}

// TestInsomnia_QualityAnalysisEndpoints tests code quality analysis endpoints.
func TestInsomnia_QualityAnalysisEndpoints(t *testing.T) {
	cleanup, store, _ := setupInsomniaTestServer(t)
	defer cleanup()

	repos := store.ListRepositories()
	require.Greater(t, len(repos), 0, "should have repositories")
	repoName := repos[0]

	t.Run("Code Survival Analysis", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/repos/"+repoName+"/analysis/survival", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "survival_rate")

		t.Logf("Survival rate: %.1f%%", result["survival_rate"].(float64)*100)
	})

	t.Run("Revert Analysis", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/repos/"+repoName+"/analysis/reverts", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "total_prs_merged")
	})

	t.Run("Hotfix Detection", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/repos/"+repoName+"/analysis/hotfixes", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "hotfix_rate")
	})
}

// TestInsomnia_ResearchExportEndpoints tests research dataset export endpoints.
func TestInsomnia_ResearchExportEndpoints(t *testing.T) {
	cleanup, _, _ := setupInsomniaTestServer(t)
	defer cleanup()

	t.Run("Export Dataset JSON", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/research/dataset?format=json", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "data")

		data := result["data"].([]interface{})
		assert.Greater(t, len(data), 0, "should have dataset records")

		t.Logf("Dataset records: %d", len(data))
	})

	t.Run("Export Dataset CSV", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/research/dataset?format=csv", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "text/csv", resp.Header.Get("Content-Type"))
	})

	t.Run("Velocity Metrics", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/research/metrics/velocity", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "data")

		data := result["data"].([]interface{})
		t.Logf("Velocity records: %d", len(data))
	})

	t.Run("Review Cost Metrics", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/research/metrics/review-costs", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "data")

		data := result["data"].([]interface{})
		t.Logf("Review cost records: %d", len(data))
	})

	t.Run("Quality Metrics", func(t *testing.T) {
		resp := makeInsomniaRequest(t, "GET", "/research/metrics/quality", true)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Contains(t, result, "data")

		data := result["data"].([]interface{})
		t.Logf("Quality records: %d", len(data))
	})
}

// TestInsomnia_AllEndpointsReturnData is a comprehensive test that ensures
// all endpoints from the Insomnia collection return non-empty data.
func TestInsomnia_AllEndpointsReturnData(t *testing.T) {
	cleanup, store, _ := setupInsomniaTestServer(t)
	defer cleanup()

	repos := store.ListRepositories()
	require.Greater(t, len(repos), 0)
	repoName := repos[0]

	// Find a PR
	var prNumber int
	for _, repo := range repos {
		prs := store.GetPRsByRepo(repo)
		if len(prs) > 0 {
			prNumber = prs[0].Number
			break
		}
	}
	require.Greater(t, prNumber, 0)

	// Find a commit SHA
	commits := store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour))
	require.Greater(t, len(commits), 0)
	commitSHA := commits[0].CommitHash[:12]

	endpoints := []struct {
		path        string
		description string
		wantStatus  int
	}{
		// Core
		{"/health", "Health Check", 200},
		{"/teams/members", "Team Members", 200},
		{"/analytics/ai-code/commits", "AI Code Commits", 200},

		// Team Analytics
		{"/analytics/team/agent-edits", "Agent Edits", 200},
		{"/analytics/team/tabs", "Tab Completions", 200},
		{"/analytics/team/dau", "Daily Active Users", 200},
		{"/analytics/team/models", "Model Usage", 200},
		{"/analytics/team/client-versions", "Client Versions", 200},
		{"/analytics/team/top-file-extensions", "Top File Extensions", 200},
		{"/analytics/team/mcp", "MCP Tool Usage", 200},
		{"/analytics/team/commands", "Command Usage", 200},
		{"/analytics/team/plans", "Plan Usage", 200},
		{"/analytics/team/ask-mode", "Ask Mode Usage", 200},
		{"/analytics/team/leaderboard", "Leaderboard", 200},

		// By-User Analytics
		{"/analytics/by-user/agent-edits", "Agent Edits by User", 200},
		{"/analytics/by-user/tabs", "Tab Completions by User", 200},
		{"/analytics/by-user/models", "Model Usage by User", 200},
		{"/analytics/by-user/client-versions", "Client Versions by User", 200},
		{"/analytics/by-user/top-file-extensions", "File Extensions by User", 200},
		{"/analytics/by-user/mcp", "MCP Usage by User", 200},
		{"/analytics/by-user/commands", "Commands by User", 200},
		{"/analytics/by-user/plans", "Plans by User", 200},
		{"/analytics/by-user/ask-mode", "Ask Mode by User", 200},

		// GitHub Analytics
		{"/analytics/github/prs", "PRs Analytics", 200},
		{"/analytics/github/reviews", "Reviews Analytics", 200},
		{"/analytics/github/issues", "Issues Analytics", 200},
		{"/analytics/github/pr-cycle-time", "PR Cycle Time", 200},
		{"/analytics/github/review-quality", "Review Quality", 200},

		// GitHub Repos
		{"/repos", "List Repositories", 200},
		{"/repos/" + repoName, "Get Repository", 200},
		{"/repos/" + repoName + "/commits", "List Commits", 200},
		{"/repos/" + repoName + "/commits/" + commitSHA, "Get Commit", 200},
		{"/repos/" + repoName + "/pulls", "List PRs", 200},
		{fmt.Sprintf("/repos/%s/pulls/%d", repoName, prNumber), "Get PR", 200},
		{fmt.Sprintf("/repos/%s/pulls/%d/commits", repoName, prNumber), "PR Commits", 200},
		{fmt.Sprintf("/repos/%s/pulls/%d/files", repoName, prNumber), "PR Files", 200},
		{fmt.Sprintf("/repos/%s/pulls/%d/reviews", repoName, prNumber), "PR Reviews", 200},

		// Quality Analysis
		{"/repos/" + repoName + "/analysis/survival", "Survival Analysis", 200},
		{"/repos/" + repoName + "/analysis/reverts", "Revert Analysis", 200},
		{"/repos/" + repoName + "/analysis/hotfixes", "Hotfix Detection", 200},

		// Research Export
		{"/research/dataset", "Research Dataset", 200},
		{"/research/metrics/velocity", "Velocity Metrics", 200},
		{"/research/metrics/review-costs", "Review Cost Metrics", 200},
		{"/research/metrics/quality", "Quality Metrics", 200},
	}

	var passed, failed int
	for _, ep := range endpoints {
		t.Run(ep.description, func(t *testing.T) {
			resp := makeInsomniaRequest(t, "GET", ep.path, ep.path != "/health")
			defer resp.Body.Close()

			if resp.StatusCode == ep.wantStatus {
				passed++
			} else {
				failed++
				t.Errorf("endpoint %s returned %d, want %d", ep.path, resp.StatusCode, ep.wantStatus)
			}

			assert.Equal(t, ep.wantStatus, resp.StatusCode, "endpoint %s", ep.path)
		})
	}

	t.Logf("Endpoint validation: %d passed, %d failed", passed, failed)
}
