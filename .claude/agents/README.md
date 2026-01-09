# Custom Subagents

Specialized agents for parallel development with isolated scope constraints.

---

## Available Agents

| Agent | Service | Phase | Model | Scope | Key Skills |
|-------|---------|-------|-------|-------|-----------|
| `planning-dev` | Any | Any | **Opus** | Research, design, task breakdown | sdd-checklist |
| `data-tier-dev` | data-tier | P8 | **Sonnet** | ETL, dbt, DuckDB/Snowflake | api-contract, sdd-checklist |
| `streamlit-dev` | streamlit-dashboard | P9 | **Sonnet** | Dashboard, Plotly | sdd-checklist |
| `quick-fix` | Any | - | **Haiku** | Small fixes only | sdd-checklist |
| `cursor-sim-cli-dev` | cursor-sim | P4 | Sonnet | CLI only | go-best-practices, api-contract |
| `cursor-sim-infra-dev` | cursor-sim | P7 | Sonnet | Docker, Cloud Run | - |
| `analytics-core-dev` | analytics-core | P5 | Sonnet | GraphQL service | typescript-graphql-patterns |
| `viz-spa-dev` | viz-spa | P6 | Sonnet | React dashboard | react-vite-patterns |

### Model Selection

| Model | Use Case | Speed | Cost |
|-------|----------|-------|------|
| **Haiku** | Quick fixes, simple tasks | Fast | Low |
| **Sonnet** | Feature implementation, complex tasks | Medium | Medium |
| **Opus** | Architecture, orchestration, review | Slower | Higher |

---

## Planning Agent (Opus)

### planning-dev

Research and design specialist using Opus model.

**Purpose**:
- Research external APIs and documentation
- Create work item files (user-story.md, design.md, task.md)
- Break down features into tasks with subagent assignments
- Design API contracts and data models

**Key Responsibility**:
- Every task created must include a **recommended subagent**
- Follow SDD methodology for all work items
- Fetch and analyze external documentation before designing

**Output Format**:
```markdown
#### TASK-XX-##: Task Name (Est: X.Xh)
**Assigned Subagent**: `{subagent-name}`
**Goal**: One sentence describing the deliverable
```

**NEVER**:
- Write implementation code (design only)
- Skip researching external API documentation
- Create tasks without subagent assignments

---

## New Agents (P8/P9)

### data-tier-dev (P8)

Python/dbt specialist for the data tier.

**ONLY work on**:
- `tools/api-loader/` - Python extraction scripts
- `dbt/` - dbt models, macros, tests

**NEVER touch**:
- `services/cursor-sim/` - API is source of truth
- Other services

**Key Responsibility**:
- Extract data from cursor-sim REST API
- Load to DuckDB (dev) / Snowflake (prod)
- Transform via dbt into analytics marts

**Critical**: cursor-sim returns RAW ARRAYS, not wrapper objects!

### streamlit-dev (P9)

Streamlit/Plotly specialist for the dashboard.

**ONLY work on**:
- `services/streamlit-dashboard/`

**Depends on**:
- P8 dbt mart tables must exist

**Key Responsibility**:
- Database connector abstraction (DuckDB/Snowflake)
- Dashboard pages with Plotly visualizations
- Embedded refresh pipeline (dev mode)

### quick-fix (Haiku)

Fast agent for small, independent tasks.

**Ideal for**:
- Typos, simple bug fixes
- Config changes
- Missing imports
- Documentation updates

**NOT for**:
- New features
- Multi-file refactoring
- Tasks requiring tests

---

## Legacy Agents (P4-P7)

### cursor-sim-cli-dev
**ONLY work on**:
- `services/cursor-sim/internal/cli/`
- `services/cursor-sim/cmd/simulator/`

**NEVER touch**:
- `internal/api/` (protects API contracts)
- `internal/generator/` (protects data generation)

### analytics-core-dev (DEPRECATED)
**Status**: P5 deprecated in favor of P8 data tier

**ONLY work on**:
- `services/cursor-analytics-core/`

### viz-spa-dev
**ONLY work on**:
- `services/cursor-viz-spa/`

**Must align with**:
- P8 dbt mart tables (formerly P5 GraphQL)

### cursor-sim-infra-dev
**Scope**:
- Docker containerization
- GCP Cloud Run deployment
- Infrastructure configuration

---

## Orchestration Model

For subagent coordination, see `.work-items/P0-F01-sdd-subagent-orchestration/design.md`:

1. **Master Agent** (Opus - Chief Architect):
   - Delegates tasks to appropriate subagents
   - Reviews cross-service code quality
   - Handles E2E testing
   - Updates DEVELOPMENT.md

