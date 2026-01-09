# Technical Design: Streamlit Analytics Dashboard

**Feature ID**: P9-F01-streamlit-dashboard
**Phase**: P9 (Streamlit Dashboard)
**Created**: January 9, 2026
**Status**: NOT_STARTED

## Overview

This feature implements a production-ready Streamlit dashboard that visualizes AI coding analytics from the P8 data tier. It supports both DuckDB (local development) and Snowflake (production) backends.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      P9: STREAMLIT DASHBOARD ARCHITECTURE                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                           PRESENTATION LAYER                            │ │
│  │  ┌──────────────────────────────────────────────────────────────────┐ │ │
│  │  │  Streamlit App (services/streamlit-dashboard/)                    │ │ │
│  │  │  ├── app.py                 # Main entry point                    │ │ │
│  │  │  ├── pages/                                                       │ │ │
│  │  │  │   ├── 1_velocity.py      # Velocity metrics                   │ │ │
│  │  │  │   ├── 2_ai_impact.py     # AI impact analysis                 │ │ │
│  │  │  │   ├── 3_quality.py       # Code quality metrics               │ │ │
│  │  │  │   └── 4_review_costs.py  # Review cost analysis               │ │ │
│  │  │  └── components/                                                  │ │ │
│  │  │      ├── sidebar.py         # Shared sidebar                     │ │ │
│  │  │      ├── metrics.py         # KPI cards                          │ │ │
│  │  │      └── charts.py          # Chart helpers                      │ │ │
│  │  └──────────────────────────────────────────────────────────────────┘ │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                       │
│                                      │ SQL Queries                           │
│                                      ▼                                       │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                           DATA ACCESS LAYER                             │ │
│  │  ┌──────────────────────────────────────────────────────────────────┐ │ │
│  │  │  db/connector.py                                                  │ │ │
│  │  │  ├── get_connection()      # DuckDB or Snowflake based on mode   │ │ │
│  │  │  ├── query(sql)            # Execute query, return DataFrame     │ │ │
│  │  │  └── refresh_data()        # Trigger pipeline (dev mode only)    │ │ │
│  │  └──────────────────────────────────────────────────────────────────┘ │ │
│  │                                                                         │ │
│  │  ┌──────────────────────────────────────────────────────────────────┐ │ │
│  │  │  queries/                                                         │ │ │
│  │  │  ├── velocity.py           # Velocity SQL queries                │ │ │
│  │  │  ├── ai_impact.py          # AI impact SQL queries               │ │ │
│  │  │  ├── quality.py            # Quality SQL queries                 │ │ │
│  │  │  └── review_costs.py       # Review costs SQL queries            │ │ │
│  │  └──────────────────────────────────────────────────────────────────┘ │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                       │
│                   ┌──────────────────┴──────────────────┐                   │
│                   │                                      │                   │
│                   ▼                                      ▼                   │
│  ┌─────────────────────────────┐    ┌─────────────────────────────┐        │
│  │  DuckDB (Dev Mode)          │    │  Snowflake (Prod Mode)      │        │
│  │  data/analytics.duckdb      │    │  CURSOR_ANALYTICS.MART.*    │        │
│  │                             │    │                             │        │
│  │  Tables:                    │    │  Tables:                    │        │
│  │  - mart.velocity            │    │  - MART.VELOCITY            │        │
│  │  - mart.ai_impact           │    │  - MART.AI_IMPACT           │        │
│  │  - mart.quality             │    │  - MART.QUALITY             │        │
│  │  - mart.review_costs        │    │  - MART.REVIEW_COSTS        │        │
│  └─────────────────────────────┘    └─────────────────────────────┘        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Directory Structure

