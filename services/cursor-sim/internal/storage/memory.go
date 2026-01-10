package storage

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// MemoryStore is a thread-safe in-memory implementation of Store.
// It uses multiple indexes for efficient queries:
// - commitsByHash: O(1) lookup by commit hash
// - commitsByUser: O(1) lookup by user, then O(log n) time filter
// - commitsByRepo: O(1) lookup by repo, then O(log n) time filter
// - commits: time-sorted slice for O(log n) range queries
type MemoryStore struct {
	mu sync.RWMutex

	// Developer data
	developers      map[string]*seed.Developer // by user_id
	developerEmails map[string]string          // email -> user_id

	// Commit data with multiple indexes
	commits       []*models.Commit            // time-sorted for range queries
	commitsByHash map[string]*models.Commit   // hash lookup
	commitsByUser map[string][]*models.Commit // user index
	commitsByRepo map[string][]*models.Commit // repo index
	needsSort     bool                        // flag to track if commits need sorting

	// PR data with indexes
	prsByRepo   map[string]map[int]*models.PullRequest // repo -> number -> PR
	prsByAuthor map[string][]*models.PullRequest       // author index
	prsByID     map[int]*models.PullRequest            // ID index for GitHub simulation
	prsByEmail  map[string][]*models.PullRequest       // email index for GitHub simulation

	// Review data (GitHub Simulation - P2-F01)
	reviewsByID       map[int]*models.Review      // review ID index
	reviewsByPRID     map[int][]*models.Review    // PR ID index
	reviewsByReviewer map[string][]*models.Review // reviewer email index

	// Issue data (GitHub Simulation - P2-F01)
	issuesByRepo  map[string]map[int]*models.Issue                 // repo -> number -> Issue
	issuesByState map[string]map[models.IssueState][]*models.Issue // repo -> state -> issues

	// ReviewComment data
	reviewComments map[string]map[int][]*models.ReviewComment // repo -> pr_number -> comments

	// Model usage data
	modelUsage []*models.ModelUsageEvent // time-sorted for range queries

	// Client version data
	clientVersions []*models.ClientVersionEvent // time-sorted for range queries

	// File extension data
	fileExtensions []*models.FileExtensionEvent // time-sorted for range queries

	// Feature usage data
	mcpTools []*models.MCPToolEvent // MCP tool usage events
	commands []*models.CommandEvent // Command usage events
	plans    []*models.PlanEvent    // Plan usage events
	askModes []*models.AskModeEvent // Ask mode usage events

	// Counters for auto-generating IDs
	nextPRID int // Auto-incrementing PR ID counter
}

// NewMemoryStore creates a new thread-safe in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		developers:        make(map[string]*seed.Developer),
		developerEmails:   make(map[string]string),
		commits:           make([]*models.Commit, 0, 1000),
		commitsByHash:     make(map[string]*models.Commit),
		commitsByUser:     make(map[string][]*models.Commit),
		commitsByRepo:     make(map[string][]*models.Commit),
		prsByRepo:         make(map[string]map[int]*models.PullRequest),
		prsByAuthor:       make(map[string][]*models.PullRequest),
		prsByID:           make(map[int]*models.PullRequest),
		prsByEmail:        make(map[string][]*models.PullRequest),
		reviewsByID:       make(map[int]*models.Review),
		reviewsByPRID:     make(map[int][]*models.Review),
		reviewsByReviewer: make(map[string][]*models.Review),
		issuesByRepo:      make(map[string]map[int]*models.Issue),
		issuesByState:     make(map[string]map[models.IssueState][]*models.Issue),
		reviewComments:    make(map[string]map[int][]*models.ReviewComment),
		modelUsage:        make([]*models.ModelUsageEvent, 0, 1000),
		clientVersions:    make([]*models.ClientVersionEvent, 0, 1000),
		fileExtensions:    make([]*models.FileExtensionEvent, 0, 5000), // Higher capacity for file events
		mcpTools:          make([]*models.MCPToolEvent, 0, 2000),
		commands:          make([]*models.CommandEvent, 0, 2000),
		plans:             make([]*models.PlanEvent, 0, 1500),
		askModes:          make([]*models.AskModeEvent, 0, 1500),
		nextPRID:          1, // Start PR IDs at 1
	}
}

