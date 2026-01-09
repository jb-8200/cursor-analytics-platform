-- Mart: Review cost analysis
-- Estimates time/effort spent in code review

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
        COUNT(DISTINCT reviewer_email) AS reviewer_count
    FROM {{ ref('stg_reviews') }}
    GROUP BY repo_name, pr_number
)

SELECT
    DATE_TRUNC('week', p.merged_at) AS week,
    p.repo_name,

    -- Volume
    COUNT(*) AS total_prs,

    -- Review effort (simple time-based estimates)
    AVG(p.review_lead_time_hours) AS avg_review_cycle_time,
    AVG(COALESCE(r.review_count, 0)) AS avg_review_rounds,
    AVG(COALESCE(r.reviewer_count, 0)) AS avg_reviewers_per_pr,

    -- Estimated review hours (rough heuristic: 1 review = 15 min base + 5 min per 100 LOC)
    AVG(
        COALESCE(r.review_count, 0) * (0.25 + (p.total_loc / 100.0 * 0.08))
    ) AS estimated_review_hours_per_pr,

    SUM(
        COALESCE(r.review_count, 0) * (0.25 + (p.total_loc / 100.0 * 0.08))
    ) AS estimated_total_review_hours,

    -- PR size distribution
    AVG(p.total_loc) AS avg_pr_size,
    SUM(CASE WHEN p.total_loc > 500 THEN 1 ELSE 0 END) AS large_prs,
    AVG(CASE WHEN p.total_loc > 500 THEN p.total_loc ELSE NULL END) AS avg_large_pr_size

FROM pr_data p
LEFT JOIN review_data r
    ON p.repo_name = r.repo_name
    AND p.pr_number = r.pr_number
GROUP BY 1, 2
ORDER BY 1 DESC, 2
