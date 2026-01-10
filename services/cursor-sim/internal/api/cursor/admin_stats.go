package cursor

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"runtime"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/models"
	domainModels "github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// GetStats returns a handler that provides comprehensive statistics about the generated simulation data.
//
// Query Parameters:
//   - include_timeseries (bool): Include time series data (commits per day, PRs per day, cycle times)
//
// Response: models.StatsResponse with:
//   - Generation: Overall counts and data size
//   - Developers: Breakdown by seniority, region, team, activity
//   - Quality: Quality metrics (revert rate, hotfix rate, code survival, review thoroughness)
//   - Variance: Standard deviation metrics
//   - Performance: Generation time, memory usage, storage efficiency
//   - Organization: Teams, divisions, repositories
//   - TimeSeries: Optional time series data
func GetStats(store storage.Store, seedData *seed.SeedData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		includeTimeSeries := r.URL.Query().Get("include_timeseries") == "true"

		// Get seed structure
		developers := seedData.Developers
		repositories := seedData.Repositories

		// Get commits from storage (use all commits for stats)
		now := time.Now()
		past := now.AddDate(-10, 0, 0) // 10 years ago to get all data
		commits := store.GetCommitsByTimeRange(past, now)

		// Get PRs from all repos
		var allPRs []domainModels.PullRequest
		for _, repo := range repositories {
			prs := store.GetPRsByRepo(repo.RepoName)
			allPRs = append(allPRs, prs...)
		}

		// Get reviews from all PRs
		var allReviews []domainModels.Review
		for _, pr := range allPRs {
			reviews, _ := store.GetReviewsByPRID(int64(pr.ID))
			allReviews = append(allReviews, reviews...)
		}

		// Get issues from all repos
		var allIssues []domainModels.Issue
		for _, repo := range repositories {
			issues, _ := store.GetIssuesByRepo(repo.RepoName)
			allIssues = append(allIssues, issues...)
		}

		// Build response
		response := models.StatsResponse{
			Generation: models.Generation{
				TotalCommits:    len(commits),
				TotalPRs:        len(allPRs),
				TotalReviews:    len(allReviews),
				TotalIssues:     len(allIssues),
				TotalDevelopers: len(developers),
				DataSize:        formatBytes(estimateDataSize(len(commits), len(allPRs), len(allReviews), len(allIssues))),
			},
			Developers: models.Developers{
				BySeniority: groupBySeniority(seedData),
				ByRegion:    groupByRegion(seedData),
				ByTeam:      groupByTeam(seedData),
				ByActivity:  groupByActivity(seedData),
			},
			Quality:     calculateQualityMetrics(allPRs, allReviews),
			Variance:    calculateVariance(commits, allPRs),
			Performance: calculatePerformance(),
			Organization: models.Organization{
				Teams:        extractUniqueTeams(seedData),
				Divisions:    extractUniqueDivisions(seedData),
				Repositories: extractUniqueRepos(seedData),
			},
		}

		// Add time series data if requested
		if includeTimeSeries {
			response.TimeSeries = &models.TimeSeries{
				CommitsPerDay: calculateCommitsPerDay(commits),
				PRsPerDay:     calculatePRsPerDay(allPRs),
				AvgCycleTime:  calculateAvgCycleTime(allPRs),
			}
		}

		// Write JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// calculateQualityMetrics calculates quality metrics from PR and review data.
// For now, returns placeholder mock data. Future implementation would analyze actual data.
func calculateQualityMetrics(prs []domainModels.PullRequest, reviews []domainModels.Review) models.Quality {
	if len(prs) == 0 {
		return models.Quality{
			AvgRevertRate:         0.0,
			AvgHotfixRate:         0.0,
			AvgCodeSurvival:       0.0,
			AvgReviewThoroughness: 0.0,
			AvgIterations:         0.0,
		}
	}

	// Calculate average review thoroughness
	// Thoroughness = avg comments per review (normalized to 0-1 scale, assuming max 10 comments)
	var totalComments int
	for _, review := range reviews {
		totalComments += review.CommentCount()
	}
	avgCommentsPerReview := 0.0
	if len(reviews) > 0 {
		avgCommentsPerReview = float64(totalComments) / float64(len(reviews))
	}
	// Normalize to 0-1 scale (assume max 10 comments for thoroughness = 1.0)
	reviewThoroughness := math.Min(avgCommentsPerReview/10.0, 1.0)

	// Calculate average iterations per PR
	// Estimate iterations from review count per PR
	prReviewCounts := make(map[int]int)
	for _, review := range reviews {
		prReviewCounts[int(review.PRID)]++
	}
	totalIterations := 0
	for _, count := range prReviewCounts {
		totalIterations += count
	}
	avgIterations := 0.0
	if len(prReviewCounts) > 0 {
		avgIterations = float64(totalIterations) / float64(len(prReviewCounts))
	}

	// Placeholder values for revert rate, hotfix rate, code survival
	// These would require additional analysis services in a full implementation
	return models.Quality{
		AvgRevertRate:         0.02, // Mock: 2% revert rate
		AvgHotfixRate:         0.08, // Mock: 8% hotfix rate
		AvgCodeSurvival:       0.85, // Mock: 85% code survival at 30 days
		AvgReviewThoroughness: reviewThoroughness,
		AvgIterations:         avgIterations,
	}
}

// calculateVariance calculates variance metrics (standard deviation) from commit and PR data.
func calculateVariance(commits []domainModels.Commit, prs []domainModels.PullRequest) models.Variance {
	if len(commits) == 0 {
		return models.Variance{
			CommitsStdDev:   0.0,
			PRSizeStdDev:    0.0,
			CycleTimeStdDev: 0.0,
		}
	}

	// Calculate commits per developer standard deviation
	commitsPerDev := make(map[string]int)
	for _, commit := range commits {
		commitsPerDev[commit.UserID]++
	}
	commitCounts := make([]float64, 0, len(commitsPerDev))
	for _, count := range commitsPerDev {
		commitCounts = append(commitCounts, float64(count))
	}
	commitsStdDev := calculateStdDev(commitCounts)

	// Calculate PR size standard deviation (in lines of code)
	prSizes := make([]float64, len(prs))
	for i, pr := range prs {
		prSizes[i] = float64(pr.Additions + pr.Deletions)
	}
	prSizeStdDev := calculateStdDev(prSizes)

	// Calculate cycle time standard deviation (in days)
	cycleTimes := make([]float64, 0, len(prs))
	for _, pr := range prs {
		if pr.MergedAt != nil && !pr.CreatedAt.IsZero() {
			cycleTime := pr.MergedAt.Sub(pr.CreatedAt).Hours() / 24.0 // Convert to days
			cycleTimes = append(cycleTimes, cycleTime)
		}
	}
	cycleTimeStdDev := calculateStdDev(cycleTimes)

	return models.Variance{
		CommitsStdDev:   commitsStdDev,
		PRSizeStdDev:    prSizeStdDev,
		CycleTimeStdDev: cycleTimeStdDev,
	}
}

// calculateStdDev calculates the standard deviation of a slice of float64 values.
func calculateStdDev(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	// Calculate mean
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	// Calculate variance
	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))

	// Standard deviation is square root of variance
	return math.Sqrt(variance)
}

