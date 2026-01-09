package github

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// PRsAnalyticsResponse is the response for GET /analytics/github/prs.
type PRsAnalyticsResponse struct {
	Data       []models.PullRequest `json:"data"`
	Pagination PRsPagination        `json:"pagination"`
	Params     PRsParams            `json:"params"`
}

// PRsPagination contains pagination metadata.
type PRsPagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

// PRsParams contains the request parameters used.
type PRsParams struct {
	Status    string `json:"status,omitempty"`
	Author    string `json:"author,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

// ListPRsAnalytics returns an HTTP handler for GET /analytics/github/prs.
// It returns a paginated list of PRs with optional filtering by status, author, and date range.
func ListPRsAnalytics(store storage.Store) http.Handler {
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
		status := query.Get("status")
		author := query.Get("author")
		startDateStr := query.Get("start_date")
		endDateStr := query.Get("end_date")

		var startDate, endDate time.Time
		var err error

		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				api.RespondError(w, http.StatusBadRequest, "invalid start_date format (use YYYY-MM-DD)")
				return
			}
		}

		if endDateStr != "" {
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				api.RespondError(w, http.StatusBadRequest, "invalid end_date format (use YYYY-MM-DD)")
				return
			}
			// Include the entire end day
			endDate = endDate.Add(24*time.Hour - time.Second)
		}

		// Get PRs based on filters
		var allPRs []models.PullRequest

		// Apply filters in sequence
		if status != "" {
			prState := models.PRState(status)
			prs, err := store.GetPRsByStatus(prState)
			if err != nil {
				api.RespondError(w, http.StatusInternalServerError, "failed to get PRs by status")
				return
			}
			allPRs = prs
		} else if author != "" {
			prs, err := store.GetPRsByAuthorEmail(author)
			if err != nil {
				api.RespondError(w, http.StatusInternalServerError, "failed to get PRs by author")
				return
			}
			allPRs = prs
		} else {
			// Get all PRs from all repositories
			repos := store.ListRepositories()
			for _, repo := range repos {
				repoPRs := store.GetPRsByRepo(repo)
				allPRs = append(allPRs, repoPRs...)
			}
		}

		// Apply additional filters
		var filteredPRs []models.PullRequest
		for _, pr := range allPRs {
			// Status filter (if not already applied)
			if status != "" && string(pr.State) != status {
				continue
			}

			// Author filter (if not already applied)
			if author != "" && pr.AuthorEmail != author {
				continue
			}

			// Date range filter
			if !startDate.IsZero() && pr.CreatedAt.Before(startDate) {
				continue
			}
			if !endDate.IsZero() && pr.CreatedAt.After(endDate) {
				continue
			}

			filteredPRs = append(filteredPRs, pr)
		}

		// Handle nil slice
		if filteredPRs == nil {
			filteredPRs = []models.PullRequest{}
		}

		// Apply pagination
		total := len(filteredPRs)
		start := (page - 1) * pageSize
		end := start + pageSize

		if start > total {
			start = total
		}
		if end > total {
			end = total
		}

		pagePRs := filteredPRs[start:end]

		// Build response
		response := PRsAnalyticsResponse{
			Data: pagePRs,
			Pagination: PRsPagination{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
			Params: PRsParams{
				Status:    status,
				Author:    author,
				StartDate: startDateStr,
				EndDate:   endDateStr,
			},
		}

		api.RespondJSON(w, http.StatusOK, response)
	})
}
