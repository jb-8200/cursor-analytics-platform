---
description: Spawn data-tier-dev subagent with P8 Python/dbt scope constraints
argument-hint: "[feature-id] [task-id]"
allowed-tools: Task
---

# Spawn data-tier-dev Subagent

Delegate implementation task to data-tier-dev subagent with P8 data tier scope constraints.

**Feature**: $1
**Task**: $2

## Objective

Implement the specified task following Test-Driven Development within P8 data tier scope.

## Scope Constraints

**Work ONLY on**:
- `tools/api-loader/` - Python extraction scripts
- `dbt/` - dbt models, macros, tests

**NEVER touch**:
- `services/cursor-sim/` - API is source of truth
- `services/streamlit-dashboard/` - P9 scope
- `services/cursor-analytics-core/` - Deprecated

## Context Files

Provide to subagent:
- `.work-items/$1/user-story.md` - Requirements
- `.work-items/$1/design.md` - Technical approach
- `.work-items/$1/task.md` - Task details
- `services/cursor-sim/SPEC.md` - Upstream API contract

## Tech Stack

- Python 3.11+, pytest
- dbt-core 1.7+, dbt-duckdb, dbt-snowflake
- DuckDB, Snowflake
- pandas, pyarrow

## API Contract CRITICAL

cursor-sim returns RAW ARRAYS, not wrapper objects:

```python
# CORRECT
resp = requests.get(f"{base_url}/repos")
repos = resp.json()  # Returns list directly

# WRONG - will fail
repos = resp.json()["repositories"]  # No wrapper!
```

## SDD Workflow

Subagent follows:

1. **SPEC**: Read requirements and cursor-sim SPEC.md
2. **TEST**: Write failing tests (RED)
3. **CODE**: Minimal implementation (GREEN)
4. **REFACTOR**: Clean up while tests pass
5. **REFLECT**: Check `dependency-reflection`
6. **SYNC**: Run `spec-sync-check` for SPEC.md updates
7. **COMMIT**: Include all changes with descriptive message

## dbt Guidelines

- Use cross-engine macros for DuckDB/Snowflake parity
- Write tests for all models
- Calculate cycle times in staging (not from API)
- Rename `was_reverted` to `is_reverted` in staging

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

dbt Changes:
- Models: {list}
- Tests: {list}

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

- P8 consumes cursor-sim REST API (RAW ARRAYS!)
- P8 produces mart tables for P9 Streamlit
- Run pytest before committing
- Run `dbt test` for model validation
- Update task.md with results

---

See also:
- `.claude/agents/data-tier-dev.md` - Full agent definition
- `.claude/skills/api-contract/` - cursor-sim API reference
- `.work-items/P8-F01-data-tier/design.md` - Architecture details
