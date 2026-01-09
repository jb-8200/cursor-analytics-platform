-- Staging model for pull requests from raw data
-- Cleans and calculates cycle time metrics

{{ config(
    materialized='view',
    schema='staging'
) }}

WITH source AS (
    SELECT * FROM {{ source('raw', 'pull_requests') }}
),

calculated AS (
    SELECT
        -- Primary keys
        number AS pr_number,
        repo_name,
        author_email,

        -- PR metadata
        state,
        additions,
        deletions,
        (additions + deletions) AS total_loc,
        changed_files,
        ai_ratio,

        -- Quality flags (rename was_reverted to is_reverted for consistency)
        was_reverted AS is_reverted,
        is_bug_fix,

        -- Timestamps
        created_at,
        merged_at,
        first_commit_at,
        first_review_at,

        -- Calculate cycle time metrics (not from API)
        -- Coding lead time: Time from first commit to PR creation
        CASE
            WHEN first_commit_at IS NOT NULL AND created_at IS NOT NULL
            THEN EXTRACT(EPOCH FROM (created_at - first_commit_at)) / 3600.0
            ELSE NULL
        END AS coding_lead_time_hours,

        -- Pickup time: Time from PR creation to first review
        CASE
            WHEN first_review_at IS NOT NULL AND created_at IS NOT NULL
            THEN EXTRACT(EPOCH FROM (first_review_at - created_at)) / 3600.0
            ELSE NULL
        END AS pickup_time_hours,

        -- Review lead time: Time from first review to merge
        CASE
            WHEN merged_at IS NOT NULL AND first_review_at IS NOT NULL
            THEN EXTRACT(EPOCH FROM (merged_at - first_review_at)) / 3600.0
            ELSE NULL
        END AS review_lead_time_hours

    FROM source
    WHERE number IS NOT NULL
      AND repo_name IS NOT NULL
      AND created_at IS NOT NULL
)

SELECT * FROM calculated
