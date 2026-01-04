# Task Breakdown: Component Integration Tests

**Feature ID**: P6-F04
**Epic**: P6 - cursor-viz-spa (Testing Enhancement)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Progress Tracker

| Task | Description | Status | Time Est | Time Actual |
|------|-------------|--------|----------|-------------|
| TASK01 | Create test utility (renderWithProviders) | TODO | 0.5h | - |
| TASK02 | Create mock data factory | TODO | 1.0h | - |
| TASK03 | Write Dashboard integration tests | TODO | 2.0h | - |
| TASK04 | Write Teams integration tests | TODO | 1.0h | - |
| TASK05 | Write Developers integration tests | TODO | 1.0h | - |
| TASK06 | Add npm scripts for integration tests | TODO | 0.25h | - |
| TASK07 | Update test documentation | TODO | 0.5h | - |

**Total Estimated**: 6.25 hours

---

## Task Details

### TASK01: Create Test Utility

**Objective**: Create renderWithProviders helper for integration tests.

**File**: `src/test-utils/integration.tsx`

**Features**:
- Wrap components with MockedProvider
- Wrap with MemoryRouter
- Accept custom mocks and route

**Acceptance Criteria**:
- [ ] Helper created
- [ ] Works with any page component
- [ ] TypeScript types correct

---

### TASK02: Create Mock Data Factory

**Objective**: Create reusable mock data matching P5 schema.

**File**: `src/test-utils/mocks.ts`

**Mock Types**:
- Dashboard success mock
- Dashboard error mock
- Dashboard loading mock
- Teams mocks
- Developers mocks

**Acceptance Criteria**:
- [ ] All mock data typed with generated types
- [ ] Success, error, loading variants
- [ ] Matches P5 schema exactly

---

### TASK03: Write Dashboard Integration Tests

**Objective**: Test Dashboard page integration.

**File**: `src/pages/__tests__/Dashboard.integration.test.tsx`

**Test Cases**:
1. KPI cards display correct data
2. VelocityHeatmap renders with data
3. TeamRadarChart renders with data
4. DeveloperTable renders with data
5. Loading state displays loading indicator
6. Error state displays error message
7. Accessibility (ARIA labels present)

**Acceptance Criteria**:
- [ ] All test cases pass
- [ ] 80%+ coverage of Dashboard page
- [ ] Tests run in < 5 seconds

---

### TASK04: Write Teams Integration Tests

**Objective**: Test Teams page integration.

**File**: `src/pages/__tests__/Teams.integration.test.tsx`

**Test Cases**:
1. Team list renders with data
2. Team details display correctly
3. Loading and error states
4. Team filtering works

**Acceptance Criteria**:
- [ ] All test cases pass
- [ ] 80%+ coverage of Teams page

---

### TASK05: Write Developers Integration Tests

**Objective**: Test Developers page integration.

**File**: `src/pages/__tests__/Developers.integration.test.tsx`

**Test Cases**:
1. Developer list renders with data
2. Pagination works
3. Search/filter works
4. Loading and error states

**Acceptance Criteria**:
- [ ] All test cases pass
- [ ] 80%+ coverage of Developers page

---

### TASK06: Add npm Scripts

**Objective**: Add scripts to run integration tests separately.

**Scripts**:
```json
{
  "scripts": {
    "test:integration": "vitest run --include '**/*.integration.test.*'"
  }
}
```

**Acceptance Criteria**:
- [ ] Script added
- [ ] `npm run test:integration` runs only integration tests
- [ ] Integration tests excluded from regular `npm test` (optional)

---

### TASK07: Update Documentation

**Objective**: Document integration testing approach.

**Files to Update**:
- `services/cursor-viz-spa/README.md`
- Create `src/test-utils/README.md`

**Content**:
- How to write integration tests
- How to use mock data factory
- Test naming conventions

**Acceptance Criteria**:
- [ ] README updated
- [ ] Example tests documented
- [ ] Mock data usage documented

---

## Dependencies

| Task | Depends On |
|------|------------|
| TASK02 | TASK01 |
| TASK03 | TASK02 |
| TASK04 | TASK02 |
| TASK05 | TASK02 |
| TASK06 | TASK03-05 |
| TASK07 | All previous |

---

## Notes

- Phase 1.1 of E2E testing strategy
- Mock data must align with P5 schema (use generated types from P6-F02)
- Integration tests are slower than unit tests but faster than E2E
