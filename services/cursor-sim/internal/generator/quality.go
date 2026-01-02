package generator

import (
	"math/rand"
	"strings"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// AIRatioCategory represents the AI usage level of a PR.
type AIRatioCategory string

const (
	AIRatioCategoryLow    AIRatioCategory = "low"
	AIRatioCategoryMedium AIRatioCategory = "medium"
	AIRatioCategoryHigh   AIRatioCategory = "high"
)

// QualityStore defines the interface for quality outcome storage operations.
type QualityStore interface {
	GetPR(repoName string, number int) (*models.PullRequest, error)
	GetPRsByRepo(repoName string) []models.PullRequest
	UpdatePR(pr models.PullRequest) error
}

// QualityGenerator generates quality outcomes for pull requests.
type QualityGenerator struct {
	seed  *seed.SeedData
	store QualityStore
	rng   *rand.Rand
}

// NewQualityGenerator creates a new quality generator with a random seed.
func NewQualityGenerator(seedData *seed.SeedData, store QualityStore) *QualityGenerator {
	return NewQualityGeneratorWithSeed(seedData, store, time.Now().UnixNano())
}

// NewQualityGeneratorWithSeed creates a new quality generator with a specific seed for reproducibility.
func NewQualityGeneratorWithSeed(seedData *seed.SeedData, store QualityStore, randSeed int64) *QualityGenerator {
	return &QualityGenerator{
		seed:  seedData,
		store: store,
		rng:   rand.New(rand.NewSource(randSeed)),
	}
}

// CategorizeAIRatio returns the AI ratio category based on the configured bands.
func (g *QualityGenerator) CategorizeAIRatio(aiRatio float64) AIRatioCategory {
	if g.seed == nil {
		// Default bands
		if aiRatio < 0.3 {
			return AIRatioCategoryLow
		} else if aiRatio < 0.7 {
			return AIRatioCategoryMedium
		}
		return AIRatioCategoryHigh
	}

	bands := g.seed.Correlations.AIRatioBands

	// Check bands in order
	if aiRatio >= bands.Low.Min && aiRatio < bands.Low.Max {
		return AIRatioCategoryLow
	}
	if aiRatio >= bands.Medium.Min && aiRatio < bands.Medium.Max {
		return AIRatioCategoryMedium
	}
	return AIRatioCategoryHigh
}

// CalculateRevertProbability calculates the probability of a PR being reverted.
func (g *QualityGenerator) CalculateRevertProbability(aiRatio float64, reviewIterations int) float64 {
	baseProb := 0.05 // Default 5%

	if g.seed != nil && g.seed.PRLifecycle.QualityOutcomes.RevertProbability.Base > 0 {
		baseProb = g.seed.PRLifecycle.QualityOutcomes.RevertProbability.Base
	}

	// Apply AI ratio modifier
	category := g.CategorizeAIRatio(aiRatio)
	modifier := 1.0

	if g.seed != nil {
		modifiers := g.seed.PRLifecycle.QualityOutcomes.RevertProbability.Modifiers.ByAIRatio
		if modifiers != nil {
			if m, ok := modifiers[string(category)]; ok {
				modifier = m
			}
		}
	}

	prob := baseProb * modifier

	// More review iterations reduce revert probability
	if reviewIterations > 1 {
		// Each additional iteration reduces probability by 10%
		iterationReduction := 1.0 - float64(reviewIterations-1)*0.1
		if iterationReduction < 0.5 {
			iterationReduction = 0.5 // Cap at 50% reduction
		}
		prob *= iterationReduction
	}

	return prob
}

// CalculateHotfixProbability calculates the probability of needing a hotfix after merge.
func (g *QualityGenerator) CalculateHotfixProbability(pr *models.PullRequest) float64 {
	baseProb := 0.1 // Default 10%

	if g.seed != nil && g.seed.PRLifecycle.QualityOutcomes.HotfixProbability.Base > 0 {
		baseProb = g.seed.PRLifecycle.QualityOutcomes.HotfixProbability.Base
	}

	// Bug fixes are more likely to need hotfixes (ironic but realistic)
	if pr.IsBugFix {
		baseProb *= 1.3
	}

	// High AI ratio increases hotfix probability
	if pr.AIRatio > 0.7 {
		baseProb *= 1.2
	}

	// Large PRs are more likely to need hotfixes
	if pr.Additions > 500 {
		baseProb *= 1.5
	}

	return baseProb
}

// IsBugFix determines if a PR is a bug fix based on title and labels.
func (g *QualityGenerator) IsBugFix(pr *models.PullRequest) bool {
	// Check labels
	for _, label := range pr.Labels {
		labelLower := strings.ToLower(label)
		if labelLower == "bug" || labelLower == "bugfix" || labelLower == "hotfix" {
			return true
		}
	}

	// Check title prefix
	titleLower := strings.ToLower(pr.Title)
	bugPrefixes := []string{"fix:", "fix(", "bugfix:", "bugfix(", "hotfix:", "hotfix("}

	for _, prefix := range bugPrefixes {
		if strings.HasPrefix(titleLower, prefix) {
			return true
		}
	}

	return false
}

// ApplyQualityOutcomes applies quality signals (revert, hotfix) to merged PRs in a repo.
func (g *QualityGenerator) ApplyQualityOutcomes(repoName string) error {
	prs := g.store.GetPRsByRepo(repoName)

	for _, pr := range prs {
		// Only apply to merged PRs
		if pr.State != models.PRStateMerged {
			continue
		}

		// Skip if already processed
		if pr.WasReverted {
			continue
		}

		// Determine if this is a bug fix
		pr.IsBugFix = g.IsBugFix(&pr)

		// Calculate and apply revert probability
		revertProb := g.CalculateRevertProbability(pr.AIRatio, pr.CommitCount)
		if g.rng.Float64() < revertProb {
			pr.WasReverted = true
		}

		// Update the PR
		if err := g.store.UpdatePR(pr); err != nil {
			return err
		}
	}

	return nil
}
