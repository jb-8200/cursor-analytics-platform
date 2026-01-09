# cursor-sim v2 Specification

**Version**: 2.0.0
**Status**: Phase 4 Complete (CLI Enhancements Done) + Phase 3 Features + P2-F01 (GitHub Analytics) Complete + P4-F04 (External Data Sources) In Progress
**Last Updated**: January 9, 2026 (TASK-GH-15: GitHub Analytics Documentation)

## Overview

cursor-sim is a high-fidelity Cursor Business API simulator that generates synthetic developer usage data. It produces data that exactly matches the Cursor API schema, enabling:

- Testing analytics pipelines without production API access
- Generating correlated datasets for SDLC research
- Drop-in replacement for cursor-analytics-core development

## Implementation Status

| Phase | Features | Status |
|-------|----------|--------|
| Phase 1 (P0) | Seed loading, CLI, 29 API endpoints, commit generation | **COMPLETE** ✅ |
| Phase 2 (P1) | GitHub PR simulation, review cycles, quality outcomes | **COMPLETE** ✅ |
| Phase 3 (P2) | Research export, code survival, quality analysis | **MOSTLY COMPLETE** ✅ |
| Phase 4 (P4) | CLI Enhancements: Flags (F01), Interactive Prompts (F02), TUI (F03) | **COMPLETE** ✅ |
| Phase 3D (Deferred) | Replay mode from corpus files | DEFERRED |

---

## Quick Start

```bash
# Build
go build -o bin/cursor-sim ./cmd/simulator

# Preview mode: Quick seed validation (< 5 seconds)
./bin/cursor-sim -mode preview -seed testdata/valid_seed.yaml
./bin/cursor-sim -mode preview -seed testdata/valid_seed.json

# Run with seed file
./bin/cursor-sim -mode runtime -seed testdata/valid_seed.json -port 8080 -days 90 -velocity high

# Run with interactive configuration (P4-F02)
./bin/cursor-sim -mode runtime -seed testdata/valid_seed.json -interactive

# Test health
curl http://localhost:8080/health

# Query with auth
curl -u cursor-sim-dev-key: http://localhost:8080/analytics/ai-code/commits
```

---

## CLI Configuration

```
cursor-sim [flags]

Flags:
  --mode string        Operation mode: runtime, preview, or replay (default "runtime")
  --seed string        Path to seed file (.json, .yaml, or .yml) (required for runtime/preview mode)
  --corpus string      Path to events.parquet (required for replay mode)
  --port int           HTTP server port (default 8080)
  --days int           Days of history to generate (default 90)
  --velocity string    Event generation rate: low, medium, high (default "medium")

  Interactive Mode (P4-F02):
  --interactive        Enable interactive configuration prompts

  Non-Interactive Mode (P4-F02):
  --developers int     Number of developers (replicates from seed if > seed count)
  --months int         Period in months (converted to days automatically)
  --max-commits int    Maximum commits per developer (0 = unlimited)
```

### Environment Variables

| Variable | Flag Equivalent | Default |
|----------|-----------------|---------|
| CURSOR_SIM_MODE | --mode | runtime |
| CURSOR_SIM_SEED | --seed | (required) |
| CURSOR_SIM_PORT | --port | 8080 |
| CURSOR_SIM_DAYS | --days | 90 |
| CURSOR_SIM_VELOCITY | --velocity | medium |

---

## TUI Features (P4-F03)

### Overview

cursor-sim includes a comprehensive Terminal User Interface (TUI) for improved user experience, built with the Charmbracelet stack (Bubble Tea, Bubbles, Lipgloss).

**Architecture**: Event-based Observer pattern decouples business logic from UI, enabling seamless migration to web interface without code changes.

### Components

#### 1. DOXAPI ASCII Banner

- **Display**: Animated ASCII art "DOXAPI" with purple→pink gradient
- **Shown**: Runtime and interactive modes only
- **Hidden**: Preview mode, `-help` flag, non-TTY environments
- **Fallback**: Plain text "DOXAPI v2.0.0" in non-TTY

**Example Output**:
```
 ____   ___  ___  ___  ____  _____
|  _ \ / _ \|  \/  _ \/ _  \/  _  \
| | | | | | | |_| | | \  __/| | | |
| | | | | | |  _  | | |  __\| | | |
| |_| | |_| | | | |_| | |  /| |_| |
|____/ \___/|_| |_|___/|_|   \_____/

v2.0.0
```

#### 2. Spinners for Loading Phases

