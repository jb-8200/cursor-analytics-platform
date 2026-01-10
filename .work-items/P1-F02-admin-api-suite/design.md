# Design Document: P1-F02 Admin API Suite

**Feature**: P1-F02 Admin API Suite
**Phase**: P1 (Foundation)
**Created**: January 10, 2026
**Status**: Design Approved

---

## Architecture Overview

The Admin API Suite extends cursor-sim with 5 integrated components providing configuration, seed management, and observability:

```
┌─────────────────────────────────────────────────────────────┐
│                  Admin API Suite (P1-F02)                   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Part 1: Environment Variables (Startup Configuration)     │
│    • CURSOR_SIM_DEVELOPERS, CURSOR_SIM_MONTHS,            │
│      CURSOR_SIM_MAX_COMMITS                               │
│    • Precedence: CLI flags > Env vars > Defaults          │
│                                                             │
│  Part 2: Runtime Admin APIs                                │
│    • POST /admin/regenerate (append/override modes)        │
│    • POST /admin/seed (JSON/YAML/CSV upload)              │
│    • GET /admin/seed/presets (predefined configs)          │
│    • GET /admin/config (inspect configuration)             │
│    • GET /admin/stats (comprehensive analytics)            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## API Endpoints

### 1. POST /admin/regenerate

**Purpose**: Regenerate simulation data with new parameters without restarting

**Request**:
```json
{
  "mode": "override",        // "append" or "override"
  "days": 400,               // 1-3650
  "velocity": "high",        // "low", "medium", "high"
  "developers": 1200,        // 0-10000 (0 = use seed count)
  "max_commits": 500         // 0-100000 (0 = unlimited)
}
```

**Response**:
```json
{
  "status": "success",
  "mode": "override",
  "data_cleaned": true,
  "commits_added": 600000,
  "prs_added": 60000,
  "reviews_added": 120000,
  "issues_added": 20000,
  "total_commits": 600000,
  "total_prs": 60000,
  "total_developers": 1200,
  "duration": "8.5s",
  "config": {
    "days": 400,
    "velocity": "high",
    "developers": 1200,
    "max_commits": 500
  }
}
```

**Validation**:
- `mode`: Must be "append" or "override"
- `days`: 1-3650
- `velocity`: "low", "medium", "high"
- `developers`: 0-10000
- `max_commits`: 0-100000

**Error Responses**:
- `400 Bad Request`: Invalid parameters
- `401 Unauthorized`: Missing/invalid API key
- `405 Method Not Allowed`: Non-POST request
- `500 Internal Server Error`: Generation failed

---

### 2. POST /admin/seed

**Purpose**: Upload/swap seed file dynamically

**Request**:
```json
{
  "format": "json",          // "json", "yaml", or "csv"
  "data": "{...JSON...}",    // Seed data as string
  "regenerate": true,        // Auto-regenerate after upload
  "regenerate_config": {     // Optional (if regenerate=true)
    "mode": "override",
    "days": 180,
    "velocity": "high"
  }
}
```

**Response**:
```json
{
  "status": "success",
  "seed_loaded": true,
  "developers": 100,
  "repositories": 20,
  "teams": ["Backend", "Frontend", "DevOps"],
  "divisions": ["Engineering", "Infrastructure"],
  "organizations": ["acme-corp"],
  "regenerated": true,
  "generate_stats": {
    "status": "success",
    "commits_added": 45000,
    "total_commits": 45000
  }
}
```

**Validation**:
- Seed data must have valid `version`, `developers[]`, `repositories[]` fields
- Org hierarchy: Organization → Division → Team → Developer
- Supported formats: JSON, YAML, CSV

---

### 3. GET /admin/seed/presets

**Purpose**: List predefined seed configurations

**Response**:
```json
{
  "presets": [
    {
      "name": "small-team",
      "description": "Small team (2 developers, 2 repos, 1 region)",
      "developers": 2,
      "teams": 2,
      "regions": ["US"]
    },
    {
      "name": "enterprise",
      "description": "Enterprise (100 developers, 20 repos, 3 regions)",
      "developers": 100,
      "teams": 10,
      "regions": ["US", "EU", "APAC"]
    }
  ]
}
```

---

### 4. GET /admin/config

**Purpose**: Inspect current runtime configuration and seed structure

**Response**:
```json
{
  "generation": {
    "days": 90,
    "velocity": "medium",
    "developers": 50,
    "max_commits": 1000
  },
  "seed": {
    "version": "1.0",
    "developers": 2,
    "repositories": 2,
    "organizations": ["acme-corp"],
    "divisions": ["Engineering"],
    "teams": ["Backend", "Frontend"],
    "regions": ["US", "EU"],
    "by_seniority": {"senior": 1, "mid": 1},
    "by_region": {"US": 1, "EU": 1},
    "by_team": {"Backend": 1, "Frontend": 1}
  },
  "external_sources": {
    "harvey": {"enabled": true, "models": ["gpt-4", "claude-3-sonnet"]},
    "copilot": {"enabled": true, "total_licenses": 50, "active_users": 35},
    "qualtrics": {"enabled": true, "survey_id": "SV_aitools_q1_2026", "response_count": 150}
  },
  "server": {
    "port": 8080,
    "version": "2.0.0",
    "uptime": "5m30s"
  }
}
```

---

### 5. GET /admin/stats

**Purpose**: Retrieve comprehensive statistics about generated data

**Query Parameters**:
- `include_timeseries=true` (optional): Include time series data

**Response**:
```json
{
  "generation": {
    "total_commits": 4500,
    "total_prs": 450,
    "total_reviews": 900,
    "total_issues": 150,
    "total_developers": 100,
    "data_size": "5.2 MB"
  },
  "developers": {
    "by_seniority": {"junior": 20, "mid": 50, "senior": 30},
    "by_region": {"US": 50, "EU": 30, "APAC": 20},
    "by_team": {"Backend": 40, "Frontend": 35, "DevOps": 25},
    "by_activity": {"low": 15, "medium": 50, "high": 35}
  },
  "quality": {
    "avg_revert_rate": 0.02,
    "avg_hotfix_rate": 0.08,
    "avg_code_survival_30d": 0.85,
    "avg_review_thoroughness": 0.75,
    "avg_pr_iterations": 1.5
  },
  "variance": {
    "commits_std_dev": 15.2,
    "pr_size_std_dev": 75.5,
    "cycle_time_std_dev": 2.3
  },
  "performance": {
    "last_generation_time": "2.34s",
    "memory_usage": "125 MB",
    "storage_efficiency": "95%"
  },
  "organization": {
    "teams": ["Backend", "Frontend", "DevOps"],
    "divisions": ["Engineering", "Infrastructure"],
    "repositories": ["acme-corp/payment-service", "acme-corp/web-app"]
  },
  "time_series": {  // Only if ?include_timeseries=true
    "commits_per_day": [15, 18, 12, 20, 16],
    "prs_per_day": [3, 2, 4, 3, 2],
    "avg_cycle_time": [4.5, 3.2, 5.1, 4.0, 3.8]
  }
}
```

---

## Data Models

### Request Models

```go
// RegenerateRequest (internal/api/models/regenerate.go)
type RegenerateRequest struct {
    Mode        string `json:"mode"`         // "append" or "override"
    Days        int    `json:"days"`         // 1-3650
    Velocity    string `json:"velocity"`     // "low", "medium", "high"
    Developers  int    `json:"developers"`   // 0-10000
    MaxCommits  int    `json:"max_commits"`  // 0-100000
}

