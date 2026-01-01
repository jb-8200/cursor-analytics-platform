# Feature F002: Aggregator Ingestion and Analytics Engine

**Feature ID:** F002  
**Service:** cursor-analytics-core  
**Priority:** P0 (Critical Path)  
**Status:** Specification Complete

---

## 1. Overview

The Aggregator Service acts as the data processing layer between the simulator and the dashboard. It polls the simulator API, normalizes event data into a relational schema, computes key performance indicators, and exposes a GraphQL API for flexible querying by the frontend.

### 1.1 Business Value

Raw event data from the Cursor API is voluminous and requires aggregation to derive meaningful insights. The aggregator transforms thousands of individual events into actionable metrics like acceptance rates, AI velocity scores, and team comparisons that directly inform engineering decisions about AI tool adoption.

### 1.2 Success Criteria

The feature is complete when the aggregator reliably syncs data from the simulator every 60 seconds without data loss, when all KPIs match manual calculations within 0.1% precision, and when GraphQL queries return correct aggregated data within 200ms for typical dashboard requests.

---

## 2. Functional Requirements

### 2.1 Data Ingestion Worker (FR-AGG-001)

The aggregator must run a background job that periodically fetches data from the simulator and persists it to PostgreSQL. This worker should be resilient to transient failures and maintain consistency.

**Polling Behavior:**
- Default interval: 60 seconds (configurable via `POLL_INTERVAL_MS`)
- On failure: Exponential backoff starting at 5 seconds, max 5 minutes
- Deduplication: Events are identified by composite key (developer_id, date) for daily stats
- Incremental sync: Only fetch data since last successful sync

**Ingestion Flow:**

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│                 │     │                  │     │                 │
│  Timer Trigger  │────▶│  Fetch Simulator │────▶│  Validate Data  │
│  (every 60s)    │     │  GET /v1/stats   │     │  Schema Check   │
│                 │     │                  │     │                 │
└─────────────────┘     └──────────────────┘     └────────┬────────┘
                                                          │
                                                          ▼
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│                 │     │                  │     │                 │
│  Emit Updated   │◀────│  Calculate KPIs  │◀────│  Upsert to DB   │
│  Event          │     │  Team Stats      │     │  (Transaction)  │
│                 │     │                  │     │                 │
└─────────────────┘     └──────────────────┘     └─────────────────┘
```

**Acceptance Criteria:**
- AC1: Worker starts automatically on service startup
- AC2: Consecutive failures do not crash the service
- AC3: Duplicate data does not create duplicate database records
- AC4: Full sync of 100 developers × 90 days completes in under 30 seconds
- AC5: Worker logs progress and errors with structured logging

### 2.2 Data Normalization (FR-AGG-002)

Raw data from the simulator must be transformed into normalized relational tables that support efficient querying and accurate aggregation.

**Source to Target Mapping:**

| Simulator Field | Database Table | Database Column |
|----------------|----------------|-----------------|
| `email` | developers | email |
| `date` | daily_stats | date |
| `totalTabsShown` | daily_stats | total_tabs_shown |
| `totalTabsAccepted` | daily_stats | total_tabs_accepted |
| `totalLinesAdded` | daily_stats | lines_added |
| `totalLinesDeleted` | daily_stats | lines_deleted |
| `acceptedLinesAdded` | daily_stats | accepted_lines_added |
| `composerRequests` | daily_stats | composer_requests |
| `chatRequests` | daily_stats | chat_requests |
| `agentRequests` | daily_stats | agent_requests |
| `cmdkUsages` | daily_stats | cmdk_usages |
| `mostUsedModel` | daily_stats | most_used_model |

**Derived Fields:**
- `acceptance_rate`: Calculated as `total_tabs_accepted / total_tabs_shown`
- `is_active`: Boolean, true if any non-zero activity metric

**Developer Upsert Logic:**
- If developer email exists, update name/team/role if changed
- If developer email is new, create new developer record
- Maintain referential integrity with daily_stats

**Acceptance Criteria:**
- AC1: All fields map correctly from simulator to database
- AC2: Null/missing fields in source default to 0 or appropriate value
- AC3: Acceptance rate calculation handles division by zero
- AC4: Developer records are created before daily_stats to maintain FK integrity
- AC5: Unicode characters in names are preserved correctly

### 2.3 KPI Calculation Engine (FR-AGG-003)

The aggregator must compute several key performance indicators from the normalized data. These calculations must be consistent and match defined formulas exactly.

**Individual Developer KPIs:**

| KPI | Formula | Description |
|-----|---------|-------------|
| Acceptance Rate | `(accepted / shown) × 100` | Percentage of AI suggestions accepted |
| AI Velocity | `(acceptedLines / totalLines) × 100` | AI contribution to code output |
| Chat Dependency | `chatRequests / (completionEvents + 1)` | Reliance on chat vs. inline |
| Productivity Score | `(acceptedLines × acceptanceRate) / 100` | Combined efficiency metric |

**Team Aggregate KPIs:**

| KPI | Formula | Description |
|-----|---------|-------------|
| Avg Acceptance Rate | `mean(developerAcceptanceRates)` | Team average acceptance |
| Total AI Lines | `sum(developerAcceptedLines)` | Total AI-generated code |
| Active Developers | `count(isActive = true)` | Developers with activity |
| AI Velocity Score | `mean(developerAiVelocity)` | Team's AI efficiency |

**Time-Series Aggregations:**
- Daily: Raw data from simulator (no aggregation needed)
- Weekly: Sum/average of daily values, Monday-Sunday weeks
- Monthly: Sum/average of daily values, calendar months
- Quarterly: Sum/average of monthly values

**Acceptance Criteria:**
- AC1: Acceptance rate returns 0 when no suggestions shown (not NaN or error)
- AC2: Team averages exclude inactive developers
- AC3: Weekly aggregations correctly handle weeks spanning month boundaries
- AC4: All KPIs are recalculated on each sync, not incrementally (for consistency)
- AC5: KPI values match manual spreadsheet calculations within 0.1%

### 2.4 Team Statistics Aggregation (FR-AGG-004)

Beyond individual KPIs, the system must maintain pre-computed team-level statistics for fast dashboard rendering.

**Team Stats Table Updates:**

After each sync, update the `team_stats` table with:

```sql
INSERT INTO team_stats (team_name, date, avg_acceptance_rate, total_ai_lines, 
                        total_chat_requests, active_developers, ai_velocity_score)
