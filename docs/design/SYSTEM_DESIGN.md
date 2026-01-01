# System Design Document

## Cursor Usage Analytics Platform

**Version:** 1.0.0  
**Last Updated:** 2025-01-01  
**Status:** Design Complete

---

## 1. Executive Summary

The Cursor Usage Analytics Platform provides engineering teams with visibility into their AI-assisted development workflows. By simulating, collecting, aggregating, and visualizing usage data from Cursor IDE, teams can understand adoption patterns, identify productivity bottlenecks, and optimize their use of AI coding assistance.

This document describes the architecture of three decoupled microservices that together form an ETL pipeline: a Go-based simulator generating realistic telemetry, a TypeScript aggregator computing KPIs and exposing GraphQL, and a React dashboard providing interactive visualizations.

---

## 2. Goals and Non-Goals

### 2.1 Goals

**G1: Realistic Simulation** - Generate synthetic usage data that statistically resembles real Cursor Business API output, including Poisson-distributed events, developer skill-based acceptance rates, and temporal patterns (work hours, sprints).

**G2: Accurate Analytics** - Compute meaningful KPIs including acceptance rates, AI velocity scores, chat dependency ratios, and team comparisons with correct aggregation logic.

**G3: Interactive Visualization** - Provide a responsive dashboard with multiple chart types (heatmaps, radar charts, tables) that update in near-real-time.

**G4: Developer Experience** - Enable rapid local development using Docker Compose with hot-reload, comprehensive testing, and clear documentation.

**G5: Spec-Driven Development** - Structure the project to support AI-assisted coding through detailed specifications that translate directly to test cases and implementation tasks.

### 2.2 Non-Goals

**NG1:** Production deployment to cloud environments (this is a local development tool)

**NG2:** Authentication and authorization (single-user local application)

**NG3:** Integration with actual Cursor Business API (simulator replaces real data)

**NG4:** Historical data persistence beyond container lifecycle (volatile by design)

**NG5:** Multi-tenancy or organization management

---

## 3. System Architecture

### 3.1 High-Level Architecture

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                            Docker Compose Network                             │
│                                                                              │
│  ┌─────────────────┐    ┌─────────────────────────┐    ┌─────────────────┐  │
│  │                 │    │                         │    │                 │  │
│  │   cursor-sim    │───▶│  cursor-analytics-core  │───▶│ cursor-viz-spa  │  │
│  │                 │    │                         │    │                 │  │
│  │   Go 1.22+      │    │   Node.js 20+           │    │  React 18       │  │
│  │   Port: 8080    │    │   TypeScript 5.x        │    │  Vite 5.x       │  │
│  │                 │    │   Port: 4000            │    │  Port: 3000     │  │
│  │   In-Memory DB  │    │                         │    │                 │  │
│  │   (SQLite)      │    │   ┌─────────────────┐   │    │  TanStack Query │  │
│  │                 │    │   │   PostgreSQL    │   │    │  Recharts       │  │
│  │   Poisson Gen   │    │   │   Port: 5432    │   │    │                 │  │
│  │                 │    │   └─────────────────┘   │    │                 │  │
│  └─────────────────┘    │                         │    └─────────────────┘  │
│         │               │   Apollo Server         │            │            │
│         │               │   GraphQL               │            │            │
│         │               └─────────────────────────┘            │            │
│         │                         │                            │            │
│         └─────────────────────────┴────────────────────────────┘            │
│                              REST + GraphQL                                  │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Data Flow

The system implements an ETL pattern with clear boundaries between services:

