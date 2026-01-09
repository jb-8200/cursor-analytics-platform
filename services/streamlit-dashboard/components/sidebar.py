"""
Shared sidebar component for all dashboard pages.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-05 Shared Sidebar Component

This module provides a reusable sidebar with:
- Repository filter dropdown (populated from data)
- Date range filter
- Refresh button (dev mode only)
- Environment indicator
- Session state management for filters

Usage:
    from components.sidebar import render_sidebar

    render_sidebar()
    # Access filters via session state
    where_clause = get_filter_where_clause()
"""

import streamlit as st
import os
from db.connector import query, refresh_data

# Determine database mode
DB_MODE = os.getenv("DB_MODE", "duckdb")


def render_sidebar():
    """
    Render the shared sidebar with filters and controls.

    This function populates the sidebar with:
    - Title
    - Repository filter (from database)
    - Date range filter
    - Refresh button (dev mode)
    - Environment indicator

    Filters are stored in st.session_state for access by pages.
    """
    with st.sidebar:
        # Title
        st.title("📊 DOXAPI Analytics")
        st.divider()

        # Filters section
        st.subheader("🔧 Filters")

        # Repository filter
        try:
            repos_df = query("SELECT DISTINCT repo_name FROM mart.velocity ORDER BY repo_name")
            repo_options = ["All"] + repos_df["repo_name"].tolist()
        except Exception as e:
            st.error(f"Failed to load repositories: {e}")
            repo_options = ["All"]

        selected_repo = st.selectbox(
            "Repository",
            repo_options,
            key="repo_selectbox"
        )

        # Date range filter
        date_options = ["Last 7 days", "Last 30 days", "Last 90 days", "All time"]
        selected_range = st.selectbox(
            "Date Range",
            date_options,
            index=2,  # Default to Last 90 days
            key="date_selectbox"
        )

        # Store in session state
        st.session_state["filter_repo"] = selected_repo
        st.session_state["filter_date_range"] = selected_range

        st.divider()

        # Refresh section
        st.subheader("🔄 Data")

        if DB_MODE == "duckdb":
            # Dev mode: Show refresh button
            if st.button("Refresh Data", use_container_width=True):
                with st.spinner("Refreshing data..."):
                    success = refresh_data()
                    if success:
                        st.success("Data refreshed!")
                        # Clear caches
                        st.cache_data.clear()
                        st.cache_resource.clear()
                        st.rerun()
                    else:
                        st.error("Refresh failed. Check logs.")
        else:
            # Production mode: Info message
            st.info("📅 Data updates every 15 min")

        st.divider()

        # Environment indicator
        st.subheader("🌐 Environment")
        if DB_MODE == "snowflake":
            st.success("🟢 Production")
        else:
            st.warning("🟡 Development")


def get_filter_where_clause() -> str:
    """
    Build SQL WHERE clause from current filter state.

    Returns:
        SQL WHERE clause string based on session state filters.
        Returns empty string if no filters applied (All repos, All time).

    Example:
        >>> where = get_filter_where_clause()
        >>> df = get_velocity_data(where)

        # With repo filter:
        "WHERE repo_name = 'acme/platform' AND week >= CURRENT_DATE - INTERVAL '90 days'"

        # All repos, last 30 days:
        "WHERE week >= CURRENT_DATE - INTERVAL '30 days'"

        # All repos, all time:
        ""
    """
    conditions = []

    # Repository filter
    repo = st.session_state.get("filter_repo", "All")
    if repo != "All":
        conditions.append(f"repo_name = '{repo}'")

    # Date range filter
    date_range = st.session_state.get("filter_date_range", "All time")
    if date_range != "All time":
        # Extract number of days from date range string
        days_map = {
            "Last 7 days": 7,
            "Last 30 days": 30,
            "Last 90 days": 90,
        }
        days = days_map.get(date_range)
        if days:
            conditions.append(f"week >= CURRENT_DATE - INTERVAL '{days} days'")

    # Build WHERE clause
    if conditions:
        return "WHERE " + " AND ".join(conditions)
    return ""
