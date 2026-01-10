-- Intermediate model: Enriched PR data with commit aggregations
-- Ephemeral model used by marts for enriched PR data
-- Note: Commit aggregations come from PR API response, not individual commits
--       (cursor-sim API doesn't include pull_request_number in commits)

{{ config(
    materialized='ephemeral'
) }}

WITH prs AS (
    SELECT * FROM {{ ref('stg_pull_requests') }}
)

SELECT
    -- PR fields
    pr_number,
    repo_name,
    author_email,
    state,
    additions,
    deletions,
    total_loc,
    changed_files,
    ai_ratio,
    is_reverted,
    is_bug_fix,
    created_at,
    merged_at,
    total_cycle_time_hours,
    reviewer_count,

    -- Commit aggregations (from PR API response)
    commit_count,
    tab_lines AS pr_tab_lines,
    composer_lines AS pr_composer_lines,
    additions AS pr_total_lines,  -- Use additions from PR
    (tab_lines + composer_lines) AS pr_ai_lines,
    (additions - tab_lines - composer_lines) AS pr_non_ai_lines,

    -- Use AI ratio from API
    ai_ratio AS final_ai_ratio

FROM prs
