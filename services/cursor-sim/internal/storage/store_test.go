package storage

import (
	"sync"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore_LoadDevelopers(t *testing.T) {
	store := NewMemoryStore()

	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob"},
	}

	err := store.LoadDevelopers(developers)
	require.NoError(t, err)

	// Verify retrieval by user ID
	dev, err := store.GetDeveloper("user_001")
	require.NoError(t, err)
	assert.Equal(t, "alice@example.com", dev.Email)

	// Verify retrieval by email
	dev, err = store.GetDeveloperByEmail("bob@example.com")
	require.NoError(t, err)
	assert.Equal(t, "user_002", dev.UserID)

	// Verify list all
	all := store.ListDevelopers()
	assert.Len(t, all, 2)
}

func TestMemoryStore_GetDeveloper_NotFound(t *testing.T) {
	store := NewMemoryStore()

	dev, err := store.GetDeveloper("nonexistent")
	require.Error(t, err)
	assert.Nil(t, dev)
	assert.Contains(t, err.Error(), "developer not found")
}

func TestMemoryStore_AddCommit(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	commit := models.Commit{
		CommitHash:      "abc123",
		UserID:          "user_001",
		UserEmail:       "test@example.com",
		RepoName:        "test/repo",
		TotalLinesAdded: 100,
		CommitTs:        now,
		CreatedAt:       now,
	}

	err := store.AddCommit(commit)
	require.NoError(t, err)

	// Verify retrieval by hash
	retrieved, err := store.GetCommitByHash("abc123")
	require.NoError(t, err)
	assert.Equal(t, "user_001", retrieved.UserID)
	assert.Equal(t, 100, retrieved.TotalLinesAdded)
}

func TestMemoryStore_GetCommitByHash_NotFound(t *testing.T) {
	store := NewMemoryStore()

	commit, err := store.GetCommitByHash("nonexistent")
	require.Error(t, err)
	assert.Nil(t, commit)
}

func TestMemoryStore_GetCommitsByTimeRange(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	// Add commits at different times
	commits := []models.Commit{
		{CommitHash: "c1", UserID: "user_001", CommitTs: now.Add(-2 * time.Hour), CreatedAt: now},
		{CommitHash: "c2", UserID: "user_001", CommitTs: now.Add(-1 * time.Hour), CreatedAt: now},
		{CommitHash: "c3", UserID: "user_001", CommitTs: now, CreatedAt: now},
		{CommitHash: "c4", UserID: "user_001", CommitTs: now.Add(1 * time.Hour), CreatedAt: now},
	}

	for _, c := range commits {
		err := store.AddCommit(c)
		require.NoError(t, err)
	}

	// Query range: -90 minutes to +30 minutes
	from := now.Add(-90 * time.Minute)
	to := now.Add(30 * time.Minute)

	results := store.GetCommitsByTimeRange(from, to)

	// Should return c2 and c3
	assert.Len(t, results, 2)
	hashes := []string{results[0].CommitHash, results[1].CommitHash}
	assert.Contains(t, hashes, "c2")
	assert.Contains(t, hashes, "c3")
}

func TestMemoryStore_GetCommitsByUser(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	commits := []models.Commit{
		{CommitHash: "c1", UserID: "user_001", CommitTs: now, CreatedAt: now},
		{CommitHash: "c2", UserID: "user_002", CommitTs: now, CreatedAt: now},
		{CommitHash: "c3", UserID: "user_001", CommitTs: now, CreatedAt: now},
	}

	for _, c := range commits {
		err := store.AddCommit(c)
		require.NoError(t, err)
	}

	// Get commits for user_001
	from := now.Add(-1 * time.Hour)
	to := now.Add(1 * time.Hour)
	results := store.GetCommitsByUser("user_001", from, to)

	assert.Len(t, results, 2)
	for _, c := range results {
		assert.Equal(t, "user_001", c.UserID)
	}
}

func TestMemoryStore_GetCommitsByRepo(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	commits := []models.Commit{
		{CommitHash: "c1", RepoName: "acme/api", UserID: "user_001", CommitTs: now, CreatedAt: now},
		{CommitHash: "c2", RepoName: "acme/web", UserID: "user_001", CommitTs: now, CreatedAt: now},
		{CommitHash: "c3", RepoName: "acme/api", UserID: "user_002", CommitTs: now, CreatedAt: now},
	}

	for _, c := range commits {
		err := store.AddCommit(c)
		require.NoError(t, err)
	}

	from := now.Add(-1 * time.Hour)
	to := now.Add(1 * time.Hour)
	results := store.GetCommitsByRepo("acme/api", from, to)

	assert.Len(t, results, 2)
	for _, c := range results {
		assert.Equal(t, "acme/api", c.RepoName)
	}
}

func TestMemoryStore_ConcurrentAccess(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	// Concurrent writes
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			commit := models.Commit{
				CommitHash: string(rune('a' + idx)),
				UserID:     "user_001",
				CommitTs:   now,
				CreatedAt:  now,
			}
			if err := store.AddCommit(commit); err != nil {
				errChan <- err
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			from := now.Add(-1 * time.Hour)
			to := now.Add(1 * time.Hour)
			_ = store.GetCommitsByTimeRange(from, to)
		}()
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		t.Errorf("concurrent access error: %v", err)
	}
}

func TestMemoryStore_TimeRangePerformance(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	// Add 1000 commits over 30 days
	for i := 0; i < 1000; i++ {
		commit := models.Commit{
			CommitHash: string(rune(i)),
			UserID:     "user_001",
			CommitTs:   now.Add(-time.Duration(i) * time.Hour),
			CreatedAt:  now,
		}
		err := store.AddCommit(commit)
		require.NoError(t, err)
	}

	// Query a 7-day range (should be fast)
	start := time.Now()
	from := now.Add(-7 * 24 * time.Hour)
	to := now
	results := store.GetCommitsByTimeRange(from, to)
	elapsed := time.Since(start)

	// Should complete in < 10ms
	assert.Less(t, elapsed, 10*time.Millisecond, "query should be fast")
	assert.Greater(t, len(results), 0, "should find commits in range")
}

func TestMemoryStore_EmptyStore(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	// All queries should return empty results
	assert.Len(t, store.ListDevelopers(), 0)
	assert.Len(t, store.GetCommitsByTimeRange(now, now.Add(1*time.Hour)), 0)
	assert.Len(t, store.GetCommitsByUser("user_001", now, now.Add(1*time.Hour)), 0)
	assert.Len(t, store.GetCommitsByRepo("test/repo", now, now.Add(1*time.Hour)), 0)
}

func TestMemoryStore_DuplicateCommitHash(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	commit1 := models.Commit{
		CommitHash: "abc123",
		UserID:     "user_001",
		CommitTs:   now,
		CreatedAt:  now,
	}

	commit2 := models.Commit{
		CommitHash: "abc123", // Same hash
		UserID:     "user_002",
		CommitTs:   now,
		CreatedAt:  now,
	}

	err := store.AddCommit(commit1)
	require.NoError(t, err)

	// Second commit with same hash should be rejected or overwrite
	err = store.AddCommit(commit2)
	require.NoError(t, err) // We allow overwrite

	// Should get the latest commit
	retrieved, err := store.GetCommitByHash("abc123")
	require.NoError(t, err)
	assert.Equal(t, "user_002", retrieved.UserID)
}
