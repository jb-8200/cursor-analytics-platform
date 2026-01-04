# Design Document: Component Integration Tests

**Feature ID**: P6-F04
**Epic**: P6 - cursor-viz-spa (Testing Enhancement)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

Implement page-level integration tests that verify components correctly integrate with GraphQL hooks and child components, using MockedProvider for controlled data flow.

---

## Architecture

### Testing Pyramid Context

```
                     ┌─────────────┐
                     │    E2E      │  ← P6-F05 (Playwright)
                     │  Full Stack │
                     └─────────────┘
                   ┌─────────────────┐
                   │   INTEGRATION   │  ← THIS FEATURE
                   │  Page + Hooks   │
                   └─────────────────┘
              ┌─────────────────────────┐
              │      COMPONENT          │  ← Exists (91.68%)
              │   Individual Charts     │
              └─────────────────────────┘
         ┌───────────────────────────────────┐
         │           UNIT                    │  ← Exists (91.68%)
         │   Hooks, Utils, Pure Functions    │
         └───────────────────────────────────┘
```

### Integration Test Structure

```
┌─────────────────────────────────────────────────────────────────────────┐
│  Integration Test Setup                                                │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  MockedProvider                                                        │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │ mocks: [                                                          │  │
│  │   {                                                               │  │
│  │     request: { query: GET_DASHBOARD_SUMMARY },                    │  │
│  │     result: { data: { dashboardSummary: {...} } }                │  │
│  │   }                                                               │  │
│  │ ]                                                                 │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│           │                                                             │
│           ▼                                                             │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  MemoryRouter (for routing context)                               │  │
│  │  ├── Dashboard Page                                               │  │
│  │  │   ├── useDashboard() hook  ← Receives mock data               │  │
│  │  │   ├── KPI Cards  ← Renders with mock data                     │  │
│  │  │   ├── VelocityHeatmap  ← Renders with mock data               │  │
│  │  │   ├── TeamRadarChart  ← Renders with mock data                │  │
│  │  │   └── DeveloperTable  ← Renders with mock data                │  │
│  │  └───────────────────────────────────────────────────────────────│  │
│  └───────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Test Implementation

### Test Utility

```typescript
// src/test-utils/integration.tsx
import { render, RenderOptions } from '@testing-library/react';
import { MockedProvider, MockedResponse } from '@apollo/client/testing';
import { MemoryRouter } from 'react-router-dom';

interface IntegrationRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  mocks?: MockedResponse[];
  route?: string;
}

export function renderWithProviders(
  ui: React.ReactElement,
  { mocks = [], route = '/', ...options }: IntegrationRenderOptions = {}
) {
  return render(ui, {
    wrapper: ({ children }) => (
      <MockedProvider mocks={mocks} addTypename={false}>
        <MemoryRouter initialEntries={[route]}>
          {children}
        </MemoryRouter>
      </MockedProvider>
    ),
    ...options,
  });
}
```

### Mock Data Factory

```typescript
// src/test-utils/mocks.ts
import { GET_DASHBOARD_SUMMARY } from '../graphql/queries';

export const mockDashboardData = {
  dashboardSummary: {
    totalDevelopers: 10,
    activeDevelopers: 8,
    overallAcceptanceRate: 78.5,
    teamComparison: [
      {
        teamName: 'Backend',
        memberCount: 5,
        topPerformer: {
          id: '1',
          name: 'Alice',
          email: 'alice@test.com',
          seniority: 'senior',
        },
        avgAcceptanceRate: 82.3,
        totalCommits: 150,
      },
    ],
    dailyTrend: [
      {
        date: '2026-01-01',
        suggestionsAccepted: 45,
        suggestionsRejected: 12,
        linesAdded: 500,
        aiLinesAdded: 350,
      },
    ],
  },
};

export const dashboardMock = {
  request: { query: GET_DASHBOARD_SUMMARY },
  result: { data: mockDashboardData },
};

export const dashboardErrorMock = {
  request: { query: GET_DASHBOARD_SUMMARY },
  error: new Error('Failed to load dashboard data'),
};

export const dashboardLoadingMock = {
  request: { query: GET_DASHBOARD_SUMMARY },
  delay: Infinity, // Never resolves
};
```

---

## Test Cases

### Dashboard Integration Tests

```typescript
// src/pages/__tests__/Dashboard.integration.test.tsx
import { screen, waitFor } from '@testing-library/react';
import { renderWithProviders } from '../../test-utils/integration';
import { dashboardMock, dashboardErrorMock, dashboardLoadingMock } from '../../test-utils/mocks';
import Dashboard from '../Dashboard';

