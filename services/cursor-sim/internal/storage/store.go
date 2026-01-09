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

	// PR operations (GitHub Simulation - P2-F01)
	StorePR(pr models.PullRequest) error
	GetPRByID(id int) (*models.PullRequest, error)
	GetPRsByStatus(status models.PRState) ([]models.PullRequest, error)
	GetPRsByAuthorEmail(authorEmail string) ([]models.PullRequest, error)
	GetPRsByRepoWithPagination(repoName string, state string, page, pageSize int) ([]models.PullRequest, int, error)

	// Review comment operations
	AddReviewComment(comment models.ReviewComment) error
	GetReviewComments(repoName string, prNumber int) []models.ReviewComment

	// Review operations (GitHub Simulation - P2-F01)
	StoreReview(review models.Review) error
	GetReviewsByPRID(prID int64) ([]models.Review, error)
	GetReviewsByReviewer(reviewerEmail string) ([]models.Review, error)
	GetReviewsByRepoPR(repoName string, prNumber int) ([]models.Review, error)

	// Issue operations (GitHub Simulation - P2-F01)
	StoreIssue(issue models.Issue) error
	GetIssueByNumber(repoName string, number int) (*models.Issue, error)
	GetIssuesByState(repoName string, state models.IssueState) ([]models.Issue, error)
	GetIssuesByRepo(repoName string) ([]models.Issue, error)

	// Model usage operations
	AddModelUsage(usage models.ModelUsageEvent) error
	GetModelUsageByTimeRange(from, to time.Time) []models.ModelUsageEvent

	// Client version operations
	AddClientVersion(event models.ClientVersionEvent) error
	GetClientVersionsByTimeRange(from, to time.Time) []models.ClientVersionEvent

	// File extension operations
	AddFileExtension(event models.FileExtensionEvent) error
	GetFileExtensionsByTimeRange(from, to time.Time) []models.FileExtensionEvent

	// MCP tool operations
	AddMCPTool(event models.MCPToolEvent) error
	GetMCPToolsByTimeRange(from, to time.Time) []models.MCPToolEvent

	// Command operations
	AddCommand(event models.CommandEvent) error
	GetCommandsByTimeRange(from, to time.Time) []models.CommandEvent

	// Plan operations
	AddPlan(event models.PlanEvent) error
	GetPlansByTimeRange(from, to time.Time) []models.PlanEvent

	// Ask mode operations
	AddAskMode(event models.AskModeEvent) error
	GetAskModeByTimeRange(from, to time.Time) []models.AskModeEvent
}
