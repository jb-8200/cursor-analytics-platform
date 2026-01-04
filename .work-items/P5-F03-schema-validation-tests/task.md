# Task Breakdown: Schema Validation Tests

**Feature ID**: P5-F03
**Epic**: P5 - cursor-analytics-core (Testing Enhancement)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Progress Tracker

| Task | Description | Status | Time Est | Time Actual |
|------|-------------|--------|----------|-------------|
| TASK01 | Create schema validation test file | TODO | 1.0h | - |
| TASK02 | Add type existence tests | TODO | 0.5h | - |
| TASK03 | Add field structure tests | TODO | 1.0h | - |
| TASK04 | Add query validation tests | TODO | 0.5h | - |
| TASK05 | Add breaking change prevention tests | TODO | 0.5h | - |
| TASK06 | Add npm script for schema tests | TODO | 0.25h | - |
| TASK07 | Document schema testing approach | TODO | 0.5h | - |

**Total Estimated**: 4.25 hours

---

## Task Details

### TASK01: Create Schema Validation Test File

**Objective**: Set up test file structure for schema validation.

**File**: `src/graphql/__tests__/schema.validation.test.ts`

**Setup**:
- Import schema
- Import graphql for query execution
- Set up describe blocks

**Acceptance Criteria**:
- [ ] Test file created
- [ ] Schema imported successfully
- [ ] Test file runs without errors

---

### TASK02: Add Type Existence Tests

**Objective**: Verify all required types exist in schema.

**Types to Test**:
- Query
- Developer
- DeveloperStats
- DailyStats
- TeamStats
- DashboardKPI
- Commit
- PageInfo
- DeveloperConnection
- CommitConnection

**Acceptance Criteria**:
- [ ] All type tests pass
- [ ] Test fails if type removed

---

### TASK03: Add Field Structure Tests

**Objective**: Verify critical fields have correct types.

**Fields to Test**:
- TeamStats.topPerformer (Developer, not [Developer])
- DailyStats.linesAdded (Int)
- Developer.stats (DeveloperStats)
- DashboardKPI.teamComparison ([TeamStats])

**Acceptance Criteria**:
- [ ] All field tests pass
- [ ] Singular vs array distinction tested
- [ ] Test fails if field type changed

---

### TASK04: Add Query Validation Tests

**Objective**: Verify P6's queries execute against schema.

**Queries to Test**:
- dashboardSummary (full query from P6)
- developers (with pagination)
- commits (with filters)
- teamStats

**Acceptance Criteria**:
- [ ] All query tests pass
- [ ] Queries match P6's actual queries

---

### TASK05: Add Breaking Change Prevention Tests

**Objective**: Prevent specific known breaking changes.

**Changes to Prevent**:
- topPerformers (plural) should not exist
- humanLinesAdded should not exist
- aiLinesDeleted should not exist

**Acceptance Criteria**:
- [ ] Tests verify fields don't exist
- [ ] Test fails if deprecated field re-added

---

### TASK06: Add npm Script

**Objective**: Add script to run schema tests separately.

**Script**:
```json
{
  "scripts": {
    "test:schema": "vitest run --include '**/schema.*.test.ts'"
  }
}
```

**Acceptance Criteria**:
- [ ] Script added
- [ ] `npm run test:schema` works

---

### TASK07: Document Schema Testing

**Objective**: Document schema testing approach.

**Files to Update**:
- `services/cursor-analytics-core/README.md`
- Add comments to test file explaining purpose

**Acceptance Criteria**:
- [ ] README updated
- [ ] Test file has clear comments

---

## Dependencies

| Task | Depends On |
|------|------------|
| TASK02 | TASK01 |
| TASK03 | TASK01 |
| TASK04 | TASK01 |
| TASK05 | TASK01 |
| TASK06 | TASK01-05 |
| TASK07 | All previous |

---

## Notes

- Phase 1.2 of E2E testing strategy
- These tests run against schema definition, not live server
- Critical for preventing P5+P6 integration issues