// SeedUploadRequest (internal/api/models/seed.go)
type SeedUploadRequest struct {
    Data             string              `json:"data"`
    Format           string              `json:"format"`  // "json", "yaml", "csv"
    Regenerate       bool                `json:"regenerate"`
    RegenerateConfig *RegenerateRequest  `json:"regenerate_config,omitempty"`
}
```

### Response Models

```go
// RegenerateResponse (internal/api/models/regenerate.go)
type RegenerateResponse struct {
    Status          string `json:"status"`
    Mode            string `json:"mode"`
    DataCleaned     bool   `json:"data_cleaned"`
    CommitsAdded    int    `json:"commits_added"`
    PRsAdded        int    `json:"prs_added"`
    ReviewsAdded    int    `json:"reviews_added"`
    IssuesAdded     int    `json:"issues_added"`
    TotalCommits    int    `json:"total_commits"`
    TotalPRs        int    `json:"total_prs"`
    TotalDevelopers int    `json:"total_developers"`
    Duration        string `json:"duration"`
    Config          struct {
        Days       int    `json:"days"`
        Velocity   string `json:"velocity"`
        Developers int    `json:"developers"`
        MaxCommits int    `json:"max_commits"`
    } `json:"config"`
}

// ConfigResponse (internal/api/models/config.go)
type ConfigResponse struct {
    Generation struct {
        Days       int    `json:"days"`
        Velocity   string `json:"velocity"`
        Developers int    `json:"developers"`
        MaxCommits int    `json:"max_commits"`
    } `json:"generation"`

    Seed struct {
        Version       string         `json:"version"`
        Developers    int            `json:"developers"`
        Repositories  int            `json:"repositories"`
        Organizations []string       `json:"organizations"`
        Divisions     []string       `json:"divisions"`
        Teams         []string       `json:"teams"`
        Regions       []string       `json:"regions"`
        BySeniority   map[string]int `json:"by_seniority"`
        ByRegion      map[string]int `json:"by_region"`
        ByTeam        map[string]int `json:"by_team"`
    } `json:"seed"`

    ExternalSources struct {
        Harvey struct {
            Enabled bool     `json:"enabled"`
            Models  []string `json:"models"`
        } `json:"harvey"`
        Copilot struct {
            Enabled       bool `json:"enabled"`
            TotalLicenses int  `json:"total_licenses"`
            ActiveUsers   int  `json:"active_users"`
        } `json:"copilot"`
        Qualtrics struct {
            Enabled       bool   `json:"enabled"`
            SurveyID      string `json:"survey_id"`
            ResponseCount int    `json:"response_count"`
        } `json:"qualtrics"`
    } `json:"external_sources"`

    Server struct {
        Port    int    `json:"port"`
        Version string `json:"version"`
        Uptime  string `json:"uptime"`
    } `json:"server"`
}

