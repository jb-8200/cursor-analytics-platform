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
        default_branch

    FROM source
    WHERE full_name IS NOT NULL
)

SELECT * FROM cleaned
