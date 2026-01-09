---
description: Spawn streamlit-dev subagent with P9 Streamlit dashboard scope constraints
argument-hint: "[feature-id] [task-id]"
allowed-tools: Task
---

# Spawn streamlit-dev Subagent

Delegate implementation task to streamlit-dev subagent with P9 dashboard scope constraints.

**Feature**: $1
**Task**: $2

## Objective

Implement the specified task following Test-Driven Development within P9 Streamlit dashboard scope.

## Scope Constraints

**Work ONLY on**:
- `services/streamlit-dashboard/` - Dashboard application

**NEVER touch**:
- `services/cursor-sim/` - API is source of truth
- `tools/api-loader/` - P8 scope
- `dbt/` - P8 scope

**Depends on**:
- P8 dbt mart tables must exist for queries to work

## Context Files

Provide to subagent:
- `.work-items/$1/user-story.md` - Requirements
- `.work-items/$1/design.md` - Technical approach
- `.work-items/$1/task.md` - Task details
- `.work-items/P8-F01-data-tier/design.md` - Upstream mart schemas

## Tech Stack

- Python 3.11+, pytest
- Streamlit 1.30+
- Plotly, pandas
- DuckDB (dev), Snowflake (prod)

## Key Patterns

### Database Connector

```python
DB_MODE = os.getenv("DB_MODE", "duckdb")

@st.cache_resource
def get_connection():
    if DB_MODE == "snowflake":
        return _get_snowflake_connection()
    return _get_duckdb_connection()

@st.cache_data(ttl=300)
def query(sql: str) -> pd.DataFrame:
    conn = get_connection()
    return conn.execute(sql).df()
```

### Session State for Filters

```python
st.session_state["filter_repo"] = selected_repo
repo = st.session_state.get("filter_repo", "All")
```

## SDD Workflow

Subagent follows:

1. **SPEC**: Read requirements and design.md
2. **TEST**: Write failing tests (RED)
3. **CODE**: Minimal implementation (GREEN)
4. **REFACTOR**: Clean up while tests pass
5. **REFLECT**: Check `dependency-reflection`
6. **SYNC**: Run `spec-sync-check` for SPEC.md updates
7. **COMMIT**: Include all changes with descriptive message

## Streamlit Guidelines

- Use `@st.cache_data` for query results (5 min TTL)
- Use `@st.cache_resource` for connections
- Render sidebar on all pages
- Support wide layout
- Use Plotly for charts

## Completion Report

Subagent reports completion as:

```
TASK COMPLETE: $2
Status: PASSED
Commit: {hash}
Tests: {count} passing
Coverage: {percent}%

Changes:
- {file list}

Pages:
- {page changes}

Notes: {context for master agent}
```

If blocked:

```
TASK BLOCKED: $2
Blocker: {issue description}
Impact: {what cannot be completed}
Needs: {what is needed to unblock}
```

## Remember

- P9 consumes P8 dbt mart tables
- P9 supports DuckDB (dev) and Snowflake (prod)
- Refresh button only in dev mode
- Run pytest before committing
- Update task.md with results

---

See also:
- `.claude/agents/streamlit-dev.md` - Full agent definition
- `.work-items/P9-F01-streamlit-dashboard/design.md` - Architecture details
- `.work-items/P8-F01-data-tier/design.md` - Upstream data models
