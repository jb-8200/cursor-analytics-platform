# Empty Dataset Fix Plan

**Created**: January 3, 2026
**Priority**: HIGH (52% of endpoints broken)

---

## Root Cause Summary

All empty dataset issues stem from **missing generator/loader calls in `main.go`**:

```go
// cmd/simulator/main.go:79-88
// ONLY THIS EXISTS:
gen := generator.NewCommitGenerator(seedData, store, cfg.Velocity)
if err := gen.GenerateCommits(ctx, cfg.Days); err != nil {
    return fmt.Errorf("failed to generate commits: %w", err)
}
```

**Missing calls:**
1. ❌ `store.LoadDevelopers(seedData.Developers)` - Never called
2. ❌ `ModelGenerator.GenerateModelUsage()` - Never called
3. ❌ `VersionGenerator.GenerateVersions()` - Never called
4. ❌ `ExtensionGenerator.GenerateExtensions()` - Never called
5. ❌ `FeatureGenerator.GenerateFeatures()` - Never called

---

## Affected Endpoints (15 total)

### Critical: Team Members (1 endpoint)
| Endpoint | Current | Fix |
|----------|---------|-----|
| `/teams/members` | `{"teamMembers":[]}` | Call `store.LoadDevelopers(seedData.Developers)` |

### Team Analytics (7 endpoints)
| Endpoint | Current | Generator Needed |
|----------|---------|------------------|
| `/analytics/team/models` | `{"data":[]}` | `ModelGenerator.GenerateModelUsage()` |
| `/analytics/team/client-versions` | `{"data":[]}` | `VersionGenerator.GenerateVersions()` |
| `/analytics/team/top-file-extensions` | `{"data":[]}` | `ExtensionGenerator.GenerateExtensions()` |
| `/analytics/team/mcp` | `{"data":[]}` | `FeatureGenerator.GenerateFeatures()` |
| `/analytics/team/commands` | `{"data":[]}` | `FeatureGenerator.GenerateFeatures()` |
| `/analytics/team/plans` | `{"data":[]}` | `FeatureGenerator.GenerateFeatures()` |
| `/analytics/team/ask-mode` | `{"data":[]}` | `FeatureGenerator.GenerateFeatures()` |

### By-User Analytics (6 endpoints)
| Endpoint | Current | Fix |
|----------|---------|-----|
| `/analytics/by-user/models` | `{"data":{}}` | Same as team models |
| `/analytics/by-user/client-versions` | `{"data":{}}` | Same as team versions |
| `/analytics/by-user/top-file-extensions` | `{"data":{}}` | Same as team extensions |
| `/analytics/by-user/mcp` | `{"data":{}}` | Same as team MCP |
| `/analytics/by-user/commands` | `{"data":{}}` | Same as team commands |
| `/analytics/by-user/plans` | `{"data":{}}` | Same as team plans |
| `/analytics/by-user/ask-mode` | `{"data":{}}` | Same as team ask-mode |

---

## Fix Strategy

### Approach: Add All Missing Generator Calls to `main.go`

**Modified flow in `run()` function:**

```go
// 1. Load seed data (EXISTING)
seedData, err := seed.LoadSeed(cfg.SeedPath)

// 2. Initialize storage (EXISTING)
store := storage.NewMemoryStore()

// 3. ✅ NEW: Load developers into storage
if err := store.LoadDevelopers(seedData.Developers); err != nil {
    return fmt.Errorf("failed to load developers: %w", err)
}
log.Printf("Loaded %d developers into storage\n", len(seedData.Developers))

// 4. Generate commits (EXISTING)
commitGen := generator.NewCommitGenerator(seedData, store, cfg.Velocity)
if err := commitGen.GenerateCommits(ctx, cfg.Days); err != nil {
    return fmt.Errorf("failed to generate commits: %w", err)
}

// 5. ✅ NEW: Generate model usage events
modelGen := generator.NewModelGenerator(seedData, store, cfg.Velocity)
if err := modelGen.GenerateModelUsage(ctx, cfg.Days); err != nil {
    return fmt.Errorf("failed to generate model usage: %w", err)
}

// 6. ✅ NEW: Generate client version events
versionGen := generator.NewVersionGenerator(seedData, store, cfg.Velocity)
if err := versionGen.GenerateVersions(ctx, cfg.Days); err != nil {
    return fmt.Errorf("failed to generate versions: %w", err)
}

// 7. ✅ NEW: Generate file extension events
extensionGen := generator.NewExtensionGenerator(seedData, store, cfg.Velocity)
if err := extensionGen.GenerateExtensions(ctx, cfg.Days); err != nil {
    return fmt.Errorf("failed to generate file extensions: %w", err)
}

// 8. ✅ NEW: Generate feature events (MCP, Commands, Plans, AskMode)
featureGen := generator.NewFeatureGenerator(seedData, store, cfg.Velocity)
if err := featureGen.GenerateFeatures(ctx, cfg.Days); err != nil {
    return fmt.Errorf("failed to generate features: %w", err)
}

// 9. Create HTTP router (EXISTING)
router := server.NewRouter(store, seedData, DefaultAPIKey)
```

