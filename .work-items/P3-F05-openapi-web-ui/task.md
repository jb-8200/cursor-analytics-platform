# P3-F05: OpenAPI Web UI - Task Breakdown

## Status: NOT STARTED

## Progress Tracker

| Task | Description | Status | Time Est |
|------|-------------|--------|----------|
| TASK01 | Session Manager | NOT STARTED | 1h |
| TASK02 | Login Handler | NOT STARTED | 1h |
| TASK03 | Login Template | NOT STARTED | 0.5h |
| TASK04 | Docs Handler | NOT STARTED | 1h |
| TASK05 | Static File Embedding | NOT STARTED | 1h |
| TASK06 | Scalar Template | NOT STARTED | 0.5h |
| TASK07 | Router Integration | NOT STARTED | 1h |
| TASK08 | Dockerfile Update | NOT STARTED | 0.5h |
| TASK09 | E2E Tests | NOT STARTED | 1h |
| TASK10 | SPEC.md Update | NOT STARTED | 0.5h |

**Total Estimated: 8-10 hours**

---

## TASK01: Session Manager

**Goal**: Implement in-memory session management

**Files**:
- CREATE: `internal/api/docs/session.go`
- CREATE: `internal/api/docs/session_test.go`

**TDD Steps**:
1. RED: Write test for session creation
2. GREEN: Implement CreateSession
3. RED: Write test for session validation
4. GREEN: Implement ValidateSession
5. RED: Write test for session deletion
6. GREEN: Implement DeleteSession
7. RED: Write test for session expiration
8. GREEN: Implement expiration check

**Deliverables**:
- [ ] SessionManager struct with Create/Validate/Delete
- [ ] 8-hour session expiry
- [ ] Thread-safe with sync.RWMutex
- [ ] Unit tests with >80% coverage

---

## TASK02: Login Handler

**Goal**: Implement login page and authentication

**Files**:
- CREATE: `internal/api/docs/handler.go`
- CREATE: `internal/api/docs/handler_test.go`

**TDD Steps**:
1. RED: Test login form rendering (GET /docs/login)
2. GREEN: Implement ServeLogin GET handler
3. RED: Test successful authentication
4. GREEN: Implement ServeLogin POST handler
5. RED: Test invalid credentials
6. GREEN: Add error handling
7. RED: Test logout clears session
8. GREEN: Implement ServeLogout handler

**Deliverables**:
- [ ] GET /docs/login renders form
- [ ] POST /docs/login validates credentials
- [ ] GET /docs/logout clears session
- [ ] Cookie setting/clearing
- [ ] Unit tests

---

## TASK03: Login Template

**Goal**: Create styled login form

**Files**:
- CREATE: `internal/api/docs/templates/login.html`

**Deliverables**:
- [ ] Dark theme matching Scalar
- [ ] Username/password form fields
- [ ] Error message display area
- [ ] Submit button
- [ ] Responsive layout

---

## TASK04: Docs Handler

**Goal**: Implement session-protected documentation page

**Files**:
- UPDATE: `internal/api/docs/handler.go`
- UPDATE: `internal/api/docs/handler_test.go`

**TDD Steps**:
1. RED: Test unauthenticated redirect
2. GREEN: Implement session check middleware
3. RED: Test authenticated access
4. GREEN: Implement docs page serving

**Deliverables**:
- [ ] Session validation middleware
- [ ] Redirect to login if no session
- [ ] Render Scalar page when authenticated
- [ ] Unit tests

---

## TASK05: Static File Embedding

**Goal**: Embed Scalar assets and OpenAPI specs

**Files**:
- CREATE: `internal/docs/static/openapi/cursor-api.yaml` (copy)
- CREATE: `internal/docs/static/openapi/github-sim-api.yaml` (copy)
- CREATE: `internal/docs/embed.go`

**TDD Steps**:
1. RED: Test OpenAPI spec serving
2. GREEN: Implement embed directive
3. RED: Test static file handler
4. GREEN: Implement file serving

**Deliverables**:
- [ ] `//go:embed` directive for static files
- [ ] OpenAPI specs copied from /specs/openapi/
- [ ] Static file handler
- [ ] Unit tests

---

## TASK06: Scalar Template

**Goal**: Create Scalar API reference page

**Files**:
- CREATE: `internal/api/docs/templates/docs.html`

**Deliverables**:
- [ ] Scalar CDN integration
- [ ] Spec selector (Cursor API / GitHub API)
- [ ] Dark theme configuration
- [ ] Responsive layout

---

## TASK07: Router Integration

**Goal**: Add docs routes to main router

**Files**:
- UPDATE: `internal/server/router.go`
- UPDATE: `internal/server/router_test.go`

**TDD Steps**:
1. RED: Integration test for /docs route
2. GREEN: Register docs handlers
3. RED: Test /docs exempt from API auth
4. GREEN: Update authProtectedRoutes

**Deliverables**:
- [ ] Docs routes registered
- [ ] /docs excluded from API key auth
- [ ] All existing tests pass
- [ ] Integration tests

---

## TASK08: Dockerfile Update

**Goal**: Ensure Docker build includes embedded files

**Files**:
- UPDATE: `Dockerfile`

**Deliverables**:
- [ ] Copy OpenAPI specs before go build
- [ ] Verify binary includes embedded files
- [ ] Test Docker build succeeds
- [ ] Test Docker run serves /docs

---

## TASK09: E2E Tests

**Goal**: Comprehensive end-to-end testing

**Files**:
- CREATE: `test/e2e/docs_test.go`

**Test Cases**:
- [ ] Unauthenticated access redirects to login
- [ ] Valid credentials grant access
- [ ] Invalid credentials show error
- [ ] Logout clears session
- [ ] OpenAPI specs served correctly
- [ ] Scalar UI loads
- [ ] Session expiration works

---

## TASK10: SPEC.md Update

**Goal**: Document the new feature

**Files**:
- UPDATE: `services/cursor-sim/SPEC.md`

**Deliverables**:
- [ ] Add "Documentation UI" section
- [ ] Document /docs routes
- [ ] Document authentication flow
- [ ] Update "Last Updated" date

---

## Implementation Notes

### Credentials (Hardcoded)
- Username: `dox`
- Password: `dox-a3`

### Session Cookie
- Name: `docs_session`
- HttpOnly: true
- SameSite: Lax
- MaxAge: 8 hours

### Dependencies
- No external Go dependencies required
- Scalar loaded via CDN
