# Task Breakdown: Rules Layer Implementation

**Feature ID**: P0-F02
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Progress Tracker

| Task | Description | Status | Time Est | Time Actual |
|------|-------------|--------|----------|-------------|
| TASK01 | Create rules directory and README | TODO | 0.25h | - |
| TASK02 | Create 01-security.md | TODO | 0.25h | - |
| TASK03 | Create 02-repo-guardrails.md | TODO | 0.25h | - |
| TASK04 | Create 03-coding-standards.md | TODO | 0.5h | - |
| TASK05 | Create 04-sdd-process.md (extract from skill) | TODO | 0.5h | - |
| TASK06 | Create service/cursor-sim.md | TODO | 0.25h | - |
| TASK07 | Create service/analytics-core.md | TODO | 0.25h | - |
| TASK08 | Create service/viz-spa.md | TODO | 0.25h | - |
| TASK09 | Update sdd-checklist skill (remove enforcement) | TODO | 0.25h | - |
| TASK10 | Verify rules load correctly | TODO | 0.25h | - |

**Total Estimated**: 3.0 hours

---

## Task Details

### TASK01: Create Rules Directory and README

**Objective**: Set up rules directory structure with documentation.

**Deliverables**:
- Create `.claude/rules/` directory
- Create `.claude/rules/service/` subdirectory
- Create `.claude/rules/README.md` explaining rules

**Acceptance Criteria**:
- [ ] Directory structure exists
- [ ] README explains when to use rules vs skills
- [ ] README lists all available rules

---

### TASK02: Create Security Rules

**Objective**: Define always-on security guardrails.

**Deliverables**:
- `.claude/rules/01-security.md`

**Content**:
- NEVER: secrets, destructive ops, PII, untrusted scripts
- ALWAYS: validate paths, use env vars, sandbox mode

---

### TASK03: Create Repository Guardrails

**Objective**: Define repository safety rules.

**Deliverables**:
- `.claude/rules/02-repo-guardrails.md`

**Content**:
- NEVER: modify outside root, force push, skip hooks
- ALWAYS: absolute paths, confirm bulk ops, preserve encoding

---

### TASK04: Create Coding Standards

**Objective**: Define shared coding standards.

**Deliverables**:
- `.claude/rules/03-coding-standards.md`

**Content**:
- Go standards (gofmt, error handling, tests)
- TypeScript standards (strict, explicit types)
- React standards (functional, Tailwind, a11y)

---

### TASK05: Create SDD Process Rules

**Objective**: Extract enforcement from sdd-checklist skill.

**Deliverables**:
- `.claude/rules/04-sdd-process.md`

**Content**:
- Post-task MUST-DO list
- NEVER proceed without completing steps

**Note**: This extracts enforcement; guidance stays in skill.

---

### TASK06: Create cursor-sim Rules

**Objective**: Service-specific rules for P4.

**Deliverables**:
- `.claude/rules/service/cursor-sim.md`

**Content**:
- API contract protection
- CLI isolation constraints
- paths: services/cursor-sim/**

---

### TASK07: Create analytics-core Rules

**Objective**: Service-specific rules for P5.

**Deliverables**:
- `.claude/rules/service/analytics-core.md`

**Content**:
- GraphQL schema backward compatibility
- Prisma database access
- paths: services/cursor-analytics-core/**

---

### TASK08: Create viz-spa Rules

**Objective**: Service-specific rules for P6.

**Deliverables**:
- `.claude/rules/service/viz-spa.md`

**Content**:
- GraphQL codegen usage
- Component testing requirements
- paths: services/cursor-viz-spa/**

---

### TASK09: Update sdd-checklist Skill

**Objective**: Remove enforcement, keep guidance.

**Deliverables**:
- Updated `.claude/skills/sdd-checklist/SKILL.md`

**Changes**:
- Remove MUST-DO enforcement (now in rules)
- Keep detailed how-to, examples, templates
- Add reference to rules/04-sdd-process.md

---

### TASK10: Verify Rules Load Correctly

**Objective**: Test that rules are discovered and applied.

**Test Steps**:
1. Restart Claude Code
2. Ask "What rules are loaded?"
3. Verify global rules appear
4. Work in services/cursor-sim/ and verify path rules apply
5. Document any issues

---

## Dependencies

- P0-F01 (Subagent Orchestration) - COMPLETE

---

## Risks

| Risk | Mitigation |
|------|------------|
| Rules too restrictive | Start permissive, tighten based on feedback |
| Rules conflict with skills | Clear separation: rules=enforcement, skills=guidance |
| Path matching issues | Test with actual file operations |
