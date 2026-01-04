# viz-spa-dev Subagent Prompt Template

## Task Context

**Feature**: {feature-id} - {feature-name}
**Task**: {task-id} - {task-description}
**Estimated Time**: {time-estimate}

## Objective

{Clear, specific objective for this task}

## Scope Constraints

**YOU MUST**:
- Work ONLY on P6 code in `services/cursor-viz-spa/`
- Follow React/Vite best practices (use `react-vite-patterns` skill)
- Consume analytics-core GraphQL API (check schema alignment)
- Write tests first (TDD: RED → GREEN → REFACTOR)
- Update `.work-items/{feature}/task.md` with progress
- Run `sdd-checklist` before committing

**YOU MUST NEVER**:
- Modify cursor-sim code (P4)
- Modify analytics-core code (P5)
- Manually define GraphQL types (use codegen when available)
- Update `.claude/DEVELOPMENT.md` (master agent only)
- Modify plan folder symlinks
- Make cross-service changes

## Context Files

Relevant files for this task:
- `.work-items/{feature}/user-story.md` - Requirements
- `.work-items/{feature}/design.md` - Technical approach
- `.work-items/{feature}/task.md` - Task details
- `services/cursor-analytics-core/src/graphql/schema.graphql` - Upstream GraphQL schema

## Tech Stack

- **Framework**: React 18+, TypeScript 5.x
- **Build**: Vite 5.x
- **GraphQL**: Apollo Client 3.x
- **Styling**: Tailwind CSS 3.x
- **Charts**: Recharts 2.x
- **Testing**: Vitest + React Testing Library (target 80%+ coverage)

## Acceptance Criteria

{List acceptance criteria from task.md}

## SDD Workflow

Follow this workflow:

1. **SPEC**: Read task requirements and P5 GraphQL schema
2. **TEST**: Write failing tests (RED)
3. **CODE**: Minimal implementation (GREEN)
4. **REFACTOR**: Clean up while tests pass
5. **REFLECT**: Run `dependency-reflection` skill
6. **SYNC**: Run `spec-sync-check` skill (if applicable)
7. **COMMIT**: Create git commit with descriptive message

## React/GraphQL Guidelines

- Use custom hooks for GraphQL queries (useDashboard, useDevelopers, etc.)
- Implement proper loading and error states
- Use cache-and-network fetch policy for data freshness
- Follow accessibility guidelines (WCAG 2.1 AA)
- Use Tailwind utility classes (avoid custom CSS)
- Write integration tests for page components

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
- Components: {component-list}
- Tests: {test-file-path}

GraphQL Queries:
- {list any queries or mutations used}

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
- Components: {component-list}
- Tests: {test-file-path}

**Commit**: {commit-hash}

**Next Steps**:
- Report completion to master agent
```

## Remember

- P6 consumes analytics-core GraphQL API (verify schema alignment)
- Use GraphQL Code Generator when available (avoid manual type definitions)
- Target 80%+ test coverage
- Ensure WCAG 2.1 AA accessibility compliance
- Use Tailwind CSS (no custom CSS files)
- Document all changes in task.md
