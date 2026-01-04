# User Story: Hooks Configuration

**Feature ID**: P0-F05
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: HIGH

---

## User Story

**As a** developer using Claude Code,
**I want** hooks automatically configured and running,
**So that** SDD reminders, markdown formatting, and commit validation happen automatically.

---

## Background

Hook scripts exist in `.claude/hooks/` but are NOT configured:
- `pre_commit.py` - SDD reminder on git commit
- `markdown_formatter.py` - Auto-format markdown files
- `sdd_reminder.py` - Post-task workflow reminder

The `hooks/README.md` says "Needs `/hooks` setup" for all three.

---

## Acceptance Criteria

### AC-01: Hooks Configured
- **Given** Claude Code starts
- **When** hooks are loaded
- **Then** all three hooks are active

### AC-02: Pre-Commit Hook Works
- **Given** a git commit command
- **When** Claude runs it
- **Then** SDD reminder is shown

### AC-03: Markdown Formatter Works
- **Given** Edit/Write on a .md file
- **When** operation completes
- **Then** markdown is auto-formatted

### AC-04: SDD Reminder Works
- **Given** Claude finishes responding
- **When** Stop event fires
- **Then** SDD checklist reminder is shown

---

## Dependencies

- None (scripts already exist)
