# SPEC.md - Frontend Dashboard (cursor-viz-spa)

**Service**: cursor-viz-spa  
**Type**: Single Page Application  
**Language**: TypeScript  
**Framework**: React 18+ with Vite  
**State Management**: TanStack Query (React Query)  
**Visualization**: Recharts  
**Port**: 3000 (development)  

## Overview

The Frontend Dashboard is a React-based single page application that visualizes AI coding assistant usage analytics. It consumes the GraphQL API provided by cursor-analytics-core and presents the data through interactive charts, tables, and KPI displays.

## Purpose

This service transforms raw analytics data into actionable visual insights. Engineering managers and team leads use it to understand AI tool adoption patterns, identify training opportunities, and measure the productivity impact of AI assistance across their teams.

## User Interface Design

### Layout Structure

The dashboard uses a responsive grid layout with three primary regions.

**Header Region**: Displays the organization name, date range selector, and summary KPI cards. KPIs include total developers, active developers, overall acceptance rate, and today's suggestion counts.

**Main Content Region**: Contains the primary visualizations arranged in a responsive grid. The velocity heatmap occupies the full width on larger screens. Below it, the team radar chart and developer efficiency table sit side by side on desktop and stack vertically on mobile.

**Sidebar Region** (optional): Provides navigation to detailed views and filter controls. Collapses to a hamburger menu on mobile viewports.

### Responsive Breakpoints

| Breakpoint | Width | Layout |
|------------|-------|--------|
| Mobile | < 768px | Single column, stacked components |
| Tablet | 768px - 1024px | Two-column grid, some stacking |
| Desktop | > 1024px | Full grid layout with sidebar |

## Components

### VelocityHeatmap

A GitHub-style contribution graph showing AI code acceptance intensity over time.

**Props**:
```typescript
interface VelocityHeatmapProps {
    data: DailyStats[];           // Array of daily statistics
    weeks?: number;               // Number of weeks to display (default: 52)
    colorScale?: string[];        // Custom color gradient
    onCellClick?: (date: Date) => void;
}
```

**Behavior**:
- Renders a grid with 7 rows (days of week) and N columns (weeks)
- Color intensity maps to `suggestionsAccepted` count
- Tooltips show date and exact count on hover
- Day labels appear on left edge (Mon, Wed, Fri)
- Month labels appear above at month boundaries
- Clicking a cell can trigger navigation to detailed view

**Color Scale**:
```typescript
const defaultColorScale = [
    '#ebedf0',  // 0 (no activity)
    '#9be9a8',  // 1-25th percentile
    '#40c463',  // 25-50th percentile
    '#30a14e',  // 50-75th percentile
    '#216e39',  // 75-100th percentile
];
```

### TeamRadarChart

A multi-axis radar chart comparing teams across different metrics.

**Props**:
```typescript
interface TeamRadarChartProps {
    data: TeamStats[];
    selectedTeams: string[];       // Team names to display
    onTeamSelect: (teams: string[]) => void;
    metrics?: MetricConfig[];      // Custom metric configuration
}

interface MetricConfig {
    key: keyof TeamStats;
    label: string;
    max?: number;                  // Max value for normalization
}
```

**Default Metrics**:
- Chat Usage (chatInteractions)
- Code Completion (totalSuggestions)
- Acceptance Rate (averageAcceptanceRate)
- AI Velocity (aiVelocity)
- Cmd+K Usage (derived from detailed stats)

**Behavior**:
- Displays 2-5 teams simultaneously as overlapping polygons
- Each team has a distinct color with semi-transparent fill
- Axes are normalized to 0-100 scale for comparability
- Legend identifies team colors
- Team selector dropdown enables choosing which teams to compare

### DeveloperTable

A sortable, filterable table displaying individual developer metrics.

**Props**:
```typescript
interface DeveloperTableProps {
    data: Developer[];
    onSort: (column: string, direction: 'asc' | 'desc') => void;
    onSearch: (term: string) => void;
    onPageChange: (page: number) => void;
    pageSize?: number;             // Default: 25
    highlightThreshold?: number;   // Acceptance rate below this is highlighted (default: 20)
}
```

**Columns**:
| Column | Field | Sortable | Description |
|--------|-------|----------|-------------|
| Name | name | Yes | Developer full name |
| Team | team | Yes | Team assignment |
| Suggestions | stats.totalSuggestions | Yes | Total shown |
| Accepted | stats.acceptedSuggestions | Yes | Total accepted |
| Rate | stats.acceptanceRate | Yes | Acceptance percentage |
| AI Lines | stats.aiLinesAdded | Yes | Lines from AI |

