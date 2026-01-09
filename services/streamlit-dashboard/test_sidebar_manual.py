"""
Manual test for sidebar component.

Run this with: streamlit run test_sidebar_manual.py

This will verify:
- Sidebar renders without errors
- Repository filter populates
- Date range filter works
- Refresh button appears (dev mode)
- Environment indicator shows
"""

import streamlit as st
from components.sidebar import render_sidebar, get_filter_where_clause

st.set_page_config(
    page_title="Sidebar Test",
    page_icon="🧪",
    layout="wide"
)

# Render the sidebar
render_sidebar()

# Main content
st.title("🧪 Sidebar Component Test")

st.markdown("""
This page tests the sidebar component.

**Check the sidebar** (left) for:
1. ✅ Title: "DOXAPI Analytics"
2. ✅ Repository filter dropdown
3. ✅ Date range filter dropdown
4. ✅ Refresh button (dev mode) OR info message (prod mode)
5. ✅ Environment indicator (Development/Production)
""")

st.divider()

# Show current filter state
st.subheader("Current Filter State")

col1, col2 = st.columns(2)

with col1:
    st.write("**Repository:**", st.session_state.get("filter_repo", "Not set"))
    st.write("**Date Range:**", st.session_state.get("filter_date_range", "Not set"))

with col2:
    where_clause = get_filter_where_clause()
    st.write("**SQL WHERE Clause:**")
    if where_clause:
        st.code(where_clause, language="sql")
    else:
        st.code("(no filter - all data)", language="text")

st.divider()

# Test instructions
st.subheader("Test Instructions")
st.markdown("""
1. **Repository Filter**: Select different repositories and verify the WHERE clause updates
2. **Date Range Filter**: Change the date range and verify the clause includes correct interval
3. **Refresh Button**: Click it in dev mode (should trigger pipeline)
4. **Environment**: Verify shows "Development" for DuckDB mode

**Expected Behavior:**
- Selecting "All" repos → no repo_name filter
- Selecting specific repo → `WHERE repo_name = 'repo/name'`
- Selecting "All time" → no date filter
- Selecting "Last N days" → `WHERE week >= CURRENT_DATE - INTERVAL 'N days'`
""")