// LoadDevelopers loads developer data into the store.
func (m *MemoryStore) LoadDevelopers(developers []seed.Developer) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range developers {
		dev := &developers[i]
		m.developers[dev.UserID] = dev
		m.developerEmails[dev.Email] = dev.UserID
	}

	return nil
}

// GetDeveloper retrieves a developer by user ID.
func (m *MemoryStore) GetDeveloper(userID string) (*seed.Developer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	dev, ok := m.developers[userID]
	if !ok {
		return nil, fmt.Errorf("developer not found: %s", userID)
	}

	return dev, nil
}

// GetDeveloperByEmail retrieves a developer by email address.
func (m *MemoryStore) GetDeveloperByEmail(email string) (*seed.Developer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userID, ok := m.developerEmails[email]
	if !ok {
		return nil, fmt.Errorf("developer not found: %s", email)
	}

	dev := m.developers[userID]
	return dev, nil
}

// ListDevelopers returns all developers sorted by UserID.
func (m *MemoryStore) ListDevelopers() []seed.Developer {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]seed.Developer, 0, len(m.developers))
	for _, dev := range m.developers {
		result = append(result, *dev)
	}

	// Sort by UserID for deterministic ordering
	sort.Slice(result, func(i, j int) bool {
		return result[i].UserID < result[j].UserID
	})

	return result
}

// AddCommit adds a commit to the store with all indexes updated.
func (m *MemoryStore) AddCommit(commit models.Commit) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Store pointer to commit
	commitPtr := &commit

	// Add to hash index
	m.commitsByHash[commit.CommitHash] = commitPtr

	// Add to time-sorted slice
	m.commits = append(m.commits, commitPtr)
	m.needsSort = true // Mark that sorting is needed

	// Add to user index
	m.commitsByUser[commit.UserID] = append(m.commitsByUser[commit.UserID], commitPtr)

	// Add to repo index
	m.commitsByRepo[commit.RepoName] = append(m.commitsByRepo[commit.RepoName], commitPtr)

	return nil
}

// GetCommitByHash retrieves a commit by its hash.
func (m *MemoryStore) GetCommitByHash(hash string) (*models.Commit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	commit, ok := m.commitsByHash[hash]
	if !ok {
		return nil, fmt.Errorf("commit not found: %s", hash)
	}

	return commit, nil
}

// GetCommitsByTimeRange returns all commits within the time range.
// Uses binary search on time-sorted slice for O(log n) performance.
func (m *MemoryStore) GetCommitsByTimeRange(from, to time.Time) []models.Commit {
	m.mu.Lock()
	if m.needsSort {
		m.sortCommits()
	}
	m.mu.Unlock()

	m.mu.RLock()
	defer m.mu.RUnlock()

	// Binary search for start index
	startIdx := sort.Search(len(m.commits), func(i int) bool {
		return !m.commits[i].CommitTs.Before(from)
	})

	// Collect commits in range
	result := make([]models.Commit, 0)
	for i := startIdx; i < len(m.commits); i++ {
		if m.commits[i].CommitTs.After(to) {
			break
		}
		result = append(result, *m.commits[i])
	}

	return result
}

// GetCommitsByUser returns commits by a specific user within the time range.
func (m *MemoryStore) GetCommitsByUser(userID string, from, to time.Time) []models.Commit {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userCommits, ok := m.commitsByUser[userID]
	if !ok {
		return []models.Commit{}
	}

	result := make([]models.Commit, 0)
	for _, commit := range userCommits {
		if !commit.CommitTs.Before(from) && !commit.CommitTs.After(to) {
			result = append(result, *commit)
		}
	}

	return result
}

// GetCommitsByRepo returns commits for a specific repository within the time range.
func (m *MemoryStore) GetCommitsByRepo(repoName string, from, to time.Time) []models.Commit {
	m.mu.RLock()
	defer m.mu.RUnlock()

	repoCommits, ok := m.commitsByRepo[repoName]
	if !ok {
		return []models.Commit{}
	}

	result := make([]models.Commit, 0)
	for _, commit := range repoCommits {
		if !commit.CommitTs.Before(from) && !commit.CommitTs.After(to) {
			result = append(result, *commit)
		}
	}

	return result
}

