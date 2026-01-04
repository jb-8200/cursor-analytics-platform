# User Story: Skills Cleanup & Catalog

**Feature ID**: P0-F03
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: MEDIUM

---

## User Story

**As a** developer using Claude Code,
**I want** a well-organized skills catalog with clear trigger descriptions,
**So that** I know what skills exist, when they activate, and how they differ from rules.

---

## Background

Current issues:
- 14 skills with no catalog README
- Flat directory (commands README references nonexistent groupings)
- `sdd-checklist` mixes enforcement (now in rules) with guidance
- Inconsistent frontmatter across skills

---

## Acceptance Criteria

### AC-01: Skills Catalog README
- **Given** the skills directory
- **When** I look for documentation
- **Then** README.md lists all skills with triggers and service mapping

### AC-02: sdd-checklist Updated
- **Given** P0-F02 (Rules Layer) is complete
- **When** sdd-checklist is loaded
- **Then** it contains only guidance (enforcement in rules)

### AC-03: Consistent Frontmatter
- **Given** any skill file
- **When** Claude loads the skill
- **Then** it has proper name, description, and optional allowed-tools

---

## Dependencies

- P0-F02 (Rules Layer) - for sdd-checklist split
