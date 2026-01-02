# Task Breakdown: cursor-sim Phase 3 (Research Framework & Completeness)

## Overview

**Feature**: cursor-sim v2 Phase 3
**Total Estimated Hours**: 35-45
**Number of Steps**: 18
**Current Step**: B02 - Model Usage Generator & Handler

This phase completes cursor-sim with:
- **Part A**: Research Framework (SIM-R013 → SIM-R015) - COMPLETE
- **Part B**: Stub Endpoint Completion (17 endpoints)
- **Part C**: Enhanced Code Quality Analysis
- **Part D**: GitHub Router Wiring (OPTIONAL - documented only)

---

## Progress Tracker

### Part A: Research Framework (15-20h) - COMPLETE

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| A01 | Research Data Models | 2.0 | DONE | 0.25 |
| A02 | Research Dataset Generator | 3.0 | DONE | 0.25 |
| A03 | Parquet/CSV Export | 2.5 | DONE | 0.25 |
| A04 | Research Metrics Service | 3.0 | DONE | 0.25 |
| A05 | Research API Handlers | 2.5 | DONE | 0.25 |
| A06 | Replay Mode Infrastructure | 3.0 | DONE | 0.25 |
| A07 | Part A Integration Tests | 2.0 | DONE | 0.25 |

### Part B: Stub Completion (10-14h) - REVISED

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| B00 | Fix Analytics Response Format | 2.0 | DONE | 1.5 |
| B01 | Update Analytics Data Models | 1.5 | DONE | 1.0 |
| B02 | Model Usage Generator & Handler | 1.5 | IN_PROGRESS | - |
| B03 | Client Version Generator & Handler | 1.0 | NOT_STARTED | - |
| B04 | File Extension Analytics Handler | 1.5 | NOT_STARTED | - |
| B05 | MCP/Commands/Plans/Ask-Mode Handlers | 2.0 | NOT_STARTED | - |
| B06 | Leaderboard Handler | 1.5 | NOT_STARTED | - |
| B07 | By-User Endpoint Handlers | 2.5 | NOT_STARTED | - |
| B08 | Part B Integration Tests | 2.0 | NOT_STARTED | - |

### Part C: Code Quality Analysis (10-15h)

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| C01 | Code Survival Tracking Models | 2.0 | NOT_STARTED | - |
| C02 | Survival Rate Calculator | 3.0 | NOT_STARTED | - |
| C03 | Revert Chain Analysis | 2.5 | NOT_STARTED | - |
| C04 | Hotfix Tracking | 2.0 | NOT_STARTED | - |
| C05 | Part C Integration Tests | 2.0 | NOT_STARTED | - |

### Part D: GitHub Router (OPTIONAL - Documented Only)

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| D01 | Wire GitHub Handlers to Router | 2.0 | OPTIONAL | - |
| D02 | GitHub API E2E Tests | 1.5 | OPTIONAL | - |

---

## Documentation Reference

**IMPORTANT**: All Part B implementations MUST match the Cursor API documentation:

| API | Documentation | Response Format |
|-----|---------------|-----------------|
| Analytics Team-Level | `docs/api-reference/cursor_analytics.md` | `{ data: [...], params: {...} }` |
| Analytics By-User | `docs/api-reference/cursor_analytics.md` | `{ data: { email: [...] }, pagination: {...}, params: {...} }` |
| Admin API | `docs/api-reference/cursor_admin.md` | Varies by endpoint |

**Skill Reference**: `.claude/skills/cursor-api-patterns.md`

---

## Part B: Stub Completion Details (REVISED)

### Step B00: Fix Analytics Response Format

**Priority**: HIGH - Must be done first

**Problem**: Current team-level handlers use `BuildPaginatedResponse` which wraps data in:
```json
{ "data": [...], "pagination": {...}, "params": {...} }
```

But the actual Cursor Analytics API team-level endpoints use:
```json
{ "data": [...], "params": {...} }
```

