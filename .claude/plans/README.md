# Plans Directory

This directory tracks the active work item using symlinks to `.work-items/` entries.

## Purpose

The `.claude/plans/active` symlink points to the current work item being implemented, making it easy for AI assistants to:

1. Quickly identify what's being worked on
2. Access user stories, design docs, and task lists
3. Maintain context across sessions

## Structure

```
.claude/plans/
├── README.md           # This file
├── active -> ../../.work-items/01-p0-scaffolding/   # Symlink to active work item
└── .gitignore          # Ignores the 'active' symlink
```

## Usage

### Creating a New Work Item

When starting a new feature with `/start-feature [name]`:

```bash
# Creates work item directory
.work-items/{nn}-{name}/
├── user-story.md       # Acceptance criteria (Given-When-Then)
├── design.md           # Technical design decisions
├── task.md             # Implementation checklist
└── test-plan.md        # Test cases (TDD)

# Updates active symlink
.claude/plans/active -> ../../.work-items/{nn}-{name}/
```

### Checking Active Work

AI assistants check `.claude/plans/active` to see:

```bash
# What am I working on?
ls -la .claude/plans/active

# Read active task list
cat .claude/plans/active/task.md

# Read acceptance criteria
cat .claude/plans/active/user-story.md
```

### Completing a Work Item

When work is done:

1. Mark all tasks complete in `task.md`
2. Remove the symlink: `rm .claude/plans/active`
3. Archive the work item (optional): move to `.work-items/archive/`

## Integration with Commands

- `/start-feature [name]` - Creates work item and sets as active
- `/next-task` - Reads from active work item if exists
- `/verify [service]` - Checks if active work matches service specs

## Git Tracking

- `.work-items/` directory is tracked in git
- `.claude/plans/active` symlink is ignored (in .gitignore)
- Each developer can work on different active items locally

## Best Practices

1. **One active item at a time** - Focus on completing before starting new work
2. **Keep task.md updated** - Check off items as you complete them
3. **TDD workflow** - Write tests first (test-plan.md → test files → implementation)
4. **Spec alignment** - Reference service SPEC.md in design.md
5. **Clear acceptance criteria** - User stories should have testable scenarios

## Example Workflow

```bash
# Start new feature
/start-feature basic-auth-middleware

# Creates and activates:
# .work-items/03-basic-auth-middleware/
# .claude/plans/active -> ../../.work-items/03-basic-auth-middleware/

# Work on it (TDD)
# 1. Read .claude/plans/active/test-plan.md
# 2. Write failing tests
# 3. Implement to make tests pass
# 4. Check off tasks in .claude/plans/active/task.md

# When done
rm .claude/plans/active
git commit -am "feat(cursor-sim): implement basic auth middleware"
```

## Directory State

- **Empty** (no `active` symlink): No active work item, ready to start new feature
- **Has `active` symlink**: Currently working on the linked work item
- **Broken `active` symlink**: Work item was deleted or moved (clean up with `rm active`)

---

**Remember**: The active symlink is a development-time tracking tool. The actual work artifacts live in `.work-items/` and are committed to git.
