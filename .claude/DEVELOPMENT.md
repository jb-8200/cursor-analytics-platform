# Development Session Context

**Last Updated**: January 2026
**Current Phase**: Phase 1 - Complete Cursor API (v2.0)
**Primary Focus**: cursor-sim v2 (Major Rewrite)

---

## Current Status

### Project State

| Component | Status | Notes |
|-----------|--------|-------|
| Documentation | COMPLETE | DESIGN.md, FEATURES.md, TASKS.md updated for v2 |
| cursor-sim v1 | ARCHIVED | To be moved to `services/cursor-sim-v1/` |
| cursor-sim v2 | NOT_STARTED | 16 tasks, ~45 hours estimated |
| cursor-analytics-core | NOT_STARTED | Unchanged from v1 plan |
| cursor-viz-spa | NOT_STARTED | Unchanged from v1 plan |
| OpenAPI Specs | PENDING | To be copied from research zip |
| DataDesigner | PENDING | To be set up in `tools/data-designer/` |

### Major Revision: v1.0 → v2.0

The project underwent a major architectural revision:

| Area | v1.0 | v2.0 |
|------|------|------|
| Data Generation | Internal random | Seed-based from DataDesigner |
| API Surface | Generic (~5 endpoints) | Exact Cursor + GitHub (49 endpoints) |
| Research Support | None | Full SDLC framework |
| Operation Modes | Single | Runtime + Replay |

### What Was Preserved from v1

- Poisson distribution timing (proven, tested)
- Go project structure patterns
- TDD workflow and test infrastructure
- Makefile and Docker configuration

### What Was Replaced

- Developer generation → loads from seed.json
- Event models → Cursor-exact schema with camelCase
- API handlers → all 29 Cursor endpoints
- CLI flags → new --mode, --seed, --corpus flags

---

## Active Work: Pre-Implementation Setup

Before beginning TASK-R001, complete these setup tasks:

### Setup Checklist

- [ ] Archive v1: `mv services/cursor-sim services/cursor-sim-v1`
- [ ] Copy OpenAPI specs to `specs/openapi/`
- [ ] Copy DataDesigner to `tools/data-designer/`
- [ ] Create sample seed.json for testing
- [ ] Validate SDD structure still works

---

## Phase 1: Complete Cursor API (MVP)

**Goal**: Working cursor-sim that exactly matches Cursor Business API

### Tasks Overview (16 total)

