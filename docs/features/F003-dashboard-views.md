# Feature F003: Dashboard Visualizations

**Feature ID:** F003  
**Service:** cursor-viz-spa  
**Priority:** P0 (Critical Path)  
**Status:** Specification Complete

---

## 1. Overview

The Dashboard provides interactive visualizations that transform aggregated usage data into actionable insights. It presents multiple views including a velocity heatmap, team comparison radar charts, and developer efficiency tables, all powered by real-time GraphQL queries.

### 1.1 Business Value

Data without visualization remains underutilized. The dashboard translates complex metrics into intuitive visual representations that enable engineering managers to quickly identify adoption patterns, compare team performance, and spot developers who may benefit from additional AI tooling training.

### 1.2 Success Criteria

The feature is complete when the dashboard loads within 2 seconds on a standard development machine, when all charts accurately reflect the underlying data without visual distortion, when the interface is responsive and usable on screens from 1280px to 4K resolution, and when users can navigate between views without page reloads.

---

## 2. Functional Requirements

### 2.1 Dashboard Layout (FR-VIZ-001)

The dashboard must present a cohesive layout that surfaces key metrics at a glance while providing drill-down capabilities for detailed analysis.

**Layout Structure:**

```
┌──────────────────────────────────────────────────────────────────────┐
│  HEADER: Logo | Navigation | Date Range Selector | Refresh Status   │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐        │
│  │ Active     │ │ Avg Accept │ │ AI Lines   │ │ Top        │        │
│  │ Developers │ │ Rate       │ │ This Week  │ │ Performer  │        │
│  │     42     │ │   72.3%    │ │   15,234   │ │ Alice S.   │        │
│  └────────────┘ └────────────┘ └────────────┘ └────────────┘        │
│                                                                      │
│  KPI CARDS ROW                                                       │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────────────────────────┐  ┌────────────────────────┐│
│  │                                     │  │                        ││
│  │         VELOCITY HEATMAP            │  │    TEAM RADAR CHART    ││
│  │    (GitHub-style contribution)      │  │   (Multi-axis compare) ││
│  │                                     │  │                        ││
│  │  Jan Feb Mar Apr May Jun Jul Aug    │  │     Backend            ││
│  │  ░░▓▓▓░░▓▓▓▓▓░▓▓░░▓▓▓░░▓▓▓▓▓░░▓▓   │  │    /    \              ││
│  │  ▓▓▓░▓▓▓▓░▓▓▓▓▓▓▓▓░░▓▓▓▓▓░▓▓▓▓▓▓   │  │   /      \             ││
│  │                                     │  │  Frontend---Platform   ││
│  └─────────────────────────────────────┘  └────────────────────────┘│
│                                                                      │
│  MAIN CHART AREA                                                     │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────────┐│
│  │  DEVELOPER EFFICIENCY TABLE                                      ││
│  │                                                                  ││
│  │  Developer     │ Team     │ AI Lines │ Accept Rate │ Trend      ││
│  │  ─────────────────────────────────────────────────────────────  ││
│  │  Alice Smith   │ Backend  │  1,234   │   92.3%     │  ↑ +5%     ││
│  │  Bob Johnson   │ Frontend │    987   │   78.1%     │  → 0%      ││
│  │  Carol White   │ Platform │    456   │   15.2%     │  ↓ -8%  ⚠ ││
│  │                                                                  ││
│  └──────────────────────────────────────────────────────────────────┘│
│                                                                      │
│  DATA TABLE AREA                                                     │
└──────────────────────────────────────────────────────────────────────┘
```

**Responsive Breakpoints:**
- Desktop (≥1280px): Full 3-column layout
- Tablet (768-1279px): 2-column layout, stacked cards
- Mobile (≤767px): Single column, simplified charts

**Acceptance Criteria:**
- AC1: Layout renders correctly at all breakpoints without horizontal scroll
- AC2: KPI cards display current values with appropriate formatting
- AC3: Charts resize proportionally when container changes
- AC4: Navigation between views does not cause full page reload
- AC5: Loading states show skeleton UI, not blank space

### 2.2 KPI Summary Cards (FR-VIZ-002)

The dashboard header must display four key performance indicators in card format, providing an at-a-glance summary of system-wide metrics.

