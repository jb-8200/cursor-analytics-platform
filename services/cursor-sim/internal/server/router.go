package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/cursor"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/github"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/harvey"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/microsoft"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/qualtrics"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/research"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/services"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// NewRouter creates and configures the HTTP router with all endpoints and middleware.
func NewRouter(store storage.Store, seedData interface{}, apiKey string) http.Handler {
	mux := http.NewServeMux()

	// Health check (no auth required)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		api.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Team Management API
	mux.Handle("/teams/members", cursor.TeamMembers(store))

	// AI Code Tracking API
	mux.Handle("/analytics/ai-code/commits", cursor.AICodeCommits(store))
	mux.Handle("/analytics/ai-code/commits.csv", cursor.AICodeCommitsCSV(store))

	// Team Analytics API (11 endpoints)
	mux.Handle("/analytics/team/agent-edits", cursor.TeamAgentEdits(store))
	mux.Handle("/analytics/team/tabs", cursor.TeamTabs(store))
	mux.Handle("/analytics/team/dau", cursor.TeamDAU(store))
	mux.Handle("/analytics/team/models", cursor.TeamModels(store))
	mux.Handle("/analytics/team/client-versions", cursor.TeamClientVersions(store))
	mux.Handle("/analytics/team/top-file-extensions", cursor.TeamTopFileExtensions(store))
	mux.Handle("/analytics/team/mcp", cursor.TeamMCP(store))
	mux.Handle("/analytics/team/commands", cursor.TeamCommands(store))
	mux.Handle("/analytics/team/plans", cursor.TeamPlans(store))
	mux.Handle("/analytics/team/ask-mode", cursor.TeamAskMode(store))
	mux.Handle("/analytics/team/leaderboard", cursor.TeamLeaderboard(store))

	// By-User Analytics API (9 endpoints)
	mux.Handle("/analytics/by-user/agent-edits", cursor.ByUserAgentEdits(store))
	mux.Handle("/analytics/by-user/tabs", cursor.ByUserTabs(store))
	mux.Handle("/analytics/by-user/models", cursor.ByUserModels(store))
	mux.Handle("/analytics/by-user/client-versions", cursor.ByUserClientVersions(store))
	mux.Handle("/analytics/by-user/top-file-extensions", cursor.ByUserTopFileExtensions(store))
	mux.Handle("/analytics/by-user/mcp", cursor.ByUserMCP(store))
	mux.Handle("/analytics/by-user/commands", cursor.ByUserCommands(store))
	mux.Handle("/analytics/by-user/plans", cursor.ByUserPlans(store))
	mux.Handle("/analytics/by-user/ask-mode", cursor.ByUserAskMode(store))

	// GitHub Simulation API (12 endpoints)
	mux.Handle("/repos", github.ListRepos(store))
	mux.Handle("/repos/", github.RepoRouter(store))

	// GitHub Analytics API (P2-F01)
	mux.Handle("/analytics/github/prs", github.ListPRsAnalytics(store))
	mux.Handle("/analytics/github/reviews", github.ListReviewsAnalytics(store))
	mux.Handle("/analytics/github/issues", github.ListIssuesAnalytics(store))
	mux.Handle("/analytics/github/pr-cycle-time", github.PRCycleTimeAnalytics(store))
	mux.Handle("/analytics/github/review-quality", github.ReviewQualityAnalytics(store))

	// Research API (5 endpoints)
	// Create research generator from seed data
	researchGen := generator.NewResearchGenerator(seedData.(*seed.SeedData), store)
	mux.Handle("/research/dataset", research.DatasetHandler(researchGen))
	mux.Handle("/research/metrics/velocity", research.VelocityMetricsHandler(researchGen))
	mux.Handle("/research/metrics/review-costs", research.ReviewCostMetricsHandler(researchGen))
	mux.Handle("/research/metrics/quality", research.QualityMetricsHandler(researchGen))

	// Harvey API (External Data Source - P4-F04)
	// Only register if Harvey is enabled in seed data
	sd := seedData.(*seed.SeedData)
	if sd.ExternalDataSources != nil && sd.ExternalDataSources.Harvey != nil && sd.ExternalDataSources.Harvey.Enabled {
		// Create external memory store
		externalStore := storage.NewExternalMemoryStore()
		mux.Handle("/harvey/api/v1/history/usage", harvey.UsageHandler(externalStore.Harvey()))
	}

	// Microsoft 365 Copilot API (External Data Source - P4-F04)
	// Only register if Copilot is enabled in seed data
	if sd.ExternalDataSources != nil && sd.ExternalDataSources.Copilot != nil && sd.ExternalDataSources.Copilot.Enabled {
		// Create external memory store for Copilot
		externalStore := storage.NewExternalMemoryStore()
		// Initialize Copilot generator from seed data
		copilotGen := generator.NewCopilotGenerator(sd)
		// Register Copilot routes under /reports/ prefix
		// Pattern: /reports/getMicrosoft365CopilotUsageUserDetail(period='D30')
		// Note: We use a custom handler wrapper to validate the specific endpoint
		copilotHandler := microsoft.CopilotUsageHandler(externalStore.Copilot(), copilotGen)
		mux.HandleFunc("/reports/", func(w http.ResponseWriter, r *http.Request) {
			// Only handle Copilot endpoint, return 404 for others
			const copilotPrefix = "/reports/getMicrosoft365CopilotUsageUserDetail"
			if r.URL.Path == copilotPrefix || (len(r.URL.Path) > len(copilotPrefix) && r.URL.Path[:len(copilotPrefix)+1] == copilotPrefix+"(") {
				copilotHandler.ServeHTTP(w, r)
				return
			}
			http.NotFound(w, r)
		})
	}

	// Qualtrics API (External Data Source - P4-F04)
	// Only register if Qualtrics is enabled in seed data
	if sd.ExternalDataSources != nil && sd.ExternalDataSources.Qualtrics != nil && sd.ExternalDataSources.Qualtrics.Enabled {
		// Initialize survey generator from seed data
		surveyGen := generator.NewSurveyGenerator(sd)
		// Create export job manager
		exportManager := services.NewExportJobManager(surveyGen)
		// Create handlers
		qualtricsHandlers := qualtrics.NewExportHandlers(exportManager)

		// Register routes under /API/v3/surveys/ prefix
		// Pattern:
		// - POST /API/v3/surveys/{surveyId}/export-responses - Start export
		// - GET /API/v3/surveys/{surveyId}/export-responses/{progressId} - Check progress
		// - GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file - Download file
		mux.HandleFunc("/API/v3/surveys/", func(w http.ResponseWriter, r *http.Request) {
			// Extract path components
			path := r.URL.Path
			const prefix = "/API/v3/surveys/"

			// Check if path starts with the prefix
			if !strings.HasPrefix(path, prefix) {
				http.NotFound(w, r)
				return
			}

			// Remove prefix to get the remainder
			remainder := strings.TrimPrefix(path, prefix)
			parts := strings.Split(remainder, "/")

			// Validate path structure: {surveyId}/export-responses/...
			if len(parts) < 2 || parts[1] != "export-responses" {
				http.NotFound(w, r)
				return
			}

			// Route based on method and path structure
			if r.Method == http.MethodPost && len(parts) == 2 {
				// POST /API/v3/surveys/{surveyId}/export-responses
				qualtricsHandlers.StartExportHandler().ServeHTTP(w, r)
				return
			}

			if r.Method == http.MethodGet && len(parts) == 3 {
				// GET /API/v3/surveys/{surveyId}/export-responses/{progressId}
				qualtricsHandlers.ProgressHandler().ServeHTTP(w, r)
				return
			}

			if r.Method == http.MethodGet && len(parts) == 4 && parts[3] == "file" {
				// GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file
				qualtricsHandlers.FileDownloadHandler().ServeHTTP(w, r)
				return
			}

			// No matching route
			http.NotFound(w, r)
		})
	}

	// Apply middleware (reverse order: Logger wraps RateLimit wraps BasicAuth wraps mux)
	limiter := api.NewRateLimiter(100, time.Minute)

	// Create auth-protected handler
	authHandler := authProtectedRoutes(mux, apiKey)

	// Apply rate limiting
	rateLimitedHandler := api.RateLimit(limiter)(authHandler)

	// Apply logging
	return api.Logger(rateLimitedHandler)
}

// authProtectedRoutes applies BasicAuth to all routes except /health.
func authProtectedRoutes(handler http.Handler, apiKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health endpoint
		if r.URL.Path == "/health" {
			handler.ServeHTTP(w, r)
			return
		}

		// Apply BasicAuth for all other routes
		api.BasicAuth(apiKey)(handler).ServeHTTP(w, r)
	})
}
