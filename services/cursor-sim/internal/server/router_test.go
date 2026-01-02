package server

import (
	"net/http/httptest"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	store := storage.NewMemoryStore()
	apiKey := "test-key"

	router := NewRouter(store, apiKey)

	assert.NotNil(t, router)
}

func TestRouter_HealthEndpoint(t *testing.T) {
	store := storage.NewMemoryStore()
	router := NewRouter(store, "test-key")

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestRouter_TeamsMembers(t *testing.T) {
	store := storage.NewMemoryStore()
	router := NewRouter(store, "test-key")

	req := httptest.NewRequest("GET", "/teams/members", nil)
	req.SetBasicAuth("test-key", "")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestRouter_AICodeCommits(t *testing.T) {
	store := storage.NewMemoryStore()
	router := NewRouter(store, "test-key")

	req := httptest.NewRequest("GET", "/analytics/ai-code/commits", nil)
	req.SetBasicAuth("test-key", "")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestRouter_WithoutAuth(t *testing.T) {
	store := storage.NewMemoryStore()
	router := NewRouter(store, "test-key")

	req := httptest.NewRequest("GET", "/teams/members", nil)
	// No auth header
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should return 401 Unauthorized
	assert.Equal(t, 401, rec.Code)
}

func TestRouter_NotFound(t *testing.T) {
	store := storage.NewMemoryStore()
	router := NewRouter(store, "test-key")

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	req.SetBasicAuth("test-key", "")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should return 404 Not Found
	assert.Equal(t, 404, rec.Code)
}
