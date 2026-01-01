# Feature 2: Core Data Models - Development Tasks

## Feature Overview
Define and implement all core data models that represent developers, commits, changes, and metrics. These models are the foundation for data generation and storage, and must be validated, serializable, and thread-safe.

---

## Task 2.1: Developer Model Definition & Validation

**User Story**:
> As a developer, I want a well-defined Developer model with validation, so that all synthetic developers have consistent and realistic attributes.

**Description**:
Implement the Developer struct with methods for creation, validation, and attribute access. Include support for realistic name generation and profile completeness.

**Acceptance Criteria**:
- [ ] `Developer` struct defined with all 12 fields (ID, Email, Name, Region, Division, Group, Team, Skills, ClientVersion, IsActive, CreatedAt, LastActiveAt)
- [ ] All fields have JSON tags for serialization
- [ ] ID is a unique identifier (UUID or similar)
- [ ] Email follows valid format (user@domain.tld)
- [ ] Region is one of: US, EU, APAC
- [ ] Division is one of: AGS, AT, ST
- [ ] Group is one of: TMOBILE, ATANT
- [ ] Team is one of: dev, support
- [ ] Skills is non-empty list of skill strings
- [ ] ClientVersion is valid semantic version (e.g., "0.35.0")
- [ ] Timestamps use time.Time with UTC timezone
- [ ] `NewDeveloper()` constructor creates valid developer with sensible defaults
- [ ] `ValidateDeveloper(d *Developer)` returns validation errors
- [ ] Struct implements custom JSON marshaling if needed

**Success Criteria**:
- Developer can be created and serialized
- All validation rules enforced
- Invalid developers are rejected
- Serialization round-trip preserves all data

**Files to Create**:
- `internal/models/developer.go` - Developer struct and methods
- `internal/models/developer_test.go` - Developer tests

**Example Test**:
```go
func TestNewDeveloper(t *testing.T) {
    dev := NewDeveloper("alice@example.com", "Alice")
    assert.NotEmpty(t, dev.ID)
    assert.Equal(t, dev.Email, "alice@example.com")
    assert.False(t, dev.IsActive)
}

func TestValidateDeveloper(t *testing.T) {
    dev := &Developer{Email: "invalid-email", Region: "INVALID"}
    errs := ValidateDeveloper(dev)
    assert.True(t, len(errs) >= 2)
}

func TestDeveloperSerialization(t *testing.T) {
    dev := NewDeveloper("bob@example.com", "Bob")
    json, err := json.Marshal(dev)
    assert.NoError(t, err)
    var dev2 Developer
    err = json.Unmarshal(json, &dev2)
    assert.NoError(t, err)
    assert.Equal(t, dev.ID, dev2.ID)
}
```

---

## Task 2.2: Commit Model Definition & Validation

**User Story**:
> As a developer, I want to represent Git commits with AI-assisted metrics, so that I can track TAB vs Composer contributions and line counts.

**Description**:
Implement the Commit struct with validation and methods for realistic commit representation.

**Acceptance Criteria**:
- [ ] `Commit` struct defined with 10 fields (Hash, Timestamp, Message, DeveloperID, Repository, Branch, LinesDelta, LinesFromTAB, LinesFromComposer, LinesNonAI, IngestionTime)
- [ ] All fields have JSON tags
- [ ] Hash is unique SHA-like identifier
- [ ] Timestamp and IngestionTime use UTC
- [ ] Repository follows naming convention (e.g., "team/project-name")
- [ ] Branch is valid git branch name
- [ ] Line counts are non-negative integers
- [ ] Sum of AI lines + NonAI lines ≤ LinesDelta
- [ ] TAB + Composer ≤ AI total lines
- [ ] `NewCommit()` constructor creates valid commit
- [ ] `ValidateCommit(c *Commit)` validates all constraints

**Success Criteria**:
- Valid commits pass validation
- Invalid commits are rejected with specific errors
- Line count constraints enforced
- Serialization works correctly

**Files to Create**:
- `internal/models/commit.go` - Commit struct and methods
- `internal/models/commit_test.go` - Commit tests