```
                                    ETL Pipeline Flow
                                    
    ┌─────────┐     ┌──────────────────────────────────────────────────────┐
    │         │     │                                                      │
    │ EXTRACT │     │   cursor-sim generates events in real-time           │
    │         │     │   └── Poisson distribution models request patterns   │
    │         │     │   └── Developer profiles determine acceptance rates  │
    │         │     │   └── REST API exposes /v1/stats/activity endpoint   │
    │         │     │                                                      │
    └────┬────┘     └──────────────────────────────────────────────────────┘
         │
         ▼
    ┌─────────┐     ┌──────────────────────────────────────────────────────┐
    │         │     │                                                      │
    │TRANSFORM│     │   cursor-analytics-core polls and processes          │
    │         │     │   └── Cron job fetches from simulator every N mins   │
    │         │     │   └── Raw events normalized into DailyStats          │
    │         │     │   └── KPIs computed: acceptance rate, AI velocity    │
    │         │     │   └── Team aggregations calculated                   │
    │         │     │                                                      │
    └────┬────┘     └──────────────────────────────────────────────────────┘
         │
         ▼
    ┌─────────┐     ┌──────────────────────────────────────────────────────┐
    │         │     │                                                      │
    │  LOAD   │     │   cursor-viz-spa queries and renders                 │
    │         │     │   └── GraphQL queries via TanStack Query             │
    │         │     │   └── Real-time updates via polling (30s default)    │
    │         │     │   └── Charts rendered with Recharts                  │
    │         │     │                                                      │
    └─────────┘     └──────────────────────────────────────────────────────┘
```

### 3.3 Service Responsibilities

**Service A: cursor-sim (Simulator)**

The simulator serves as a mock implementation of the Cursor Business Admin API. It generates synthetic but statistically realistic usage data that mimics what a real Cursor for Business deployment would produce.

Core responsibilities include generating developer profiles with varied characteristics (seniority, team, working patterns), producing usage events following Poisson distributions to model realistic request patterns, exposing REST endpoints that match Cursor Business API conventions, and maintaining in-memory state for fast response times.

**Service B: cursor-analytics-core (Aggregator)**

The aggregator serves as the data processing engine and API gateway for the frontend. It transforms raw events into actionable metrics and provides a GraphQL interface for flexible querying.

Core responsibilities include polling the simulator on a configurable interval, normalizing heterogeneous event data into a consistent schema, computing derived metrics like acceptance rates and AI velocity scores, storing aggregated data in PostgreSQL for historical queries, and exposing a GraphQL API with optimized resolvers.

**Service C: cursor-viz-spa (Dashboard)**

The dashboard provides interactive visualizations for understanding team and individual AI usage patterns. It focuses on clear data presentation and responsive user experience.

Core responsibilities include fetching data via GraphQL with intelligent caching, rendering multiple chart types (heatmaps, radar charts, tables), providing filtering and drill-down capabilities, and updating automatically to reflect new data.

---

## 4. Service Specifications

### 4.1 cursor-sim (Go)

#### 4.1.1 Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Runtime | Go 1.22+ | High-performance concurrent simulation |
| HTTP Server | Chi Router | Lightweight, idiomatic middleware support |
| In-Memory DB | go-sqlite3 (memory mode) | Fast volatile storage with SQL queries |
| CLI Framework | Cobra | Industry-standard Go CLI |
| Testing | testify + mockery | Assertions and mock generation |
| Random | math/rand/v2 | Deterministic seeding for reproducibility |

#### 4.1.2 CLI Configuration

```go
// cmd/root.go
type SimulatorConfig struct {
    Port         int     `mapstructure:"port"`        // Default: 8080
    Developers   int     `mapstructure:"developers"`  // Default: 50
    Velocity     string  `mapstructure:"velocity"`    // low|medium|high
    Fluctuation  float64 `mapstructure:"fluctuation"` // 0.0-1.0
    Seed         int64   `mapstructure:"seed"`        // For reproducibility
    HistoryDays  int     `mapstructure:"history"`     // Days of data to generate
}
```

#### 4.1.3 Data Generation Model

The simulator creates developer profiles with statistical distributions that model real-world patterns:

