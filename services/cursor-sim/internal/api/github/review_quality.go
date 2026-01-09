package github

import (
	"net/http"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// ReviewQualityResponse is the response for GET /analytics/github/review-quality.
type ReviewQualityResponse struct {
	Data   ReviewQualityData   `json:"data"`
	Params ReviewQualityParams `json:"params"`
}

// ReviewQualityData contains review quality metrics.
type ReviewQualityData struct {
	ApprovalRate           float64 `json:"approval_rate"`
	AvgReviewersPerPR      float64 `json:"avg_reviewers_per_pr"`
	AvgCommentsPerReview   float64 `json:"avg_comments_per_review"`
	ChangesRequestedRate   float64 `json:"changes_requested_rate"`
	PendingRate            float64 `json:"pending_rate"`
	TotalReviews           int     `json:"total_reviews"`
	TotalPRsReviewed       int     `json:"total_prs_reviewed"`
}

// ReviewQualityParams contains the request parameters used.
type ReviewQualityParams struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// ReviewQualityAnalytics returns an HTTP handler for GET /analytics/github/review-quality.
// It calculates review quality metrics including approval rate and average reviewers per PR.
func ReviewQualityAnalytics(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		query := r.URL.Query()
		fromStr := query.Get("from")
		toStr := query.Get("to")

		// Parse date range
		var fromTime, toTime time.Time
		var err error

		if fromStr != "" {
			fromTime, err = time.Parse("2006-01-02", fromStr)
			if err != nil {
				api.RespondError(w, http.StatusBadRequest, "invalid from date format (use YYYY-MM-DD)")
				return
			}
		} else {
			fromTime = time.Time{} // Zero time (no lower bound)
		}

		if toStr != "" {
			toTime, err = time.Parse("2006-01-02", toStr)
			if err != nil {
				api.RespondError(w, http.StatusBadRequest, "invalid to date format (use YYYY-MM-DD)")
				return
			}
			// Set to end of day
			toTime = toTime.Add(24*time.Hour - time.Nanosecond)
		} else {
			toTime = time.Now() // Current time
		}

		// Get all merged PRs across all repositories
		var allMergedPRs []models.PullRequest
		repos := store.ListRepositories()
		for _, repo := range repos {
			prs := store.GetPRsByRepoAndState(repo, models.PRStateMerged)
			allMergedPRs = append(allMergedPRs, prs...)
		}

		// Filter by date range (based on merge date)
		var filteredPRs []models.PullRequest
		for _, pr := range allMergedPRs {
			if pr.MergedAt == nil {
				continue
			}

			// Check if within date range
			if !fromTime.IsZero() && pr.MergedAt.Before(fromTime) {
				continue
			}
			if pr.MergedAt.After(toTime) {
				continue
			}

			filteredPRs = append(filteredPRs, pr)
		}

		// Calculate metrics
		data := calculateReviewQualityMetrics(store, filteredPRs)

		// Build response
		response := ReviewQualityResponse{
			Data: data,
			Params: ReviewQualityParams{
				From: fromStr,
				To:   toStr,
			},
		}

		api.RespondJSON(w, http.StatusOK, response)
	})
}

// calculateReviewQualityMetrics computes review quality metrics from a list of PRs.
func calculateReviewQualityMetrics(store storage.Store, prs []models.PullRequest) ReviewQualityData {
	if len(prs) == 0 {
		return ReviewQualityData{
			TotalReviews:     0,
			TotalPRsReviewed: 0,
		}
	}

	var allReviews []models.Review
	var totalComments int
	var approvalCount int
	var changesRequestedCount int
	var pendingCount int

	// Track unique reviewers per PR
	prReviewerMap := make(map[int]map[string]bool) // PR ID -> set of reviewers

	// Collect reviews for all filtered PRs
	for _, pr := range prs {
		reviews, err := store.GetReviewsByPRID(int64(pr.ID))
		if err != nil {
			continue
		}

		if prReviewerMap[pr.ID] == nil {
			prReviewerMap[pr.ID] = make(map[string]bool)
		}

		for _, review := range reviews {
			allReviews = append(allReviews, review)
			totalComments += review.CommentCount()

			// Count states
			switch review.State {
			case models.ReviewStateApproved:
				approvalCount++
			case models.ReviewStateChangesRequested:
				changesRequestedCount++
			case models.ReviewStatePending:
				pendingCount++
			}

			// Track unique reviewers per PR
			prReviewerMap[pr.ID][review.Reviewer] = true
		}
	}

	// Calculate rates
	totalReviews := len(allReviews)
	approvalRate := 0.0
	changesRequestedRate := 0.0
	pendingRate := 0.0

	if totalReviews > 0 {
		approvalRate = float64(approvalCount) / float64(totalReviews)
		changesRequestedRate = float64(changesRequestedCount) / float64(totalReviews)
		pendingRate = float64(pendingCount) / float64(totalReviews)
	}

	// Calculate average comments per review
	avgCommentsPerReview := 0.0
	if totalReviews > 0 {
		avgCommentsPerReview = float64(totalComments) / float64(totalReviews)
	}

	// Calculate average reviewers per PR
	totalReviewers := 0
	for _, reviewerSet := range prReviewerMap {
		totalReviewers += len(reviewerSet)
	}

	avgReviewersPerPR := 0.0
	if len(prs) > 0 {
		avgReviewersPerPR = float64(totalReviewers) / float64(len(prs))
	}

	return ReviewQualityData{
		ApprovalRate:           approvalRate,
		AvgReviewersPerPR:      avgReviewersPerPR,
		AvgCommentsPerReview:   avgCommentsPerReview,
		ChangesRequestedRate:   changesRequestedRate,
		PendingRate:            pendingRate,
		TotalReviews:           totalReviews,
		TotalPRsReviewed:       len(prs),
	}
}
