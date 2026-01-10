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
        "commitHash" AS commit_hash,
        "repoName" AS repo_name,
        "userEmail" AS user_email,

        -- AI metrics
        "tabLinesAdded" AS tab_lines_added,
        "composerLinesAdded" AS composer_lines_added,
        "nonAiLinesAdded" AS non_ai_lines_added,
        ("tabLinesAdded" + "composerLinesAdded") AS ai_lines_added,
        ("tabLinesAdded" + "composerLinesAdded" + "nonAiLinesAdded") AS total_lines_added,

        -- Timestamps (cast VARCHAR to TIMESTAMP)
        TRY_CAST("commitTs" AS TIMESTAMP) AS committed_at,

        -- Optional fields
        "branchName" AS branch_name

    FROM source
    WHERE "commitHash" IS NOT NULL
      AND "commitTs" IS NOT NULL
)

SELECT * FROM cleaned
