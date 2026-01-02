# User Stories: Cursor Usage Analytics Platform

> **📚 REFERENCE DOCUMENT**
> This is a project-level user story overview for orientation purposes.
> **Source of truth**: `.work-items/{feature}/user-story.md` for active feature requirements.

**Version**: 1.0.0
**Last Updated**: January 2026

This document contains user stories organized by service and feature area. Each story follows the standard format and includes detailed acceptance criteria that serve as the basis for test cases in our TDD workflow.

## Story Format

Each user story follows this structure to ensure consistency and completeness. The story statement uses the classic "As a [role], I want [capability], so that [benefit]" format. Acceptance criteria use the Given-When-Then format to make test case derivation straightforward.

## Service A: cursor-sim User Stories

### Epic: Simulator Configuration

#### US-SIM-001: Configure Simulator via CLI

**As a** developer setting up a demo environment  
**I want to** configure the simulator through command-line flags  
**So that** I can quickly adjust simulation parameters without editing configuration files

**Priority**: P0  
**Story Points**: 3  
**Feature**: SIM-004

**Acceptance Criteria:**

**Scenario 1: Default configuration**  
Given I start the simulator without any flags, when the application initializes, then it should use default values of port 8080, 50 developers, high velocity, and 0.2 fluctuation.

**Scenario 2: Custom port configuration**  
Given I provide the flag "--port=9000", when the application initializes, then it should listen on port 9000.

**Scenario 3: Invalid port rejected**  
Given I provide the flag "--port=80" (privileged port), when the application attempts to initialize, then it should exit with error code 1 and display a message explaining that ports below 1024 require root privileges.

**Scenario 4: Developer count configuration**  
Given I provide the flag "--developers=100", when the application initializes, then it should generate exactly 100 developer profiles.

**Scenario 5: Help text displayed**  
Given I provide the flag "--help", when the application runs, then it should display documentation for all available flags and exit with code 0.

**Test Cases:**
```go
func TestCLI_DefaultConfiguration(t *testing.T)
func TestCLI_CustomPort(t *testing.T)
func TestCLI_InvalidPortRejected(t *testing.T)
func TestCLI_DeveloperCount(t *testing.T)
func TestCLI_HelpText(t *testing.T)
```

---

### Epic: Developer Simulation

#### US-SIM-002: Generate Realistic Developer Profiles

**As a** QA engineer testing the analytics dashboard  
**I want** the simulator to generate realistic developer profiles  
**So that** my test data resembles production data patterns

**Priority**: P0  
**Story Points**: 5  
**Feature**: SIM-001

**Acceptance Criteria:**

**Scenario 1: Profile uniqueness**  
Given the simulator is configured with 50 developers, when profiles are generated, then each profile should have a unique ID, unique email address, and unique name.

**Scenario 2: Seniority distribution**  
Given the simulator generates 100 developer profiles with default settings, when I analyze the seniority distribution, then approximately 20% should be junior, 50% should be mid-level, and 30% should be senior (within 10% variance).

**Scenario 3: Acceptance rate correlation**  
Given a developer profile is generated with seniority "senior", when I check their acceptance rate, then it should be between 0.85 and 0.95.

**Scenario 4: Reproducible generation**  
Given I start the simulator twice with the same seed value, when profiles are generated, then the same set of profiles should be produced in both runs.

**Scenario 5: Team assignment**  
Given the simulator generates 50 developers, when I group them by team, then no team should have fewer than 3 members and no team should have more than 15 members.

**Test Cases:**
```go
func TestDeveloperGeneration_UniqueProfiles(t *testing.T)
func TestDeveloperGeneration_SeniorityDistribution(t *testing.T)
func TestDeveloperGeneration_AcceptanceRateByLevel(t *testing.T)
func TestDeveloperGeneration_Reproducibility(t *testing.T)
func TestDeveloperGeneration_TeamAssignment(t *testing.T)
```

---

#### US-SIM-003: Generate Usage Events

**As a** data analyst evaluating AI tool adoption  
**I want** the simulator to generate realistic usage events  
**So that** I can test analytics calculations with data that matches real patterns

**Priority**: P0  
**Story Points**: 8  
**Feature**: SIM-002

**Acceptance Criteria:**

