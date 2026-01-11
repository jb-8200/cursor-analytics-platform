# REST API Reference: cursor-sim

This document provides the complete REST API specification for the cursor-sim service. The simulator exposes 19 endpoints that mimic the Cursor Business Activity API, enabling the aggregator to fetch developer activity data without requiring access to production Cursor credentials.

## Source of Truth

**For accurate API specifications, always reference these documents:**

| Document | Description |
|----------|-------------|
| **specs/openapi/cursor-api.yaml** | OpenAPI 3.1.0 specification (19 endpoints) - CANONICAL |
| [cursor_overview.md](cursor_overview.md) | Authentication, rate limits, caching, error handling |
| [cursor_admin.md](cursor_admin.md) | Admin API - Team management, usage data, spending |
| [cursor_analytics.md](cursor_analytics.md) | Analytics API - Team metrics, DAU, model usage |
| [cursor_codetrack.md](cursor_codetrack.md) | AI Code Tracking API - Per-commit metrics (Enterprise) |
| [cursor_agents.md](cursor_agents.md) | Cloud Agents API - Programmatic agent management |

**Implementation Status:** Last updated January 10, 2026. Includes P1-F02 (Admin Configuration) and P4-F05 (External Data Sources).

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

### Admin Configuration API (P1-F02)

These endpoints enable runtime configuration of the simulator without restarting. They support data regeneration, seed management, and statistics retrieval.

#### GET /admin/config

Retrieve current simulator configuration including generation parameters and enabled features.

**Response:**

```json
{
  "mode": "runtime",
  "days": 90,
  "velocity": "medium",
  "developers": 50,
  "max_commits": 1000,
  "external_sources": {
    "harvey": {
      "enabled": true,
      "models": ["gpt-4", "claude-3-sonnet"]
    },
    "copilot": {
      "enabled": true
    },
    "qualtrics": {
      "enabled": true,
      "surveyId": "SV_survey_id"
    }
  }
}
```

#### GET /admin/stats

Retrieve simulator statistics including data counts and generation metrics.

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `time_series` | boolean | false | Include time series data if true |

**Response:**

```json
{
  "commits_count": 45000,
  "prs_count": 4500,
  "reviews_count": 9000,
  "issues_count": 1350,
  "developers_count": 50,
  "last_generation_time": "5.2s",
  "time_series": [
    {
      "date": "2026-01-10",
      "commits": 450,
      "prs": 45
    }
  ]
}
```

#### POST /admin/regenerate

Regenerate simulation data with new parameters. Supports two modes:
- **append**: Adds new data to existing storage
- **override**: Clears all data and generates fresh dataset

**Request Body:**

```json
{
  "mode": "override",
  "days": 90,
  "velocity": "medium",
  "developers": 50,
  "max_commits": 1000
}
```

**Response:**

```json
{
  "status": "success",
  "mode": "override",
  "data_cleaned": true,
  "commits_added": 45000,
  "prs_added": 4500,
  "reviews_added": 9000,
  "issues_added": 1350,
  "total_commits": 45000,
  "total_prs": 4500,
  "total_developers": 50,
  "duration": "5.2s",
  "config": {
    "mode": "override",
    "days": 90,
    "velocity": "medium",
    "developers": 50,
    "max_commits": 1000
  }
}
```

#### POST /admin/seed

Upload a new seed file (JSON, YAML, or CSV) to configure the simulator.

**Request Body:**

```json
{
  "seed_data": "{...seed file content...}",
  "format": "json"
}
```

**Response:**

```json
{
  "status": "success",
  "developers_count": 50,
  "repos_count": 10
}
```

#### GET /admin/seed/presets

List all available seed presets for quick configuration.

**Response:**

```json
{
  "presets": [
    {
      "name": "small",
      "description": "Small team (10 devs, 30 days)",
      "developers": 10,
      "days": 30
    },
    {
      "name": "medium",
      "description": "Medium team (50 devs, 90 days)",
      "developers": 50,
      "days": 90
    },
    {
      "name": "large",
      "description": "Large team (500 devs, 180 days)",
      "developers": 500,
      "days": 180
    }
  ]
}
```

---

### External Data Sources API (P4-F05)

These endpoints provide integrations with third-party data sources. They are only active when configured in the seed file.

#### GET /harvey/api/v1/history/usage

Returns Harvey AI legal document analysis usage events. Requires Harvey to be enabled in seed configuration.

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `from` | string | Yes | Start date (YYYY-MM-DD or RFC3339 format) |
| `to` | string | Yes | End date (YYYY-MM-DD or RFC3339 format) |
| `user` | string | No | Filter by user email |
| `task` | string | No | Filter by task type (e.g., legal_review, contract_analysis) |
| `page` | integer | No | Page number (default 1) |
| `page_size` | integer | No | Items per page (default 50, max 100) |

**Response:**

```json
{
  "data": [
    {
      "id": "event_001",
      "timestamp": "2026-01-10T15:30:00Z",
      "user_email": "dev@example.com",
      "task_type": "legal_review",
      "document_name": "contract.pdf",
      "duration_seconds": 180,
      "status": "completed"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 50,
    "totalCount": 156,
    "totalPages": 4,
    "hasNextPage": true
  }
}
```

