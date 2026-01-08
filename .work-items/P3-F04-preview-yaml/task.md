# Task Breakdown: Preview Mode & YAML Seed Support

**Feature ID**: P3-F04-preview-yaml
**Created**: January 8, 2026
**Status**: Ready to Start
**Inspired By**: NVIDIA NeMo DataDesigner

---

## Progress Tracker

| Phase | Tasks | Status | Estimated | Actual |
|-------|-------|--------|-----------|--------|
| **Setup** | 1 | ⏳ TODO | 0.5h | - |
| **Feature 1: YAML Support** | 3 | ⏳ TODO | 2.5h | - |
| **Feature 2: Preview Mode Core** | 4 | ⏳ TODO | 4.0h | - |
| **Feature 3: Validation Framework** | 2 | ⏳ TODO | 2.0h | - |
| **Feature 4: Integration & Polish** | 1 | ⏳ TODO | 1.5h | - |
| **TOTAL** | **11** | **0/11** | **10.5h** | **-** |

---

## Feature Breakdown

### SETUP

#### TASK-PREV-00: Initialize Work Item

**Goal**: Set up work item structure and planning documents

**Acceptance Criteria**:
- ✅ Created `.work-items/P3-F04-preview-yaml/`
- ✅ Written `user-story.md` with EARS format
- ✅ Written `design.md` with technical approach
- ✅ Written `task.md` with atomic tasks

**Estimated**: 0.5h
**Status**: ✅ COMPLETE (just completed)

---

### FEATURE 1: YAML Seed File Support

#### TASK-PREV-01: Add YAML Dependency and Basic Parser (RED)

**Goal**: Add yaml.v3 library and enhance LoadSeed function

**TDD Approach**:
```go
// Test FIRST (RED)
func TestLoadSeed_YAML(t *testing.T) {
    seed, err := LoadSeed("testdata/valid_seed.yaml")
    require.NoError(t, err)
    assert.Equal(t, 2, len(seed.Developers))
    assert.Equal(t, "alice", seed.Developers[0].UserID)
}

func TestLoadSeed_YAMLWithComments(t *testing.T) {
    seed, err := LoadSeed("testdata/commented.yaml")
    require.NoError(t, err)
    assert.Equal(t, 3, len(seed.Developers))
}

func TestLoadSeed_InvalidExtension(t *testing.T) {
    _, err := LoadSeed("testdata/seed.csv")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "unsupported seed file format")
}

func TestLoadSeed_MalformedYAML(t *testing.T) {
    _, err := LoadSeed("testdata/malformed.yaml")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "failed to parse YAML")
}
```

**Implementation Steps**:
1. Add `gopkg.in/yaml.v3` to go.mod
2. Write tests for YAML parsing
3. Modify `LoadSeed` to detect extension (.json, .yaml, .yml)
4. Implement YAML parsing branch
5. Run tests (GREEN)
6. Refactor for clarity

**Files**:
- MODIFY: `services/cursor-sim/go.mod`
- MODIFY: `services/cursor-sim/internal/seed/loader.go`
- MODIFY: `services/cursor-sim/internal/seed/loader_test.go`
- NEW: `services/cursor-sim/testdata/valid_seed.yaml`
- NEW: `services/cursor-sim/testdata/commented.yaml`
- NEW: `services/cursor-sim/testdata/malformed.yaml`

**Acceptance Criteria**:
- ✅ yaml.v3 dependency added
- ✅ LoadSeed detects .yaml and .yml extensions
- ✅ YAML files parse correctly
- ✅ Comments in YAML are ignored (YAML spec behavior)
- ✅ JSON parsing unchanged (backward compatibility)
- ✅ Clear error messages for malformed YAML
- ✅ All tests pass

**Estimated**: 1.5h
**Status**: ⏳ TODO

---

#### TASK-PREV-02: Add YAML Struct Tags (GREEN)

**Goal**: Update SeedData structs to support both JSON and YAML

**TDD Approach**:
```go
func TestSeedData_UnmarshalJSON(t *testing.T) {
    jsonData := `{"developers":[{"user_id":"alice"}]}`
    var seed SeedData
    err := json.Unmarshal([]byte(jsonData), &seed)
    require.NoError(t, err)
    assert.Equal(t, "alice", seed.Developers[0].UserID)
}

func TestSeedData_UnmarshalYAML(t *testing.T) {
    yamlData := `
