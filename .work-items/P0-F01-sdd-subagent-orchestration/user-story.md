# User Story: SDD Subagent Orchestration Protocol

**Feature ID**: P0-F01
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED

---

## User Story

**As a** development team using Claude Code with SDD workflow,
**I want** a well-defined protocol for how master agents and subagents coordinate work,
**So that** parallel development is efficient, code quality is maintained, and documentation stays current.

---

## Background

Currently, Claude Code subagents operate independently but lack a formal protocol for:
- How subagents report task completion
- When and how the master agent performs code review
- How E2E testing catches cross-subagent integration issues
- When documentation (DEVELOPMENT.md, plan folder) gets updated

This leads to:
- Inconsistent task.md updates across features
- Delayed code reviews discovering issues late
- Cross-service integration bugs slipping through
- Stale documentation

---

## Acceptance Criteria

### AC-01: Subagent Task Update Protocol
- **Given** a subagent is working on a feature task
- **When** the subagent completes a task (tests pass + commit)
- **Then** the subagent MUST update only the feature-level `task.md` with status and notes

### AC-02: Subagent Completion Reporting
- **Given** a subagent has completed its task and updated task.md
- **When** the subagent finishes its work session
- **Then** the subagent MUST report completion status to the master agent
- **And** include any blockers, dependencies, or follow-up items

### AC-03: Master Agent Code Review
- **Given** all active subagents have reported completion
- **When** the master agent reviews the changes
- **Then** the master agent MUST:
  - Review code quality across all changed files
  - Check for cross-service consistency
  - Create new tasks if issues found
  - Delegate issue resolution to appropriate subagent

### AC-04: Master Agent E2E Testing
- **Given** code review passes or issues are resolved
- **When** the master agent performs E2E testing
- **Then** the master agent MUST:
  - Run full stack integration tests
  - Fix cross-subagent issues directly (to avoid coordination overhead)
  - Document any E2E fixes made

### AC-05: Documentation Update Protocol
- **Given** E2E tests pass
- **When** the master agent finalizes the work session
- **Then** the master agent MUST:
  - Update DEVELOPMENT.md with session summary
  - Update plan folder (active symlink if applicable)
  - Commit all remaining uncommitted changes

---

## Out of Scope

- Subagent creation/configuration (handled by Claude Code system)
- Individual task implementation details
- Service-specific testing strategies

---

## Dependencies

- Claude Code subagent infrastructure
- Existing SDD workflow skills (spec-process-core, sdd-checklist)
- Work item folder structure

---

## Priority

**Priority**: HIGH
**Rationale**: Foundation for all parallel development. Without clear protocol, subagent coordination is ad-hoc and error-prone.

---

## Notes

This protocol is documentation-only (no code implementation). It defines how Claude Code agents should behave during parallel development sessions.
