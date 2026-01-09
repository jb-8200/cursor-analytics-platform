{% macro percentile_cont(percentile, expression) %}
    {%- if target.type == 'duckdb' -%}
        PERCENTILE_CONT({{ percentile }}) WITHIN GROUP (ORDER BY {{ expression }})
    {%- elif target.type == 'snowflake' -%}
        PERCENTILE_CONT({{ percentile }}) WITHIN GROUP (ORDER BY {{ expression }})
    {%- else -%}
        PERCENTILE_CONT({{ percentile }}) WITHIN GROUP (ORDER BY {{ expression }})
    {%- endif -%}
{% endmacro %}