**Example Test**:
```go
func TestNewCommit(t *testing.T) {
    dev := NewDeveloper("dev@example.com", "Dev")
    commit := NewCommit("sha123abc...", dev.ID, "repo/project")
    assert.Equal(t, len(commit.Hash), 40)
    assert.Equal(t, commit.DeveloperID, dev.ID)
}

func TestValidateCommitLineCounts(t *testing.T) {
    commit := &Commit{
        LinesDelta:        100,
        LinesFromTAB:      60,
        LinesFromComposer: 50,  // Invalid: 60+50 > 100
    }
    errs := ValidateCommit(commit)
    assert.NotEmpty(t, errs)
}
```

---

## Task 2.3: Change Model Definition & Validation

**User Story**:
> As a developer, I want to track individual AI-suggested changes, so that I can analyze which models and sources generate the most accepted suggestions.

**Description**:
Implement the Change struct representing individual AI-suggested code changes accepted in commits.

**Acceptance Criteria**:
- [ ] `Change` struct defined with 11 fields (ChangeID, CommitHash, DeveloperID, Timestamp, Source, Model, FilesDelta, LinesAdded, LinesModified, IngestionTime)
- [ ] All fields have JSON tags
- [ ] ChangeID is unique identifier
- [ ] Source is one of: TAB, COMPOSER
- [ ] Model is one of: claude-3.5-sonnet, claude-opus-4.5, gpt-4, etc.
- [ ] FilesDelta is positive integer
- [ ] LinesAdded and LinesModified are non-negative
- [ ] Timestamps use UTC
- [ ] `NewChange()` constructor creates valid change
- [ ] `ValidateChange(c *Change)` validates all constraints
- [ ] ChangeID generation is deterministic (reproducible from inputs)

**Success Criteria**:
- Changes are uniquely identifiable
- All validation rules enforced
- Serialization works correctly
- Change IDs can be reproduced for consistency

**Files to Create**:
- `internal/models/change.go` - Change struct and methods
- `internal/models/change_test.go` - Change tests

**Example Test**:
```go
func TestNewChange(t *testing.T) {
    change := NewChange("sha123", "dev-id", "TAB", "claude-3.5-sonnet")
    assert.NotEmpty(t, change.ChangeID)
    assert.Equal(t, change.Source, "TAB")
}

func TestChangeValidation(t *testing.T) {
    change := &Change{Source: "INVALID", Files: -1}
    errs := ValidateChange(change)
    assert.True(t, len(errs) >= 2)
}

func TestChangeIDDeterministic(t *testing.T) {
    id1 := NewChange("sha", "dev", "TAB", "model").ChangeID
    id2 := NewChange("sha", "dev", "TAB", "model").ChangeID
    assert.Equal(t, id1, id2)
}
```

---

## Task 2.4: Daily Metrics Model Definition

**User Story**:
> As an analyst, I want to query aggregated daily metrics, so that I can see team-wide trends and performance indicators.

**Description**:
Implement the DailyMetrics struct for storing aggregated daily statistics.

**Acceptance Criteria**:
- [ ] `DailyMetrics` struct defined with 8 fields (Date, DAU, AgentEditsAccepted, TabCompletions, ComposerEdits, ModelUsage map, TopFileExtensions map, CommandUsage map)
- [ ] Date is stored without time (midnight UTC)
- [ ] All count fields are non-negative integers
- [ ] Maps are nil-safe and handle empty cases
- [ ] `NewDailyMetrics()` constructor initializes empty collections
- [ ] `ValidateDailyMetrics(m *DailyMetrics)` validates all constraints
- [ ] Methods for adding metrics (e.g., `AddAgentEdit()`, `AddTabCompletion()`)

**Success Criteria**:
- Metrics aggregate correctly
- Addition methods update maps properly
- Serialization preserves map data
- Validation enforces non-negative counts

**Files to Create**:
- `internal/models/metrics.go` - DailyMetrics and related structs
- `internal/models/metrics_test.go` - Metrics tests

**Example Test**:
```go
func TestNewDailyMetrics(t *testing.T) {
    metrics := NewDailyMetrics(time.Now())
    assert.Equal(t, metrics.DAU, 0)
    assert.Empty(t, metrics.ModelUsage)
}

func TestAddMetrics(t *testing.T) {
    metrics := NewDailyMetrics(time.Now())
    metrics.IncrementDAU()
    metrics.AddTabCompletion()
    assert.Equal(t, metrics.DAU, 1)
    assert.Equal(t, metrics.TabCompletions, 1)
}

func TestModelUsageTracking(t *testing.T) {
    metrics := NewDailyMetrics(time.Now())
    metrics.IncrementModelUsage("claude-3.5-sonnet")
    metrics.IncrementModelUsage("claude-3.5-sonnet")
    assert.Equal(t, metrics.ModelUsage["claude-3.5-sonnet"], 2)
}
```

