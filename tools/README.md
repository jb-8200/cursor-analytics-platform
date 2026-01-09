# Data Pipeline Tools

This directory contains the ETL pipeline orchestration and data extraction tools for the Cursor Analytics Platform.

## Overview

The data tier implements a three-stage ETL pipeline:

1. **Extract**: Pull data from cursor-sim REST API to Parquet files
2. **Load**: Load Parquet files into DuckDB (dev) or Snowflake (prod)
3. **Transform**: Run dbt models to create analytics-ready tables

## Quick Start

### Run Full Pipeline

```bash
# Using make (recommended)
make pipeline

# Or directly
./tools/run_pipeline.sh
```

### Run Individual Steps

```bash
# Extract only
make extract

# Load only
make load

# Transform only (dbt)
make dbt-build
```

## Configuration

All steps support environment variable configuration:

| Variable | Default | Description |
|----------|---------|-------------|
| `CURSOR_SIM_URL` | `http://localhost:8080` | Base URL for cursor-sim API |
| `DATA_DIR` | `./data` | Base directory for data files |
| `API_KEY` | `cursor-sim-dev-key` | API key for authentication |
| `START_DATE` | `90d` | Start date for commits filter |
| `CONTINUE_ON_ERROR` | `false` | Continue on step failures |

### Example with Custom Configuration

```bash
# Override cursor-sim URL
CURSOR_SIM_URL=http://staging:8080 make pipeline

# Use custom data directory
DATA_DIR=/tmp/analytics-data make pipeline

# Continue on errors
CONTINUE_ON_ERROR=true ./tools/run_pipeline.sh
```

## Pipeline Stages

### Stage 1: Extract

**What it does**: Fetches data from cursor-sim API endpoints and writes to Parquet files

**Endpoints**:
- `/repos` - Repository list
- `/analytics/ai-code/commits` - Commit-level AI telemetry
- `/repos/{owner}/{repo}/pulls` - Pull request data
- `/repos/{owner}/{repo}/pulls/{number}/reviews` - Review events

**Output**: `data/raw/*.parquet`

**CLI**:
```bash
python tools/api-loader/loader.py \
    --url http://localhost:8080 \
    --output data/raw \
    --api-key cursor-sim-dev-key \
    --start-date 90d
```

### Stage 2: Load

**What it does**: Loads Parquet files into DuckDB `raw` schema

**Input**: `data/raw/*.parquet`

**Output**: `data/analytics.duckdb` with `raw.repos`, `raw.commits`, `raw.pull_requests`, `raw.reviews`

**CLI**:
```bash
python tools/api-loader/duckdb_loader.py \
    --parquet-dir data/raw \
    --db-path data/analytics.duckdb
```

**Modes**:
- **Full refresh** (default): Replace tables with new data
- **Incremental**: Append to existing tables (`--incremental` flag)

### Stage 3: Transform

**What it does**: Runs dbt models to create analytics marts

**Models**:
- **Staging**: Clean and normalize raw data (`staging.stg_*`)
- **Intermediate**: Join PRs with commits (`int_pr_with_commits`)
- **Marts**: Pre-aggregated analytics tables (`mart.*`)
  - `mart.velocity` - Weekly PR cycle times
  - `mart.ai_impact` - Metrics by AI usage band
  - `mart.quality` - Revert and bug fix rates
  - `mart.review_costs` - Review iteration costs

**CLI**:
```bash
cd dbt
dbt deps
dbt build --target dev
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make pipeline` | Run full ETL pipeline |
| `make extract` | Extract data from cursor-sim API |
| `make load` | Load Parquet to DuckDB |
| `make dbt-deps` | Install dbt dependencies |
| `make dbt-build` | Run dbt models (build all) |
| `make dbt-test` | Run dbt tests only |
| `make dbt-run` | Run dbt models without tests |
| `make dbt-docs` | Generate and serve dbt docs |
| `make ci-local` | Full pipeline for local CI |
| `make clean-data` | Remove generated data files |
| `make query-duckdb` | Open DuckDB CLI |

