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
}

// NewMemoryStore creates a new thread-safe in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		developers:      make(map[string]*seed.Developer),
		developerEmails: make(map[string]string),
		commits:         make([]*models.Commit, 0, 1000),
		commitsByHash:   make(map[string]*models.Commit),
		commitsByUser:   make(map[string][]*models.Commit),
		commitsByRepo:   make(map[string][]*models.Commit),
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

// ListDevelopers returns all developers.
func (m *MemoryStore) ListDevelopers() []seed.Developer {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]seed.Developer, 0, len(m.developers))
	for _, dev := range m.developers {
		result = append(result, *dev)
	}

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
