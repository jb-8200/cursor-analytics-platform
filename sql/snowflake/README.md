# Snowflake Production Loading Scripts

SQL scripts for loading data into Snowflake in production environment.

## Prerequisites

1. Snowflake account with proper permissions
2. GCS storage integration configured (see setup below)
3. Parquet files uploaded to GCS bucket

## Setup

### 1. Create Storage Integration (One-time setup)

```sql
-- Run this as ACCOUNTADMIN
USE ROLE ACCOUNTADMIN;

CREATE STORAGE INTEGRATION IF NOT EXISTS gcs_integration
  TYPE = EXTERNAL_STAGE
  STORAGE_PROVIDER = 'GCS'
  ENABLED = TRUE
  STORAGE_ALLOWED_LOCATIONS = ('gcs://cursor-analytics-data/raw/');

-- Grant usage to SYSADMIN role
GRANT USAGE ON INTEGRATION gcs_integration TO ROLE SYSADMIN;
```

### 2. Get Service Account for GCS

```sql
-- Retrieve the service account email
DESC STORAGE INTEGRATION gcs_integration;
```

Copy the `STORAGE_GCP_SERVICE_ACCOUNT` value and grant it permissions in GCP:
- IAM: Storage Object Viewer on the bucket

### 3. Create Database and Schema

```sql
CREATE DATABASE IF NOT EXISTS cursor_analytics;
CREATE SCHEMA IF NOT EXISTS cursor_analytics.raw;
```

## Execution Order

Run scripts in the following order:

### Step 1: Create Stage
```bash
snowsql -f setup_stages.sql
```

This creates the GCS external stage pointing to the bucket.

### Step 2: Create Raw Tables
```bash
snowsql -f setup_raw_tables.sql
```

This creates the raw schema tables (repos, commits, pull_requests, reviews).

### Step 3: Load Data
```bash
snowsql -f copy_raw_tables.sql
```

This loads Parquet files from GCS into Snowflake tables using COPY INTO.

## Idempotency

All scripts are idempotent and can be run multiple times:
- `CREATE OR REPLACE TABLE` for tables
- `CREATE STAGE IF NOT EXISTS` for stages
- COPY INTO is automatically idempotent (checks file metadata)

## Incremental Loading

For incremental loads after the initial load:

```sql
-- Loads only new files not previously loaded
COPY INTO cursor_analytics.raw.commits
  FROM @cursor_analytics.raw.gcs_stage/commits/
  FILE_FORMAT = (TYPE = PARQUET)
  MATCH_BY_COLUMN_NAME = CASE_INSENSITIVE;
```

Snowflake tracks which files have been loaded and skips duplicates.

## Verification

After loading, verify data:

```sql
-- Check row counts
SELECT COUNT(*) FROM cursor_analytics.raw.repos;
SELECT COUNT(*) FROM cursor_analytics.raw.commits;
SELECT COUNT(*) FROM cursor_analytics.raw.pull_requests;
SELECT COUNT(*) FROM cursor_analytics.raw.reviews;

-- Sample data
SELECT * FROM cursor_analytics.raw.commits LIMIT 10;
SELECT * FROM cursor_analytics.raw.pull_requests LIMIT 10;
```

## Troubleshooting

### Error: Cannot find file

Check that files exist in GCS:
```sql
LIST @cursor_analytics.raw.gcs_stage/commits/;
```

### Error: Column mismatch

Verify Parquet schema matches table DDL:
```sql
-- Inspect Parquet file structure
SELECT $1 FROM @cursor_analytics.raw.gcs_stage/commits/ (FILE_FORMAT => 'PARQUET') LIMIT 1;
```

### Error: Permission denied

Check storage integration permissions:
```sql
DESC STORAGE INTEGRATION gcs_integration;
```

Ensure the service account has Storage Object Viewer role in GCP.

## Environment Variables

For automated runs with snowsql:

```bash
export SNOWFLAKE_ACCOUNT="your_account"
export SNOWFLAKE_USER="your_user"
export SNOWFLAKE_PASSWORD="your_password"
export SNOWFLAKE_WAREHOUSE="your_warehouse"
export SNOWFLAKE_DATABASE="cursor_analytics"
export SNOWFLAKE_SCHEMA="raw"

snowsql -f setup_stages.sql
```

## Notes

- MATCH_BY_COLUMN_NAME=CASE_INSENSITIVE allows flexible column name matching
- Parquet column names are case-sensitive, but this setting makes Snowflake more forgiving
- All timestamps are stored as TIMESTAMP_NTZ (no timezone) for consistency
- Primary keys are enforced for data quality (repos and commits)
- Composite keys for pull_requests (repo_name, number) and reviews (repo_name, pr_number, id)