// StatsResponse (internal/api/models/stats.go)
type StatsResponse struct {
    Generation struct {
        TotalCommits    int    `json:"total_commits"`
        TotalPRs        int    `json:"total_prs"`
        TotalReviews    int    `json:"total_reviews"`
        TotalIssues     int    `json:"total_issues"`
        TotalDevelopers int    `json:"total_developers"`
        DataSize        string `json:"data_size"`
    } `json:"generation"`

    Developers struct {
        BySeniority map[string]int `json:"by_seniority"`
        ByRegion    map[string]int `json:"by_region"`
        ByTeam      map[string]int `json:"by_team"`
        ByActivity  map[string]int `json:"by_activity"`
    } `json:"developers"`

    Quality struct {
        AvgRevertRate         float64 `json:"avg_revert_rate"`
        AvgHotfixRate         float64 `json:"avg_hotfix_rate"`
        AvgCodeSurvival       float64 `json:"avg_code_survival_30d"`
        AvgReviewThoroughness float64 `json:"avg_review_thoroughness"`
        AvgIterations         float64 `json:"avg_pr_iterations"`
    } `json:"quality"`

    Variance struct {
        CommitsStdDev   float64 `json:"commits_std_dev"`
        PRSizeStdDev    float64 `json:"pr_size_std_dev"`
        CycleTimeStdDev float64 `json:"cycle_time_std_dev"`
    } `json:"variance"`

    Performance struct {
        LastGenerationTime string `json:"last_generation_time"`
        MemoryUsage        string `json:"memory_usage"`
        StorageEfficiency  string `json:"storage_efficiency"`
    } `json:"performance"`

    Organization struct {
        Teams        []string `json:"teams"`
        Divisions    []string `json:"divisions"`
        Repositories []string `json:"repositories"`
    } `json:"organization"`

    TimeSeries *struct {
        CommitsPerDay []int     `json:"commits_per_day,omitempty"`
        PRsPerDay     []int     `json:"prs_per_day,omitempty"`
        AvgCycleTime  []float64 `json:"avg_cycle_time,omitempty"`
    } `json:"time_series,omitempty"`
}
```

---

## Storage Layer Changes

### Interface Extension

```go
// store.go
type Store interface {
    // ... existing methods ...

    // NEW: Clear all data (for override mode)
    ClearAllData() error

    // NEW: Get current storage statistics
    GetStats() StorageStats
}

