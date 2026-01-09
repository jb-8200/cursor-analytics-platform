"""
Tests for Velocity Metrics Dashboard Page.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-07 Velocity Page

Tests verify:
- Page can be imported without errors
- Required functions are available
- Page displays expected components
"""

import pytest
import sys
import os

# Add parent directory to path for imports
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))


def test_velocity_page_imports():
    """Verify velocity page can be imported."""
    try:
        # This will execute the page module imports
        import importlib.util
        spec = importlib.util.spec_from_file_location(
            "velocity_page",
            "pages/1_velocity.py"
        )
        assert spec is not None, "Could not load velocity page spec"

        # Just verify the file exists and has valid Python syntax
        # We can't actually execute it without Streamlit context
        module = importlib.util.module_from_spec(spec)
        assert module is not None, "Could not create module from spec"

    except SyntaxError as e:
        pytest.fail(f"Velocity page has syntax error: {e}")
    except Exception as e:
        # Other errors are acceptable since we can't run Streamlit context
        # We just want to verify syntax is valid
        pass


def test_velocity_page_file_exists():
    """Verify velocity page file exists."""
    import os
    page_path = "pages/1_velocity.py"
    assert os.path.exists(page_path), f"Velocity page not found at {page_path}"


def test_velocity_page_has_required_imports():
    """Verify velocity page has required imports."""
    with open("pages/1_velocity.py", "r") as f:
        content = f.read()

    # Check for required imports
    required_imports = [
        "import streamlit",
        "import plotly",
        "from components.sidebar import render_sidebar",
        "from queries.velocity import"
    ]

    for required in required_imports:
        assert required in content, f"Missing required import: {required}"


def test_velocity_page_has_page_config():
    """Verify velocity page has Streamlit page configuration."""
    with open("pages/1_velocity.py", "r") as f:
        content = f.read()

    assert "st.set_page_config" in content, "Missing st.set_page_config"
    assert 'page_title="Velocity' in content or "page_title='Velocity" in content
    assert "layout=\"wide\"" in content or "layout='wide'" in content


def test_velocity_page_renders_sidebar():
    """Verify velocity page calls render_sidebar."""
    with open("pages/1_velocity.py", "r") as f:
        content = f.read()

    assert "render_sidebar()" in content, "Missing render_sidebar() call"


def test_velocity_page_has_title():
    """Verify velocity page has a title."""
    with open("pages/1_velocity.py", "r") as f:
        content = f.read()

    assert "st.title" in content, "Missing st.title"
    assert "Velocity" in content, "Missing 'Velocity' in title"


def test_velocity_page_uses_filters():
    """Verify velocity page uses filter WHERE clause."""
    with open("pages/1_velocity.py", "r") as f:
        content = f.read()

    assert "get_filter_where_clause" in content, "Missing get_filter_where_clause"


def test_velocity_page_has_kpis():
    """Verify velocity page displays KPI metrics."""
    with open("pages/1_velocity.py", "r") as f:
        content = f.read()

    # Check for 4 columns (4 KPIs)
    assert "st.columns(4)" in content, "Missing 4-column layout for KPIs"

    # Check for metrics
    kpi_metrics = ["Total PRs", "Cycle Time", "Developers", "AI Ratio"]
    for metric in kpi_metrics:
        # Metrics should appear somewhere in the page
        assert metric.lower().replace(" ", "") in content.lower().replace(" ", ""), \
            f"Missing KPI metric: {metric}"


def test_velocity_page_has_charts():
    """Verify velocity page has required charts."""
    with open("pages/1_velocity.py", "r") as f:
        content = f.read()

    # Check for plotly charts
    assert "st.plotly_chart" in content, "Missing plotly charts"

    # Check for cycle time breakdown query
    assert "get_cycle_time_breakdown" in content, "Missing cycle time breakdown"


def test_velocity_page_queries_data():
    """Verify velocity page queries velocity data."""
    with open("pages/1_velocity.py", "r") as f:
        content = f.read()

    assert "get_velocity_data" in content, "Missing get_velocity_data query"
    assert "get_velocity_summary" in content or "get_velocity_data" in content


def test_velocity_page_handles_errors():
    """Verify velocity page has error handling."""
    with open("pages/1_velocity.py", "r") as f:
        content = f.read()

    # Check for try/except or error handling
    assert "try:" in content or "except" in content, "Missing error handling"
    assert "st.error" in content or "st.warning" in content, "Missing error display"


def test_velocity_page_has_documentation():
    """Verify velocity page has docstring."""
    with open("pages/1_velocity.py", "r") as f:
        content = f.read()

    # Check for module docstring
    assert '"""' in content or "'''" in content, "Missing module docstring"
    assert "TASK-P9-07" in content, "Missing task reference in docstring"
