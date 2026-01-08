# Task Breakdown: Interactive CLI Configuration

**Feature ID**: P4-F02-cli-enhancement
**Created**: January 3, 2026
**Status**: Ready to Start

**Note**: Phase 4A (Empty Dataset Fixes) was completed separately. This task list is for P4-F02 only.

---

## Progress Tracker

| Phase | Tasks | Status | Estimated | Actual |
|-------|-------|--------|-----------|--------|
| **Setup** | 1 | ✅ DONE | 0.5h | 0.5h |
| **Feature 1: Interactive Prompts** | 3 | ✅ DONE | 4.0h | 3.0h |
| **Feature 2: Developer Replication** | 3 | ✅ DONE | 3.0h | 2.5h |
| **Feature 3: Commit Limit** | 3 | ✅ DONE | 2.5h | 2.5h |
| **Feature 4: Integration** | 2 | ⏳ IN PROGRESS | 2.0h | 1.0h |
| **Feature 5: Empty Dataset Fix** | 2 | ⏳ TODO | 2.0h | - |
| **TOTAL** | **14** | **10/14** | **14.0h** | **9.5h** |

---

## Feature Breakdown

### SETUP

#### TASK-CLI-00: Initialize Work Item

**Goal**: Set up work item structure and planning documents

**Acceptance Criteria**:
- ✅ Created `.work-items/cursor-sim-phase4-cli-enhancement/`
- ✅ Written `user-story.md` with EARS format
- ✅ Written `design.md` with technical approach
- ✅ Written `task.md` with atomic tasks

**Estimated**: 0.5h
**Status**: ✅ COMPLETE

---

### FEATURE 1: Interactive Prompt Module

#### TASK-CLI-01: Create Interactive Prompt Infrastructure (RED)

**Goal**: Implement `PromptForInt` function with validation and tests

**TDD Approach**:
```go
// Test FIRST (RED)
func TestPromptForInt_ValidInput(t *testing.T) {
    // Simulate user typing "5\n"
    input := strings.NewReader("5\n")
    // ... assert returns 5
}

func TestPromptForInt_DefaultOnEmpty(t *testing.T) {
    // Simulate user pressing Enter
    input := strings.NewReader("\n")
    // ... assert returns default value
}

func TestPromptForInt_InvalidInputRetry(t *testing.T) {
    // Simulate "abc\n5\n"
    input := strings.NewReader("abc\n5\n")
    // ... assert retries and accepts 5
}
```

**Implementation Steps**:
1. Write tests for:
   - Valid integer input
   - Empty input (use default)
   - Invalid input with retry
   - Out-of-range input with retry
   - Max retries exceeded (use default)
2. Create `internal/config/interactive.go`
3. Implement `PromptConfig` struct
4. Implement `PromptForInt` function
5. Run tests (GREEN)
6. Refactor for clarity

**Files**:
- NEW: `internal/config/interactive.go`
- NEW: `internal/config/interactive_test.go`

**Acceptance Criteria**:
- ✅ All 5 test cases pass
- ✅ Regex validation works for `^\d+$`
- ✅ Range validation enforces min/max
- ✅ Empty input uses default
- ✅ Max retries defaults to 3

**Estimated**: 2.0h
**Status**: ✅ COMPLETE
**Actual**: 1.5h

---

#### TASK-CLI-02: Implement InteractiveConfig Function (GREEN)

**Goal**: Wire up 3 prompts (developers, months, max commits)

**TDD Approach**:
```go
func TestInteractiveConfig_AllDefaults(t *testing.T) {
    // Simulate pressing Enter 3 times
    input := strings.NewReader("\n\n\n")
    // ... assert returns GenerationParams{2, 90, 1000}
}

func TestInteractiveConfig_CustomValues(t *testing.T) {
    // Simulate "5\n6\n2500\n"
    input := strings.NewReader("5\n6\n2500\n")
    // ... assert returns GenerationParams{5, 180, 2500}
}
```

