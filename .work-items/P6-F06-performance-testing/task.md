# Task Breakdown: Performance Testing with Lighthouse

**Feature ID**: P6-F06
**Epic**: P6 - cursor-viz-spa (Testing Enhancement)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Progress Tracker

| Task | Description | Status | Time Est | Time Actual |
|------|-------------|--------|----------|-------------|
| TASK01 | Install Lighthouse CI | TODO | 0.25h | - |
| TASK02 | Create Lighthouse config | TODO | 0.5h | - |
| TASK03 | Define performance budget | TODO | 0.5h | - |
| TASK04 | Add npm scripts | TODO | 0.25h | - |
| TASK05 | Run initial audit and baseline | TODO | 0.5h | - |
| TASK06 | Configure CI workflow | TODO | 1.0h | - |
| TASK07 | Document performance testing | TODO | 0.5h | - |

**Total Estimated**: 3.5 hours

---

## Task Details

### TASK01: Install Lighthouse CI

**Objective**: Add Lighthouse CI CLI.

**Commands**:
```bash
cd services/cursor-viz-spa
npm install -D @lhci/cli
```

**Acceptance Criteria**:
- [ ] Package installed
- [ ] `npx lhci --version` works

---

### TASK02: Create Lighthouse Config

**Objective**: Configure Lighthouse CI.

**File**: `services/cursor-viz-spa/lighthouserc.js`

**Configuration**:
- Collect from /dashboard
- 3 runs for consistency
- Desktop preset
- Filesystem upload

**Acceptance Criteria**:
- [ ] Config file created
- [ ] `npx lhci autorun` works locally

---

### TASK03: Define Performance Budget

**Objective**: Set performance thresholds.

**File**: `services/cursor-viz-spa/lighthouse-budget.json`

**Thresholds**:
- FCP: < 2000ms
- TTI: < 5000ms
- Speed Index: < 3500ms
- Performance Score: > 80

**Acceptance Criteria**:
- [ ] Budget file created
- [ ] Thresholds are achievable for current build

---

### TASK04: Add npm Scripts

**Objective**: Add Lighthouse scripts.

**Scripts**:
```json
{
  "scripts": {
    "lighthouse": "lhci autorun",
    "preview": "vite preview --port 4173"
  }
}
```

**Acceptance Criteria**:
- [ ] Scripts added
- [ ] `npm run lighthouse` works after build

---

### TASK05: Run Initial Audit

**Objective**: Establish performance baseline.

**Steps**:
1. `npm run build`
2. `npm run lighthouse`
3. Review results
4. Adjust thresholds if needed

**Acceptance Criteria**:
- [ ] Initial audit complete
- [ ] Results saved in lighthouse-results/
- [ ] Baseline performance documented

---

### TASK06: Configure CI Workflow

**Objective**: Run Lighthouse in CI.

**File**: `.github/workflows/lighthouse.yml`

**Requirements**:
- Build P6
- Run Lighthouse audit
- Assert against budget
- Upload results

**Acceptance Criteria**:
- [ ] Workflow file created
- [ ] CI runs successfully
- [ ] Results uploaded as artifacts

---

### TASK07: Document Performance Testing

**Objective**: Document performance testing.

**Files to Update**:
- `services/cursor-viz-spa/README.md`

**Content**:
- How to run Lighthouse locally
- Performance targets
- How to interpret results
- How to update budget

**Acceptance Criteria**:
- [ ] README updated
- [ ] Performance targets documented

---

## Dependencies

| Task | Depends On |
|------|------------|
| TASK02 | TASK01 |
| TASK03 | TASK02 |
| TASK04 | TASK02 |
| TASK05 | TASK03, TASK04 |
| TASK06 | TASK05 |
| TASK07 | All previous |

---

## Notes

- Phase 4 of E2E testing strategy (lowest priority)
- Tests production build, not dev server
- Can be run independently of P5 (static content)
- Consider adding to pre-deploy workflow