**Card Specifications:**

| Card | Metric | Format | Color Logic |
|------|--------|--------|-------------|
| Active Developers | Count of developers with activity in range | Integer | Green if > 80% of total |
| Avg Acceptance Rate | Mean acceptance rate across all active devs | Percentage (1 decimal) | Green > 70%, Yellow 50-70%, Red < 50% |
| AI Lines This Week | Sum of accepted AI lines in last 7 days | Number with K/M suffix | Always neutral |
| Top Performer | Developer with highest acceptance rate | Name + small avatar | Always neutral |

**Data Binding:**
```typescript
// hooks/useDashboardKPIs.ts
const GET_KPI_SUMMARY = gql`
  query GetDashboardSummary($range: DateRange!) {
    getDashboardSummary(range: $range) {
      totalDevelopers
      activeDevelopers
      avgAcceptanceRate
      totalAiLinesThisWeek
      topPerformer {
        id
        name
        stats { acceptanceRate }
      }
    }
  }
`;
```

**Acceptance Criteria:**
- AC1: Cards update within 5 seconds of new data being available
- AC2: Large numbers display with appropriate suffixes (1,234 → 1.2K)
- AC3: Percentage values show exactly one decimal place
- AC4: Top performer name truncates with ellipsis if > 20 characters
- AC5: Color transitions smoothly when thresholds are crossed

### 2.3 Velocity Heatmap (FR-VIZ-003)

The heatmap provides a GitHub-style contribution graph showing AI code acceptance intensity over time. Each cell represents one day, with color intensity reflecting the volume of accepted suggestions.

**Visual Specification:**

```
┌─────────────────────────────────────────────────────────────┐
│                  Velocity Heatmap                           │
│                                                             │
│     Jan      Feb      Mar      Apr      May      Jun       │
│  M  ░░▓▓▓░░▓▓▓▓▓░░▓▓▓░░▓▓▓░░▓▓▓▓▓░░▓▓▓░░▓▓▓░░▓▓▓▓▓░       │
│  T  ▓▓▓░▓▓▓▓░▓▓▓▓▓▓▓░▓▓▓░▓▓▓▓░▓▓▓▓▓▓▓░▓▓▓░▓▓▓▓░▓▓▓▓       │
│  W  ░▓▓▓▓▓░▓▓▓▓░▓▓▓▓▓░▓▓▓▓▓░▓▓▓▓░▓▓▓▓▓░▓▓▓▓▓░▓▓▓▓░▓       │
│  T  ▓▓░▓▓▓▓▓░▓▓▓▓▓▓▓▓░▓▓▓▓▓▓░▓▓▓▓▓▓▓▓░▓▓▓▓▓▓░▓▓▓▓▓▓       │
│  F  ░░░▓▓░░░▓▓▓░░░░░░▓▓░░░▓▓▓░░░░░░▓▓░░░▓▓▓░░░░░░▓▓       │
│  S  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░       │
│  S  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░       │
│                                                             │
│  Legend: ░ No activity  ▒ Low  ▓ Medium  █ High            │
└─────────────────────────────────────────────────────────────┘
```

**Color Scale (Green Theme):**
- Level 0: `#ebedf0` (no activity)
- Level 1: `#9be9a8` (1-25th percentile)
- Level 2: `#40c463` (25-50th percentile)
- Level 3: `#30a14e` (50-75th percentile)
- Level 4: `#216e39` (75-100th percentile)

**Interaction:**
- Hover: Tooltip shows date, exact value, and comparison to average
- Click: Drill down to that day's developer breakdown
- Scroll: Horizontal scroll for data spanning > 6 months

**GraphQL Query:**
```typescript
const GET_HEATMAP = gql`
  query GetVelocityHeatmap($developerId: ID, $teamName: String, $days: Int) {
    getVelocityHeatmap(developerId: $developerId, teamName: $teamName, days: $days) {
      date
      weekday
      week
      intensity
      value
    }
  }
`;
```

**Acceptance Criteria:**
- AC1: Heatmap displays 52 weeks of data by default
- AC2: Color intensity accurately reflects value percentile
- AC3: Weekends are visually distinct (can show less activity naturally)
- AC4: Tooltip appears within 100ms of hover
- AC5: Click drill-down navigates without page reload
- AC6: Empty days render as Level 0, not as missing cells

