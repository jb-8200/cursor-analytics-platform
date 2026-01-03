package generator

import (
	"math/rand"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// ResearchStore defines the interface for research data generation.
type ResearchStore interface {
	GetCommitsByTimeRange(from, to time.Time) []models.Commit
	GetPRsByRepo(repoName string) []models.PullRequest
	GetPRsByRepoAndState(repoName string, state models.PRState) []models.PullRequest
	GetReviewComments(repoName string, prNumber int) []models.ReviewComment
	GetDeveloper(userID string) (*seed.Developer, error)
	ListRepositories() []string
}

// ResearchGenerator generates research datasets from simulated data.
type ResearchGenerator struct {
	seed  *seed.SeedData
	store ResearchStore
	rng   *rand.Rand
}

// NewResearchGenerator creates a new research generator with a random seed.
func NewResearchGenerator(seedData *seed.SeedData, store ResearchStore) *ResearchGenerator {
	return NewResearchGeneratorWithSeed(seedData, store, time.Now().UnixNano())
}

// NewResearchGeneratorWithSeed creates a new research generator with a specific seed for reproducibility.
func NewResearchGeneratorWithSeed(seedData *seed.SeedData, store ResearchStore, randSeed int64) *ResearchGenerator {
	return &ResearchGenerator{
		seed:  seedData,
		store: store,
		rng:   rand.New(rand.NewSource(randSeed)),
	}
}

// GenerateDataset generates a research dataset for the given time range.
func (g *ResearchGenerator) GenerateDataset(from, to time.Time) ([]models.ResearchDataPoint, error) {
	return g.JoinCommitPRData(from, to)
}

// JoinCommitPRData correlates commits with PRs to create research data points.
func (g *ResearchGenerator) JoinCommitPRData(from, to time.Time) ([]models.ResearchDataPoint, error) {
	commits := g.store.GetCommitsByTimeRange(from, to)
	if len(commits) == 0 {
		return nil, nil
	}

	// Build a map of PR by (repo, branch) for quick lookup
	prMap := g.buildPRMap()

	var dataPoints []models.ResearchDataPoint
	for _, commit := range commits {
		dp := g.createDataPointFromCommit(commit, prMap)
		if dp != nil {
			g.ApplyControlVariables(dp)
			dataPoints = append(dataPoints, *dp)
		}
	}

	return dataPoints, nil
}

// buildPRMap creates a map of PRs indexed by "repo:branch" for efficient lookup.
func (g *ResearchGenerator) buildPRMap() map[string]*models.PullRequest {
	prMap := make(map[string]*models.PullRequest)

	for _, repo := range g.store.ListRepositories() {
		prs := g.store.GetPRsByRepo(repo)
		for i := range prs {
			key := prs[i].RepoName + ":" + prs[i].HeadBranch
			prMap[key] = &prs[i]
		}
	}

	return prMap
}

// createDataPointFromCommit creates a ResearchDataPoint from a commit, correlating with PR if available.
func (g *ResearchGenerator) createDataPointFromCommit(commit models.Commit, prMap map[string]*models.PullRequest) *models.ResearchDataPoint {
	// Calculate AI-attributed lines
	aiLinesAdded := commit.TabLinesAdded + commit.ComposerLinesAdded
	nonAILinesAdded := commit.TotalLinesAdded - aiLinesAdded
	if nonAILinesAdded < 0 {
		nonAILinesAdded = 0
	}

	// Estimate AI-deleted lines (proportional to AI ratio)
	aiRatio := commit.AIRatio()
	aiLinesDeleted := int(float64(commit.TotalLinesDeleted) * aiRatio)

	dp := &models.ResearchDataPoint{
		// Identifiers
		CommitHash:  commit.CommitHash,
		AuthorID:    commit.UserID,
		AuthorEmail: commit.UserEmail,
		RepoName:    commit.RepoName,

		// AI Metrics
		AIRatio:         aiRatio,
		AILinesAdded:    aiLinesAdded,
		AILinesDeleted:  aiLinesDeleted,
		NonAILinesAdded: nonAILinesAdded,
		TabLines:        commit.TabLinesAdded,
		ComposerLines:   commit.ComposerLinesAdded,

		// PR Metrics
		PRVolume:     commit.TotalLinesAdded + commit.TotalLinesDeleted,
		Additions:    commit.TotalLinesAdded,
		Deletions:    commit.TotalLinesDeleted,
		PRScatter:    0, // Will be set from PR if available
		FilesChanged: 0, // Will be set from PR if available

		Timestamp: commit.CommitTs,
	}

	// Try to find associated PR
	prKey := commit.RepoName + ":" + commit.BranchName
	if pr, ok := prMap[prKey]; ok {
		dp.PRNumber = pr.Number
		dp.FilesChanged = pr.ChangedFiles
		dp.PRScatter = pr.ChangedFiles
		dp.AIRatio = pr.AIRatio // Use PR-level AI ratio if available

		// Calculate cycle times
		commits := g.getCommitsForPR(pr)
		dp.CodingLeadTimeHours = g.CalculateCodingLeadTime(commits)
		dp.ReviewLeadTimeHours = g.CalculateReviewLeadTime(*pr)

		// Calculate pickup time (first review delay)
		dp.PickupTimeHours = g.CalculatePickupTime(*pr)

		// Get review comments for metrics
		reviews := g.store.GetReviewComments(pr.RepoName, pr.Number)
		dp.ReviewIterations = g.CountReviewIterations(reviews)
		dp.IterationCount = dp.ReviewIterations // Canonical field

		// Review cost metrics
		dp.ReviewDensity = g.CalculateReviewDensity(reviews, dp.PRVolume)
		dp.ReviewerCount = g.CountUniqueReviewers(reviews)
		dp.ReworkRatio = g.CalculateReworkRatio(*pr)
		dp.ScopeCreep = g.CalculateScopeCreep(*pr)

		// Quality outcomes
		dp.WasReverted = pr.WasReverted
		dp.IsReverted = pr.WasReverted // Canonical field
		dp.RequiredHotfix = pr.IsBugFix && pr.WasReverted
		dp.HasHotfixFollowup = g.HasHotfixFollowup(*pr) // Detect actual hotfixes

		// Greenfield metrics
		dp.IsGreenfield = g.IsGreenfield(*pr)
		dp.GreenfieldIndex = g.CalculateGreenfieldIndex(*pr)

		// Code survival
		dp.SurvivalRate30d = g.CalculateSurvivalRate30d(*pr)
	}

	return dp
}

// getCommitsForPR retrieves commits associated with a PR based on time range.
func (g *ResearchGenerator) getCommitsForPR(pr *models.PullRequest) []models.Commit {
	// Get commits in the PR's time range
	endTime := pr.CreatedAt
	if pr.MergedAt != nil {
		endTime = *pr.MergedAt
	}
	startTime := pr.CreatedAt.Add(-7 * 24 * time.Hour) // Look back 7 days

	commits := g.store.GetCommitsByTimeRange(startTime, endTime)

	// Filter to commits on the same branch and repo
	var prCommits []models.Commit
	for _, c := range commits {
		if c.RepoName == pr.RepoName && c.BranchName == pr.HeadBranch {
			prCommits = append(prCommits, c)
		}
	}

	return prCommits
}

// CalculateCodingLeadTime calculates the time from first to last commit in hours.
func (g *ResearchGenerator) CalculateCodingLeadTime(commits []models.Commit) float64 {
	if len(commits) < 2 {
		return 0
	}

	// Find first and last commit times
	firstTime := commits[0].CommitTs
	lastTime := commits[0].CommitTs

	for _, c := range commits {
		if c.CommitTs.Before(firstTime) {
			firstTime = c.CommitTs
		}
		if c.CommitTs.After(lastTime) {
			lastTime = c.CommitTs
		}
	}

	return lastTime.Sub(firstTime).Hours()
}

// CalculateReviewLeadTime calculates the time from PR creation to merge in hours.
func (g *ResearchGenerator) CalculateReviewLeadTime(pr models.PullRequest) float64 {
	if pr.MergedAt == nil {
		return 0
	}
	return pr.MergedAt.Sub(pr.CreatedAt).Hours()
}

// CalculateMergeLeadTime calculates the time from approval to merge in hours.
func (g *ResearchGenerator) CalculateMergeLeadTime(approvedAt, mergedAt time.Time) float64 {
	return mergedAt.Sub(approvedAt).Hours()
}

// IsGreenfield determines if a PR represents greenfield (new) code.
// A PR is considered greenfield if additions significantly outweigh deletions.
func (g *ResearchGenerator) IsGreenfield(pr models.PullRequest) bool {
	if pr.Additions == 0 {
		return false
	}
	if pr.Deletions == 0 {
		return true
	}

	// Consider greenfield if additions are more than 80% of total changes
	additionRatio := float64(pr.Additions) / float64(pr.Additions+pr.Deletions)
	return additionRatio > 0.80
}

// CalculatePickupTime calculates the time from PR creation to first review in hours.
func (g *ResearchGenerator) CalculatePickupTime(pr models.PullRequest) float64 {
	if pr.FirstReviewAt == nil {
		return 0
	}
	return pr.FirstReviewAt.Sub(pr.CreatedAt).Hours()
}

// CalculateReviewDensity calculates review comments per line of code.
func (g *ResearchGenerator) CalculateReviewDensity(reviews []models.ReviewComment, prVolume int) float64 {
	if prVolume == 0 {
		return 0
	}
	return float64(len(reviews)) / float64(prVolume)
}

// CountUniqueReviewers counts the number of unique reviewers for a PR.
func (g *ResearchGenerator) CountUniqueReviewers(reviews []models.ReviewComment) int {
	reviewers := make(map[string]bool)
	for _, r := range reviews {
		reviewers[r.AuthorID] = true
	}
	return len(reviewers)
}

// CalculateReworkRatio calculates the ratio of lines changed during review to initial lines.
func (g *ResearchGenerator) CalculateReworkRatio(pr models.PullRequest) float64 {
	if pr.InitialAdditions == 0 {
		return 0
	}
	// Rework = difference between final and initial additions
	rework := abs(pr.Additions - pr.InitialAdditions)
	return float64(rework) / float64(pr.InitialAdditions)
}

// CalculateScopeCreep calculates the change in scope during review.
func (g *ResearchGenerator) CalculateScopeCreep(pr models.PullRequest) float64 {
	finalLoC := pr.Additions
	if finalLoC == 0 {
		return 0
	}
	initialLoC := pr.InitialAdditions
	return float64(finalLoC-initialLoC) / float64(finalLoC)
}

// HasHotfixFollowup determines if the PR has a hotfix follow-up within 48 hours.
func (g *ResearchGenerator) HasHotfixFollowup(pr models.PullRequest) bool {
	if pr.MergedAt == nil {
		return false
	}

	// Get all PRs in the same repo
	prs := g.store.GetPRsByRepoAndState(pr.RepoName, models.PRStateMerged)

	// Look for hotfix PRs within 48 hours
	windowEnd := pr.MergedAt.Add(48 * time.Hour)
	for _, otherPR := range prs {
		if otherPR.Number == pr.Number {
			continue
		}
		if otherPR.MergedAt == nil {
			continue
		}

		// Check if within window
		if otherPR.MergedAt.After(*pr.MergedAt) && otherPR.MergedAt.Before(windowEnd) {
			// Check if it's a fix PR
			if otherPR.IsBugFix {
				return true
			}
		}
	}

	return false
}

// CalculateGreenfieldIndex calculates the percentage of PR lines in new files.
func (g *ResearchGenerator) CalculateGreenfieldIndex(pr models.PullRequest) float64 {
	// For simulation purposes, estimate based on PR age and additions/deletions ratio
	// A PR with high additions and low deletions is likely greenfield
	if pr.Additions == 0 {
		return 0.0
	}

	totalChanges := pr.Additions + pr.Deletions
	if totalChanges == 0 {
		return 0.0
	}

	// Calculate ratio of additions to total changes
	additionRatio := float64(pr.Additions) / float64(totalChanges)

	// PRs with >80% additions are likely greenfield
	if additionRatio > 0.80 {
		return additionRatio
	}

	// PRs with balanced changes are likely brownfield
	return additionRatio * 0.5 // Reduce estimate for balanced changes
}

// CalculateSurvivalRate30d calculates the 30-day code survival rate for a PR.
func (g *ResearchGenerator) CalculateSurvivalRate30d(pr models.PullRequest) float64 {
	if pr.MergedAt == nil {
		return 0
	}

	// Check if enough time has passed to measure survival
	thirtyDaysLater := pr.MergedAt.Add(30 * 24 * time.Hour)
	if time.Now().Before(thirtyDaysLater) {
		return 0 // Not enough time has passed
	}

	// Probabilistic survival based on AI ratio and revert status
	// High AI code tends to have lower survival
	if pr.WasReverted {
		return 0 // Code that was reverted has 0% survival
	}

	// Base survival rate decreases with AI ratio
	baseSurvival := 1.0 - (pr.AIRatio * 0.3) // High AI reduces survival by up to 30%

	// Add some randomness (reproducible via seeded RNG)
	variance := (g.rng.Float64() - 0.5) * 0.2 // +/- 10% variance
	survival := baseSurvival + variance

	// Clamp to [0, 1]
	if survival < 0 {
		survival = 0
	}
	if survival > 1 {
		survival = 1
	}

	return survival
}

// ApplyControlVariables applies control variables from seed data to a data point.
func (g *ResearchGenerator) ApplyControlVariables(dp *models.ResearchDataPoint) {
	// Apply author seniority
	if dev, err := g.store.GetDeveloper(dp.AuthorID); err == nil && dev != nil {
		dp.AuthorSeniority = dev.Seniority
	}

	// Apply repo maturity and age
	dp.RepoMaturity = g.getRepoMaturity(dp.RepoName)
	dp.RepoAgeDays = g.getRepoAgeDays(dp.RepoName)
	dp.PrimaryLanguage = g.getRepoPrimaryLanguage(dp.RepoName)
}

// getRepoMaturity returns a string representation of repo maturity.
func (g *ResearchGenerator) getRepoMaturity(repoName string) string {
	if g.seed == nil {
		return "unknown"
	}

	for _, repo := range g.seed.Repositories {
		if repo.RepoName == repoName {
			ageDays := repo.Maturity.AgeDays
			switch {
			case ageDays < 90:
				return "greenfield"
			case ageDays < 180:
				return "developing"
			default:
				return "mature"
			}
		}
	}

	return "unknown"
}

// CountReviewIterations counts the number of review iterations (changes requested cycles).
func (g *ResearchGenerator) CountReviewIterations(reviews []models.ReviewComment) int {
	iterations := 0
	for _, r := range reviews {
		if r.State == models.ReviewStateChangesRequested {
			iterations++
		}
	}
	return iterations
}

// getRepoAgeDays returns the age of the repository in days.
func (g *ResearchGenerator) getRepoAgeDays(repoName string) int {
	if g.seed == nil {
		return 0
	}

	for _, repo := range g.seed.Repositories {
		if repo.RepoName == repoName {
			return repo.Maturity.AgeDays
		}
	}

	return 0
}

// getRepoPrimaryLanguage returns the primary programming language for a repository.
func (g *ResearchGenerator) getRepoPrimaryLanguage(repoName string) string {
	if g.seed == nil {
		return "unknown"
	}

	for _, repo := range g.seed.Repositories {
		if repo.RepoName == repoName {
			return repo.PrimaryLanguage
		}
	}

	return "unknown"
}

// abs returns the absolute value of an integer.
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
