# System Design Document: Cursor Usage Analytics Platform

**Version**: 1.0.0  
**Last Updated**: January 2026  
**Status**: Draft  

## 1. Executive Summary

The Cursor Usage Analytics Platform is a microservices-based system designed to simulate, aggregate, and visualize AI coding assistant usage metrics. The platform enables engineering teams to understand how developers interact with Cursor IDE's AI features, measure adoption rates, and identify opportunities for improving developer productivity.

The system consists of three decoupled services following the ETL (Extract, Transform, Load) pattern, where the Simulator generates realistic synthetic data, the Aggregator processes and calculates KPIs, and the Dashboard provides interactive visualizations.

## 2. Problem Statement

Organizations adopting AI coding assistants like Cursor lack visibility into how effectively their teams utilize these tools. The official Cursor for Business API provides some metrics, but teams need a way to analyze usage patterns at scale, compare team performance, and identify developers who may benefit from additional training.

This platform addresses these needs by providing a complete analytics solution that can work with both simulated data (for development and demos) and real Cursor API data (for production use).

## 3. System Architecture

### 3.1 High-Level Architecture

The system implements a three-tier architecture with clear separation of concerns between data generation, processing, and presentation. Each service operates independently and communicates through well-defined API contracts.

```
┌────────────────────────────────────────────────────────────────────────────┐
│                           Docker Network (cursor-net)                       │
├────────────────────────────────────────────────────────────────────────────┤
│                                                                            │
│  ┌──────────────────┐    REST     ┌──────────────────┐                    │
│  │                  │   (JSON)    │                  │                    │
│  │   cursor-sim     │────────────▶│ cursor-analytics │                    │
│  │   (Go)           │             │     -core        │                    │
│  │                  │             │   (TypeScript)   │                    │
│  │   Port: 8080     │             │   Port: 4000     │                    │
│  └──────────────────┘             └────────┬─────────┘                    │
│         │                                  │                              │
│         │ In-Memory                        │ PostgreSQL                   │
│         │ Storage                          │                              │
│         ▼                                  ▼                              │
│  ┌──────────────────┐             ┌──────────────────┐                    │
│  │   sync.Map       │             │   cursor_db      │                    │
│  │   (volatile)     │             │   (persistent)   │                    │
│  └──────────────────┘             └──────────────────┘                    │
│                                            │                              │
│                                   GraphQL  │                              │
│                                            ▼                              │
│                               ┌──────────────────┐                        │
│                               │                  │                        │
│                               │  cursor-viz-spa  │                        │
│                               │  (React)         │                        │
│                               │                  │                        │
│                               │   Port: 3000     │                        │
│                               └──────────────────┘                        │
│                                                                           │
└───────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Data Flow Sequence

The following sequence diagram illustrates the complete data flow from simulation to visualization:

```
User          cursor-sim      cursor-analytics-core     PostgreSQL     cursor-viz-spa
 │                │                    │                     │               │
 │  CLI Start     │                    │                     │               │
 │───────────────▶│                    │                     │               │
 │                │                    │                     │               │
 │                │  Generate          │                     │               │
 │                │  Synthetic Data    │                     │               │
 │                │  (Poisson dist.)   │                     │               │
 │                │────────────┐       │                     │               │
 │                │            │       │                     │               │
 │                │◀───────────┘       │                     │               │
 │                │                    │                     │               │
 │                │    Cron Poll       │                     │               │
 │                │◀───────────────────│                     │               │
 │                │                    │                     │               │
 │                │  GET /v1/stats/    │                     │               │
 │                │     activity       │                     │               │
 │                │───────────────────▶│                     │               │
 │                │                    │                     │               │
 │                │  JSON Events       │                     │               │
 │                │◀───────────────────│                     │               │
 │                │                    │                     │               │
 │                │                    │  Normalize &        │               │
 │                │                    │  Calculate KPIs     │               │
 │                │                    │────────────────────▶│               │
 │                │                    │                     │               │
 │                │                    │  Store Metrics      │               │
 │                │                    │────────────────────▶│               │
 │                │                    │                     │               │
 │                │                    │     GraphQL Query   │               │
 │                │                    │◀────────────────────│───────────────│
 │                │                    │                     │               │
 │                │                    │     Query DB        │               │
 │                │                    │────────────────────▶│               │
 │                │                    │                     │               │
 │                │                    │     Results         │               │
 │                │                    │◀────────────────────│               │
 │                │                    │                     │               │
 │                │                    │     GraphQL Response│               │
 │                │                    │─────────────────────│──────────────▶│
 │                │                    │                     │               │
 │                │                    │                     │     Render    │
 │◀──────────────────────────────────────────────────────────│───────────────│
 │                                                                  Dashboard