// calculatePerformance calculates performance metrics.
func calculatePerformance() models.Performance {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return models.Performance{
		LastGenerationTime: "N/A", // Would need to track this in generator
		MemoryUsage:        formatBytes(int(m.Alloc)),
		StorageEfficiency:  "95%", // Mock value
	}
}

// calculateCommitsPerDay calculates commits per day time series.
func calculateCommitsPerDay(commits []domainModels.Commit) []int {
	if len(commits) == 0 {
		return []int{}
	}

	// Find date range
	minDate, maxDate := commits[0].CommitTs, commits[0].CommitTs
	for _, commit := range commits {
		if commit.CommitTs.Before(minDate) {
			minDate = commit.CommitTs
		}
		if commit.CommitTs.After(maxDate) {
			maxDate = commit.CommitTs
		}
	}

	// Create buckets for each day
	days := int(maxDate.Sub(minDate).Hours()/24) + 1
	if days > 365 {
		days = 365 // Limit to 1 year of data
	}

	commitCounts := make([]int, days)
	for _, commit := range commits {
		dayIndex := int(commit.CommitTs.Sub(minDate).Hours() / 24)
		if dayIndex >= 0 && dayIndex < days {
			commitCounts[dayIndex]++
		}
	}

	return commitCounts
}

// calculatePRsPerDay calculates PRs per day time series.
func calculatePRsPerDay(prs []domainModels.PullRequest) []int {
	if len(prs) == 0 {
		return []int{}
	}

	// Find date range
	minDate, maxDate := prs[0].CreatedAt, prs[0].CreatedAt
	for _, pr := range prs {
		if pr.CreatedAt.Before(minDate) {
			minDate = pr.CreatedAt
		}
		if pr.CreatedAt.After(maxDate) {
			maxDate = pr.CreatedAt
		}
	}

	// Create buckets for each day
	days := int(maxDate.Sub(minDate).Hours()/24) + 1
	if days > 365 {
		days = 365 // Limit to 1 year of data
	}

	prCounts := make([]int, days)
	for _, pr := range prs {
		dayIndex := int(pr.CreatedAt.Sub(minDate).Hours() / 24)
		if dayIndex >= 0 && dayIndex < days {
			prCounts[dayIndex]++
		}
	}

	return prCounts
}

