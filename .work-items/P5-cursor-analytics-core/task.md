# Task Breakdown: cursor-analytics-core

## Overview

**Feature**: cursor-analytics-core (GraphQL Aggregator)
**Total Estimated Hours**: 25-30
**Number of Steps**: 10
**Current Step**: Step 10 - COMPLETE ✅

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
| 08 | Metrics Service | 3.0 | DONE | 3.0 |
| 09 | Dashboard Summary | 2.5 | DONE | 2.5 |
| 10 | Integration & E2E Tests | 3.0 | DONE | 3.0 |

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

### Step 09: Dashboard Summary ✅ COMPLETE

**Status**: DONE
**Actual Time**: 2.5h

**Tasks**:
- [x] Implement dashboardSummary resolver
- [x] Aggregate all KPIs
- [x] Top contributors query (top performer)
- [x] Team stats resolver
- [x] Teams query resolver
- [x] Unit tests

**Files Created**:
```
src/graphql/resolvers/dashboard.ts
src/graphql/resolvers/__tests__/dashboard.test.ts
```

**Files Updated**:
```
src/graphql/resolvers/index.ts (exported dashboard resolvers)
SPEC.md (Step 09 marked complete)
```

**Acceptance Criteria**:
- ✅ Dashboard query returns all KPIs (totalDevelopers, activeDevelopers, overallAcceptanceRate, etc.)
- ✅ Top performer identified by AI lines contributed
- ✅ Team stats work with weighted averages
- ✅ Teams query returns all teams with statistics
- ✅ All 14 unit tests passing
- ✅ Total test count: 101 passed
- ✅ Build successful
- ✅ Lint clean (dashboard.ts)

**Key Features Implemented**:
- Query resolvers: dashboardSummary(), teamStats(), teams()
- Date range support with presets (TODAY, LAST_7_DAYS, etc.)
- Weighted team acceptance rate (activity-based, not simple average)
- Weighted team AI velocity calculation
- Daily trend calculation for time-series visualization
- Team comparison across organization
- Top performer identification per team

---

### Step 10: Integration & E2E Tests ✅ COMPLETE

**Status**: DONE
**Actual Time**: 3.0h

**Tasks**:
- [x] Integration test infrastructure (test DB setup, utilities)
- [x] Integration test: full GraphQL queries with real database
- [x] E2E test: cursor-sim → analytics-core → queries
- [x] Performance test: 10k commits (10,500+ events, 50 devs, 30 days)
- [x] Document API with comprehensive examples

**Files Created**:
```
src/__tests__/integration/setup.ts                  # Test infrastructure
src/__tests__/integration/graphql.test.ts           # 13 integration tests
src/__tests__/e2e/full-pipeline.test.ts             # 6 E2E tests
src/__tests__/performance/large-dataset.test.ts     # 7 performance tests
API.md                                              # Complete API documentation
```

**Files Updated**:
```
jest.config.js                                      # Coverage exclusions
src/services/__tests__/metrics.test.ts              # Added preset tests
SPEC.md                                             # Testing section added
```

**Acceptance Criteria**:
- ✅ All integration tests pass (13 tests)
- ✅ E2E pipeline works (6 tests with realistic data)
- ✅ Performance tests verify < 2000ms for dashboard (7 tests)
- ✅ Coverage > 80% achieved (91.49%)
- ✅ API documentation complete with examples

**Test Summary**:
- Unit tests: 107 tests, 91.49% coverage
- Integration tests: 13 tests (GraphQL + DB)
- E2E tests: 6 tests (full pipeline)
- Performance tests: 7 tests (large datasets)
- Total: 133+ tests passing

**Key Features Implemented**:
- Test setup utilities (createTestDb, seedTestData, cleanupDb)
- Integration tests for all GraphQL queries
- E2E tests for complete data flow
- Performance benchmarks with timing assertions
- Comprehensive API documentation with query examples
- Testing section in SPEC.md documenting all test types

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
