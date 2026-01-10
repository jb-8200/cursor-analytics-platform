package storage

import (
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMemoryStore_ClearAllData tests that ClearAllData resets all data structures.
func TestMemoryStore_ClearAllData(t *testing.T) {
	store := NewMemoryStore()

	// Populate store with test data
	developers := []seed.Developer{
		{UserID: "user1", Email: "user1@example.com", Name: "User 1"},
		{UserID: "user2", Email: "user2@example.com", Name: "User 2"},
	}
	err := store.LoadDevelopers(developers)
	require.NoError(t, err)

	// Add commits
	err = store.AddCommit(models.Commit{
		CommitHash: "abc123",
		UserID:     "user1",
		RepoName:   "test/repo",
		CommitTs:   time.Now(),
	})
	require.NoError(t, err)

	// Add PRs
	err = store.AddPR(models.PullRequest{
		ID:          1,
		Number:      100,
		RepoName:    "test/repo",
		AuthorID:    "user1",
		AuthorEmail: "user1@example.com",
		State:       models.PRStateOpen,
	})
	require.NoError(t, err)

	// Add reviews
	err = store.StoreReview(models.Review{
		ID:       1,
		PRID:     1,
		Reviewer: "user2@example.com",
		State:    "approved",
	})
	require.NoError(t, err)

	// Add issues
	err = store.StoreIssue(models.Issue{
		Number:   1,
		RepoName: "test/repo",
		State:    models.IssueStateOpen,
	})
	require.NoError(t, err)

	// Add events
	err = store.AddModelUsage(models.ModelUsageEvent{
		UserID:    "user1",
		ModelName: "claude-sonnet-4",
		UsageType: "code",
		Timestamp: time.Now(),
	})
	require.NoError(t, err)

	err = store.AddMCPTool(models.MCPToolEvent{
		UserID:        "user1",
		ToolName:      "read_file",
		MCPServerName: "filesystem",
		Timestamp:     time.Now(),
	})
	require.NoError(t, err)

	// Verify data exists
	statsBefore := store.GetStats()
	assert.Equal(t, 2, statsBefore.Developers)
	assert.Equal(t, 1, statsBefore.Commits)
	assert.Equal(t, 1, statsBefore.PullRequests)
	assert.Equal(t, 1, statsBefore.Reviews)
	assert.Equal(t, 1, statsBefore.Issues)
	assert.Equal(t, 1, statsBefore.ModelUsage)
	assert.Equal(t, 1, statsBefore.MCPTools)

	// Clear all data
	err = store.ClearAllData()
	require.NoError(t, err)

	// Verify all data is cleared
	statsAfter := store.GetStats()
	assert.Equal(t, 0, statsAfter.Developers, "Developers should be cleared")
	assert.Equal(t, 0, statsAfter.Commits, "Commits should be cleared")
	assert.Equal(t, 0, statsAfter.PullRequests, "PRs should be cleared")
	assert.Equal(t, 0, statsAfter.Reviews, "Reviews should be cleared")
	assert.Equal(t, 0, statsAfter.Issues, "Issues should be cleared")
	assert.Equal(t, 0, statsAfter.ModelUsage, "Model usage should be cleared")
	assert.Equal(t, 0, statsAfter.ClientVersions, "Client versions should be cleared")
	assert.Equal(t, 0, statsAfter.FileExtensions, "File extensions should be cleared")
	assert.Equal(t, 0, statsAfter.MCPTools, "MCP tools should be cleared")
	assert.Equal(t, 0, statsAfter.Commands, "Commands should be cleared")
	assert.Equal(t, 0, statsAfter.Plans, "Plans should be cleared")
	assert.Equal(t, 0, statsAfter.AskModes, "Ask modes should be cleared")
}

// TestMemoryStore_GetStats tests that GetStats returns accurate counts.
func TestMemoryStore_GetStats(t *testing.T) {
	store := NewMemoryStore()

	// Empty store
	stats := store.GetStats()
	assert.Equal(t, 0, stats.Developers)
	assert.Equal(t, 0, stats.Commits)
	assert.Equal(t, 0, stats.PullRequests)

	// Add developers
	developers := []seed.Developer{
		{UserID: "user1", Email: "user1@example.com"},
		{UserID: "user2", Email: "user2@example.com"},
		{UserID: "user3", Email: "user3@example.com"},
	}
	err := store.LoadDevelopers(developers)
	require.NoError(t, err)

	stats = store.GetStats()
	assert.Equal(t, 3, stats.Developers)

	// Add commits
	for i := 0; i < 10; i++ {
		err = store.AddCommit(models.Commit{
			CommitHash: "commit" + string(rune(i)),
			UserID:     "user1",
			CommitTs:   time.Now(),
		})
		require.NoError(t, err)
	}

	stats = store.GetStats()
	assert.Equal(t, 10, stats.Commits)

	// Add PRs
	for i := 0; i < 5; i++ {
		err = store.AddPR(models.PullRequest{
			Number:      i + 1,
			RepoName:    "test/repo",
			AuthorID:    "user1",
			AuthorEmail: "user1@example.com",
		})
		require.NoError(t, err)
	}

	stats = store.GetStats()
	assert.Equal(t, 5, stats.PullRequests)

	// Add reviews (2 per PR)
	for prID := 1; prID <= 5; prID++ {
		for reviewID := 0; reviewID < 2; reviewID++ {
			err = store.StoreReview(models.Review{
				ID:       prID*10 + reviewID,
				PRID:     prID,
				Reviewer: "user2@example.com",
			})
			require.NoError(t, err)
		}
	}

	stats = store.GetStats()
	assert.Equal(t, 10, stats.Reviews)

	// Add issues (2 repos with 3 issues each)
	for repo := 1; repo <= 2; repo++ {
		for issueNum := 1; issueNum <= 3; issueNum++ {
			err = store.StoreIssue(models.Issue{
				Number:   issueNum,
				RepoName: "test/repo" + string(rune(repo)),
				State:    models.IssueStateOpen,
			})
			require.NoError(t, err)
		}
	}

	stats = store.GetStats()
	assert.Equal(t, 6, stats.Issues)

	// Add various events
	for i := 0; i < 8; i++ {
		err = store.AddModelUsage(models.ModelUsageEvent{
			UserID:    "user1",
			Timestamp: time.Now(),
		})
		require.NoError(t, err)
	}

	for i := 0; i < 12; i++ {
		err = store.AddClientVersion(models.ClientVersionEvent{
			UserID:    "user1",
			Timestamp: time.Now(),
		})
		require.NoError(t, err)
	}

	for i := 0; i < 15; i++ {
		err = store.AddFileExtension(models.FileExtensionEvent{
			UserID:    "user1",
			Timestamp: time.Now(),
		})
		require.NoError(t, err)
	}

	for i := 0; i < 4; i++ {
		err = store.AddMCPTool(models.MCPToolEvent{
			UserID:    "user1",
			Timestamp: time.Now(),
		})
		require.NoError(t, err)
	}

	for i := 0; i < 6; i++ {
		err = store.AddCommand(models.CommandEvent{
			UserID:    "user1",
			Timestamp: time.Now(),
		})
		require.NoError(t, err)
	}

	for i := 0; i < 3; i++ {
		err = store.AddPlan(models.PlanEvent{
			UserID:    "user1",
			Timestamp: time.Now(),
		})
		require.NoError(t, err)
	}

	for i := 0; i < 5; i++ {
		err = store.AddAskMode(models.AskModeEvent{
			UserID:    "user1",
			Timestamp: time.Now(),
		})
		require.NoError(t, err)
	}

	// Verify all counts
	stats = store.GetStats()
	assert.Equal(t, 3, stats.Developers)
	assert.Equal(t, 10, stats.Commits)
	assert.Equal(t, 5, stats.PullRequests)
	assert.Equal(t, 10, stats.Reviews)
	assert.Equal(t, 6, stats.Issues)
	assert.Equal(t, 8, stats.ModelUsage)
	assert.Equal(t, 12, stats.ClientVersions)
	assert.Equal(t, 15, stats.FileExtensions)
	assert.Equal(t, 4, stats.MCPTools)
	assert.Equal(t, 6, stats.Commands)
	assert.Equal(t, 3, stats.Plans)
	assert.Equal(t, 5, stats.AskModes)
}

// TestMemoryStore_ClearAllData_ThreadSafe tests that ClearAllData is thread-safe.
func TestMemoryStore_ClearAllData_ThreadSafe(t *testing.T) {
	store := NewMemoryStore()

	// Add some initial data
	developers := []seed.Developer{
		{UserID: "user1", Email: "user1@example.com"},
	}
	err := store.LoadDevelopers(developers)
	require.NoError(t, err)

	// Clear data concurrently
	done := make(chan bool, 2)

	go func() {
		err := store.ClearAllData()
		assert.NoError(t, err)
		done <- true
	}()

	go func() {
		err := store.ClearAllData()
		assert.NoError(t, err)
		done <- true
	}()

	<-done
	<-done

	// Verify data is cleared
	stats := store.GetStats()
	assert.Equal(t, 0, stats.Developers)
	assert.Equal(t, 0, stats.Commits)
}

// TestMemoryStore_GetStats_ThreadSafe tests that GetStats is thread-safe.
func TestMemoryStore_GetStats_ThreadSafe(t *testing.T) {
	store := NewMemoryStore()

	// Add some data
	developers := []seed.Developer{
		{UserID: "user1", Email: "user1@example.com"},
		{UserID: "user2", Email: "user2@example.com"},
	}
	err := store.LoadDevelopers(developers)
	require.NoError(t, err)

	// Read stats concurrently while adding data
	done := make(chan bool, 10)

	// Multiple readers
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				stats := store.GetStats()
				// Stats should be consistent (developers don't change)
				assert.Equal(t, 2, stats.Developers)
			}
			done <- true
		}()
	}

	// Writer adding commits
	go func() {
		for i := 0; i < 50; i++ {
			err := store.AddCommit(models.Commit{
				CommitHash: "commit" + string(rune(i)),
				UserID:     "user1",
				CommitTs:   time.Now(),
			})
			assert.NoError(t, err)
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 6; i++ {
		<-done
	}

	// Final stats should be accurate
	stats := store.GetStats()
	assert.Equal(t, 2, stats.Developers)
	assert.Equal(t, 50, stats.Commits)
}
