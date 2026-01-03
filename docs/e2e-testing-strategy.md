# E2E and Integration Testing Enhancement Strategy

**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: HIGH

## Executive Summary

Current testing focuses on unit and component tests, but lacks integration and E2E coverage. The P5+P6 integration issues (Jan 4, 2026) revealed critical gaps:
- Dashboard integration untested
- GraphQL query compatibility unchecked
- Full data flow (P4 → P5 → P6) never validated

This document proposes a comprehensive testing pyramid with emphasis on integration and E2E tests.

---

## Current Testing State

### cursor-sim (P4)
- ✅ Unit tests: 80%+ coverage
- ✅ Integration tests: API endpoints tested with real storage
- ✅ E2E tests: Seed data generation validated
- ⚠️ No contract tests with P5

### cursor-analytics-core (P5)
- ✅ Unit tests: 91.49% coverage
- ✅ Integration tests: GraphQL resolvers with test database
- ✅ Resolver tests: All queries/mutations tested
- ⚠️ No E2E tests with P4 data ingestion
- ❌ No schema validation tests

### cursor-viz-spa (P6)
- ✅ Unit tests: 91.68% coverage
- ✅ Component tests: Chart components tested with mock data
- ⚠️ Component integration incomplete: Dashboard never tested with real hooks
- ❌ No E2E tests with P5 GraphQL
- ❌ No visual regression tests

---

## Testing Pyramid (Proposed)

```
                         ┌─────────────┐
                         │   E2E Tests │  ← 5% (NEW)
                         │  (Playwright)│
                         └─────────────┘
                       ┌────────────────────┐
                       │ Integration Tests  │  ← 15% (EXPAND)
                       │ (GraphQL, Docker)  │
                       └────────────────────┘
                  ┌──────────────────────────────┐
                  │   Component Tests            │  ← 30% (CURRENT)
                  │   (Vitest, React Testing Lib)│
                  └──────────────────────────────┘
            ┌────────────────────────────────────────┐
            │   Unit Tests                           │  ← 50% (CURRENT)
            │   (Vitest, Go test)                    │
            └────────────────────────────────────────┘
```

**Target Distribution:**
- **Unit Tests** (50%): Pure functions, utilities, models
- **Component Tests** (30%): React components, hooks (isolated)
- **Integration Tests** (15%): Service boundaries, GraphQL, database
- **E2E Tests** (5%): Full stack, critical user journeys

---

## Phase 1: Integration Testing Enhancements

### 1.1 P6 Component Integration Tests

**Goal**: Test that pages correctly integrate hooks and components.

#### Test: Dashboard Integration

```typescript
// services/cursor-viz-spa/src/pages/__tests__/Dashboard.integration.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import { MockedProvider } from '@apollo/client/testing';
import Dashboard from '../Dashboard';
import { GET_DASHBOARD_SUMMARY } from '../../graphql/queries';

describe('Dashboard Integration', () => {
  const mocks = [
    {
      request: {
        query: GET_DASHBOARD_SUMMARY,
      },
      result: {
        data: {
          dashboardSummary: {
            totalDevelopers: 3,
            activeDevelopers: 3,
            overallAcceptanceRate: 85.5,
            teamComparison: [
              {
                teamName: 'Backend',
                topPerformer: {
                  id: '1',
                  name: 'Alice',
                  email: 'alice@example.com',
                  seniority: 'senior',
                },
              },
            ],
            dailyTrend: [
              {
                date: '2026-01-01',
                suggestionsAccepted: 10,
                linesAdded: 100,
              },
            ],
          },
        },
      },
    },
  ];

  it('should fetch and display dashboard data', async () => {
    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <Dashboard />
      </MockedProvider>
    );

    // Loading state
    expect(screen.getByText(/loading dashboard data/i)).toBeInTheDocument();

    // Wait for data
    await waitFor(() => {
      expect(screen.queryByText(/loading/i)).not.toBeInTheDocument();
    });

    // KPI cards
    expect(screen.getByText('3')).toBeInTheDocument(); // Total devs
    expect(screen.getByText('85.5%')).toBeInTheDocument(); // Acceptance rate

    // Components rendered
    expect(screen.getByText(/velocity heatmap/i)).toBeInTheDocument();
    expect(screen.getByText(/team radar/i)).toBeInTheDocument();
    expect(screen.getByText(/developer table/i)).toBeInTheDocument();
  });

  it('should handle GraphQL errors gracefully', async () => {
    const errorMocks = [
      {
        request: {
          query: GET_DASHBOARD_SUMMARY,
        },
        error: new Error('Network error'),
      },
    ];

    render(
      <MockedProvider mocks={errorMocks} addTypename={false}>
        <Dashboard />
      </MockedProvider>
    );

    await waitFor(() => {
      expect(screen.getByText(/error loading dashboard/i)).toBeInTheDocument();
    });
  });

  it('should pass correct props to chart components', async () => {
    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <Dashboard />
      </MockedProvider>
    );

    await waitFor(() => {
      expect(screen.queryByText(/loading/i)).not.toBeInTheDocument();
    });

    // Verify VelocityHeatmap receives DailyStats[]
    const heatmap = screen.getByTestId('velocity-heatmap');
    expect(heatmap).toHaveAttribute('data-count', '1'); // 1 daily trend item

    // Verify TeamRadarChart receives selectedTeams
    const radar = screen.getByTestId('team-radar');
    expect(radar).toHaveAttribute('data-teams', 'Backend');
  });
});
```

