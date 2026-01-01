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

### Spec-Driven Development Process
1. **Read the spec**: Check `specs/` for feature specifications
2. **Write tests first**: Follow TDD - tests define the contract
3. **Implement minimally**: Write just enough code to pass tests
4. **Refactor**: Clean up while keeping tests green
5. **Document**: Update specs if behavior changes

### Before Implementing Any Feature
1. Read the relevant `SPEC.md` in the service directory
2. Check `docs/USER_STORIES.md` for acceptance criteria
3. Review `docs/TASKS.md` for implementation checklist
4. Write failing tests based on the spec
5. Implement to make tests pass

### Claude Code Integration

This project uses Claude Code's native features for enhanced development workflow:

#### Available Skills
Skills provide specialized knowledge for specific technical domains. Use these when working on related code:

- **cursor-api-patterns** - Comprehensive guide for implementing Cursor Business API-compatible endpoints (response formats, auth, pagination, rate limiting, CSV exports)
- **go-best-practices** - Go coding standards for cursor-sim (project structure, error handling, concurrency, testing, HTTP handlers)

To use a skill, simply reference it in conversation (Claude loads skills automatically).

#### Available Slash Commands
Commands automate common development workflows:

- `/spec [service-name]` - Display service specification (cursor-sim, cursor-analytics-core, cursor-viz-spa)
- `/start-feature [name]` - Create SDD work item structure (.work-items/{nn}-{name}/ with user-story.md, design.md, task.md, test-plan.md)
- `/verify [service-name]` - Check spec-test-code alignment, run coverage, find gaps
- `/next-task [service-name]` - Show next actionable task from implementation plan

#### Session Context
- `.claude/DEVELOPMENT.md` - Current development status, recent work, active focus, next steps, open questions
- Read this file at session start to understand project state

## Key Documentation Files

### Project Documentation
| File | Purpose |
|------|---------|
| `.claude/DEVELOPMENT.md` | **START HERE** - Current session context, active work, next steps |
| `docs/DESIGN.md` | System architecture and technical decisions |
| `docs/FEATURES.md` | Feature breakdown by service |
| `docs/USER_STORIES.md` | User stories with acceptance criteria (712 lines) |
| `docs/TASKS.md` | Implementation tasks and checklists (895 lines) |
| `docs/TESTING_STRATEGY.md` | TDD approach and testing guidelines |
| `docs/API_REFERENCE.md` | Cursor Business API documentation |
| `PROJECT_REVIEW.md` | Comprehensive gap analysis and project status |
| `P0_MAKERUNNABLE.md` | 8 tasks to make repo runnable (MVP scaffolding) |

### Service Specifications
| File | Purpose |
|------|---------|
| `services/cursor-sim/SPEC.md` | Complete cursor-sim specification (v2.0.0, 1018 lines) |
| `services/cursor-analytics-core/SPEC.md` | Aggregator service spec (pending) |
| `services/cursor-viz-spa/SPEC.md` | Frontend dashboard spec (pending) |
| `specs/api/graphql-schema.graphql` | GraphQL schema contract (392 lines) |

### Claude Code Integration
| File | Purpose |
|------|---------|
| `.claude/skills/cursor-api-patterns.md` | Cursor API implementation guide |
| `.claude/skills/go-best-practices.md` | Go coding standards |
| `.claude/commands/spec.md` | `/spec` command definition |
| `.claude/commands/start-feature.md` | `/start-feature` command definition |
| `.claude/commands/verify.md` | `/verify` command definition |
| `.claude/commands/next-task.md` | `/next-task` command definition |

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
