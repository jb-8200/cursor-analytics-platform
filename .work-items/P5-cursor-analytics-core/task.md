# Task Breakdown: cursor-analytics-core

## Overview

**Feature**: cursor-analytics-core (GraphQL Aggregator)
**Total Estimated Hours**: 25-30
**Number of Steps**: 10
**Current Step**: Step 03 - COMPLETE

## Progress Tracker

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| 01 | Project Setup | 2.0 | DONE | 1.5 |
| 02 | Database Schema & Migrations | 3.0 | DONE | 2.5 |
| 03 | cursor-sim REST Client | 2.5 | DONE | 2.0 |
| 04 | Ingestion Worker | 3.5 | NOT_STARTED | - |
| 05 | GraphQL Schema | 2.0 | NOT_STARTED | - |
| 06 | Developer Resolvers | 2.5 | NOT_STARTED | - |
| 07 | Commit Resolvers | 2.0 | NOT_STARTED | - |
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

### Step 05: GraphQL Schema

**Tasks**:
- [ ] Define type definitions
- [ ] Set up Apollo Server
- [ ] Configure context
- [ ] Add health check query
- [ ] Generate TypeScript types

**Files to Create**:
```
src/graphql/schema.ts
src/graphql/context.ts
src/graphql/server.ts
```

**Acceptance Criteria**:
- Apollo Server starts on port 4000
- GraphQL Playground accessible
- Health query returns ok

---

### Step 06: Developer Resolvers

**Tasks**:
- [ ] Implement developer query
- [ ] Implement developers query with filters
- [ ] Add stats field resolver
- [ ] Pagination support
- [ ] Unit tests

**Files to Create**:
```
src/graphql/resolvers/developer.ts
tests/unit/graphql/developer.test.ts
```

**Acceptance Criteria**:
- Single developer query works
- List query with filters works
- Stats calculated correctly

---

### Step 07: Commit Resolvers

**Tasks**:
- [ ] Implement commits query
- [ ] Add filters (user, repo, date range)
- [ ] Pagination (cursor-based)
- [ ] Unit tests

**Files to Create**:
```
src/graphql/resolvers/commit.ts
tests/unit/graphql/commit.test.ts
```

**Acceptance Criteria**:
- Commits query returns data
- All filters work correctly
- Pagination works

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