---

## Task 2.5: Developer Metrics Model Definition

**User Story**:
> As an analyst, I want per-developer daily metrics, so that I can track individual contributor productivity and tool adoption.

**Description**:
Implement the DeveloperMetrics struct for per-developer aggregated statistics.

**Acceptance Criteria**:
- [ ] `DeveloperMetrics` struct defined with 7 fields (DeveloperID, Date, AgentEdits, TabCompletions, ComposerEdits, ModelsUsed map, FileExtensions map)
- [ ] Date is stored without time
- [ ] All counts are non-negative
- [ ] Maps are nil-safe
- [ ] `NewDeveloperMetrics()` creates new instance
- [ ] `ValidateDeveloperMetrics()` validates constraints
- [ ] Methods for incrementing counts

**Success Criteria**:
- Metrics track per-developer data correctly
- Maps aggregate properly
- Validation works as expected

**Files to Use**:
- `internal/models/metrics.go` - Add DeveloperMetrics

**Example Test**:
```go
func TestNewDeveloperMetrics(t *testing.T) {
    metrics := NewDeveloperMetrics("dev-id", time.Now())
    assert.Equal(t, metrics.DeveloperID, "dev-id")
    assert.Equal(t, metrics.AgentEdits, 0)
}

func TestDeveloperMetricsIncrement(t *testing.T) {
    metrics := NewDeveloperMetrics("dev-id", time.Now())
    metrics.IncrementAgentEdits()
    metrics.AddFileExtension(".go")
    assert.Equal(t, metrics.AgentEdits, 1)
    assert.Equal(t, metrics.FileExtensions[".go"], 1)
}
```

---

## Task 2.6: Time Helper Utilities

**User Story**:
> As a developer, I want consistent timestamp handling across all models, so that dates are always in UTC and consistent.

**Description**:
Implement utility functions for time handling to ensure consistency.

**Acceptance Criteria**:
- [ ] `NowUTC()` returns current time in UTC
- [ ] `StartOfDay(t time.Time)` returns midnight UTC for given date
- [ ] `EndOfDay(t time.Time)` returns 23:59:59 UTC for given date
- [ ] `ParseDateRange(start, end string)` parses date strings including relative dates like "7d", "today"
- [ ] All functions handle edge cases (invalid input, timezones)
- [ ] Functions are well-tested with various inputs

**Success Criteria**:
- Time handling is consistent
- Date parsing supports all specified formats
- Timezone issues are prevented
- Edge cases handled gracefully

**Files to Create**:
- `internal/models/time_helpers.go` - Time utility functions
- `internal/models/time_helpers_test.go` - Time tests

**Example Test**:
```go
func TestNowUTC(t *testing.T) {
    now := NowUTC()
    assert.Equal(t, now.Location(), time.UTC)
}

func TestStartOfDay(t *testing.T) {
    t1 := time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)
    start := StartOfDay(t1)
    assert.Equal(t, start.Hour(), 0)
    assert.Equal(t, start.Minute(), 0)
    assert.Equal(t, start.Day(), 15)
}

func TestParseDateRangeRelative(t *testing.T) {
    start, end, err := ParseDateRange("7d", "today")
    assert.NoError(t, err)
    assert.True(t, start.Before(end))
}
```

---

## Task 2.7: Model Constants & Enums

**User Story**:
> As a developer, I want defined constants for all enum values, so that I can avoid string duplication and catch typos at compile time.

**Description**:
Define constants for all enumerated values used throughout the models.

**Acceptance Criteria**:
- [ ] Region constants: RegionUS, RegionEU, RegionAPAC
- [ ] Division constants: DivisionAGS, DivisionAT, DivisionST
- [ ] Group constants: GroupTMOBILE, GroupATANT
- [ ] Team constants: TeamDev, TeamSupport
- [ ] Source constants: SourceTAB, SourceCOMPOSER
- [ ] Model constants for AI models: ModelClaude35Sonnet, ModelClaudeOpus45, ModelGPT4, etc.
- [ ] Velocity constants: VelocityLow, VelocityMedium, VelocityHigh
- [ ] Slice of valid values for each enum
- [ ] Validation functions like `IsValidRegion(r string) bool`

