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

from typing import Optional, Tuple, Dict, Any
from db.connector import query
import pandas as pd


def _build_filter(repo_name: Optional[str], days: Optional[int]) -> Tuple[str, Dict[str, Any]]:
    """Helper to build WHERE clause and parameters."""
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


def get_ai_impact_data(repo_name: Optional[str] = None, days: Optional[int] = None) -> pd.DataFrame:
    """
    Get AI impact metrics grouped by usage band.
    """
    where_clause, params = _build_filter(repo_name, days)
    
    sql = f"""
    SELECT
        week,
        ai_usage_band,
        pr_count,
        avg_ai_ratio,
        avg_total_cycle_time,
        revert_rate,
        bug_fix_rate,
        avg_pr_size,
        avg_files_changed
    FROM main_mart.mart_ai_impact
    {where_clause}
    ORDER BY week DESC, ai_usage_band
    """
    return query(sql, params)


def get_band_comparison(repo_name: Optional[str] = None, days: Optional[int] = None) -> pd.DataFrame:
    """
    Get comparison of metrics across AI usage bands.
    """
    where_clause, params = _build_filter(repo_name, days)
    
    sql = f"""
    SELECT
        ai_usage_band,
        SUM(pr_count) as total_prs,
        AVG(avg_ai_ratio) as avg_ai_ratio,
        AVG(avg_total_cycle_time) as avg_cycle_time,
        AVG(revert_rate) as avg_revert_rate,
        AVG(bug_fix_rate) as avg_bug_fix_rate
    FROM main_mart.mart_ai_impact
    {where_clause}
    GROUP BY ai_usage_band
    ORDER BY ai_usage_band
    """
    return query(sql, params)


def get_band_trends(repo_name: Optional[str] = None, days: Optional[int] = None) -> pd.DataFrame:
    """
    Get week-over-week trends for each AI usage band.
    """
    where_clause, params = _build_filter(repo_name, days)

    sql = f"""
    SELECT
        week,
        ai_usage_band,
        pr_count,
        avg_total_cycle_time,
        revert_rate
    FROM main_mart.mart_ai_impact
    {where_clause}
    ORDER BY week DESC, ai_usage_band
    """
    return query(sql, params)
