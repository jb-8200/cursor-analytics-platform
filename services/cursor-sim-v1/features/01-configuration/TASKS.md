# Feature 1: Configuration & Initialization - Development Tasks

## Feature Overview
Accept and validate JSON configuration to control simulation parameters. This feature is the entry point for the simulator - it reads user-provided configuration, applies defaults, validates bounds, and initializes the application state.

---

## Task 1.1: Configuration Data Structure Definition

**User Story**:
> As a developer, I want to have well-defined Go structs that represent all configuration options, so that I can parse, validate, and use them throughout the application.

**Description**:
Define the complete configuration data structures that will hold all simulation parameters, API settings, and export options.

**Acceptance Criteria**:
- [ ] `Config` struct defined with auth, simulation, api, and output sections
- [ ] All nested structs properly defined (`AuthConfig`, `SimulationConfig`, `APIConfig`, `OutputConfig`)
- [ ] All fields have JSON tags for unmarshaling
- [ ] Documentation strings explain each field
- [ ] Type system prevents invalid enum values (use custom types if needed)
- [ ] Struct supports round-trip JSON marshaling/unmarshaling
- [ ] All 15+ configuration parameters are present

**Success Criteria**:
- Unit tests verify struct marshaling/unmarshaling
- Tests confirm all JSON tags are correct
- Example config JSON can be parsed without errors

**Files to Create**:
- `internal/config/types.go` - Configuration structs
- `internal/config/types_test.go` - Struct tests

**Example Test**:
```go
func TestConfigMarshalUnmarshal(t *testing.T) {
    // Test that config can be marshaled to JSON and back
    // Test that all fields are preserved
}
```

---

## Task 1.2: Configuration File Loader

**User Story**:
> As a user, I want to load configuration from a JSON file, so that I can provide custom simulation parameters without editing code.

**Description**:
Implement a configuration loader that reads JSON files and parses them into the Config struct.

**Acceptance Criteria**:
- [ ] `LoadConfigFromFile(path string)` function reads JSON file
- [ ] Returns `*Config` and `error` tuple
- [ ] Returns helpful error message if file not found
- [ ] Returns parsing error with line number if JSON is invalid
- [ ] Handles file read permissions errors gracefully
- [ ] Returns parsed Config with all fields populated

**Success Criteria**:
- Loads valid JSON config file successfully
- Rejects invalid JSON with descriptive error
- Handles missing file gracefully
- Handles permission denied gracefully
- No panics on edge cases

**Files to Create**:
- `internal/config/loader.go` - File loading logic
- `internal/config/loader_test.go` - Loader tests

**Test Fixtures**:
- `tests/fixtures/config-valid.json` - Valid sample config
- `tests/fixtures/config-invalid.json` - Malformed JSON
- `tests/fixtures/config-missing.json` - File that doesn't exist

**Example Test**:
```go
func TestLoadValidConfigFile(t *testing.T) {
    // Load valid config and verify all fields
}

func TestLoadInvalidJSON(t *testing.T) {
    // Load malformed JSON and check error message
}

func TestLoadMissingFile(t *testing.T) {
    // Try to load non-existent file, verify error
}
```

---

## Task 1.3: Configuration Defaults & Merging

**User Story**:
> As a developer, I want to apply sensible defaults to user configuration, so that users don't need to specify every parameter.

**Description**:
Implement logic to apply default values to configuration fields that weren't provided by the user.

**Acceptance Criteria**:
- [ ] `ApplyDefaults(cfg *Config)` function applies all defaults
- [ ] Existing non-zero values are not overwritten
- [ ] Returns modified Config
- [ ] All 10+ defaults are applied (see configuration spec)
- [ ] Defaults include: developers=50, velocity=medium, volatility=0.2, port=8080, etc.
- [ ] Region/division/team distributions sum to 1.0

**Success Criteria**:
- Empty Config gets all defaults applied
- Partially filled Config preserves user values
- Defaults match documentation
- All defaults within valid bounds

