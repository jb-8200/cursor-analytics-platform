package generator

import (
	"context"
	"crypto/sha1"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// Store defines the interface for storing generated commits.
type Store interface {
	AddCommit(commit models.Commit) error
}

// CommitGenerator generates synthetic commits based on seed data.
type CommitGenerator struct {
	seed     *seed.SeedData
	store    Store
	velocity *VelocityConfig
	rng      *rand.Rand
}

// NewCommitGenerator creates a new commit generator with a random seed.
func NewCommitGenerator(seedData *seed.SeedData, store Store, velocity string) *CommitGenerator {
	return NewCommitGeneratorWithSeed(seedData, store, velocity, time.Now().UnixNano())
}

// NewCommitGeneratorWithSeed creates a new commit generator with a specific seed for reproducibility.
func NewCommitGeneratorWithSeed(seedData *seed.SeedData, store Store, velocity string, randSeed int64) *CommitGenerator {
	return &CommitGenerator{
		seed:     seedData,
		store:    store,
		velocity: NewVelocityConfig(velocity),
		rng:      rand.New(rand.NewSource(randSeed)),
	}
}

// GenerateCommits generates commits for the specified number of days.
// Uses Poisson process for timing and lognormal distribution for commit sizes.
func (g *CommitGenerator) GenerateCommits(ctx context.Context, days int) error {
	startTime := time.Now().AddDate(0, 0, -days)

	for _, dev := range g.seed.Developers {
		if err := g.generateForDeveloper(ctx, dev, startTime); err != nil {
			return err
		}
	}

	return nil
}

// generateForDeveloper generates commits for a single developer over the time range.
func (g *CommitGenerator) generateForDeveloper(ctx context.Context, dev seed.Developer, startTime time.Time) error {
	// Calculate commit rate for this developer
	commitsPerDay := g.velocity.CommitsPerDay(dev.PRBehavior.PRsPerWeek)
	if commitsPerDay <= 0 {
		return nil // Developer doesn't commit
	}

	// Convert to rate per hour for finer granularity
	commitsPerHour := commitsPerDay / 24.0

	// Generate commits using Poisson process (exponential inter-arrival times)
	current := startTime
	now := time.Now()

	for current.Before(now) {
		// Select from context
		if err := ctx.Err(); err != nil {
			return err
		}

		// Wait time follows exponential distribution
		// For Poisson process with rate λ, inter-arrival times are Exp(λ)
		waitHours := g.exponential(commitsPerHour)
		current = current.Add(time.Duration(waitHours * float64(time.Hour)))

		if current.After(now) {
			break
		}

		// Generate and store the commit
		commit := g.generateCommit(dev, current)
		if err := g.store.AddCommit(commit); err != nil {
			return err
		}
	}

	return nil
}

// generateCommit creates a single commit for the developer at the given time.
func (g *CommitGenerator) generateCommit(dev seed.Developer, timestamp time.Time) models.Commit {
	// Generate commit size using lognormal distribution
	// Mean size from developer's PR behavior, with variability
	meanSize := float64(dev.PRBehavior.AvgPRSizeLOC) / 3.0 // Divide by ~commits per PR
	linesAdded := g.lognormal(meanSize, meanSize*0.5)
	if linesAdded < 1 {
		linesAdded = 1
	}

	// Split lines between AI and non-AI based on acceptance rate
	// Higher acceptance rate = more AI usage
	aiRatio := dev.AcceptanceRate
	aiLines := int(float64(linesAdded) * aiRatio)

	// Split AI lines between Tab and Composer
	// Tab typically 60-80% of AI lines, Composer 20-40%
	tabRatio := 0.6 + g.rng.Float64()*0.2 // 60-80%
	tabLines := int(float64(aiLines) * tabRatio)
	composerLines := aiLines - tabLines

	// Remaining lines are non-AI
	nonAILines := linesAdded - aiLines

	// Ensure non-negative values
	if tabLines < 0 {
		tabLines = 0
	}
	if composerLines < 0 {
		composerLines = 0
	}
	if nonAILines < 0 {
		nonAILines = 0
	}

	// Generate deletions (typically 10-30% of additions)
	deletionRatio := 0.1 + g.rng.Float64()*0.2
	linesDeleted := int(float64(linesAdded) * deletionRatio)

	// Split deletions proportionally
	tabDeleted := int(float64(linesDeleted) * tabRatio)
	composerDeleted := int(float64(linesDeleted) * (1.0 - tabRatio))
	nonAIDeleted := linesDeleted - tabDeleted - composerDeleted
	if nonAIDeleted < 0 {
		nonAIDeleted = 0
	}

	// Select repository (prefer ones matching developer's team)
	repo := g.selectRepository(dev)

	// Generate commit message from templates
	message := g.generateCommitMessage()

	// Generate commit hash
	commitHash := g.generateCommitHash(dev, timestamp)

	return models.Commit{
		CommitHash:           commitHash,
		UserID:               dev.UserID,
		UserEmail:            dev.Email,
		UserName:             dev.Name,
		RepoName:             repo.RepoName,
		BranchName:           g.selectBranch(repo),
		IsPrimaryBranch:      g.rng.Float64() < 0.3, // 30% on primary branch
		TotalLinesAdded:      linesAdded,
		TotalLinesDeleted:    linesDeleted,
		TabLinesAdded:        tabLines,
		TabLinesDeleted:      tabDeleted,
		ComposerLinesAdded:   composerLines,
		ComposerLinesDeleted: composerDeleted,
		NonAILinesAdded:      nonAILines,
		NonAILinesDeleted:    nonAIDeleted,
		Message:              message,
		CommitTs:             timestamp,
		CreatedAt:            timestamp,
	}
}

// exponential generates a random value from an exponential distribution with the given rate.
func (g *CommitGenerator) exponential(rate float64) float64 {
	if rate <= 0 {
		return 0
	}
	// Inverse transform: -ln(U) / λ where U ~ Uniform(0,1)
	return -math.Log(1.0-g.rng.Float64()) / rate
}

// lognormal generates a random value from a lognormal distribution.
func (g *CommitGenerator) lognormal(mean, stddev float64) int {
	// Generate from normal distribution
	normal := g.rng.NormFloat64()*stddev + mean

	// Convert to lognormal
	value := math.Exp(normal / (mean + 1.0))

	return int(value * mean)
}

// selectRepository selects a repository for the commit, preferring repos matching the developer's team.
func (g *CommitGenerator) selectRepository(dev seed.Developer) seed.Repository {
	if len(g.seed.Repositories) == 0 {
		// Return a default repo if none defined
		return seed.Repository{
			RepoName:      "default/repo",
			DefaultBranch: "main",
		}
	}

	// Try to find a repo matching the developer's team
	for _, repo := range g.seed.Repositories {
		for _, team := range repo.Teams {
			if team == dev.Team {
				return repo
			}
		}
	}

	// Fall back to random repository
	return g.seed.Repositories[g.rng.Intn(len(g.seed.Repositories))]
}

// selectBranch selects a branch name for the commit.
func (g *CommitGenerator) selectBranch(repo seed.Repository) string {
	// 30% chance of primary branch, 70% feature branch
	if g.rng.Float64() < 0.3 {
		if repo.DefaultBranch != "" {
			return repo.DefaultBranch
		}
		return "main"
	}

	// Generate feature branch name
	branchTypes := []string{"feature", "bugfix", "refactor", "chore"}
	branchType := branchTypes[g.rng.Intn(len(branchTypes))]
	return fmt.Sprintf("%s/task-%d", branchType, g.rng.Intn(1000))
}

// generateCommitMessage generates a commit message from templates.
func (g *CommitGenerator) generateCommitMessage() string {
	if len(g.seed.TextTemplates.CommitMessages.Feature) == 0 {
		return "Update code"
	}

	// Select message type
	messageType := g.rng.Intn(4)

	var templates []string
	switch messageType {
	case 0:
		templates = g.seed.TextTemplates.CommitMessages.Feature
	case 1:
		templates = g.seed.TextTemplates.CommitMessages.Bugfix
	case 2:
		templates = g.seed.TextTemplates.CommitMessages.Refactor
	case 3:
		templates = g.seed.TextTemplates.CommitMessages.Chore
	}

	if len(templates) == 0 {
		return "Update code"
	}

	return templates[g.rng.Intn(len(templates))]
}

// generateCommitHash generates a SHA-1 like commit hash.
func (g *CommitGenerator) generateCommitHash(dev seed.Developer, timestamp time.Time) string {
	data := fmt.Sprintf("%s-%s-%d-%d",
		dev.UserID,
		dev.Email,
		timestamp.Unix(),
		g.rng.Int63())

	hash := sha1.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)[:12] // First 12 chars like short git hash
}
