package server

import (
	"net/http"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/cursor"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/github"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/research"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
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

	// Research API (5 endpoints)
	// Create research generator from seed data
	researchGen := generator.NewResearchGenerator(seedData.(*seed.SeedData), store)
	mux.Handle("/research/dataset", research.DatasetHandler(researchGen))
	mux.Handle("/research/metrics/velocity", research.VelocityMetricsHandler(researchGen))
	mux.Handle("/research/metrics/review-costs", research.ReviewCostMetricsHandler(researchGen))
	mux.Handle("/research/metrics/quality", research.QualityMetricsHandler(researchGen))

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
