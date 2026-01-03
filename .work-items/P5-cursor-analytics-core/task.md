# Task Breakdown: cursor-analytics-core

## Overview

**Feature**: cursor-analytics-core (GraphQL Aggregator)
**Total Estimated Hours**: 25-30
**Number of Steps**: 10
**Current Step**: Step 07 - COMPLETE

## Progress Tracker

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| 01 | Project Setup | 2.0 | DONE | 1.5 |
| 02 | Database Schema & Migrations | 3.0 | DONE | 2.5 |
| 03 | cursor-sim REST Client | 2.5 | DONE | 2.0 |
| 04 | Ingestion Worker | 3.5 | NOT_STARTED | - |
| 05 | GraphQL Schema | 2.0 | DONE | 2.0 |
| 06 | Developer Resolvers | 2.5 | DONE | 2.5 |
| 07 | Commit Resolvers | 2.0 | DONE | 2.0 |
| 08 | Metrics Service | 3.0 | NOT_STARTED | - |
| 09 | Dashboard Summary | 2.5 | NOT_STARTED | - |
| 10 | Integration & E2E Tests | 3.0 | NOT_STARTED | - |

## Dependency Graph

```
Step 01 (Project Setup)
    │
    ├── Step 02 (Database)
    │       │
    │       └── Step 04 (Ingestion) ◄── Step 03 (REST Client)
    │               │
    │               ▼
    │       Step 08 (Metrics)
    │               │
    └── Step 05 (GraphQL Schema)
            │
            ├── Step 06 (Developer Resolvers)
            │
            ├── Step 07 (Commit Resolvers)
            │
            └── Step 09 (Dashboard Summary)
                    │
                    └── Step 10 (E2E Tests)
```

## Critical Path

01 → 02 → 03 → 04 → 08 → 06 → 09 → 10

## Step Details

### Step 01: Project Setup ✅ COMPLETE

**Status**: DONE
**Actual Time**: 1.5h

**Tasks**:
- [x] Initialize Node.js project with TypeScript
- [x] Configure ESLint + Prettier
- [x] Set up Jest testing
- [x] Create project structure
- [x] Configure environment variables
- [x] Add npm scripts

**Files Created**:
```
eslint.config.js
.prettierrc
jest.config.js
src/config/index.ts
src/config/index.test.ts
```

**Acceptance Criteria**:
- ✅ `npm run build` compiles successfully
- ✅ `npm run lint` passes
- ✅ `npm test` runs (with placeholder test passing)

---

### Step 02: Database Schema & Migrations ✅ COMPLETE

**Status**: DONE
**Actual Time**: 2.5h

**Tasks**:
- [x] Install and configure Prisma
- [x] Define schema.prisma with all tables
- [x] Create initial migration
- [x] Write seed script for testing
- [x] Test database connection

**Files Created**:
```
prisma/schema.prisma
prisma/migrations/20260103_init/migration.sql
prisma/seed.ts
prisma.config.ts
src/db/client.ts
src/db/__tests__/client.test.ts
.env.example
.gitignore
```

**Acceptance Criteria**:
- ✅ Migration SQL created successfully
- ✅ Prisma Client generated
- ✅ Seed script creates 5 developers + 7 days of events
- ✅ Database client wrapper with health check
- ✅ Tests passing (6 passed, 1 skipped)
- ✅ Build successful
- ✅ Lint clean

---

### Step 03: cursor-sim REST Client ✅ COMPLETE

**Status**: DONE
**Actual Time**: 2.0h

**Tasks**:
- [x] Create typed client for cursor-sim API
- [x] Implement getCommits with pagination
- [x] Implement getTeamMembers
- [x] Handle authentication
- [x] Handle errors and retries
- [x] Unit tests with mocked responses

**Files Created**:
```
src/ingestion/types.ts           # TypeScript type definitions
src/ingestion/client.ts           # CursorSimClient implementation
src/ingestion/__tests__/client.test.ts  # Comprehensive unit tests
```