SELECT 
    d.team,
    ds.date,
    AVG(ds.acceptance_rate),
    SUM(ds.accepted_lines_added),
    SUM(ds.chat_requests),
    COUNT(CASE WHEN ds.total_tabs_shown > 0 OR ds.chat_requests > 0 THEN 1 END),
    AVG(CASE WHEN ds.lines_added > 0 
        THEN (ds.accepted_lines_added::float / ds.lines_added) * 100 
        ELSE 0 END)
FROM daily_stats ds
JOIN developers d ON ds.developer_id = d.id
WHERE ds.date = :today
GROUP BY d.team, ds.date
ON CONFLICT (team_name, date) DO UPDATE SET
    avg_acceptance_rate = EXCLUDED.avg_acceptance_rate,
    total_ai_lines = EXCLUDED.total_ai_lines,
    total_chat_requests = EXCLUDED.total_chat_requests,
    active_developers = EXCLUDED.active_developers,
    ai_velocity_score = EXCLUDED.ai_velocity_score;
```

**Acceptance Criteria:**
- AC1: Team stats are updated within 5 seconds of data ingestion completing
- AC2: Historical team stats are preserved when new data arrives
- AC3: Team names with special characters are handled correctly
- AC4: Queries on team_stats return in under 50ms for 100 days of data

### 2.5 GraphQL Schema and Resolvers (FR-AGG-005)

The aggregator exposes a GraphQL API that the frontend consumes. The schema must support flexible querying with appropriate depth limits and pagination.

**Core Queries:**

```graphql
type Query {
  # Single developer with nested stats
  developer(id: ID!): Developer
  
  # List developers with optional filtering
  developers(team: String, limit: Int, offset: Int): DeveloperConnection!
  
  # Single developer's usage statistics
  getDevStats(id: ID!, range: DateRange): UsageStats
  
  # Team-level statistics
  getTeamStats(teamName: String!): TeamStats
  getAllTeamStats: [TeamStats!]!
  
  # Dashboard aggregates
  getDashboardSummary(range: DateRange): DashboardKPI!
  
  # Visualization-specific queries
  getVelocityHeatmap(developerId: ID, teamName: String, days: Int): [HeatmapCell!]!
  getTeamRadarData: [TeamRadarPoint!]!
  getDeveloperEfficiencyTable(sortBy: String, sortOrder: String, limit: Int): [DeveloperEfficiencyRow!]!
  
  # Time series data
  getDailyTrend(range: DateRange): [DailyTrendPoint!]!
}

