# cursor-analytics-core

GraphQL API service that transforms cursor-sim data into aggregated analytics for the visualization dashboard.

## Overview

cursor-analytics-core (P5) is the middle tier of the Cursor Analytics Platform, providing a GraphQL API that:
- Fetches historical data from cursor-sim REST API (P4)
- Aggregates and transforms data into analytics
- Exposes GraphQL schema for cursor-viz-spa (P6)
- Persists processed data in PostgreSQL

## Tech Stack

- **Framework**: Node.js 18+ with TypeScript
- **GraphQL**: Apollo Server 4
- **Database**: PostgreSQL 14+ with Prisma ORM
- **Testing**: Jest
- **Deployment**: Docker Compose (GraphQL + PostgreSQL)
- **Port**: 4000 (GraphQL), 5432 (PostgreSQL)

## Architecture

```
┌─────────────────────────────────────────┐
│  cursor-analytics-core (P5)             │
│  ┌───────────────────────────────────┐  │
│  │  Apollo Server 4                  │  │
│  │  GraphQL API (Port 4000)          │  │
│  │  ┌─────────────────────────────┐  │  │
│  │  │ Resolvers                   │  │  │
│  │  │ - dashboardSummary          │  │  │
│  │  │ - teamComparison            │  │  │
│  │  │ - developerPerformance      │  │  │
│  │  └─────────────────────────────┘  │  │
│  └───────────────────────────────────┘  │
│                 │                       │
│                 ▼                       │
│  ┌───────────────────────────────────┐  │
│  │  Prisma ORM                       │  │
│  │  Database Access Layer            │  │
│  └───────────────────────────────────┘  │
│                 │                       │
│                 ▼                       │
│  ┌───────────────────────────────────┐  │
│  │  PostgreSQL (Port 5432)           │  │
│  │  Persistent Storage               │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
                 ▲
                 │ REST API
                 │
        ┌────────┴────────┐
        │  cursor-sim (P4) │
        │  Port 8080       │
        └──────────────────┘
```

## Getting Started

### Prerequisites

- Node.js 18+
- Docker & Docker Compose
- cursor-sim running on port 8080

### Installation

```bash
cd services/cursor-analytics-core
npm install
```

### Database Setup

```bash
# Generate Prisma Client
npm run db:generate

# Run migrations
npm run db:migrate:deploy

# (Optional) Seed database
npm run db:seed
```

### Development

#### Option 1: Docker Compose (Recommended)

```bash
# Start GraphQL + PostgreSQL in Docker
docker-compose up

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

The GraphQL server will be available at `http://localhost:4000/graphql`.

#### Option 2: Local Development

```bash
# Start PostgreSQL in Docker
docker-compose up postgres

# Start GraphQL server locally
npm run dev

# In another terminal, run tests
npm test

# Run tests with coverage
npm run test:coverage

# Type check
npm run type-check

# Lint code
npm run lint
```

### Build

```bash
npm run build
npm start
```

## GraphQL Schema

### Queries

```graphql
type Query {
  # Dashboard summary with KPIs and trends
  dashboardSummary(range: DateRangeInput): DashboardSummary!

  # Team comparison metrics
  teamComparison(range: DateRangeInput): [TeamStats!]!

  # Individual developer performance
  developerPerformance(developerId: ID!, range: DateRangeInput): DeveloperStats
}
```

### Key Types

```graphql
type DashboardSummary {
  totalDevelopers: Int!
  activeDevelopers: Int!
  overallAcceptanceRate: Float
  aiVelocityToday: Int!
  dailyTrend: [DailyStats!]!
  teamComparison: [TeamStats!]!
}

type TeamStats {
  teamName: String!
  aiAcceptanceRate: Float
  topPerformer: Developer  # Note: Singular, optional
}

type DailyStats {
  date: String!
  linesAdded: Int!  # Note: Not 'humanLinesAdded'
  aiLinesAdded: Int!
}
```

**IMPORTANT**: The actual schema is the source of truth. See `src/graphql/schema.ts` for complete schema definition.

## GraphQL Playground

