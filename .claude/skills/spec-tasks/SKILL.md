---
name: spec-tasks
description: Create or revise task breakdowns for feature implementation. Use when decomposing features into work units, planning sprints, or tracking implementation progress. Produces task.md files with progress tracker and TDD approach. (project)
---

# Task Breakdown Standard

Task breakdowns decompose a feature into small, independently completable work units. Each task should be 1-4 hours and result in a single commit.

## Hierarchy

```
Phase (P#) = Epic level
  └── Feature (F##) = Work item directory
       └── Task (TASK##) = Implementation step
```

## Output Location

`.work-items/{P#-F##-feature-name}/task.md`

## Template

```markdown
# Task Breakdown: {Feature Name}

**Feature ID**: P#-F##-feature-name
**Phase**: P# (Phase Name)
**Created**: {Date}
**Status**: IN PROGRESS

## Summary

**User Story**: `.work-items/{feature}/user-story.md`
**Design Doc**: `.work-items/{feature}/design.md`
**Service**: {cursor-sim | cursor-analytics-core | cursor-viz-spa}

**Total Estimated Hours**: {sum of all tasks}
**Total Tasks**: {count}

## Progress Tracker

| Task ID | Task | Hours | Status | Actual |
|---------|------|-------|--------|--------|
| TASK01 | {Task name} | 2.0 | DONE | 1.5 |
| TASK02 | {Task name} | 1.5 | IN_PROGRESS | - |
| TASK03 | {Task name} | 2.5 | NOT_STARTED | - |

**Current Task**: TASK##

---

## Task Details

### TASK01: {Task Name}

**Estimated**: {hours}h
**Status**: {NOT_STARTED | IN_PROGRESS | DONE}
**Actual**: {hours when done}

**Files**:
- `path/to/file1.go`
- `path/to/file2.go`

**Tasks**:
- [ ] Subtask 1
- [ ] Subtask 2
- [ ] Write tests
- [ ] Run tests

**Acceptance Criteria**:
- {Testable criterion 1}
- {Testable criterion 2}

**TDD Approach**:
1. RED: {What test to write first}
2. GREEN: {Minimal implementation}
3. REFACTOR: {What to clean up}
```

## Task Sizing

**Target: 1-4 hours per task**

| Size | Description | Example |
|------|-------------|---------|
| Too small (< 1h) | Combine with related work | "Add import statement" |
| Just right (1-4h) | Clear scope, one commit | "Implement health check endpoint" |
| Too large (> 4h) | Break into smaller tasks | "Implement full API" |

## Task Naming

Use action verbs + specific target:

Good:
- "Implement ResearchRow struct with Parquet tags"
- "Add GetDatasetByTimeRange to storage interface"
- "Create /research/dataset handler"

Bad:
- "Work on research export"
- "Fix stuff"
- "Finish implementation"

## Status Values

| Status | Meaning |
|--------|---------|
| `NOT_STARTED` | Task not yet begun |
| `IN_PROGRESS` | Currently working on |
| `BLOCKED` | Waiting on dependency |
| `DONE` | Committed and documented |

**Only ONE task should be IN_PROGRESS at a time.**

## Task Completion Workflow

After completing each task:

1. All tests pass
2. Stage changes: `git add {files}`
3. Commit with descriptive message
4. Update task.md: Status → DONE, add Actual hours
5. Update DEVELOPMENT.md
6. Start next task

**See `sdd-checklist` skill for enforcement.**

## Red Flags

| Flag | Fix |
|------|-----|
| Task > 4 hours | Break into smaller tasks |
| Multiple tasks IN_PROGRESS | Complete one before starting another |
| No acceptance criteria | Add testable criteria |
| Vague task names | Use action verb + specific target |
| Missing test subtasks | Add "Write tests" and "Run tests" |