```

### 3.3 Service Boundaries

Each service has clear responsibilities and owns its data domain. The boundaries are designed to allow independent scaling, deployment, and technology evolution.

**cursor-sim (Data Generator)**
The simulator owns the domain of synthetic data generation, including developer profiles, team structures, and usage event creation. It provides a read-only API that mimics the Cursor Business API contract. The service is stateless between restarts, generating fresh data on each startup based on CLI parameters.

**cursor-analytics-core (Data Processor)**
The aggregator owns the domain of data persistence, normalization, and metric calculation. It maintains the source of truth for all historical analytics data. The service exposes a GraphQL API that serves as the single entry point for all analytical queries.

**cursor-viz-spa (Data Consumer)**
The dashboard owns the domain of data visualization and user interaction. It is a pure consumer of the GraphQL API, maintaining no persistent state beyond user preferences. The service handles all presentation logic and client-side caching.

## 4. Service Specifications

### 4.1 Service A: cursor-sim (Cursor API Simulator)

The simulator creates realistic synthetic data patterns using statistical distributions that mimic actual developer behavior. Rather than generating purely random data, it employs Poisson distributions for event timing and configurable acceptance rates based on developer experience levels.

**Technology Choices:**

The Go programming language was selected for its excellent concurrency primitives, fast startup time, and low memory footprint, making it ideal for a CLI tool that may generate hundreds of concurrent developer simulations. The built-in `net/http` package provides all necessary HTTP server capabilities without external dependencies.

**Data Model:**

The simulator maintains an in-memory representation of the organization structure, mapping developers to teams and generating events that reflect realistic usage patterns. Senior developers show higher suggestion acceptance rates (approximately 90%) compared to junior developers (approximately 60%), reflecting real-world observations that experienced developers better understand how to leverage AI assistance.

```go
// Core domain types for the simulator
type Developer struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    Email         string    `json:"email"`
    Team          string    `json:"team"`
    Seniority     string    `json:"seniority"` // "junior", "mid", "senior"
    AcceptanceRate float64  `json:"acceptanceRate"`
    CreatedAt     time.Time `json:"createdAt"`
}

type UsageEvent struct {
    ID          string    `json:"id"`
    DeveloperID string    `json:"developerId"`
    EventType   string    `json:"eventType"`
    Timestamp   time.Time `json:"timestamp"`
    Metadata    EventMeta `json:"metadata"`
}

type EventMeta struct {
    LinesAdded    int    `json:"linesAdded,omitempty"`
    LinesDeleted  int    `json:"linesDeleted,omitempty"`
    ModelUsed     string `json:"modelUsed,omitempty"`
    Accepted      bool   `json:"accepted,omitempty"`
    TokensInput   int    `json:"tokensInput,omitempty"`
    TokensOutput  int    `json:"tokensOutput,omitempty"`
}
```

**Event Types:**

The simulator generates four primary event types aligned with the Cursor telemetry schema discovered through API research:

| Event Type | Description | Generation Rate |
|------------|-------------|-----------------|
| `cpp_suggestion_shown` | Tab completion displayed | Base rate × velocity |
| `cpp_suggestion_accepted` | Tab completion accepted | Shown × acceptance_rate |
| `chat_message` | Chat interaction started | 15% of suggestion rate |
| `cmd_k_prompt` | Inline edit command used | 10% of suggestion rate |

**API Contract:**

The simulator exposes REST endpoints that mirror the Cursor Business API structure, enabling the aggregator to work with either simulated or real data without code changes.

```
GET  /v1/org/users                          # List all developers
GET  /v1/org/users/:id                      # Get single developer
GET  /v1/stats/activity?from={ts}&to={ts}   # Activity events
GET  /v1/stats/daily-usage                  # Aggregated daily metrics
GET  /health                                # Service health check
```

### 4.2 Service B: cursor-analytics-core (Aggregator Service)

The aggregator serves as the analytical engine, transforming raw events into actionable insights. It implements a background polling mechanism to fetch data from the simulator and maintains a PostgreSQL database for persistent storage of normalized metrics.

**Technology Choices:**

TypeScript with Node.js was selected for its strong typing support, excellent GraphQL tooling (Apollo Server), and rich ecosystem of data processing libraries. PostgreSQL provides reliable persistent storage with excellent support for time-series queries and aggregations.

**Database Schema:**

The database schema is optimized for analytical queries, with proper indexing on frequently filtered columns like timestamps and team identifiers.

```sql
-- Developers table stores profile information
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

