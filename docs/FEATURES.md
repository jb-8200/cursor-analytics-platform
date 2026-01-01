# Feature Breakdown: Cursor Usage Analytics Platform

**Version**: 1.0.0  
**Last Updated**: January 2026  

This document provides a comprehensive breakdown of all features across the three services, organized by priority and implementation phase. Each feature includes its scope, dependencies, and acceptance criteria summary.

## Feature Overview Matrix

The following matrix shows all features across services with their implementation priority and estimated complexity.

| Feature ID | Service | Feature Name | Priority | Complexity | Phase |
|------------|---------|--------------|----------|------------|-------|
| SIM-001 | cursor-sim | Developer Profile Generation | P0 | Medium | 1 |
| SIM-002 | cursor-sim | Event Generation Engine | P0 | High | 1 |
| SIM-003 | cursor-sim | REST API Server | P0 | Low | 1 |
| SIM-004 | cursor-sim | CLI Configuration | P0 | Low | 1 |
| SIM-005 | cursor-sim | Team Structure Simulation | P1 | Medium | 2 |
| SIM-006 | cursor-sim | Realistic Time Patterns | P2 | Medium | 3 |
| CORE-001 | cursor-analytics-core | Data Ingestion Worker | P0 | Medium | 1 |
| CORE-002 | cursor-analytics-core | Database Schema & Migrations | P0 | Medium | 1 |
| CORE-003 | cursor-analytics-core | GraphQL API Server | P0 | Medium | 1 |
| CORE-004 | cursor-analytics-core | Metric Calculations | P0 | High | 1 |
| CORE-005 | cursor-analytics-core | Developer Queries | P0 | Medium | 1 |
| CORE-006 | cursor-analytics-core | Team Aggregations | P1 | Medium | 2 |
| CORE-007 | cursor-analytics-core | Dashboard KPIs | P1 | Medium | 2 |
| CORE-008 | cursor-analytics-core | Time Range Filtering | P1 | Low | 2 |
| VIZ-001 | cursor-viz-spa | Dashboard Layout | P0 | Medium | 1 |
| VIZ-002 | cursor-viz-spa | Velocity Heatmap | P0 | High | 1 |
| VIZ-003 | cursor-viz-spa | Team Radar Chart | P1 | High | 2 |
| VIZ-004 | cursor-viz-spa | Developer Efficiency Table | P0 | Medium | 1 |
| VIZ-005 | cursor-viz-spa | Date Range Picker | P1 | Low | 2 |
| VIZ-006 | cursor-viz-spa | Loading States | P1 | Low | 2 |
| VIZ-007 | cursor-viz-spa | Error Handling | P1 | Low | 2 |

## Service A: cursor-sim Features

### SIM-001: Developer Profile Generation

**Priority**: P0 (Critical)  
**Complexity**: Medium  
**Phase**: 1  

This feature generates realistic developer profiles with varying characteristics that influence their simulated behavior patterns.

**Scope:**
The feature creates developer profiles with unique identifiers, names, email addresses, team assignments, and seniority levels. Each profile includes a calculated acceptance rate based on seniority, reflecting the observation that more experienced developers tend to accept a higher percentage of AI suggestions.

**Technical Requirements:**
The system should support generating between 1 and 1000 developer profiles on startup. Names and emails should be generated using a deterministic algorithm based on a seed value, ensuring reproducibility. The seniority distribution should default to 20% junior, 50% mid-level, and 30% senior developers, with acceptance rates of 60%, 75%, and 90% respectively.

**Dependencies:**
None - this is a foundational feature.

**Acceptance Criteria Summary:**
Profiles are generated on startup based on CLI parameters. Each profile has a unique identifier that persists for the session duration. Profiles are retrievable via the REST API endpoints.

### SIM-002: Event Generation Engine

**Priority**: P0 (Critical)  
**Complexity**: High  
**Phase**: 1  

This feature implements the core simulation logic that generates usage events following statistical distributions that mimic real developer behavior.

**Scope:**
The engine generates four event types (cpp_suggestion_shown, cpp_suggestion_accepted, chat_message, cmd_k_prompt) at rates determined by the velocity parameter. Events follow a Poisson distribution to create natural clustering rather than uniform spacing. The fluctuation parameter adds per-developer variance so not all developers exhibit identical patterns.

