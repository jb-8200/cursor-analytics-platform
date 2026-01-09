{% macro array_length(array_field) %}
    {%- if target.type == 'duckdb' -%}
        ARRAY_LENGTH({{ array_field }})
    {%- elif target.type == 'snowflake' -%}
        ARRAY_SIZE({{ array_field }})
    {%- else -%}
        ARRAY_LENGTH({{ array_field }}, 1)
    {%- endif -%}
{% endmacro %}