```
Developer Profile Generation
├── Name: Faker library (realistic names)
├── Team: Random from ["Backend", "Frontend", "Platform", "Mobile", "Data"]
├── Role: Weighted random (20% Senior, 50% Mid, 30% Junior)
├── Acceptance Rate Baseline:
│   ├── Senior: 85-95% (experienced, better prompts)
│   ├── Mid: 65-80% (learning optimal usage)
│   └── Junior: 40-60% (still adapting)
└── Activity Pattern:
    ├── Work hours: Normal distribution centered on 10am-6pm
    ├── Days: Higher activity Mon-Thu, lower Fri
    └── Sprint effects: +20% activity in sprint weeks
```

Event generation uses Poisson distribution to model realistic request patterns:

```go
// internal/generator/poisson.go
func generateEvents(dev Developer, velocity Velocity, duration time.Duration) []Event {
    // λ (lambda) = expected events per hour based on velocity
    // low=10, medium=50, high=100
    lambda := velocityToLambda(velocity)
    
    // Apply developer-specific modifier
    lambda *= dev.ActivityMultiplier
    
    // Generate inter-arrival times using exponential distribution
    // This creates Poisson-distributed event counts per interval
    for t := 0; t < duration; t += interval {
        count := rand.Poisson(lambda * interval.Hours())
        for i := 0; i < count; i++ {
            events = append(events, generateSingleEvent(dev, t))
        }
    }
    return events
}
```

#### 4.1.4 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/org/users` | GET | List all simulated developers |
| `/v1/stats/activity` | GET | Usage events with time range filter |
| `/v1/stats/daily` | GET | Pre-aggregated daily statistics |
| `/v1/events` | GET | Raw event stream |
| `/health` | GET | Health check |

**GET /v1/stats/activity Response:**

```json
{
  "data": [
    {
      "email": "alice@example.com",
      "date": "2025-01-01",
      "isActive": true,
      "totalTabsShown": 142,
      "totalTabsAccepted": 98,
      "totalLinesAdded": 487,
      "totalLinesDeleted": 123,
      "acceptedLinesAdded": 312,
      "composerRequests": 8,
      "chatRequests": 23,
      "agentRequests": 2,
      "cmdkUsages": 15,
      "mostUsedModel": "claude-sonnet-4"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 100,
    "totalPages": 3,
    "hasNextPage": true
  }
}
```

### 4.2 cursor-analytics-core (TypeScript)

#### 4.2.1 Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Runtime | Node.js 20+ | LTS with native TypeScript support |
| Framework | Express.js | Lightweight, well-known |
| GraphQL | Apollo Server 4 | Industry standard, excellent tooling |
| Database | PostgreSQL 15 | Robust relational with JSONB support |
| ORM | Prisma | Type-safe queries, migrations |
| Scheduler | node-cron | Simple cron-like scheduling |
| Testing | Jest + ts-jest | Full TypeScript support |
| Testing DB | pg-mem | In-memory PostgreSQL for tests |

#### 4.2.2 Database Schema

```sql
-- migrations/001_initial.sql

CREATE TABLE developers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    team VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE usage_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id UUID REFERENCES developers(id),
    event_type VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_events_developer_time ON usage_events(developer_id, timestamp);
CREATE INDEX idx_events_type ON usage_events(event_type);

CREATE TABLE daily_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id UUID REFERENCES developers(id),
    date DATE NOT NULL,
    total_tabs_shown INT DEFAULT 0,
    total_tabs_accepted INT DEFAULT 0,
    lines_added INT DEFAULT 0,
    lines_deleted INT DEFAULT 0,
    accepted_lines_added INT DEFAULT 0,
    chat_requests INT DEFAULT 0,
    composer_requests INT DEFAULT 0,
    agent_requests INT DEFAULT 0,
    cmdk_usages INT DEFAULT 0,
    most_used_model VARCHAR(50),
    acceptance_rate FLOAT GENERATED ALWAYS AS (
        CASE WHEN total_tabs_shown > 0 
        THEN total_tabs_accepted::FLOAT / total_tabs_shown 
        ELSE 0 END
    ) STORED,
    UNIQUE(developer_id, date)
);

CREATE TABLE team_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_name VARCHAR(100) NOT NULL,
    date DATE NOT NULL,
    avg_acceptance_rate FLOAT,
    total_ai_lines INT,
    total_chat_requests INT,
    active_developers INT,
    ai_velocity_score FLOAT,
    UNIQUE(team_name, date)
);
```

