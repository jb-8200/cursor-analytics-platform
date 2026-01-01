# Implementation Tasks: Cursor Usage Analytics Platform

**Version**: 1.0.0  
**Last Updated**: January 2026  

This document breaks down features into actionable implementation tasks. Each task is designed to be completable in 2-4 hours and includes a clear definition of done. Tasks are organized by phase and service, with dependencies explicitly noted. The task structure supports both human developers and AI coding assistants following spec-driven development.

## Task Status Legend

Each task uses one of the following status indicators to track progress across the project.

**NOT_STARTED** indicates work has not begun on this task. **IN_PROGRESS** indicates active development is underway. **BLOCKED** indicates the task cannot proceed due to an unresolved dependency. **REVIEW** indicates implementation is complete and awaiting code review. **DONE** indicates the task is complete and merged.

## Phase 1: Core Functionality (MVP)

Phase 1 delivers a working end-to-end system. The goal is to have all three services communicating and basic visualizations functional by the end of this phase.

### Service A: cursor-sim Tasks

#### TASK-SIM-001: Initialize Go Project Structure

**Status**: NOT_STARTED  
**Feature**: SIM-004  
**Estimated Hours**: 2  
**Assignee**: Unassigned  

This task creates the foundational Go project with proper module structure, linting configuration, and initial dependencies.

**Implementation Steps:**

The developer should start by creating the project directory and initializing the Go module with `go mod init github.com/org/cursor-sim`. The directory structure should follow Go conventions with a `cmd/` directory for the main entry point, an `internal/` directory for private packages, and a `pkg/` directory for any packages intended for external use.

The essential directories to create include `cmd/cursor-sim/` for the main application, `internal/config/` for configuration handling, `internal/models/` for domain types, `internal/generator/` for data generation logic, `internal/api/` for HTTP handlers, and `internal/storage/` for the in-memory data store.

Configuration files needed include a `.golangci.yml` for linting rules, a `Makefile` with common commands, and a `Dockerfile` for containerization.

**Definition of Done:**

The project compiles successfully with `go build ./...`. The linter passes with `golangci-lint run`. A minimal main.go file exists that prints a version string and exits.

**Dependencies**: None

**Files to Create:**
- cmd/cursor-sim/main.go
- internal/config/config.go
- go.mod
- .golangci.yml
- Makefile
- Dockerfile

---

#### TASK-SIM-002: Implement CLI Flag Parsing

**Status**: NOT_STARTED  
**Feature**: SIM-004  
**Estimated Hours**: 2  
**Assignee**: Unassigned  

This task implements command-line argument parsing for all simulator configuration options.

**Implementation Steps:**

Create a `Config` struct in `internal/config/config.go` that holds all configuration values including Port (int), Developers (int), Velocity (string), Fluctuation (float64), and Seed (int64 for reproducibility).

Implement a `ParseFlags()` function that uses the standard `flag` package to define and parse all flags. The function should return a `Config` struct and an error for invalid configurations.

Add validation logic that ensures Port is between 1024 and 65535, Developers is between 1 and 1000, Velocity is one of "low", "medium", or "high", and Fluctuation is between 0.0 and 1.0.

Implement `--help` output that clearly documents each flag, its default value, and valid options.

**Definition of Done:**

Running `./cursor-sim --help` displays all options with descriptions. Invalid flag values produce descriptive error messages. Valid configurations are correctly parsed into the Config struct. Unit tests cover all validation scenarios.

**Dependencies**: TASK-SIM-001

**Files to Create/Modify:**
- internal/config/config.go
- internal/config/config_test.go
- cmd/cursor-sim/main.go

**Test Cases to Write:**
```go
func TestParseFlags_Defaults(t *testing.T)
func TestParseFlags_CustomValues(t *testing.T)
func TestParseFlags_InvalidPort(t *testing.T)
func TestParseFlags_InvalidVelocity(t *testing.T)
func TestParseFlags_InvalidFluctuation(t *testing.T)
```

