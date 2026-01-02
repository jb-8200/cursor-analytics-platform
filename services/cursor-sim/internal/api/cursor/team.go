package cursor

import (
	"net/http"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// TeamAgentEdits returns handler for GET /analytics/team/agent-edits.
// Aggregates AI-generated code edits by day.
func TeamAgentEdits(store storage.Store) http.Handler {
	return teamMetricHandler(store, "agent-edits", func(commits []models.Commit, params models.Params) interface{} {
		// Group by date and aggregate
		dayMap := make(map[string]*models.AgentEditsDay)

		for _, c := range commits {
			date := c.CommitTs.Format("2006-01-02")
			if _, exists := dayMap[date]; !exists {
				dayMap[date] = &models.AgentEditsDay{
					EventDate: date,
				}
			}

			// Aggregate edits (using commits as proxy for diffs)
			aiLines := c.TabLinesAdded + c.ComposerLinesAdded
			dayMap[date].TotalSuggestedDiffs++
			if aiLines > 0 {
				dayMap[date].TotalAcceptedDiffs++
				dayMap[date].TotalGreenLinesAccepted += aiLines
			}
		}

		// Convert to array
		result := make([]models.AgentEditsDay, 0, len(dayMap))
		for _, day := range dayMap {
			result = append(result, *day)
		}

		return result
	})
}

// TeamTabs returns handler for GET /analytics/team/tabs.
// Aggregates tab completion metrics by day.
func TeamTabs(store storage.Store) http.Handler {
	return teamMetricHandler(store, "tabs", func(commits []models.Commit, params models.Params) interface{} {
		dayMap := make(map[string]*models.TabUsageDay)

		for _, c := range commits {
			date := c.CommitTs.Format("2006-01-02")
			if _, exists := dayMap[date]; !exists {
				dayMap[date] = &models.TabUsageDay{
					EventDate: date,
				}
			}

			// Aggregate tab completions (using commit lines as proxy for tab completions)
			if c.TabLinesAdded > 0 {
				dayMap[date].TotalSuggestions++
				dayMap[date].TotalAccepts++
				dayMap[date].TotalGreenLinesAccepted += c.TabLinesAdded
				dayMap[date].TotalGreenLinesSuggested += c.TabLinesAdded
				dayMap[date].TotalLinesSuggested += c.TabLinesAdded
				dayMap[date].TotalLinesAccepted += c.TabLinesAdded
			}
		}

		result := make([]models.TabUsageDay, 0, len(dayMap))
		for _, day := range dayMap {
			result = append(result, *day)
		}

		return result
	})
}

// TeamDAU returns handler for GET /analytics/team/dau.
// Counts distinct active users per day.
func TeamDAU(store storage.Store) http.Handler {
	return teamMetricHandler(store, "dau", func(commits []models.Commit, params models.Params) interface{} {
		dayMap := make(map[string]map[string]bool)

		for _, c := range commits {
			date := c.CommitTs.Format("2006-01-02")
			if _, exists := dayMap[date]; !exists {
				dayMap[date] = make(map[string]bool)
			}
			dayMap[date][c.UserID] = true
		}

		result := make([]models.DAUDay, 0, len(dayMap))
		for date, users := range dayMap {
			result = append(result, models.DAUDay{
				Date: date, // Changed from EventDate
				DAU:  len(users), // Changed from UniqueUsers
				// CLI, CloudAgent, and Bugbot DAU are not tracked yet (will be 0)
			})
		}

		return result
	})
}

// Stub endpoints - return empty data for metrics we don't track yet

func TeamModels(store storage.Store) http.Handler {
	return stubHandler("models")
}

func TeamClientVersions(store storage.Store) http.Handler {
	return stubHandler("client-versions")
}

func TeamTopFileExtensions(store storage.Store) http.Handler {
	return stubHandler("top-file-extensions")
}

func TeamMCP(store storage.Store) http.Handler {
	return stubHandler("mcp")
}

func TeamCommands(store storage.Store) http.Handler {
	return stubHandler("commands")
}

func TeamPlans(store storage.Store) http.Handler {
	return stubHandler("plans")
}

func TeamAskMode(store storage.Store) http.Handler {
	return stubHandler("ask-mode")
}

func TeamLeaderboard(store storage.Store) http.Handler {
	return stubHandler("leaderboard")
}

// Helper functions

// teamMetricHandler creates a handler with common pattern for team metrics.
// Uses Analytics API team-level response format (no pagination wrapper).
// Reference: docs/api-reference/cursor_analytics.md
func teamMetricHandler(store storage.Store, metric string, extract func([]models.Commit, models.Params) interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters (startDate, endDate, etc.)
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse date range from validated params
		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)

		// Extend to include full day
		to = to.Add(24*time.Hour - time.Second)

		// Get commits in range
		commits := store.GetCommitsByTimeRange(from, to)

		// Extract metric-specific data
		data := extract(commits, params)

		// Build response using Analytics API format: { data: [...], params: {...} }
		response := api.BuildAnalyticsTeamResponse(data, metric, params)

		// Send JSON response
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// stubHandler returns a handler that responds with empty data for unimplemented team-level metrics.
// Uses Analytics API team-level response format (no pagination wrapper).
// Reference: docs/api-reference/cursor_analytics.md (Team-Level Endpoints)
func stubHandler(metric string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, _ := api.ParseQueryParams(r)
		response := api.BuildAnalyticsTeamResponse([]interface{}{}, metric, params)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// stubHandlerByUser returns a handler that responds with empty data for unimplemented by-user metrics.
// Uses Analytics API by-user response format (with pagination wrapper).
// Reference: docs/api-reference/cursor_analytics.md (By-User Endpoints)
func stubHandlerByUser(metric string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, _ := api.ParseQueryParams(r)
		response := api.BuildPaginatedResponse([]interface{}{}, params, 0)
		api.RespondJSON(w, http.StatusOK, response)
	})
}
