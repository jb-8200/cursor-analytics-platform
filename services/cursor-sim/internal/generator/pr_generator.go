package generator

import (
	"math/rand"
	"sort"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// PRStore defines the interface for PR storage operations.
type PRStore interface {
	GetCommitsByTimeRange(from, to time.Time) []models.Commit
	GetDeveloper(userID string) (*seed.Developer, error)
	AddPR(pr models.PullRequest) error
	GetNextPRNumber(repoName string) int
	AddReviewComment(comment models.ReviewComment) error
}

// CommitCluster represents a group of commits that form a PR.
type CommitCluster struct {
	Commits   []models.Commit
	AuthorID  string
	Branch    string
	RepoName  string
	StartTime time.Time
	EndTime   time.Time
}

// PRGenerator generates pull requests from commit clusters.
type PRGenerator struct {
	seed  *seed.SeedData
	store PRStore
	rng   *rand.Rand
}

// NewPRGenerator creates a new PR generator with a random seed.
func NewPRGenerator(seedData *seed.SeedData, store PRStore) *PRGenerator {
	return NewPRGeneratorWithSeed(seedData, store, time.Now().UnixNano())
}

// NewPRGeneratorWithSeed creates a new PR generator with a specific seed for reproducibility.
func NewPRGeneratorWithSeed(seedData *seed.SeedData, store PRStore, randSeed int64) *PRGenerator {
	return &PRGenerator{
		seed:  seedData,
		store: store,
		rng:   rand.New(rand.NewSource(randSeed)),
	}
}

// GeneratePRsFromCommits generates PRs from commits in the given time range.
func (g *PRGenerator) GeneratePRsFromCommits(from, to time.Time) error {
	commits := g.store.GetCommitsByTimeRange(from, to)
	if len(commits) == 0 {
		return nil
	}

	clusters := g.ClusterCommits(commits)

	for _, cluster := range clusters {
		prNumber := g.store.GetNextPRNumber(cluster.RepoName)
		pr := g.createPRFromCluster(cluster, prNumber)

		if err := g.store.AddPR(pr); err != nil {
			return err
		}
	}

	return nil
}

// ClusterCommits groups commits by author, branch, and time window.
func (g *PRGenerator) ClusterCommits(commits []models.Commit) []CommitCluster {
	// Sort commits by time
	sorted := make([]models.Commit, len(commits))
	copy(sorted, commits)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CommitTs.Before(sorted[j].CommitTs)
	})

	// Cluster key: author:branch:repo
	clusters := make(map[string]*CommitCluster)

	// Maximum time gap between commits in the same PR (24 hours)
	maxGap := 24 * time.Hour

	for _, commit := range sorted {
		key := commit.UserID + ":" + commit.BranchName + ":" + commit.RepoName

		existing, ok := clusters[key]
		if ok {
			// Check if this commit is within the time window of the existing cluster
			if commit.CommitTs.Sub(existing.EndTime) <= maxGap {
				// Add to existing cluster
				existing.Commits = append(existing.Commits, commit)
				existing.EndTime = commit.CommitTs
				continue
			}
		}

		// Start a new cluster (either no existing or gap too large)
		// For simplicity, we overwrite - in practice, we'd want to keep both
		clusters[key] = &CommitCluster{
			Commits:   []models.Commit{commit},
			AuthorID:  commit.UserID,
			Branch:    commit.BranchName,
			RepoName:  commit.RepoName,
			StartTime: commit.CommitTs,
			EndTime:   commit.CommitTs,
		}
	}

	// Convert map to slice
	result := make([]CommitCluster, 0, len(clusters))
	for _, cluster := range clusters {
		result = append(result, *cluster)
	}

	return result
}