// sortCommits sorts the commits slice by timestamp.
// Must be called with write lock held.
func (m *MemoryStore) sortCommits() {
	sort.Slice(m.commits, func(i, j int) bool {
		return m.commits[i].CommitTs.Before(m.commits[j].CommitTs)
	})
	m.needsSort = false
}

// PR Storage Methods

// AddPR adds a pull request to the store.
func (m *MemoryStore) AddPR(pr models.PullRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Auto-generate ID if not set
	if pr.ID == 0 {
		pr.ID = m.nextPRID
		m.nextPRID++
	}

	prPtr := &pr

	// Ensure repo map exists
	if m.prsByRepo[pr.RepoName] == nil {
		m.prsByRepo[pr.RepoName] = make(map[int]*models.PullRequest)
	}

	// Add to repo index
	m.prsByRepo[pr.RepoName][pr.Number] = prPtr

	// Add to author index
	m.prsByAuthor[pr.AuthorID] = append(m.prsByAuthor[pr.AuthorID], prPtr)

	// Add to ID index for GitHub simulation
	m.prsByID[pr.ID] = prPtr

	// Add to email index for GitHub simulation
	m.prsByEmail[pr.AuthorEmail] = append(m.prsByEmail[pr.AuthorEmail], prPtr)

	return nil
}

// UpdatePR updates an existing pull request.
func (m *MemoryStore) UpdatePR(pr models.PullRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.prsByRepo[pr.RepoName] == nil {
		return fmt.Errorf("PR not found: %s#%d", pr.RepoName, pr.Number)
	}

	existing, ok := m.prsByRepo[pr.RepoName][pr.Number]
	if !ok {
		return fmt.Errorf("PR not found: %s#%d", pr.RepoName, pr.Number)
	}

	// Update in place
	*existing = pr
	return nil
}

// GetPR retrieves a pull request by repo and number.
func (m *MemoryStore) GetPR(repoName string, number int) (*models.PullRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	repoPRs, ok := m.prsByRepo[repoName]
	if !ok {
		return nil, fmt.Errorf("PR not found: %s#%d", repoName, number)
	}

	pr, ok := repoPRs[number]
	if !ok {
		return nil, fmt.Errorf("PR not found: %s#%d", repoName, number)
	}

	return pr, nil
}

// GetPRsByRepo returns all PRs for a repository.
func (m *MemoryStore) GetPRsByRepo(repoName string) []models.PullRequest {
	m.mu.RLock()
	defer m.mu.RUnlock()

	repoPRs, ok := m.prsByRepo[repoName]
	if !ok {
		return []models.PullRequest{}
	}

	result := make([]models.PullRequest, 0, len(repoPRs))
	for _, pr := range repoPRs {
		result = append(result, *pr)
	}

	return result
}

// GetPRsByRepoAndState returns PRs for a repository filtered by state.
func (m *MemoryStore) GetPRsByRepoAndState(repoName string, state models.PRState) []models.PullRequest {
	m.mu.RLock()
	defer m.mu.RUnlock()

	repoPRs, ok := m.prsByRepo[repoName]
	if !ok {
		return []models.PullRequest{}
	}

	result := make([]models.PullRequest, 0)
	for _, pr := range repoPRs {
		if pr.State == state {
			result = append(result, *pr)
		}
	}

	return result
}

// GetPRsByAuthor returns all PRs by a specific author.
func (m *MemoryStore) GetPRsByAuthor(authorID string) []models.PullRequest {
	m.mu.RLock()
	defer m.mu.RUnlock()

	authorPRs, ok := m.prsByAuthor[authorID]
	if !ok {
		return []models.PullRequest{}
	}

	result := make([]models.PullRequest, 0, len(authorPRs))
	for _, pr := range authorPRs {
		result = append(result, *pr)
	}

	return result
}