**Run**:
```bash
cd services/cursor-viz-spa
npm run test:integration
```

---

### 1.2 P5 GraphQL Schema Validation Tests

**Goal**: Ensure P5 schema matches expected structure.

#### Test: Schema Structure

```typescript
// services/cursor-analytics-core/src/graphql/__tests__/schema.test.ts
import { graphql, buildSchema } from 'graphql';
import { typeDefs } from '../schema';

describe('GraphQL Schema', () => {
  let schema: any;

  beforeAll(() => {
    schema = buildSchema(typeDefs);
  });

  it('should have TeamStats type with correct fields', () => {
    const teamStatsType = schema.getType('TeamStats');
    expect(teamStatsType).toBeDefined();

    const fields = teamStatsType.getFields();
    expect(fields.topPerformer).toBeDefined();
    expect(fields.topPerformer.type.toString()).toBe('Developer');

    // Ensure old field doesn't exist
    expect(fields.topPerformers).toBeUndefined();
  });

  it('should have DailyStats type with correct fields', () => {
    const dailyStatsType = schema.getType('DailyStats');
    expect(dailyStatsType).toBeDefined();

    const fields = dailyStatsType.getFields();
    expect(fields.linesAdded).toBeDefined();
    expect(fields.aiLinesAdded).toBeDefined();

    // Ensure wrong field doesn't exist
    expect(fields.humanLinesAdded).toBeUndefined();
  });

  it('should validate DashboardKPI query', async () => {
    const query = `
      query {
        dashboardSummary {
          totalDevelopers
          teamComparison {
            topPerformer {
              name
            }
          }
        }
      }
    `;

    const result = await graphql({
      schema,
      source: query,
    });

    expect(result.errors).toBeUndefined();
  });

  it('should reject invalid field names', async () => {
    const query = `
      query {
        dashboardSummary {
          teamComparison {
            topPerformers {  # WRONG: Should be topPerformer
              name
            }
          }
        }
      }
    `;

    const result = await graphql({
      schema,
      source: query,
    });

    expect(result.errors).toBeDefined();
    expect(result.errors![0].message).toContain('Cannot query field "topPerformers"');
  });
});
```

**Run**:
```bash
cd services/cursor-analytics-core
npm run test:schema
```

---

### 1.3 P5 Integration Tests with Docker PostgreSQL

**Goal**: Test resolvers with real database (not mocks).

#### Test: Dashboard Summary Resolver

