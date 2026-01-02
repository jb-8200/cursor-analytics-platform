package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListRepos(t *testing.T) {
	store := setupTestStore()

	t.Run("list all repositories", func(t *testing.T) {
		handler := ListRepos(store)

		req := httptest.NewRequest(http.MethodGet, "/repos", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []RepoInfo
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(response), 2, "should return at least 2 repos")

		// Verify repo structure
		for _, repo := range response {
			assert.NotEmpty(t, repo.FullName, "repo should have full_name")
		}
	})
}

func TestGetRepo(t *testing.T) {
	store := setupTestStore()

	t.Run("get existing repo", func(t *testing.T) {
		handler := GetRepo(store)

		req := httptest.NewRequest(http.MethodGet, "/repos/acme/api", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response RepoInfo
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "acme/api", response.FullName)
		assert.GreaterOrEqual(t, response.OpenPRs, 0)
	})

	t.Run("repo not found returns 404", func(t *testing.T) {
		handler := GetRepo(store)

		req := httptest.NewRequest(http.MethodGet, "/repos/nonexistent/repo", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