// GetNextPRNumber returns the next available PR number for a repository.
func (m *MemoryStore) GetNextPRNumber(repoName string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	repoPRs, ok := m.prsByRepo[repoName]
	if !ok || len(repoPRs) == 0 {
		return 1
	}

	maxNum := 0
	for num := range repoPRs {
		if num > maxNum {
			maxNum = num
		}
	}

	return maxNum + 1
}

// ListRepositories returns all repository names that have PRs or issues.
func (m *MemoryStore) ListRepositories() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a map to deduplicate repos from both PRs and issues
	repoSet := make(map[string]bool)
	for repo := range m.prsByRepo {
		repoSet[repo] = true
	}
	for repo := range m.issuesByRepo {
		repoSet[repo] = true
	}

	// Convert to slice
	repos := make([]string, 0, len(repoSet))
	for repo := range repoSet {
		repos = append(repos, repo)
	}

	return repos
}

// ReviewComment Storage Methods

// AddReviewComment adds a review comment to the store.
func (m *MemoryStore) AddReviewComment(comment models.ReviewComment) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	commentPtr := &comment

	// Ensure repo map exists
	if m.reviewComments[comment.RepoName] == nil {
		m.reviewComments[comment.RepoName] = make(map[int][]*models.ReviewComment)
	}

	// Add to PR's comments
	m.reviewComments[comment.RepoName][comment.PRNumber] = append(
		m.reviewComments[comment.RepoName][comment.PRNumber],
		commentPtr,
	)

	return nil
}

// GetReviewComments returns all comments for a PR.
func (m *MemoryStore) GetReviewComments(repoName string, prNumber int) []models.ReviewComment {
	m.mu.RLock()
	defer m.mu.RUnlock()

	repoComments, ok := m.reviewComments[repoName]
	if !ok {
		return []models.ReviewComment{}
	}

	prComments, ok := repoComments[prNumber]
	if !ok {
		return []models.ReviewComment{}
	}

	result := make([]models.ReviewComment, 0, len(prComments))
	for _, c := range prComments {
		result = append(result, *c)
	}

	return result
}

// AddModelUsage stores a model usage event.
func (m *MemoryStore) AddModelUsage(usage models.ModelUsageEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.modelUsage = append(m.modelUsage, &usage)
	return nil
}

// GetModelUsageByTimeRange retrieves all model usage events within a time range.
// Returns events sorted by timestamp.
func (m *MemoryStore) GetModelUsageByTimeRange(from, to time.Time) []models.ModelUsageEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]models.ModelUsageEvent, 0)
	for _, usage := range m.modelUsage {
		if !usage.Timestamp.Before(from) && usage.Timestamp.Before(to) {
			result = append(result, *usage)
		}
	}

	// Sort by timestamp
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})

	return result
}

// AddClientVersion stores a client version event.
func (m *MemoryStore) AddClientVersion(event models.ClientVersionEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clientVersions = append(m.clientVersions, &event)
	return nil
}

// GetClientVersionsByTimeRange retrieves all client version events within a time range.
// Returns events sorted by timestamp.
func (m *MemoryStore) GetClientVersionsByTimeRange(from, to time.Time) []models.ClientVersionEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]models.ClientVersionEvent, 0)
	for _, event := range m.clientVersions {
		if !event.Timestamp.Before(from) && event.Timestamp.Before(to) {
			result = append(result, *event)
		}
	}

	// Sort by timestamp
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})

	return result
}

// AddFileExtension stores a file extension event.
func (m *MemoryStore) AddFileExtension(event models.FileExtensionEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.fileExtensions = append(m.fileExtensions, &event)
	return nil
}

// GetFileExtensionsByTimeRange retrieves all file extension events within a time range.
// Returns events sorted by timestamp.
func (m *MemoryStore) GetFileExtensionsByTimeRange(from, to time.Time) []models.FileExtensionEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]models.FileExtensionEvent, 0)
	for _, event := range m.fileExtensions {
		if !event.Timestamp.Before(from) && event.Timestamp.Before(to) {
			result = append(result, *event)
		}
	}

	// Sort by timestamp
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})

	return result
}

// AddMCPTool stores an MCP tool event.
func (m *MemoryStore) AddMCPTool(event models.MCPToolEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mcpTools = append(m.mcpTools, &event)
	return nil
}