```typescript
// services/cursor-analytics-core/src/resolvers/__tests__/dashboardSummary.integration.test.ts
import { ApolloServer } from '@apollo/server';
import { PrismaClient } from '../generated/prisma';
import { typeDefs } from '../graphql/schema';
import { resolvers } from '../graphql/resolvers';

describe('DashboardSummary Resolver Integration', () => {
  let server: ApolloServer;
  let prisma: PrismaClient;

  beforeAll(async () => {
    // Use test database
    prisma = new PrismaClient({
      datasources: {
        db: {
          url: process.env.TEST_DATABASE_URL,
        },
      },
    });

    server = new ApolloServer({
      typeDefs,
      resolvers,
      context: () => ({ prisma }),
    });

    // Seed test data
    await prisma.developer.createMany({
      data: [
        { externalId: 'alice', name: 'Alice', email: 'alice@test.com', team: 'Backend', seniority: 'senior' },
        { externalId: 'bob', name: 'Bob', email: 'bob@test.com', team: 'Frontend', seniority: 'mid' },
      ],
    });
  });

  afterAll(async () => {
    await prisma.$executeRaw`TRUNCATE TABLE developers CASCADE`;
    await prisma.$disconnect();
  });

  it('should return correct developer counts', async () => {
    const result = await server.executeOperation({
      query: `
        query {
          dashboardSummary {
            totalDevelopers
            activeDevelopers
          }
        }
      `,
    });

    expect(result.body.kind).toBe('single');
    if (result.body.kind === 'single') {
      expect(result.body.singleResult.data?.dashboardSummary).toEqual({
        totalDevelopers: 2,
        activeDevelopers: 2,
      });
    }
  });

  it('should return team comparison with topPerformer', async () => {
    const result = await server.executeOperation({
      query: `
        query {
          dashboardSummary {
            teamComparison {
              teamName
              memberCount
              topPerformer {
                name
                seniority
              }
            }
          }
        }
      `,
    });

    expect(result.body.kind).toBe('single');
    if (result.body.kind === 'single') {
      const teams = result.body.singleResult.data?.dashboardSummary.teamComparison;
      expect(teams).toHaveLength(2);
      expect(teams.find((t: any) => t.teamName === 'Backend')?.topPerformer?.name).toBe('Alice');
    }
  });
});
```

**Run**:
```bash
cd services/cursor-analytics-core
docker-compose -f docker-compose.test.yml up -d
npm run test:integration
docker-compose -f docker-compose.test.yml down
```

---

## Phase 2: End-to-End Testing with Playwright

### 2.1 Setup Playwright

```bash
cd services/cursor-viz-spa
npm install -D @playwright/test
npx playwright install
```

#### Configuration

```typescript
// services/cursor-viz-spa/playwright.config.ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
  },
});
```

---

### 2.2 E2E Test: Full Stack Data Flow

**Goal**: Validate P4 → P5 → P6 integration.

#### Test: Dashboard E2E

```typescript
// services/cursor-viz-spa/tests/e2e/dashboard.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Dashboard E2E', () => {
  test.beforeAll(async () => {
    // Assume P4 and P5 are running in Docker
    // Check health endpoints
    const p4Response = await fetch('http://localhost:8080/health');
    expect(p4Response.status).toBe(200);

    const p5Response = await fetch('http://localhost:4000/graphql', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query: '{ health { status } }' }),
    });
    expect(p5Response.status).toBe(200);
  });

  test('should display dashboard with real data from P5', async ({ page }) => {
    // Navigate to dashboard
    await page.goto('/dashboard');

    // Wait for GraphQL request to complete
    const graphqlResponse = await page.waitForResponse(
      (resp) => resp.url().includes('/graphql') && resp.status() === 200
    );

    // Verify response data
    const responseBody = await graphqlResponse.json();
    expect(responseBody.data.dashboardSummary.totalDevelopers).toBeGreaterThan(0);

    // Verify UI updates with data
    await expect(page.locator('[data-testid="total-devs"]')).toContainText(/\d+/);

    // Verify charts rendered (not placeholders)
    await expect(page.locator('[data-testid="velocity-heatmap"]')).toBeVisible();
    await expect(page.locator('.recharts-surface')).toBeVisible(); // TeamRadarChart

    // Verify developer table has rows
    const tableRows = page.locator('[data-testid="developer-table"] tbody tr');
    await expect(tableRows).not.toHaveCount(0);

    // Take screenshot for visual regression
    await expect(page).toHaveScreenshot('dashboard-loaded.png');
  });

  test('should handle API errors gracefully', async ({ page, context }) => {
    // Block GraphQL requests to simulate server error
    await context.route('**/graphql', (route) => {
      route.fulfill({
        status: 500,
        body: JSON.stringify({ errors: [{ message: 'Internal server error' }] }),
      });
    });

    await page.goto('/dashboard');

    // Verify error message displayed
    await expect(page.getByText(/error loading dashboard/i)).toBeVisible();

    // Verify no chart placeholders shown
    await expect(page.getByText(/chart placeholder/i)).not.toBeVisible();
  });

  test('should make only one GraphQL request on load', async ({ page }) => {
    const graphqlRequests: any[] = [];

    page.on('request', (request) => {
      if (request.url().includes('/graphql')) {
        graphqlRequests.push({
          method: request.method(),
          postData: request.postDataJSON(),
        });
      }
    });

    await page.goto('/dashboard');
    await page.waitForLoadState('networkidle');

    // Should make exactly 1 request for dashboardSummary
    expect(graphqlRequests).toHaveLength(1);
    expect(graphqlRequests[0].postData.query).toContain('dashboardSummary');
  });

  test('should navigate to teams and developers pages', async ({ page }) => {
    await page.goto('/dashboard');

    // Navigate to Teams
    await page.click('a[href="/teams"]');
    await expect(page).toHaveURL('/teams');
    await expect(page.locator('h1')).toContainText(/teams/i);

    // Navigate to Developers
    await page.click('a[href="/developers"]');
    await expect(page).toHaveURL('/developers');
    await expect(page.locator('h1')).toContainText(/developers/i);
  });
});
```

