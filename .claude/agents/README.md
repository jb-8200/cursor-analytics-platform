# Custom Subagents

Specialized agents for parallel development with isolated scope constraints.

## Available Agents

| Agent | Service | Phase | Scope | Key Skills |
|-------|---------|-------|-------|-----------|
| `cursor-sim-cli-dev` | cursor-sim | P4 | CLI only | go-best-practices, api-contract |
| `cursor-sim-infra-dev` | cursor-sim | P7 | Docker, Cloud Run | - |
| `analytics-core-dev` | analytics-core | P5 | GraphQL service | typescript-graphql-patterns, api-contract |
| `viz-spa-dev` | viz-spa | P6 | React dashboard | react-vite-patterns |

---

## Scope Constraints

### cursor-sim-cli-dev
**ONLY work on**:
- `services/cursor-sim/internal/cli/`
- `services/cursor-sim/cmd/simulator/`

**NEVER touch**:
- `internal/api/` (protects API contracts for P5/P6)
- `internal/generator/` (protects data generation)

### analytics-core-dev
**ONLY work on**:
- `services/cursor-analytics-core/`

**Must align with**:
- cursor-sim API contract (`api-contract` skill)
- GraphQL schema for P6 consumers

### viz-spa-dev
**ONLY work on**:
- `services/cursor-viz-spa/`

**Must align with**:
- analytics-core GraphQL schema (P5)
- Verify schema with `react-vite-patterns` skill

### cursor-sim-infra-dev
**Scope**:
- Docker containerization
- GCP Cloud Run deployment
- Infrastructure configuration

---

## Orchestration Model

For subagent coordination, see `.work-items/P0-F01-sdd-subagent-orchestration/design.md`:

1. **Master Agent** (Chief Architect):
   - Delegates tasks to appropriate subagents
   - Reviews cross-service code quality
   - Handles E2E testing
   - Updates DEVELOPMENT.md

2. **Subagents** (Service-Specialized):
   - Work within assigned scope only
   - Update `.work-items/{feature}/task.md` with progress
   - Report completion to master agent
   - NEVER update DEVELOPMENT.md or plan folder symlinks

3. **Task.md Update Format**:
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
/subagent/cursor-sim-cli [feature-id] [task-id]
/subagent/analytics-core [feature-id] [task-id]
/subagent/viz-spa [feature-id] [task-id]
```

**Example**:
```bash
/subagent/cursor-sim-cli P4-F02 TASK07
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

---

## Key Rules

### All Subagents
- ✅ Follow the 7-step SDD workflow
- ✅ Run `sdd-checklist` before committing
- ✅ Write tests first (TDD)
- ✅ Target 80%+ test coverage
- ✅ Update task.md with progress
- ✅ Report completion to master agent

### Never:
- ❌ Update `.claude/DEVELOPMENT.md` (master agent only)
- ❌ Modify plan folder symlinks
- ❌ Make cross-service changes without coordination
- ❌ Skip tests before committing

---

## See Also

- **SDD Protocol**: `.work-items/P0-F01-sdd-subagent-orchestration/design.md`
- **Rules**: `.claude/rules/` (enforcement constraints)
- **Skills**: `.claude/skills/` (guidance and patterns)
- **Specifications**: `services/{service}/SPEC.md`

---

**Last Updated**: January 4, 2026
