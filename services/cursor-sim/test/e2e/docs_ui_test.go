package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/server"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createDocsTestConfig() *config.Config {
	return &config.Config{
		Mode:       "runtime",
		Port:       8080,
		Days:       1,
		Velocity:   "medium",
		SeedPath:   "",
		CorpusPath: "",
	}
}

const docsTestVersion = "1.0.0"

func createDocsTestSeedData() *seed.SeedData {
	return &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:    "user_001",
				Email:     "test@example.com",
				Name:      "Test Developer",
				Seniority: "mid",
				WorkingHoursBand: seed.WorkingHours{
					Start: 9,
					End:   18,
				},
			},
		},
		Repositories: []seed.Repository{
			{
				RepoName:      "test/repo",
				DefaultBranch: "main",
			},
		},
	}
}

// TestDocsUI_LoginPageWithoutAuth verifies login page is accessible without authentication
func TestDocsUI_LoginPageWithoutAuth(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createDocsTestSeedData()
	router := server.NewRouter(store, seedData, "test-key", createDocsTestConfig(), docsTestVersion)

	req := httptest.NewRequest("GET", "/docs/login", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")

	// Verify form elements are present
	body := rec.Body.String()
	assert.Contains(t, body, "username")
	assert.Contains(t, body, "password")
	assert.Contains(t, body, "dox")
	assert.Contains(t, body, "dox-a3")
}

// TestDocsUI_UnauthorizedRedirectsToLogin verifies /docs redirects to login without session
func TestDocsUI_UnauthorizedRedirectsToLogin(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createDocsTestSeedData()
	router := server.NewRouter(store, seedData, "test-key", createDocsTestConfig(), docsTestVersion)

	req := httptest.NewRequest("GET", "/docs", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusFound, rec.Code)
	location := rec.Header().Get("Location")
	assert.Contains(t, location, "/docs/login")
}

// TestDocsUI_InvalidCredentials shows error message
func TestDocsUI_InvalidCredentials(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createDocsTestSeedData()
	router := server.NewRouter(store, seedData, "test-key", createDocsTestConfig(), docsTestVersion)

	// POST with invalid credentials
	req := httptest.NewRequest("POST", "/docs/login", strings.NewReader("username=wrong&password=wrong"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "Invalid") // Should show error message
}

// TestDocsUI_ValidCredentialsGrantAccess verifies valid credentials authenticate
func TestDocsUI_ValidCredentialsGrantAccess(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createDocsTestSeedData()
	router := server.NewRouter(store, seedData, "test-key", createDocsTestConfig(), docsTestVersion)

	// POST with valid credentials
	req := httptest.NewRequest("POST", "/docs/login", strings.NewReader("username=dox&password=dox-a3"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should redirect to /docs
	assert.Equal(t, http.StatusFound, rec.Code)
	location := rec.Header().Get("Location")
	assert.Contains(t, location, "/docs")

	// Should set session cookie
	cookies := rec.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "docs_session" && cookie.Value != "" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected docs_session cookie to be set")
}

// TestDocsUI_AuthenticatedAccessToDocsPage verifies authenticated users can access docs
func TestDocsUI_AuthenticatedAccessToDocsPage(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createDocsTestSeedData()
	router := server.NewRouter(store, seedData, "test-key", createDocsTestConfig(), docsTestVersion)

	// First, authenticate and get session cookie
	loginReq := httptest.NewRequest("POST", "/docs/login", strings.NewReader("username=dox&password=dox-a3"))
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginRec := httptest.NewRecorder()
	router.ServeHTTP(loginRec, loginReq)

	// Extract session cookie
	var sessionCookie *http.Cookie
	for _, cookie := range loginRec.Result().Cookies() {
		if cookie.Name == "docs_session" {
			sessionCookie = cookie
			break
		}
	}
	require.NotNil(t, sessionCookie, "Expected session cookie")

	// Now access /docs with session cookie
	docsReq := httptest.NewRequest("GET", "/docs", nil)
	docsReq.AddCookie(sessionCookie)
	docsRec := httptest.NewRecorder()

	router.ServeHTTP(docsRec, docsReq)

	assert.Equal(t, http.StatusOK, docsRec.Code)
	body := docsRec.Body.String()
	assert.Contains(t, body, "Scalar")      // Should contain Scalar reference
	assert.Contains(t, body, "Cursor API")  // Should contain API title
}

// TestDocsUI_OpenAPISpecServed verifies cursor-api.yaml is served
func TestDocsUI_OpenAPISpecServed(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createDocsTestSeedData()
	router := server.NewRouter(store, seedData, "test-key", createDocsTestConfig(), docsTestVersion)

	req := httptest.NewRequest("GET", "/docs/openapi/cursor-api.yaml", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/yaml", rec.Header().Get("Content-Type"))

	// Verify it's valid YAML by checking for common OpenAPI fields
	body := rec.Body.String()
	assert.Contains(t, body, "openapi:")
	assert.Contains(t, body, "info:")
	assert.Contains(t, body, "paths:")
}

// TestDocsUI_GitHubAPISpecServed verifies github-sim-api.yaml is served
func TestDocsUI_GitHubAPISpecServed(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createDocsTestSeedData()
	router := server.NewRouter(store, seedData, "test-key", createDocsTestConfig(), docsTestVersion)

	req := httptest.NewRequest("GET", "/docs/openapi/github-sim-api.yaml", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/yaml", rec.Header().Get("Content-Type"))

	// Verify it's valid YAML
	body := rec.Body.String()
	assert.Contains(t, body, "openapi:")
	assert.Contains(t, body, "info:")
	assert.Contains(t, body, "paths:")
}

// TestDocsUI_NotFoundSpec returns 404 for missing spec
func TestDocsUI_NotFoundSpec(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createDocsTestSeedData()
	router := server.NewRouter(store, seedData, "test-key", createDocsTestConfig(), docsTestVersion)

	req := httptest.NewRequest("GET", "/docs/openapi/nonexistent.yaml", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestDocsUI_LogoutClearsSession verifies logout clears session
func TestDocsUI_LogoutClearsSession(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createDocsTestSeedData()
	router := server.NewRouter(store, seedData, "test-key", createDocsTestConfig(), docsTestVersion)

	// First, authenticate
	loginReq := httptest.NewRequest("POST", "/docs/login", strings.NewReader("username=dox&password=dox-a3"))
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginRec := httptest.NewRecorder()
	router.ServeHTTP(loginRec, loginReq)

	// Extract session cookie
	var sessionCookie *http.Cookie
	for _, cookie := range loginRec.Result().Cookies() {
		if cookie.Name == "docs_session" {
			sessionCookie = cookie
			break
		}
	}
	require.NotNil(t, sessionCookie)

	// Now logout
	logoutReq := httptest.NewRequest("GET", "/docs/logout", nil)
	logoutReq.AddCookie(sessionCookie)
	logoutRec := httptest.NewRecorder()

	router.ServeHTTP(logoutRec, logoutReq)

	// Should redirect to login
	assert.Equal(t, http.StatusFound, logoutRec.Code)
	location := logoutRec.Header().Get("Location")
	assert.Contains(t, location, "/docs/login")

	// After logout, accessing /docs should redirect to login
	docsReq := httptest.NewRequest("GET", "/docs", nil)
	docsReq.AddCookie(sessionCookie) // Old session cookie
	docsRec := httptest.NewRecorder()

	router.ServeHTTP(docsRec, docsReq)

	// Should redirect back to login (session is invalid)
	assert.Equal(t, http.StatusFound, docsRec.Code)
}

// TestDocsUI_EndToEndWorkflow tests complete user journey
func TestDocsUI_EndToEndWorkflow(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createDocsTestSeedData()
	router := server.NewRouter(store, seedData, "test-key", createDocsTestConfig(), docsTestVersion)

	// Step 1: Visit /docs without session -> redirect to login
	req1 := httptest.NewRequest("GET", "/docs", nil)
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	assert.Equal(t, http.StatusFound, rec1.Code)

	// Step 2: Visit login page -> see form
	req2 := httptest.NewRequest("GET", "/docs/login", nil)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	assert.Equal(t, http.StatusOK, rec2.Code)
	assert.Contains(t, rec2.Body.String(), "username")

	// Step 3: Submit login with valid credentials
	req3 := httptest.NewRequest("POST", "/docs/login", strings.NewReader("username=dox&password=dox-a3"))
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec3 := httptest.NewRecorder()
	router.ServeHTTP(rec3, req3)
	assert.Equal(t, http.StatusFound, rec3.Code)

	// Extract session cookie
	var sessionCookie *http.Cookie
	for _, cookie := range rec3.Result().Cookies() {
		if cookie.Name == "docs_session" {
			sessionCookie = cookie
			break
		}
	}
	require.NotNil(t, sessionCookie)

	// Step 4: Access docs page with session -> see Scalar UI
	req4 := httptest.NewRequest("GET", "/docs", nil)
	req4.AddCookie(sessionCookie)
	rec4 := httptest.NewRecorder()
	router.ServeHTTP(rec4, req4)
	assert.Equal(t, http.StatusOK, rec4.Code)
	assert.Contains(t, rec4.Body.String(), "Scalar")

	// Step 5: Get OpenAPI spec
	req5 := httptest.NewRequest("GET", "/docs/openapi/cursor-api.yaml", nil)
	rec5 := httptest.NewRecorder()
	router.ServeHTTP(rec5, req5)
	assert.Equal(t, http.StatusOK, rec5.Code)
	assert.Contains(t, rec5.Body.String(), "openapi:")
}

// TestDocsUI_DocsDoesntRequireAPIAuth verifies /docs has its own auth, not API key auth
func TestDocsUI_DocsDoesntRequireAPIAuth(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createDocsTestSeedData()
	router := server.NewRouter(store, seedData, "test-key", createDocsTestConfig(), docsTestVersion)

	// Visit login without API key - should work
	req := httptest.NewRequest("GET", "/docs/login", nil)
	// No basic auth set
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should succeed (200), not 401 (unauthorized)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "username")
}

// TestDocsUI_APIAuthStillRequiredForAPIs verifies API endpoints still require auth
func TestDocsUI_APIAuthStillRequiredForAPIs(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := createDocsTestSeedData()
	router := server.NewRouter(store, seedData, "test-key", createDocsTestConfig(), docsTestVersion)

	// Try to access API without auth - should fail
	req := httptest.NewRequest("GET", "/teams/members", nil)
	// No basic auth set
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should fail with 401
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// Now try with correct auth - should work
	req2 := httptest.NewRequest("GET", "/teams/members", nil)
	req2.SetBasicAuth("test-key", "")
	rec2 := httptest.NewRecorder()

	router.ServeHTTP(rec2, req2)

	// Should succeed
	assert.Equal(t, http.StatusOK, rec2.Code)

	// Verify it's valid JSON
	var result map[string]interface{}
	err := json.NewDecoder(rec2.Body).Decode(&result)
	assert.NoError(t, err)
}
