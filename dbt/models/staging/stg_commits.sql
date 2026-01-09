-- Staging model for commits from raw data
-- Cleans and normalizes commit-level AI telemetry data

{{ config(
    materialized='view',
    schema='staging'
) }}

WITH source AS (
    SELECT * FROM {{ source('raw', 'commits') }}
),

cleaned AS (
    SELECT
        -- Primary keys
        commit_hash,
        repo_name,
        user_email,

        -- AI metrics
        tab_lines_added,
        composer_lines_added,
        non_ai_lines_added,
        (tab_lines_added + composer_lines_added) AS ai_lines_added,
        (tab_lines_added + composer_lines_added + non_ai_lines_added) AS total_lines_added,

        -- Timestamps
        commit_ts AS committed_at,

        -- Optional fields
        pull_request_number,
        branch_name

    FROM source
    WHERE commit_hash IS NOT NULL
      AND commit_ts IS NOT NULL
)

SELECT * FROM cleaned