- **Usage**: Loading seed data, generating events, creating indexes
- **Display**: Animated spinner with message in TTY
- **Fallback**: Text-based status "⏳ Loading..." in non-TTY
- **Methods**:
  - `Start()` - Begin animation
  - `Stop(message)` - End with completion message ✅
  - `UpdateMessage(text)` - Change message while running

#### 3. Progress Bar for Generation

- **Tracks**: Commit generation progress by day
- **Display**: ASCII bar with percentage (e.g., `[████░░░░░] 40%`)
- **Range**: 0% (start) to 100% (complete)
- **Updates**: Real-time progress via events
- **Methods**:
  - `Update(current)` - Update progress
  - `GetProgress()` - Current count
  - `GetPercentage()` - Percentage 0-100
  - `Render()` - ASCII bar output

#### 4. Interactive Form with Bubble Tea

- **Fields**:
  - Developers (1-100 people)
  - Period (1-24 months)
  - Max Commits (100-2000 per developer)
- **Navigation**: Tab/Shift+Tab between fields, arrow keys
- **Validation**: Real-time range checking with error messages
- **Submit**: Enter key on last field (if valid)
- **Cancel**: ESC key

**Constraints**:
```
Developers:   1-100 (default 10)
Months:       1-24  (default 6)
Max Commits:  100-2000 (default 500)
```

### Event-Based Architecture

**Decoupling Pattern**: Generators emit events; UI subscribes independently.

```
Generator              Event Emitter            TUI Renderer
─────────              ─────────────            ────────────
GenerateCommits() ──→ Emit(ProgressEvent) ──→ HandleEvent()
                                               ├─ Update spinner
                                               └─ Update progress bar
```

**Event Types**:
- `PhaseStartEvent` - Loading/generating phase begins
- `PhaseCompleteEvent` - Phase finished successfully
- `ProgressEvent` - Progress update (current/total)
- `WarningEvent` - Non-fatal issue
- `ErrorEvent` - Fatal error

**Benefits**:
- Generators test without UI
- Web interface can replace TUI without generator changes
- Multiple UIs (CLI, web, API) can consume same events
- Logging/metrics can subscribe without coupling

### Terminal Capability Detection

- **Color Support**: Checks `termenv`, respects `NO_COLOR` env var
- **TTY Detection**: Distinguishes interactive terminal from piped output
- **Graceful Fallback**: Text-only output in CI/CD, non-TTY environments

**Functions**:
- `SupportsColor()` - Color capabilities
- `IsTTY()` - Interactive terminal
- `ShouldUseTUI()` - Use animated UI

### Usage Examples

```bash
# Interactive configuration with TUI
./bin/cursor-sim -mode runtime -seed seed.json -interactive

# Runtime mode (shows banner, spinners, progress)
./bin/cursor-sim -mode runtime -seed seed.json -days 90

# Non-TTY (piped to file)
./bin/cursor-sim -mode runtime -seed seed.json | tee output.log

# Disable colors
NO_COLOR=1 ./bin/cursor-sim -mode runtime -seed seed.json

# Preview mode (no banner, TUI suppressed)
./bin/cursor-sim -mode preview -seed seed.json
```

### Testing

- **Unit Tests**: 92 tests covering all components
- **E2E Tests**: 11 integration tests verifying full workflow
- **Manual Verification**:
  - [ ] TTY with colors: Spinner animates, progress bar visible
  - [ ] Non-TTY: Text fallback, no spinner animation
  - [ ] NO_COLOR: Plain ASCII, no color codes
  - [ ] Interactive: Form navigation and validation work
  - [ ] Preview: No banner displayed

---

## API Reference

### Authentication

All endpoints except `/health` require Basic Auth:
- **Username**: API key (default: `cursor-sim-dev-key`)
- **Password**: empty string

```bash
curl -u API_KEY: http://localhost:8080/teams/members
```

### Rate Limiting

| Endpoint Group | Limit | Response on Exceed |
|----------------|-------|-------------------|
| Team Analytics | 100 req/min | 429 Too Many Requests |
| By-User Analytics | 50 req/min | 429 Too Many Requests |

### Endpoints (30 Total)

#### Health Check

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | No | Health check, returns `{"status":"ok"}` |

#### Team Management API (1 endpoint implemented)

| Method | Path | Auth | Status |
|--------|------|------|--------|
| GET | `/teams/members` | Yes | ✅ Implemented |

#### AI Code Tracking API (2 endpoints implemented)

| Method | Path | Auth | Status |
|--------|------|------|--------|
| GET | `/analytics/ai-code/commits` | Yes | ✅ Implemented |
| GET | `/analytics/ai-code/commits.csv` | Yes | ✅ Implemented |

