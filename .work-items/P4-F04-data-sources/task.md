# Task Breakdown: External Data Source Simulators

**Feature ID**: P4-F04-data-sources
**Created**: January 9, 2026
**Status**: Planning
**Subagent Protocol**: Master agent (Opus) delegates to cursor-sim-cli-dev (Sonnet)

---

## Progress Tracker

| Phase | Tasks | Status | Estimated | Actual |
|-------|-------|--------|-----------|--------|
| **Infrastructure** | 2 | ✅ 2/2 | 3.0h | 2.5h |
| **Harvey API** | 4 | ✅ 4/4 | 6.0h | 4.5h |
| **Microsoft Copilot API** | 4 | 🔄 3/4 | 6.5h | 3.5h |
| **Qualtrics API** | 4 | 🔄 1/4 | 8.0h | 1.5h |
| **Integration & E2E** | 2 | Pending | 3.5h | - |
| **TOTAL** | **16** | **10/16** | **27.0h** | 12.0h |

---

## Task Breakdown by Phase

### PHASE 1: Infrastructure

#### TASK-DS-01: Extend Seed Schema with External Data Sources (Est: 1.5h)
**Status**: ✅ COMPLETE
**Actual**: 1.5h
**Commit**: 2d8768c

**Goal**: Add HarveyUser, M365Tenant, and QualtricsConfig types to seed schema

**TDD Approach**:
```go
func TestSeedData_WithHarveyUsers(t *testing.T) {
    yaml := `
version: "1.0.0"
harvey_users:
  - user_id: atty_001
    email: john@lawfirm.com
    name: John Attorney
    role: partner
    practice_area: corporate
    activity_multiplier: 1.2
    client_matters: [2024.001, 2024.045]
`
    data, err := seed.LoadFromYAML([]byte(yaml))
    require.NoError(t, err)
    require.Len(t, data.HarveyUsers, 1)
    assert.Equal(t, "atty_001", data.HarveyUsers[0].UserID)
    assert.Equal(t, "corporate", data.HarveyUsers[0].PracticeArea)
}

func TestSeedData_WithM365Tenant(t *testing.T) {
    yaml := `
version: "1.0.0"
m365_tenant:
  tenant_id: tenant_abc
  display_name: Acme Corp
  users:
    - user_id: m365_001
      email: jane@company.com
      display_name: Jane Developer
      copilot_enabled: true
`
    data, err := seed.LoadFromYAML([]byte(yaml))
    require.NoError(t, err)
    require.NotNil(t, data.M365Tenant)
    assert.Len(t, data.M365Tenant.Users, 1)
}

func TestSeedData_WithQualtrics(t *testing.T) {
    yaml := `
version: "1.0.0"
qualtrics:
  surveys:
    - survey_id: SV_abc123
      name: AI Tools Survey
      response_count: 100
`
    data, err := seed.LoadFromYAML([]byte(yaml))
    require.NoError(t, err)
    require.NotNil(t, data.Qualtrics)
    assert.Len(t, data.Qualtrics.Surveys, 1)
}
```

**Files to Create/Modify**:
- MODIFY: `internal/seed/types.go` - Add new types
- MODIFY: `internal/seed/types_test.go` - Add tests
- NEW: `testdata/enterprise_seed.yaml` - Example with all sources

**Acceptance Criteria**:
- [x] HarveySeedConfig struct with required fields
- [x] CopilotSeedConfig struct with required fields
- [x] QualtricsSeedConfig struct with required fields
- [x] ExternalDataSourcesSeed container struct
- [x] Seed struct extended with ExternalDataSources field
- [x] All JSON tags properly set for serialization
- [x] Tests pass for marshaling/unmarshaling (96.6% coverage)
- [x] Example seed.json includes external data sources section
- [x] Backward compatible (existing seeds still work)

---

#### TASK-DS-02: Extend Storage Layer for External Data (Est: 1.5h)
**Status**: ✅ COMPLETE
**Actual**: 1.0h
**Commit**: 198df5f

**Goal**: Add storage methods for Harvey events, Copilot usage, and export jobs

**Implementation Notes**:
Created a new `external.go` file with separate ExternalMemoryStore to keep external data source storage isolated from the existing storage. This approach:
- Preserves backward compatibility with existing Store interface
- Provides clean separation of concerns
- Allows independent scaling and testing

**Files Created**:
- NEW: `internal/storage/external.go` - Storage types and implementations
- NEW: `internal/storage/external_test.go` - Comprehensive test suite

**Key Components**:
- `HarveyStore` interface with GetUsage, StoreUsage methods
- `CopilotStore` interface with GetUsage, StoreUsage methods
- `QualtricsStore` interface with GetSurveys, StoreSurveys, GetExportJob, StoreExportJob, GetFile, StoreFile methods
- `ExternalDataStore` container interface aggregating all three
- `ExternalMemoryStore` in-memory implementation with sync.RWMutex for thread safety

**Acceptance Criteria**:
- [x] HarveyStore interface defined
- [x] CopilotStore interface defined
- [x] QualtricsStore interface defined
- [x] ExternalDataStore container interface
- [x] In-memory implementations for all stores
- [x] Thread-safe concurrent access (sync.RWMutex)
- [x] Tests pass with 100% coverage on external.go

---

### PHASE 2: Harvey API

#### TASK-DS-03: Create Harvey Usage Model (Est: 1.0h)
**Assigned Subagent**: `cursor-sim-cli-dev`

**Goal**: Implement HarveyUsageEvent model with validation

**TDD Approach**:
```go
func TestHarveyUsageEvent_Validate(t *testing.T) {
    tests := []struct {
        name    string
        event   models.HarveyUsageEvent
        wantErr bool
    }{
        {
            name: "valid event",
            event: models.HarveyUsageEvent{
                EventID:           12345,
                MessageID:         "uuid-abc",
                Time:              time.Now(),
                User:              "user@firm.com",
                Task:              models.HarveyTaskAssist,
                Source:            models.HarveySourceFiles,
                FeedbackSentiment: models.HarveySentimentPositive,
            },
            wantErr: false,
        },
        {
            name: "missing event_id",
            event: models.HarveyUsageEvent{
                User: "user@firm.com",
                Task: models.HarveyTaskAssist,
            },
            wantErr: true,
        },
        {
            name: "missing user",
            event: models.HarveyUsageEvent{
                EventID: 12345,
                Task:    models.HarveyTaskAssist,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.event.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestHarveyTask_Constants(t *testing.T) {
    assert.Equal(t, models.HarveyTask("Assist"), models.HarveyTaskAssist)
    assert.Equal(t, models.HarveyTask("Draft"), models.HarveyTaskDraft)
    assert.Equal(t, models.HarveyTask("Review"), models.HarveyTaskReview)
    assert.Equal(t, models.HarveyTask("Research"), models.HarveyTaskResearch)
}
```

