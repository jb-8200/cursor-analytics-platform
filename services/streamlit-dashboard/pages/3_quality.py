"""
Code Quality Metrics Dashboard Page.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-09 Quality Page

This page displays:
- KPI cards: Revert Rate, Hotfix Rate, Avg Time to Revert, Total Reverted PRs
- Line chart: Quality metrics trend over time
- Bar chart: Revert rates by repository
- Weekly quality data table

Usage:
    Navigate to "Quality" page in Streamlit sidebar navigation
"""

import streamlit as st
import plotly.express as px
import plotly.graph_objects as go
from components.sidebar import render_sidebar, get_filter_params
from queries.quality import (
    get_quality_data,
    get_quality_summary,
    get_revert_trends,
    get_quality_by_ai_band
)

# Page configuration
st.set_page_config(
    page_title="Quality Metrics | DOXAPI Analytics",
    page_icon="✅",
    layout="wide"
)

# Render sidebar
render_sidebar()

# Page title
st.title("✅ Code Quality Metrics")
st.markdown("Track revert rates, bug fixes, and quality trends over time.")

st.divider()

# Get filter params from sidebar state
repo_name, date_range, days = get_filter_params()

# Fetch data
try:
    df = get_quality_data(repo_name=repo_name, days=days)
    summary = get_quality_summary(repo_name=repo_name, days=days)

    if df.empty:
        st.warning("No data available for the selected filters.")
        st.stop()

    # --- KPI Cards ---
    st.subheader("📊 Key Metrics")

    col1, col2, col3, col4 = st.columns(4)

    with col1:
        avg_revert_rate = summary.get('avg_revert_rate', 0) or 0
        st.metric(
            "Revert Rate",
            f"{avg_revert_rate:.1%}",
            help="Average percentage of PRs that were reverted"
        )

    with col2:
        avg_bug_fix_rate = summary.get('avg_bug_fix_rate', 0) or 0
        st.metric(
            "Bug Fix Rate",
            f"{avg_bug_fix_rate:.1%}",
            help="Average percentage of PRs that are bug fixes"
        )

    with col3:
        avg_reviews = summary.get('avg_reviews_per_pr', 0) or 0
        st.metric(
            "Avg Reviews per PR",
            f"{avg_reviews:.1f}",
            help="Average number of reviews per pull request"
        )

    with col4:
        total_reverted = int(summary.get('total_reverted', 0) or 0)
        st.metric(
            "Total Reverted PRs",
            f"{total_reverted:,}",
            help="Total number of PRs that were reverted"
        )

    st.divider()

    # --- Revert Rate Trend Over Time ---
    st.subheader("📈 Revert Rate Trend")
    st.markdown("Track how revert rates change week over week")

    trends = get_revert_trends(repo_name=repo_name, days=days)

    if not trends.empty:
        trends_sorted = trends.sort_values("week")

        # Aggregate by week if multiple repos
        weekly_trend = trends_sorted.groupby("week").agg({
            "revert_rate": "mean",
            "reverted_prs": "sum"
        }).reset_index()

        fig_trend = go.Figure()

        fig_trend.add_trace(go.Scatter(
            x=weekly_trend["week"],
            y=weekly_trend["revert_rate"],
            mode='lines+markers',
            name='Revert Rate',
            line=dict(color='#e74c3c', width=2),
            hovertemplate='<b>Revert Rate</b><br>%{y:.2%}<extra></extra>'
        ))

        fig_trend.update_layout(
            title="Weekly Revert Rate Trend",
            xaxis_title="Week",
            yaxis_title="Revert Rate",
            yaxis_tickformat=".1%",
            hovermode='x unified',
            margin=dict(l=0, r=0, t=30, b=0)
        )

        st.plotly_chart(fig_trend, use_container_width=True)
    else:
        st.info("No trend data available.")

    st.divider()

    # --- Revert Rate by Repository ---
    st.subheader("📊 Revert Rate by Repository")
    st.markdown("Compare quality metrics across different repositories")

    # Aggregate by repository
    repo_summary = df.groupby("repo_name").agg({
        "total_prs": "sum",
        "reverted_prs": "sum",
        "revert_rate": "mean"
    }).reset_index()

    if not repo_summary.empty:
        # Sort by revert rate descending
        repo_summary_sorted = repo_summary.sort_values("revert_rate", ascending=False)

        fig_bar = px.bar(
            repo_summary_sorted,
            x="repo_name",
            y="revert_rate",
            color="revert_rate",
            color_continuous_scale=["#2ecc71", "#f39c12", "#e74c3c"],
            text="revert_rate",
            title="Average Revert Rate by Repository",
            labels={
                "repo_name": "Repository",
                "revert_rate": "Revert Rate"
            }
        )

        fig_bar.update_traces(
            texttemplate='%{text:.1%}',
            textposition='outside'
        )

        fig_bar.update_layout(
            xaxis_title="Repository",
            yaxis_title="Revert Rate",
            yaxis_tickformat=".1%",
            showlegend=False,
            coloraxis_showscale=False,
            margin=dict(l=0, r=0, t=30, b=0)
        )

        st.plotly_chart(fig_bar, use_container_width=True)
    else:
        st.info("No repository data available.")

    st.divider()

    # --- Quality by AI Usage Band ---
    st.subheader("🤖 Quality by AI Usage Band")
    st.markdown("Compare revert rates across different AI usage levels")

    quality_by_band = get_quality_by_ai_band(repo_name=repo_name, days=days)

    if not quality_by_band.empty:
        # Sort for proper ordering
        band_order = {"low": 0, "medium": 1, "high": 2}
        quality_by_band["order"] = quality_by_band["ai_usage_band"].map(band_order)
        quality_by_band = quality_by_band.sort_values("order")

        fig_band = px.bar(
            quality_by_band,
            x="ai_usage_band",
            y="avg_revert_rate",
            color="ai_usage_band",
            title="Revert Rate by AI Usage Band",
            category_orders={"ai_usage_band": ["low", "medium", "high"]},
            color_discrete_map={
                "low": "#e74c3c",
                "medium": "#f39c12",
                "high": "#2ecc71"
            },
            text="avg_revert_rate",
            labels={
                "ai_usage_band": "AI Usage Band",
                "avg_revert_rate": "Revert Rate"
            }
        )

        fig_band.update_traces(
            texttemplate='%{text:.1%}',
            textposition='outside'
        )

        fig_band.update_layout(
            xaxis_title="AI Usage Band",
            yaxis_title="Revert Rate",
            yaxis_tickformat=".1%",
            showlegend=False,
            margin=dict(l=0, r=0, t=30, b=0)
        )

        st.plotly_chart(fig_band, use_container_width=True)
    else:
        st.info("No AI band data available.")

    st.divider()

    # --- Weekly Quality Data Table ---
    st.subheader("📋 Weekly Quality Data")
    st.markdown("Detailed weekly metrics (most recent first)")

    # Format the dataframe for display
    display_df = df.copy()
    display_df["week"] = display_df["week"].astype(str)
    display_df["revert_rate"] = display_df["revert_rate"].apply(lambda x: f"{x:.1%}")
    display_df["bug_fix_rate"] = display_df["bug_fix_rate"].apply(lambda x: f"{x:.1%}")
    display_df["avg_reviews_per_pr"] = display_df["avg_reviews_per_pr"].apply(lambda x: f"{x:.1f}")

    # Select and rename columns for better display
    display_columns = {
        "week": "Week",
        "repo_name": "Repository",
        "total_prs": "Total PRs",
        "reverted_prs": "Reverted",
        "revert_rate": "Revert Rate",
        "bug_fix_prs": "Bug Fixes",
        "bug_fix_rate": "Bug Fix Rate",
        "avg_reviews_per_pr": "Avg Reviews/PR"
    }

    display_df = display_df[list(display_columns.keys())].rename(columns=display_columns)

    st.dataframe(
        display_df,
        use_container_width=True,
        hide_index=True
    )

except Exception as e:
    st.error(f"Error loading quality data: {e}")
    st.exception(e)