---

## Task Breakdown (TDD Approach)

### TASK-FIX-01: Load Developers into Storage (RED → GREEN → REFACTOR)

**Goal**: Fix `/teams/members` endpoint

**TDD Steps:**
1. **RED**: Write E2E test that fails
   ```go
   func TestE2E_TeamMembers(t *testing.T) {
       // Start server
       // GET /teams/members
       // Assert: len(teamMembers) == 2
       // Assert: contains alice@example.com
   }
   ```

2. **GREEN**: Add `store.LoadDevelopers()` call in main.go after line 77

3. **REFACTOR**: None needed

**Acceptance Criteria**:
- ✅ E2E test passes
- ✅ `/teams/members` returns 2 developers
- ✅ Developer names and emails match seed file

**Estimated**: 0.5h
**Actual**: 0.5h
**Status**: ✅ COMPLETE (Commit: 8357fc1)
**Files**: `cmd/simulator/main.go`, `test/e2e/team_members_test.go` (NEW)

---

### TASK-FIX-02: Generate Model Usage Events (RED → GREEN → REFACTOR)

**Goal**: Fix `/analytics/team/models` and `/analytics/by-user/models`

**TDD Steps:**
1. **RED**: Write E2E test
   ```go
   func TestE2E_TeamModels(t *testing.T) {
       // Start server
       // GET /analytics/team/models
       // Assert: len(data) > 0
       // Assert: contains "gpt-4-turbo" or "claude-3-sonnet"
   }
   ```

2. **GREEN**: Add `ModelGenerator.GenerateModelUsage()` call in main.go

3. **REFACTOR**: Extract generator initialization into helper function (optional)

**Acceptance Criteria**:
- ✅ Team models endpoint returns data
- ✅ By-user models endpoint returns data
- ✅ Models match seed file preferences ("gpt-4-turbo", "claude-3-sonnet")

**Estimated**: 1.0h
**Actual**: 0.75h
**Status**: ✅ COMPLETE (Commit: a017287)
**Files**: `cmd/simulator/main.go`, `test/e2e/model_analytics_test.go` (NEW)

---

### TASK-FIX-03: Generate Client Version Events (RED → GREEN → REFACTOR)

**Goal**: Fix `/analytics/team/client-versions` and by-user variant

**TDD Steps:**
1. **RED**: Write E2E test for client versions

2. **GREEN**: Add `VersionGenerator.GenerateVersions()` call

3. **REFACTOR**: None

**Acceptance Criteria**:
- ✅ Endpoints return non-empty data
- ✅ Versions are realistic (e.g., "0.41.0", "0.42.1")

**Estimated**: 0.5h
**Actual**: 0.25h
**Status**: ✅ COMPLETE (Commit: 7bdadca)
**Files**: `cmd/simulator/main.go`, `test/e2e/version_analytics_test.go` (NEW)

---

### TASK-FIX-04: Generate File Extension Events (RED → GREEN → REFACTOR)

**Goal**: Fix `/analytics/team/top-file-extensions` and by-user variant

**TDD Steps:**
1. **RED**: Write E2E test for file extensions

2. **GREEN**: Add `ExtensionGenerator.GenerateFileExtensions()` call

3. **REFACTOR**: None

**Acceptance Criteria**:
- ✅ Endpoints return data
- ✅ Extensions match repo languages (go, ts, tsx, py, sql, etc.)

**Estimated**: 0.5h
**Actual**: 0.5h
**Status**: ✅ COMPLETE (Commit: 0408b60)
**Files**: `cmd/simulator/main.go`, `test/e2e/extension_analytics_test.go` (NEW)

---

### TASK-FIX-05: Generate Feature Events (MCP, Commands, Plans, AskMode) (RED → GREEN → REFACTOR)