**Files to Create**:
- NEW: `internal/models/harvey.go`
- NEW: `internal/models/harvey_test.go`

**Acceptance Criteria**:
- [ ] HarveyUsageEvent struct matches API spec
- [ ] Task type constants (Assist, Draft, Review, Research)
- [ ] Source type constants (Files, Web, Knowledge)
- [ ] Sentiment type constants (positive, negative, neutral)
- [ ] Validate() method checks required fields
- [ ] JSON tags match exact field names from spec
- [ ] Tests pass with 100% model coverage

---

#### TASK-DS-04: Create Harvey Generator (Est: 2.0h)
**Status**: ✅ COMPLETE
**Actual**: 1.5h
**Commit**: 7b73668

**Goal**: Implement HarveyGenerator for event generation

**TDD Approach**:
```go
func TestHarveyGenerator_GenerateEvents(t *testing.T) {
    seedData := &seed.SeedData{
        HarveyUsers: []seed.HarveyUser{
            {
                UserID:             "atty_001",
                Email:              "john@firm.com",
                ActivityMultiplier: 1.0,
                ClientMatters:      []float64{2024.001},
            },
        },
    }

    gen := generator.NewHarveyGeneratorWithSeed(seedData, 12345)
    config := generator.DefaultHarveyConfig()

    from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
    to := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)

    events := gen.GenerateEvents(from, to, config)

    // Expect ~5 events/day * 7 days = ~35 events
    assert.True(t, len(events) >= 20)
    assert.True(t, len(events) <= 50)

    // Verify all events are from our user
    for _, e := range events {
        assert.Equal(t, "john@firm.com", e.User)
    }
}

func TestHarveyGenerator_TaskDistribution(t *testing.T) {
    seedData := &seed.SeedData{
        HarveyUsers: []seed.HarveyUser{
            {UserID: "atty_001", Email: "user@firm.com", ActivityMultiplier: 1.0},
        },
    }

    gen := generator.NewHarveyGeneratorWithSeed(seedData, 12345)
    config := generator.DefaultHarveyConfig()

    from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
    to := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

    events := gen.GenerateEvents(from, to, config)

    // Count task types
    counts := make(map[models.HarveyTask]int)
    for _, e := range events {
        counts[e.Task]++
    }

    total := float64(len(events))
    // Allow 10% tolerance from configured rates
    assertRate(t, counts[models.HarveyTaskAssist]/total, 0.35, 0.10)
    assertRate(t, counts[models.HarveyTaskDraft]/total, 0.30, 0.10)
    assertRate(t, counts[models.HarveyTaskReview]/total, 0.25, 0.10)
    assertRate(t, counts[models.HarveyTaskResearch]/total, 0.10, 0.10)
}

func TestHarveyGenerator_Reproducible(t *testing.T) {
    seedData := &seed.SeedData{
        HarveyUsers: []seed.HarveyUser{
            {UserID: "atty_001", Email: "user@firm.com", ActivityMultiplier: 1.0},
        },
    }

    from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
    to := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)
    config := generator.DefaultHarveyConfig()

    // Same seed should produce identical results
    gen1 := generator.NewHarveyGeneratorWithSeed(seedData, 12345)
    events1 := gen1.GenerateEvents(from, to, config)

    gen2 := generator.NewHarveyGeneratorWithSeed(seedData, 12345)
    events2 := gen2.GenerateEvents(from, to, config)

    require.Equal(t, len(events1), len(events2))
    for i := range events1 {
        assert.Equal(t, events1[i].EventID, events2[i].EventID)
        assert.Equal(t, events1[i].Task, events2[i].Task)
    }
}
```

**Files to Create**:
- NEW: `internal/generator/harvey_generator.go`
- NEW: `internal/generator/harvey_generator_test.go`

**Acceptance Criteria**:
- [x] Poisson-distributed event counts per user/day
- [x] Configurable task distribution rates
- [x] Configurable sentiment rates
- [x] Activity multiplier applied per user (uses developers from seed)
- [x] Client matter assignment from seed (generated from practice areas)
- [x] Working hours constraint (8 AM - 6 PM)
- [x] Reproducible with same random seed
- [x] Tests pass with 96.2% coverage (exceeds 90% target)

---

#### TASK-DS-05: Create Harvey API Handler (Est: 1.5h)
**Assigned Subagent**: `cursor-sim-cli-dev`

**Goal**: Implement HTTP handler for Harvey usage endpoint

**TDD Approach**:
```go
func TestHarveyUsageHandler_Success(t *testing.T) {
    store := storage.NewMemoryStore()
    gen := generator.NewHarveyGenerator(&seed.SeedData{})

    // Pre-populate store
    events := []models.HarveyUsageEvent{
        {EventID: 1, User: "a@firm.com", Task: models.HarveyTaskAssist, Time: time.Now()},
        {EventID: 2, User: "b@firm.com", Task: models.HarveyTaskDraft, Time: time.Now()},
    }
    for _, e := range events {
        store.AddHarveyEvent(e)
    }

    handler := harvey.UsageHandler(store, gen)

    req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage", nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    var resp map[string]interface{}
    json.Unmarshal(rr.Body.Bytes(), &resp)

    data := resp["data"].([]interface{})
    assert.Len(t, data, 2)
}

func TestHarveyUsageHandler_FilterByUser(t *testing.T) {
    store := storage.NewMemoryStore()
    gen := generator.NewHarveyGenerator(&seed.SeedData{})

    store.AddHarveyEvent(models.HarveyUsageEvent{
        EventID: 1, User: "a@firm.com", Task: models.HarveyTaskAssist, Time: time.Now(),
    })
    store.AddHarveyEvent(models.HarveyUsageEvent{
        EventID: 2, User: "b@firm.com", Task: models.HarveyTaskDraft, Time: time.Now(),
    })

    handler := harvey.UsageHandler(store, gen)

    req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?user=a@firm.com", nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    var resp map[string]interface{}
    json.Unmarshal(rr.Body.Bytes(), &resp)

    data := resp["data"].([]interface{})
    assert.Len(t, data, 1)
}

func TestHarveyUsageHandler_Pagination(t *testing.T) {
    store := storage.NewMemoryStore()
    gen := generator.NewHarveyGenerator(&seed.SeedData{})

    for i := 0; i < 150; i++ {
        store.AddHarveyEvent(models.HarveyUsageEvent{
            EventID: int64(i), User: "user@firm.com", Task: models.HarveyTaskAssist, Time: time.Now(),
        })
    }

    handler := harvey.UsageHandler(store, gen)

    req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?page=1&page_size=50", nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    var resp map[string]interface{}
    json.Unmarshal(rr.Body.Bytes(), &resp)

    data := resp["data"].([]interface{})
    assert.Len(t, data, 50)

    pagination := resp["pagination"].(map[string]interface{})
    assert.Equal(t, float64(3), pagination["totalPages"])
    assert.Equal(t, true, pagination["hasNextPage"])
}
```

