# SPEC.md - Cursor API Simulator (cursor-sim)

**Service**: cursor-sim
**Type**: CLI Tool + REST API Server + GraphQL API
**Language**: Go 1.21+
**Port**: 8080 (configurable)
**Version**: 2.0.0
**Last Updated**: January 2026

---

## Overview

The Cursor API Simulator is a mock server that generates synthetic usage data **mimicking the actual Cursor Business API**. It enables development and testing of the analytics dashboard without requiring access to real Cursor data or an enterprise Cursor subscription.

The simulator provides:
- **Cursor-compatible REST endpoints** matching the official Analytics API and AI Code Tracking API
- **GraphQL API** for flexible querying with filters
- **In-memory database** with query support
- **Interactive CLI controls** for managing data generation
- **JSON configuration** for flexible setup

---

## Purpose

This service enables:
1. **Dashboard Development**: Frontend developers can build against realistic API responses
2. **Integration Testing**: Test the full pipeline without Cursor enterprise access
3. **Demo Environments**: Showcase the analytics platform with synthetic data
4. **Load Testing**: Generate high volumes of data to test performance

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         cursor-sim                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────────┐     ┌──────────────────────────────────────┐  │
│  │  JSON Config    │────▶│  Configuration Manager               │  │
│  │  (stdin/file)   │     │  - Validates input                   │  │
│  └─────────────────┘     │  - Applies defaults                  │  │
│                          └──────────────────────────────────────┘  │
│                                         │                           │
│                                         ▼                           │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │                    Data Generator Engine                      │  │
│  │  ┌─────────────────────┐  ┌─────────────────────────────┐   │  │
│  │  │ Developer Generator │  │ Commit/Change Generator     │   │  │
│  │  │ - Regions           │  │ - Poisson timing            │   │  │
│  │  │ - Divisions         │  │ - TAB vs COMPOSER           │   │  │
│  │  │ - Groups            │  │ - AI model rotation         │   │  │
│  │  │ - Teams             │  │ - Line count simulation     │   │  │
│  │  └─────────────────────┘  └─────────────────────────────┘   │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                         │                           │
│                                         ▼                           │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │                   In-Memory Database                          │  │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌─────────────┐  │  │
│  │  │Developers │ │ Commits   │ │ Changes   │ │Daily Metrics│  │  │
│  │  └───────────┘ └───────────┘ └───────────┘ └─────────────┘  │  │
│  │                                                               │  │
│  │  Indexes: by Region, Division, Group, Team, Date             │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                          │                    │                     │
│            ┌─────────────┴──────┐    ┌───────┴───────┐             │
│            ▼                    ▼    ▼               ▼             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐    │
│  │   REST API      │  │  GraphQL API    │  │  CLI Dashboard  │    │
│  │   (Cursor-like) │  │  (Flexible)     │  │  (Interactive)  │    │
│  │   Port: 8080    │  │  /graphql       │  │  Stats + Ctrl   │    │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘    │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## JSON Configuration

The simulator accepts JSON configuration via file or stdin.

### Configuration Schema

```json
{
  "auth": {
    "api_key": "cursor_sim_key_abc123",
    "api_secret": "cursor_sim_secret_xyz789",
    "team_id": "team_demo_001"
  },
  "organization": {
    "developers": 100,
    "regions": {
      "US": 0.50,
      "EU": 0.35,
      "APAC": 0.15
    },
    "divisions": {
      "AGS": 0.40,
      "AT": 0.35,
      "ST": 0.25
    },
    "groups": {
      "TMOBILE": 0.60,
      "ATANT": 0.40
    },
    "teams": {
      "dev": 0.75,
      "support": 0.25
    }
  },
  "generation": {
    "velocity": "high",
    "volatility": 0.25,
    "seed": 12345,
    "ai_models": ["claude-3.5-sonnet", "claude-opus-4", "gpt-4-turbo"],
    "tab_vs_composer_ratio": 0.70
  },
  "break_condition": {
    "type": "pr_count",
    "value": 100000
  },
  "server": {
    "port": 8080,
    "host": "0.0.0.0",
    "enable_graphql": true
  },
  "export": {
    "format": "json",
    "path": "./exports"
  }
}
```

### Configuration Parameters

#### auth
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `api_key` | string | Yes | Fake API key for Basic Auth simulation |
| `api_secret` | string | Yes | Fake API secret for Basic Auth |
| `team_id` | string | Yes | Team identifier for API responses |

