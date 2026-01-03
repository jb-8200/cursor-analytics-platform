# Development Session Context

**Last Updated**: January 3, 2026
**Current Phase**: Phase 3 Active (Part B Complete)
**Primary Focus**: cursor-sim v2 Phase 3 - Part C: Code Quality Analysis

---

## Current Status

### Project State

| Component | Status | Notes |
|-----------|--------|-------|
| **cursor-sim v2 Phase 1** | **COMPLETE** ✅ | 16 tasks, 10.75h actual vs 44.5h estimated |
| **cursor-sim v2 Phase 3 Part A** | **COMPLETE** ✅ | Research Framework (7 steps, 1.75h actual vs 15-20h est) |
| **cursor-sim v2 Phase 3 Part B** | **COMPLETE** ✅ | Stub Completion (8 steps, 11.9h actual vs 12.5h est) |
| cursor-sim v2 Phase 3 Part C | NEXT | Code Quality Analysis (5 steps, 10-15h est) |
| cursor-analytics-core | NOT_STARTED | GraphQL aggregator |
| cursor-viz-spa | NOT_STARTED | React dashboard |
| **SDD Infrastructure** | **ENHANCED** ✅ | genai-specs model adopted for Claude Code |

### cursor-sim v2 Phase 1 Completion Summary

**Time Efficiency**: 10.75h actual / 44.5h estimated = **76% faster than planned**

| Category | Estimated | Actual | Delta |
|----------|-----------|--------|-------|
| Foundation (01-05) | 11.0h | 2.5h | -77% |
| Generation (06-07) | 9.0h | 2.25h | -75% |
| Endpoints (08-13) | 16.5h | 4.25h | -74% |
| Integration (14-16) | 8.0h | 1.75h | -78% |

**Test Coverage**: 90.3% average across all packages

---

## Active Work Items

### Recently Completed (January 2, 2026)

**SDD Infrastructure Enhancement** (Commit: 446fdad)
- ✅ Adopted genai-specs model for Claude Code
- ✅ Created `docs/spec-driven-design.md` - master methodology doc
- ✅ Reorganized skills into 4 categories: process/standards/guidelines/operational
- ✅ Created new skills: spec-process-core, spec-process-dev, spec-user-story, spec-design, spec-tasks
- ✅ Slimmed CLAUDE.md from ~320 to ~140 lines (minimal spine)
- ✅ Updated hooks README clarifying they don't execute in Claude Code
- 📝 Approach: Specs + Skills + Discipline (not automated hooks)

**Phase 3 Part B: Steps B00-B08** (Commits: dadc124, 33842cf, f1c937b, 789554b, 0601d1b, 26b69ab, b59399e, 32820ae, 1274cc1)
- ✅ B00: Fixed Analytics Response Format (team vs by-user)
- ✅ B01: Updated 14 data models to match Cursor API exactly
- ✅ B02: Model Usage Generator & Handler
- ✅ B03: Client Version Generator & Handler
- ✅ B04: File Extension Analytics Generator & Handler
- ✅ B05: MCP/Commands/Plans/Ask-Mode Feature Generators & Handlers
- ✅ B06: Leaderboard Handler with Dual Rankings (Tab + Agent)
- ✅ B07: All 9 By-User Analytics Endpoint Handlers
- ✅ B08: Comprehensive E2E Tests for All 20 Analytics Endpoints
- ✅ All tests passing (15/15 packages, 62 E2E test cases)
- ⏱️ Actual: 11.9h / Estimated: 12.5h (5% under budget)

### Current Focus: Part C - GitHub Simulation + Quality Analysis

> **Design Decisions Finalized** (January 3, 2026):
> - PR Generation: Session-based, on-the-fly from commit groupings ✅ COMPLETE
> - Greenfield: First commit timestamp for file
> - Quality Correlations: Probabilistic with sigmoid risk score
> - Code Survival: File-level tracking
> - Replay Mode: Deferred to Phase 3D