**Implementation Steps**:
1. Write tests for:
   - All defaults (3x Enter)
   - Custom values (valid inputs)
   - Mixed (some defaults, some custom)
   - Month-to-day conversion (6 months → 180 days)
2. Implement `InteractiveConfig()` function
3. Add summary output (validation display)
4. Run tests (GREEN)

**Files**:
- MODIFY: `internal/config/interactive.go`
- MODIFY: `internal/config/interactive_test.go`

**Acceptance Criteria**:
- ✅ 3 prompts execute in sequence
- ✅ Months converted to days (months * 30)
- ✅ Summary displays validated values
- ✅ Tests cover all scenarios

**Estimated**: 1.5h
**Status**: ✅ COMPLETE
**Actual**: 1.0h

---

#### TASK-CLI-03: Add Interactive Flag to Config Struct (REFACTOR)

**Goal**: Update `Config` struct to support interactive mode

**TDD Approach**:
```go
func TestParseFlags_InteractiveMode(t *testing.T) {
    args := []string{"-interactive"}
    cfg, err := parseFlagsWithArgs(args)
    assert.NoError(t, err)
    assert.True(t, cfg.Interactive)
}
```

**Implementation Steps**:
1. Add `Interactive bool` to `Config` struct
2. Add `GenParams GenerationParams` to `Config`
3. Add `-interactive` flag to `ParseFlags`
4. Add `-developers`, `-max-commits` flags for non-interactive mode
5. Write tests for flag parsing
6. Run tests (GREEN)

**Files**:
- MODIFY: `internal/config/config.go`
- MODIFY: `internal/config/config_test.go`