#### organization
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `developers` | int | 50 | Number of developers to generate (1-10000) |
| `regions` | object | {US:0.5, EU:0.35, APAC:0.15} | Region distribution (must sum to 1.0) |
| `divisions` | object | {AGS:0.4, AT:0.35, ST:0.25} | Division distribution |
| `groups` | object | {TMOBILE:0.6, ATANT:0.4} | Group distribution |
| `teams` | object | {dev:0.75, support:0.25} | Team distribution |

#### generation
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `velocity` | string | "high" | Event rate: "low" (10/hr), "medium" (50/hr), "high" (100/hr) |
| `volatility` | float | 0.2 | Per-developer variance (0.0-1.0) |
| `seed` | int64 | timestamp | Random seed for reproducibility |
| `ai_models` | array | [...] | AI models to rotate through |
| `tab_vs_composer_ratio` | float | 0.7 | Ratio of TAB to COMPOSER suggestions |

#### break_condition
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `type` | string | "none" | "pr_count", "duration", "none" |
| `value` | int | 100000 | Threshold for break (PR count or seconds) |

---

## CLI Interface

### Starting the Simulator

```bash
# With config file
./cursor-sim --config config.json

# With config via stdin
cat config.json | ./cursor-sim --config -

# With minimal flags (uses defaults)
./cursor-sim --developers 100 --velocity high

# Show help
./cursor-sim --help
```

### CLI Flags

```
Usage: cursor-sim [options]

Configuration:
  --config string      Path to JSON config file (use "-" for stdin)
  --developers int     Number of developers (default: 50)
  --velocity string    Event rate: low|medium|high (default: high)
  --volatility float   Per-developer variance 0.0-1.0 (default: 0.2)
  --seed int64         Random seed for reproducibility
  --port int           HTTP server port (default: 8080)

Flags override config file values.
```

### Interactive Controls

While the simulator is running, the following keyboard controls are available:

| Key | Action | Description |
|-----|--------|-------------|
| `Ctrl+S` | Soft Stop | Pause data generation, keep API running |
| `Ctrl+R` | Resume | Resume data generation after soft stop |
| `Ctrl+E` | Export | Dump in-memory DB to JSON file |
| `Ctrl+Q` | Quick Stats | Show current statistics snapshot |
| `Ctrl+C` | Shutdown | Graceful shutdown (stop API, exit) |

### CLI Dashboard Display

```
╔══════════════════════════════════════════════════════════════════════╗
║            Cursor API Simulator v2.0.0 - RUNNING                     ║
╠══════════════════════════════════════════════════════════════════════╣
║ Status: GENERATING          Uptime: 00:15:32                         ║
║ API: http://localhost:8080  GraphQL: http://localhost:8080/graphql   ║
╠══════════════════════════════════════════════════════════════════════╣
║ DEVELOPERS                                                           ║
║   Total: 142                                                         ║
║   ├── US: 71 (50%)  EU: 50 (35%)  APAC: 21 (15%)                   ║
║   ├── AGS: 57  AT: 50  ST: 35                                       ║
║   └── TMOBILE: 85  ATANT: 57                                        ║
╠══════════════════════════════════════════════════════════════════════╣
║ COMMITS & CHANGES                                                    ║
║   Total Commits:     12,847 cumulative                              ║
║   Total Changes:     28,391 cumulative                              ║
║   Current Rate:      156 commits/min                                ║
║   ├── From TAB:      19,874 (70%)                                   ║
║   └── From COMPOSER: 8,517 (30%)                                    ║
╠══════════════════════════════════════════════════════════════════════╣
║ BREAK CONDITION                                                      ║
║   Type: PR Count                                                     ║
║   Progress: 12,847 / 100,000 (12.8%)                                ║
║   ████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 12.8%  ║
╠══════════════════════════════════════════════════════════════════════╣
║ Controls: [S]oft Stop | [R]esume | [E]xport | [Q]uick Stats | Ctrl+C ║
╚══════════════════════════════════════════════════════════════════════╝
```

---

## Data Models

### Developer