#### 4.2.3 GraphQL Schema

```graphql
# schema.graphql

scalar DateTime
scalar Date

enum DateRange {
  DAY
  WEEK
  MONTH
  QUARTER
}

type Developer {
  id: ID!
  email: String!
  name: String!
  team: String!
  role: String!
  stats(range: DateRange): UsageStats
  dailyStats(startDate: Date, endDate: Date): [DailyStats!]!
}

type UsageStats {
  totalSuggestions: Int!
  acceptedSuggestions: Int!
  acceptanceRate: Float!
  totalLinesAdded: Int!
  acceptedLinesAdded: Int!
  chatInteractions: Int!
  composerRequests: Int!
  agentRequests: Int!
  cmdkUsages: Int!
  aiVelocity: Float!
  chatDependency: Float!
}

type DailyStats {
  date: Date!
  totalTabsShown: Int!
  totalTabsAccepted: Int!
  linesAdded: Int!
  linesDeleted: Int!
  acceptedLinesAdded: Int!
  chatRequests: Int!
  composerRequests: Int!
  acceptanceRate: Float!
  mostUsedModel: String
}

type TeamStats {
  teamName: String!
  avgAcceptanceRate: Float!
  totalAiLines: Int!
  totalChatRequests: Int!
  activeDevelopers: Int!
  aiVelocityScore: Float!
  developers: [Developer!]!
}

type CommitBreakdown {
  type: String!
  count: Int!
  aiUtilization: Float!
}

type DashboardKPI {
  totalDevelopers: Int!
  activeDevelopers: Int!
  avgAcceptanceRate: Float!
  totalAiLinesThisWeek: Int!
  topPerformer: Developer
  teamComparison: [TeamStats!]!
  dailyTrend: [DailyTrendPoint!]!
}

type DailyTrendPoint {
  date: Date!
  totalEvents: Int!
  acceptanceRate: Float!
  aiLinesAdded: Int!
}

type HeatmapCell {
  date: Date!
  weekday: Int!
  week: Int!
  intensity: Float!
  value: Int!
}

type Query {
  # Developer queries
  developer(id: ID!): Developer
  developers(team: String, limit: Int): [Developer!]!
  
  # Statistics queries
  getDevStats(id: ID!): UsageStats
  getTeamStats(teamName: String!): TeamStats
  getAllTeamStats: [TeamStats!]!
  
  # Dashboard queries
  getDashboardSummary(range: DateRange): DashboardKPI!
  getVelocityHeatmap(developerId: ID, teamName: String): [HeatmapCell!]!
  getTeamRadarData: [TeamRadarPoint!]!
  
  # Table queries
  getDeveloperEfficiencyTable(
    sortBy: String
    sortOrder: String
    limit: Int
  ): [DeveloperEfficiencyRow!]!
}

type TeamRadarPoint {
  teamName: String!
  chatUsage: Float!
  codeCompletion: Float!
  refactoringPrompts: Float!
  agentUsage: Float!
  acceptanceRate: Float!
}

type DeveloperEfficiencyRow {
  developer: Developer!
  totalAiLines: Int!
  acceptanceRate: Float!
  isLowPerformer: Boolean!
  trend: String!
}

type Subscription {
  statsUpdated: DashboardKPI!
}
```

#### 4.2.4 KPI Calculation Logic

