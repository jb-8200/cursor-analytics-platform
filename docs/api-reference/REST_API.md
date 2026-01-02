# REST API Reference: cursor-sim

This document provides the REST API specification for the cursor-sim service. The simulator exposes endpoints that mimic the Cursor Business Activity API, enabling the aggregator to fetch developer activity data without requiring access to production Cursor credentials.

## Source of Truth

**For accurate API specifications, always reference the Cursor API documentation:**

| Document | Description |
|----------|-------------|
| [cursor_overview.md](cursor_overview.md) | Authentication, rate limits, caching, error handling |
| [cursor_admin.md](cursor_admin.md) | Admin API - Team management, usage data, spending |
| [cursor_analytics.md](cursor_analytics.md) | Analytics API - Team metrics, DAU, model usage |
| [cursor_codetrack.md](cursor_codetrack.md) | AI Code Tracking API - Per-commit metrics (Enterprise) |
| [cursor_agents.md](cursor_agents.md) | Cloud Agents API - Programmatic agent management |

**Claude Code Integration:** When implementing endpoints, use the `.claude/skills/cursor-api-patterns.md` skill for quick reference patterns.

---

## Base URL

When running locally via Docker Compose, the simulator is accessible at `http://localhost:8080`. When services communicate within the Docker network, use `http://cursor-sim:8080`.

---

## Authentication

The simulator supports Basic Authentication to match the production Cursor API:

```bash
curl https://localhost:8080/analytics/ai-code/commits \
  -u YOUR_API_KEY:
```

Or with the Authorization header:

```bash
Authorization: Basic {base64_encode('YOUR_API_KEY:')}
```

Note: Password is empty. The API key is used as the username.

---

## Implemented Endpoints

### AI Code Tracking API (Enterprise)

These endpoints match the Cursor AI Code Tracking API specification in [cursor_codetrack.md](cursor_codetrack.md).

#### GET /analytics/ai-code/commits

Retrieve per-commit AI usage metrics with line attribution.

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `startDate` | string | 7d ago | ISO date, "now", or relative (7d, 30d) |
| `endDate` | string | now | ISO date, "now", or relative (0d) |
| `user` | string | - | Filter by email or user_id |
| `page` | int | 1 | Page number (1-indexed) |
| `pageSize` | int | 100 | Max: 1000 |

**Response:**

```json
{
  "items": [
    {
      "commitHash": "a1b2c3d4",
      "userId": "user_001",
      "userEmail": "developer@company.com",
      "repoName": "acme/api",
      "branchName": "main",
      "isPrimaryBranch": true,
      "totalLinesAdded": 120,
      "totalLinesDeleted": 30,
      "tabLinesAdded": 50,
      "tabLinesDeleted": 10,
      "composerLinesAdded": 40,
      "composerLinesDeleted": 5,
      "nonAiLinesAdded": 30,
      "nonAiLinesDeleted": 15,
      "message": "Refactor: extract analytics client",
      "commitTs": "2025-07-30T14:12:03.000Z",
      "createdAt": "2025-07-30T14:12:30.000Z"
    }
  ],
  "totalCount": 42,
  "page": 1,
  "pageSize": 100
}
```

#### GET /analytics/ai-code/commits.csv

Download commit metrics as CSV for large data extractions.

**Response Headers:**
- `Content-Type: text/csv; charset=utf-8`
- `Content-Disposition: attachment; filename="commits-YYYYMMDD.csv"`

#### GET /analytics/ai-code/changes

Retrieve granular accepted AI changes grouped by changeId.

**Response:**

```json
{
  "items": [
    {
      "changeId": "749356201",
      "userId": "user_001",
      "userEmail": "developer@company.com",
      "source": "COMPOSER",
      "model": "gpt-4o",
      "totalLinesAdded": 18,
      "totalLinesDeleted": 4,
      "createdAt": "2025-07-30T15:10:12.000Z",
      "metadata": [
        {
          "fileName": "src/analytics/report.ts",
          "fileExtension": "ts",
          "linesAdded": 12,
          "linesDeleted": 3
        }
      ]
    }
  ],
  "totalCount": 128,
  "page": 1,
  "pageSize": 100
}
```

#### GET /analytics/ai-code/changes.csv

Download change metrics as CSV.

---

### Admin API

These endpoints match the Cursor Admin API specification in [cursor_admin.md](cursor_admin.md).

#### GET /teams/members

Retrieve all team members.

**Response:**

```json
{
  "teamMembers": [
    {
      "name": "Alex Chen",
      "email": "alex@company.com",
      "role": "member"
    },
    {
      "name": "Sam Admin",
      "email": "sam@company.com",
      "role": "owner"
    }
  ]
}
```

---

### Health Check

#### GET /v1/health

Check simulator health status. This endpoint does not require authentication.

**Response:**

```json
{
  "status": "healthy",
  "mode": "runtime",
  "seed_loaded": true,
  "developers_count": 50,
  "commits_count": 1250,
  "uptime_seconds": 3600
}
```

---

## Error Responses

All endpoints return consistent error responses:

### 400 Bad Request

```json
{
  "error": "Bad Request",
  "message": "Invalid startDate format: expected YYYY-MM-DD or relative (7d, 30d)"
}
```

### 401 Unauthorized

```json
{
  "error": "Unauthorized",
  "message": "Invalid API key"
}
```

### 429 Too Many Requests

```json
{
  "error": "Too Many Requests",
  "message": "Rate limit exceeded. Please try again later."
}
```

### 500 Internal Server Error

```json
{
  "error": "Internal Server Error",
  "message": "An unexpected error occurred"
}
```

---

## Rate Limiting

The simulator implements rate limiting matching the Cursor API:

| Endpoint Type | Limit |
|---------------|-------|
| AI Code Tracking | 20 requests/minute per endpoint |
| Admin API | 20 requests/minute |

---

## Date Formats

All date parameters support multiple formats:

| Format | Example | Description |
|--------|---------|-------------|
| ISO 8601 | `2025-01-15T10:30:00Z` | Full timestamp |
| Date only | `2025-01-15` | Recommended for caching |
| Relative | `7d`, `30d` | Days ago |
| Keywords | `now`, `today`, `yesterday` | Current/previous day |

---

## Comparison with Cursor Business API

The simulator endpoints are designed to be API-compatible with the production Cursor Business API. Key differences:

| Aspect | Simulator | Production |
|--------|-----------|------------|
| Authentication | Optional (configurable) | Required |
| Data source | Synthetic generated data | Real usage data |
| Rate limits | Configurable | Fixed per tier |
| Enterprise features | All available | Plan-dependent |

When migrating from the simulator to the production API, ensure you:
1. Configure proper API key authentication
2. Handle Enterprise-only endpoint access restrictions
3. Adjust for any response timing differences

---

## OpenAPI Specification

The complete OpenAPI 3.1 specification is available at:
- `specs/openapi/cursor-api.yaml`

Use this for code generation and detailed schema validation.
