---
name: cursor-api-patterns
description: Cursor Business API implementation patterns. Use when implementing API endpoints, writing HTTP handlers, handling authentication, pagination, error responses, or CSV exports. Covers response formats for AI Code Tracking, Analytics, and Admin APIs. (project)
---

# Cursor API Patterns

This skill provides expertise on implementing Cursor Business API-compatible endpoints.

## Documentation Source of Truth

**Always reference these files for accurate API specifications:**

| Document | Path | Description |
|----------|------|-------------|
| Overview | `docs/api-reference/cursor_overview.md` | Authentication, rate limits, caching |
| Admin API | `docs/api-reference/cursor_admin.md` | Team management, usage data |
| Analytics API | `docs/api-reference/cursor_analytics.md` | Team metrics, DAU, model usage |
| AI Code Tracking | `docs/api-reference/cursor_codetrack.md` | Per-commit metrics (Enterprise) |
| Cloud Agents | `docs/api-reference/cursor_agents.md` | Agent management |

## Response Format Standards

### AI Code Tracking API (Enterprise)

Endpoints: `/analytics/ai-code/commits`, `/analytics/ai-code/changes`

```json
{
  "items": [...],
  "totalCount": 42,
  "page": 1,
  "pageSize": 100
}
```

### Analytics API (Team-Level)

Endpoints: `/analytics/team/*`

```json
{
  "data": [...],
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
  "data": { "alice@example.com": [...] },
  "pagination": {
    "page": 1,
    "pageSize": 100,
    "totalUsers": 250,
    "hasNextPage": true
  },
  "params": {...}
}
```

### Admin API - Team Members

Endpoint: `GET /teams/members`

```json
{
  "teamMembers": [
    { "name": "Alex", "email": "dev@company.com", "role": "member" }
  ]
}
```

### Error Response Structure

```json
{
  "error": "Bad Request",
  "message": "Human-readable error message"
}
```

## Authentication Pattern

All Cursor API endpoints use Basic Authentication:

```
Authorization: Basic {base64(api_key:)}
```

Note: Password is empty. The API key is used as the username.

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

## Query Parameters

### AI Code Tracking API

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `startDate` | string | 7d ago | ISO date or relative (7d, 30d) |
| `endDate` | string | now | ISO date or relative |
| `user` | string | - | Filter by email or user_id |
| `page` | int | 1 | Page number (1-indexed) |
| `pageSize` | int | 100 | Max: 1000 |

### Analytics API

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `startDate` | string | 7d ago | YYYY-MM-DD or shortcuts |
| `endDate` | string | today | YYYY-MM-DD or shortcuts |
| `users` | string | - | Comma-separated emails |
| `page` | int | 1 | For by-user endpoints |
| `pageSize` | int | 100 | Max: 500 |

## Rate Limits

| API | Limit |
|-----|-------|
| Admin API | 20 req/min |
| Analytics Team | 100 req/min |
| Analytics By-user | 50 req/min |
| AI Code Tracking | 20 req/min |

## CSV Export Pattern

For endpoints ending in `.csv`:

```
Content-Type: text/csv; charset=utf-8
Content-Disposition: attachment; filename="commits-2025-01-15.csv"
```

- Use snake_case for CSV headers
- JSON uses camelCase
- Quote fields containing commas

## HTTP Status Codes

| Code | When to Use |
|------|-------------|
| 200 | Successful request |
| 400 | Invalid parameters |
| 401 | Authentication failed |
| 403 | Enterprise feature on wrong plan |
| 404 | Resource not found |
| 429 | Rate limit exceeded |
| 500 | Unexpected server error |

## API Response Summary

| API | Endpoint Pattern | Response Wrapper |
|-----|------------------|------------------|
| AI Code Tracking | `/analytics/ai-code/*` | `{ items, totalCount, page, pageSize }` |
| Analytics Team | `/analytics/team/*` | `{ data, params }` |
| Analytics By-User | `/analytics/by-user/*` | `{ data, pagination, params }` |
| Admin Members | `/teams/members` | `{ teamMembers }` |

## Common Mistakes to Avoid

1. **Don't confuse response formats** - AI Code Tracking uses `items`, Analytics uses `data`
2. **Don't use wrong parameter names** - Use `startDate`/`endDate`, not `from`/`to`
3. **Don't return arrays directly** - Always use proper envelope structure
4. **Don't skip pagination metadata** - Include even if only one page
5. **Don't use 200 for errors** - Use proper HTTP status codes
