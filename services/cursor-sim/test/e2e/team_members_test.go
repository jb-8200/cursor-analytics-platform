package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/cursor"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/server"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_TeamMembers verifies that /teams/members returns developers loaded from seed data.
func TestE2E_TeamMembers(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers into storage (this is the fix we're testing)
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err, "Failed to load developers")

	// Generate minimal data (commits) to simulate real environment
	ctx := context.Background()
	commitGen := generator.NewCommitGenerator(seedData, store, "medium")
	err = commitGen.GenerateCommits(ctx, 7, 0) // Just 1 week of data for speed
	require.NoError(t, err, "Failed to generate commits")

	// Create router with loaded data
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/teams/members", nil)
	req.SetBasicAuth("test-api-key", "")

	// Execute request
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code, "Expected 200 OK")

	// Parse response
	var response cursor.TeamMembersResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")

	// Verify team members are present
	assert.Len(t, response.TeamMembers, 2, "Expected 2 team members from seed file")

	// Verify first member (alice@example.com)
	assert.Equal(t, "Alice Developer", response.TeamMembers[0].Name)
	assert.Equal(t, "alice@example.com", response.TeamMembers[0].Email)
	assert.Equal(t, "member", response.TeamMembers[0].Role)

	// Verify second member (bob@example.com)
	assert.Equal(t, "Bob Developer", response.TeamMembers[1].Name)
	assert.Equal(t, "bob@example.com", response.TeamMembers[1].Email)
	assert.Equal(t, "member", response.TeamMembers[1].Role)
}

// TestE2E_TeamMembers_EmptyWithoutLoad verifies that endpoint returns empty when developers not loaded.
func TestE2E_TeamMembers_EmptyWithoutLoad(t *testing.T) {
	// This test documents the bug: without LoadDevelopers(), endpoint returns empty

	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err, "Failed to load seed data")

	// Initialize storage
	store := storage.NewMemoryStore()

	// NOTE: Intentionally NOT calling store.LoadDevelopers() to show the bug

	// Create router
	router := server.NewRouter(store, seedData, "test-api-key")

	// Create test request
	req := httptest.NewRequest("GET", "/teams/members", nil)
	req.SetBasicAuth("test-api-key", "")

	// Execute request
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code, "Expected 200 OK")

	// Parse response
	var response cursor.TeamMembersResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")

	// This is the BUG: returns empty array
	assert.Len(t, response.TeamMembers, 0, "BUG: Without LoadDevelopers(), returns empty")
}
