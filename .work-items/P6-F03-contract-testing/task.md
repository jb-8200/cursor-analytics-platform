# Task Breakdown: Contract Testing with GraphQL Inspector

**Feature ID**: P6-F03
**Epic**: P6 - cursor-viz-spa (Service Decoupling Phase 3)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Progress Tracker

| Task | Description | Status | Time Est | Time Actual |
|------|-------------|--------|----------|-------------|
| TASK01 | Install GraphQL Inspector | TODO | 0.25h | - |
| TASK02 | Create validation script | TODO | 0.5h | - |
| TASK03 | Add npm scripts | TODO | 0.25h | - |
| TASK04 | Configure pre-commit hook | TODO | 0.5h | - |
| TASK05 | Add CI validation step | TODO | 0.5h | - |
| TASK06 | Test with intentional errors | TODO | 0.5h | - |
| TASK07 | Document validation workflow | TODO | 0.5h | - |

**Total Estimated**: 3.0 hours

---

## Task Details

### TASK01: Install GraphQL Inspector

**Objective**: Add GraphQL Inspector CLI to P6.

**Commands**:
```bash
cd services/cursor-viz-spa
npm install -D @graphql-inspector/cli
```

**Acceptance Criteria**:
- [ ] Package installed
- [ ] `npx graphql-inspector --version` works

---

### TASK02: Create Validation Script

**Objective**: Create shell script for query validation.

**File**: `services/cursor-viz-spa/scripts/validate-queries.sh`

**Features**:
- Check if P5 is running
- Run graphql-inspector validate
- Clear error messages

**Acceptance Criteria**:
- [ ] Script created and executable
- [ ] Fails gracefully if P5 not running
- [ ] Validates all .ts files with GraphQL queries

---

### TASK03: Add npm Scripts

**Objective**: Add schema validation scripts to package.json.

**Scripts**:
```json
{
  "scripts": {
    "schema:validate": "bash scripts/validate-queries.sh"
  }
}
```

**Acceptance Criteria**:
- [ ] `npm run schema:validate` works
- [ ] Clear output on success and failure

---

### TASK04: Configure Pre-Commit Hook

**Objective**: Block commits with invalid queries.

**Implementation**:
- Use existing Husky setup (from P6-F02)
- Only validate if GraphQL files changed

**Acceptance Criteria**:
- [ ] Hook triggers on GraphQL file changes
- [ ] Commit blocked on validation failure
- [ ] Commit proceeds on success
- [ ] Hook skipped if no GraphQL changes

---

### TASK05: Add CI Validation Step

**Objective**: Validate queries in CI pipeline.

**File**: Update `.github/workflows/ci.yml`

**Steps**:
1. Start P5 in background
2. Run validation
3. Report errors in PR

**Acceptance Criteria**:
- [ ] CI step added
- [ ] P5 starts successfully in CI
- [ ] Validation runs
- [ ] Clear error messages in CI logs

---

### TASK06: Test with Intentional Errors

**Objective**: Verify validation catches real errors.

**Test Cases**:
1. Invalid field name (topPerformers → topPerformer)
2. Missing required variable
3. Wrong variable type
4. Invalid fragment

**Acceptance Criteria**:
- [ ] All error types caught
- [ ] Error messages are helpful
- [ ] File and line numbers shown

---

### TASK07: Document Validation Workflow

**Objective**: Document the validation workflow.

**Files to Update**:
- `services/cursor-viz-spa/README.md`
- `docs/INTEGRATION.md`

**Content**:
- How to run validation manually
- What errors look like
- How to fix common errors

**Acceptance Criteria**:
- [ ] README updated
- [ ] Common errors documented
- [ ] Troubleshooting guide added

---

## Dependencies

| Task | Depends On |
|------|------------|
| TASK02 | TASK01 |
| TASK03 | TASK02 |
| TASK04 | TASK03, P6-F02 TASK07 (Husky) |
| TASK05 | TASK03 |
| TASK06 | TASK03 |
| TASK07 | All previous |

---

## Notes

- Phase 3 of data contract testing mitigation plan
- Complements TypeScript type checking (P6-F02)
- Requires P5 to be running for validation