**Query Parameters:**
- `from` (string): Start date YYYY-MM-DD
- `to` (string): End date YYYY-MM-DD
- `page` (int): Page number (default 1)
- `page_size` (int): Items per page (default 100, max 500)
- `user_id` (string): Filter by user email
- `repo_name` (string): Filter by repository

#### Team Analytics API (11 endpoints)

| Method | Path | Auth | Status |
|--------|------|------|--------|
| GET | `/analytics/team/agent-edits` | Yes | ✅ Implemented |
| GET | `/analytics/team/tabs` | Yes | ✅ Implemented |
| GET | `/analytics/team/dau` | Yes | ✅ Implemented |
| GET | `/analytics/team/models` | Yes | ✅ Implemented |
| GET | `/analytics/team/client-versions` | Yes | ✅ Implemented |
| GET | `/analytics/team/top-file-extensions` | Yes | ✅ Implemented |
| GET | `/analytics/team/mcp` | Yes | ✅ Implemented |
| GET | `/analytics/team/commands` | Yes | ✅ Implemented |
| GET | `/analytics/team/plans` | Yes | ✅ Implemented |
| GET | `/analytics/team/ask-mode` | Yes | ✅ Implemented |
| GET | `/analytics/team/leaderboard` | Yes | ✅ Implemented |

#### By-User Analytics API (9 endpoints)

| Method | Path | Auth | Status |
|--------|------|------|--------|
| GET | `/analytics/by-user/agent-edits` | Yes | ✅ Implemented |
| GET | `/analytics/by-user/tabs` | Yes | ✅ Implemented |
| GET | `/analytics/by-user/models` | Yes | ✅ Implemented |
| GET | `/analytics/by-user/client-versions` | Yes | ✅ Implemented |
| GET | `/analytics/by-user/top-file-extensions` | Yes | ✅ Implemented |
| GET | `/analytics/by-user/mcp` | Yes | ✅ Implemented |
| GET | `/analytics/by-user/commands` | Yes | ✅ Implemented |
| GET | `/analytics/by-user/plans` | Yes | ✅ Implemented |
| GET | `/analytics/by-user/ask-mode` | Yes | ✅ Implemented |

#### GitHub Analytics API (5 endpoints) - P2-F01

| Method | Path | Auth | Status |
|--------|------|------|--------|
| GET | `/analytics/github/prs` | Yes | ✅ Implemented |
| GET | `/analytics/github/reviews` | Yes | ✅ Implemented |
| GET | `/analytics/github/issues` | Yes | ✅ Implemented |
| GET | `/analytics/github/pr-cycle-time` | Yes | ✅ Implemented |
| GET | `/analytics/github/review-quality` | Yes | ✅ Implemented |

**Query Parameters:**

**/analytics/github/prs**
- `status` (string): Filter by PR state (open, merged, closed)
- `author` (string): Filter by author email
- `start_date` (string): Start date YYYY-MM-DD
- `end_date` (string): End date YYYY-MM-DD
- `page` (int): Page number (default 1)
- `page_size` (int): Items per page (default 20, max 100)

**/analytics/github/reviews**
- `pr_id` (int): Filter by PR number
- `reviewer` (string): Filter by reviewer email
- `page` (int): Page number (default 1)
- `page_size` (int): Items per page (default 20, max 100)

**/analytics/github/issues**
- `state` (string): Filter by issue state (open, closed)
- `labels` (string): Comma-separated labels (AND logic)
- `page` (int): Page number (default 1)
- `page_size` (int): Items per page (default 20, max 100)

**/analytics/github/pr-cycle-time**
- `from` (string): Start date YYYY-MM-DD
- `to` (string): End date YYYY-MM-DD

**/analytics/github/review-quality**
- `from` (string): Start date YYYY-MM-DD
- `to` (string): End date YYYY-MM-DD

**Response Schemas:**

*Standard List Response (PRs, Reviews, Issues)*:
```json
{
  "data": [
    {
      "number": 123,
      "title": "feat: add authentication handler",
      "state": "merged",
      "author_email": "alice@example.com",
      "author_name": "Alice Developer",
      "repo_name": "acme/platform",
      "base_branch": "main",
      "head_branch": "feature/auth-login",
      "additions": 250,
      "deletions": 45,
      "commit_count": 8,
      "ai_ratio": 0.68,
      "created_at": "2026-01-01T10:00:00Z",
      "merged_at": "2026-01-05T14:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 150
  },
  "params": {
    "status": "merged",
    "author": "alice@example.com"
  }
}
```

