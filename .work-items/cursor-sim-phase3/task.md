# Task Breakdown: cursor-sim Phase 3 (Research Framework & Completeness)

## Overview

**Feature**: cursor-sim v2 Phase 3
**Total Estimated Hours**: 35-45
**Number of Steps**: 18
**Current Step**: A07 - Part A Integration Tests

This phase completes cursor-sim with:
- **Part A**: Research Framework (SIM-R013 → SIM-R015)
- **Part B**: Stub Endpoint Completion (17 endpoints)
- **Part C**: Enhanced Code Quality Analysis
- **Part D**: GitHub Router Wiring (OPTIONAL - documented only)

---

## Progress Tracker

### Part A: Research Framework (15-20h)

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| A01 | Research Data Models | 2.0 | DONE | 0.25 |
| A02 | Research Dataset Generator | 3.0 | DONE | 0.25 |
| A03 | Parquet/CSV Export | 2.5 | DONE | 0.25 |
| A04 | Research Metrics Service | 3.0 | DONE | 0.25 |
| A05 | Research API Handlers | 2.5 | DONE | 0.25 |
| A06 | Replay Mode Infrastructure | 3.0 | DONE | 0.25 |
| A07 | Part A Integration Tests | 2.0 | NOT_STARTED | - |

### Part B: Stub Completion (8-12h)

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| B01 | Model Usage Generator | 1.5 | NOT_STARTED | - |
| B02 | Client Version Generator | 1.0 | NOT_STARTED | - |
| B03 | File Extension Analytics | 1.5 | NOT_STARTED | - |
| B04 | MCP/Commands/Plans Generators | 2.0 | NOT_STARTED | - |
| B05 | Ask Mode & Leaderboard | 1.5 | NOT_STARTED | - |
| B06 | By-User Endpoint Handlers | 2.5 | NOT_STARTED | - |
| B07 | Part B Integration Tests | 1.5 | NOT_STARTED | - |

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

## Dependency Graph

```
PART A: Research Framework
A01 (Models)
 │
 ├── A02 (Dataset Generator)
 │    │
 │    └── A03 (Export) ──► A05 (API Handlers)
 │
 ├── A04 (Metrics Service) ──► A05 (API Handlers)
 │
 └── A06 (Replay Mode)
      │
      └── A07 (Integration Tests)

PART B: Stub Completion
B01 ─┬─► B06 (By-User Handlers)
B02 ─┤        │
B03 ─┤        └── B07 (Tests)
B04 ─┤
B05 ─┘

PART C: Code Quality
C01 (Models)
 │
 ├── C02 (Survival Calculator)
 │
 ├── C03 (Revert Analysis)
 │
 └── C04 (Hotfix Tracking)
      │
      └── C05 (Tests)
```

---

## Part A: Research Framework Details

### Step A01: Research Data Models

**File**: `internal/models/research.go`

**Tasks**:
- [ ] Define ResearchDataPoint struct (pre-joined row)
- [ ] Define VelocityMetrics struct
- [ ] Define ReviewCostMetrics struct
- [ ] Define QualityMetrics struct
- [ ] Define CodeSurvivalRecord struct
- [ ] Add JSON/Parquet tags

**Fields for ResearchDataPoint**:
```go
type ResearchDataPoint struct {
    // Identifiers
    CommitHash    string
    PRNumber      int
    AuthorID      string
    RepoName      string

    // AI Metrics (Independent Variables)
    AIRatio       float64
    TabLines      int
    ComposerLines int

    // PR Metrics
    Additions     int
    Deletions     int
    FilesChanged  int

    // Cycle Times (Dependent Variables)
    CodingLeadTimeHours  float64
    ReviewLeadTimeHours  float64
    MergeLeadTimeHours   float64

    // Quality Outcomes (Dependent Variables)
    WasReverted   bool
    RequiredHotfix bool
    ReviewIterations int

    // Controls
    AuthorSeniority string
    RepoMaturity    string
    IsGreenfield    bool

    Timestamp     time.Time
}
```

**Acceptance Criteria**:
- All research variables from DESIGN.md represented
- Proper JSON serialization
- Unit tests for model validation

---

### Step A02: Research Dataset Generator

**File**: `internal/generator/research_generator.go`

**Tasks**:
- [ ] Implement JoinCommitPRData() - correlate commits with PRs
- [ ] Calculate cycle times from timestamps
- [ ] Determine greenfield vs legacy code
- [ ] Apply control variables from seed
- [ ] Unit tests with deterministic seed

