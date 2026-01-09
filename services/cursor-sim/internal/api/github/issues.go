package github

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// IssuesAnalyticsResponse is the response for GET /analytics/github/issues.
type IssuesAnalyticsResponse struct {
	Data       []models.Issue   `json:"data"`
	Pagination IssuesPagination `json:"pagination"`
	Params     IssuesParams     `json:"params"`
}

// IssuesPagination contains pagination metadata.
type IssuesPagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

// IssuesParams contains the request parameters used.
type IssuesParams struct {
	State  string `json:"state,omitempty"`
	Labels string `json:"labels,omitempty"`
}

// ListIssuesAnalytics returns an HTTP handler for GET /analytics/github/issues.
// It returns a paginated list of issues with optional filtering by state and labels.
func ListIssuesAnalytics(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		query := r.URL.Query()

		// Parse pagination
		page := 1
		pageSize := 20

		if pageStr := query.Get("page"); pageStr != "" {
			p, err := strconv.Atoi(pageStr)
			if err != nil || p < 1 {
				api.RespondError(w, http.StatusBadRequest, "invalid page parameter")
				return
			}
			page = p
		}

		if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
			ps, err := strconv.Atoi(pageSizeStr)
			if err != nil || ps < 1 {
				api.RespondError(w, http.StatusBadRequest, "invalid page_size parameter")
				return
			}
			pageSize = ps
			if pageSize > 100 {
				pageSize = 100 // Max page size
			}
		}

		// Parse filters
		stateStr := query.Get("state")
		labelsStr := query.Get("labels")

		// Validate state if provided
		if stateStr != "" && stateStr != string(models.IssueStateOpen) && stateStr != string(models.IssueStateClosed) {
			api.RespondError(w, http.StatusBadRequest, "invalid state parameter (must be 'open' or 'closed')")
			return
		}

		// Parse labels (comma-separated)
		var requestedLabels []string
		if labelsStr != "" {
			requestedLabels = strings.Split(labelsStr, ",")
			// Trim whitespace from each label
			for i, label := range requestedLabels {
				requestedLabels[i] = strings.TrimSpace(label)
			}
		}

		// Get all issues from all repositories
		var allIssues []models.Issue
		repos := store.ListRepositories()
		for _, repo := range repos {
			var repoIssues []models.Issue
			var err error

			// If state filter is provided, use GetIssuesByState
			if stateStr != "" {
				repoIssues, err = store.GetIssuesByState(repo, models.IssueState(stateStr))
			} else {
				repoIssues, err = store.GetIssuesByRepo(repo)
			}

			if err != nil {
				continue // Skip on error
			}
			allIssues = append(allIssues, repoIssues...)
		}

		// Apply label filters
		var filteredIssues []models.Issue
		for _, issue := range allIssues {
			// Skip if doesn't match all requested labels
			if len(requestedLabels) > 0 {
				hasAllLabels := true
				for _, reqLabel := range requestedLabels {
					found := false
					for _, issueLabel := range issue.Labels {
						if issueLabel == reqLabel {
							found = true
							break
						}
					}
					if !found {
						hasAllLabels = false
						break
					}
				}
				if !hasAllLabels {
					continue
				}
			}

			filteredIssues = append(filteredIssues, issue)
		}

		// Handle nil slice
		if filteredIssues == nil {
			filteredIssues = []models.Issue{}
		}

		// Apply pagination
		total := len(filteredIssues)
		start := (page - 1) * pageSize
		end := start + pageSize

		if start > total {
			start = total
		}
		if end > total {
			end = total
		}

		pageIssues := filteredIssues[start:end]

		// Build response
		response := IssuesAnalyticsResponse{
			Data: pageIssues,
			Pagination: IssuesPagination{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
			Params: IssuesParams{
				State:  stateStr,
				Labels: labelsStr,
			},
		}

		api.RespondJSON(w, http.StatusOK, response)
	})
}