When running in development mode, access the GraphQL Playground at:

```
http://localhost:4000/graphql
```

Example query:
```graphql
query GetDashboard {
  dashboardSummary(range: { preset: LAST_30_DAYS }) {
    totalDevelopers
    activeDevelopers
    overallAcceptanceRate
    aiVelocityToday
    dailyTrend {
      date
      linesAdded
      aiLinesAdded
    }
    teamComparison {
      teamName
      aiAcceptanceRate
      topPerformer {
        id
        name
        email
        seniority
      }
    }
  }
}
```

## Testing

### Test Coverage (as of January 4, 2026)

- **Unit tests**: Resolver and service testing
- **Integration tests**: Full GraphQL query testing
- **Coverage target**: 80%

### Run Tests

```bash
# All tests
npm test

# Watch mode
npm run test:watch

# Coverage report
npm run test:coverage

# Specific test file
npm test -- src/graphql/resolvers.test.ts
```

## Configuration

Create a `.env` file (see `.env.example`):

```bash
# Database
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/cursor_analytics

# cursor-sim API
CURSOR_SIM_API_URL=http://localhost:8080

# Server
PORT=4000
NODE_ENV=development

# GraphQL
GRAPHQL_INTROSPECTION=true
GRAPHQL_PLAYGROUND=true
```

## Project Structure

```
src/
├── graphql/
│   ├── schema.ts         # GraphQL type definitions (SOURCE OF TRUTH)
│   ├── resolvers.ts      # Query/mutation resolvers
│   └── context.ts        # GraphQL context (Prisma, dataloaders)
├── services/
│   ├── dashboard.ts      # Dashboard aggregation logic
│   ├── team.ts           # Team analytics
│   └── developer.ts      # Developer analytics
├── datasources/
│   └── cursor-sim.ts     # REST API client for cursor-sim
├── prisma/
│   ├── schema.prisma     # Database schema
│   └── migrations/       # Database migrations
└── __tests__/            # Test files
```

## Platform Integration

cursor-analytics-core (P5) is the middle tier of the **GraphQL Analytics Path** in the Cursor Analytics Platform:

### GraphQL Analytics Path (Original)

```
cursor-sim (P4) → cursor-analytics-core (P5) → cursor-viz-spa (P6)
  Docker          Docker Compose             Local npm dev
  Port 8080       Port 4000 (GraphQL)        Port 3000
                  Port 5432 (PostgreSQL)
```

### Alternative: dbt + Streamlit Analytics Path (New - January 2026)

The platform now includes an alternative analytics path using dbt + Streamlit:

```
cursor-sim (P4) → streamlit-dashboard (P9)
  Docker          Docker (includes dbt P8 + DuckDB)
  Port 8080       Port 8501
```

**Key Differences:**

| Feature | GraphQL Path (P5+P6) | Streamlit Path (P9) |
|---------|---------------------|---------------------|
| **Backend** | PostgreSQL + Prisma | DuckDB (dev) or Snowflake (prod) |
| **Transform** | TypeScript services | dbt SQL models |
| **API** | GraphQL | Direct database queries |
| **Frontend** | React + Apollo Client | Streamlit Python |
| **Use Case** | Interactive dashboards, complex queries | Data analysis, reporting, embedded ETL |

**Both paths are maintained** - choose based on your use case:
- **GraphQL path**: Interactive web dashboards, mobile apps, external API consumers
- **Streamlit path**: Internal analytics, data exploration, ad-hoc reporting

### Integration Status (January 4, 2026)

✅ **P4+P5 Integration**: Complete
- Fetches data from cursor-sim REST API
- Transforms and aggregates into GraphQL

✅ **P5+P6 Integration**: Complete
- Serves GraphQL schema to cursor-viz-spa
- All schema mismatches resolved (see MITIGATION-PLAN.md)

### Data Contract Testing

**CRITICAL**: The GraphQL schema in `src/graphql/schema.ts` is the source of truth for P6.

To prevent schema drift:
1. Always update `schema.ts` first
2. Run tests to validate schema changes
3. P6 should run `npm run codegen` to regenerate types
4. Never manually update P6 types without schema changes

