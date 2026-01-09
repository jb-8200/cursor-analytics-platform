"""
Tests for sidebar component.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-05 Shared Sidebar Component

This test suite verifies the sidebar component behavior:
- Repository filter dropdown populated from data
- Date range filter options
- Refresh button visibility (dev mode only)
- Environment indicator (dev vs production)
- Session state updates
"""

import pytest
from unittest.mock import Mock, patch, MagicMock
import pandas as pd
import os


@pytest.fixture
def mock_repos_df():
    """Mock repository list from database."""
    return pd.DataFrame({
        "repo_name": ["acme/platform", "acme/api", "acme/web"]
    })


@pytest.fixture
def mock_streamlit():
    """Mock Streamlit module for testing."""
    with patch("components.sidebar.st") as mock_st:
        # Mock session_state as a dictionary
        mock_st.session_state = {}

        # Mock sidebar context manager
        sidebar_context = MagicMock()
        mock_st.sidebar.__enter__ = Mock(return_value=sidebar_context)
        mock_st.sidebar.__exit__ = Mock(return_value=False)

        yield mock_st


@pytest.fixture
def mock_query(mock_repos_df):
    """Mock query function."""
    with patch("components.sidebar.query") as mock_q:
        mock_q.return_value = mock_repos_df
        yield mock_q


def test_render_sidebar_creates_title(mock_streamlit, mock_query):
    """Verify sidebar renders title."""
    from components.sidebar import render_sidebar

    render_sidebar()

    # Should call st.sidebar context manager
    assert mock_streamlit.sidebar.__enter__.called


def test_render_sidebar_populates_repo_filter(mock_streamlit, mock_query, mock_repos_df):
    """Verify repository filter is populated from database."""
    from components.sidebar import render_sidebar

    render_sidebar()

    # Should query for distinct repos
    mock_query.assert_called_once()
    assert "DISTINCT repo_name" in mock_query.call_args[0][0]
    assert "mart.velocity" in mock_query.call_args[0][0]


def test_render_sidebar_stores_filters_in_session(mock_streamlit, mock_query):
    """Verify filters are stored in session state."""
    from components.sidebar import render_sidebar

    # Mock selectbox to return specific values
    mock_streamlit.selectbox.side_effect = ["acme/platform", "Last 30 days"]

    render_sidebar()

    # Should store in session_state
    assert "filter_repo" in mock_streamlit.session_state
    assert "filter_date_range" in mock_streamlit.session_state


def test_render_sidebar_shows_refresh_in_dev_mode(mock_streamlit, mock_query, monkeypatch):
    """Verify refresh button is visible in DuckDB mode."""
    monkeypatch.setenv("DB_MODE", "duckdb")

    from components.sidebar import render_sidebar

    render_sidebar()

    # Should call st.button for refresh
    button_calls = [call for call in mock_streamlit.method_calls if call[0] == "button"]
    assert len(button_calls) > 0

    # At least one button should be "Refresh Data"
    button_texts = [call[1][0] if call[1] else "" for call in button_calls]
    assert any("Refresh" in text for text in button_texts)


def test_render_sidebar_hides_refresh_in_prod_mode(mock_streamlit, mock_query, monkeypatch):
    """Verify refresh button is hidden in Snowflake mode."""
    monkeypatch.setenv("DB_MODE", "snowflake")

    from components.sidebar import render_sidebar

    render_sidebar()

    # Should show info message instead
    info_calls = [call for call in mock_streamlit.method_calls if call[0] == "info"]
    assert len(info_calls) > 0


def test_render_sidebar_shows_dev_indicator(mock_streamlit, mock_query, monkeypatch):
    """Verify environment indicator shows Development for DuckDB."""
    monkeypatch.setenv("DB_MODE", "duckdb")

    from components.sidebar import render_sidebar

    render_sidebar()

    # Should call st.warning for dev mode
    warning_calls = [call for call in mock_streamlit.method_calls if call[0] == "warning"]
    assert len(warning_calls) > 0


def test_render_sidebar_shows_prod_indicator(mock_streamlit, mock_query, monkeypatch):
    """Verify environment indicator shows Production for Snowflake."""
    monkeypatch.setenv("DB_MODE", "snowflake")

    from components.sidebar import render_sidebar

    render_sidebar()

    # Should call st.success for prod mode
    success_calls = [call for call in mock_streamlit.method_calls if call[0] == "success"]
    assert len(success_calls) > 0


def test_get_filter_where_clause_all_repos(monkeypatch):
    """Verify get_filter_where_clause returns empty string for 'All' repos."""
    from components.sidebar import get_filter_where_clause

    # Mock session_state
    with patch("components.sidebar.st") as mock_st:
        mock_st.session_state = {
            "filter_repo": "All",
            "filter_date_range": "Last 90 days"
        }

        result = get_filter_where_clause()

        # Should include date filter but not repo filter
        assert "WHERE" in result
        assert "repo_name" not in result
        assert "week >=" in result


def test_get_filter_where_clause_specific_repo(monkeypatch):
    """Verify get_filter_where_clause includes repo filter for specific repo."""
    from components.sidebar import get_filter_where_clause

    # Mock session_state
    with patch("components.sidebar.st") as mock_st:
        mock_st.session_state = {
            "filter_repo": "acme/platform",
            "filter_date_range": "All time"
        }

        result = get_filter_where_clause()

        # Should include repo filter
        assert "WHERE" in result
        assert "repo_name = 'acme/platform'" in result


def test_get_filter_where_clause_date_ranges():
    """Verify get_filter_where_clause handles all date range options."""
    from components.sidebar import get_filter_where_clause

    test_cases = [
        ("Last 7 days", 7),
        ("Last 30 days", 30),
        ("Last 90 days", 90),
        ("All time", None),
    ]

    for range_option, expected_days in test_cases:
        with patch("components.sidebar.st") as mock_st:
            mock_st.session_state = {
                "filter_repo": "All",
                "filter_date_range": range_option
            }

            result = get_filter_where_clause()

            if expected_days:
                assert "week >=" in result
                assert f"INTERVAL '{expected_days} days'" in result or f"interval '{expected_days} days'" in result
            else:
                # All time should not have date filter
                assert "week >=" not in result or result == ""