**C00 Complete** ✅ (4.0h estimated / actual TBD)
- Session model with seniority-based parameters (junior/mid/senior)
- PR grouping with inactivity gap (15-60m), max commits rules
- 58 generator tests passing (0 failures)
- Backwards-compatible API for e2e tests
- Files: `session.go`, `session_test.go`, `pr_generator.go`, `pr_generator_test.go`, `pr_generator_integration_test.go`
- Commit: 005cb62

**C01 Complete** ✅ (2.0h estimated / actual TBD)
- All 12 GitHub routes wired to router (ListRepos, RepoRouter with nested routes)
- All 5 Research routes wired to router (DatasetHandler, VelocityMetricsHandler, ReviewCostMetricsHandler, QualityMetricsHandler)
- Implemented missing handlers: ListCommits, GetCommit, ListPullCommits, ListPullFiles
- Implemented RepoRouter for dynamic GitHub API routing
- Greenfield index calculation for PR files
- All tests passing (15/15 packages)
- Files: `commits.go`, `commits_test.go`, `files.go`, `files_test.go`, `router.go`
- Commit: 68a81aa

**C02 Complete** ✅ (3.0h estimated / 3.0h actual)
- FileSurvival model with file lifecycle tracking
- SurvivalService with cohort-based survival calculation
- GET /repos/{owner}/{repo}/analysis/survival endpoint
- File birth/death tracking via commit patterns
- Developer breakdown with individual survival rates
- Probabilistic deletion model (sigmoid-based)
- Reproducible results with seeded RNG
- All tests passing (15/15 packages, 9 new tests)
- Files: `quality.go`, `survival.go`, `survival_test.go`, `analysis.go`, `analysis_test.go`
- Commit: 77e5025
- **SDD Violation**: Fixed task.md update (commit e8f4a31)

**Next Task**: C03 - Revert Chain Analysis
- Estimated: 2.5h
- Status: Ready to start
- Detect revert commits via pattern matching
- Calculate revert risk score (sigmoid function)
- Handler for GET /repos/{owner}/{repo}/analysis/reverts

### Active Symlink

```
.claude/plans/active -> ../../.work-items/cursor-sim-phase3/task.md
```

**Currently working on**: cursor-sim Phase 3 Part C (GitHub Simulation + Quality Analysis)

---

## Documentation Hierarchy (SDD Compliant)

### Source of Truth

```
services/{service}/SPEC.md      ← Technical specification
.work-items/{feature}/          ← Active work tracking
├── user-story.md
├── design.md
├── task.md
└── {NN}_step.md
.claude/plans/active            ← Symlink to current work
```

### Reference Documents

```
docs/                           ← Project-level overview (REFERENCE ONLY)
├── DESIGN.md                   ← System architecture
├── FEATURES.md                 ← Feature breakdown
├── TASKS.md                    ← Task overview
└── USER_STORIES.md             ← User stories
```

---

## Next Steps (Choose One)

### Option A: cursor-sim Phase 2 (GitHub PR Simulation)

**Scope**: SIM-R009 through SIM-R012
- PR generation pipeline
- Review simulation
- GitHub Repos/PRs API
- Quality outcomes

**Estimated**: 20-25 hours

### Option B: cursor-analytics-core

**Scope**: CORE-001 through CORE-005
- Data ingestion worker
- PostgreSQL schema
- GraphQL API server
- Metric calculations
- Developer queries

**Estimated**: 25-30 hours

### Recommendation

Start **cursor-sim Phase 2** to complete the simulator before moving downstream:
```
cursor-sim Phase 1 (MVP)     → cursor-sim Phase 2      → cursor-analytics-core
     ✅ DONE                      Next                      After
     (29 endpoints)               (GitHub PR sim)           (ETL pipeline)
```

Phase 1 only provides basic commit data. Phase 2 adds PR lifecycle, reviews, and quality outcomes needed for meaningful SDLC research.

