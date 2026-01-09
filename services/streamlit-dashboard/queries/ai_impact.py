"""
AI Impact analysis SQL queries.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-04 SQL Query Modules

This module provides parameterized SQL queries for AI impact analysis:
- Metrics grouped by AI usage bands (low, medium, high)
- Comparison of cycle times and quality by AI usage
- Band-specific trends

Depends on: mart.ai_impact (P8 dbt mart)
"""

from db.connector import query
import pandas as pd


def get_ai_impact_data(where_clause: str = "") -> pd.DataFrame:
    """
    Get AI impact metrics grouped by usage band.

    Args:
        where_clause: Optional SQL WHERE clause for filtering (e.g., "WHERE week >= '2026-01-01'")

    Returns:
        DataFrame with AI impact metrics including:
        - week: Week start date
        - ai_usage_band: AI usage band (low, medium, high)
        - pr_count: Number of PRs in this band
        - avg_ai_ratio: Average AI contribution ratio
        - avg_coding_lead_time: Average coding lead time (days)
        - avg_review_cycle_time: Average review cycle time (days)
        - revert_rate: Revert rate for this band

    Example:
        >>> df = get_ai_impact_data()
        >>> df = get_ai_impact_data("WHERE week >= '2026-01-01'")
    """
    sql = f"""
    SELECT
        week,
        ai_usage_band,
        pr_count,
        avg_ai_ratio,
        avg_coding_lead_time,
        avg_review_cycle_time,
        revert_rate
    FROM mart.ai_impact
    {where_clause}
    ORDER BY week DESC, ai_usage_band
    """
    return query(sql)


def get_band_comparison(where_clause: str = "") -> pd.DataFrame:
    """
    Get comparison of metrics across AI usage bands.

    Args:
        where_clause: Optional SQL WHERE clause for filtering

    Returns:
        DataFrame aggregated by band with average metrics:
        - ai_usage_band: Band name (low, medium, high)
        - total_prs: Total PRs in band
        - avg_ai_ratio: Average AI ratio
        - avg_coding_lead_time: Average coding lead time
        - avg_review_cycle_time: Average review cycle time
        - avg_revert_rate: Average revert rate

    Example:
        >>> comparison = get_band_comparison()
        >>> print(comparison[comparison['ai_usage_band'] == 'high'])
    """
    sql = f"""
    SELECT
        ai_usage_band,
        SUM(pr_count) as total_prs,
        AVG(avg_ai_ratio) as avg_ai_ratio,
        AVG(avg_coding_lead_time) as avg_coding_lead_time,
        AVG(avg_review_cycle_time) as avg_review_cycle_time,
        AVG(revert_rate) as avg_revert_rate
    FROM mart.ai_impact
    {where_clause}
    GROUP BY ai_usage_band
    ORDER BY ai_usage_band
    """
    return query(sql)


def get_band_trends(where_clause: str = "") -> pd.DataFrame:
    """
    Get week-over-week trends for each AI usage band.

    Args:
        where_clause: Optional SQL WHERE clause for filtering

    Returns:
        DataFrame with weekly trends by band:
        - week: Week start date
        - ai_usage_band: Band name
        - pr_count: PRs in this week/band
        - avg_coding_lead_time: Coding lead time for this week/band
        - revert_rate: Revert rate for this week/band

    Example:
        >>> trends = get_band_trends()
        >>> high_ai = trends[trends['ai_usage_band'] == 'high']
    """
    sql = f"""
    SELECT
        week,
        ai_usage_band,
        pr_count,
        avg_coding_lead_time,
        revert_rate
    FROM mart.ai_impact
    {where_clause}
    ORDER BY week DESC, ai_usage_band
    """
    return query(sql)