// NEW: Storage statistics structure
type StorageStats struct {
    Commits        int `json:"commits"`
    PullRequests   int `json:"pull_requests"`
    Reviews        int `json:"reviews"`
    Issues         int `json:"issues"`
    Developers     int `json:"developers"`
    ModelUsage     int `json:"model_usage"`
    ClientVersions int `json:"client_versions"`
    FileExtensions int `json:"file_extensions"`
    MCPTools       int `json:"mcp_tools"`
    Commands       int `json:"commands"`
    Plans          int `json:"plans"`
    AskModes       int `json:"ask_modes"`
}
```

### Implementation

```go
// memory.go
func (m *MemoryStore) ClearAllData() error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Reset all maps and slices to initial state
    m.commits = make([]*models.Commit, 0, 10000)
    m.commitsByHash = make(map[string]*models.Commit)
    m.commitsByUser = make(map[string][]*models.Commit)
    m.commitsByRepo = make(map[string][]*models.Commit)
    m.prsByRepo = make(map[string]map[int]*models.PullRequest)
    m.prsByID = make(map[int]*models.PullRequest)
    m.reviews = make([]*models.Review, 0, 1000)
    m.issues = make([]*models.Issue, 0, 1000)
    m.modelUsage = make([]*models.ModelUsageEvent, 0, 1000)
    m.clientVersions = make([]*models.ClientVersionEvent, 0, 1000)
    m.fileExtensions = make([]*models.FileExtensionEvent, 0, 1000)
    m.mcpTools = make([]*models.MCPToolEvent, 0, 1000)
    m.commands = make([]*models.CommandEvent, 0, 1000)
    m.plans = make([]*models.PlanEvent, 0, 1000)
    m.askModes = make([]*models.AskModeEvent, 0, 1000)
    m.developers = make([]*models.Developer, 0, 100)

    return nil
}

func (m *MemoryStore) GetStats() StorageStats {
    m.mu.RLock()
    defer m.mu.RUnlock()

    return StorageStats{
        Commits:        len(m.commits),
        PullRequests:   len(m.prsByID),
        Reviews:        len(m.reviews),
        Issues:         len(m.issues),
        Developers:     len(m.developers),
        ModelUsage:     len(m.modelUsage),
        ClientVersions: len(m.clientVersions),
        FileExtensions: len(m.fileExtensions),
        MCPTools:       len(m.mcpTools),
        Commands:       len(m.commands),
        Plans:          len(m.plans),
        AskModes:       len(m.askModes),
    }
}
```

---

## Configuration Changes

### Environment Variable Precedence

```
CLI flags > Environment Variables > Defaults

Example:
1. Default: days=90
2. CURSOR_SIM_DAYS=180 → days=180
3. -days 30 → days=30 (CLI flag wins)
```

### config.go Changes

```go
// After line 108 in parseFlagsWithArgs()

// Apply environment variable overrides for GenerationParams
if v := os.Getenv("CURSOR_SIM_DEVELOPERS"); v != "" {
    if dev, err := strconv.Atoi(v); err == nil {
        cfg.GenParams.Developers = dev
    }
}
if v := os.Getenv("CURSOR_SIM_MONTHS"); v != "" {
    if months, err := strconv.Atoi(v); err == nil {
        cfg.GenParams.Days = months * 30
        cfg.Days = months * 30  // Keep cfg.Days in sync
    }
}
if v := os.Getenv("CURSOR_SIM_MAX_COMMITS"); v != "" {
    if maxCommits, err := strconv.Atoi(v); err == nil {
        cfg.GenParams.MaxCommits = maxCommits
    }
}
```

---

## File Structure

### New Files

```
services/cursor-sim/internal/api/models/
├── regenerate.go         # Regenerate request/response
├── seed.go               # Seed upload request/response
├── config.go             # Config inspection response
└── stats.go              # Statistics response

