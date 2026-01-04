# Task Breakdown: Commands/Prompts Restructure

**Feature ID**: P0-F04
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: COMPLETE ✅

---

## Progress Tracker

| Task | Description | Status | Time Est | Time Actual |
|------|-------------|--------|----------|-------------|
| TASK01 | Create commands/subagent/ directory | COMPLETE | 0.1h | 0.1h |
| TASK02 | Convert cursor-sim-cli template to command | COMPLETE | 0.25h | 0.25h |
| TASK03 | Convert analytics-core template to command | COMPLETE | 0.25h | 0.25h |
| TASK04 | Convert viz-spa template to command | COMPLETE | 0.25h | 0.25h |
| TASK05 | Remove .claude/prompts/ directory | COMPLETE | 0.1h | 0.1h |
| TASK06 | Update DEVELOPMENT.md references | COMPLETE | 0.1h | 0.1h |

**Total Estimated**: 1.0 hours
**Total Actual**: 1.0 hours
**Status**: Feature COMPLETE (6/6 tasks)

---

## Task Details

### TASK01: Create Subagent Commands Directory

**Objective**: Set up new command location.

**Deliverables**:
- `.claude/commands/subagent/` directory

---

### TASK02-04: Convert Templates to Commands

**Objective**: Add frontmatter, convert placeholders.

**Changes per file**:
- Add frontmatter: description, argument-hint, allowed-tools
- Replace `{feature-id}` → `$1`
- Replace `{task-id}` → `$2`
- Reference rules instead of inline constraints
- Slim down (link to rules/skills)

---

### TASK05: Remove Prompts Directory

**Objective**: Clean up old structure.

**Command**:
```bash
rm -rf .claude/prompts/
```

---

### TASK06: Update References

**Objective**: Fix documentation pointing to old location.

**Files to update**:
- `.work-items/P0-F01-sdd-subagent-orchestration/design.md`
- `.claude/README.md` (if referenced)

---

## Dependencies

- P0-F02 (Rules Layer) - service rules must exist for reference
