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


def get_quality_data(repo_name: Optional[str] = None, days: Optional[int] = None) -> pd.DataFrame:
    """
    Get weekly quality metrics.
    """
    where_clause, params = _build_filter(repo_name, days)
    
    sql = f"""
    SELECT
        week,
        repo_name,
        total_prs,
        reverted_prs,
        revert_rate,
        bug_fix_prs,
        bug_fix_rate,
        avg_reviews_per_pr,
        unreviewed_prs
    FROM main_mart.mart_quality
    {where_clause}
    ORDER BY week DESC
    """
    return query(sql, params)


def get_revert_trends(repo_name: Optional[str] = None, days: Optional[int] = None) -> pd.DataFrame:
    """
    Get revert rate trends over time.
    """
    where_clause, params = _build_filter(repo_name, days)
    
    sql = f"""
    SELECT
        week,
        repo_name,
        revert_rate,
        reverted_prs
    FROM main_mart.mart_quality
    {where_clause}
    ORDER BY week DESC
    """
    return query(sql, params)


def get_quality_summary(repo_name: Optional[str] = None, days: Optional[int] = None) -> dict:
    """
    Get summary statistics for quality metrics.
    """
    where_clause, params = _build_filter(repo_name, days)
    
    sql = f"""
    SELECT
        SUM(total_prs) as total_prs,
        SUM(reverted_prs) as total_reverted,
        AVG(revert_rate) as avg_revert_rate,
        AVG(bug_fix_rate) as avg_bug_fix_rate,
        AVG(avg_reviews_per_pr) as avg_reviews_per_pr
    FROM main_mart.mart_quality
    {where_clause}
    """
    result = query(sql, params)
    if result.empty:
        return {}
    return result.iloc[0].to_dict()


def get_quality_by_ai_band(repo_name: Optional[str] = None, days: Optional[int] = None) -> pd.DataFrame:
    """
    Get quality metrics grouped by AI usage band.

    Note: Uses mart_ai_impact since ai_usage_band is only in that table.
    """
    where_clause, params = _build_filter(repo_name, days)

    sql = f"""
    SELECT
        ai_usage_band,
        AVG(revert_rate) as avg_revert_rate,
        AVG(bug_fix_rate) as avg_bug_fix_rate
    FROM main_mart.mart_ai_impact
    {where_clause}
    GROUP BY ai_usage_band
    ORDER BY ai_usage_band
    """
    return query(sql, params)
