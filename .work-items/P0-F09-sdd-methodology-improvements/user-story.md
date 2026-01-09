# User Story: SDD Methodology Improvements

**Feature ID**: P0-F09-sdd-methodology-improvements
**Created**: January 9, 2026
**Status**: DOCUMENTING (Retroactive)
**Type**: Process Documentation

---

## As a/I want/So that

**As a** development team using AI-assisted coding,
**I want** a well-documented subagent orchestration system with clear SDD workflows,
**So that** parallel development is efficient, scope is protected, and methodology is consistently followed.

---

## Context

This feature retroactively documents all improvements made to the subagent orchestration system and Spec-Driven Development (SDD) methodology during the January 9, 2026 development session.

---

## Key Improvements Made

### 1. Subagent Orchestration

- **planning-dev Agent (Opus)**: Research and design specialist that creates work items
- **cursor-sim-api-dev Agent (Sonnet)**: NEW - Backend specialist for models, generators, API, storage
- **cursor-sim-cli-dev Agent (Sonnet)**: CLI specialist with clear scope separation from api-dev
- **data-tier-dev Agent (Sonnet)**: Python/dbt ETL specialist
- **streamlit-dev Agent (Sonnet)**: Dashboard specialist
- **quick-fix Agent (Haiku)**: Fast fixes for simple tasks

### 2. Rules Compliance

- All agents now explicitly follow `.claude/rules/` directory
- Security, repo guardrails, coding standards, SDD process enforced

### 3. Question Escalation Protocol

- Subagents ask orchestrator when unclear (don't assume)
- Orchestrator relays questions to user
- Answers passed back through orchestrator
- Clarifications documented in design.md or task.md

### 4. CLI Delegation Pattern

- Non-CLI agents cannot run shell commands in background mode
- Agents escalate CLI needs to orchestrator
- Orchestrator delegates to quick-fix (simple) or infra-dev (complex)
- Maintains context integrity across agents

### 5. Scope Separation

- cursor-sim split into API (backend) and CLI (frontend) agents
- Clear file ownership prevents conflicts
- Parallel development without stepping on each other

---

## Acceptance Criteria

### Documentation
- [x] All agent definitions documented in `.claude/agents/`
- [x] README.md updated with agent table and orchestration model
- [x] CLI delegation pattern documented
- [x] Question escalation protocol documented
- [ ] This P0-F09 work item documents all improvements

### Agent Files
- [x] planning-dev.md with rules compliance and escalation
- [x] cursor-sim-api-dev.md (NEW) with backend scope
- [x] cursor-sim-cli-dev.md with CLI-only scope
- [x] data-tier-dev.md with CLI delegation
- [x] streamlit-dev.md with CLI delegation

### Process
- [x] SDD 7-step workflow enforced (SPEC-TEST-CODE-REFACTOR-REFLECT-SYNC-COMMIT)
- [x] dependency-reflection skill integrated
- [x] spec-sync-check skill integrated
- [x] sdd-checklist skill required before commits
