# Claude Code Integration

This directory contains Claude Code configuration for the Cursor Analytics Platform.

---

## Quick Start

1. **Read session context**: `DEVELOPMENT.md` (current state, active work)
2. **Check active work**: `ls -la plans/active`
3. **Reference skills** when needed (see below)
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
├── skills/                     # Knowledge guides (categorized)
│   ├── process/                # Workflow stages
│   │   ├── spec-process-core.md    # Core SDD principles
│   │   └── spec-process-dev.md     # TDD development workflow
│   │
│   ├── standards/              # Artifact templates
│   │   ├── spec-user-story.md      # User story format
│   │   ├── spec-design.md          # Design doc format
│   │   └── spec-tasks.md           # Task breakdown format
│   │
│   ├── guidelines/             # Technology-specific
│   │   ├── go-best-practices.md    # Go patterns (cursor-sim)
│   │   └── cursor-api-patterns.md  # Cursor API implementation
│   │
│   └── operational/            # Day-to-day enforcement
│       ├── sdd-checklist.md        # Post-task commit enforcement
│       ├── sdd-workflow.md         # Full workflow reference
│       ├── model-selection-guide.md # Model optimization
│       └── spec-driven-development.md # SDD methodology
│
├── commands/                   # Custom slash commands
│   ├── start-feature.md        # Initialize feature context
│   ├── complete-feature.md     # Verify and close feature
│   ├── implement.md            # TDD implementation workflow
│   ├── status.md               # Current project state
│   ├── next-task.md            # Find next work item
│   └── spec.md                 # Display service specification
│
├── hooks/                      # Documentation only (NOT EXECUTED)
│   ├── README.md               # Explains limitations + alternatives
│   ├── pre_prompt.py           # Context injection (doc only)
│   ├── pre_commit.py           # Test enforcement (doc only)
│   └── pre_patch.py            # Lint enforcement (doc only)
│
└── plans/
    ├── README.md               # Symlink mechanism docs
    └── active -> ...           # Symlink to current work item
```

---

## Skills (Knowledge Guides)

Skills provide specialized knowledge. They activate based on semantic match or explicit reference.

### Process Skills

For workflow guidance:

| Skill | Use When |
|-------|----------|
| `spec-process-core` | Starting work, understanding workflow |
| `spec-process-dev` | Implementing with TDD |

### Standards Skills

For creating artifacts:

| Skill | Use When |
|-------|----------|
| `spec-user-story` | Creating user stories, PRDs |
| `spec-design` | Creating design docs, ADRs |
| `spec-tasks` | Creating task breakdowns |

### Guidelines Skills

For technology-specific patterns:

| Skill | Use When |
|-------|----------|
| `go-best-practices` | Writing Go code |
| `cursor-api-patterns` | Implementing API endpoints |

### Operational Skills

For day-to-day workflow:

| Skill | Use When |
|-------|----------|
| `sdd-checklist` | After completing any task (CRITICAL) |
| `model-selection-guide` | Choosing which model to use |

### Triggering Skills

Reference explicitly for reliable activation:

```
"Following spec-process-core, let's plan this feature"
"Using spec-user-story, create a user story for login"
"Apply go-best-practices to implement the handler"
"Following sdd-checklist, commit the changes"
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

## Hooks (Not Executed)

**IMPORTANT**: The Python hooks in `hooks/` are documentation-only. Claude Code does not execute them.

See `hooks/README.md` for:
- What hooks were designed to do
- Alternative implementations
- How to enforce workflow manually

**Alternative**: Use `sdd-checklist` skill + TodoWrite for enforcement.

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

See `skills/operational/model-selection-guide.md` for details.

---

## Best Practices

1. **Read DEVELOPMENT.md first** each session
2. **Reference skills explicitly** in requests
3. **Use TodoWrite** for multi-step work
4. **Commit after every task** (sdd-checklist)
5. **Keep CLAUDE.md minimal** - heavy content goes in skills