**Files to Use**:
- `internal/config/defaults.go` - Default values and logic

**Example Test**:
```go
func TestApplyDefaultsEmpty(t *testing.T) {
    cfg := &Config{}
    ApplyDefaults(cfg)
    assert.Equal(t, cfg.Simulation.Developers, 50)
    assert.Equal(t, cfg.API.Port, 8080)
}

func TestApplyDefaultsPreserveUserValues(t *testing.T) {
    cfg := &Config{
        Simulation: SimulationConfig{Developers: 100},
    }
    ApplyDefaults(cfg)
    assert.Equal(t, cfg.Simulation.Developers, 100)
}
```

---

## Task 1.4: Configuration Validation

**User Story**:
> As a user, I want clear validation errors when my configuration is invalid, so that I can quickly fix problems.

**Description**:
Implement comprehensive validation of all configuration parameters against their allowed bounds and constraints.

**Acceptance Criteria**:
- [ ] `ValidateConfig(cfg *Config)` function validates all parameters
- [ ] Returns `[]ValidationError` (list of all validation issues)
- [ ] Checks developers count is between 1-10000
- [ ] Checks PR velocity is one of: low, medium, high
- [ ] Checks volatility is between 0.0 and 1.0
- [ ] Checks port is between 1024-65535
- [ ] Checks regions, divisions, groups, teams are valid
- [ ] Checks region/division distributions sum to 1.0 (within 0.01 tolerance)
- [ ] Checks break_condition.value is positive if set
- [ ] Validates API credentials are not empty
- [ ] Returns empty list if all validations pass

**Success Criteria**:
- Valid config returns no errors
- Each invalid parameter produces specific error message
- Multiple errors collected and returned together
- Error messages guide user to fix (e.g., "port must be 1024-65535")

**Files to Use**:
- `internal/config/validator.go` - Validation logic
- `internal/config/validator_test.go` - Validator tests

**Example Test**:
```go
func TestValidateDevelopersOutOfBounds(t *testing.T) {
    cfg := &Config{
        Simulation: SimulationConfig{Developers: 20000},
    }
    errs := ValidateConfig(cfg)
    assert.Len(t, errs, 1)
    assert.Contains(t, errs[0].Message, "developers")
}

func TestValidateMultipleErrors(t *testing.T) {
    cfg := &Config{
        Simulation: SimulationConfig{
            Developers: 0,
            PRVelocity: "invalid",
            Volatility: 1.5,
        },
    }
    errs := ValidateConfig(cfg)
    assert.True(t, len(errs) >= 3)
}
```

---

## Task 1.5: Configuration Initialization Pipeline

**User Story**:
> As a developer, I want a single function to load, apply defaults, and validate configuration, so that I can initialize the app with one call.

**Description**:
Implement the complete initialization pipeline that combines loader, defaults, and validation into a single workflow.

**Acceptance Criteria**:
- [ ] `InitializeConfig(configPath string)` function orchestrates the pipeline
- [ ] Loads config from file
- [ ] Applies defaults if loading fails but fallback is allowed
- [ ] Validates config
- [ ] Returns `(*Config, error)` tuple
- [ ] Returns error if file missing AND no defaults allowed
- [ ] Returns error if validation fails (includes all validation errors)
- [ ] Logs each step of initialization
- [ ] Successful initialization returns fully validated Config

**Success Criteria**:
- Complete pipeline works end-to-end
- Error handling is clean and informative
- Logs indicate what step failed
- Invalid configs are rejected at initialization

**Files to Use**:
- `internal/config/init.go` - Initialization pipeline
- `internal/config/init_test.go` - Integration tests

**Example Test**:
```go
func TestInitializeConfigSuccess(t *testing.T) {
    cfg, err := InitializeConfig("tests/fixtures/config-valid.json")
    assert.NoError(t, err)
    assert.NotNil(t, cfg)
}

func TestInitializeConfigInvalid(t *testing.T) {
    cfg, err := InitializeConfig("tests/fixtures/config-invalid.json")
    assert.Error(t, err)
    assert.Nil(t, cfg)
    assert.Contains(t, err.Error(), "validation")
}
```

