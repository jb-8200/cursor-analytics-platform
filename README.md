# cursor-analytics-platform

A microservices platform for simulating, aggregating, and visualising AI coding-assistant usage metrics. Built to support SDLC research on the real impact of tools like Cursor — measuring velocity, review cost, and code quality outcomes across correlated datasets.

## Services

| Service | Description |
|---|---|
| `cursor-sim` | High-fidelity API simulator — exact Cursor (29 endpoints) + GitHub PR lifecycle (20 endpoints) |
| `cursor-analytics-core` | GraphQL aggregator with PostgreSQL persistence |
| `cursor-viz-spa` | React dashboard for interactive metric exploration |
| `streamlit-dashboard` | Streamlit data app for research queries |

## Data pipeline

```
cursor-sim → api-loader (Python) → Parquet landing → dbt transforms → DuckDB (dev) / Snowflake (prod) → Dashboard
```

- **Dev**: DuckDB with the `main_` schema prefix (`main_raw`, `main_staging`, `main_mart`)
- **Prod**: Snowflake via SnapLogic extraction; same dbt SQL, Snowflake dialect
- **Data generation**: Seed-based deterministic generation via DataDesigner (NVIDIA) — reproducible across environments

## Research questions addressed

- **Velocity**: Does AI assistance accelerate coding and review cycles?
- **Review cost**: Does AI-generated code require more review iterations?
- **Quality**: Do AI-assisted commits have higher revert or survival rates?

## Requirements

- Docker + Docker Compose (runs all services locally)
- Node.js 20+ (cursor-viz-spa)
- Python 3.11+ (api-loader, dbt, streamlit)
- dbt-duckdb (dev) / dbt-snowflake (prod)

## Quick start

```bash
# Start the simulator and core services
docker compose up cursor-sim cursor-analytics-core

# Load a seed dataset and run transforms
python tools/api-loader/main.py --seed data/seeds/default.json
cd dbt && dbt run

# Launch the dashboard
streamlit run streamlit-dashboard/app.py
```

## Documentation reading order

1. `docs/DESIGN.md` — overall architecture and design decisions
2. `services/cursor-sim/SPEC.md` — API contract and data formats
3. `docs/design/new_data_architecture.md` — ETL pipeline details
4. `docs/TESTING_STRATEGY.md` — data contract testing approach
5. `services/{service}/README.md` — per-service setup

## Contact

[jb@jishutech.io](mailto:jb@jishutech.io) · [jishutech.io](https://jishutech.io)
