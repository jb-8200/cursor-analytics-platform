package storage

import (
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// Store defines the interface for storing and querying simulation data.
// Implementations must be thread-safe.
type Store interface {
	// Developer operations
	LoadDevelopers(developers []seed.Developer) error
	GetDeveloper(userID string) (*seed.Developer, error)
	GetDeveloperByEmail(email string) (*seed.Developer, error)
	ListDevelopers() []seed.Developer

	// Commit operations
	AddCommit(commit models.Commit) error
	GetCommitByHash(hash string) (*models.Commit, error)
	GetCommitsByTimeRange(from, to time.Time) []models.Commit
	GetCommitsByUser(userID string, from, to time.Time) []models.Commit
	GetCommitsByRepo(repoName string, from, to time.Time) []models.Commit

	// PR operations
	AddPR(pr models.PullRequest) error
	UpdatePR(pr models.PullRequest) error
	GetPR(repoName string, number int) (*models.PullRequest, error)
	GetPRsByRepo(repoName string) []models.PullRequest
	GetPRsByRepoAndState(repoName string, state models.PRState) []models.PullRequest
	GetPRsByAuthor(authorID string) []models.PullRequest
	GetNextPRNumber(repoName string) int
	ListRepositories() []string

	// Review comment operations
	AddReviewComment(comment models.ReviewComment) error
	GetReviewComments(repoName string, prNumber int) []models.ReviewComment

	// Model usage operations
	AddModelUsage(usage models.ModelUsageEvent) error
	GetModelUsageByTimeRange(from, to time.Time) []models.ModelUsageEvent
}
