# Claude Code Integration

This directory contains Claude Code-specific configuration and guidance.

## What Actually Works in Claude Code

### ✅ Skills (Knowledge Guides)
Skills are loaded automatically when relevant topics are discussed. They provide specialized knowledge:

- **`skills/cursor-api-patterns.md`** - Cursor Business API implementation patterns
- **`skills/go-best-practices.md`** - Go coding standards for this project
- **`skills/model-selection-guide.md`** - Which model to use for each task

**How to use:** Just reference them in conversation:
```
"Following cursor-api-patterns.md, implement the /v1/analytics/ai-code/commits endpoint"
"Using go-best-practices.md patterns, create the Developer struct"
"Based on model-selection-guide.md, which model should I use for TASK-SIM-004?"
```

### ✅ DEVELOPMENT.md (Session Context)
Read `DEVELOPMENT.md` at the start of each session to understand:
- Current project status
- Recent work completed
- Active focus area
- Next steps

### ✅ Model Selection
Use the model-selection-guide.md to optimize cost/performance:

```
"Use Haiku to implement TASK-SIM-003 (well-specified struct)"
"Use Sonnet to implement TASK-SIM-004 (complex Poisson logic)"
```

## What Doesn't Work

### ❌ Custom Slash Commands
Claude Code doesn't support custom slash commands via `.md` files. Instead:
- Use direct requests: "Implement TASK-SIM-001"
- Reference the skills: "Following model-selection-guide.md..."

### ❌ Automated Workflows
There's no automatic `/implement` or `/start-feature` command. Instead:
- Ask Claude directly to implement tasks
- Reference the patterns in skills/

## Recommended Workflow

### Starting a New Task

```
"I want to implement TASK-SIM-003: Developer Profile Generator

1. Read services/cursor-sim/SPEC.md lines 145-250
2. Follow go-best-practices.md for struct patterns
3. Use Haiku since it's well-specified
4. Write tests first (TDD)"
```

### Getting Model Recommendations

```
"What model should I use for TASK-SIM-004 according to model-selection-guide.md?"
```

### Following Best Practices

```
"Implement the health check endpoint following cursor-api-patterns.md response format"
```

## Directory Structure

```
.claude/
├── README.md                          # This file
├── DEVELOPMENT.md                     # Session context (read first!)
├── MODEL_SELECTION_SUMMARY.md         # Model optimization guide
├── settings.local.json                # Claude Code settings
├── skills/                            # Knowledge guides
│   ├── cursor-api-patterns.md         # API implementation patterns
│   ├── go-best-practices.md           # Go coding standards
│   ├── model-selection-guide.md       # Task-to-model mapping
│   └── spec-driven-development.md     # SDD methodology
└── hooks/                             # Hook descriptions (documentation only)
    └── README.md
```

## Best Practices

1. **Read DEVELOPMENT.md first** each session
2. **Reference skills explicitly** in your requests
3. **Specify model preference** when you know what you want
4. **Use TodoWrite** to track multi-step work
5. **Commit frequently** as you complete tasks

## Example Session

```
# Session start
User: "Read DEVELOPMENT.md and tell me what to work on next"

Claude: [Reads DEVELOPMENT.md] "You should implement TASK-SIM-001 next..."

User: "Use Sonnet to implement TASK-SIM-001 following the SPEC.md"

Claude: [Implements with Sonnet following specs and best practices]

User: "Now use Haiku to implement TASK-SIM-003"

Claude: [Implements with Haiku since it's well-specified]
```

## Integration with Project Docs

This `.claude/` folder works together with:
- `CLAUDE.md` - Project instructions (root level)
- `services/*/SPEC.md` - Service specifications
- `docs/TASKS.md` - Implementation task list
- `docs/USER_STORIES.md` - Acceptance criteria

## Summary

**Keep it simple:**
- Skills = Knowledge that Claude loads automatically
- DEVELOPMENT.md = Current state
- Direct requests = How you interact with Claude
- No magic commands needed!