See [Data Contract Testing Strategy](../../docs/data-contract-testing.md) for the comprehensive mitigation plan.

## Common Tasks

### Update GraphQL Schema

```bash
# 1. Edit src/graphql/schema.ts
# 2. Update resolvers if needed
# 3. Run tests
npm test

# 4. Notify P6 team to run codegen
# (In P6) npm run codegen
```

### Database Migrations

```bash
# Create migration
npx prisma migrate dev --name add_new_field

# Apply migrations (production)
npm run db:migrate:deploy

# Reset database (development only)
npx prisma migrate reset
```

### Debugging

```bash
# Enable debug logging
DEBUG=* npm run dev

# Check database connection
npx prisma db pull

# View database in Prisma Studio
npx prisma studio
```

## Docker Deployment

### Build Image

```bash
docker build -t cursor-analytics-core:latest .
```

### Run with Docker Compose

```bash
# Development
docker-compose up

# Production
docker-compose -f docker-compose.prod.yml up -d
```

### Environment Variables (Docker)

```yaml
# docker-compose.yml
environment:
  DATABASE_URL: postgresql://postgres:postgres@postgres:5432/cursor_analytics
  CURSOR_SIM_API_URL: http://cursor-sim:8080
  PORT: 4000
  NODE_ENV: production
```

## Troubleshooting

### GraphQL 400 Bad Request Errors

**Symptom**: P6 shows "Cannot query field X on type Y"

**Root Cause**: Schema mismatch between P5 and P6

**Solution**:
1. Check `src/graphql/schema.ts` for actual field names
2. Update P6 queries to match P5 schema
3. In P6: `npm run codegen` to regenerate types
4. See [docs/MITIGATION-PLAN.md](../../docs/MITIGATION-PLAN.md) for prevention

### Database Connection Errors

**Symptom**: "Can't reach database server"

**Solutions**:
```bash
# Check PostgreSQL is running
docker-compose ps

# Check connection string
echo $DATABASE_URL

# Restart PostgreSQL
docker-compose restart postgres

# Check logs
docker-compose logs postgres
```

### Prisma Client Out of Sync

**Symptom**: "Prisma Client does not match schema"

**Solution**:
```bash
# Regenerate Prisma Client
npm run db:generate

# If migrations are pending
npm run db:migrate:deploy
```

## Performance

### Response Time Targets

| Query | Target | Notes |
|-------|--------|-------|
| `dashboardSummary` | < 500ms | Cached for 1 minute |
| `teamComparison` | < 300ms | Indexed by teamName |
| `developerPerformance` | < 200ms | Indexed by developerId |

### Optimization Strategies

- **Caching**: Apollo Server response cache
- **Batching**: DataLoader for N+1 query prevention
- **Indexing**: Database indexes on common queries
- **Pagination**: Cursor-based pagination for large datasets

## Documentation

### Service Documentation

- **Specification**: `SPEC.md` (to be created)
- **User Story**: `.work-items/P5-cursor-analytics-core/user-story.md`
- **Design**: `.work-items/P5-cursor-analytics-core/design.md`
- **Tasks**: `.work-items/P5-cursor-analytics-core/task.md`

### Platform Documentation

- **Architecture**: [docs/DESIGN.md](../../docs/DESIGN.md)
- **Integration Guide**: [docs/INTEGRATION.md](../../docs/INTEGRATION.md)
- **Data Contract Testing**: [docs/data-contract-testing.md](../../docs/data-contract-testing.md)
- **E2E Testing Strategy**: [docs/e2e-testing-strategy.md](../../docs/e2e-testing-strategy.md)
- **Mitigation Plan**: [docs/MITIGATION-PLAN.md](../../docs/MITIGATION-PLAN.md)

## Support

- **Issues**: Report bugs in project issue tracker
- **Schema Changes**: Follow breaking change policy in MITIGATION-PLAN.md
- **Questions**: See platform documentation in `docs/`

## License

Internal development tool - not for external distribution.