describe('Dashboard Integration', () => {
  describe('Success State', () => {
    it('displays KPI cards with correct data', async () => {
      renderWithProviders(<Dashboard />, { mocks: [dashboardMock] });

      await waitFor(() => {
        expect(screen.getByText('10')).toBeInTheDocument(); // Total devs
        expect(screen.getByText('78.5%')).toBeInTheDocument(); // Acceptance rate
      });
    });

    it('renders VelocityHeatmap with daily trend data', async () => {
      renderWithProviders(<Dashboard />, { mocks: [dashboardMock] });

      await waitFor(() => {
        expect(screen.getByTestId('velocity-heatmap')).toBeInTheDocument();
      });
    });

    it('renders TeamRadarChart with team comparison data', async () => {
      renderWithProviders(<Dashboard />, { mocks: [dashboardMock] });

      await waitFor(() => {
        expect(screen.getByTestId('team-radar')).toBeInTheDocument();
        expect(screen.getByText('Backend')).toBeInTheDocument();
      });
    });

    it('renders DeveloperTable with team data', async () => {
      renderWithProviders(<Dashboard />, { mocks: [dashboardMock] });

      await waitFor(() => {
        expect(screen.getByText('Alice')).toBeInTheDocument();
      });
    });
  });

  describe('Loading State', () => {
    it('displays loading indicator', () => {
      renderWithProviders(<Dashboard />, { mocks: [dashboardLoadingMock] });

      expect(screen.getByText(/loading/i)).toBeInTheDocument();
    });

    it('does not display data while loading', () => {
      renderWithProviders(<Dashboard />, { mocks: [dashboardLoadingMock] });

      expect(screen.queryByTestId('velocity-heatmap')).not.toBeInTheDocument();
    });
  });

  describe('Error State', () => {
    it('displays error message', async () => {
      renderWithProviders(<Dashboard />, { mocks: [dashboardErrorMock] });

      await waitFor(() => {
        expect(screen.getByText(/error/i)).toBeInTheDocument();
      });
    });

    it('does not display charts on error', async () => {
      renderWithProviders(<Dashboard />, { mocks: [dashboardErrorMock] });

      await waitFor(() => {
        expect(screen.queryByTestId('velocity-heatmap')).not.toBeInTheDocument();
      });
    });
  });
});
```

---

## File Structure

```
src/
├── test-utils/
│   ├── integration.tsx     # renderWithProviders helper
│   └── mocks.ts           # Mock data factories
├── pages/
│   └── __tests__/
│       ├── Dashboard.integration.test.tsx
│       ├── Teams.integration.test.tsx
│       └── Developers.integration.test.tsx
```

---

## Test Configuration

### Vitest Config Update

```typescript
// vitest.config.ts
export default defineConfig({
  test: {
    // Separate integration tests
    include: ['src/**/*.{test,spec}.{ts,tsx}'],
    // Integration tests may take longer
    testTimeout: 10000,
    // Coverage
    coverage: {
      exclude: ['src/test-utils/**'],
    },
  },
});
```

### NPM Scripts

```json
{
  "scripts": {
    "test": "vitest run",
    "test:unit": "vitest run --exclude '**/*.integration.test.*'",
    "test:integration": "vitest run --include '**/*.integration.test.*'"
  }
}
```

---

## Mock Data Alignment

Critical: Mock data MUST match P5 GraphQL schema exactly.

**Source of Truth**:
- After P6-F02: Use generated types from `src/graphql/generated.ts`
- Mock data should be typed with generated interfaces

```typescript
import type { DashboardKPI } from '../graphql/generated';

export const mockDashboardData: { dashboardSummary: DashboardKPI } = {
  dashboardSummary: {
    totalDevelopers: 10,
    // TypeScript validates this matches P5 schema!
  },
};
```

---

## Success Metrics

| Metric | Before | After |
|--------|--------|-------|
| Page integration coverage | 0% | 80%+ |
| Hook integration coverage | 0% | 80%+ |
| Error state coverage | Partial | Full |
| Loading state coverage | Partial | Full |

---

## References

- [Apollo Client Testing](https://www.apollographql.com/docs/react/development-testing/testing/)
- [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
- `docs/e2e-testing-strategy.md` (Phase 1.1)