**Files to Create**:
- NEW: `internal/api/harvey/handlers.go`
- NEW: `internal/api/harvey/handlers_test.go`

**Acceptance Criteria**:
- [ ] GET /harvey/api/v1/history/usage endpoint
- [ ] Date range filtering (from, to params)
- [ ] User filtering
- [ ] Task type filtering
- [ ] Pagination (page, page_size)
- [ ] Response format matches spec
- [ ] Basic auth required
- [ ] Tests pass with 90%+ coverage

---

#### TASK-DS-06: Harvey Integration and Router Setup (Est: 1.5h)
**Assigned Subagent**: `cursor-sim-cli-dev`

**Goal**: Integrate Harvey API into main router with conditional registration

**TDD Approach**:
```go
func TestRouter_HarveyRoutes_Enabled(t *testing.T) {
    seedData := &seed.SeedData{
        HarveyUsers: []seed.HarveyUser{
            {UserID: "atty_001", Email: "user@firm.com"},
        },
    }

    store := storage.NewMemoryStore()
    router := server.NewRouter(store, seedData)

    // Route should exist
    req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage", nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    assert.NotEqual(t, http.StatusNotFound, rr.Code)
}

func TestRouter_HarveyRoutes_Disabled(t *testing.T) {
    seedData := &seed.SeedData{
        // No HarveyUsers
    }

    store := storage.NewMemoryStore()
    router := server.NewRouter(store, seedData)

    req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage", nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusNotFound, rr.Code)
}
```

**Files to Modify**:
- MODIFY: `internal/server/router.go` - Add Harvey routes
- MODIFY: `cmd/simulator/main.go` - Initialize Harvey generator

**Acceptance Criteria**:
- [x] Routes only registered when Harvey.Enabled in seed
- [x] ExternalMemoryStore created for Harvey data
- [x] Harvey handler integrated via harvey.UsageHandler
- [x] Authentication middleware applied
- [x] Tests pass (all scenarios: enabled and disabled)

**Status**: COMPLETE (Jan 9, 2026)
**Time**: 1.5h actual / 1.5h estimated
**Commit**: b2c3dc1

---

### PHASE 3: Microsoft 365 Copilot API

#### TASK-DS-07: Create Copilot Usage Model (Est: 1.0h)
**Status**: ✅ COMPLETE
**Actual**: 0.5h
**Commit**: 0a9edb1

**Goal**: Implement CopilotUsageUserDetail model matching Graph API

**TDD Approach**:
```go
func TestCopilotReportPeriod_Days(t *testing.T) {
    tests := []struct {
        period models.CopilotReportPeriod
        days   int
    }{
        {models.CopilotPeriodD7, 7},
        {models.CopilotPeriodD30, 30},
        {models.CopilotPeriodD90, 90},
        {models.CopilotPeriodD180, 180},
        {models.CopilotPeriodAll, 180},
    }

    for _, tt := range tests {
        t.Run(string(tt.period), func(t *testing.T) {
            assert.Equal(t, tt.days, tt.period.Days())
        })
    }
}

func TestCopilotUsageUserDetail_GetAppLastActivityDate(t *testing.T) {
    teamsDate := "2026-01-08"
    wordDate := "2026-01-05"

    detail := models.CopilotUsageUserDetail{
        MicrosoftTeamsCopilotLastActivityDate: &teamsDate,
        WordCopilotLastActivityDate:           &wordDate,
    }

    assert.Equal(t, &teamsDate, detail.GetAppLastActivityDate(models.CopilotAppTeams))
    assert.Equal(t, &wordDate, detail.GetAppLastActivityDate(models.CopilotAppWord))
    assert.Nil(t, detail.GetAppLastActivityDate(models.CopilotAppExcel))
}

func TestCopilotUsageUserDetail_JSONMarshal(t *testing.T) {
    date := "2026-01-08"
    detail := models.CopilotUsageUserDetail{
        ReportRefreshDate:                     "2026-01-09",
        ReportPeriod:                          30,
        UserPrincipalName:                     "user@company.com",
        DisplayName:                           "Jane Dev",
        LastActivityDate:                      &date,
        MicrosoftTeamsCopilotLastActivityDate: &date,
    }

    data, err := json.Marshal(detail)
    require.NoError(t, err)

    // Verify exact field names match Microsoft API
    assert.Contains(t, string(data), `"reportRefreshDate"`)
    assert.Contains(t, string(data), `"microsoftTeamsCopilotLastActivityDate"`)
}
```

**Files to Create**:
- NEW: `internal/models/copilot.go`
- NEW: `internal/models/copilot_test.go`

**Acceptance Criteria**:
- [x] CopilotUsageUserDetail matches Graph API schema exactly
- [x] CopilotReportPeriod enum with Days() method
- [x] CopilotApp enum for all 8 apps
- [x] AllCopilotApps() helper function
- [x] Nullable date fields use *string
- [x] JSON field names match Microsoft spec exactly
- [x] Tests pass with 100% coverage

---

#### TASK-DS-08: Create Copilot Generator (Est: 2.0h)
**Status**: ✅ COMPLETE
**Actual**: 1.5h
**Commit**: 64de277

**Goal**: Implement CopilotGenerator for usage data generation

