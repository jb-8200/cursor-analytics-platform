# Task Breakdown: cursor-sim Phase 3 (Research Framework & Completeness)

## Overview

**Feature**: cursor-sim v2 Phase 3
**Total Estimated Hours**: 35-45
**Number of Steps**: 18
**Current Step**: C03 - Revert Chain Analysis

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
| B02 | Model Usage Generator & Handler | 1.5 | DONE | 1.5 |
| B03 | Client Version Generator & Handler | 1.0 | DONE | 1.0 |
| B04 | File Extension Analytics Handler | 1.5 | DONE | 1.2 |
| B05 | MCP/Commands/Plans/Ask-Mode Handlers | 2.0 | DONE | 1.5 |
| B06 | Leaderboard Handler | 1.5 | DONE | 1.2 |
| B07 | By-User Endpoint Handlers | 2.5 | DONE | 2.0 |
| B08 | Part B Integration Tests | 2.0 | DONE | 1.5 |

### Part C: GitHub Simulation + Quality Analysis (18.5h) - REVISED

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| C00 | PR Generation Pipeline | 4.0 | DONE | 4.0 |
| C01 | Wire GitHub/Research Routes | 2.0 | DONE | 2.0 |
| C02 | Code Survival Calculator | 3.0 | DONE | 3.0 |
| C03 | Revert Chain Analysis | 2.5 | NOT_STARTED | - |
| C04 | Hotfix Tracking | 2.0 | NOT_STARTED | - |
| C05 | Research Dataset Enhancement | 2.5 | NOT_STARTED | - |
| C06 | Part C Integration Tests | 2.5 | NOT_STARTED | - |

### Part D: Replay Mode (DEFERRED to Phase 3D)

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| D01 | Replay Mode Infrastructure | 3.0 | DEFERRED | - |
| D02 | Corpus Loading/Saving | 2.0 | DEFERRED | - |

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

## Part C: GitHub Simulation + Quality Analysis (REVISED)

> **Design Decisions** (January 3, 2026):
> 1. **PR Generation**: Derive on-the-fly from commit groupings with session-based parameters
> 2. **Greenfield**: First commit timestamp for the file (not OS file creation)
> 3. **Replay Mode**: Deferred to Phase 3D (use seeded RNG for reproducibility)
> 4. **Quality Correlations**: Probabilistic with sigmoid risk score
> 5. **Code Survival**: File-level tracking (simple, fast, good enough)

### Research Framework Reference

**Primary Reference**: `docs/design/External - Methods Proposal - AI on SDLC Study.md`

This document defines the scientific framework for measuring AI impact on SDLC:

| Category | Metrics to Implement |
|----------|---------------------|
| **Independent (Table 1)** | AI Usage Intensity, PR Volume, PR Scatter, Greenfield Index, Repo Maturity |
| **Velocity (Table 2)** | Coding Lead Time, Pickup Time, Review Lead Time, Volume Throughput, Merge Rate |
| **Review Costs (Table 3)** | Review Density, Iteration Count, Rework Ratio, Scope Creep, Reviewer Count |
| **Quality (Table 4)** | Revert Rate, Code Survival Rate, Hotfix Follow-up |

**Excluded Metrics** (per methodology doc):
- Repository-local code duplication (cross-repo in microservices)
- Test coverage (CI-pipeline specific, incomplete)
- Deep static/architectural metrics (high compute overhead)

---

### Step C00: PR Generation Pipeline

**Priority**: CRITICAL - All other C steps depend on this

**Estimated Hours**: 4.0

**Files**:
- `internal/generator/pr_generator.go` (new)
- `internal/generator/session.go` (new)
- `internal/models/pr.go` (update)
- `internal/storage/memory.go` (update)
- `internal/storage/interface.go` (update)

**Description**:
Generate PRs on-the-fly from commit groupings using session-based parameters.
PRs emerge naturally from "work sessions" with developer-specific characteristics.

