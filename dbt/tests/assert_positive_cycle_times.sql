-- Test: Cycle times should be non-negative
-- If any PRs have negative cycle times, something is wrong with timestamp calculations

SELECT
    pr_number,
    repo_name,
    coding_lead_time_hours,
    pickup_time_hours,
    review_lead_time_hours
FROM {{ ref('stg_pull_requests') }}
WHERE (coding_lead_time_hours IS NOT NULL AND coding_lead_time_hours < 0)
   OR (pickup_time_hours IS NOT NULL AND pickup_time_hours < 0)
   OR (review_lead_time_hours IS NOT NULL AND review_lead_time_hours < 0)
