# Task Breakdown: SDD Subagent Orchestration Protocol

**Feature ID**: P0-F01
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Progress Tracker

| Task | Description | Status | Time Est | Time Actual |
|------|-------------|--------|----------|-------------|
| TASK01 | Define subagent task.md update guidelines | COMPLETE | 0.5h | 0.25h |
| TASK02 | Define subagent completion reporting protocol | COMPLETE | 0.5h | 0h* |
| TASK03 | Define master agent code review workflow | COMPLETE | 0.5h | 0h* |
| TASK04 | Define master agent E2E testing protocol | COMPLETE | 0.5h | 0h* |
| TASK05 | Define documentation update protocol | COMPLETE | 0.5h | 0h* |
| TASK06 | Update CLAUDE.md with protocol reference | COMPLETE | 0.5h | 0.25h |
| TASK07 | Create agent prompt templates | COMPLETE | 1.0h | 0.5h |

**Total Estimated**: 4.0 hours
**Total Actual**: 1.0 hours (Tasks 02-05 were pre-documented in design.md)

---

## Task Details

### TASK01: Define Subagent task.md Update Guidelines

**Objective**: Create clear guidelines for how subagents update task.md files.

**Deliverables**:
- Document in design.md: task.md update format
- Define required fields for status updates
- Define timing (when to update)

**Acceptance Criteria**:
- [ ] Format template documented
- [ ] Required vs optional fields defined
- [ ] Update triggers clearly specified

---

### TASK02: Define Subagent Completion Reporting Protocol

**Objective**: Standardize how subagents report task completion to master agent.

**Deliverables**:
- Completion report format/template
- Status categories (PASSED, BLOCKED, PARTIAL)
- Required information for each status

**Acceptance Criteria**:
- [ ] Report format documented
- [ ] All status types defined
- [ ] Example reports provided

---

### TASK03: Define Master Agent Code Review Workflow

**Objective**: Document the code review process master agent performs.

**Deliverables**:
- Code review checklist
- Issue categorization (minor, cross-service, blocking)
- Escalation protocol for each category

**Acceptance Criteria**:
- [ ] Checklist complete and actionable
- [ ] Issue categories with examples
- [ ] Clear escalation paths

---

### TASK04: Define Master Agent E2E Testing Protocol

**Objective**: Document E2E testing and issue resolution workflow.

**Deliverables**:
- E2E test scope definition
- Issue resolution protocol (master fixes directly)
- Documentation requirements for E2E fixes

**Acceptance Criteria**:
- [ ] Test categories defined
- [ ] Resolution workflow documented
- [ ] Fix documentation format specified

---

### TASK05: Define Documentation Update Protocol

**Objective**: Document when and how master agent updates project documentation.

**Deliverables**:
- DEVELOPMENT.md update template
- Plan folder management rules
- Final commit checklist

**Acceptance Criteria**:
- [ ] Update template documented
- [ ] Symlink management rules clear
- [ ] Commit checklist complete

---

### TASK06: Update CLAUDE.md with Protocol Reference

**Objective**: Add reference to this protocol in CLAUDE.md.

**Deliverables**:
- Add Subagent Orchestration section to CLAUDE.md
- Link to P0-F01 design document
- Summary of key rules

**Acceptance Criteria**:
- [ ] CLAUDE.md updated
- [ ] Key rules summarized
- [ ] Link to full protocol included

---

### TASK07: Create Agent Prompt Templates

**Objective**: Create reusable prompt templates for spawning subagents.

**Deliverables**:
- Template for P4 (cursor-sim-cli-dev) tasks
- Template for P5 (analytics-core-dev) tasks
- Template for P6 (viz-spa-dev) tasks
- Include protocol reminders in templates

**Acceptance Criteria**:
- [ ] All 3 templates created
- [ ] Templates include scope constraints
- [ ] Templates include reporting requirements

---

## Completion Summary

### TASK06: Update CLAUDE.md with Protocol Reference (January 4, 2026)

**Status**: COMPLETE
**Time**: 0.25h / 0.5h

**Completed**:
- Added "Subagent Orchestration Protocol" section to CLAUDE.md
- Documented master agent responsibilities
- Documented subagent constraints
- Added completion flow diagram
- Linked to full protocol in P0-F01 design.md

**Changes**:
- Modified: CLAUDE.md (lines 93-120)

---

### TASK07: Create Agent Prompt Templates (January 4, 2026)

**Status**: COMPLETE
**Time**: 0.5h / 1.0h

**Completed**:
- Created `.claude/prompts/` directory
- Created cursor-sim-cli-dev-template.md (CLI-only scope)
- Created analytics-core-dev-template.md (GraphQL/TypeScript)
- Created viz-spa-dev-template.md (React/Vite/Apollo Client)
- All templates include scope constraints and protocol reminders

**Changes**:
- Added: .claude/prompts/cursor-sim-cli-dev-template.md
- Added: .claude/prompts/analytics-core-dev-template.md
- Added: .claude/prompts/viz-spa-dev-template.md

---

## Notes

- This is a documentation-only feature (no code implementation)
- All tasks produce documentation artifacts
- Protocol is now active and referenced in CLAUDE.md

---

## Dependencies

None - this is foundational documentation.

---

## Risks

| Risk | Mitigation |
|------|------------|
| Agents may not follow protocol | Include reminders in prompts |
| Protocol may be too rigid | Allow for reasonable exceptions |
| Documentation may become stale | Review protocol quarterly |
