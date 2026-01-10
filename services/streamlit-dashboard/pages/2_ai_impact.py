"""
AI Impact Analysis Dashboard Page.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-08 AI Impact Page

This page displays:
- KPI cards: Total PRs, Avg AI Ratio, High AI PRs %, Composer vs Tab ratio
- Box plot: Cycle times across AI bands (low/medium/high)
- Bar chart: AI ratio distribution by repository
- Band comparison table

Usage:
    Navigate to "AI Impact" page in Streamlit sidebar navigation
"""

import streamlit as st
import plotly.express as px
import plotly.graph_objects as go
from components.sidebar import render_sidebar, get_filter_params
from queries.ai_impact import get_ai_impact_data, get_band_comparison, get_band_trends

# Page configuration
st.set_page_config(
    page_title="AI Impact | DOXAPI Analytics",
    page_icon="🤖",
    layout="wide"
)

# Render sidebar
render_sidebar()

# Page title
st.title("🤖 AI Impact Analysis")
st.markdown("Analyze how AI-assisted development affects cycle times, quality, and productivity.")

st.divider()

# Get filter params from sidebar state
repo_name, date_range, days = get_filter_params()

# Fetch data
try:
    df = get_ai_impact_data(repo_name=repo_name, days=days)
    band_comparison = get_band_comparison(repo_name=repo_name, days=days)

    if df.empty:
        st.warning("No data available for the selected filters.")
        st.stop()

    # --- KPI Cards ---
    st.subheader("📊 Key Metrics")

    col1, col2, col3, col4 = st.columns(4)

    with col1:
        total_prs = int(df["pr_count"].sum())
        st.metric(
            "Total PRs",
            f"{total_prs:,}",
            help="Total number of pull requests analyzed"
        )

    with col2:
        avg_ai_ratio = df["avg_ai_ratio"].mean()
        st.metric(
            "Avg AI Ratio",
            f"{avg_ai_ratio:.0%}",
            help="Average AI contribution ratio across all PRs"
        )

    with col3:
        # Calculate high AI PRs percentage
        if not band_comparison.empty:
            high_ai_prs = band_comparison[band_comparison["ai_usage_band"] == "high"]["total_prs"].sum()
            total = band_comparison["total_prs"].sum()
            high_ai_pct = high_ai_prs / total if total > 0 else 0
        else:
            high_ai_pct = 0
        st.metric(
            "High AI PRs",
            f"{high_ai_pct:.0%}",
            help="Percentage of PRs with high AI usage (>50%)"
        )

    with col4:
        # Average review cycle time
        avg_review_time = df["avg_review_cycle_time"].mean()
        st.metric(
            "Avg Review Time",
            f"{avg_review_time:.1f} days",
            help="Average review cycle time across all bands"
        )

    st.divider()

    # --- Band Comparison Table ---
    st.subheader("📋 Metrics by AI Usage Band")
    st.markdown("Summary statistics grouped by AI usage level (low <25%, medium 25-50%, high >50%)")

    if not band_comparison.empty:
        # Format for display
        display_df = band_comparison.copy()
        display_df["ai_usage_band"] = display_df["ai_usage_band"].str.capitalize()
        display_df["avg_ai_ratio"] = display_df["avg_ai_ratio"].apply(lambda x: f"{x:.1%}")
        display_df["avg_coding_lead_time"] = display_df["avg_coding_lead_time"].apply(lambda x: f"{x:.1f} days")
        display_df["avg_review_cycle_time"] = display_df["avg_review_cycle_time"].apply(lambda x: f"{x:.1f} days")
        display_df["avg_revert_rate"] = display_df["avg_revert_rate"].apply(lambda x: f"{x:.1%}")

        display_columns = {
            "ai_usage_band": "AI Usage Band",
            "total_prs": "PRs",
            "avg_ai_ratio": "Avg AI Ratio",
            "avg_coding_lead_time": "Coding Lead Time",
            "avg_review_cycle_time": "Review Cycle Time",
            "avg_revert_rate": "Revert Rate"
        }

        display_df = display_df[list(display_columns.keys())].rename(columns=display_columns)

        st.dataframe(
            display_df,
            use_container_width=True,
            hide_index=True
        )
    else:
        st.info("No band comparison data available.")

    st.divider()

    # --- Box Plot: Cycle Time by AI Band ---
    st.subheader("📦 Cycle Time Distribution by AI Band")
    st.markdown("Compare coding lead time distribution across AI usage levels")

    # Create box plot
    fig_box = px.box(
        df,
        x="ai_usage_band",
        y="avg_coding_lead_time",
        color="ai_usage_band",
        title="Coding Lead Time by AI Usage Band",
        category_orders={"ai_usage_band": ["low", "medium", "high"]},
        color_discrete_map={
            "low": "#e74c3c",
            "medium": "#f39c12",
            "high": "#2ecc71"
        },
        labels={
            "ai_usage_band": "AI Usage Band",
            "avg_coding_lead_time": "Coding Lead Time (days)"
        }
    )

    fig_box.update_layout(
        xaxis_title="AI Usage Band",
        yaxis_title="Coding Lead Time (days)",
        showlegend=False,
        margin=dict(l=0, r=0, t=30, b=0)
    )

    st.plotly_chart(fig_box, use_container_width=True)

    st.divider()

    # --- Bar Chart: Revert Rate by AI Band ---
    st.subheader("📊 Revert Rate by AI Usage Band")
    st.markdown("Compare quality metrics (revert rates) across AI usage levels")

    if not band_comparison.empty:
        # Sort for proper ordering
        band_order = {"low": 0, "medium": 1, "high": 2}
        band_comparison_sorted = band_comparison.copy()
        band_comparison_sorted["order"] = band_comparison_sorted["ai_usage_band"].map(band_order)
        band_comparison_sorted = band_comparison_sorted.sort_values("order")

        fig_bar = px.bar(
            band_comparison_sorted,
            x="ai_usage_band",
            y="avg_revert_rate",
            color="ai_usage_band",
            title="Average Revert Rate by AI Band",
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

        fig_bar.update_traces(
            texttemplate='%{text:.1%}',
            textposition='outside'
        )

        fig_bar.update_layout(
            xaxis_title="AI Usage Band",
            yaxis_title="Revert Rate",
            yaxis_tickformat=".1%",
            showlegend=False,
            margin=dict(l=0, r=0, t=30, b=0)
        )

        st.plotly_chart(fig_bar, use_container_width=True)
    else:
        st.info("No data available for revert rate chart.")

    st.divider()

    # --- AI Usage Trend Over Time ---
    st.subheader("📈 AI Usage Trend Over Time")
    st.markdown("Track how AI usage distribution changes week over week")

    trends = get_band_trends(repo_name=repo_name, days=days)

    if not trends.empty:
        # Pivot for stacked area chart
        trends_sorted = trends.sort_values("week")

        fig_trend = px.area(
            trends_sorted,
            x="week",
            y="pr_count",
            color="ai_usage_band",
            title="PR Count by AI Usage Band Over Time",
            category_orders={"ai_usage_band": ["low", "medium", "high"]},
            color_discrete_map={
                "low": "#e74c3c",
                "medium": "#f39c12",
                "high": "#2ecc71"
            },
            labels={
                "week": "Week",
                "pr_count": "Number of PRs",
                "ai_usage_band": "AI Usage Band"
            }
        )

        fig_trend.update_layout(
            xaxis_title="Week",
            yaxis_title="Number of PRs",
            legend=dict(
                orientation="h",
                yanchor="bottom",
                y=1.02,
                xanchor="right",
                x=1
            ),
            margin=dict(l=0, r=0, t=30, b=0)
        )

        st.plotly_chart(fig_trend, use_container_width=True)
    else:
        st.info("No trend data available.")

except Exception as e:
    st.error(f"Error loading AI impact data: {e}")
    st.exception(e)
