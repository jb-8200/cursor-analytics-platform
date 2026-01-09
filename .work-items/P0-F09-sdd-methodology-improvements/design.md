# Technical Design: SDD Methodology Improvements

**Feature ID**: P0-F09-sdd-methodology-improvements
**Created**: January 9, 2026
**Status**: DOCUMENTING (Retroactive)

---

## Overview

This document captures all technical decisions and patterns established for subagent orchestration and SDD methodology improvements.

---

## Architecture

### Subagent Hierarchy

```
┌─────────────────────────────────────────────────────────────────┐
│                    ORCHESTRATOR (Opus)                          │
│  - Delegates tasks to specialized subagents                     │
│  - Reviews cross-service code quality                           │
│  - Handles E2E testing and fixes                                │
│  - Updates DEVELOPMENT.md                                       │
│  - Relays questions between subagents and user                  │
└─────────────────────────────────────────────────────────────────┘
          │              │              │              │
          ▼              ▼              ▼              ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│ planning-dev│ │cursor-sim-  │ │data-tier-dev│ │streamlit-dev│
│   (Opus)    │ │ api-dev     │ │  (Sonnet)   │ │  (Sonnet)   │
│             │ │  (Sonnet)   │ │             │ │             │
│ Research +  │ │ Models +    │ │ Python ETL  │ │ Dashboard   │
│ Design only │ │ Generators  │ │ + dbt       │ │ + Plotly    │
└─────────────┘ │ + API       │ └─────────────┘ └─────────────┘
                └─────────────┘
                      │
          ┌───────────┴───────────┐
          ▼                       ▼
┌─────────────────┐     ┌─────────────────┐
│cursor-sim-cli-  │     │cursor-sim-infra-│
│     dev         │     │     dev         │
│   (Sonnet)      │     │   (Sonnet)      │
│                 │     │                 │
│ CLI only        │     │ Docker + Cloud  │
│ (cmd/, config/) │     │ Run deployment  │
└─────────────────┘     └─────────────────┘
          │
          ▼
┌─────────────────┐
│   quick-fix     │
│   (Haiku)       │
│                 │
│ Simple tasks,   │
│ CLI commands    │
└─────────────────┘
```

---

## Agent Scope Definitions

### cursor-sim Service Split

The cursor-sim service is split between two complementary agents:

| Agent | Scope | Files |
|-------|-------|-------|
| **cursor-sim-api-dev** | Backend | `internal/models/`, `internal/generator/`, `internal/storage/`, `internal/api/`, `internal/seed/`, `internal/services/` |
| **cursor-sim-cli-dev** | Frontend | `cmd/simulator/`, `internal/config/`, `internal/cli/` |

**Rationale**: Enables parallel development without conflicts. API contracts remain stable for P5/P6 integration.

### Non-CLI Agents

| Agent | Primary Work | Cannot Modify |
|-------|--------------|---------------|
| planning-dev | user-story.md, design.md, task.md | Implementation code |
| data-tier-dev | `tools/api-loader/`, `dbt/` | `services/cursor-sim/` |
| streamlit-dev | `services/streamlit-dashboard/` | Other services |
| analytics-core-dev | `services/cursor-analytics-core/` | Other services |
| viz-spa-dev | `services/cursor-viz-spa/` | Other services |

---

## Communication Protocols

### Question Escalation Protocol

When a subagent encounters ambiguity:

```
1. Subagent → Orchestrator: "QUESTION for orchestrator: [topic, question, impact]"
2. Orchestrator → User: Relays question
3. User → Orchestrator: Provides answer
4. Orchestrator → Subagent: Relays answer (via resume)
5. Subagent: Documents clarification and continues
```

**Example Format**:
```
QUESTION for orchestrator:
- Topic: Microsoft 365 Copilot API authentication
- Question: Should we simulate OAuth token flow or use simplified Bearer auth?
- Impact: Affects complexity of auth handler and test setup
```

### CLI Delegation Pattern

Non-CLI agents cannot run shell commands in background mode. When needed:

```
1. Subagent → Orchestrator: "CLI ACTION NEEDED: [command, purpose, context]"
2. Orchestrator → quick-fix/infra-dev: Delegates CLI task
3. CLI Agent: Executes command, reports result
4. Orchestrator → Subagent: Reports completion
```

**Example Format**:
```
CLI ACTION NEEDED:
- Command: pip install -r requirements.txt && pytest tests/ -v
- Purpose: Validate loader tests pass
- Context: Completed TASK-P8-03, need test verification
```

### Completion Reporting

Standard completion report format:

```
TASK COMPLETE: TASK-{ID}
Status: PASSED
Commit: {hash}
Tests: {count} passing
Coverage: {percent}%

Changes:
- {file list}

Notes: {context for master agent}
```

Blocked report format:

```
TASK BLOCKED: TASK-{ID}
Blocker: {issue description}
Impact: {what cannot be completed}
Needs: {what is needed to unblock}
```

---

## SDD Workflow Integration

### Enhanced 7-Step Workflow

```
1. SPEC     → Read specification before coding
2. TEST     → Write failing tests (RED)
3. CODE     → Minimal implementation (GREEN)
4. REFACTOR → Clean up while tests pass
5. REFLECT  → Check dependency reflections (dependency-reflection skill)
6. SYNC     → Update SPEC.md if triggered (spec-sync-check skill)
7. COMMIT   → Every task = commit (code + docs)
```

### Skills Integration

| Step | Skill | Purpose |
|------|-------|---------|
| 1. SPEC | spec-process-core | Core SDD principles |
| 2. TEST | spec-process-dev | TDD RED-GREEN-REFACTOR |
| 5. REFLECT | dependency-reflection | Check cross-file impacts |
| 6. SYNC | spec-sync-check | Trigger SPEC.md updates |
| 7. COMMIT | sdd-checklist | Pre-commit verification |

### Mandatory Checks

Before every commit, subagents must verify:

- [ ] Tests pass
- [ ] Dependency reflections checked
- [ ] SPEC.md updated if triggered
- [ ] Only in-scope files modified
- [ ] task.md updated with progress

---

## Rules Compliance

All agents follow `.claude/rules/`:

| Rule | Enforcement |
|------|-------------|
| `01-security.md` | Never commit secrets, validate paths |
| `02-repo-guardrails.md` | Git safety, file protection |
| `03-coding-standards.md` | Go/TS/React patterns, testing |
| `04-sdd-process.md` | 7-step workflow enforcement |

---

## Model Selection Strategy

| Model | Use Case | Cost | Speed |
|-------|----------|------|-------|
| **Opus** | Planning, orchestration, code review | Higher | Slower |
| **Sonnet** | Feature implementation, complex tasks | Medium | Medium |
| **Haiku** | Quick fixes, simple CLI tasks | Low | Fast |

**Decision Matrix**:
- Research/design/architecture → Opus (planning-dev)
- Multi-file feature work → Sonnet (service-specific agents)
- Single-file fixes, CLI commands → Haiku (quick-fix)

---

## Files Modified/Created

### Agent Definitions
- `.claude/agents/planning-dev.md` - Added rules compliance, escalation protocol
- `.claude/agents/cursor-sim-api-dev.md` - NEW: Backend specialist
- `.claude/agents/cursor-sim-cli-dev.md` - Clarified CLI-only scope
- `.claude/agents/data-tier-dev.md` - Added CLI delegation
- `.claude/agents/streamlit-dev.md` - Added CLI delegation
- `.claude/agents/README.md` - Updated with all agents, delegation pattern

### Skills
- `.claude/skills/dependency-reflection/` - Cross-file impact checking
- `.claude/skills/spec-sync-check/` - SPEC.md update triggers
- `.claude/skills/sdd-checklist/` - Pre-commit verification

### Rules
- `.claude/rules/04-sdd-process.md` - 7-step workflow enforcement

---

## Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-01-09 | Split cursor-sim into api-dev and cli-dev | Enable parallel development, protect API contracts |
| 2026-01-09 | Add CLI delegation pattern | Background agents can't run shell, need orchestrator delegation |
| 2026-01-09 | Add question escalation protocol | Prevent assumptions, ensure user input on ambiguous requirements |
| 2026-01-09 | Require rules compliance for all agents | Consistent enforcement of security and process |
| 2026-01-09 | Use Opus for planning-dev | Research and design requires higher reasoning capability |
| 2026-01-09 | Add REFLECT and SYNC steps to SDD | Catch dependency drift and keep SPEC.md current |

---

## Future Improvements

### Potential Enhancements
1. **Agent Registry**: Auto-register new agents in Task tool
2. **Progress Dashboard**: Visual tracking of subagent tasks
3. **Automated Testing**: CI integration for agent output validation
4. **Context Compression**: Better handling of long-running agent sessions

### Known Limitations
1. cursor-sim-api-dev not yet registered in Task tool (uses general-purpose)
2. CLI delegation requires manual orchestrator intervention
3. No automatic progress aggregation across agents

---

**Last Updated**: January 9, 2026
