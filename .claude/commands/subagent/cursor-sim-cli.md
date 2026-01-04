---
description: Spawn cursor-sim-cli-dev subagent with CLI scope constraints for P4 feature implementation
argument-hint: "[feature-id] [task-id]"
allowed-tools: Task
---

# Spawn cursor-sim-cli-dev Subagent

Delegate implementation task to cursor-sim-cli-dev subagent with CLI-only scope constraints.

**Feature**: $1
**Task**: $2

## Objective

Implement the specified task following Test-Driven Development within CLI scope only.

## Scope Constraints

**See** `.claude/rules/service/cursor-sim.md` for always-on enforcement.

**Additional for this task**:
- Work ONLY on CLI code: `services/cursor-sim/internal/cli/` and `cmd/simulator/`
- Follow `go-best-practices` skill for Go patterns
- Update `.work-items/$1/task.md` with progress on completion

## Context Files

Provide to subagent:
- `.work-items/$1/user-story.md` - Requirements
- `.work-items/$1/design.md` - Technical approach
- `.work-items/$1/task.md` - Task details
- `services/cursor-sim/SPEC.md` - Service specification

## SDD Workflow

Subagent follows:

1. **SPEC**: Read requirements
2. **TEST**: Write failing tests (RED)
3. **CODE**: Minimal implementation (GREEN)
4. **REFACTOR**: Clean up while tests pass
5. **REFLECT**: Check `dependency-reflection`
6. **SYNC**: Run `spec-sync-check` for SPEC.md updates
7. **COMMIT**: Include all changes with descriptive message

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

- CLI scope ONLY (protects API contracts for P5/P6)
- Never modify `internal/api/` or `internal/generator/`
- All tests must pass before commit
- Update task.md with results

---

See also:
- `.claude/rules/service/cursor-sim.md` - API contract protection
- `.claude/skills/go-best-practices/` - Go coding patterns
- `services/cursor-sim/SPEC.md` - Complete specification
