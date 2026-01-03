# SPEC.md - Aggregator Service (cursor-analytics-core)

**Service**: cursor-analytics-core
**Type**: Backend API Service
**Language**: TypeScript (Node.js 20+)
**Framework**: Apollo Server 4 + Express
**Database**: PostgreSQL 15+
**Port**: 4000 (configurable)
**Last Updated**: January 3, 2026 (Step 10 - Integration & E2E Tests complete)

## Implementation Status

| Step | Description | Status |
|------|-------------|--------|
| 01 | Project Setup | ✅ COMPLETE |
| 02 | Database Schema & Migrations | ✅ COMPLETE |
| 03 | cursor-sim REST Client | ✅ COMPLETE |
| 04 | Ingestion Worker | NOT_STARTED |
| 05 | GraphQL Schema | ✅ COMPLETE |
| 06 | Developer Resolvers | ✅ COMPLETE |
| 07 | Commit Resolvers | ✅ COMPLETE |
| 08 | Metrics Service | ✅ COMPLETE |
| 09 | Dashboard Summary | ✅ COMPLETE |
| 10 | Integration & E2E Tests | ✅ COMPLETE |

## Overview

The Aggregator Service is the analytical engine of the Cursor Usage Analytics Platform. It ingests raw usage events from the simulator (or real Cursor API), normalizes the data into a relational structure, calculates key performance indicators, and exposes the results through a GraphQL API consumed by the frontend dashboard.

## Purpose

This service bridges the gap between raw telemetry data and actionable insights. It handles the complexity of data transformation, metric calculation, and efficient querying so that the frontend can focus purely on visualization.

## Architecture

The service follows a layered architecture with clear separation of concerns.

**API Layer**: Apollo Server handles GraphQL requests, validates queries against the schema, and delegates to resolvers.

**Resolver Layer**: Resolvers orchestrate data fetching using DataLoaders to batch and cache database queries, preventing N+1 query problems.

**Service Layer**: Business logic services implement metric calculations and complex aggregations.

**Data Access Layer**: Database client and query builders handle PostgreSQL interactions.

**Worker Layer**: Background workers handle scheduled tasks like data ingestion and materialized view refresh.

## Database Schema

### developers Table

Stores developer profile information synced from the simulator.

```sql
CREATE TABLE developers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    team VARCHAR(255) NOT NULL,
    seniority VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_developers_team ON developers(team);
CREATE INDEX idx_developers_external_id ON developers(external_id);
```

### usage_events Table

Stores individual usage events ingested from the simulator.

```sql
CREATE TABLE usage_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(255) UNIQUE NOT NULL,
    developer_id UUID NOT NULL REFERENCES developers(id),
    event_type VARCHAR(100) NOT NULL,
    event_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    lines_added INTEGER DEFAULT 0,
    lines_deleted INTEGER DEFAULT 0,
    model_used VARCHAR(100),
    accepted BOOLEAN,
    tokens_input INTEGER DEFAULT 0,
    tokens_output INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_events_developer ON usage_events(developer_id);
CREATE INDEX idx_events_timestamp ON usage_events(event_timestamp);
CREATE INDEX idx_events_type ON usage_events(event_type);
CREATE INDEX idx_events_developer_timestamp ON usage_events(developer_id, event_timestamp);
```

### daily_stats Materialized View

Pre-aggregated daily statistics for performance.

```sql
CREATE MATERIALIZED VIEW daily_stats AS
SELECT 
    developer_id,
    DATE(event_timestamp) as stat_date,
    COUNT(*) FILTER (WHERE event_type = 'cpp_suggestion_shown') as suggestions_shown,
    COUNT(*) FILTER (WHERE event_type = 'cpp_suggestion_accepted') as suggestions_accepted,
    COUNT(*) FILTER (WHERE event_type = 'chat_message') as chat_interactions,
    COUNT(*) FILTER (WHERE event_type = 'cmd_k_prompt') as cmd_k_usages,
    SUM(lines_added) as total_lines_added,
    SUM(lines_deleted) as total_lines_deleted,
    SUM(lines_added) FILTER (WHERE accepted = true) as ai_lines_added
FROM usage_events
GROUP BY developer_id, DATE(event_timestamp);

CREATE UNIQUE INDEX idx_daily_stats_pk ON daily_stats(developer_id, stat_date);
CREATE INDEX idx_daily_stats_date ON daily_stats(stat_date);
```

## GraphQL Schema

### Types