**Session Model**:
```go
type Session struct {
    Developer     seed.Developer
    Repo          seed.Repository
    Branch        string
    StartTime     time.Time
    MaxCommits    int           // Seniority-based: seniors 5-12, juniors 2-5
    TargetLoC     int           // Affects commit sizes
    InactivityGap time.Duration // From working hours
    Commits       []models.Commit
}
```

**Grouping Rules**:
1. Open PR when work session starts (first commit on new branch)
2. Keep adding commits until:
   - Inactivity gap > N minutes (developer-specific)
   - Max commits per PR reached (seniority-based)
   - Random early close triggered (volatility)
3. Finalize PR metrics and store

**Acceptance Criteria**:
- [ ] `Session` struct with seniority-based `MaxCommits` sampling
- [ ] `sampleMaxCommits(seniority)`: junior=2-5, mid=4-8, senior=5-12
- [ ] `sampleTargetLoC(seniority)`: junior=50-150, mid=100-300, senior=150-500
- [ ] `sampleInactivityGap(workingHours)`: 15-60 minutes
- [ ] Commits grouped by `(repo, branch, author)` with time-based session boundaries
- [ ] PR envelope stores: number, timestamps, author, repo, branch, commit list
- [ ] PR metrics calculated: additions, deletions, changed_files, ai_ratio
- [ ] `first_commit_at`, `created_at` derived from commit timestamps
- [ ] Storage methods: `AddPR()`, `GetPRsByRepo()`, `GetPRsByRepoAndState()`
- [ ] Unit tests with seeded RNG for reproducibility
- [ ] Memory bounded (persist PR envelope, not full commit copies)

**Correlation Enforcement** (via session params):
```go
func StartSession(dev seed.Developer, repo seed.Repository) *Session {
    return &Session{
        Developer:     dev,
        Repo:          repo,
        MaxCommits:    sampleMaxCommits(dev.Seniority),
        TargetLoC:     sampleTargetLoC(dev.Seniority),
        InactivityGap: sampleGap(dev.WorkingHoursBand),
    }
}
```

---

### Step C01: Wire GitHub/Research Routes

**Estimated Hours**: 2.0

**Files**:
- `internal/server/router.go` (update)
- `internal/api/github/commits.go` (new)
- `internal/api/github/files.go` (new)

**Description**:
Wire all GitHub and Research endpoints to the main router.
Implement missing handlers for commits and files.

**New Handlers Needed**:
- `GET /repos/{owner}/{repo}/pulls/{n}/commits` - List commits in PR
- `GET /repos/{owner}/{repo}/pulls/{n}/files` - List files changed in PR
- `GET /repos/{owner}/{repo}/commits` - List commits in repo
- `GET /repos/{owner}/{repo}/commits/{sha}` - Get commit details

**Acceptance Criteria**:
- [ ] All 12 GitHub routes wired to router
- [ ] All 5 Research routes wired to router
- [ ] `ListPullCommits` handler returns commits linked to PR
- [ ] `ListPullFiles` handler returns files with `greenfield_index`
- [ ] `ListCommits` handler with pagination and filtering
- [ ] `GetCommit` handler with AI telemetry fields
- [ ] Greenfield calculation: `file_created_at` = first commit timestamp for file
- [ ] Files < 30 days old marked as `is_greenfield: true`
- [ ] `greenfield_index` = % of PR lines in greenfield files

**Route Additions to `router.go`**:
```go
// GitHub Simulation API
mux.Handle("/repos", github.ListRepos(store))
mux.Handle("/repos/", github.RepoRouter(store))

// Research API
mux.Handle("/research/dataset", research.DatasetHandler(researchGen))
mux.Handle("/research/metrics/velocity", research.VelocityMetricsHandler(researchGen))
mux.Handle("/research/metrics/review-costs", research.ReviewCostMetricsHandler(researchGen))
mux.Handle("/research/metrics/quality", research.QualityMetricsHandler(researchGen))
```

---

### Step C02: Code Survival Calculator (File-Level)

**Estimated Hours**: 3.0