**Technical Requirements:**
Events should be generated continuously in background goroutines, one per simulated developer. The base event rate should be configurable through the velocity parameter, where "low" generates approximately 10 events per hour and "high" generates approximately 100 events per hour. Each developer's actual rate should vary by plus or minus the fluctuation percentage. Suggestion acceptance events should only be generated following suggestion shown events, with the probability determined by the developer's acceptance rate.

**Dependencies:**
SIM-001 (Developer Profile Generation)

**Acceptance Criteria Summary:**
Events are generated continuously while the server runs. Event timing follows a Poisson distribution. Acceptance rates match developer profiles within statistical variance.

### SIM-003: REST API Server

**Priority**: P0 (Critical)  
**Complexity**: Low  
**Phase**: 1  

This feature exposes the simulated data through REST endpoints that mirror the Cursor Business API contract.

**Scope:**
The server implements endpoints for listing developers, retrieving individual developer details, and querying activity events with time range filters. All responses use JSON format with schemas compatible with the real Cursor API.

**Technical Requirements:**
The server should listen on a configurable port (default 8080). Endpoints should support CORS for local development. The activity endpoint should accept "from" and "to" query parameters in ISO 8601 format. Response pagination should be implemented for the activity endpoint when results exceed 1000 events.

**Dependencies:**
SIM-001, SIM-002

**Acceptance Criteria Summary:**
All endpoints respond with valid JSON matching the documented schema. The health check endpoint returns within 100ms. Time range filtering correctly bounds returned events.

### SIM-004: CLI Configuration

**Priority**: P0 (Critical)  
**Complexity**: Low  
**Phase**: 1  

This feature provides command-line interface options for configuring the simulator's behavior.

**Scope:**
The CLI accepts flags for port, number of developers, velocity level, and fluctuation percentage. Help text documents all available options. Invalid configurations are rejected with helpful error messages.

**Technical Requirements:**
Flags should use the standard Go flag package. The port flag should validate that the value is between 1024 and 65535. The developers flag should accept values between 1 and 1000. The velocity flag should accept "low", "medium", or "high" as values. The fluctuation flag should accept decimal values between 0 and 1.

**Dependencies:**
None

**Acceptance Criteria Summary:**
Running with --help displays all options. Invalid flag values produce descriptive error messages. Defaults are applied when flags are omitted.

### SIM-005: Team Structure Simulation

**Priority**: P1 (Important)  
**Complexity**: Medium  
**Phase**: 2  

This feature adds organizational structure to the simulation by grouping developers into teams with distinct characteristics.

**Scope:**
The feature creates between 3 and 10 teams (configurable), distributing developers across teams. Each team can have different baseline characteristics affecting the average behavior of its members. Team names can be provided via a configuration file or generated automatically.

**Technical Requirements:**
Teams should be approximately equal in size with some variance. A configuration file format should allow specifying team names and their baseline acceptance rate modifiers. The API should support filtering developers by team and aggregating statistics at the team level.

**Dependencies:**
SIM-001, SIM-003

**Acceptance Criteria Summary:**
Developers are distributed across teams on startup. Team-level aggregations are available via the API. Custom team configurations are loaded from file when provided.

### SIM-006: Realistic Time Patterns

**Priority**: P2 (Nice to Have)  
**Complexity**: Medium  
**Phase**: 3  

This feature adds time-based patterns to the simulation, reflecting that developers have different activity levels during working hours versus nights and weekends.

**Scope:**
The feature modifies event generation rates based on simulated time of day and day of week. Configurable parameters control the work day start and end times, weekend reduction factor, and timezone offset.

**Technical Requirements:**
Event rates should be reduced by 80% outside of configured work hours. Weekend rates should be further reduced by the weekend factor (default 90% reduction). The simulation should optionally accelerate time to generate weeks of data quickly for demo purposes.

**Dependencies:**
SIM-002

**Acceptance Criteria Summary:**
Event density visibly varies by time of day in the data. Weekend data shows significantly reduced activity. Time acceleration mode generates historical data on demand.

## Service B: cursor-analytics-core Features

### CORE-001: Data Ingestion Worker

**Priority**: P0 (Critical)  
**Complexity**: Medium  
**Phase**: 1  

This feature implements the background process that fetches data from the simulator and stores it in the database.

**Scope:**
The worker polls the simulator's activity endpoint at configurable intervals, transforms the JSON response into database records, and handles deduplication to prevent storing duplicate events on overlapping queries.