-- Events table stores raw usage events
CREATE TABLE usage_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id UUID REFERENCES developers(id),
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

-- Materialized view for daily aggregations
CREATE MATERIALIZED VIEW daily_stats AS
SELECT 
    developer_id,
    DATE(event_timestamp) as date,
    COUNT(*) FILTER (WHERE event_type = 'cpp_suggestion_shown') as suggestions_shown,
    COUNT(*) FILTER (WHERE event_type = 'cpp_suggestion_accepted') as suggestions_accepted,
    COUNT(*) FILTER (WHERE event_type = 'chat_message') as chat_interactions,
    COUNT(*) FILTER (WHERE event_type = 'cmd_k_prompt') as cmd_k_usages,
    SUM(lines_added) as total_lines_added,
    SUM(lines_deleted) as total_lines_deleted
FROM usage_events
GROUP BY developer_id, DATE(event_timestamp);

-- Indexes for query performance
CREATE INDEX idx_events_developer ON usage_events(developer_id);
CREATE INDEX idx_events_timestamp ON usage_events(event_timestamp);
CREATE INDEX idx_events_type ON usage_events(event_type);
CREATE INDEX idx_daily_stats_date ON daily_stats(date);
```

**GraphQL Schema:**

The GraphQL API provides a flexible query interface for the frontend dashboard, supporting various aggregation levels and time ranges.

```graphql
scalar DateTime

type Developer {
    id: ID!
    externalId: String!
    name: String!
    email: String!
    team: String!
    seniority: String
    stats(range: DateRangeInput): UsageStats
    dailyStats(range: DateRangeInput): [DailyStats!]!
}

type UsageStats {
    totalSuggestions: Int!
    acceptedSuggestions: Int!
    acceptanceRate: Float!
    chatInteractions: Int!
    cmdKUsages: Int!
    totalLinesAdded: Int!
    totalLinesDeleted: Int!
    aiVelocity: Float!
}

type DailyStats {
    date: DateTime!
    suggestionsShown: Int!
    suggestionsAccepted: Int!
    acceptanceRate: Float!
    chatInteractions: Int!
    linesAdded: Int!
    linesDeleted: Int!
}

type TeamStats {
    teamName: String!
    memberCount: Int!
    averageAcceptanceRate: Float!
    totalSuggestions: Int!
    totalAccepted: Int!
    chatInteractions: Int!
    topPerformer: Developer
}

type DashboardKPI {
    totalDevelopers: Int!
    activeDevelopers: Int!
    overallAcceptanceRate: Float!
    totalSuggestionsToday: Int!
    totalAcceptedToday: Int!
    teamComparison: [TeamStats!]!
    dailyTrend: [DailyStats!]!
}

input DateRangeInput {
    from: DateTime!
    to: DateTime!
}

type Query {
    developer(id: ID!): Developer
    developers(team: String, limit: Int, offset: Int): [Developer!]!
    teamStats(teamName: String!): TeamStats
    teams: [TeamStats!]!
    dashboardSummary(range: DateRangeInput): DashboardKPI!
}
```

**Metric Calculations:**

The aggregator computes several key performance indicators that provide insight into AI tool effectiveness:

| Metric | Formula | Purpose |
|--------|---------|---------|
| Acceptance Rate | `(accepted / shown) × 100` | Measures suggestion quality |
| AI Velocity | `(AI lines / total lines) × 100` | Measures AI contribution |
| Chat Dependency | `chat_count / code_events` | Measures reliance on chat |
| Productivity Index | Composite score | Overall effectiveness |

### 4.3 Service C: cursor-viz-spa (Frontend Dashboard)

The dashboard provides an intuitive interface for exploring analytics data through interactive visualizations. Built with React and modern tooling, it offers responsive charts and tables that update in real-time as new data arrives.

**Technology Choices:**

React with Vite was selected for fast development iteration and optimal production builds. TanStack Query (formerly React Query) handles server state management and caching, eliminating the need for complex global state solutions. Recharts provides declarative charting components that integrate naturally with React's component model.

**Key Dashboard Views:**

The dashboard presents three primary visualization panels, each designed to answer specific questions about AI tool usage:

**Velocity Heatmap** - This GitHub-style contribution graph uses color intensity to show when teams rely most heavily on Cursor. Darker cells indicate higher AI code acceptance, helping identify patterns in AI adoption over time. The heatmap answers questions like "Which days see the highest AI assistance?" and "Are there weekly patterns in tool usage?"

**Team Comparison Radar** - This multi-axis chart enables side-by-side comparison of teams across dimensions like Chat Usage, Code Completion rate, and Refactoring prompt frequency. Engineering managers can quickly identify which teams are underutilizing AI features or may benefit from additional training.

**Developer Efficiency Table** - This sortable table displays individual metrics with visual indicators for performance thresholds. Rows showing acceptance rates below 20% are highlighted to flag potential issues such as poor suggestion quality, developer unfamiliarity with the tool, or misaligned AI model settings.

**Component Architecture:**

```
src/
├── components/
│   ├── charts/
│   │   ├── VelocityHeatmap.tsx      # GitHub-style contribution graph
│   │   ├── TeamRadarChart.tsx       # Multi-axis team comparison
│   │   └── TrendLineChart.tsx       # Time-series trend display
│   ├── tables/
│   │   └── DeveloperTable.tsx       # Sortable efficiency table
│   ├── layout/
│   │   ├── Header.tsx               # KPI summary bar
│   │   ├── Sidebar.tsx              # Navigation and filters
│   │   └── DashboardGrid.tsx        # Responsive layout container
│   └── common/
│       ├── LoadingState.tsx         # Skeleton loaders
│       ├── ErrorBoundary.tsx        # Error handling wrapper
│       └── DateRangePicker.tsx      # Time range selection
├── hooks/
│   ├── useDevStats.ts               # Developer statistics query
│   ├── useTeamStats.ts              # Team statistics query
│   └── useDashboard.ts              # Dashboard summary query
├── graphql/
│   ├── queries.ts                   # GraphQL query definitions
│   └── fragments.ts                 # Reusable query fragments
└── pages/
    ├── Dashboard.tsx                # Main dashboard page
    ├── TeamView.tsx                 # Team detail page
    └── DeveloperView.tsx            # Developer detail page