#### GET /reports/getMicrosoft365CopilotUsageUserDetail(period='...')

Returns Microsoft 365 Copilot usage metrics. OData-compliant endpoint. Supports JSON or CSV export. Requires Copilot to be enabled in seed configuration.

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `period` | string | Yes | Report period - D7, D30, D90, or D180 |
| `$format` | string | No | Response format (application/json or text/csv) |

**Response (JSON):**

```json
{
  "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#reports.getM365CopilotUsageUserDetail()",
  "value": [
    {
      "reportRefreshDate": "2026-01-10",
      "userPrincipalName": "user@company.com",
      "displayName": "User Name",
      "reportPeriod": 30,
      "copilotCompletionEventsCount": 145,
      "copilotCompletionTokenCount": 8234,
      "copilotCitations": 23
    }
  ]
}
```

**Response (CSV):**

```csv
Report Refresh Date,User Principal Name,Display Name,Report Period,Copilot Completion Events,Copilot Completion Tokens,Copilot Citations
2026-01-10,user@company.com,User Name,30,145,8234,23
```

#### POST /API/v3/surveys/{surveyId}/export-responses

Starts a Qualtrics survey response export job. Requires Qualtrics to be enabled in seed configuration. Returns immediately with a progress ID for polling.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `surveyId` | string | Qualtrics survey ID |

**Response:**

```json
{
  "result": {
    "progressId": "ES_abc123def456",
    "status": "inProgress",
    "percentComplete": 0,
    "estimatedSeconds": 120
  }
}
```

#### GET /API/v3/surveys/{surveyId}/export-responses/{progressId}

Polls the status of an export job. Response includes completion percentage and file ID when ready.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `surveyId` | string | Qualtrics survey ID |
| `progressId` | string | Progress ID from export start response |

**Response:**

```json
{
  "result": {
    "progressId": "ES_abc123def456",
    "status": "complete",
    "percentComplete": 100,
    "fileId": "FILE_xyz789abc123"
  }
}
```

**Possible Status Values:**
- `inProgress`: Export job still running
- `complete`: Export ready for download
- `failed`: Export job failed

#### GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file

Downloads the exported survey responses as a ZIP file containing CSV data.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `surveyId` | string | Qualtrics survey ID |
| `fileId` | string | File ID from progress poll response |

**Response:**

- Content-Type: `application/zip`
- Body: Binary ZIP file containing `survey_responses.csv`

---

### Health Check

#### GET /health

Check simulator health status. This endpoint does not require authentication.

**Response:**

```json
{
  "status": "healthy",
  "mode": "runtime",
  "seed_loaded": true,
  "developers_count": 50,
  "commits_count": 45000,
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

## Endpoint Summary

**Total Endpoints: 19**

| Category | Count | Endpoints |
|----------|-------|-----------|
| **Admin API** | 4 | GET /teams/members, POST /teams/daily-usage-data, POST /teams/filtered-usage-events, POST /teams/spend |
| **Admin Configuration** | 5 | GET /admin/config, GET /admin/stats, POST /admin/regenerate, POST /admin/seed, GET /admin/seed/presets |
| **AI Code Tracking** | 4 | GET /analytics/ai-code/commits, GET /analytics/ai-code/commits.csv, GET /analytics/ai-code/changes, GET /analytics/ai-code/changes.csv |
| **External Data Sources** | 5 | GET /harvey/api/v1/history/usage, GET /reports/getMicrosoft365CopilotUsageUserDetail, POST /API/v3/surveys/{surveyId}/export-responses, GET /API/v3/surveys/{surveyId}/export-responses/{progressId}, GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file |
| **Health** | 1 | GET /health |

---

## OpenAPI Specification

The complete OpenAPI 3.1.0 specification is the canonical source of truth and is available at:

- **`specs/openapi/cursor-api.yaml`** - Main API specification (19 endpoints, all schemas)
- **`specs/openapi/github-sim-api.yaml`** - GitHub simulation API (16 endpoints)

Use the OpenAPI spec for:
- Code generation (TypeScript, Python, Go, etc.)
- Automated documentation
- Contract testing
- API client integration

---

## Implementation Notes

### Admin Configuration API (P1-F02)
The Admin Configuration API enables runtime management of the simulator without restarting. This is useful for:
- Scaling data generation (changing developer count, history length)
- Switching between generation modes (append vs override)
- Configuring external data sources
- Monitoring simulator health and statistics

### External Data Sources API (P4-F05)
The External Data Sources API provides simulated integrations with third-party tools:
- **Harvey AI**: Legal document analysis usage tracking
- **Microsoft 365 Copilot**: OData-compliant usage metrics (JSON or CSV)
- **Qualtrics**: Three-step survey export workflow (start → poll → download)

All external data sources are optional and configured via the seed file.

### Authentication
All endpoints except `/health` require HTTP Basic Authentication:
```
Authorization: Basic {base64('YOUR_API_KEY:')}
```
The API key is used as the username with an empty password.