---

#### TASK-SIM-003: Implement Developer Profile Generator

**Status**: NOT_STARTED  
**Feature**: SIM-001  
**Estimated Hours**: 4  
**Assignee**: Unassigned  

This task implements the logic for generating realistic developer profiles with varied characteristics.

**Implementation Steps:**

Create a `Developer` struct in `internal/models/developer.go` with fields for ID, Name, Email, Team, Seniority, AcceptanceRate, and CreatedAt.

Implement a name generator that produces realistic names using a deterministic algorithm based on the random seed. The generator should use common first and last name lists to create combinations.

Create a team assignment algorithm that distributes developers across a configurable number of teams (default 5) with some variance in team sizes.

Implement seniority assignment following the distribution of 20% junior, 50% mid-level, and 30% senior, with acceptance rates of 60%, 75%, and 90% respectively (with some per-developer variance).

Create a `GenerateDevelopers(count int, seed int64) []Developer` function that produces the specified number of profiles.

**Definition of Done:**

Generated profiles have unique IDs and emails. Seniority distribution matches target percentages within 10% variance. Acceptance rates correlate correctly with seniority. The same seed produces identical profiles across runs. All unit tests pass.

**Dependencies**: TASK-SIM-001

**Files to Create:**
- internal/models/developer.go
- internal/generator/developer_generator.go
- internal/generator/developer_generator_test.go
- internal/generator/names.go (name data)

**Test Cases to Write:**
```go
func TestGenerateDevelopers_Count(t *testing.T)
func TestGenerateDevelopers_UniqueIDs(t *testing.T)
func TestGenerateDevelopers_SeniorityDistribution(t *testing.T)
func TestGenerateDevelopers_AcceptanceRates(t *testing.T)
func TestGenerateDevelopers_Reproducibility(t *testing.T)
```

---

#### TASK-SIM-004: Implement Event Generation Engine

**Status**: NOT_STARTED  
**Feature**: SIM-002  
**Estimated Hours**: 6  
**Assignee**: Unassigned  

This task implements the core simulation logic for generating usage events with realistic timing patterns.

**Implementation Steps:**

Create a `UsageEvent` struct in `internal/models/event.go` with fields for ID, DeveloperID, EventType, Timestamp, and Metadata (containing LinesAdded, LinesDeleted, ModelUsed, Accepted, TokensInput, TokensOutput).

Implement a Poisson distribution event timer that generates events with natural clustering rather than uniform spacing. The mean interval should be derived from the velocity configuration.

Create event type generators for each of the four types: cpp_suggestion_shown, cpp_suggestion_accepted, chat_message, and cmd_k_prompt. The accepted events should only follow shown events with probability matching the developer's acceptance rate.

Implement per-developer fluctuation by adjusting each developer's base rate by a random factor within the configured fluctuation range.

Create an `EventGenerator` that manages goroutines for each developer, coordinating event generation and providing access to generated events.

**Definition of Done:**

Events are generated continuously with Poisson-distributed timing. Event types appear in correct proportions. Acceptance rates match developer profiles within statistical variance. The generator can be started and stopped cleanly. All unit tests pass including statistical distribution tests.

**Dependencies**: TASK-SIM-003

**Files to Create:**
- internal/models/event.go
- internal/generator/event_generator.go
- internal/generator/event_generator_test.go
- internal/generator/poisson.go

**Test Cases to Write:**
```go
func TestEventGenerator_AllEventTypes(t *testing.T)
func TestEventGenerator_AcceptanceCorrelation(t *testing.T)
func TestEventGenerator_VelocityImpact(t *testing.T)
func TestEventGenerator_PoissonDistribution(t *testing.T)
func TestEventGenerator_StartStop(t *testing.T)
```

---

#### TASK-SIM-005: Implement In-Memory Storage

**Status**: NOT_STARTED  
**Feature**: SIM-003  
**Estimated Hours**: 3  
**Assignee**: Unassigned  

This task implements thread-safe in-memory storage for developers and events.