**Acceptance Criteria**:
- Dataset rows properly join commit + PR + review data
- Cycle times calculated correctly
- Control variables populated from seed

---

### Step A03: Parquet/CSV Export

**Files**:
- `internal/export/csv.go`
- `internal/export/parquet.go`

**Tasks**:
- [ ] Implement CSV export with proper headers
- [ ] Implement Parquet export using parquet-go
- [ ] Add streaming for large datasets
- [ ] Support filtering by date range
- [ ] Unit tests

**Dependencies**: `github.com/xitongsys/parquet-go`

**Acceptance Criteria**:
- CSV exports load correctly in pandas/R
- Parquet files readable by pyarrow
- 50MB/s export throughput

---

### Step A04: Research Metrics Service

**File**: `internal/services/research_metrics.go`

**Tasks**:
- [ ] CalculateVelocityMetrics(from, to time.Time)
- [ ] CalculateReviewCostMetrics(from, to time.Time)
- [ ] CalculateQualityMetrics(from, to time.Time)
- [ ] Group by AI ratio bands (low/medium/high)
- [ ] Statistical aggregations (mean, median, std)
- [ ] Unit tests

**Acceptance Criteria**:
- Metrics grouped correctly by AI ratio
- Statistical calculations verified
- Edge cases handled (empty data)

---

### Step A05: Research API Handlers

**File**: `internal/api/research/handlers.go`

**Tasks**:
- [ ] GET /research/dataset - Export endpoint
- [ ] GET /research/metrics/velocity
- [ ] GET /research/metrics/review-costs
- [ ] GET /research/metrics/quality
- [ ] Support format query param (csv, parquet, json)
- [ ] Register routes in router
- [ ] Handler tests

**Acceptance Criteria**:
- All endpoints return correct schema
- Format parameter works
- Auth required on all endpoints

---

### Step A06: Replay Mode Infrastructure

**Files**:
- `internal/replay/loader.go`
- `internal/replay/server.go`

**Tasks**:
- [ ] Load events from Parquet corpus file
- [ ] Index events by time for range queries
- [ ] Implement replay HTTP handlers
- [ ] Support --mode=replay --corpus=path CLI
- [ ] Unit tests

**Acceptance Criteria**:
- Replay mode serves pre-generated data
- Same API responses as runtime mode
- Corpus file format documented

---

### Step A07: Part A Integration Tests

**File**: `test/e2e/research_test.go`

**Tasks**:
- [ ] E2E: Generate data → Export → Verify format
- [ ] E2E: Research metrics endpoints
- [ ] E2E: Replay mode serves correct data
- [ ] Verify Parquet readable by pyarrow
- [ ] Performance test: 10k rows export

**Acceptance Criteria**:
- All E2E tests pass
- Export files valid
- Performance meets target

---

## Part B: Stub Completion Details

### Step B01: Model Usage Generator

**File**: `internal/generator/model_generator.go`

**Tasks**:
- [ ] Define ModelUsage struct
- [ ] Generate model usage events from seed preferences
- [ ] Track model distribution per developer
- [ ] Store in memory
- [ ] Update /analytics/team/models endpoint
- [ ] Unit tests

**Models to Support**: gpt-4, gpt-3.5-turbo, claude-3-opus, claude-3-sonnet

---

### Step B02: Client Version Generator

**File**: `internal/generator/version_generator.go`

**Tasks**:
- [ ] Define ClientVersion struct
- [ ] Generate version distribution (realistic semver)
- [ ] Track version adoption over time
- [ ] Update /analytics/team/client-versions endpoint
- [ ] Unit tests

---

### Step B03: File Extension Analytics

**File**: `internal/generator/extension_generator.go`

**Tasks**:
- [ ] Track file extensions from commits
- [ ] Aggregate by extension type
- [ ] Calculate AI ratio per extension
- [ ] Update /analytics/team/top-file-extensions
- [ ] Unit tests

---

### Step B04: MCP/Commands/Plans Generators

**File**: `internal/generator/feature_generators.go`

**Tasks**:
- [ ] Generate MCP (Model Context Protocol) usage
- [ ] Generate command usage patterns
- [ ] Generate plan mode usage
- [ ] Update respective endpoints
- [ ] Unit tests

---

### Step B05: Ask Mode & Leaderboard

**File**: `internal/generator/engagement_generator.go`

**Tasks**:
- [ ] Generate ask-mode interactions
- [ ] Calculate leaderboard rankings
- [ ] Rank by: lines accepted, commits, AI ratio
- [ ] Update endpoints
- [ ] Unit tests

