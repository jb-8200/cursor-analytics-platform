# Technical Design: Data Tier (ETL Pipeline)

**Feature ID**: P8-F01-data-tier
**Phase**: P8 (Data Tier)
**Created**: January 9, 2026
**Status**: NOT_STARTED

## Overview

This feature implements a "Modern Data Stack in a Box" pattern for the Cursor Analytics Platform, replacing the deprecated P5 (cursor-analytics-core) with a dbt-based transformation layer.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           P8: DATA TIER ARCHITECTURE                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌────────────────────┐                                                     │
│  │    cursor-sim      │  P4: Source of Truth (UNCHANGED)                    │
│  │    (Go REST API)   │  Port 8080                                          │
│  │                    │                                                      │
│  │  Endpoints:        │                                                      │
│  │  /analytics/*      │  Cursor API (commits, team stats)                   │
│  │  /repos/*          │  GitHub API (PRs, reviews)                          │
│  └─────────┬──────────┘                                                     │
│            │                                                                 │
│            │ REST API (JSON)                                                 │
│            ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         EXTRACT LAYER                                │   │
│  │  ┌──────────────────────────────────────────────────────────────┐  │   │
│  │  │  tools/api-loader/loader.py                                   │  │   │
│  │  │  - Paginated extraction from cursor-sim                       │  │   │
│  │  │  - Handles raw array responses (cursor-sim contract)          │  │   │
│  │  │  - Outputs: data/raw/*.parquet                                │  │   │
│  │  └──────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│            │                                                                 │
│            │ Parquet Files                                                   │
│            ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         LOAD LAYER                                   │   │
│  │  ┌──────────────────────────────────────────────────────────────┐  │   │
│  │  │  data/analytics.duckdb (dev) │ Snowflake (prod)              │  │   │
│  │  │  Schema: raw.*                                                │  │   │
│  │  │  - raw.commits                                                │  │   │
│  │  │  - raw.pull_requests                                          │  │   │
│  │  │  - raw.reviews                                                │  │   │
│  │  │  - raw.repos                                                  │  │   │
│  │  └──────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│            │                                                                 │
│            │ SQL                                                             │
│            ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         TRANSFORM LAYER (dbt)                        │   │
│  │  ┌──────────────────────────────────────────────────────────────┐  │   │
│  │  │  dbt/models/                                                  │  │   │
│  │  │  ├── staging/       (views)     stg_commits, stg_prs, etc.   │  │   │
│  │  │  ├── intermediate/  (ephemeral) int_pr_with_commits          │  │   │
│  │  │  └── marts/         (tables)    mart_velocity, mart_quality  │  │   │
│  │  └──────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│            │                                                                 │
│            │ SQL Queries                                                     │
│            ▼                                                                 │
│  ┌────────────────────┐                                                     │
│  │  P9: Streamlit     │  Dashboard (separate phase)                         │
│  └────────────────────┘                                                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Component Details

### 1. API Loader (`tools/api-loader/`)

**Purpose**: Extract data from cursor-sim REST API into Parquet files.

**Key Design Decisions**:
- **Simple extraction only**: No business logic in Python. All transformations in dbt.
- **Handle cursor-sim contract**: Raw arrays for GitHub-style endpoints, paginated response for commits.
- **Idempotent**: Full refresh on each run (no incremental for MVP).

```
tools/api-loader/
├── loader.py           # Main extraction script
├── extractors/
│   ├── __init__.py
│   ├── cursor_api.py   # /analytics/*, /teams/*
│   └── github_api.py   # /repos/*
├── schemas/            # Validation schemas
│   ├── commits.json
│   ├── pull_requests.json
│   └── reviews.json
├── requirements.txt
└── Dockerfile          # For production job
```

**API Response Format Handling** (CRITICAL):

cursor-sim uses **two response formats**. The loader MUST handle both:

**Format 1: Cursor Analytics Style (Pagination Wrapper)**
```json
{
  "items": [...],      // Use this key
  "totalCount": 1000,
  "page": 1,
  "pageSize": 500
}
```

**Format 2: GitHub Style (Raw Arrays)**
```json
[...]  // Raw array with no wrapper
```

**Implementation** (`tools/api-loader/extractors/base.py`):
```python
def fetch_cursor_style_paginated(self, endpoint, ...):
    response = resp.json()

    # Handle both formats
    if "items" in response:
        data = response.get("items", [])
        total_count = response.get("totalCount", 0)
    else:
        data = response.get("data", [])
        pagination = response.get("pagination", {})
```

**Endpoint Mapping**:

| Endpoint | Response Format | Loader Path |
|----------|-----------------|------------|
| `GET /analytics/ai-code/commits` | `{items:[]}` | Extract `items` array |
| `GET /repos` | Raw array | Use directly |
| `GET /repos/{o}/{r}/pulls` | Raw array | Use directly |
| `GET /repos/{o}/{r}/pulls/{n}/reviews` | Raw array | Use directly |

**Column Mapping Contract** (API → dbt):

| cursor-sim Field (camelCase) | dbt Staging Column (snake_case) | Type | Notes |
|-----|-----|------|------|
| `commitHash` | `commit_hash` | VARCHAR | Primary key |
| `userEmail` | `user_email` | VARCHAR | Developer lookup |
| `repoName` | `repo_name` | VARCHAR | Repository lookup |
| `tabLinesAdded` | `tab_lines_added` | INTEGER | TAB AI completions |
| `composerLinesAdded` | `composer_lines_added` | INTEGER | Composer AI edits |
| `nonAiLinesAdded` | `non_ai_lines_added` | INTEGER | Human-only lines |
| `commitTs` | `committed_at` | TIMESTAMP | UTC timestamp |

**Mapping Implementation**: `dbt/models/staging/stg_commits.sql` and siblings

---

### 2. DuckDB/Snowflake (`data/` and Snowflake)

**Development**: DuckDB file at `data/analytics.duckdb`
**Production**: Snowflake with schemas `RAW`, `STAGING`, `MART`

**DuckDB Schema Naming Convention** (CRITICAL):

DuckDB uses `main_` prefix for custom schemas. **All references MUST include the prefix**:

```sql
-- CORRECT (DuckDB)
SELECT * FROM main_raw.commits
SELECT * FROM main_staging.stg_commits
SELECT * FROM main_mart.mart_velocity

-- INCORRECT (fails with "Catalog Error")
SELECT * FROM raw.commits
SELECT * FROM mart.velocity
```

**Schema Design**:

```sql
-- RAW schema: Landing zone, minimal transformation
CREATE SCHEMA main_raw;

-- STAGING schema: dbt views for cleaned data
CREATE SCHEMA main_staging;

-- MART schema: dbt tables for analytics
CREATE SCHEMA main_mart;
```

**Dashboard Queries MUST Use**: `FROM main_mart.mart_*` (not `mart.*`)

---

### 3. dbt Project (`dbt/`)

**Purpose**: All transformation logic in declarative SQL.

```
dbt/
├── dbt_project.yml
├── profiles.yml              # Multi-target: dev (DuckDB), prod (Snowflake)
├── packages.yml              # dbt_utils, etc.
├── models/
│   ├── sources.yml           # Define raw sources
│   ├── staging/
│   │   ├── _staging.yml      # Model documentation
│   │   ├── stg_commits.sql
│   │   ├── stg_pull_requests.sql
│   │   └── stg_reviews.sql
│   ├── intermediate/
│   │   ├── _intermediate.yml
│   │   └── int_pr_with_commits.sql
│   └── marts/
│       ├── _marts.yml
│       ├── mart_velocity.sql
│       ├── mart_review_costs.sql
│       ├── mart_quality.sql
│       └── mart_ai_impact.sql
├── macros/
│   └── cross_engine.sql      # DuckDB/Snowflake compatibility
├── tests/
│   └── assert_positive_cycle_times.sql
└── seeds/
    └── ai_ratio_bands.csv    # Reference data
```

---

## Data Models

### Raw Layer (from cursor-sim)

**raw.commits** (from `/analytics/ai-code/commits`):
```sql
commit_hash         VARCHAR PRIMARY KEY
user_id             VARCHAR
user_email          VARCHAR
repo_name           VARCHAR
branch_name         VARCHAR
is_primary_branch   BOOLEAN
total_lines_added   INTEGER
total_lines_deleted INTEGER
tab_lines_added     INTEGER
tab_lines_deleted   INTEGER
composer_lines_added INTEGER
composer_lines_deleted INTEGER
non_ai_lines_added  INTEGER
non_ai_lines_deleted INTEGER
commit_ts           TIMESTAMP
created_at          TIMESTAMP
pull_request_number INTEGER  -- FK to PRs
```

**raw.pull_requests** (from `/repos/{o}/{r}/pulls`):
```sql
number              INTEGER
repo_name           VARCHAR
title               VARCHAR
state               VARCHAR  -- open, closed, merged
author_email        VARCHAR
additions           INTEGER
deletions           INTEGER
changed_files       INTEGER
ai_ratio            FLOAT
was_reverted        BOOLEAN  -- cursor-sim field name
is_bug_fix          BOOLEAN
created_at          TIMESTAMP
merged_at           TIMESTAMP
first_commit_at     TIMESTAMP
first_review_at     TIMESTAMP
last_commit_at      TIMESTAMP
reviewers           VARCHAR[]  -- Array of reviewer IDs

PRIMARY KEY (repo_name, number)
```

**raw.reviews** (from `/repos/{o}/{r}/pulls/{n}/reviews`):
```sql
id                  INTEGER PRIMARY KEY
pr_number           INTEGER
repo_name           VARCHAR
author_id           VARCHAR
body                VARCHAR
state               VARCHAR  -- pending, approved, changes_requested
created_at          TIMESTAMP
```

---

### Staging Layer (dbt views)

**staging.stg_commits**:
- Renames columns to snake_case
- Parses timestamps
- Adds `ai_lines_added = tab_lines_added + composer_lines_added`

**staging.stg_pull_requests**:
- **Calculates** cycle times from timestamps (not from API):
  - `coding_lead_time_hours = EXTRACT(EPOCH FROM (created_at - first_commit_at)) / 3600`
  - `pickup_time_hours = EXTRACT(EPOCH FROM (first_review_at - created_at)) / 3600`
  - `review_lead_time_hours = EXTRACT(EPOCH FROM (merged_at - first_review_at)) / 3600`
- Renames `was_reverted` → `is_reverted`
- Calculates `reviewer_count` from array length

**staging.stg_reviews**:
- Cleans and normalizes review data
- Adds `is_approval = (state = 'approved')`

---

### Mart Layer (dbt tables)

**mart.velocity**:
```sql
week                DATE        -- Week start (Monday)
repo_name           VARCHAR
active_developers   INTEGER
total_prs           INTEGER
avg_pr_size         FLOAT       -- Lines changed
coding_lead_time    FLOAT       -- Average hours
pickup_time         FLOAT       -- Average hours
review_lead_time    FLOAT       -- Average hours
total_cycle_time    FLOAT       -- Sum of above
p50_cycle_time      FLOAT       -- Median
p90_cycle_time      FLOAT       -- 90th percentile
avg_ai_ratio        FLOAT       -- 0.0 - 1.0
```

**mart.ai_impact**:
```sql
week                DATE
ai_usage_band       VARCHAR     -- 'low', 'medium', 'high'
pr_count            INTEGER
avg_ai_ratio        FLOAT
avg_coding_lead_time FLOAT
avg_review_cycle_time FLOAT
avg_review_density  FLOAT       -- Comments per LoC
avg_iterations      FLOAT
revert_rate         FLOAT       -- 0.0 - 1.0
hotfix_rate         FLOAT       -- 0.0 - 1.0
```

**mart.quality**:
```sql
week                DATE
repo_name           VARCHAR
total_prs           INTEGER
reverted_prs        INTEGER
revert_rate         FLOAT
bug_fix_prs         INTEGER
bug_fix_rate        FLOAT
ai_ratio_high_reverts INTEGER   -- Reverts with AI ratio > 0.6
ai_ratio_low_reverts INTEGER    -- Reverts with AI ratio < 0.3
```

---

## Cross-Engine Compatibility

**dbt Macros** for DuckDB/Snowflake parity:

```sql
-- dbt/macros/cross_engine.sql

{% macro date_trunc_week(column) %}
    {% if target.type == 'duckdb' %}
        DATE_TRUNC('week', {{ column }})
    {% else %}
        DATE_TRUNC('WEEK', {{ column }})
    {% endif %}
{% endmacro %}

{% macro array_length(column) %}
    {% if target.type == 'duckdb' %}
        ARRAY_LENGTH({{ column }})
    {% else %}
        ARRAY_SIZE({{ column }})
    {% endif %}
{% endmacro %}

{% macro percentile_cont(p, column) %}
    {% if target.type == 'duckdb' %}
        PERCENTILE_CONT({{ p }}) WITHIN GROUP (ORDER BY {{ column }})
    {% else %}
        PERCENTILE_CONT({{ p }}) WITHIN GROUP (ORDER BY {{ column }})
    {% endif %}
{% endmacro %}
```

---

## Pipeline Orchestration

### Local Development

```bash
# tools/run_pipeline.sh
#!/bin/bash
set -e

CURSOR_SIM_URL=${CURSOR_SIM_URL:-"http://localhost:8080"}

echo "=== Step 1: Extract from cursor-sim ==="
python tools/api-loader/loader.py \
    --url "$CURSOR_SIM_URL" \
    --output data/raw

echo "=== Step 2: Load to DuckDB ==="
python -c "
import duckdb
from pathlib import Path

conn = duckdb.connect('data/analytics.duckdb')
conn.execute('CREATE SCHEMA IF NOT EXISTS raw')

for parquet in Path('data/raw').glob('*.parquet'):
    table = parquet.stem
    conn.execute(f'''
        CREATE OR REPLACE TABLE raw.{table} AS
        SELECT * FROM read_parquet('{parquet}')
    ''')
    print(f'  Loaded raw.{table}')

conn.close()
"

echo "=== Step 3: Run dbt ==="
cd dbt && dbt deps && dbt build --target dev

echo "=== Pipeline complete ==="
```

### Production Jobs

**Cloud Run Job: data-loader**
- Triggered: Every 15 minutes via Cloud Scheduler
- Extracts from cursor-sim Cloud Run service
- Writes Parquet to GCS bucket
- Triggers dbt-runner job on completion

**Cloud Run Job: dbt-runner**
- Triggered: After data-loader completes
- Runs dbt against Snowflake
- Uses `dbt run --target prod`

---

## Directory Structure (Final)

```
cursor-analytics-platform/
├── tools/
│   └── api-loader/               # P8: Extraction layer
│       ├── loader.py
│       ├── extractors/
│       ├── schemas/
│       ├── requirements.txt
│       └── Dockerfile
│
├── dbt/                          # P8: Transform layer
│   ├── dbt_project.yml
│   ├── profiles.yml
│   ├── models/
│   │   ├── staging/
│   │   ├── intermediate/
│   │   └── marts/
│   ├── macros/
│   ├── tests/
│   └── seeds/
│
├── jobs/                         # P8: Production job containers
│   ├── data-loader/
│   │   ├── Dockerfile
│   │   └── requirements.txt
│   └── dbt-runner/
│       ├── Dockerfile
│       └── run.sh
│
├── data/                         # Local data (gitignored)
│   ├── raw/                      # Parquet from loader
│   └── analytics.duckdb         # DuckDB database
│
└── services/
    ├── cursor-sim/               # P4: Keep as-is
    └── streamlit-dashboard/      # P9: Consumes P8 data
```

---

## Testing Strategy

### Loader Tests

```python
# tools/api-loader/tests/test_loader.py

def test_extract_commits_pagination():
    """Verify loader handles paginated commits response."""

def test_extract_prs_raw_array():
    """Verify loader handles raw array (not wrapper object)."""

def test_extract_reviews_raw_array():
    """Verify loader handles raw array for reviews."""

def test_parquet_schema_validation():
    """Verify output Parquet matches expected schema."""
```

### dbt Tests

```yaml
# dbt/models/staging/_staging.yml
models:
  - name: stg_commits
    columns:
      - name: commit_sha
        tests: [unique, not_null]
      - name: ai_lines_added
        tests:
          - dbt_utils.accepted_range:
              min_value: 0

  - name: stg_pull_requests
    columns:
      - name: coding_lead_time_hours
        tests:
          - dbt_utils.accepted_range:
              min_value: 0
      - name: ai_ratio
        tests:
          - dbt_utils.accepted_range:
              min_value: 0
              max_value: 1
```

### Integration Tests

```bash
# Test full pipeline
make ci-local

# Verify mart tables populated
duckdb data/analytics.duckdb <<EOF
SELECT COUNT(*) FROM mart.velocity;
SELECT COUNT(*) FROM mart.ai_impact;
EOF
```

---

## Production Deployment

### GCS Bucket Structure

```
gs://cursor-analytics-data/
├── raw/
│   ├── commits/
│   │   └── commits_2026-01-09.parquet
│   ├── pull_requests/
│   │   └── prs_2026-01-09.parquet
│   └── reviews/
│       └── reviews_2026-01-09.parquet
└── dbt-artifacts/
    └── manifest.json
```

### Snowflake Setup

```sql
-- Create warehouse
CREATE WAREHOUSE TRANSFORM_WH
  WAREHOUSE_SIZE = 'XSMALL'
  AUTO_SUSPEND = 60
  AUTO_RESUME = TRUE;

-- Create database and schemas
CREATE DATABASE CURSOR_ANALYTICS;
CREATE SCHEMA CURSOR_ANALYTICS.RAW;
CREATE SCHEMA CURSOR_ANALYTICS.STAGING;
CREATE SCHEMA CURSOR_ANALYTICS.MART;

-- Create storage integration for GCS
CREATE STORAGE INTEGRATION gcs_integration
  TYPE = EXTERNAL_STAGE
  STORAGE_PROVIDER = GCS
  ENABLED = TRUE
  STORAGE_ALLOWED_LOCATIONS = ('gcs://cursor-analytics-data/');
```

---

## Performance Considerations

| Operation | Target | Notes |
|-----------|--------|-------|
| Loader extraction (90 days) | < 60s | Parallel extraction, connection pooling |
| DuckDB load (100k rows) | < 10s | Direct Parquet read |
| dbt run (all models) | < 2min | Incremental in future |
| dbt test | < 30s | Focus on critical tests |

---

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| cursor-sim API changes | High | Contract tests, version pinning |
| DuckDB/Snowflake SQL drift | Medium | Cross-engine macros, test both |
| Snowflake costs | Medium | X-Small warehouse, auto-suspend |
| Data volume growth | Medium | Incremental loading (future) |

---

**Next**: See `task.md` for implementation breakdown.