**Implementation Steps:**

Create a `Store` interface in `internal/storage/store.go` that defines methods for storing and retrieving developers and events.

Implement `MemoryStore` that uses `sync.Map` or mutex-protected maps for thread-safe access. The store should support adding developers, listing all developers, getting a developer by ID, adding events, and querying events by time range.

Implement efficient time-range queries by maintaining events in a time-sorted structure. Consider using a skip list or B-tree for efficient range queries, or accept O(n) scanning for simplicity given the in-memory nature.

Add a method to clear all data for testing purposes.

**Definition of Done:**

Concurrent access to the store is safe. Time-range queries return correct results. The store can handle 1000 developers and 100,000 events without performance degradation. All unit tests pass including concurrency tests.

**Dependencies**: TASK-SIM-003, TASK-SIM-004

**Files to Create:**
- internal/storage/store.go
- internal/storage/memory_store.go
- internal/storage/memory_store_test.go

**Test Cases to Write:**
```go
func TestMemoryStore_AddDeveloper(t *testing.T)
func TestMemoryStore_GetDeveloper(t *testing.T)
func TestMemoryStore_AddEvent(t *testing.T)
func TestMemoryStore_QueryEventsByTimeRange(t *testing.T)
func TestMemoryStore_ConcurrentAccess(t *testing.T)
```

---

#### TASK-SIM-006: Implement REST API Handlers

**Status**: NOT_STARTED  
**Feature**: SIM-003  
**Estimated Hours**: 4  
**Assignee**: Unassigned  

This task implements the REST API endpoints that expose simulated data.

**Implementation Steps:**

Create HTTP handlers in `internal/api/handlers.go` using the standard `net/http` package. Each handler should interact with the storage layer to retrieve data.

Implement `GET /v1/org/users` to return all developers as a JSON array. Implement `GET /v1/org/users/:id` to return a single developer or 404 if not found. Implement `GET /v1/stats/activity` with `from` and `to` query parameters for time filtering. Implement `GET /health` to return service status.

Add middleware for CORS headers to allow cross-origin requests during development. Add middleware for request logging.

Implement pagination for the activity endpoint, returning a maximum of 1000 events per request with a `nextCursor` field when more results exist.

**Definition of Done:**

All endpoints return correct JSON responses matching the documented schema. Time filtering correctly bounds results. CORS headers are present on responses. Invalid requests return appropriate error codes. All integration tests pass.

**Dependencies**: TASK-SIM-005

**Files to Create:**
- internal/api/handlers.go
- internal/api/handlers_test.go
- internal/api/middleware.go
- internal/api/router.go

**Test Cases to Write:**
```go
func TestHandler_ListDevelopers(t *testing.T)
func TestHandler_GetDeveloper(t *testing.T)
func TestHandler_GetDeveloper_NotFound(t *testing.T)
func TestHandler_QueryActivity(t *testing.T)
func TestHandler_QueryActivity_TimeFilter(t *testing.T)
func TestHandler_QueryActivity_Pagination(t *testing.T)
func TestHandler_Health(t *testing.T)
```

---

#### TASK-SIM-007: Wire Up Main Application

**Status**: NOT_STARTED  
**Feature**: SIM-004  
**Estimated Hours**: 2  
**Assignee**: Unassigned  

This task connects all components in the main function to create the working simulator.

**Implementation Steps:**

Update `cmd/cursor-sim/main.go` to parse configuration, initialize the storage, generate developers, start the event generator, set up the HTTP router, and start the server.

Implement graceful shutdown handling using signals (SIGINT, SIGTERM) to cleanly stop event generation and close the HTTP server.

Add startup logging that reports the configuration being used.

Create a Makefile target for building and running the simulator with common configurations.

**Definition of Done:**

Running `./cursor-sim` starts the server with default configuration. The server responds to health checks immediately. Events begin generating after startup. SIGINT causes graceful shutdown. Integration tests verify end-to-end functionality.