**Scenario 1: Event type variety**  
Given the simulator runs for 1 minute with 10 developers, when I collect all generated events, then I should see all four event types (cpp_suggestion_shown, cpp_suggestion_accepted, chat_message, cmd_k_prompt).

**Scenario 2: Suggestion acceptance follows profile rate**  
Given a developer has an acceptance rate of 0.9, when 100 cpp_suggestion_shown events are generated for them, then approximately 90 cpp_suggestion_accepted events should also be generated (within 15% variance).

**Scenario 3: Velocity affects event rate**  
Given the simulator is configured with velocity "high", when I count events generated in one minute for a single developer, then the count should be approximately 100 events per hour (1.67 events per minute, within 50% variance due to Poisson distribution).

**Scenario 4: Fluctuation adds variance**  
Given the simulator is configured with fluctuation 0.3 and 20 developers, when I compare event rates across developers, then the standard deviation of rates should be at least 15% of the mean rate.

**Scenario 5: Events have valid timestamps**  
Given events are being generated, when I examine event timestamps, then all timestamps should be in the past relative to current time and no two events from the same developer should have identical timestamps.

**Test Cases:**
```go
func TestEventGeneration_AllEventTypes(t *testing.T)
func TestEventGeneration_AcceptanceRateCorrelation(t *testing.T)
func TestEventGeneration_VelocityImpact(t *testing.T)
func TestEventGeneration_FluctuationVariance(t *testing.T)
func TestEventGeneration_ValidTimestamps(t *testing.T)
```

---

### Epic: API Endpoints

#### US-SIM-004: List Simulated Developers

**As a** backend developer integrating with the simulator  
**I want** an endpoint that returns all simulated developers  
**So that** I can fetch the developer roster for processing

**Priority**: P0  
**Story Points**: 2  
**Feature**: SIM-003

**Acceptance Criteria:**

**Scenario 1: Successful response**  
Given the simulator is running with 50 developers, when I call GET /v1/org/users, then I should receive a 200 status code with a JSON array of 50 developer objects.

**Scenario 2: Response schema**  
Given I call GET /v1/org/users, when I examine a developer object in the response, then it should contain fields: id (string), name (string), email (string), team (string), and seniority (string).

**Scenario 3: CORS headers**  
Given I call GET /v1/org/users with an Origin header, when I examine the response headers, then Access-Control-Allow-Origin should be present.

**Test Cases:**
```go
func TestAPI_ListDevelopers_Success(t *testing.T)
func TestAPI_ListDevelopers_Schema(t *testing.T)
func TestAPI_ListDevelopers_CORS(t *testing.T)
```

---

#### US-SIM-005: Query Activity Events

**As a** backend developer building the aggregator  
**I want** an endpoint to query activity events by time range  
**So that** I can fetch only new events since my last poll

**Priority**: P0  
**Story Points**: 5  
**Feature**: SIM-003

**Acceptance Criteria:**

**Scenario 1: Time range filtering**  
Given events exist from the last hour, when I call GET /v1/stats/activity?from=<1-hour-ago>&to=<now>, then I should only receive events within that time range.

**Scenario 2: Empty result for future range**  
Given the current time is 2026-01-15T10:00:00Z, when I call GET /v1/stats/activity?from=2026-01-15T11:00:00Z&to=2026-01-15T12:00:00Z, then I should receive an empty array.

**Scenario 3: Invalid date format rejected**  
Given I call GET /v1/stats/activity?from=invalid-date, when the server processes the request, then I should receive a 400 status code with an error message explaining the expected date format.

**Scenario 4: Pagination for large results**  
Given more than 1000 events exist in the queried range, when I call the endpoint without pagination parameters, then the response should include the first 1000 events and a "nextCursor" field for pagination.

**Scenario 5: Event object schema**  
Given I call GET /v1/stats/activity with a valid range, when I examine an event in the response, then it should contain: id (string), developerId (string), eventType (string), timestamp (ISO 8601 string), and metadata (object).

**Test Cases:**
```go
func TestAPI_ActivityEvents_TimeFiltering(t *testing.T)
func TestAPI_ActivityEvents_EmptyRange(t *testing.T)
func TestAPI_ActivityEvents_InvalidDate(t *testing.T)
func TestAPI_ActivityEvents_Pagination(t *testing.T)
func TestAPI_ActivityEvents_Schema(t *testing.T)
```

