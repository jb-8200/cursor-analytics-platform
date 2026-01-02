# Development Session Context

**Last Updated**: January 2, 2026
**Current Phase**: P1 - Core Functionality
**Primary Focus**: cursor-sim (Cursor API Simulator)

---

## Current Status

### Project State
- **Implementation**: 25% (TASK-SIM-001 ✓, TASK-SIM-002 ✓)
- **Specifications**: 100% complete
- **Testing**: 80.4% coverage (config package)
- **Infrastructure**: Docker Compose ready, multi-stage Dockerfiles complete

### Recent Work

#### Completed (January 2, 2026)

1. **TASK-SIM-001: Initialize Go Project Structure** ✓
   - Created standard Go project layout (cmd/, internal/)
   - Implemented configuration package with validation
   - Added basic domain models (Developer, errors)
   - Created Makefile with build automation
   - Added .golangci.yml linter configuration
   - Multi-stage Dockerfile for optimized builds
   - Comprehensive README documentation
   - Test coverage: 93.3% (initial) → 80.4% (with flag parsing)

2. **TASK-SIM-002: Implement CLI Flag Parsing** ✓
   - ParseFlags() using standard flag package
   - JSON configuration file support
   - Comprehensive validation with helpful error messages
   - Custom --help output with examples
   - Integrated with main.go
   - Created config.example.json template
   - All tests passing (80.4% coverage)

3. **Previous Milestones**
   - Comprehensive Specifications (SPEC.md v2.0.0)
   - Claude Code SDD Structure (skills, commands)
   - Project Review (gap analysis complete)

#### In Progress
- None (ready for TASK-SIM-003)

---

## Active Work Item

**Next**: TASK-SIM-003 - Implement Developer Profile Generator

**Recommended Model**: Haiku ⚡ (well-specified struct from SPEC.md:145-250)
**Estimated Time**: 4 hours
**Dependencies**: TASK-SIM-001 ✓, TASK-SIM-002 ✓

**Objective**: Implement realistic developer profile generation with:
- Organizational hierarchy (Region → Division → Group → Team)
- Seniority distribution (20% junior, 50% mid, 30% senior)
- Acceptance rate correlation with seniority
- Deterministic name generation using seed

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

### P1: cursor-sim Core Features (In Progress)
From docs/TASKS.md:

- ✅ TASK-SIM-001: Initialize Go Project Structure
- ✅ TASK-SIM-002: Implement CLI Flag Parsing
- **→ TASK-SIM-003: Implement Developer Profile Generator** (Next)
- TASK-SIM-004: Implement Event Generation Engine
- TASK-SIM-005: Implement In-Memory Storage
- TASK-SIM-006: Implement REST API Handlers
- TASK-SIM-007: Wire Up Main Application

**Progress**: 2/7 tasks complete (29%)

### P0: Make Runnable (Partial)
Basic scaffolding complete:
- ✅ P0.1: Go scaffolding (go.mod, main.go, Dockerfile)
- ⏸️ P0.2-P0.8: Other services (deferred to focus on cursor-sim)

### P2: Integration & Polish (Future)
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

1. **TDD Workflow Established** (January 2, 2026)
   - Successfully followed Red-Green-Refactor cycle
   - Wrote tests first, implemented to pass
   - Achieved 80.4% test coverage on config package
   - All validation scenarios covered

2. **Go Best Practices Applied**
   - Standard project layout (cmd/, internal/)
   - Table-driven tests with testify
   - Error wrapping with context
   - Proper package organization

3. **Configuration System Complete**
   - CLI flags using standard library
   - JSON file support implemented
   - Validation with helpful error messages
   - --help output with examples

4. **Build Infrastructure Ready**
   - Makefile with all necessary targets
   - golangci-lint configuration
   - Multi-stage Docker builds
   - Comprehensive README

### Earlier Clarifications

1. **API Endpoint Alignment** (RESOLVED)
   - User confirmed simulator must match actual Cursor API
   - SPEC.md v2.0.0 now uses correct endpoints

2. **SDD Methodology** (RESOLVED)
   - Confirmed Claude Code uses Skills (not .mdc rules)
   - Skills and commands now implemented

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