```
services/streamlit-dashboard/
├── app.py                      # Main Streamlit entrypoint
├── Dockerfile                  # Production container
├── requirements.txt            # Python dependencies
│
├── pages/                      # Multi-page Streamlit app
│   ├── 1_velocity.py          # Velocity metrics page
│   ├── 2_ai_impact.py         # AI impact analysis page
│   ├── 3_quality.py           # Code quality page
│   └── 4_review_costs.py      # Review costs page
│
├── components/                 # Reusable UI components
│   ├── __init__.py
│   ├── sidebar.py             # Shared sidebar with filters
│   ├── metrics.py             # KPI card components
│   └── charts.py              # Chart helper functions
│
├── db/                         # Database layer
│   ├── __init__.py
│   └── connector.py           # DuckDB/Snowflake abstraction
│
├── queries/                    # SQL queries (parameterized)
│   ├── __init__.py
│   ├── velocity.py
│   ├── ai_impact.py
│   ├── quality.py
│   └── review_costs.py
│
├── pipeline/                   # Embedded pipeline (dev mode)
│   ├── __init__.py
│   ├── loader.py              # Symlink to tools/api-loader/loader.py
│   └── run_dbt.py             # dbt execution wrapper
│
└── tests/
    ├── __init__.py
    ├── test_connector.py
    ├── test_queries.py
    └── test_components.py
```

---

## Component Details

### 1. Database Connector (`db/connector.py`)

**Purpose**: Abstract DuckDB/Snowflake connection based on environment.

```python
# db/connector.py
import os
from functools import lru_cache
import streamlit as st
import pandas as pd

DB_MODE = os.getenv("DB_MODE", "duckdb")


@st.cache_resource
def get_connection():
    """Get database connection based on DB_MODE environment variable."""
    if DB_MODE == "snowflake":
        return _get_snowflake_connection()
    else:
        return _get_duckdb_connection()


def _get_duckdb_connection():
    """Local DuckDB connection."""
    import duckdb
    db_path = os.getenv("DUCKDB_PATH", "/data/analytics.duckdb")
    return duckdb.connect(db_path, read_only=False)


def _get_snowflake_connection():
    """Production Snowflake connection."""
    import snowflake.connector
    return snowflake.connector.connect(
        account=os.getenv("SNOWFLAKE_ACCOUNT"),
        user=os.getenv("SNOWFLAKE_USER"),
        password=os.getenv("SNOWFLAKE_PASSWORD"),
        database=os.getenv("SNOWFLAKE_DATABASE", "CURSOR_ANALYTICS"),
        schema=os.getenv("SNOWFLAKE_SCHEMA", "MART"),
        warehouse=os.getenv("SNOWFLAKE_WAREHOUSE", "TRANSFORM_WH"),
    )


@st.cache_data(ttl=300)  # 5 minute cache
def query(sql: str, params: dict = None) -> pd.DataFrame:
    """Execute SQL query and return DataFrame. Cached for 5 minutes."""
    conn = get_connection()

    if DB_MODE == "snowflake":
        cursor = conn.cursor()
        cursor.execute(sql, params or {})
        columns = [desc[0] for desc in cursor.description]
        return pd.DataFrame(cursor.fetchall(), columns=columns)
    else:
        if params:
            # DuckDB parameterized query
            return conn.execute(sql, list(params.values())).df()
        return conn.execute(sql).df()


def refresh_data():
    """
    Trigger ETL pipeline (dev mode only).
    In production, data updates via scheduled Cloud Run Jobs.
    """
    if DB_MODE == "snowflake":
        st.warning("Refresh not available in production. Data updates every 15 minutes via scheduled jobs.")
        return False

    from pipeline.loader import DataLoader
    from pipeline.run_dbt import run_dbt_build

    cursor_sim_url = os.getenv("CURSOR_SIM_URL", "http://localhost:8080")

    with st.spinner("Extracting data from cursor-sim..."):
        loader = DataLoader(cursor_sim_url)
        loader.run("/data/raw")

    with st.spinner("Loading to DuckDB..."):
        from pipeline.duckdb_loader import load_parquet_to_duckdb
        load_parquet_to_duckdb("/data/raw", "/data/analytics.duckdb")

    with st.spinner("Running dbt transforms..."):
        run_dbt_build(target="dev")

    # Clear all caches
    st.cache_data.clear()
    st.cache_resource.clear()

    st.success("Data refreshed successfully!")
    return True
```

