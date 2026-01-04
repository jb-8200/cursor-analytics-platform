# User Story: Agents Documentation

**Feature ID**: P0-F06
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: LOW

---

## User Story

**As a** developer using subagents,
**I want** documentation explaining available agents and their scope,
**So that** I know which agent to use for each task.

---

## Background

`.claude/agents/` has 4 subagent definitions:
- cursor-sim-cli-dev.md
- cursor-sim-infra-dev.md
- analytics-core-dev.md
- viz-spa-dev.md

But no README.md explaining:
- What each agent is for
- What skills they should use
- Scope constraints
- When to use which

---

## Acceptance Criteria

### AC-01: Agents README
- **Given** the agents directory
- **When** I look for documentation
- **Then** README.md explains all agents

### AC-02: Skills Mapping
- **Given** an agent definition
- **When** reviewed
- **Then** skills field lists appropriate skills

---

## Dependencies

- None
