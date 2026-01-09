"""
DOXAPI Analytics Dashboard - Main Entry Point

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-06 Home Page
Author: Claude (streamlit-dev agent)
Created: 2026-01-09
"""

import streamlit as st

# Configure page
st.set_page_config(
    page_title="DOXAPI Analytics",
    page_icon="📊",
    layout="wide",
    initial_sidebar_state="expanded"
)

# Main content
st.title("📊 DOXAPI Analytics Dashboard")

# KPI metrics row
col1, col2, col3, col4 = st.columns(4)

# KPI placeholders (will connect to dbt marts when available)
with col1:
    st.metric(label="Total PRs", value="--", delta=None)

with col2:
    st.metric(label="Avg Cycle Time", value="-- days", delta=None)

with col3:
    st.metric(label="Avg Revert Rate", value="--%", delta=None)

with col4:
    st.metric(label="Avg AI Ratio", value="--%", delta=None)

st.markdown("---")

st.markdown("""
## Welcome to the AI Code Analytics Dashboard

This dashboard provides insights into AI-assisted coding impact on:

- **Velocity**: PR cycle times and throughput
- **AI Impact**: Metrics by AI usage bands
- **Quality**: Revert rates and quality trends
- **Review Costs**: Code review burden analysis

### Getting Started

Navigate using the sidebar to explore different metrics.
""")

# Status info
st.info("📦 P9-F01: Home page with KPI metrics. Connect to dbt marts for live data.")