**Acceptance Criteria**:
- ✅ Client correctly calls cursor-sim endpoints (/teams/members, /analytics/ai-code/commits)
- ✅ Pagination handled correctly with query string parameters
- ✅ Errors handled gracefully with retry logic and client error detection
- ✅ Basic Auth implemented using API key
- ✅ Exponential backoff for retries (capped at 30s)
- ✅ Timeout handling with AbortController
- ✅ All 19 unit tests passing
- ✅ Build successful
- ✅ Lint clean

**Key Features Implemented**:
- Typed responses matching cursor-sim API contract
- Automatic retry on 5xx errors and 429 rate limit
- No retry on 4xx client errors (except 429)
- Query parameter builder for filters and pagination
- Configurable timeout, retry attempts, and retry delay
- Comprehensive test coverage with mocked fetch

---

### Step 04: Ingestion Worker

**Tasks**:
- [ ] Implement polling mechanism
- [ ] Transform cursor-sim data to DB schema
- [ ] Handle incremental syncs
- [ ] Implement deduplication
- [ ] Track sync state
- [ ] Unit tests

**Files to Create**:
```
src/ingestion/worker.ts
src/ingestion/transformer.ts
tests/unit/ingestion/worker.test.ts
```

**Acceptance Criteria**:
- Worker polls at configured interval
- Data correctly transformed and stored
- Duplicates handled

---

### Step 05: GraphQL Schema ✅ COMPLETE

**Status**: DONE
**Actual Time**: 2.0h

**Tasks**:
- [x] Define type definitions
- [x] Set up Apollo Server
- [x] Configure context
- [x] Add health check query
- [x] Write comprehensive tests
- [x] Update index.ts to use new schema

**Files Created**:
```
src/graphql/schema.ts                        # Full GraphQL schema
src/graphql/context.ts                       # Context with DB and REST client
src/graphql/server.ts                        # Apollo Server setup with health resolver
src/graphql/__tests__/server.test.ts         # 14 comprehensive tests
src/graphql/__tests__/context.test.ts        # Context creation tests
```

**Acceptance Criteria**:
- ✅ Apollo Server configured with complete schema
- ✅ Health check query returns status for DB and simulator
- ✅ All GraphQL types defined (Developer, DailyStats, TeamStats, DashboardKPI)
- ✅ DateTime scalar implemented
- ✅ Context includes PrismaClient and CursorSimClient
- ✅ All 14 tests passing
- ✅ Build successful
- ✅ Lint clean

**Key Features Implemented**:
- Complete GraphQL schema with all query types
- Health check resolver with DB and simulator status
- DateTime scalar for proper date handling
- Context factory for request-scoped data
- Comprehensive test coverage (14 tests)
- Error handling and introspection enabled

---

### Step 06: Developer Resolvers ✅ COMPLETE

**Status**: DONE
**Actual Time**: 2.5h

**Tasks**:
- [x] Implement developer query
- [x] Implement developers query with filters
- [x] Add stats field resolver
- [x] Add dailyStats field resolver
- [x] Pagination support
- [x] Unit tests

**Files Created**:
```
src/graphql/resolvers/developer.ts               # Developer resolver implementation
src/graphql/resolvers/__tests__/developer.test.ts # Comprehensive unit tests (13 tests)
src/graphql/resolvers/index.ts                   # Resolver aggregator
```

**Acceptance Criteria**:
- ✅ Single developer query works
- ✅ List query with filters works (team, seniority)
- ✅ Stats calculated correctly (acceptance rate, AI velocity)
- ✅ Daily stats aggregation working
- ✅ Pagination with hasNextPage/hasPreviousPage
- ✅ All 13 unit tests passing
- ✅ Build successful
- ✅ Lint clean

**Key Features Implemented**:
- Query resolvers: `developer(id)`, `developers()` with filtering
- Field resolvers: `Developer.stats()`, `Developer.dailyStats()`
- Date range filtering for stats
- Acceptance rate calculation (null when no suggestions)
- AI velocity calculation (null when no lines)
- Cursor-based pagination with pageInfo
- Comprehensive test coverage with mocked Prisma