**Run**:
```bash
# Start full stack
cd services/cursor-sim && docker run --rm -d -p 8080:8080 cursor-sim:latest
cd services/cursor-analytics-core && docker-compose up -d
cd services/cursor-viz-spa && npm run dev &

# Run E2E tests
npx playwright test

# Generate HTML report
npx playwright show-report
```

---

### 2.3 Visual Regression Tests

**Goal**: Detect UI changes automatically.

#### Test: Dashboard Visual Regression

```typescript
// services/cursor-viz-spa/tests/e2e/visual.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Visual Regression', () => {
  test('dashboard should match baseline', async ({ page }) => {
    await page.goto('/dashboard');
    await page.waitForResponse((resp) => resp.url().includes('/graphql'));

    // Wait for all charts to render
    await page.waitForSelector('[data-testid="velocity-heatmap"]');
    await page.waitForSelector('.recharts-surface');

    // Full page screenshot
    await expect(page).toHaveScreenshot('dashboard-full.png', {
      fullPage: true,
      animations: 'disabled',
    });
  });

  test('KPI cards should match baseline', async ({ page }) => {
    await page.goto('/dashboard');
    await page.waitForResponse((resp) => resp.url().includes('/graphql'));

    // Specific component screenshot
    const kpiCards = page.locator('[data-testid="kpi-cards"]');
    await expect(kpiCards).toHaveScreenshot('kpi-cards.png');
  });

  test('velocity heatmap should match baseline', async ({ page }) => {
    await page.goto('/dashboard');
    await page.waitForResponse((resp) => resp.url().includes('/graphql'));

    const heatmap = page.locator('[data-testid="velocity-heatmap"]');
    await expect(heatmap).toHaveScreenshot('velocity-heatmap.png');
  });
});
```

**Update on Intentional Changes**:
```bash
# Update baseline screenshots
npx playwright test --update-snapshots
```

---

## Phase 3: CI/CD Integration

### 3.1 GitHub Actions Workflow

```yaml
# .github/workflows/e2e.yml
name: E2E Tests

on:
  pull_request:
    paths:
      - 'services/cursor-analytics-core/**'
      - 'services/cursor-viz-spa/**'
  push:
    branches: [main]

jobs:
  e2e:
    runs-on: ubuntu-latest
    timeout-minutes: 15

    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'

      - name: Start cursor-sim (P4)
        run: |
          cd services/cursor-sim
          docker build -t cursor-sim:test .
          docker run -d --name cursor-sim -p 8080:8080 cursor-sim:test

      - name: Start cursor-analytics-core (P5)
        run: |
          cd services/cursor-analytics-core
          docker-compose up -d
          sleep 10  # Wait for PostgreSQL to be ready

      - name: Seed P5 database
        run: |
          docker exec cursor-analytics-postgres psql -U cursor -d cursor_analytics << 'EOF'
          INSERT INTO developers (id, external_id, name, email, team, seniority)
          VALUES
            (gen_random_uuid(), 'alice', 'Alice', 'alice@test.com', 'Backend', 'senior'),
            (gen_random_uuid(), 'bob', 'Bob', 'bob@test.com', 'Frontend', 'mid');
          EOF

      - name: Install P6 dependencies
        run: |
          cd services/cursor-viz-spa
          npm ci

      - name: Install Playwright
        run: |
          cd services/cursor-viz-spa
          npx playwright install --with-deps

      - name: Run E2E tests
        run: |
          cd services/cursor-viz-spa
          npm run test:e2e

      - name: Upload Playwright report
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: services/cursor-viz-spa/playwright-report/
          retention-days: 7

      - name: Upload visual regression diffs
        if: failure()
        uses: actions/upload-artifact@v3
        with:
          name: visual-diffs
          path: services/cursor-viz-spa/test-results/
          retention-days: 7

      - name: Cleanup
        if: always()
        run: |
          docker stop cursor-sim || true
          docker rm cursor-sim || true
          cd services/cursor-analytics-core && docker-compose down
```