**Behavior**:
- Click column header to sort (toggle direction on re-click)
- Search input filters by developer name (debounced 300ms)
- Rows with acceptance rate below threshold show warning background
- Pagination controls show page numbers and prev/next buttons
- Row click navigates to developer detail view

### DateRangePicker

A dropdown selector for controlling the time range of all dashboard data.

**Props**:
```typescript
interface DateRangePickerProps {
    value: DateRange;
    onChange: (range: DateRange) => void;
    presets?: DateRangePreset[];
}

interface DateRange {
    from: Date;
    to: Date;
    preset?: string;
}
```

**Presets**:
- Today
- This Week (Monday to today)
- This Month (1st to today)
- Last 7 Days
- Last 30 Days
- Last 90 Days
- Custom

**Behavior**:
- Dropdown shows preset options
- Selecting "Custom" opens a calendar for start/end selection
- Selected range is displayed in human-readable format
- Range is persisted in URL query parameters
- Changing range triggers data refetch via React Query

### LoadingState

Skeleton loaders that maintain layout while data loads.

**Variants**:
- `HeatmapSkeleton`: Grid of pulsing rectangles matching heatmap shape
- `RadarSkeleton`: Circular pulsing shape
- `TableSkeleton`: Rows of pulsing rectangles for table
- `KPISkeleton`: Card-shaped pulsing rectangles

### ErrorState

Consistent error display with retry capability.

**Props**:
```typescript
interface ErrorStateProps {
    title?: string;
    message: string;
    onRetry?: () => void;
    showDetails?: boolean;
    error?: Error;
}
```

**Behavior**:
- Displays user-friendly error message
- Shows "Retry" button when `onRetry` provided
- Optionally shows technical details for debugging
- Logs error to console in development, to error service in production

## State Management

### React Query Configuration

All server state is managed through TanStack Query (React Query):

```typescript
const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            staleTime: 30 * 1000,        // 30 seconds
            cacheTime: 5 * 60 * 1000,    // 5 minutes
            refetchOnWindowFocus: true,
            retry: 3,
            retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
        },
    },
});
```

### Query Hooks

Custom hooks encapsulate GraphQL queries:

```typescript
// hooks/useDashboard.ts
export function useDashboard(range?: DateRange) {
    return useQuery({
        queryKey: ['dashboard', range?.from, range?.to],
        queryFn: () => fetchDashboardSummary(range),
    });
}

// hooks/useDevelopers.ts
export function useDevelopers(options: DeveloperQueryOptions) {
    return useQuery({
        queryKey: ['developers', options],
        queryFn: () => fetchDevelopers(options),
    });
}

// hooks/useTeamStats.ts
export function useTeamStats(teamName: string) {
    return useQuery({
        queryKey: ['team', teamName],
        queryFn: () => fetchTeamStats(teamName),
    });
}
```

### URL State

The date range selection is synchronized with URL query parameters for shareability:

```typescript
// URL: /dashboard?from=2026-01-01&to=2026-01-15
const [searchParams, setSearchParams] = useSearchParams();

const dateRange = useMemo(() => ({
    from: searchParams.get('from') || defaultFrom,
    to: searchParams.get('to') || defaultTo,
}), [searchParams]);
```

## GraphQL Integration

### Apollo Client Setup

```typescript
const client = new ApolloClient({
    uri: import.meta.env.VITE_GRAPHQL_URL,
    cache: new InMemoryCache({
        typePolicies: {
            Developer: {
                keyFields: ['id'],
            },
            DailyStats: {
                keyFields: ['date'],
            },
        },
    }),
});
```

### Query Definitions

```typescript
// graphql/queries.ts
export const GET_DASHBOARD_SUMMARY = gql`
    query GetDashboardSummary($range: DateRangeInput) {
        dashboardSummary(range: $range) {
            totalDevelopers
            activeDevelopers
            overallAcceptanceRate
            totalSuggestionsToday
            totalAcceptedToday
            aiVelocityToday
            teamComparison {
                teamName
                memberCount
                averageAcceptanceRate
                totalSuggestions
                aiVelocity
            }
            dailyTrend {
                date
                suggestionsShown
                suggestionsAccepted
                acceptanceRate
            }
        }
    }
