-- Intermediate model: Join PRs with aggregated commit data
-- Ephemeral model used by marts for enriched PR data

{{ config(
    materialized='ephemeral'
) }}

WITH prs AS (
    SELECT * FROM {{ ref('stg_pull_requests') }}
),

commit_agg AS (
    SELECT
        pull_request_number,
        repo_name,
        COUNT(*) AS commit_count,
        SUM(ai_lines_added) AS total_ai_lines,
        SUM(total_lines_added) AS total_lines,
        SUM(tab_lines_added) AS total_tab_lines,
        SUM(composer_lines_added) AS total_composer_lines,
        SUM(non_ai_lines_added) AS total_non_ai_lines,
        MIN(committed_at) AS first_commit_at_calc,
        MAX(committed_at) AS last_commit_at
    FROM {{ ref('stg_commits') }}
    WHERE pull_request_number IS NOT NULL
    GROUP BY pull_request_number, repo_name
)

SELECT
    -- PR fields
    p.*,

    -- Commit aggregations
    COALESCE(c.commit_count, 0) AS commit_count,
    COALESCE(c.total_ai_lines, 0) AS pr_ai_lines,
    COALESCE(c.total_lines, 0) AS pr_total_lines,
    COALESCE(c.total_tab_lines, 0) AS pr_tab_lines,
    COALESCE(c.total_composer_lines, 0) AS pr_composer_lines,
    COALESCE(c.total_non_ai_lines, 0) AS pr_non_ai_lines,

    -- Calculated AI ratio with fallbacks
    -- Prefer API-provided ai_ratio, fall back to calculated from commits
    COALESCE(
        p.ai_ratio,
        CASE
            WHEN c.total_lines > 0
            THEN c.total_ai_lines::FLOAT / c.total_lines
            ELSE 0
        END,
        0
    ) AS final_ai_ratio,

    -- Commit timestamps
    c.first_commit_at_calc,
    c.last_commit_at

FROM prs p
LEFT JOIN commit_agg c
    ON p.pr_number = c.pull_request_number
    AND p.repo_name = c.repo_name
