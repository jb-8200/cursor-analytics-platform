package api

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// BasicAuth creates middleware that validates API key via HTTP Basic Auth.
// Username should be the API key, password should be empty (as per Cursor API spec).
func BasicAuth(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, _, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", "Basic")
				RespondError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			if username != apiKey {
				RespondError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimiter implements a token bucket rate limiter.
type RateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

// NewRateLimiter creates a rate limiter with the given capacity and window.
// Example: NewRateLimiter(100, time.Minute) allows 100 requests per minute.
func NewRateLimiter(capacity int, window time.Duration) *RateLimiter {
	refillRate := float64(capacity) / window.Seconds()
	return &RateLimiter{
		tokens:     float64(capacity),
		maxTokens:  float64(capacity),
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request can proceed and consumes a token if so.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()

	// Refill tokens based on time elapsed
	rl.tokens = min(rl.maxTokens, rl.tokens+elapsed*rl.refillRate)
	rl.lastRefill = now

	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}

	return false
}

// RetryAfter returns the number of seconds to wait before retrying.
func (rl *RateLimiter) RetryAfter() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	tokensNeeded := 1.0 - rl.tokens
	if tokensNeeded <= 0 {
		return 0
	}

	seconds := tokensNeeded / rl.refillRate
	return int(seconds) + 1 // Round up
}

// RateLimit creates middleware that enforces rate limiting.
func RateLimit(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				retryAfter := limiter.RetryAfter()
				w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
				RespondError(w, http.StatusTooManyRequests, "Rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Logger creates middleware that logs HTTP requests.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(lw, r)

		duration := time.Since(start)
		log.Printf(
			"%s %s %d %s",
			r.Method,
			r.URL.Path,
			lw.statusCode,
			duration,
		)
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

// min returns the minimum of two float64 values.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