2. **Subagents** (Sonnet - Service-Specialized):
   - Work within assigned scope only
   - Update `.work-items/{feature}/task.md` with progress
   - Report completion to master agent
   - NEVER update DEVELOPMENT.md or plan folder symlinks

3. **Quick-fix** (Haiku - Fast Fixes):
   - Small, independent tasks only
   - Escalate if task is complex
   - Minimal overhead

4. **Task.md Update Format**:
   ```markdown
   ### TASK##: {Task Name}
   **Status**: COMPLETE
   **Time**: {actual}h / {estimated}h
   **Commit**: {hash}
   **Changes**: [file list]
   ```

---

## Spawning Subagents

Use slash commands in `.claude/commands/subagent/`:

```bash
# P8 Data Tier
/subagent/data-tier P8-F01-data-tier TASK-P8-01

# P9 Streamlit Dashboard
/subagent/streamlit P9-F01-streamlit-dashboard TASK-P9-01

# Quick Fix (Haiku)
/subagent/quick-fix "Fix typo in README.md"

# Legacy
/subagent/cursor-sim-cli P4-F02 TASK07
/subagent/analytics-core P5-F01 TASK01
/subagent/viz-spa P6-F01 TASK01
```

---

## Background Execution

Run subagents in background for parallel development:

```python
# Master agent spawns multiple subagents
Task(
    subagent_type="data-tier-dev",
    model="sonnet",
    run_in_background=True,
    prompt="Implement TASK-P8-01..."
)

Task(
    subagent_type="streamlit-dev",
    model="sonnet",
    run_in_background=True,
    prompt="Implement TASK-P9-01..."  # Can run in parallel for infra tasks
)
```

---

## SDD Workflow for Subagents

All subagents follow Spec-Driven Development:

1. **SPEC**: Read SPEC.md and task requirements
2. **TEST**: Write failing tests (RED)
3. **CODE**: Minimal implementation (GREEN)
4. **REFACTOR**: Clean up while tests pass
5. **REFLECT**: Run `dependency-reflection` skill
6. **SYNC**: Run `spec-sync-check` skill
7. **COMMIT**: Create commit with descriptive message

---

## Completion Reporting

When subagent completes task:

```
TASK COMPLETE: {task-id}
Status: PASSED
Commit: {commit-hash}
Tests: {count} passing
Coverage: {percent}%

Changes:
- {file-path}

Notes: {any blockers or follow-up}
```

If blocked:

```
TASK BLOCKED: {task-id}
Blocker: {issue description}
Impact: {what cannot be completed}
Needs: {what is needed to unblock}
```

Quick-fix format (shorter):

```
FIXED: {what was fixed}
File: {path}
Commit: {hash}
```

---

## Key Rules

### All Subagents
- Follow the 7-step SDD workflow
- Run `sdd-checklist` before committing
- Write tests first (TDD)
- Target 80%+ test coverage
- Update task.md with progress
- Report completion to master agent

### Never:
- Update `.claude/DEVELOPMENT.md` (master agent only)
- Modify plan folder symlinks
- Make cross-service changes without coordination
- Skip tests before committing

### Escalation
- If task is too complex for quick-fix, escalate
- If blocked by dependency, report to master agent
- If scope violation needed, ask master agent first

### CLI/Shell Command Delegation

Non-CLI agents (data-tier-dev, streamlit-dev, viz-spa-dev, analytics-core-dev) **cannot run shell commands in background mode**.

When CLI actions are needed:
1. Agent escalates to orchestrator with command details
2. Orchestrator delegates to appropriate CLI agent:
   - `quick-fix` (Haiku) - Simple tasks: pip install, pytest, basic commands
   - `cursor-sim-infra-dev` (Sonnet) - Complex tasks: Docker, Cloud Run, infrastructure

This maintains context integrity while allowing necessary CLI operations.

---

## Dependency Graph

```
P4 (cursor-sim) ──► P8 (data-tier) ──► P9 (streamlit)
     API                ETL/dbt            Dashboard
   [source]           [transform]         [visualize]

P4 also ──► P5 (analytics-core) ──► P6 (viz-spa)
            [DEPRECATED]            [keep for now]
```

---

## See Also

- **SDD Protocol**: `.work-items/P0-F01-sdd-subagent-orchestration/design.md`
- **Rules**: `.claude/rules/` (enforcement constraints)
- **Skills**: `.claude/skills/` (guidance and patterns)
- **Specifications**: `services/{service}/SPEC.md`
- **Work Items**: `.work-items/P8-F01-data-tier/`, `.work-items/P9-F01-streamlit-dashboard/`

---

**Last Updated**: January 9, 2026
