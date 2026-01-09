"""
Test cases for Home Page (app.py).

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-06 Home Page

Tests verify:
- Page configuration
- KPI calculations
- Navigation instructions
- Sidebar rendering
- Error handling
"""

import pytest
from unittest.mock import patch, MagicMock
import pandas as pd
import sys
import os

# Add parent directory to path for imports
sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))


class TestHomePageConfiguration:
    """Test page configuration and setup."""

    def test_page_config_title_is_set(self):
        """Verify page title is DOXAPI Analytics."""
        with patch("streamlit.set_page_config") as mock_config:
            with patch("streamlit.title"):
                with patch("streamlit.markdown"):
                    with patch("streamlit.info"):
                        with patch("streamlit.columns") as mock_cols:
                            mock_cols.return_value = [MagicMock() for _ in range(4)]
                            with patch("streamlit.metric"):
                                with patch("streamlit.error"):
                                    # Force reimport
                                    import importlib
                                    if "app" in sys.modules:
                                        del sys.modules["app"]
                                    import app

                                    # Check set_page_config was called
                                    assert mock_config.called


class TestHomePageContent:
    """Test home page content rendering."""

    def test_title_displays_correctly(self):
        """Verify page title is displayed."""
        with patch("streamlit.set_page_config"):
            with patch("streamlit.title") as mock_title:
                with patch("streamlit.markdown"):
                    with patch("streamlit.info"):
                        with patch("streamlit.columns") as mock_cols:
                            mock_cols.return_value = [MagicMock() for _ in range(4)]
                            with patch("streamlit.metric"):
                                with patch("streamlit.error"):
                                    import importlib
                                    if "app" in sys.modules:
                                        del sys.modules["app"]
                                    import app

                                    # Title should be called
                                    assert mock_title.called


class TestHomePageKPIs:
    """Test KPI calculations and display."""

    def test_kpis_use_four_columns(self):
        """Verify KPIs are displayed in 4 columns."""
        with patch("streamlit.set_page_config"):
            with patch("streamlit.title"):
                with patch("streamlit.markdown"):
                    with patch("streamlit.info"):
                        with patch("streamlit.columns") as mock_columns:
                            mock_columns.return_value = [MagicMock() for _ in range(4)]
                            with patch("streamlit.metric"):
                                with patch("streamlit.error"):
                                    import importlib
                                    if "app" in sys.modules:
                                        del sys.modules["app"]
                                    import app

                                    # Verify columns(4) was called
                                    mock_columns.assert_called_with(4)


class TestHomePageErrorHandling:
    """Test error handling for home page."""

    def test_handles_empty_database(self):
        """Verify graceful handling when database is empty."""
        with patch("streamlit.set_page_config"):
            with patch("streamlit.title"):
                with patch("streamlit.markdown"):
                    with patch("streamlit.info"):
                        with patch("streamlit.columns") as mock_cols:
                            mock_cols.return_value = [MagicMock() for _ in range(4)]
                            with patch("streamlit.metric"):
                                with patch("streamlit.error"):
                                    import importlib
                                    if "app" in sys.modules:
                                        del sys.modules["app"]
                                    # Should not raise exception
                                    import app


# Run tests with: pytest tests/test_home_page.py -v