```

## 5. Technical Decisions

### 5.1 Why Microservices?

The microservices architecture was chosen for several reasons specific to this project's requirements:

**Independent Development** - Each service can be developed by different team members or contractors without tight coordination, as long as API contracts are respected.

**Technology Flexibility** - Go's performance characteristics make it ideal for data generation, while TypeScript's type safety benefits the complex business logic in the aggregator.

**Selective Replacement** - The simulator can be swapped for real Cursor API integration by simply changing the aggregator's data source configuration, without touching the frontend.

**Development Workflow** - Frontend developers can work against the simulator while backend changes are in progress, and vice versa.

### 5.2 Why GraphQL?

GraphQL was selected over REST for the aggregator-to-frontend communication for the following reasons:

**Query Flexibility** - Dashboard components can request exactly the data they need, reducing over-fetching common with REST endpoints that return fixed response shapes.

**Type Safety** - The GraphQL schema serves as a contract that can be used to generate TypeScript types, ensuring type safety across the API boundary.

**Introspection** - Developers can explore the available data through GraphQL Playground, accelerating frontend development and debugging.

**Caching** - Apollo Client provides sophisticated caching that reduces unnecessary network requests and improves perceived performance.

### 5.3 Why In-Memory Storage for Simulator?

The simulator uses in-memory storage rather than a persistent database because:

**Simplicity** - No database setup required to run the simulator, reducing the barrier to entry for new developers.

**Reproducibility** - Each run generates fresh data based on the same random seed, making it easier to reproduce and debug issues.

**Performance** - In-memory operations are extremely fast, allowing the simulator to generate events for hundreds of developers without performance concerns.

**Appropriate for Purpose** - The simulator's data is ephemeral by design; historical accuracy is handled by the aggregator's persistent storage.

## 6. Deployment Architecture

### 6.1 Docker Compose Configuration

The platform runs on developer workstations using Docker Compose, providing a consistent environment across different operating systems.

```yaml
version: '3.8'

services:
  cursor-sim:
    build: ./services/cursor-sim
    ports:
      - "8080:8080"
    environment:
      - CURSOR_SIM_PORT=8080
      - CURSOR_SIM_DEVELOPERS=50
      - CURSOR_SIM_VELOCITY=high
      - CURSOR_SIM_FLUCTUATION=0.2
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3

  postgres:
    image: postgres:15-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=cursor
      - POSTGRES_PASSWORD=cursor_dev
      - POSTGRES_DB=cursor_analytics
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U cursor"]
      interval: 10s
      timeout: 5s
      retries: 5

  cursor-analytics-core:
    build: ./services/cursor-analytics-core
    ports:
      - "4000:4000"
    environment:
      - DATABASE_URL=postgresql://cursor:cursor_dev@postgres:5432/cursor_analytics
      - SIMULATOR_URL=http://cursor-sim:8080
      - POLL_INTERVAL_MS=60000
      - NODE_ENV=development
    depends_on:
      postgres:
        condition: service_healthy
      cursor-sim:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:4000/health"]
      interval: 10s
      timeout: 5s
      retries: 3

  cursor-viz-spa:
    build: ./services/cursor-viz-spa
    ports:
      - "3000:3000"
    environment:
      - VITE_GRAPHQL_URL=http://localhost:4000/graphql
    depends_on:
      cursor-analytics-core:
        condition: service_healthy

