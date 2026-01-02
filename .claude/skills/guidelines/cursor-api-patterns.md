# Cursor API Patterns Skill

This skill provides expertise on implementing Cursor Business API-compatible endpoints.

## Documentation Source of Truth

**Always reference these files for accurate API specifications:**

| Document | Path | Description |
|----------|------|-------------|
| Overview | `docs/api-reference/cursor_overview.md` | Authentication, rate limits, caching, best practices |
| Admin API | `docs/api-reference/cursor_admin.md` | Team management, usage data, spending |
| Analytics API | `docs/api-reference/cursor_analytics.md` | Team metrics, DAU, model usage, leaderboards |
| AI Code Tracking | `docs/api-reference/cursor_codetrack.md` | Per-commit metrics, code changes (Enterprise) |
| Cloud Agents | `docs/api-reference/cursor_agents.md` | Programmatic agent management |

**Before implementing any endpoint, read the relevant documentation file first.**

---

## Response Format Standards

### AI Code Tracking API (Enterprise)

Endpoints: `/analytics/ai-code/commits`, `/analytics/ai-code/changes`

```json
{
  "items": [
    {
      "commitHash": "a1b2c3d4",
      "userId": "user_3k9x8q...",
      "userEmail": "developer@company.com",
      "repoName": "company/repo",
      "totalLinesAdded": 120,
      "tabLinesAdded": 50,
      "composerLinesAdded": 40,
      "nonAiLinesAdded": 30,
      "commitTs": "2025-07-30T14:12:03.000Z",
      "createdAt": "2025-07-30T14:12:30.000Z"
    }
  ],
  "totalCount": 42,
  "page": 1,
  "pageSize": 100
}
```

### Analytics API (Team-Level)

Endpoints: `/analytics/team/*`

```json
{
  "data": [
    {
      "event_date": "2025-01-15",
      "total_suggested_diffs": 145,
      "total_accepted_diffs": 98
    }
  ],
  "params": {
    "metric": "agent-edits",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31"
  }
}
```

### Analytics API (By-User)

Endpoints: `/analytics/by-user/*`

```json
{
  "data": {
    "alice@example.com": [
      { "event_date": "2025-01-15", "suggested_lines": 125 }
    ],
    "bob@example.com": [
      { "event_date": "2025-01-15", "suggested_lines": 95 }
    ]
  },
  "pagination": {
    "page": 1,
    "pageSize": 100,
    "totalUsers": 250,
    "totalPages": 3,
    "hasNextPage": true,
    "hasPreviousPage": false
  },
  "params": {
    "metric": "agent-edits",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31",
    "userMappings": [
      { "id": "user_abc123", "email": "alice@example.com" }
    ]
  }
}
```

### Admin API - Team Members

Endpoint: `GET /teams/members`

```json
{
  "teamMembers": [
    {
      "name": "Alex",
      "email": "developer@company.com",
      "role": "member"
    }
  ]
}
```

### Error Response Structure

All error responses MUST use:

```json
{
  "error": "Bad Request",
  "message": "Human-readable error message"
}
```

---

## Authentication Pattern

### Basic Authentication

All Cursor API endpoints use Basic Authentication:

```
Authorization: Basic {base64(api_key:)}
```

Note: Password is empty. The API key is used as the username.

**Implementation Rules:**
1. Validate credentials on EVERY request
2. Return `401 Unauthorized` if invalid
3. Use constant-time comparison for credentials