```graphql
scalar DateTime

type Developer {
    id: ID!
    externalId: String!
    name: String!
    email: String!
    team: String!
    seniority: String
    createdAt: DateTime!
    stats(range: DateRangeInput): UsageStats
    dailyStats(range: DateRangeInput): [DailyStats!]!
}

type UsageStats {
    totalSuggestions: Int!
    acceptedSuggestions: Int!
    acceptanceRate: Float
    chatInteractions: Int!
    cmdKUsages: Int!
    totalLinesAdded: Int!
    totalLinesDeleted: Int!
    aiLinesAdded: Int!
    aiVelocity: Float
}

type DailyStats {
    date: DateTime!
    suggestionsShown: Int!
    suggestionsAccepted: Int!
    acceptanceRate: Float
    chatInteractions: Int!
    cmdKUsages: Int!
    linesAdded: Int!
    linesDeleted: Int!
    aiLinesAdded: Int!
}

type TeamStats {
    teamName: String!
    memberCount: Int!
    activeMemberCount: Int!
    averageAcceptanceRate: Float
    totalSuggestions: Int!
    totalAccepted: Int!
    chatInteractions: Int!
    aiVelocity: Float
    topPerformer: Developer
}

type DashboardKPI {
    totalDevelopers: Int!
    activeDevelopers: Int!
    overallAcceptanceRate: Float
    totalSuggestionsToday: Int!
    totalAcceptedToday: Int!
    aiVelocityToday: Float
    teamComparison: [TeamStats!]!
    dailyTrend: [DailyStats!]!
}

input DateRangeInput {
    from: DateTime!
    to: DateTime!
}

enum DateRangePreset {
    TODAY
    THIS_WEEK
    THIS_MONTH
    LAST_7_DAYS
    LAST_30_DAYS
    LAST_90_DAYS
}
```

### Queries

```graphql
type Query {
    # Single developer by ID
    developer(id: ID!): Developer

    # List developers with optional filtering
    developers(
        team: String
        seniority: String
        limit: Int = 50
        offset: Int = 0
        sortBy: String = "name"
        sortOrder: String = "asc"
    ): DeveloperConnection!

    # List commits (accepted AI suggestions) with optional filtering
    commits(
        userId: ID
        team: String
        dateRange: DateRangeInput
        sortBy: String = "timestamp"
        sortOrder: String = "desc"
        limit: Int = 50
        offset: Int = 0
    ): CommitConnection!

    # Team statistics
    teamStats(teamName: String!): TeamStats
    teams: [TeamStats!]!

    # Dashboard summary - optimized for main dashboard view
    dashboardSummary(
        range: DateRangeInput
        preset: DateRangePreset
    ): DashboardKPI!

    # Health check
    health: HealthStatus!
}

type DeveloperConnection {
    nodes: [Developer!]!
    totalCount: Int!
    pageInfo: PageInfo!
}

type Commit {
    id: ID!
    externalId: String!
    timestamp: DateTime!
    linesAdded: Int!
    linesDeleted: Int!
    modelUsed: String
    tokensInput: Int!
    tokensOutput: Int!
    author: Developer!
}

type CommitConnection {
    nodes: [Commit!]!
    totalCount: Int!
    pageInfo: PageInfo!
}

type PageInfo {
    hasNextPage: Boolean!
    hasPreviousPage: Boolean!
    startCursor: String
    endCursor: String
}

type HealthStatus {
    status: String!
    database: String!
    simulator: String!
    lastIngestion: DateTime
    version: String!
}
```

## Metric Calculations

### Acceptance Rate

Formula: `(acceptedSuggestions / totalSuggestions) × 100`

Returns `null` when `totalSuggestions` is 0 to avoid division by zero.

Implementation:
```typescript
function calculateAcceptanceRate(accepted: number, shown: number): number | null {
    if (shown === 0) return null;
    return Math.round((accepted / shown) * 10000) / 100; // Round to 2 decimal places
}
```

### AI Velocity

Formula: `(aiLinesAdded / totalLinesAdded) × 100`

Measures the percentage of code that originated from AI suggestions.

Returns `null` when `totalLinesAdded` is 0.

### Team Aggregations

Team metrics use **weighted averages** based on activity, not simple averages of percentages. This prevents teams with one low-activity member from having skewed statistics.

```typescript
function calculateTeamAcceptanceRate(members: DeveloperStats[]): number | null {
    const totalShown = members.reduce((sum, m) => sum + m.suggestionsShown, 0);
    const totalAccepted = members.reduce((sum, m) => sum + m.suggestionsAccepted, 0);
    
    if (totalShown === 0) return null;
    return (totalAccepted / totalShown) * 100;
}
```

