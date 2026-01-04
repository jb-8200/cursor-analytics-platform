---
description: Spawn viz-spa-dev subagent with P6 React/Vite scope constraints
argument-hint: "[feature-id] [task-id]"
allowed-tools: Task
---

# Spawn viz-spa-dev Subagent

Delegate implementation task to viz-spa-dev subagent with P6 service scope constraints.

**Feature**: $1
**Task**: $2

## Objective

Implement the specified task following Test-Driven Development within P6 React/Vite service scope.

## Scope Constraints

**See** `.claude/rules/service/viz-spa.md` for always-on enforcement.

**Additional for this task**:
- Work ONLY on P6 code: `services/cursor-viz-spa/`
- Follow `react-vite-patterns` skill for React/Vite patterns
- Verify schema alignment with `services/cursor-analytics-core/src/graphql/schema.graphql`
- Update `.work-items/$1/task.md` with progress on completion

## Context Files

Provide to subagent:
- `.work-items/$1/user-story.md` - Requirements
- `.work-items/$1/design.md` - Technical approach
- `.work-items/$1/task.md` - Task details
- `services/cursor-analytics-core/src/graphql/schema.graphql` - Upstream GraphQL schema

## Tech Stack

- React 18+, TypeScript 5.x
- Vite 5.x, Tailwind CSS 3.x
- Apollo Client 3.x (GraphQL)
- Recharts 2.x (charts)
- Vitest + React Testing Library (target 80%+ coverage)

## SDD Workflow

Subagent follows:

1. **SPEC**: Read requirements and P5 GraphQL schema
2. **TEST**: Write failing tests (RED)
3. **CODE**: Minimal implementation (GREEN)
4. **REFACTOR**: Clean up while tests pass
5. **REFLECT**: Check `dependency-reflection`
6. **SYNC**: Run `spec-sync-check` if applicable
7. **COMMIT**: Include all changes with descriptive message

## React/GraphQL Guidelines

- Use custom hooks for queries (useDashboard, useDevelopers, etc.)
- Implement loading and error states for all queries
- Use cache-and-network fetch policy for data freshness
- Follow WCAG 2.1 AA accessibility guidelines
- Use Tailwind utilities only (no custom CSS)
- Write integration tests for page components

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

Components:
- {component list}

GraphQL Queries:
- {queries/mutations used}

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

- P6 consumes analytics-core GraphQL (verify schema alignment)
- Use GraphQL Code Generator (avoid manual type definitions)
- Target 80%+ test coverage
- Ensure WCAG 2.1 AA accessibility
- Use Tailwind CSS (no custom CSS)
- Update task.md with results

---

See also:
- `.claude/rules/service/viz-spa.md` - GraphQL codegen and component constraints
- `.claude/skills/react-vite-patterns/` - React/Vite patterns
- `services/cursor-analytics-core/src/graphql/schema.graphql` - GraphQL schema reference
