# CLAUDE.md - Cursor Usage Analytics Platform

## Project Overview

This is a **Cursor Usage Analytics Platform** consisting of three decoupled microservices following an ETL (Extract, Transform, Load) architecture pattern. The platform simulates and analyzes AI coding assistant usage metrics to provide insights into developer productivity and AI tool adoption.

## Architecture Summary

```
┌─────────────────┐     ┌──────────────────────┐     ┌─────────────────┐
│   cursor-sim    │────▶│ cursor-analytics-core│────▶│  cursor-viz-spa │
│   (Go + REST)   │     │   (TS + GraphQL)     │     │  (React + Vite) │
│   Port: 8080    │     │   Port: 4000         │     │   Port: 3000    │
└─────────────────┘     └──────────────────────┘     └─────────────────┘
     Simulator              Aggregator                  Dashboard
     (Extract)              (Transform)                  (Load/View)
```

## Services

### Service A: cursor-sim (Cursor API Simulator)
- **Language**: Go 1.21+
- **Type**: CLI Tool + REST API Server
- **Storage**: In-memory (sync.Map or go-memdb)
- **Port**: 8080
- **Purpose**: Generates synthetic Cursor usage data mimicking real developer behavior

### Service B: cursor-analytics-core (Aggregator Service)
- **Language**: TypeScript (Node.js 20+)
- **Framework**: Apollo Server (GraphQL)
- **Database**: PostgreSQL 15+
- **Port**: 4000
- **Purpose**: Ingests, normalizes, and calculates KPIs from simulator data

### Service C: cursor-viz-spa (Frontend Dashboard)
- **Language**: TypeScript
- **Framework**: React 18+ with Vite
- **State**: TanStack Query (React Query)
- **Charts**: Recharts
- **Port**: 3000
- **Purpose**: Visualizes analytics through interactive dashboards

## Technology Stack

### Languages & Frameworks
- Go 1.21+ (cursor-sim)
- TypeScript 5.3+ (cursor-analytics-core, cursor-viz-spa)
- Node.js 20 LTS
- React 18.2+

### Testing Stack (Recommended)
- **Go**: `go test` + `testify` + `mockery`
- **TypeScript Backend**: `jest` + `ts-jest` + `supertest` + `@graphql-tools/mock`
- **React Frontend**: `vitest` + `@testing-library/react` + `msw` (Mock Service Worker)
- **E2E**: `playwright`
- **Coverage**: 80% minimum for all services

### Infrastructure
- Docker + Docker Compose (local development)
- PostgreSQL 15 (analytics-core database)
- No external cloud dependencies

## Quick Start Commands

```bash
# Start all services
docker-compose up -d

# Run cursor-sim standalone
cd services/cursor-sim && go run . --port=8080 --developers=50 --velocity=high

# Run analytics-core
cd services/cursor-analytics-core && npm run dev

# Run frontend
cd services/cursor-viz-spa && npm run dev

# Run all tests
make test

# Run tests with coverage
make test-coverage
```

## Development Workflow

### Documentation Hierarchy (SDD Compliant)

**Source of Truth** (in priority order):

| Location | Purpose | When to Update |
|----------|---------|----------------|
| `services/{service}/SPEC.md` | Technical specification | When API/behavior changes |
| `.work-items/{feature}/` | Active work tracking | During feature development |
| `.claude/DEVELOPMENT.md` | Session context | Each session start/end |
| `.claude/plans/active` | Symlink to current task | When switching features |

**Reference Documents** (read-only context):

| Location | Purpose |
|----------|---------|
| `docs/DESIGN.md` | System architecture overview |
| `docs/FEATURES.md` | Feature breakdown |
| `docs/TASKS.md` | Task overview |
| `docs/USER_STORIES.md` | User story reference |

### Work Items Structure

Each feature has a dedicated work-items directory:

```
.work-items/{feature}/
├── user-story.md      # Requirements and acceptance criteria
├── design.md          # Technical design decisions
├── task.md            # Task breakdown with status tracking
└── {NN}_step.md       # (optional) Detailed step files
```