---

## Service B: cursor-analytics-core User Stories

### Epic: Data Ingestion

#### US-CORE-001: Poll Simulator for New Events

**As a** system administrator  
**I want** the aggregator to automatically fetch new events from the simulator  
**So that** analytics stay up-to-date without manual intervention

**Priority**: P0  
**Story Points**: 5  
**Feature**: CORE-001

**Acceptance Criteria:**

**Scenario 1: Automatic polling**  
Given the aggregator is running and connected to a simulator, when 60 seconds pass, then the aggregator should make a request to the simulator's activity endpoint.

**Scenario 2: Incremental fetching**  
Given the aggregator successfully fetched events up to timestamp T1, when the next poll occurs, then the request should use T1 as the "from" parameter.

**Scenario 3: Retry on failure**  
Given the simulator is temporarily unavailable, when a poll attempt fails, then the aggregator should retry with exponential backoff (1s, 2s, 4s, max 30s).

**Scenario 4: Deduplication**  
Given an event with ID "evt-123" was already stored, when a poll returns the same event, then no duplicate record should be created in the database.

**Test Cases:**
```typescript
describe('DataIngestionWorker', () => {
  it('should poll at configured interval')
  it('should use last timestamp for incremental fetch')
  it('should retry with exponential backoff')
  it('should deduplicate events')
})
```

---

#### US-CORE-002: Store Events in Database

**As a** data analyst  
**I want** all usage events to be persisted in a database  
**So that** I can query historical data for trend analysis

**Priority**: P0  
**Story Points**: 5  
**Feature**: CORE-001, CORE-002

**Acceptance Criteria:**

**Scenario 1: Event persistence**  
Given an event is received from the simulator, when the ingestion process completes, then the event should be queryable from the database.

**Scenario 2: Developer auto-creation**  
Given an event references a developer ID not in the database, when the event is processed, then a new developer record should be created from the simulator's developer endpoint.

**Scenario 3: Referential integrity**  
Given an event is stored in the database, when I query the event, then the developer_id foreign key should reference a valid developer record.

**Scenario 4: Transaction safety**  
Given a batch of 100 events is being processed and an error occurs on event 50, when the transaction rolls back, then none of the 100 events should be in the database.

**Test Cases:**
```typescript
describe('EventStorage', () => {
  it('should persist events to database')
  it('should auto-create missing developers')
  it('should maintain referential integrity')
  it('should rollback on partial failure')
})
```

---

### Epic: Metric Calculations

#### US-CORE-003: Calculate Acceptance Rate

**As a** engineering manager  
**I want** to see each developer's suggestion acceptance rate  
**So that** I can identify who is effectively using AI assistance

**Priority**: P0  
**Story Points**: 3  
**Feature**: CORE-004

**Acceptance Criteria:**

**Scenario 1: Basic calculation**  
Given a developer has 100 cpp_suggestion_shown events and 75 cpp_suggestion_accepted events, when I query their acceptance rate, then the result should be 75.0.

**Scenario 2: No suggestions shown**  
Given a developer has 0 cpp_suggestion_shown events, when I query their acceptance rate, then the result should be null (not 0 or error).

**Scenario 3: Time-bounded calculation**  
Given a developer had different behavior last week vs this week, when I query their acceptance rate for this week only, then the calculation should only include this week's events.

**Scenario 4: Precision**  
Given a developer has 99 shown and 33 accepted events, when I query their acceptance rate, then the result should be 33.33 (rounded to 2 decimal places).

**Test Cases:**
```typescript
describe('AcceptanceRateCalculation', () => {
  it('should calculate rate as (accepted/shown)*100')
  it('should return null when no suggestions shown')
  it('should respect time boundaries')
  it('should round to 2 decimal places')
})
```

---

#### US-CORE-004: Calculate AI Velocity

**As a** engineering leader measuring AI adoption  
**I want** to see what percentage of code is AI-generated  
**So that** I can track our return on investment in AI tools

**Priority**: P0  
**Story Points**: 5  
**Feature**: CORE-004

**Acceptance Criteria:**

**Scenario 1: Basic calculation**  
Given events show 500 AI lines added and 2000 total lines added, when I query AI velocity, then the result should be 25.0.