---

### Step B06: By-User Endpoint Handlers

**File**: `internal/api/cursor/byuser_full.go`

**Tasks**:
- [ ] Implement all 9 by-user endpoints with real data
- [ ] Filter by userId parameter
- [ ] Apply pagination
- [ ] Unit tests for each endpoint

**Endpoints**:
- /analytics/by-user/agent-edits
- /analytics/by-user/tabs
- /analytics/by-user/models
- /analytics/by-user/client-versions
- /analytics/by-user/top-file-extensions
- /analytics/by-user/mcp
- /analytics/by-user/commands
- /analytics/by-user/plans
- /analytics/by-user/ask-mode

---

### Step B07: Part B Integration Tests

**File**: `test/e2e/analytics_complete_test.go`

**Tasks**:
- [ ] E2E: All team analytics return data
- [ ] E2E: All by-user analytics return data
- [ ] Verify response schemas match Cursor API
- [ ] No more stub responses

---

## Part C: Code Quality Analysis Details

### Step C01: Code Survival Tracking Models

**File**: `internal/models/survival.go`

**Tasks**:
- [ ] Define CodeBlock struct (file, start_line, end_line, hash)
- [ ] Define SurvivalRecord struct
- [ ] Define SurvivalWindow enum (30d, 60d, 90d)
- [ ] Unit tests

---

### Step C02: Survival Rate Calculator

**File**: `internal/generator/survival_calculator.go`

**Tasks**:
- [ ] Track code blocks through time
- [ ] Calculate survival at 30/60/90 day windows
- [ ] Correlate survival with AI ratio
- [ ] Store results in memory
- [ ] Unit tests with time mocking

---

### Step C03: Revert Chain Analysis

**File**: `internal/generator/revert_analyzer.go`

**Tasks**:
- [ ] Identify revert commits (message patterns)
- [ ] Link reverts to original commits
- [ ] Calculate revert chains (revert of revert)
- [ ] Track AI ratio correlation
- [ ] Unit tests

---

### Step C04: Hotfix Tracking

**File**: `internal/generator/hotfix_tracker.go`

**Tasks**:
- [ ] Identify hotfix PRs (branch patterns, labels)
- [ ] Link hotfixes to original PRs
- [ ] Calculate hotfix rate by AI ratio
- [ ] Time-to-hotfix metrics
- [ ] Unit tests

---

### Step C05: Part C Integration Tests

**File**: `test/e2e/quality_analysis_test.go`

**Tasks**:
- [ ] E2E: Survival rates calculated correctly
- [ ] E2E: Revert chains identified
- [ ] E2E: Hotfix correlations
- [ ] Verify statistical validity

---

## Part D: GitHub Router (OPTIONAL)

> **Note**: This section is documented for future implementation but marked OPTIONAL.
> The GitHub API handlers exist in `internal/api/github/` but are not wired to the router.

### Step D01: Wire GitHub Handlers to Router (OPTIONAL)

**File**: `internal/server/router.go`

**Tasks**:
- [ ] Register /repos routes
- [ ] Register /repos/{owner}/{repo}/pulls routes
- [ ] Add authentication middleware
- [ ] Update SPEC.md with new endpoints

### Step D02: GitHub API E2E Tests (OPTIONAL)

**File**: `test/e2e/github_api_test.go`

**Tasks**:
- [ ] E2E: List repos
- [ ] E2E: Get PRs with pagination
- [ ] E2E: Get PR reviews

---

## Model Recommendations

| Step | Model | Rationale |
|------|-------|-----------|
| A01, A05, B06 | Haiku | Well-specified, standard patterns |
| A02, A03, A04 | Sonnet | Data processing complexity |
| A06, A07 | Sonnet | Replay mode, integration |
| B01-B05 | Haiku | Simple generators |
| B07, C05 | Sonnet | Integration testing |
| C01-C04 | Sonnet | Algorithm complexity |

---

## TDD Checklist (Per Step)

- [ ] Read step details and acceptance criteria
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
1. A01 → A02 → A03 → A04 → A05 → A06 → A07 (commit: "feat: Phase 3 Part A")
2. B01 → B02 → B03 → B04 → B05 → B06 → B07 (commit: "feat: Phase 3 Part B")
3. C01 → C02 → C03 → C04 → C05 (commit: "feat: Phase 3 Part C")
4. D01 → D02 (OPTIONAL - separate commit if done)
```
