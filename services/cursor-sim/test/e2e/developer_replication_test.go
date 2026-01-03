package e2e

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/cursor"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/server"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_DeveloperReplication_Downsample verifies that when requesting fewer developers than seed count,
// the system randomly samples from the seed developers.
func TestE2E_DeveloperReplication_Downsample(t *testing.T) {
	// Load test seed data (has 2 developers)
	seedPath := "../../testdata/valid_seed.json"
	targetCount := 1 // Request fewer than seed count

	// Use fixed RNG for reproducibility
	rng := rand.New(rand.NewSource(42))

	// Load seed with replication
	seedData, developers, err := seed.LoadSeedWithReplication(seedPath, targetCount, rng)
	require.NoError(t, err, "Failed to load seed with replication")

	// Verify we got exactly 1 developer
	require.Len(t, developers, targetCount, "Should have exactly %d developer(s)", targetCount)

	// Initialize storage
	store := storage.NewMemoryStore()
	err = store.LoadDevelopers(developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate minimal data to simulate real environment
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 1) // 1 day of data
	require.NoError(t, err, "Failed to generate commits")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Query /teams/members endpoint
	req := httptest.NewRequest("GET", "/teams/members", nil)
	req.SetBasicAuth("test-api-key", "")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, 200, rec.Code, "Expected 200 OK")

	// Parse response
	var response cursor.TeamMembersResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")

	// Verify exactly 1 developer returned
	assert.Len(t, response.TeamMembers, targetCount, "Expected %d team member(s)", targetCount)

	// Verify member is one of the seed developers
	validEmails := map[string]bool{
		"alice@example.com": true,
		"bob@example.com":   true,
	}
	assert.True(t, validEmails[response.TeamMembers[0].Email], "Member should be from seed data")
}

// TestE2E_DeveloperReplication_ExactMatch verifies that when requesting the same number of developers
// as the seed count, all seed developers are returned.
func TestE2E_DeveloperReplication_ExactMatch(t *testing.T) {
	// Load test seed data (has 2 developers)
	seedPath := "../../testdata/valid_seed.json"
	targetCount := 2 // Request exact match

	// Use fixed RNG for reproducibility
	rng := rand.New(rand.NewSource(42))

	// Load seed with replication
	seedData, developers, err := seed.LoadSeedWithReplication(seedPath, targetCount, rng)
	require.NoError(t, err, "Failed to load seed with replication")

	// Verify we got exactly 2 developers
	require.Len(t, developers, targetCount, "Should have exactly %d developers", targetCount)

	// Initialize storage
	store := storage.NewMemoryStore()
	err = store.LoadDevelopers(developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate minimal data
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 1)
	require.NoError(t, err, "Failed to generate commits")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Query /teams/members endpoint
	req := httptest.NewRequest("GET", "/teams/members", nil)
	req.SetBasicAuth("test-api-key", "")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, 200, rec.Code, "Expected 200 OK")

	// Parse response
	var response cursor.TeamMembersResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")

	// Verify exactly 2 developers returned
	assert.Len(t, response.TeamMembers, targetCount, "Expected %d team members", targetCount)

	// Verify all emails are present
	emails := make(map[string]bool)
	for _, member := range response.TeamMembers {
		emails[member.Email] = true
	}
	assert.True(t, emails["alice@example.com"], "Should have alice@example.com")
	assert.True(t, emails["bob@example.com"], "Should have bob@example.com")
}

// TestE2E_DeveloperReplication_ScaleUp verifies that when requesting more developers than seed count,
// developers are cloned with unique IDs and emails.
func TestE2E_DeveloperReplication_ScaleUp(t *testing.T) {
	// Load test seed data (has 2 developers)
	seedPath := "../../testdata/valid_seed.json"
	targetCount := 5 // Request more than seed count to trigger cloning

	// Use fixed RNG for reproducibility
	rng := rand.New(rand.NewSource(42))

	// Load seed with replication
	seedData, developers, err := seed.LoadSeedWithReplication(seedPath, targetCount, rng)
	require.NoError(t, err, "Failed to load seed with replication")

	// Verify we got exactly 5 developers
	require.Len(t, developers, targetCount, "Should have exactly %d developers", targetCount)

	// Initialize storage
	store := storage.NewMemoryStore()
	err = store.LoadDevelopers(developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate minimal data
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 1)
	require.NoError(t, err, "Failed to generate commits")

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Query /teams/members endpoint
	req := httptest.NewRequest("GET", "/teams/members", nil)
	req.SetBasicAuth("test-api-key", "")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, 200, rec.Code, "Expected 200 OK")

	// Parse response
	var response cursor.TeamMembersResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")

	// Verify exactly 5 developers returned
	assert.Len(t, response.TeamMembers, targetCount, "Expected %d team members", targetCount)

	// Verify all emails are unique
	emails := make(map[string]bool)
	for _, member := range response.TeamMembers {
		assert.False(t, emails[member.Email], "Email %s should be unique", member.Email)
		emails[member.Email] = true
	}

	// Verify we have cloned developers (should have "clone" in email)
	cloneCount := 0
	for email := range emails {
		if containsCloneMarker(email) {
			cloneCount++
		}
	}
	assert.Equal(t, 3, cloneCount, "Should have 3 cloned developers (5 total - 2 original)")

	// Verify original developers are present
	assert.True(t, emails["alice@example.com"] || containsAliceClone(emails), "Should have alice or alice clone")
	assert.True(t, emails["bob@example.com"] || containsBobClone(emails), "Should have bob or bob clone")
}

// TestE2E_DeveloperReplication_UniqueIDs verifies that all developer IDs are unique
// when using replication, including clones.
func TestE2E_DeveloperReplication_UniqueIDs(t *testing.T) {
	seedPath := "../../testdata/valid_seed.json"
	targetCount := 10 // Large number to ensure multiple clone iterations

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Load seed with replication
	_, developers, err := seed.LoadSeedWithReplication(seedPath, targetCount, rng)
	require.NoError(t, err, "Failed to load seed with replication")

	// Verify all IDs are unique
	ids := make(map[string]bool)
	emails := make(map[string]bool)

	for _, dev := range developers {
		assert.False(t, ids[dev.UserID], "User ID %s should be unique", dev.UserID)
		assert.False(t, emails[dev.Email], "Email %s should be unique", dev.Email)
		ids[dev.UserID] = true
		emails[dev.Email] = true
	}

	assert.Len(t, ids, targetCount, "Should have %d unique user IDs", targetCount)
	assert.Len(t, emails, targetCount, "Should have %d unique emails", targetCount)
}

// Helper functions

func containsCloneMarker(email string) bool {
	return len(email) > 5 && email[0:5] == "clone"
}

func containsAliceClone(emails map[string]bool) bool {
	for email := range emails {
		if containsCloneMarker(email) && len(email) > 12 && email[6:11] == "alice" {
			return true
		}
	}
	return false
}

func containsBobClone(emails map[string]bool) bool {
	for email := range emails {
		if containsCloneMarker(email) && len(email) > 10 && email[6:9] == "bob" {
			return true
		}
	}
	return false
}
