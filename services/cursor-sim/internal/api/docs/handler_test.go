package docs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeLoginFormGET(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest("GET", "/docs/login", nil)
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "login") && !strings.Contains(body, "Login") {
		t.Error("expected login form in response")
	}
}

func TestAuthenticateValidCredentials(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest("POST", "/docs/login", strings.NewReader("username=dox&password=dox-a3"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	// Should set a cookie for authenticated session
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Error("expected session cookie to be set")
	}

	// Find docs_session cookie
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "docs_session" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected docs_session cookie")
	}
}

func TestAuthenticateInvalidCredentials(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest("POST", "/docs/login", strings.NewReader("username=wrong&password=wrong"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Invalid") && !strings.Contains(body, "invalid") {
		t.Error("expected error message in response")
	}
}

func TestLogout(t *testing.T) {
	handler := NewHandler()

	// Create a session
	sessionID := handler.sessions.CreateSession("testuser")

	// Create logout request with session cookie
	req := httptest.NewRequest("GET", "/docs/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "docs_session",
		Value: sessionID,
	})
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	// Should redirect
	if w.Code != http.StatusFound && w.Code != http.StatusSeeOther {
		t.Errorf("expected redirect status, got %d", w.Code)
	}

	// Session should be deleted
	_, valid := handler.sessions.ValidateSession(sessionID)
	if valid {
		t.Error("expected session to be deleted after logout")
	}
}

func TestDocsWithoutSession(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest("GET", "/docs", nil)
	w := httptest.NewRecorder()

	handler.Index(w, req)

	// Should redirect to login
	if w.Code != http.StatusFound {
		t.Errorf("expected redirect, got status %d", w.Code)
	}
}

func TestDocsWithValidSession(t *testing.T) {
	handler := NewHandler()

	// Create a session
	sessionID := handler.sessions.CreateSession("testuser")

	// Request /docs with valid session
	req := httptest.NewRequest("GET", "/docs", nil)
	req.AddCookie(&http.Cookie{
		Name:  "docs_session",
		Value: sessionID,
	})
	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Scalar") && !strings.Contains(body, "scalar") {
		t.Error("expected Scalar content in response")
	}
}

func TestDocsWithExpiredSession(t *testing.T) {
	handler := NewHandler()

	// Create a session and manually expire it
	sessionID := handler.sessions.CreateSession("testuser")
	handler.sessions.mu.Lock()
	session := handler.sessions.sessions[sessionID]
	session.ExpiresAt = session.CreatedAt // Expired
	handler.sessions.mu.Unlock()

	// Request /docs with expired session
	req := httptest.NewRequest("GET", "/docs", nil)
	req.AddCookie(&http.Cookie{
		Name:  "docs_session",
		Value: sessionID,
	})
	w := httptest.NewRecorder()

	handler.Index(w, req)

	// Should redirect to login
	if w.Code != http.StatusFound {
		t.Errorf("expected redirect, got status %d", w.Code)
	}
}

func TestServeOpenAPISpec(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest("GET", "/docs/openapi/cursor-api.yaml", nil)
	w := httptest.NewRecorder()

	handler.Static(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/yaml" {
		t.Errorf("expected content type 'application/yaml', got '%s'", contentType)
	}

	// Check that content is served
	body := w.Body.String()
	if len(body) == 0 {
		t.Error("expected non-empty response body")
	}
}

func TestServeNotFound(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest("GET", "/docs/nonexistent.yaml", nil)
	w := httptest.NewRecorder()

	handler.Static(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
