# Technical Design: Preview Mode & YAML Seed Support

**Feature ID**: P3-F04-preview-yaml
**Created**: January 8, 2026
**Status**: Draft
**Architecture Level**: Service Enhancement (cursor-sim)

---

## Overview

Implement preview mode for quick seed validation and add YAML seed file support, inspired by NVIDIA DataDesigner patterns. This enhances cursor-sim's usability without changing core generation logic.

---

## Goals

1. **Fast Feedback Loop**: Preview mode completes in < 5 seconds
2. **Validation Early**: Catch seed issues before full generation
3. **Readable Config**: YAML format with comments for complex seeds
4. **Backward Compatible**: JSON seeds continue working identically
5. **Zero Runtime Impact**: Preview mode doesn't affect runtime mode performance

---

## Architecture

### Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│  (cmd/simulator/main.go)                                    │
│                                                              │
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │ -mode flag  │  │ -seed flag   │  │ Flag Parser     │   │
│  │ runtime     │  │ .json/.yaml  │  │ (existing)      │   │
│  │ preview ◄───┼──┼──────────────┼──► New: preview    │   │
│  └─────────────┘  └──────────────┘  └─────────────────┘   │
└──────────────────────────┬───────────────────────────────┘
                           │
                ┌──────────▼──────────┐
                │  Mode Router (NEW)  │
                │  mode == "preview"? │
                └──────┬──────┬───────┘
                       │      │
         ┌─────────────┘      └────────────┐
         │                                  │
┌────────▼────────┐              ┌─────────▼──────────┐
│ Preview Mode    │              │  Runtime Mode      │
│ (NEW)           │              │  (existing)        │
│                 │              │                    │
│ • Load seed     │              │ • Load seed        │
│ • Generate 7d   │              │ • Generate full    │
│ • Format output │              │ • Start server     │
│ • Exit          │              │ • Serve endpoints  │
└────────┬────────┘              └─────────┬──────────┘
         │                                  │
         │                                  │
    ┌────▼──────────────────────────────────▼────┐
    │      Seed Loader (ENHANCED)                │
    │                                             │
    │  ┌──────────────┐    ┌──────────────┐     │
    │  │ JSON Parser  │    │ YAML Parser  │     │
    │  │ (existing)   │    │ (NEW)        │     │
    │  └──────┬───────┘    └──────┬───────┘     │
    │         │                    │             │
    │         └────────┬───────────┘             │
    │                  │                         │
    │        ┌─────────▼──────────┐             │
    │        │  SeedData Struct   │             │
    │        │  (unchanged)       │             │
    │        └─────────┬──────────┘             │
    │                  │                         │
    │        ┌─────────▼──────────┐             │
    │        │  Validator (NEW)   │             │
    │        │  • Model names     │             │
    │        │  • Working hours   │             │
    │        │  • Velocity        │             │
    │        └────────────────────┘             │
    └─────────────────────────────────────────┘
```

---

## Component Design

### 1. Seed Loader Enhancement

**File**: `internal/seed/loader.go`

#### Current Implementation
```go
func LoadSeed(path string) (*SeedData, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var seed SeedData
    if err := json.Unmarshal(data, &seed); err != nil {
        return nil, err
    }

    return &seed, nil
}
```

#### Enhanced Implementation
```go
func LoadSeed(path string) (*SeedData, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read seed file: %w", err)
    }

    // Detect format by extension
    var seed SeedData
    ext := filepath.Ext(path)

    switch ext {
    case ".json":
        if err := json.Unmarshal(data, &seed); err != nil {
            return nil, fmt.Errorf("failed to parse JSON: %w", err)
        }
    case ".yaml", ".yml":
        if err := yaml.Unmarshal(data, &seed); err != nil {
            return nil, fmt.Errorf("failed to parse YAML: %w", err)
        }
    default:
        return nil, fmt.Errorf("unsupported seed file format: %s (use .json or .yaml)", ext)
    }

    return &seed, nil
}
```

**Dependencies**: `gopkg.in/yaml.v3`

**Testing Strategy**:
- Test JSON parsing (existing behavior)
- Test YAML parsing (new behavior)
- Test YAML with comments
- Test invalid extensions
- Test malformed JSON/YAML

---

### 2. Preview Mode Implementation

**File**: `internal/preview/preview.go` (NEW)

#### Interface
```go
package preview

