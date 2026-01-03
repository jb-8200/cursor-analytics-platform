# Task Breakdown: cursor-viz-spa Dashboard

**Feature ID**: P6-cursor-viz-spa
**Phase**: P6 (Visualization SPA)
**Created**: January 3, 2026
**Status**: IN_PROGRESS

## Summary

**User Story**: `.work-items/P6-cursor-viz-spa/user-story.md`
**Design Doc**: `.work-items/P6-cursor-viz-spa/design.md`
**Service**: cursor-viz-spa
**Specification**: `services/cursor-viz-spa/SPEC.md`

**Total Estimated Hours**: 24.0h
**Total Tasks**: 7

## Progress Tracker

| Task ID | Task | Hours | Status | Actual |
|---------|------|-------|--------|--------|
| TASK01 | Project Setup (Dependencies, Tailwind, ESLint, Vitest) | 3.0 | DONE | 2.5 |
| TASK02 | Apollo Client Configuration & GraphQL Setup | 3.0 | DONE | 2.0 |
| TASK03 | Core Layout Components & Routing | 3.5 | DONE | 3.0 |
| TASK04 | Chart Components (Heatmap, Radar, Table) | 5.0 | NOT_STARTED | - |
| TASK05 | Filter Controls & Date Picker | 3.0 | NOT_STARTED | - |
| TASK06 | GraphQL Queries & Custom Hooks | 4.0 | NOT_STARTED | - |
| TASK07 | Testing Setup & Initial Tests | 2.5 | NOT_STARTED | - |

**Current Task**: TASK03 (Complete)

---

## Task Details

### TASK01: Project Setup (Dependencies, Tailwind, ESLint, Vitest)

**Estimated**: 3.0h
**Status**: DONE
**Actual**: 2.5h

**Objective**: Complete the foundational setup for the React/Vite project with all necessary tooling.

**Files**:
- `services/cursor-viz-spa/package.json`
- `services/cursor-viz-spa/tailwind.config.js`
- `services/cursor-viz-spa/postcss.config.js`
- `services/cursor-viz-spa/.eslintrc.cjs`
- `services/cursor-viz-spa/.prettierrc`
- `services/cursor-viz-spa/vitest.config.ts`
- `services/cursor-viz-spa/src/test/setup.ts`
- `services/cursor-viz-spa/.env.example`

**Tasks**:
- [ ] Install missing dependencies (Tailwind, router, date-fns, testing libs)
- [ ] Configure Tailwind CSS with custom theme from SPEC
- [ ] Set up ESLint + Prettier for code quality
- [ ] Configure Vitest for unit and integration tests
- [ ] Create directory structure (components/, hooks/, graphql/, utils/)
- [ ] Add environment variable configuration
- [ ] Update npm scripts for dev, build, test, lint
- [ ] Create test setup with MSW for API mocking
- [ ] Write setup validation tests

**Acceptance Criteria**:
- All dependencies installed and listed in package.json
- `npm run dev` starts development server on port 3000
- `npm run test` runs Vitest tests
- `npm run lint` checks code with ESLint
- `npm run build` produces optimized build
- Tailwind utilities are available in components
- TypeScript strict mode enabled and passing

**TDD Approach**:
1. RED: Write test for Tailwind config existence
2. GREEN: Create Tailwind config with custom theme
3. REFACTOR: Validate theme colors match SPEC

---

### TASK02: Apollo Client Configuration & GraphQL Setup

**Estimated**: 3.0h
**Status**: DONE
**Actual**: 2.0h

**Objective**: Set up Apollo Client for GraphQL communication with cursor-analytics-core.

**Files**:
- `services/cursor-viz-spa/src/graphql/client.ts`
- `services/cursor-viz-spa/src/graphql/queries.ts`
- `services/cursor-viz-spa/src/graphql/types.ts`
- `services/cursor-viz-spa/src/App.tsx`
- `services/cursor-viz-spa/codegen.yml`

**Tasks**:
- [x] Create Apollo Client instance with cache policies
- [x] Define initial GraphQL query stubs
- [x] Set up GraphQL Code Generator
- [x] Wrap app with ApolloProvider
- [x] Create mock GraphQL server for testing
- [x] Write tests for client configuration

**Acceptance Criteria**:
- Apollo Client connects to VITE_GRAPHQL_URL
- Cache policies configured for Developer and DailyStats
- GraphQL queries are type-safe
- Mock server available for tests

**TDD Approach**:
1. RED: Write test for Apollo Client initialization
2. GREEN: Create client with minimal config
3. REFACTOR: Add cache policies and error handling

---

### TASK03: Core Layout Components & Routing

**Estimated**: 3.5h
**Status**: DONE
**Actual**: 3.0h

**Objective**: Build responsive layout structure and routing.

**Files**:
- `services/cursor-viz-spa/src/components/layout/Header.tsx`
- `services/cursor-viz-spa/src/components/layout/Sidebar.tsx`
- `services/cursor-viz-spa/src/pages/Dashboard.tsx`
- `services/cursor-viz-spa/src/pages/TeamList.tsx`
- `services/cursor-viz-spa/src/pages/DeveloperList.tsx`
- `services/cursor-viz-spa/src/App.tsx`

