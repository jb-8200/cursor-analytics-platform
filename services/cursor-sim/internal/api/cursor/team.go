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

// TeamModels returns handler for GET /analytics/team/models.
// Aggregates model usage by day with breakdown by model.
func TeamModels(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse date range
		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)
		to = to.Add(24*time.Hour - time.Second)

		// Get model usage events in range
		events := store.GetModelUsageByTimeRange(from, to)

		// Aggregate by date and model
		dayMap := make(map[string]map[string]map[string]bool) // date -> model -> users

		for _, event := range events {
			date := event.EventDate
			if _, exists := dayMap[date]; !exists {
				dayMap[date] = make(map[string]map[string]bool)
			}
			if _, exists := dayMap[date][event.ModelName]; !exists {
				dayMap[date][event.ModelName] = make(map[string]bool)
			}
			dayMap[date][event.ModelName][event.UserID] = true
		}

		// Count messages per model per day (approximate: 1 event = 1 message)
		messageCountMap := make(map[string]map[string]int) // date -> model -> count
		for _, event := range events {
			date := event.EventDate
			if _, exists := messageCountMap[date]; !exists {
				messageCountMap[date] = make(map[string]int)
			}
			messageCountMap[date][event.ModelName]++
		}

		// Build response
		result := make([]models.ModelUsageDay, 0, len(dayMap))
		for date, modelData := range dayMap {
			breakdown := make(map[string]models.ModelBreakdownItem)
			for model, users := range modelData {
				breakdown[model] = models.ModelBreakdownItem{
					Messages: messageCountMap[date][model],
					Users:    len(users),
				}
			}
			result = append(result, models.ModelUsageDay{
				Date:           date,
				ModelBreakdown: breakdown,
			})
		}

		// Build response using Analytics API format
		response := api.BuildAnalyticsTeamResponse(result, "models", params)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// TeamClientVersions returns handler for GET /analytics/team/client-versions.
// Aggregates client version distribution by day.
func TeamClientVersions(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse date range
		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)
		to = to.Add(24*time.Hour - time.Second)

		// Get client version events in range
		events := store.GetClientVersionsByTimeRange(from, to)

		// Aggregate by date and version
		// Map: date -> version -> set of user IDs
		dayVersionUsers := make(map[string]map[string]map[string]bool)

		for _, event := range events {
			date := event.EventDate
			if _, exists := dayVersionUsers[date]; !exists {
				dayVersionUsers[date] = make(map[string]map[string]bool)
			}
			if _, exists := dayVersionUsers[date][event.ClientVersion]; !exists {
				dayVersionUsers[date][event.ClientVersion] = make(map[string]bool)
			}
			dayVersionUsers[date][event.ClientVersion][event.UserID] = true
		}

		// Build response data
		result := make([]models.ClientVersionDay, 0)
		for date, versionUsers := range dayVersionUsers {
			// Calculate total users for this date
			totalUsers := make(map[string]bool)
			for _, users := range versionUsers {
				for userID := range users {
					totalUsers[userID] = true
				}
			}
			totalCount := len(totalUsers)

			// Create entries for each version
			for version, users := range versionUsers {
				userCount := len(users)
				percentage := 0.0
				if totalCount > 0 {
					percentage = float64(userCount) / float64(totalCount)
				}

				result = append(result, models.ClientVersionDay{
					EventDate:     date,
					ClientVersion: version,
					UserCount:     userCount,
					Percentage:    percentage,
				})
			}
		}

		// Build response using Analytics API format
		response := api.BuildAnalyticsTeamResponse(result, "client-versions", params)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// TeamTopFileExtensions returns handler for GET /analytics/team/top-file-extensions.
// Aggregates file extension usage by day, showing top 5 by suggestion volume.
func TeamTopFileExtensions(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse date range
		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)
		to = to.Add(24*time.Hour - time.Second)

		// Get file extension events in range
		events := store.GetFileExtensionsByTimeRange(from, to)

		// Aggregate by date and extension
		// Map: date -> extension -> aggregated stats
		dayExtStats := make(map[string]map[string]*models.FileExtensionDay)

		for _, event := range events {
			date := event.EventDate
			if _, exists := dayExtStats[date]; !exists {
				dayExtStats[date] = make(map[string]*models.FileExtensionDay)
			}

			if _, exists := dayExtStats[date][event.FileExtension]; !exists {
				dayExtStats[date][event.FileExtension] = &models.FileExtensionDay{
					EventDate:     date,
					FileExtension: event.FileExtension,
				}
			}

			// Aggregate stats
			stat := dayExtStats[date][event.FileExtension]
			stat.TotalFiles++
			stat.TotalAccepts += event.LinesAccepted
			stat.TotalRejects += event.LinesRejected
			stat.TotalLinesSuggested += event.LinesSuggested
			stat.TotalLinesAccepted += event.LinesAccepted
			stat.TotalLinesRejected += event.LinesRejected
		}

		// Build response: top 5 extensions per day by suggestion volume
		result := make([]models.FileExtensionDay, 0)
		for _, extStats := range dayExtStats {
			// Sort by total lines suggested (descending)
			var extensions []*models.FileExtensionDay
			for _, stat := range extStats {
				extensions = append(extensions, stat)
			}

			// Sort by lines suggested
			for i := 0; i < len(extensions); i++ {
				for j := i + 1; j < len(extensions); j++ {
					if extensions[j].TotalLinesSuggested > extensions[i].TotalLinesSuggested {
						extensions[i], extensions[j] = extensions[j], extensions[i]
					}
				}
			}

			// Take top 5
			limit := 5
			if len(extensions) < limit {
				limit = len(extensions)
			}

			for i := 0; i < limit; i++ {
				result = append(result, *extensions[i])
			}
		}

		// Build response using Analytics API format
		response := api.BuildAnalyticsTeamResponse(result, "top-file-extensions", params)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// TeamMCP returns handler for GET /analytics/team/mcp.