**Technical Requirements:**
The polling interval should be configurable via environment variable, defaulting to 60 seconds. The worker should implement retry logic with exponential backoff when the simulator is unavailable. Event deduplication should use the event ID as the unique key. The worker should track the last successfully processed timestamp to optimize subsequent queries.

**Dependencies:**
CORE-002

**Acceptance Criteria Summary:**
Events are fetched and stored automatically when the service starts. Duplicate events are not created on overlapping polls. Service recovers gracefully from simulator outages.

### CORE-002: Database Schema & Migrations

**Priority**: P0 (Critical)  
**Complexity**: Medium  
**Phase**: 1  

This feature defines the PostgreSQL database schema and provides migration tooling for schema evolution.

**Scope:**
The schema includes tables for developers, usage events, and materialized views for aggregated statistics. Migrations are versioned and can be applied forward and rolled back. Initial seed data is provided for development purposes.

**Technical Requirements:**
Migrations should use a versioning tool such as knex, Prisma, or node-pg-migrate. The schema should include appropriate indexes for query performance. Foreign key constraints should ensure referential integrity. A materialized view should pre-aggregate daily statistics for performance.

**Dependencies:**
None

**Acceptance Criteria Summary:**
Migrations run successfully on a fresh database. Schema matches the documented design. Rollback restores the previous schema state.

### CORE-003: GraphQL API Server

**Priority**: P0 (Critical)  
**Complexity**: Medium  
**Phase**: 1  

This feature implements the GraphQL server that exposes analytics data to the frontend.

**Scope:**
The server provides a GraphQL endpoint with queries for developers, teams, and dashboard statistics. The schema is typed and includes resolver implementations for all queries. GraphQL Playground is available in development mode for API exploration.

**Technical Requirements:**
The server should use Apollo Server 4 with Express. The schema should be defined using SDL (Schema Definition Language). Resolvers should use DataLoader to batch database queries and avoid N+1 problems. The server should implement query complexity analysis to prevent expensive queries.

**Dependencies:**
CORE-002

**Acceptance Criteria Summary:**
GraphQL endpoint responds to documented queries. DataLoader batches related queries efficiently. Query complexity limits reject overly expensive operations.

### CORE-004: Metric Calculations

**Priority**: P0 (Critical)  
**Complexity**: High  
**Phase**: 1  

This feature implements the business logic for calculating key performance indicators from raw event data.

**Scope:**
The feature calculates acceptance rate, AI velocity, chat dependency ratio, and composite productivity scores. Calculations can be performed at the developer, team, and organization levels. Time-windowed calculations support day, week, and month granularities.

**Technical Requirements:**
Calculations should be performed in PostgreSQL where possible for efficiency. The acceptance rate formula is (accepted suggestions divided by shown suggestions) multiplied by 100. AI velocity formula is (AI lines added divided by total lines added) multiplied by 100. Null handling should avoid division by zero errors. Results should be cached in the materialized view and refreshed periodically.

**Dependencies:**
CORE-002

**Acceptance Criteria Summary:**
Metrics match expected values for test data sets. Division by zero scenarios return null rather than errors. Cached calculations refresh at configured intervals.

### CORE-005: Developer Queries

**Priority**: P0 (Critical)  
**Complexity**: Medium  
**Phase**: 1  

This feature implements GraphQL queries for retrieving developer information and statistics.

**Scope:**
Queries support fetching a single developer by ID, listing all developers with pagination, and filtering developers by team. Each developer query can include nested statistics with optional time range parameters.

**Technical Requirements:**
The developer query should accept an ID parameter and return null if not found. The developers query should support limit, offset, and team filter parameters. Statistics should be calculated lazily when requested to avoid unnecessary computation.

**Dependencies:**
CORE-003, CORE-004

**Acceptance Criteria Summary:**
Developer queries return expected data shape. Pagination correctly limits results. Team filtering returns only matching developers.

### CORE-006: Team Aggregations

**Priority**: P1 (Important)  
**Complexity**: Medium  
**Phase**: 2  

This feature adds queries for team-level statistics aggregated from individual developer data.

**Scope:**
Queries return team statistics including member count, average acceptance rate, total suggestions, and identification of top performers. Comparison data enables ranking teams against each other.

**Technical Requirements:**
Team statistics should be calculated from developer statistics, not stored separately. The top performer identification should use acceptance rate as the primary metric. Team comparison should return all teams sorted by the specified metric.