// GetMCPToolsByTimeRange retrieves all MCP tool events within a time range.
func (m *MemoryStore) GetMCPToolsByTimeRange(from, to time.Time) []models.MCPToolEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]models.MCPToolEvent, 0)
	for _, event := range m.mcpTools {
		if !event.Timestamp.Before(from) && event.Timestamp.Before(to) {
			result = append(result, *event)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result
}

// AddCommand stores a command event.
func (m *MemoryStore) AddCommand(event models.CommandEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.commands = append(m.commands, &event)
	return nil
}

// GetCommandsByTimeRange retrieves all command events within a time range.
func (m *MemoryStore) GetCommandsByTimeRange(from, to time.Time) []models.CommandEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]models.CommandEvent, 0)
	for _, event := range m.commands {
		if !event.Timestamp.Before(from) && event.Timestamp.Before(to) {
			result = append(result, *event)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result
}

// AddPlan stores a plan event.
func (m *MemoryStore) AddPlan(event models.PlanEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.plans = append(m.plans, &event)
	return nil
}

// GetPlansByTimeRange retrieves all plan events within a time range.
func (m *MemoryStore) GetPlansByTimeRange(from, to time.Time) []models.PlanEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]models.PlanEvent, 0)
	for _, event := range m.plans {
		if !event.Timestamp.Before(from) && event.Timestamp.Before(to) {
			result = append(result, *event)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result
}

// AddAskMode stores an ask mode event.
func (m *MemoryStore) AddAskMode(event models.AskModeEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.askModes = append(m.askModes, &event)
	return nil
}

// GetAskModeByTimeRange retrieves all ask mode events within a time range.
func (m *MemoryStore) GetAskModeByTimeRange(from, to time.Time) []models.AskModeEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]models.AskModeEvent, 0)
	for _, event := range m.askModes {
		if !event.Timestamp.Before(from) && event.Timestamp.Before(to) {
			result = append(result, *event)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result
}

// GitHub Simulation Storage Methods (P2-F01)

// StorePR stores a pull request by ID for GitHub simulation.
func (m *MemoryStore) StorePR(pr models.PullRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	prPtr := &pr

	// Store by ID
	m.prsByID[pr.ID] = prPtr

	// Store by email for author lookup
	m.prsByEmail[pr.AuthorEmail] = append(m.prsByEmail[pr.AuthorEmail], prPtr)

	// Also maintain repo/number index (compatibility with existing AddPR)
	if m.prsByRepo[pr.RepoName] == nil {
		m.prsByRepo[pr.RepoName] = make(map[int]*models.PullRequest)
	}
	m.prsByRepo[pr.RepoName][pr.Number] = prPtr

	// Maintain author index
	m.prsByAuthor[pr.AuthorID] = append(m.prsByAuthor[pr.AuthorID], prPtr)

	return nil
}

// GetPRByID retrieves a pull request by its internal ID.
func (m *MemoryStore) GetPRByID(id int) (*models.PullRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pr, ok := m.prsByID[id]
	if !ok {
		return nil, fmt.Errorf("PR not found: %d", id)
	}

	return pr, nil
}

// GetPRsByStatus returns all PRs with the given state across all repos.
func (m *MemoryStore) GetPRsByStatus(status models.PRState) ([]models.PullRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]models.PullRequest, 0)
	for _, pr := range m.prsByID {
		if pr.State == status {
			result = append(result, *pr)
		}
	}

	return result, nil
}

// GetPRsByAuthorEmail returns all PRs by a specific author email.
func (m *MemoryStore) GetPRsByAuthorEmail(authorEmail string) ([]models.PullRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	prs, ok := m.prsByEmail[authorEmail]
	if !ok {
		return []models.PullRequest{}, nil
	}

	result := make([]models.PullRequest, 0, len(prs))
	for _, pr := range prs {
		result = append(result, *pr)
	}

	return result, nil
}

