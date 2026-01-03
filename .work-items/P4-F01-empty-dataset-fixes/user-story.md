# User Story: Fix Empty Dataset Issues

**Feature ID**: P4-F01-empty-dataset-fixes
**Created**: January 3, 2026
**Status**: ✅ COMPLETE

---

## Executive Summary

As a **SDLC researcher** using cursor-sim to generate synthetic developer data, I want **all API endpoints to return populated data** so that I can **reliably analyze the full range of developer behaviors and metrics** without encountering empty datasets.

---

## User Story (EARS Format)

### Story: Fix Empty Analytics Endpoints

**As a** researcher using cursor-sim to study developer behavior
**I want** all analytics endpoints to return non-empty data
**So that** I can analyze the complete dataset without encountering empty arrays

**Acceptance Criteria**:

```gherkin
Given I start cursor-sim with a valid seed file
When I query any analytics endpoint
Then I receive non-empty data arrays
And all 15 previously-empty endpoints now return data
And the data matches the expected response format
And no endpoint returns {"data": []} or {"data": {}}
```

---

## Problem Statement

### Current Issues

**52% of endpoints broken** (15 out of 29 total):

**Affected Endpoints**:
1. `/teams/members` - Empty team members array
2. `/analytics/team/models` - Empty model usage data
3. `/analytics/by-user/models` - Empty per-user model data
4. `/analytics/team/client-versions` - Empty version data
5. `/analytics/by-user/client-versions` - Empty per-user versions
6. `/analytics/team/top-file-extensions` - Empty extensions
7. `/analytics/by-user/top-file-extensions` - Empty per-user extensions
8. `/analytics/team/mcp` - Empty MCP tool usage
9. `/analytics/by-user/mcp` - Empty per-user MCP
10. `/analytics/team/commands` - Empty command usage
11. `/analytics/by-user/commands` - Empty per-user commands
12. `/analytics/team/plans` - Empty plan usage
13. `/analytics/by-user/plans` - Empty per-user plans
14. `/analytics/team/ask-mode` - Empty ask mode data
15. `/analytics/by-user/ask-mode` - Empty per-user ask mode

### Root Cause

All generators exist and work correctly, but **were never called in main.go**:

```go
// cmd/simulator/main.go - BEFORE P4-F01
// Only commit generator was called:
gen := generator.NewCommitGenerator(seedData, store, cfg.Velocity)
gen.GenerateCommits(ctx, cfg.Days)

// Missing:
// - store.LoadDevelopers()
// - ModelGenerator.GenerateModelUsage()
// - VersionGenerator.GenerateClientVersions()
// - ExtensionGenerator.GenerateFileExtensions()
// - FeatureGenerator.GenerateFeatures()
```

---

## Goals & Non-Goals

### Goals

- ✅ Fix all 15 empty endpoints
- ✅ Call all 5 generators in proper sequence
- ✅ Create comprehensive E2E tests for each generator
- ✅ Create integration test validating all endpoints
- ✅ Achieve 0% empty endpoint rate
- ✅ Follow TDD approach (RED → GREEN → REFACTOR)

### Non-Goals

- ❌ Change API response formats
- ❌ Add new endpoints
- ❌ Modify generator logic (already working)
- ❌ Change seed file structure

---

## Success Metrics

| Metric | Before | After | Target |
|--------|--------|-------|--------|
| Working Endpoints | 14/29 (48%) | 29/29 (100%) | 100% ✅ |
| Empty Endpoints | 15/29 (52%) | 0/29 (0%) | 0% ✅ |
| Test Coverage | 62 tests | 89 tests | +27 E2E ✅ |
| Time Efficiency | - | 4.5h / 5.0h | 10% under ✅ |

---

## Solution Summary

### 5 Generator Calls Added to main.go

```go
// 1. Load developers into storage
store.LoadDevelopers(seedData.Developers)

// 2. Generate commits (existing)
commitGen.GenerateCommits(ctx, cfg.Days)

// 3. Generate model usage events
modelGen.GenerateModelUsage(ctx, cfg.Days)

// 4. Generate client version events
versionGen.GenerateClientVersions(ctx, cfg.Days)

// 5. Generate file extension events
extensionGen.GenerateFileExtensions(ctx, cfg.Days)

// 6. Generate feature events (MCP, Commands, Plans, AskMode)
featureGen.GenerateFeatures(ctx, cfg.Days)
```

### Test Coverage Added

- 27 new E2E test cases
- 6 test files created
- Integration test validating all 16 endpoints
- 100% pass rate

---

## Implementation Summary

**6 Tasks Completed**:
1. TASK-FIX-01: Load developers (0.5h)
2. TASK-FIX-02: Model analytics (0.75h)
3. TASK-FIX-03: Client versions (0.25h)
4. TASK-FIX-04: File extensions (0.5h)
5. TASK-FIX-05: Features - MCP/Commands/Plans/AskMode (1.5h)
6. TASK-FIX-06: Integration test (1.0h)

**Total Time**: 4.5h (10% under 5.0h estimate)

---

## Related Work Items

- `.work-items/cursor-sim-v2/` - Phase 1 foundation
- `.work-items/cursor-sim-phase3/` - Part C quality analysis
- `.work-items/cursor-sim-phase4b-cli-enhancement/` - Next: Interactive CLI
- `services/cursor-sim/SPEC.md` - Technical specification

---

**Status**: ✅ COMPLETE (January 3, 2026)
**Result**: 100% success - All 15 endpoints fixed, 0 empty datasets
