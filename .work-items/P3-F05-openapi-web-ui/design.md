# P3-F05: OpenAPI Web UI - Technical Design

## Overview

Add an interactive web-based API documentation interface to cursor-sim Docker container, accessible at `/docs` with session-based authentication.

## Tool Selection: Scalar API Reference

| Tool | Size | Interactive Testing | Decision |
|------|------|---------------------|----------|
| Swagger UI | ~3MB | Yes | Too heavy |
| Redoc | ~1MB | No | No "Try It" |
| **Scalar** | ~500KB | Yes | **Selected** |

**Justification**: Scalar provides modern UI, interactive API testing, minimal footprint (~6% Docker image increase), and no build step required.

## Architecture

### Route Structure

```
/docs           -> Scalar API Reference (session-protected)
/docs/login     -> Login page (GET: form, POST: authenticate)
/docs/logout    -> Clear session
/docs/openapi/* -> Serve OpenAPI YAML specs
```

### Authentication Flow

```
1. GET /docs -> Check session cookie
2. No session -> Redirect to /docs/login
3. POST /docs/login (user=dox, pass=dox-a3) -> Set session cookie
4. Session valid 8 hours -> Access /docs
```

### Session Management

```go
type Session struct {
    ID        string
    Username  string
    CreatedAt time.Time
    ExpiresAt time.Time
}

type SessionManager struct {
    sessions map[string]*Session
    mu       sync.RWMutex
}
```

- In-memory storage (sufficient for simulator)
- 8-hour session expiry
- Secure cookie with HttpOnly flag

## File Structure

```
services/cursor-sim/
├── internal/
│   └── api/
│       └── docs/                    # NEW
│           ├── handler.go           # Docs router, static serving
│           ├── handler_test.go
│           ├── session.go           # Session management
│           ├── session_test.go
│           └── templates/
│               ├── login.html       # Login form
│               └── docs.html        # Scalar page
│   └── docs/
│       └── static/                  # NEW - embedded assets
│           └── openapi/
│               ├── cursor-api.yaml
│               └── github-sim-api.yaml
```

## Router Integration

Modify `internal/server/router.go`:

```go
// Documentation UI (P3-F05) - with separate session auth
docsHandler := docs.NewHandler()
mux.Handle("/docs", docsHandler.Index())
mux.Handle("/docs/", docsHandler.Static())
```

Update `authProtectedRoutes` to skip `/docs` paths (like `/health`):

```go
func authProtectedRoutes(handler http.Handler, apiKey string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip auth for health and docs endpoints
        if r.URL.Path == "/health" || strings.HasPrefix(r.URL.Path, "/docs") {
            handler.ServeHTTP(w, r)
            return
        }
        api.BasicAuth(apiKey)(handler).ServeHTTP(w, r)
    })
}
```

## Static File Embedding

Use Go 1.16+ embed directive:

```go
//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS
```

Benefits:
- Single binary deployment
- No external file dependencies
- Works with existing multi-stage Docker build

## Dockerfile Changes

```dockerfile
# In builder stage, copy OpenAPI specs before build
COPY specs/openapi/*.yaml /src/internal/docs/static/openapi/

# Build includes embedded files automatically
RUN go build -ldflags="-s -w" -o /out/cursor-sim ./cmd/simulator
```

## Security Considerations

1. **Session Security**: HttpOnly cookie, secure flag in production
2. **Credential Storage**: Hardcoded for simplicity (simulator use case)
3. **CSRF Protection**: Not required (no state-changing operations)
4. **Rate Limiting**: Inherits existing rate limiter

## Alternatives Considered

### Option A: Basic Auth for /docs
- Pro: No session management
- Con: Poor UX, browser caches credentials

### Option B: CDN-only Scalar
- Pro: Zero Docker size impact
- Con: Requires internet, no offline support

### Option C: Separate docs container
- Pro: Full separation of concerns
- Con: Complex deployment, not single binary

**Decision**: Session-based auth with embedded Scalar provides best balance.