### 2.4 Team Comparison Radar Chart (FR-VIZ-004)

The radar chart enables comparison of multiple teams across five dimensions simultaneously, revealing relative strengths and improvement areas.

**Dimensions:**

| Axis | Source Metric | Normalization |
|------|---------------|---------------|
| Chat Usage | chatRequests / activeDevelopers | 0-100 scale per max |
| Code Completion | totalTabsAccepted / activeDevelopers | 0-100 scale per max |
| Refactoring Prompts | cmdkUsages / activeDevelopers | 0-100 scale per max |
| Agent Usage | agentRequests / activeDevelopers | 0-100 scale per max |
| Acceptance Rate | avgAcceptanceRate | Already 0-100% |

**Visual Specification:**

```
              Chat Usage
                  │
                  │  100
                 ╱│╲
              ╱   │   ╲
           ╱      │      ╲
Acceptance │──────┼──────│ Code Completion
   Rate    │      │      │
           ╲      │      ╱
              ╲   │   ╱
                 ╲│╱
                  │
           Agent Usage ─── Refactoring
                          Prompts
                          
  ─── Backend Team (blue)
  ─── Frontend Team (orange)
  ─── Platform Team (green)
```

**Interaction:**
- Hover on axis: Highlight that dimension's values for all teams
- Click team legend: Toggle team visibility
- Hover on area: Show exact values for that team

**Acceptance Criteria:**
- AC1: All five axes are evenly distributed (72° apart)
- AC2: Values are normalized so teams can be fairly compared
- AC3: Legend clearly identifies each team's color
- AC4: At least 3 distinct teams can be displayed simultaneously
- AC5: Overlapping areas use transparency (0.3 alpha) for visibility

### 2.5 Developer Efficiency Table (FR-VIZ-005)

The data table provides sortable, filterable access to individual developer metrics with visual indicators for performance outliers.

**Table Columns:**

| Column | Type | Sortable | Description |
|--------|------|----------|-------------|
| Developer | String + Avatar | Yes | Name with team badge |
| Team | String | Yes | Team name |
| AI Lines | Number | Yes | Total accepted AI lines |
| Accept Rate | Percentage | Yes | Suggestion acceptance rate |
| Trend | Indicator | No | Week-over-week change |
| Status | Icon | No | Warning if rate < 20% |

**Row Styling:**
- Normal: Default styling
- Low Performer: Red left border, red background tint (acceptance < 20%)
- Top Performer: Green left border (top 10%)

**Filtering:**
- Team dropdown filter
- Search by developer name
- Toggle to hide low performers

**Pagination:**
- 25 rows per page default
- Show total count
- Keyboard navigation support

**GraphQL Query:**
```typescript
const GET_EFFICIENCY_TABLE = gql`
  query GetDeveloperEfficiencyTable($sortBy: String, $sortOrder: String, $limit: Int) {
    getDeveloperEfficiencyTable(sortBy: $sortBy, sortOrder: $sortOrder, limit: $limit) {
      developer {
        id
        name
        team
      }
      totalAiLines
      acceptanceRate
      isLowPerformer
      trend
    }
  }
`;
```

**Acceptance Criteria:**
- AC1: Table sorts correctly by all sortable columns
- AC2: Low performers (< 20% acceptance) have clear visual indicator
- AC3: Clicking a row navigates to developer detail view
- AC4: Table handles 500+ rows without performance degradation
- AC5: Empty state shows helpful message, not blank table

### 2.6 Date Range Selection (FR-VIZ-006)

Users must be able to filter all dashboard data by date range, with sensible presets and custom range support.

**Preset Ranges:**
- Today
- Last 7 Days
- Last 30 Days (default)
- Last 90 Days
- This Quarter
- Custom Range (date picker)

**Behavior:**
- Range selection updates all charts simultaneously
- URL reflects current range for shareability
- Range persists in local storage for return visits

