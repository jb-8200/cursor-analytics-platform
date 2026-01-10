# Streamlit Analytics Dashboard

**Feature ID**: P9-F01-streamlit-dashboard
**Created**: 2026-01-09
**Status**: COMPLETE ✅

## Overview

Production-ready Streamlit dashboard for visualizing AI coding analytics. Supports both DuckDB (local development) and Snowflake (production) backends.

### Platform Architecture

This service (P9) provides an alternative analytics path using **dbt + DuckDB/Snowflake**, complementing the original GraphQL path:

```
Path 1 (GraphQL):   cursor-sim → cursor-analytics-core → cursor-viz-spa
                    (P4)         (P5 GraphQL)            (P6 React)

Path 2 (dbt):       cursor-sim → streamlit-dashboard
                    (P4)         (P9 Streamlit + P8 dbt)
```

**Key differences:**
- **Streamlit path (this service)**: Embedded dbt transformations, direct SQL queries, Python-based analytics
- **GraphQL path**: TypeScript services, GraphQL API, React frontend

Both paths consume the same cursor-sim data but use different transformation and presentation layers.

## Features

- **Velocity Metrics**: PR cycle times, throughput, and developer activity
- **AI Impact Analysis**: Metrics grouped by AI usage bands (low/medium/high)
- **Quality Dashboard**: Revert rates and code quality trends
- **Review Costs**: Code review burden and efficiency metrics

## Architecture

```
streamlit-dashboard/
├── app.py                      # Main entry point
├── pages/                      # Multi-page dashboard
├── components/                 # Reusable UI components
├── db/                         # Database connectors
├── queries/                    # SQL query modules
├── pipeline/                   # Embedded ETL (dev mode)
└── tests/                      # Test suite
```

## Data Contract

This dashboard is a **consumer** in the data contract hierarchy. It never touches raw API data:

```
cursor-sim API → api-loader → dbt → Dashboard
(source)        (extract)    (transform)  (consume)
```

**Important Constraints**:
1. Dashboard queries `main_mart.*` tables, not raw API data
2. All user inputs are parameterized (SQL injection prevention)
3. Column names follow dbt schema, not API schema

### Available Columns

Queries return these columns from dbt marts:

**mart_velocity**:
- week, repo_name, active_developers, total_prs, total_commits
- avg_pr_size, avg_files_changed
- **avg_total_cycle_time** (hours)
- avg_ai_ratio, total_ai_lines, total_lines

**mart_ai_impact**:
- week, ai_usage_band, pr_count
- avg_ai_ratio, **avg_total_cycle_time**
- revert_rate, bug_fix_rate
- avg_pr_size, avg_files_changed

**mart_quality**:
- week, repo_name, total_prs, reverted_prs, revert_rate
- bug_fix_prs, bug_fix_rate, avg_reviews_per_pr, unreviewed_prs

**mart_review_costs**:
- week, repo_name, total_prs
- avg_review_rounds, avg_reviewers_per_pr
- **avg_review_cycle_time**
- estimated_review_hours_per_pr, estimated_total_review_hours, large_prs

**Columns NOT available** (removed from queries):
- `p50_cycle_time` - Use `avg_total_cycle_time` instead
- `avg_coding_lead_time` - Not computed in mart
- `avg_review_iterations` - Not computed in mart

## Security

### SQL Injection Prevention

All dashboard queries use parameterized binding:

```python
# SECURE: Parameters passed to DuckDB
sql = "SELECT * FROM main_mart.mart_velocity WHERE repo_name = $repo"
params = {"repo": repo_name}
query(sql, params)

# VULNERABLE (never do this):
sql = f"SELECT * FROM main_mart.mart_velocity WHERE repo_name = '{repo_name}'"
```

### INTERVAL Syntax

DuckDB doesn't support parameterized INTERVAL expressions. Use f-string interpolation:

```python
# CORRECT: f-string for INTERVAL (days is validated as integer)
f"WHERE week >= CURRENT_DATE - INTERVAL '{days}' DAY"

# INCORRECT: Parameter in INTERVAL (fails)
"WHERE week >= CURRENT_DATE - INTERVAL $days DAY"
```

### Schema Naming

DuckDB requires `main_*` prefix for schema-qualified table names:

```sql
-- CORRECT
FROM main_mart.mart_velocity

-- INCORRECT (fails with "Catalog Error")
FROM mart.velocity
```

## Setup

### Local Development (DuckDB)

```bash
# Install dependencies
pip install -r requirements.txt

# Create local data directories
mkdir -p data/raw

# Run dashboard (uses default relative paths)
streamlit run app.py

# OR with explicit configuration
export DB_MODE=duckdb
export DUCKDB_PATH=data/analytics.duckdb
export RAW_DATA_PATH=data/raw
export DBT_PROJECT_DIR=../../dbt
export CURSOR_SIM_URL=http://localhost:8080
streamlit run app.py
```

