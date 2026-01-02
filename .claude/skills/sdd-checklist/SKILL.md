---
name: sdd-checklist
description: Post-task completion checklist for Spec-Driven Development. Use after completing any implementation task to ensure proper commit hygiene. Triggers on "task complete", "step done", "ready to commit", or moving to next task.
---

# SDD Checklist

**Purpose**: Enforce SDD methodology after every completed task.

## CRITICAL: After Every Task Completion

When you complete a task and all tests pass, **you MUST**:

### 1. Verify Tests Pass

```bash
go test ./...  # or appropriate test command
```

**Status**: All packages ok

### 2. Stage Changes

```bash
git add <files-related-to-task>
```

**Only stage files for THIS task** - don't include unrelated changes

### 3. Commit with Descriptive Message

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

### 4. Update Progress Tracking

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

### 5. Only Then Proceed to Next Task

**NEVER** move to the next task before:
1. Tests passing
2. Code committed
3. Progress documented

## Why This Matters

**Without commits:**
- Work can be lost
- History is unclear
- Collaboration breaks
- Can't track progress
- Hard to debug issues

**SDD Flow:**

```
Spec → Tests → Implementation → Tests Pass → COMMIT → Next Task
                                              ^^^^^^
                                           YOU ARE HERE
```

## Red Flags

If you catch yourself saying any of these **WITHOUT committing first**:
- "Now let's move to Step B02..."
- "Ready for the next step?"
- "Step B01 complete! Would you like to continue..."

**STOP** and commit first!

## Correct Pattern

1. "Step B01 complete. Let me commit these changes..."
2. [Stages files, commits, updates docs]
3. "Commit complete. Ready to start Step B02?"

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

**Every completed task = 5 steps:**
1. Tests pass
2. Git commit
3. Update task.md
4. Update DEVELOPMENT.md
5. Proceed to next

**No exceptions.** This is the SDD way.
