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
	prs := make([]models.PullRequest, 0)
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
	commitIDs := make([]string, len(commits))

	for i, commit := range commits {
		totalAdditions += commit.TotalLinesAdded
		totalDeletions += commit.TotalLinesDeleted
		totalTabLines += commit.TabLinesAdded
		totalComposerLines += commit.ComposerLinesAdded
		commitIDs[i] = commit.CommitHash

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

	// Assign PR status: 85% merged, 10% closed, 5% open
	state := g.assignPRStatus()

	// Generate PR timestamps
	createdAt := firstCommitTime.Add(-time.Duration(g.rng.Intn(60)) * time.Minute) // PR created up to 1h before first commit
	var mergedAt, closedAt *time.Time

	if state == models.PRStateMerged {
		// Merge happens 1-7 days after creation
		daysToMerge := 1 + g.rng.Intn(7)
		merged := createdAt.Add(time.Duration(daysToMerge) * 24 * time.Hour)
		mergedAt = &merged
	} else if state == models.PRStateClosed {
		// Close happens 1-14 days after creation
		daysToClosed := 1 + g.rng.Intn(14)
		closed := createdAt.Add(time.Duration(daysToClosed) * 24 * time.Hour)
		closedAt = &closed
	}

	// Generate PR title from branch name and commits
	title := g.generatePRTitle(session.Branch, commits)

	pr := models.PullRequest{
		Number:        prNumber,
		Title:         title,
		Body:          fmt.Sprintf("Auto-generated PR from %d commits", len(commits)),
		State:         state,
		AuthorID:      session.Developer.UserID,
		AuthorEmail:   session.Developer.Email,
		AuthorName:    session.Developer.Name,
		RepoName:      session.Repo.RepoName,
		BaseBranch:    session.Repo.DefaultBranch,
		HeadBranch:    session.Branch,
		CommitIDs:     commitIDs,
		Additions:     totalAdditions,
		Deletions:     totalDeletions,
		CommitCount:   len(commits),
		AIRatio:       aiRatio,
		TabLines:      totalTabLines,
		ComposerLines: totalComposerLines,
		CreatedAt:     createdAt,
		UpdatedAt:     lastCommitTime,
		MergedAt:      mergedAt,
		ClosedAt:      closedAt,
	}

	return pr
}

// assignPRStatus assigns a PR status based on distribution: 85% merged, 10% closed, 5% open.
func (g *PRGenerator) assignPRStatus() models.PRState {
	roll := g.rng.Float64()
	if roll < 0.85 {
		return models.PRStateMerged
	} else if roll < 0.95 {
		return models.PRStateClosed
	}
	return models.PRStateOpen
}

// generatePRTitle creates a descriptive PR title from branch name and commit messages.
func (g *PRGenerator) generatePRTitle(branch string, commits []models.Commit) string {
	// Extract feature name from branch (e.g., "feature/auth-login" -> "Auth login")
	// If branch has format "type/name", use the name part
	parts := []rune(branch)
	var featureName string

	slashIdx := -1
	for i, ch := range parts {
		if ch == '/' {
			slashIdx = i
			break
		}
	}

	if slashIdx >= 0 && slashIdx+1 < len(parts) {
		featureName = string(parts[slashIdx+1:])
		// Replace dashes/underscores with spaces and capitalize first letter
		nameRunes := []rune(featureName)
		if len(nameRunes) > 0 {
			nameRunes[0] = []rune(string(nameRunes[0]))[0] // Keep original case
			for i, ch := range nameRunes {
				if ch == '-' || ch == '_' {
					nameRunes[i] = ' '
				}
			}
			featureName = string(nameRunes)
		}
	} else {
		featureName = branch
	}

	// Use first commit message as base if available
	if len(commits) > 0 && commits[0].Message != "" {
		return commits[0].Message
	}

	return fmt.Sprintf("Implement %s", featureName)
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
