-- Mart: AI impact analysis by usage band
-- Compares outcomes across low/medium/high AI usage

{{ config(
    materialized='table',
    schema='mart'
) }}

WITH pr_data AS (
    SELECT
        *,
        CASE
            WHEN final_ai_ratio < 0.3 THEN 'low'
            WHEN final_ai_ratio < 0.6 THEN 'medium'
            ELSE 'high'
        END AS ai_usage_band
    FROM {{ ref('int_pr_with_commits') }}
    WHERE merged_at IS NOT NULL
)

SELECT
    ai_usage_band,
    {{ date_trunc_week('merged_at') }} AS week,

    -- Volume
    COUNT(*) AS pr_count,
    AVG(final_ai_ratio) AS avg_ai_ratio,

    -- Cycle times
    AVG(total_cycle_time_hours) AS avg_total_cycle_time,

    -- Quality indicators
    AVG(CASE WHEN is_reverted THEN 1 ELSE 0 END) AS revert_rate,
    AVG(CASE WHEN is_bug_fix THEN 1 ELSE 0 END) AS bug_fix_rate,

    -- Size metrics
    AVG(total_loc) AS avg_pr_size,
    AVG(changed_files) AS avg_files_changed

FROM pr_data
GROUP BY 1, 2
ORDER BY 2 DESC, 1