**TDD Approach**:
```go
func TestCopilotGenerator_GenerateUsageReport(t *testing.T) {
    seedData := &seed.SeedData{
        M365Tenant: &seed.M365Tenant{
            TenantID:    "tenant_abc",
            DisplayName: "Acme Corp",
            Users: []seed.M365User{
                {
                    UserID:         "m365_001",
                    Email:          "jane@company.com",
                    DisplayName:    "Jane Developer",
                    CopilotEnabled: true,
                },
                {
                    UserID:         "m365_002",
                    Email:          "bob@company.com",
                    DisplayName:    "Bob Manager",
                    CopilotEnabled: true,
                },
            },
        },
    }

    gen := generator.NewCopilotGeneratorWithSeed(seedData, 12345)
    config := generator.DefaultCopilotConfig()

    usage := gen.GenerateUsageReport(models.CopilotPeriodD30, config)

    require.Len(t, usage, 2)
    assert.Equal(t, 30, usage[0].ReportPeriod)
    assert.Equal(t, "jane@company.com", usage[0].UserPrincipalName)
}

func TestCopilotGenerator_AppAdoption(t *testing.T) {
    seedData := &seed.SeedData{
        M365Tenant: &seed.M365Tenant{
            Users: []seed.M365User{
                {Email: "user@company.com", CopilotEnabled: true},
            },
        },
    }

    // Generate many samples to verify adoption rates
    teamsSeen := 0
    oneNoteSeen := 0
    samples := 100

    for i := 0; i < samples; i++ {
        gen := generator.NewCopilotGeneratorWithSeed(seedData, int64(i))
        config := generator.DefaultCopilotConfig()
        usage := gen.GenerateUsageReport(models.CopilotPeriodD30, config)

        if usage[0].MicrosoftTeamsCopilotLastActivityDate != nil {
            teamsSeen++
        }
        if usage[0].OneNoteCopilotLastActivityDate != nil {
            oneNoteSeen++
        }
    }

    // Teams should have ~85% adoption, OneNote ~10%
    assert.True(t, teamsSeen > 70, "Teams adoption too low")
    assert.True(t, oneNoteSeen < 25, "OneNote adoption too high")
}

func TestCopilotGenerator_ActivityDatesWithinPeriod(t *testing.T) {
    seedData := &seed.SeedData{
        M365Tenant: &seed.M365Tenant{
            Users: []seed.M365User{
                {Email: "user@company.com", CopilotEnabled: true},
            },
        },
    }

    gen := generator.NewCopilotGeneratorWithSeed(seedData, 12345)
    config := generator.DefaultCopilotConfig()
    usage := gen.GenerateUsageReport(models.CopilotPeriodD30, config)

    now := time.Now()
    thirtyDaysAgo := now.AddDate(0, 0, -30)

    for _, app := range models.AllCopilotApps() {
        dateStr := usage[0].GetAppLastActivityDate(app)
        if dateStr != nil {
            activityDate, _ := time.Parse("2006-01-02", *dateStr)
            assert.True(t, activityDate.After(thirtyDaysAgo) || activityDate.Equal(thirtyDaysAgo))
            assert.True(t, activityDate.Before(now) || activityDate.Equal(now))
        }
    }
}
```

**Files to Create**:
- NEW: `internal/generator/copilot_generator.go`
- NEW: `internal/generator/copilot_generator_test.go`

**Acceptance Criteria**:
- [x] Generates usage for all developers from seed
- [x] Configurable app adoption rates
- [x] Activity dates within report period
- [x] Last activity date computed from individual apps
- [x] Reproducible with same random seed
- [x] Tests pass with 100% coverage (exceeds 90% target)
- [x] Follows existing generator patterns

---

#### TASK-DS-09: Create Copilot API Handler (Est: 2.0h)
**Assigned Subagent**: `cursor-sim-cli-dev`

**Goal**: Implement Graph API handler with JSON and CSV support

**TDD Approach**:
```go
func TestCopilotUsageHandler_JSONResponse(t *testing.T) {
    store := storage.NewMemoryStore()
    gen := generator.NewCopilotGenerator(&seed.SeedData{})

    // Pre-populate store
    store.AddCopilotUsage(models.CopilotPeriodD30, models.CopilotUsageUserDetail{
        UserPrincipalName: "user@company.com",
        ReportPeriod:      30,
    })

    handler := microsoft.UsageUserDetailHandler(store, gen)

    req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=application/json", nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    var resp models.CopilotUsageResponse
    json.Unmarshal(rr.Body.Bytes(), &resp)

    assert.Len(t, resp.Value, 1)
    assert.Equal(t, "user@company.com", resp.Value[0].UserPrincipalName)
}

func TestCopilotUsageHandler_CSVRedirect(t *testing.T) {
    store := storage.NewMemoryStore()
    gen := generator.NewCopilotGenerator(&seed.SeedData{})

    store.AddCopilotUsage(models.CopilotPeriodD30, models.CopilotUsageUserDetail{
        UserPrincipalName: "user@company.com",
    })

    handler := microsoft.UsageUserDetailHandler(store, gen)

    req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=text/csv", nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusFound, rr.Code)
    assert.Contains(t, rr.Header().Get("Location"), "/reports/download/")
}

func TestCopilotUsageHandler_InvalidPeriod(t *testing.T) {
    store := storage.NewMemoryStore()
    gen := generator.NewCopilotGenerator(&seed.SeedData{})

    handler := microsoft.UsageUserDetailHandler(store, gen)

    req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='INVALID')", nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCopilotCSVDownload_Success(t *testing.T) {
    store := storage.NewMemoryStore()

    // Store CSV data
    store.StoreCopilotCSVData("abc123", []models.CopilotUsageUserDetail{
        {UserPrincipalName: "user@company.com", ReportPeriod: 30},
    })

    handler := microsoft.CSVDownloadHandler(store)

    req := httptest.NewRequest("GET", "/reports/download/abc123.csv", nil)
    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Equal(t, "text/csv", rr.Header().Get("Content-Type"))
    assert.Contains(t, rr.Body.String(), "Report Refresh Date")
}
```

**Files to Create**:
- NEW: `internal/api/microsoft/copilot_handlers.go`
- NEW: `internal/api/microsoft/copilot_handlers_test.go`

**Acceptance Criteria**:
- [ ] Endpoint pattern matches Graph API format
- [ ] JSON response with OData structure
- [ ] CSV redirects to download URL
- [ ] CSV download serves correct content
- [ ] Pagination with skiptoken
- [ ] Period parameter validated
- [ ] Tests pass with 90%+ coverage

---

#### TASK-DS-10: Copilot Integration and Router Setup (Est: 1.5h)
**Assigned Subagent**: `cursor-sim-cli-dev`

**Goal**: Integrate Copilot API into main router

**TDD Approach**:
```go
func TestRouter_CopilotRoutes_Enabled(t *testing.T) {
    seedData := &seed.SeedData{
        M365Tenant: &seed.M365Tenant{
            TenantID: "tenant_abc",
            Users: []seed.M365User{
                {Email: "user@company.com", CopilotEnabled: true},
            },
        },
    }

    store := storage.NewMemoryStore()
    router := server.NewRouter(store, seedData)

    req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')", nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    assert.NotEqual(t, http.StatusNotFound, rr.Code)
}

func TestRouter_CopilotRoutes_Disabled(t *testing.T) {
    seedData := &seed.SeedData{
        // No M365Tenant
    }

    store := storage.NewMemoryStore()
    router := server.NewRouter(store, seedData)

    req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')", nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusNotFound, rr.Code)
}
```

