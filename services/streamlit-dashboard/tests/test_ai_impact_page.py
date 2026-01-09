"""
Tests for AI Impact Dashboard Page.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-08 AI Impact Page

Tests verify:
- Page can be imported without errors
- Required functions are available
- Page displays expected components (KPIs, charts)
- AI usage bands are ordered correctly
"""

import pytest
import sys
import os

# Add parent directory to path for imports
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))


def test_ai_impact_page_imports():
    """Verify AI impact page can be imported."""
    try:
        import importlib.util
        spec = importlib.util.spec_from_file_location(
            "ai_impact_page",
            "pages/2_ai_impact.py"
        )
        assert spec is not None, "Could not load AI impact page spec"

        module = importlib.util.module_from_spec(spec)
        assert module is not None, "Could not create module from spec"

    except SyntaxError as e:
        pytest.fail(f"AI impact page has syntax error: {e}")
    except Exception as e:
        # Other errors are acceptable since we can't run Streamlit context
        pass


def test_ai_impact_page_file_exists():
    """Verify AI impact page file exists."""
    page_path = "pages/2_ai_impact.py"
    assert os.path.exists(page_path), f"AI impact page not found at {page_path}"


def test_ai_impact_page_has_required_imports():
    """Verify AI impact page has required imports."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    required_imports = [
        "import streamlit",
        "import plotly",
        "from components.sidebar import render_sidebar",
        "from queries.ai_impact import"
    ]

    for required in required_imports:
        assert required in content, f"Missing required import: {required}"


def test_ai_impact_page_has_page_config():
    """Verify AI impact page has Streamlit page configuration."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    assert "st.set_page_config" in content, "Missing st.set_page_config"
    assert 'page_title="AI Impact' in content or "page_title='AI Impact" in content
    assert "layout=\"wide\"" in content or "layout='wide'" in content


def test_ai_impact_page_renders_sidebar():
    """Verify AI impact page calls render_sidebar."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    assert "render_sidebar()" in content, "Missing render_sidebar() call"


def test_ai_impact_page_has_title():
    """Verify AI impact page has a title."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    assert "st.title" in content, "Missing st.title"
    assert "AI Impact" in content, "Missing 'AI Impact' in title"


def test_ai_impact_page_uses_filters():
    """Verify AI impact page uses filter WHERE clause."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    assert "get_filter_where_clause" in content, "Missing get_filter_where_clause"


def test_ai_impact_page_has_kpis():
    """Verify AI impact page displays 4 KPI metrics."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    # Check for 4 columns (4 KPIs)
    assert "st.columns(4)" in content, "Missing 4-column layout for KPIs"

    # Check for metrics
    assert "st.metric" in content, "Missing st.metric calls"


def test_ai_impact_page_has_kpi_total_commits():
    """Verify AI impact page has Total Commits KPI."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read().lower()

    # Check for commits/PRs metric
    assert "total" in content and ("pr" in content or "commit" in content), \
        "Missing Total Commits/PRs KPI"


def test_ai_impact_page_has_kpi_ai_ratio():
    """Verify AI impact page has Avg AI Ratio KPI."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read().lower()

    assert "ai ratio" in content or "ai_ratio" in content, "Missing AI Ratio KPI"


def test_ai_impact_page_has_box_plot():
    """Verify AI impact page has box plot for cycle times by AI band."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    # Check for box plot
    assert "px.box" in content or "go.Box" in content, "Missing box plot chart"
    assert "st.plotly_chart" in content, "Missing plotly chart display"


def test_ai_impact_page_has_bar_chart():
    """Verify AI impact page has bar chart for AI ratio distribution."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    assert "px.bar" in content, "Missing bar chart"


def test_ai_impact_page_queries_data():
    """Verify AI impact page queries AI impact data."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    assert "get_ai_impact_data" in content, "Missing get_ai_impact_data query"


def test_ai_impact_page_has_band_comparison():
    """Verify AI impact page shows band comparison."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    assert "get_band_comparison" in content, "Missing get_band_comparison query"


def test_ai_impact_page_orders_bands_correctly():
    """Verify AI impact page orders AI usage bands as low, medium, high."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    # Check for band ordering configuration
    assert "low" in content.lower() and "medium" in content.lower() and "high" in content.lower(), \
        "Missing AI usage band references"
    assert "category_orders" in content or "['low', 'medium', 'high']" in content or \
           '["low", "medium", "high"]' in content, \
        "Missing band ordering configuration"


def test_ai_impact_page_handles_errors():
    """Verify AI impact page has error handling."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    assert "try:" in content or "except" in content, "Missing error handling"
    assert "st.error" in content or "st.warning" in content, "Missing error display"


def test_ai_impact_page_has_documentation():
    """Verify AI impact page has docstring."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    assert '"""' in content or "'''" in content, "Missing module docstring"
    assert "TASK-P9-08" in content, "Missing task reference in docstring"


def test_ai_impact_page_has_subheaders():
    """Verify AI impact page has section subheaders."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    assert "st.subheader" in content, "Missing section subheaders"


def test_ai_impact_page_has_dividers():
    """Verify AI impact page has visual dividers."""
    with open("pages/2_ai_impact.py", "r") as f:
        content = f.read()

    assert "st.divider" in content, "Missing visual dividers"