**Dependencies**: TASK-SIM-002, TASK-SIM-004, TASK-SIM-006

**Files to Modify:**
- cmd/cursor-sim/main.go
- Makefile

---

### Service B: cursor-analytics-core Tasks

#### TASK-CORE-001: Initialize TypeScript Project Structure

**Status**: NOT_STARTED  
**Feature**: CORE-002  
**Estimated Hours**: 2  
**Assignee**: Unassigned  

This task creates the Node.js/TypeScript project with Apollo Server and database tooling.

**Implementation Steps:**

Initialize the project with `npm init` and configure TypeScript with strict mode enabled. Install core dependencies including Apollo Server 4, Express, PostgreSQL client (pg or Prisma), and GraphQL tools.

Create the directory structure with `src/` containing subdirectories for `resolvers/`, `models/`, `services/`, `workers/`, and `db/`.

Configure ESLint with TypeScript support and Prettier for code formatting. Set up Jest with ts-jest for testing.

Create a Dockerfile for containerization and configure the build process.

**Definition of Done:**

The project compiles with `npm run build`. Linting passes with `npm run lint`. Jest is configured and a sample test runs. The Docker image builds successfully.

**Dependencies**: None

**Files to Create:**
- package.json
- tsconfig.json
- .eslintrc.js
- .prettierrc
- jest.config.js
- Dockerfile
- src/index.ts (minimal)

---

#### TASK-CORE-002: Define Database Schema and Migrations

**Status**: NOT_STARTED  
**Feature**: CORE-002  
**Estimated Hours**: 3  
**Assignee**: Unassigned  

This task creates the PostgreSQL schema with migration support.

**Implementation Steps:**

Choose a migration tool (Prisma recommended for TypeScript projects, or knex/node-pg-migrate for raw SQL control). Create the initial migration that defines the developers table with columns for id (UUID primary key), external_id (unique, from simulator), name, email, team, seniority, and timestamps.

Create the usage_events table with columns for id (UUID primary key), developer_id (foreign key), event_type, event_timestamp, lines_added, lines_deleted, model_used, accepted (boolean), tokens_input, tokens_output, and created_at.

Add indexes on developer_id, event_timestamp, and event_type for query performance.

Create a materialized view for daily_stats that pre-aggregates daily statistics per developer.

**Definition of Done:**

Migrations run successfully against a fresh PostgreSQL instance. Schema matches the design document. Indexes are created for performance-critical columns. A rollback migration exists and successfully reverses the changes.

**Dependencies**: TASK-CORE-001

