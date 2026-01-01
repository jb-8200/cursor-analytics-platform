# Cursor API Patterns Skill

This skill provides expertise on implementing Cursor Business API-compatible endpoints.

## Overview

When implementing endpoints that mimic Cursor's actual APIs, follow these patterns exactly to ensure compatibility.

---

## Response Format Standard

### Success Response Structure

All successful responses MUST use this envelope:

```json
{
  "data": [...],           // Array or object with actual results
  "pagination": {          // Present for paginated endpoints
    "page": 1,
    "pageSize": 100,
    "total": 5847,
    "hasMore": true
  },
  "params": {              // Echo of request parameters
    "startDate": "2026-01-08",
    "endDate": "2026-01-15",
    "user": null
  }
}
```

### Error Response Structure

All error responses MUST use:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {}          // Optional additional context
  }
}
```

---

## Authentication Pattern

### Basic Authentication

All Cursor API endpoints use Basic Authentication:

```
Authorization: Basic {base64(api_key:api_secret)}
```

**Implementation Rules:**
1. Validate credentials on EVERY request
2. Return `401 Unauthorized` if invalid
3. Use constant-time comparison for credentials
4. Error response for 401:
   ```json
   {
     "error": {
       "code": "UNAUTHORIZED",
       "message": "Invalid API credentials"
     }
   }
   ```

**Go Implementation:**
```go
func BasicAuthMiddleware(expectedKey, expectedSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            username, password, ok := r.BasicAuth()
            if !ok || username != expectedKey || password != expectedSecret {
                w.Header().Set("WWW-Authenticate", `Basic realm="Cursor API"`)
                w.WriteHeader(http.StatusUnauthorized)
                json.NewEncoder(w).Encode(map[string]interface{}{
                    "error": map[string]string{
                        "code":    "UNAUTHORIZED",
                        "message": "Invalid API credentials",
                    },
                })
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## Date Parameter Handling

### Supported Formats

Cursor API accepts multiple date formats:

1. **ISO 8601**: `2026-01-15T10:30:00Z`
2. **Date only**: `2026-01-15`
3. **Relative shortcuts**:
   - `7d` - 7 days ago
   - `30d` - 30 days ago
   - `90d` - 90 days ago
   - `today` - Start of current UTC day
   - `now` - Current timestamp

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

    if param == "today" {
        now := time.Now().UTC()
        return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
    }

    if param == "now" {
        return time.Now().UTC(), nil
    }

    // Try ISO 8601
    t, err := time.Parse(time.RFC3339, param)
    if err == nil {
        return t, nil
    }

    // Try date-only format
    return time.Parse("2006-01-02", param)
}
```

---

## Pagination Pattern

### Query Parameters

| Parameter | Type | Default | Max |
|-----------|------|---------|-----|
| `page` | int | 1 | - |
| `pageSize` | int | 100 | 1000 |

### Response

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "pageSize": 100,
    "total": 5847,
    "hasMore": true
  }
}
```

**Calculation:**
```go
type PaginationInfo struct {
    Page     int  `json:"page"`
    PageSize int  `json:"pageSize"`
    Total    int  `json:"total"`
    HasMore  bool `json:"hasMore"`
}

func CalculatePagination(page, pageSize, totalCount int) PaginationInfo {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 1000 {
        pageSize = 100
    }

    offset := (page - 1) * pageSize
    hasMore := offset+pageSize < totalCount

    return PaginationInfo{
        Page:     page,
        PageSize: pageSize,
        Total:    totalCount,
        HasMore:  hasMore,
    }
}
```

---

## Rate Limiting

### Cursor API Limits

| Endpoint Type | Limit |
|---------------|-------|
| `/analytics/team/*` | 100 requests/minute |
| `/analytics/by-user/*` | 50 requests/minute |

### Rate Limit Response

**HTTP Status**: 429 Too Many Requests

```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Retry after 60 seconds.",
    "retry_after": 60
  }
}
```

**Headers to Include:**
```
Retry-After: 60
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1705329600
```

---

## Endpoint Patterns

### AI Code Tracking Pattern

**Endpoint**: `GET /v1/analytics/ai-code/commits`

**Query Parameters:**
- `startDate` (string)
- `endDate` (string)
- `user` (string, optional)
- `page` (int)
- `pageSize` (int)

**Response:**
```json
{
  "data": [
    {
      "commit_hash": "abc123...",
      "timestamp": "2026-01-15T10:30:00Z",
      "user_id": "user_abc123",
      "user_email": "jane@example.com",
      "repository": "frontend-app",
      "branch": "main",
      "total_lines": 145,
      "lines_from_tab": 87,
      "lines_from_composer": 35,
      "lines_non_ai": 23,
      "ingestion_time": "2026-01-15T10:31:00Z"
    }
  ],
  "pagination": {...},
  "params": {...}
}
```

### Team Analytics Pattern

**Endpoint**: `GET /v1/analytics/team/agent-edits`

**Response:**
```json
{
  "data": [
    {
      "date": "2026-01-15",
      "total_edits": 1247,
      "edits_from_tab": 873,
      "edits_from_composer": 374,
      "unique_users": 45
    }
  ],
  "pagination": {...}
}
```

---

## CSV Export Pattern

For endpoints ending in `.csv`:

**Headers:**
```
Content-Type: text/csv
Content-Disposition: attachment; filename="commits-2026-01-15.csv"
```

**Format:**
```csv
commit_hash,timestamp,user_email,repository,total_lines,lines_from_tab,lines_from_composer,lines_non_ai
abc123,2026-01-15T10:30:00Z,jane@example.com,frontend-app,145,87,35,23
```

**Rules:**
1. First row is header
2. No spaces after commas
3. Use RFC3339 for timestamps
4. Quote fields containing commas

---

## Testing Requirements

### Response Schema Validation

Every endpoint test MUST verify:

```go
func TestEndpointResponseSchema(t *testing.T) {
    // 1. Status code
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // 2. Content-Type
    assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

    // 3. Response structure
    var response struct {
        Data       interface{} `json:"data"`
        Pagination interface{} `json:"pagination,omitempty"`
        Params     interface{} `json:"params,omitempty"`
    }
    err := json.NewDecoder(resp.Body).Decode(&response)
    assert.NoError(t, err)
    assert.NotNil(t, response.Data)
}
```

### Auth Test

```go
func TestUnauthorizedAccess(t *testing.T) {
    req, _ := http.NewRequest("GET", "/v1/analytics/team/dau", nil)
    // No auth header

    resp := executeRequest(req)
    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

    var errResp ErrorResponse
    json.NewDecoder(resp.Body).Decode(&errResp)
    assert.Equal(t, "UNAUTHORIZED", errResp.Error.Code)
}
```

---

## Common Mistakes to Avoid

1. **Don't return arrays directly** - Always wrap in `{ "data": [...] }`
2. **Don't skip pagination envelope** - Include even if only one page
3. **Don't use 200 for errors** - Use proper HTTP status codes
4. **Don't ignore date format validation** - Return 400 for invalid dates
5. **Don't hardcode credentials** - Always use config/environment
6. **Don't skip rate limiting** - Implement even in simulator
7. **Don't forget CORS headers** - Include `Access-Control-Allow-Origin`

---

## Quick Reference

### HTTP Status Codes

| Code | When to Use |
|------|-------------|
| 200 | Successful request |
| 400 | Invalid parameters (bad date, invalid filter) |
| 401 | Authentication failed |
| 404 | Resource not found |
| 429 | Rate limit exceeded |
| 500 | Unexpected server error |

### Error Codes

| Code | HTTP Status | Meaning |
|------|-------------|---------|
| `UNAUTHORIZED` | 401 | Invalid API credentials |
| `INVALID_DATE_RANGE` | 400 | Start date after end date |
| `INVALID_PARAMETER` | 400 | Invalid query parameter |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |

---

## When to Use This Skill

Invoke this skill when:
- Implementing any `/v1/analytics/*` endpoint
- Writing API response handlers
- Debugging API response format issues
- Writing integration tests for API endpoints
- Reviewing API implementation code