---

### Step 07: Commit Resolvers ✅ COMPLETE

**Status**: DONE
**Actual Time**: 2.0h

**Tasks**:
- [x] Implement commits query resolver
- [x] Add filters (userId, team, date range)
- [x] Implement cursor-based pagination
- [x] Support sorting (by timestamp, author)
- [x] Unit tests with mocked Prisma
- [x] Update GraphQL schema with Commit type
- [x] Export commit resolvers in index

**Files Created**:
```
src/graphql/resolvers/commit.ts                     # Commit resolver implementation
src/graphql/resolvers/__tests__/commit.test.ts      # 11 comprehensive unit tests
```

**Files Updated**:
```
src/graphql/schema.ts                               # Added Commit type and commits query
src/graphql/resolvers/index.ts                      # Export commit resolvers
```

**Acceptance Criteria**:
- ✅ Commits query returns paginated usage events (accepted suggestions)
- ✅ All filters work correctly (userId, team, dateRange)
- ✅ Pagination works with hasNextPage/hasPreviousPage
- ✅ Sorting works (timestamp desc/asc, author name)
- ✅ All 11 unit tests passing
- ✅ Total test count: 63 passed (7 test suites)
- ✅ Build successful
- ✅ Lint clean

**Key Features Implemented**:
- Query resolver: `commits()` with filtering and pagination
- Field resolvers: `Commit.timestamp()`, `Commit.author()`
- Filters: userId, team, dateRange (from/to)
- Sorting: by timestamp (default desc) or author name
- Cursor-based pagination with pageInfo
- Maps usage_events (eventType = 'cpp_suggestion_accepted') to Commit type
- Comprehensive test coverage with mocked Prisma client

---

### Step 08: Metrics Service

**Tasks**:
- [ ] Implement AI velocity calculation
- [ ] Implement acceptance rate calculation
- [ ] Team-level aggregations
- [ ] Time-range filtering
- [ ] Unit tests with fixture data

**Files to Create**:
```
src/services/metrics.ts
tests/unit/services/metrics.test.ts
tests/fixtures/commits.ts
```

**Acceptance Criteria**:
- AI velocity calculated correctly
- Edge cases handled (division by zero)
- Time filtering works

---

### Step 09: Dashboard Summary

**Tasks**:
- [ ] Implement dashboardSummary resolver
- [ ] Aggregate all KPIs
- [ ] Top contributors query
- [ ] Team stats resolver
- [ ] Unit tests

**Files to Create**:
```
src/graphql/resolvers/dashboard.ts
tests/unit/graphql/dashboard.test.ts
```

**Acceptance Criteria**:
- Dashboard query returns all KPIs
- Top contributors sorted correctly
- Team stats work

---

### Step 10: Integration & E2E Tests

**Tasks**:
- [ ] Integration test: ingestion worker with test DB
- [ ] Integration test: full GraphQL queries
- [ ] E2E test: cursor-sim → core → queries
- [ ] Performance test: 10k commits
- [ ] Document API with examples

**Files to Create**:
```
tests/integration/ingestion.test.ts
tests/integration/graphql.test.ts
tests/e2e/full-pipeline.test.ts
```

**Acceptance Criteria**:
- All integration tests pass
- E2E pipeline works
- Coverage > 80%

---

## Model Recommendations

| Step | Model | Rationale |
|------|-------|-----------|
| 01, 02, 05 | Haiku | Boilerplate, well-specified |
| 03, 04 | Sonnet | API client, polling complexity |
| 06, 07 | Haiku | Standard resolver patterns |
| 08, 09, 10 | Sonnet | Metrics logic, integration |

## TDD Checklist (Per Step)

- [ ] Read step details and acceptance criteria
- [ ] Write failing test(s) for the step
- [ ] Run tests, confirm RED
- [ ] Implement minimal code to pass
- [ ] Run tests, confirm GREEN
- [ ] Refactor while green
- [ ] Run linter (npm run lint)
- [ ] Update step status to DONE
- [ ] Commit with time tracking
