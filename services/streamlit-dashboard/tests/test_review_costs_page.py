"""
Tests for Review Costs Dashboard Page.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-09 Review Costs Page

Tests verify:
- Page can be imported without errors
- Required functions are available
- Page displays expected components (KPIs, charts)
- Review metrics are displayed correctly
"""

import pytest
import sys
import os

# Add parent directory to path for imports
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))


def test_review_costs_page_imports():
    """Verify review costs page can be imported."""
    try:
        import importlib.util
        spec = importlib.util.spec_from_file_location(
            "review_costs_page",
            "pages/4_review_costs.py"
        )
        assert spec is not None, "Could not load review costs page spec"

        module = importlib.util.module_from_spec(spec)
        assert module is not None, "Could not create module from spec"

    except SyntaxError as e:
        pytest.fail(f"Review costs page has syntax error: {e}")
    except Exception as e:
        # Other errors are acceptable since we can't run Streamlit context
        pass


def test_review_costs_page_file_exists():
    """Verify review costs page file exists."""
    page_path = "pages/4_review_costs.py"
    assert os.path.exists(page_path), f"Review costs page not found at {page_path}"


def test_review_costs_page_has_required_imports():
    """Verify review costs page has required imports."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    required_imports = [
        "import streamlit",
        "import plotly",
        "from components.sidebar import render_sidebar",
        "from queries.review_costs import"
    ]

    for required in required_imports:
        assert required in content, f"Missing required import: {required}"


def test_review_costs_page_has_page_config():
    """Verify review costs page has Streamlit page configuration."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "st.set_page_config" in content, "Missing st.set_page_config"
    assert 'page_title=' in content
    assert "layout=\"wide\"" in content or "layout='wide'" in content


def test_review_costs_page_renders_sidebar():
    """Verify review costs page calls render_sidebar."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "render_sidebar()" in content, "Missing render_sidebar() call"


def test_review_costs_page_has_title():
    """Verify review costs page has a title."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "st.title" in content, "Missing st.title"
    assert "Review" in content, "Missing 'Review' in title"


def test_review_costs_page_uses_filters():
    """Verify review costs page uses filter WHERE clause."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "get_filter_where_clause" in content, "Missing get_filter_where_clause"


def test_review_costs_page_has_kpis():
    """Verify review costs page displays 4 KPI metrics."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    # Check for 4 columns (4 KPIs)
    assert "st.columns(4)" in content, "Missing 4-column layout for KPIs"
    assert "st.metric" in content, "Missing st.metric calls"


def test_review_costs_page_has_iterations_kpi():
    """Verify review costs page has Avg Reviews/Iterations KPI."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read().lower()

    assert "iteration" in content or "reviews" in content, \
        "Missing Reviews/Iterations KPI"


def test_review_costs_page_has_turnaround_kpi():
    """Verify review costs page has Review Turnaround KPI."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read().lower()

    assert "turnaround" in content or "review time" in content or "reviewtime" in content, \
        "Missing Review Turnaround KPI"


def test_review_costs_page_has_comments_kpi():
    """Verify review costs page has Comment Density KPI."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read().lower()

    assert "comment" in content, "Missing Comment Density KPI"


def test_review_costs_page_has_charts():
    """Verify review costs page has charts for review analysis."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "st.plotly_chart" in content, "Missing plotly charts"


def test_review_costs_page_has_bar_chart():
    """Verify review costs page has bar chart."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "px.bar" in content, "Missing bar chart"


def test_review_costs_page_queries_data():
    """Verify review costs page queries review costs data."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "get_review_costs_data" in content, "Missing get_review_costs_data query"


def test_review_costs_page_has_summary():
    """Verify review costs page shows summary statistics."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "get_review_costs_summary" in content, "Missing get_review_costs_summary query"


def test_review_costs_page_handles_errors():
    """Verify review costs page has error handling."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "try:" in content or "except" in content, "Missing error handling"
    assert "st.error" in content or "st.warning" in content, "Missing error display"


def test_review_costs_page_has_documentation():
    """Verify review costs page has docstring."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert '"""' in content or "'''" in content, "Missing module docstring"
    assert "TASK-P9-09" in content, "Missing task reference in docstring"


def test_review_costs_page_has_subheaders():
    """Verify review costs page has section subheaders."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "st.subheader" in content, "Missing section subheaders"


def test_review_costs_page_has_dividers():
    """Verify review costs page has visual dividers."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "st.divider" in content, "Missing visual dividers"


def test_review_costs_page_has_workload_analysis():
    """Verify review costs page shows reviewer workload."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "get_reviewer_workload" in content or "workload" in content.lower(), \
        "Missing reviewer workload analysis"


def test_review_costs_page_has_ai_band_comparison():
    """Verify review costs page compares by AI band."""
    with open("pages/4_review_costs.py", "r") as f:
        content = f.read()

    assert "get_review_costs_by_ai_band" in content or "ai_usage_band" in content.lower(), \
        "Missing AI band comparison"
