# Cursor Business API Reference

**Version**: 2.0.0
**Last Updated**: January 11, 2026

This document provides a reference for the Cursor for Business API based on official documentation research. The Cursor Analytics Platform simulator is designed to produce data compatible with these API patterns.

**Note**: This document includes both the production Cursor Business API contract (Admin, AI Code Tracking, Background Agents) AND simulator-specific management APIs (Admin Configuration, External Data Sources) that extend the contract for development and testing purposes.

## API Overview

Cursor provides a three-tier API architecture for business and enterprise customers. The Admin API enables team management and usage analytics. The Background Agents API provides programmatic control over AI coding agents. The Enterprise-exclusive AI Code Tracking API offers commit-level metrics for detailed code contribution analysis.

## Authentication

All Cursor APIs use HTTP Basic Authentication. The API key serves as the username, and no password is required. API keys follow the format `key_` followed by 64 hexadecimal characters.

To authenticate, include an Authorization header with the Base64-encoded credentials where the format is `key_xxx:` (note the trailing colon for empty password).

```
Authorization: Basic a2V5X3h4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4Og==
```

API keys are generated from the Cursor dashboard under Settings and then Cursor Admin API Keys.

## Admin API

The Admin API provides endpoints for team management and usage analytics.

### Base URL

```
https://api.cursor.com
```

### Get Team Members

Returns a list of all members in the organization.

**Request**

```
GET /teams/members
```

**Response**

```json
{
    "members": [
        {
            "email": "developer@company.com",
            "name": "Jane Developer",
            "role": "member",
            "dateAdded": "2025-06-15T10:30:00Z"
        }
    ],
    "total": 50,
    "numPages": 1,
    "currentPage": 1,
    "pageSize": 100
}
```

**Pagination Parameters**

The `page` parameter accepts 1-indexed page numbers. The `pageSize` parameter accepts values of 10, 25, or 100. The `searchTerm` parameter filters results by email.

### Daily Usage Data

Returns per-user usage metrics aggregated by day.

**Request**

```
POST /teams/daily-usage-data
Content-Type: application/json

{
    "startDate": "2026-01-01",
    "endDate": "2026-01-31",
    "page": 1,
    "pageSize": 100
}
```

**Response**

```json
{
    "usageData": [
        {
            "email": "developer@company.com",
            "date": "2026-01-15",
            "isActive": true,
            "totalTabsShown": 98,
            "totalTabsAccepted": 22,
            "totalLinesAdded": 310,
            "totalLinesDeleted": 48,
            "acceptedLinesAdded": 105,
            "composerRequests": 6,
            "chatRequests": 12,
            "agentRequests": 3,
            "cmdkUsages": 8,
            "mostUsedModel": "claude-4-opus"
        }
    ],
    "numPages": 5,
    "currentPage": 1,
    "pageSize": 100,
    "hasNextPage": true,
    "hasPreviousPage": false
}
```

**Field Descriptions**

The `totalTabsShown` field represents Tab completion suggestions displayed to the user. The `totalTabsAccepted` field represents Tab completions the user accepted. The `acceptedLinesAdded` field counts lines of code from accepted AI suggestions. The `composerRequests` field counts uses of the Composer feature for multi-file edits. The `chatRequests` field counts conversation interactions with the AI. The `agentRequests` field counts background agent task submissions. The `cmdkUsages` field counts inline edit command uses (Cmd+K or Ctrl+K).

**Date Range Limitation**

The maximum date range is 90 days. Requests exceeding this limit receive a 400 error.

### Team Spending

Returns cost information aggregated by user.

**Request**

```
POST /teams/spend
Content-Type: application/json

{
    "startDate": "2026-01-01",
    "endDate": "2026-01-31",
    "page": 1,
    "pageSize": 100
}
```

**Response**

```json
{
    "spending": [
        {
            "email": "developer@company.com",
            "spendCents": 4523,
            "fastPremiumRequests": 156,
            "spendingLimit": 10000
        }
    ],
    "totalSpendCents": 45230,
    "numPages": 1,
    "currentPage": 1
}
```

### Usage Events

Returns granular event-level data including token consumption.

**Request**

```
POST /teams/usage-events
Content-Type: application/json

{
    "startDate": "2026-01-01",
    "endDate": "2026-01-07",
    "page": 1,
    "pageSize": 100
}
```

**Response**

