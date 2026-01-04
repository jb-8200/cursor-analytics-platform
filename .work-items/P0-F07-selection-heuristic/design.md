# Design Document: Selection Heuristic Guide

**Feature ID**: P0-F07
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

Add selection heuristic to `.claude/README.md` explaining when to use each Claude Code feature.

---

## Selection Heuristic Table

```markdown
## When to Use What

| Feature | Location | Behavior | Use For |
|---------|----------|----------|---------|
| **Rules** | `.claude/rules/` | Always loaded, always applied | Security, coding standards, SDD enforcement |
| **Commands** | `.claude/commands/` | Invoked with `/command` | Workflows, deployments, status checks |
| **Skills** | `.claude/skills/` | Auto-triggered by context | Knowledge, patterns, guidance |
| **Agents** | `.claude/agents/` | Spawned for isolated tasks | Parallel service development |
| **Hooks** | `settings.local.json` | Run on tool events | Formatting, validation, reminders |
| **Memory** | `CLAUDE.md` | Always included | Project overview, key instructions |
```

---

## Decision Tree

```markdown
### Decision Tree

When adding a new instruction:

1. **Must ALWAYS happen, no exceptions?**
   → **Rule** in `.claude/rules/`

2. **User explicitly invokes with `/command`?**
   → **Command** in `.claude/commands/`

3. **Provides knowledge when context matches?**
   → **Skill** in `.claude/skills/`

4. **Needs isolated context for parallel work?**
   → **Agent** in `.claude/agents/`

5. **Runs automatically on tool events?**
   → **Hook** in `settings.local.json`

6. **Project-wide context always needed?**
   → **Memory** in `CLAUDE.md`
```

---

## Examples Section

```markdown
### Examples

| Instruction | Type | Location |
|-------------|------|----------|
| "Never commit secrets" | Rule | rules/01-security.md |
| "Show project status" | Command | commands/status.md |
| "Go error handling patterns" | Skill | skills/go-best-practices/ |
| "Work on P4 CLI only" | Agent | agents/cursor-sim-cli-dev.md |
| "Format markdown after edit" | Hook | settings.local.json |
| "This is a monorepo with 3 services" | Memory | CLAUDE.md |
```
