package cursor

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// ByUserAgentEdits returns handler for GET /analytics/by-user/agent-edits.
// Returns agent edits broken down by user.
func ByUserAgentEdits(store storage.Store) http.Handler {
	return byUserHandler(store, "agent-edits", func(commits []models.Commit) map[string]interface{} {
		// Group by user email -> date -> metrics
		userdata := make(map[string]map[string]*models.AgentEditDay)

		for _, commit := range commits {
			email := commit.UserEmail
			if email == "" {
				continue
			}

			if userdata[email] == nil {
				userdata[email] = make(map[string]*models.AgentEditDay)
			}

			date := commit.CommitTs.Format("2006-01-02")
			if userdata[email][date] == nil {
				userdata[email][date] = &models.AgentEditDay{EventDate: date}
			}

			userdata[email][date].SuggestedLines += commit.ComposerLinesAdded
			userdata[email][date].AcceptedLines += commit.ComposerLinesAdded // All suggested lines are accepted in current model
		}

		// Convert to arrays sorted by date
		result := make(map[string]interface{})
		for email, dates := range userdata {
			var days []models.AgentEditDay
			for _, day := range dates {
				days = append(days, *day)
			}
			sort.Slice(days, func(i, j int) bool {
				return days[i].EventDate < days[j].EventDate
			})
			result[email] = days
		}

		return result
	})
}

// ByUserTabs returns handler for GET /analytics/by-user/tabs.
// Returns tab completion metrics by user.
func ByUserTabs(store storage.Store) http.Handler {
	return byUserHandler(store, "tabs", func(commits []models.Commit) map[string]interface{} {
		// Group by user email -> date -> metrics
		userdata := make(map[string]map[string]*models.TabDay)

		for _, commit := range commits {
			email := commit.UserEmail
			if email == "" {
				continue
			}

			if userdata[email] == nil {
				userdata[email] = make(map[string]*models.TabDay)
			}

			date := commit.CommitTs.Format("2006-01-02")
			if userdata[email][date] == nil {
				userdata[email][date] = &models.TabDay{EventDate: date}
			}

			userdata[email][date].SuggestedLines += commit.TabLinesAdded
			userdata[email][date].AcceptedLines += commit.TabLinesAdded
		}

		// Convert to arrays sorted by date
		result := make(map[string]interface{})
		for email, dates := range userdata {
			var days []models.TabDay
			for _, day := range dates {
				days = append(days, *day)
			}
			sort.Slice(days, func(i, j int) bool {
				return days[i].EventDate < days[j].EventDate
			})
			result[email] = days
		}

		return result
	})
}

// ByUserModels returns handler for GET /analytics/by-user/models.
func ByUserModels(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)
		to = to.Add(24*time.Hour - time.Second)

		events := store.GetModelUsageByTimeRange(from, to)

		// Group by user -> date -> model -> count
		userdata := make(map[string]map[string]map[string]int)

		for _, event := range events {
			email := event.UserEmail
			if email == "" {
				continue
			}

			if userdata[email] == nil {
				userdata[email] = make(map[string]map[string]int)
			}

			date := event.Timestamp.Format("2006-01-02")
			if userdata[email][date] == nil {
				userdata[email][date] = make(map[string]int)
			}

			userdata[email][date][event.ModelName]++
		}

		// Convert to response format
		result := make(map[string]interface{})
		for email, dates := range userdata {
			var days []models.ModelDay
			for date, modelCounts := range dates {
				var breakdown []models.ModelBreakdown
				for model, count := range modelCounts {
					breakdown = append(breakdown, models.ModelBreakdown{
						Model:        model,
						MessageCount: count,
					})
				}
				sort.Slice(breakdown, func(i, j int) bool {
					return breakdown[i].Model < breakdown[j].Model
				})
				days = append(days, models.ModelDay{
					EventDate: date,
					Breakdown: breakdown,
				})
			}
			sort.Slice(days, func(i, j int) bool {
				return days[i].EventDate < days[j].EventDate
			})
			result[email] = days
		}

		// Apply user filtering and pagination
		response := buildByUserResponse(result, "models", params, store)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// ByUserClientVersions returns handler for GET /analytics/by-user/client-versions.
