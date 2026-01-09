-- ============================================================================
-- Snowflake COPY INTO Script
-- ============================================================================
-- Purpose: Load Parquet files from GCS external stage into raw tables
-- Run: snowsql -f copy_raw_tables.sql
-- Idempotency: COPY INTO is idempotent (tracks file metadata)
-- ============================================================================

USE DATABASE cursor_analytics;
USE SCHEMA raw;

-- ============================================================================
-- COPY INTO repos
-- ============================================================================
-- Load repository metadata from Parquet files
COPY INTO cursor_analytics.raw.repos (
    full_name,
    name,
    owner,
    default_branch,
    description,
    language,
    created_at,
    updated_at,
    pushed_at,
    size_kb,
    stargazers_count,
    watchers_count,
    forks_count,
    open_issues_count,
    private,
    fork,
    archived,
    disabled
)
FROM (
    SELECT
        $1:full_name::VARCHAR,
        $1:name::VARCHAR,
        $1:owner::VARCHAR,
        $1:default_branch::VARCHAR,
        $1:description::VARCHAR,
        $1:language::VARCHAR,
        $1:created_at::TIMESTAMP_NTZ,
        $1:updated_at::TIMESTAMP_NTZ,
        $1:pushed_at::TIMESTAMP_NTZ,
        $1:size_kb::INTEGER,
        $1:stargazers_count::INTEGER,
        $1:watchers_count::INTEGER,
        $1:forks_count::INTEGER,
        $1:open_issues_count::INTEGER,
        $1:private::BOOLEAN,
        $1:fork::BOOLEAN,
        $1:archived::BOOLEAN,
        $1:disabled::BOOLEAN
    FROM @cursor_analytics.raw.gcs_stage/repos/
)
FILE_FORMAT = (TYPE = PARQUET)
MATCH_BY_COLUMN_NAME = CASE_INSENSITIVE
ON_ERROR = CONTINUE;

-- Report results
SELECT 'repos' AS table_name, COUNT(*) AS row_count
FROM cursor_analytics.raw.repos;

-- ============================================================================
-- COPY INTO commits
-- ============================================================================
-- Load commit-level AI telemetry from Parquet files
COPY INTO cursor_analytics.raw.commits (
    commit_hash,
    user_id,
    user_email,
    user_name,
    repo_name,
    branch_name,
    is_primary_branch,
    total_lines_added,
    total_lines_deleted,
    tab_lines_added,
    tab_lines_deleted,
    composer_lines_added,
    composer_lines_deleted,
    non_ai_lines_added,
    non_ai_lines_deleted,
    message,
    commit_ts,
    created_at
)
FROM (
    SELECT
        $1:commit_hash::VARCHAR,
        $1:user_id::VARCHAR,
        $1:user_email::VARCHAR,
        $1:user_name::VARCHAR,
        $1:repo_name::VARCHAR,
        $1:branch_name::VARCHAR,
        $1:is_primary_branch::BOOLEAN,
        $1:total_lines_added::INTEGER,
        $1:total_lines_deleted::INTEGER,
        $1:tab_lines_added::INTEGER,
        $1:tab_lines_deleted::INTEGER,
        $1:composer_lines_added::INTEGER,
        $1:composer_lines_deleted::INTEGER,
        $1:non_ai_lines_added::INTEGER,
        $1:non_ai_lines_deleted::INTEGER,
        $1:message::VARCHAR,
        $1:commit_ts::TIMESTAMP_NTZ,
        $1:created_at::TIMESTAMP_NTZ
    FROM @cursor_analytics.raw.gcs_stage/commits/
)
FILE_FORMAT = (TYPE = PARQUET)
MATCH_BY_COLUMN_NAME = CASE_INSENSITIVE
ON_ERROR = CONTINUE;

-- Report results
SELECT 'commits' AS table_name, COUNT(*) AS row_count
FROM cursor_analytics.raw.commits;

