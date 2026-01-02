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
}
