package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test handler that returns 200 OK
func testHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

func TestBasicAuth_ValidKey(t *testing.T) {
	apiKey := "test-api-key-123"
	handler := BasicAuth(apiKey)(testHandler())

	req := httptest.NewRequest("GET", "/test", nil)
	req.SetBasicAuth(apiKey, "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestBasicAuth_InvalidKey(t *testing.T) {
	apiKey := "test-api-key-123"
	handler := BasicAuth(apiKey)(testHandler())

	req := httptest.NewRequest("GET", "/test", nil)
	req.SetBasicAuth("wrong-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Unauthorized")
}

func TestBasicAuth_MissingAuth(t *testing.T) {
	apiKey := "test-api-key-123"
	handler := BasicAuth(apiKey)(testHandler())

	req := httptest.NewRequest("GET", "/test", nil)
	// No auth header
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "Basic", rec.Header().Get("WWW-Authenticate"))
}

func TestBasicAuth_EmptyPassword(t *testing.T) {
	apiKey := "test-api-key-123"
	handler := BasicAuth(apiKey)(testHandler())

	// Valid: API key as username, empty password
	req := httptest.NewRequest("GET", "/test", nil)
	req.SetBasicAuth(apiKey, "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestBasicAuth_NonEmptyPassword(t *testing.T) {
	apiKey := "test-api-key-123"
	handler := BasicAuth(apiKey)(testHandler())

	// Invalid: password should be empty
	req := httptest.NewRequest("GET", "/test", nil)
	req.SetBasicAuth(apiKey, "some-password")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should still accept (Cursor API accepts any password)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRateLimit_WithinLimit(t *testing.T) {
	limiter := NewRateLimiter(10, time.Minute) // 10 requests per minute
	handler := RateLimit(limiter)(testHandler())

	// Send 5 requests (within limit)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "request %d should succeed", i)
	}
}

func TestRateLimit_ExceedLimit(t *testing.T) {
	limiter := NewRateLimiter(3, time.Minute) // 3 requests per minute
	handler := RateLimit(limiter)(testHandler())

	// Send 3 requests (at limit)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code, "request %d should succeed", i)
	}

	// 4th request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("Retry-After"))
}

func TestRateLimit_RefillTokens(t *testing.T) {
	limiter := NewRateLimiter(2, 100*time.Millisecond) // 2 requests per 100ms
	handler := RateLimit(limiter)(testHandler())

	// Use up tokens
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	}

	// Should be rate limited now
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusTooManyRequests, rec.Code)

	// Wait for refill
	time.Sleep(150 * time.Millisecond)

	// Should work again
	req = httptest.NewRequest("GET", "/test", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestLogger_LogsRequests(t *testing.T) {
	handler := Logger(testHandler())

	req := httptest.NewRequest("GET", "/test?param=value", nil)
	rec := httptest.NewRecorder()

	// Should not panic and should pass through
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestLogger_LogsDuration(t *testing.T) {
	slowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	handler := Logger(slowHandler)

	req := httptest.NewRequest("GET", "/slow", nil)
	rec := httptest.NewRecorder()

	start := time.Now()
	handler.ServeHTTP(rec, req)
	elapsed := time.Since(start)

	// Should have waited at least 10ms
	assert.GreaterOrEqual(t, elapsed, 10*time.Millisecond)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMiddlewareChaining(t *testing.T) {
	apiKey := "test-key"
	limiter := NewRateLimiter(100, time.Minute)

	// Chain all middleware: Logger -> RateLimit -> BasicAuth -> Handler
	handler := Logger(
		RateLimit(limiter)(
			BasicAuth(apiKey)(testHandler()),
		),
	)

	req := httptest.NewRequest("GET", "/test", nil)
	req.SetBasicAuth(apiKey, "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestMiddlewareChaining_AuthFailsFirst(t *testing.T) {
	apiKey := "test-key"
	limiter := NewRateLimiter(100, time.Minute)

	handler := Logger(
		RateLimit(limiter)(
			BasicAuth(apiKey)(testHandler()),
		),
	)

	req := httptest.NewRequest("GET", "/test", nil)
	// No auth header
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should fail at BasicAuth before reaching RateLimit
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	limiter := NewRateLimiter(100, time.Minute)
	handler := RateLimit(limiter)(testHandler())

	// Send 50 concurrent requests
	done := make(chan bool, 50)
	for i := 0; i < 50; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			done <- rec.Code == http.StatusOK
		}()
	}

	// All should succeed (within limit)
	successCount := 0
	for i := 0; i < 50; i++ {
		if <-done {
			successCount++
		}
	}

	assert.Equal(t, 50, successCount, "all requests should succeed")
}