---

### 2. Main App (`app.py`)

**Purpose**: Streamlit entry point with sidebar and navigation.

```python
# app.py
import streamlit as st
from components.sidebar import render_sidebar

st.set_page_config(
    page_title="DOXAPI Analytics",
    page_icon="📊",
    layout="wide",
    initial_sidebar_state="expanded"
)

# Render shared sidebar
render_sidebar()

# Home page content
st.title("📊 DOXAPI Analytics Dashboard")

st.markdown("""
Welcome to the AI Code Analytics Dashboard. Use the sidebar to navigate:

- **Velocity**: PR cycle times and throughput metrics
- **AI Impact**: Compare metrics across AI usage bands
- **Quality**: Revert rates and code quality trends
- **Review Costs**: Code review burden analysis
""")

# Quick stats on home page
from db.connector import query

col1, col2, col3, col4 = st.columns(4)

with col1:
    total_prs = query("SELECT COUNT(*) FROM mart.velocity")[0][0]
    st.metric("Total PRs", f"{total_prs:,}")

with col2:
    avg_cycle = query("SELECT AVG(total_cycle_time) FROM mart.velocity")[0][0]
    st.metric("Avg Cycle Time", f"{avg_cycle:.1f} days")

with col3:
    revert_rate = query("SELECT AVG(revert_rate) FROM mart.quality")[0][0]
    st.metric("Avg Revert Rate", f"{revert_rate:.1%}")

with col4:
    avg_ai = query("SELECT AVG(avg_ai_ratio) FROM mart.velocity")[0][0]
    st.metric("Avg AI Ratio", f"{avg_ai:.0%}")
```

---

### 3. Sidebar Component (`components/sidebar.py`)

**Purpose**: Shared sidebar with filters and refresh button.

```python
# components/sidebar.py
import streamlit as st
import os
from db.connector import query, refresh_data, DB_MODE


def render_sidebar():
    """Render shared sidebar with filters and refresh button."""
    with st.sidebar:
        st.image("assets/logo.png", width=200)  # Optional logo
        st.title("DOXAPI Analytics")

        st.divider()

        # Filters
        st.subheader("🔧 Filters")

        # Repository filter
        repos = query("SELECT DISTINCT repo_name FROM mart.velocity ORDER BY repo_name")
        repo_options = ["All"] + repos["repo_name"].tolist()
        selected_repo = st.selectbox("Repository", repo_options)

        # Date range filter
        date_range = st.selectbox(
            "Date Range",
            ["Last 7 days", "Last 30 days", "Last 90 days", "All time"],
            index=2
        )

        # Store filter state in session
        st.session_state["filter_repo"] = selected_repo
        st.session_state["filter_date_range"] = date_range

        st.divider()

        # Refresh button (dev mode only)
        st.subheader("🔄 Data")

        if DB_MODE == "duckdb":
            if st.button("Refresh Data", use_container_width=True):
                refresh_data()
                st.rerun()
        else:
            st.info("Data updates every 15 minutes via scheduled jobs.")

        # Last updated timestamp
        last_updated = query("""
            SELECT MAX(week) as last_week FROM mart.velocity
        """)["last_week"][0]
        st.caption(f"Data through: {last_updated}")

        st.divider()

        # Environment indicator
        if DB_MODE == "snowflake":
            st.success("🟢 Production (Snowflake)")
        else:
            st.warning("🟡 Development (DuckDB)")
```

---

### 4. Velocity Page (`pages/1_velocity.py`)

**Purpose**: Display velocity metrics and cycle time trends.

