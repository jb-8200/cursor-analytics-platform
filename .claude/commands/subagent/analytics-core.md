---
description: Spawn analytics-core-dev subagent with P5 GraphQL/TypeScript scope constraints
argument-hint: "[feature-id] [task-id]"
allowed-tools: Task
---

# Spawn analytics-core-dev Subagent

Delegate implementation task to analytics-core-dev subagent with P5 service scope constraints.

**Feature**: $1
**Task**: $2

## Objective

Implement the specified task following Test-Driven Development within P5 GraphQL service scope.

## Scope Constraints

**See** `.claude/rules/service/analytics-core.md` for always-on enforcement.

**Additional for this task**:
- Work ONLY on P5 code: `services/cursor-analytics-core/`
- Follow `typescript-graphql-patterns` skill for TypeScript/GraphQL patterns
- Consume cursor-sim API via `api-contract` skill
- Update `.work-items/$1/task.md` with progress on completion

## Context Files

Provide to subagent:
- `.work-items/$1/user-story.md` - Requirements
- `.work-items/$1/design.md` - Technical approach
- `.work-items/$1/task.md` - Task details
- `services/cursor-analytics-core/SPEC.md` - Service specification
- `services/cursor-sim/SPEC.md` - Upstream API contract

## Tech Stack

- Node.js 18+, TypeScript 5.x
- Apollo Server 4, GraphQL
- PostgreSQL + Prisma ORM
- Jest (target 80%+ coverage)

## SDD Workflow

Subagent follows:

1. **SPEC**: Read requirements and SPEC.md
2. **TEST**: Write failing tests (RED)
3. **CODE**: Minimal implementation (GREEN)
4. **REFACTOR**: Clean up while tests pass
5. **REFLECT**: Check `dependency-reflection`
6. **SYNC**: Run `spec-sync-check` for SPEC.md updates
7. **COMMIT**: Include all changes with descriptive message

## GraphQL Guidelines

- Keep schema backward compatible
- Document breaking changes for P6 consumers
- Use proper GraphQL types (not just String)
- Include field descriptions
- Use cursor-based pagination for lists

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

GraphQL Changes:
- {schema/resolver changes}

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

- P5 consumes cursor-sim REST API (check `api-contract`)
- P5 exposes GraphQL for P6 (document schema changes)
- Use Prisma for all database access
- Target 80%+ test coverage
- Update task.md with results

---

See also:
- `.claude/rules/service/analytics-core.md` - Schema and database constraints
- `.claude/skills/typescript-graphql-patterns/` - GraphQL patterns
- `.claude/skills/api-contract/` - cursor-sim API reference
- `services/cursor-analytics-core/SPEC.md` - Complete specification
