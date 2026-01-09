"""
Code Review Costs Dashboard Page.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-09 Review Costs Page

This page displays:
- KPI cards: Avg Reviews/PR, Review Turnaround, Approval Rate, Comment Density
- Charts for review burden analysis
- Reviewer workload trends
- Review costs by AI usage band

Usage:
    Navigate to "Review Costs" page in Streamlit sidebar navigation
"""

import streamlit as st
import plotly.express as px
import plotly.graph_objects as go
from components.sidebar import render_sidebar, get_filter_where_clause
from queries.review_costs import (
    get_review_costs_data,
    get_review_costs_summary,
    get_reviewer_workload,
    get_review_costs_by_ai_band
)

# Page configuration
st.set_page_config(
    page_title="Review Costs | DOXAPI Analytics",
    page_icon="👥",
    layout="wide"
)

# Render sidebar
render_sidebar()

# Page title
st.title("👥 Code Review Costs")
st.markdown("Analyze review iterations, reviewer workload, and review time metrics.")

st.divider()

# Get filter WHERE clause from sidebar state
where = get_filter_where_clause()

# Fetch data
try:
    df = get_review_costs_data(where)
    summary = get_review_costs_summary(where)

    if df.empty:
        st.warning("No data available for the selected filters.")
        st.stop()

    # --- KPI Cards ---
    st.subheader("📊 Key Metrics")

    col1, col2, col3, col4 = st.columns(4)

    with col1:
        avg_iterations = summary.get('avg_iterations', 0) or 0
        st.metric(
            "Avg Iterations",
            f"{avg_iterations:.1f}",
            help="Average number of review cycles per PR"
        )

    with col2:
        avg_reviewers = summary.get('avg_reviewers', 0) or 0
        st.metric(
            "Avg Reviewers/PR",
            f"{avg_reviewers:.1f}",
            help="Average number of reviewers per pull request"
        )

    with col3:
        avg_comments = summary.get('avg_comments', 0) or 0
        st.metric(
            "Avg Comments/PR",
            f"{avg_comments:.1f}",
            help="Average number of review comments per PR"
        )

    with col4:
        total_hours = int(summary.get('total_hours', 0) or 0)
        st.metric(
            "Total Review Hours",
            f"{total_hours:,}h",
            help="Total hours spent in code review"
        )

    st.divider()

    # --- Review Time Trend ---
    st.subheader("📈 Review Time Trend")
    st.markdown("Track how review turnaround time changes over time")

    workload = get_reviewer_workload(where)

    if not workload.empty:
        workload_sorted = workload.sort_values("week")

        fig_trend = go.Figure()

        fig_trend.add_trace(go.Scatter(
            x=workload_sorted["week"],
            y=workload_sorted["avg_review_time"],
            mode='lines+markers',
            name='Avg Review Time',
            line=dict(color='#3498db', width=2),
            hovertemplate='<b>Review Time</b><br>%{y:.2f} days<extra></extra>'
        ))

        fig_trend.update_layout(
            title="Weekly Average Review Time",
            xaxis_title="Week",
            yaxis_title="Review Time (days)",
            hovermode='x unified',
            margin=dict(l=0, r=0, t=30, b=0)
        )

        st.plotly_chart(fig_trend, use_container_width=True)
    else:
        st.info("No workload trend data available.")

    st.divider()

    # --- Review Hours by Week ---
    st.subheader("⏱️ Total Review Hours by Week")
    st.markdown("Understand the aggregate review burden on your team")

    if not workload.empty:
        workload_sorted = workload.sort_values("week")

        fig_hours = px.bar(
            workload_sorted,
            x="week",
            y="total_review_hours",
            color="total_review_hours",
            color_continuous_scale="Blues",
            title="Weekly Total Review Hours",
            labels={
                "week": "Week",
                "total_review_hours": "Review Hours"
            }
        )

        fig_hours.update_layout(
            xaxis_title="Week",
            yaxis_title="Total Review Hours",
            showlegend=False,
            coloraxis_showscale=False,
            margin=dict(l=0, r=0, t=30, b=0)
        )

        st.plotly_chart(fig_hours, use_container_width=True)
    else:
        st.info("No review hours data available.")

    st.divider()

    # --- Review Costs by AI Usage Band ---
    st.subheader("🤖 Review Costs by AI Usage Band")
    st.markdown("Compare review burden across different AI usage levels")

    costs_by_band = get_review_costs_by_ai_band(where)

    if not costs_by_band.empty:
        # Sort for proper ordering
        band_order = {"low": 0, "medium": 1, "high": 2}
        costs_by_band["order"] = costs_by_band["ai_usage_band"].map(band_order)
        costs_by_band = costs_by_band.sort_values("order")

        # Create grouped bar chart
        fig_band = go.Figure()

        fig_band.add_trace(go.Bar(
            x=costs_by_band["ai_usage_band"],
            y=costs_by_band["avg_review_iterations"],
            name='Avg Iterations',
            marker_color='#3498db',
            text=costs_by_band["avg_review_iterations"].apply(lambda x: f"{x:.1f}"),
            textposition='outside'
        ))

        fig_band.add_trace(go.Bar(
            x=costs_by_band["ai_usage_band"],
            y=costs_by_band["avg_review_comments"],
            name='Avg Comments',
            marker_color='#2ecc71',
            text=costs_by_band["avg_review_comments"].apply(lambda x: f"{x:.1f}"),
            textposition='outside'
        ))

        fig_band.update_layout(
            title="Review Metrics by AI Usage Band",
            xaxis_title="AI Usage Band",
            yaxis_title="Count",
            barmode='group',
            xaxis=dict(categoryorder='array', categoryarray=['low', 'medium', 'high']),
            legend=dict(
                orientation="h",
                yanchor="bottom",
                y=1.02,
                xanchor="right",
                x=1
            ),
            margin=dict(l=0, r=0, t=30, b=0)
        )

        st.plotly_chart(fig_band, use_container_width=True)

        # Review time comparison
        fig_time_band = px.bar(
            costs_by_band,
            x="ai_usage_band",
            y="avg_review_time",
            color="ai_usage_band",
            title="Average Review Time by AI Usage Band",
            category_orders={"ai_usage_band": ["low", "medium", "high"]},
            color_discrete_map={
                "low": "#e74c3c",
                "medium": "#f39c12",
                "high": "#2ecc71"
            },
            text="avg_review_time",
            labels={
                "ai_usage_band": "AI Usage Band",
                "avg_review_time": "Review Time (days)"
            }
        )

        fig_time_band.update_traces(
            texttemplate='%{text:.1f} days',
            textposition='outside'
        )

        fig_time_band.update_layout(
            xaxis_title="AI Usage Band",
            yaxis_title="Review Time (days)",
            showlegend=False,
            margin=dict(l=0, r=0, t=30, b=0)
        )

        st.plotly_chart(fig_time_band, use_container_width=True)
    else:
        st.info("No AI band data available.")

    st.divider()

    # --- Review Costs by Repository ---
    st.subheader("📊 Review Costs by Repository")
    st.markdown("Compare review burden across different repositories")

    # Aggregate by repository
    repo_summary = df.groupby("repo_name").agg({
        "total_prs": "sum",
        "avg_review_iterations": "mean",
        "avg_reviewers_per_pr": "mean",
        "avg_review_comments": "mean",
        "total_review_hours": "sum"
    }).reset_index()

    if not repo_summary.empty:
        # Sort by total review hours descending
        repo_summary_sorted = repo_summary.sort_values("total_review_hours", ascending=False)

        fig_repo = px.bar(
            repo_summary_sorted,
            x="repo_name",
            y="total_review_hours",
            color="avg_review_iterations",
            color_continuous_scale="Oranges",
            text="total_review_hours",
            title="Total Review Hours by Repository",
            labels={
                "repo_name": "Repository",
                "total_review_hours": "Total Review Hours",
                "avg_review_iterations": "Avg Iterations"
            }
        )

        fig_repo.update_traces(
            texttemplate='%{text:.0f}h',
            textposition='outside'
        )

        fig_repo.update_layout(
            xaxis_title="Repository",
            yaxis_title="Total Review Hours",
            coloraxis_colorbar_title="Avg Iterations",
            margin=dict(l=0, r=0, t=30, b=0)
        )

        st.plotly_chart(fig_repo, use_container_width=True)
    else:
        st.info("No repository data available.")

    st.divider()

    # --- Weekly Review Data Table ---
    st.subheader("📋 Weekly Review Data")
    st.markdown("Detailed weekly metrics (most recent first)")

    # Format the dataframe for display
    display_df = df.copy()
    display_df["week"] = display_df["week"].astype(str)
    display_df["avg_review_iterations"] = display_df["avg_review_iterations"].apply(lambda x: f"{x:.1f}")
    display_df["avg_reviewers_per_pr"] = display_df["avg_reviewers_per_pr"].apply(lambda x: f"{x:.1f}")
    display_df["avg_review_comments"] = display_df["avg_review_comments"].apply(lambda x: f"{x:.1f}")
    display_df["avg_review_time"] = display_df["avg_review_time"].apply(lambda x: f"{x:.1f}")

    # Select and rename columns for better display
    display_columns = {
        "week": "Week",
        "repo_name": "Repository",
        "total_prs": "PRs",
        "avg_review_iterations": "Iterations",
        "avg_reviewers_per_pr": "Reviewers/PR",
        "avg_review_comments": "Comments/PR",
        "avg_review_time": "Review Time (days)",
        "total_review_hours": "Total Hours"
    }

    display_df = display_df[list(display_columns.keys())].rename(columns=display_columns)

    st.dataframe(
        display_df,
        use_container_width=True,
        hide_index=True
    )

except Exception as e:
    st.error(f"Error loading review costs data: {e}")
    st.exception(e)
