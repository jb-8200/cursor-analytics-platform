package services

import (
	"regexp"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

// HotfixPRStore defines the storage interface for accessing PRs.
type HotfixPRStore interface {
	GetPRsByRepo(repoName string) []models.PullRequest
	GetPRsByRepoAndState(repoName string, state models.PRState) []models.PullRequest
}

// HotfixService detects hotfix PRs that follow merged PRs.
type HotfixService struct {
	store HotfixPRStore
}

// NewHotfixService creates a new hotfix analysis service.
func NewHotfixService(store HotfixPRStore) *HotfixService {
	return &HotfixService{
		store: store,
	}
}

// IsHotfixPR determines if a PR is a hotfix based on title/body content.
// Checks for patterns: "fix", "hotfix", "urgent", "patch" (as whole words)
func IsHotfixPR(pr models.PullRequest) bool {
	patterns := []string{
		`(?i)\bfix\b`,
		`(?i)\bhotfix\b`,
		`(?i)\burgent\b`,
		`(?i)\bpatch\b`,
	}

	content := pr.Title + " " + pr.Body

	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, content)
		if matched {
			return true
		}
	}

	return false
}

// FilesOverlap returns the common file paths between two PRs' changed files.
// For simulation, we generate synthetic file paths from PR data.
func FilesOverlap(pr1, pr2 models.PullRequest) []string {
	// Extract file extensions from repos (simplified for simulation)
	// In reality, we'd parse the actual diffs
	// For now, return files_in_common based on simple heuristic:
	// If PRs have similar changed_files count and same repo, they likely touch same files

	if pr1.RepoName != pr2.RepoName || pr1.ChangedFiles == 0 || pr2.ChangedFiles == 0 {
		return []string{}
	}

	// Simulate overlap: common files are subset of min(changed_files)
	overlapCount := (pr1.ChangedFiles + pr2.ChangedFiles) / 4 // Assume 25% overlap on average
	if overlapCount == 0 {
		overlapCount = 1
	}

	files := []string{}
	fileTypes := []string{".ts", ".tsx", ".go", ".py", ".js", ".jsx"}

	for i := 0; i < overlapCount && i < len(fileTypes); i++ {
		files = append(files, "src/module"+fileTypes[i])
	}

	return files
}

// GetHotfixes analyzes a repository for hotfix patterns within a time window.
// windowHours specifies how many hours after a merge to check for hotfixes.
// since and until define the time range to analyze.
func (h *HotfixService) GetHotfixes(
	repoName string,
	windowHours int,
	since, until time.Time,
) *models.HotfixAnalysis {
	// Get all merged PRs
	mergedPRs := h.store.GetPRsByRepoAndState(repoName, models.PRStateMerged)

	// Filter PRs merged within analysis period
	var analyzedPRs []models.PullRequest
	for _, pr := range mergedPRs {
		if pr.MergedAt != nil &&
			!pr.MergedAt.Before(since) &&
			pr.MergedAt.Before(until) {
			analyzedPRs = append(analyzedPRs, pr)
		}
	}

	if len(analyzedPRs) == 0 {
		return &models.HotfixAnalysis{
			WindowHours:    windowHours,
			TotalPRsMerged: 0,
			PRsWithHotfix:  0,
			HotfixRate:     0.0,
			HotfixPRs:      []models.HotfixPRInfo{},
		}
	}

	// Find hotfix relationships
	var hotfixes []models.HotfixPRInfo
	windowDuration := time.Duration(windowHours) * time.Hour

	for i, originalPR := range analyzedPRs {
		if originalPR.MergedAt == nil {
			continue
		}

		// Look at all subsequent PRs for hotfixes
		for j := i + 1; j < len(analyzedPRs); j++ {
			potentialHotfix := analyzedPRs[j]

			if potentialHotfix.MergedAt == nil {
				continue
			}

			// Check if within window
			timeSinceMerge := potentialHotfix.MergedAt.Sub(*originalPR.MergedAt)
			if timeSinceMerge < 0 || timeSinceMerge > windowDuration {
				continue
			}

			// Check if it's a hotfix PR (title/body pattern)
			if !IsHotfixPR(potentialHotfix) {
				continue
			}

			// Check for file overlap
			filesInCommon := FilesOverlap(originalPR, potentialHotfix)
			if len(filesInCommon) == 0 {
				continue
			}

			// Found a hotfix!
			hoursBetween := timeSinceMerge.Hours()
			hotfixes = append(hotfixes, models.HotfixPRInfo{
				OriginalPR:    originalPR.Number,
				HotfixPR:      potentialHotfix.Number,
				HoursBetween:  hoursBetween,
				FilesInCommon: filesInCommon,
			})
		}
	}

	// Calculate hotfix rate
	prsWithHotfix := len(hotfixes)
	hotfixRate := 0.0
	if len(analyzedPRs) > 0 {
		hotfixRate = float64(prsWithHotfix) / float64(len(analyzedPRs))
	}

	return &models.HotfixAnalysis{
		WindowHours:    windowHours,
		TotalPRsMerged: len(analyzedPRs),
		PRsWithHotfix:  prsWithHotfix,
		HotfixRate:     hotfixRate,
		HotfixPRs:      hotfixes,
	}
}
