# Task Breakdown: Skills Cleanup & Catalog

**Feature ID**: P0-F03
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Progress Tracker

| Task | Description | Status | Time Est | Time Actual |
|------|-------------|--------|----------|-------------|
| TASK01 | Create skills/README.md catalog | TODO | 0.5h | - |
| TASK02 | Update sdd-checklist (remove enforcement) | TODO | 0.25h | - |
| TASK03 | Audit and fix skill frontmatter | TODO | 0.5h | - |
| TASK04 | Update commands/README.md (fix wrong references) | TODO | 0.25h | - |

**Total Estimated**: 1.5 hours

---

## Task Details

### TASK01: Create Skills Catalog README

**Objective**: Document all skills with triggers and service mapping.

**Deliverables**:
- `.claude/skills/README.md`

**Content**:
- Table of all 14 skills
- Trigger keywords for each
- Service mapping (P4/P5/P6/All)
- Skill vs Rule distinction

---

### TASK02: Update sdd-checklist

**Objective**: Remove enforcement, keep guidance.

**Prerequisite**: P0-F02 TASK05 complete (rules/04-sdd-process.md exists)

**Deliverables**:
- Updated `.claude/skills/sdd-checklist/SKILL.md`

**Changes**:
- Remove MUST/NEVER enforcement language
- Add reference to rules for enforcement
- Keep examples, templates, TodoWrite integration

---

### TASK03: Audit Skill Frontmatter

**Objective**: Ensure consistent frontmatter across all skills.

**Audit Checklist**:
- [ ] name matches directory
- [ ] description has trigger keywords
- [ ] allowed-tools if needed

**Skills to check**: All 14 skills

---

### TASK04: Fix Commands README

**Objective**: Remove references to nonexistent skill groupings.

**Deliverables**:
- Updated `.claude/commands/README.md`

**Changes**:
- Remove references to `process/`, `standards/`, `guidelines/`, `operational/`
- Update to match flat skills structure
- Add correct skill references

---

## Dependencies

- P0-F02 (Rules Layer) - TASK05 must complete before TASK02