```typescript
// src/services/kpi-calculator.ts

export interface KPICalculator {
  calculateAcceptanceRate(shown: number, accepted: number): number;
  calculateAIVelocity(acceptedLines: number, totalLines: number): number;
  calculateChatDependency(chatRequests: number, codeEvents: number): number;
  calculateTeamScore(developers: DeveloperStats[]): TeamScore;
}

export const kpiCalculator: KPICalculator = {
  // Acceptance Rate = (Accepted Suggestions / Shown Suggestions) × 100
  calculateAcceptanceRate(shown: number, accepted: number): number {
    if (shown === 0) return 0;
    return (accepted / shown) * 100;
  },

  // AI Velocity = (AI-Accepted Lines / Total Lines Written) × 100
  // Measures what percentage of code output comes from AI assistance
  calculateAIVelocity(acceptedLines: number, totalLines: number): number {
    if (totalLines === 0) return 0;
    return (acceptedLines / totalLines) * 100;
  },

  // Chat Dependency = Chat Requests / (Code Completion Events + 1)
  // Higher values indicate reliance on chat over inline completions
  calculateChatDependency(chatRequests: number, codeEvents: number): number {
    return chatRequests / (codeEvents + 1);
  },

  // Team Score aggregates individual developer metrics
  calculateTeamScore(developers: DeveloperStats[]): TeamScore {
    const active = developers.filter(d => d.isActive);
    return {
      avgAcceptanceRate: mean(active.map(d => d.acceptanceRate)),
      totalAiLines: sum(active.map(d => d.acceptedLinesAdded)),
      activeDevelopers: active.length,
      aiVelocityScore: mean(active.map(d => d.aiVelocity)),
    };
  },
};
```

### 4.3 cursor-viz-spa (React)

#### 4.3.1 Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Framework | React 18 | Component-based UI with hooks |
| Build Tool | Vite 5 | Fast HMR, optimized builds |
| GraphQL Client | Apollo Client 3 | Integrated cache, hooks |
| State | TanStack Query | Server state management |
| Charts | Recharts | Composable React charts |
| Styling | Tailwind CSS | Utility-first, rapid development |
| Testing | Vitest + RTL | Fast, React-focused testing |
| E2E | Playwright | Cross-browser testing |
| Mocking | MSW | Network-level API mocking |

#### 4.3.2 Component Architecture

```
src/
├── components/
│   ├── layout/
│   │   ├── Header.tsx
│   │   ├── Sidebar.tsx
│   │   └── DashboardGrid.tsx
│   ├── charts/
│   │   ├── VelocityHeatmap.tsx
│   │   ├── TeamRadarChart.tsx
│   │   ├── TrendLineChart.tsx
│   │   └── AcceptanceGauge.tsx
│   ├── tables/
│   │   ├── DeveloperEfficiencyTable.tsx
│   │   └── TeamComparisonTable.tsx
│   ├── cards/
│   │   ├── KPICard.tsx
│   │   └── DeveloperCard.tsx
│   └── common/
│       ├── SkeletonLoader.tsx
│       ├── ErrorBoundary.tsx
│       └── Tooltip.tsx
├── pages/
│   ├── Dashboard.tsx
│   ├── TeamView.tsx
│   └── DeveloperView.tsx
├── hooks/
│   ├── useDashboardData.ts
│   ├── useTeamStats.ts
│   └── usePolling.ts
├── graphql/
│   ├── queries.ts
│   ├── mutations.ts
│   └── fragments.ts
└── utils/
    ├── formatters.ts
    ├── colorScales.ts
    └── dateUtils.ts
```

#### 4.3.3 Key Visualizations

**Velocity Heatmap (GitHub-style contribution graph)**

This visualization shows AI code acceptance intensity over time. Each cell represents a day, colored by the volume of accepted AI suggestions. This helps teams identify patterns in AI usage, such as sprint effects or seasonal variations.

```typescript
// components/charts/VelocityHeatmap.tsx
interface HeatmapProps {
  data: HeatmapCell[];
  colorScale?: 'green' | 'blue' | 'purple';
  onCellClick?: (cell: HeatmapCell) => void;
}

// Color intensity maps to acceptedLinesAdded
// 0-25%: lightest, 25-50%: light, 50-75%: medium, 75-100%: dark
```

**Team Comparison Radar Chart**

This radar chart compares teams across multiple dimensions simultaneously: Chat Usage, Code Completion, Refactoring Prompts, Agent Usage, and Acceptance Rate. Teams can quickly identify their strengths and areas for improvement relative to other teams.

