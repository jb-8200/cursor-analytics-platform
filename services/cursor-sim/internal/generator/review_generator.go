package generator

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// ReviewStore defines the interface for review storage operations.
type ReviewStore interface {
	GetPR(repoName string, number int) (*models.PullRequest, error)
	GetPRsByRepo(repoName string) []models.PullRequest
	UpdatePR(pr models.PullRequest) error
	AddReviewComment(comment models.ReviewComment) error
	GetReviewComments(repoName string, prNumber int) []models.ReviewComment
	ListDevelopers() []seed.Developer
	GetDeveloper(userID string) (*seed.Developer, error)
}

// ReviewGenerator generates code reviews for pull requests.
type ReviewGenerator struct {
	seed          *seed.SeedData
	store         ReviewStore
	rng           *rand.Rand
	reviewCounter int // Auto-incrementing ID for reviews
}

// NewReviewGenerator creates a new review generator with seedData and RNG (no store required).
// This constructor is used for the simpler GenerateReviewsForPR(pr) []Review pattern.
func NewReviewGenerator(seedData *seed.SeedData, rng *rand.Rand) *ReviewGenerator {
	return &ReviewGenerator{
		seed:          seedData,
		store:         nil,
		rng:           rng,
		reviewCounter: 1,
	}
}

// NewReviewGeneratorWithStore creates a new review generator with a store for persistent operations.
func NewReviewGeneratorWithStore(seedData *seed.SeedData, store ReviewStore) *ReviewGenerator {
	return NewReviewGeneratorWithSeed(seedData, store, time.Now().UnixNano())
}

// NewReviewGeneratorWithSeed creates a new review generator with a specific seed for reproducibility.
func NewReviewGeneratorWithSeed(seedData *seed.SeedData, store ReviewStore, randSeed int64) *ReviewGenerator {
	return &ReviewGenerator{
		seed:          seedData,
		store:         store,
		rng:           rand.New(rand.NewSource(randSeed)),
		reviewCounter: 1,
	}
}

// SelectReviewers selects appropriate reviewers for a PR.
// Prefers developers from the same team, excludes the author.
func (g *ReviewGenerator) SelectReviewers(authorID, authorTeam string) []string {
	developers := g.store.ListDevelopers()

	// Filter eligible reviewers (same team, not author)
	var sameTeam []string
	var otherTeam []string

	for _, dev := range developers {
		if dev.UserID == authorID {
			continue // Exclude author
		}
		if dev.Team == authorTeam {
			sameTeam = append(sameTeam, dev.UserID)
		} else {
			otherTeam = append(otherTeam, dev.UserID)
		}
	}

	// Sort for reproducibility
	sort.Strings(sameTeam)
	sort.Strings(otherTeam)

	// Determine how many reviewers to select
	maxReviewers := 2
	if g.seed != nil && g.seed.PRLifecycle.ReviewPatterns.ReviewerCount.Base > 0 {
		maxReviewers = int(g.seed.PRLifecycle.ReviewPatterns.ReviewerCount.Base)
	}

	// Select reviewers, preferring same team
	var reviewers []string

	// Shuffle candidates using our seeded RNG
	g.rng.Shuffle(len(sameTeam), func(i, j int) {
		sameTeam[i], sameTeam[j] = sameTeam[j], sameTeam[i]
	})

	for _, id := range sameTeam {
		if len(reviewers) >= maxReviewers {
			break
		}
		reviewers = append(reviewers, id)
	}

	// If we need more reviewers, pick from other teams
	if len(reviewers) < maxReviewers {
		g.rng.Shuffle(len(otherTeam), func(i, j int) {
			otherTeam[i], otherTeam[j] = otherTeam[j], otherTeam[i]
		})

		for _, id := range otherTeam {
			if len(reviewers) >= maxReviewers {
				break
			}
			reviewers = append(reviewers, id)
		}
	}

	return reviewers
}