**Scenario 2: No code written**  
Given there are no lines added in the time period, when I query AI velocity, then the result should be null.

**Scenario 3: Team aggregation**  
Given Team Alpha has 3 developers with AI velocities of 20%, 30%, and 40%, when I query team AI velocity, then the result should be the weighted average based on total lines, not a simple average of percentages.

**Test Cases:**
```typescript
describe('AIVelocityCalculation', () => {
  it('should calculate as (AI lines / total lines) * 100')
  it('should return null when no code written')
  it('should use weighted average for team aggregation')
})
```

---

### Epic: GraphQL API

#### US-CORE-005: Query Developer Statistics

**As a** frontend developer building the dashboard  
**I want** a GraphQL query to fetch developer statistics  
**So that** I can display individual performance metrics

**Priority**: P0  
**Story Points**: 5  
**Feature**: CORE-003, CORE-005

**Acceptance Criteria:**

**Scenario 1: Query single developer**  
Given a developer with ID "dev-123" exists, when I execute `query { developer(id: "dev-123") { name, team } }`, then I should receive that developer's data.

**Scenario 2: Include nested statistics**  
Given a developer exists, when I execute `query { developer(id: "dev-123") { stats { acceptanceRate, aiVelocity } } }`, then I should receive calculated statistics.

**Scenario 3: Developer not found**  
Given no developer with ID "dev-999" exists, when I query for that developer, then I should receive null (not an error).

**Scenario 4: Time range parameter**  
Given I query for statistics with a specific date range, when the results are returned, then only events within that range should be included in calculations.

**Test Cases:**
```typescript
describe('GraphQL Developer Queries', () => {
  it('should return developer by ID')
  it('should include nested statistics when requested')
  it('should return null for non-existent developer')
  it('should filter statistics by date range')
})
```

---

#### US-CORE-006: Query Dashboard Summary

**As a** engineering manager viewing the dashboard  
**I want** a single query that returns all dashboard KPIs  
**So that** the dashboard loads efficiently with one request

**Priority**: P1  
**Story Points**: 8  
**Feature**: CORE-007

**Acceptance Criteria:**

**Scenario 1: Complete KPI response**  
Given the database contains analytics data, when I execute `query { dashboardSummary { totalDevelopers, activeDevelopers, overallAcceptanceRate } }`, then all requested fields should be populated.

**Scenario 2: Team comparison inclusion**  
Given multiple teams exist, when I query `dashboardSummary { teamComparison { teamName, averageAcceptanceRate } }`, then all teams should be included in the comparison array.

**Scenario 3: Daily trend data**  
Given a 7-day range is specified, when I query `dashboardSummary(range: {...}) { dailyTrend { date, suggestionsShown } }`, then exactly 7 data points should be returned (one per day).

**Scenario 4: Performance requirement**  
Given 500 developers and 100,000 events exist, when I execute the dashboard summary query, then the response should return within 500ms.

**Test Cases:**
```typescript
describe('GraphQL Dashboard Summary', () => {
  it('should return all KPI fields')
  it('should include all teams in comparison')
  it('should return correct number of trend data points')
  it('should meet performance requirements')
})
```

---

## Service C: cursor-viz-spa User Stories

### Epic: Dashboard Visualization

#### US-VIZ-001: View Dashboard Overview

**As a** engineering manager  
**I want** to see a dashboard with key metrics at a glance  
**So that** I can quickly assess AI tool adoption across my organization

**Priority**: P0  
**Story Points**: 8  
**Feature**: VIZ-001, VIZ-002, VIZ-004

**Acceptance Criteria:**

**Scenario 1: Dashboard loads**  
Given the backend is running and has data, when I navigate to the dashboard URL, then I should see the velocity heatmap and developer efficiency table.

**Scenario 2: KPI header displays**  
Given analytics data exists, when the dashboard loads, then the header should display total developers, active developers, and overall acceptance rate.

**Scenario 3: Responsive layout**  
Given I am viewing on a mobile device (viewport < 768px), when the dashboard renders, then charts should stack vertically and remain readable.

**Scenario 4: Loading state**  
Given the GraphQL query is in progress, when the dashboard is rendering, then skeleton loaders should appear in place of charts.