```json
{
    "events": [
        {
            "email": "developer@company.com",
            "timestamp": "2026-01-15T10:30:00Z",
            "eventType": "chat",
            "model": "claude-4-opus",
            "inputTokens": 1500,
            "outputTokens": 450,
            "cacheWriteTokens": 200,
            "cacheReadTokens": 800,
            "totalCents": 12
        }
    ],
    "numPages": 10,
    "currentPage": 1
}
```

## Background Agents API

The Background Agents API enables programmatic creation and management of AI coding agents.

### Base URL

```
https://api.cursor.com/v0
```

### Authentication

Background Agents API uses Bearer token authentication.

```
Authorization: Bearer <token>
```

### Create Agent

Creates a new background agent to work on a coding task.

**Request**

```
POST /v0/agents
Content-Type: application/json

{
    "prompt": "Fix the bug in the authentication module",
    "repository": {
        "provider": "github",
        "owner": "company",
        "name": "backend-api",
        "branch": "main"
    },
    "webhookUrl": "https://company.com/webhooks/cursor"
}
```

**Response**

```json
{
    "agentId": "agent-abc123",
    "status": "CREATING",
    "createdAt": "2026-01-15T10:30:00Z"
}
```

### Get Agent Status

Retrieves the current status of an agent.

**Request**

```
GET /v0/agents/{agentId}
```

**Response**

```json
{
    "agentId": "agent-abc123",
    "status": "FINISHED",
    "prompt": "Fix the bug in the authentication module",
    "repository": {
        "provider": "github",
        "owner": "company",
        "name": "backend-api",
        "branch": "main"
    },
    "targetBranch": "cursor/fix-auth-bug",
    "pullRequestUrl": "https://github.com/company/backend-api/pull/123",
    "summary": "Fixed null pointer exception in JWT validation",
    "createdAt": "2026-01-15T10:30:00Z",
    "completedAt": "2026-01-15T10:45:00Z"
}
```

**Agent Status Values**

The `CREATING` status indicates the agent environment is being initialized. The `RUNNING` status indicates the agent is actively working on the task. The `FINISHED` status indicates the agent completed successfully. The `ERROR` status indicates the agent encountered an error and stopped.

### List Agents

Returns all agents for the organization.

**Request**

```
GET /v0/agents
```

**Response**

```json
{
    "agents": [
        {
            "agentId": "agent-abc123",
            "status": "FINISHED",
            "prompt": "Fix the bug...",
            "createdAt": "2026-01-15T10:30:00Z"
        }
    ]
}
```

### Agent Webhooks

Webhooks notify external systems when agent status changes.

**Webhook Payload**

```json
{
    "agentId": "agent-abc123",
    "status": "FINISHED",
    "repository": {
        "provider": "github",
        "owner": "company",
        "name": "backend-api"
    },
    "targetBranch": "cursor/fix-auth-bug",
    "pullRequestUrl": "https://github.com/company/backend-api/pull/123",
    "summary": "Fixed null pointer exception in JWT validation"
}
```

**Signature Verification**

Webhooks include an `X-Webhook-Signature` header containing an HMAC-SHA256 signature. The format is `sha256=<hex_digest>`. Verify by computing the HMAC of the raw request body using your webhook secret.

## AI Code Tracking API (Enterprise Only)

The AI Code Tracking API provides commit-level attribution of AI-generated code.

### Get Commit Statistics

Returns AI contribution metrics at the commit level.

**Request**

```
GET /analytics/ai-code/commits?startDate=7d&endDate=now&page=1&pageSize=100
```

**Response**

```json
{
    "commits": [
        {
            "commitHash": "abc123def456",
            "userId": "encrypted-user-id",
            "repoName": "backend-api",
            "branchName": "feature/auth",
            "isPrimaryBranch": false,
            "timestamp": "2026-01-15T10:30:00Z",
            "tabLinesAdded": 45,
            "tabLinesDeleted": 12,
            "composerLinesAdded": 23,
            "composerLinesDeleted": 5,
            "nonAiLinesAdded": 78,
            "nonAiLinesDeleted": 34
        }
    ],
    "total": 1500,
    "page": 1,
    "pageSize": 100,
    "hasNextPage": true
}
```

**Field Descriptions**

The `tabLinesAdded` and `tabLinesDeleted` fields represent code from inline Tab completions. The `composerLinesAdded` and `composerLinesDeleted` fields represent code from Composer or chat-applied diffs. The `nonAiLinesAdded` and `nonAiLinesDeleted` fields represent manually written code.

### CSV Export

For large data exports, a streaming CSV endpoint is available.

**Request**

```
GET /analytics/ai-code/commits.csv?startDate=30d&endDate=now
```

This endpoint streams results with 10,000 records per page.

### Code Change Metrics

