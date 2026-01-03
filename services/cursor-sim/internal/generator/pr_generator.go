package generator

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// PRStore defines the storage interface needed by PRGenerator for backwards compatibility.
type PRStore interface {
	GetCommitsByTimeRange(from, to time.Time) []models.Commit
	AddPR(pr models.PullRequest) error
}

// PRGenerator groups commits into PRs using session-based rules.
type PRGenerator struct {
	seed      *seed.SeedData
	store     PRStore  // Optional: for backwards-compatible GeneratePRsFromCommits
	rng       *rand.Rand
	prCounter int
}

// NewPRGenerator creates a new PR generator with seed data and RNG.
func NewPRGenerator(seedData *seed.SeedData, rng *rand.Rand) *PRGenerator {
	return &PRGenerator{
		seed:      seedData,
		store:     nil,  // No automatic storage
		rng:       rng,
		prCounter: 1,
	}
}

// NewPRGeneratorWithSeed creates a new PR generator with a specific random seed for reproducibility.
// The store parameter enables backwards-compatible GeneratePRsFromCommits(from, to) method.
func NewPRGeneratorWithSeed(seedData *seed.SeedData, store PRStore, randSeed int64) *PRGenerator {
	return &PRGenerator{
		seed:      seedData,
		store:     store,
		rng:       rand.New(rand.NewSource(randSeed)),
		prCounter: 1,
	}
}

// GeneratePRsFromCommits generates PRs from commits in the given time range and stores them.
// This is the backwards-compatible API that requires a store to be provided via NewPRGeneratorWithSeed.
func (g *PRGenerator) GeneratePRsFromCommits(from, to time.Time) error {
	if g.store == nil {
		return fmt.Errorf("store is required for GeneratePRsFromCommits - use NewPRGeneratorWithSeed")
	}

	// Fetch commits from store
	commits := g.store.GetCommitsByTimeRange(from, to)

	// Generate PRs
	prs := g.GroupCommitsIntoPRs(commits)

	// Store PRs
	for _, pr := range prs {
		if err := g.store.AddPR(pr); err != nil {
			return fmt.Errorf("failed to store PR %d: %w", pr.Number, err)
		}
	}

	return nil
}

// GroupCommitsIntoPRs groups commits into PRs based on (repo, branch, author) and session rules.
// Returns a list of PR envelopes with aggregated metrics.
func (g *PRGenerator) GroupCommitsIntoPRs(commits []models.Commit) []models.PullRequest {
	// Sort commits by timestamp
	sorted := make([]models.Commit, len(commits))
	copy(sorted, commits)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CommitTs.Before(sorted[j].CommitTs)
	})

	// Group by (repo, branch, author)
	groups := make(map[string][]models.Commit)
	for _, commit := range sorted {
		key := fmt.Sprintf("%s:%s:%s", commit.RepoName, commit.BranchName, commit.UserID)
		groups[key] = append(groups[key], commit)
	}

	// Generate PRs for each group
	var prs []models.PullRequest
	for _, groupCommits := range groups {
		if len(groupCommits) == 0 {
			continue
		}

		// Get developer info
		firstCommit := groupCommits[0]
		developer := g.getDeveloper(firstCommit.UserID)
		repository := g.getRepository(firstCommit.RepoName)

		// Create sessions and PRs
		sessionPRs := g.createPRsFromCommits(groupCommits, developer, repository)
		prs = append(prs, sessionPRs...)
	}

	return prs
}

// createPRsFromCommits creates PRs from a group of commits using session-based rules.
func (g *PRGenerator) createPRsFromCommits(commits []models.Commit, developer seed.Developer, repository seed.Repository) []models.PullRequest {
	if len(commits) == 0 {
		return nil
	}

	var prs []models.PullRequest
	var currentSession *Session

	for i, commit := range commits {
		// Start new session if needed
		if currentSession == nil {
			currentSession = StartSession(
				developer,
				repository,
				commit.BranchName,
				commit.CommitTs,
				g.rng,
			)
		}

		// Add commit to session
		currentSession.AddCommit(commit)

		// Check if session should close
		shouldClose := false
		if i < len(commits)-1 {
			nextCommitTime := commits[i+1].CommitTs
			shouldClose = currentSession.ShouldClose(nextCommitTime, g.rng)
		} else {
			// Last commit always closes the session
			shouldClose = true
		}

		if shouldClose {
			// Finalize PR from session
			pr := g.finalizePR(currentSession)
			prs = append(prs, pr)
			currentSession = nil
		}
	}

	return prs
}

// finalizePR creates a PR envelope from a completed session.
func (g *PRGenerator) finalizePR(session *Session) models.PullRequest {
	commits := session.Commits

	// Aggregate metrics
	var totalAdditions, totalDeletions, totalTabLines, totalComposerLines int
	var firstCommitTime, lastCommitTime time.Time

	for i, commit := range commits {
		totalAdditions += commit.TotalLinesAdded
		totalDeletions += commit.TotalLinesDeleted
		totalTabLines += commit.TabLinesAdded
		totalComposerLines += commit.ComposerLinesAdded

		if i == 0 {
			firstCommitTime = commit.CommitTs
		}
		lastCommitTime = commit.CommitTs
	}

	// Calculate AI ratio
	var aiRatio float64
	if totalAdditions > 0 {
		aiRatio = float64(totalTabLines+totalComposerLines) / float64(totalAdditions)
	}

	// Create PR
	prNumber := g.prCounter
	g.prCounter++

	pr := models.PullRequest{
		Number:        prNumber,
		Title:         fmt.Sprintf("PR #%d: %s", prNumber, session.Branch),
		Body:          fmt.Sprintf("Auto-generated PR from %d commits", len(commits)),
		State:         models.PRStateOpen,
		AuthorID:      session.Developer.UserID,
		AuthorEmail:   session.Developer.Email,
		AuthorName:    session.Developer.Name,
		RepoName:      session.Repo.RepoName,
		BaseBranch:    session.Repo.DefaultBranch,
		HeadBranch:    session.Branch,
		Additions:     totalAdditions,
		Deletions:     totalDeletions,
		CommitCount:   len(commits),
		AIRatio:       aiRatio,
		TabLines:      totalTabLines,
		ComposerLines: totalComposerLines,
		CreatedAt:     firstCommitTime,
		UpdatedAt:     lastCommitTime,
	}

	return pr
}

// getDeveloper returns the developer with the given userID.
func (g *PRGenerator) getDeveloper(userID string) seed.Developer {
	for _, dev := range g.seed.Developers {
		if dev.UserID == userID {
			return dev
		}
	}
	// Return default developer if not found
	return seed.Developer{
		UserID:    userID,
		Seniority: "mid",
		WorkingHoursBand: seed.WorkingHours{
			Start: 9,
			End:   18,
		},
	}
}

// getRepository returns the repository with the given name.
func (g *PRGenerator) getRepository(repoName string) seed.Repository {
	for _, repo := range g.seed.Repositories {
		if repo.RepoName == repoName {
			return repo
		}
	}
	// Return default repository if not found
	return seed.Repository{
		RepoName:      repoName,
		DefaultBranch: "main",
	}
}
