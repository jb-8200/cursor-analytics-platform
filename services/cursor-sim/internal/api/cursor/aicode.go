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
func AICodeCommits(store storage.Store) http.Handler {
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
			api.RespondError(w, http.StatusBadRequest, "invalid from date: must be YYYY-MM-DD format")
			return
		}
		to, err := time.Parse("2006-01-02", params.To)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid to date: must be YYYY-MM-DD format")
			return
		}

		// Add time to include full day
		to = to.Add(24*time.Hour - time.Second)

		// Query commits based on filters
		var commits []models.Commit
		if params.UserID != "" {
			// Filter by user
			commits = store.GetCommitsByUser(params.UserID, from, to)
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

		// Build paginated response
		response := api.BuildPaginatedResponse(pageCommits, params, totalCount)

		// Send JSON response
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// AICodeCommitsCSV returns an HTTP handler for GET /analytics/ai-code/commits.csv.
// It reuses the query logic from AICodeCommits and exports results as CSV.
func AICodeCommitsCSV(store storage.Store) http.Handler {
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
			api.RespondError(w, http.StatusBadRequest, "invalid from date: must be YYYY-MM-DD format")
			return
		}
		to, err := time.Parse("2006-01-02", params.To)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid to date: must be YYYY-MM-DD format")
			return
		}

		// Add time to include full day
		to = to.Add(24*time.Hour - time.Second)

		// Query commits based on filters (same logic as JSON endpoint)
		var commits []models.Commit
		if params.UserID != "" {
			commits = store.GetCommitsByUser(params.UserID, from, to)
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
