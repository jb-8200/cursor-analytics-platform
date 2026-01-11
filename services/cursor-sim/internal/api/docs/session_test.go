package docs

import (
	"testing"
	"time"
)

func TestCreateSession(t *testing.T) {
	sm := NewSessionManager()
	sessionID := sm.CreateSession("testuser")

	if sessionID == "" {
		t.Error("expected non-empty session ID")
	}
}

func TestValidateSession(t *testing.T) {
	sm := NewSessionManager()
	sessionID := sm.CreateSession("testuser")

	username, valid := sm.ValidateSession(sessionID)
	if !valid {
		t.Error("expected session to be valid")
	}
	if username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", username)
	}
}

func TestValidateSessionInvalid(t *testing.T) {
	sm := NewSessionManager()

	username, valid := sm.ValidateSession("invalid-id")
	if valid {
		t.Error("expected session to be invalid")
	}
	if username != "" {
		t.Errorf("expected empty username for invalid session, got '%s'", username)
	}
}

func TestDeleteSession(t *testing.T) {
	sm := NewSessionManager()
	sessionID := sm.CreateSession("testuser")

	sm.DeleteSession(sessionID)

	username, valid := sm.ValidateSession(sessionID)
	if valid {
		t.Error("expected session to be deleted")
	}
	if username != "" {
		t.Errorf("expected empty username after deletion, got '%s'", username)
	}
}

func TestSessionExpiration(t *testing.T) {
	sm := NewSessionManager()
	sessionID := sm.CreateSession("testuser")

	// Manually expire the session
	sm.mu.Lock()
	session := sm.sessions[sessionID]
	session.ExpiresAt = time.Now().Add(-1 * time.Second)
	sm.mu.Unlock()

	username, valid := sm.ValidateSession(sessionID)
	if valid {
		t.Error("expected expired session to be invalid")
	}
	if username != "" {
		t.Errorf("expected empty username for expired session, got '%s'", username)
	}
}

func TestSessionExpiryConcurrency(t *testing.T) {
	sm := NewSessionManager()
	sessionID := sm.CreateSession("testuser")

	// Ensure 8-hour expiry
	sm.mu.RLock()
	session := sm.sessions[sessionID]
	expiryTime := session.ExpiresAt
	createdTime := session.CreatedAt
	sm.mu.RUnlock()

	expectedDuration := 8 * time.Hour
	actualDuration := expiryTime.Sub(createdTime)

	if actualDuration < expectedDuration-time.Second || actualDuration > expectedDuration+time.Second {
		t.Errorf("expected session duration ~%v, got %v", expectedDuration, actualDuration)
	}
}

func TestConcurrentSessionOperations(t *testing.T) {
	sm := NewSessionManager()

	// Create multiple sessions concurrently
	sessionIDs := make(chan string, 100)
	for i := 0; i < 100; i++ {
		go func(id int) {
			sid := sm.CreateSession("user" + string(rune(id)))
			sessionIDs <- sid
		}(i)
	}

	// Validate sessions concurrently
	errCount := 0
	for i := 0; i < 100; i++ {
		sessionID := <-sessionIDs
		_, valid := sm.ValidateSession(sessionID)
		if !valid {
			errCount++
		}
	}

	if errCount > 0 {
		t.Errorf("expected all sessions valid, got %d invalid", errCount)
	}
}