```python
# pages/1_velocity.py
import streamlit as st
import plotly.express as px
from db.connector import query
from components.metrics import render_kpi_row
from queries.velocity import get_velocity_data, get_cycle_time_breakdown

st.set_page_config(page_title="Velocity", page_icon="🚀", layout="wide")

st.title("🚀 Velocity Metrics")

# Get filter state from sidebar
repo_filter = st.session_state.get("filter_repo", "All")
date_filter = st.session_state.get("filter_date_range", "Last 90 days")

# Build SQL filter clause
where_clause = ""
if repo_filter != "All":
    where_clause = f"WHERE repo_name = '{repo_filter}'"

# Get data
df = get_velocity_data(where_clause)

# KPI Row
render_kpi_row([
    {"label": "Total PRs", "value": df["total_prs"].sum(), "format": ","},
    {"label": "Avg Cycle Time", "value": df["total_cycle_time"].mean(), "format": ".1f", "suffix": " days"},
    {"label": "Active Devs", "value": df["active_developers"].max(), "format": ","},
    {"label": "Avg AI Ratio", "value": df["avg_ai_ratio"].mean(), "format": ".0%"},
])

st.divider()

# Cycle Time Trend Chart
st.subheader("Cycle Time Trend")
fig = px.line(
    df.sort_values("week"),
    x="week",
    y=["coding_lead_time", "pickup_time", "review_lead_time"],
    title="Weekly Cycle Time Components",
    labels={"value": "Days", "week": "Week"},
)
st.plotly_chart(fig, use_container_width=True)

# Cycle Time Breakdown
st.subheader("Cycle Time Breakdown")
breakdown = get_cycle_time_breakdown(where_clause)
fig2 = px.bar(
    breakdown,
    x="component",
    y="hours",
    title="Average Cycle Time by Component",
    color="component",
)
st.plotly_chart(fig2, use_container_width=True)

# P50/P90 Comparison
st.subheader("Cycle Time Percentiles")
col1, col2 = st.columns(2)
with col1:
    st.metric("P50 Cycle Time", f"{df['p50_cycle_time'].mean():.1f} days")
with col2:
    st.metric("P90 Cycle Time", f"{df['p90_cycle_time'].mean():.1f} days")
```

---

### 5. SQL Queries (`queries/velocity.py`)

**Purpose**: Parameterized SQL queries for velocity metrics.

```python
# queries/velocity.py
from db.connector import query


def get_velocity_data(where_clause: str = "") -> "pd.DataFrame":
    """Get weekly velocity metrics."""
    sql = f"""
    SELECT
        week,
        repo_name,
        active_developers,
        total_prs,
        avg_pr_size,
        coding_lead_time,
        pickup_time,
        review_lead_time,
        total_cycle_time,
        p50_cycle_time,
        p90_cycle_time,
        avg_ai_ratio
    FROM mart.velocity
    {where_clause}
    ORDER BY week DESC
    """
    return query(sql)


def get_cycle_time_breakdown(where_clause: str = "") -> "pd.DataFrame":
    """Get cycle time breakdown by component."""
    sql = f"""
    SELECT
        'Coding' as component,
        AVG(coding_lead_time) * 24 as hours
    FROM mart.velocity {where_clause}
    UNION ALL
    SELECT
        'Pickup' as component,
        AVG(pickup_time) * 24 as hours
    FROM mart.velocity {where_clause}
    UNION ALL
    SELECT
        'Review' as component,
        AVG(review_lead_time) * 24 as hours
    FROM mart.velocity {where_clause}
    """
    return query(sql)


def get_velocity_summary(where_clause: str = "") -> dict:
    """Get summary statistics."""
    sql = f"""
    SELECT
        SUM(total_prs) as total_prs,
        AVG(total_cycle_time) as avg_cycle_time,
        MAX(active_developers) as max_developers,
        AVG(avg_ai_ratio) as avg_ai_ratio
    FROM mart.velocity
    {where_clause}
    """
    result = query(sql)
    return result.iloc[0].to_dict()
```

---

## Docker Configuration

### Dockerfile

