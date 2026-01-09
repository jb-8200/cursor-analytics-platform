"""
Streamlit-specific refresh functionality with UI feedback.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-10 Refresh Pipeline

This module provides Streamlit-aware refresh functions with progress indicators.
"""

import os
from typing import Optional


def refresh_data_with_ui() -> bool:
    """
    Trigger ETL pipeline refresh with Streamlit UI feedback.

    This function is designed to be called from Streamlit pages and provides
    visual feedback via spinners and status messages.

    Returns:
        bool: True if refresh succeeded, False otherwise

    Note:
        Only works in DuckDB mode. Shows warning in Snowflake/production mode.
    """
    try:
        import streamlit as st
    except ImportError:
        # Fallback to non-UI version if streamlit not available
        from db.connector import refresh_data
        return refresh_data()

    from db.connector import refresh_data, DB_MODE

    if DB_MODE == "snowflake":
        st.warning("⚠️ Data refresh not available in production mode.")
        st.info("Data is updated automatically every 15 minutes via scheduled jobs.")
        return False

    try:
        # Step 1: Extract
        with st.spinner("📥 Extracting data from cursor-sim..."):
            cursor_sim_url = os.getenv("CURSOR_SIM_URL", "http://localhost:8080")
            st.caption(f"Fetching from {cursor_sim_url}")

            # Trigger the refresh (which includes all steps)
            success = refresh_data()

            if not success:
                st.error("❌ Refresh failed. Check logs for details.")
                return False

        # Clear caches after successful refresh
        st.cache_data.clear()
        st.cache_resource.clear()

        st.success("✅ Data refreshed successfully!")
        st.info("♻️ Caches cleared. Data will reload on next query.")

        return True

    except Exception as e:
        st.error(f"❌ Refresh failed: {str(e)}")
        return False


def is_refresh_available() -> bool:
    """
    Check if data refresh is available in current environment.

    Returns:
        bool: True if refresh is available (DuckDB mode), False otherwise
    """
    from db.connector import DB_MODE
    return DB_MODE == "duckdb"


def get_last_refresh_time() -> Optional[str]:
    """
    Get timestamp of last data refresh.

    Returns:
        Optional[str]: ISO timestamp of last refresh, or None if not available

    Note:
        This is a placeholder for future implementation.
        Would require storing refresh timestamps in a metadata table.
    """
    # TODO: Implement metadata tracking
    return None