*PR Cycle Time Response*:
```json
{
  "data": {
    "avgTimeToFirstReview": 172800,
    "avgTimeToMerge": 518400,
    "medianTimeToMerge": 432000,
    "p50TimeToMerge": 432000,
    "p75TimeToMerge": 604800,
    "p90TimeToMerge": 777600,
    "totalPRsAnalyzed": 150
  },
  "params": {
    "from": "2026-01-01",
    "to": "2026-01-31"
  }
}
```

*Review Quality Response*:
```json
{
  "data": {
    "approval_rate": 0.70,
    "avg_reviewers_per_pr": 2.1,
    "avg_comments_per_review": 2.8,
    "changes_requested_rate": 0.20,
    "pending_rate": 0.10,
    "total_reviews": 315,
    "total_prs_reviewed": 150
  },
  "params": {
    "from": "2026-01-01",
    "to": "2026-01-31"
  }
}
```

**Legend**: ✅ Fully implemented

---

## Response Format

All endpoints return JSON with this structure:

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "pageSize": 100,
    "totalPages": 5,
    "hasNextPage": true,
    "hasPreviousPage": false
  },
  "params": {
    "from": "2026-01-01",
    "to": "2026-01-07",
    "page": 1,
    "pageSize": 100
  }
}
```

### Commit Schema

```json
{
  "commitHash": "abc123def456...",
  "userId": "user_001",
  "userEmail": "dev@example.com",
  "userName": "Jane Developer",
  "repoName": "acme/platform",
  "branchName": "feature/auth",
  "isPrimaryBranch": false,
  "totalLinesAdded": 150,
  "totalLinesDeleted": 45,
  "tabLinesAdded": 90,
  "tabLinesDeleted": 20,
  "composerLinesAdded": 35,
  "composerLinesDeleted": 10,
  "nonAiLinesAdded": 25,
  "nonAiLinesDeleted": 15,
  "message": "feat: add authentication handler",
  "commitTs": "2026-01-02T14:30:00Z",
  "createdAt": "2026-01-02T14:30:00Z"
}
```

**Invariant**: `totalLinesAdded = tabLinesAdded + composerLinesAdded + nonAiLinesAdded`

---

## Seed File Schema

The seed file defines developers, repositories, and generation parameters. Both JSON (.json) and YAML (.yaml, .yml) formats are supported:

```json
{
  "version": "1.0.0",
  "developers": [{
    "user_id": "user_001",
    "email": "dev@example.com",
    "name": "Jane Developer",
    "org": "Acme Corp",
    "division": "Engineering",
    "team": "Platform",
    "role": "engineer",
    "region": "us-west",
    "timezone": "America/Los_Angeles",
    "locale": "en-US",
    "seniority": "senior",
    "activity_level": "high",
    "acceptance_rate": 0.85,
    "pr_behavior": {
      "prs_per_week": 2.5,
      "avg_pr_size_loc": 150,
      "avg_files_per_pr": 5,
      "review_thoroughness": 0.8,
      "iteration_tolerance": 3
    },
    "coding_speed": {"mean": 4.0, "std_dev": 1.5},
    "preferred_models": ["claude-sonnet-4", "gpt-4o"],
    "chat_vs_code_ratio": {"chat": 0.3, "code": 0.7},
    "working_hours_band": {"start": 9, "end": 18}
  }],
  "repositories": [{
    "repo_name": "acme/platform",
    "primary_language": "go",
    "age_days": 730,
    "maturity": "mature",
    "teams": ["Platform", "API"]
  }],
  "text_templates": {
    "commit_messages": ["feat: {action} {component}", "fix: {bug} in {component}"],
    "pr_titles": ["[{team}] {action}: {description}"]
  },
  "correlations": {
    "seniority_acceptance_rate": {"junior": 0.65, "mid": 0.75, "senior": 0.85},
    "ai_ratio_revert_rate": {"high": 0.05, "medium": 0.08, "low": 0.12}
  },
  "pr_lifecycle": {
    "review_time_hours": {"mean": 24, "std_dev": 12},
    "iterations": {"mean": 2, "max": 5}
  }
}
```

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         cursor-sim v2                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────────┐  │
│  │ CLI/Config   │───▶│ Seed Loader  │───▶│ Commit Generator     │  │
│  │ (flag pkg)   │    │ (JSON parse) │    │ (Poisson timing)     │  │
│  └──────────────┘    └──────────────┘    └──────────────────────┘  │
│                              │                      │               │
│                              ▼                      ▼               │
│                       ┌──────────────────────────────────┐         │
│                       │      In-Memory Storage           │         │
│                       │  ┌────────────┐ ┌─────────────┐  │         │
│                       │  │ Developers │ │ Commits     │  │         │
│                       │  │ (by ID)    │ │ (by time)   │  │         │
│                       │  └────────────┘ └─────────────┘  │         │
│                       └──────────────────────────────────┘         │
│                                      │                              │
│  ┌───────────────────────────────────┴────────────────────────┐   │
│  │                    HTTP Router                              │   │
│  │  ┌─────────────────────────────────────────────────────┐   │   │
│  │  │ Middleware: Logger → RateLimit → BasicAuth          │   │   │
│  │  └─────────────────────────────────────────────────────┘   │   │
│  │                                                             │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │   │
│  │  │ Admin API   │ │ AI Code API │ │ Analytics   │          │   │
│  │  │ /teams/*    │ │ /analytics/ │ │ /analytics/ │          │   │
│  │  │ (1 ep)      │ │ ai-code/*   │ │ team/* (11) │          │   │
│  │  └─────────────┘ │ (2 eps)     │ │ by-user (9) │          │   │
│  │                  └─────────────┘ └─────────────┘          │   │
│  └────────────────────────────────────────────────────────────┘   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### Package Structure

```
services/cursor-sim/
├── cmd/simulator/
│   ├── main.go           # Entry point
│   └── main_test.go      # Integration tests
├── internal/
│   ├── config/           # CLI flags, env vars
│   ├── seed/             # Seed loading and validation (JSON/YAML)
│   ├── preview/          # Preview mode for seed validation
│   ├── models/           # Cursor API types (13 files)
│   ├── generator/        # Event generation (25 files: commits, PRs, reviews, quality)
│   ├── storage/          # In-memory storage
│   ├── api/              # Middleware, response helpers
│   │   ├── cursor/       # Cursor API handlers
│   │   └── github/       # GitHub API handlers
│   ├── server/           # HTTP router
│   ├── services/         # Business logic (code survival, hotfix, revert analysis)
│   ├── export/           # Research data export (Parquet, CSV)
│   └── replay/           # Replay mode infrastructure (deferred)
├── test/e2e/             # End-to-end tests
├── testdata/             # Test seed files
├── bin/                  # Build output
├── go.mod
└── SPEC.md               # This file
```

---

## Generation Algorithm

### Commit Generation (Poisson Process)

1. For each developer in seed:
   - Calculate commit rate: `commits_per_day = prs_per_week * 2 / 7`
   - Scale by velocity multiplier (low: 0.5, medium: 1.0, high: 2.0)

2. For each commit:
   - Inter-arrival time: exponential distribution with rate λ
   - Commit size: lognormal(mean=avg_pr_size_loc, σ=0.5)
   - AI ratio: based on developer's `acceptance_rate`
   - TAB/COMPOSER split: 60-80% TAB, 20-40% COMPOSER

### PR Generation (Session-Based Grouping)

**Algorithm**: Groups commits into PRs using session-based rules by (repo, branch, author).

1. **Session Creation**:
   - Commits on same branch by same author form sessions
   - Session closes after 4+ hour gap or based on developer behavior
   - Each session becomes one PR

2. **PR Status Distribution**:
   - 85% merged (merged_at = created_at + 1-7 days)
   - 10% closed (closed_at = created_at + 1-14 days)
   - 5% remain open

3. **PR Metrics**:
   - Aggregates all commits in session (additions, deletions, AI ratio)
   - Commit count: 3-10 commits per PR (typical)
   - Branch names: feature/*, bugfix/* patterns
   - PR titles: Generated from branch name and commit messages

### Review Generation (Per-PR)

**Algorithm**: Generates reviews for PRs based on state and review patterns.

1. **Reviewer Selection**:
   - 1-3 reviewers per PR (configurable via seed)
   - Prefers same-team reviewers
   - Excludes PR author from reviewer pool

2. **Review State Distribution**:
   - 70% approved (LGTM, Ship it!)
   - 20% changes_requested (with 0-5 inline comments)
   - 10% pending (incomplete reviews)

3. **Review Timing**:
   - Review submitted between PR creation and merge/close
   - First review typically within 1-2 days
   - Multiple review iterations for non-approved reviews

4. **Review Comments**:
   - Approved: Short positive messages
   - Changes requested: 0-5 inline comments with suggestions
   - Pending: No comments yet

### Issue Generation (PR-Linked)

**Algorithm**: Generates issues linked to merged PRs.

1. **Issue Creation Rate**:
   - 40% of merged PRs close an issue
   - 10% of generated issues remain open

2. **Issue Timing**:
   - Issues created 1-7 days before PR creation
   - Closed when PR merges (if linked)

3. **Issue Properties**:
   - Labels: bug, feature, enhancement (1-2 per issue)
   - Title: Derived from PR title
   - State: open or closed based on PR merge status

### Reproducibility

Same seed file + same random seed = identical output:

```go
rng := rand.New(rand.NewSource(time.Now().UnixNano()))
// For reproducible runs:
rng := rand.New(rand.NewSource(12345))
```

---

## Performance Targets

| Metric | Target | Actual |
|--------|--------|--------|
| Startup time | < 2s | ~500ms |
| API p99 latency | < 50ms | ~10ms |
| Memory (10k commits) | < 100MB | ~50MB |
| Generation rate | 1000+ commits/sec | ~5000/sec |

---

## Test Coverage

| Package | Coverage | Notes |
|---------|----------|-------|
| seed | 96.2% | Comprehensive validation |
| storage | 98.7% | Including concurrency |
| server | 100.0% | Router tests |
| api | 91.3% | Middleware, response |
| api/cursor | 87.5% | Handler tests |
| generator | 87.0% | Including performance |
| config | 89.4% | Flag parsing |
| cmd/simulator | 61.7% | Signal handling hard to test |

**Overall**: 90.3% average (exceeds 80% target)

---

## Phase 2 Features (Implemented) ✅

### SIM-R009: PR Generation Pipeline ✅
- Generates full PR lifecycle from commits
- Links commits to PRs with proper foreign keys
- Tracks PR state (open, review, merged, closed)
- Implements realistic PR timelines and scatter patterns

### SIM-R010: Review Simulation ✅
- Generates review comments from templates
- Simulates review iterations based on correlations
- Tracks approval/rejection cycles
- Models reviewer assignment and thoroughness

### SIM-R011: GitHub Repos/PRs API ✅
- `GET /repos` - List repositories
- `GET /repos/{owner}/{repo}/pulls` - List PRs
- `GET /repos/{owner}/{repo}/pulls/{number}` - PR details
- `GET /repos/{owner}/{repo}/pulls/{number}/commits` - PR commits
- `GET /repos/{owner}/{repo}/analysis/reverts` - Revert analysis
- `GET /repos/{owner}/{repo}/analysis/hotfixes` - Hotfix detection

### SIM-R011a: GitHub Analytics API (P2-F01) ✅
- `GET /analytics/github/prs` - PR analytics with filtering (status, author, date range) and pagination ✅
- `GET /analytics/github/reviews` - Review analytics with filtering (pr_id, reviewer) and pagination ✅
- `GET /analytics/github/issues` - Issue analytics with filtering (state, labels) and pagination ✅
- `GET /analytics/github/pr-cycle-time` - PR lifecycle metrics (time to first review, time to merge, percentiles) ✅
- `GET /analytics/github/review-quality` - Review quality metrics (approval rate, avg reviewers, avg comments) ✅

**E2E Test Coverage** (TASK-GH-14): 11 test scenarios covering full pipeline (commits → PRs → reviews → issues), filtering, pagination, authentication, and error handling ✅

**Query Parameters**:
- `page` (int): Page number (default 1)
- `page_size` (int): Items per page (default 20, max 100)
- PRs: `status`, `author`, `start_date`, `end_date`
- Reviews: `pr_id`, `reviewer`
- Issues: `state` (open/closed), `labels` (comma-separated)
- PR Cycle Time: `from` (YYYY-MM-DD), `to` (YYYY-MM-DD)
- Review Quality: `from` (YYYY-MM-DD), `to` (YYYY-MM-DD)

**Response Formats**:

*Standard List Response (PRs, Reviews, Issues)*:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  },
  "params": {
    "status": "merged",
    "author": "alice@example.com"
  }
}
```

