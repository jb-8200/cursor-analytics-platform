# REST API Reference: cursor-sim

This document provides the complete REST API specification for the cursor-sim service. The simulator exposes endpoints that mimic the Cursor Business Activity API, enabling the aggregator to fetch developer activity data without requiring access to production Cursor credentials.

---

## Base URL

When running locally via Docker Compose, the simulator is accessible at `http://localhost:8080`. When services communicate within the Docker network, use `http://cursor-sim:8080`.

All endpoints are prefixed with `/v1` to indicate the API version.

---

## Authentication

The simulator does not require authentication in the current development configuration. Future versions may implement API key authentication for parity with the real Cursor Business API.

---

## Common Response Patterns

All endpoints return JSON responses with consistent structure. Successful responses include a 2xx status code and the requested data. Error responses include a 4xx or 5xx status code along with an error object containing a message and optional details.

Pagination is implemented consistently across list endpoints. Paginated responses include a `pagination` object with `total` (total number of records), `page` (current page number, 1-indexed), `pageSize` (number of records per page), `hasNextPage` (boolean indicating more pages exist), and `hasPreviousPage` (boolean indicating previous pages exist).

---

## Endpoints

### GET /v1/health

The health check endpoint returns the current status of the simulator service. Container orchestration systems use this endpoint to determine if the service is ready to accept traffic.

**Request**

This endpoint accepts no parameters.

```http
GET /v1/health HTTP/1.1
Host: localhost:8080
```

**Response**

A successful response returns HTTP 200 with the service status information.

```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": 3600,
  "developers": 50,
  "eventsGenerated": 125000,
  "configuration": {
    "velocity": "high",
    "fluctuation": 0.2,
    "seed": 1704067200
  }
}
```

The `status` field is either "healthy" or "degraded". The `uptime` field indicates seconds since service start. The `developers` field shows the configured number of simulated developers. The `eventsGenerated` field is the total count of events currently stored.

**Error Responses**

If the service is starting up or experiencing issues, it returns HTTP 503 Service Unavailable with a message explaining the condition.

---

### GET /v1/org/users

This endpoint returns the list of simulated developers. Each developer has a unique profile with name, email, team assignment, and seniority level. These profiles remain consistent for the lifetime of the simulator process.

**Request**

The endpoint accepts optional query parameters for pagination.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | integer | 1 | Page number (1-indexed) |
| pageSize | integer | 50 | Records per page (max: 100) |

```http
GET /v1/org/users?page=1&pageSize=25 HTTP/1.1
Host: localhost:8080
```

**Response**

A successful response returns HTTP 200 with the list of developers.

```json
{
  "users": [
    {
      "id": "dev-a1b2c3d4",
      "name": "Alice Chen",
      "email": "alice.chen@example.com",
      "team": "Backend",
      "seniority": "senior",
      "isActive": true,
      "createdAt": "2024-10-01T00:00:00Z"
    },
    {
      "id": "dev-e5f6g7h8",
      "name": "Bob Martinez",
      "email": "bob.martinez@example.com",
      "team": "Frontend",
      "seniority": "mid",
      "isActive": true,
      "createdAt": "2024-10-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 50,
    "page": 1,
    "pageSize": 25,
    "hasNextPage": true,
    "hasPreviousPage": false
  }
}
```

The `seniority` field is one of "junior", "mid", or "senior". The `isActive` field indicates whether the developer has generated events in the current simulation.

**Error Responses**

If pagination parameters are invalid, the endpoint returns HTTP 400 Bad Request with details about the validation failure.

---

### GET /v1/stats/activity

This endpoint returns activity events within a specified time range. Events represent developer interactions with Cursor AI features including code completions, chat messages, and inline edit prompts.

This is the primary endpoint that the cursor-analytics-core service polls to ingest data for aggregation.

**Request**

The endpoint requires time range parameters and accepts optional pagination.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| from | ISO 8601 datetime | Yes | Start of time range (inclusive) |
| to | ISO 8601 datetime | Yes | End of time range (inclusive) |
| page | integer | No | Page number (default: 1) |
| pageSize | integer | No | Records per page (default: 100, max: 1000) |
| developerId | string | No | Filter to specific developer |
| eventType | string | No | Filter to specific event type |

```http
GET /v1/stats/activity?from=2025-01-01T00:00:00Z&to=2025-01-01T23:59:59Z&pageSize=100 HTTP/1.1
Host: localhost:8080
```

**Response**

A successful response returns HTTP 200 with the list of events.