func ByUserClientVersions(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)
		to = to.Add(24*time.Hour - time.Second)

		events := store.GetClientVersionsByTimeRange(from, to)

		// Group by user -> date -> version
		userdata := make(map[string]map[string]string)

		for _, event := range events {
			email := event.UserEmail
			if email == "" {
				continue
			}

			if userdata[email] == nil {
				userdata[email] = make(map[string]string)
			}

			date := event.Timestamp.Format("2006-01-02")
			// Store the latest version seen on this date
			userdata[email][date] = event.ClientVersion
		}

		// Convert to response format
		result := make(map[string]interface{})
		for email, dates := range userdata {
			var days []models.ByUserClientVersionDay
			for date, version := range dates {
				days = append(days, models.ByUserClientVersionDay{
					EventDate: date,
					Version:   version,
				})
			}
			sort.Slice(days, func(i, j int) bool {
				return days[i].EventDate < days[j].EventDate
			})
			result[email] = days
		}

		// Apply user filtering and pagination
		response := buildByUserResponse(result, "client-versions", params, store)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// ByUserTopFileExtensions returns handler for GET /analytics/by-user/top-file-extensions.
func ByUserTopFileExtensions(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)
		to = to.Add(24*time.Hour - time.Second)

		events := store.GetFileExtensionsByTimeRange(from, to)

		// Group by user -> date -> extension -> metrics
		userdata := make(map[string]map[string]map[string]*models.ExtensionMetrics)

		for _, event := range events {
			email := event.UserEmail
			if email == "" {
				continue
			}

			if userdata[email] == nil {
				userdata[email] = make(map[string]map[string]*models.ExtensionMetrics)
			}

			date := event.Timestamp.Format("2006-01-02")
			if userdata[email][date] == nil {
				userdata[email][date] = make(map[string]*models.ExtensionMetrics)
			}

			if userdata[email][date][event.FileExtension] == nil {
				userdata[email][date][event.FileExtension] = &models.ExtensionMetrics{}
			}

			userdata[email][date][event.FileExtension].SuggestedLines += event.LinesSuggested
			userdata[email][date][event.FileExtension].AcceptedLines += event.LinesAccepted
		}

		// Convert to response format (top 5 per day)
		result := make(map[string]interface{})
		for email, dates := range userdata {
			var days []models.ByUserFileExtensionDay
			for date, extMap := range dates {
				// Convert to slice and sort by suggested lines
				var exts []models.ExtensionStats
				for ext, metrics := range extMap {
					exts = append(exts, models.ExtensionStats{
						Extension:      ext,
						SuggestedLines: metrics.SuggestedLines,
						AcceptedLines:  metrics.AcceptedLines,
					})
				}
				sort.Slice(exts, func(i, j int) bool {
					return exts[i].SuggestedLines > exts[j].SuggestedLines
				})

				// Take top 5
				if len(exts) > 5 {
					exts = exts[:5]
				}

				days = append(days, models.ByUserFileExtensionDay{
					EventDate:  date,
					Extensions: exts,
				})
			}
			sort.Slice(days, func(i, j int) bool {
				return days[i].EventDate < days[j].EventDate
			})
			result[email] = days
		}

		// Apply user filtering and pagination
		response := buildByUserResponse(result, "top-file-extensions", params, store)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// ByUserMCP returns handler for GET /analytics/by-user/mcp.
