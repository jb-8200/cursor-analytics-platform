# Claude Code Rules

Always-on guardrails that enforce constraints automatically.

---

## Overview

Rules are loaded automatically and always apply. They define:
- **NEVER**: Prohibited actions (safety, security)
- **ALWAYS**: Required actions (process, standards)

---

## Available Rules

| Rule | Scope | Purpose |
|------|-------|---------|
| `01-security.md` | Global | Secrets, PII, destructive ops, untrusted scripts |
| `02-repo-guardrails.md` | Global | Git safety, file protection, write scope |
| `03-coding-standards.md` | Global | Go/TS/React patterns, testing, error handling |
| `04-sdd-process.md` | Global | SDD post-task enforcement |
| `service/cursor-sim.md` | `services/cursor-sim/**` | API contract protection, CLI isolation |
| `service/analytics-core.md` | `services/cursor-analytics-core/**` | GraphQL schema, Prisma database |
| `service/viz-spa.md` | `services/cursor-viz-spa/**` | GraphQL codegen, component testing |

---

## Rule vs Skill vs Command

| Type | Behavior | Use For |
|------|----------|---------|
| **Rule** | Always loaded, always enforced | MUST/NEVER constraints |
| **Skill** | Auto-triggered by context | Knowledge, patterns, how-to |
| **Command** | Invoked with `/command` | Workflows, scripts |

**Decision**:
- Must ALWAYS happen? → Rule
- Provides helpful guidance? → Skill
- User invokes explicitly? → Command

---

## Examples

**Rule**: "Never commit secrets" → Always enforced
**Skill**: "Go error handling patterns" → Triggered when working with Go
**Command**: "/status" → Run when user types it

---

## Rule File Format

### Global Rules

```markdown
---
description: What this rule enforces
---

# Rule Name

## NEVER
- Prohibited action 1
- Prohibited action 2

## ALWAYS
- Required action 1
- Required action 2
```

### Path-Scoped Rules

```markdown
---
description: Service-specific constraints
paths: services/cursor-sim/**
---

# cursor-sim Rules

[Content applies only when working in matched paths]
```

---

## See Also

- `.claude/skills/` - Auto-triggered knowledge
- `.claude/commands/` - User-invoked workflows
- `docs/claude-docs/memory.md` - Rules in Claude Code
