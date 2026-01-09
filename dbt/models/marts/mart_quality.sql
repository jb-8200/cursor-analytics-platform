-- Mart: Quality metrics by repo and week
-- Tracks reverts, bug fixes, and review engagement

{{ config(
    materialized='table',
    schema='mart'
) }}

WITH pr_data AS (
    SELECT * FROM {{ ref('int_pr_with_commits') }}
    WHERE merged_at IS NOT NULL
),

review_data AS (
    SELECT
        repo_name,
        pr_number,
        COUNT(*) AS review_count,
        SUM(CASE WHEN is_approval THEN 1 ELSE 0 END) AS approval_count
    FROM {{ ref('stg_reviews') }}
    GROUP BY repo_name, pr_number
)

SELECT
    DATE_TRUNC('week', p.merged_at) AS week,
    p.repo_name,

    -- Volume
    COUNT(*) AS total_prs,

    -- Quality indicators
    SUM(CASE WHEN p.is_reverted THEN 1 ELSE 0 END) AS reverted_prs,
    AVG(CASE WHEN p.is_reverted THEN 1 ELSE 0 END) AS revert_rate,
    SUM(CASE WHEN p.is_bug_fix THEN 1 ELSE 0 END) AS bug_fix_prs,
    AVG(CASE WHEN p.is_bug_fix THEN 1 ELSE 0 END) AS bug_fix_rate,

    -- Review engagement
    AVG(COALESCE(r.review_count, 0)) AS avg_reviews_per_pr,
    AVG(COALESCE(r.approval_count, 0)) AS avg_approvals_per_pr,
    SUM(CASE WHEN COALESCE(r.review_count, 0) = 0 THEN 1 ELSE 0 END) AS unreviewed_prs,

    -- Size correlation
    AVG(CASE WHEN p.is_reverted THEN p.total_loc ELSE NULL END) AS avg_reverted_pr_size,
    AVG(CASE WHEN NOT p.is_reverted THEN p.total_loc ELSE NULL END) AS avg_successful_pr_size

FROM pr_data p
LEFT JOIN review_data r
    ON p.repo_name = r.repo_name
    AND p.pr_number = r.pr_number
GROUP BY 1, 2
ORDER BY 1 DESC, 2
