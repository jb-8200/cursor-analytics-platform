# Design Document: Skills Cleanup & Catalog

**Feature ID**: P0-F03
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

Create skills catalog README, update sdd-checklist (after rules extraction), ensure consistent frontmatter.

---

## Skills Catalog README

Create `.claude/skills/README.md`:

```markdown
# Available Skills

Skills auto-activate when your request matches their description.

## Catalog

| Skill | Triggers | Service | Purpose |
|-------|----------|---------|---------|
| api-contract | "API", "endpoint", "contract" | P4 | cursor-sim API reference |
| go-best-practices | Go code, handlers, structs | P4 | Go coding standards |
| cursor-api-patterns | API implementation | P4 | Cursor API patterns |
| typescript-graphql-patterns | GraphQL, resolvers | P5 | TypeScript/GraphQL |
| react-vite-patterns | React, Vite, hooks | P6 | React/Vite patterns |
| sdd-checklist | "task complete", "commit" | All | Post-task workflow guidance |
| dependency-reflection | refactoring, model changes | All | Dependency checking |
| spec-sync-check | phase complete, endpoint added | All | SPEC.md update triggers |
| spec-process-core | SDD workflow | All | Core SDD principles |
| spec-process-dev | TDD implementation | All | TDD workflow |
| spec-user-story | user story creation | All | User story format |
| spec-design | design document | All | Design doc format |
| spec-tasks | task breakdown | All | Task format |
| model-selection-guide | model choice, cost | All | Model optimization |

## Skill vs Rule

| Type | Behavior | Use For |
|------|----------|---------|
| **Rule** | Always loaded, always enforced | Security, standards, process enforcement |
| **Skill** | Auto-triggered when relevant | Knowledge, patterns, how-to guidance |

If something MUST happen every time → Rule (`.claude/rules/`)
If something provides helpful guidance → Skill (`.claude/skills/`)
```

---

## sdd-checklist Update

After P0-F02 extracts enforcement to rules:

**Remove**:
- "you MUST" enforcement language
- "NEVER proceed" prohibitions

**Keep**:
- Detailed step-by-step guidance
- Examples and templates
- Integration with TodoWrite
- Commit message format

**Add**:
- Reference to `rules/04-sdd-process.md` for enforcement
- Note: "For enforcement, see `.claude/rules/04-sdd-process.md`"

---

## Frontmatter Audit

Ensure all skills have:

```yaml
---
name: skill-name
description: Clear description with trigger keywords
allowed-tools: [optional, only if restricting]
---
```

Skills to audit:
- [ ] api-contract
- [ ] cursor-api-patterns
- [ ] dependency-reflection
- [ ] go-best-practices
- [ ] model-selection-guide
- [ ] react-vite-patterns
- [ ] sdd-checklist
- [ ] spec-design
- [ ] spec-process-core
- [ ] spec-process-dev
- [ ] spec-sync-check
- [ ] spec-tasks
- [ ] spec-user-story
- [ ] typescript-graphql-patterns
