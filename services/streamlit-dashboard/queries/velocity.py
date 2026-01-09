"""
Velocity metrics SQL queries.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-04 SQL Query Modules

This module provides parameterized SQL queries for velocity metrics:
- Weekly PR throughput and cycle times
- Developer activity
- AI usage ratios
- Cycle time breakdown by component

Depends on: mart.velocity (P8 dbt mart)
"""

from db.connector import query
import pandas as pd


def get_velocity_data(where_clause: str = "") -> pd.DataFrame:
    """
    Get weekly velocity metrics.

    Args:
        where_clause: Optional SQL WHERE clause for filtering (e.g., "WHERE repo_name = 'acme/platform'")

    Returns:
        DataFrame with velocity metrics including:
        - week: Week start date
        - repo_name: Repository name
        - active_developers: Count of active developers
        - total_prs: Total PRs merged
        - avg_pr_size: Average PR size (lines changed)
        - coding_lead_time: Average time from first commit to PR open (days)
        - pickup_time: Average time from PR open to first review (days)
        - review_lead_time: Average time from first review to merge (days)
        - total_cycle_time: Total cycle time (days)
        - p50_cycle_time: Median cycle time (days)
        - p90_cycle_time: 90th percentile cycle time (days)
        - avg_ai_ratio: Average AI contribution ratio

    Example:
        >>> df = get_velocity_data()
        >>> df = get_velocity_data("WHERE repo_name = 'acme/platform'")
    """
    sql = f"""
    SELECT
        week,
        repo_name,
        active_developers,
        total_prs,
        avg_pr_size,
        coding_lead_time,
        pickup_time,
        review_lead_time,
        total_cycle_time,
        p50_cycle_time,
        p90_cycle_time,
        avg_ai_ratio
    FROM mart.velocity
    {where_clause}
    ORDER BY week DESC
    """
    return query(sql)


def get_cycle_time_breakdown(where_clause: str = "") -> pd.DataFrame:
    """
    Get cycle time breakdown by component.

    Args:
        where_clause: Optional SQL WHERE clause for filtering

    Returns:
        DataFrame with columns:
        - component: Component name (Coding, Pickup, Review)
        - hours: Average hours for this component

    Example:
        >>> breakdown = get_cycle_time_breakdown()
        >>> breakdown = get_cycle_time_breakdown("WHERE repo_name = 'acme/platform'")
    """
    sql = f"""
    SELECT
        'Coding' as component,
        AVG(coding_lead_time) * 24 as hours
    FROM mart.velocity {where_clause}
    UNION ALL
    SELECT
        'Pickup' as component,
        AVG(pickup_time) * 24 as hours
    FROM mart.velocity {where_clause}
    UNION ALL
    SELECT
        'Review' as component,
        AVG(review_lead_time) * 24 as hours
    FROM mart.velocity {where_clause}
    """
    return query(sql)


def get_velocity_summary(where_clause: str = "") -> dict:
    """
    Get summary statistics for velocity metrics.

    Args:
        where_clause: Optional SQL WHERE clause for filtering

    Returns:
        Dictionary with summary stats:
        - total_prs: Total number of PRs
        - avg_cycle_time: Average cycle time
        - max_developers: Maximum active developers
        - avg_ai_ratio: Average AI ratio

    Example:
        >>> summary = get_velocity_summary()
        >>> print(summary['total_prs'])
    """
    sql = f"""
    SELECT
        SUM(total_prs) as total_prs,
        AVG(total_cycle_time) as avg_cycle_time,
        MAX(active_developers) as max_developers,
        AVG(avg_ai_ratio) as avg_ai_ratio
    FROM mart.velocity
    {where_clause}
    """
    result = query(sql)
    return result.iloc[0].to_dict()