**Acceptance Criteria**:
- ✅ `Config` struct has `Interactive` and `GenParams` fields
- ✅ `-interactive` flag parses correctly
- ✅ Non-interactive flags work (`-developers`, `-months`, `-max-commits`)
- ✅ Backward compatible (existing flags unchanged)
- ✅ Mixed mode validation (can't use both interactive and non-interactive)
- ✅ Months-to-days conversion (months * 30)
- ✅ All tests pass

**Estimated**: 0.5h
**Status**: ✅ COMPLETE
**Actual**: 0.5h

---

### FEATURE 2: Developer Replication

#### TASK-CLI-04: Create Developer Replicator Module (RED)

**Goal**: Implement `ReplicateDevelopers` function with sampling/cloning

**TDD Approach**:
```go
func TestReplicateDevelopers_Downsample(t *testing.T) {
    // Seed: 5 developers, Request: 2
    // ... assert returns 2 random developers from seed
}

func TestReplicateDevelopers_ExactMatch(t *testing.T) {
    // Seed: 3 developers, Request: 3
    // ... assert returns all 3 developers
}

func TestReplicateDevelopers_Replicate(t *testing.T) {
    // Seed: 2 developers, Request: 5
    // ... assert returns 5 developers with cloned IDs/emails
}
```

**Implementation Steps**:
1. Write tests for:
   - Downsample (N < seed count)
   - Exact match (N == seed count)
   - Replicate (N > seed count)
   - Unique IDs for clones
   - Clone naming convention
2. Create `internal/seed/replicator.go`
3. Implement `ReplicateDevelopers` function
4. Run tests (GREEN)

**Files**:
- NEW: `internal/seed/replicator.go`
- NEW: `internal/seed/replicator_test.go`

**Acceptance Criteria**:
- ✅ Downsampling uses random selection
- ✅ Replication clones developers with unique IDs
- ✅ Clone naming: `user_001_clone1`, `alice_clone1@example.com`
- ✅ All tests pass

**Estimated**: 1.5h
**Status**: ✅ COMPLETE
**Actual**: 1.5h

---

#### TASK-CLI-05: Integrate Replicator into Seed Loading (GREEN)

**Goal**: Modify seed loading to replicate developers based on config

**TDD Approach**:
```go
func TestLoadSeedWithReplication(t *testing.T) {
    seedPath := "testdata/valid_seed.json"  // 2 developers
    developers := 5
    // ... assert seed.Developers has 5 entries
}
```

**Implementation Steps**:
1. Add `ReplicatedDevelopers` field to `SeedData`
2. Modify seed loader to accept developer count
3. Call `ReplicateDevelopers` if count > 0
4. Write integration test
5. Run tests (GREEN)

**Files**:
- MODIFY: `internal/seed/loader.go`
- MODIFY: `internal/seed/loader_test.go`

**Acceptance Criteria**:
- ✅ `LoadSeedWithReplication` accepts optional developer count
- ✅ Returns replicated developers when count specified
- ✅ Original seed data preserved
- ✅ Integration test validates replication

**Estimated**: 1.0h
**Status**: ✅ COMPLETE
**Actual**: 0.5h

---

#### TASK-CLI-06: Add E2E Test for Developer Replication (REFACTOR)

**Goal**: End-to-end test with full generation pipeline

**TDD Approach**:
```go
func TestE2E_DeveloperReplication(t *testing.T) {
    // Start cursor-sim with 5 developers
    // Query /teams/members endpoint
    // ... assert 5 developers returned
}
```

**Implementation Steps**:
1. Write E2E test in `test/e2e/`
2. Start server with replicated developers
3. Query team members endpoint
4. Assert correct count and unique IDs
5. Run test (GREEN)

**Files**:
- NEW: `test/e2e/developer_replication_test.go`

**Acceptance Criteria**:
- ✅ E2E test spawns server with N developers
- ✅ API returns exactly N developers
- ✅ All developer IDs are unique
- ✅ Test passes

**Estimated**: 0.5h
**Status**: ✅ COMPLETE
**Actual**: 0.5h

---

### FEATURE 3: Commit Limit

#### TASK-CLI-07: Add Max Commit Tracking to Generator (RED)

**Goal**: Modify `GenerateCommits` to stop at max commits

**TDD Approach**:
```go
func TestGenerateCommits_MaxLimit(t *testing.T) {
    // 2 developers, 90 days, max 10 commits
    // ... assert exactly 10 commits generated
}

func TestGenerateCommits_NoLimit(t *testing.T) {
    // 2 developers, 90 days, max 0 (unlimited)
    // ... assert commits generated based on Poisson process
}
```

**Implementation Steps**:
1. Write tests for:
   - Max limit reached
   - No limit (0 = unlimited)
   - Limit distributes across developers
2. Modify `GenerateCommits` signature: add `maxCommits int`
3. Add tracking counter
4. Stop generation when limit reached
5. Run tests (GREEN)

**Files**:
- MODIFY: `internal/generator/commit_generator.go`
- MODIFY: `internal/generator/commit_generator_test.go`

**Acceptance Criteria**:
- ✅ `GenerateCommits(ctx, days, maxCommits int)` signature
- ✅ Generation stops at max commits
- ✅ 0 = unlimited (existing behavior)
- ✅ Log message when limit reached early

**Estimated**: 1.5h
**Status**: ✅ COMPLETE
**Actual**: 1.5h
**Commit**: 0618a1b

---

#### TASK-CLI-08: Update Main to Pass Max Commits (GREEN)

**Goal**: Wire interactive params to generator

**TDD Approach**:
```go
func TestMain_InteractiveWithMaxCommits(t *testing.T) {
    // Simulate interactive: "2\n3\n500\n"
    // ... assert generator called with maxCommits=500
}
```

**Implementation Steps**:
1. Modify `run()` in `main.go` to use `cfg.GenParams.MaxCommits`
2. Pass to `GenerateCommits(ctx, days, maxCommits)`
3. Write integration test
4. Run test (GREEN)

**Files**:
- MODIFY: `cmd/simulator/main.go`
- MODIFY: `cmd/simulator/main_test.go`

**Acceptance Criteria**:
- ✅ Main passes max commits to generator
- ✅ Interactive mode parameters flow correctly
- ✅ Non-interactive mode still works
- ✅ Tests pass

**Estimated**: 0.5h
**Status**: ✅ COMPLETE
**Actual**: 0.5h
**Commit**: 871f3a1

---

#### TASK-CLI-09: Add E2E Test for Commit Limit (REFACTOR)

**Goal**: End-to-end test verifying commit limit enforcement

**TDD Approach**:
```go
func TestE2E_CommitLimit(t *testing.T) {
    // Start with max 50 commits
    // Query /analytics/ai-code/commits
    // ... assert exactly 50 commits returned (across all pages)
}
```

**Implementation Steps**:
1. Write E2E test
2. Start server with max commits
3. Query all commits (paginated)
4. Assert exact count
5. Run test (GREEN)

**Files**:
- NEW: `test/e2e/commit_limit_test.go`

**Acceptance Criteria**:
- ✅ E2E test verifies commit count
- ✅ Pagination handled correctly
- ✅ Test passes

**Estimated**: 0.5h
**Status**: ✅ COMPLETE
**Actual**: 0.5h
**Commit**: 11f574b

---

### FEATURE 4: Integration & CLI UX

#### TASK-CLI-10: Wire Interactive Mode into Main Entry Point (GREEN)

**Goal**: Connect interactive prompts to startup flow

**Implementation Steps**:
1. Modify `main.go` to check `cfg.Interactive`
2. Call `InteractiveConfig()` if true
3. Override `cfg.GenParams` with results
4. Display startup summary
5. Test manually

**Files**:
- MODIFY: `cmd/simulator/main.go`

**Acceptance Criteria**:
- ✅ `-interactive` flag triggers prompts
- ✅ Non-interactive mode skips prompts
- ✅ Parameters flow to generator correctly
- ✅ Startup summary displays chosen values

**Estimated**: 1.0h
**Status**: ✅ COMPLETE
**Actual**: 1.0h
**Commit**: 61206b2

---

#### TASK-CLI-11: Manual Testing & UX Polish (REFACTOR)

**Goal**: Polish interactive UX and error messages

**Steps**:
1. Run interactive mode manually
2. Test all input scenarios:
   - All defaults
   - Custom values
   - Invalid inputs
   - Retry behavior
3. Refine error messages
4. Add colored output (optional)
5. Update SPEC.md with new flags

**Files**:
- MODIFY: `internal/config/interactive.go`
- MODIFY: `services/cursor-sim/SPEC.md`

**Acceptance Criteria**:
- ✅ UX feels smooth and intuitive
- ✅ Error messages are clear
- ✅ SPEC.md updated with `-interactive`, `-developers`, `-max-commits`
- ✅ Manual test checklist completed

**Estimated**: 1.0h
**Status**: ⏳ TODO

---

### FEATURE 5: Fix Empty Dataset Issues

#### TASK-CLI-12: Investigate Empty Dataset Root Cause (RED)

**Goal**: Identify why `/analytics/team/models` returns empty data

**Steps**:
1. Run generator with current seed
2. Query problematic endpoints:
   - `/analytics/team/models`
   - `/analytics/team/mcp`
   - `/analytics/team/commands`
3. Trace through handler code
4. Identify missing generator calls
5. Document findings

**Files**:
- NOTES: Document root cause in design.md

**Acceptance Criteria**:
- ✅ Root cause identified
- ✅ Missing generators documented
- ✅ Fix strategy determined

**Estimated**: 1.0h
**Status**: ⏳ TODO

---

#### TASK-CLI-13: Implement Missing Generators (GREEN)

**Goal**: Ensure all analytics endpoints populate data

**TDD Approach**:
```go
func TestModelGenerator_PopulatesData(t *testing.T) {
    // Generate commits + models
    // ... assert model analytics has data
}
```

**Implementation Steps**:
1. Write tests for missing generators
2. Implement or fix generator calls
3. Run E2E tests for all endpoints
4. Assert no endpoints return empty data
5. Run tests (GREEN)

**Files**:
- MODIFY: `internal/generator/*.go` (as needed)
- MODIFY: `test/e2e/*_test.go`

**Acceptance Criteria**:
- ✅ All 29 endpoints return non-empty data
- ✅ Model, MCP, Commands analytics populated
- ✅ E2E tests pass for all endpoints
- ✅ No empty `data` arrays in responses

**Estimated**: 1.0h
**Status**: ⏳ TODO

---

## Testing Strategy Summary

### Unit Tests (Go)

| Package | Test Count | Coverage Target |
|---------|------------|-----------------|
| `config` | 15+ | 95% |
| `seed` | 10+ | 90% |
| `generator` | 5+ (modified) | 85% |

### Integration Tests

| Test | Scope |
|------|-------|
| Interactive config flow | Full prompt sequence |
| Developer replication | Seed → Replicate → Storage |
| Commit limit enforcement | Generate → Query → Assert count |

### E2E Tests

| Test | Endpoint |
|------|----------|
| Developer count | `/teams/members` |
| Commit limit | `/analytics/ai-code/commits` |
| Empty dataset fix | All 29 endpoints |

---

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Stdin blocking in tests | Medium | High | Mock stdin with `strings.Reader` |
| Developer cloning breaks uniqueness | Low | Medium | Enforce unique IDs in tests |
| Max commits too restrictive | Low | Low | Document in SPEC.md |
| Backward compatibility broken | Low | High | Keep all existing flags |

---

## Dependency Graph

```
TASK-CLI-00 (Setup)
    │
    ├──▶ TASK-CLI-01 (PromptForInt)
    │         │
    │         └──▶ TASK-CLI-02 (InteractiveConfig)
    │                   │
    │                   └──▶ TASK-CLI-03 (Config Struct)
    │                             │
    │                             └──▶ TASK-CLI-10 (Wire to Main)
    │
    ├──▶ TASK-CLI-04 (Replicator)
    │         │
    │         └──▶ TASK-CLI-05 (Integrate)
    │                   │
    │                   └──▶ TASK-CLI-06 (E2E Test)
    │
    ├──▶ TASK-CLI-07 (Max Commits)
    │         │
    │         └──▶ TASK-CLI-08 (Wire to Main)
    │                   │
    │                   └──▶ TASK-CLI-09 (E2E Test)
    │
    ├──▶ TASK-CLI-12 (Investigate Empty Datasets)
    │         │
    │         └──▶ TASK-CLI-13 (Fix Generators)
    │
    └──▶ TASK-CLI-11 (UX Polish)
```

---

## Definition of Done (Per Task)

- ✅ Tests written BEFORE implementation (TDD)
- ✅ All tests pass (unit + integration)
- ✅ Code coverage meets target (>85%)
- ✅ No linting errors (`go vet`, `gofmt`)
- ✅ Documentation updated (SPEC.md, comments)
- ✅ Git commit with descriptive message
- ✅ Dependency reflections checked
- ✅ SPEC.md synced if needed

---

## Rollout Plan

### Phase 1: Core Implementation (Tasks 01-09)
- Interactive prompts
- Developer replication
- Commit limits
- **Estimated**: 10.0h

### Phase 2: Integration & Polish (Tasks 10-11)
- Wire to main entry point
- UX refinements
- **Estimated**: 2.0h

### Phase 3: Empty Dataset Fix (Tasks 12-13)
- Root cause analysis
- Generator fixes
- **Estimated**: 2.0h

### Total Estimated Effort: 14.0 hours

---

## Success Criteria (Phase Completion)

- ✅ All 14 tasks completed
- ✅ All tests passing (15/15 packages)
- ✅ 0 endpoints return empty data
- ✅ Interactive mode works flawlessly
- ✅ Backward compatibility maintained
- ✅ SPEC.md updated with new features
- ✅ DEVELOPMENT.md updated with Phase 4 completion

---

**Next Action**: Present plan to user for approval → Begin TASK-CLI-01
