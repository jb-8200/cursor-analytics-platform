-- ============================================================================
-- Snowflake Raw Tables DDL Script
-- ============================================================================
-- Purpose: Create raw schema tables for data from cursor-sim extraction
-- Run: snowsql -f setup_raw_tables.sql
-- Idempotency: CREATE OR REPLACE TABLE (safe to run multiple times)
-- ============================================================================

USE DATABASE cursor_analytics;
CREATE SCHEMA IF NOT EXISTS raw;
USE SCHEMA raw;

-- ============================================================================
-- repos: Repository metadata
-- Source: GET /repos endpoint
-- ============================================================================
CREATE OR REPLACE TABLE cursor_analytics.raw.repos (
    -- Identifiers
    full_name VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    owner VARCHAR(255),

    -- Repository details
    default_branch VARCHAR(100),
    description VARCHAR(1000),
    language VARCHAR(100),

    -- Metadata
    created_at TIMESTAMP_NTZ,
    updated_at TIMESTAMP_NTZ,
    pushed_at TIMESTAMP_NTZ,

    -- Statistics
    size_kb INTEGER,
    stargazers_count INTEGER,
    watchers_count INTEGER,
    forks_count INTEGER,
    open_issues_count INTEGER,

    -- Flags
    private BOOLEAN,
    fork BOOLEAN,
    archived BOOLEAN,
    disabled BOOLEAN,

    -- Load metadata
    loaded_at TIMESTAMP_NTZ DEFAULT CURRENT_TIMESTAMP(),

    PRIMARY KEY (full_name)
)
COMMENT = 'Repository metadata from /repos endpoint';

-- ============================================================================
-- commits: AI code tracking commit data
-- Source: GET /analytics/ai-code/commits endpoint (Cursor Analytics-style)
-- ============================================================================
CREATE OR REPLACE TABLE cursor_analytics.raw.commits (
    -- Primary identifiers
    commit_hash VARCHAR(64) NOT NULL,

    -- User information
    user_id VARCHAR(255) NOT NULL,
    user_email VARCHAR(255) NOT NULL,
    user_name VARCHAR(255),

    -- Repository context
    repo_name VARCHAR(255) NOT NULL,
    branch_name VARCHAR(255),
    is_primary_branch BOOLEAN,

    -- Line counts (total)
    total_lines_added INTEGER,
    total_lines_deleted INTEGER,

    -- Line counts (by AI source)
    tab_lines_added INTEGER,
    tab_lines_deleted INTEGER,
    composer_lines_added INTEGER,
    composer_lines_deleted INTEGER,
    non_ai_lines_added INTEGER,
    non_ai_lines_deleted INTEGER,

    -- Commit metadata
    message VARCHAR(5000),
    commit_ts TIMESTAMP_NTZ,
    created_at TIMESTAMP_NTZ,

    -- Load metadata
    loaded_at TIMESTAMP_NTZ DEFAULT CURRENT_TIMESTAMP(),

    PRIMARY KEY (commit_hash)
)
COMMENT = 'Commit-level AI code tracking data from /analytics/ai-code/commits';

-- ============================================================================
-- pull_requests: PR lifecycle data
-- Source: GET /repos/{owner}/{repo}/pulls endpoint (GitHub-style)
-- ============================================================================
CREATE OR REPLACE TABLE cursor_analytics.raw.pull_requests (
    -- Composite key
    number INTEGER NOT NULL,
    repo_name VARCHAR(255) NOT NULL,

    -- PR details
    title VARCHAR(1000),
    state VARCHAR(50),
    author_email VARCHAR(255),
    author_login VARCHAR(255),

    -- Code changes
    additions INTEGER,
    deletions INTEGER,
    changed_files INTEGER,

    -- AI metrics
    ai_ratio FLOAT,
    ai_lines_total INTEGER,
    human_lines_total INTEGER,

    -- Quality indicators
    was_reverted BOOLEAN,
    is_bug_fix BOOLEAN,
    is_hotfix BOOLEAN,

    -- Lifecycle timestamps
    created_at TIMESTAMP_NTZ,
    updated_at TIMESTAMP_NTZ,
    closed_at TIMESTAMP_NTZ,
    merged_at TIMESTAMP_NTZ,
    first_commit_at TIMESTAMP_NTZ,
    first_review_at TIMESTAMP_NTZ,

    -- Branch information
    head_ref VARCHAR(255),
    base_ref VARCHAR(255),

    -- Metadata
    commit_count INTEGER,
    review_count INTEGER,
    reviewers ARRAY,
    labels ARRAY,

    -- URLs
    html_url VARCHAR(500),
    diff_url VARCHAR(500),
    patch_url VARCHAR(500),

    -- Load metadata
    loaded_at TIMESTAMP_NTZ DEFAULT CURRENT_TIMESTAMP(),

    PRIMARY KEY (repo_name, number)
)
COMMENT = 'Pull request lifecycle data from /repos/{owner}/{repo}/pulls';

-- ============================================================================
-- reviews: PR review events
-- Source: GET /repos/{owner}/{repo}/pulls/{number}/reviews endpoint
-- ============================================================================
CREATE OR REPLACE TABLE cursor_analytics.raw.reviews (
    -- Composite key
    id INTEGER NOT NULL,
    repo_name VARCHAR(255) NOT NULL,
    pr_number INTEGER NOT NULL,

    -- Review details
    state VARCHAR(50),
    body VARCHAR(5000),

    -- Reviewer information
    user_login VARCHAR(255),
    user_email VARCHAR(255),

    -- Review metadata
    submitted_at TIMESTAMP_NTZ,
    commit_id VARCHAR(64),

    -- URLs
    html_url VARCHAR(500),
    pull_request_url VARCHAR(500),

    -- Load metadata
    loaded_at TIMESTAMP_NTZ DEFAULT CURRENT_TIMESTAMP(),

    PRIMARY KEY (repo_name, pr_number, id)
)
COMMENT = 'PR review events from /repos/{owner}/{repo}/pulls/{number}/reviews';

-- ============================================================================
-- Verification Queries
-- ============================================================================
-- Uncomment to verify table creation:
-- SHOW TABLES IN SCHEMA cursor_analytics.raw;
-- DESCRIBE TABLE cursor_analytics.raw.commits;
-- DESCRIBE TABLE cursor_analytics.raw.pull_requests;
