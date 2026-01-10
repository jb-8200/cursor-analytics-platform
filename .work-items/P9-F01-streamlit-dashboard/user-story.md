# User Story: Streamlit Analytics Dashboard

**Feature ID**: P9-F01-streamlit-dashboard
**Phase**: P9 (Streamlit Dashboard)
**Created**: January 9, 2026
**Status**: COMPLETE ✅ (12/12 tasks)

## Overview

As a **engineering manager**, I want a **Streamlit analytics dashboard** so that I can **visualize AI coding impact on velocity, review costs, and code quality** with data from DuckDB (dev) or Snowflake (prod).

## Context

This dashboard replaces/complements the P6 (cursor-viz-spa) React SPA with a Python-native visualization layer. It consumes data from the P8 data tier (dbt marts).

**Key Principle**: The dashboard is production-ready, deployable to Cloud Run, and can connect to either DuckDB (local dev) or Snowflake (production).

## Data Flow Philosophy

This dashboard is the **consumer** in the data contract hierarchy. It never touches raw API data:

```
┌──────────────────────────────────────────────────────────────────────────┐
│                       DATA CONTRACT HIERARCHY                             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  1. API LAYER (cursor-sim) - SOURCE OF TRUTH                            │
│     ├── Response formats: {items:[]} or raw arrays                      │
│     ├── Field names: camelCase (commitHash, userEmail)                  │
│     ├── Data types: Defined by cursor-sim SPEC.md                       │
│     └── Contract: services/cursor-sim/SPEC.md                           │
│                                                                          │
│  2. DATA TIER (P8: dbt) - TRANSFORMATION CONTRACT                       │
│     ├── Column mapping: camelCase → snake_case                          │
│     ├── Aggregations: Weekly metrics, AI usage bands                    │
│     ├── Materializations: Tables in main_mart schema                    │
│     └── Contract: dbt/models/schema.yml                                 │
│                                                                          │
│  3. KPI DASHBOARD (P9: Streamlit) - VISUALIZATION CONSUMER              │
│     ├── Queries: SELECT from main_mart.mart_* tables                    │
│     ├── Filters: Parameterized repo_name and days                       │
│     ├── Charts: Pre-aggregated data from dbt marts                      │
│     └── Contract: services/streamlit-dashboard/queries/*.py             │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**Key Insight**: Dashboard queries NEVER access raw API data. All data flows through dbt transformations.

**Security**: All user inputs are parameterized to prevent SQL injection:
```python
# SECURE: Parameters passed to database
sql = "WHERE repo_name = $repo AND week >= CURRENT_DATE - INTERVAL '{days}' DAY"
params = {"repo": repo_name}
query(sql, params)

# VULNERABLE: Never concatenate user input
sql = f"WHERE repo_name = '{repo_name}'"  # SQL injection risk!
```

---

## User Stories

### US-P9-001: Database Connection Abstraction

**As a** developer
**I want** a database connector that works with both DuckDB and Snowflake
**So that** I can develop locally with DuckDB and deploy to production with Snowflake

**Acceptance Criteria**:
```gherkin
Given the environment variable DB_MODE is set to "duckdb"
When the dashboard starts
Then it connects to the local DuckDB file at /data/analytics.duckdb
And queries return data from dbt mart tables

Given the environment variable DB_MODE is set to "snowflake"
When the dashboard starts
Then it connects to Snowflake using SNOWFLAKE_* credentials
And queries return data from the same dbt mart tables
```

---

### US-P9-002: Velocity Dashboard Page

**As a** engineering manager
**I want** a velocity metrics dashboard
**So that** I can track PR cycle times and throughput by team

**Acceptance Criteria**:
```gherkin
Given the mart.velocity table is populated
When I navigate to the Velocity page
Then I see a line chart of weekly cycle time trends
And I see a breakdown of coding, pickup, and review lead times
And I can filter by repository
And I see P50 and P90 cycle time metrics
And the data refreshes when I click the Refresh button (dev mode only)
```

**Metrics Displayed**:
- Weekly PR count
- Average cycle time (coding + pickup + review)
- P50/P90 cycle times
- Active developer count
- Average PR size (lines changed)

---

### US-P9-003: AI Impact Dashboard Page

**As a** SDLC researcher
**I want** an AI impact analysis dashboard
**So that** I can compare metrics across AI usage bands (low/medium/high)

**Acceptance Criteria**:
```gherkin
Given the mart.ai_impact table is populated
When I navigate to the AI Impact page
Then I see metrics grouped by AI usage band (low <30%, medium 30-60%, high >60%)
And I see a comparison of cycle times across bands
And I see a comparison of revert rates across bands
And visualizations highlight correlations between AI usage and outcomes
```

**Metrics Displayed**:
- PR count by AI band
- Average cycle time by AI band
- Revert rate by AI band
- Review density by AI band

---

### US-P9-004: Quality Dashboard Page

**As a** engineering manager
**I want** a code quality dashboard
**So that** I can track revert rates and identify quality trends

**Acceptance Criteria**:
```gherkin
Given the mart.quality table is populated
When I navigate to the Quality page
Then I see weekly revert rate trends
And I see bug fix rate trends
And I can filter by repository
And I see a breakdown of reverts by AI ratio bands
```

**Metrics Displayed**:
- Weekly revert rate
- Weekly bug fix rate
- Reverts by AI ratio band
- Quality trend over time

---

### US-P9-005: Review Costs Dashboard Page

**As a** engineering manager
**I want** a review costs dashboard
**So that** I can understand code review burden and efficiency

**Acceptance Criteria**:
```gherkin
Given the mart.review_costs table is populated
When I navigate to the Review Costs page
Then I see average review iterations per PR
And I see review comment density (comments per line of code)
And I see average reviewer count per PR
And I can compare review costs by AI usage band
```

---

### US-P9-006: Data Refresh Button (Dev Mode)

**As a** developer
**I want** a "Refresh Data" button in the sidebar
**So that** I can trigger the ETL pipeline and see fresh data

**Acceptance Criteria**:
```gherkin
Given DB_MODE is set to "duckdb" (dev mode)
When I click the "Refresh Data" button in the sidebar
Then the loader extracts fresh data from cursor-sim
And dbt runs to update mart tables
And the dashboard refreshes with new data
And I see a success message