// GetPRsByRepoWithPagination returns PRs for a repo with pagination support.
// state can be empty string for all states, or "open", "closed", "merged".
func (m *MemoryStore) GetPRsByRepoWithPagination(repoName string, state string, page, pageSize int) ([]models.PullRequest, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	repoPRs, ok := m.prsByRepo[repoName]
	if !ok {
		return []models.PullRequest{}, 0, nil
	}

	// Collect all PRs (optionally filtered by state)
	allPRs := make([]models.PullRequest, 0)
	for _, pr := range repoPRs {
		if state == "" || string(pr.State) == state {
			allPRs = append(allPRs, *pr)
		}
	}

	total := len(allPRs)

	// Calculate pagination
	startIdx := (page - 1) * pageSize
	endIdx := startIdx + pageSize

	if startIdx >= total {
		return []models.PullRequest{}, total, nil
	}

	if endIdx > total {
		endIdx = total
	}

	result := allPRs[startIdx:endIdx]
	return result, total, nil
}

// StoreReview stores a review.
func (m *MemoryStore) StoreReview(review models.Review) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reviewPtr := &review

	// Store by review ID
	m.reviewsByID[review.ID] = reviewPtr

	// Store by PR ID
	m.reviewsByPRID[review.PRID] = append(m.reviewsByPRID[review.PRID], reviewPtr)

	// Store by reviewer email
	m.reviewsByReviewer[review.Reviewer] = append(m.reviewsByReviewer[review.Reviewer], reviewPtr)

	return nil
}

// GetReviewsByPRID retrieves all reviews for a specific PR by PR ID.
func (m *MemoryStore) GetReviewsByPRID(prID int64) ([]models.Review, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	reviews, ok := m.reviewsByPRID[int(prID)]
	if !ok {
		return []models.Review{}, nil
	}

	result := make([]models.Review, 0, len(reviews))
	for _, review := range reviews {
		result = append(result, *review)
	}

	return result, nil
}

// GetReviewsByReviewer retrieves all reviews by a specific reviewer email.
func (m *MemoryStore) GetReviewsByReviewer(reviewerEmail string) ([]models.Review, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	reviews, ok := m.reviewsByReviewer[reviewerEmail]
	if !ok {
		return []models.Review{}, nil
	}

	result := make([]models.Review, 0, len(reviews))
	for _, review := range reviews {
		result = append(result, *review)
	}

	return result, nil
}

// GetReviewsByRepoPR retrieves all reviews for a specific PR by repo name and PR number.
func (m *MemoryStore) GetReviewsByRepoPR(repoName string, prNumber int) ([]models.Review, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// First, find the PR by repo and number to get its ID
	repoPRs, ok := m.prsByRepo[repoName]
	if !ok {
		return []models.Review{}, nil
	}

	pr, ok := repoPRs[prNumber]
	if !ok {
		return []models.Review{}, nil
	}

	// Get reviews by PR ID
	reviews, ok := m.reviewsByPRID[pr.ID]
	if !ok {
		return []models.Review{}, nil
	}

	result := make([]models.Review, 0, len(reviews))
	for _, review := range reviews {
		result = append(result, *review)
	}

	return result, nil
}

// StoreIssue stores an issue.
func (m *MemoryStore) StoreIssue(issue models.Issue) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	issuePtr := &issue

	// Ensure repo map exists
	if m.issuesByRepo[issue.RepoName] == nil {
		m.issuesByRepo[issue.RepoName] = make(map[int]*models.Issue)
	}

	// Store by repo and number
	m.issuesByRepo[issue.RepoName][issue.Number] = issuePtr

	// Ensure state map exists
	if m.issuesByState[issue.RepoName] == nil {
		m.issuesByState[issue.RepoName] = make(map[models.IssueState][]*models.Issue)
	}

	// Store by state
	m.issuesByState[issue.RepoName][issue.State] = append(
		m.issuesByState[issue.RepoName][issue.State],
		issuePtr,
	)

	return nil
}

// GetIssueByNumber retrieves an issue by repo name and number.
func (m *MemoryStore) GetIssueByNumber(repoName string, number int) (*models.Issue, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	repoIssues, ok := m.issuesByRepo[repoName]
	if !ok {
		return nil, fmt.Errorf("issue not found: %s#%d", repoName, number)
	}

	issue, ok := repoIssues[number]
	if !ok {
		return nil, fmt.Errorf("issue not found: %s#%d", repoName, number)
	}

	return issue, nil
}