**Goal**: Fix 4 team endpoints + 4 by-user endpoints (8 total)

**TDD Steps:**
1. **RED**: Write E2E tests for all 4 feature types
   ```go
   func TestE2E_TeamMCP(t *testing.T) { ... }
   func TestE2E_TeamCommands(t *testing.T) { ... }
   func TestE2E_TeamPlans(t *testing.T) { ... }
   func TestE2E_TeamAskMode(t *testing.T) { ... }
   ```

2. **GREEN**: Add `FeatureGenerator.GenerateFeatures()` call (single call handles all 4)

3. **REFACTOR**: None

**Acceptance Criteria**:
- ✅ All 8 endpoints return data (4 team + 4 by-user)
- ✅ MCP tools are realistic ("read_file", "write_file", "search_web", "execute_command")
- ✅ Commands are valid ("explain", "refactor", "fix", "optimize", "test")

**Estimated**: 1.5h
**Actual**: 1.5h
**Status**: ✅ COMPLETE (Commit: a122e88)
**Files**: `cmd/simulator/main.go`, `test/e2e/feature_analytics_test.go` (NEW, 9 tests)

---

### TASK-FIX-06: Integration Test - All Endpoints (REFACTOR)

**Goal**: Comprehensive test ensuring NO empty datasets

**TDD Steps:**
1. Create `test/e2e/all_endpoints_test.go`
2. Test all 29 endpoints systematically
3. Assert: No endpoint returns empty `data` array or `{}`

**Acceptance Criteria**:
- ✅ All 29 endpoints return non-empty data
- ✅ Test runs in < 5 seconds
- ✅ Can be run in CI

**Estimated**: 1.0h
**Files**: `test/e2e/all_endpoints_test.go` (NEW)

---

## Timeline

| Task | Description | Estimated | Dependencies |
|------|-------------|-----------|--------------|
| FIX-01 | Load developers | 0.5h | None |
| FIX-02 | Model usage | 1.0h | None |
| FIX-03 | Client versions | 0.5h | None |
| FIX-04 | File extensions | 0.5h | None |
| FIX-05 | Features (MCP/Commands/Plans/AskMode) | 1.5h | None |
| FIX-06 | Integration test all endpoints | 1.0h | FIX-01 to FIX-05 |
| **TOTAL** | **6 tasks** | **5.0h** | - |

All tasks can be done **in parallel** (no dependencies between FIX-01 to FIX-05).

---

## Testing Strategy

### Unit Tests
- ✅ Existing generator tests already pass
- ✅ Storage tests already cover `LoadDevelopers()`

### E2E Tests (NEW)
| Test File | Endpoints Covered | Lines |
|-----------|-------------------|-------|
| `team_members_test.go` | 1 | ~50 |
| `model_analytics_test.go` | 2 (team + by-user) | ~80 |
| `version_analytics_test.go` | 2 | ~80 |
| `extension_analytics_test.go` | 2 | ~80 |
| `feature_analytics_test.go` | 8 (4 features × 2 modes) | ~200 |
| `all_endpoints_test.go` | All 29 | ~300 |

**Total New Test Code**: ~790 lines

---

## Success Metrics

| Metric | Before | After | Target |
|--------|--------|-------|--------|
| Working endpoints | 5/29 (17%) | 29/29 | 100% |
| Empty endpoints | 15/29 (52%) | 0/29 | 0% |
| Test coverage (cmd/simulator) | 61.7% | >75% | 75%+ |
| E2E test runtime | ~2s | ~4s | <5s |

---

## Risks & Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Generators are slow | High startup time | Optimize if > 3s |
| Storage memory usage | OOM for large datasets | Profile and optimize |
| Generator bugs | Incorrect data | Existing unit tests catch bugs |

---

## Rollout Plan

1. **Phase 1**: FIX-01 (Team members) - Quick win
2. **Phase 2**: FIX-02 to FIX-05 (All generators) - Core fix
3. **Phase 3**: FIX-06 (Integration test) - Verification

**Total Estimated Time**: 5.0 hours

---

## Next Steps

1. ✅ Get user approval for plan
2. Start with TASK-FIX-01 (quickest win)
3. Run tests after each task
4. Update SPEC.md after all fixes complete
5. Update DEVELOPMENT.md with completion status

---

**Ready to implement?** All generators are already written and tested. We just need to wire them up in `main.go`!
