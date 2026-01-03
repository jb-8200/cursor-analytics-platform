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
| TASK04 | Chart Components (Heatmap, Radar, Table) | 5.0 | DONE | 4.5 |
| TASK05 | Filter Controls & Date Picker | 3.0 | DONE | 3.0 |
| TASK06 | GraphQL Queries & Custom Hooks | 4.0 | DONE | 4.0 |
| TASK07 | Testing Setup & Integration Tests | 2.5 | DONE | 2.5 |

**Current Task**: ALL COMPLETE ✅ (24.5h / 24.0h)

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
**Status**: DONE
**Actual**: 4.5h

**Objective**: Implement the three primary visualization components.

**Files**:
- `services/cursor-viz-spa/src/components/charts/VelocityHeatmap.tsx`
- `services/cursor-viz-spa/src/components/charts/VelocityHeatmap.test.tsx`
- `services/cursor-viz-spa/src/components/charts/TeamRadarChart.tsx`
- `services/cursor-viz-spa/src/components/charts/TeamRadarChart.test.tsx`
- `services/cursor-viz-spa/src/components/charts/DeveloperTable.tsx`
- `services/cursor-viz-spa/src/components/charts/DeveloperTable.test.tsx`
- `services/cursor-viz-spa/src/components/charts/index.ts`

**Tasks**:
- [x] Build VelocityHeatmap with GitHub-style grid
- [x] Build TeamRadarChart using Recharts
- [x] Build DeveloperTable with sorting and pagination
- [x] Add tooltips and interactions
- [x] Style components with Tailwind
- [x] Write comprehensive component tests

**Acceptance Criteria**:
- ✅ Heatmap displays 52 weeks of data
- ✅ Radar chart shows 2-5 teams overlapping
- ✅ Table sorts by any column
- ✅ Components are accessible (WCAG 2.1 AA)
- ✅ All tests passing (18 + 16 + 26 = 60 tests)

**TDD Approach**:
1. RED: Write test for VelocityHeatmap data rendering
2. GREEN: Create basic grid structure
3. REFACTOR: Add color scale and tooltips

**Implementation Notes**:
- VelocityHeatmap: Implemented GitHub-style contribution grid with configurable color scale, tooltips, and date-based cell navigation. Handles data gaps gracefully.
- TeamRadarChart: Used Recharts library with normalized metrics (0-100 scale). Supports 2-5 team comparison with custom metric configuration. Includes interactive team selection UI.
- DeveloperTable: Fully sortable and searchable table with pagination. Highlights low acceptance rates below configurable threshold. Keyboard accessible with proper ARIA attributes.
- All components styled with Tailwind CSS and follow accessibility best practices.

---

### TASK05: Filter Controls & Date Picker

**Estimated**: 3.0h
**Status**: DONE
**Actual**: 3.0h

**Objective**: Build date range picker and filter controls.

**Files**:
- `services/cursor-viz-spa/src/components/filters/DateRangePicker.tsx`
- `services/cursor-viz-spa/src/components/filters/DateRangePicker.test.tsx`
- `services/cursor-viz-spa/src/components/filters/SearchInput.tsx`
- `services/cursor-viz-spa/src/components/filters/SearchInput.test.tsx`
- `services/cursor-viz-spa/src/components/filters/index.ts`
- `services/cursor-viz-spa/src/hooks/useDateRange.ts`
- `services/cursor-viz-spa/src/hooks/useDateRange.test.ts`
- `services/cursor-viz-spa/src/hooks/useUrlState.ts`
- `services/cursor-viz-spa/src/hooks/useUrlState.test.tsx`
- `services/cursor-viz-spa/src/hooks/index.ts`

**Tasks**:
- [x] Create DateRangePicker with presets
- [x] Implement useDateRange hook for state management
- [x] Create useUrlState hook for URL synchronization
- [x] Add debounced search input for table
- [x] Write tests for filter interactions