// GetIssuesByState retrieves all issues for a repo filtered by state.
func (m *MemoryStore) GetIssuesByState(repoName string, state models.IssueState) ([]models.Issue, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	repoStates, ok := m.issuesByState[repoName]
	if !ok {
		return []models.Issue{}, nil
	}

	issues, ok := repoStates[state]
	if !ok {
		return []models.Issue{}, nil
	}

	result := make([]models.Issue, 0, len(issues))
	for _, issue := range issues {
		result = append(result, *issue)
	}

	return result, nil
}

// GetIssuesByRepo retrieves all issues for a repo (all states).
func (m *MemoryStore) GetIssuesByRepo(repoName string) ([]models.Issue, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	repoIssues, ok := m.issuesByRepo[repoName]
	if !ok {
		return []models.Issue{}, nil
	}

	result := make([]models.Issue, 0, len(repoIssues))
	for _, issue := range repoIssues {
		result = append(result, *issue)
	}

	return result, nil
}

// Admin API Methods (P1-F02)

// ClearAllData removes all data from the store.
// Used by regenerate endpoint in override mode.
// This method is thread-safe and resets all data structures to their initial state.
func (m *MemoryStore) ClearAllData() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Reset developer data
	m.developers = make(map[string]*seed.Developer)
	m.developerEmails = make(map[string]string)

	// Reset commit data
	m.commits = make([]*models.Commit, 0, 10000)
	m.commitsByHash = make(map[string]*models.Commit)
	m.commitsByUser = make(map[string][]*models.Commit)
	m.commitsByRepo = make(map[string][]*models.Commit)
	m.needsSort = false

	// Reset PR data
	m.prsByRepo = make(map[string]map[int]*models.PullRequest)
	m.prsByAuthor = make(map[string][]*models.PullRequest)
	m.prsByID = make(map[int]*models.PullRequest)
	m.prsByEmail = make(map[string][]*models.PullRequest)
	m.nextPRID = 1

	// Reset review data
	m.reviewsByID = make(map[int]*models.Review)
	m.reviewsByPRID = make(map[int][]*models.Review)
	m.reviewsByReviewer = make(map[string][]*models.Review)

	// Reset issue data
	m.issuesByRepo = make(map[string]map[int]*models.Issue)
	m.issuesByState = make(map[string]map[models.IssueState][]*models.Issue)

	// Reset review comments
	m.reviewComments = make(map[string]map[int][]*models.ReviewComment)

	// Reset event data
	m.modelUsage = make([]*models.ModelUsageEvent, 0, 1000)
	m.clientVersions = make([]*models.ClientVersionEvent, 0, 1000)
	m.fileExtensions = make([]*models.FileExtensionEvent, 0, 5000)
	m.mcpTools = make([]*models.MCPToolEvent, 0, 2000)
	m.commands = make([]*models.CommandEvent, 0, 2000)
	m.plans = make([]*models.PlanEvent, 0, 1500)
	m.askModes = make([]*models.AskModeEvent, 0, 1500)

	return nil
}

// GetStats returns current storage statistics.
// Used by Admin API for observability and to calculate deltas during regeneration.
// This method is thread-safe.
func (m *MemoryStore) GetStats() StorageStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Count issues across all repos
	issueCount := 0
	for _, repoIssues := range m.issuesByRepo {
		issueCount += len(repoIssues)
	}

	// Count reviews by ID (most accurate)
	reviewCount := len(m.reviewsByID)

	return StorageStats{
		Commits:        len(m.commits),
		PullRequests:   len(m.prsByID),
		Reviews:        reviewCount,
		Issues:         issueCount,
		Developers:     len(m.developers),
		ModelUsage:     len(m.modelUsage),
		ClientVersions: len(m.clientVersions),
		FileExtensions: len(m.fileExtensions),
		MCPTools:       len(m.mcpTools),
		Commands:       len(m.commands),
		Plans:          len(m.plans),
		AskModes:       len(m.askModes),
	}
}