developers:
  - user_id: alice
`
    var seed SeedData
    err := yaml.Unmarshal([]byte(yamlData), &seed)
    require.NoError(t, err)
    assert.Equal(t, "alice", seed.Developers[0].UserID)
}

func TestSeedData_BothFormatsEquivalent(t *testing.T) {
    jsonSeed, _ := LoadSeed("testdata/valid_seed.json")
    yamlSeed, _ := LoadSeed("testdata/valid_seed.yaml")

    // Both should produce identical structures
    assert.Equal(t, len(jsonSeed.Developers), len(yamlSeed.Developers))
    assert.Equal(t, jsonSeed.Developers[0].UserID, yamlSeed.Developers[0].UserID)
}
```

**Implementation Steps**:
1. Add `yaml:` struct tags to SeedData, Developer, WorkingHours
2. Create equivalent YAML test files
3. Write equivalence tests
4. Run tests (GREEN)

**Files**:
- MODIFY: `services/cursor-sim/internal/seed/seed.go`
- MODIFY: `services/cursor-sim/internal/seed/loader_test.go`

**Acceptance Criteria**:
- ✅ All structs have both `json:` and `yaml:` tags
- ✅ JSON and YAML produce identical SeedData structures
- ✅ Field names match (snake_case in both formats)
- ✅ All tests pass

**Estimated**: 0.5h
**Status**: ⏳ TODO

---

#### TASK-PREV-03: Add E2E Test for YAML in Runtime Mode (REFACTOR)

**Goal**: Verify YAML seeds work end-to-end in runtime mode

**TDD Approach**:
```go
func TestE2E_YAMLSeed(t *testing.T) {
    // Start server with YAML seed
    cmd := exec.Command(
        "./bin/cursor-sim",
        "-mode", "runtime",
        "-seed", "testdata/valid_seed.yaml",
        "-port", "19020",
        "-days", "7",
    )

    // Start in background
    // ...

    // Query /teams/members
    resp, err := http.Get("http://localhost:19020/teams/members")
    require.NoError(t, err)

    // Verify response
    var result struct {
        Data []map[string]interface{} `json:"data"`
    }
    json.NewDecoder(resp.Body).Decode(&result)

    assert.Greater(t, len(result.Data), 0)
    // Verify developers from YAML loaded correctly
}
```

**Implementation Steps**:
1. Create E2E test file
2. Copy valid_seed.json → valid_seed.yaml
3. Start server with YAML seed
4. Query endpoint
5. Verify data generated correctly
6. Run test (GREEN)

**Files**:
- NEW: `services/cursor-sim/test/e2e/yaml_seed_test.go`

**Acceptance Criteria**:
- ✅ Server starts with YAML seed
- ✅ Data generation works identically to JSON
- ✅ API endpoints return correct data
- ✅ E2E test passes

**Estimated**: 0.5h
**Status**: ⏳ TODO

---

### FEATURE 2: Preview Mode Core

#### TASK-PREV-04: Create Preview Package and Config (RED)

**Goal**: Implement preview package structure and configuration

**TDD Approach**:
```go
func TestPreview_New(t *testing.T) {
    seedData := &seed.SeedData{
        Developers: []seed.Developer{
            {UserID: "alice", Email: "alice@example.com"},
        },
    }

    var buf bytes.Buffer
    cfg := preview.Config{Days: 7, MaxCommits: 50}
    p := preview.New(seedData, cfg, &buf)

    assert.NotNil(t, p)
}
```

**Implementation Steps**:
1. Create `internal/preview/` package
2. Define Config struct
3. Define Preview struct
4. Implement New constructor
5. Write basic tests
6. Run tests (GREEN)

**Files**:
- NEW: `services/cursor-sim/internal/preview/preview.go`
- NEW: `services/cursor-sim/internal/preview/preview_test.go`

**Acceptance Criteria**:
- ✅ Preview package created
- ✅ Config struct defined (Days, MaxCommits, MaxEvents)
- ✅ Preview struct defined (seedData, config, writer)
- ✅ New constructor works
- ✅ Tests pass

**Estimated**: 0.5h
**Status**: ⏳ TODO

---

#### TASK-PREV-05: Implement Preview Run Method (GREEN)

**Goal**: Core preview execution logic

