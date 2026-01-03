---
name: dependency-reflection
description: Identifies required updates to related files after primary changes. Use after completing any code changes to detect documentation drift, missing test updates, and cross-file dependencies. Triggers on refactoring, model changes, or interface modifications.
---

# Dependency Reflection Check

**Purpose**: Detect when changes in one file require updates to related files (reflections).

## When to Use

Run this check after ANY of:
- Completing a task that modifies multiple files
- Refactoring existing code
- Adding new packages or services
- Modifying data models or interfaces
- Making changes to storage, generators, or API handlers

## What Are "Dependency Reflections"?

**Dependency reflections** are changes that ripple from a primary file to related files:

```
Primary Change: Add field to Commit model
       │
       ├──► Generator must populate new field
       ├──► SPEC.md schema must document new field
       ├──► Tests must validate new field
       └──► Handlers may need to filter/transform field
```

## Reflection Type Taxonomy

```
              Primary Change
                    │
    ┌───────────────┼───────────────┐
    │               │               │
Documentation   Code Sync       Test Sync
Reflections     Reflections     Reflections
```

### 1. Documentation Reflections

Changes that require documentation updates:

| Primary File | Required Doc Updates |
|-------------|---------------------|
| `internal/models/*.go` | SPEC.md → Response Format, Schema examples |
| `internal/api/**/*.go` | SPEC.md → Endpoints table |
| `internal/services/*.go` | SPEC.md → Phase Features, Architecture |
| `cmd/**/*.go` | SPEC.md → CLI Configuration, README |
| `internal/**/*.go` | SPEC.md → Package Structure |
| `.work-items/**/task.md` | DEVELOPMENT.md → Current work status |
| Phase completion | SPEC.md → Implementation Status |

### 2. Code Synchronization Reflections

Changes that require updates to related code files:

| Primary File | Required Code Updates |
|-------------|----------------------|
| `internal/models/*.go` | Generators that create these models |
| `internal/storage/store.go` (interface) | All handlers using storage methods |
| `internal/generator/*.go` | Storage methods that persist generated data |
| `internal/seed/*.go` | Generators that consume seed data |
| `internal/api/response.go` | All handlers using response helpers |

### 3. Test Synchronization Reflections

Changes that require new or updated tests:

| Primary File | Required Test Updates |
|-------------|---------------------|
| `internal/models/*.go` | Model validation tests, marshal/unmarshal tests |
| `internal/api/**/*.go` | Handler unit tests, E2E tests |
| `internal/services/*.go` | Service unit tests, integration tests |
| `internal/generator/*.go` | Generator tests with edge cases |
| `internal/storage/*.go` | Storage tests, concurrency tests |

## Reflection Detection Matrix

Use this matrix to detect reflections:

| Files Changed | Check These Files | Reflection Type |
|--------------|-------------------|-----------------|
| **Models** (`internal/models/*.go`) | | |
| → Add/modify struct field | Generators creating this model | Code Sync |
| → Add/modify struct field | Handlers returning this model | Code Sync |
| → Add/modify struct field | SPEC.md schema section | Documentation |
| → Add/modify struct field | Model tests | Test Sync |
| **API Handlers** (`internal/api/**/*.go`) | | |
| → Add new endpoint | SPEC.md endpoints table | Documentation |
| → Add new endpoint | E2E tests | Test Sync |
| → Modify handler logic | Handler unit tests | Test Sync |
| → Change response format | SPEC.md examples | Documentation |
| **Services** (`internal/services/*.go`) | | |
| → Add new service | SPEC.md phase features | Documentation |
| → Add new service | Service unit tests | Test Sync |
| → Modify service logic | Integration tests | Test Sync |
| **Generators** (`internal/generator/*.go`) | | |
| → Add new generator | Storage interface for persistence | Code Sync |
| → Modify generation logic | Generator tests | Test Sync |
| → Add new event type | Models for new event | Code Sync |
| **Storage** (`internal/storage/*.go`) | | |
| → Modify interface | ALL handlers using that method | Code Sync |
| → Add new method | Handlers that need that method | Code Sync |
| → Modify storage logic | Storage tests, concurrency tests | Test Sync |
| **Config/CLI** (`internal/config/*.go`, `cmd/**/*.go`) | | |
| → Add/modify flag | SPEC.md CLI section | Documentation |
| → Add/modify flag | README usage examples | Documentation |
| → Add/modify flag | Config tests | Test Sync |
| **Seed** (`internal/seed/*.go`) | | |
| → Modify schema | Generators consuming seed | Code Sync |
| → Add validation | Seed loader tests | Test Sync |
| **Work Items** (`.work-items/**/task.md`) | | |
| → Complete step | DEVELOPMENT.md current work | Documentation |
| → Complete phase | SPEC.md implementation status | Documentation |

## Reflection Checklist

After modifying files, verify:

### Documentation Sync
- [ ] SPEC.md reflects new/changed endpoints
- [ ] SPEC.md reflects new/changed models (schema section)
- [ ] SPEC.md reflects phase completion status
- [ ] SPEC.md Package Structure includes new directories
- [ ] DEVELOPMENT.md reflects current work state
- [ ] task.md reflects step completion status