```go
type Developer struct {
    ID             string    `json:"id"`              // UUID: "user_abc123"
    Email          string    `json:"email"`           // "jane.smith@company.com"
    Name           string    `json:"name"`            // "Jane Smith"
    Region         string    `json:"region"`          // "US", "EU", "APAC"
    Division       string    `json:"division"`        // "AGS", "AT", "ST"
    Group          string    `json:"group"`           // "TMOBILE", "ATANT"
    Team           string    `json:"team"`            // "dev", "support"
    Seniority      string    `json:"seniority"`       // "junior", "mid", "senior"
    ClientVersion  string    `json:"client_version"`  // "0.43.6"
    AcceptanceRate float64   `json:"acceptance_rate"` // 0.0-1.0
    IsActive       bool      `json:"is_active"`
    CreatedAt      time.Time `json:"created_at"`
    LastActiveAt   time.Time `json:"last_active_at"`
}
```

### Commit (AI Code Tracking)

```go
type Commit struct {
    Hash              string    `json:"commit_hash"`
    Timestamp         time.Time `json:"timestamp"`
    Message           string    `json:"message"`
    UserID            string    `json:"user_id"`
    UserEmail         string    `json:"user_email"`
    Repository        string    `json:"repository"`
    Branch            string    `json:"branch"`
    TotalLines        int       `json:"total_lines"`
    LinesFromTAB      int       `json:"lines_from_tab"`
    LinesFromComposer int       `json:"lines_from_composer"`
    LinesNonAI        int       `json:"lines_non_ai"`
    IngestionTime     time.Time `json:"ingestion_time"`
}
```

### Change (AI Code Tracking)

```go
type Change struct {
    ChangeID      string    `json:"change_id"`      // Deterministic ID
    CommitHash    string    `json:"commit_hash"`
    UserID        string    `json:"user_id"`
    Timestamp     time.Time `json:"timestamp"`
    Source        string    `json:"source"`         // "TAB" or "COMPOSER"
    Model         string    `json:"model"`          // "claude-3.5-sonnet"
    FilePath      string    `json:"file_path"`
    FileExtension string    `json:"file_extension"` // ".ts", ".go", ".py"
    LinesAdded    int       `json:"lines_added"`
    LinesRemoved  int       `json:"lines_removed"`
    IngestionTime time.Time `json:"ingestion_time"`
}
```

### Daily Metrics (Team Analytics)

```go
type DailyMetrics struct {
    Date              time.Time      `json:"date"`
    DAU               int            `json:"dau"`
    AgentEdits        int            `json:"agent_edits"`
    TabCompletions    int            `json:"tab_completions"`
    ComposerEdits     int            `json:"composer_edits"`
    TotalLinesAdded   int            `json:"total_lines_added"`
    ModelUsage        map[string]int `json:"model_usage"`
    TopFileExtensions map[string]int `json:"top_file_extensions"`
    CommandUsage      map[string]int `json:"command_usage"`
}
```

---

## REST API Endpoints (Cursor-Compatible)

All endpoints require **Basic Authentication** with the configured `api_key` and `api_secret`.

```
Authorization: Basic {base64(api_key:api_secret)}
```

### Authentication

