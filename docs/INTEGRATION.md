# P5/P6 Integration Testing Guide

Complete guide for running the cursor-analytics-platform stack locally with full integration between cursor-sim (P4), cursor-analytics-core (P5), and cursor-viz-spa (P6).

## Architecture Overview

```
┌─────────────────┐      ┌──────────────────────┐      ┌─────────────────────┐
│   cursor-sim    │      │ cursor-analytics-core│      │  cursor-viz-spa     │
│   (P4 - Go)     │─────▶│   (P5 - GraphQL)     │─────▶│  (P6 - React/Vite)  │
│   Port 8080     │ REST │   Port 4000          │ GQL  │   Port 3000         │
└─────────────────┘      └──────────────────────┘      └─────────────────────┘
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │    PostgreSQL       │
                         │    Port 5432        │
                         └─────────────────────┘
```

## Quick Start (All Services)

### Prerequisites

- Docker installed and running
- Node.js 18+ installed
- npm installed
- Go 1.22+ installed (optional, only if running cursor-sim locally without Docker)

### 1. Start All Services

```bash
# From project root
./tools/start-integration.sh
```

Or manually:

```bash
# Terminal 1: Start cursor-sim (Docker)
DETACH=true ./tools/docker-local.sh

# Terminal 2: Start PostgreSQL + cursor-analytics-core
cd services/cursor-analytics-core
docker-compose up -d
npm run db:generate
npm run dev

# Terminal 3: Start cursor-viz-spa
cd services/cursor-viz-spa
npm run dev
```

### 2. Verify Services

```bash
# Check cursor-sim
curl http://localhost:8080/health

# Check PostgreSQL
docker exec cursor-analytics-postgres psql -U cursor -d cursor_analytics -c '\l'

# Check P5 GraphQL (requires headers for CSRF protection)
curl -H "Content-Type: application/json" \
  -d '{"query":"{ __typename }"}' \
  http://localhost:4000/graphql

# Open P6 Dashboard in browser
open http://localhost:3000
```

## Service Details

### cursor-sim (P4 - Data Source)

**Purpose**: Simulates Cursor IDE usage data (commits, suggestions, chat interactions)

**Tech Stack**: Go 1.22+
**Port**: 8080
**Mode**: Docker container (recommended) or local binary

**Configuration**:
- DAYS: 90 (default) - Days of historical data
- VELOCITY: medium (default) - Commit frequency (low/medium/high)
- MODE: runtime (default) - Simulation mode

**Endpoints**:
- `GET /health` - Health check
- `GET /teams/members` - List all team members
- `GET /commits` - List commit history
- `GET /stats/daily` - Daily statistics

**Authentication**: Basic Auth with API key `cursor-sim-dev-key`

**Start**:
```bash
# Docker (recommended)
DETACH=true ./tools/docker-local.sh

# Local binary
cd services/cursor-sim
go build -o bin/cursor-sim ./cmd/simulator
./bin/cursor-sim -mode runtime -seed testdata/valid_seed.json -port 8080
```

**Stop**:
```bash
docker stop cursor-sim-local
```

---

### cursor-analytics-core (P5 - GraphQL API)

**Purpose**: GraphQL API for querying analytics data, aggregating metrics

**Tech Stack**: TypeScript, Apollo Server 4, Prisma, PostgreSQL
**Port**: 4000
**Database**: PostgreSQL 15 (Docker)

**Configuration**:
- DATABASE_URL: `postgresql://cursor:cursor_dev@localhost:5432/cursor_analytics`
- NODE_ENV: development

**GraphQL Schema**:
- `health` - System health check
- `developer(id)` - Get developer by ID
- `developers(team, seniority, limit, offset)` - List developers with filtering
- `teamStats(team, dateRange)` - Team-level statistics
- `teams` - List all teams with stats
- `dashboardSummary(dateRange)` - Dashboard KPIs and trends
- `commits(userId, team, dateRange, limit, cursor)` - Paginated commit history

