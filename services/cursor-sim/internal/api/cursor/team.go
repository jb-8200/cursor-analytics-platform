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
	return teamMetricHandler(store, func(commits []models.Commit, params models.Params) interface{} {
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
	return teamMetricHandler(store, func(commits []models.Commit, params models.Params) interface{} {
		dayMap := make(map[string]*models.TabCompletionDay)

		for _, c := range commits {
			date := c.CommitTs.Format("2006-01-02")
			if _, exists := dayMap[date]; !exists {
				dayMap[date] = &models.TabCompletionDay{
					EventDate: date,
				}
			}

			// Aggregate tab completions (lines as proxy for completions)
			if c.TabLinesAdded > 0 {
				dayMap[date].TotalSuggests++
				dayMap[date].TotalAccepts++
			}
		}

		result := make([]models.TabCompletionDay, 0, len(dayMap))
		for _, day := range dayMap {
			result = append(result, *day)
		}

		return result
	})
}

// TeamDAU returns handler for GET /analytics/team/dau.
// Counts distinct active users per day.
func TeamDAU(store storage.Store) http.Handler {
	return teamMetricHandler(store, func(commits []models.Commit, params models.Params) interface{} {
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
				EventDate:   date,
				UniqueUsers: len(users),
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
func teamMetricHandler(store storage.Store, extract func([]models.Commit, models.Params) interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse date range
		from, err := time.Parse("2006-01-02", params.From)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid from date")
			return
		}
		to, err := time.Parse("2006-01-02", params.To)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid to date")
			return
		}

		// Extend to include full day
		to = to.Add(24*time.Hour - time.Second)

		// Get commits in range
		commits := store.GetCommitsByTimeRange(from, to)

		// Extract metric-specific data
		data := extract(commits, params)

		// Build response
		response := api.BuildPaginatedResponse(data, params, 0)

		// Send JSON response
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// stubHandler returns a handler that responds with empty data for unimplemented metrics.
func stubHandler(metric string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, _ := api.ParseQueryParams(r)
		response := api.BuildPaginatedResponse([]interface{}{}, params, 0)
		api.RespondJSON(w, http.StatusOK, response)
	})
}