services/cursor-sim/internal/api/cursor/
├── admin_regenerate.go   # POST /admin/regenerate handler
├── admin_seed.go         # POST /admin/seed + GET /admin/seed/presets
├── admin_config.go       # GET /admin/config handler
└── admin_stats.go        # GET /admin/stats handler

Test files:
├── admin_regenerate_test.go
├── admin_seed_test.go
├── admin_config_test.go
└── admin_stats_test.go
```

### Modified Files

```
services/cursor-sim/internal/config/
├── config.go             # Add env var parsing (lines 110-127)
└── config_test.go        # Add env var tests (6 new tests)

services/cursor-sim/internal/storage/
├── store.go              # Add ClearAllData(), GetStats() to interface
└── memory.go             # Implement new methods (~60 lines)

services/cursor-sim/internal/server/
└── router.go             # Register 5 new admin endpoints

services/cursor-sim/cmd/simulator/
└── main.go               # Pass seedData to router

Docker/deployment:
├── Dockerfile            # Remove hardcoded -days, -velocity
├── docker-compose.yml    # Fix env var list
├── .env.example          # Document all env vars
└── tools/deploy-cursor-sim.sh  # Add new env vars
```

---

## Security Considerations

### Authentication

All Admin API endpoints use existing Basic Auth with API key:
- Same authentication mechanism as existing endpoints
- No additional API keys or roles required
- Thread-safe operations using existing storage mutex locks

### Validation

Strict parameter validation prevents abuse:
- Days: 1-3650 (max ~10 years)
- Developers: 0-10000
- MaxCommits: 0-100000
- Velocity: whitelist (low/medium/high)
- Mode: whitelist (append/override)

### Thread Safety

- Use existing `sync.RWMutex` in MemoryStore
- ClearAllData() acquires write lock
- GetStats() acquires read lock
- Regenerate handler serializes generation operations

### Future Enhancements (Out of Scope)

- Rate limiting for admin endpoints
- Audit logging of admin operations
- Role-based access control
- Admin operation webhooks

---

## Performance Considerations

### Memory Usage

Large datasets (1200 developers × 400 days) estimate:
- Commits: ~600,000 × 500 bytes = ~300 MB
- PRs: ~60,000 × 1KB = ~60 MB
- Reviews: ~120,000 × 300 bytes = ~36 MB
- Total: ~400-500 MB in-memory

### Generation Time

Estimated for 1200 developers × 400 days:
- Commit generation: ~5-8 seconds
- PR generation: ~2-3 seconds
- Review generation: ~1-2 seconds
- Feature events: ~1-2 seconds
- Total: ~10-15 seconds

### Concurrent Reads

- Regeneration blocks all reads during override mode
- Consider progress streaming for very large datasets (future enhancement)

---

## Testing Strategy

### Unit Tests

1. **config_test.go** (6 new tests):
   - TestParseFlags_EnvironmentOverrides_Developers
   - TestParseFlags_EnvironmentOverrides_Months
   - TestParseFlags_EnvironmentOverrides_MaxCommits
   - TestParseFlags_EnvironmentOverrides_AllGenParams
   - TestParseFlags_CLIOverridesEnvironment
   - TestParseFlags_EnvVarsDoNotTriggerMixedMode

2. **admin_regenerate_test.go** (4 tests):
   - TestRegenerateAppendMode
   - TestRegenerateOverrideMode
   - TestRegenerateInvalidMode
   - TestRegenerateInvalidVelocity

3. **admin_seed_test.go** (6 tests):
   - TestSeedUpload_JSON
   - TestSeedUpload_YAML
   - TestSeedUpload_CSV
   - TestSeedUpload_InvalidFormat
   - TestSeedUpload_WithRegenerate
   - TestGetSeedPresets

4. **admin_config_test.go** (2 tests):
   - TestGetConfig
   - TestGetConfig_ExternalSources

5. **admin_stats_test.go** (3 tests):
   - TestGetStats_Basic
   - TestGetStats_WithTimeSeries
   - TestGetStats_Calculations

### E2E Tests

1. Test environment variable override (Docker)
2. Test override mode with 1200 developers, 400 days
3. Test append mode (cumulative data)
4. Test seed upload (JSON/YAML/CSV)
5. Test config inspection
6. Test stats retrieval
7. Test parameter validation
8. Test authentication

### Coverage Target

- New code: 80%+ coverage
- Critical paths (regenerate, seed upload): 90%+ coverage

---

## Deployment Strategy

### Docker Compose

1. Update `.env` file with desired parameters
2. Run `docker-compose up -d cursor-sim`
3. Verify via `GET /admin/config`

### GCP Cloud Run

1. Set environment variables in deployment script
2. Run `tools/deploy-cursor-sim.sh production`
3. Verify via `GET /admin/config` using service URL

### Rollback Plan

If issues arise:
1. No schema changes, so rollback is safe
2. Existing endpoints unchanged
3. New endpoints can be disabled via router if needed

---

## Dependencies

### Internal

- `internal/config` - Environment variable parsing
- `internal/storage` - ClearAllData(), GetStats()
- `internal/generator` - Commit, PR, Review generators
- `internal/seed` - Seed data loading (JSON/YAML/CSV)
- `internal/api` - Handler factory pattern
- `internal/server` - Router registration

### External

- `gopkg.in/yaml.v3` - YAML parsing (already used)
- `encoding/json` - JSON parsing (stdlib)
- `encoding/csv` - CSV parsing (stdlib)
- `sync` - Thread-safe storage operations (stdlib)

---

## Success Criteria

1. ✅ Environment variables configure generation parameters
2. ✅ Runtime reconfiguration without container restart
3. ✅ Seed file upload/swap (JSON/YAML/CSV)
4. ✅ Configuration inspection via GET /admin/config
5. ✅ Comprehensive statistics via GET /admin/stats
6. ✅ Thread-safe operations
7. ✅ 80%+ test coverage
8. ✅ Complete API documentation in SPEC.md
9. ✅ Works in Docker Compose and GCP Cloud Run

---

## ADRs (Architecture Decision Records)

### ADR-001: Use Append/Override Modes Instead of Replace

**Decision**: Implement two distinct modes (append/override) instead of a single replace operation.

**Rationale**:
- **Append**: Allows extending time series data without losing historical context
- **Override**: Clean slate for testing different scenarios
- Users can choose based on use case

**Alternatives Considered**:
- Single "regenerate" operation (always override) - Less flexible
- Snapshot/restore mechanism - Too complex for MVP

### ADR-002: Thread-Safe Operations Using Existing Mutex

**Decision**: Use existing `sync.RWMutex` in MemoryStore for thread safety.

**Rationale**:
- Reuses proven concurrency pattern
- No additional complexity
- Blocks reads during override mode (acceptable trade-off)

**Alternatives Considered**:
- Copy-on-write - Complex, more memory
- Lock-free data structures - Overkill for current scale

### ADR-003: Basic Auth for Admin Endpoints

**Decision**: Use existing Basic Auth with API key for admin endpoints.

**Rationale**:
- Consistent with existing endpoints
- Simple deployment (no additional secrets)
- Sufficient for private deployments

**Alternatives Considered**:
- Separate admin API key - Unnecessary complexity
- OAuth/JWT - Overkill for current use case

### ADR-004: Support JSON, YAML, and CSV Seed Formats

**Decision**: Accept three formats for maximum flexibility.

**Rationale**:
- **JSON**: Comprehensive (full org structure, external data sources)
- **YAML**: Human-friendly for manual editing
- **CSV**: Simple for basic use cases (user_id, email, name)

**Alternatives Considered**:
- JSON only - Less user-friendly
- Add TOML - Not common in this domain

---

**Design Status**: Approved for Implementation
**Next Step**: Create task breakdown (task.md) with subagent assignments
