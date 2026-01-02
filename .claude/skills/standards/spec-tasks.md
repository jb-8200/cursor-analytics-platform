# Task Breakdown Standard

**Trigger**: Creating task breakdowns, implementation plans, work decomposition, sprint planning

---

## Purpose

Task breakdowns decompose a feature into small, independently completable work units. Each task should be 1-4 hours and result in a single commit.

## Output Location

`.work-items/{feature-name}/task.md`

---

## Template

```markdown
# Task Breakdown: {Feature Name}

## Summary

**User Story**: `.work-items/{feature}/user-story.md`
**Design Doc**: `.work-items/{feature}/design.md`
**Service**: {cursor-sim | cursor-analytics-core | cursor-viz-spa}

**Total Estimated Hours**: {sum of all tasks}
**Total Tasks**: {count}

## Progress Tracker

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| 01 | {Task name} | 2.0 | DONE | 1.5 |
| 02 | {Task name} | 1.5 | IN_PROGRESS | - |
| 03 | {Task name} | 2.5 | NOT_STARTED | - |
| 04 | {Task name} | 1.0 | NOT_STARTED | - |

**Current Step**: {NN}

---

## Detailed Steps

### Step 01: {Task Name}

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

---

### Step 02: {Task Name}

{Same structure as Step 01...}

---

## Dependencies

{Show task dependencies if any}

```
Step 01 ──┬──▶ Step 02 ──▶ Step 04
          │
          └──▶ Step 03 ──┘
```

## Testing Milestones

| After Step | Coverage Target | Test Type |
|------------|-----------------|-----------|
| 01 | 80% | Unit |
| 02 | 85% | Unit + Integration |
| Final | 90% | Unit + Integration + E2E |

## Notes

{Any implementation notes, gotchas, or reminders}

## Related Documents

- `.work-items/{feature}/user-story.md` - Requirements
- `.work-items/{feature}/design.md` - Technical design
- `services/{service}/SPEC.md` - Service specification
```

---

## Writing Guidelines

### Task Sizing

**Target: 1-4 hours per task**

| Size | Description | Example |
|------|-------------|---------|
| Too small (< 1h) | Combine with related work | "Add import statement" |
| Just right (1-4h) | Clear scope, one commit | "Implement health check endpoint" |
| Too large (> 4h) | Break into smaller tasks | "Implement full API" |

### Task Naming

Use action verbs + specific target:

Good:
- "Implement ResearchRow struct with Parquet tags"
- "Add GetDatasetByTimeRange to storage interface"
- "Create /research/dataset handler"

Bad:
- "Work on research export"
- "Fix stuff"
- "Finish implementation"

### Status Values

| Status | Meaning |
|--------|---------|
| `NOT_STARTED` | Task not yet begun |
| `IN_PROGRESS` | Currently working on |
| `BLOCKED` | Waiting on dependency |
| `DONE` | Committed and documented |

**Only ONE task should be IN_PROGRESS at a time.**

### Subtask Checklist

Each step should have:
- [ ] Implementation subtasks
- [ ] Write tests
- [ ] Run tests
- [ ] Update documentation (if needed)

### Acceptance Criteria

Each task needs clear "done" criteria:

Bad:
```
Acceptance Criteria:
- Works correctly
```

Good:
```
Acceptance Criteria:
- ResearchRow struct has all 22 fields from design.md
- All fields have correct Parquet tags
- JSON marshaling produces snake_case keys
- Unit test covers serialization round-trip
```

---

## Example

```markdown
# Task Breakdown: Research Dataset Export

## Summary

**User Story**: `.work-items/cursor-sim-phase3/user-story.md`
**Design Doc**: `.work-items/cursor-sim-phase3/design.md`
**Service**: cursor-sim

**Total Estimated Hours**: 12.5h
**Total Tasks**: 6

## Progress Tracker

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| 01 | ResearchRow struct + Parquet tags | 1.5 | DONE | 1.0 |
| 02 | DatasetBuilder implementation | 2.5 | DONE | 2.0 |
| 03 | Parquet exporter | 2.0 | IN_PROGRESS | - |
| 04 | CSV/JSON exporters | 1.5 | NOT_STARTED | - |
| 05 | HTTP handlers | 2.5 | NOT_STARTED | - |
| 06 | Integration tests | 2.5 | NOT_STARTED | - |

**Current Step**: 03

---

### Step 03: Parquet Exporter

**Estimated**: 2.0h
**Status**: IN_PROGRESS
**Actual**: -

**Files**:
- `internal/export/parquet.go`
- `internal/export/parquet_test.go`

**Tasks**:
- [ ] Add parquet-go dependency
- [ ] Implement ParquetExporter struct
- [ ] Add Export(w io.Writer, rows []ResearchRow) method
- [ ] Add SNAPPY compression option
- [ ] Write unit tests
- [ ] Run tests

**Acceptance Criteria**:
- Export produces valid Parquet file readable by pandas
- File uses SNAPPY compression
- All 22 columns present with correct types
- Unit test verifies round-trip: write → read → compare

**TDD Approach**:
1. RED: Test that Export() produces readable Parquet
2. GREEN: Implement minimal Export with parquet-go
3. REFACTOR: Add compression option, clean up
```

---

## Process

1. **Read user story** for requirements
2. **Read design doc** for technical approach
3. **Identify major components** from design
4. **Break into 1-4 hour tasks**
5. **Order by dependencies**
6. **Add TDD approach** for each task
7. **Create file** in `.work-items/{feature}/task.md`
8. **Begin Step 01** with `/implement` or TDD cycle

---

## Task Completion Workflow

After completing each task:

1. ✅ All tests pass
2. ✅ Stage changes: `git add {files}`
3. ✅ Commit with descriptive message
4. ✅ Update task.md: Status → DONE, add Actual hours
5. ✅ Update DEVELOPMENT.md
6. ✅ Start next task

**See `sdd-checklist` skill for enforcement.**

---

## Red Flags

| Flag | Fix |
|------|-----|
| Task > 4 hours | Break into smaller tasks |
| Multiple tasks IN_PROGRESS | Complete one before starting another |
| No acceptance criteria | Add testable criteria |
| Vague task names | Use action verb + specific target |
| Missing test subtasks | Add "Write tests" and "Run tests" |
