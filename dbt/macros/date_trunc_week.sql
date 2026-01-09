{% macro date_trunc_week(timestamp_field) %}
    DATE_TRUNC('week', {{ timestamp_field }})
{% endmacro %}
