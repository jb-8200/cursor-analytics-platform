-- Test: Cycle times should be non-negative
-- If any PRs have negative cycle times, something is wrong with timestamp calculations

SELECT
    pr_number,
    repo_name,
    total_cycle_time_hours
FROM {{ ref('stg_pull_requests') }}
WHERE (total_cycle_time_hours IS NOT NULL AND total_cycle_time_hours < 0)
