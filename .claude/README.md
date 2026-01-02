# Claude Code Integration

This directory contains Claude Code configuration for the Cursor Analytics Platform.

---

## Quick Start

1. **Read session context**: `DEVELOPMENT.md` (current state, active work)
2. **Check active work**: `ls -la plans/active`
3. **Skills activate automatically** based on your request
4. **Follow SDD workflow**: Spec → Tests → Code → Commit

---

## Directory Structure

```
.claude/
├── DEVELOPMENT.md              # Session context (READ FIRST)
├── README.md                   # This file
├── MODEL_SELECTION_SUMMARY.md  # Model cost optimization
├── settings.local.json         # Claude Code settings
│
├── skills/                     # Knowledge guides (auto-discovered)
│   ├── go-best-practices/      # Go patterns (cursor-sim)
│   │   └── SKILL.md
│   ├── cursor-api-patterns/    # Cursor API implementation
│   │   └── SKILL.md
│   ├── sdd-checklist/          # Post-task commit enforcement
│   │   └── SKILL.md
│   ├── model-selection-guide/  # Model cost optimization
│   │   └── SKILL.md
│   ├── spec-process-core/      # Core SDD principles
│   │   └── SKILL.md
│   ├── spec-process-dev/       # TDD development workflow
│   │   └── SKILL.md
│   ├── spec-user-story/        # User story format
│   │   └── SKILL.md
│   ├── spec-design/            # Design doc format
│   │   └── SKILL.md
│   ├── spec-tasks/             # Task breakdown format
│   │   └── SKILL.md
│   ├── sdd-workflow/           # Full workflow reference
│   │   └── SKILL.md
│   └── spec-driven-development/ # SDD methodology
│       └── SKILL.md
│
├── commands/                   # Custom slash commands
│   ├── start-feature.md        # Initialize feature context
│   ├── complete-feature.md     # Verify and close feature
│   ├── implement.md            # TDD implementation workflow
│   ├── status.md               # Current project state
│   ├── next-task.md            # Find next work item
│   └── spec.md                 # Display service specification
│
├── hooks/                      # Claude Code hooks
│   ├── README.md               # Hook setup instructions
│   ├── pre_commit.py           # SDD checklist on commit
│   ├── markdown_formatter.py   # Auto-format markdown
│   └── sdd_reminder.py         # Post-task reminder
│
└── plans/
    ├── README.md               # Symlink mechanism docs
    └── active -> ...           # Symlink to current work item
```

---

## Skills (Auto-Discovered)

Skills provide specialized knowledge. Claude automatically discovers and activates them based on your request.

### How Skills Work

1. At startup, Claude loads skill names and descriptions
2. When your request matches a skill's description, Claude asks to use it
3. You approve, and Claude loads the full skill content

### Available Skills

| Skill | Description |
|-------|-------------|
| `go-best-practices` | Go coding standards, error handling, testing |
| `cursor-api-patterns` | Cursor API response formats, auth, pagination |
| `sdd-checklist` | Post-task commit enforcement (CRITICAL) |
| `model-selection-guide` | Choose Haiku vs Sonnet vs Opus |
| `spec-process-core` | Core SDD workflow principles |
| `spec-process-dev` | TDD RED-GREEN-REFACTOR cycle |
| `spec-user-story` | Create user stories with EARS format |
| `spec-design` | Create design docs with ADRs |
| `spec-tasks` | Create task breakdowns |
| `sdd-workflow` | Feature lifecycle and phase gates |
| `spec-driven-development` | Full SDD methodology reference |

### Triggering Skills

Skills activate automatically, but you can also reference explicitly:

```
"Help me implement a Go HTTP handler"
  → go-best-practices activates

"Create a user story for the export feature"
  → spec-user-story activates

"What should I do after completing this task?"
  → sdd-checklist activates
```

---

## Commands (Slash Commands)

Custom slash commands for workflow automation:

| Command | Purpose |
|---------|---------|
| `/start-feature {name}` | Initialize feature, create symlink, load context |
| `/complete-feature {name}` | Verify completion, remove symlink |
| `/implement {task-id}` | TDD implementation workflow |
| `/status` | Show current project state |
| `/next-task` | Find next work item |
| `/spec {service}` | Display service specification |

---

## Hooks

Claude Code hooks run shell commands at specific lifecycle events.

### Available Hooks

| Hook | Event | Purpose |
|------|-------|---------|
| `pre_commit.py` | PreToolUse (Bash) | SDD reminder on git commit |
| `markdown_formatter.py` | PostToolUse (Edit\|Write) | Auto-format markdown |
| `sdd_reminder.py` | Stop | Post-task workflow reminder |

### Setup

Run `/hooks` to configure, or add to `.claude/settings.local.json`.

See `hooks/README.md` for detailed setup instructions.

---

## SDD Workflow Quick Reference

### Session Start

1. Read `DEVELOPMENT.md`
2. Check `plans/active` symlink
3. Review current task status
4. Continue with TDD

### Feature Lifecycle

```
/start-feature {name}
    ↓
Read user-story.md + design.md + task.md
    ↓
For each task:
    RED → GREEN → REFACTOR → COMMIT
    ↓
/complete-feature {name}
```

### Task Completion (CRITICAL)

After every task, follow `sdd-checklist`:

1. Tests pass
2. Git commit
3. Update task.md
4. Update DEVELOPMENT.md
5. Next task

---

## Integration with Project

| File | Purpose |
|------|---------|
| `CLAUDE.md` | Always-included project context |
| `docs/spec-driven-design.md` | Full SDD methodology |
| `services/{service}/SPEC.md` | Technical specifications |
| `.work-items/{feature}/` | Active work tracking |

---

## Model Selection

| Task Type | Model | Why |
|-----------|-------|-----|
| Spec writing, architecture | **Opus** | Complex reasoning |
| Well-specified implementation | **Haiku** | Cost-effective |
| Complex implementation | **Sonnet** | Balanced capability |

See `model-selection-guide` skill for details.

---

## Best Practices

1. **Read DEVELOPMENT.md first** each session
2. **Let skills activate naturally** based on your requests
3. **Use TodoWrite** for multi-step work
4. **Commit after every task** (sdd-checklist)
5. **Keep CLAUDE.md minimal** - heavy content goes in skills