| Task | Description | Hours | Model | Status |
|------|-------------|-------|-------|--------|
| TASK-R001 | Project structure | 1 | Haiku | NOT_STARTED |
| TASK-R002 | Seed schema types | 2 | Haiku | NOT_STARTED |
| TASK-R003 | Seed loader + validation | 3 | Sonnet | NOT_STARTED |
| TASK-R004 | CLI v2 flags | 2 | Haiku | NOT_STARTED |
| TASK-R005 | Cursor data models | 3 | Haiku | NOT_STARTED |
| TASK-R006 | Commit generation | 5 | Sonnet | NOT_STARTED |
| TASK-R007 | Storage v2 | 4 | Sonnet | NOT_STARTED |
| TASK-R008 | API infrastructure | 2 | Haiku | NOT_STARTED |
| TASK-R009 | /teams/members | 1.5 | Haiku | NOT_STARTED |
| TASK-R010 | /ai-code/commits | 2 | Sonnet | NOT_STARTED |
| TASK-R011 | /ai-code/commits.csv | 1 | Haiku | NOT_STARTED |
| TASK-R012 | /team/* (11 endpoints) | 6 | Sonnet | NOT_STARTED |
| TASK-R013 | /by-user/* (9 endpoints) | 4 | Sonnet | NOT_STARTED |
| TASK-R014 | Router | 2 | Haiku | NOT_STARTED |
| TASK-R015 | Main application | 2 | Haiku | NOT_STARTED |
| TASK-R016 | E2E tests | 4 | Sonnet | NOT_STARTED |

**Total**: 44.5 hours

### Critical Path

```
R001 → R002 → R003 → R006 → R007 → Endpoints → R014 → R015 → R016
```

---

## Cursor API Surface (29 endpoints)

### Admin API (4)
```
GET  /teams/members
POST /teams/daily-usage-data
POST /teams/filtered-usage-events
POST /teams/spend
```

### AI Code Tracking (4)
```
GET /analytics/ai-code/commits
GET /analytics/ai-code/commits.csv
GET /analytics/ai-code/changes
GET /analytics/ai-code/changes.csv
```

### Team Analytics (11)
```
GET /analytics/team/agent-edits
GET /analytics/team/tabs
GET /analytics/team/dau
GET /analytics/team/client-versions
GET /analytics/team/models
GET /analytics/team/top-file-extensions
GET /analytics/team/mcp
GET /analytics/team/commands
GET /analytics/team/plans
GET /analytics/team/ask-mode
GET /analytics/team/leaderboard
```

### By-User Analytics (9)
```
GET /analytics/by-user/agent-edits
GET /analytics/by-user/tabs
GET /analytics/by-user/models
GET /analytics/by-user/top-file-extensions
GET /analytics/by-user/client-versions
GET /analytics/by-user/mcp
GET /analytics/by-user/commands
GET /analytics/by-user/plans
GET /analytics/by-user/ask-mode
```

### Health (1)
```
GET /health
```

---

## Key Documentation Files

### v2.0 Documentation

| File | Purpose | Status |
|------|---------|--------|
| docs/DESIGN.md | System architecture v2.0 | COMPLETE |
| docs/FEATURES.md | Feature breakdown v2.0 | COMPLETE |
| docs/TASKS.md | Implementation tasks v2.0 | COMPLETE |
| .claude/DEVELOPMENT.md | Session context (this file) | COMPLETE |

### Specifications

| File | Purpose | Status |
|------|---------|--------|
| specs/openapi/cursor-api.yaml | Cursor API schema | PENDING |
| specs/openapi/github-sim-api.yaml | GitHub simulation API | PENDING |
| tools/data-designer/seed_schema.yaml | Seed file schema | PENDING |

### Claude Code Integration

| File | Purpose |
|------|---------|
| CLAUDE.md | Project instructions |
| .claude/skills/cursor-api-patterns.md | API implementation guide |
| .claude/skills/go-best-practices.md | Go coding standards |
| .claude/skills/model-selection-guide.md | Task-to-model mapping |
| .claude/commands/implement.md | /implement command |
| .claude/commands/next-task.md | /next-task command |

---

## TDD Workflow

### Red-Green-Refactor Cycle

1. **RED**: Write failing test
   ```bash
   go test ./... -v
   # Test should FAIL
   ```

2. **GREEN**: Write minimal code to pass
   ```bash
   go test ./... -v
   # Test should PASS
   ```

3. **REFACTOR**: Clean up while green
   ```bash
   go test ./... -v
   gofmt -s -w .
   golangci-lint run
   ```

### Coverage Target
- **Minimum**: 80%
- **Check**: `go test ./... -cover`

---

## Development Commands

### cursor-sim v2
```bash
cd services/cursor-sim

# Tests
go test ./... -v -cover
make test

# Linting
golangci-lint run
make lint

# Build
go build -o bin/cursor-sim cmd/simulator/main.go
make build

# Run (after implementation)
./bin/cursor-sim --mode=runtime --seed=seed.json --port=8080
```

### Docker
```bash
docker-compose up -d
docker-compose logs -f cursor-sim
docker-compose down
```

---

## Model Selection Guide

| Task Type | Model | Rationale |
|-----------|-------|-----------|
| Type definitions | Haiku | Well-specified, low complexity |
| Validation logic | Sonnet | Requires careful edge case handling |
| Generation engine | Sonnet | Complex Poisson/statistical logic |
| Simple endpoints | Haiku | Straightforward CRUD |
| Complex aggregations | Sonnet | Multiple transformations |
| E2E tests | Sonnet | Integration complexity |

---

## Session Checklist

When starting a new session:

1. [ ] Read this file (DEVELOPMENT.md)
2. [ ] Check current task status in docs/TASKS.md
3. [ ] Identify next task to implement
4. [ ] Select appropriate model per guide
5. [ ] Follow TDD workflow
6. [ ] Update task status when complete

---

## Reference Links

### External
- [Cursor Analytics API](https://cursor.com/docs/account/teams/analytics-api)
- [Cursor AI Code Tracking](https://docs.cursor.com/business/api-reference/ai-code-tracking)

### Internal
- docs/DESIGN.md - Architecture
- docs/FEATURES.md - Feature specs
- docs/TASKS.md - Task breakdown

---

## Notes

### Architecture Decisions

1. **Seed-based generation**: DataDesigner generates dimension data, cursor-sim generates time-series events

2. **Exact API matching**: Response schemas match Cursor documentation exactly for drop-in replacement

3. **Two operation modes**:
   - Runtime: Generate from seed.json
   - Replay: Serve from pre-generated Parquet

4. **JOIN key consistency**: All APIs share commit_sha, user_email, repo_name keys

### Open Items

1. NVIDIA API access for DataDesigner (scipy/faker fallback available)
2. Parquet library selection for replay mode (Phase 3)
3. PostgreSQL schema for cursor-analytics-core (unchanged)

---

**Remember**: Specifications → Tests → Implementation → Refactor

This is the way.
