-- Mart: Weekly velocity metrics by repo
-- Pre-aggregated for dashboard consumption

{{ config(
    materialized='table',
    schema='mart'
) }}

WITH pr_data AS (
    SELECT * FROM {{ ref('int_pr_with_commits') }}
    WHERE merged_at IS NOT NULL
)

SELECT
    -- Time dimension (week starting Monday)
    DATE_TRUNC('week', merged_at) AS week,
    repo_name,

    -- Team productivity
    COUNT(DISTINCT author_email) AS active_developers,
    COUNT(*) AS total_prs,
    SUM(commit_count) AS total_commits,

    -- Size metrics
    AVG(total_loc) AS avg_pr_size,
    AVG(changed_files) AS avg_files_changed,

    -- Cycle time metrics (in hours)
    AVG(coding_lead_time_hours) AS avg_coding_lead_time,
    AVG(pickup_time_hours) AS avg_pickup_time,
    AVG(review_lead_time_hours) AS avg_review_lead_time,
    AVG(
        COALESCE(coding_lead_time_hours, 0) +
        COALESCE(pickup_time_hours, 0) +
        COALESCE(review_lead_time_hours, 0)
    ) AS avg_total_cycle_time,

    -- AI adoption
    AVG(final_ai_ratio) AS avg_ai_ratio,
    SUM(pr_ai_lines) AS total_ai_lines,
    SUM(pr_total_lines) AS total_lines

FROM pr_data
GROUP BY 1, 2
ORDER BY 1 DESC, 2
