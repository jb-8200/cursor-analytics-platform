# Design: cursor-analytics-core

**Feature**: cursor-analytics-core (GraphQL Aggregator)
**Status**: NOT_STARTED
**Estimated Hours**: 25-30

---

## Overview

cursor-analytics-core is a TypeScript service that acts as the data aggregation layer between cursor-sim and cursor-viz-spa. It polls cursor-sim for commit data, stores it in PostgreSQL, calculates metrics, and exposes a GraphQL API for the frontend.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        cursor-analytics-core                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────────────┐    ┌──────────────────┐    ┌─────────────────┐   │
│  │ Ingestion Worker │───▶│ PostgreSQL       │◄───│ GraphQL Server  │   │
│  │   (Polls sim)    │    │   (Data store)   │    │   (Apollo)      │   │
│  └──────────────────┘    └──────────────────┘    └─────────────────┘   │
│         │                        ▲                       │              │
│         │                        │                       │              │
│         ▼                        │                       ▼              │
│  ┌──────────────────┐    ┌──────────────────┐    ┌─────────────────┐   │
│  │ cursor-sim       │    │ Metric Calculator│    │ cursor-viz-spa  │   │
│  │ :8080            │    │ (Derived KPIs)   │    │ :3000           │   │
│  └──────────────────┘    └──────────────────┘    └─────────────────┘   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Technology Stack

- **Runtime**: Node.js 20 LTS
- **Language**: TypeScript 5.3+
- **GraphQL**: Apollo Server 4.x
- **Database**: PostgreSQL 15 with Prisma ORM
- **Testing**: Jest + ts-jest + supertest
- **Linting**: ESLint + Prettier

## Package Structure

```
services/cursor-analytics-core/
├── src/
│   ├── index.ts              # Entry point
│   ├── config/
│   │   └── index.ts          # Environment config
│   ├── db/
│   │   ├── client.ts         # Prisma client
│   │   └── migrations/       # Database migrations
│   ├── ingestion/
│   │   ├── worker.ts         # Polling worker
│   │   ├── client.ts         # cursor-sim REST client
│   │   └── transformer.ts    # Data normalization
│   ├── graphql/
│   │   ├── schema.ts         # Type definitions
│   │   ├── resolvers/
│   │   │   ├── developer.ts
│   │   │   ├── commit.ts
│   │   │   └── dashboard.ts
│   │   └── context.ts        # Request context
│   ├── services/
│   │   ├── metrics.ts        # Metric calculations
│   │   └── developer.ts      # Developer queries
│   └── types/
│       └── index.ts          # TypeScript types
├── prisma/
│   └── schema.prisma         # Database schema
├── tests/
│   ├── unit/
│   ├── integration/
│   └── fixtures/
├── package.json
├── tsconfig.json
└── SPEC.md
```

## Database Schema

### PostgreSQL Tables

```sql
-- Developers table (denormalized from cursor-sim)
CREATE TABLE developers (
  id VARCHAR(255) PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  org VARCHAR(255),
  division VARCHAR(255),
  team VARCHAR(255),
  role VARCHAR(50),
  region VARCHAR(50),
  timezone VARCHAR(50),
  seniority VARCHAR(50),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Commits table
CREATE TABLE commits (
  commit_hash VARCHAR(64) PRIMARY KEY,
  user_id VARCHAR(255) REFERENCES developers(id),
  repo_name VARCHAR(255) NOT NULL,
  branch_name VARCHAR(255) NOT NULL,
  is_primary_branch BOOLEAN DEFAULT false,
  total_lines_added INT NOT NULL,
  total_lines_deleted INT NOT NULL,
  tab_lines_added INT NOT NULL,
  tab_lines_deleted INT NOT NULL,
  composer_lines_added INT NOT NULL,
  composer_lines_deleted INT NOT NULL,
  non_ai_lines_added INT NOT NULL,
  non_ai_lines_deleted INT NOT NULL,
  message TEXT,
  commit_ts TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),

  -- Indexes
  INDEX idx_commits_user (user_id),
  INDEX idx_commits_repo (repo_name),
  INDEX idx_commits_ts (commit_ts)
);

-- Sync state table
CREATE TABLE sync_state (
  id SERIAL PRIMARY KEY,
  last_sync_ts TIMESTAMP NOT NULL,
  records_synced INT NOT NULL,
  status VARCHAR(50) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);
```

## GraphQL Schema

