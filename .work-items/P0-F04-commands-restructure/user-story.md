# User Story: Commands/Prompts Restructure

**Feature ID**: P0-F04
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: MEDIUM

---

## User Story

**As a** developer spawning subagents,
**I want** proper slash commands with frontmatter for subagent templates,
**So that** I can invoke them with `/subagent/cursor-sim-cli P4-F02 TASK07`.

---

## Background

Current `.claude/prompts/` templates:
- Have `{placeholders}` instead of `$ARGUMENTS`
- No frontmatter (description, argument-hint, allowed-tools)
- Not integrated with Claude Code command system
- Require manual copy/paste

---

## Acceptance Criteria

### AC-01: Commands Structure
- **Given** prompts directory
- **When** restructured
- **Then** `.claude/commands/subagent/` contains proper commands

### AC-02: Frontmatter
- **Given** a subagent command
- **When** loaded by Claude Code
- **Then** it has description, argument-hint, allowed-tools

### AC-03: Arguments Work
- **Given** `/subagent/cursor-sim-cli P4-F02 TASK07`
- **When** command is invoked
- **Then** $1 and $2 are substituted correctly

---

## Dependencies

- P0-F02 (Rules Layer) - service constraints referenced