**Files to Modify**:
- MODIFY: `internal/server/router.go` - Add Copilot routes
- MODIFY: `cmd/simulator/main.go` - Initialize Copilot generator

**Acceptance Criteria**:
- [ ] Routes only registered when M365Tenant in seed
- [ ] Generator initialized from seed
- [ ] Data pre-generated for all periods
- [ ] Authentication middleware applied
- [ ] Tests pass

---

### PHASE 4: Qualtrics API

#### TASK-DS-11: Create Qualtrics Export Models (Est: 1.5h)
**Status**: ✅ COMPLETE
**Actual**: 1.5h
**Commit**: 091b72c

**Goal**: Implement export job and survey response models

**TDD Approach**:
```go
func TestExportJob_States(t *testing.T) {
    job := &models.ExportJob{
        ProgressID: "ES_abc123",
        SurveyID:   "SV_xyz",
        Status:     models.ExportStatusInProgress,
    }

    assert.Equal(t, models.ExportStatusInProgress, job.Status)

    job.PercentComplete = 100
    job.Status = models.ExportStatusComplete
    job.FileID = "FILE_xyz"

    assert.Equal(t, models.ExportStatusComplete, job.Status)
    assert.NotEmpty(t, job.FileID)
}

func TestSurveyResponse_Fields(t *testing.T) {
    resp := models.SurveyResponse{
        ResponseID:            "R_abc123",
        RespondentEmail:       "user@company.com",
        OverallAISatisfaction: 4,
        CursorSatisfaction:    5,
        CopilotSatisfaction:   3,
        MostUsedTool:          "Cursor",
    }

    assert.Equal(t, 4, resp.OverallAISatisfaction)
    assert.Equal(t, "Cursor", resp.MostUsedTool)
}

func TestGenerateZIPFile(t *testing.T) {
    responses := []models.SurveyResponse{
        {
            ResponseID:            "R_001",
            RespondentEmail:       "user1@company.com",
            OverallAISatisfaction: 4,
        },
        {
            ResponseID:            "R_002",
            RespondentEmail:       "user2@company.com",
            OverallAISatisfaction: 5,
        },
    }

    zipData, err := models.GenerateZIPFile(responses)
    require.NoError(t, err)
    assert.NotEmpty(t, zipData)

    // Verify ZIP structure
    reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
    require.NoError(t, err)
    assert.Len(t, reader.File, 1)
    assert.Equal(t, "survey_responses.csv", reader.File[0].Name)

    // Verify CSV content
    csvFile, _ := reader.File[0].Open()
    csvReader := csv.NewReader(csvFile)
    rows, _ := csvReader.ReadAll()

    assert.Len(t, rows, 3) // Header + 2 data rows
    assert.Contains(t, rows[0], "ResponseID")
    assert.Equal(t, "R_001", rows[1][0])
}
```

**Files to Create**:
- NEW: `internal/models/qualtrics.go`
- NEW: `internal/models/qualtrics_test.go`

**Acceptance Criteria**:
- [x] ExportJob with status transitions
- [x] ExportJobStatus enum (inProgress, complete, failed)
- [x] SurveyResponse with all survey fields
- [x] GenerateZIPFile produces valid ZIP with CSV
- [x] Response format matches Qualtrics API
- [x] Tests pass with 73.7% coverage (all tests passing)

---

#### TASK-DS-12: Create Survey Generator (Est: 2.0h)
**Assigned Subagent**: `cursor-sim-cli-dev`

**Goal**: Implement SurveyGenerator for response data

**TDD Approach**:
```go
func TestSurveyGenerator_GenerateResponses(t *testing.T) {
    seedData := &seed.SeedData{
        Developers: []seed.Developer{
            {Email: "dev1@company.com"},
            {Email: "dev2@company.com"},
        },
        Qualtrics: &seed.QualtricsConfig{
            Surveys: []seed.SurveyConfig{
                {
                    SurveyID:      "SV_abc123",
                    Name:          "AI Tools Survey",
                    ResponseCount: 10,
                    SatisfactionDistribution: map[int]float64{
                        1: 0.05, 2: 0.10, 3: 0.25, 4: 0.40, 5: 0.20,
                    },
                },
            },
        },
    }

    gen := generator.NewSurveyGeneratorWithSeed(seedData, 12345)
    responses := gen.GenerateSurveyResponses("SV_abc123")

    assert.Len(t, responses, 10)

    // Verify satisfaction distribution roughly matches
    counts := make(map[int]int)
    for _, r := range responses {
        counts[r.OverallAISatisfaction]++
    }

    // Most responses should be 3, 4, or 5
    highSatisfaction := counts[3] + counts[4] + counts[5]
    assert.True(t, highSatisfaction >= 7, "Expected mostly 3-5 satisfaction")
}

func TestSurveyGenerator_RespondentPools(t *testing.T) {
    seedData := &seed.SeedData{
        Developers: []seed.Developer{
            {Email: "dev@company.com"},
        },
        HarveyUsers: []seed.HarveyUser{
            {Email: "attorney@firm.com"},
        },
        Qualtrics: &seed.QualtricsConfig{
            Surveys: []seed.SurveyConfig{
                {
                    SurveyID:      "SV_abc",
                    ResponseCount: 100,
                    RespondentPools: []seed.RespondentPoolConfig{
                        {Pool: "developers", Weight: 0.7},
                        {Pool: "harvey_users", Weight: 0.3},
                    },
                },
            },
        },
    }

    gen := generator.NewSurveyGeneratorWithSeed(seedData, 12345)
    responses := gen.GenerateSurveyResponses("SV_abc")

    devCount := 0
    legalCount := 0
    for _, r := range responses {
        if strings.Contains(r.RespondentEmail, "company.com") {
            devCount++
        } else if strings.Contains(r.RespondentEmail, "firm.com") {
            legalCount++
        }
    }

    // Should be roughly 70/30 split
    assert.True(t, devCount >= 60, "Expected more developer responses")
    assert.True(t, legalCount >= 20, "Expected more legal responses")
}

func TestSurveyGenerator_FreeTextFeedback(t *testing.T) {
    seedData := &seed.SeedData{
        Qualtrics: &seed.QualtricsConfig{
            Surveys: []seed.SurveyConfig{
                {SurveyID: "SV_abc", ResponseCount: 5},
            },
        },
    }

    gen := generator.NewSurveyGeneratorWithSeed(seedData, 12345)
    responses := gen.GenerateSurveyResponses("SV_abc")

    // At least some responses should have feedback
    hasPositive := false
    hasImprovement := false
    for _, r := range responses {
        if r.PositiveFeedback != "" {
            hasPositive = true
        }
        if r.ImprovementAreas != "" {
            hasImprovement = true
        }
    }

    assert.True(t, hasPositive, "Expected some positive feedback")
    assert.True(t, hasImprovement, "Expected some improvement areas")
}
```