**TDD Approach**:
```go
func TestPreview_Run(t *testing.T) {
    seedData := &seed.SeedData{
        Developers: []seed.Developer{
            {
                UserID: "alice",
                Email:  "alice@example.com",
                WorkingHoursBand: seed.WorkingHours{Start: 9, End: 17},
                PreferredModels: []string{"claude-sonnet-4.5"},
            },
        },
    }

    var buf bytes.Buffer
    cfg := preview.Config{Days: 7, MaxCommits: 10, MaxEvents: 5}
    p := preview.New(seedData, cfg, &buf)

    err := p.Run(context.Background())
    require.NoError(t, err)

    output := buf.String()
    assert.Contains(t, output, "PREVIEW MODE")
    assert.Contains(t, output, "alice")
    assert.Contains(t, output, "Sample Commits")
}

func TestPreview_RunWithTimeout(t *testing.T) {
    // Test that preview respects context timeout
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    p := preview.New(largeSeedData, cfg, &buf)
    err := p.Run(ctx)

    // Should complete within timeout or return context error
    if err != nil && !errors.Is(err, context.DeadlineExceeded) {
        t.Errorf("unexpected error: %v", err)
    }
}
```

**Implementation Steps**:
1. Implement Run method skeleton
2. Add developer summary display
3. Add sample data generation (limited)
4. Add sample data display
5. Add statistics display
6. Write tests
7. Run tests (GREEN)

**Files**:
- MODIFY: `services/cursor-sim/internal/preview/preview.go`
- MODIFY: `services/cursor-sim/internal/preview/preview_test.go`

**Acceptance Criteria**:
- ✅ Run method executes successfully
- ✅ Output includes developer summary
- ✅ Output includes sample commits
- ✅ Output includes statistics
- ✅ Respects context timeout
- ✅ All tests pass

**Estimated**: 2.0h
**Status**: ⏳ TODO

---

#### TASK-PREV-06: Implement Preview Output Formatters (REFACTOR)

**Goal**: Create clean, readable preview output

**TDD Approach**:
```go
func TestPreview_DisplayDeveloperSummary(t *testing.T) {
    var buf bytes.Buffer
    p := preview.New(seedData, cfg, &buf)

    p.displayDeveloperSummary()

    output := buf.String()
    assert.Contains(t, output, "Developers: 2")
    assert.Contains(t, output, "alice")
    assert.Contains(t, output, "Working Hours: 09:00 - 17:00")
}

func TestPreview_DisplayStatistics(t *testing.T) {
    // Mock store with 10 commits, 2 developers
    store := &mockStore{commits: generateMockCommits(10, 2)}

    var buf bytes.Buffer
    p := preview.New(seedData, cfg, &buf)
    p.displayStatistics(store)

    output := buf.String()
    assert.Contains(t, output, "Total commits: 10")
    assert.Contains(t, output, "Developers: 2")
    assert.Contains(t, output, "Avg commits/dev: 5.0")
}

func TestPreview_TruncateLongMessages(t *testing.T) {
    longMsg := "This is a very long commit message that should be truncated to fit nicely"
    truncated := truncate(longMsg, 50)

    assert.LessOrEqual(t, len(truncated), 53) // 50 + "..."
    assert.Contains(t, truncated, "...")
}
```

**Implementation Steps**:
1. Implement displayDeveloperSummary
2. Implement displaySampleData
3. Implement displayStatistics
4. Implement truncate helper
5. Write tests for each formatter
6. Run tests (GREEN)
7. Refactor for clean output

**Files**:
- MODIFY: `services/cursor-sim/internal/preview/preview.go`
- MODIFY: `services/cursor-sim/internal/preview/preview_test.go`

**Acceptance Criteria**:
- ✅ Developer summary is readable
- ✅ Sample commits formatted nicely
- ✅ Statistics clear and informative
- ✅ Long messages truncated
- ✅ Output fits 80-column terminal
- ✅ All tests pass

**Estimated**: 1.0h
**Status**: ⏳ TODO

---

#### TASK-PREV-07: Wire Preview Mode into Main (GREEN)

**Goal**: Add preview mode routing to main entry point

