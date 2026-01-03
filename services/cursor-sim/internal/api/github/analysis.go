package github

import (
	"net/http"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/services"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// SurvivalAnalysisHandler returns an HTTP handler for GET /repos/{owner}/{repo}/analysis/survival.
// It calculates file-level code survival metrics for a repository.
func SurvivalAnalysisHandler(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repoName := parseRepoFromPath(r.URL.Path)
		if repoName == "" {
			respondError(w, http.StatusBadRequest, "invalid repository path")
			return
		}

		// Parse query parameters
		cohortStartStr := r.URL.Query().Get("cohort_start")
		cohortEndStr := r.URL.Query().Get("cohort_end")
		observationDateStr := r.URL.Query().Get("observation_date")

		// Default to last 30 days cohort if not specified
		now := time.Now()
		cohortEnd := now
		cohortStart := now.AddDate(0, 0, -30)
		observationDate := now

		if cohortStartStr != "" {
			if parsed, err := time.Parse("2006-01-02", cohortStartStr); err == nil {
				cohortStart = parsed
			}
		}

		// Store the original end date string for response formatting
		cohortEndFormatted := cohortEnd.Format("2006-01-02")
		observationDateFormatted := observationDate.Format("2006-01-02")

		if cohortEndStr != "" {
			if parsed, err := time.Parse("2006-01-02", cohortEndStr); err == nil {
				cohortEndFormatted = cohortEndStr
				// Add 24 hours to include the full end date in the range query
				cohortEnd = parsed.Add(24 * time.Hour)
			}
		} else {
			// When no end is specified, add 24h to include full current day
			cohortEnd = cohortEnd.Add(24 * time.Hour)
		}

		if observationDateStr != "" {
			if parsed, err := time.Parse("2006-01-02", observationDateStr); err == nil {
				observationDateFormatted = observationDateStr
				observationDate = parsed.Add(24 * time.Hour) // Include full day
			}
		} else {
			// When no observation date is specified, add 24h
			observationDate = observationDate.Add(24 * time.Hour)
		}

		// Calculate survival metrics
		svc := services.NewSurvivalService(store)
		analysis, err := svc.CalculateSurvival(repoName, cohortStart, cohortEnd, observationDate)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "failed to calculate survival: "+err.Error())
			return
		}

		// Override the formatted dates to match the original query parameters
		analysis.CohortEnd = cohortEndFormatted
		analysis.ObservationDate = observationDateFormatted

		respondJSON(w, http.StatusOK, analysis)
	})
}
