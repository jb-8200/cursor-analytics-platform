# User Story: Rules Layer Implementation

**Feature ID**: P0-F02
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: CRITICAL

---

## User Story

**As a** developer using Claude Code with SDD workflow,
**I want** always-on guardrails that enforce security, coding standards, and process constraints,
**So that** enforcement happens automatically without relying on skills/prompts.

---

## Background

Currently, always-on constraints are scattered across:
- CLAUDE.md (bloated with workflow AND constraints)
- Skills like `sdd-checklist` (mix enforcement with guidance)
- No path-specific rules (same rules for Go/TS/React)

Claude Code supports `.claude/rules/*.md` with optional `paths:` frontmatter for scoped rules.

---

## Acceptance Criteria

### AC-01: Security Rules
- **Given** any Claude Code session
- **When** Claude operates on the codebase
- **Then** security rules (secrets, destructive ops, PII) are always enforced

### AC-02: Repository Guardrails
- **Given** any file operation
- **When** Claude attempts to modify protected files
- **Then** repo guardrails (git safety, file protection) apply

### AC-03: Coding Standards
- **Given** code being written
- **When** Claude generates code
- **Then** shared coding standards (Go/TS/React) are applied

### AC-04: SDD Process Enforcement
- **Given** a task completion scenario
- **When** Claude finishes a task
- **Then** SDD must-dos are enforced (not just suggested)

### AC-05: Service-Specific Rules
- **Given** work in a specific service directory
- **When** Claude operates in `services/cursor-sim/**`
- **Then** cursor-sim-specific rules apply (API contract protection)

---

## Out of Scope

- Changing how skills work
- Modifying commands
- Hook configuration (separate feature)

---

## Dependencies

- None (foundational feature)

---

## Notes

This is the CRITICAL foundation that enables proper separation of concerns.
Rules = enforcement (always-on). Skills = guidance (when relevant).