### Code Sync
- [ ] Generators produce instances of new/modified models
- [ ] Handlers use correct storage methods (if storage changed)
- [ ] Models have all required JSON tags
- [ ] Interfaces are implemented by all consumers
- [ ] Response helpers used consistently

### Test Sync
- [ ] New code paths have test coverage
- [ ] Modified behavior has updated assertions
- [ ] E2E tests cover new endpoints
- [ ] Edge cases are tested
- [ ] Concurrency tests for storage changes

## Regression Test Protocol

After identifying reflections, run appropriate regression tests:

### After Unit Changes (handlers, services)
```bash
go test ./internal/... -v
```

### After Handler Changes
```bash
go test ./internal/api/... -v
go test ./test/e2e/... -v
```

### After Model Changes
```bash
go test ./... -cover
# Verify coverage >= 80%
```

### After Storage Interface Changes
```bash
go test ./internal/... -v -race
# Check for race conditions
```

### After Major Refactor
```bash
go test ./... -v -race -count=5
# Run multiple times to catch flaky tests
```

## Integration with SDD Workflow

This check happens at **Step 5: REFLECT** in the enhanced SDD cycle:

```
3. CODE     → Minimal implementation (GREEN)
4. REFACTOR → Clean up while tests pass
5. REFLECT  → Check dependency reflections ← YOU ARE HERE
6. SYNC     → Update SPEC.md if triggered
7. COMMIT   → Commit code + docs together
```

## Example Scenarios

### Scenario 1: Added New Field to Commit Model

**Primary Change**: Added `prNumber` field to `internal/models/commit.go`

**Detected Reflections**:
1. **Code Sync**:
   - `internal/generator/commit_generator.go` must populate `prNumber`
   - `internal/generator/pr_generator.go` must assign PRs to commits
2. **Documentation**:
   - `SPEC.md` lines 175-200: Update Commit Schema JSON example
3. **Test Sync**:
   - `internal/models/commit_test.go`: Add test for `prNumber` marshaling
   - `internal/generator/commit_generator_test.go`: Verify `prNumber` populated

**Actions**:
1. Update generators to populate `prNumber`
2. Update SPEC.md Commit schema
3. Add/update tests
4. Run `go test ./... -cover` to verify
5. Commit all changes together

### Scenario 2: Added New Endpoint `/analytics/team/models`

**Primary Change**: Created `internal/api/cursor/team.go` handler

**Detected Reflections**:
1. **Code Sync**:
   - `internal/generator/model_usage.go` must generate model usage events
   - `internal/storage/store.go` may need new method for model data
2. **Documentation**:
   - `SPEC.md` lines 120-132: Add endpoint to Team Analytics table
3. **Test Sync**:
   - `internal/api/cursor/team_test.go`: Add handler test
   - `test/e2e/team_test.go`: Add E2E test

**Actions**:
1. Verify generator produces model usage data
2. Update SPEC.md endpoints table
3. Write handler and E2E tests
4. Run `go test ./internal/api/... ./test/e2e/...`
5. Commit all changes together

### Scenario 3: Completed Phase 3 Part C (Step C06)

**Primary Change**: Completed last step of Phase 3 Part C

**Detected Reflections**:
1. **Documentation**:
   - `SPEC.md` lines 17-22: Update Phase 3 status to "MOSTLY COMPLETE ✅"
   - `SPEC.md` lines 404-431: Update Phase 3 Features section
   - `.work-items/cursor-sim-phase3/task.md`: Mark C06 as DONE
   - `DEVELOPMENT.md`: Update current work status
2. **Test Sync**:
   - Verify all Part C tests pass

**Actions**:
1. Update SPEC.md Implementation Status table
2. Update SPEC.md Phase 3 Features section
3. Update task.md progress tracker
4. Update DEVELOPMENT.md
5. Run full regression: `go test ./...`
6. Commit all documentation updates together

## Quick Reference Card

**After making changes, ask yourself:**

### Did I modify models?
→ Check: Generators, Handlers, SPEC.md schema, Tests

### Did I add/modify endpoints?
→ Check: SPEC.md endpoints table, E2E tests

### Did I change storage interface?
→ Check: All handlers using storage, Tests

### Did I complete a phase step?
→ Check: SPEC.md status, task.md, DEVELOPMENT.md

### Did I refactor code?
→ Check: All tests still pass, Docs still accurate

**If answer is YES** → Run reflection check before committing!

## Automation Suggestion

Consider adding a pre-commit reminder:
```bash
# Before every commit, run:
echo "Reflection Check:"
echo "1. Documentation sync? (SPEC.md, task.md, DEVELOPMENT.md)"
echo "2. Code sync? (Generators, handlers, storage)"
echo "3. Test sync? (Unit tests, E2E tests, coverage)"
```

---

**Remember**: Code changes rarely happen in isolation. Check reflections before every commit!
