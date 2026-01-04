---
name: api-contract
description: cursor-sim API contract reference. Use when implementing services that consume cursor-sim data (analytics-core, viz-spa), ensuring API alignment, or verifying response formats. Source of truth for data models and endpoint contracts. (project)
---

# cursor-sim API Contract

This skill provides the API contract for cursor-sim, the data source for all analytics services.

**Source of Truth**: `services/cursor-sim/SPEC.md`

## Quick Reference

### Base Configuration

- **Port**: 8080
- **Auth**: Basic Authentication (API key as username)
- **Format**: JSON

### Endpoint Groups

| Group | Prefix | Description |
|-------|--------|-------------|
| Health | `/health` | Service health check |
| Admin | `/admin/*` | Team and member management |
| Cursor Analytics | `/analytics/*` | AI code tracking metrics |
| GitHub | `/repos/*` | Repository and PR data |
| Research | `/research/*` | Research dataset export |

## Response Formats

### Team-Level Analytics Response

```json
{
  "data": { ... },
  "params": {
    "start_date": "2025-01-01",
    "end_date": "2025-01-31"
  }
}
```

### By-User Analytics Response

```json
{
  "data": {
    "user-id-1": { ... },
    "user-id-2": { ... }
  },
  "pagination": {
    "total": 100,
    "page": 1,
    "page_size": 20
  },
  "params": { ... }
}
```

### Paginated List Response

```json
{
  "data": [ ... ],
  "pagination": {
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

## Key Endpoints

### Admin Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/teams/members` | GET | List team members |
| `/admin/users` | GET | List all users (paginated) |
| `/admin/users/{id}` | GET | Get user by ID |

### Analytics Endpoints (Team)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/analytics/team/overview` | GET | Overall AI usage metrics |
| `/analytics/team/total-usage` | GET | Total AI code statistics |
| `/analytics/team/models` | GET | Model usage breakdown |
| `/analytics/team/client-versions` | GET | Client version distribution |
| `/analytics/team/top-file-extensions` | GET | File type breakdown |
| `/analytics/team/leaderboard` | GET | Developer rankings |
| `/analytics/team/mcp` | GET | MCP tool usage |
| `/analytics/team/commands` | GET | Command usage stats |
| `/analytics/team/plans` | GET | Plan usage stats |
| `/analytics/team/ask-mode` | GET | Ask mode usage |

### Analytics Endpoints (By-User)

Same as team endpoints but with `/by-user/` prefix:
- `/analytics/by-user/overview`
- `/analytics/by-user/models`
- etc.

### GitHub Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/repos` | GET | List repositories |
| `/repos/{owner}/{repo}` | GET | Get repository |
| `/repos/{owner}/{repo}/commits` | GET | List commits |
| `/repos/{owner}/{repo}/commits/{sha}` | GET | Get commit |
| `/repos/{owner}/{repo}/pulls` | GET | List pull requests |
| `/repos/{owner}/{repo}/pulls/{number}` | GET | Get PR details |
| `/repos/{owner}/{repo}/pulls/{number}/commits` | GET | PR commits |
| `/repos/{owner}/{repo}/pulls/{number}/files` | GET | PR files |
| `/repos/{owner}/{repo}/pulls/{number}/reviews` | GET | PR reviews |
| `/repos/{owner}/{repo}/analysis/survival` | GET | Code survival |
| `/repos/{owner}/{repo}/analysis/reverts` | GET | Revert chains |
| `/repos/{owner}/{repo}/analysis/hotfixes` | GET | Hotfix patterns |

### Research Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/research/dataset` | GET | Full research dataset (CSV/JSON) |
| `/research/metrics/velocity` | GET | Velocity by AI band |
| `/research/metrics/review-costs` | GET | Review costs by AI band |
| `/research/metrics/quality` | GET | Quality by AI band |

## Common Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `start_date` | string | ISO date (YYYY-MM-DD) |
| `end_date` | string | ISO date (YYYY-MM-DD) |
| `page` | int | Page number (1-based) |
| `page_size` | int | Items per page (default: 20) |
| `format` | string | Response format (csv, json) |

## Data Models

### Developer

```typescript
interface Developer {
  user_id: string;
  name: string;
  email: string;
  seniority: "junior" | "mid" | "senior";
  ai_preference: number; // 0.0-1.0
  active: boolean;
}
```

### Commit

```typescript
interface Commit {
  sha: string;
  message: string;
  author_id: string;
  timestamp: string; // ISO 8601
  ai_lines_added: number;
  ai_lines_deleted: number;
  human_lines_added: number;
  human_lines_deleted: number;
  ai_ratio: number; // 0.0-1.0
  files: CommitFile[];
}
```

### Pull Request

```typescript
interface PullRequest {
  number: number;
  title: string;
  state: "open" | "closed" | "merged";
  author_id: string;
  created_at: string;
  merged_at?: string;
  ai_lines_total: number;
  human_lines_total: number;
  ai_ratio: number;
  commit_count: number;
  file_count: number;
}
```

## Integration Guide

### For analytics-core (GraphQL aggregator)

1. Fetch data from cursor-sim REST endpoints
2. Aggregate across time periods
3. Calculate derived metrics
4. Expose via GraphQL schema

### For viz-spa (React dashboard)

1. Connect to analytics-core GraphQL
2. Display visualizations of cursor-sim data
3. Support filtering by date range, team, user

## Full Specification

For complete details, read the full specification:

```bash
cat services/cursor-sim/SPEC.md
```

This includes:
- All endpoint details with request/response schemas
- Data model definitions
- Authentication requirements
- Rate limiting
- Error responses