**Files**:
- `internal/api/response.go` - Add `BuildAnalyticsTeamResponse` function
- `internal/api/cursor/team.go` - Update handlers to use new function

**Tasks**:
- [ ] Create `AnalyticsTeamResponse` struct in `models/response.go`
- [ ] Add `BuildAnalyticsTeamResponse(data, params)` in `api/response.go`
- [ ] Update `teamMetricHandler` to use new function
- [ ] Update tests in `api/response_test.go`

**Cursor API Reference**: `docs/api-reference/cursor_analytics.md` (Team-Level Endpoints section)

**Expected Response Format**:
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

---

### Step B01: Update Analytics Data Models

**File**: `internal/models/team_stats.go`

**Problem**: Current models don't match Cursor API field names and structures.

**Tasks**:
- [ ] Update `AgentEditsDay` to match cursor_analytics.md schema
- [ ] Update `TabCompletionDay` to match cursor_analytics.md schema
- [ ] Update `DAUDay` to match cursor_analytics.md schema (field: `date` not `event_date`)
- [ ] Add `ModelUsageDay` struct with nested `model_breakdown`
- [ ] Add `ClientVersionDay` struct
- [ ] Add `FileExtensionDay` struct
- [ ] Add `MCPUsageDay` struct
- [ ] Add `CommandUsageDay` struct
- [ ] Add `PlanUsageDay` struct
- [ ] Add `AskModeDay` struct
- [ ] Add `LeaderboardResponse` struct (two leaderboards)
- [ ] Update existing handler tests

**Cursor API Reference**: `docs/api-reference/cursor_analytics.md`

**Required Model Schemas** (from cursor_analytics.md):

