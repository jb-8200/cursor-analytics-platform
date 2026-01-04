# Design Document: Rules Layer Implementation

**Feature ID**: P0-F02
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

Implement `.claude/rules/` directory with modular, always-on guardrails following Claude Code's recommended patterns.

---

## Architecture

### Directory Structure

```
.claude/rules/
├── README.md                    # Explain rules + when to use
├── 01-security.md               # Secrets, PII, auth, destructive ops
├── 02-repo-guardrails.md        # Git safety, file protection, write scope
├── 03-coding-standards.md       # Shared Go/TS/React defaults
├── 04-sdd-process.md            # MUST-DO items (extracted from sdd-checklist)
└── service/
    ├── cursor-sim.md            # paths: services/cursor-sim/**
    ├── analytics-core.md        # paths: services/cursor-analytics-core/**
    └── viz-spa.md               # paths: services/cursor-viz-spa/**
```

---

## Rule File Format

### Global Rules (No paths)

```markdown
---
description: Brief description of what this rule enforces
---

# Rule Name

## NEVER
- List of prohibited actions

## ALWAYS
- List of required actions

## Context
Additional guidance when needed
```

### Path-Scoped Rules

```markdown
---
description: Service-specific guardrails
paths: services/cursor-sim/**
---

# cursor-sim Rules

[Content only applies when working in matched paths]
```

---

## Rule Content

### 01-security.md

**Purpose**: Security and compliance guardrails

**NEVER**:
- Commit secrets (.env, credentials.json, API keys)
- Run destructive git commands (force push main, hard reset)
- Execute untrusted scripts without sandbox
- Expose internal paths or credentials in output
- Store PII in logs or comments

**ALWAYS**:
- Validate file paths before operations
- Use environment variables for secrets
- Prefer sandbox mode for unknown scripts
- Review hook code before enabling

### 02-repo-guardrails.md

**Purpose**: Repository safety

**NEVER**:
- Modify files outside project root
- Delete files without explicit user request
- Overwrite user changes without confirmation
- Use `git push --force` to main/master
- Skip pre-commit hooks (--no-verify)

**ALWAYS**:
- Use absolute paths for clarity
- Confirm before bulk operations
- Preserve file encoding (UTF-8)
- Stage only task-related files

### 03-coding-standards.md

**Purpose**: Shared coding standards

**Go (cursor-sim)**:
- Use `gofmt` formatting
- Error handling: return errors, don't panic
- Test files: `*_test.go` with table-driven tests

**TypeScript (analytics-core)**:
- Strict mode enabled
- Explicit return types
- Use interfaces over types for objects

**React (viz-spa)**:
- Functional components with hooks
- Tailwind for styling (no custom CSS)
- WCAG 2.1 AA accessibility

### 04-sdd-process.md

**Purpose**: SDD enforcement (extracted from sdd-checklist)

**After EVERY Task Completion, you MUST**:
1. Verify tests pass
2. Check dependency reflections
3. Update SPEC.md if triggered
4. Git commit (code + docs)
5. Update task.md status
6. Update DEVELOPMENT.md

**NEVER proceed to next task without completing all steps.**

### service/cursor-sim.md

**Purpose**: cursor-sim specific guardrails
**Paths**: `services/cursor-sim/**`

**API Contract Protection**:
- NEVER modify `internal/api/` without updating SPEC.md
- ALL endpoint changes require E2E tests
- Keep response schemas backward-compatible for P5/P6

**CLI Isolation**:
- CLI subagent MUST NOT touch `internal/api/` or `internal/generator/`

### service/analytics-core.md

**Purpose**: analytics-core specific guardrails
**Paths**: `services/cursor-analytics-core/**`

**GraphQL Schema**:
- ALL schema changes must be backward compatible
- Document breaking changes for P6
- Use cursor-based pagination for lists

**Database**:
- Use Prisma for all DB access
- Never raw SQL without parameterization

### service/viz-spa.md

**Purpose**: viz-spa specific guardrails
**Paths**: `services/cursor-viz-spa/**`

**GraphQL Client**:
- Use generated types from codegen (when available)
- NEVER manually define GraphQL types
- Verify queries match P5 schema

**Components**:
- All components must have tests
- Use renderWithProviders for testing

---

## Integration

### CLAUDE.md Update

Slim down CLAUDE.md to reference rules:

```markdown
## Guardrails

Always-on rules in `.claude/rules/`:
- Security: secrets, PII, destructive ops
- Repo: git safety, file protection
- Standards: Go/TS/React patterns
- SDD Process: task completion enforcement
- Service-specific: API contracts, isolation
```

### sdd-checklist Skill Update

Remove enforcement items, keep guidance:
- Move MUST-DO list to `rules/04-sdd-process.md`
- Keep detailed how-to, examples, templates in skill

---

## Testing

Verify rules load correctly:
1. Start Claude Code in project
2. Ask "What rules are loaded?"
3. Verify all 7 rules appear
4. Test path-specific rules by working in service directories
