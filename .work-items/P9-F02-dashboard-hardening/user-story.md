# User Story: P9-F02 Streamlit Dashboard Hardening

## Title
As a Platform Engineer, I want the Streamlit Dashboard to be secure, reliable, and deployable so that I can safely expose it to users and ensure data freshness.

## Context
The initial implementation of the Streamlit dashboard (P9-F01) successfully established the UI skeleton but suffers from critical issues hindering production use:
1.  **Security**: SQL injection vulnerabilities in filter logic.
2.  **Stability**: Data pipeline is broken due to missing `dbt` dependencies in the container.
3.  **Portability**: Hardcoded paths prevent local execution outside specific Docker contexts.

## Acceptance Criteria
- [ ] **No SQL Injection**: All database queries must use parameterized binding (no f-strings for values).
- [ ] **Working Pipeline**: The "Refresh Data" button must successfully run the extraction and dbt transformation pipeline within the Docker container.
- [ ] **Portability**: The service must run locally (`streamlit run app.py`) without requiring root-level `/data` directories.
- [ ] **Dependency Management**: `dbt-duckdb` must be installed and functional in the container.

## Reference
- Design: `docs/design/new_data_architecture.md` (Visualize Layer)
- Previous Review: `dashboard_review.md` (Artifact from analysis)
