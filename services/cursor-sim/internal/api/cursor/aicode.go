package cursor

import (
	"net/http"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// AICodeCommits returns an HTTP handler for GET /analytics/ai-code/commits.
// It retrieves commits from storage with time range and user filtering.
// Response format matches the Cursor API CommitsResponse schema.
func AICodeCommits(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters (startDate, endDate, user, page, pageSize)
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse date range from already-validated params
		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)

		// Add time to include full day
		to = to.Add(24*time.Hour - time.Second)

		// Query commits based on filters
		var commits []models.Commit
		if params.User != "" {
			// Filter by user (email or user_id)
			commits = store.GetCommitsByUser(params.User, from, to)
		} else if params.RepoName != "" {
			// Filter by repo
			commits = store.GetCommitsByRepo(params.RepoName, from, to)
		} else {
			// Get all commits in range
			commits = store.GetCommitsByTimeRange(from, to)
		}

		// Apply pagination
		totalCount := len(commits)
		start := (params.Page - 1) * params.PageSize
		end := start + params.PageSize

		// Bounds checking
		if start > totalCount {
			start = totalCount
		}
		if end > totalCount {
			end = totalCount
		}

		// Extract page of commits
		pageCommits := commits[start:end]

		// Build response matching CommitsResponse schema from OpenAPI spec
		response := models.CommitsResponse{
			Items:      pageCommits,
			TotalCount: totalCount,
			Page:       params.Page,
			PageSize:   params.PageSize,
		}

		// Send JSON response
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// AICodeCommitsCSV returns an HTTP handler for GET /analytics/ai-code/commits.csv.
// It reuses the query logic from AICodeCommits and exports results as CSV.
// CSV exports don't use pagination - all results are streamed.
func AICodeCommitsCSV(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters (startDate, endDate, user)
		params, err := api.ParseQueryParams(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse date range from already-validated params
		from, _ := time.Parse("2006-01-02", params.StartDate)
		to, _ := time.Parse("2006-01-02", params.EndDate)

		// Add time to include full day
		to = to.Add(24*time.Hour - time.Second)

		// Query commits based on filters (same logic as JSON endpoint)
		var commits []models.Commit
		if params.User != "" {
			commits = store.GetCommitsByUser(params.User, from, to)
		} else if params.RepoName != "" {
			commits = store.GetCommitsByRepo(params.RepoName, from, to)
		} else {
			commits = store.GetCommitsByTimeRange(from, to)
		}

		// No pagination for CSV export - return all results
		// (CSV exports are typically used for full data dumps)

		// Send CSV response
		api.RespondCSV(w, commits)
	})
}