## Querying Data

### DuckDB CLI

```bash
# Open DuckDB CLI
make query-duckdb

# Or directly
duckdb data/analytics.duckdb
```

### Example Queries

```sql
-- View velocity metrics
SELECT * FROM mart.velocity
WHERE week >= '2025-01-01'
ORDER BY week DESC;

-- View AI impact by usage band
SELECT
    ai_usage_band,
    COUNT(*) as pr_count,
    AVG(avg_coding_lead_time) as avg_coding_time
FROM mart.ai_impact
WHERE week >= '2025-01-01'
GROUP BY ai_usage_band;

-- View quality metrics
SELECT * FROM mart.quality
WHERE week >= '2025-01-01'
ORDER BY week DESC;
```

## Testing

### Run Pipeline Tests

```bash
# Test script structure and configuration
./tools/test_pipeline.sh
```

### Test Individual Components

```bash
# Test extractors
cd tools/api-loader
PYTHONPATH=. python -m pytest tests/test_base.py -v
PYTHONPATH=. python -m pytest tests/test_specific_extractors.py -v
PYTHONPATH=. python -m pytest tests/test_loader.py -v

# Test DuckDB loader
PYTHONPATH=. python -m pytest tests/test_duckdb_loader.py -v

# Test dbt models
cd dbt
dbt test --target dev
```

## Directory Structure

```
tools/
├── README.md                    # This file
├── run_pipeline.sh              # Main orchestration script
├── test_pipeline.sh             # Pipeline test suite
└── api-loader/                  # Data extraction
    ├── loader.py                # Main loader CLI
    ├── duckdb_loader.py         # DuckDB loading
    ├── extractors/              # API extractors
    │   ├── base.py              # Base extractor
    │   ├── repos.py             # Repos extractor
    │   ├── commits.py           # Commits extractor
    │   ├── prs.py               # PRs extractor
    │   └── reviews.py           # Reviews extractor
    └── tests/                   # Test suite

dbt/
├── dbt_project.yml              # dbt configuration
├── profiles.yml                 # Connection profiles
├── models/                      # dbt models
│   ├── sources.yml              # Raw data sources
│   ├── staging/                 # Staging models
│   ├── intermediate/            # Intermediate models
│   └── marts/                   # Analytics marts
└── macros/                      # Cross-engine macros

data/
├── raw/                         # Parquet files (gitignored)
└── analytics.duckdb             # DuckDB database (gitignored)
```

## Troubleshooting

### cursor-sim not responding

```bash
# Check if cursor-sim is running
curl http://localhost:8080/health

# Start cursor-sim if needed
cd services/cursor-sim
go run . --port 8080
```

### Empty Parquet files

Check cursor-sim has data:

```bash
# List repositories
curl http://localhost:8080/repos

# Check commits
curl http://localhost:8080/analytics/ai-code/commits?startDate=90d
```

### dbt fails

```bash
# Install dependencies
make dbt-deps

# Check DuckDB has raw tables
duckdb data/analytics.duckdb -c "SHOW TABLES FROM raw;"

# Run with verbose output
cd dbt && dbt build --target dev --debug
```

### Permission denied

```bash
# Make scripts executable
chmod +x tools/run_pipeline.sh
chmod +x tools/test_pipeline.sh
```

## Production Deployment

For production deployment to Cloud Run Jobs, see:
- Snowflake loading: `sql/snowflake/README.md`
- Docker jobs: `jobs/data-loader/` and `jobs/dbt-runner/`

## See Also

- [Data Tier Specification](../services/cursor-sim/SPEC.md)
- [dbt Documentation](../dbt/README.md)
- [Snowflake Setup](../sql/snowflake/README.md)
- [Task Breakdown](../.work-items/P8-F01-data-tier/task.md)