```dockerfile
# services/streamlit-dashboard/Dockerfile
FROM python:3.11-slim

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Install Python dependencies
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . ./

# Copy dbt project (for embedded pipeline in dev)
COPY ../../dbt/ ./dbt/

# Create data directory
RUN mkdir -p /data/raw

# Non-root user
RUN useradd -m -u 1000 streamlit && chown -R streamlit:streamlit /app /data
USER streamlit

# Expose port
EXPOSE 8501

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:8501/_stcore/health || exit 1

# Entrypoint
ENTRYPOINT ["streamlit", "run", "app.py", \
    "--server.port=8501", \
    "--server.address=0.0.0.0", \
    "--server.headless=true"]
```

### Requirements

```txt
# requirements.txt
streamlit>=1.30.0
pandas>=2.0.0
plotly>=5.18.0
duckdb>=0.9.0
snowflake-connector-python>=3.5.0

# For embedded pipeline (dev mode)
requests>=2.31.0
pyarrow>=14.0.0

# Testing
pytest>=7.4.0
pytest-cov>=4.1.0
```

---

## Environment Configuration

### Development (.env.dev)

```bash
DB_MODE=duckdb
DUCKDB_PATH=/data/analytics.duckdb
CURSOR_SIM_URL=http://cursor-sim:8080
```

### Production (.env.prod)

```bash
DB_MODE=snowflake
SNOWFLAKE_ACCOUNT=xxx.us-central1.gcp
SNOWFLAKE_USER=STREAMLIT_USER
SNOWFLAKE_PASSWORD=***  # From Secret Manager
SNOWFLAKE_DATABASE=CURSOR_ANALYTICS
SNOWFLAKE_SCHEMA=MART
SNOWFLAKE_WAREHOUSE=TRANSFORM_WH
```

---

## Caching Strategy

### Query Caching

```python
@st.cache_data(ttl=300)  # 5 minute TTL
def query(sql: str, params: dict = None) -> pd.DataFrame:
    """Cached query execution."""
    ...
```

### Connection Caching

```python
@st.cache_resource  # Persistent across reruns
def get_connection():
    """Cached database connection."""
    ...
```

### Cache Invalidation

```python
def refresh_data():
    """Clear all caches after data refresh."""
    st.cache_data.clear()
    st.cache_resource.clear()
```

---

## Performance Optimization

| Optimization | Implementation |
|--------------|----------------|
| Query caching | `@st.cache_data(ttl=300)` |
| Connection pooling | `@st.cache_resource` |
| Lazy loading | Load data only when page accessed |
| Minimal queries | Pre-aggregated mart tables |
| Column selection | Only SELECT needed columns |

---

## Security

| Concern | Mitigation |
|---------|------------|
| Credentials | Environment variables, Secret Manager |
| SQL injection | Parameterized queries |
| Access control | Cloud Run IAM |
| Data exposure | Read-only database user |

---

## Testing Strategy

### Unit Tests

```python
# tests/test_connector.py
def test_get_connection_duckdb(monkeypatch):
    monkeypatch.setenv("DB_MODE", "duckdb")
    conn = get_connection()
    assert conn is not None

def test_query_caching():
    result1 = query("SELECT 1")
    result2 = query("SELECT 1")
    # Should return cached result
    assert result1.equals(result2)
```

### Integration Tests

```python
# tests/test_pages.py
def test_velocity_page_loads():
    """Verify velocity page renders without error."""
    from pages.velocity import main
    # Should not raise
    main()
```

---

## Deployment

### Local Development

```bash
# Start with docker-compose
docker-compose up streamlit-dashboard

# Or run directly
cd services/streamlit-dashboard
streamlit run app.py
```

### Production (Cloud Run)

```bash
# Build and push
gcloud builds submit --tag gcr.io/PROJECT/streamlit-dashboard

# Deploy
gcloud run deploy streamlit-dashboard \
    --image gcr.io/PROJECT/streamlit-dashboard \
    --port 8501 \
    --set-env-vars DB_MODE=snowflake \
    --set-secrets SNOWFLAKE_PASSWORD=snowflake-password:latest \
    --allow-unauthenticated
```

---

**Next**: See `task.md` for implementation breakdown.