import (
    "context"
    "fmt"
    "io"
    "github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
    "github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
    "github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// Config defines preview mode configuration
type Config struct {
    Days       int  // Default: 7
    MaxCommits int  // Default: 50
    MaxEvents  int  // Per event type, default: 10
}

// Preview generates sample data and displays formatted output
type Preview struct {
    seedData *seed.SeedData
    config   Config
    writer   io.Writer
}

func New(seedData *seed.SeedData, config Config, writer io.Writer) *Preview {
    return &Preview{
        seedData: seedData,
        config:   config,
        writer:   writer,
    }
}

// Run executes preview mode
func (p *Preview) Run(ctx context.Context) error {
    // 1. Validate seed data
    if err := p.validateSeed(); err != nil {
        return fmt.Errorf("seed validation failed: %w", err)
    }

    // 2. Display developer summary
    p.displayDeveloperSummary()

    // 3. Generate sample data
    store := storage.NewMemoryStore()
    if err := p.generateSampleData(ctx, store); err != nil {
        return fmt.Errorf("sample generation failed: %w", err)
    }

    // 4. Display sample data
    p.displaySampleData(store)

    // 5. Display statistics
    p.displayStatistics(store)

    // 6. Display warnings
    p.displayWarnings()

    return nil
}
```

#### Key Methods

**Validation**:
```go
func (p *Preview) validateSeed() error {
    warnings := []string{}

    // Check developers
    if len(p.seedData.Developers) == 0 {
        return fmt.Errorf("no developers defined in seed")
    }

    validModels := map[string]bool{
        "claude-sonnet-4.5": true,
        "claude-opus-4":     true,
        "gpt-4o":           true,
        "gpt-4-turbo":      true,
    }

    for _, dev := range p.seedData.Developers {
        // Validate working hours
        if dev.WorkingHoursBand.Start < 0 || dev.WorkingHoursBand.Start > 23 {
            warnings = append(warnings, fmt.Sprintf(
                "Developer %s: Invalid start hour %d (must be 0-23)",
                dev.UserID, dev.WorkingHoursBand.Start))
        }

        // Validate models
        for _, model := range dev.PreferredModels {
            if !validModels[model] {
                warnings = append(warnings, fmt.Sprintf(
                    "Developer %s: Unknown model '%s'", dev.UserID, model))
            }
        }
    }

    p.warnings = warnings
    return nil
}
```

**Sample Generation**:
```go
func (p *Preview) generateSampleData(ctx context.Context, store storage.Store) error {
    // Use shorter time period and limits
    days := p.config.Days
    maxCommits := p.config.MaxCommits

    // Generate commits
    commitGen := generator.NewCommitGenerator(p.seedData, store, "medium", 12345)
    if err := commitGen.GenerateCommits(ctx, days, maxCommits); err != nil {
        return fmt.Errorf("commit generation: %w", err)
    }

    // Generate model usage (smaller sample)
    modelGen := generator.NewModelGenerator(p.seedData, store, "medium", 12345)
    if err := modelGen.GenerateModelUsage(ctx, days); err != nil {
        return fmt.Errorf("model usage generation: %w", err)
    }

    // Skip or minimize other generators for speed
    return nil
}
```

**Output Formatting**:
```go
func (p *Preview) displayDeveloperSummary() {
    fmt.Fprintf(p.writer, "\n═══ PREVIEW MODE ═══\n\n")
    fmt.Fprintf(p.writer, "Developers: %d\n\n", len(p.seedData.Developers))

    for _, dev := range p.seedData.Developers {
        fmt.Fprintf(p.writer, "  • %s (%s)\n", dev.UserID, dev.Email)
        fmt.Fprintf(p.writer, "    Working Hours: %02d:00 - %02d:00\n",
            dev.WorkingHoursBand.Start, dev.WorkingHoursBand.End)
        fmt.Fprintf(p.writer, "    Models: %v\n", dev.PreferredModels)
    }
    fmt.Fprintf(p.writer, "\n")
}

func (p *Preview) displaySampleData(store storage.Store) {
    commits := store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour))

    fmt.Fprintf(p.writer, "Sample Commits (first 10):\n\n")
    for i, commit := range commits {
        if i >= 10 {
            break
        }
        fmt.Fprintf(p.writer, "  %s | %s | %s | %s\n",
            commit.Timestamp.Format("2006-01-02 15:04"),
            commit.UserID,
            commit.Model,
            truncate(commit.CommitMessage, 50))
    }
    fmt.Fprintf(p.writer, "\n")
}