```go
// AgentEditsDay matches Cursor Analytics API /analytics/team/agent-edits
type AgentEditsDay struct {
    EventDate               string `json:"event_date"`
    TotalSuggestedDiffs     int    `json:"total_suggested_diffs"`
    TotalAcceptedDiffs      int    `json:"total_accepted_diffs"`
    TotalRejectedDiffs      int    `json:"total_rejected_diffs"`
    TotalGreenLinesAccepted int    `json:"total_green_lines_accepted"`
    TotalRedLinesAccepted   int    `json:"total_red_lines_accepted"`
    TotalGreenLinesRejected int    `json:"total_green_lines_rejected"`
    TotalRedLinesRejected   int    `json:"total_red_lines_rejected"`
    TotalGreenLinesSuggested int   `json:"total_green_lines_suggested"`
    TotalRedLinesSuggested  int    `json:"total_red_lines_suggested"`
    TotalLinesSuggested     int    `json:"total_lines_suggested"`
    TotalLinesAccepted      int    `json:"total_lines_accepted"`
}

// TabUsageDay matches Cursor Analytics API /analytics/team/tabs
type TabUsageDay struct {
    EventDate               string `json:"event_date"`
    TotalSuggestions        int    `json:"total_suggestions"`
    TotalAccepts            int    `json:"total_accepts"`
    TotalRejects            int    `json:"total_rejects"`
    TotalGreenLinesAccepted int    `json:"total_green_lines_accepted"`
    TotalRedLinesAccepted   int    `json:"total_red_lines_accepted"`
    TotalGreenLinesRejected int    `json:"total_green_lines_rejected"`
    TotalRedLinesRejected   int    `json:"total_red_lines_rejected"`
    TotalGreenLinesSuggested int   `json:"total_green_lines_suggested"`
    TotalRedLinesSuggested  int    `json:"total_red_lines_suggested"`
    TotalLinesSuggested     int    `json:"total_lines_suggested"`
    TotalLinesAccepted      int    `json:"total_lines_accepted"`
}

// DAUDay matches Cursor Analytics API /analytics/team/dau
type DAUDay struct {
    Date          string `json:"date"`  // NOT event_date!
    DAU           int    `json:"dau"`
    CLIDAU        int    `json:"cli_dau"`
    CloudAgentDAU int    `json:"cloud_agent_dau"`
    BugbotDAU     int    `json:"bugbot_dau"`
}

// ModelUsageDay matches Cursor Analytics API /analytics/team/models
type ModelUsageDay struct {
    Date           string                       `json:"date"`
    ModelBreakdown map[string]ModelBreakdownItem `json:"model_breakdown"`
}

type ModelBreakdownItem struct {
    Messages int `json:"messages"`
    Users    int `json:"users"`
}

// ClientVersionDay matches Cursor Analytics API /analytics/team/client-versions
type ClientVersionDay struct {
    EventDate     string  `json:"event_date"`
    ClientVersion string  `json:"client_version"`
    UserCount     int     `json:"user_count"`
    Percentage    float64 `json:"percentage"`
}

// FileExtensionDay matches Cursor Analytics API /analytics/team/top-file-extensions
type FileExtensionDay struct {
    EventDate           string `json:"event_date"`
    FileExtension       string `json:"file_extension"`
    TotalFiles          int    `json:"total_files"`
    TotalAccepts        int    `json:"total_accepts"`
    TotalRejects        int    `json:"total_rejects"`
    TotalLinesSuggested int    `json:"total_lines_suggested"`
    TotalLinesAccepted  int    `json:"total_lines_accepted"`
    TotalLinesRejected  int    `json:"total_lines_rejected"`
}

// MCPUsageDay matches Cursor Analytics API /analytics/team/mcp
type MCPUsageDay struct {
    EventDate     string `json:"event_date"`
    ToolName      string `json:"tool_name"`
    MCPServerName string `json:"mcp_server_name"`
    Usage         int    `json:"usage"`
}

// CommandUsageDay matches Cursor Analytics API /analytics/team/commands
type CommandUsageDay struct {
    EventDate   string `json:"event_date"`
    CommandName string `json:"command_name"`
    Usage       int    `json:"usage"`
}

// PlanUsageDay matches Cursor Analytics API /analytics/team/plans
type PlanUsageDay struct {
    EventDate string `json:"event_date"`
    Model     string `json:"model"`
    Usage     int    `json:"usage"`
}

// AskModeDay matches Cursor Analytics API /analytics/team/ask-mode
type AskModeDay struct {
    EventDate string `json:"event_date"`
    Model     string `json:"model"`
    Usage     int    `json:"usage"`
}

// LeaderboardResponse matches Cursor Analytics API /analytics/team/leaderboard
type LeaderboardResponse struct {
    TabLeaderboard   LeaderboardSection `json:"tab_leaderboard"`
    AgentLeaderboard LeaderboardSection `json:"agent_leaderboard"`
}

type LeaderboardSection struct {
    Data       []LeaderboardEntry `json:"data"`
    TotalUsers int                `json:"total_users"`
}

type LeaderboardEntry struct {
    Email               string  `json:"email"`
    UserID              string  `json:"user_id"`
    TotalAccepts        int     `json:"total_accepts"`
    TotalLinesAccepted  int     `json:"total_lines_accepted"`
    TotalLinesSuggested int     `json:"total_lines_suggested"`
    LineAcceptanceRatio float64 `json:"line_acceptance_ratio"`
    AcceptRatio         float64 `json:"accept_ratio,omitempty"`
    FavoriteModel       string  `json:"favorite_model,omitempty"`
    Rank                int     `json:"rank"`
}
```

---

### Step B02: Model Usage Generator & Handler

**Files**:
- `internal/generator/model_generator.go` (exists, may need updates)
- `internal/api/cursor/team.go`
- `internal/storage/memory.go`

**Tasks**:
- [ ] Verify `ModelUsageEvent` struct matches needs
- [ ] Update generator to produce realistic model distribution
- [ ] Add `GetModelUsageByTimeRange` to storage interface
- [ ] Implement `TeamModels` handler with correct response format
- [ ] Unit tests

