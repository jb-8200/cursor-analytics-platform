-- ============================================================================
-- Snowflake Stage Setup Script
-- ============================================================================
-- Purpose: Create external stage for GCS bucket containing Parquet files
-- Run: snowsql -f setup_stages.sql
-- Idempotency: CREATE STAGE IF NOT EXISTS (safe to run multiple times)
-- ============================================================================

USE DATABASE cursor_analytics;
USE SCHEMA raw;

-- Create external stage pointing to GCS bucket
-- Note: Requires gcs_integration to be created first (see README.md)
CREATE STAGE IF NOT EXISTS cursor_analytics.raw.gcs_stage
  STORAGE_INTEGRATION = gcs_integration
  URL = 'gcs://cursor-analytics-data/raw/'
  COMMENT = 'External stage for Parquet files from cursor-sim extraction';

-- Verify stage creation
SHOW STAGES IN SCHEMA cursor_analytics.raw;

-- Optional: List files in stage to verify access
-- Uncomment to test connectivity:
-- LIST @cursor_analytics.raw.gcs_stage;