// createPRFromCluster creates a PR from a commit cluster.
func (g *PRGenerator) createPRFromCluster(cluster CommitCluster, prNumber int) models.PullRequest {
	// Aggregate metrics from commits
	var totalAdded, totalDeleted int
	var tabLines, composerLines int
	var changedFiles int

	for _, commit := range cluster.Commits {
		totalAdded += commit.TotalLinesAdded
		totalDeleted += commit.TotalLinesDeleted
		tabLines += commit.TabLinesAdded
		composerLines += commit.ComposerLinesAdded
		changedFiles++ // Simplified: count commits as proxy for changed files
	}

	// Calculate AI ratio
	var aiRatio float64
	if totalAdded > 0 {
		aiRatio = float64(tabLines+composerLines) / float64(totalAdded)
	}

	// Get author info from first commit
	var authorEmail, authorName string
	if len(cluster.Commits) > 0 {
		authorEmail = cluster.Commits[0].UserEmail
		authorName = cluster.Commits[0].UserName
	}

	// Determine PR state based on age
	now := time.Now()
	prAge := now.Sub(cluster.EndTime)

	state := models.PRStateOpen
	var mergedAt *time.Time
	var closedAt *time.Time

	// Simulate PR lifecycle based on age
	if prAge > 48*time.Hour {
		// Old PRs are likely merged
		state = models.PRStateMerged
		mergeTime := cluster.EndTime.Add(g.sampleReviewTime())
		mergedAt = &mergeTime
	} else if prAge > 24*time.Hour {
		// PRs between 24-48 hours may be merged (70% chance)
		if g.rng.Float64() < 0.7 {
			state = models.PRStateMerged
			mergeTime := cluster.EndTime.Add(g.sampleReviewTime())
			mergedAt = &mergeTime
		}
	}

	// Generate title
	title := g.generatePRTitle(cluster)

	// Generate body from commit messages
	body := g.generatePRBody(cluster)

	// Determine base branch
	baseBranch := "main"
	if g.seed != nil && len(g.seed.Repositories) > 0 {
		for _, repo := range g.seed.Repositories {
			if repo.RepoName == cluster.RepoName && repo.DefaultBranch != "" {
				baseBranch = repo.DefaultBranch
				break
			}
		}
	}

	return models.PullRequest{
		Number:        prNumber,
		Title:         title,
		Body:          body,
		State:         state,
		AuthorID:      cluster.AuthorID,
		AuthorEmail:   authorEmail,
		AuthorName:    authorName,
		RepoName:      cluster.RepoName,
		BaseBranch:    baseBranch,
		HeadBranch:    cluster.Branch,
		Reviewers:     []string{}, // To be filled by review generator
		Labels:        g.inferLabels(cluster),
		Additions:     totalAdded,
		Deletions:     totalDeleted,
		ChangedFiles:  changedFiles,
		CommitCount:   len(cluster.Commits),
		AIRatio:       aiRatio,
		TabLines:      tabLines,
		ComposerLines: composerLines,
		CreatedAt:     cluster.StartTime,
		UpdatedAt:     cluster.EndTime,
		MergedAt:      mergedAt,
		ClosedAt:      closedAt,
		WasReverted:   false, // To be set by quality outcomes
		IsBugFix:      g.isBugFix(cluster),
	}
}

// sampleReviewTime generates a random review lead time.
func (g *PRGenerator) sampleReviewTime() time.Duration {
	// Default to 8 hours mean if not configured
	meanHours := 8.0
	stdHours := 4.0

	if g.seed != nil && g.seed.PRLifecycle.CycleTimes.ReviewLeadTime.Params.Mean > 0 {
		meanHours = g.seed.PRLifecycle.CycleTimes.ReviewLeadTime.Params.Mean
		stdHours = g.seed.PRLifecycle.CycleTimes.ReviewLeadTime.Params.Std
	}

	// Sample from lognormal distribution
	hours := g.rng.NormFloat64()*stdHours + meanHours
	if hours < 1 {
		hours = 1
	}

	return time.Duration(hours * float64(time.Hour))
}

// generatePRTitle generates a title for the PR.
func (g *PRGenerator) generatePRTitle(cluster CommitCluster) string {
	// Use first commit message as base if available
	if len(cluster.Commits) > 0 && cluster.Commits[0].Message != "" {
		return cluster.Commits[0].Message
	}

	// Use template if available
	if g.seed != nil && len(g.seed.TextTemplates.PRTitles) > 0 {
		template := g.seed.TextTemplates.PRTitles[g.rng.Intn(len(g.seed.TextTemplates.PRTitles))]
		return template
	}

	// Default title based on branch
	return "Update " + cluster.Branch
}

// generatePRBody generates a description for the PR.
func (g *PRGenerator) generatePRBody(cluster CommitCluster) string {
	body := "## Changes\n\n"

	for _, commit := range cluster.Commits {
		if commit.Message != "" {
			body += "- " + commit.Message + "\n"
		}
	}

	body += "\n## Metrics\n\n"
	body += "- Commits: " + string(rune('0'+len(cluster.Commits))) + "\n"

	return body
}

// inferLabels infers labels from the PR content.
func (g *PRGenerator) inferLabels(cluster CommitCluster) []string {
	labels := []string{}

	// Check commit messages for type hints
	for _, commit := range cluster.Commits {
		msg := commit.Message
		if len(msg) > 5 {
			prefix := msg[:5]
			switch {
			case prefix == "feat:" || prefix == "feat(" || prefix == "featu":
				if !contains(labels, "enhancement") {
					labels = append(labels, "enhancement")
				}
			case prefix == "fix:" || prefix == "fix(":
				if !contains(labels, "bug") {
					labels = append(labels, "bug")
				}
			case prefix == "docs:" || prefix == "docs(":
				if !contains(labels, "documentation") {
					labels = append(labels, "documentation")
				}
			}
		}
	}

	return labels
}

// isBugFix determines if the PR is a bug fix.
func (g *PRGenerator) isBugFix(cluster CommitCluster) bool {
	for _, commit := range cluster.Commits {
		msg := commit.Message
		if len(msg) >= 4 && msg[:4] == "fix:" {
			return true
		}
		if len(msg) >= 4 && msg[:4] == "fix(" {
			return true
		}
	}
	return false
}

// contains checks if a slice contains a string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
