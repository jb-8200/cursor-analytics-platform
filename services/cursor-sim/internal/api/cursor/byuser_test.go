package cursor

import (
	"net/http/httptest"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestByUserAgentEdits_Success(t *testing.T) {
	store := setupTeamTestStore() // Reuse team test setup
	handler := ByUserAgentEdits(store)

	req := httptest.NewRequest("GET", "/analytics/by-user/agent-edits", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestByUserTabs_Success(t *testing.T) {
	store := setupTeamTestStore()
	handler := ByUserTabs(store)

	req := httptest.NewRequest("GET", "/analytics/by-user/tabs", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

// Stub endpoints
func TestByUserModels_Stub(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := ByUserModels(store)

	req := httptest.NewRequest("GET", "/analytics/by-user/models", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestByUserClientVersions_Stub(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := ByUserClientVersions(store)

	req := httptest.NewRequest("GET", "/analytics/by-user/client-versions", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestByUserTopFileExtensions_Stub(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := ByUserTopFileExtensions(store)

	req := httptest.NewRequest("GET", "/analytics/by-user/top-file-extensions", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestByUserMCP_Stub(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := ByUserMCP(store)

	req := httptest.NewRequest("GET", "/analytics/by-user/mcp", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestByUserCommands_Stub(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := ByUserCommands(store)

	req := httptest.NewRequest("GET", "/analytics/by-user/commands", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestByUserPlans_Stub(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := ByUserPlans(store)

	req := httptest.NewRequest("GET", "/analytics/by-user/plans", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestByUserAskMode_Stub(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := ByUserAskMode(store)

	req := httptest.NewRequest("GET", "/analytics/by-user/ask-mode", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}