**TDD Approach**:
```go
func TestMain_PreviewMode(t *testing.T) {
    cfg := &config.Config{
        Mode:     "preview",
        SeedPath: "../../testdata/valid_seed.yaml",
    }

    seedData, err := seed.LoadSeed(cfg.SeedPath)
    require.NoError(t, err)

    err = runPreviewMode(seedData, cfg)
    assert.NoError(t, err)
}

func TestMain_InvalidMode(t *testing.T) {
    cfg := &config.Config{
        Mode:     "invalid",
        SeedPath: "../../testdata/valid_seed.json",
    }

    err := run()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid mode: 'invalid'")
}
```

**Implementation Steps**:
1. Add mode router to main.go
2. Implement runPreviewMode function
3. Update runRuntimeMode function name (if needed)
4. Add mode validation
5. Write integration tests
6. Run tests (GREEN)

**Files**:
- MODIFY: `services/cursor-sim/cmd/simulator/main.go`
- MODIFY: `services/cursor-sim/cmd/simulator/main_test.go`

**Acceptance Criteria**:
- ✅ `-mode preview` triggers preview mode
- ✅ `-mode runtime` triggers runtime mode (existing)
- ✅ Invalid mode shows error
- ✅ Preview mode exits cleanly (code 0)
- ✅ Integration tests pass

**Estimated**: 0.5h
**Status**: ⏳ TODO

---

### FEATURE 3: Validation Framework

#### TASK-PREV-08: Implement Seed Validators (RED)

**Goal**: Add validation logic for seed data

**TDD Approach**:
```go
func TestPreview_ValidateSeed_ValidData(t *testing.T) {
    seedData := &seed.SeedData{
        Developers: []seed.Developer{
            {
                UserID: "alice",
                Email:  "alice@example.com",
                WorkingHoursBand: seed.WorkingHours{Start: 9, End: 17},
                PreferredModels: []string{"claude-sonnet-4.5"},
            },
        },
    }

    p := preview.New(seedData, cfg, &buf)
    err := p.validateSeed()

    assert.NoError(t, err)
    assert.Empty(t, p.warnings)
}

func TestPreview_ValidateSeed_InvalidModel(t *testing.T) {
    seedData := &seed.SeedData{
        Developers: []seed.Developer{
            {
                UserID: "alice",
                PreferredModels: []string{"gpt-5000"}, // Invalid
            },
        },
    }

    p := preview.New(seedData, cfg, &buf)
    err := p.validateSeed()

    assert.NoError(t, err) // Warnings, not errors
    assert.NotEmpty(t, p.warnings)
    assert.Contains(t, p.warnings[0], "Unknown model 'gpt-5000'")
}

func TestPreview_ValidateSeed_InvalidWorkingHours(t *testing.T) {
    seedData := &seed.SeedData{
        Developers: []seed.Developer{
            {
                UserID: "alice",
                WorkingHoursBand: seed.WorkingHours{Start: 25, End: 30}, // Invalid
            },
        },
    }

    p := preview.New(seedData, cfg, &buf)
    err := p.validateSeed()

    assert.NoError(t, err)
    assert.Contains(t, p.warnings[0], "Invalid start hour 25")
}

func TestPreview_ValidateSeed_NoDevelopers(t *testing.T) {
    seedData := &seed.SeedData{
        Developers: []seed.Developer{}, // Empty
    }

    p := preview.New(seedData, cfg, &buf)
    err := p.validateSeed()

    assert.Error(t, err) // Fatal error, not warning
    assert.Contains(t, err.Error(), "no developers defined")
}
```

**Implementation Steps**:
1. Define valid model names
2. Implement working hours validation (0-23)
3. Implement model name validation
4. Implement velocity validation (low, medium, high)
5. Implement email format validation (optional)
6. Write tests for each validator
7. Run tests (GREEN)

**Files**:
- MODIFY: `services/cursor-sim/internal/preview/preview.go`
- MODIFY: `services/cursor-sim/internal/preview/preview_test.go`

**Acceptance Criteria**:
- ✅ Valid seeds pass without warnings
- ✅ Invalid models trigger warnings
- ✅ Invalid working hours trigger warnings
- ✅ Empty developer list triggers error
- ✅ Multiple issues accumulate in warnings list
- ✅ All tests pass

**Estimated**: 1.5h
**Status**: ⏳ TODO

---

#### TASK-PREV-09: Display Validation Warnings (REFACTOR)

**Goal**: Format and display validation warnings

