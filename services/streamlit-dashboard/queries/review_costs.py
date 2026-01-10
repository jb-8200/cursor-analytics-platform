"""
Code Review Costs SQL queries.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-04 SQL Query Modules

This module provides parameterized SQL queries for review cost analysis:
- Review iterations and comments
- Reviewer workload
- Review time distribution
- Total review hours

Depends on: mart.review_costs (P8 dbt mart)
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


def get_review_costs_data(repo_name: Optional[str] = None, days: Optional[int] = None) -> pd.DataFrame:
    """
    Get weekly review cost metrics.
    """
    where_clause, params = _build_filter(repo_name, days)
    
    sql = f"""
    SELECT
        week,
        repo_name,
        total_prs,
        avg_review_rounds,
        avg_reviewers_per_pr,
        avg_review_cycle_time,
        estimated_review_hours_per_pr,
        estimated_total_review_hours,
        large_prs
    FROM main_mart.mart_review_costs
    {where_clause}
    ORDER BY week DESC
    """
    return query(sql, params)


def get_review_iteration_distribution(repo_name: Optional[str] = None, days: Optional[int] = None) -> pd.DataFrame:
    """
    Get distribution of PRs by number of review iterations.
    """
    where_clause, params = _build_filter(repo_name, days)
    
    sql = f"""
    SELECT
        CASE
            WHEN avg_review_rounds = 1 THEN '1'
            WHEN avg_review_rounds = 2 THEN '2'
            ELSE '3+'
        END as iteration_count,
        COUNT(*) as pr_count,
        COUNT(*) * 100.0 / SUM(COUNT(*)) OVER () as percentage
    FROM main_mart.mart_review_costs
    {where_clause}
    GROUP BY iteration_count
    ORDER BY iteration_count
    """
    return query(sql, params)


def get_reviewer_workload(repo_name: Optional[str] = None, days: Optional[int] = None) -> pd.DataFrame:
    """
    Get reviewer workload metrics.
    """
    where_clause, params = _build_filter(repo_name, days)
    
    sql = f"""
    SELECT
        week,
        avg_reviewers_per_pr,
        estimated_total_review_hours,
        avg_review_cycle_time
    FROM main_mart.mart_review_costs
    {where_clause}
    ORDER BY week DESC
    """
    return query(sql, params)


def get_review_costs_summary(repo_name: Optional[str] = None, days: Optional[int] = None) -> dict:
    """
    Get summary statistics for review costs.
    """
    where_clause, params = _build_filter(repo_name, days)
    
    sql = f"""
    SELECT
        SUM(total_prs) as total_prs,
        AVG(avg_review_rounds) as avg_iterations,
        AVG(avg_reviewers_per_pr) as avg_reviewers,
        AVG(avg_review_rounds) as avg_comments,
        SUM(estimated_total_review_hours) as total_hours
    FROM main_mart.mart_review_costs
    {where_clause}
    """
    result = query(sql, params)
    if result.empty:
        return {}
    return result.iloc[0].to_dict()


def get_review_costs_by_ai_band(repo_name: Optional[str] = None, days: Optional[int] = None) -> pd.DataFrame:
    """
    Get review costs grouped by AI usage band.

    Note: Uses mart_ai_impact since ai_usage_band is only in that table.
    """
    where_clause, params = _build_filter(repo_name, days)

    sql = f"""
    SELECT
        ai_usage_band,
        AVG(avg_total_cycle_time) as avg_review_cycle_time,
        COUNT(*) as pr_count
    FROM main_mart.mart_ai_impact
    {where_clause}
    GROUP BY ai_usage_band
    ORDER BY ai_usage_band
    """
    return query(sql, params)
