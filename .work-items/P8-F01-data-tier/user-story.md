# User Story: Data Tier (ETL Pipeline)

**Feature ID**: P8-F01-data-tier
**Phase**: P8 (Data Tier)
**Created**: January 9, 2026
**Status**: NOT_STARTED

## Overview

As a **platform operator**, I want a **modern data stack ETL pipeline** so that I can **extract data from cursor-sim, transform it with dbt, and serve it to analytics dashboards** with dev/prod parity using DuckDB locally and Snowflake in production.

## Context

This phase replaces the deprecated P5 (cursor-analytics-core) Node.js/GraphQL backend with a "Modern Data Stack in a Box" pattern:

- **Extract**: Python loader pulls data from cursor-sim REST API
- **Transform**: dbt models handle all business logic in SQL
- **Load**: DuckDB (dev) or Snowflake (prod)
- **Serve**: SQL queries from Streamlit dashboard (P9)

**Key Principle**: cursor-sim (P4) is the source of truth. We do NOT modify its API. We build everything on top of its existing contract.

---

## User Stories

### US-P8-001: API Data Extraction

**As a** data engineer
**I want** a Python loader that extracts data from cursor-sim API
**So that** I can land raw data in Parquet format for downstream processing

**Acceptance Criteria**:
```gherkin
Given cursor-sim is running at http://localhost:8080
When I run the loader with --url http://localhost:8080 --output data/raw
Then Parquet files are created for commits, pull_requests, reviews, and repos
And the loader handles pagination correctly
And the loader handles cursor-sim's raw array responses (not wrapper objects)
And extraction completes in under 60 seconds for 90 days of data
```

**API Contract Notes** (cursor-sim returns raw arrays, NOT wrapper objects):
- `GET /analytics/ai-code/commits` → `{"items": [...], "totalCount": N}`
- `GET /repos` → `[...]` (raw array)
- `GET /repos/{owner}/{repo}/pulls` → `[...]` (raw array)
- `GET /repos/{owner}/{repo}/pulls/{n}/reviews` → `[...]` (raw array)

---

### US-P8-002: DuckDB Raw Data Loading

**As a** data engineer
**I want** to load Parquet files into DuckDB raw schema
**So that** dbt can transform them using standard SQL

**Acceptance Criteria**:
```gherkin
Given Parquet files exist in data/raw/
When I run the DuckDB loader
Then tables are created in the raw schema: raw.commits, raw.pull_requests, raw.reviews
And data types are correctly inferred from Parquet
And the process is idempotent (CREATE OR REPLACE)
And loading completes in under 10 seconds for 100k rows
```

---

### US-P8-003: dbt Staging Models

**As a** analytics engineer
**I want** dbt staging models that clean and normalize raw data
**So that** downstream models have consistent, well-typed data

**Acceptance Criteria**:
```gherkin
Given raw tables exist in DuckDB
When I run dbt run --target dev
Then staging views are created: staging.stg_commits, staging.stg_pull_requests, staging.stg_reviews
And column names are snake_case
And timestamps are properly parsed
And null values are handled appropriately
And calculated fields are added (e.g., cycle times from timestamps)
```

**Calculated Fields** (computed in dbt, not expected from cursor-sim):
- `coding_lead_time_hours` = `created_at - first_commit_at`
- `pickup_time_hours` = `first_review_at - created_at`
- `review_lead_time_hours` = `merged_at - first_review_at`
- `reviewer_count` = array length of reviewers
- `is_reverted` = renamed from cursor-sim's `was_reverted`

---

### US-P8-004: dbt Mart Models

**As a** analytics engineer
**I want** dbt mart models that aggregate metrics by AI usage bands
**So that** the Streamlit dashboard can query pre-computed analytics

**Acceptance Criteria**:
```gherkin
Given staging models are populated
When I run dbt run --target dev
Then mart tables are created:
  - mart.velocity (weekly metrics by repo/developer)
  - mart.review_costs (review density, iterations)
  - mart.quality (revert rate, bug fix rate)
  - mart.ai_impact (metrics grouped by AI ratio bands)
And mart tables are materialized as tables (not views)
And dbt tests pass for data quality
```

---

### US-P8-005: DuckDB/Snowflake Parity

**As a** data engineer
**I want** dbt macros that handle SQL dialect differences
**So that** the same models work on both DuckDB (dev) and Snowflake (prod)

**Acceptance Criteria**:
```gherkin
Given dbt models use cross-engine macros
When I run dbt run --target dev
Then models execute successfully on DuckDB
When I run dbt run --target prod
Then models execute successfully on Snowflake
And DATE_TRUNC, PERCENTILE_CONT work correctly on both engines
And array functions are handled appropriately per engine
```

---

### US-P8-006: Pipeline Orchestration

**As a** platform operator
**I want** a single command to run the complete pipeline
**So that** I can refresh data with one step

**Acceptance Criteria**:
```gherkin
Given cursor-sim is running
When I run make pipeline or ./tools/run_pipeline.sh
Then the loader extracts fresh data from cursor-sim
And Parquet files are loaded into DuckDB
And dbt runs all models successfully
And dbt tests pass
And the pipeline completes in under 5 minutes
```

---

### US-P8-007: Production Job Containers

**As a** DevOps engineer
**I want** Docker containers for loader and dbt runner
**So that** I can deploy them as Cloud Run Jobs

**Acceptance Criteria**:
```gherkin
Given Dockerfiles exist for jobs/data-loader and jobs/dbt-runner
When I build the containers
Then the loader container can extract to GCS
And the dbt-runner container can run against Snowflake
And containers start in under 30 seconds
And containers use non-root users
And health checks are implemented
```

---

## Dependencies

- **P4 (cursor-sim)**: Must be running and healthy
- **cursor-sim API contract**: Raw array responses for GitHub-style endpoints
- **DuckDB**: Local analytics database
- **dbt-duckdb**: dbt adapter for DuckDB
- **dbt-snowflake**: dbt adapter for production
- **Snowflake account**: For production deployment

---

## Out of Scope

- Real-time streaming (this is batch ETL)
- Data lineage tracking (future enhancement)
- Data catalog (future enhancement)
- CI/CD for dbt (separate task)
- Incremental loading (full refresh for MVP)

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| cursor-sim API changes | Contract tests in loader, pin cursor-sim version |
| DuckDB/Snowflake SQL drift | Cross-engine macros, test on both engines |
| Large data volumes | Pagination in loader, batch processing |
| Snowflake costs | Use X-Small warehouse, auto-suspend |

---

**Next**: See `design.md` for technical architecture and `task.md` for implementation breakdown.