**Dependencies:**
CORE-004, CORE-005

**Acceptance Criteria Summary:**
Team statistics accurately aggregate developer data. Top performer is correctly identified. All teams are included in comparison queries.

### CORE-007: Dashboard KPIs

**Priority**: P1 (Important)  
**Complexity**: Medium  
**Phase**: 2  

This feature implements the comprehensive dashboard query that returns all key performance indicators in a single request.

**Scope:**
The dashboard query returns organization-wide statistics, team comparisons, and trend data optimized for the main dashboard view. The query is designed to minimize database round trips while providing all necessary data.

**Technical Requirements:**
The query should return total and active developer counts, overall acceptance rate, today's suggestion counts, team comparison array, and daily trend data for the specified range. Active developers are those with at least one event in the last 7 days. The query should complete within 500ms for organizations with up to 500 developers.

**Dependencies:**
CORE-004, CORE-005, CORE-006

**Acceptance Criteria Summary:**
Dashboard query returns complete KPI data. Active developer count reflects recent activity. Query performance meets latency requirements.

### CORE-008: Time Range Filtering

**Priority**: P1 (Important)  
**Complexity**: Low  
**Phase**: 2  

This feature adds time range support to all statistics queries for historical analysis.

**Scope:**
All statistics queries accept optional from and to parameters that bound the time window for calculations. Preset ranges (today, this week, this month, last 30 days) are supported as shortcuts.

**Technical Requirements:**
Date parameters should be parsed as ISO 8601 timestamps. Missing parameters default to appropriate boundaries (start of epoch for "from", current time for "to"). Preset range names should expand to their corresponding timestamps server-side.

**Dependencies:**
CORE-004

**Acceptance Criteria Summary:**
Time ranges correctly filter event data. Preset ranges map to expected boundaries. Invalid date formats return descriptive errors.

## Service C: cursor-viz-spa Features

### VIZ-001: Dashboard Layout

**Priority**: P0 (Critical)  
**Complexity**: Medium  
**Phase**: 1  

This feature implements the main dashboard page layout with responsive grid arrangement of visualization components.

**Scope:**
The layout includes a header with summary KPIs, a main content area with chart components arranged in a responsive grid, and a sidebar for navigation and filters. The layout adapts to desktop, tablet, and mobile viewport sizes.

**Technical Requirements:**
The layout should use CSS Grid or Flexbox for responsive behavior. Breakpoints should be set at 768px and 1024px for tablet and desktop transitions. The sidebar should collapse to a hamburger menu on mobile. Chart containers should maintain aspect ratios when resizing.

**Dependencies:**
None

**Acceptance Criteria Summary:**
Layout renders correctly at all viewport sizes. Components reflow appropriately at breakpoints. No horizontal scrolling on mobile.

### VIZ-002: Velocity Heatmap

**Priority**: P0 (Critical)  
**Complexity**: High  
**Phase**: 1  

This feature implements the GitHub-style contribution graph showing AI code acceptance intensity over time.

**Scope:**
The heatmap displays a grid of cells representing days, colored by the intensity of AI-accepted code on that day. Tooltips show exact values on hover. The display shows the most recent 52 weeks by default with options to scroll back further.

**Technical Requirements:**
Cells should use a gradient from light to dark based on the acceptance count. The color scale should be configurable. Week labels should appear on the left edge. Month labels should appear above the appropriate week boundaries. The component should accept daily statistics data and calculate the color mapping internally.

**Dependencies:**
VIZ-001, CORE-007

**Acceptance Criteria Summary:**
Heatmap renders with correct date alignment. Colors accurately represent value intensity. Tooltips display on hover.

### VIZ-003: Team Radar Chart

**Priority**: P1 (Important)  
**Complexity**: High  
**Phase**: 2  

This feature implements the multi-axis radar chart for comparing teams across different metrics.

**Scope:**
The radar chart displays multiple teams as overlapping polygons on axes representing different metrics. Selectable teams allow focusing the comparison on specific groups. Axes include Chat Usage, Code Completion, Refactoring Prompts, Acceptance Rate, and AI Velocity.

**Technical Requirements:**
The chart should support displaying 2-5 teams simultaneously. Team selection should use a multi-select dropdown. Axis labels should be positioned outside the chart area. Each team should have a distinct color with semi-transparent fill. The chart should include a legend mapping colors to team names.