*PR Cycle Time Response*:
```json
{
  "data": {
    "avgTimeToFirstReview": 172800,    // seconds (2 days)
    "avgTimeToMerge": 518400,          // seconds (6 days)
    "medianTimeToMerge": 432000,       // seconds (5 days)
    "p50TimeToMerge": 432000,
    "p75TimeToMerge": 604800,          // 7 days
    "p90TimeToMerge": 777600,          // 9 days
    "totalPRsAnalyzed": 150
  },
  "params": {
    "from": "2025-01-01",
    "to": "2025-01-31"
  }
}
```

*Review Quality Response*:
```json
{
  "data": {
    "approval_rate": 0.75,              // % of reviews that are approvals
    "avg_reviewers_per_pr": 2.3,        // Avg number of reviewers per merged PR
    "avg_comments_per_review": 3.1,     // Avg comment count per review
    "changes_requested_rate": 0.15,     // % of reviews requesting changes
    "pending_rate": 0.10,               // % of pending reviews
    "total_reviews": 150,               // Total reviews in period
    "total_prs_reviewed": 85            // Total PRs with reviews
  },
  "params": {
    "from": "2026-01-01",
    "to": "2026-01-09"
  }
}
```

### SIM-R012: Quality Outcomes ✅
- Revert tracking with chain analysis
- Hotfix detection (48-hour window)
- Code survival metrics (7, 14, 30 day tracking)
- Risk scoring for revert chains