`;

export const GET_DEVELOPERS = gql`
    query GetDevelopers($team: String, $limit: Int, $offset: Int) {
        developers(team: $team, limit: $limit, offset: $offset) {
            nodes {
                id
                name
                email
                team
                seniority
                stats {
                    totalSuggestions
                    acceptedSuggestions
                    acceptanceRate
                    aiLinesAdded
                }
            }
            totalCount
            pageInfo {
                hasNextPage
                hasPreviousPage
            }
        }
    }
`;
```

### Type Generation

GraphQL types are generated from the schema using GraphQL Code Generator:

```yaml
# codegen.yml
schema: "http://localhost:4000/graphql"
documents: "src/graphql/**/*.ts"
generates:
  src/generated/graphql.ts:
    plugins:
      - typescript
      - typescript-operations
      - typescript-react-apollo
```

## Routing

The application uses React Router for navigation:

```typescript
const routes = [
    { path: '/', element: <Navigate to="/dashboard" /> },
    { path: '/dashboard', element: <Dashboard /> },
    { path: '/teams', element: <TeamList /> },
    { path: '/teams/:teamName', element: <TeamDetail /> },
    { path: '/developers', element: <DeveloperList /> },
    { path: '/developers/:id', element: <DeveloperDetail /> },
];
```

## Styling

### Tailwind CSS Configuration

The project uses Tailwind CSS with a custom theme:

```javascript
// tailwind.config.js
module.exports = {
    theme: {
        extend: {
            colors: {
                primary: {
                    50: '#f0f9ff',
                    500: '#0ea5e9',
                    900: '#0c4a6e',
                },
                success: '#22c55e',
                warning: '#f59e0b',
                danger: '#ef4444',
            },
        },
    },
};
```

### Component Styling Conventions

- Use Tailwind utility classes for most styling
- Extract repeated patterns to component classes
- Use CSS modules for complex component-specific styles
- Maintain consistent spacing using Tailwind's spacing scale

## Accessibility

The dashboard follows WCAG 2.1 AA guidelines:

- All interactive elements are keyboard accessible
- Color is not the only indicator of information
- Charts include text alternatives in tooltips
- Focus indicators are visible
- Proper heading hierarchy is maintained
- ARIA labels on interactive elements

## Performance

### Bundle Optimization

- Code splitting by route using React.lazy()
- Tree shaking enabled via Vite
- Chart libraries loaded dynamically
- Images optimized and lazy-loaded

### Render Optimization

- Memoization of expensive computations
- Virtualized lists for large datasets
- Debounced search inputs
- Skeleton loaders for perceived performance

### Target Metrics

| Metric | Target |
|--------|--------|
| First Contentful Paint | < 1.5s |
| Largest Contentful Paint | < 2.5s |
| Time to Interactive | < 3.5s |
| Cumulative Layout Shift | < 0.1 |

## Configuration

### Environment Variables

```bash
VITE_GRAPHQL_URL=http://localhost:4000/graphql
VITE_APP_TITLE=Cursor Analytics Dashboard
VITE_DEFAULT_DATE_RANGE=LAST_30_DAYS
```

### Build Configuration

```typescript
// vite.config.ts
export default defineConfig({
    plugins: [react()],
    build: {
        target: 'es2020',
        sourcemap: true,
        rollupOptions: {
            output: {
                manualChunks: {
                    vendor: ['react', 'react-dom'],
                    charts: ['recharts'],
                    apollo: ['@apollo/client'],
                },
            },
        },
    },
});
```

## Testing Requirements

Unit tests must cover:
- [ ] All custom hooks with mocked queries
- [ ] Component rendering with various props
- [ ] User interactions (clicks, typing)
- [ ] Error state display
- [ ] Loading state display

Integration tests must cover:
- [ ] Full page rendering with mocked API
- [ ] Navigation between routes
- [ ] Date range changes and data refetch
- [ ] Table sorting and filtering

E2E tests must cover:
- [ ] Dashboard loads with data
- [ ] Date range selection updates charts
- [ ] Search filters developer table
- [ ] Navigation to detail pages

## Dependencies

Production:
- `react` + `react-dom` - UI framework
- `react-router-dom` - Routing
- `@apollo/client` - GraphQL client
- `@tanstack/react-query` - Server state (alternative to Apollo)
- `recharts` - Charting library
- `date-fns` - Date manipulation
- `tailwindcss` - Styling

Development:
- `typescript` - Type checking
- `vite` - Build tool
- `vitest` - Test runner
- `@testing-library/react` - Component testing
- `msw` - API mocking
- `@graphql-codegen/cli` - Type generation
