-- Staging model for PR reviews from raw data
-- Cleans and normalizes review event data

{{ config(
    materialized='view',
    schema='staging'
) }}

WITH source AS (
    SELECT * FROM {{ source('raw', 'reviews') }}
),

cleaned AS (
    SELECT
        -- Primary keys
        id AS review_id,
        repo_name,
        pr_number,

        -- Reviewer information
        reviewer,
        state,

        -- Derived flags
        CASE
            WHEN UPPER(state) = 'APPROVED' THEN TRUE
            ELSE FALSE
        END AS is_approval,

        -- Timestamps (cast VARCHAR to TIMESTAMP)
        TRY_CAST(submitted_at AS TIMESTAMP) AS submitted_at,

        -- Metadata
        body AS review_comment

    FROM source
    WHERE id IS NOT NULL
      AND repo_name IS NOT NULL
      AND pr_number IS NOT NULL
)

SELECT * FROM cleaned