### SIM-R015: External Data Sources API (P4-F04) ✅

#### Microsoft 365 Copilot Usage API

**Endpoint**: `GET /reports/getMicrosoft365CopilotUsageUserDetail(period='D30')`

Simulates the Microsoft Graph API endpoint for Copilot usage tracking.

**Configuration**:
Routes are conditionally registered only when Copilot is enabled in seed data:
```json
{
  "external_data_sources": {
    "copilot": {
      "enabled": true,
      "total_licenses": 100,
      "active_users": 85,
      "adoption_percentage": 0.85,
      "top_apps": ["Teams", "Word", "Outlook"]
    }
  }
}
```

**Supported Periods**:
- `D7` - Last 7 days
- `D30` - Last 30 days (default)
- `D90` - Last 90 days
- `D180` - Last 180 days

**Query Parameters**:
- `$format` (string): Response format (`application/json` or `text/csv`)

**Response Format (JSON)**:
```json
{
  "@odata.context": "https://graph.microsoft.com/beta/$metadata#reports/getMicrosoft365CopilotUsageUserDetail(period='D30')",
  "value": [
    {
      "reportRefreshDate": "2026-01-09",
      "reportPeriod": 30,
      "userPrincipalName": "user@example.com",
      "displayName": "Jane Developer",
      "lastActivityDate": "2026-01-08",
      "microsoftTeamsCopilotLastActivityDate": "2026-01-08",
      "wordCopilotLastActivityDate": "2026-01-07",
      "excelCopilotLastActivityDate": null,
      "powerPointCopilotLastActivityDate": "2026-01-05",
      "outlookCopilotLastActivityDate": "2026-01-08",
      "oneNoteCopilotLastActivityDate": null,
      "loopCopilotLastActivityDate": null,
      "copilotChatLastActivityDate": "2026-01-08"
    }
  ]
}
```