```typescript
// components/charts/TeamRadarChart.tsx
interface RadarProps {
  teams: TeamRadarPoint[];
  dimensions: string[];
  highlightTeam?: string;
}
```

**Developer Efficiency Table**

A sortable data table showing individual developer metrics. Rows with acceptance rates below 20% are highlighted in red as potential indicators of low-quality suggestions or suboptimal usage patterns.

```typescript
// components/tables/DeveloperEfficiencyTable.tsx
interface TableProps {
  data: DeveloperEfficiencyRow[];
  lowPerformerThreshold?: number; // Default: 20%
  onRowClick?: (developer: Developer) => void;
}
```

---

## 5. Infrastructure

### 5.1 Docker Compose Configuration

```yaml
# docker-compose.yml
version: '3.8'

services:
  cursor-sim:
    build: ./services/cursor-sim
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - DEVELOPER_COUNT=50
      - VELOCITY=high
      - FLUCTUATION=0.2
      - SEED=42
      - HISTORY_DAYS=30
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: cursor
      POSTGRES_PASSWORD: analytics
      POSTGRES_DB: cursor_analytics
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U cursor"]
      interval: 5s
      timeout: 3s
      retries: 5

  cursor-analytics-core:
    build: ./services/cursor-analytics-core
    ports:
      - "4000:4000"
    environment:
      - PORT=4000
      - DATABASE_URL=postgres://cursor:analytics@db:5432/cursor_analytics
      - SIMULATOR_URL=http://cursor-sim:8080
      - POLL_INTERVAL_MS=60000
    depends_on:
      db:
        condition: service_healthy
      cursor-sim:
        condition: service_healthy

  cursor-viz-spa:
    build: ./services/cursor-viz-spa
    ports:
      - "3000:3000"
    environment:
      - VITE_GRAPHQL_URL=http://localhost:4000/graphql
      - VITE_POLL_INTERVAL=30000
    depends_on:
      - cursor-analytics-core

volumes:
  postgres_data:
```

### 5.2 Development Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                    Development Commands                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  make dev          Start all services with hot-reload            │
│  make test         Run all tests across services                 │
│  make test-sim     Run cursor-sim tests only                     │
│  make test-core    Run cursor-analytics-core tests only          │
│  make test-spa     Run cursor-viz-spa tests only                 │
│  make lint         Run linters across all services               │
│  make build        Build production Docker images                │
│  make clean        Stop containers and clean volumes             │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 6. Testing Strategy

### 6.1 cursor-sim (Go)

**Unit Tests** cover the event generation logic, Poisson distribution implementation, developer profile generation, and data model serialization. Tests use table-driven patterns with deterministic seeding for reproducibility.

**Integration Tests** verify HTTP handlers return correct responses, query parameters filter data appropriately, and pagination works correctly.

```go
// internal/generator/poisson_test.go
func TestPoissonDistribution(t *testing.T) {
    tests := []struct {
        name     string
        lambda   float64
        samples  int
        wantMean float64
        tolerance float64
    }{
        {"low velocity", 10.0, 10000, 10.0, 0.5},
        {"high velocity", 100.0, 10000, 100.0, 2.0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gen := NewPoissonGenerator(tt.lambda, 42)
            values := make([]int, tt.samples)
            for i := range values {
                values[i] = gen.Next()
            }
            gotMean := mean(values)
            assert.InDelta(t, tt.wantMean, gotMean, tt.tolerance)
        })
    }
}
```

### 6.2 cursor-analytics-core (TypeScript)

**Unit Tests** cover KPI calculation functions, data transformation logic, and utility functions. These tests are fast and isolated.

**Integration Tests** use pg-mem for an in-memory PostgreSQL database and Apollo's test client for GraphQL resolver testing. These verify the full request-response cycle without external dependencies.