// GenerateReviewForPR generates a review with comments for a PR.
func (g *ReviewGenerator) GenerateReviewForPR(repoName string, prNumber int, reviewerID string) error {
	pr, err := g.store.GetPR(repoName, prNumber)
	if err != nil {
		return err
	}
	if pr == nil {
		return fmt.Errorf("PR not found: %s#%d", repoName, prNumber)
	}

	// Calculate number of comments based on LOC and density
	numComments := g.calculateCommentCount(pr)

	// Generate comments
	now := time.Now()
	for i := 0; i < numComments; i++ {
		comment := g.generateComment(pr, reviewerID, now, i)
		if err := g.store.AddReviewComment(comment); err != nil {
			return err
		}
	}

	return nil
}

// calculateCommentCount determines how many comments to generate.
func (g *ReviewGenerator) calculateCommentCount(pr *models.PullRequest) int {
	baseDensity := 2.0 // Default: 2 comments per 100 LOC
	if g.seed != nil && g.seed.PRLifecycle.ReviewPatterns.CommentsPer100LOC.Base > 0 {
		baseDensity = g.seed.PRLifecycle.ReviewPatterns.CommentsPer100LOC.Base
	}

	// Calculate expected comments based on LOC
	loc := pr.Additions + pr.Deletions
	expected := float64(loc) * baseDensity / 100.0

	// Add some variance
	variance := g.rng.Float64() * 0.5 // 0-50% variance
	count := int(expected * (1 + variance - 0.25))

	// Ensure at least 1 comment
	if count < 1 {
		count = 1
	}

	return count
}

// generateComment generates a single review comment.
func (g *ReviewGenerator) generateComment(pr *models.PullRequest, reviewerID string, baseTime time.Time, index int) models.ReviewComment {
	// Get comment templates
	body := g.selectCommentTemplate()

	// Stagger comment times
	commentTime := baseTime.Add(time.Duration(index*5) * time.Minute)

	return models.ReviewComment{
		ID:        g.rng.Intn(1000000),
		PRNumber:  pr.Number,
		RepoName:  pr.RepoName,
		AuthorID:  reviewerID,
		Body:      body,
		State:     models.ReviewStatePending,
		CreatedAt: commentTime,
	}
}

// selectCommentTemplate selects a random comment template.
func (g *ReviewGenerator) selectCommentTemplate() string {
	if g.seed == nil {
		return "LGTM"
	}

	// Collect all templates
	templates := []string{}
	templates = append(templates, g.seed.TextTemplates.ReviewComments.Style...)
	templates = append(templates, g.seed.TextTemplates.ReviewComments.Logic...)
	templates = append(templates, g.seed.TextTemplates.ReviewComments.Suggestion...)
	templates = append(templates, g.seed.TextTemplates.ReviewComments.Approval...)

	if len(templates) == 0 {
		return "LGTM"
	}

	return templates[g.rng.Intn(len(templates))]
}

// GenerateApprovalDecision generates an approval/rejection decision for a PR.
func (g *ReviewGenerator) GenerateApprovalDecision(pr *models.PullRequest) models.ReviewState {
	// Base approval rate (default 80%)
	approvalRate := 0.8

	// Adjust based on AI ratio if configured
	if g.seed != nil {
		// Higher AI ratio might lead to more scrutiny
		if pr.AIRatio > 0.7 {
			approvalRate *= 0.9 // 10% reduction for high AI ratio
		}

		// Larger PRs might need more iterations
		if pr.Additions > 300 {
			approvalRate *= 0.85 // 15% reduction for large PRs
		}
	}

	// Make decision
	if g.rng.Float64() < approvalRate {
		return models.ReviewStateApproved
	}

	// 70% changes requested, 30% just comment
	if g.rng.Float64() < 0.7 {
		return models.ReviewStateChangesRequested
	}

	return models.ReviewStatePending
}

