# User Story: Selection Heuristic Guide

**Feature ID**: P0-F07
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: LOW

---

## User Story

**As a** developer configuring Claude Code,
**I want** clear guidance on when to use rules vs commands vs skills vs agents,
**So that** I put instructions in the right place.

---

## Background

docs/claude-docs/ has comprehensive reference docs but:
- No "when to use what" guide specific to this project
- No mapping from guidance type to project location
- New developers don't know where to add instructions

---

## Acceptance Criteria

### AC-01: README Updated
- **Given** .claude/README.md
- **When** I read it
- **Then** selection heuristic table is included

### AC-02: Decision Tree
- **Given** a new instruction to add
- **When** I consult the guide
- **Then** I know which type (rule/command/skill/agent/hook) to use

---

## Dependencies

- P0-F02 through P0-F06 complete (all types exist)