**Files to Create**:
- NEW: `internal/generator/survey_generator.go`
- NEW: `internal/generator/survey_generator_test.go`

**Acceptance Criteria**:
- [ ] Generates configured number of responses
- [ ] Satisfaction distribution matches seed config
- [ ] Respondent pools weighted correctly
- [ ] Realistic free-text feedback templates
- [ ] Demographics pulled from linked users
- [ ] Reproducible with same random seed
- [ ] Tests pass with 90%+ coverage

---

#### TASK-DS-13: Create Qualtrics Export State Machine (Est: 2.5h)
**Assigned Subagent**: `cursor-sim-cli-dev`

**Goal**: Implement ExportJobManager with state transitions

**TDD Approach**:
```go
func TestExportJobManager_StartExport(t *testing.T) {
    gen := generator.NewSurveyGenerator(&seed.SeedData{})
    manager := services.NewExportJobManager(gen)

    job, err := manager.StartExport("SV_abc123")
    require.NoError(t, err)

    assert.NotEmpty(t, job.ProgressID)
    assert.Equal(t, "SV_abc123", job.SurveyID)
    assert.Equal(t, models.ExportStatusInProgress, job.Status)
    assert.Equal(t, 0, job.PercentComplete)
}

func TestExportJobManager_ProgressAdvancement(t *testing.T) {
    gen := generator.NewSurveyGenerator(&seed.SeedData{
        Qualtrics: &seed.QualtricsConfig{
            Surveys: []seed.SurveyConfig{
                {SurveyID: "SV_abc123", ResponseCount: 10},
            },
        },
    })
    manager := services.NewExportJobManager(gen)

    job, _ := manager.StartExport("SV_abc123")

    // Poll until complete
    for job.Status == models.ExportStatusInProgress {
        job, _ = manager.GetProgress(job.ProgressID)
    }

    assert.Equal(t, models.ExportStatusComplete, job.Status)
    assert.Equal(t, 100, job.PercentComplete)
    assert.NotEmpty(t, job.FileID)
}

func TestExportJobManager_FileDownload(t *testing.T) {
    gen := generator.NewSurveyGenerator(&seed.SeedData{
        Qualtrics: &seed.QualtricsConfig{
            Surveys: []seed.SurveyConfig{
                {SurveyID: "SV_abc123", ResponseCount: 5},
            },
        },
    })
    manager := services.NewExportJobManager(gen)

    job, _ := manager.StartExport("SV_abc123")

    // Poll until complete
    for job.Status == models.ExportStatusInProgress {
        job, _ = manager.GetProgress(job.ProgressID)
    }

    // Download file
    data, err := manager.GetFile(job.FileID)
    require.NoError(t, err)
    assert.NotEmpty(t, data)

    // Verify it's a valid ZIP
    reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
    require.NoError(t, err)
    assert.NotEmpty(t, reader.File)
}

func TestExportJobManager_NotFound(t *testing.T) {
    gen := generator.NewSurveyGenerator(&seed.SeedData{})
    manager := services.NewExportJobManager(gen)

    _, err := manager.GetProgress("nonexistent")
    assert.Error(t, err)

    _, err = manager.GetFile("nonexistent")
    assert.Error(t, err)
}

func TestExportJobManager_ConcurrentExports(t *testing.T) {
    gen := generator.NewSurveyGenerator(&seed.SeedData{
        Qualtrics: &seed.QualtricsConfig{
            Surveys: []seed.SurveyConfig{
                {SurveyID: "SV_1", ResponseCount: 5},
                {SurveyID: "SV_2", ResponseCount: 5},
            },
        },
    })
    manager := services.NewExportJobManager(gen)

    var wg sync.WaitGroup
    wg.Add(2)

    var job1, job2 *models.ExportJob

    go func() {
        defer wg.Done()
        job1, _ = manager.StartExport("SV_1")
        for job1.Status == models.ExportStatusInProgress {
            job1, _ = manager.GetProgress(job1.ProgressID)
        }
    }()

    go func() {
        defer wg.Done()
        job2, _ = manager.StartExport("SV_2")
        for job2.Status == models.ExportStatusInProgress {
            job2, _ = manager.GetProgress(job2.ProgressID)
        }
    }()

    wg.Wait()

    assert.Equal(t, models.ExportStatusComplete, job1.Status)
    assert.Equal(t, models.ExportStatusComplete, job2.Status)
    assert.NotEqual(t, job1.FileID, job2.FileID)
}
```

**Files to Create**:
- NEW: `internal/services/qualtrics_export.go`
- NEW: `internal/services/qualtrics_export_test.go`

**Acceptance Criteria**:
- [ ] StartExport creates job with progressId
- [ ] GetProgress advances progress on each poll
- [ ] Status transitions: inProgress -> complete
- [ ] FileID assigned on completion
- [ ] GetFile returns valid ZIP
- [ ] Thread-safe concurrent exports
- [ ] Job expiry/cleanup (optional)
- [ ] Tests pass with 95%+ coverage

---

#### TASK-DS-14: Create Qualtrics API Handlers (Est: 2.0h)
**Assigned Subagent**: `cursor-sim-cli-dev`

**Goal**: Implement all three Qualtrics endpoints