**Files**:
- `internal/services/survival.go` (new)
- `internal/models/quality.go` (new)
- `internal/api/github/analysis.go` (new)
- `internal/storage/memory.go` (update)

**Description**:
Track file-level code survival across commits.
Implement the `/repos/{owner}/{repo}/analysis/survival` endpoint.

**File Survival Model**:
```go
type FileSurvival struct {
    FilePath        string    `json:"file_path"`
    RepoName        string    `json:"repo_name"`
    CreatedAt       time.Time `json:"created_at"`      // First commit timestamp
    LastModifiedAt  time.Time `json:"last_modified_at"`
    AILinesAdded    int       `json:"ai_lines_added"`
    HumanLinesAdded int       `json:"human_lines_added"`
    TotalLines      int       `json:"total_lines"`
    RevertEvents    int       `json:"revert_events"`
    IsDeleted       bool      `json:"is_deleted"`
    DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}
```

**Survival Calculation**:
- Track each file from first appearance to deletion or observation date
- `survival_rate` = files_surviving / files_added_in_cohort
- Support cohort windows: 30d, 60d, 90d

**Acceptance Criteria**:
- [ ] `FileSurvival` model with all fields from schema
- [ ] `SurvivalService.CalculateSurvival(repoName, cohortStart, cohortEnd, observationDate)`
- [ ] Track file birth as first commit containing that file path
- [ ] Track file death as commit that deletes the file
- [ ] Aggregate AI vs human lines per file
- [ ] Handler: `GET /repos/{owner}/{repo}/analysis/survival`
- [ ] Response matches `github-sim-api.yaml` `SurvivalAnalysis` schema
- [ ] `by_developer` breakdown in response
- [ ] Unit tests for survival rate calculation
- [ ] E2E test: 100 files created, 20 deleted → 80% survival

**Response Format** (from github-sim-api.yaml):
```json
{
  "cohort_start": "2025-12-01",
  "cohort_end": "2025-12-31",
  "observation_date": "2026-01-31",
  "total_lines_added": 15000,
  "lines_surviving": 12500,
  "survival_rate": 0.833,
  "by_developer": [
    { "email": "alice@example.com", "lines_added": 5000, "lines_surviving": 4500, "survival_rate": 0.90 }
  ]
}
```

---

### Step C03: Revert Chain Analysis

**Estimated Hours**: 2.5

**Files**:
- `internal/services/reverts.go` (new)
- `internal/models/quality.go` (update)
- `internal/api/github/analysis.go` (update)

**Description**:
Detect revert commits and link them to original PRs.
Implement the `/repos/{owner}/{repo}/analysis/reverts` endpoint.

**Revert Detection**:
1. Message pattern matching: `revert`, `Revert`, `rollback`
2. Link to original PR via commit SHA or PR number in message
3. Calculate `days_to_revert`

**Risk Score Model** (probabilistic enforcement):
```go
func CalculateRevertRisk(pr models.PullRequest, dev seed.Developer) float64 {
    // Sigmoid: high AI + low seniority + high volatility → higher risk
    rawScore := a*pr.AIRatio + b*volatility + c*seniorityPenalty(dev.Seniority)
    return 1.0 / (1.0 + math.Exp(-rawScore))
}

func ShouldRevert(pr models.PullRequest, dev seed.Developer, rng *rand.Rand) bool {
    risk := CalculateRevertRisk(pr, dev)
    return rng.Float64() < risk
}
```

**Acceptance Criteria**:
- [ ] `RevertEvent` model linking revert commit to original PR
- [ ] Pattern matching for revert detection in commit messages
- [ ] `CalculateRevertRisk()` with sigmoid function
- [ ] Seed correlation: `ai_ratio_revert_rate` applied probabilistically
- [ ] `RevertService.GetReverts(repoName, windowDays, since, until)`
- [ ] Handler: `GET /repos/{owner}/{repo}/analysis/reverts`
- [ ] Response matches `github-sim-api.yaml` `RevertAnalysis` schema
- [ ] Default `window_days` = 7
- [ ] Unit tests: verify correlation holds at population level (1000 PRs)
- [ ] E2E test: high-AI PRs have statistically higher revert rate