**Tasks**:
- [x] Create responsive Header with user menu placeholder
- [x] Create collapsible Sidebar for navigation
- [x] Set up React Router with routes from SPEC
- [x] Build Dashboard page layout
- [x] Create page stubs for Teams and Developers
- [x] Write component tests

**Acceptance Criteria**:
- Header displays on all pages
- Sidebar collapses on mobile (< 768px)
- Routes navigate correctly
- Loading and error states render properly

**TDD Approach**:
1. RED: Write test for Header rendering
2. GREEN: Create Header component
3. REFACTOR: Extract KPI card into separate component

---

### TASK04: Chart Components (Heatmap, Radar, Table)

**Estimated**: 5.0h
**Status**: NOT_STARTED
**Actual**: -

**Objective**: Implement the three primary visualization components.

**Files**:
- `services/cursor-viz-spa/src/components/charts/VelocityHeatmap.tsx`
- `services/cursor-viz-spa/src/components/charts/TeamRadarChart.tsx`
- `services/cursor-viz-spa/src/components/charts/DeveloperTable.tsx`

**Tasks**:
- [ ] Build VelocityHeatmap with GitHub-style grid
- [ ] Build TeamRadarChart using Recharts
- [ ] Build DeveloperTable with sorting and pagination
- [ ] Add tooltips and interactions
- [ ] Style components with Tailwind
- [ ] Write comprehensive component tests

**Acceptance Criteria**:
- Heatmap displays 52 weeks of data
- Radar chart shows 2-5 teams overlapping
- Table sorts by any column
- Components are accessible (WCAG 2.1 AA)

**TDD Approach**:
1. RED: Write test for VelocityHeatmap data rendering
2. GREEN: Create basic grid structure
3. REFACTOR: Add color scale and tooltips

---

### TASK05: Filter Controls & Date Picker

**Estimated**: 3.0h
**Status**: NOT_STARTED
**Actual**: -

**Objective**: Build date range picker and filter controls.

**Files**:
- `services/cursor-viz-spa/src/components/filters/DateRangePicker.tsx`
- `services/cursor-viz-spa/src/hooks/useDateRange.ts`
- `services/cursor-viz-spa/src/hooks/useUrlState.ts`

**Tasks**:
- [ ] Create DateRangePicker with presets
- [ ] Implement useDateRange hook for state management
- [ ] Create useUrlState hook for URL synchronization
- [ ] Add debounced search input for table
- [ ] Write tests for filter interactions

**Acceptance Criteria**:
- Date range syncs with URL query params
- Preset selection updates all charts
- Custom range selection works
- Search filters developer table

**TDD Approach**:
1. RED: Write test for useDateRange hook
2. GREEN: Implement hook with initial state
3. REFACTOR: Add URL persistence

---

### TASK06: GraphQL Queries & Custom Hooks

**Estimated**: 4.0h
**Status**: NOT_STARTED
**Actual**: -

**Objective**: Define all GraphQL queries and create custom hooks.

**Files**:
- `services/cursor-viz-spa/src/hooks/useDashboard.ts`
- `services/cursor-viz-spa/src/hooks/useDevelopers.ts`
- `services/cursor-viz-spa/src/hooks/useTeamStats.ts`
- `services/cursor-viz-spa/src/graphql/queries.ts`

**Tasks**:
- [ ] Define GET_DASHBOARD_SUMMARY query
- [ ] Define GET_DEVELOPERS query with pagination
- [ ] Define GET_TEAM_STATS query
- [ ] Create custom hooks wrapping React Query
- [ ] Configure retry and cache policies
- [ ] Write tests for hooks with mocked data

**Acceptance Criteria**:
- All queries are type-safe
- Hooks handle loading, error, and success states
- Cache invalidation works correctly
- Retry logic prevents excessive requests

**TDD Approach**:
1. RED: Write test for useDashboard hook
2. GREEN: Create hook with query
3. REFACTOR: Add error handling and retries

---

### TASK07: Testing Setup & Initial Tests

**Estimated**: 2.5h
**Status**: NOT_STARTED
**Actual**: -

**Objective**: Complete test infrastructure and write initial test suites.

**Files**:
- `services/cursor-viz-spa/src/test/setup.ts`
- `services/cursor-viz-spa/src/test/mocks/handlers.ts`
- `services/cursor-viz-spa/src/test/utils.tsx`
- `services/cursor-viz-spa/src/__tests__/Dashboard.test.tsx`

**Tasks**:
- [ ] Set up MSW for GraphQL mocking
- [ ] Create test utilities and render helpers
- [ ] Write integration test for Dashboard page
- [ ] Write unit tests for all hooks
- [ ] Achieve 80% code coverage
- [ ] Set up coverage reporting

**Acceptance Criteria**:
- All tests pass with `npm run test`
- Coverage report shows >= 80%
- Integration tests verify full page rendering
- MSW mocks all GraphQL endpoints

**TDD Approach**:
1. RED: Write test for Dashboard with mock data
2. GREEN: Make test pass with minimal implementation
3. REFACTOR: Add edge cases and error scenarios

---

**Notes**:
- Full implementation depends on P5 (cursor-analytics-core) GraphQL schema
- Can proceed with mock data for initial setup
- Real integration testing will occur after P5 is deployed