**Test Cases:**
```typescript
describe('Dashboard Overview', () => {
  it('should render velocity heatmap')
  it('should render developer efficiency table')
  it('should display KPIs in header')
  it('should be responsive on mobile')
  it('should show loading skeletons')
})
```

---

#### US-VIZ-002: View Velocity Heatmap

**As a** engineering manager analyzing adoption trends  
**I want** to see a heatmap of AI code acceptance over time  
**So that** I can identify patterns in when developers use AI assistance most

**Priority**: P0  
**Story Points**: 8  
**Feature**: VIZ-002

**Acceptance Criteria:**

**Scenario 1: Grid rendering**  
Given 52 weeks of data exist, when the heatmap renders, then I should see a grid with 7 rows (days) and 52 columns (weeks).

**Scenario 2: Color intensity**  
Given a day has high AI acceptance count, when I view that cell, then it should appear darker than cells with lower counts.

**Scenario 3: Tooltip on hover**  
Given I hover over a cell, when the tooltip appears, then it should display the date and exact acceptance count.

**Scenario 4: Week labels**  
Given the heatmap renders, when I look at the left edge, then I should see day-of-week labels (Mon, Wed, Fri minimum).

**Scenario 5: Month boundaries**  
Given the heatmap renders, when I look at the top edge, then I should see month labels positioned at month boundaries.

**Test Cases:**
```tsx
describe('VelocityHeatmap', () => {
  it('should render correct grid dimensions')
  it('should apply color intensity based on value')
  it('should display tooltip on hover')
  it('should show day-of-week labels')
  it('should show month labels at boundaries')
})
```

---

#### US-VIZ-003: View Developer Efficiency Table

**As a** team lead reviewing individual performance  
**I want** to see a table of developers with their AI metrics  
**So that** I can identify who may need additional training

**Priority**: P0  
**Story Points**: 5  
**Feature**: VIZ-004

**Acceptance Criteria:**

**Scenario 1: Table columns**  
Given the table renders, when I examine the headers, then I should see columns for Name, Team, Total Suggestions, Accepted, Acceptance Rate, and AI Lines.

**Scenario 2: Sorting**  
Given the table is displaying 20 developers, when I click on the "Acceptance Rate" column header, then rows should reorder by acceptance rate descending.

**Scenario 3: Low performance highlighting**  
Given a developer has an acceptance rate below 20%, when their row renders, then it should have a red background or other visual warning indicator.

**Scenario 4: Pagination**  
Given 100 developers exist, when the table renders with default page size 25, then I should see 25 rows and pagination controls to access other pages.

**Scenario 5: Search filter**  
Given I type "John" in the search box, when the filter applies, then only developers with "John" in their name should appear.

**Test Cases:**
```tsx
describe('DeveloperEfficiencyTable', () => {
  it('should display all required columns')
  it('should sort by column on header click')
  it('should highlight low acceptance rates')
  it('should paginate large datasets')
  it('should filter by search term')
})
```

---

#### US-VIZ-004: Compare Teams with Radar Chart

**As a** VP of Engineering comparing team performance  
**I want** to see a radar chart comparing teams across multiple metrics  
**So that** I can identify which teams may benefit from AI tool training

**Priority**: P1  
**Story Points**: 8  
**Feature**: VIZ-003

**Acceptance Criteria:**

**Scenario 1: Multi-axis display**  
Given the radar chart renders, when I examine the axes, then I should see at least 5 axes: Chat Usage, Code Completion, Acceptance Rate, AI Velocity, and Cmd+K Usage.

**Scenario 2: Team selection**  
Given 5 teams exist, when I select 3 teams from the dropdown, then the chart should display exactly 3 overlapping polygons.

**Scenario 3: Legend**  
Given 3 teams are selected, when I view the legend, then each team should have a distinct color and label.

**Scenario 4: Axis normalization**  
Given teams have very different absolute values, when the chart renders, then axes should be normalized so all teams are comparable.

**Test Cases:**
```tsx
describe('TeamRadarChart', () => {
  it('should render 5 metric axes')
  it('should display selected teams only')
  it('should show legend with team colors')
  it('should normalize axis values')
})
```

---

### Epic: Filtering and Interaction