**Acceptance Criteria:**
- AC1: Selecting a range triggers data refresh within 1 second
- AC2: Custom date picker enforces max 90-day range
- AC3: URL updates to include range (e.g., `?range=30d`)
- AC4: Invalid date ranges show validation error
- AC5: "Last 7 Days" adapts to current date

### 2.7 Real-Time Updates (FR-VIZ-007)

The dashboard must periodically poll for new data and update visualizations without user intervention.

**Polling Configuration:**
- Default interval: 30 seconds
- Configurable via `VITE_POLL_INTERVAL` environment variable
- Visual indicator shows last update time and next update countdown

**Update Behavior:**
- Charts animate transitions when data changes
- No full re-render; incremental updates only
- Stale data indicator if > 5 minutes since last successful fetch

**Acceptance Criteria:**
- AC1: Dashboard updates automatically every 30 seconds
- AC2: Update indicator shows accurate countdown
- AC3: Failed polls do not cause UI to break
- AC4: User interaction (hover, scroll) does not interrupt polls
- AC5: Multiple rapid polls do not cause race conditions

---

## 3. Non-Functional Requirements

### 3.1 Performance

- Initial load: Time to Interactive under 2 seconds
- Chart render: Under 100ms for datasets < 1000 points
- Bundle size: Under 500KB gzipped for initial load

### 3.2 Accessibility

- WCAG 2.1 AA compliance
- Keyboard navigation for all interactive elements
- Color blind-friendly chart palettes (with patterns)
- Screen reader announcements for data updates

### 3.3 Browser Support

- Chrome (last 2 versions)
- Firefox (last 2 versions)
- Safari (last 2 versions)
- Edge (last 2 versions)

### 3.4 Error States

- GraphQL errors: Toast notification with retry button
- Empty data: Friendly empty state with suggestions
- Partial data: Render what's available, indicate missing

---

## 4. Technical Design Notes

### 4.1 Component Hierarchy

```
<App>
  <ApolloProvider>
    <QueryClientProvider>
      <Router>
        <DashboardLayout>
          <Header>
            <DateRangePicker />
            <UpdateIndicator />
          </Header>
          <KPICardGrid>
            <KPICard metric="activeDevelopers" />
            <KPICard metric="avgAcceptanceRate" />
            <KPICard metric="aiLinesThisWeek" />
            <KPICard metric="topPerformer" />
          </KPICardGrid>
          <ChartGrid>
            <VelocityHeatmap />
            <TeamRadarChart />
          </ChartGrid>
          <DeveloperEfficiencyTable />
        </DashboardLayout>
      </Router>
    </QueryClientProvider>
  </ApolloProvider>
</App>
```

### 4.2 State Management Strategy

The application uses a layered state management approach. Server state such as developer data, statistics, and KPIs is managed via TanStack Query with automatic caching and revalidation. UI state including selected date range, sort preferences, and filters lives in React Context. Derived state is computed using useMemo to prevent unnecessary recalculations.

```typescript
// contexts/DashboardContext.tsx
interface DashboardState {
  dateRange: DateRange;
  selectedTeam: string | null;
  sortConfig: SortConfig;
  highlightedDeveloper: string | null;
}

// Changes to context trigger re-fetch of relevant queries
```

### 4.3 Chart Library Integration

Recharts is used for its React-first design. The library provides composable chart components that integrate naturally with React's component model, support for responsive containers, and animation capabilities.

```typescript
// components/charts/VelocityHeatmap.tsx
import { ResponsiveContainer, ScatterChart, Scatter, XAxis, YAxis, Cell, Tooltip } from 'recharts';

export const VelocityHeatmap: React.FC<HeatmapProps> = ({ data }) => {
  const colorScale = useColorScale();
  
  return (
    <ResponsiveContainer width="100%" height={200}>
      <ScatterChart>
        <XAxis dataKey="week" type="number" />
        <YAxis dataKey="weekday" type="number" />
        <Tooltip content={<HeatmapTooltip />} />
        <Scatter data={data}>
          {data.map((cell, index) => (
            <Cell key={index} fill={colorScale(cell.intensity)} />
          ))}
        </Scatter>
      </ScatterChart>
    </ResponsiveContainer>
  );
};
```

### 4.4 GraphQL Caching Strategy

Apollo Client's normalized cache is configured to optimize for dashboard patterns. Queries that return the same developer from different views will share cached data, reducing network requests.