**Response Format**:
```json
{
  "window_days": 7,
  "total_prs_merged": 500,
  "total_prs_reverted": 12,
  "revert_rate": 0.024,
  "reverted_prs": [
    { "pr_number": 123, "merged_at": "2026-01-10T10:00:00Z", "reverted_at": "2026-01-12T15:00:00Z", "days_to_revert": 2.2 }
  ]
}
```

---

### Step C04: Hotfix Tracking

**Estimated Hours**: 2.0

**Files**:
- `internal/services/hotfixes.go` (new)
- `internal/models/quality.go` (update)
- `internal/api/github/analysis.go` (update)

**Description**:
Detect fix-PRs that follow a merged PR within 48 hours to the same files.
Implement the `/repos/{owner}/{repo}/analysis/hotfixes` endpoint.

**Hotfix Detection**:
1. For each merged PR, find subsequent PRs within `window_hours`
2. Check for overlapping file paths
3. Mark as hotfix if title/body contains: `fix`, `hotfix`, `urgent`, `patch`

**Acceptance Criteria**:
- [ ] `HotfixEvent` model linking original PR to hotfix PR
- [ ] File path overlap detection
- [ ] Title/body pattern matching for fix indicators
- [ ] `HotfixService.GetHotfixes(repoName, windowHours, since, until)`
- [ ] Handler: `GET /repos/{owner}/{repo}/analysis/hotfixes`
- [ ] Response matches `github-sim-api.yaml` `HotfixAnalysis` schema
- [ ] Default `window_hours` = 48
- [ ] `files_in_common` array in response
- [ ] Unit tests for overlap detection
- [ ] E2E test: simulate 3 hotfix scenarios

**Response Format**:
```json
{
  "window_hours": 48,
  "total_prs_merged": 500,
  "prs_with_hotfix": 25,
  "hotfix_rate": 0.05,
  "hotfix_prs": [
    { "original_pr": 120, "hotfix_pr": 125, "hours_between": 4.5, "files_in_common": ["src/auth.ts"] }
  ]
}
```

---

### Step C05: Research Dataset Enhancement

**Estimated Hours**: 2.5

**Files**:
- `internal/generator/research_generator.go` (update)
- `internal/models/research.go` (update)
- `internal/services/research_metrics.go` (update)

**Description**:
Enhance the research dataset with all missing fields from the experimental design framework.
Ensure JOIN keys work across Cursor + GitHub endpoints.

**Metric Definitions** (from Methods Proposal Table 1-4):

| Metric | Formula | Table |
|--------|---------|-------|
| `greenfield_index` | % of PR lines in files created < X days ago | Table 1 |
| `pickup_time_hours` | First Review Timestamp - PR Open Timestamp | Table 2 |
| `coding_lead_time_hours` | PR Open Timestamp - First Commit Timestamp | Table 2 |
| `review_lead_time_hours` | Merge Timestamp - First Review Timestamp | Table 2 |
| `merge_rate` | Merged PRs / (Merged + Closed PRs) | Table 2 |
| `review_density` | Total Review Comments / PR Volume (LoC) | Table 3 |
| `iteration_count` | Count of "Review Requested" → "New Commit" cycles | Table 3 |
| `rework_ratio` | Total LoC Changed During Review / Total LoC in First Draft | Table 3 |
| `scope_creep` | (Final LoC - Initial LoC) / Final LoC | Table 3 |
| `reviewer_count` | Count of unique users who commented or approved | Table 3 |
| `revert_rate` | % of Merged PRs reverted within X days | Table 4 |
| `survival_rate` | % of lines added in Month M still present in M+X | Table 4 |
| `hotfix_rate` | % of PRs followed by fix-PR to same files within Xh | Table 4 |

