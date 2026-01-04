# Design Document: E2E Testing with Playwright

**Feature ID**: P6-F05
**Epic**: P6 - cursor-viz-spa (Testing Enhancement)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

Implement Playwright E2E tests for cursor-viz-spa that test the full P4→P5→P6 data flow through a real browser.

---

## Architecture

### E2E Test Environment

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           E2E Test Execution                           │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  Playwright Test Runner                                                │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │ • Chromium browser (headless in CI)                               │  │
│  │ • Page navigation and interactions                                │  │
│  │ • Network request interception                                    │  │
│  │ • Screenshot comparison                                           │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│           │                                                             │
│           ▼                                                             │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  cursor-viz-spa (P6)      http://localhost:3000                   │  │
│  │  ├── Dashboard page                                               │  │
│  │  ├── Teams page                                                   │  │
│  │  └── Developers page                                              │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│           │                                                             │
│           │ GraphQL Requests                                           │
│           ▼                                                             │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  cursor-analytics-core (P5)   http://localhost:4000               │  │
│  │  ├── GraphQL API                                                  │  │
│  │  └── PostgreSQL Database                                          │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│           │                                                             │
│           │ REST API Requests                                          │
│           ▼                                                             │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  cursor-sim (P4)              http://localhost:8080               │  │
│  │  └── Simulation Data API                                          │  │
│  └───────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Configuration

### Playwright Config

```typescript
// services/cursor-viz-spa/playwright.config.ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: false,  // Sequential for deterministic data
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: [
    ['html', { open: 'never' }],
    ['json', { outputFile: 'test-results/results.json' }],
  ],
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'on-first-retry',
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
    timeout: 120000,
  },
});
```

---

## Test Structure

```
tests/e2e/
├── fixtures/
│   └── test-data.ts         # Test data helpers
├── dashboard.spec.ts        # Dashboard tests
├── navigation.spec.ts       # Navigation tests
├── error-handling.spec.ts   # Error state tests
└── visual.spec.ts           # Visual regression tests
```

---

## Test Implementation

### Dashboard E2E Test

```typescript
// tests/e2e/dashboard.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Dashboard E2E', () => {
  test.beforeAll(async () => {
    // Verify services are running
    const p5Response = await fetch('http://localhost:4000/graphql', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query: '{ health { status } }' }),
    });
    expect(p5Response.ok).toBe(true);
  });

  test('loads with real data from P5', async ({ page }) => {
    await page.goto('/dashboard');

    // Wait for GraphQL request
    const response = await page.waitForResponse(
      resp => resp.url().includes('/graphql') && resp.status() === 200
    );
    const data = await response.json();
    expect(data.data.dashboardSummary.totalDevelopers).toBeGreaterThan(0);

    // Verify UI shows data
    await expect(page.locator('[data-testid="total-devs"]')).toContainText(/\d+/);
    await expect(page.locator('[data-testid="velocity-heatmap"]')).toBeVisible();
  });

  test('displays KPI cards', async ({ page }) => {
    await page.goto('/dashboard');
    await page.waitForLoadState('networkidle');

    // Check KPI cards exist
    await expect(page.locator('[data-testid="kpi-total-devs"]')).toBeVisible();
    await expect(page.locator('[data-testid="kpi-active-devs"]')).toBeVisible();
    await expect(page.locator('[data-testid="kpi-acceptance-rate"]')).toBeVisible();
  });

  test('renders charts without placeholders', async ({ page }) => {
    await page.goto('/dashboard');
    await page.waitForLoadState('networkidle');

    // Verify no placeholder text
    await expect(page.getByText(/chart placeholder/i)).not.toBeVisible();
    await expect(page.locator('.recharts-surface')).toBeVisible();
  });
});
```

### Visual Regression Test

```typescript
// tests/e2e/visual.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Visual Regression', () => {
  test('dashboard matches baseline', async ({ page }) => {
    await page.goto('/dashboard');
    await page.waitForLoadState('networkidle');

    // Wait for all charts to render
    await page.waitForSelector('[data-testid="velocity-heatmap"]');
    await page.waitForSelector('.recharts-surface');

    // Full page screenshot
    await expect(page).toHaveScreenshot('dashboard.png', {
      fullPage: true,
      animations: 'disabled',
    });
  });

  test('KPI cards match baseline', async ({ page }) => {
    await page.goto('/dashboard');
    await page.waitForLoadState('networkidle');

    const kpiCards = page.locator('[data-testid="kpi-cards"]');
    await expect(kpiCards).toHaveScreenshot('kpi-cards.png');
  });
});
```

### Error Handling Test

```typescript
// tests/e2e/error-handling.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Error Handling', () => {
  test('shows error when P5 unavailable', async ({ page, context }) => {
    // Block GraphQL requests
    await context.route('**/graphql', route => {
      route.fulfill({ status: 500, body: 'Internal Server Error' });
    });

    await page.goto('/dashboard');

    // Verify error message
    await expect(page.getByText(/error/i)).toBeVisible();
    await expect(page.locator('[data-testid="velocity-heatmap"]')).not.toBeVisible();
  });
});
```

---

## CI Integration

### GitHub Actions

```yaml
# .github/workflows/e2e.yml
name: E2E Tests

on:
  pull_request:
    paths:
      - 'services/cursor-viz-spa/**'
      - 'services/cursor-analytics-core/**'

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

      - name: Start P4 (cursor-sim)
        run: |
          cd services/cursor-sim
          docker build -t cursor-sim:test .
          docker run -d --name cursor-sim -p 8080:8080 cursor-sim:test

      - name: Start P5 (cursor-analytics-core)
        run: |
          cd services/cursor-analytics-core
          docker-compose up -d
          sleep 10

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

      - name: Upload artifacts
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: services/cursor-viz-spa/playwright-report/
          retention-days: 7
```

---

## NPM Scripts

```json
{
  "scripts": {
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "test:e2e:debug": "playwright test --debug",
    "test:visual:update": "playwright test visual.spec.ts --update-snapshots"
  }
}
```

---

## Success Metrics

| Metric | Before | After |
|--------|--------|-------|
| Full stack test coverage | 0% | 100% (critical paths) |
| Visual regression coverage | 0% | 100% (key pages) |
| Integration issue detection | Manual | Automated |
| Time to detect issues | Days | Minutes |

---

## References

- [Playwright Documentation](https://playwright.dev)
- [Visual Comparisons](https://playwright.dev/docs/test-snapshots)
- `docs/e2e-testing-strategy.md` (Phase 2)