Returns granular accepted AI changes grouped by changeId.

**Request**

```
GET /analytics/ai-code/changes?startDate=7d&endDate=now&page=1&pageSize=100
```

**Response**

```json
{
    "changes": [
        {
            "changeId": "749356201",
            "userId": "encrypted-user-id",
            "source": "COMPOSER",
            "model": "claude-4o",
            "totalLinesAdded": 18,
            "totalLinesDeleted": 4,
            "timestamp": "2026-01-15T10:30:00Z",
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
    "total": 2500,
    "page": 1,
    "pageSize": 100,
    "hasNextPage": true
}
```

---

## Admin Configuration API (Simulator Extension)

The Admin Configuration API enables runtime management of the simulator without restart. This is a simulator-specific feature for development and testing.

### Base URL

```
http://localhost:8080  (local simulator)
https://cursor-sim-xxxxx.a.run.app  (GCP Cloud Run)
```

### Get Configuration

Retrieves current simulator configuration including generation parameters and enabled features.

**Request**

```
GET /admin/config
Authorization: Basic {base64('API_KEY:')}
```

**Response**

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

### Get Statistics

Retrieves simulator statistics including data counts and generation metrics.

**Request**

```
GET /admin/stats?time_series=true
Authorization: Basic {base64('API_KEY:')}
```

**Response**

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

### Regenerate Data

Regenerates simulation data with new parameters. Supports append (add to existing) or override (replace all) modes.

**Request**

```
POST /admin/regenerate
Content-Type: application/json
Authorization: Basic {base64('API_KEY:')}

{
    "mode": "override",
    "days": 90,
    "velocity": "medium",
    "developers": 50,
    "max_commits": 1000
}
```

**Response**

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
    "total_developers": 50,
    "duration": "5.2s"
}
```

### Upload Seed File

Configure the simulator using a seed file (JSON, YAML, or CSV).

**Request**

```
POST /admin/seed
Content-Type: application/json
Authorization: Basic {base64('API_KEY:')}