**Dependencies:**
VIZ-001, CORE-006

**Acceptance Criteria Summary:**
Radar chart renders with correct axis positions. Team polygons accurately represent metric values. Legend clearly identifies each team.

### VIZ-004: Developer Efficiency Table

**Priority**: P0 (Critical)  
**Complexity**: Medium  
**Phase**: 1  

This feature implements the sortable table displaying individual developer metrics with visual performance indicators.

**Scope:**
The table displays columns for developer name, team, total AI suggestions, accepted suggestions, acceptance rate, and AI lines written. Rows are sortable by any column. Rows with acceptance rates below 20% are highlighted in red to flag potential issues.

**Technical Requirements:**
Sorting should be client-side for tables under 100 rows, server-side for larger datasets. The acceptance rate column should include a visual bar indicator in addition to the numeric value. Pagination should display 25 rows per page by default. A search filter should allow finding specific developers by name.

**Dependencies:**
VIZ-001, CORE-005

**Acceptance Criteria Summary:**
Table sorts correctly by all columns. Low acceptance rate rows are visually distinct. Search filter matches partial names.

### VIZ-005: Date Range Picker

**Priority**: P1 (Important)  
**Complexity**: Low  
**Phase**: 2  

This feature provides a date range selection component that controls the time window for all dashboard visualizations.

**Scope:**
The picker includes preset options (Today, This Week, This Month, Last 30 Days, Custom) and a custom date range selector. Selection updates all dashboard components through shared state.

**Technical Requirements:**
The component should use a dropdown for preset selection. Custom range selection should open a calendar popup allowing start and end date selection. The selected range should be displayed in a human-readable format. Range changes should trigger re-fetching of dashboard data through React Query.

**Dependencies:**
CORE-008

**Acceptance Criteria Summary:**
Preset selections apply correct date boundaries. Custom range selector validates that start is before end. Dashboard updates when range changes.

### VIZ-006: Loading States

**Priority**: P1 (Important)  
**Complexity**: Low  
**Phase**: 2  

This feature implements consistent loading indicators across all dashboard components during data fetching.

**Scope:**
Each visualization component displays a skeleton loader while its data is loading. The skeleton matches the approximate shape of the loaded content to minimize layout shift. A global loading indicator appears in the header during any active fetch.

**Technical Requirements:**
Skeleton loaders should use CSS animation for subtle pulsing effect. The loading state should be determined by React Query's isLoading flag. Components should display cached data with a subtle staleness indicator while revalidating.

**Dependencies:**
VIZ-001

**Acceptance Criteria Summary:**
Skeletons appear during initial data load. Layout does not shift when data arrives. Stale data indicator appears during revalidation.

### VIZ-007: Error Handling

**Priority**: P1 (Important)  
**Complexity**: Low  
**Phase**: 2  

This feature implements consistent error display and recovery mechanisms across the dashboard.

**Scope:**
Components display helpful error messages when queries fail. A retry button allows manual re-fetch. The error boundary prevents individual component failures from crashing the entire dashboard.

**Technical Requirements:**
Error messages should be user-friendly, not raw error strings. The retry button should be prominently displayed within the error state. Error boundaries should log errors to console in development and to an error service in production. Network errors should suggest checking connectivity.

**Dependencies:**
VIZ-001

**Acceptance Criteria Summary:**
Errors display user-friendly messages. Retry successfully re-fetches data. Single component errors do not affect siblings.

## Phase Summary

### Phase 1: Core Functionality (MVP)

Phase 1 delivers a working end-to-end system with basic simulation, data processing, and visualization. Upon completion, users can run the platform locally, see synthetic developer data flowing through the pipeline, and view basic analytics on the dashboard.

Features included: SIM-001 through SIM-004, CORE-001 through CORE-005, VIZ-001, VIZ-002, VIZ-004.

### Phase 2: Enhanced Analytics

Phase 2 adds team-level analytics, improved filtering, and polished user experience. Upon completion, users can compare teams, filter by date range, and experience smooth loading and error states.

Features included: SIM-005, CORE-006 through CORE-008, VIZ-003, VIZ-005 through VIZ-007.

### Phase 3: Advanced Simulation

Phase 3 adds sophisticated simulation features for more realistic demos and testing scenarios. Upon completion, the simulator produces data patterns that closely match real-world usage.

Features included: SIM-006.