### Production (Snowflake)

```bash
# Set environment variables
export DB_MODE=snowflake
export SNOWFLAKE_ACCOUNT=xxx.us-central1.gcp
export SNOWFLAKE_USER=STREAMLIT_USER
export SNOWFLAKE_PASSWORD=***
export SNOWFLAKE_DATABASE=CURSOR_ANALYTICS
export SNOWFLAKE_SCHEMA=MART
export SNOWFLAKE_WAREHOUSE=TRANSFORM_WH

# Run dashboard
streamlit run app.py
```

## Docker

### Docker Compose (Recommended)

Run with cursor-sim using Docker Compose:

```bash
# Start cursor-sim + streamlit-dashboard
docker-compose up -d cursor-sim streamlit-dashboard

# Access dashboard
open http://localhost:8501

# View logs
docker-compose logs -f streamlit-dashboard

# Stop services
docker-compose down cursor-sim streamlit-dashboard
```

The docker-compose.yml automatically:
- Starts cursor-sim on port 8080
- Starts streamlit-dashboard on port 8501
- Mounts dbt project at `/app/dbt`
- Persists DuckDB data in `analytics_data` volume
- Configures environment variables

### Manual Docker Build

```bash
# Build image
docker build -t streamlit-dashboard .

# Run container (requires cursor-sim running)
docker run -p 8501:8501 \
  -e DB_MODE=duckdb \
  -e DUCKDB_PATH=/data/analytics.duckdb \
  -e CURSOR_SIM_URL=http://host.docker.internal:8080 \
  -v analytics_data:/data \
  -v $(pwd)/../../dbt:/app/dbt:ro \
  streamlit-dashboard
```

## Testing

```bash
# Run tests
pytest tests/

# Run tests with coverage
pytest --cov=. --cov-report=html tests/
```

### Security Testing

Verify SQL injection protection:
1. Launch dashboard
2. Modify sidebar selectbox value via dev tools to: `'; DROP TABLE mart.velocity; --`
3. Verify no SQL errors occur and data remains intact

## Dependencies

- **P8 (Data Tier)**: Requires dbt mart tables (mart.velocity, mart.ai_impact, mart.quality, mart.review_costs)
- **cursor-sim (P4)**: Required for dev mode refresh
- **DuckDB**: Local analytics database
- **Snowflake**: Production data warehouse (optional)

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DB_MODE` | No | `duckdb` | Database mode: `duckdb` or `snowflake` |
| `DUCKDB_PATH` | No | `data/analytics.duckdb` | Path to DuckDB file (relative to app) |
| `RAW_DATA_PATH` | No | `data/raw` | Directory for raw extracted data |
| `DBT_PROJECT_DIR` | No | `/app/dbt` | Path to dbt project directory |
| `CURSOR_SIM_URL` | No | `http://localhost:8080` | cursor-sim API URL |
| `SNOWFLAKE_ACCOUNT` | Yes (prod) | - | Snowflake account identifier |
| `SNOWFLAKE_USER` | Yes (prod) | - | Snowflake username |
| `SNOWFLAKE_PASSWORD` | Yes (prod) | - | Snowflake password |
| `SNOWFLAKE_DATABASE` | No | `CURSOR_ANALYTICS` | Snowflake database name |
| `SNOWFLAKE_SCHEMA` | No | `MART` | Snowflake schema name |
| `SNOWFLAKE_WAREHOUSE` | No | `TRANSFORM_WH` | Snowflake warehouse name |

## Development Status

| Task | Status | Description |
|------|--------|-------------|
| TASK-P9-01 | ✅ COMPLETE | Infrastructure setup |
| TASK-P9-02 | ✅ COMPLETE | Streamlit config |
| TASK-P9-03 | ✅ COMPLETE | Database connector |
| TASK-P9-04 | ✅ COMPLETE | SQL query modules |
| TASK-P9-05 | ✅ COMPLETE | Sidebar component |
| TASK-P9-06 | ✅ COMPLETE | Home page |
| TASK-P9-07 | ✅ COMPLETE | Velocity page |
| TASK-P9-08 | ✅ COMPLETE | AI Impact page |
| TASK-P9-09 | ✅ COMPLETE | Quality/Review pages |
| TASK-P9-10 | ✅ COMPLETE | Refresh pipeline |
| TASK-P9-11 | ✅ COMPLETE | Dockerfile |
| TASK-P9-12 | ✅ COMPLETE | Docker Compose |

**All 12 tasks completed on January 9, 2026.**

## Documentation

- [User Story](/.work-items/P9-F01-streamlit-dashboard/user-story.md)
- [Technical Design](/.work-items/P9-F01-streamlit-dashboard/design.md)
- [Task Breakdown](/.work-items/P9-F01-streamlit-dashboard/task.md)

## License

Copyright 2026 DOXAPI. All rights reserved.
