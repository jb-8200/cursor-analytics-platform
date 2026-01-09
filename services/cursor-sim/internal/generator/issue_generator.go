package generator

import (
	"math/rand"
	"strings"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// IssueStore defines the interface for storing issues.
type IssueStore interface {
	StoreIssue(issue models.Issue) error
}

// IssueGenerator generates issues linked to pull requests.
// - 40% of merged PRs close an issue
// - 10% of issues remain open
// - Issue created before PR
// - Labels: bug, feature, enhancement
type IssueGenerator struct {
	seed       *seed.SeedData
	rng        *rand.Rand
	nextNumber int
	store      IssueStore
}

// NewIssueGenerator creates a new IssueGenerator with seed data and RNG.
func NewIssueGenerator(seedData *seed.SeedData, rng *rand.Rand) *IssueGenerator {
	return &IssueGenerator{
		seed:       seedData,
		rng:        rng,
		nextNumber: 1,
	}
}

// NewIssueGeneratorWithStore creates a new IssueGenerator with storage capability.
func NewIssueGeneratorWithStore(seedData *seed.SeedData, store IssueStore, randSeed int64) *IssueGenerator {
	return &IssueGenerator{
		seed:       seedData,
		rng:        rand.New(rand.NewSource(randSeed)),
		nextNumber: 1,
		store:      store,
	}
}

// GenerateAndStoreIssuesForPRs generates issues for PRs and stores them.
// Returns the count of issues generated and any error encountered.
func (g *IssueGenerator) GenerateAndStoreIssuesForPRs(prs []models.PullRequest, repoName string) (int, error) {
	if g.store == nil {
		return 0, nil
	}

	issues := g.GenerateIssuesForPRs(prs, repoName)
	for _, issue := range issues {
		if err := g.store.StoreIssue(issue); err != nil {
			return 0, err
		}
	}
	return len(issues), nil
}

// GenerateIssuesForPRs generates issues for a list of PRs.
// - 40% of merged PRs will have an associated issue
// - 10% of generated issues remain open
func (g *IssueGenerator) GenerateIssuesForPRs(prs []models.PullRequest, repoName string) []models.Issue {
	if len(prs) == 0 {
		return []models.Issue{}
	}

	var issues []models.Issue

	// Filter to merged PRs only
	var mergedPRs []models.PullRequest
	for _, pr := range prs {
		if pr.State == models.PRStateMerged {
			mergedPRs = append(mergedPRs, pr)
		}
	}

	// 40% of merged PRs close an issue
	targetIssueCount := int(float64(len(mergedPRs)) * 0.40)
	if targetIssueCount < 1 && len(mergedPRs) > 0 {
		// Ensure at least 1 issue if there are merged PRs
		if g.rng.Float64() < 0.40 {
			targetIssueCount = 1
		}
	}

	// Shuffle merged PRs to randomly select which ones get issues
	shuffledPRs := make([]models.PullRequest, len(mergedPRs))
	copy(shuffledPRs, mergedPRs)
	g.rng.Shuffle(len(shuffledPRs), func(i, j int) {
		shuffledPRs[i], shuffledPRs[j] = shuffledPRs[j], shuffledPRs[i]
	})

	// Generate issues for selected PRs
	for i := 0; i < targetIssueCount && i < len(shuffledPRs); i++ {
		pr := shuffledPRs[i]
		issue := g.generateIssueForPR(pr, repoName)
		issues = append(issues, issue)
	}

	return issues
}

// generateIssueForPR creates an issue linked to a specific PR.
func (g *IssueGenerator) generateIssueForPR(pr models.PullRequest, repoName string) models.Issue {
	// 10% of issues remain open
	isOpen := g.rng.Float64() < 0.10

	// Issue created 1-7 days before PR created_at
	daysBeforePR := time.Duration(g.rng.Intn(7)+1) * 24 * time.Hour
	issueCreatedAt := pr.CreatedAt.Add(-daysBeforePR)

	// Determine issue state and closed time
	var state models.IssueState
	var closedAt *time.Time
	var closedByPRID *int

	if isOpen {
		state = models.IssueStateOpen
	} else {
		state = models.IssueStateClosed
		// Closed when PR was merged
		if pr.MergedAt != nil {
			closedAt = pr.MergedAt
		}
		prNumber := pr.Number
		closedByPRID = &prNumber
	}

	// Generate title from PR title
	title := g.generateIssueTitleFromPR(pr.Title)

	// Select random author from developers
	author := g.selectRandomDeveloper()

	// Assign labels
	labels := g.assignLabels()

	issue := models.Issue{
		Number:       g.nextNumber,
		Title:        title,
		Body:         g.generateIssueBody(title),
		State:        state,
		AuthorID:     author,
		RepoName:     repoName,
		Labels:       labels,
		Assignees:    []string{},
		CreatedAt:    issueCreatedAt,
		UpdatedAt:    issueCreatedAt,
		ClosedAt:     closedAt,
		ClosedByPRID: closedByPRID,
	}

	g.nextNumber++
	return issue
}

// generateIssueTitleFromPR creates an issue title based on PR title.
func (g *IssueGenerator) generateIssueTitleFromPR(prTitle string) string {
	// Common prefixes to try to extract the core issue
	prefixes := []string{"Fix ", "Add ", "Update ", "Implement ", "Refactor "}

	for _, prefix := range prefixes {
		if strings.HasPrefix(prTitle, prefix) {
			// Remove the prefix and return the rest
			return strings.TrimPrefix(prTitle, prefix)
		}
	}

	// If no prefix match, return the title as-is
	return prTitle
}

// generateIssueBody creates a simple issue body.
func (g *IssueGenerator) generateIssueBody(title string) string {
	templates := []string{
		"## Description\n\nThis issue tracks: " + title,
		"## Problem\n\n" + title + " needs to be addressed.",
		"## Request\n\nPlease implement: " + title,
	}
	return templates[g.rng.Intn(len(templates))]
}

// selectRandomDeveloper picks a random developer from seed data.
func (g *IssueGenerator) selectRandomDeveloper() string {
	if len(g.seed.Developers) == 0 {
		return "unknown@example.com"
	}
	dev := g.seed.Developers[g.rng.Intn(len(g.seed.Developers))]
	return dev.Email
}

// assignLabels randomly assigns 1-2 labels from bug, feature, enhancement.
func (g *IssueGenerator) assignLabels() []string {
	allLabels := []string{"bug", "feature", "enhancement"}

	// Shuffle and pick 1-2 labels
	g.rng.Shuffle(len(allLabels), func(i, j int) {
		allLabels[i], allLabels[j] = allLabels[j], allLabels[i]
	})

	numLabels := g.rng.Intn(2) + 1 // 1 or 2 labels
	return allLabels[:numLabels]
}