Given DB_MODE is set to "snowflake" (prod mode)
When I view the sidebar
Then the "Refresh Data" button shows a tooltip saying "Data updates via scheduled jobs"
And the button is disabled
```

---

### US-P9-007: Production Deployment (Cloud Run)

**As a** DevOps engineer
**I want** the Streamlit dashboard deployed to Cloud Run
**So that** stakeholders can access it via a public URL

**Acceptance Criteria**:
```gherkin
Given the Dockerfile exists for the dashboard
When I deploy to Cloud Run
Then the dashboard is accessible via HTTPS URL
And it connects to Snowflake using secrets
And it auto-scales based on traffic
And health checks pass
```

---

### US-P9-008: Caching for Performance

**As a** user
**I want** the dashboard to load quickly
**So that** I don't wait for slow database queries

**Acceptance Criteria**:
```gherkin
Given I navigate to any dashboard page
When data is queried from the database
Then results are cached using @st.cache_data
And subsequent page loads use cached data
And cache expires after 5 minutes (configurable)
And I can manually clear cache via Refresh button
```

---

## Non-Functional Requirements

### Performance

| Metric | Target |
|--------|--------|
| Initial page load | < 3 seconds |
| Chart render | < 1 second |
| Cached query | < 100ms |
| Uncached query (DuckDB) | < 500ms |
| Uncached query (Snowflake) | < 2 seconds |

### Accessibility

- WCAG 2.1 AA compliance where possible
- Keyboard navigation for filters
- Alt text for charts (via Streamlit's built-in features)

### Security

- No credentials in code
- Snowflake credentials via environment variables
- Cloud Run service account for GCP resources

---

## Dependencies

- **P8 (Data Tier)**: dbt mart tables must exist
- **cursor-sim (P4)**: Must be running for dev mode refresh
- **DuckDB**: Local analytics database
- **Snowflake**: Production data warehouse
- **Streamlit**: Dashboard framework
- **Plotly/Altair**: Charting library

---

## Out of Scope

- User authentication (use Cloud Run IAM for access control)
- Custom theming beyond Streamlit defaults
- Export to PDF/Excel (future enhancement)
- Real-time streaming updates
- Mobile-optimized layout (desktop-first)

---

## Wireframes

### Sidebar
```
┌─────────────────────────┐
│      DOXAPI             │
│   Analytics Dashboard   │
│                         │
│  ───────────────────    │
│                         │
│  📊 Pages               │
│  • Velocity             │
│  • AI Impact            │
│  • Quality              │
│  • Review Costs         │
│                         │
│  ───────────────────    │
│                         │
│  🔧 Filters             │
│  Repository: [All ▼]    │
│  Date Range: [90 days ▼]│
│                         │
│  ───────────────────    │
│                         │
│  [🔄 Refresh Data]      │
│                         │
│  Last updated:          │
│  2026-01-09 10:30 AM    │
└─────────────────────────┘
```

### Velocity Page
```
┌─────────────────────────────────────────────────────────────────┐
│  Velocity Metrics                                    [Filter ▼] │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐             │
│  │  142 PRs     │ │  4.2 days    │ │  32 devs     │             │
│  │  This Week   │ │  Avg Cycle   │ │  Active      │             │
│  └──────────────┘ └──────────────┘ └──────────────┘             │
│                                                                  │
│  Cycle Time Trend                                                │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                                              ___            │ │
│  │                    ___                   ___/   \___       │ │
│  │      ___      ___/   \___           ___/            \      │ │
│  │  ___/   \____/           \_________/                 \___  │ │
│  │                                                            │ │
│  └────────────────────────────────────────────────────────────┘ │
│   Jan      Feb      Mar      Apr      May      Jun              │
│                                                                  │
│  Cycle Time Breakdown                                            │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  ████████████░░░░░░░░░░░░░░░░░░  Coding: 2.1 days (50%)    │ │
│  │  ████░░░░░░░░░░░░░░░░░░░░░░░░░░  Pickup: 0.8 days (19%)    │ │
│  │  █████░░░░░░░░░░░░░░░░░░░░░░░░░  Review: 1.3 days (31%)    │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

**Next**: See `design.md` for technical architecture and `task.md` for implementation breakdown.