```json
{
  "events": [
    {
      "id": "evt-uuid-12345",
      "type": "cpp_suggestion_shown",
      "timestamp": "2025-01-01T10:30:00Z",
      "developerId": "dev-a1b2c3d4",
      "metadata": {
        "linesAffected": 5,
        "model": "claude-4-opus",
        "promptTokens": 150,
        "completionTokens": 45,
        "language": "typescript",
        "sessionId": "sess-xyz789"
      }
    },
    {
      "id": "evt-uuid-12346",
      "type": "cpp_suggestion_accepted",
      "timestamp": "2025-01-01T10:30:02Z",
      "developerId": "dev-a1b2c3d4",
      "metadata": {
        "linesAffected": 5,
        "model": "claude-4-opus",
        "acceptanceLatencyMs": 2000,
        "sessionId": "sess-xyz789"
      }
    },
    {
      "id": "evt-uuid-12347",
      "type": "chat_message",
      "timestamp": "2025-01-01T10:35:00Z",
      "developerId": "dev-a1b2c3d4",
      "metadata": {
        "model": "claude-4-opus",
        "promptTokens": 500,
        "completionTokens": 300,
        "conversationId": "conv-abc123"
      }
    },
    {
      "id": "evt-uuid-12348",
      "type": "cmd_k_prompt",
      "timestamp": "2025-01-01T10:40:00Z",
      "developerId": "dev-a1b2c3d4",
      "metadata": {
        "linesAffected": 12,
        "model": "claude-4-opus",
        "promptTokens": 200,
        "completionTokens": 150,
        "editType": "refactor"
      }
    }
  ],
  "pagination": {
    "total": 5000,
    "page": 1,
    "pageSize": 100,
    "hasNextPage": true,
    "hasPreviousPage": false
  },
  "meta": {
    "queryTimeMs": 45,
    "rangeStart": "2025-01-01T00:00:00Z",
    "rangeEnd": "2025-01-01T23:59:59Z"
  }
}
```

**Event Types**

The `type` field indicates the kind of developer activity.

The `cpp_suggestion_shown` type indicates that a Tab completion suggestion was displayed to the developer. This happens when Cursor's autocomplete feature generates a code suggestion. The metadata includes `linesAffected` (number of lines in the suggestion), `model` (the AI model used), `language` (programming language), and `sessionId` (groups related events).

The `cpp_suggestion_accepted` type indicates that the developer accepted a previously shown Tab completion. This event always follows a corresponding `cpp_suggestion_shown` event with the same `sessionId`. The metadata includes `acceptanceLatencyMs` (time between shown and accepted).

The `chat_message` type represents an interaction with Cursor's chat panel. This includes asking questions, requesting explanations, or generating code through conversation. The metadata includes `conversationId` (groups messages in the same conversation thread).

The `cmd_k_prompt` type represents an inline edit triggered by the Cmd+K shortcut. The developer selects code and provides a prompt to modify it. The metadata includes `editType` (categorizes the edit as "refactor", "fix", "explain", "generate", or "other").

**Error Responses**

If the `from` or `to` parameters are missing, the endpoint returns HTTP 400 Bad Request with a message indicating which parameter is required.

If the date format is invalid, the endpoint returns HTTP 400 Bad Request with details about the expected format.

If the time range exceeds 90 days, the endpoint returns HTTP 400 Bad Request indicating the maximum allowed range.

---

### POST /v1/admin/reset

This administrative endpoint resets the simulator state, clearing all generated events and optionally reconfiguring simulation parameters. This is useful during development and testing to start with a fresh dataset.

**Request**

The request body is optional. If provided, it can override simulation parameters.

```http
POST /v1/admin/reset HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "developers": 100,
  "velocity": "medium",
  "fluctuation": 0.3,
  "seed": 12345
}
```

| Field | Type | Description |
|-------|------|-------------|
| developers | integer | Number of developers to simulate |
| velocity | string | Event rate: "low", "medium", "high" |
| fluctuation | float | Variance coefficient (0.0-1.0) |
| seed | integer | Random seed for reproducibility |

**Response**

A successful response returns HTTP 200 with the new configuration.

```json
{
  "status": "reset_complete",
  "configuration": {
    "developers": 100,
    "velocity": "medium",
    "fluctuation": 0.3,
    "seed": 12345
  },
  "eventsGenerated": 450000,
  "generationTimeMs": 2500
}
```

**Error Responses**

If configuration values are invalid, the endpoint returns HTTP 400 Bad Request with validation details.

---

## Rate Limiting

The simulator does not implement rate limiting since it is designed for local development use. However, callers should implement reasonable polling intervals (recommended: 60 seconds minimum) to avoid unnecessary load.

---

## Event Generation Details

The simulator generates events following statistical patterns to mimic real developer behavior. Understanding these patterns helps interpret the data correctly.

Events follow a Poisson distribution for timing, which means events are independent and occur at a constant average rate. The rate varies by velocity setting and is modified by a circadian rhythm function that produces more events during typical working hours (9 AM to 6 PM) and fewer events overnight.

The acceptance rate for `cpp_suggestion_accepted` events correlates with developer seniority. Senior developers accept approximately 90% of suggestions, mid-level developers accept 75%, and junior developers accept 60%. These rates have some random variance controlled by the fluctuation parameter.

Event metadata is generated with realistic values. Token counts follow typical distributions for code completions. Line counts vary based on edit type and developer seniority. Session IDs group related events to simulate realistic editing sessions.

---

## Comparison with Cursor Business API

The simulator endpoints are designed to be compatible with the real Cursor Business API. The main differences are that the simulator uses simplified authentication (none required), the simulator returns all events immediately rather than requiring webhook subscriptions, and some metadata fields that contain sensitive information in production are replaced with simulated values.

When migrating from the simulator to the real API, the primary code changes involve adding authentication headers and adjusting for any schema differences in the metadata fields. The core data structures (events, developers, pagination) are intentionally aligned.
