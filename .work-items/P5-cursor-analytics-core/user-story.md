# User Stories: cursor-analytics-core

**Feature**: cursor-analytics-core (GraphQL Aggregator)
**Priority**: P0
**Status**: NOT_STARTED

---

## Epic: Data Ingestion

### US-CORE-001: Poll Simulator for New Events

**As a** system administrator
**I want** the aggregator to automatically fetch new events from the simulator
**So that** analytics stay up-to-date without manual intervention

**Priority**: P0
**Story Points**: 5
**Feature**: CORE-001

**Acceptance Criteria:**

**Scenario 1: Automatic polling**
Given the aggregator is running and connected to cursor-sim
When 60 seconds pass
Then the aggregator should make a request to the simulator's commits endpoint

**Scenario 2: Incremental fetching**
Given the aggregator has previously synced data
When it polls again
Then it should only request events since the last sync timestamp

**Scenario 3: Error handling**
Given the simulator is temporarily unavailable
When a poll request fails
Then the aggregator should retry with exponential backoff

---

### US-CORE-002: Store Events in Database

**As a** data analyst
**I want** all usage events to be persisted in a database
**So that** I can query historical data for trend analysis

**Priority**: P0
**Story Points**: 5
**Feature**: CORE-001, CORE-002

**Acceptance Criteria:**

**Scenario 1: Event persistence**
Given an event is received from the simulator
When the ingestion process completes
Then the event should be queryable from the database

**Scenario 2: Deduplication**
Given an event with the same commit hash already exists
When attempting to insert it again
Then the system should ignore the duplicate

**Scenario 3: Schema conformance**
Given commit data from cursor-sim
When stored in PostgreSQL
Then all fields should map correctly to the schema

---

## Epic: Metric Calculations

### US-CORE-003: Calculate Acceptance Rate

**As a** engineering manager
**I want** to see each developer's suggestion acceptance rate
**So that** I can identify who is effectively using AI assistance

**Priority**: P0
**Story Points**: 5
**Feature**: CORE-004

**Acceptance Criteria:**

**Scenario 1: Basic calculation**
Given a developer has 100 AI suggestions shown and 75 accepted
When I query their acceptance rate
Then the result should be 75.0%

**Scenario 2: No suggestions**
Given a developer has no recorded suggestions
When I query their acceptance rate
Then the result should be null (not zero)

---

### US-CORE-004: Calculate AI Velocity

**As a** engineering leader measuring AI adoption
**I want** to see what percentage of code is AI-generated
**So that** I can track our return on investment in AI tools

**Priority**: P0
**Story Points**: 5
**Feature**: CORE-004

**Acceptance Criteria:**

**Scenario 1: Basic calculation**
Given events show 500 AI lines added and 2000 total lines added
When I query AI velocity
Then the result should be 25.0%

**Scenario 2: Breakdown by source**
Given commits have tab and composer lines
When I query AI velocity with breakdown
Then I should see separate TAB and COMPOSER percentages

---

## Epic: GraphQL API

### US-CORE-005: Query Developer Statistics

**As a** frontend developer building the dashboard
**I want** a GraphQL query to fetch developer statistics
**So that** I can display individual performance metrics

**Priority**: P0
**Story Points**: 5
**Feature**: CORE-003, CORE-005

**Acceptance Criteria:**

**Scenario 1: Query single developer**
Given a developer with ID "user_001" exists
When I execute `query { developer(id: "user_001") { name, team } }`
Then I should receive that developer's data

**Scenario 2: List developers with filters**
Given multiple developers exist
When I execute `query { developers(team: "Platform") { id, name } }`
Then I should receive only developers on that team

**Scenario 3: Developer metrics**
Given a developer has commit data
When I query their statistics
Then I should receive: acceptanceRate, aiVelocity, totalCommits, totalLinesAdded

---

### US-CORE-006: Query Dashboard Summary

**As a** engineering manager viewing the dashboard
**I want** a single query that returns all dashboard KPIs
**So that** the dashboard loads efficiently with one request

**Priority**: P0
**Story Points**: 3
**Feature**: CORE-007

**Acceptance Criteria:**

**Scenario 1: Complete KPI response**
Given the database contains analytics data
When I execute `query { dashboardSummary { totalDevelopers, activeDevelopers, overallAcceptanceRate } }`
Then all requested fields should be populated

**Scenario 2: Time-filtered summary**
Given I want metrics for a specific period
When I execute with dateRange parameter
Then KPIs should be calculated for that period only

---

## Dependencies

- **Requires**: cursor-sim Phase 1 (COMPLETE)
- **Enables**: cursor-viz-spa dashboard implementation
