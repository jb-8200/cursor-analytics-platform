package services

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// PRStore defines the interface for accessing PRs needed by revert analysis.
type PRStore interface {
	GetPRsByRepo(repoName string) []models.PullRequest
	GetPRsByRepoAndState(repoName string, state models.PRState) []models.PullRequest
	GetDeveloper(userID string) (*seed.Developer, error)
}

// RevertService calculates revert metrics for PRs in a repository.
type RevertService struct {
	store PRStore
	rng   *rand.Rand
}

// NewRevertService creates a new revert analysis service.
func NewRevertService(store PRStore) *RevertService {
	return &RevertService{
		store: store,
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewRevertServiceWithSeed creates a revert service with a specific random seed for reproducibility.
func NewRevertServiceWithSeed(store PRStore, randSeed int64) *RevertService {
	return &RevertService{
		store: store,
		rng:   rand.New(rand.NewSource(randSeed)),
	}
}

// CalculateRevertRisk calculates the probability that a PR will be reverted.
// Uses sigmoid function: higher AI ratio + lower seniority → higher risk
func CalculateRevertRisk(pr models.PullRequest, dev *seed.Developer) float64 {
	// Baseline risk parameters
	const (
		aiWeight        = 1.8  // Weight for AI ratio impact
		seniorityWeight = 0.8  // Seniority modifier weight
		activityWeight  = 0.3  // Activity level affects risk
		baselineRisk    = -3.5 // Shifts sigmoid center (lower = lower overall risk)
	)

	// Seniority penalty (0-1, higher is worse - junior has higher penalty)
	seniorityPenalty := 0.0
	if dev != nil {
		switch dev.Seniority {
		case "junior":
			seniorityPenalty = 1.0 // Highest penalty
		case "mid":
			seniorityPenalty = 0.5
		case "senior":
			seniorityPenalty = -0.5 // Negative penalty (actually reduces risk)
		}
	}

	// Activity level modifier (higher activity can mean more churn)
	activityModifier := 0.5 // default
	if dev != nil {
		switch dev.ActivityLevel {
		case "high":
			activityModifier = 0.8
		case "medium":
			activityModifier = 0.5
		case "low":
			activityModifier = 0.2
		}
	}

	// Calculate raw score
	rawScore := baselineRisk +
		aiWeight*pr.AIRatio +
		seniorityWeight*seniorityPenalty +
		activityWeight*activityModifier

	// Apply sigmoid: 1 / (1 + e^(-rawScore))
	risk := 1.0 / (1.0 + math.Exp(-rawScore))

	// Cap at reasonable maximum (15%)
	if risk > 0.15 {
		risk = 0.15
	}

	return risk
}

// ShouldRevert determines if a PR should be reverted based on risk probability.
func (s *RevertService) ShouldRevert(pr models.PullRequest, dev *seed.Developer) bool {
	risk := CalculateRevertRisk(pr, dev)
	return s.rng.Float64() < risk
}

// IsRevertMessage detects if a commit message indicates a revert.
// Pattern matches: "revert", "Revert", "rollback", "Rollback"
func IsRevertMessage(message string) bool {
	patterns := []string{
		`(?i)revert`,
		`(?i)rollback`,
		`(?i)backing out`,
		`(?i)reverting`,
	}

	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, message)
		if matched {
			return true
		}
	}

	return false
}

// ExtractPRNumber extracts a PR number from a commit message.
// Looks for patterns like "#123", "PR 123", "pull request #123"
func ExtractPRNumber(message string) (int, bool) {
	patterns := []string{
		`#(\d+)`,
		`PR\s+(\d+)`,
		`pull request\s+#(\d+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(message)
		if len(matches) > 1 {
			var prNum int
			fmt.Sscanf(matches[1], "%d", &prNum)
			if prNum > 0 {
				return prNum, true
			}
		}
	}

	return 0, false
}

// GetReverts analyzes a repository for revert events within a time window.
// windowDays specifies how many days after merge to check for reverts.
// since and until define the time range to analyze.
func (s *RevertService) GetReverts(
	repoName string,
	windowDays int,
	since, until time.Time,
) (*models.RevertAnalysis, error) {
	// Get all merged PRs
	allPRs := s.store.GetPRsByRepoAndState(repoName, models.PRStateMerged)

	// Filter PRs merged within the analysis period
	var mergedPRs []models.PullRequest
	for _, pr := range allPRs {
		if pr.MergedAt != nil &&
			!pr.MergedAt.Before(since) &&
			pr.MergedAt.Before(until) {
			mergedPRs = append(mergedPRs, pr)
		}
	}

	if len(mergedPRs) == 0 {
		return &models.RevertAnalysis{
			WindowDays:       windowDays,
			TotalPRsMerged:   0,
			TotalPRsReverted: 0,
			RevertRate:       0.0,
			RevertedPRs:      []models.RevertedPR{},
		}, nil
	}

	// Simulate reverts for merged PRs based on risk calculation
	var revertedPRs []models.RevertedPR

	for _, pr := range mergedPRs {
		// Get developer info
		dev, _ := s.store.GetDeveloper(pr.AuthorID)

		// Check if PR should be reverted
		if s.ShouldRevert(pr, dev) {
			// Calculate revert timing (within window)
			daysToRevert := s.rng.Float64() * float64(windowDays)
			revertTime := pr.MergedAt.Add(time.Duration(daysToRevert*24) * time.Hour)

			// Only include if revert happens before analysis end
			if revertTime.Before(until) {
				revertedPRs = append(revertedPRs, models.RevertedPR{
					PRNumber:     pr.Number,
					MergedAt:     pr.MergedAt.Format(time.RFC3339),
					RevertedAt:   revertTime.Format(time.RFC3339),
					DaysToRevert: daysToRevert,
				})
			}
		}
	}

	// Calculate revert rate
	revertRate := 0.0
	if len(mergedPRs) > 0 {
		revertRate = float64(len(revertedPRs)) / float64(len(mergedPRs))
	}

	return &models.RevertAnalysis{
		WindowDays:       windowDays,
		TotalPRsMerged:   len(mergedPRs),
		TotalPRsReverted: len(revertedPRs),
		RevertRate:       revertRate,
		RevertedPRs:      revertedPRs,
	}, nil
}

// GenerateRevertMessage generates a realistic revert commit message.
func GenerateRevertMessage(prNumber int, originalTitle string) string {
	templates := []string{
		"Revert \"{}\" (#{})",
		"Rollback PR #{}: {}",
		"Reverting changes from #{} - {}",
		"Backing out #{}: {}",
	}

	template := templates[rand.Intn(len(templates))]
	template = strings.ReplaceAll(template, "{}", "%s")
	return fmt.Sprintf(template, originalTitle, prNumber)
}