**TDD Approach**:
```go
func TestPreview_DisplayWarnings_NoWarnings(t *testing.T) {
    var buf bytes.Buffer
    p := preview.New(validSeedData, cfg, &buf)
    p.warnings = []string{}

    p.displayWarnings()

    output := buf.String()
    assert.Contains(t, output, "✅ No validation warnings")
}

func TestPreview_DisplayWarnings_MultipleWarnings(t *testing.T) {
    var buf bytes.Buffer
    p := preview.New(invalidSeedData, cfg, &buf)
    p.warnings = []string{
        "Developer alice: Unknown model 'gpt-5000'",
        "Developer bob: Invalid start hour 25",
    }

    p.displayWarnings()

    output := buf.String()
    assert.Contains(t, output, "⚠️  Validation Warnings")
    assert.Contains(t, output, "gpt-5000")
    assert.Contains(t, output, "Invalid start hour")
}
```

**Implementation Steps**:
1. Implement displayWarnings method
2. Add formatting for warnings list
3. Add emoji indicators (✅, ⚠️)
4. Write tests
5. Run tests (GREEN)

**Files**:
- MODIFY: `services/cursor-sim/internal/preview/preview.go`
- MODIFY: `services/cursor-sim/internal/preview/preview_test.go`

**Acceptance Criteria**:
- ✅ No warnings shows green checkmark
- ✅ Warnings displayed with ⚠️ indicator
- ✅ Each warning on separate line
- ✅ Clear and readable format
- ✅ Tests pass

**Estimated**: 0.5h
**Status**: ⏳ TODO

---

### FEATURE 4: Integration & Polish

#### TASK-PREV-10: Add E2E Test for Preview Mode (REFACTOR)

**Goal**: End-to-end validation of preview mode

**TDD Approach**:
```go
func TestE2E_PreviewMode(t *testing.T) {
    cmd := exec.Command(
        "./bin/cursor-sim",
        "-mode", "preview",
        "-seed", "testdata/valid_seed.yaml",
    )

    output, err := cmd.CombinedOutput()
    require.NoError(t, err)

    outputStr := string(output)
    assert.Contains(t, outputStr, "PREVIEW MODE")
    assert.Contains(t, outputStr, "Sample Commits")
    assert.Contains(t, outputStr, "Statistics")
    assert.Contains(t, outputStr, "Preview complete")
}

func TestE2E_PreviewWithWarnings(t *testing.T) {
    cmd := exec.Command(
        "./bin/cursor-sim",
        "-mode", "preview",
        "-seed", "testdata/invalid_models.yaml",
    )

    output, err := cmd.CombinedOutput()
    require.NoError(t, err) // Exit 0 even with warnings

    outputStr := string(output)
    assert.Contains(t, outputStr, "⚠️  Validation Warnings")
}

func TestE2E_PreviewTimeout(t *testing.T) {
    // Test with very large seed file
    cmd := exec.Command(
        "./bin/cursor-sim",
        "-mode", "preview",
        "-seed", "testdata/large_seed.yaml",
    )

    // Should complete within 10 seconds
    cmd.Start()
    done := make(chan error)
    go func() { done <- cmd.Wait() }()

    select {
    case <-time.After(10 * time.Second):
        cmd.Process.Kill()
        t.Fatal("preview mode timed out")
    case err := <-done:
        assert.NoError(t, err)
    }
}
```

**Implementation Steps**:
1. Build binary
2. Write E2E test for basic preview
3. Write E2E test with warnings
4. Write E2E test for performance (< 5s)
5. Run tests (GREEN)

**Files**:
- NEW: `services/cursor-sim/test/e2e/preview_test.go`
- NEW: `services/cursor-sim/testdata/invalid_models.yaml`

**Acceptance Criteria**:
- ✅ Preview mode runs end-to-end
- ✅ Exit code 0 on success
- ✅ Warnings display correctly
- ✅ Completes within 5 seconds
- ✅ E2E tests pass

**Estimated**: 1.0h
**Status**: ⏳ TODO

---

#### TASK-PREV-11: Update Documentation (REFACTOR)

**Goal**: Document preview mode and YAML support

**Implementation Steps**:
1. Update SPEC.md:
   - Add preview mode to Quick Start
   - Add YAML seed file example
   - Update CLI Configuration section
   - Add validation warnings documentation
2. Update README.md:
   - Add preview mode section
   - Add YAML example
   - Update usage instructions
