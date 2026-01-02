# Development Session Context

**Last Updated**: January 2, 2026
**Current Phase**: Phase 1 Complete - Deciding Next Steps
**Primary Focus**: cursor-sim v2 Phase 2 OR cursor-analytics-core

---

## Current Status

### Project State

| Component | Status | Notes |
|-----------|--------|-------|
| **cursor-sim v2 Phase 1** | **COMPLETE** ✅ | 16 tasks, 10.75h actual vs 44.5h estimated |
| cursor-sim v2 Phase 2 | NOT_STARTED | GitHub PR simulation |
| cursor-sim v2 Phase 3 | NOT_STARTED | Replay mode, research export |
| cursor-analytics-core | NOT_STARTED | GraphQL aggregator |
| cursor-viz-spa | NOT_STARTED | React dashboard |
| Documentation Cleanup | IN_PROGRESS | SDD alignment |

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

### Current Focus: Ready for Next Feature

SDD alignment complete. Documentation structure now complies with SDD methodology:

- [x] Created `services/cursor-sim/SPEC.md`
- [x] Updated `.claude/DEVELOPMENT.md` (this file)
- [x] Marked `docs/` as reference-only
- [x] Created `.work-items/cursor-sim-phase2/`
- [x] Created `.work-items/cursor-analytics-core/`
- [x] Updated `CLAUDE.md` with hierarchy

### Active Symlink

```
.claude/plans/active -> (none - no feature currently active)
```

To start a new feature, run `/start-feature <feature-name>` where feature-name is one of:
- `cursor-sim-phase2` - GitHub PR simulation
- `cursor-analytics-core` - GraphQL aggregator

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

Start **cursor-analytics-core** to complete the ETL pipeline:
```
cursor-sim (8080) → cursor-analytics-core (4000) → cursor-viz-spa (3000)
     ✅ DONE              Next                      After
```

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
| `CLAUDE.md` | Project instructions |
| `.claude/DEVELOPMENT.md` | This file - session context |
| `.claude/skills/` | Skill definitions |
| `.claude/commands/` | Command definitions |

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