-- ============================================================================
-- COPY INTO pull_requests
-- ============================================================================
-- Load PR lifecycle data from Parquet files
COPY INTO cursor_analytics.raw.pull_requests (
    number,
    repo_name,
    title,
    state,
    author_email,
    author_login,
    additions,
    deletions,
    changed_files,
    ai_ratio,
    ai_lines_total,
    human_lines_total,
    was_reverted,
    is_bug_fix,
    is_hotfix,
    created_at,
    updated_at,
    closed_at,
    merged_at,
    first_commit_at,
    first_review_at,
    head_ref,
    base_ref,
    commit_count,
    review_count,
    reviewers,
    labels,
    html_url,
    diff_url,
    patch_url
)
FROM (
    SELECT
        $1:number::INTEGER,
        $1:repo_name::VARCHAR,
        $1:title::VARCHAR,
        $1:state::VARCHAR,
        $1:author_email::VARCHAR,
        $1:author_login::VARCHAR,
        $1:additions::INTEGER,
        $1:deletions::INTEGER,
        $1:changed_files::INTEGER,
        $1:ai_ratio::FLOAT,
        $1:ai_lines_total::INTEGER,
        $1:human_lines_total::INTEGER,
        $1:was_reverted::BOOLEAN,
        $1:is_bug_fix::BOOLEAN,
        $1:is_hotfix::BOOLEAN,
        $1:created_at::TIMESTAMP_NTZ,
        $1:updated_at::TIMESTAMP_NTZ,
        $1:closed_at::TIMESTAMP_NTZ,
        $1:merged_at::TIMESTAMP_NTZ,
        $1:first_commit_at::TIMESTAMP_NTZ,
        $1:first_review_at::TIMESTAMP_NTZ,
        $1:head_ref::VARCHAR,
        $1:base_ref::VARCHAR,
        $1:commit_count::INTEGER,
        $1:review_count::INTEGER,
        $1:reviewers::ARRAY,
        $1:labels::ARRAY,
        $1:html_url::VARCHAR,
        $1:diff_url::VARCHAR,
        $1:patch_url::VARCHAR
    FROM @cursor_analytics.raw.gcs_stage/pull_requests/
)
FILE_FORMAT = (TYPE = PARQUET)
MATCH_BY_COLUMN_NAME = CASE_INSENSITIVE
ON_ERROR = CONTINUE;

-- Report results
SELECT 'pull_requests' AS table_name, COUNT(*) AS row_count
FROM cursor_analytics.raw.pull_requests;

-- ============================================================================
-- COPY INTO reviews
-- ============================================================================
-- Load PR review events from Parquet files
COPY INTO cursor_analytics.raw.reviews (
    id,
    repo_name,
    pr_number,
    state,
    body,
    user_login,
    user_email,
    submitted_at,
    commit_id,
    html_url,
    pull_request_url
)
FROM (
    SELECT
        $1:id::INTEGER,
        $1:repo_name::VARCHAR,
        $1:pr_number::INTEGER,
        $1:state::VARCHAR,
        $1:body::VARCHAR,
        $1:user_login::VARCHAR,
        $1:user_email::VARCHAR,
        $1:submitted_at::TIMESTAMP_NTZ,
        $1:commit_id::VARCHAR,
        $1:html_url::VARCHAR,
        $1:pull_request_url::VARCHAR
    FROM @cursor_analytics.raw.gcs_stage/reviews/
)
FILE_FORMAT = (TYPE = PARQUET)
MATCH_BY_COLUMN_NAME = CASE_INSENSITIVE
ON_ERROR = CONTINUE;

-- Report results
SELECT 'reviews' AS table_name, COUNT(*) AS row_count
FROM cursor_analytics.raw.reviews;

-- ============================================================================
-- Summary: All tables loaded
-- ============================================================================
SELECT 'COPY INTO completed' AS status,
       (SELECT COUNT(*) FROM cursor_analytics.raw.repos) AS repos_count,
       (SELECT COUNT(*) FROM cursor_analytics.raw.commits) AS commits_count,
       (SELECT COUNT(*) FROM cursor_analytics.raw.pull_requests) AS prs_count,
       (SELECT COUNT(*) FROM cursor_analytics.raw.reviews) AS reviews_count;

-- ============================================================================
-- Optional: Check for errors during load
-- ============================================================================
-- Uncomment to view load history and any errors:
-- SELECT * FROM TABLE(INFORMATION_SCHEMA.COPY_HISTORY(
--     TABLE_NAME => 'CURSOR_ANALYTICS.RAW.COMMITS',
--     START_TIME => DATEADD(HOURS, -1, CURRENT_TIMESTAMP())
-- ));
