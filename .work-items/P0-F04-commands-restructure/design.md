# Design Document: Commands/Prompts Restructure

**Feature ID**: P0-F04
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

Convert `.claude/prompts/` templates to proper slash commands in `.claude/commands/subagent/`.

---

## Structure Change

**Before**:
```
.claude/prompts/
├── cursor-sim-cli-dev-template.md
├── analytics-core-dev-template.md
└── viz-spa-dev-template.md
```

**After**:
```
.claude/commands/subagent/
├── cursor-sim-cli.md
├── analytics-core.md
└── viz-spa.md
```

---

## Command Format

### Example: cursor-sim-cli.md

```markdown
---
description: Spawn cursor-sim-cli-dev subagent with scope constraints
argument-hint: [feature-id] [task-id]
allowed-tools: Task
---

# Spawn cursor-sim-cli-dev Subagent

**Feature**: $1
**Task**: $2

## Objective

Implement the specified task following TDD within CLI scope only.

## Scope Constraints

See `.claude/rules/service/cursor-sim.md` for always-on enforcement.

**Additional for this task**:
- Work ONLY on CLI code in `services/cursor-sim/internal/cli/` and `cmd/simulator/`
- Follow `go-best-practices` skill
- Update `.work-items/$1/task.md` with progress

## Context Files

- `.work-items/$1/user-story.md`
- `.work-items/$1/design.md`
- `.work-items/$1/task.md`
- `services/cursor-sim/SPEC.md`

## SDD Workflow

1. SPEC: Read requirements
2. TEST: Write failing tests (RED)
3. CODE: Minimal implementation (GREEN)
4. REFACTOR: Clean up
5. REFLECT: Check dependencies
6. SYNC: Update SPEC.md if needed
7. COMMIT: Include all changes

## Completion Report

Report in this format:
```
TASK COMPLETE: $2
Status: PASSED
Commit: {hash}
Tests: {count} passing
```
```

---

## Cleanup

After conversion:
1. Remove `.claude/prompts/` directory
2. Update P0-F01 design.md to reference new location
3. Update any documentation referencing prompts