**TDD Approach**:
```go
func TestQualtricsStartExport_Success(t *testing.T) {
    gen := generator.NewSurveyGenerator(&seed.SeedData{})
    manager := services.NewExportJobManager(gen)
    handlers := qualtrics.NewExportHandlers(manager)

    req := httptest.NewRequest("POST", "/API/v3/surveys/SV_abc123/export-responses", nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    handlers.StartExportHandler().ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    var resp models.ExportStartResponse
    json.Unmarshal(rr.Body.Bytes(), &resp)

    assert.NotEmpty(t, resp.Result.ProgressID)
    assert.Equal(t, "inProgress", resp.Result.Status)
    assert.Equal(t, 0, resp.Result.PercentComplete)
}

func TestQualtricsProgress_InProgress(t *testing.T) {
    gen := generator.NewSurveyGenerator(&seed.SeedData{})
    manager := services.NewExportJobManager(gen)
    handlers := qualtrics.NewExportHandlers(manager)

    // Start export
    job, _ := manager.StartExport("SV_abc123")

    req := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/"+job.ProgressID, nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    handlers.ProgressHandler().ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    var resp models.ExportProgressResponse
    json.Unmarshal(rr.Body.Bytes(), &resp)

    // Progress should have advanced
    assert.True(t, resp.Result.PercentComplete > 0)
}

func TestQualtricsProgress_Complete(t *testing.T) {
    gen := generator.NewSurveyGenerator(&seed.SeedData{
        Qualtrics: &seed.QualtricsConfig{
            Surveys: []seed.SurveyConfig{
                {SurveyID: "SV_abc123", ResponseCount: 5},
            },
        },
    })
    manager := services.NewExportJobManager(gen)
    handlers := qualtrics.NewExportHandlers(manager)

    job, _ := manager.StartExport("SV_abc123")

    // Poll until complete
    var resp models.ExportProgressResponse
    for {
        req := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/"+job.ProgressID, nil)
        req.SetBasicAuth("api-key", "")

        rr := httptest.NewRecorder()
        handlers.ProgressHandler().ServeHTTP(rr, req)

        json.Unmarshal(rr.Body.Bytes(), &resp)
        if resp.Result.Status == "complete" {
            break
        }
    }

    assert.Equal(t, 100, resp.Result.PercentComplete)
    assert.NotEmpty(t, resp.Result.FileID)
}

func TestQualtricsFileDownload_Success(t *testing.T) {
    gen := generator.NewSurveyGenerator(&seed.SeedData{
        Qualtrics: &seed.QualtricsConfig{
            Surveys: []seed.SurveyConfig{
                {SurveyID: "SV_abc123", ResponseCount: 5},
            },
        },
    })
    manager := services.NewExportJobManager(gen)
    handlers := qualtrics.NewExportHandlers(manager)

    job, _ := manager.StartExport("SV_abc123")
    for job.Status == models.ExportStatusInProgress {
        job, _ = manager.GetProgress(job.ProgressID)
    }

    req := httptest.NewRequest("GET", "/API/v3/surveys/SV_abc123/export-responses/"+job.FileID+"/file", nil)
    req.SetBasicAuth("api-key", "")

    rr := httptest.NewRecorder()
    handlers.FileDownloadHandler().ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Equal(t, "application/zip", rr.Header().Get("Content-Type"))
    assert.True(t, len(rr.Body.Bytes()) > 0)
}
```

**Files to Create**:
- NEW: `internal/api/qualtrics/handlers.go`
- NEW: `internal/api/qualtrics/handlers_test.go`

**Acceptance Criteria**:
- [ ] POST start-export returns progressId
- [ ] GET progress returns current status
- [ ] GET file returns ZIP download
- [ ] Response format matches Qualtrics API
- [ ] Authentication required
- [ ] Error handling for invalid IDs
- [ ] Tests pass with 90%+ coverage

---

### PHASE 5: Integration & E2E

#### TASK-DS-15: Router Integration and Main Initialization (Est: 2.0h)
**Assigned Subagent**: `cursor-sim-cli-dev`

**Goal**: Wire all APIs into router with conditional registration

**TDD Approach**:
```go
func TestRouter_AllAPIsEnabled(t *testing.T) {
    seedData := &seed.SeedData{
        Developers: []seed.Developer{{Email: "dev@company.com"}},
        HarveyUsers: []seed.HarveyUser{{Email: "atty@firm.com"}},
        M365Tenant: &seed.M365Tenant{
            Users: []seed.M365User{{Email: "m365@company.com", CopilotEnabled: true}},
        },
        Qualtrics: &seed.QualtricsConfig{
            Surveys: []seed.SurveyConfig{{SurveyID: "SV_abc", ResponseCount: 10}},
        },
    }

    store := storage.NewMemoryStore()
    router := server.NewRouter(store, seedData)

    // Test all endpoints exist
    endpoints := []struct {
        method string
        path   string
    }{
        {"GET", "/harvey/api/v1/history/usage"},
        {"GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')"},
        {"POST", "/API/v3/surveys/SV_abc/export-responses"},
    }

    for _, ep := range endpoints {
        req := httptest.NewRequest(ep.method, ep.path, nil)
        req.SetBasicAuth("api-key", "")

        rr := httptest.NewRecorder()
        router.ServeHTTP(rr, req)

        assert.NotEqual(t, http.StatusNotFound, rr.Code, "Endpoint %s %s should exist", ep.method, ep.path)
    }
}

func TestRouter_MinimalSeed(t *testing.T) {
    // Minimal seed with only developers (no external sources)
    seedData := &seed.SeedData{
        Developers: []seed.Developer{{Email: "dev@company.com"}},
    }

    store := storage.NewMemoryStore()
    router := server.NewRouter(store, seedData)

    // Existing endpoints should work
    req := httptest.NewRequest("GET", "/health", nil)
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)
    assert.Equal(t, http.StatusOK, rr.Code)

    // New endpoints should NOT exist
    req = httptest.NewRequest("GET", "/harvey/api/v1/history/usage", nil)
    req.SetBasicAuth("api-key", "")
    rr = httptest.NewRecorder()
    router.ServeHTTP(rr, req)
    assert.Equal(t, http.StatusNotFound, rr.Code)
}
```

**Files to Modify**:
- MODIFY: `internal/server/router.go` - Add all routes
- MODIFY: `cmd/simulator/main.go` - Initialize all generators

**Acceptance Criteria**:
- [ ] Routes registered conditionally based on seed
- [ ] All generators initialized from seed
- [ ] Data pre-generated on startup
- [ ] Existing functionality unchanged
- [ ] Tests pass

---

#### TASK-DS-16: E2E Tests and SPEC.md Update (Est: 1.5h)
**Assigned Subagent**: `cursor-sim-cli-dev`

**Goal**: Comprehensive E2E tests and documentation update

