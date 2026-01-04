# Design Document: Agents Documentation

**Feature ID**: P0-F06
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

Create README.md for agents directory and ensure agents have proper skills mapping.

---

## Agents README

Create `.claude/agents/README.md`:

```markdown
# Custom Subagents

Specialized agents for parallel development with isolated scope.

## Available Agents

| Agent | Service | Phase | Scope | Skills |
|-------|---------|-------|-------|--------|
| cursor-sim-cli-dev | cursor-sim | P4 | CLI only | go-best-practices |
| cursor-sim-infra-dev | cursor-sim | P7 | Docker, Cloud Run | - |
| analytics-core-dev | analytics-core | P5 | GraphQL service | typescript-graphql-patterns, api-contract |
| viz-spa-dev | viz-spa | P6 | React dashboard | react-vite-patterns |

## Scope Constraints

### cursor-sim-cli-dev
- ONLY: `internal/cli/`, `cmd/simulator/`
- NEVER: `internal/api/`, `internal/generator/`
- Protects API contracts for P5/P6

### analytics-core-dev
- ONLY: `services/cursor-analytics-core/`
- Must align with cursor-sim API (api-contract skill)

### viz-spa-dev
- ONLY: `services/cursor-viz-spa/`
- Must align with analytics-core GraphQL schema

## Orchestration

See `.work-items/P0-F01-sdd-subagent-orchestration/design.md`:
- Master agent delegates tasks
- Subagents update task.md only
- Master handles DEVELOPMENT.md and E2E fixes

## Spawning Agents

Use commands in `.claude/commands/subagent/`:
```
/subagent/cursor-sim-cli P4-F02 TASK07
/subagent/analytics-core P5-F02 TASK01
/subagent/viz-spa P6-F02 TASK01
```
```

---

## Agent Files Update

Add `skills:` field to agent definitions if missing:

```yaml
---
name: cursor-sim-cli-dev
skills: go-best-practices, sdd-checklist
---
```
