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
	dp := &models.ResearchDataPoint{
		CommitHash:    commit.CommitHash,
		AuthorID:      commit.UserID,
		RepoName:      commit.RepoName,
		AIRatio:       commit.AIRatio(),
		TabLines:      commit.TabLinesAdded,
		ComposerLines: commit.ComposerLinesAdded,
		Additions:     commit.TotalLinesAdded,
		Deletions:     commit.TotalLinesDeleted,
		Timestamp:     commit.CommitTs,
	}

	// Try to find associated PR
	prKey := commit.RepoName + ":" + commit.BranchName
	if pr, ok := prMap[prKey]; ok {
		dp.PRNumber = pr.Number
		dp.FilesChanged = pr.ChangedFiles
		dp.AIRatio = pr.AIRatio // Use PR-level AI ratio if available

		// Calculate cycle times
		commits := g.getCommitsForPR(pr)
		dp.CodingLeadTimeHours = g.CalculateCodingLeadTime(commits)
		dp.ReviewLeadTimeHours = g.CalculateReviewLeadTime(*pr)

		// Get review iterations
		reviews := g.store.GetReviewComments(pr.RepoName, pr.Number)
		dp.ReviewIterations = g.CountReviewIterations(reviews)

		// Quality outcomes
		dp.WasReverted = pr.WasReverted
		dp.RequiredHotfix = pr.IsBugFix && pr.WasReverted

		// Determine greenfield status
		dp.IsGreenfield = g.IsGreenfield(*pr)
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

// ApplyControlVariables applies control variables from seed data to a data point.
func (g *ResearchGenerator) ApplyControlVariables(dp *models.ResearchDataPoint) {
	// Apply author seniority
	if dev, err := g.store.GetDeveloper(dp.AuthorID); err == nil && dev != nil {
		dp.AuthorSeniority = dev.Seniority
	}

	// Apply repo maturity
	dp.RepoMaturity = g.getRepoMaturity(dp.RepoName)
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
