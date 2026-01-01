# Development Session Context

**Last Updated**: January 2026
**Current Phase**: P0 - Make Runnable
**Primary Focus**: cursor-sim (Cursor API Simulator)

---

## Current Status

### Project State
- **Implementation**: 0% (All specs, no code)
- **Specifications**: 100% complete
- **Testing**: Not started
- **Infrastructure**: Docker Compose defined, Dockerfiles pending

### Recent Work

#### Completed
1. **Comprehensive Specifications**
   - services/cursor-sim/SPEC.md v2.0.0 (1018 lines) - matches actual Cursor API
   - specs/api/graphql-schema.graphql (392 lines)
   - docs/DESIGN.md, docs/TASKS.md, docs/USER_STORIES.md

2. **Claude Code SDD Structure**
   - Skills: cursor-api-patterns.md, go-best-practices.md
   - Commands: /spec, /start-feature, /verify, /next-task
   - This DEVELOPMENT.md file

3. **Project Review**
   - PROJECT_REVIEW.md - comprehensive gap analysis
   - Identified simulator API mismatch with Cursor (now fixed in SPEC v2.0.0)
   - Documented .claude folder structure needs

#### In Progress
- Implementing Claude Code native SDD methodology
- Setting up .claude/plans/ directory

---

## Active Work Item

**None** - Ready to start P0 tasks

Suggested next step: Execute P0_MAKERUNNABLE.md tasks to create minimal scaffolding

---

## Service Focus: cursor-sim

### Purpose
Go-based CLI tool and REST API server that simulates Cursor Business API endpoints with realistic synthetic data.

### Key Decisions Made

1. **API Compatibility** (CRITICAL)
   - Simulator MUST match actual Cursor API endpoints:
     - `/v1/analytics/ai-code/commits`
     - `/v1/analytics/ai-code/changes`
     - `/v1/analytics/team/agent-edits`
     - `/v1/analytics/team/tabs`
     - `/v1/analytics/team/dau`
     - `/v1/analytics/team/models`
     - `/v1/teams/members`

2. **Organizational Hierarchy**
   - Developer model includes: Region → Division → Group → Team
   - Regions: US (50%), EU (35%), APAC (15%)
   - Divisions: AGS, AT, ST
   - Groups: TMOBILE, ATANT
   - Teams: dev (75%), support (25%)

3. **Configuration**
   - Accepts JSON config file (not just CLI flags)
   - Generates fake API credentials (Basic Auth)
   - Supports break conditions (PR count limit or infinite)

4. **Interactive CLI**
   - Live dashboard with stats
   - Ctrl+S: Soft stop (finish in-flight events)
   - Ctrl+E: Export database to JSON
   - Ctrl+C: Immediate quit

5. **Storage**
   - In-memory only (sync.Map or go-memdb)
   - Acceptable data loss on restart for MVP
   - Thread-safe concurrent access

### Tech Stack
- **Language**: Go 1.21+
- **HTTP**: net/http (standard library)
- **Testing**: go test + testify + mockery
- **Storage**: sync.Map or go-memdb
- **Distribution**: Poisson for realistic event timing

---

## Key Documentation Files

### Specifications (Single Source of Truth)
| File | Purpose | Status |
|------|---------|--------|
| services/cursor-sim/SPEC.md | Complete service specification | ✓ v2.0.0 |
| specs/api/graphql-schema.graphql | GraphQL contract | ✓ Complete |
| services/cursor-analytics-core/SPEC.md | Aggregator spec | Pending |
| services/cursor-viz-spa/SPEC.md | Frontend spec | Pending |

### Design Documents
| File | Purpose | Status |
|------|---------|--------|
| docs/DESIGN.md | System architecture | ✓ v1.0.0 |
| docs/FEATURES.md | Feature breakdown | ✓ Complete |
| docs/TESTING_STRATEGY.md | TDD approach | ✓ Complete |

### Implementation Guides
| File | Purpose | Status |
|------|---------|--------|
| docs/TASKS.md | Task breakdown (895 lines) | ✓ Complete |
| docs/USER_STORIES.md | Acceptance criteria (712 lines) | ✓ Complete |
| P0_MAKERUNNABLE.md | 8 tasks to make repo runnable | ✓ Complete |
| PROJECT_REVIEW.md | Gap analysis | ✓ Complete |

### Claude Code Integration
| File | Purpose | Status |
|------|---------|--------|
| CLAUDE.md | Project instructions for AI | ✓ Complete |
| .claude/skills/cursor-api-patterns.md | API implementation guide | ✓ Complete |
| .claude/skills/go-best-practices.md | Go coding standards | ✓ Complete |
| .claude/commands/spec.md | /spec command | ✓ Complete |
| .claude/commands/start-feature.md | /start-feature command | ✓ Complete |
| .claude/commands/verify.md | /verify command | ✓ Complete |
| .claude/commands/next-task.md | /next-task command | ✓ Complete |
| .claude/DEVELOPMENT.md | Session context (this file) | ✓ Complete |

---

## Next Steps (Priority Order)

### P0: Make Runnable (Est: 2 hours)
These 8 tasks from P0_MAKERUNNABLE.md must be completed first:

1. **P0.1**: Create Go scaffolding (go.mod, main.go, Dockerfile)
2. **P0.2**: Create TypeScript scaffolding (package.json, tsconfig.json, Dockerfile)
3. **P0.3**: Create React scaffolding (Vite setup, Dockerfile)
4. **P0.4**: Create .env.example
5. **P0.5**: Update docker-compose.yml with build contexts
6. **P0.6**: Update Makefile with new targets
7. **P0.7**: Verify `docker-compose up` works
8. **P0.8**: Implement minimal happy path (sim → core → viz)

### P1: cursor-sim Core Features (Est: 1-2 weeks)
From docs/TASKS.md:

- TASK-SIM-001: Initialize Go Project Structure
- TASK-SIM-002: Implement CLI Flag Parsing
- TASK-SIM-003: Implement Developer Profile Generator
- TASK-SIM-004: Implement Event Generation Engine
- TASK-SIM-005: Implement In-Memory Storage
- TASK-SIM-006: Implement REST API Handlers
- TASK-SIM-007: Wire Up Main Application

### P2: Integration & Polish
- Implement actual .claude/hooks/ scripts
- Add rate limiting
- Add Basic Auth simulation
- Add OpenTelemetry traces
- Add Prometheus metrics

---

## TDD Workflow Reminder

### Red-Green-Refactor Cycle

1. **RED**: Write a failing test
   ```bash
   go test ./... -v
   # Test should FAIL with clear error message
   ```

2. **GREEN**: Write minimal code to pass
   ```bash
   go test ./... -v
   # Test should PASS
   ```

3. **REFACTOR**: Clean up while tests stay green
   ```bash
   go test ./... -v
   # Tests should STILL PASS
   gofmt -s -w .
   golangci-lint run
   ```

### Coverage Target
- **Minimum**: 80% for all services
- **Check**: `go test ./... -cover`
- **Report**: `go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out`

---

## Common Development Commands

### cursor-sim (Go)
```bash
# Navigate to service
cd services/cursor-sim

# Initialize module (first time)
go mod init github.com/yourusername/cursor-analytics-platform/services/cursor-sim

# Run tests
go test ./...
go test ./... -v -cover

# Run linter
golangci-lint run

# Format code
gofmt -s -w .

# Run locally
go run cmd/simulator/main.go --config=config.example.json

# Build binary
go build -o bin/cursor-sim cmd/simulator/main.go
```

### Docker (All Services)
```bash
# Build all services
docker-compose build

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f cursor-sim
docker-compose logs -f cursor-analytics-core
docker-compose logs -f cursor-viz-spa

# Stop all services
docker-compose down

# Restart single service
docker-compose restart cursor-sim
```

### Makefile Shortcuts (when P0.6 complete)
```bash
make build          # Build all Docker images
make up             # Start all services
make down           # Stop all services
make test           # Run all tests
make test-coverage  # Run tests with coverage
make logs           # Tail all logs
```

---

## Open Questions / Decisions Needed

1. **Database Schema**: PostgreSQL schema for cursor-analytics-core not yet defined
2. **Authentication**: Should simulator validate credentials or just echo them?
3. **Rate Limiting**: Implement immediately or defer to P2?
4. **Observability**: OpenTelemetry setup - P0, P1, or P2?
5. **CI/CD**: GitHub Actions or other? (Deferred to P2)

---

## Reference Links

### External Documentation
- [Cursor Business API - AI Code Tracking](https://docs.cursor.com/business/api-reference/ai-code-tracking)
- [Cursor Business API - Analytics](https://docs.cursor.com/business/api-reference/analytics)

### Internal Specs
- Simulator: services/cursor-sim/SPEC.md
- GraphQL: specs/api/graphql-schema.graphql
- Tasks: docs/TASKS.md
- Stories: docs/USER_STORIES.md

---

## Session Notes

### Recent Clarifications

1. **API Endpoint Alignment** (RESOLVED)
   - User confirmed simulator must match actual Cursor API
   - SPEC.md v2.0.0 now uses correct endpoints

2. **SDD Methodology** (RESOLVED)
   - Confirmed Claude Code uses Skills (not .mdc rules)
   - Skills and commands now implemented

3. **Configuration Approach** (RESOLVED)
   - Accept JSON config file (not just CLI flags)
   - Example config in SPEC.md lines 287-332

4. **Interactive CLI** (RESOLVED)
   - Ctrl+S, Ctrl+E, Ctrl+C controls required
   - Live dashboard with org hierarchy stats

---

## Tips for AI Assistants

When starting a new session:

1. **Read this file first** to understand current state
2. **Check `.claude/plans/active`** symlink for active work item
3. **Run `/next-task`** to see what to work on
4. **Read relevant SPEC.md** before writing any code
5. **Write tests first** - this is non-negotiable
6. **Use Skills** - invoke cursor-api-patterns or go-best-practices as needed
7. **Update this file** when making significant progress or decisions

### Quick Context Check
```bash
# What am I working on?
ls -la .claude/plans/active

# What's next?
/next-task cursor-sim

# What does the spec say?
/spec cursor-sim
```

---

**Remember**: Specifications → Tests → Implementation → Refactor → Documentation

This is the way.