**Models to Generate**: claude-sonnet-4.5, gpt-4o, claude-opus-4, o3, claude-4-sonnet-thinking

**Response Format** (from cursor_analytics.md):
```json
{
  "data": [
    {
      "date": "2025-01-15",
      "model_breakdown": {
        "claude-sonnet-4.5": { "messages": 1250, "users": 28 },
        "gpt-4o": { "messages": 450, "users": 15 }
      }
    }
  ],
  "params": {...}
}
```

---

### Step B03: Client Version Generator & Handler

**Files**:
- `internal/generator/version_generator.go` (new)
- `internal/models/team_stats.go`
- `internal/api/cursor/team.go`
- `internal/storage/memory.go`

**Tasks**:
- [ ] Create `ClientVersionEvent` model
- [ ] Create version generator (realistic semver: 0.42.x, 0.43.x)
- [ ] Add `GetClientVersionsByTimeRange` to storage
- [ ] Implement `TeamClientVersions` handler
- [ ] Unit tests

**Response Format** (from cursor_analytics.md):
```json
{
  "data": [
    {
      "event_date": "2025-01-01",
      "client_version": "0.42.3",
      "user_count": 35,
      "percentage": 0.833
    }
  ],
  "params": {...}
}
```

---

### Step B04: File Extension Analytics Handler

**Files**:
- `internal/api/cursor/team.go`

**Tasks**:
- [ ] Extract file extension from commits (use existing commit data)
- [ ] Aggregate by extension and date
- [ ] Calculate top 5 extensions per day by suggestion volume
- [ ] Implement `TeamTopFileExtensions` handler
- [ ] Unit tests

**Note**: Can derive from existing commit data - check if `RepoName` contains file info or need to generate synthetic extension data.

**Response Format** (from cursor_analytics.md):
```json
{
  "data": [
    {
      "event_date": "2025-01-15",
      "file_extension": "tsx",
      "total_files": 156,
      "total_accepts": 98,
      "total_rejects": 45,
      "total_lines_suggested": 3230,
      "total_lines_accepted": 2340,
      "total_lines_rejected": 890
    }
  ],
  "params": {...}
}
```

---

### Step B05: MCP/Commands/Plans/Ask-Mode Handlers

**Files**:
- `internal/generator/feature_generators.go` (new)
- `internal/models/team_stats.go`
- `internal/api/cursor/team.go`
- `internal/storage/memory.go`

**Tasks**:
- [ ] Create event models for each feature type
- [ ] Create generators for each feature (based on developer seed preferences)
- [ ] Add storage methods for each
- [ ] Implement 4 handlers: `TeamMCP`, `TeamCommands`, `TeamPlans`, `TeamAskMode`
- [ ] Unit tests for each

**MCP Tools to Generate**: read_file, search_web, execute_command, etc.
**Commands to Generate**: explain, refactor, fix, test, etc.
**Plans/AskMode**: Use model preferences from seed

---

### Step B06: Leaderboard Handler

**File**: `internal/api/cursor/team.go`

**Tasks**:
- [ ] Calculate tab leaderboard from commits (TabLinesAdded, acceptance rates)
- [ ] Calculate agent leaderboard from commits (ComposerLinesAdded, acceptance rates)
- [ ] Support pagination (`page`, `pageSize` params)
- [ ] Support user filtering (`users` param)
- [ ] Implement ranking logic
- [ ] Implement `TeamLeaderboard` handler
- [ ] Unit tests

**Response Format** (from cursor_analytics.md):
```json
{
  "data": {
    "tab_leaderboard": {
      "data": [
        {
          "email": "alice@example.com",
          "user_id": "user_abc123",
          "total_accepts": 1334,
          "total_lines_accepted": 3455,
          "total_lines_suggested": 15307,
          "line_acceptance_ratio": 0.226,
          "accept_ratio": 0.233,
          "rank": 1
        }
      ],
      "total_users": 142
    },
    "agent_leaderboard": {
      "data": [...],
      "total_users": 142
    }
  },
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "totalUsers": 142,
    "totalPages": 15,
    "hasNextPage": true,
    "hasPreviousPage": false
  },
  "params": {...}
}
```