### Spec-Driven Development Process
1. **Read the spec**: Check `services/{service}/SPEC.md` first
2. **Check work items**: Review `.work-items/{feature}/task.md` for current status
3. **Write tests first**: Follow TDD - tests define the contract
4. **Implement minimally**: Write just enough code to pass tests
5. **Refactor**: Clean up while keeping tests green
6. **Update status**: Mark tasks complete in task.md

### Before Implementing Any Feature
1. Read `services/{service}/SPEC.md` for technical specification
2. Check `.work-items/{feature}/task.md` for implementation status
3. Review `.work-items/{feature}/design.md` for architecture decisions
4. Write failing tests based on acceptance criteria
5. Implement to make tests pass

### Claude Code Integration

This project uses Claude Code's Skills feature for enhanced development workflow.

#### Available Skills
Skills provide specialized knowledge. Reference them explicitly in your requests:

- **cursor-api-patterns** - Cursor Business API implementation guide (response formats, auth, pagination, rate limiting, CSV exports)
- **go-best-practices** - Go coding standards (project structure, error handling, concurrency, testing, HTTP handlers)
- **model-selection-guide** - Task-to-model mapping for cost optimization (Haiku vs Sonnet vs Opus)
- **sdd-checklist** - **CRITICAL**: Spec-Driven Development enforcement checklist (commit after tests pass, update progress docs)

**Usage:** Reference skills directly:
```
"Following cursor-api-patterns.md, implement the health check endpoint"
"Using go-best-practices.md, create the Developer struct"
"Based on model-selection-guide.md, use Haiku for TASK-SIM-003"
```

#### Session Context
- `.claude/DEVELOPMENT.md` - Current development status, recent work, active focus, next steps
- `.claude/MODEL_SELECTION_SUMMARY.md` - Cost optimization guide
- Read these files at session start to understand project state

See `.claude/README.md` for complete usage guide.

## Key Documentation Files

### Source of Truth (Read First!)
| File | Purpose |
|------|---------|
| `.claude/DEVELOPMENT.md` | **START HERE** - Current session context, active work, next steps |
| `services/cursor-sim/SPEC.md` | cursor-sim v2 specification (Phase 1 Complete) |
| `.work-items/cursor-sim-v2/task.md` | Phase 1 task tracking (COMPLETE) |
| `.work-items/cursor-sim-phase2/task.md` | Phase 2 task tracking (NOT_STARTED) |
| `.work-items/cursor-analytics-core/task.md` | Aggregator task tracking (NOT_STARTED) |

### Service Specifications
| File | Purpose |
|------|---------|
| `services/cursor-sim/SPEC.md` | Complete cursor-sim specification (v2.0.0) |
| `services/cursor-analytics-core/SPEC.md` | Aggregator service spec (pending) |
| `services/cursor-viz-spa/SPEC.md` | Frontend dashboard spec (pending) |
| `specs/api/graphql-schema.graphql` | GraphQL schema contract |

### Reference Documents (Background Context)
| File | Purpose |
|------|---------|
| `docs/DESIGN.md` | System architecture overview |
| `docs/FEATURES.md` | Feature breakdown by service |
| `docs/USER_STORIES.md` | User stories reference |
| `docs/TASKS.md` | Task overview reference |
| `docs/TESTING_STRATEGY.md` | TDD approach and testing guidelines |

### Cursor API Reference (Design & Testing Source of Truth)
| File | Purpose |
|------|---------|
| `docs/api-reference/cursor_overview.md` | Authentication, rate limits, caching, best practices |
| `docs/api-reference/cursor_admin.md` | Admin API - Team management, usage data, spending |
| `docs/api-reference/cursor_analytics.md` | Analytics API - Team metrics, DAU, model usage |
| `docs/api-reference/cursor_codetrack.md` | AI Code Tracking API - Per-commit metrics (Enterprise) |
| `docs/api-reference/cursor_agents.md` | Cloud Agents API - Programmatic agent management |
| `docs/api-reference/REST_API.md` | Simulator implementation reference |
| `specs/openapi/cursor-api.yaml` | OpenAPI 3.1 specification for code generation |