type DeveloperConnection {
  nodes: [Developer!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}
```

**Resolver Performance Requirements:**
- All resolvers must include DataLoader for N+1 prevention
- Complex aggregations should use materialized views or pre-computed tables
- Maximum query depth: 5 levels
- Maximum result set: 1000 items without pagination

**Acceptance Criteria:**
- AC1: All queries return data matching the schema exactly
- AC2: N+1 query problem is avoided (verified via query logging)
- AC3: Null handling follows GraphQL best practices (nullable vs. non-nullable)
- AC4: Query complexity limits prevent denial of service
- AC5: GraphQL Playground is available in development mode

### 2.6 Error Handling and Resilience (FR-AGG-006)

The aggregator must handle errors gracefully and provide meaningful feedback without crashing or losing data.

**Error Categories:**

| Category | Handling | User Feedback |
|----------|----------|---------------|
| Simulator Unreachable | Retry with backoff | GraphQL returns stale data flag |
| Invalid Data Format | Log and skip record | Partial success, list skipped |
| Database Connection Lost | Retry 3 times, then fail | Service unhealthy status |
| Query Timeout | Cancel and return partial | Error extension with timeout hint |

**Circuit Breaker Pattern:**

The ingestion worker implements a circuit breaker to prevent overwhelming a failing simulator:

```typescript
// services/ingestion-worker.ts
const circuitBreaker = new CircuitBreaker({
  failureThreshold: 5,      // Open after 5 consecutive failures
  successThreshold: 2,      // Close after 2 successes
  timeout: 30000,           // Request timeout
  resetTimeout: 60000,      // Time before half-open state
});
```

**Acceptance Criteria:**
- AC1: Service continues running when simulator is unavailable
- AC2: Partial data ingestion does not corrupt existing data
- AC3: Health endpoint accurately reflects connection status
- AC4: Error messages do not expose internal implementation details
- AC5: Circuit breaker prevents cascading failures

---

## 3. Non-Functional Requirements

### 3.1 Performance

- Ingestion: Process 100 developers × 90 days in under 30 seconds
- Query Response: P95 under 200ms for dashboard summary
- Database: Indexed queries only; no full table scans for typical operations

### 3.2 Reliability

- Automatic reconnection to PostgreSQL on connection loss
- Transaction rollback on partial failure
- Idempotent ingestion (re-running sync produces same result)

### 3.3 Observability

- Structured JSON logging with correlation IDs
- Prometheus metrics endpoint (/metrics)
- GraphQL query performance tracing

### 3.4 Data Integrity

- Foreign key constraints enforced at database level
- Unique constraints on (developer_id, date) for daily_stats
- No orphaned records allowed

---

## 4. Technical Design Notes

### 4.1 Database Connection Pooling

Use a connection pool to efficiently manage PostgreSQL connections:

```typescript
// config/database.ts
import { Pool } from 'pg';

export const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
  max: 20,                    // Maximum pool size
  idleTimeoutMillis: 30000,   // Close idle connections after 30s
  connectionTimeoutMillis: 5000,
});
```

### 4.2 DataLoader Pattern

Prevent N+1 queries by batching database requests:

```typescript
// loaders/developer-loader.ts
import DataLoader from 'dataloader';

export const createDeveloperLoader = (db: Pool) => 
  new DataLoader<string, Developer>(async (ids) => {
    const result = await db.query(
      'SELECT * FROM developers WHERE id = ANY($1)',
      [ids]
    );
    const map = new Map(result.rows.map(r => [r.id, r]));
    return ids.map(id => map.get(id) ?? null);
  });
```

### 4.3 Incremental Sync Strategy

Track sync state to avoid re-processing all data:

```typescript
// services/sync-state.ts
interface SyncState {
  lastSyncTime: Date;
  lastSuccessfulSync: Date;
  consecutiveFailures: number;
  totalRecordsSynced: number;
}