func (p *Preview) displayStatistics(store storage.Store) {
    commits := store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour))
    developers := store.ListDevelopers()

    fmt.Fprintf(p.writer, "Statistics:\n")
    fmt.Fprintf(p.writer, "  • Total commits: %d\n", len(commits))
    fmt.Fprintf(p.writer, "  • Developers: %d\n", len(developers))
    fmt.Fprintf(p.writer, "  • Avg commits/dev: %.1f\n",
        float64(len(commits))/float64(len(developers)))
    fmt.Fprintf(p.writer, "\n")
}

func (p *Preview) displayWarnings() {
    if len(p.warnings) == 0 {
        fmt.Fprintf(p.writer, "✅ No validation warnings\n\n")
        return
    }

    fmt.Fprintf(p.writer, "⚠️  Validation Warnings:\n")
    for _, warning := range p.warnings {
        fmt.Fprintf(p.writer, "  - %s\n", warning)
    }
    fmt.Fprintf(p.writer, "\n")
}
```

---

### 3. Mode Router

**File**: `cmd/simulator/main.go`

#### Enhanced Main Function
```go
func run() error {
    cfg, err := config.ParseFlags()
    if err != nil {
        return fmt.Errorf("config parsing: %w", err)
    }

    // Load seed (supports JSON and YAML)
    seedData, err := seed.LoadSeed(cfg.SeedPath)
    if err != nil {
        return fmt.Errorf("seed loading: %w", err)
    }

    // Route based on mode
    switch cfg.Mode {
    case "preview":
        return runPreviewMode(seedData, cfg)
    case "runtime":
        return runRuntimeMode(seedData, cfg)
    default:
        return fmt.Errorf("invalid mode: '%s' (valid: runtime, preview)", cfg.Mode)
    }
}

func runPreviewMode(seedData *seed.SeedData, cfg *config.Config) error {
    previewCfg := preview.Config{
        Days:       7,
        MaxCommits: 50,
        MaxEvents:  10,
    }

    p := preview.New(seedData, previewCfg, os.Stdout)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := p.Run(ctx); err != nil {
        return fmt.Errorf("preview failed: %w", err)
    }

    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println("Preview complete. Use -mode runtime for full generation.")
    return nil
}

func runRuntimeMode(seedData *seed.SeedData, cfg *config.Config) error {
    // Existing implementation unchanged
    // ...
}
```

---

## Data Structures

### SeedData (No Changes Required)
```go
type SeedData struct {
    Developers []Developer `json:"developers" yaml:"developers"`
}

type Developer struct {
    UserID           string       `json:"user_id" yaml:"user_id"`
    Email            string       `json:"email" yaml:"email"`
    WorkingHoursBand WorkingHours `json:"working_hours" yaml:"working_hours"`
    PreferredModels  []string     `json:"preferred_models" yaml:"preferred_models"`
    Velocity         string       `json:"velocity" yaml:"velocity"`
}

type WorkingHours struct {
    Start int `json:"start" yaml:"start"`
    End   int `json:"end" yaml:"end"`
}
```

**Key Point**: Using struct tags `json:` AND `yaml:` allows same structure for both formats.

---

## YAML Example

**File**: `testdata/preview_example.yaml`
```yaml
# cursor-sim seed file
# YAML format with comments for readability

developers:
  # Engineering team lead
  - user_id: alice
    email: alice@example.com
    working_hours:
      start: 9   # PST timezone (09:00)
      end: 17    # 5 PM
    preferred_models:
      - claude-sonnet-4.5
      - gpt-4o
    velocity: high  # Commits 2x more frequently

  # Backend engineer
  - user_id: bob
    email: bob@example.com
    working_hours:
      start: 10  # Later start
      end: 18
    preferred_models:
      - claude-opus-4
    velocity: medium

  # Frontend engineer
  - user_id: charlie
    email: charlie@example.com
    working_hours:
      start: 8   # Early bird
      end: 16
    preferred_models:
      - gpt-4-turbo
    velocity: medium
```

---

## Testing Strategy

### Unit Tests

**File**: `internal/seed/loader_test.go`
```go
func TestLoadSeed_JSON(t *testing.T) {
    seed, err := LoadSeed("testdata/valid_seed.json")
    require.NoError(t, err)
    assert.Equal(t, 2, len(seed.Developers))
}