**TDD Approach**:
```go
func TestE2E_HarveyFullWorkflow(t *testing.T) {
    // Start server with enterprise seed
    seedData := loadTestSeed(t, "enterprise_seed.yaml")
    server := startTestServer(t, seedData)
    defer server.Close()

    // Query Harvey usage
    resp, err := authenticatedGet(server.URL + "/harvey/api/v1/history/usage")
    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    data := result["data"].([]interface{})
    assert.True(t, len(data) > 0, "Expected Harvey events")
}

func TestE2E_CopilotFullWorkflow(t *testing.T) {
    seedData := loadTestSeed(t, "enterprise_seed.yaml")
    server := startTestServer(t, seedData)
    defer server.Close()

    // Test JSON response
    resp, err := authenticatedGet(server.URL + "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')")
    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // Test CSV redirect
    resp, err = authenticatedGetNoFollow(server.URL + "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=text/csv")
    require.NoError(t, err)
    assert.Equal(t, http.StatusFound, resp.StatusCode)

    // Follow redirect and download CSV
    csvURL := resp.Header.Get("Location")
    resp, err = http.Get(server.URL + csvURL)
    require.NoError(t, err)
    assert.Contains(t, resp.Header.Get("Content-Type"), "text/csv")
}

func TestE2E_QualtricsFullWorkflow(t *testing.T) {
    seedData := loadTestSeed(t, "enterprise_seed.yaml")
    server := startTestServer(t, seedData)
    defer server.Close()

    // Step 1: Start export
    resp, err := authenticatedPost(server.URL + "/API/v3/surveys/SV_aitools_q1_2026/export-responses")
    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var startResult models.ExportStartResponse
    json.NewDecoder(resp.Body).Decode(&startResult)
    progressID := startResult.Result.ProgressID

    // Step 2: Poll until complete
    var progressResult models.ExportProgressResponse
    for i := 0; i < 10; i++ {
        resp, err = authenticatedGet(server.URL + "/API/v3/surveys/SV_aitools_q1_2026/export-responses/" + progressID)
        require.NoError(t, err)

        json.NewDecoder(resp.Body).Decode(&progressResult)
        if progressResult.Result.Status == "complete" {
            break
        }
        time.Sleep(100 * time.Millisecond)
    }

    assert.Equal(t, "complete", progressResult.Result.Status)
    fileID := progressResult.Result.FileID

    // Step 3: Download file
    resp, err = authenticatedGet(server.URL + "/API/v3/surveys/SV_aitools_q1_2026/export-responses/" + fileID + "/file")
    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    assert.Equal(t, "application/zip", resp.Header.Get("Content-Type"))

    // Verify ZIP contents
    body, _ := io.ReadAll(resp.Body)
    reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
    require.NoError(t, err)
    assert.Equal(t, "survey_responses.csv", reader.File[0].Name)
}

func TestE2E_CrossSourceUserCorrelation(t *testing.T) {
    seedData := loadTestSeed(t, "enterprise_seed.yaml")
    server := startTestServer(t, seedData)
    defer server.Close()

    // Get user from Harvey
    harveyResp, _ := authenticatedGet(server.URL + "/harvey/api/v1/history/usage")
    var harveyResult map[string]interface{}
    json.NewDecoder(harveyResp.Body).Decode(&harveyResult)

    harveyUsers := make(map[string]bool)
    for _, e := range harveyResult["data"].([]interface{}) {
        event := e.(map[string]interface{})
        harveyUsers[event["User"].(string)] = true
    }

    // Get user from Copilot
    copilotResp, _ := authenticatedGet(server.URL + "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')")
    var copilotResult models.CopilotUsageResponse
    json.NewDecoder(copilotResp.Body).Decode(&copilotResult)

    copilotUsers := make(map[string]bool)
    for _, u := range copilotResult.Value {
        copilotUsers[u.UserPrincipalName] = true
    }

    // Verify at least one user exists in both
    // (depends on seed configuration)
    t.Log("Harvey users:", len(harveyUsers))
    t.Log("Copilot users:", len(copilotUsers))
}
```

**Files to Create**:
- NEW: `test/e2e/harvey_test.go`
- NEW: `test/e2e/copilot_test.go`
- NEW: `test/e2e/qualtrics_test.go`
- NEW: `testdata/enterprise_seed.yaml`

**Files to Modify**:
- MODIFY: `services/cursor-sim/SPEC.md` - Document new endpoints

**Acceptance Criteria**:
- [ ] E2E tests for all three APIs
- [ ] Full workflow tests (start to finish)
- [ ] Cross-source correlation verified
- [ ] SPEC.md updated with new endpoints
- [ ] Example seed file with all sources
- [ ] All tests pass

---

## Dependency Graph

```
TASK-DS-01 (Seed Schema) ─┬─► TASK-DS-02 (Storage)
                          │
                          ├─► TASK-DS-03 (Harvey Model) ─► TASK-DS-04 (Harvey Gen) ─► TASK-DS-05 (Harvey Handler) ─► TASK-DS-06 (Harvey Router)
                          │
                          ├─► TASK-DS-07 (Copilot Model) ─► TASK-DS-08 (Copilot Gen) ─► TASK-DS-09 (Copilot Handler) ─► TASK-DS-10 (Copilot Router)
                          │
                          └─► TASK-DS-11 (Qualtrics Model) ─► TASK-DS-12 (Survey Gen) ─► TASK-DS-13 (State Machine) ─► TASK-DS-14 (Qualtrics Handler)
                                                                                                                                   │
TASK-DS-06 ─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼─► TASK-DS-15 (Router) ─► TASK-DS-16 (E2E)
TASK-DS-10 ─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
TASK-DS-14 ─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Subagent Orchestration

### Assignment Summary

| Task Range | Subagent | Reason |
|------------|----------|--------|
| TASK-DS-01 to TASK-DS-16 | `cursor-sim-cli-dev` | All work is within cursor-sim Go codebase |

### Subagent Instructions

When spawning `cursor-sim-cli-dev`, provide:

```
Context: P4-F04-data-sources - External Data Source Simulators
Task: {TASK-DS-XX}
Goal: {Goal from task definition}
Files: {Files to create/modify}
TDD: Write tests first, implement to pass tests
Constraints:
  - Follow existing cursor-sim patterns
  - 90%+ test coverage
  - JSON field names must match API specs exactly
  - Thread-safe storage operations
```

### Parallel Execution Opportunities

After TASK-DS-01 and TASK-DS-02 complete, the following can run in parallel:
- TASK-DS-03 to TASK-DS-06 (Harvey API)
- TASK-DS-07 to TASK-DS-10 (Copilot API)
- TASK-DS-11 to TASK-DS-14 (Qualtrics API)

Maximum parallelism: 3 subagents (one per API)

---

## Definition of Done (Per Task)

- [ ] Tests written BEFORE implementation (TDD)
- [ ] All tests pass
- [ ] Coverage meets target (90%+ for generators, 95%+ for state machine)
- [ ] No linting errors (`go vet`, `gofmt`)
- [ ] Code follows existing patterns
- [ ] Thread-safety verified for concurrent operations
- [ ] Git commit with descriptive message
- [ ] task.md updated with status

---

## Definition of Done (Feature Complete)

- [ ] All 16 tasks completed
- [ ] Harvey API functional with E2E tests
- [ ] Microsoft Copilot API functional with E2E tests
- [ ] Qualtrics API functional with 3-step workflow
- [ ] Cross-source user correlation verified
- [ ] SPEC.md updated with all endpoints
- [ ] Example enterprise seed documented
- [ ] All tests passing (unit + E2E)
- [ ] Master agent code review complete
- [ ] Final commit with P4-F04 documentation

---

**Next Action**: Master agent (Opus) reviews and approves plan, then spawns cursor-sim-cli-dev for TASK-DS-01