// Stored in database table 'sync_state'
// Queried on startup to determine where to resume
```

### 4.4 GraphQL Context Pattern

Pass loaders and database through context:

```typescript
// server.ts
const server = new ApolloServer({
  typeDefs,
  resolvers,
  context: ({ req }) => ({
    db: pool,
    loaders: {
      developer: createDeveloperLoader(pool),
      dailyStats: createDailyStatsLoader(pool),
    },
    user: req.user, // For future auth
  }),
});
```

---

## 5. Dependencies

### 5.1 External Libraries

| Library | Version | Purpose |
|---------|---------|---------|
| express | ^4.18 | HTTP server |
| @apollo/server | ^4.0 | GraphQL server |
| graphql | ^16.0 | GraphQL reference implementation |
| pg | ^8.11 | PostgreSQL client |
| dataloader | ^2.2 | Batch loading |
| node-cron | ^3.0 | Scheduled tasks |
| zod | ^3.22 | Runtime type validation |
| pino | ^8.0 | Structured logging |

### 5.2 Internal Dependencies

- **cursor-sim**: Source of usage data (HTTP REST API)
- **PostgreSQL**: Persistent data storage

---

## 6. Test Cases

### 6.1 Unit Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| UT-AGG-001 | Calculate acceptance rate with zero shown | Returns 0, not NaN |
| UT-AGG-002 | Calculate AI velocity with zero total lines | Returns 0 |
| UT-AGG-003 | Transform simulator response to DB format | All fields mapped correctly |
| UT-AGG-004 | Handle missing optional fields in source | Defaults applied |
| UT-AGG-005 | Team average excludes inactive developers | Only active in calculation |
| UT-AGG-006 | DateRange enum maps to correct day counts | DAY=1, WEEK=7, MONTH=30, QUARTER=90 |

### 6.2 Integration Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| IT-AGG-001 | Ingest data from simulator | Data persisted to PostgreSQL |
| IT-AGG-002 | Duplicate ingestion is idempotent | No duplicate records |
| IT-AGG-003 | GraphQL getDashboardSummary returns data | Valid KPIs in response |
| IT-AGG-004 | DataLoader batches developer queries | Single SQL query for N developers |
| IT-AGG-005 | Team stats aggregation is correct | Matches manual calculation |
| IT-AGG-006 | Pagination returns correct pages | Consistent ordering, no duplicates |

### 6.3 GraphQL Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| GQL-AGG-001 | Query developer by ID | Returns matching developer |
| GQL-AGG-002 | Query nonexistent developer | Returns null, no error |
| GQL-AGG-003 | Query with excessive depth | Rejected with depth limit error |
| GQL-AGG-004 | Query with invalid DateRange | Validation error returned |
| GQL-AGG-005 | Nested stats resolve correctly | Child data matches parent filter |

---

## 7. Related User Stories

- [US-AGG-001](../user-stories/US-AGG-001-sync-data.md): Sync Data from Simulator
- [US-AGG-002](../user-stories/US-AGG-002-calculate-kpis.md): Calculate KPIs
- [US-AGG-003](../user-stories/US-AGG-003-query-stats.md): Query Statistics via GraphQL
- [US-AGG-004](../user-stories/US-AGG-004-team-comparison.md): Compare Teams

---

## 8. Implementation Tasks

- [TASK-010](../tasks/TASK-010-agg-project-setup.md): Set up Node.js TypeScript project
- [TASK-011](../tasks/TASK-011-agg-prisma.md): Configure Prisma and database schema
- [TASK-012](../tasks/TASK-012-agg-ingestion.md): Implement ingestion worker
- [TASK-013](../tasks/TASK-013-agg-kpi-calculator.md): Implement KPI calculation service
- [TASK-014](../tasks/TASK-014-agg-graphql-schema.md): Define GraphQL schema
- [TASK-015](../tasks/TASK-015-agg-resolvers.md): Implement GraphQL resolvers
- [TASK-016](../tasks/TASK-016-agg-dataloaders.md): Implement DataLoaders
- [TASK-017](../tasks/TASK-017-agg-docker.md): Create Dockerfile and migrations