**CSV Export**: When `$format=text/csv` is specified, returns CSV with all fields:
- Content-Type: `text/csv`
- Content-Disposition: `attachment; filename=copilot-usage-{period}.csv`

**Authentication**: Requires Basic Authentication (same as other cursor-sim endpoints)

**Data Generation**:
- Uses `CopilotGenerator` with app adoption rates:
  - Teams: 85% (most popular)
  - Word: 70%
  - Outlook: 65%
  - PowerPoint: 50%
  - Excel: 40%
  - Copilot Chat: 75%
  - Loop: 20%
  - OneNote: 10%
- Activity dates are randomly distributed within the requested period
- Generated data is stored in memory for consistency across requests

#### Qualtrics Survey Export API

**Endpoints**:
- `POST /API/v3/surveys/{surveyId}/export-responses` - Start export
- `GET /API/v3/surveys/{surveyId}/export-responses/{progressId}` - Check progress
- `GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file` - Download file

Simulates the Qualtrics Survey API v3 export workflow for retrieving survey responses.

**Configuration**:
Routes are conditionally registered only when Qualtrics is enabled in seed data:
```json
{
  "external_data_sources": {
    "qualtrics": {
      "enabled": true,
      "survey_id": "SV_abc123",
      "survey_name": "Developer Satisfaction Survey",
      "response_count": 50
    }
  }
}
```

**Workflow**:
1. Client initiates export with `POST /API/v3/surveys/{surveyId}/export-responses`
2. Server returns `progressId` with status `inProgress` at 0%
3. Client polls `GET /API/v3/surveys/{surveyId}/export-responses/{progressId}`
4. Progress advances by 20% per poll until 100% complete
5. When complete, server provides `fileId` for download
6. Client downloads ZIP file via `GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file`

