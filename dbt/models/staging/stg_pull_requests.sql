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

        -- Commit aggregations (from API)
        commit_count,
        tab_lines,
        composer_lines,

        -- Quality flags (rename was_reverted to is_reverted for consistency)
        was_reverted AS is_reverted,
        is_bug_fix,

        -- Timestamps (cast VARCHAR to TIMESTAMP)
        TRY_CAST(created_at AS TIMESTAMP) AS created_at,
        TRY_CAST(merged_at AS TIMESTAMP) AS merged_at,

        -- Calculate cycle time metrics (not from API)
        -- Total cycle time: Time from PR creation to merge
        CASE
            WHEN merged_at IS NOT NULL AND created_at IS NOT NULL
            THEN EXTRACT(EPOCH FROM (TRY_CAST(merged_at AS TIMESTAMP) - TRY_CAST(created_at AS TIMESTAMP))) / 3600.0
            ELSE NULL
        END AS total_cycle_time_hours,

        -- Reviewer count (already an integer in API response)
        reviewers AS reviewer_count

    FROM source
    WHERE number IS NOT NULL
      AND repo_name IS NOT NULL
      AND created_at IS NOT NULL
)

SELECT * FROM calculated
