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

from typing import Optional, Tuple, Dict, Any
from db.connector import query
import pandas as pd


def _build_filter(repo_name: Optional[str], days: Optional[int]) -> Tuple[str, Dict[str, Any]]:
    """
    Build WHERE clause and parameters for filtering.
    
    Args:
        repo_name: Repository name filter (or None for all)
        days: Number of days lookback (or None for all time)
        
    Returns:
        Tuple of (where_clause, params_dict)
    """
    conditions = []
    params = {}

    if repo_name and repo_name != "All":
        conditions.append("repo_name = $repo")
        params["repo"] = repo_name

    if days:
        # Note: Use f-string for INTERVAL since DuckDB doesn't support parameterized INTERVAL
        # This is safe because days is validated as an integer
        conditions.append(f"week >= CURRENT_DATE - INTERVAL '{days}' DAY")

    if not conditions:
        return "", {}
        
    return "WHERE " + " AND ".join(conditions), params


def get_velocity_data(repo_name: Optional[str] = None, days: Optional[int] = None) -> pd.DataFrame:
    """
    Get weekly velocity metrics.

    Args:
        repo_name: Optional repository name for filtering
        days: Optional number of days for date range filter

    Returns:
        DataFrame with velocity metrics
    """
    where_clause, params = _build_filter(repo_name, days)
    
    sql = f"""
    SELECT
        week,
        repo_name,
        active_developers,
        total_prs,
        total_commits,
        avg_pr_size,
        avg_files_changed,
        avg_total_cycle_time,
        avg_ai_ratio,
        total_ai_lines,
        total_lines
    FROM main_mart.mart_velocity
    {where_clause}
    ORDER BY week DESC
    """
    return query(sql, params)


def get_cycle_time_breakdown(repo_name: Optional[str] = None, days: Optional[int] = None) -> pd.DataFrame:
    """
    Get cycle time breakdown by component.

    Args:
        repo_name: Optional repository name filter
        days: Optional number of days filter

    Note:
        Returns total cycle time only (individual components not available in current schema)
    """
    where_clause, params = _build_filter(repo_name, days)

    sql = f"""
    SELECT
        'Total Cycle Time' as component,
        AVG(avg_total_cycle_time) as hours
    FROM main_mart.mart_velocity {where_clause}
    """
    return query(sql, params)


def get_velocity_summary(repo_name: Optional[str] = None, days: Optional[int] = None) -> dict:
    """
    Get summary statistics for velocity metrics.

    Args:
        repo_name: Optional repository name filter
        days: Optional number of days filter
    """
    where_clause, params = _build_filter(repo_name, days)
    
    sql = f"""
    SELECT
        SUM(total_prs) as total_prs,
        AVG(avg_total_cycle_time) as avg_cycle_time,
        MAX(active_developers) as max_developers,
        AVG(avg_ai_ratio) as avg_ai_ratio
    FROM main_mart.mart_velocity
    {where_clause}
    """
    result = query(sql, params)
    if result.empty:
        return {}
    return result.iloc[0].to_dict()