**Files to Create:**
- prisma/schema.prisma (or migrations/*.sql)
- src/db/client.ts
- src/db/migrations/*

---

#### TASK-CORE-003: Implement GraphQL Schema

**Status**: NOT_STARTED  
**Feature**: CORE-003  
**Estimated Hours**: 3  
**Assignee**: Unassigned  

This task defines the GraphQL schema using SDL and sets up Apollo Server.

**Implementation Steps:**

Create the schema definition in `src/schema.graphql` with all types defined in the design document. Define the Developer type with fields for id, externalId, name, email, team, seniority, stats (with range argument), and dailyStats.

Define the UsageStats type with fields for totalSuggestions, acceptedSuggestions, acceptanceRate, chatInteractions, cmdKUsages, totalLinesAdded, totalLinesDeleted, and aiVelocity.

Define the Query type with developer, developers, teamStats, teams, and dashboardSummary queries.

Create input types for DateRangeInput with from and to DateTime fields.

Set up Apollo Server with the schema and placeholder resolvers that return mock data.

**Definition of Done:**

GraphQL Playground is accessible and shows the schema. All types and queries are defined matching the specification. Mock resolvers return valid response shapes. Type generation produces TypeScript types from the schema.

**Dependencies**: TASK-CORE-001

**Files to Create:**
- src/schema.graphql
- src/resolvers/index.ts (placeholder)
- src/server.ts
- src/generated/types.ts (from codegen)

---

#### TASK-CORE-004: Implement Data Ingestion Worker

**Status**: NOT_STARTED  
**Feature**: CORE-001  
**Estimated Hours**: 5  
**Assignee**: Unassigned  

This task implements the background worker that polls the simulator and stores data.

**Implementation Steps:**

Create a `DataIngestionWorker` class in `src/workers/ingestion.ts` that manages the polling lifecycle. The worker should track the last successful fetch timestamp and use it for incremental queries.

Implement HTTP client logic to fetch from the simulator's `/v1/org/users` and `/v1/stats/activity` endpoints. Handle the JSON response and transform it into database records.

Add retry logic with exponential backoff for failed requests. Start with a 1-second delay, double after each failure, and cap at 30 seconds.

Implement deduplication by checking event IDs before insertion. Consider using upsert operations for efficiency.

Create a worker manager that starts the ingestion on application startup and stops it gracefully on shutdown.

**Definition of Done:**

The worker polls at the configured interval. New events are stored in the database. Duplicate events are not created. The worker recovers gracefully from simulator outages. Unit tests cover all scenarios including failure cases.

**Dependencies**: TASK-CORE-002

**Files to Create:**
- src/workers/ingestion.ts
- src/workers/ingestion.test.ts
- src/services/simulatorClient.ts
- src/services/simulatorClient.test.ts

**Test Cases to Write:**
```typescript
describe('DataIngestionWorker', () => {
  it('should poll at configured interval')
  it('should use last timestamp for incremental fetch')
  it('should retry with exponential backoff')
  it('should deduplicate events')
  it('should stop gracefully on shutdown')
})
```

---

#### TASK-CORE-005: Implement Metric Calculation Service

**Status**: NOT_STARTED  
**Feature**: CORE-004  
**Estimated Hours**: 4  
**Assignee**: Unassigned  

This task implements the business logic for calculating KPIs from raw event data.

**Implementation Steps:**

Create a `MetricsService` class in `src/services/metrics.ts` with methods for calculating each KPI.

Implement `calculateAcceptanceRate(developerId, dateRange)` that queries suggestion events and computes the percentage. Return null when no suggestions exist to avoid division by zero.

Implement `calculateAIVelocity(developerId, dateRange)` that computes the ratio of AI lines to total lines.

Implement aggregate methods for team and organization-level calculations. Team calculations should use weighted averages based on total activity, not simple averages of percentages.

Implement `refreshDailyStats()` to update the materialized view for pre-computed statistics.

**Definition of Done:**

All metric calculations match the formulas in the design document. Edge cases (zero denominators) are handled correctly. Team aggregations use weighted averages. Unit tests verify calculations with known test data.

**Dependencies**: TASK-CORE-002

**Files to Create:**
- src/services/metrics.ts
- src/services/metrics.test.ts

**Test Cases to Write:**
```typescript
describe('MetricsService', () => {
  it('should calculate acceptance rate correctly')
  it('should return null for zero suggestions')
  it('should calculate AI velocity correctly')
  it('should use weighted average for team stats')
  it('should respect date range boundaries')
})
```

---

#### TASK-CORE-006: Implement Developer Resolvers

**Status**: NOT_STARTED  
**Feature**: CORE-005  
**Estimated Hours**: 4  
**Assignee**: Unassigned  

This task implements GraphQL resolvers for developer queries.

**Implementation Steps:**

Create resolver functions in `src/resolvers/developer.ts` for the developer and developers queries.

Implement the `developer(id)` resolver to fetch a single developer by ID, returning null if not found.

Implement the `developers(team, limit, offset)` resolver with filtering and pagination support. Use DataLoader to batch database queries for efficiency.

Implement the nested `stats` resolver on the Developer type that calls the metrics service to calculate statistics. Accept the optional range argument for time-bounded calculations.

Implement the `dailyStats` resolver to return the pre-aggregated daily statistics from the materialized view.

**Definition of Done:**

All developer queries return correct data matching the schema. Pagination works correctly with limit and offset. Team filtering returns only matching developers. Nested stats resolve correctly with optional time range. DataLoader prevents N+1 query problems.

**Dependencies**: TASK-CORE-003, TASK-CORE-005

**Files to Create:**
- src/resolvers/developer.ts
- src/resolvers/developer.test.ts
- src/dataloaders/developerLoader.ts

**Test Cases to Write:**
```typescript
describe('Developer Resolvers', () => {
  it('should return developer by ID')
  it('should return null for non-existent developer')
  it('should paginate developers list')
  it('should filter by team')
  it('should calculate nested stats')
})
```

---

### Service C: cursor-viz-spa Tasks

#### TASK-VIZ-001: Initialize React Project with Vite

**Status**: NOT_STARTED  
**Feature**: VIZ-001  
**Estimated Hours**: 2  
**Assignee**: Unassigned  

This task creates the React application with modern tooling.

**Implementation Steps:**

Create the project using `npm create vite@latest cursor-viz-spa -- --template react-ts`. Install dependencies including Apollo Client, TanStack Query, Recharts, and Tailwind CSS.

Configure TypeScript with strict mode. Set up ESLint with React and TypeScript rules. Configure Prettier for consistent formatting.

Create the directory structure with `src/components/`, `src/hooks/`, `src/graphql/`, `src/pages/`, and `src/utils/`.

Set up Vitest for unit testing and create a sample test to verify the configuration. Configure MSW (Mock Service Worker) for API mocking in tests.

Create a Dockerfile for containerization.

**Definition of Done:**

The application runs with `npm run dev`. Linting passes with `npm run lint`. A sample component test runs successfully. The Docker image builds and serves the application.

**Dependencies**: None

**Files to Create:**
- package.json
- vite.config.ts
- tsconfig.json
- tailwind.config.js
- .eslintrc.cjs
- Dockerfile
- src/main.tsx
- src/App.tsx

---

#### TASK-VIZ-002: Implement Dashboard Layout

**Status**: NOT_STARTED  
**Feature**: VIZ-001  
**Estimated Hours**: 4  
**Assignee**: Unassigned  

This task creates the responsive dashboard layout structure.

**Implementation Steps:**

Create a `DashboardGrid` component using CSS Grid that defines areas for the header, main content, and optional sidebar.

Create a `Header` component that displays the KPI summary bar with placeholders for total developers, active developers, and acceptance rate.

Implement responsive breakpoints at 768px and 1024px using Tailwind's responsive utilities. On mobile, the layout should stack vertically. On tablet and desktop, content should use the grid layout.

Create placeholder components for `VelocityHeatmap`, `DeveloperTable`, and `TeamRadarChart` that will be implemented in subsequent tasks.

**Definition of Done:**

The dashboard renders with proper layout at all viewport sizes. Components reflow appropriately at breakpoints. No horizontal scrolling occurs on mobile. Placeholder components are visible in their grid positions.

**Dependencies**: TASK-VIZ-001

**Files to Create:**
- src/components/layout/DashboardGrid.tsx
- src/components/layout/Header.tsx
- src/components/layout/Sidebar.tsx
- src/pages/Dashboard.tsx

**Test Cases to Write:**
```tsx
describe('DashboardGrid', () => {
  it('should render all layout sections')
  it('should apply responsive classes')
})
```

---

#### TASK-VIZ-003: Implement GraphQL Client Setup

**Status**: NOT_STARTED  
**Feature**: VIZ-001  
**Estimated Hours**: 2  
**Assignee**: Unassigned  

This task configures Apollo Client for GraphQL communication.

**Implementation Steps:**

Install and configure Apollo Client with the GraphQL endpoint from environment variables. Create the client instance in `src/graphql/client.ts`.

Define GraphQL queries in `src/graphql/queries.ts` for dashboard summary, developer list, and team statistics.

Create query fragments in `src/graphql/fragments.ts` for reusable type selections.

Set up GraphQL code generation to produce TypeScript types from the schema. Configure the codegen to run as part of the build process.

**Definition of Done:**

Apollo Client connects to the backend GraphQL endpoint. Queries are defined with proper TypeScript types. Code generation produces accurate types from the schema. The client handles errors gracefully.

**Dependencies**: TASK-VIZ-001

**Files to Create:**
- src/graphql/client.ts
- src/graphql/queries.ts
- src/graphql/fragments.ts
- codegen.ts

---

#### TASK-VIZ-004: Implement Velocity Heatmap

**Status**: NOT_STARTED  
**Feature**: VIZ-002  
**Estimated Hours**: 6  
**Assignee**: Unassigned  

This task implements the GitHub-style contribution heatmap.

**Implementation Steps:**

Create a `VelocityHeatmap` component that accepts daily statistics data as props. The component should render a grid of cells representing days, with 7 rows for days of the week and columns for weeks.

Implement color mapping logic that converts acceptance counts to color intensity. Use a gradient from light (low activity) to dark (high activity). Make the color scale configurable.

Add day-of-week labels on the left edge (at minimum Mon, Wed, Fri). Add month labels above the appropriate week boundaries.

Implement hover tooltips using a lightweight tooltip library or custom CSS. The tooltip should display the date and exact acceptance count.

Use Recharts or custom SVG rendering for the grid visualization.

**Definition of Done:**

The heatmap renders with correct date alignment. Colors accurately represent value intensity. Tooltips display on hover. Labels are positioned correctly. The component is responsive and resizes appropriately.

**Dependencies**: TASK-VIZ-002, TASK-VIZ-003

**Files to Create:**
- src/components/charts/VelocityHeatmap.tsx
- src/components/charts/VelocityHeatmap.test.tsx
- src/components/charts/Tooltip.tsx

**Test Cases to Write:**
```tsx
describe('VelocityHeatmap', () => {
  it('should render correct number of cells')
  it('should apply color intensity based on value')
  it('should show day labels')
  it('should show month labels')
})
```

---

#### TASK-VIZ-005: Implement Developer Efficiency Table

**Status**: NOT_STARTED  
**Feature**: VIZ-004  
**Estimated Hours**: 4  
**Assignee**: Unassigned  

This task implements the sortable developer metrics table.

**Implementation Steps:**

Create a `DeveloperTable` component that displays columns for Name, Team, Total Suggestions, Accepted, Acceptance Rate, and AI Lines.

Implement client-side sorting by adding click handlers to column headers. Track the current sort column and direction in component state. Sort the data array before rendering.

Add visual indicators for sort direction (up/down arrows) on the active column header.

Implement conditional row styling that applies a warning background to rows where acceptance rate is below 20%.

Add a search input that filters the displayed rows by developer name. Implement debouncing to avoid excessive re-renders during typing.

Implement pagination with configurable page size (default 25). Add navigation controls for first, previous, next, and last pages.

**Definition of Done:**

The table displays all required columns. Sorting works correctly for all columns. Low acceptance rate rows are visually distinct. Search filtering matches partial names. Pagination displays the correct subset of rows.

**Dependencies**: TASK-VIZ-002, TASK-VIZ-003

**Files to Create:**
- src/components/tables/DeveloperTable.tsx
- src/components/tables/DeveloperTable.test.tsx
- src/components/tables/TablePagination.tsx

**Test Cases to Write:**
```tsx
describe('DeveloperTable', () => {
  it('should display all columns')
  it('should sort by column on click')
  it('should highlight low acceptance rates')
  it('should filter by search term')
  it('should paginate correctly')
})
```

---

### Infrastructure Tasks

#### TASK-INFRA-001: Create Docker Compose Configuration

**Status**: NOT_STARTED  
**Feature**: All  
**Estimated Hours**: 3  
**Assignee**: Unassigned  

This task creates the Docker Compose file that orchestrates all services.

**Implementation Steps:**

Create `docker-compose.yml` in the project root that defines services for cursor-sim, postgres, cursor-analytics-core, and cursor-viz-spa.

Configure health checks for each service. The simulator should check `/health`, the core should verify database connectivity and simulator reachability, and the SPA should check that the dev server is responding.

Set up service dependencies so that postgres starts first, then cursor-sim, then cursor-analytics-core (after both are healthy), and finally cursor-viz-spa.

Configure environment variables for each service using the values defined in the design document.

Create a `docker-compose.override.yml` for development-specific settings like volume mounts for hot reloading.

**Definition of Done:**

Running `docker-compose up` starts all services in the correct order. Health checks pass for all services. The dashboard is accessible at localhost:3000 and displays data from the simulator.

**Dependencies**: TASK-SIM-007, TASK-CORE-006, TASK-VIZ-005

**Files to Create:**
- docker-compose.yml
- docker-compose.override.yml
- .env.example

---

## Phase 2: Enhanced Analytics

Phase 2 adds team-level analytics, improved filtering, and polished user experience.

### TASK-CORE-007: Implement Team Statistics Resolvers

**Status**: NOT_STARTED  
**Feature**: CORE-006  
**Estimated Hours**: 4  
**Assignee**: Unassigned  

Implements GraphQL resolvers for team-level aggregations including member count, average acceptance rate, and top performer identification.

**Dependencies**: TASK-CORE-006

---

### TASK-CORE-008: Implement Dashboard Summary Query

**Status**: NOT_STARTED  
**Feature**: CORE-007  
**Estimated Hours**: 5  
**Assignee**: Unassigned  

Implements the comprehensive dashboard query that returns all KPIs efficiently in a single request.

**Dependencies**: TASK-CORE-007

---

### TASK-VIZ-006: Implement Team Radar Chart

**Status**: NOT_STARTED  
**Feature**: VIZ-003  
**Estimated Hours**: 6  
**Assignee**: Unassigned  

Implements the multi-axis radar chart for team comparison using Recharts.

**Dependencies**: TASK-CORE-007

---

### TASK-VIZ-007: Implement Date Range Picker

**Status**: NOT_STARTED  
**Feature**: VIZ-005  
**Estimated Hours**: 3  
**Assignee**: Unassigned  

Implements the date range selection component with presets and custom range support.

**Dependencies**: TASK-VIZ-003

---

### TASK-VIZ-008: Implement Loading and Error States

**Status**: NOT_STARTED  
**Feature**: VIZ-006, VIZ-007  
**Estimated Hours**: 4  
**Assignee**: Unassigned  

Implements skeleton loaders, error displays, and retry mechanisms across all components.

**Dependencies**: TASK-VIZ-004, TASK-VIZ-005

---

## Task Dependency Graph

The following shows the critical path through the implementation tasks.

```
TASK-SIM-001 (Project Init)
    │
    ├── TASK-SIM-002 (CLI) ────────────┐
    │                                   │
    └── TASK-SIM-003 (Developers)       │
            │                           │
            └── TASK-SIM-004 (Events)   │
                    │                   │
                    └── TASK-SIM-005 (Storage)
                            │
                            └── TASK-SIM-006 (API)
                                    │
                                    └── TASK-SIM-007 (Wire Up)
                                            │
                                            ▼
TASK-CORE-001 (Project Init)
    │
    ├── TASK-CORE-002 (Schema)
    │       │
    │       ├── TASK-CORE-004 (Ingestion)
    │       │
    │       └── TASK-CORE-005 (Metrics)
    │               │
    │               └── TASK-CORE-006 (Resolvers)
    │                       │
    └── TASK-CORE-003 (GraphQL)
                            │
                            ▼
TASK-VIZ-001 (Project Init)
    │
    ├── TASK-VIZ-002 (Layout)
    │       │
    │       ├── TASK-VIZ-004 (Heatmap)
    │       │
    │       └── TASK-VIZ-005 (Table)
    │
    └── TASK-VIZ-003 (Client)
                            │
                            ▼
                    TASK-INFRA-001 (Docker)
```
