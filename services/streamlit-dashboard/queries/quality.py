"""
Code Quality metrics SQL queries.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-04 SQL Query Modules

This module provides parameterized SQL queries for quality metrics:
- Revert rates and trends
- Time to revert
- Bug fix rates
- Quality trends over time

Depends on: mart.quality (P8 dbt mart)
"""

from db.connector import query
import pandas as pd


def get_quality_data(where_clause: str = "") -> pd.DataFrame:
    """
    Get weekly quality metrics.

    Args:
        where_clause: Optional SQL WHERE clause for filtering (e.g., "WHERE repo_name = 'acme/platform'")

    Returns:
        DataFrame with quality metrics including:
        - week: Week start date
        - repo_name: Repository name
        - total_prs: Total PRs merged
        - reverted_prs: Number of PRs that were reverted
        - revert_rate: Revert rate (reverted_prs / total_prs)
        - avg_time_to_revert: Average time from merge to revert (days)
        - bug_fix_prs: Number of PRs marked as bug fixes
        - bug_fix_rate: Bug fix rate (bug_fix_prs / total_prs)

    Example:
        >>> df = get_quality_data()
        >>> df = get_quality_data("WHERE repo_name = 'acme/platform'")
    """
    sql = f"""
    SELECT
        week,
        repo_name,
        total_prs,
        reverted_prs,
        revert_rate,
        avg_time_to_revert,
        bug_fix_prs,
        bug_fix_rate
    FROM mart.quality
    {where_clause}
    ORDER BY week DESC
    """
    return query(sql)


def get_revert_trends(where_clause: str = "") -> pd.DataFrame:
    """
    Get revert rate trends over time.

    Args:
        where_clause: Optional SQL WHERE clause for filtering

    Returns:
        DataFrame with weekly revert trends:
        - week: Week start date
        - repo_name: Repository name
        - revert_rate: Revert rate for this week
        - reverted_prs: Number of reverted PRs

    Example:
        >>> trends = get_revert_trends()
        >>> recent = trends.head(12)  # Last 12 weeks
    """
    sql = f"""
    SELECT
        week,
        repo_name,
        revert_rate,
        reverted_prs
    FROM mart.quality
    {where_clause}
    ORDER BY week DESC
    """
    return query(sql)


def get_quality_summary(where_clause: str = "") -> dict:
    """
    Get summary statistics for quality metrics.

    Args:
        where_clause: Optional SQL WHERE clause for filtering

    Returns:
        Dictionary with summary stats:
        - total_prs: Total number of PRs
        - total_reverted: Total reverted PRs
        - avg_revert_rate: Average revert rate
        - avg_bug_fix_rate: Average bug fix rate
        - avg_time_to_revert: Average time to revert

    Example:
        >>> summary = get_quality_summary()
        >>> print(f"Revert rate: {summary['avg_revert_rate']:.2%}")
    """
    sql = f"""
    SELECT
        SUM(total_prs) as total_prs,
        SUM(reverted_prs) as total_reverted,
        AVG(revert_rate) as avg_revert_rate,
        AVG(bug_fix_rate) as avg_bug_fix_rate,
        AVG(avg_time_to_revert) as avg_time_to_revert
    FROM mart.quality
    {where_clause}
    """
    result = query(sql)
    return result.iloc[0].to_dict()


def get_quality_by_ai_band(where_clause: str = "") -> pd.DataFrame:
    """
    Get quality metrics grouped by AI usage band.

    Args:
        where_clause: Optional SQL WHERE clause for filtering

    Returns:
        DataFrame with quality by AI band:
        - ai_usage_band: AI usage band (low, medium, high)
        - avg_revert_rate: Average revert rate for this band
        - avg_bug_fix_rate: Average bug fix rate for this band

    Example:
        >>> by_band = get_quality_by_ai_band()
        >>> high_ai = by_band[by_band['ai_usage_band'] == 'high']
    """
    sql = f"""
    SELECT
        ai_usage_band,
        AVG(revert_rate) as avg_revert_rate,
        AVG(bug_fix_rate) as avg_bug_fix_rate
    FROM mart.quality
    {where_clause}
    GROUP BY ai_usage_band
    ORDER BY ai_usage_band
    """
    return query(sql)