```graphql
type Query {
  # Developer queries
  developer(id: ID!): Developer
  developers(
    team: String
    division: String
    limit: Int = 50
    offset: Int = 0
  ): [Developer!]!

  # Commit queries
  commits(
    userId: String
    repoName: String
    from: DateTime
    to: DateTime
    limit: Int = 100
    offset: Int = 0
  ): CommitConnection!

  # Dashboard summary
  dashboardSummary(dateRange: DateRangeInput): DashboardSummary!

  # Team analytics
  teamStats(team: String, dateRange: DateRangeInput): TeamStats!
}

type Developer {
  id: ID!
  email: String!
  name: String!
  team: String
  division: String
  seniority: String
  stats(dateRange: DateRangeInput): DeveloperStats!
}

type DeveloperStats {
  totalCommits: Int!
  totalLinesAdded: Int!
  totalLinesDeleted: Int!
  tabLinesAdded: Int!
  composerLinesAdded: Int!
  aiVelocity: Float!
  acceptanceRate: Float
}

type DashboardSummary {
  totalDevelopers: Int!
  activeDevelopers: Int!
  totalCommits: Int!
  totalLinesAdded: Int!
  overallAIVelocity: Float!
  overallAcceptanceRate: Float
  topContributors: [Developer!]!
}

type TeamStats {
  team: String!
  memberCount: Int!
  totalCommits: Int!
  aiVelocity: Float!
  topRepositories: [RepoStats!]!
}
```

## Ingestion Worker

### Polling Logic

```typescript
class IngestionWorker {
  private pollIntervalMs: number;
  private cursorSimClient: CursorSimClient;
  private db: PrismaClient;

  async start(): Promise<void> {
    setInterval(() => this.poll(), this.pollIntervalMs);
  }

  async poll(): Promise<void> {
    const lastSync = await this.getLastSyncTimestamp();
    const commits = await this.cursorSimClient.getCommits({
      from: lastSync,
      to: new Date(),
    });

    await this.db.$transaction(async (tx) => {
      // Upsert developers
      for (const commit of commits) {
        await tx.developer.upsert({...});
        await tx.commit.upsert({...});
      }

      // Update sync state
      await tx.syncState.create({
        data: {
          lastSyncTs: new Date(),
          recordsSynced: commits.length,
          status: 'success',
        },
      });
    });
  }
}
```

## Metric Calculations

```typescript
class MetricsService {
  // AI Velocity = (tabLines + composerLines) / totalLines * 100
  calculateAIVelocity(commits: Commit[]): number {
    const totalAdded = sum(commits.map(c => c.totalLinesAdded));
    const aiAdded = sum(commits.map(c => c.tabLinesAdded + c.composerLinesAdded));
    return totalAdded > 0 ? (aiAdded / totalAdded) * 100 : 0;
  }

  // Aggregate stats for a developer
  async getDeveloperStats(developerId: string, dateRange?: DateRange): Promise<DeveloperStats> {
    const commits = await this.db.commit.findMany({
      where: {
        userId: developerId,
        commitTs: dateRange ? {
          gte: dateRange.from,
          lte: dateRange.to,
        } : undefined,
      },
    });

    return {
      totalCommits: commits.length,
      totalLinesAdded: sum(commits.map(c => c.totalLinesAdded)),
      totalLinesDeleted: sum(commits.map(c => c.totalLinesDeleted)),
      tabLinesAdded: sum(commits.map(c => c.tabLinesAdded)),
      composerLinesAdded: sum(commits.map(c => c.composerLinesAdded)),
      aiVelocity: this.calculateAIVelocity(commits),
      acceptanceRate: null, // Calculated from events, not commits
    };
  }
}
```

## Environment Variables

```bash
# Database
DATABASE_URL=postgresql://cursor:cursor_dev@localhost:5432/cursor_analytics

# cursor-sim connection
SIMULATOR_URL=http://localhost:8080
SIMULATOR_API_KEY=cursor-sim-dev-key

# Polling
POLL_INTERVAL_MS=60000

# Server
GRAPHQL_PORT=4000
```

## Testing Strategy

### Unit Tests
- Metric calculations
- Data transformation
- GraphQL resolvers (mocked DB)

### Integration Tests
- Full GraphQL queries with test database
- Ingestion worker with mocked cursor-sim

### E2E Tests
- Full pipeline: cursor-sim → ingestion → GraphQL
- Performance under load

## Decision Log

| Decision | Rationale |
|----------|-----------|
| Prisma ORM | Type-safe database access, migrations |
| Apollo Server 4 | Modern GraphQL implementation |
| Polling vs WebSocket | Simpler implementation, adequate for analytics |
| Denormalized developers | Avoid joins, faster queries |