**Acceptance Criteria**:
- [ ] `ResearchDataPoint` updated with all fields from `github-sim-api.yaml`
- [ ] `greenfield_index` calculated from file creation dates
- [ ] `pickup_time_hours` requires `first_review_at` on PR
- [ ] `scope_creep` requires `initial_additions` tracking
- [ ] `rework_ratio` requires diff between first and final commit
- [ ] `review_density` = comments / total_lines
- [ ] `/research/dataset` returns all 20+ columns
- [ ] CSV export includes all columns with headers
- [ ] JOIN key validation: `commit_hash` matches `sha`, `user_email` matches `author.email`
- [ ] Unit tests for each derived field calculation
- [ ] E2E test: export dataset, load in pandas, verify schema

**Updated ResearchDataPoint**:
```go
type ResearchDataPoint struct {
    // Identifiers
    CommitHash  string `json:"commit_hash"`
    PRNumber    int    `json:"pr_number"`
    AuthorID    string `json:"author_id"`
    AuthorEmail string `json:"author_email"`  // NEW: for JOIN
    RepoName    string `json:"repo_name"`

    // AI Metrics (Independent Variables)
    AIRatio        float64 `json:"ai_ratio"`
    AILinesAdded   int     `json:"ai_lines_added"`    // NEW
    AILinesDeleted int     `json:"ai_lines_deleted"`  // NEW
    NonAILinesAdded int    `json:"non_ai_lines_added"` // NEW

    // PR Metrics (Controls)
    PRVolume       int     `json:"pr_volume"`        // additions + deletions
    PRScatter      int     `json:"pr_scatter"`       // changed_files
    GreenfieldIndex float64 `json:"greenfield_index"` // NEW: % lines in new files

    // Cycle Times (Velocity Outcomes)
    CodingLeadTimeHours float64 `json:"coding_lead_time_hours"`
    PickupTimeHours     float64 `json:"pickup_time_hours"`      // NEW
    ReviewLeadTimeHours float64 `json:"review_lead_time_hours"`

    // Review Costs (Outcomes)
    ReviewDensity    float64 `json:"review_density"`     // NEW
    IterationCount   int     `json:"iteration_count"`
    ReworkRatio      float64 `json:"rework_ratio"`       // NEW
    ScopeCreep       float64 `json:"scope_creep"`        // NEW
    ReviewerCount    int     `json:"reviewer_count"`     // NEW

    // Quality Outcomes
    IsReverted        bool    `json:"is_reverted"`
    HasHotfixFollowup bool    `json:"has_hotfix_followup"` // NEW
    SurvivalRate30d   float64 `json:"survival_rate_30d"`   // NEW

    // Control Variables
    AuthorSeniority   string `json:"author_seniority"`
    RepoMaturity      string `json:"repo_maturity"`
    RepoAgeDays       int    `json:"repo_age_days"`        // NEW
    PrimaryLanguage   string `json:"primary_language"`     // NEW

    Timestamp time.Time `json:"timestamp"`
}
```

---

### Step C06: Part C Integration Tests

**Estimated Hours**: 2.5

**Files**:
- `test/e2e/github_api_test.go` (new)
- `test/e2e/research_api_test.go` (new)
- `test/e2e/quality_analysis_test.go` (new)

**Description**:
Comprehensive E2E tests for all GitHub and Research endpoints.

**Test Scenarios**:

**GitHub API Tests**:
1. `GET /repos` returns repository list
2. `GET /repos/{owner}/{repo}` returns repo details with maturity metrics
3. `GET /repos/{owner}/{repo}/pulls` returns PRs with all cycle time fields
4. `GET /repos/{owner}/{repo}/pulls/{n}` returns full PR with AI summary
5. `GET /repos/{owner}/{repo}/pulls/{n}/commits` returns commits with JOIN keys
6. `GET /repos/{owner}/{repo}/pulls/{n}/files` returns files with greenfield_index
7. `GET /repos/{owner}/{repo}/pulls/{n}/reviews` returns reviews with states
8. `GET /repos/{owner}/{repo}/commits` returns commits with pagination
9. `GET /repos/{owner}/{repo}/commits/{sha}` returns commit with AI contribution
10. `GET /repos/{owner}/{repo}/analysis/survival` returns survival metrics
11. `GET /repos/{owner}/{repo}/analysis/reverts` returns revert analysis
12. `GET /repos/{owner}/{repo}/analysis/hotfixes` returns hotfix analysis