```typescript
// src/services/__tests__/kpi-calculator.test.ts
describe('KPICalculator', () => {
  describe('calculateAcceptanceRate', () => {
    it('returns percentage of accepted suggestions', () => {
      expect(kpiCalculator.calculateAcceptanceRate(100, 75)).toBe(75);
    });

    it('returns 0 when no suggestions shown', () => {
      expect(kpiCalculator.calculateAcceptanceRate(0, 0)).toBe(0);
    });

    it('handles edge case of more accepted than shown', () => {
      // This shouldn't happen but should not throw
      expect(kpiCalculator.calculateAcceptanceRate(50, 60)).toBe(120);
    });
  });
});
```

### 6.3 cursor-viz-spa (React)

**Component Tests** use React Testing Library to verify components render correctly, handle loading and error states, respond to user interactions, and display data accurately.

**Integration Tests** use MSW (Mock Service Worker) to intercept GraphQL requests, allowing tests to verify the full data flow from query to render.

**E2E Tests** use Playwright to verify critical user journeys like viewing the dashboard, filtering by team, and drilling into developer details.

```typescript
// src/components/charts/__tests__/VelocityHeatmap.test.tsx
describe('VelocityHeatmap', () => {
  it('renders cells for each day in the data', () => {
    const data = generateMockHeatmapData(30);
    render(<VelocityHeatmap data={data} />);
    expect(screen.getAllByRole('cell')).toHaveLength(30);
  });

  it('applies correct color intensity based on value', () => {
    const data = [
      { date: '2025-01-01', value: 10, intensity: 0.25 },
      { date: '2025-01-02', value: 50, intensity: 0.75 },
    ];
    render(<VelocityHeatmap data={data} />);
    const cells = screen.getAllByRole('cell');
    expect(cells[0]).toHaveClass('intensity-low');
    expect(cells[1]).toHaveClass('intensity-high');
  });
});
```

---

## 7. Security Considerations

While this is a local development tool without production security requirements, several practices are followed for good hygiene:

**Input Validation**: All query parameters are validated and sanitized to prevent injection attacks.

**Error Handling**: Errors are logged internally but not exposed in detail to API responses.

**Dependencies**: Regular `npm audit` and `go mod verify` checks are part of the CI pipeline.

**No Secrets**: No real API keys or credentials are used; all data is synthetic.

---

## 8. Future Considerations

These items are out of scope for the current version but could be added later:

**Real Cursor API Integration**: Replace the simulator with actual Cursor Business API calls once access is available.

**Persistent Storage**: Add S3 or local file export for historical analysis beyond container lifecycle.

**Authentication**: Add basic auth or OAuth if the tool is deployed for team access.

**Alerting**: Add threshold-based alerts when acceptance rates drop below configurable levels.

**CI/CD Pipeline**: Add GitHub Actions for automated testing and Docker Hub publishing.

---

## 9. Appendix

### 9.1 Glossary

| Term | Definition |
|------|------------|
| **Acceptance Rate** | Percentage of shown AI suggestions that were accepted by the developer |
| **AI Velocity** | Ratio of AI-generated code lines to total code lines written |
| **Chat Dependency** | Measure of how much a developer relies on chat vs. inline completions |
| **cpp_suggestion** | "Cursor Plus Plus" - Cursor's inline code completion feature |
| **cmd_k** | Keyboard shortcut for Cursor's inline edit feature |
| **Tab Completion** | Accepting an AI suggestion by pressing Tab |

### 9.2 Reference Documents

- `/docs/features/` - Feature specifications
- `/docs/user-stories/` - User story definitions
- `/docs/tasks/` - Implementation task breakdowns
- `/docs/api-reference/` - API documentation

### 9.3 Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2025-01-01 | Use Go for simulator | High-performance concurrent simulation |
| 2025-01-01 | Use TypeScript for aggregator | Type safety with GraphQL codegen |
| 2025-01-01 | Use Recharts over D3 | Easier integration with React components |
| 2025-01-01 | In-memory SQLite for simulator | Volatile storage matches requirements |
| 2025-01-01 | PostgreSQL for aggregator | Robust for relational queries and JSONB |