// Aggregates MCP tool usage by day.
func TeamMCP(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)
		to = to.Add(24*time.Hour - time.Second)

		events := store.GetMCPToolsByTimeRange(from, to)
		dayToolUsage := make(map[string]map[string]int) // date -> (toolName + server) -> count

		for _, event := range events {
			key := event.ToolName + ":" + event.MCPServerName
			if _, exists := dayToolUsage[event.EventDate]; !exists {
				dayToolUsage[event.EventDate] = make(map[string]int)
			}
			dayToolUsage[event.EventDate][key]++
		}

		result := make([]models.MCPUsageDay, 0)
		for date, tools := range dayToolUsage {
			for key, usage := range tools {
				parts := splitKey(key)
				result = append(result, models.MCPUsageDay{
					EventDate:     date,
					ToolName:      parts[0],
					MCPServerName: parts[1],
					Usage:         usage,
				})
			}
		}

		response := api.BuildAnalyticsTeamResponse(result, "mcp", params)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// TeamCommands returns handler for GET /analytics/team/commands.
// Aggregates command usage by day.
func TeamCommands(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)
		to = to.Add(24*time.Hour - time.Second)

		events := store.GetCommandsByTimeRange(from, to)
		dayCommandUsage := make(map[string]map[string]int) // date -> command -> count

		for _, event := range events {
			if _, exists := dayCommandUsage[event.EventDate]; !exists {
				dayCommandUsage[event.EventDate] = make(map[string]int)
			}
			dayCommandUsage[event.EventDate][event.CommandName]++
		}

		result := make([]models.CommandUsageDay, 0)
		for date, commands := range dayCommandUsage {
			for cmd, usage := range commands {
				result = append(result, models.CommandUsageDay{
					EventDate:   date,
					CommandName: cmd,
					Usage:       usage,
				})
			}
		}

		response := api.BuildAnalyticsTeamResponse(result, "commands", params)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// TeamPlans returns handler for GET /analytics/team/plans.
// Aggregates plan usage by day and model.
func TeamPlans(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)
		to = to.Add(24*time.Hour - time.Second)

		events := store.GetPlansByTimeRange(from, to)
		dayModelUsage := make(map[string]map[string]int) // date -> model -> count

		for _, event := range events {
			if _, exists := dayModelUsage[event.EventDate]; !exists {
				dayModelUsage[event.EventDate] = make(map[string]int)
			}
			dayModelUsage[event.EventDate][event.Model]++
		}

		result := make([]models.PlanUsageDay, 0)
		for date, modelUsage := range dayModelUsage {
			for model, usage := range modelUsage {
				result = append(result, models.PlanUsageDay{
					EventDate: date,
					Model:     model,
					Usage:     usage,
				})
			}
		}

		response := api.BuildAnalyticsTeamResponse(result, "plans", params)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// TeamAskMode returns handler for GET /analytics/team/ask-mode.
// Aggregates ask mode usage by day and model.
func TeamAskMode(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)
		to = to.Add(24*time.Hour - time.Second)

		events := store.GetAskModeByTimeRange(from, to)
		dayModelUsage := make(map[string]map[string]int) // date -> model -> count

		for _, event := range events {
			if _, exists := dayModelUsage[event.EventDate]; !exists {
				dayModelUsage[event.EventDate] = make(map[string]int)
			}
			dayModelUsage[event.EventDate][event.Model]++
		}

		result := make([]models.AskModeDay, 0)
		for date, modelUsage := range dayModelUsage {
			for model, usage := range modelUsage {
				result = append(result, models.AskModeDay{
					EventDate: date,
					Model:     model,
					Usage:     usage,
				})
			}
		}

		response := api.BuildAnalyticsTeamResponse(result, "ask-mode", params)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// Helper function to split key in format "tool:server"
func splitKey(key string) []string {
	for i := 0; i < len(key); i++ {
		if key[i] == ':' {
			return []string{key[:i], key[i+1:]}
		}
	}
	return []string{key, ""}
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
