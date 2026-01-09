# Streamlit Analytics Dashboard

**Feature ID**: P9-F01-streamlit-dashboard
**Created**: 2026-01-09
**Status**: COMPLETE ✅

## Overview

Production-ready Streamlit dashboard for visualizing AI coding analytics. Supports both DuckDB (local development) and Snowflake (production) backends.

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

## Setup

### Local Development (DuckDB)

```bash
# Install dependencies
pip install -r requirements.txt

# Set environment variables
export DB_MODE=duckdb
export DUCKDB_PATH=/data/analytics.duckdb
export CURSOR_SIM_URL=http://localhost:8080

# Run dashboard
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

```bash
# Build image
docker build -t streamlit-dashboard .

# Run container
docker run -p 8501:8501 \
  -e DB_MODE=duckdb \
  -e DUCKDB_PATH=/data/analytics.duckdb \
  -v /path/to/data:/data \
  streamlit-dashboard
```

## Testing

```bash
# Run tests
pytest tests/

# Run tests with coverage
pytest --cov=. --cov-report=html tests/
```

## Dependencies

- **P8 (Data Tier)**: Requires dbt mart tables (mart.velocity, mart.ai_impact, mart.quality, mart.review_costs)
- **cursor-sim (P4)**: Required for dev mode refresh
- **DuckDB**: Local analytics database
- **Snowflake**: Production data warehouse (optional)

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DB_MODE` | No | `duckdb` | Database mode: `duckdb` or `snowflake` |
| `DUCKDB_PATH` | No | `/data/analytics.duckdb` | Path to DuckDB file |
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