// SimulateReviewIterations simulates review iterations for a PR.
func (g *ReviewGenerator) SimulateReviewIterations(repoName string, prNumber int) (int, error) {
	pr, err := g.store.GetPR(repoName, prNumber)
	if err != nil {
		return 0, err
	}
	if pr == nil {
		return 0, fmt.Errorf("PR not found: %s#%d", repoName, prNumber)
	}

	// Get max iterations
	maxIterations := 3
	lambda := 1.5 // Average iterations

	if g.seed != nil {
		if params := g.seed.PRLifecycle.ReviewPatterns.Iterations.Params; params != nil {
			if l, ok := params["lambda"]; ok {
				lambda = l
			}
		}
	}

	// Sample from Poisson-like distribution (simplified)
	iterations := 1
	for i := 1; i < maxIterations; i++ {
		// Probability of another iteration decreases
		if g.rng.Float64() < lambda/(float64(i)+lambda) {
			iterations++

			// Select reviewers
			dev, _ := g.store.GetDeveloper(pr.AuthorID)
			team := ""
			if dev != nil {
				team = dev.Team
			}
			reviewers := g.SelectReviewers(pr.AuthorID, team)

			// Generate review for each reviewer
			for _, reviewerID := range reviewers {
				if err := g.GenerateReviewForPR(repoName, prNumber, reviewerID); err != nil {
					return iterations, err
				}
			}
		}
	}

	return iterations, nil
}

// GenerateReviewsForRepo generates reviews for all open PRs in a repository.
func (g *ReviewGenerator) GenerateReviewsForRepo(repoName string) (int, error) {
	prs := g.store.GetPRsByRepo(repoName)

	reviewedCount := 0
	for _, pr := range prs {
		// Skip non-open PRs
		if pr.State != models.PRStateOpen {
			continue
		}

		// Get author's team
		dev, _ := g.store.GetDeveloper(pr.AuthorID)
		team := ""
		if dev != nil {
			team = dev.Team
		}

		// Select reviewers
		reviewers := g.SelectReviewers(pr.AuthorID, team)

		// Generate reviews
		for _, reviewerID := range reviewers {
			if err := g.GenerateReviewForPR(repoName, pr.Number, reviewerID); err != nil {
				return reviewedCount, err
			}
		}

		reviewedCount++
	}

	return reviewedCount, nil
}

// GenerateReviewsForPR generates reviews for a given pull request.
// This method doesn't require a store - it uses seed data developers directly.
// Returns a slice of Review objects with timing based on PR lifecycle.
// State distribution: 70% approved, 20% changes_requested, 10% pending.
func (g *ReviewGenerator) GenerateReviewsForPR(pr models.PullRequest) []models.Review {
	// Get available reviewers from seed data (excluding PR author)
	availableReviewers := g.getAvailableReviewers(pr.AuthorEmail)
	if len(availableReviewers) == 0 {
		return []models.Review{}
	}

	// Determine number of reviewers (1-3, capped by available)
	numReviewers := 1 + g.rng.Intn(3) // 1-3 reviewers
	if numReviewers > len(availableReviewers) {
		numReviewers = len(availableReviewers)
	}

	// Shuffle and select reviewers
	g.rng.Shuffle(len(availableReviewers), func(i, j int) {
		availableReviewers[i], availableReviewers[j] = availableReviewers[j], availableReviewers[i]
	})
	selectedReviewers := availableReviewers[:numReviewers]

	// Calculate PR end time for review timing
	prEndTime := g.calculatePREndTime(pr)

	reviews := make([]models.Review, 0, numReviewers)
	for _, reviewerEmail := range selectedReviewers {
		review := g.generateReview(pr, reviewerEmail, prEndTime)
		reviews = append(reviews, review)
	}

	return reviews
}

// getAvailableReviewers returns emails of developers who can review (not the author).
func (g *ReviewGenerator) getAvailableReviewers(authorEmail string) []string {
	if g.seed == nil || len(g.seed.Developers) == 0 {
		return nil
	}

	var reviewers []string
	for _, dev := range g.seed.Developers {
		if dev.Email != authorEmail {
			reviewers = append(reviewers, dev.Email)
		}
	}
	return reviewers
}

