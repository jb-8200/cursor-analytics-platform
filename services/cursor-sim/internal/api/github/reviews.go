package github

import (
	"net/http"
	"strconv"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// ReviewsAnalyticsResponse is the response for GET /analytics/github/reviews.
type ReviewsAnalyticsResponse struct {
	Data       []models.Review   `json:"data"`
	Pagination ReviewsPagination `json:"pagination"`
	Params     ReviewsParams     `json:"params"`
}

// ReviewsPagination contains pagination metadata.
type ReviewsPagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

// ReviewsParams contains the request parameters used.
type ReviewsParams struct {
	PRID     int    `json:"pr_id,omitempty"`
	Reviewer string `json:"reviewer,omitempty"`
}

// ListReviewsAnalytics returns an HTTP handler for GET /analytics/github/reviews.
// It returns a paginated list of reviews with optional filtering by pr_id and reviewer.
func ListReviewsAnalytics(store storage.Store) http.Handler {
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
		prIDStr := query.Get("pr_id")
		reviewer := query.Get("reviewer")

		var prID int
		var err error

		if prIDStr != "" {
			prID, err = strconv.Atoi(prIDStr)
			if err != nil || prID < 1 {
				api.RespondError(w, http.StatusBadRequest, "invalid pr_id parameter")
				return
			}
		}

		// Get reviews based on filters
		var allReviews []models.Review

		// Apply filters in priority order
		if prID > 0 && reviewer != "" {
			// Both filters - get by PR and filter by reviewer
			reviews, err := store.GetReviewsByPRID(int64(prID))
			if err != nil {
				api.RespondError(w, http.StatusInternalServerError, "failed to get reviews by PR")
				return
			}
			// Filter by reviewer
			for _, review := range reviews {
				if review.Reviewer == reviewer {
					allReviews = append(allReviews, review)
				}
			}
		} else if prID > 0 {
			// Filter by PR only
			reviews, err := store.GetReviewsByPRID(int64(prID))
			if err != nil {
				api.RespondError(w, http.StatusInternalServerError, "failed to get reviews by PR")
				return
			}
			allReviews = reviews
		} else if reviewer != "" {
			// Filter by reviewer only
			reviews, err := store.GetReviewsByReviewer(reviewer)
			if err != nil {
				api.RespondError(w, http.StatusInternalServerError, "failed to get reviews by reviewer")
				return
			}
			allReviews = reviews
		} else {
			// Get all reviews - iterate through all PRs
			// Since there's no GetAllReviews method, we need to get reviews for each PR
			// For now, return empty if no filters provided (or implement GetAllReviews)
			// Let's implement a simple approach: get all PRs and their reviews
			repos := store.ListRepositories()
			for _, repo := range repos {
				prs := store.GetPRsByRepo(repo)
				for _, pr := range prs {
					reviews, err := store.GetReviewsByPRID(int64(pr.ID))
					if err != nil {
						continue // Skip on error
					}
					allReviews = append(allReviews, reviews...)
				}
			}
		}

		// Handle nil slice
		if allReviews == nil {
			allReviews = []models.Review{}
		}

		// Apply pagination
		total := len(allReviews)
		start := (page - 1) * pageSize
		end := start + pageSize

		if start > total {
			start = total
		}
		if end > total {
			end = total
		}

		pageReviews := allReviews[start:end]

		// Build response
		response := ReviewsAnalyticsResponse{
			Data: pageReviews,
			Pagination: ReviewsPagination{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
			Params: ReviewsParams{
				PRID:     prID,
				Reviewer: reviewer,
			},
		}

		api.RespondJSON(w, http.StatusOK, response)
	})
}