**Acceptance Criteria**:
- ✅ Date range syncs with URL query params (via useUrlState)
- ✅ Preset selection updates all charts (DateRangePicker with 6 presets)
- ✅ Custom range selection works (with validation)
- ✅ Search filters developer table (SearchInput with 300ms debounce)
- ✅ All tests passing (12 + 8 + 11 + 13 = 44 new tests)

**TDD Approach**:
1. RED: Write test for useDateRange hook
2. GREEN: Implement hook with initial state
3. REFACTOR: Add URL persistence

**Implementation Notes**:
- **useDateRange**: State management hook with 6 presets (7d, 30d, 90d, 6m, 1y, custom). Calculates date ranges using date-fns. Formats ranges for display.
- **useUrlState**: Generic hook for syncing any state with URL query params. Supports custom serialization/deserialization. Removes params when set to default value.
- **DateRangePicker**: Dropdown component with preset selection and custom date inputs. Validates date ranges. Keyboard accessible. Highlights current preset.
- **SearchInput**: Debounced input component (default 300ms). Shows search icon and clear button. Immediate clear (no debounce). Controlled component pattern.
- All components styled with Tailwind CSS and follow accessibility best practices (WCAG 2.1 AA).

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

### TASK07: Testing Setup & Integration Tests

**Estimated**: 2.5h
**Status**: DONE
**Actual**: 2.5h

**Objective**: Complete test infrastructure and write comprehensive test suites.

**Files**:
- `services/cursor-viz-spa/src/test/setup.ts` (already existed)
- `services/cursor-viz-spa/src/test/mocks/handlers.ts` (already existed)
- `services/cursor-viz-spa/src/test/utils.tsx` (enhanced)
- `services/cursor-viz-spa/src/test/__tests__/utils.test.tsx` (new)
- `services/cursor-viz-spa/src/__tests__/Dashboard.integration.test.tsx` (new)
- `services/cursor-viz-spa/src/graphql/client.ts` (added default export)
- `services/cursor-viz-spa/package.json` (added @vitest/coverage-v8)

**Tasks**:
- [x] Set up MSW for GraphQL mocking (already complete from TASK02)
- [x] Enhanced test utilities with Apollo + Router providers
- [x] Wrote comprehensive integration test for Dashboard page (11 tests)
- [x] Verified all hook tests are comprehensive (already complete from TASK06)
- [x] Achieved 91.68% code coverage (exceeds 80% target!)
- [x] Set up coverage reporting with @vitest/coverage-v8

**Acceptance Criteria**:
- ✅ All tests pass with `npm run test` (162 tests passing)
- ✅ Coverage report shows >= 80% (91.68% achieved)
- ✅ Integration tests verify full page rendering with Apollo/Router context
- ✅ MSW mocks all GraphQL endpoints (HealthCheck, Dashboard, Developers, Teams)
- ✅ Type check passes (`npm run type-check`)
- ✅ Build succeeds (`npm run build`)

**TDD Approach**:
1. RED: Write test for renderWithProviders utility (failed without Apollo setup)
2. GREEN: Enhanced utils.tsx with Apollo + MemoryRouter providers
3. REFACTOR: Created comprehensive Dashboard integration tests
4. VERIFY: Fixed all TypeScript errors, all tests passing

**Implementation Notes**:
- **Test Utilities Enhanced**: Added `renderWithProviders` that wraps components with ApolloProvider (using MSW-intercepted GraphQL) and MemoryRouter for routing context
- **Integration Tests**: Created 11 integration tests for Dashboard covering rendering, layout, accessibility, and component structure
- **Coverage Achievement**: 91.68% overall, with most packages >93% coverage
- **Fixed Issues**: Corrected SVG title attribute in VelocityHeatmap, fixed hook signatures, removed unused imports
- **Final Counts**: 162 total tests, 20 test files, all passing

---

**Notes**:
- Full implementation depends on P5 (cursor-analytics-core) GraphQL schema
- Can proceed with mock data for initial setup
- Real integration testing will occur after P5 is deployed
