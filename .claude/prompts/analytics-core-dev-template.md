# analytics-core-dev Subagent Prompt Template

## Task Context

**Feature**: {feature-id} - {feature-name}
**Task**: {task-id} - {task-description}
**Estimated Time**: {time-estimate}

## Objective

{Clear, specific objective for this task}

## Scope Constraints

**YOU MUST**:
- Work ONLY on P5 code in `services/cursor-analytics-core/`
- Follow TypeScript/GraphQL best practices (use `typescript-graphql-patterns` skill)
- Consume cursor-sim API via `api-contract` skill
- Write tests first (TDD: RED → GREEN → REFACTOR)
- Update `.work-items/{feature}/task.md` with progress
- Run `sdd-checklist` before committing

**YOU MUST NEVER**:
- Modify cursor-sim code (P4)
- Modify viz-spa code (P6)
- Update `.claude/DEVELOPMENT.md` (master agent only)
- Modify plan folder symlinks
- Make cross-service changes

## Context Files

Relevant files for this task:
- `.work-items/{feature}/user-story.md` - Requirements
- `.work-items/{feature}/design.md` - Technical approach
- `.work-items/{feature}/task.md` - Task details
- `services/cursor-analytics-core/SPEC.md` - Service specification (if exists)
- `services/cursor-sim/SPEC.md` - Upstream API contract

## Tech Stack

- **Runtime**: Node.js 18+, TypeScript 5.x
- **GraphQL**: Apollo Server 4
- **Database**: PostgreSQL + Prisma ORM
- **Testing**: Jest with coverage target 80%+

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

## GraphQL Schema Guidelines

- All schema changes MUST be backward compatible
- Document breaking changes clearly for P6
- Use proper GraphQL types (avoid String for everything)
- Include field descriptions
- Use cursor-based pagination for lists

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
- Schema: {schema-changes}
- Tests: {test-file-path}

GraphQL Changes:
- {list any schema or resolver changes}

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
- Schema: {schema-changes}
- Tests: {test-file-path}

**Commit**: {commit-hash}

**Next Steps**:
- Report completion to master agent
```

## Remember

- P5 consumes cursor-sim REST API (check `api-contract` skill)
- P5 exposes GraphQL for P6 (document schema changes)
- Use Prisma for database access
- Target 80%+ test coverage
- Document all changes in task.md