---

## Phase 4: Performance Testing

### 4.1 Lighthouse CI

**Goal**: Ensure dashboard loads fast.

```yaml
# .github/workflows/lighthouse.yml
name: Lighthouse CI

on:
  pull_request:
    paths:
      - 'services/cursor-viz-spa/**'

jobs:
  lighthouse:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: treosh/lighthouse-ci-action@v9
        with:
          urls: |
            http://localhost:3000/dashboard
          budgetPath: ./lighthouse-budget.json
          uploadArtifacts: true
```

**Budget**:
```json
// services/cursor-viz-spa/lighthouse-budget.json
{
  "ci": {
    "assert": {
      "preset": "lighthouse:recommended",
      "assertions": {
        "first-contentful-paint": ["warn", { "maxNumericValue": 2000 }],
        "interactive": ["error", { "maxNumericValue": 5000 }],
        "speed-index": ["warn", { "maxNumericValue": 3500 }]
      }
    }
  }
}
```

---

## Testing Matrix

| Service | Unit | Component | Integration | E2E | Visual | Contract |
|---------|------|-----------|-------------|-----|--------|----------|
| **P4 (cursor-sim)** | ✅ Go test | N/A | ✅ API tests | ⚠️ Planned | N/A | ⚠️ Planned |
| **P5 (cursor-analytics-core)** | ✅ Vitest | N/A | ✅ DB tests | ⚠️ Planned | N/A | ✅ Schema tests |
| **P6 (cursor-viz-spa)** | ✅ Vitest | ✅ RTL | ✅ Page tests | ✅ Playwright | ✅ Snapshots | ✅ Codegen |

---

## Rollout Timeline

### Week 1: Foundation
- [ ] Add P6 component integration tests (Dashboard, TeamList, DeveloperList)
- [ ] Add P5 schema validation tests
- [ ] Document test patterns in README

### Week 2: E2E Setup
- [ ] Install and configure Playwright
- [ ] Write first E2E test (Dashboard happy path)
- [ ] Add visual regression baseline screenshots
- [ ] Set up CI workflow for E2E tests

### Week 3: Coverage Expansion
- [ ] Add E2E tests for Teams and Developers pages
- [ ] Add error scenario tests (network errors, empty states)
- [ ] Add performance tests with Lighthouse
- [ ] Add contract tests with GraphQL Inspector

### Week 4: Production Readiness
- [ ] Run E2E tests nightly against staging
- [ ] Set up visual regression monitoring
- [ ] Add Slack alerts for test failures
- [ ] Document troubleshooting guide

---

## Success Metrics

### Coverage Targets

| Test Type | Current | Target (6 months) |
|-----------|---------|-------------------|
| Unit Test Coverage | 91% | 90%+ (maintain) |
| Component Integration | 0% | 80%+ |
| E2E Test Coverage (critical paths) | 0% | 100% |
| Visual Regression (pages) | 0% | 100% |

### Quality Metrics

| Metric | Baseline (Before) | Target (After) |
|--------|-------------------|----------------|
| Integration bugs found in production | 4 (Jan 2026) | 0 |
| Time to detect integration issues | Manual testing (days) | CI failure (minutes) |
| GraphQL schema mismatches | 100% at runtime | 0% (caught at compile-time) |
| Visual regressions | Manual QA | Automated detection |

---

## Appendix: Test Commands Reference

```bash
# P6: Run all tests
cd services/cursor-viz-spa
npm run test              # Unit + component tests
npm run test:integration  # Component integration tests
npm run test:e2e          # Playwright E2E tests
npm run test:visual       # Visual regression tests
npm run test:all          # All of the above

# P5: Run all tests
cd services/cursor-analytics-core
npm run test              # Unit tests
npm run test:integration  # Integration tests with DB
npm run test:schema       # Schema validation tests

# P4: Run all tests
cd services/cursor-sim
go test ./...             # All Go tests
go test -v ./test/e2e/... # E2E API tests
```
