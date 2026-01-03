# cursor-sim v2 Specification

**Version**: 2.0.0
**Status**: Phase 3 Mostly Complete (Parts A, B, C Done)
**Last Updated**: January 3, 2026

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
| Phase 3D (Deferred) | Replay mode from corpus files | DEFERRED |

---

## Quick Start

```bash
# Build
go build -o bin/cursor-sim ./cmd/simulator

# Run with seed file
./bin/cursor-sim -mode runtime -seed testdata/valid_seed.json -port 8080 -days 90 -velocity high

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
  --mode string        Operation mode: runtime or replay (default "runtime")
  --seed string        Path to seed.json file (required for runtime mode)
  --corpus string      Path to events.parquet (required for replay mode)
  --port int           HTTP server port (default 8080)
  --days int           Days of history to generate (default 90)
  --velocity string    Event generation rate: low, medium, high (default "medium")
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

### Endpoints (29 Total)

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

The seed file defines developers, repositories, and generation parameters:

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
│   ├── seed/             # Seed loading and validation
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

### SIM-R012: Quality Outcomes ✅
- Revert tracking with chain analysis
- Hotfix detection (48-hour window)
- Code survival metrics (7, 14, 30 day tracking)
- Risk scoring for revert chains

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
