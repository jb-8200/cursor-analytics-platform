-- Staging model for repositories from raw data
-- Cleans and normalizes repository metadata

{{ config(
    materialized='view',
    schema='staging'
) }}

WITH source AS (
    SELECT * FROM {{ source('raw', 'repos') }}
),

cleaned AS (
    SELECT
        -- Primary key
        full_name AS repo_name,

        -- Repository metadata
        default_branch,
        language,
        is_private,

        -- Timestamps
        created_at,
        updated_at

    FROM source
    WHERE full_name IS NOT NULL
)

SELECT * FROM cleaned