### Active Developer Definition

A developer is considered "active" if they have at least one event in the specified time range. For the dashboard summary without a range, the default is the last 7 days.

## Data Ingestion

### Polling Worker

The ingestion worker runs as a background process that periodically fetches data from the simulator.

Configuration:
```typescript
interface IngestionConfig {
    simulatorUrl: string;      // URL of the simulator API
    pollIntervalMs: number;    // Polling interval (default: 60000)
    batchSize: number;         // Events per batch (default: 1000)
    retryAttempts: number;     // Max retries on failure (default: 3)
    retryDelayMs: number;      // Initial retry delay (default: 1000)
}
```

Workflow:
1. Query the last successful ingestion timestamp from the database
2. Fetch events from simulator using `from` = last timestamp
3. Transform events to database records
4. Upsert developers (create if not exists)
5. Insert events with deduplication by `external_id`
6. Update last ingestion timestamp
7. Schedule next poll

### Deduplication

Events are deduplicated using the `external_id` field (the event ID from the simulator). The database constraint prevents duplicates, and upsert operations handle conflicts gracefully.

```sql
INSERT INTO usage_events (external_id, developer_id, event_type, ...)
VALUES ($1, $2, $3, ...)
ON CONFLICT (external_id) DO NOTHING;
```

### Retry Logic

Failed ingestion attempts use exponential backoff:

```typescript
const delays = [1000, 2000, 4000, 8000, 16000, 30000]; // ms
// Capped at 30 seconds
```

After exhausting retries, the worker logs an error and continues with the next scheduled poll.

## Error Handling

### GraphQL Errors

Errors are returned in standard GraphQL format:

```json
{
    "errors": [
        {
            "message": "Developer not found",
            "extensions": {
                "code": "NOT_FOUND",
                "id": "dev-123"
            }
        }
    ]
}
```

Error codes:
- `NOT_FOUND`: Requested resource doesn't exist
- `VALIDATION_ERROR`: Invalid input parameters
- `DATABASE_ERROR`: Database operation failed
- `SIMULATOR_ERROR`: Failed to fetch from simulator

### Input Validation

All inputs are validated before processing:
- Date ranges: `from` must be before `to`, both must be valid ISO 8601
- IDs: Must match UUID format
- Pagination: `limit` must be 1-100, `offset` must be >= 0

## Configuration

### Environment Variables

```bash
# Database
DATABASE_URL=postgresql://user:password@host:5432/database

# Simulator
SIMULATOR_URL=http://localhost:8080
POLL_INTERVAL_MS=60000

# Server
PORT=4000
NODE_ENV=development

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Performance
QUERY_COMPLEXITY_LIMIT=1000
DATALOADER_BATCH_SIZE=100
```

### Query Complexity

GraphQL queries are analyzed for complexity to prevent expensive operations:

```typescript
const complexityRules = {
    developer: 1,
    developers: (args) => args.limit * 2,
    stats: 5,
    dailyStats: (args) => daysBetween(args.range) * 2,
    dashboardSummary: 50
};
```

Queries exceeding `QUERY_COMPLEXITY_LIMIT` are rejected with an error.

## Performance Requirements

| Operation | Target Latency | Notes |
|-----------|---------------|-------|
| developer(id) | < 50ms | Single row lookup |
| developers(limit: 50) | < 100ms | Paginated list |
| teamStats | < 200ms | Aggregation query |
| dashboardSummary | < 500ms | Complex aggregations |

### Caching Strategy

**DataLoader**: Batches and caches database queries within a single request. Prevents N+1 queries when resolving nested fields.

**Materialized View**: The `daily_stats` view pre-computes daily aggregations. Refreshed every 5 minutes by a scheduled job.

**Response Caching**: Consider adding Redis caching for `dashboardSummary` with 30-second TTL in production.

## Health Checks

The `/health` endpoint (also available as GraphQL query) reports:

```json
{
    "status": "healthy",
    "database": "connected",
    "simulator": "reachable",
    "lastIngestion": "2026-01-15T10:30:00Z",
    "version": "1.0.0"
}
```

Unhealthy states:
- `database: "disconnected"` - Cannot reach PostgreSQL
- `simulator: "unreachable"` - Cannot reach simulator API
- `lastIngestion: null` - No successful ingestion yet

## Testing Requirements

### Unit Tests ✅ COMPLETE