// calculatePREndTime returns the end time to use for review timing.
// Uses MergedAt for merged PRs, ClosedAt for closed PRs, or 7 days from creation for open PRs.
func (g *ReviewGenerator) calculatePREndTime(pr models.PullRequest) time.Time {
	if pr.MergedAt != nil {
		return *pr.MergedAt
	}
	if pr.ClosedAt != nil {
		return *pr.ClosedAt
	}
	// For open PRs, use 7 days from creation as max review time
	return pr.CreatedAt.Add(7 * 24 * time.Hour)
}

// generateReview creates a single review for a PR.
func (g *ReviewGenerator) generateReview(pr models.PullRequest, reviewerEmail string, prEndTime time.Time) models.Review {
	// Generate review timing between PR creation and end time
	duration := prEndTime.Sub(pr.CreatedAt)
	if duration <= 0 {
		duration = time.Hour // Minimum 1 hour window
	}

	// Reviews typically happen within the first half of the PR lifecycle
	reviewOffset := time.Duration(g.rng.Int63n(int64(duration)))
	submittedAt := pr.CreatedAt.Add(reviewOffset)

	// Generate review state: 70% approved, 20% changes_requested, 10% pending
	state := g.generateReviewState()

	// Generate review ID
	reviewID := g.reviewCounter
	g.reviewCounter++

	// Create the review
	review := models.Review{
		ID:          reviewID,
		PRID:        pr.Number,
		Reviewer:    reviewerEmail,
		State:       state,
		SubmittedAt: submittedAt,
	}

	// Add body and comments for non-approved reviews
	if state != models.ReviewStateApproved {
		review.Body = g.generateReviewBody(state)
		review.Comments = g.generateReviewComments(pr, reviewerEmail, submittedAt)
	} else {
		// Approved reviews often have a short positive message
		review.Body = g.generateApprovalBody()
	}

	return review
}

// generateReviewState generates a review state based on distribution:
// 70% approved, 20% changes_requested, 10% pending.
func (g *ReviewGenerator) generateReviewState() models.ReviewState {
	roll := g.rng.Float64()
	if roll < 0.70 {
		return models.ReviewStateApproved
	} else if roll < 0.90 {
		return models.ReviewStateChangesRequested
	}
	return models.ReviewStatePending
}

// generateReviewBody generates appropriate body text based on review state.
func (g *ReviewGenerator) generateReviewBody(state models.ReviewState) string {
	if g.seed == nil || len(g.seed.TextTemplates.ReviewComments.Suggestion) == 0 {
		if state == models.ReviewStateChangesRequested {
			return "Please address the requested changes."
		}
		return "Reviewing..."
	}

	templates := g.seed.TextTemplates.ReviewComments.Suggestion
	if state == models.ReviewStateChangesRequested && len(g.seed.TextTemplates.ReviewComments.Logic) > 0 {
		templates = g.seed.TextTemplates.ReviewComments.Logic
	}

	return templates[g.rng.Intn(len(templates))]
}

// generateApprovalBody generates body text for approved reviews.
func (g *ReviewGenerator) generateApprovalBody() string {
	if g.seed == nil || len(g.seed.TextTemplates.ReviewComments.Approval) == 0 {
		return "LGTM!"
	}

	templates := g.seed.TextTemplates.ReviewComments.Approval
	return templates[g.rng.Intn(len(templates))]
}

// generateReviewComments generates inline comments for non-approved reviews.
func (g *ReviewGenerator) generateReviewComments(pr models.PullRequest, reviewerEmail string, baseTime time.Time) []models.ReviewComment {
	// Generate 0-5 comments for non-approved reviews
	numComments := g.rng.Intn(6) // 0-5 comments
	if numComments == 0 {
		return nil
	}

	comments := make([]models.ReviewComment, numComments)
	for i := 0; i < numComments; i++ {
		comments[i] = models.ReviewComment{
			ID:        g.rng.Intn(1000000),
			PRNumber:  pr.Number,
			RepoName:  pr.RepoName,
			AuthorID:  reviewerEmail,
			Body:      g.selectCommentTemplate(),
			State:     models.ReviewStateChangesRequested,
			CreatedAt: baseTime.Add(time.Duration(i) * time.Minute),
		}
	}

	return comments
}