```typescript
// apollo/cache.ts
const cache = new InMemoryCache({
  typePolicies: {
    Developer: {
      keyFields: ['id'],
    },
    DailyStats: {
      keyFields: ['date', 'developerId'],
    },
    Query: {
      fields: {
        getDashboardSummary: {
          merge: true, // Deep merge updates
        },
      },
    },
  },
});
```

---

## 5. Dependencies

### 5.1 External Libraries

| Library | Version | Purpose |
|---------|---------|---------|
| react | ^18.2 | UI framework |
| vite | ^5.0 | Build tool |
| @apollo/client | ^3.8 | GraphQL client |
| @tanstack/react-query | ^5.0 | Server state management |
| recharts | ^2.10 | Charting library |
| tailwindcss | ^3.4 | Utility CSS framework |
| date-fns | ^3.0 | Date utilities |
| react-router-dom | ^6.0 | Client-side routing |

### 5.2 Development Dependencies

| Library | Version | Purpose |
|---------|---------|---------|
| vitest | ^1.0 | Test runner |
| @testing-library/react | ^14.0 | Component testing |
| msw | ^2.0 | API mocking |
| playwright | ^1.40 | E2E testing |
| storybook | ^8.0 | Component documentation |

---

## 6. Test Cases

### 6.1 Component Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| CT-VIZ-001 | KPICard renders value correctly | Displays formatted value |
| CT-VIZ-002 | KPICard applies color based on threshold | Correct color class applied |
| CT-VIZ-003 | Heatmap renders all cells | Cell count matches data length |
| CT-VIZ-004 | Heatmap tooltip shows on hover | Tooltip visible with correct content |
| CT-VIZ-005 | RadarChart shows all teams | Legend matches team count |
| CT-VIZ-006 | Table sorts on column click | Rows reorder correctly |
| CT-VIZ-007 | Table highlights low performers | Red styling on < 20% rows |

### 6.2 Integration Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| IT-VIZ-001 | Dashboard loads with mock data | All components render |
| IT-VIZ-002 | Date range change updates all charts | New queries fired, data refreshes |
| IT-VIZ-003 | Polling updates data periodically | New values appear without refresh |
| IT-VIZ-004 | Error state shows retry option | Toast notification with button |
| IT-VIZ-005 | Empty state renders correctly | Helpful message displayed |

### 6.3 E2E Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| E2E-VIZ-001 | User views dashboard summary | All KPIs visible and populated |
| E2E-VIZ-002 | User hovers heatmap cell | Tooltip appears with date/value |
| E2E-VIZ-003 | User clicks developer row | Navigates to developer detail |
| E2E-VIZ-004 | User changes date range | All charts update with new data |
| E2E-VIZ-005 | Dashboard auto-updates | New data appears after poll |

---

## 7. Related User Stories

- [US-VIZ-001](../user-stories/US-VIZ-001-view-dashboard.md): View Dashboard Summary
- [US-VIZ-002](../user-stories/US-VIZ-002-explore-heatmap.md): Explore Velocity Heatmap
- [US-VIZ-003](../user-stories/US-VIZ-003-compare-teams.md): Compare Teams
- [US-VIZ-004](../user-stories/US-VIZ-004-identify-outliers.md): Identify Performance Outliers

---

## 8. Implementation Tasks

- [TASK-018](../tasks/TASK-018-viz-project-setup.md): Set up React Vite project
- [TASK-019](../tasks/TASK-019-viz-apollo.md): Configure Apollo Client
- [TASK-020](../tasks/TASK-020-viz-layout.md): Implement dashboard layout
- [TASK-021](../tasks/TASK-021-viz-kpi-cards.md): Implement KPI summary cards
- [TASK-022](../tasks/TASK-022-viz-heatmap.md): Implement velocity heatmap
- [TASK-023](../tasks/TASK-023-viz-radar.md): Implement team radar chart
- [TASK-024](../tasks/TASK-024-viz-table.md): Implement efficiency table
- [TASK-025](../tasks/TASK-025-viz-polling.md): Implement real-time polling
- [TASK-026](../tasks/TASK-026-viz-docker.md): Create Dockerfile
