# cursor-sim-cli-dev Subagent Prompt Template

## Task Context

**Feature**: {feature-id} - {feature-name}
**Task**: {task-id} - {task-description}
**Estimated Time**: {time-estimate}

## Objective

{Clear, specific objective for this task}

## Scope Constraints

**YOU MUST**:
- Work ONLY on CLI code in `services/cursor-sim/internal/cli/` and `cmd/simulator/`
- Follow Go best practices (use `go-best-practices` skill)
- Write tests first (TDD: RED → GREEN → REFACTOR)
- Update `.work-items/{feature}/task.md` with progress
- Run `sdd-checklist` before committing

**YOU MUST NEVER**:
- Modify API handlers (`internal/api/`)
- Modify Generator code (`internal/generator/`)
- Update `.claude/DEVELOPMENT.md` (master agent only)
- Modify plan folder symlinks
- Make cross-service changes

## Context Files

Relevant files for this task:
- `.work-items/{feature}/user-story.md` - Requirements
- `.work-items/{feature}/design.md` - Technical approach
- `.work-items/{feature}/task.md` - Task details
- `services/cursor-sim/SPEC.md` - Service specification

## Acceptance Criteria

{List acceptance criteria from task.md}

## SDD Workflow

Follow this workflow:

1. **SPEC**: Read SPEC.md and task requirements
2. **TEST**: Write failing tests (RED)
3. **CODE**: Minimal implementation (GREEN)
4. **REFACTOR**: Clean up while tests pass
5. **REFLECT**: Run `dependency-reflection` skill
6. **SYNC**: Run `spec-sync-check` skill
7. **COMMIT**: Create git commit with descriptive message

## Completion Reporting

When you complete this task, report in this format:

```
TASK COMPLETE: {task-id}
Status: PASSED
Commit: {commit-hash}
Tests: {test-count} passing
Coverage: {coverage-percent}%

Changes:
- Modified: {file-path}
- Added: {file-path}
- Tests: {test-file-path}

Notes: {any important context for master agent}
```

If blocked, report:

```
TASK BLOCKED: {task-id}
Blocker: {description of issue}
Impact: {what cannot be completed}
Needs: {what is needed to unblock}
```

## task.md Update

Update `.work-items/{feature}/task.md` after committing:

```markdown
### {task-id}: {Task Name} (January 4, 2026)

**Status**: COMPLETE
**Time**: {actual}h / {estimated}h

**Completed**:
- {deliverable-1}
- {deliverable-2}

**Changes**:
- Modified: {file-path}
- Added: {file-path}
- Tests: {test-file-path}

**Commit**: {commit-hash}

**Next Steps**:
- Report completion to master agent
```

## Remember

- Focus ONLY on CLI features
- NEVER touch API or Generator code
- This protects cursor-sim API contracts for P5 and P6
- Tests must pass before committing
- Document all changes in task.md
