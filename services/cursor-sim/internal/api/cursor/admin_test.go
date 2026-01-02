package cursor

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeamMembers_Success(t *testing.T) {
	store := storage.NewMemoryStore()
	developers := []seed.Developer{
		{
			UserID: "user_001",
			Email:  "alice@example.com",
			Name:   "Alice Chen",
		},
		{
			UserID: "user_002",
			Email:  "bob@example.com",
			Name:   "Bob Smith",
		},
	}
	err := store.LoadDevelopers(developers)
	require.NoError(t, err)

	handler := TeamMembers(store)

	req := httptest.NewRequest("GET", "/teams/members", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response TeamMembersResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.TeamMembers, 2)
	assert.Equal(t, "Alice Chen", response.TeamMembers[0].Name)
	assert.Equal(t, "alice@example.com", response.TeamMembers[0].Email)
	assert.Equal(t, "member", response.TeamMembers[0].Role)
	assert.Equal(t, "Bob Smith", response.TeamMembers[1].Name)
}

func TestTeamMembers_EmptyStore(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := TeamMembers(store)

	req := httptest.NewRequest("GET", "/teams/members", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	var response TeamMembersResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.TeamMembers, 0)
}

func TestTeamMembers_MultipleMembers(t *testing.T) {
	store := storage.NewMemoryStore()

	// Create 10 developers
	developers := make([]seed.Developer, 10)
	for i := 0; i < 10; i++ {
		developers[i] = seed.Developer{
			UserID: "user_" + string(rune('0'+i)),
			Email:  "dev" + string(rune('0'+i)) + "@example.com",
			Name:   "Developer " + string(rune('0'+i)),
		}
	}
	err := store.LoadDevelopers(developers)
	require.NoError(t, err)

	handler := TeamMembers(store)

	req := httptest.NewRequest("GET", "/teams/members", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	var response TeamMembersResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.TeamMembers, 10)
}

func TestTeamMembers_RoleMapping(t *testing.T) {
	store := storage.NewMemoryStore()
	developers := []seed.Developer{
		{UserID: "user_001", Email: "dev1@example.com", Name: "Dev One"},
		{UserID: "user_002", Email: "dev2@example.com", Name: "Dev Two"},
	}
	err := store.LoadDevelopers(developers)
	require.NoError(t, err)

	handler := TeamMembers(store)

	req := httptest.NewRequest("GET", "/teams/members", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var response TeamMembersResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// All members should have "member" role
	for _, member := range response.TeamMembers {
		assert.Equal(t, "member", member.Role)
	}
}