func TestLoadSeed_YAML(t *testing.T) {
    seed, err := LoadSeed("testdata/valid_seed.yaml")
    require.NoError(t, err)
    assert.Equal(t, 2, len(seed.Developers))
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
```

**File**: `internal/preview/preview_test.go`
```go
func TestPreview_Run(t *testing.T) {
    seedData := &seed.SeedData{
        Developers: []seed.Developer{
            {UserID: "alice", Email: "alice@example.com", ...},
        },
    }

    var buf bytes.Buffer
    cfg := Config{Days: 7, MaxCommits: 10}
    p := New(seedData, cfg, &buf)

    err := p.Run(context.Background())
    require.NoError(t, err)

    output := buf.String()
    assert.Contains(t, output, "PREVIEW MODE")
    assert.Contains(t, output, "alice")
    assert.Contains(t, output, "Sample Commits")
}

func TestPreview_Validation(t *testing.T) {
    seedData := &seed.SeedData{
        Developers: []seed.Developer{
            {
                UserID: "alice",
                PreferredModels: []string{"invalid-model"},
                WorkingHoursBand: seed.WorkingHours{Start: 25, End: 30}, // Invalid
            },
        },
    }

    var buf bytes.Buffer
    p := New(seedData, Config{}, &buf)

    err := p.Run(context.Background())
    require.NoError(t, err) // Warnings, not errors

    output := buf.String()
    assert.Contains(t, output, "⚠️  Validation Warnings")
    assert.Contains(t, output, "Unknown model")
    assert.Contains(t, output, "Invalid start hour")
}
```

### Integration Tests

**File**: `cmd/simulator/main_test.go`
```go
func TestRun_PreviewMode(t *testing.T) {
    // Set flags programmatically
    cfg := &config.Config{
        Mode:     "preview",
        SeedPath: "../../testdata/valid_seed.yaml",
    }

    err := runPreviewMode(seedData, cfg)
    assert.NoError(t, err)
}
```

### E2E Tests

**File**: `test/e2e/preview_test.go`
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
    assert.Contains(t, outputStr, "Preview complete")
}
```

---

## Performance Considerations

### Preview Mode Optimization

1. **Limited Generation**:
   - 7 days max (vs 90-180 days in runtime)
   - 50 commits max (vs 500-2500)
   - Minimal event types (commits + model usage only)

2. **Fast Exit**:
   - No server startup overhead
   - No endpoint initialization
   - Direct stdout output

3. **Target**: < 5 seconds for typical seed

### YAML Parsing Overhead

- YAML parsing: ~50ms for 100KB file
- JSON parsing: ~30ms for same file
- **Acceptable**: 20ms difference negligible at startup

---

## Security Considerations

1. **File Path Validation**:
   - Validate seed path doesn't escape directory
   - Check file size < 10MB before parsing

2. **YAML Bomb Protection**:
   - Use `yaml.v3` with safe defaults
   - No custom YAML tags allowed
   - Limit document size

3. **No Credential Exposure**:
   - Seed files don't contain secrets
   - Preview output safe to log/share

---

## Backward Compatibility

### Guarantees

1. **Existing JSON Seeds**: Continue working identically
2. **Runtime Mode**: No behavior changes
3. **API Contracts**: Unaffected (preview doesn't start server)
4. **Flag Parsing**: `-mode runtime` remains default

### Migration Path

**Week 1**: Release with both formats
**Week 2**: Update documentation with YAML examples
**Week 3**: Monitor adoption, gather feedback
**Week 4+**: Consider YAML as recommended format

---

## Alternative Approaches Considered

### Alternative 1: Web UI for Preview
**Rejected**: Adds complexity, CLI is sufficient for target users

### Alternative 2: GraphQL Introspection-style Preview
**Rejected**: Overkill, simple text output is better

### Alternative 3: TOML Instead of YAML
**Rejected**: YAML more popular in data engineering, better Go library support

---

## Dependencies

| Dependency | Version | Purpose | License |
|------------|---------|---------|---------|
| `gopkg.in/yaml.v3` | v3.0.1 | YAML parsing | MIT |

**Binary Size Impact**: +200KB (acceptable)

---

## Rollout Plan

### Phase 1: YAML Support (2 hours)
- Add yaml.v3 dependency
- Enhance LoadSeed function
- Write unit tests
- Update SPEC.md

### Phase 2: Preview Mode Core (4 hours)
- Implement preview package
- Add mode router
- Write preview tests
- E2E validation

### Phase 3: Validation Framework (2 hours)
- Implement validators
- Warning display
- Test edge cases

### Phase 4: Documentation (2 hours)
- Update README.md
- Add YAML examples
- Update CLI help text
- Write usage guide

---

## Success Criteria

- [ ] YAML seeds load identically to JSON
- [ ] Preview mode completes in < 5 seconds
- [ ] All validation warnings display correctly
- [ ] 100% backward compatibility with JSON
- [ ] Zero runtime mode performance impact
- [ ] Documentation complete with examples

---

**Next Step**: Create task.md with atomic implementation tasks following TDD approach.