#### US-VIZ-005: Filter by Date Range

**As a** analyst investigating a specific time period  
**I want** to filter all dashboard views by date range  
**So that** I can analyze trends during a particular sprint or quarter

**Priority**: P1  
**Story Points**: 5  
**Feature**: VIZ-005

**Acceptance Criteria:**

**Scenario 1: Preset selection**  
Given I click the date range dropdown, when I select "Last 30 Days", then all charts should update to show only the last 30 days of data.

**Scenario 2: Custom range**  
Given I select "Custom" from the dropdown, when I pick a start and end date, then charts should update to show only that range.

**Scenario 3: Range validation**  
Given I attempt to set a start date after the end date, when I try to apply, then I should see an error message and the invalid range should not be applied.

**Scenario 4: URL persistence**  
Given I select a custom date range, when I copy the URL and open it in a new tab, then the same date range should be applied.

**Test Cases:**
```tsx
describe('DateRangePicker', () => {
  it('should apply preset ranges')
  it('should allow custom date selection')
  it('should validate date range')
  it('should persist range in URL')
})
```

---

#### US-VIZ-006: Handle Loading and Error States

**As a** user with a slow network connection  
**I want** to see appropriate loading and error feedback  
**So that** I understand the application state and can recover from errors

**Priority**: P1  
**Story Points**: 3  
**Feature**: VIZ-006, VIZ-007

**Acceptance Criteria:**

**Scenario 1: Loading skeleton**  
Given the GraphQL query is pending, when I view a chart component, then I should see a skeleton placeholder matching the chart shape.

**Scenario 2: Error display**  
Given the GraphQL query fails, when I view a chart component, then I should see a user-friendly error message.

**Scenario 3: Retry action**  
Given an error is displayed, when I click the "Retry" button, then the query should be re-executed.

**Scenario 4: Partial failure isolation**  
Given the team comparison query fails but dashboard summary succeeds, when I view the dashboard, then working components should still display data and only the failed component should show an error.

**Test Cases:**
```tsx
describe('LoadingAndErrorStates', () => {
  it('should display skeleton during loading')
  it('should show user-friendly error message')
  it('should retry on button click')
  it('should isolate component failures')
})
```

---

## Cross-Service User Stories

### Epic: End-to-End Data Flow

#### US-E2E-001: View Live Simulation Data

**As a** developer demonstrating the platform  
**I want** to start the entire stack and see data flowing  
**So that** I can show the complete functionality to stakeholders

**Priority**: P0  
**Story Points**: 5  
**Feature**: All services

**Acceptance Criteria:**

**Scenario 1: Docker Compose startup**  
Given I run `docker-compose up`, when all services are healthy, then the dashboard at localhost:3000 should display data from the simulator.

**Scenario 2: Data freshness**  
Given the simulator is generating events, when I wait 2 minutes, then new events should appear in the dashboard (aggregator poll interval is 60s).

**Scenario 3: Configuration propagation**  
Given I start the simulator with "--developers=100", when I view the dashboard, then the developer count should show 100.

**Test Cases:**
```typescript
describe('End-to-End Flow', () => {
  it('should display data after docker-compose up')
  it('should update with new data within 2 minutes')
  it('should reflect simulator configuration')
})
```

---

## Story Index by Feature

For quick reference, here is a mapping from feature IDs to related user stories.

| Feature ID | User Stories |
|------------|--------------|
| SIM-001 | US-SIM-002 |
| SIM-002 | US-SIM-003 |
| SIM-003 | US-SIM-004, US-SIM-005 |
| SIM-004 | US-SIM-001 |
| CORE-001 | US-CORE-001, US-CORE-002 |
| CORE-002 | US-CORE-002 |
| CORE-003 | US-CORE-005, US-CORE-006 |
| CORE-004 | US-CORE-003, US-CORE-004 |
| CORE-005 | US-CORE-005 |
| CORE-007 | US-CORE-006 |
| VIZ-001 | US-VIZ-001 |
| VIZ-002 | US-VIZ-002 |
| VIZ-003 | US-VIZ-004 |
| VIZ-004 | US-VIZ-003 |
| VIZ-005 | US-VIZ-005 |
| VIZ-006 | US-VIZ-006 |
| VIZ-007 | US-VIZ-006 |
