---
name: sdd-checklist
description: Post-task completion checklist for Spec-Driven Development. Use after completing any implementation task to ensure proper commit hygiene. Triggers on "task complete", "step done", "ready to commit", or moving to next task.
---

# SDD Checklist

**Purpose**: Detailed guidance for SDD task completion workflow.

**ENFORCEMENT**: See `.claude/rules/04-sdd-process.md` for MUST/NEVER requirements.

## After Every Task Completion: 7-Step Workflow

Follow these steps after completing a task and verifying tests pass:

### 1. Verify Tests Pass

```bash
go test ./...  # or appropriate test command
```

**Status**: All packages ok

### 2. Check Dependency Reflections (REFLECT)

Run the **dependency-reflection** check:

**Ask yourself:**
- Did I modify models? → Check generators, handlers, SPEC.md schemas, tests
- Did I add/modify endpoints? → Check SPEC.md endpoints table, E2E tests
- Did I change storage interface? → Check all handlers using storage
- Did I complete a phase step? → Check SPEC.md status, task.md, DEVELOPMENT.md
- Did I refactor code? → Check all tests still pass, docs still accurate

**Reflection Checklist:**
- [ ] Documentation sync verified (SPEC.md, task.md, DEVELOPMENT.md)
- [ ] Code sync verified (generators, handlers, storage match)
- [ ] Test sync verified (tests cover new/modified behavior)

**Run regression tests** if any medium/high-priority reflections detected.

### 3. Update SPEC.md if Needed (SYNC)

Run the **spec-sync-check**:

**High-Priority Triggers** (MUST update SPEC.md):
- [ ] Completed a phase step (update Implementation Status table)
- [ ] Added new endpoint (update Endpoints table)
- [ ] Added new service/package (update Package Structure)
- [ ] Modified CLI/config (update CLI Configuration)

**If triggered:**
1. Open `services/{service}/SPEC.md`
2. Update relevant sections (see spec-sync-check for details)
3. Update "Last Updated" date
4. Verify accuracy of changes

**Include SPEC.md in commit** if updated.

### 4. Stage Changes

```bash
git add <files-related-to-task>
# If SPEC.md was updated:
git add services/{service}/SPEC.md
```

**Stage files for THIS task** including:
- Code changes
- Test changes
- SPEC.md (if updated during SYNC step)
- Do NOT include unrelated changes

### 5. Commit with Descriptive Message

```bash
git commit -m "$(cat <<'EOF'
feat(service): complete Step XNN - Task Name

Brief description of what was accomplished.

## Changes
- List key changes
- With bullet points

Files changed: N files
Test status: All passing

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
)"
```

### 6. Update Progress Tracking

**a) Update task.md progress tracker:**

```markdown
| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| XNN | Task Name | 2.0 | DONE | 1.5 |
```

**b) Update .claude/DEVELOPMENT.md:**
- Current status
- Recently completed work
- Next steps

### 7. Only Then Proceed to Next Task

Verify all steps complete before moving forward.

**See `.claude/rules/04-sdd-process.md` for enforcement requirements.**

## Workflow Context

The SDD workflow ties together:
- **SPEC**: What needs to be done
- **TEST**: Verify expectations
- **CODE**: Minimal implementation
- **REFLECT**: Check dependencies (dependency-reflection skill)
- **SYNC**: Update SPEC.md (spec-sync-check skill)
- **COMMIT**: Record progress
- **NEXT**: Move forward confidently

Each step builds on the previous, ensuring quality and clarity.

## Integration with TodoWrite

Use TodoWrite to track:

```javascript
[
  {"content": "Complete Step B01", "status": "completed"},
  {"content": "Commit B01 changes", "status": "completed"},
  {"content": "Update progress docs", "status": "completed"},
  {"content": "Start Step B02", "status": "in_progress"}
]
```

## Summary

**Every completed task = 7 steps:**
1. Tests pass
2. **Check reflections (REFLECT)** ← NEW
3. **Update SPEC.md if triggered (SYNC)** ← NEW
4. Stage changes (code + SPEC.md if updated)
5. Git commit
6. Update task.md and DEVELOPMENT.md
7. Proceed to next

**No exceptions.** This is the enhanced SDD way.

## Quick Reference

Before committing, verify:
- ✅ Tests pass
- ✅ Reflections checked (use dependency-reflection)
- ✅ SPEC.md updated if needed (use spec-sync-check)
- ✅ All related files staged (code + docs)
- ✅ Commit message descriptive
- ✅ Progress documented
- ✅ Ready for next task