### Claude Code Integration
| File | Purpose |
|------|---------|
| `.claude/README.md` | Claude Code integration guide |
| `.claude/DEVELOPMENT.md` | Session context (read first!) |
| `.claude/MODEL_SELECTION_SUMMARY.md` | Model optimization strategies |
| `.claude/skills/cursor-api-patterns.md` | Cursor API implementation guide |
| `.claude/skills/go-best-practices.md` | Go coding standards |
| `.claude/skills/model-selection-guide.md` | Task-to-model mapping |
| `.claude/skills/spec-driven-development.md` | SDD methodology |

## Code Style Guidelines

### Go (cursor-sim)
- Follow `gofmt` and `golint`
- Use meaningful variable names
- Prefer composition over inheritance
- Document exported functions

### TypeScript (analytics-core, viz-spa)
- Strict mode enabled
- ESLint + Prettier
- Prefer functional components in React
- Use explicit return types on functions

## Environment Variables

```bash
# cursor-sim
CURSOR_SIM_PORT=8080
CURSOR_SIM_DEVELOPERS=50
CURSOR_SIM_VELOCITY=high
CURSOR_SIM_FLUCTUATION=0.2

# cursor-analytics-core
DATABASE_URL=postgresql://user:pass@localhost:5432/cursor_analytics
SIMULATOR_URL=http://localhost:8080
POLL_INTERVAL_MS=60000
GRAPHQL_PORT=4000

# cursor-viz-spa
VITE_GRAPHQL_URL=http://localhost:4000/graphql
```

## API Contracts

### Simulator → Aggregator (REST)
**cursor-sim mimics actual Cursor Business API endpoints:**

#### AI Code Tracking API
- `GET /v1/analytics/ai-code/commits` - AI-generated commits with line counts
- `GET /v1/analytics/ai-code/changes` - Individual code changes (tab/composer)

#### Analytics API
- `GET /v1/analytics/team/agent-edits` - Daily agent edit counts
- `GET /v1/analytics/team/tabs` - Tab completion metrics
- `GET /v1/analytics/team/dau` - Daily active users
- `GET /v1/analytics/team/models` - Model usage stats

#### Team Management API
- `GET /v1/teams/members` - List team members

#### Utilities
- `GET /v1/health` - Health check

**Authentication**: Basic Auth (simulated credentials)
**Response Format**: `{ "data": [...], "pagination": {...}, "params": {...} }`

See `services/cursor-sim/SPEC.md` for complete API documentation.

### Aggregator → Frontend (GraphQL)
- See `specs/api/graphql-schema.graphql` for full schema
- Key queries: `getDevStats`, `getTeamStats`, `getDashboardSummary`

## Common Tasks for AI Assistants

### When Asked to Implement a Feature
1. First read the spec: `cat services/{service}/SPEC.md`
2. Check user stories: `cat docs/USER_STORIES.md | grep -A 20 "{feature}"`
3. Write tests first following TDD
4. Implement the feature
5. Run tests: `make test`

### When Asked to Fix a Bug
1. Write a failing test that reproduces the bug
2. Fix the code to make the test pass
3. Verify no regression: `make test`

### When Asked to Add an API Endpoint
1. Update the spec file first
2. Add to `specs/api/` OpenAPI or GraphQL schema
3. Write integration tests
4. Implement the endpoint
5. Update documentation

## Project Conventions

### File Naming
- Go: `snake_case.go`
- TypeScript: `camelCase.ts` for utilities, `PascalCase.tsx` for components
- Tests: `*.test.ts`, `*.test.go`, `*.spec.tsx`
- Specs: `SPEC.md` in each service directory

### Git Workflow
- Feature branches: `feature/{service}/{feature-name}`
- Bug fixes: `fix/{service}/{bug-description}`
- Conventional commits: `feat:`, `fix:`, `docs:`, `test:`, `refactor:`

### PR Requirements
- All tests passing
- Coverage not decreased
- Spec updated if behavior changes
- At least one reviewer approval

## Troubleshooting

### Common Issues
1. **Port already in use**: Check `docker-compose ps` and stop conflicting containers
2. **Database connection failed**: Ensure PostgreSQL is running and migrations applied
3. **GraphQL type errors**: Regenerate types with `npm run codegen`
4. **Go module issues**: Run `go mod tidy` in cursor-sim directory