---

## Task 1.6: CLI Configuration Input Handling

**User Story**:
> As a user, I want to provide configuration via CLI arguments, stdin, or environment variables, so that I don't always need a config file.

**Description**:
Implement command-line argument parsing to accept configuration from multiple sources and merge them intelligently.

**Acceptance Criteria**:
- [ ] CLI flag `--config <path>` accepts config file path
- [ ] CLI flag `--developers <num>` overrides developers count
- [ ] CLI flag `--velocity <level>` overrides PR velocity
- [ ] CLI flag `--port <num>` overrides API port
- [ ] Environment variable `CURSOR_SIM_CONFIG` provides default config path
- [ ] Flags have higher priority than env vars
- [ ] Env vars have higher priority than defaults
- [ ] `--help` flag displays usage information
- [ ] `--version` flag displays version

**Success Criteria**:
- CLI flags parse correctly
- Invalid flag values produce error
- Help text is clear and complete
- Flag precedence works as specified

**Files to Create**:
- `internal/cli/flags.go` - Flag parsing
- `internal/cli/flags_test.go` - Flag tests

**Example Test**:
```go
func TestParseCLIFlags(t *testing.T) {
    args := []string{"--config", "test.json", "--developers", "100"}
    flags, err := ParseFlags(args)
    assert.NoError(t, err)
    assert.Equal(t, flags.ConfigPath, "test.json")
    assert.Equal(t, flags.DeveloperOverride, 100)
}
```

---

## Task 1.7: Configuration Documentation & Examples

**User Story**:
> As a new user, I want clear documentation and example configs, so that I can quickly get started.

**Description**:
Create comprehensive documentation and example configuration files for different use cases.

**Acceptance Criteria**:
- [ ] `README-CONFIG.md` documents all configuration options
- [ ] Each parameter has description, type, default, and bounds
- [ ] `config.example.json` shows a complete sample config
- [ ] `config.minimal.json` shows minimal required config
- [ ] `config.high-volume.json` shows high-velocity configuration
- [ ] `config.distributed.json` shows multi-region setup
- [ ] Each example config is valid and can be loaded

**Success Criteria**:
- Documentation is clear and complete
- All examples are valid JSON
- Examples cover common use cases
- User can use examples as templates

**Files to Create**:
- `README-CONFIG.md` - Configuration documentation
- `tests/fixtures/config.example.json` - Complete example
- `tests/fixtures/config.minimal.json` - Minimal config
- `tests/fixtures/config.high-volume.json` - High-velocity example
- `tests/fixtures/config.distributed.json` - Multi-region example

---

## Feature 1 Integration Test

**File**: `tests/integration_config_test.go`

```go
func TestConfigurationWorkflow(t *testing.T) {
    // 1. Load from file
    cfg, err := InitializeConfig("tests/fixtures/config-valid.json")
    assert.NoError(t, err)

    // 2. Verify defaults applied
    assert.Equal(t, cfg.API.Port, 8080)

    // 3. Verify validation passed
    errs := ValidateConfig(cfg)
    assert.Empty(t, errs)

    // 4. Verify all fields accessible
    assert.NotEmpty(t, cfg.Auth.APIKey)
    assert.True(t, cfg.Simulation.Developers > 0)
}
```

---

## Feature 1 Completion Checklist

- [ ] Task 1.1: Configuration structs defined and tested
- [ ] Task 1.2: Configuration file loader working
- [ ] Task 1.3: Defaults application working
- [ ] Task 1.4: Validation comprehensive and tested
- [ ] Task 1.5: Initialization pipeline complete
- [ ] Task 1.6: CLI argument parsing working
- [ ] Task 1.7: Documentation and examples complete
- [ ] Integration test passes
- [ ] Test coverage ≥85%
- [ ] All code builds without warnings
- [ ] README updated with configuration instructions

