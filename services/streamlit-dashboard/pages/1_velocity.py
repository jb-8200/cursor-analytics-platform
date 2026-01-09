"""
Velocity Metrics Dashboard Page.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-07 Velocity Page

This page displays:
- KPI cards: Total PRs, Avg Cycle Time, Active Developers, AI Ratio
- Line chart: Cycle time trend over time
- Bar chart: Cycle time breakdown by component
- Data table: Weekly velocity metrics

Usage:
    Navigate to "Velocity" page in Streamlit sidebar navigation
"""

import streamlit as st
import plotly.express as px
import plotly.graph_objects as go
from components.sidebar import render_sidebar, get_filter_where_clause
from queries.velocity import get_velocity_data, get_cycle_time_breakdown, get_velocity_summary

# Page configuration
st.set_page_config(
    page_title="Velocity Metrics",
    page_icon="🚀",
    layout="wide"
)

# Render sidebar
render_sidebar()

# Page title
st.title("🚀 Velocity Metrics")
st.markdown("Track PR throughput, cycle times, and developer activity.")

st.divider()

# Get filter WHERE clause from sidebar state
where = get_filter_where_clause()

# Fetch data
try:
    df = get_velocity_data(where)
    summary = get_velocity_summary(where)

    if df.empty:
        st.warning("No data available for the selected filters.")
        st.stop()

    # KPI Cards
    st.subheader("📊 Key Metrics")

    col1, col2, col3, col4 = st.columns(4)

    with col1:
        total_prs = int(summary.get('total_prs', 0))
        st.metric(
            "Total PRs",
            f"{total_prs:,}",
            help="Total number of pull requests merged in the selected period"
        )

    with col2:
        avg_cycle = summary.get('avg_cycle_time', 0)
        st.metric(
            "Avg Cycle Time",
            f"{avg_cycle:.1f} days",
            help="Average time from first commit to merge"
        )

    with col3:
        max_devs = int(summary.get('max_developers', 0))
        st.metric(
            "Active Developers",
            f"{max_devs:,}",
            help="Maximum number of active developers in any week"
        )

    with col4:
        avg_ai = summary.get('avg_ai_ratio', 0)
        st.metric(
            "Avg AI Ratio",
            f"{avg_ai:.0%}",
            help="Average AI contribution ratio across all PRs"
        )

    st.divider()

    # Cycle Time Trend Chart
    st.subheader("📈 Cycle Time Trend")
    st.markdown("Breakdown of cycle time components over time (Coding, Pickup, Review)")

    # Sort by week for proper time series
    df_sorted = df.sort_values("week")

    # Create line chart with multiple components
    fig = go.Figure()

    fig.add_trace(go.Scatter(
        x=df_sorted["week"],
        y=df_sorted["coding_lead_time"],
        mode='lines+markers',
        name='Coding Lead Time',
        line=dict(color='#3498db', width=2),
        hovertemplate='<b>Coding</b><br>%{y:.2f} days<extra></extra>'
    ))

    fig.add_trace(go.Scatter(
        x=df_sorted["week"],
        y=df_sorted["pickup_time"],
        mode='lines+markers',
        name='Pickup Time',
        line=dict(color='#e74c3c', width=2),
        hovertemplate='<b>Pickup</b><br>%{y:.2f} days<extra></extra>'
    ))

    fig.add_trace(go.Scatter(
        x=df_sorted["week"],
        y=df_sorted["review_lead_time"],
        mode='lines+markers',
        name='Review Lead Time',
        line=dict(color='#2ecc71', width=2),
        hovertemplate='<b>Review</b><br>%{y:.2f} days<extra></extra>'
    ))

    fig.update_layout(
        xaxis_title="Week",
        yaxis_title="Days",
        hovermode='x unified',
        legend=dict(
            orientation="h",
            yanchor="bottom",
            y=1.02,
            xanchor="right",
            x=1
        ),
        margin=dict(l=0, r=0, t=30, b=0)
    )

    st.plotly_chart(fig, use_container_width=True)

    st.divider()

    # Cycle Time Breakdown Chart
    st.subheader("⏱️ Average Cycle Time Breakdown")
    st.markdown("Average hours spent in each phase across all PRs")

    breakdown = get_cycle_time_breakdown(where)

    # Create bar chart
    fig2 = px.bar(
        breakdown,
        x="component",
        y="hours",
        color="component",
        text="hours",
        color_discrete_map={
            "Coding": "#3498db",
            "Pickup": "#e74c3c",
            "Review": "#2ecc71"
        }
    )

    fig2.update_traces(
        texttemplate='%{text:.1f}h',
        textposition='outside'
    )

    fig2.update_layout(
        xaxis_title="Component",
        yaxis_title="Hours",
        showlegend=False,
        margin=dict(l=0, r=0, t=30, b=0)
    )

    st.plotly_chart(fig2, use_container_width=True)

    st.divider()

    # Weekly Data Table
    st.subheader("📋 Weekly Velocity Data")
    st.markdown("Detailed weekly metrics (most recent first)")

    # Format the dataframe for display
    display_df = df.copy()
    display_df["week"] = display_df["week"].astype(str)
    display_df["avg_ai_ratio"] = display_df["avg_ai_ratio"].apply(lambda x: f"{x:.1%}")
    display_df["total_cycle_time"] = display_df["total_cycle_time"].apply(lambda x: f"{x:.1f}")
    display_df["p50_cycle_time"] = display_df["p50_cycle_time"].apply(lambda x: f"{x:.1f}")
    display_df["p90_cycle_time"] = display_df["p90_cycle_time"].apply(lambda x: f"{x:.1f}")

    # Select and rename columns for better display
    display_columns = {
        "week": "Week",
        "repo_name": "Repository",
        "total_prs": "PRs",
        "active_developers": "Devs",
        "total_cycle_time": "Cycle Time (days)",
        "p50_cycle_time": "P50 (days)",
        "p90_cycle_time": "P90 (days)",
        "avg_ai_ratio": "AI Ratio"
    }

    display_df = display_df[list(display_columns.keys())].rename(columns=display_columns)

    st.dataframe(
        display_df,
        use_container_width=True,
        hide_index=True
    )

except Exception as e:
    st.error(f"Error loading velocity data: {e}")
    st.exception(e)
