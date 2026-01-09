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

from db.connector import query
import pandas as pd


def get_review_costs_data(where_clause: str = "") -> pd.DataFrame:
    """
    Get weekly review cost metrics.

    Args:
        where_clause: Optional SQL WHERE clause for filtering (e.g., "WHERE repo_name = 'acme/platform'")

    Returns:
        DataFrame with review cost metrics including:
        - week: Week start date
        - repo_name: Repository name
        - total_prs: Total PRs reviewed
        - avg_review_iterations: Average number of review cycles per PR
        - avg_reviewers_per_pr: Average number of reviewers per PR
        - avg_review_comments: Average number of review comments per PR
        - avg_review_time: Average time spent in review (days)
        - total_review_hours: Total review hours for the week

    Example:
        >>> df = get_review_costs_data()
        >>> df = get_review_costs_data("WHERE repo_name = 'acme/platform'")
    """
    sql = f"""
    SELECT
        week,
        repo_name,
        total_prs,
        avg_review_iterations,
        avg_reviewers_per_pr,
        avg_review_comments,
        avg_review_time,
        total_review_hours
    FROM mart.review_costs
    {where_clause}
    ORDER BY week DESC
    """
    return query(sql)


def get_review_iteration_distribution(where_clause: str = "") -> pd.DataFrame:
    """
    Get distribution of PRs by number of review iterations.

    Args:
        where_clause: Optional SQL WHERE clause for filtering

    Returns:
        DataFrame with iteration distribution:
        - iteration_count: Number of iterations (1, 2, 3+)
        - pr_count: Number of PRs with this iteration count
        - percentage: Percentage of total PRs

    Example:
        >>> dist = get_review_iteration_distribution()
        >>> one_iteration = dist[dist['iteration_count'] == 1]
    """
    sql = f"""
    SELECT
        CASE
            WHEN avg_review_iterations = 1 THEN '1'
            WHEN avg_review_iterations = 2 THEN '2'
            ELSE '3+'
        END as iteration_count,
        COUNT(*) as pr_count,
        COUNT(*) * 100.0 / SUM(COUNT(*)) OVER () as percentage
    FROM mart.review_costs
    {where_clause}
    GROUP BY iteration_count
    ORDER BY iteration_count
    """
    return query(sql)


def get_reviewer_workload(where_clause: str = "") -> pd.DataFrame:
    """
    Get reviewer workload metrics.

    Args:
        where_clause: Optional SQL WHERE clause for filtering

    Returns:
        DataFrame with workload metrics:
        - week: Week start date
        - avg_reviewers_per_pr: Average reviewers per PR
        - total_review_hours: Total review hours
        - avg_review_time: Average time per review

    Example:
        >>> workload = get_reviewer_workload()
        >>> recent = workload.head(12)  # Last 12 weeks
    """
    sql = f"""
    SELECT
        week,
        avg_reviewers_per_pr,
        total_review_hours,
        avg_review_time
    FROM mart.review_costs
    {where_clause}
    ORDER BY week DESC
    """
    return query(sql)


def get_review_costs_summary(where_clause: str = "") -> dict:
    """
    Get summary statistics for review costs.

    Args:
        where_clause: Optional SQL WHERE clause for filtering

    Returns:
        Dictionary with summary stats:
        - total_prs: Total number of PRs reviewed
        - avg_iterations: Average review iterations
        - avg_reviewers: Average reviewers per PR
        - avg_comments: Average comments per PR
        - total_hours: Total review hours

    Example:
        >>> summary = get_review_costs_summary()
        >>> print(f"Total review hours: {summary['total_hours']:.0f}")
    """
    sql = f"""
    SELECT
        SUM(total_prs) as total_prs,
        AVG(avg_review_iterations) as avg_iterations,
        AVG(avg_reviewers_per_pr) as avg_reviewers,
        AVG(avg_review_comments) as avg_comments,
        SUM(total_review_hours) as total_hours
    FROM mart.review_costs
    {where_clause}
    """
    result = query(sql)
    return result.iloc[0].to_dict()


def get_review_costs_by_ai_band(where_clause: str = "") -> pd.DataFrame:
    """
    Get review costs grouped by AI usage band.

    Args:
        where_clause: Optional SQL WHERE clause for filtering

    Returns:
        DataFrame with costs by AI band:
        - ai_usage_band: AI usage band (low, medium, high)
        - avg_review_iterations: Average iterations for this band
        - avg_review_comments: Average comments for this band
        - avg_review_time: Average review time for this band

    Example:
        >>> by_band = get_review_costs_by_ai_band()
        >>> high_ai = by_band[by_band['ai_usage_band'] == 'high']
    """
    sql = f"""
    SELECT
        ai_usage_band,
        AVG(avg_review_iterations) as avg_review_iterations,
        AVG(avg_review_comments) as avg_review_comments,
        AVG(avg_review_time) as avg_review_time
    FROM mart.review_costs
    {where_clause}
    GROUP BY ai_usage_band
    ORDER BY ai_usage_band
    """
    return query(sql)
