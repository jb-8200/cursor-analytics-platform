# Spec-Driven Development (SDD) Checklist

**Purpose**: Enforce SDD methodology since automated hooks don't exist in Claude Code.

## CRITICAL: After Every Task Completion

When you complete a task and all tests pass, **you MUST**:

### 1. ✅ Verify Tests Pass
```bash
go test ./...  # or appropriate test command
```
**Status**: All packages ok

### 2. 📝 Stage Changes
```bash
git add <files-related-to-task>
```
**Only stage files for THIS task** - don't include unrelated changes

### 3. 💾 Commit with Descriptive Message
```bash
git commit -m "$(cat <<'EOF'
feat(service): complete Step XNN - Task Name

Brief description of what was accomplished.

## Changes
- List key changes
- With bullet points
- For clarity

Files changed: N files
Test status: All passing

Reference: path/to/spec or issue

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <model-name> <noreply@anthropic.com>
EOF
)"
```

### 4. 📊 Update Progress Tracking

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

### 5. 🔄 Only Then Proceed to Next Task

❌ **NEVER** move to the next task before:
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
- ❌ "Now let's move to Step B02..."
- ❌ "Ready for the next step?"
- ❌ "Step B01 complete! Would you like to continue..."

**STOP** and commit first!

## Correct Pattern

✅ "Step B01 complete. Let me commit these changes..."
✅ [Stages files, commits, updates docs]
✅ "Commit complete. Ready to start Step B02?"

## Exception: Documentation-Only Changes

You MAY skip commits for:
- README updates
- Comment additions
- Documentation clarifications

But **when in doubt, commit**.

## Hook Implementation Status

| Hook | Status | Alternative |
|------|--------|-------------|
| pre_prompt.py | ❌ Not implemented | Read DEVELOPMENT.md manually |
| pre_commit.py | ❌ Not implemented | Run tests before commit |
| post_test.py | ❌ Not implemented | **THIS CHECKLIST** |

## Integration with TodoWrite

Use TodoWrite to track:
```javascript
[
  {"content": "Complete Step B01", "status": "completed"},
  {"content": "Commit B01 changes", "status": "completed"},  // <-- ADD THIS
  {"content": "Update progress docs", "status": "completed"}, // <-- AND THIS
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