---

### Step B07: By-User Endpoint Handlers

**File**: `internal/api/cursor/byuser.go`

**Tasks**:
- [ ] Create `BuildAnalyticsByUserResponse` helper function
- [ ] Implement all 9 by-user handlers with real data
- [ ] Group data by user email
- [ ] Support pagination on users (not on data)
- [ ] Support `users` filter parameter (comma-separated)
- [ ] Include `userMappings` in params
- [ ] Unit tests for each

**Endpoints**:
1. `/analytics/by-user/agent-edits`
2. `/analytics/by-user/tabs`
3. `/analytics/by-user/models`
4. `/analytics/by-user/client-versions`
5. `/analytics/by-user/top-file-extensions`
6. `/analytics/by-user/mcp`
7. `/analytics/by-user/commands`
8. `/analytics/by-user/plans`
9. `/analytics/by-user/ask-mode`

**Response Format** (from cursor_analytics.md):
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

---

### Step B08: Part B Integration Tests

**File**: `test/e2e/analytics_complete_test.go`

**Tasks**:
- [ ] E2E: All team analytics endpoints return correct schema
- [ ] E2E: All by-user analytics endpoints return correct schema
- [ ] E2E: Leaderboard pagination works
- [ ] E2E: User filtering works on by-user endpoints
- [ ] E2E: Date filtering works
- [ ] Verify response schemas match `docs/api-reference/cursor_analytics.md`
- [ ] No more stub responses (all return real data)

---

## Dependency Graph (Updated)

```
PART B: Stub Completion (REVISED)

B00 (Fix Response Format)
 │
 └── B01 (Update Models)
      │
      ├── B02 (Model Usage) ───┐
      ├── B03 (Client Version) │
      ├── B04 (File Extension) ├──► B07 (By-User Handlers)
      ├── B05 (MCP/Cmd/Plans)  │         │
      └── B06 (Leaderboard) ───┘         │
                                         └── B08 (Tests)
```

---

## Part C: Code Quality Analysis Details

(Unchanged from original - see below)

### Step C01: Code Survival Tracking Models
...

### Step C02: Survival Rate Calculator
...

### Step C03: Revert Chain Analysis
...

### Step C04: Hotfix Tracking
...

### Step C05: Part C Integration Tests
...

---

## Part D: GitHub Router (OPTIONAL)

(Unchanged - OPTIONAL)

---

## Model Recommendations

| Step | Model | Rationale |
|------|-------|-----------|
| B00, B01 | Haiku | Schema alignment, straightforward |
| B02, B03 | Haiku | Simple generators |
| B04, B05 | Haiku | Data extraction from commits |
| B06 | Sonnet | Complex ranking algorithm |
| B07 | Sonnet | Multiple similar handlers |
| B08 | Sonnet | Integration testing |
| C01-C04 | Sonnet | Algorithm complexity |

---

## TDD Checklist (Per Step)

- [ ] Read step details and acceptance criteria
- [ ] Read relevant section in `docs/api-reference/cursor_analytics.md`
- [ ] Write failing test(s) for the step
- [ ] Run tests, confirm RED
- [ ] Implement minimal code to pass
- [ ] Run tests, confirm GREEN
- [ ] Refactor while green
- [ ] Run linter
- [ ] Update step status to DONE
- [ ] Commit with step reference

---

## Execution Order

```
1. B00 → B01 → B02 → B03 → B04 → B05 → B06 → B07 → B08 (commit: "feat: Phase 3 Part B")
2. C01 → C02 → C03 → C04 → C05 (commit: "feat: Phase 3 Part C")
3. D01 → D02 (OPTIONAL - separate commit if done)
```