**Start**:
```bash
cd services/cursor-analytics-core

# 1. Start PostgreSQL
docker-compose up -d

# 2. Generate Prisma client
npm run db:generate

# 3. Apply migrations (if needed)
# Migrations are applied automatically via docker exec

# 4. Seed database (optional - manual data insertion)
# See "Database Seeding" section below

# 5. Start GraphQL server
npm run dev
```

**GraphQL Playground**: http://localhost:4000/

**Stop**:
```bash
docker-compose down
pkill -f "tsx watch src/index.ts"
```

---

### cursor-viz-spa (P6 - Dashboard)

**Purpose**: React dashboard for visualizing analytics data

**Tech Stack**: React 18, Vite, Apollo Client, Tailwind CSS, Recharts
**Port**: 3000
**Connects to**: P5 GraphQL API (http://localhost:4000/graphql)

**Configuration** (`.env`):
```
VITE_GRAPHQL_URL=http://localhost:4000/graphql
VITE_APP_TITLE=Cursor Analytics Dashboard
VITE_DEFAULT_DATE_RANGE=LAST_30_DAYS
VITE_ENABLE_DEBUG=true
```

**Features**:
- Dashboard page with KPIs, charts, and trends
- Teams page with team comparisons
- Developers page with individual stats
- Interactive filters (date range, team, seniority)
- Responsive design (mobile/tablet/desktop)

**Start**:
```bash
cd services/cursor-viz-spa
npm run dev
```

**Access**: http://localhost:3000

**Stop**:
```bash
pkill -f "vite"
```

---

## Database Setup

### PostgreSQL via Docker

The `docker-compose.yml` in `services/cursor-analytics-core/` provides a pre-configured PostgreSQL instance.

```bash
cd services/cursor-analytics-core
docker-compose up -d
```

**Connection Details**:
- Host: localhost
- Port: 5432
- Database: cursor_analytics
- User: cursor
- Password: cursor_dev

### Schema Migrations

Migrations are located in `services/cursor-analytics-core/prisma/migrations/`.

**Apply migrations**:
```bash
# Via Prisma
npm run db:migrate:deploy

# Or manually via Docker
cat prisma/migrations/20260103_init/migration.sql | \
  docker exec -i cursor-analytics-postgres psql -U cursor -d cursor_analytics
```

### Database Seeding

**Option 1: Manual SQL Insertion** (Current approach for integration testing)

```bash
docker exec -i cursor-analytics-postgres psql -U cursor -d cursor_analytics << 'EOF'
INSERT INTO developers (id, external_id, name, email, team, seniority)
VALUES
  (gen_random_uuid(), 'alice', 'Alice Johnson', 'alice@example.com', 'Backend', 'senior'),
  (gen_random_uuid(), 'bob', 'Bob Smith', 'bob@example.com', 'Frontend', 'mid'),
  (gen_random_uuid(), 'carol', 'Carol Davis', 'carol@example.com', 'Backend', 'junior');
EOF
```

**Option 2: Seed Script** (Future - when Prisma connection issues resolved)

```bash
npm run seed
```

**Option 3: Populate from cursor-sim** (Future - requires Step 04 Ingestion Worker)

Automatically sync data from cursor-sim REST API into PostgreSQL.

---

## Integration Testing Workflow

### 1. Verify cursor-sim Data

```bash
# Get team members
curl -u cursor-sim-dev-key: http://localhost:8080/teams/members | jq

# Get commits
curl -u cursor-sim-dev-key: http://localhost:8080/commits?limit=5 | jq
```

### 2. Verify P5 GraphQL Queries

Open http://localhost:4000/ in your browser (GraphQL Playground).

**Example queries**:

```graphql
# Health check
query {
  health {
    status
    database {
      connected
      developerCount
      eventCount
    }
  }
}

# List developers
query {
  developers(limit: 10) {
    nodes {
      id
      name
      email
      team
      seniority
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}

# Get team stats
query {
  teamStats(team: "Backend") {
    team
    memberCount
    activeMemberCount
    acceptanceRate
    aiVelocity
  }
}

# Dashboard summary
query {
  dashboardSummary(dateRange: LAST_30_DAYS) {
    totalDevelopers
    activeDevelopers
    overallAcceptanceRate
    todayStats {
      suggestionsCount
      acceptedCount
    }
  }
}
```

### 3. Verify P6 Dashboard

Open http://localhost:3000 in your browser.

**Manual Testing Checklist**:

- [ ] Dashboard page loads without errors
- [ ] KPI cards display correct data (Total Devs, Acceptance Rate, AI Velocity)
- [ ] Velocity Heatmap renders with data
- [ ] Team Radar Chart shows team comparisons
- [ ] Developer Table displays with sorting and pagination
- [ ] Date range picker works (7d, 30d, 90d, custom)
- [ ] Search input filters developers
- [ ] Navigate to Teams page - team list and stats load
- [ ] Navigate to Developers page - developer profiles load
- [ ] Mobile responsive layout works (resize browser)
- [ ] No console errors in browser DevTools

### 4. End-to-End Data Flow Test

**Scenario**: Verify data flows from cursor-sim → P5 → P6

1. **cursor-sim** generates simulated data (90 days, medium velocity)
2. **P5 database** contains seeded developer and event data
3. **P5 GraphQL** exposes data via queries
4. **P6 Dashboard** fetches and displays data via Apollo Client

**Verification**:
```bash
# 1. Count events in cursor-sim
curl -u cursor-sim-dev-key: http://localhost:8080/commits | jq '. | length'

# 2. Count events in P5 database
docker exec cursor-analytics-postgres psql -U cursor -d cursor_analytics \
  -c 'SELECT COUNT(*) FROM usage_events;'

# 3. Query via P5 GraphQL
curl -H "Content-Type: application/json" \
  -d '{"query":"{ developers { totalCount } }"}' \
  http://localhost:4000/graphql | jq

# 4. Open P6 dashboard and verify developer count matches
open http://localhost:3000
```

---

## Troubleshooting

### cursor-sim (P4) Issues

**Port 8080 already in use**:
```bash
PORT=8081 ./tools/docker-local.sh
```

**Container won't start**:
```bash
docker logs cursor-sim-local
docker stop cursor-sim-local
docker rm cursor-sim-local
REBUILD=true ./tools/docker-local.sh
```

**No data returned from API**:
- Check health endpoint: `curl http://localhost:8080/health`
- Verify seed file loaded: Check container logs
- Try with different velocity: `VELOCITY=high ./tools/docker-local.sh`

---

### PostgreSQL (P5 Database) Issues

**Port 5432 already in use**:
```bash
# Find process using port
lsof -i :5432

# Stop conflicting PostgreSQL
brew services stop postgresql  # macOS
sudo systemctl stop postgresql  # Linux
```

**Database connection failed**:
```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check logs
docker logs cursor-analytics-postgres

# Restart
docker-compose restart
```

**Cannot connect from host machine**:
```bash
# Test connection
docker exec cursor-analytics-postgres psql -U cursor -d cursor_analytics -c 'SELECT version();'

# If works inside container but not from host, use docker exec for operations
```

**Database not found**:
```bash
# Recreate database
docker-compose down -v
docker-compose up -d
sleep 5
cat services/cursor-analytics-core/prisma/migrations/20260103_init/migration.sql | \
  docker exec -i cursor-analytics-postgres psql -U cursor -d cursor_analytics
```

---

### cursor-analytics-core (P5) Issues

**Module not found: './generated/prisma'**:
```bash
npm run db:generate
npm run build  # or use npm run dev for development
```

**CSRF errors on GraphQL queries**:
Add `Content-Type: application/json` header to all requests:
```bash
curl -H "Content-Type: application/json" \
  -d '{"query":"{ __typename }"}' \
  http://localhost:4000/graphql
```

**GraphQL server won't start**:
```bash
# Check logs
cat /tmp/claude/-Users-jbellish-VSProjects-cursor-analytics-platform/tasks/*.output

# Kill and restart
pkill -f "tsx watch"
npm run dev
```

**No data returned from queries**:
```bash
# Check database has data
docker exec cursor-analytics-postgres psql -U cursor -d cursor_analytics \
  -c 'SELECT COUNT(*) FROM developers; SELECT COUNT(*) FROM usage_events;'

# If empty, re-seed (see "Database Seeding" section)
```

---

### cursor-viz-spa (P6) Issues

**Port 3000 already in use**:
```bash
# Kill process on port 3000
lsof -ti:3000 | xargs kill -9

# Or use different port
PORT=3001 npm run dev
```

**Dashboard shows loading indefinitely**:
1. Open browser DevTools (F12)
2. Check Console for errors
3. Check Network tab for failed GraphQL requests
4. Common causes:
   - P5 GraphQL server not running → Start with `npm run dev` in P5 directory
   - Wrong GraphQL URL in `.env` → Check `VITE_GRAPHQL_URL=http://localhost:4000/graphql`
   - CORS errors → P5 should allow all origins in development mode

**GraphQL Network Errors**:
```bash
# Verify P5 is running
curl -H "Content-Type: application/json" \
  -d '{"query":"{ __typename }"}' \
  http://localhost:4000/graphql

# Check P6 .env configuration
cat services/cursor-viz-spa/.env
```

**Charts not rendering**:
- Check browser console for errors
- Verify data is returned from GraphQL (Network tab)
- Check if Recharts library is installed: `npm list recharts`

---

## Performance Testing

### Load Testing P5 GraphQL

```bash
# Install artillery (if not installed)
npm install -g artillery

# Run load test (100 requests/sec for 60 seconds)
artillery quick --count 100 --num 60 http://localhost:4000/graphql
```

### Database Query Performance

```bash
# Enable query logging in P5
# Add to .env:
# DATABASE_LOGGING=true

# Check slow queries
docker exec cursor-analytics-postgres psql -U cursor -d cursor_analytics \
  -c 'SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;'
```

---

## Cleanup

### Stop All Services

```bash
# Stop cursor-sim
docker stop cursor-sim-local

# Stop PostgreSQL
cd services/cursor-analytics-core
docker-compose down

# Stop P5 GraphQL
pkill -f "tsx watch src/index.ts"

# Stop P6 Dashboard
pkill -f "vite"
```

### Reset All Data

```bash
# Remove cursor-sim container
docker rm cursor-sim-local

# Remove PostgreSQL and data
cd services/cursor-analytics-core
docker-compose down -v

# Clear browser cache for P6
# Chrome: DevTools → Application → Clear Storage
```

---

## Next Steps

### Future Enhancements

1. **Step 04 (P5 Ingestion Worker)** - Automatic data sync from cursor-sim to PostgreSQL
2. **E2E Tests** - Automated integration tests using Playwright or Cypress
3. **Production Docker Compose** - Multi-service docker-compose for production deployment
4. **Monitoring** - Add Prometheus/Grafana for metrics collection
5. **Authentication** - Add JWT-based auth for GraphQL and dashboard

---

## Support

- **cursor-sim Documentation**: `services/cursor-sim/README.md`
- **P5 GraphQL API Docs**: `services/cursor-analytics-core/API.md`
- **P5 SPEC**: `services/cursor-analytics-core/SPEC.md`
- **P6 SPEC**: `services/cursor-viz-spa/SPEC.md`
- **Cloud Deployment**: `docs/cursor-sim-cloud-run.md`

## Service Status

| Service | Status | URL | Notes |
|---------|--------|-----|-------|
| cursor-sim (P4) | ✅ COMPLETE | http://localhost:8080 | Docker ready, REST API functional |
| cursor-analytics-core (P5) | ✅ COMPLETE | http://localhost:4000 | GraphQL API functional, 91.49% test coverage |
| cursor-viz-spa (P6) | ✅ COMPLETE | http://localhost:3000 | Dashboard ready, 91.68% test coverage |
| PostgreSQL | ✅ READY | localhost:5432 | Docker Compose configured |

**All integration components ready for testing!**
