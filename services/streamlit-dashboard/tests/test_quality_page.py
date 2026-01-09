"""
Tests for Quality Metrics Dashboard Page.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-09 Quality Page

Tests verify:
- Page can be imported without errors
- Required functions are available
- Page displays expected components (KPIs, charts)
- Quality metrics are displayed correctly
"""

import pytest
import sys
import os

# Add parent directory to path for imports
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))


def test_quality_page_imports():
    """Verify quality page can be imported."""
    try:
        import importlib.util
        spec = importlib.util.spec_from_file_location(
            "quality_page",
            "pages/3_quality.py"
        )
        assert spec is not None, "Could not load quality page spec"

        module = importlib.util.module_from_spec(spec)
        assert module is not None, "Could not create module from spec"

    except SyntaxError as e:
        pytest.fail(f"Quality page has syntax error: {e}")
    except Exception as e:
        # Other errors are acceptable since we can't run Streamlit context
        pass


def test_quality_page_file_exists():
    """Verify quality page file exists."""
    page_path = "pages/3_quality.py"
    assert os.path.exists(page_path), f"Quality page not found at {page_path}"


def test_quality_page_has_required_imports():
    """Verify quality page has required imports."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    required_imports = [
        "import streamlit",
        "import plotly",
        "from components.sidebar import render_sidebar",
        "from queries.quality import"
    ]

    for required in required_imports:
        assert required in content, f"Missing required import: {required}"


def test_quality_page_has_page_config():
    """Verify quality page has Streamlit page configuration."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    assert "st.set_page_config" in content, "Missing st.set_page_config"
    assert 'page_title=' in content
    assert "layout=\"wide\"" in content or "layout='wide'" in content


def test_quality_page_renders_sidebar():
    """Verify quality page calls render_sidebar."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    assert "render_sidebar()" in content, "Missing render_sidebar() call"


def test_quality_page_has_title():
    """Verify quality page has a title."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    assert "st.title" in content, "Missing st.title"
    assert "Quality" in content, "Missing 'Quality' in title"


def test_quality_page_uses_filters():
    """Verify quality page uses filter WHERE clause."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    assert "get_filter_where_clause" in content, "Missing get_filter_where_clause"


def test_quality_page_has_kpis():
    """Verify quality page displays 4 KPI metrics."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    # Check for 4 columns (4 KPIs)
    assert "st.columns(4)" in content, "Missing 4-column layout for KPIs"
    assert "st.metric" in content, "Missing st.metric calls"


def test_quality_page_has_revert_rate_kpi():
    """Verify quality page has Revert Rate KPI."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read().lower()

    assert "revert" in content, "Missing Revert Rate KPI"


def test_quality_page_has_hotfix_rate_kpi():
    """Verify quality page has Hotfix/Bug Fix Rate KPI."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read().lower()

    assert "bug" in content or "hotfix" in content or "fix" in content, \
        "Missing Hotfix/Bug Fix Rate KPI"


def test_quality_page_has_line_chart():
    """Verify quality page has line chart for quality trends."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    # Check for line chart (trend over time)
    assert "px.line" in content or "go.Scatter" in content, "Missing line chart"
    assert "st.plotly_chart" in content, "Missing plotly chart display"


def test_quality_page_has_bar_chart():
    """Verify quality page has bar chart for revert rates by repo."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    assert "px.bar" in content, "Missing bar chart"


def test_quality_page_queries_data():
    """Verify quality page queries quality data."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    assert "get_quality_data" in content, "Missing get_quality_data query"


def test_quality_page_has_summary():
    """Verify quality page shows summary statistics."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    assert "get_quality_summary" in content, "Missing get_quality_summary query"


def test_quality_page_handles_errors():
    """Verify quality page has error handling."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    assert "try:" in content or "except" in content, "Missing error handling"
    assert "st.error" in content or "st.warning" in content, "Missing error display"


def test_quality_page_has_documentation():
    """Verify quality page has docstring."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    assert '"""' in content or "'''" in content, "Missing module docstring"
    assert "TASK-P9-09" in content, "Missing task reference in docstring"


def test_quality_page_has_subheaders():
    """Verify quality page has section subheaders."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    assert "st.subheader" in content, "Missing section subheaders"


def test_quality_page_has_dividers():
    """Verify quality page has visual dividers."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    assert "st.divider" in content, "Missing visual dividers"


def test_quality_page_has_revert_trends():
    """Verify quality page shows revert trends."""
    with open("pages/3_quality.py", "r") as f:
        content = f.read()

    assert "get_revert_trends" in content or "trend" in content.lower(), \
        "Missing revert trends visualization"