**Coverage**: 91.49% (exceeds 80% threshold)

Implemented tests:
- ✅ All metric calculation functions (acceptance rate, AI velocity, weighted averages)
- ✅ Input validation logic
- ✅ Date range parsing and preset expansion (all 6 presets)
- ✅ Error handling for edge cases
- ✅ GraphQL resolver logic with mocked Prisma
- ✅ Context creation and request ID generation
- ✅ cursor-sim REST client with retry logic
- ✅ Database client health checks

**Test Suites**: 9 suites, 107 total tests
**Test Files**:
- `src/config/index.test.ts` - Configuration loading (3 tests)
- `src/db/__tests__/client.test.ts` - Database client (6 tests)
- `src/graphql/__tests__/context.test.ts` - GraphQL context (3 tests)
- `src/graphql/__tests__/server.test.ts` - Apollo Server setup (14 tests)
- `src/graphql/resolvers/__tests__/developer.test.ts` - Developer resolvers (13 tests)
- `src/graphql/resolvers/__tests__/commit.test.ts` - Commit resolvers (11 tests)
- `src/graphql/resolvers/__tests__/dashboard.test.ts` - Dashboard resolvers (14 tests)
- `src/services/__tests__/metrics.test.ts` - Metrics service (24 tests)
- `src/ingestion/__tests__/client.test.ts` - cursor-sim client (19 tests)

### Integration Tests ✅ COMPLETE

**Location**: `src/__tests__/integration/`

Implemented tests:
- ✅ GraphQL query execution with real Prisma and PostgreSQL
- ✅ Health check query with database status
- ✅ Developer queries with filtering (team, seniority) and pagination
- ✅ Developer stats calculation with date range filtering
- ✅ Daily stats aggregation and grouping
- ✅ Commits query with pagination and sorting
- ✅ Dashboard summary with team comparison and trends
- ✅ Team stats with weighted averages

**Test Setup**: `src/__tests__/integration/setup.ts`
- `createTestDb()` - Fresh Prisma client for testing
- `createTestSimClient()` - Mocked cursor-sim client
- `seedTestData()` - Realistic test data (3 developers, 7 days of events)
- `cleanupDb()` - Teardown after tests

**Test Suite**: 13 integration tests covering full GraphQL API

### E2E Tests ✅ COMPLETE

**Location**: `src/__tests__/e2e/full-pipeline.test.ts`

Tests the complete data flow:
- ✅ Seed database → Execute GraphQL queries → Verify aggregations
- ✅ Complex multi-query pipeline (dashboard with team breakdown)
- ✅ Developer profiles with nested stats and daily breakdown
- ✅ Commit aggregation across teams with sorting
- ✅ Team comparison with weighted metrics
- ✅ Filtered commits by user and date range
- ✅ Data consistency validation (referential integrity)
- ✅ Weighted team average calculations

**Dataset**: 5 developers, 2 teams, 7 days of realistic usage data

**Test Suite**: 6 E2E tests covering complex scenarios

### Performance Tests ✅ COMPLETE

**Location**: `src/__tests__/performance/large-dataset.test.ts`

Tests with production-scale data:
- ✅ 10,500+ events across 50 developers over 30 days
- ✅ Dashboard summary query < 2000ms
- ✅ List 100 developers with stats < 1000ms
- ✅ Paginate through 1000 commits < 500ms per page
- ✅ Team stats aggregation < 800ms
- ✅ Date range filtering < 300ms
- ✅ 10 concurrent queries < 3000ms total
- ✅ Pagination performance consistency (no degradation)

**Performance Targets**:

| Operation | Target | Actual (Verified) |
|-----------|--------|-------------------|
| Dashboard summary (50 devs) | < 2000ms | ✅ Tested |
| Developers list (100) | < 1000ms | ✅ Tested |
| Commits pagination (100) | < 500ms | ✅ Tested |
| Team aggregation (5 teams) | < 800ms | ✅ Tested |
| Date range filter | < 300ms | ✅ Tested |

**Test Suite**: 7 performance tests with timing assertions

## Dependencies

Production dependencies:
- `@apollo/server` - GraphQL server
- `express` - HTTP framework
- `pg` - PostgreSQL client (or `prisma`)
- `dataloader` - Request-scoped batching
- `graphql` - GraphQL execution engine
- `date-fns` - Date manipulation

Development dependencies:
- `typescript` - Type checking
- `jest` + `ts-jest` - Testing
- `supertest` - HTTP testing
- `pg-mem` - In-memory PostgreSQL for tests