**Research API Tests**:
1. `GET /research/dataset?format=json` returns all columns
2. `GET /research/dataset?format=csv` returns valid CSV with headers
3. `GET /research/metrics/velocity` returns velocity by AI ratio band
4. `GET /research/metrics/review-costs` returns review cost metrics
5. `GET /research/metrics/quality` returns quality metrics

**Hypothesis Validation Tests** (from Methods Proposal):

| Hypothesis | Test | Expected Direction |
|------------|------|-------------------|
| AI reduces Coding Lead Time | Compare high-AI vs low-AI PRs | High AI → Lower coding time |
| AI may increase Review Lead Time | Compare high-AI vs low-AI PRs | High AI → Higher review time |
| High AI increases Review Density | Correlation: AI ratio ↔ comments/LoC | Positive correlation |
| High AI increases Iteration Count | Correlation: AI ratio ↔ cycles | Positive correlation |
| High AI increases Revert Rate | Compare revert rates by AI band | High AI → Higher revert |
| AI code has lower Survival Rate | Cohort analysis by AI ratio | High AI → Lower survival |

**Statistical Validation**:
1. Verify correlations hold at population level (N > 1000 PRs)
2. Use appropriate tests: t-test for means, chi-square for rates
3. Report effect sizes, not just p-values
4. Seed RNG for reproducible statistical tests

**Acceptance Criteria**:
- [ ] 12 GitHub endpoint E2E tests
- [ ] 5 Research endpoint E2E tests
- [ ] Response schema validation against `github-sim-api.yaml`
- [ ] JOIN key consistency test: fetch commit from both APIs, verify match
- [ ] Correlation validation with statistical significance
- [ ] All tests pass with seeded RNG for reproducibility
- [ ] Test coverage > 80% for Part C code

---

## Part C Dependency Graph

```
C00 (PR Generator) ────────────────────────────────────────┐
         │                                                  │
         ▼                                                  │
C01 (Wire Routes) ────────────────────────────────────────┤
         │                                                  │
         ├──────► C02 (Survival) ──────────────────────────┤
         │                                                  │
         ├──────► C03 (Reverts) ───────────────────────────┤
         │                                                  │
         └──────► C04 (Hotfixes) ──────────────────────────┤
                                                            │
                                                            ▼
                                               C05 (Dataset Enhancement)
                                                            │
                                                            ▼
                                               C06 (E2E Tests)
```

---

## Part D: Replay Mode (DEFERRED)

Deferred to Phase 3D per design decision.

For Part C, reproducibility is achieved via:
- Seeded RNG in all generators
- Deterministic event generation from seed.json
- Periodic snapshots (optional Ctrl+E dump)

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
| C00 | Sonnet | Critical path, session modeling |
| C01 | Haiku | Route wiring, straightforward |
| C02 | Sonnet | Survival calculation algorithm |
| C03 | Sonnet | Risk scoring, correlation enforcement |
| C04 | Haiku | Pattern matching, simpler logic |
| C05 | Sonnet | Multiple derived fields |
| C06 | Sonnet | Complex test scenarios |

---

## TDD Checklist (Per Step)

- [ ] Read step details and acceptance criteria
- [ ] Read relevant section in `github-sim-api.yaml` or `cursor_analytics.md`
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
2. C00 → C01 → C02 → C03 → C04 → C05 → C06 (commit: "feat: Phase 3 Part C")
```

---

## Estimated Hours Summary

| Part | Steps | Estimated | Actual |
|------|-------|-----------|--------|
| Part A | A01-A07 | 15-20h | 1.75h ✅ |
| Part B | B00-B08 | 12.5h | 11.9h ✅ |
| Part C | C00-C06 | 18.5h | - |
| **Total** | | **46-51h** | **13.65h** |
