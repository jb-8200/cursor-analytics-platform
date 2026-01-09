package github

import (
	"net/http"
	"sort"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// PRCycleTimeResponse is the response for GET /analytics/github/pr-cycle-time.
type PRCycleTimeResponse struct {
	Data   PRCycleTimeData   `json:"data"`
	Params PRCycleTimeParams `json:"params"`
}

// PRCycleTimeData contains PR lifecycle metrics.
type PRCycleTimeData struct {
	AvgTimeToFirstReview float64 `json:"avgTimeToFirstReview"` // seconds
	AvgTimeToMerge       float64 `json:"avgTimeToMerge"`       // seconds
	MedianTimeToMerge    float64 `json:"medianTimeToMerge"`    // seconds
	P50TimeToMerge       float64 `json:"p50TimeToMerge"`       // seconds
	P75TimeToMerge       float64 `json:"p75TimeToMerge"`       // seconds
	P90TimeToMerge       float64 `json:"p90TimeToMerge"`       // seconds
	TotalPRsAnalyzed     int     `json:"totalPRsAnalyzed"`
}

// PRCycleTimeParams contains the request parameters used.
type PRCycleTimeParams struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// PRCycleTimeAnalytics returns an HTTP handler for GET /analytics/github/pr-cycle-time.
// It calculates PR lifecycle metrics including time to first review and time to merge.
func PRCycleTimeAnalytics(store storage.Store) http.Handler {
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
		data := calculateCycleTimeMetrics(filteredPRs)

		// Build response
		response := PRCycleTimeResponse{
			Data: data,
			Params: PRCycleTimeParams{
				From: fromStr,
				To:   toStr,
			},
		}

		api.RespondJSON(w, http.StatusOK, response)
	})
}

// calculateCycleTimeMetrics computes PR lifecycle metrics from a list of PRs.
func calculateCycleTimeMetrics(prs []models.PullRequest) PRCycleTimeData {
	if len(prs) == 0 {
		return PRCycleTimeData{
			TotalPRsAnalyzed: 0,
		}
	}

	var timeToFirstReviewSum float64
	var timeToFirstReviewCount int
	var timeToMergeValues []float64

	for _, pr := range prs {
		// Calculate time to first review (if available)
		if pr.FirstReviewAt != nil {
			timeToFirstReview := pr.FirstReviewAt.Sub(pr.CreatedAt).Seconds()
			if timeToFirstReview >= 0 {
				timeToFirstReviewSum += timeToFirstReview
				timeToFirstReviewCount++
			}
		}

		// Calculate time to merge
		if pr.MergedAt != nil {
			timeToMerge := pr.MergedAt.Sub(pr.CreatedAt).Seconds()
			if timeToMerge >= 0 {
				timeToMergeValues = append(timeToMergeValues, timeToMerge)
			}
		}
	}

	// Calculate averages
	avgTimeToFirstReview := 0.0
	if timeToFirstReviewCount > 0 {
		avgTimeToFirstReview = timeToFirstReviewSum / float64(timeToFirstReviewCount)
	}

	avgTimeToMerge := 0.0
	if len(timeToMergeValues) > 0 {
		sum := 0.0
		for _, val := range timeToMergeValues {
			sum += val
		}
		avgTimeToMerge = sum / float64(len(timeToMergeValues))
	}

	// Calculate median and percentiles
	medianTimeToMerge := 0.0
	p50 := 0.0
	p75 := 0.0
	p90 := 0.0

	if len(timeToMergeValues) > 0 {
		// Sort for percentile calculation
		sort.Float64s(timeToMergeValues)

		medianTimeToMerge = percentile(timeToMergeValues, 50)
		p50 = percentile(timeToMergeValues, 50)
		p75 = percentile(timeToMergeValues, 75)
		p90 = percentile(timeToMergeValues, 90)
	}

	return PRCycleTimeData{
		AvgTimeToFirstReview: avgTimeToFirstReview,
		AvgTimeToMerge:       avgTimeToMerge,
		MedianTimeToMerge:    medianTimeToMerge,
		P50TimeToMerge:       p50,
		P75TimeToMerge:       p75,
		P90TimeToMerge:       p90,
		TotalPRsAnalyzed:     len(prs),
	}
}

// percentile calculates the percentile value from a sorted slice.
// p should be between 0 and 100.
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0.0
	}

	if len(sorted) == 1 {
		return sorted[0]
	}

	// Calculate index
	index := (p / 100.0) * float64(len(sorted)-1)
	lowerIndex := int(index)
	upperIndex := lowerIndex + 1

	if upperIndex >= len(sorted) {
		return sorted[len(sorted)-1]
	}

	// Linear interpolation between values
	weight := index - float64(lowerIndex)
	return sorted[lowerIndex]*(1-weight) + sorted[upperIndex]*weight
}