3. Create example YAML files:
   - `testdata/preview_example.yaml` (with comments)
   - `testdata/large_team.yaml` (10+ developers)
4. Update CLI help text

**Files**:
- MODIFY: `services/cursor-sim/SPEC.md`
- MODIFY: `services/cursor-sim/README.md`
- NEW: `services/cursor-sim/testdata/preview_example.yaml`
- NEW: `services/cursor-sim/testdata/large_team.yaml`
- MODIFY: `services/cursor-sim/internal/config/config.go` (help text)

**Acceptance Criteria**:
- ✅ SPEC.md documents preview mode
- ✅ SPEC.md shows YAML example
- ✅ README.md has preview section
- ✅ Example YAML files provided
- ✅ CLI help text updated
- ✅ "Last Updated" date updated

**Estimated**: 0.5h
**Status**: ⏳ TODO

---

## Testing Strategy Summary

### Unit Tests (Go)

| Package | Test Count | Coverage Target |
|---------|------------|-----------------|
| `seed` | 10+ | 95% |
| `preview` | 15+ | 95% |
| `config` | 3+ (existing + mode validation) | 95% |

### Integration Tests

| Test | Scope |
|------|-------|
| YAML seed loading | Seed → Parse → Storage |
| Preview mode execution | Full preview flow |
| Mode routing | CLI → Router → Preview/Runtime |

### E2E Tests

| Test | Endpoint |
|------|----------|
| YAML in runtime mode | Full server with YAML seed |
| Preview mode | CLI execution and output |
| Validation warnings | Preview with invalid seed |

---

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| YAML parsing performance | Low | Low | Benchmark, limit file size to 10MB |
| Preview sample not representative | Medium | Medium | Use stratified sampling per developer |
| Binary size increase | Medium | Low | Accept 200KB increase, standard library |
| Users confused by warnings | Low | Medium | Clear messaging, examples in docs |

---

## Dependency Graph

```
TASK-PREV-00 (Setup)
    │
    ├──▶ TASK-PREV-01 (YAML Parser)
    │         │
    │         └──▶ TASK-PREV-02 (Struct Tags)
    │                   │
    │                   └──▶ TASK-PREV-03 (E2E YAML Test)
    │
    ├──▶ TASK-PREV-04 (Preview Package)
    │         │
    │         └──▶ TASK-PREV-05 (Run Method)
    │                   │
    │                   ├──▶ TASK-PREV-06 (Formatters)
    │                   │         │
    │                   │         └──▶ TASK-PREV-07 (Wire to Main)
    │                   │
    │                   └──▶ TASK-PREV-08 (Validators)
    │                             │
    │                             └──▶ TASK-PREV-09 (Warning Display)
    │
    └──▶ TASK-PREV-10 (E2E Preview)
              │
              └──▶ TASK-PREV-11 (Documentation)
```

---

## Definition of Done (Per Task)

- ✅ Tests written BEFORE implementation (TDD)
- ✅ All tests pass (unit + integration)
- ✅ Code coverage meets target (>90%)
- ✅ No linting errors (`go vet`, `gofmt`)
- ✅ Documentation updated (SPEC.md, comments)
- ✅ Git commit with descriptive message
- ✅ Dependency reflections checked
- ✅ SPEC.md synced if needed

---

## Rollout Plan

### Phase 1: YAML Support (Tasks 01-03)
- Add YAML parsing
- Update struct tags
- E2E validation
- **Estimated**: 2.5h

### Phase 2: Preview Mode Core (Tasks 04-07)
- Preview package
- Run method
- Output formatting
- Main integration
- **Estimated**: 4.0h

### Phase 3: Validation (Tasks 08-09)
- Validators
- Warning display
- **Estimated**: 2.0h

### Phase 4: Polish (Tasks 10-11)
- E2E tests
- Documentation
- **Estimated**: 1.5h

### Total Estimated Effort: 10.5 hours

---

## Success Criteria (Phase Completion)

- ✅ All 11 tasks completed
- ✅ All tests passing (100% unit, integration, E2E)
- ✅ YAML seeds work identically to JSON
- ✅ Preview mode completes in < 5 seconds
- ✅ Validation catches common seed issues
- ✅ Documentation complete with examples
- ✅ DEVELOPMENT.md updated with P3-F04 completion

---

**Next Action**: Present plan to user for approval → Begin TASK-PREV-01 (YAML Parser)