**Invalid credentials return:**
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid API credentials"
  }
}
```

---

### AI Code Tracking API

#### GET /v1/analytics/ai-code/commits

Returns paginated commit metrics with AI contribution breakdown.

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `startDate` | string | No | Start date (ISO 8601, "7d", "30d", "today") |
| `endDate` | string | No | End date (ISO 8601, "today", "now") |
| `user` | string | No | Filter by user email or ID |
| `page` | int | No | Page number (default: 1) |
| `pageSize` | int | No | Results per page (default: 100, max: 1000) |

**Response: 200 OK**
```json
{
  "data": [
    {
      "commit_hash": "abc123def456...",
      "timestamp": "2026-01-15T10:30:00Z",
      "message": "feat: add user authentication",
      "user_id": "user_abc123",
      "user_email": "jane.smith@company.com",
      "repository": "frontend-app",
      "branch": "feature/auth",
      "total_lines": 145,
      "lines_from_tab": 87,
      "lines_from_composer": 35,
      "lines_non_ai": 23,
      "ingestion_time": "2026-01-15T10:31:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 100,
    "total": 5847,
    "hasMore": true
  },
  "params": {
    "startDate": "2026-01-08",
    "endDate": "2026-01-15",
    "user": null
  }
}
```

#### GET /v1/analytics/ai-code/commits.csv

Returns commit data as CSV stream for bulk export.

**Response: 200 OK (text/csv)**
```csv
commit_hash,timestamp,user_email,repository,total_lines,lines_from_tab,lines_from_composer,lines_non_ai
abc123def456,2026-01-15T10:30:00Z,jane.smith@company.com,frontend-app,145,87,35,23
```

#### GET /v1/analytics/ai-code/changes

Returns paginated AI-suggested changes (individual accepted suggestions).

**Query Parameters:** Same as /commits

**Response: 200 OK**
```json
{
  "data": [
    {
      "change_id": "chg_abc123",
      "commit_hash": "abc123def456...",
      "user_id": "user_abc123",
      "timestamp": "2026-01-15T10:30:00Z",
      "source": "TAB",
      "model": "claude-3.5-sonnet",
      "file_path": "src/auth/login.ts",
      "file_extension": ".ts",
      "lines_added": 12,
      "lines_removed": 3,
      "ingestion_time": "2026-01-15T10:31:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 100,
    "total": 28391
  }
}
```

#### GET /v1/analytics/ai-code/changes.csv

Returns changes data as CSV stream.

---

### Team Analytics API

#### GET /v1/analytics/team/agent-edits

Returns team-wide AI-suggested code edits accepted.

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `startDate` | string | No | Start date |
| `endDate` | string | No | End date |
| `page` | int | No | Page number |
| `pageSize` | int | No | Results per page |

**Response: 200 OK**
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

#### GET /v1/analytics/team/tabs

Returns TAB autocomplete metrics.

**Response: 200 OK**
```json
{
  "data": [
    {
      "date": "2026-01-15",
      "completions_shown": 8745,
      "completions_accepted": 6121,
      "acceptance_rate": 70.0,
      "unique_users": 48
    }
  ]
}
```

#### GET /v1/analytics/team/dau

Returns daily active users count.

**Response: 200 OK**
```json
{
  "data": [
    {
      "date": "2026-01-15",
      "dau": 48,
      "total_users": 142
    }
  ]
}
```

#### GET /v1/analytics/team/models

Returns AI model usage distribution.

**Response: 200 OK**
```json
{
  "data": [
    {
      "date": "2026-01-15",
      "model_usage": {
        "claude-3.5-sonnet": 4521,
        "claude-opus-4": 2134,
        "gpt-4-turbo": 1892
      }
    }
  ]
}
```

#### GET /v1/analytics/team/client-versions

Returns Cursor IDE version distribution.

**Response: 200 OK**
```json
{
  "data": {
    "versions": {
      "0.43.6": 45,
      "0.43.5": 32,
      "0.42.0": 15
    },
    "total_users": 92
  }
}
```

#### GET /v1/analytics/team/top-file-extensions

Returns most-edited file types.

**Response: 200 OK**
```json
{
  "data": {
    "extensions": {
      ".ts": 3421,
      ".tsx": 2891,
      ".py": 1234,
      ".go": 987,
      ".js": 654
    }
  }
}
```

#### GET /v1/analytics/team/leaderboard

Returns users ranked by specified metric.

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `metric` | string | Yes | "agent_edits", "tab_completions", "lines_added" |
| `limit` | int | No | Number of results (default: 10) |

**Response: 200 OK**
```json
{
  "data": [
    {
      "rank": 1,
      "user_id": "user_abc123",
      "user_email": "jane.smith@company.com",
      "metric_value": 1247
    },
    {
      "rank": 2,
      "user_id": "user_def456",
      "user_email": "john.doe@company.com",
      "metric_value": 1189
    }
  ]
}
```

---

### By-User Analytics API

All team endpoints have by-user variants that return per-user breakdowns:

- `GET /v1/analytics/by-user/agent-edits`
- `GET /v1/analytics/by-user/tabs`
- `GET /v1/analytics/by-user/models`
- `GET /v1/analytics/by-user/top-file-extensions`

**Additional Query Parameter:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `users` | string | Comma-separated emails/IDs |

---

### Health & Info

#### GET /v1/health

Returns service health status.

**Response: 200 OK**
```json
{
  "status": "healthy",
  "version": "2.0.0",
  "uptime": "02:30:15",
  "generation": {
    "status": "running",
    "developers": 142,
    "commits_generated": 12847,
    "changes_generated": 28391
  },
  "break_condition": {
    "type": "pr_count",
    "current": 12847,
    "target": 100000,
    "progress_percent": 12.8
  }
}
```

---

## GraphQL API

The simulator also exposes a GraphQL API at `/graphql` for flexible querying.

### Schema

```graphql
type Query {
  # Developers
  developer(id: ID!): Developer
  developers(
    region: [String]
    division: [String]
    group: [String]
    team: [String]
    limit: Int
    offset: Int
  ): DeveloperConnection!

  # Commits & Changes
  commits(
    startDate: DateTime
    endDate: DateTime
    userId: String
    limit: Int
    offset: Int
  ): CommitConnection!

  changes(
    startDate: DateTime
    endDate: DateTime
    userId: String
    source: String
    limit: Int
    offset: Int
  ): ChangeConnection!

  # Metrics
  dailyMetrics(
    startDate: DateTime
    endDate: DateTime
  ): [DailyMetrics!]!

  teamStats(
    startDate: DateTime
    endDate: DateTime
  ): TeamStats!

  leaderboard(
    metric: String!
    startDate: DateTime
    endDate: DateTime
    limit: Int
  ): [LeaderboardEntry!]!

  # Health
  health: HealthStatus!
}

type Developer {
  id: ID!
  email: String!
  name: String!
  region: String!
  division: String!
  group: String!
  team: String!
  seniority: String!
  clientVersion: String!
  isActive: Boolean!
  createdAt: DateTime!
  lastActiveAt: DateTime!

  # Nested metrics
  stats(startDate: DateTime, endDate: DateTime): DeveloperStats
  commits(limit: Int): [Commit!]!
}

type DeveloperStats {
  totalCommits: Int!
  totalChanges: Int!
  linesFromTab: Int!
  linesFromComposer: Int!
  acceptanceRate: Float
  topModels: [ModelUsage!]!
  topExtensions: [ExtensionUsage!]!
}

type Commit {
  hash: String!
  timestamp: DateTime!
  message: String!
  userId: String!
  repository: String!
  branch: String!
  totalLines: Int!
  linesFromTab: Int!
  linesFromComposer: Int!
  linesNonAI: Int!
}

type Change {
  changeId: String!
  commitHash: String!
  userId: String!
  timestamp: DateTime!
  source: String!
  model: String!
  filePath: String!
  fileExtension: String!
  linesAdded: Int!
  linesRemoved: Int!
}

type DailyMetrics {
  date: DateTime!
  dau: Int!
  agentEdits: Int!
  tabCompletions: Int!
  composerEdits: Int!
  modelUsage: [ModelUsage!]!
  topFileExtensions: [ExtensionUsage!]!
}

type TeamStats {
  totalDevelopers: Int!
  activeDevelopers: Int!
  totalCommits: Int!
  totalChanges: Int!
  averageAcceptanceRate: Float
  regionBreakdown: [RegionStats!]!
  divisionBreakdown: [DivisionStats!]!
}

type LeaderboardEntry {
  rank: Int!
  developer: Developer!
  metricValue: Int!
}
```

---

## Data Generation Logic

### Developer Generation

On startup, developers are generated with the following logic:

1. **Count**: Use `organization.developers` from config
2. **Distribution**: Apply configured percentages for Region, Division, Group, Team
3. **Seniority**: 20% junior, 50% mid, 30% senior
4. **Acceptance Rates**:
   - Junior: 55-65%
   - Mid: 70-80%
   - Senior: 85-95%
5. **Client Versions**: Random from recent versions (0.43.x, 0.42.x)
6. **Names**: Deterministic from seed, realistic first/last combinations

### Commit & Change Generation

Events are generated continuously in background goroutines:

1. **Poisson Timing**: Events follow Poisson distribution based on velocity
2. **Per-Developer Variance**: Each developer's rate varies by volatility
3. **TAB vs COMPOSER**: Use configured `tab_vs_composer_ratio`
4. **AI Models**: Rotate through configured `ai_models`
5. **File Extensions**: Weighted random (.ts: 25%, .tsx: 20%, .py: 15%, etc.)
6. **Line Counts**: Realistic distribution (1-50 lines per change)

### Velocity Settings

| Level | Commits/Hour/Dev | Description |
|-------|------------------|-------------|
| low | ~5 | Light usage, occasional AI assistance |
| medium | ~25 | Regular usage, frequent suggestions |
| high | ~50 | Heavy usage, continuous AI assistance |

---

## Data Export

### Export Command (Ctrl+E)

When triggered, exports the entire in-memory database to a file.

**JSON Export** (default):
```json
{
  "export_timestamp": "2026-01-15T14:30:00Z",
  "config": { ... },
  "statistics": {
    "total_developers": 142,
    "total_commits": 12847,
    "total_changes": 28391,
    "generation_duration": "02:30:15"
  },
  "data": {
    "developers": [...],
    "commits": [...],
    "changes": [...],
    "daily_metrics": [...]
  }
}
```

**Binary Export** (optional, for large datasets):
Uses Go's gob encoding for efficient serialization.

---

## Rate Limiting

The simulator implements Cursor-compatible rate limits:

| Endpoint Type | Limit |
|---------------|-------|
| Team endpoints | 100 requests/minute |
| By-user endpoints | 50 requests/minute |

Exceeding limits returns:
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Retry after 60 seconds.",
    "retry_after": 60
  }
}
```

**HTTP Status**: 429 Too Many Requests

---

## Error Handling

### Error Response Format

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message",
    "details": { ... }
  }
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Invalid API credentials |
| `INVALID_DATE_RANGE` | 400 | Invalid date parameters |
| `INVALID_PARAMETER` | 400 | Invalid query parameter |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

---

## Testing Requirements

### Unit Tests

- [ ] JSON config parsing and validation
- [ ] Developer generation with correct distributions
- [ ] Commit/change generation with Poisson timing
- [ ] REST API response schemas match Cursor format
- [ ] GraphQL resolvers return correct data
- [ ] Rate limiting enforcement
- [ ] Basic Auth validation
- [ ] Date range parsing (ISO 8601, relative dates)

### Integration Tests

- [ ] Full config → generation → API flow
- [ ] Soft stop (Ctrl+S) pauses generation, API continues
- [ ] Export (Ctrl+E) produces valid JSON
- [ ] Graceful shutdown (Ctrl+C)
- [ ] Memory stability under load

### Coverage Target

- **Minimum**: 80% line coverage
- **Critical paths**: 95% coverage (auth, rate limiting, data generation)

---

## Dependencies

```go
// go.mod
module github.com/org/cursor-analytics-platform/services/cursor-sim

go 1.21

require (
    github.com/gorilla/mux v1.8.1
    github.com/99designs/gqlgen v0.17.45
    github.com/joho/godotenv v1.5.1
    github.com/google/uuid v1.6.0
)
```

---

## File Structure

```
services/cursor-sim/
├── cmd/
│   └── simulator/
│       └── main.go                # Entry point
├── internal/
│   ├── config/
│   │   ├── loader.go              # JSON config loading
│   │   ├── validator.go           # Config validation
│   │   └── defaults.go            # Default values
│   ├── models/
│   │   ├── developer.go
│   │   ├── commit.go
│   │   ├── change.go
│   │   └── metrics.go
│   ├── generator/
│   │   ├── developer_gen.go       # Developer creation
│   │   ├── commit_gen.go          # Commit generation
│   │   ├── change_gen.go          # Change generation
│   │   └── poisson.go             # Timing distribution
│   ├── db/
│   │   ├── store.go               # Interface
│   │   ├── memory_store.go        # In-memory impl
│   │   └── indexer.go             # Query indexes
│   ├── api/
│   │   ├── server.go              # HTTP server
│   │   ├── auth.go                # Basic Auth middleware
│   │   ├── ratelimit.go           # Rate limiting
│   │   └── handlers/
│   │       ├── ai_code.go         # /analytics/ai-code/*
│   │       ├── team.go            # /analytics/team/*
│   │       └── by_user.go         # /analytics/by-user/*
│   ├── graphql/
│   │   ├── schema.graphqls
│   │   ├── resolver.go
│   │   └── generated.go
│   ├── cli/
│   │   ├── dashboard.go           # Terminal UI
│   │   ├── signals.go             # Ctrl+S/E/C handlers
│   │   └── stats.go               # Statistics tracking
│   └── export/
│       ├── json_exporter.go
│       └── binary_exporter.go
├── go.mod
├── go.sum
├── Dockerfile
├── SPEC.md                        # This file
└── README.md
```

---

## Appendix: Cursor API Reference

Based on official Cursor documentation:

- **AI Code Tracking API**: https://cursor.com/docs/account/teams/ai-code-tracking-api
- **Analytics API**: https://cursor.com/docs/account/teams/analytics-api

The simulator aims to provide 100% API compatibility with these endpoints.