---

## cursor-sim v2 Quick Reference

### Running the Simulator

```bash
cd services/cursor-sim

# Build
go build -o bin/cursor-sim ./cmd/simulator

# Run
./bin/cursor-sim -mode runtime -seed testdata/valid_seed.json -port 8080 -days 90

# Test endpoints
curl http://localhost:8080/health
curl -u cursor-sim-dev-key: http://localhost:8080/analytics/ai-code/commits
```

### API Endpoints (29 total)

| Group | Endpoints | Status |
|-------|-----------|--------|
| Health | 1 | ✅ Implemented |
| Admin | 1 | ✅ Implemented |
| AI Code | 2 | ✅ Implemented |
| Team Analytics | 11 | ✅ 3 functional, 8 stubs |
| By-User Analytics | 9 | ⚡ All stubs |

---

## Development Commands

### cursor-sim

```bash
cd services/cursor-sim

# Run all tests
go test ./... -v

# Check coverage
go test ./... -cover

# Run E2E tests
go test ./test/e2e -v

# Build
go build -o bin/cursor-sim ./cmd/simulator
```

### SDD Workflow

```bash
# Start a feature
/start-feature cursor-analytics-core

# Check current status
/status

# Implement a task
/implement TASK-CORE-001

# Complete a feature
/complete-feature cursor-sim-v2
```

---

## Key Files

### Specifications

| File | Description |
|------|-------------|
| `services/cursor-sim/SPEC.md` | cursor-sim technical spec |
| `services/cursor-analytics-core/SPEC.md` | analytics-core spec |
| `services/cursor-viz-spa/SPEC.md` | dashboard spec |

### Work Items

| Directory | Description |
|-----------|-------------|
| `.work-items/cursor-sim-v2/` | Phase 1 work (COMPLETE) |
| `.work-items/cursor-sim-phase2/` | Phase 2 work (TODO) |
| `.work-items/cursor-analytics-core/` | Aggregator work (TODO) |

### Claude Integration

| File | Purpose |
|------|---------|
| `CLAUDE.md` | Minimal operational spine (~140 lines) |
| `.claude/DEVELOPMENT.md` | This file - session context |
| `.claude/README.md` | Claude Code integration guide |
| `docs/spec-driven-design.md` | Full SDD methodology |

### Skills (Categorized)

| Directory | Purpose |
|-----------|---------|
| `.claude/skills/process/` | Workflow stages (spec-process-core, spec-process-dev) |
| `.claude/skills/standards/` | Artifact templates (spec-user-story, spec-design, spec-tasks) |
| `.claude/skills/guidelines/` | Tech-specific (go-best-practices, cursor-api-patterns) |
| `.claude/skills/operational/` | Day-to-day (sdd-checklist, model-selection-guide) |

---

## Session Checklist

When starting a new session:

1. [x] Read DEVELOPMENT.md (this file)
2. [ ] Check active work: `ls -la .claude/plans/active`
3. [ ] Review current work item in `.work-items/`
4. [ ] Continue with next task or start new feature
5. [ ] Follow TDD: RED → GREEN → REFACTOR
6. [ ] Commit after each task

---

## Architecture Overview

```
┌─────────────────┐     ┌──────────────────────┐     ┌─────────────────┐
│   cursor-sim    │────▶│ cursor-analytics-core│────▶│  cursor-viz-spa │
│   (Go + REST)   │     │   (TS + GraphQL)     │     │  (React + Vite) │
│   Port: 8080    │     │   Port: 4000         │     │   Port: 3000    │
│   ✅ COMPLETE   │     │   ⏳ NOT_STARTED     │     │   ⏳ NOT_STARTED │
└─────────────────┘     └──────────────────────┘     └─────────────────┘
     Simulator              Aggregator                  Dashboard
     (Extract)              (Transform)                  (Load/View)
```

---

**Remember**: Specifications → Tests → Implementation → Refactor

This is the way.