func ByUserMCP(store storage.Store) http.Handler {
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

		// Group by user -> date -> tool -> count
		userdata := make(map[string]map[string]map[string]int)

		for _, event := range events {
			email := event.UserEmail
			if email == "" {
				continue
			}

			if userdata[email] == nil {
				userdata[email] = make(map[string]map[string]int)
			}

			date := event.Timestamp.Format("2006-01-02")
			if userdata[email][date] == nil {
				userdata[email][date] = make(map[string]int)
			}

			key := event.ToolName + ":" + event.MCPServerName
			userdata[email][date][key]++
		}

		// Convert to response format
		result := make(map[string]interface{})
		for email, dates := range userdata {
			var days []models.MCPToolDay
			for date, toolCounts := range dates {
				var tools []models.MCPToolUsage
				for key, count := range toolCounts {
					parts := splitKey(key)
					tools = append(tools, models.MCPToolUsage{
						ToolName:      parts[0],
						MCPServerName: parts[1],
						UsageCount:    count,
					})
				}
				sort.Slice(tools, func(i, j int) bool {
					if tools[i].ToolName != tools[j].ToolName {
						return tools[i].ToolName < tools[j].ToolName
					}
					return tools[i].MCPServerName < tools[j].MCPServerName
				})
				days = append(days, models.MCPToolDay{
					EventDate: date,
					Tools:     tools,
				})
			}
			sort.Slice(days, func(i, j int) bool {
				return days[i].EventDate < days[j].EventDate
			})
			result[email] = days
		}

		// Apply user filtering and pagination
		response := buildByUserResponse(result, "mcp", params, store)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// ByUserCommands returns handler for GET /analytics/by-user/commands.
func ByUserCommands(store storage.Store) http.Handler {
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

		// Group by user -> date -> command -> count
		userdata := make(map[string]map[string]map[string]int)

		for _, event := range events {
			email := event.UserEmail
			if email == "" {
				continue
			}

			if userdata[email] == nil {
				userdata[email] = make(map[string]map[string]int)
			}

			date := event.Timestamp.Format("2006-01-02")
			if userdata[email][date] == nil {
				userdata[email][date] = make(map[string]int)
			}

			userdata[email][date][event.CommandName]++
		}

		// Convert to response format
		result := make(map[string]interface{})
		for email, dates := range userdata {
			var days []models.CommandDay
			for date, cmdCounts := range dates {
				var commands []models.CommandUsage
				for cmd, count := range cmdCounts {
					commands = append(commands, models.CommandUsage{
						CommandName: cmd,
						UsageCount:  count,
					})
				}
				sort.Slice(commands, func(i, j int) bool {
					return commands[i].CommandName < commands[j].CommandName
				})
				days = append(days, models.CommandDay{
					EventDate: date,
					Commands:  commands,
				})
			}
			sort.Slice(days, func(i, j int) bool {
				return days[i].EventDate < days[j].EventDate
			})
			result[email] = days
		}

		// Apply user filtering and pagination
		response := buildByUserResponse(result, "commands", params, store)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// ByUserPlans returns handler for GET /analytics/by-user/plans.
func ByUserPlans(store storage.Store) http.Handler {
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

		// Group by user -> date -> model -> count
		userdata := make(map[string]map[string]map[string]int)

		for _, event := range events {
			email := event.UserEmail
			if email == "" {
				continue
			}

			if userdata[email] == nil {
				userdata[email] = make(map[string]map[string]int)
			}

			date := event.Timestamp.Format("2006-01-02")
			if userdata[email][date] == nil {
				userdata[email][date] = make(map[string]int)
			}

			userdata[email][date][event.Model]++
		}

		// Convert to response format
		result := make(map[string]interface{})
		for email, dates := range userdata {
			var days []models.ByUserPlanDay
			for date, modelCounts := range dates {
				var modelUsage []models.PlanModelUsage
				for model, count := range modelCounts {
					modelUsage = append(modelUsage, models.PlanModelUsage{
						Model:      model,
						UsageCount: count,
					})
				}
				sort.Slice(modelUsage, func(i, j int) bool {
					return modelUsage[i].Model < modelUsage[j].Model
				})
				days = append(days, models.ByUserPlanDay{
					EventDate: date,
					Models:    modelUsage,
				})
			}
			sort.Slice(days, func(i, j int) bool {
				return days[i].EventDate < days[j].EventDate
			})
			result[email] = days
		}

		// Apply user filtering and pagination
		response := buildByUserResponse(result, "plans", params, store)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// ByUserAskMode returns handler for GET /analytics/by-user/ask-mode.
func ByUserAskMode(store storage.Store) http.Handler {
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

		// Group by user -> date -> model -> count
		userdata := make(map[string]map[string]map[string]int)

		for _, event := range events {
			email := event.UserEmail
			if email == "" {
				continue
			}

			if userdata[email] == nil {
				userdata[email] = make(map[string]map[string]int)
			}

			date := event.Timestamp.Format("2006-01-02")
			if userdata[email][date] == nil {
				userdata[email][date] = make(map[string]int)
			}

			userdata[email][date][event.Model]++
		}

		// Convert to response format
		result := make(map[string]interface{})
		for email, dates := range userdata {
			var days []models.ByUserAskModeDay
			for date, modelCounts := range dates {
				var modelUsage []models.AskModeModelUsage
				for model, count := range modelCounts {
					modelUsage = append(modelUsage, models.AskModeModelUsage{
						Model:      model,
						UsageCount: count,
					})
				}
				sort.Slice(modelUsage, func(i, j int) bool {
					return modelUsage[i].Model < modelUsage[j].Model
				})
				days = append(days, models.ByUserAskModeDay{
					EventDate: date,
					Models:    modelUsage,
				})
			}
			sort.Slice(days, func(i, j int) bool {
				return days[i].EventDate < days[j].EventDate
			})
			result[email] = days
		}

		// Apply user filtering and pagination
		response := buildByUserResponse(result, "ask-mode", params, store)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// byUserHandler creates a handler for by-user endpoints that use commit data.
func byUserHandler(store storage.Store, metric string, extract func([]models.Commit) map[string]interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)
		to = to.Add(24*time.Hour - time.Second)

		commits := store.GetCommitsByTimeRange(from, to)
		data := extract(commits)

		response := buildByUserResponse(data, metric, params, store)
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// buildByUserResponse applies user filtering, pagination, and builds the final response.
func buildByUserResponse(data map[string]interface{}, metric string, params models.Params, store storage.Store) models.AnalyticsByUserResponse {
	// Apply user filtering if provided
	var filteredData map[string]interface{}
	if params.User != "" {
		filterEmails := strings.Split(params.User, ",")
		filterMap := make(map[string]bool)
		for _, email := range filterEmails {
			filterMap[strings.TrimSpace(email)] = true
		}

		filteredData = make(map[string]interface{})
		for email, userData := range data {
			if filterMap[email] {
				filteredData[email] = userData
			}
		}
	} else {
		filteredData = data
	}

	// Get sorted user emails for pagination
	var emails []string
	for email := range filteredData {
		emails = append(emails, email)
	}
	sort.Strings(emails)

	totalUsers := len(emails)

	// Apply pagination on users
	start := (params.Page - 1) * params.PageSize
	end := start + params.PageSize

	if start >= len(emails) {
		start = 0
		end = 0
	} else if end > len(emails) {
		end = len(emails)
	}

	paginatedEmails := emails[start:end]

	// Build paginated data
	paginatedData := make(map[string]interface{})
	for _, email := range paginatedEmails {
		paginatedData[email] = filteredData[email]
	}

	// Build user mappings
	var userMappings []models.UserMapping
	developers := store.ListDevelopers()
	for _, dev := range developers {
		if _, exists := paginatedData[dev.Email]; exists {
			userMappings = append(userMappings, models.UserMapping{
				ID:    dev.UserID,
				Email: dev.Email,
			})
		}
	}

	return api.BuildAnalyticsByUserResponse(paginatedData, metric, params, totalUsers, userMappings)
}
