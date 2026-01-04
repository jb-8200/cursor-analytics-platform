# Cursor Business API Reference

**Version**: 1.0.0  
**Last Updated**: January 2026  

This document provides a reference for the Cursor for Business API based on official documentation research. The Cursor Analytics Platform simulator is designed to produce data compatible with these API patterns.

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

The Cursor Analytics Platform simulator implements API endpoints that closely match this contract. The following mappings apply from simulator to Cursor API.

The simulator endpoint `GET /v1/org/users` maps to the Cursor endpoint `GET /teams/members`. The simulator endpoint `GET /v1/stats/activity` maps to `POST /teams/usage-events`. The simulator endpoint `GET /v1/stats/daily-usage` maps to `POST /teams/daily-usage-data`.

The primary differences are that the simulator uses GET with query parameters while Cursor uses POST with JSON body for some endpoints. Additionally, the simulator does not implement authentication since it runs locally for development.
