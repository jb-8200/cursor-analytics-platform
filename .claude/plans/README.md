# Active Plans Directory

This directory tracks the currently active work item using symlinks.

## Mechanism

When a feature is started with `/start-feature {name}`, a symlink is created:

```
.claude/plans/active -> ../../.work-items/{name}/task.md
```

This provides:
1. **Clear signal** of what is currently being worked on
2. **Git-tracked state** of active work
3. **Machine-readable** for hooks to discover context

## Discovery

Claude Code checks this folder at session start to understand focus:

```bash
# Check what's active
ls -la .claude/plans/active

# If symlink exists, read the task file
cat "$(readlink .claude/plans/active)"
```

## Lifecycle

1. **Start**: `/start-feature cursor-sim-v2`
   - Creates symlink: `active -> ../../.work-items/cursor-sim-v2/task.md`
   - Commits change to git

2. **Work**: Normal TDD loop on steps
   - Update step status in task.md
   - Commit with time tracking

3. **Complete**: `/complete-feature cursor-sim-v2`
   - Verifies all steps done
   - Removes symlink
   - Commits completion

## Current State

Check `active` symlink (if present) to see active feature.