{
    "seed_data": "{...seed file content...}",
    "format": "json"
}
```

**Response**

```json
{
    "status": "success",
    "developers_count": 50,
    "repos_count": 10
}
```

### Get Seed Presets

Lists available seed presets for quick configuration.

**Request**

```
GET /admin/seed/presets
Authorization: Basic {base64('API_KEY:')}
```

**Response**

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

## External Data Sources API (Simulator Extension)

The External Data Sources API provides simulated integrations with third-party tools. These endpoints are only active when configured in the seed file. This is a simulator-specific feature for development and testing.

### Harvey AI Usage

Returns Harvey AI legal document analysis usage events.

**Request**

```
GET /harvey/api/v1/history/usage?from=2026-01-01&to=2026-01-31&page=1&page_size=50
Authorization: Basic {base64('API_KEY:')}
```

**Query Parameters**
- `from`: Start date (YYYY-MM-DD format)
- `to`: End date (YYYY-MM-DD format)
- `user`: Optional email filter
- `task`: Optional task type filter (legal_review, contract_analysis)
- `page`: Page number (default 1)
- `page_size`: Items per page (default 50, max 100)

**Response**

```json
{
    "data": [
        {
            "id": "event_001",
            "timestamp": "2026-01-15T10:30:00Z",
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

### Microsoft 365 Copilot Usage

Returns Microsoft 365 Copilot usage metrics. OData-compliant endpoint supporting JSON or CSV export.

**Request**

```
GET /reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=application/json
Authorization: Basic {base64('API_KEY:')}
```

**Query Parameters**
- `period`: Report period - D7, D30, D90, or D180
- `$format`: Response format - application/json (default) or text/csv

**Response (JSON)**

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

### Qualtrics Survey Export

Qualtrics survey export follows a three-step workflow: start export, poll progress, download file.

#### Step 1: Start Export

**Request**

```
POST /API/v3/surveys/{surveyId}/export-responses
Authorization: Basic {base64('API_KEY:')}
```

**Response**

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

#### Step 2: Poll Progress

**Request**

```
GET /API/v3/surveys/{surveyId}/export-responses/{progressId}
Authorization: Basic {base64('API_KEY:')}
```

**Response**

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

**Possible Status Values**: `inProgress`, `complete`, `failed`

#### Step 3: Download File

**Request**

```
GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file
Authorization: Basic {base64('API_KEY:')}
```

**Response**: Binary ZIP file containing `survey_responses.csv`

---

## Telemetry Event Types

Based on the API response fields and internal documentation, Cursor tracks the following telemetry event categories.

### Tab Completions

Tab completions are tracked through the `totalTabsShown` and `totalTabsAccepted` metrics. The acceptance rate can be calculated as accepted divided by shown, multiplied by 100.

### Chat Interactions

Chat interactions are counted in the `chatRequests` field. Each conversation turn counts as one interaction.

### Composer Operations

Multi-file edit operations through Composer are counted in `composerRequests`.

### Inline Edit (Cmd+K)

Inline edit commands are counted in `cmdkUsages`.

### Agent Activity

Background agent submissions are counted in `agentRequests`. The separate `bugbotUsages` field may track automated bug detection runs.

## Rate Limits

Rate limits are enforced but specific thresholds are not publicly documented. When rate limited, the API returns HTTP 429 with a `Retry-After` header indicating seconds until the next request is allowed.

Best practices include implementing exponential backoff on 429 responses, caching responses when possible, and using date range queries efficiently to minimize requests.

## Error Responses

Error responses follow a consistent format.

```json
{
    "error": {
        "code": "INVALID_DATE_RANGE",
        "message": "Date range cannot exceed 90 days",
        "details": {
            "maxDays": 90,
            "requestedDays": 120
        }
    }
}
```

Common error codes include `INVALID_DATE_RANGE` for date range exceeding limits, `UNAUTHORIZED` for invalid or missing API key, `RATE_LIMITED` for too many requests, and `NOT_FOUND` for requested resource not existing.

## Simulator Compatibility

The Cursor Analytics Platform simulator implements API endpoints that closely match the production Cursor API contract. Additionally, the simulator provides management APIs (Admin Configuration and External Data Sources) that extend the contract for development and testing purposes.

### Production API Compliance

The following mappings apply from simulator to Cursor API for the core business APIs:

| Simulator Endpoint | Cursor Endpoint | Purpose |
|---|---|---|
| `GET /teams/members` | `GET /teams/members` | Team member list |
| `POST /teams/daily-usage-data` | `POST /teams/daily-usage-data` | Daily usage metrics |
| `POST /teams/filtered-usage-events` | `POST /teams/usage-events` | Granular usage events |
| `POST /teams/spend` | `POST /teams/spend` | Team spending data |
| `GET /analytics/ai-code/commits` | `GET /analytics/ai-code/commits` | Commit-level AI metrics |
| `GET /analytics/ai-code/commits.csv` | `GET /analytics/ai-code/commits.csv` | CSV export |
| `GET /analytics/ai-code/changes` | `GET /analytics/ai-code/changes` | Code change metrics |
| `GET /analytics/ai-code/changes.csv` | `GET /analytics/ai-code/changes.csv` | CSV export |

### Simulator Extensions

The simulator provides additional management APIs for development:

| Endpoint | Purpose | Feature |
|---|---|---|
| `GET /admin/config` | Retrieve configuration | Admin Configuration (P1-F02) |
| `GET /admin/stats` | Retrieve statistics | Admin Configuration (P1-F02) |
| `POST /admin/regenerate` | Regenerate data | Admin Configuration (P1-F02) |
| `POST /admin/seed` | Upload seed file | Admin Configuration (P1-F02) |
| `GET /admin/seed/presets` | List presets | Admin Configuration (P1-F02) |
| `GET /harvey/api/v1/history/usage` | Harvey AI usage | External Data Sources (P4-F05) |
| `GET /reports/getMicrosoft365CopilotUsageUserDetail` | Copilot metrics | External Data Sources (P4-F05) |
| `POST /API/v3/surveys/{id}/export-responses` | Start export | External Data Sources (P4-F05) |
| `GET /API/v3/surveys/{id}/export-responses/{progressId}` | Poll progress | External Data Sources (P4-F05) |
| `GET /API/v3/surveys/{id}/export-responses/{fileId}/file` | Download export | External Data Sources (P4-F05) |

### Key Differences

- **Authentication**: The simulator supports optional Basic Auth (configurable), while production requires API keys
- **Data Source**: Simulator uses synthetically generated data; production uses real usage data
- **Rate Limits**: Simulator rate limits are configurable; production enforces fixed limits per tier
- **Management APIs**: Simulator includes admin configuration and external data source endpoints for development
- **Hosting**: Local deployment (`localhost:8080`), Docker Compose, or GCP Cloud Run

### Version Tracking

- **Simulator OpenAPI Spec**: `specs/openapi/cursor-api.yaml` (19 endpoints)
- **Production API Contract**: Documented in this reference (Admin, AI Code Tracking, Background Agents)
- **Last Sync**: January 11, 2026