**Success Criteria**:
- No magic strings in code
- Type safety for enums
- Easy to add new enum values
- Validation functions work correctly

**Files to Create**:
- `internal/models/constants.go` - All constants and validators
- `internal/models/constants_test.go` - Constant tests

**Example Test**:
```go
func TestIsValidRegion(t *testing.T) {
    assert.True(t, IsValidRegion(RegionUS))
    assert.False(t, IsValidRegion("INVALID"))
}

func TestValidModelsList(t *testing.T) {
    assert.Greater(t, len(ValidModels), 0)
    assert.Contains(t, ValidModels, ModelClaude35Sonnet)
}
```

---

## Task 2.8: Model JSON Marshaling & Custom Types

**User Story**:
> As a developer, I want models to serialize/deserialize properly with custom JSON handling, so that we can handle special cases.

**Description**:
Implement custom JSON marshaling/unmarshaling for models where needed (e.g., time formatting, enum validation).

**Acceptance Criteria**:
- [ ] All models marshal to valid JSON
- [ ] Timestamp fields use RFC3339 format in JSON
- [ ] Enums are validated on unmarshal
- [ ] Invalid enum values produce clear error messages
- [ ] Round-trip serialization preserves all data
- [ ] Custom types don't interfere with normal operations

**Success Criteria**:
- Models serialize correctly
- Deserialization validates data
- Error messages are helpful
- No data loss in round-trip

**Files to Create**:
- `internal/models/json.go` - Custom marshaling methods
- `internal/models/json_test.go` - JSON marshaling tests

**Example Test**:
```go
func TestDeveloperJSONMarshal(t *testing.T) {
    dev := NewDeveloper("test@example.com", "Test")
    data, err := json.Marshal(dev)
    assert.NoError(t, err)

    var dev2 Developer
    err = json.Unmarshal(data, &dev2)
    assert.NoError(t, err)
    assert.Equal(t, dev.ID, dev2.ID)
}

func TestCommitJSONTimestamp(t *testing.T) {
    commit := NewCommit("sha", "dev-id", "repo")
    data, _ := json.Marshal(commit)
    assert.Contains(t, string(data), "T")  // RFC3339 format
    assert.Contains(t, string(data), "Z")  // UTC timezone
}
```

---

## Feature 2 Integration Test

**File**: `tests/integration_models_test.go`

```go
func TestModelInteroperability(t *testing.T) {
    // Create developer
    dev := NewDeveloper("dev@example.com", "Developer")
    assert.NoError(t, ValidateDeveloperErrors(ValidateDeveloper(dev)))

    // Create commit
    commit := NewCommit("sha123", dev.ID, "repo/name")
    assert.NoError(t, ValidateCommitErrors(ValidateCommit(commit)))

    // Create change
    change := NewChange("sha123", dev.ID, SourceTAB, ModelClaude35Sonnet)
    assert.NoError(t, ValidateChangeErrors(ValidateChange(change)))

    // Create metrics
    metrics := NewDailyMetrics(time.Now())
    metrics.IncrementDAU()

    // All serializable
    devJSON, _ := json.Marshal(dev)
    commitJSON, _ := json.Marshal(commit)
    changeJSON, _ := json.Marshal(change)
    metricsJSON, _ := json.Marshal(metrics)

    assert.True(t, len(devJSON) > 0)
    assert.True(t, len(commitJSON) > 0)
    assert.True(t, len(changeJSON) > 0)
    assert.True(t, len(metricsJSON) > 0)
}
```

---

## Feature 2 Completion Checklist

- [ ] Task 2.1: Developer model complete and tested
- [ ] Task 2.2: Commit model complete and tested
- [ ] Task 2.3: Change model complete and tested
- [ ] Task 2.4: Daily metrics model complete and tested
- [ ] Task 2.5: Developer metrics model complete and tested
- [ ] Task 2.6: Time helpers complete and tested
- [ ] Task 2.7: Constants and validators complete
- [ ] Task 2.8: JSON marshaling working correctly
- [ ] Integration test passes
- [ ] Test coverage ≥85%
- [ ] All models properly documented
- [ ] No magic strings in code