**Go Implementation:**
```go
func BasicAuthMiddleware(expectedKey string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            username, _, ok := r.BasicAuth()
            if !ok || username != expectedKey {
                w.Header().Set("WWW-Authenticate", `Basic realm="Cursor API"`)
                w.WriteHeader(http.StatusUnauthorized)
                json.NewEncoder(w).Encode(map[string]string{
                    "error":   "Unauthorized",
                    "message": "Invalid API key",
                })
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## Query Parameters by API

### AI Code Tracking API

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `startDate` | string | 7d ago | ISO date, "now", or relative (7d, 30d) |
| `endDate` | string | now | ISO date, "now", or relative (0d) |
| `user` | string | - | Filter by email, user_id, or numeric ID |
| `page` | int | 1 | Page number (1-indexed) |
| `pageSize` | int | 100 | Max: 1000 |

### Analytics API

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `startDate` | string | 7d ago | YYYY-MM-DD or shortcuts (7d, 30d, today) |
| `endDate` | string | today | YYYY-MM-DD or shortcuts |
| `users` | string | - | Comma-separated emails or user IDs |
| `page` | int | 1 | For by-user and leaderboard endpoints |
| `pageSize` | int | 100 | Max: 500 |

---

## Date Parameter Handling

### Supported Formats

1. **ISO 8601**: `2025-01-15T10:30:00Z`
2. **Date only**: `2025-01-15` (Recommended for caching)
3. **Relative shortcuts**:
   - `7d` - 7 days ago
   - `30d` - 30 days ago
   - `today` or `now` - Current date
   - `yesterday` - Previous day

**Parsing Logic:**
```go
func ParseDateParam(param string) (time.Time, error) {
    // Try relative shortcuts first
    if strings.HasSuffix(param, "d") {
        days, err := strconv.Atoi(strings.TrimSuffix(param, "d"))
        if err == nil {
            return time.Now().UTC().AddDate(0, 0, -days), nil
        }
    }

    switch param {
    case "today", "now":
        now := time.Now().UTC()
        return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
    case "yesterday":
        now := time.Now().UTC().AddDate(0, 0, -1)
        return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
    }

    // Try ISO 8601
    if t, err := time.Parse(time.RFC3339, param); err == nil {
        return t, nil
    }

    // Try date-only format
    return time.Parse("2006-01-02", param)
}
```

---

## Rate Limits

| API | Endpoint Type | Limit |
|-----|---------------|-------|
| Admin API | Most endpoints | 20 req/min |
| Admin API | `/teams/user-spend-limit` | 60 req/min |
| Analytics API | Team-level | 100 req/min |
| Analytics API | By-user | 50 req/min |
| AI Code Tracking | All endpoints | 20 req/min per endpoint |

### Rate Limit Response

**HTTP Status**: 429 Too Many Requests

```json
{
  "error": "Too Many Requests",
  "message": "Rate limit exceeded. Please try again later."
}
```

---

## CSV Export Pattern

For endpoints ending in `.csv`:

**Headers:**
```
Content-Type: text/csv; charset=utf-8
Content-Disposition: attachment; filename="commits-2025-01-15.csv"
```

**CSV Column Naming:**
- Use snake_case for CSV headers (e.g., `commit_hash`, `user_id`)
- JSON uses camelCase (e.g., `commitHash`, `userId`)

**Rules:**
1. First row is header
2. No spaces after commas
3. Use RFC3339 for timestamps
4. Quote fields containing commas

---

## Testing Requirements

### Response Schema Validation

Every endpoint test MUST verify against the documentation:

```go
func TestAICodeCommitsResponse(t *testing.T) {
    // Reference: docs/api-reference/cursor_codetrack.md

    // 1. Status code
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // 2. Content-Type
    assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

    // 3. Response structure per cursor_codetrack.md
    var response struct {
        Items      []interface{} `json:"items"`
        TotalCount int           `json:"totalCount"`
        Page       int           `json:"page"`
        PageSize   int           `json:"pageSize"`
    }
    err := json.NewDecoder(resp.Body).Decode(&response)
    assert.NoError(t, err)
    assert.NotNil(t, response.Items)
}
```

### Auth Test

```go
func TestUnauthorizedAccess(t *testing.T) {
    req, _ := http.NewRequest("GET", "/analytics/ai-code/commits", nil)
    // No auth header

    resp := executeRequest(req)
    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

    var errResp map[string]string
    json.NewDecoder(resp.Body).Decode(&errResp)
    assert.Equal(t, "Unauthorized", errResp["error"])
}
```

---

## Common Mistakes to Avoid

1. **Don't confuse response formats** - AI Code Tracking uses `items`, Analytics uses `data`
2. **Don't use wrong parameter names** - Use `startDate`/`endDate`/`user`, not `from`/`to`/`userId`
3. **Don't return arrays directly** - Always use the proper envelope structure
4. **Don't skip pagination metadata** - Include even if only one page
5. **Don't use 200 for errors** - Use proper HTTP status codes
6. **Don't hardcode credentials** - Always use config/environment

---

## Quick Reference

### HTTP Status Codes

| Code | When to Use |
|------|-------------|
| 200 | Successful request |
| 400 | Invalid parameters (bad date, invalid filter) |
| 401 | Authentication failed |
| 403 | Enterprise feature on non-Enterprise plan |
| 404 | Resource not found |
| 429 | Rate limit exceeded |
| 500 | Unexpected server error |

### API Response Summary

| API | Endpoint Pattern | Response Wrapper |
|-----|------------------|------------------|
| AI Code Tracking | `/analytics/ai-code/*` | `{ items, totalCount, page, pageSize }` |
| Analytics Team | `/analytics/team/*` | `{ data, params }` |
| Analytics By-User | `/analytics/by-user/*` | `{ data, pagination, params }` |
| Admin Members | `/teams/members` | `{ teamMembers }` |

---

## When to Use This Skill

Invoke this skill when:
- Implementing any `/analytics/*` or `/teams/*` endpoint
- Writing API response handlers
- Debugging API response format issues
- Writing integration tests for API endpoints
- Reviewing API implementation code

**Always cross-reference with the documentation files in `docs/api-reference/`**