**Start Export Response** (Step 1):
```json
{
  "result": {
    "progressId": "ES_a5591ddcd0b2409b",
    "status": "inProgress",
    "percentComplete": 0
  },
  "meta": {
    "httpStatus": "200 - OK",
    "requestId": "ES_a5591ddcd0b2409b"
  }
}
```

**Progress Response** (Steps 2-4):
```json
{
  "result": {
    "status": "inProgress",
    "percentComplete": 40
  },
  "meta": {
    "httpStatus": "200 - OK",
    "requestId": "ES_a5591ddcd0b2409b"
  }
}
```

**Completion Response** (Step 5):
```json
{
  "result": {
    "status": "complete",
    "percentComplete": 100,
    "fileId": "FILE_24464e06033e1808"
  },
  "meta": {
    "httpStatus": "200 - OK",
    "requestId": "ES_a5591ddcd0b2409b"
  }
}
```

**File Download Response** (Step 6):
- Content-Type: `application/zip`
- Content-Disposition: `attachment; filename="survey_responses.zip"`
- Body: ZIP file containing CSV with survey responses

**Authentication**: Requires Basic Authentication (same as other cursor-sim endpoints)

**Data Generation**:
- Uses `SurveyGenerator` to create realistic survey responses
- Generates responses from team developers based on seed configuration
- Response count controlled by `response_count` in seed data
- Survey fields:
  - `ResponseID`: Unique response identifier (R_xxx format)
  - `StartDate`: Response start timestamp
  - `EndDate`: Response completion timestamp
  - `Status`: Always "Complete"
  - `Progress`: Always 100
  - `Duration`: Time taken in seconds (60-600)
  - `Finished`: Always "True"
  - `RecordedDate`: Same as EndDate
  - `ResponseId`: Same as ResponseID
  - `DistributionChannel`: Always "anonymous"
  - `UserLanguage`: Always "EN"
  - `Q1_Satisfaction`: Satisfaction rating (1-5 scale)
  - `Q2_MostUsedTool`: Most frequently used tool (Composer, Chat, Inline Edit)
  - `Q3_FreeText`: Optional free-text feedback
- Distribution follows realistic patterns:
  - Satisfaction: 40% high (4-5), 40% medium (3), 20% low (1-2)
  - Tool usage: Composer 40%, Chat 35%, Inline Edit 25%
  - 70% of responses include free-text feedback
- Generated ZIP contains CSV file with all responses
- Export job state maintained in memory via `ExportJobManager`

---

## Phase 3 Features (Mostly Complete) ✅

### SIM-R013: Research Dataset Export ✅
- Exports to Parquet and CSV formats
- Pre-joined research tables with 38 columns
- Proper JOIN keys for SDLC analysis
- Includes: velocity metrics, review costs, quality outcomes
- API: `GET /research/dataset?format=parquet|csv`

### SIM-R014: Code Survival Tracking ✅
- Tracks AI code through refactoring
- Measures code retention at 7, 14, 30 days
- Generates survival rate metrics
- Integrated into research dataset

### SIM-R015: Replay Mode ⏸️
- **Status**: DEFERRED to Phase 3D
- Infrastructure created in `internal/replay/`
- Will support loading events from corpus files
- Future: Replay historical data with time-scaling

---

## Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| Jan 2026 | Seed-based generation | Random data lacks correlations for research |
| Jan 2026 | Exact Cursor API match | Drop-in replacement for testing |
| Jan 2026 | In-memory storage | MVP simplicity, acceptable for simulator |
| Jan 2026 | Go 1.21+ | Team expertise, proven Poisson implementation |
| Jan 2026 | Stub pattern for Phase 1 | Complete API surface quickly, upgrade later |
| Jan 2026 | Full endpoint implementation (Phase 3 Part B) | All 29 endpoints now production-ready |
| Jan 2026 | Revert chain analysis | Track cascading reverts for quality research |
| Jan 2026 | Hotfix detection (48h window) | Identify urgent fixes following merges |
| Jan 2026 | 38-column research dataset | Comprehensive SDLC metrics for analysis |
| Jan 2026 | Defer replay mode | Focus on generation quality first |

---

## Related Documentation

- `.work-items/cursor-sim-v2/` - v2 implementation tracking
- `.work-items/cursor-sim-phase2/` - PR lifecycle implementation
- `.work-items/cursor-sim-phase3/` - Research framework & quality analysis (CURRENT)
- `docs/FEATURES.md` - Project-level feature overview
- `docs/TASKS.md` - Implementation task breakdown
- `.claude/skills/cursor-api-patterns.md` - API implementation patterns
- `.claude/skills/go-best-practices.md` - Go coding standards