volumes:
  postgres_data:

networks:
  default:
    name: cursor-net
```

### 6.2 Health Checks and Resilience

Each service implements health check endpoints that Docker uses to determine service readiness:

| Service | Endpoint | Checks |
|---------|----------|--------|
| cursor-sim | `GET /health` | Server responding |
| cursor-analytics-core | `GET /health` | DB connection, simulator reachable |
| cursor-viz-spa | `GET /health` | Server responding |

The aggregator implements retry logic with exponential backoff when fetching from the simulator, ensuring temporary network issues don't cause data loss.

## 7. Security Considerations

### 7.1 Development Environment Security

Since this platform runs on developer workstations for development and demo purposes, security considerations focus on reasonable defaults:

**No External Exposure** - Docker Compose binds services to localhost only, preventing external access.

**Generated Data** - The simulator produces synthetic data, eliminating concerns about exposing real developer information during development.

**Credentials** - Default credentials in docker-compose.yml are clearly marked as development-only; production deployments would require proper secrets management.

### 7.2 Future Production Considerations

For production deployment, the following additional measures would be required:

**Authentication** - Add JWT-based authentication to the GraphQL API, integrating with the organization's identity provider.

**HTTPS** - All service-to-service communication should use TLS, with certificates managed through a proper PKI.

**API Keys** - Integration with real Cursor API would require secure storage of API keys, potentially using Vault or a cloud secrets manager.

**Audit Logging** - Track who accesses what data for compliance and security monitoring.

## 8. Observability

### 8.1 Logging

Each service implements structured JSON logging with consistent fields:

```json
{
  "timestamp": "2026-01-15T10:30:00.000Z",
  "level": "info",
  "service": "cursor-analytics-core",
  "message": "Fetched 150 events from simulator",
  "traceId": "abc123",
  "duration_ms": 45
}
```

### 8.2 Metrics

Key metrics exposed by each service for monitoring:

**cursor-sim**
- `cursor_sim_developers_total` - Number of simulated developers
- `cursor_sim_events_generated_total` - Events generated by type
- `cursor_sim_requests_total` - API requests handled

**cursor-analytics-core**
- `cursor_core_poll_duration_seconds` - Time to fetch and process data
- `cursor_core_events_processed_total` - Events ingested by type
- `cursor_core_graphql_requests_total` - GraphQL queries by operation

**cursor-viz-spa**
- `cursor_viz_page_loads_total` - Page views by route
- `cursor_viz_query_duration_seconds` - Client-side query timing

## 9. Future Enhancements

### 9.1 Planned Features

**Real Cursor API Integration** - Replace the simulator with direct integration to Cursor Business API for production analytics.

**Alert System** - Notify team leads when acceptance rates drop below thresholds or usage patterns indicate issues.

**Export Functionality** - Allow downloading of analytics data in CSV or JSON format for further analysis.

**Custom Dashboards** - Enable users to create and save custom dashboard configurations.

### 9.2 Scalability Path

The current architecture supports vertical scaling through Docker resource limits. For larger deployments, the following horizontal scaling strategies are available:

**Aggregator** - Can be scaled horizontally with a shared PostgreSQL database, using connection pooling.

**Database** - PostgreSQL can be replaced with TimescaleDB for better time-series performance at scale.

**Frontend** - Static assets can be served from a CDN, with the GraphQL endpoint load-balanced across multiple aggregator instances.

## 10. Appendix

### 10.1 Glossary

| Term | Definition |
|------|------------|
| Acceptance Rate | Percentage of AI suggestions accepted by developers |
| AI Velocity | Ratio of AI-generated code to total code written |
| CPP | Cursor Plus Plus, the autocomplete feature (not C++) |
| KPI | Key Performance Indicator |
| Tab Completion | Inline code suggestions triggered by Tab key |

### 10.2 References

- Cursor Business API Documentation: https://docs.cursor.com/account/teams/admin-api
- Apollo GraphQL Server: https://www.apollographql.com/docs/apollo-server/
- Recharts Documentation: https://recharts.org/
- TanStack Query: https://tanstack.com/query/
