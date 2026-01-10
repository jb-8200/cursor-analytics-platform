# Design: P9-F02 Streamlit Dashboard Hardening

## Goal
Harden the Streamlit Dashboard service by fixing security vulnerabilities, resolving infrastructure dependencies, and improving configuration management.

## Architecture

### 1. Data Access Layer
**Pattern**: Parameterized Queries
**Reason**: To prevent SQL injection.
**Change**:
- Refactor `db/connector.py` to accept `params` dictionary.
- Update `duckdb.execute(sql, list(params.values()))` to more robust binding if possible, or strictly validate order.
- Rewrite `queries/*.py` to pass parameters instead of injecting them.

### 2. Infrastructure
**Pattern**: Self-Contained Docker Image
**Reason**: The container must run `dbt` for data refresh.
**Change**:
- `Dockerfile`: Install `dbt-duckdb`.
- **Build Context**: The `dbt` project resides in `../../dbt`.
  - *Option A*: Change context to root.
  - *Option B (Chosen)*: Assume `docker-compose` handles the volume mount for `dbt` in development (as is standard for this architecture), but ensure the *Production* build (Dockerfile) supports the necessary dependencies. For production deployment, we might need a multi-stage build or copy strategy, but for this "hardening" phase, getting Dev mode working via `requirements.txt` is the priority.

### 3. Configuration
**Pattern**: Environment Variables with Defaults
**Reason**: Support both Container and Host execution.
**Change**:
- `DUCKDB_PATH`: `os.getenv("DUCKDB_PATH", "data/analytics.duckdb")` (Relative!)
- `DBT_PROJECT_DIR`: Configurable path to dbt project.

## Implementation Steps

### Phase 1: Security (Code)
1.  **Refactor Connector**: Update `db/connector.py` to support `query(sql, params)`.
2.  **Secure Queries**: Rewrite `velocity.py`, `ai_impact.py`, `quality.py`, `review_costs.py`, `sidebar.py`.

### Phase 2: Platform (Infra)
1.  **Dependencies**: Add `dbt-duckdb` to `services/streamlit-dashboard/requirements.txt`.
2.  **Docker**: Update `Dockerfile` to install dependencies and ensure permissions.
3.  **Paths**: Search and replace absolute `/data/` paths with relative or env-var driven paths.

## Verification
- **Test**: Run `pytest` on query modules (if tests exist) or manual smoke test.
- **Manual**: Verify "Refresh Data" button in UI.