// calculateAvgCycleTime calculates average cycle time per day.
func calculateAvgCycleTime(prs []domainModels.PullRequest) []float64 {
	if len(prs) == 0 {
		return []float64{}
	}

	// Find date range
	minDate, maxDate := prs[0].CreatedAt, prs[0].CreatedAt
	for _, pr := range prs {
		if pr.CreatedAt.Before(minDate) {
			minDate = pr.CreatedAt
		}
		if pr.CreatedAt.After(maxDate) {
			maxDate = pr.CreatedAt
		}
	}

	// Create buckets for each day
	days := int(maxDate.Sub(minDate).Hours()/24) + 1
	if days > 365 {
		days = 365 // Limit to 1 year of data
	}

	cycleTimes := make([][]float64, days)
	for _, pr := range prs {
		if pr.MergedAt != nil && !pr.CreatedAt.IsZero() {
			dayIndex := int(pr.CreatedAt.Sub(minDate).Hours() / 24)
			if dayIndex >= 0 && dayIndex < days {
				cycleTime := pr.MergedAt.Sub(pr.CreatedAt).Hours() / 24.0 // Convert to days
				cycleTimes[dayIndex] = append(cycleTimes[dayIndex], cycleTime)
			}
		}
	}

	// Calculate average cycle time for each day
	avgCycleTimes := make([]float64, days)
	for i, times := range cycleTimes {
		if len(times) > 0 {
			sum := 0.0
			for _, t := range times {
				sum += t
			}
			avgCycleTimes[i] = sum / float64(len(times))
		} else {
			avgCycleTimes[i] = 0.0
		}
	}

	return avgCycleTimes
}

// formatBytes formats a byte count as a human-readable string.
func formatBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// estimateDataSize estimates the total data size in bytes.
func estimateDataSize(commits, prs, reviews, issues int) int {
	// Rough estimates per item:
	// Commit: ~500 bytes
	// PR: ~1KB
	// Review: ~300 bytes
	// Issue: ~500 bytes
	return commits*500 + prs*1024 + reviews*300 + issues*500
}
