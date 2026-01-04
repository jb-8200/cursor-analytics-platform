# Task Breakdown: E2E Testing with Playwright

**Feature ID**: P6-F05
**Epic**: P6 - cursor-viz-spa (Testing Enhancement)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Progress Tracker

| Task | Description | Status | Time Est | Time Actual |
|------|-------------|--------|----------|-------------|
| TASK01 | Install and configure Playwright | TODO | 0.5h | - |
| TASK02 | Create test fixtures and helpers | TODO | 0.5h | - |
| TASK03 | Write Dashboard E2E tests | TODO | 2.0h | - |
| TASK04 | Write Navigation E2E tests | TODO | 1.0h | - |
| TASK05 | Write Error Handling tests | TODO | 1.0h | - |
| TASK06 | Create visual regression baselines | TODO | 1.0h | - |
| TASK07 | Configure CI workflow | TODO | 1.0h | - |
| TASK08 | Document E2E testing workflow | TODO | 0.5h | - |

**Total Estimated**: 7.5 hours

---

## Task Details

### TASK01: Install and Configure Playwright

**Objective**: Set up Playwright testing framework.

**Commands**:
```bash
cd services/cursor-viz-spa
npm install -D @playwright/test
npx playwright install
```

**Files to Create**:
- `playwright.config.ts`
- `tests/e2e/` directory

**Acceptance Criteria**:
- [ ] Playwright installed
- [ ] Browsers installed
- [ ] Config file created
- [ ] `npx playwright test` runs (may fail initially)

---

### TASK02: Create Test Fixtures and Helpers

**Objective**: Set up reusable test utilities.

**Files**:
- `tests/e2e/fixtures/test-data.ts`
- `tests/e2e/fixtures/index.ts`

**Helpers**:
- Service health check
- Wait for data load
- Screenshot helpers

**Acceptance Criteria**:
- [ ] Fixtures created
- [ ] Health check helper works
- [ ] Reusable in all test files

---

### TASK03: Write Dashboard E2E Tests

**Objective**: Test Dashboard page with real data.

**File**: `tests/e2e/dashboard.spec.ts`

**Test Cases**:
1. Page loads successfully
2. GraphQL request completes
3. KPI cards display data
4. Charts render (not placeholders)
5. Data matches GraphQL response

**Acceptance Criteria**:
- [ ] All tests pass with P5 running
- [ ] Tests verify real data display
- [ ] Timeout handling for slow loads

---

### TASK04: Write Navigation E2E Tests

**Objective**: Test page navigation.

**File**: `tests/e2e/navigation.spec.ts`

**Test Cases**:
1. Navigate from Dashboard to Teams
2. Navigate from Dashboard to Developers
3. Back button works
4. Direct URL access works

**Acceptance Criteria**:
- [ ] All navigation tests pass
- [ ] URL changes verified
- [ ] Page content loads after navigation

---

### TASK05: Write Error Handling Tests

**Objective**: Test error states.

**File**: `tests/e2e/error-handling.spec.ts`

**Test Cases**:
1. GraphQL error displays message
2. Network error handled gracefully
3. Empty data state
4. Page doesn't crash on error

**Acceptance Criteria**:
- [ ] Error messages display correctly
- [ ] No JavaScript errors in console
- [ ] App remains usable after error

---

### TASK06: Create Visual Regression Baselines

**Objective**: Create baseline screenshots.

**File**: `tests/e2e/visual.spec.ts`

**Baselines**:
- Dashboard full page
- KPI cards section
- VelocityHeatmap
- TeamRadarChart

**Commands**:
```bash
npx playwright test visual.spec.ts --update-snapshots
```

**Acceptance Criteria**:
- [ ] Baseline screenshots created
- [ ] Screenshots stored in version control
- [ ] Comparison works on subsequent runs

---

### TASK07: Configure CI Workflow

**Objective**: Run E2E tests in CI.

**File**: `.github/workflows/e2e.yml`

**Requirements**:
- Start P4 and P5 in Docker
- Wait for services to be ready
- Run Playwright tests
- Upload reports as artifacts

**Acceptance Criteria**:
- [ ] Workflow file created
- [ ] Services start successfully in CI
- [ ] Tests run and report results
- [ ] Artifacts uploaded

---

### TASK08: Document E2E Testing

**Objective**: Document E2E testing workflow.

**Files to Update**:
- `services/cursor-viz-spa/README.md`
- Create `tests/e2e/README.md`

**Content**:
- How to run E2E tests locally
- Prerequisites (P4, P5 running)
- How to update baselines
- Troubleshooting guide

**Acceptance Criteria**:
- [ ] README updated
- [ ] Local setup documented
- [ ] CI workflow documented

---

## Dependencies

| Task | Depends On |
|------|------------|
| TASK02 | TASK01 |
| TASK03 | TASK02 |
| TASK04 | TASK02 |
| TASK05 | TASK02 |
| TASK06 | TASK03-05 (pages work) |
| TASK07 | TASK06 |
| TASK08 | All previous |

---

## Prerequisites

- P4 (cursor-sim) running on port 8080
- P5 (cursor-analytics-core) running on port 4000
- Database seeded with test data

**Startup Script**:
```bash
# Terminal 1: P4
cd services/cursor-sim && docker run -p 8080:8080 cursor-sim:latest

# Terminal 2: P5
cd services/cursor-analytics-core && docker-compose up

# Terminal 3: P6 + Tests
cd services/cursor-viz-spa && npm run test:e2e
```

---

## Notes

- Phase 2 of E2E testing strategy
- E2E tests are slower than unit/integration tests
- Run only on PR changes to P5/P6
- Visual regression can be flaky; use stable selectors
