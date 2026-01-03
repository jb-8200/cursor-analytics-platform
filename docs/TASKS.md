# Implementation Tasks: Cursor Usage Analytics Platform

> **📚 REFERENCE DOCUMENT**
> This is a project-level task overview for orientation purposes.
> **Source of truth**: `.work-items/{feature}/task.md` for active task tracking and progress.

**Version**: 2.0.0
**Last Updated**: January 2026

This document breaks down features into actionable implementation tasks for cursor-sim v2. Each task is designed to be completable in 2-4 hours and follows TDD workflow.

## Task Status Legend

| Status | Description |
|--------|-------------|
| NOT_STARTED | Work has not begun |
| IN_PROGRESS | Active development |
| BLOCKED | Dependency not met |
| REVIEW | Awaiting code review |
| DONE | Complete and merged |

## v1.0 Tasks (ARCHIVED)

The v1.0 tasks (TASK-SIM-001 through TASK-SIM-007) are archived. The following v1 tasks were completed:

| v1.0 Task | Status | Reuse in v2 |
|-----------|--------|-------------|
| TASK-SIM-001: Go Project Structure | DONE | Partial - directory layout |
| TASK-SIM-002: CLI Flag Parsing | DONE | Replace - new flags |
| TASK-SIM-003: Developer Generator | DONE | Replace - seed loading |
| TASK-SIM-004: Event Generator | DONE | Partial - Poisson timing |
| TASK-SIM-005: In-Memory Storage | NOT_STARTED | Replace - new models |
| TASK-SIM-006: REST API Handlers | NOT_STARTED | Replace - Cursor API |
| TASK-SIM-007: Wire Up Main | NOT_STARTED | Replace - new modes |

---

## Phase 1: Complete Cursor API (MVP)

### Service A: cursor-sim v2 Tasks

---

#### TASK-R001: Initialize v2 Project Structure

**Status**: NOT_STARTED
**Feature**: SIM-R001, SIM-R002
**Estimated Hours**: 1
**Model Recommendation**: Haiku

Create the v2 project structure by copying from v1 and updating module paths.

**Implementation Steps:**
1. Create new v2 directory: `mkdir -p services/cursor-sim`
2. Copy reusable infrastructure: Makefile, Dockerfile, .golangci.yml
4. Initialize go.mod with same module path
5. Create v2 directory structure:
   ```
   services/cursor-sim/
   ├── cmd/simulator/main.go
   ├── internal/
   │   ├── config/
   │   ├── seed/
   │   ├── models/
   │   ├── generator/
   │   ├── storage/
   │   └── api/
   │       ├── cursor/
   │       ├── github/
   │       └── research/
   ├── go.mod
   ├── Makefile
   └── Dockerfile
   ```

**Definition of Done:**
- Project compiles with `go build ./...`
- Linter passes with `golangci-lint run`
- Main prints version and exits

**Dependencies:** None

---

#### TASK-R002: Implement Seed Schema Types

**Status**: NOT_STARTED
**Feature**: SIM-R001
**Estimated Hours**: 2
**Model Recommendation**: Haiku

Define Go types matching the seed.json schema from DataDesigner.

**Implementation Steps:**
1. Create `internal/seed/types.go` with all seed structures
2. Define Developer with all fields (user_id, email, pr_behavior, etc.)
3. Define Repository with language, maturity, teams
4. Define Correlations struct with coefficient maps
5. Define TextTemplates for commit messages, PR titles
6. Add JSON tags matching exact schema

**Types to Define:**
```go
type SeedData struct {
    Developers    []Developer    `json:"developers"`
    Repositories  []Repository   `json:"repositories"`
    Correlations  Correlations   `json:"correlations"`
    TextTemplates TextTemplates  `json:"text_templates"`
}

type Developer struct {
    UserID         string     `json:"user_id"`
    Email          string     `json:"email"`
    Name           string     `json:"name"`
    Org            string     `json:"org"`
    Division       string     `json:"division"`
    Team           string     `json:"team"`
    Seniority      string     `json:"seniority"`
    Region         string     `json:"region"`
    AcceptanceRate float64    `json:"acceptance_rate"`
    PRBehavior     PRBehavior `json:"pr_behavior"`
}

type PRBehavior struct {
    PRsPerWeek       float64 `json:"prs_per_week"`
    AvgPRSizeLoc     int     `json:"avg_pr_size_loc"`
    GreenfieldRatio  float64 `json:"greenfield_ratio"`
}

type Repository struct {
    RepoName        string   `json:"repo_name"`
    PrimaryLanguage string   `json:"primary_language"`
    AgeDays         int      `json:"age_days"`
    Maturity        string   `json:"maturity"`
    Teams           []string `json:"teams"`
}
```

**Definition of Done:**
- All seed types defined with JSON tags
- Unit tests verify JSON marshaling/unmarshaling
- Test with sample seed.json from DataDesigner

**Files to Create:**
- internal/seed/types.go
- internal/seed/types_test.go

---

#### TASK-R003: Implement Seed Loader with Validation

**Status**: NOT_STARTED
**Feature**: SIM-R001
**Estimated Hours**: 3
**Model Recommendation**: Sonnet

Implement seed.json loading and comprehensive validation.

**Implementation Steps:**
1. Create `internal/seed/loader.go` with `LoadSeed(path string) (*SeedData, error)`
2. Read and parse JSON file
3. Validate required fields are present
4. Validate user_id format (must match `user_xxx`)
5. Validate acceptance_rate is 0.0-1.0
6. Validate seniority is one of: junior, mid, senior
7. Validate region is supported
8. Validate repository languages are known
9. Validate correlations have required keys
10. Return structured errors with context

**Validation Rules:**
```go
func (s *SeedData) Validate() error {
    for i, dev := range s.Developers {
        if !strings.HasPrefix(dev.UserID, "user_") {
            return fmt.Errorf("developer[%d]: user_id must start with 'user_', got %q", i, dev.UserID)
        }
        if dev.AcceptanceRate < 0 || dev.AcceptanceRate > 1 {
            return fmt.Errorf("developer[%d]: acceptance_rate must be 0-1, got %f", i, dev.AcceptanceRate)
        }
        // ... more validations
    }
    return nil
}
```

**Definition of Done:**
- LoadSeed successfully loads valid seed files
- Invalid seeds produce descriptive errors
- Validation covers all required fields
- 90%+ test coverage

**Files to Create:**
- internal/seed/loader.go
- internal/seed/loader_test.go
- internal/seed/validation.go
- testdata/valid_seed.json
- testdata/invalid_seed_*.json (various failure cases)

**Test Cases:**
```go
func TestLoadSeed_Valid(t *testing.T)
func TestLoadSeed_FileNotFound(t *testing.T)
func TestLoadSeed_InvalidJSON(t *testing.T)
func TestLoadSeed_MissingDevelopers(t *testing.T)
func TestLoadSeed_InvalidUserID(t *testing.T)
func TestLoadSeed_AcceptanceRateOutOfRange(t *testing.T)
func TestLoadSeed_InvalidSeniority(t *testing.T)
```

**Dependencies:** TASK-R002

---

#### TASK-R004: Implement CLI v2 Flags

**Status**: NOT_STARTED
**Feature**: SIM-R002
**Estimated Hours**: 2
**Model Recommendation**: Haiku

Implement v2 CLI with mode, seed, and corpus flags.

**Implementation Steps:**
1. Create `internal/config/config.go` with v2 Config struct
2. Add Mode field (runtime/replay)
3. Add SeedPath field (for runtime mode)
4. Add CorpusPath field (for replay mode)
5. Add Port, Days, Velocity fields
6. Implement ParseFlags() with validation
7. Validate mode-specific requirements
8. Support environment variable overrides

**Config Struct:**
```go
type Config struct {
    Mode       string // "runtime" or "replay"
    SeedPath   string // required for runtime
    CorpusPath string // required for replay
    Port       int
    Days       int
    Velocity   string // "low", "medium", "high"
}

func ParseFlags() (*Config, error) {
    cfg := &Config{}
    flag.StringVar(&cfg.Mode, "mode", "runtime", "Operation mode: runtime or replay")
    flag.StringVar(&cfg.SeedPath, "seed", "", "Path to seed.json (required for runtime)")
    flag.StringVar(&cfg.CorpusPath, "corpus", "", "Path to events.parquet (required for replay)")
    flag.IntVar(&cfg.Port, "port", 8080, "HTTP server port")
    flag.IntVar(&cfg.Days, "days", 90, "Days of history to generate")
    flag.StringVar(&cfg.Velocity, "velocity", "medium", "Event rate: low, medium, high")
    flag.Parse()

    // Environment variable overrides
    if v := os.Getenv("CURSOR_SIM_MODE"); v != "" {
        cfg.Mode = v
    }
    // ... more overrides

    return cfg, cfg.Validate()
}
```

**Definition of Done:**
- All flags parse correctly
- Mode validation enforces seed/corpus requirements
- Environment variables override flags
- --help shows all options
- Unit tests cover all scenarios

**Files to Create:**
- internal/config/config.go
- internal/config/config_test.go

**Dependencies:** TASK-R001

---

#### TASK-R005: Implement Cursor Data Models

**Status**: NOT_STARTED
**Feature**: SIM-R003
**Estimated Hours**: 3
**Model Recommendation**: Haiku

Define Go types matching exact Cursor API response schemas.

**Implementation Steps:**
1. Create `internal/models/commit.go` with Commit struct
2. Create `internal/models/change.go` with Change struct
3. Create `internal/models/team_stats.go` for analytics types
4. Match exact JSON field names from Cursor API
5. Use camelCase for Cursor APIs
6. Add helper methods for common operations

**Types to Define:**
```go
// internal/models/commit.go
type Commit struct {
    CommitHash           string    `json:"commitHash"`
    UserID               string    `json:"userId"`
    UserEmail            string    `json:"userEmail"`
    UserName             string    `json:"userName"`
    RepoName             string    `json:"repoName"`
    BranchName           string    `json:"branchName"`
    IsPrimaryBranch      bool      `json:"isPrimaryBranch"`
    TotalLinesAdded      int       `json:"totalLinesAdded"`
    TotalLinesDeleted    int       `json:"totalLinesDeleted"`
    TabLinesAdded        int       `json:"tabLinesAdded"`
    TabLinesDeleted      int       `json:"tabLinesDeleted"`
    ComposerLinesAdded   int       `json:"composerLinesAdded"`
    ComposerLinesDeleted int       `json:"composerLinesDeleted"`
    NonAILinesAdded      int       `json:"nonAiLinesAdded"`
    NonAILinesDeleted    int       `json:"nonAiLinesDeleted"`
    Message              string    `json:"message"`
    CommitTs             time.Time `json:"commitTs"`
    CreatedAt            time.Time `json:"createdAt"`
}

// internal/models/team_stats.go
type AgentEditsDay struct {
    EventDate               string `json:"event_date"`
    TotalSuggestedDiffs     int    `json:"total_suggested_diffs"`
    TotalAcceptedDiffs      int    `json:"total_accepted_diffs"`
    TotalRejectedDiffs      int    `json:"total_rejected_diffs"`
    TotalGreenLinesAccepted int    `json:"total_green_lines_accepted"`
    TotalRedLinesAccepted   int    `json:"total_red_lines_accepted"`
    // ... all fields from Cursor docs
}
```

**Definition of Done:**
- All Cursor API response types defined
- JSON tags match Cursor documentation exactly
- Helper methods for calculations (e.g., AIRatio())
- Unit tests verify JSON compatibility

**Files to Create:**
- internal/models/commit.go
- internal/models/change.go
- internal/models/team_stats.go
- internal/models/user_stats.go
- internal/models/response.go (pagination, params)

**Dependencies:** None

---

#### TASK-R006: Implement Commit Generation Engine

**Status**: NOT_STARTED
**Feature**: SIM-R003
**Estimated Hours**: 5
**Model Recommendation**: Sonnet

Generate commits with AI attribution matching seed behavior.

**Implementation Steps:**
1. Create `internal/generator/commit_generator.go`
2. Copy Poisson timing from v1 (tested, working)
3. Implement developer selection weighted by prs_per_week
4. Generate commit sizes using lognormal distribution
5. Split AI lines: TAB (60-80%), COMPOSER (20-40%)
6. Calculate nonAiLines = max(0, total - tab - composer)
7. Generate commit messages from templates
8. Assign to repositories by team
9. Generate for configured time range (--days)

**Key Methods:**
```go
type CommitGenerator struct {
    seed       *seed.SeedData
    store      storage.Store
    rng        *rand.Rand
    velocity   VelocityConfig
}

func (g *CommitGenerator) GenerateCommits(ctx context.Context, days int) error {
    startTime := time.Now().AddDate(0, 0, -days)

    for _, dev := range g.seed.Developers {
        if err := g.generateForDeveloper(ctx, dev, startTime); err != nil {
            return err
        }
    }
    return nil
}

func (g *CommitGenerator) generateForDeveloper(ctx context.Context, dev seed.Developer, start time.Time) error {
    // Use Poisson distribution for commit timing
    rate := g.velocity.CommitsPerHour(dev.PRBehavior.PRsPerWeek)

    current := start
    for current.Before(time.Now()) {
        // Wait time follows exponential distribution
        waitHours := g.exponential(1.0 / rate)
        current = current.Add(time.Duration(waitHours * float64(time.Hour)))

        if current.After(time.Now()) {
            break
        }

        // Generate commit
        commit := g.generateCommit(dev, current)
        if err := g.store.AddCommit(commit); err != nil {
            return err
        }
    }
    return nil
}
```

**Definition of Done:**
- Commits generated for configured date range
- AI attribution sums correctly (tab + composer + nonAi = total)
- Developer distribution matches weights
- Reproducible with same seed
- 85%+ test coverage

**Files to Create:**
- internal/generator/commit_generator.go
- internal/generator/commit_generator_test.go
- internal/generator/poisson.go (copy from v1)
- internal/generator/velocity.go

**Test Cases:**
```go
func TestCommitGenerator_AIAttribution(t *testing.T)
func TestCommitGenerator_DeveloperDistribution(t *testing.T)
func TestCommitGenerator_TimeRange(t *testing.T)
func TestCommitGenerator_Reproducibility(t *testing.T)
func TestCommitGenerator_RegionTimezones(t *testing.T)
```

**Dependencies:** TASK-R003, TASK-R005

---

#### TASK-R007: Implement In-Memory Storage v2

**Status**: NOT_STARTED
**Feature**: SIM-R008
**Estimated Hours**: 4
**Model Recommendation**: Sonnet

Thread-safe storage for v2 models with efficient queries.

**Implementation Steps:**
1. Create `internal/storage/store.go` with interface
2. Create `internal/storage/memory.go` with implementation
3. Index commits by time for range queries
4. Index by user for by-user endpoints
5. Index by repo for GitHub API
6. Use sync.RWMutex for concurrency
7. Implement efficient time range filtering

**Interface:**
```go
type Store interface {
    // Developers (loaded from seed)
    LoadDevelopers(developers []seed.Developer) error
    GetDeveloper(userID string) (*seed.Developer, error)
    GetDeveloperByEmail(email string) (*seed.Developer, error)
    ListDevelopers() []seed.Developer

    // Commits
    AddCommit(commit models.Commit) error
    GetCommitByHash(hash string) (*models.Commit, error)
    GetCommitsByTimeRange(from, to time.Time) []models.Commit
    GetCommitsByUser(userID string, from, to time.Time) []models.Commit
    GetCommitsByRepo(repoName string, from, to time.Time) []models.Commit

    // Aggregations
    GetDailyCommitStats(from, to time.Time) []models.DailyStats
    GetUserDailyStats(userID string, from, to time.Time) []models.DailyStats
}
```

**Memory Structure:**
```go
type MemoryStore struct {
    mu sync.RWMutex

    developers      map[string]*seed.Developer // by user_id
    developerEmails map[string]string          // email -> user_id

    commits         []*models.Commit           // time-sorted
    commitsByHash   map[string]*models.Commit
    commitsByUser   map[string][]*models.Commit
    commitsByRepo   map[string][]*models.Commit
}
```

**Definition of Done:**
- Thread-safe read/write operations
- Time range queries O(log n) or better
- Memory < 500MB for 100k commits
- Query latency < 10ms
- 90%+ test coverage

**Files to Create:**
- internal/storage/store.go
- internal/storage/memory.go
- internal/storage/memory_test.go

**Test Cases:**
```go
func TestMemoryStore_LoadDevelopers(t *testing.T)
func TestMemoryStore_AddCommit(t *testing.T)
func TestMemoryStore_TimeRangeQuery(t *testing.T)
func TestMemoryStore_UserQuery(t *testing.T)
func TestMemoryStore_ConcurrentAccess(t *testing.T)
func TestMemoryStore_MemoryEfficiency(t *testing.T)
```

**Dependencies:** TASK-R002, TASK-R005

---

#### TASK-R008: Implement Common API Infrastructure

**Status**: NOT_STARTED
**Feature**: SIM-R004
**Estimated Hours**: 2
**Model Recommendation**: Haiku

Create shared API infrastructure: auth, pagination, error handling.

**Implementation Steps:**
1. Create `internal/api/middleware.go` for common middleware
2. Implement Basic Auth middleware
3. Implement rate limiting middleware
4. Create `internal/api/response.go` for response helpers
5. Implement Cursor pagination format
6. Implement error response format

**Response Helpers:**
```go
type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Pagination Pagination  `json:"pagination"`
    Params     Params      `json:"params"`
}

type Pagination struct {
    Page            int  `json:"page"`
    PageSize        int  `json:"pageSize"`
    TotalUsers      int  `json:"totalUsers,omitempty"`
    TotalPages      int  `json:"totalPages"`
    HasNextPage     bool `json:"hasNextPage"`
    HasPreviousPage bool `json:"hasPreviousPage"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func WriteCSV(w http.ResponseWriter, filename string, data [][]string) {
    w.Header().Set("Content-Type", "text/csv")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
    csv.NewWriter(w).WriteAll(data)
}
```

**Middleware:**
```go
func BasicAuth(apiKey string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, _, ok := r.BasicAuth()
            if !ok || user != apiKey {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler
```

**Definition of Done:**
- Auth middleware validates Basic Auth
- Rate limiting returns 429 when exceeded
- Pagination helper generates correct format
- CSV writer produces RFC 4180 output

**Files to Create:**
- internal/api/middleware.go
- internal/api/middleware_test.go
- internal/api/response.go
- internal/api/response_test.go

**Dependencies:** None

---

#### TASK-R009: Implement /teams/members Endpoint

**Status**: NOT_STARTED
**Feature**: SIM-R004
**Estimated Hours**: 1.5
**Model Recommendation**: Haiku

First Cursor Admin API endpoint.

**Implementation Steps:**
1. Create `internal/api/cursor/admin.go`
2. Implement GET /teams/members handler
3. Return developers from storage
4. Apply pagination
5. Match Cursor response format exactly

**Response Format:**
```json
{
  "data": [
    {
      "id": "user_001",
      "email": "dev@example.com",
      "name": "Jane Developer",
      "role": "member"
    }
  ],
  "pagination": {...}
}
```

**Definition of Done:**
- Endpoint returns all developers
- Pagination works correctly
- Response matches Cursor schema
- Auth required

**Files to Create:**
- internal/api/cursor/admin.go
- internal/api/cursor/admin_test.go

**Dependencies:** TASK-R007, TASK-R008

---

#### TASK-R010: Implement /analytics/ai-code/commits Endpoint

**Status**: NOT_STARTED
**Feature**: SIM-R005
**Estimated Hours**: 2
**Model Recommendation**: Sonnet

Primary AI Code Tracking endpoint.

**Implementation Steps:**
1. Create `internal/api/cursor/aicode.go`
2. Implement GET /analytics/ai-code/commits
3. Parse startDate/endDate query params
4. Parse users filter (comma-separated)
5. Query storage for commits in range
6. Format response with exact Cursor schema

**Query Parameters:**
- startDate (default: 7 days ago)
- endDate (default: today)
- users (optional: comma-separated emails or IDs)
- page, pageSize

**Response Format:**
```json
{
  "data": [
    {
      "commitHash": "abc123",
      "userId": "user_001",
      "userEmail": "dev@example.com",
      "repoName": "acme/platform",
      "totalLinesAdded": 150,
      "tabLinesAdded": 80,
      "composerLinesAdded": 40,
      "nonAiLinesAdded": 30,
      "commitTs": "2026-01-15T10:30:00Z"
    }
  ],
  "pagination": {...},
  "params": {
    "startDate": "2026-01-08",
    "endDate": "2026-01-15"
  }
}
```

**Definition of Done:**
- Date filtering works correctly
- User filtering supports email and user_id
- Pagination works correctly
- Response matches Cursor exactly

**Files to Create:**
- internal/api/cursor/aicode.go
- internal/api/cursor/aicode_test.go

**Dependencies:** TASK-R007, TASK-R008

---

#### TASK-R011: Implement /analytics/ai-code/commits.csv Endpoint

**Status**: NOT_STARTED
**Feature**: SIM-R005
**Estimated Hours**: 1
**Model Recommendation**: Haiku

CSV export for commits.

**Implementation Steps:**
1. Add handler for GET /analytics/ai-code/commits.csv
2. Reuse query logic from JSON endpoint
3. Format as RFC 4180 CSV
4. Set Content-Disposition header

**Definition of Done:**
- CSV is valid RFC 4180
- Headers match JSON field names
- Download filename is commits.csv

**Dependencies:** TASK-R010

---

#### TASK-R012: Implement /analytics/team/* Endpoints (11 total)

**Status**: NOT_STARTED
**Feature**: SIM-R006
**Estimated Hours**: 6
**Model Recommendation**: Sonnet

All 11 team analytics endpoints.

**Implementation Steps:**
1. Create `internal/api/cursor/team.go`
2. Implement shared handler pattern for common logic
3. Implement each endpoint with specific aggregation:
   - agent-edits: aggregate diffs suggested/accepted/rejected
   - tabs: aggregate tab completions
   - dau: count distinct active users per day
   - client-versions: group by version
   - models: group by model name
   - top-file-extensions: group by extension
   - mcp: group by MCP server
   - commands: group by command name
   - plans: group by model
   - ask-mode: group by model
   - leaderboard: rank by acceptance ratio

**Common Handler Pattern:**
```go
type teamHandler struct {
    store   storage.Store
    metric  string
    extract func(commits []models.Commit) interface{}
}

func (h *teamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    from, to := parseDateRange(r)
    commits := h.store.GetCommitsByTimeRange(from, to)

    data := h.extract(commits)

    WriteJSON(w, http.StatusOK, PaginatedResponse{
        Data:   data,
        Params: Params{Metric: h.metric, TeamID: 12345},
    })
}
```

**Definition of Done:**
- All 11 endpoints implemented
- Aggregations are accurate
- Date filtering works
- Rate limiting enforced (100/min)

**Files to Create:**
- internal/api/cursor/team.go
- internal/api/cursor/team_test.go

**Dependencies:** TASK-R007, TASK-R008

---

#### TASK-R013: Implement /analytics/by-user/* Endpoints (9 total)

**Status**: NOT_STARTED
**Feature**: SIM-R007
**Estimated Hours**: 4
**Model Recommendation**: Sonnet

All 9 by-user analytics endpoints.

**Implementation Steps:**
1. Create `internal/api/cursor/byuser.go`
2. Implement shared pattern for user grouping
3. Group data by user_email
4. Include userMappings in params
5. Apply pagination across users

**Response Format:**
```json
{
  "data": {
    "dev@example.com": [
      { "event_date": "2026-01-15", "total_accepts": 45 }
    ]
  },
  "pagination": {...},
  "params": {
    "metric": "tabs",
    "userMappings": [
      { "id": "user_001", "email": "dev@example.com" }
    ]
  }
}
```

**Definition of Done:**
- All 9 endpoints implemented
- Data grouped correctly by user
- Pagination across users works
- Rate limiting enforced (50/min)

**Files to Create:**
- internal/api/cursor/byuser.go
- internal/api/cursor/byuser_test.go

**Dependencies:** TASK-R007, TASK-R008

---

#### TASK-R014: Implement HTTP Server and Router

**Status**: NOT_STARTED
**Feature**: SIM-R004
**Estimated Hours**: 2
**Model Recommendation**: Haiku

Wire up all endpoints with routing and middleware.

**Implementation Steps:**
1. Create `internal/api/router.go`
2. Register all Cursor endpoints
3. Apply auth middleware to analytics routes
4. Apply rate limiting middleware
5. Add health check endpoint
6. Add graceful shutdown handling

**Router Setup:**
```go
func NewRouter(store storage.Store, cfg *config.Config) http.Handler {
    r := http.NewServeMux()

    // Health (no auth)
    r.HandleFunc("GET /health", healthHandler)

    // Cursor Admin API
    admin := cursor.NewAdminHandler(store)
    r.Handle("GET /teams/members", admin.Members())

    // Cursor AI Code Tracking
    aicode := cursor.NewAICodeHandler(store)
    r.Handle("GET /analytics/ai-code/commits", aicode.Commits())
    r.Handle("GET /analytics/ai-code/commits.csv", aicode.CommitsCSV())

    // ... all other endpoints

    // Apply middleware
    handler := RateLimit(100, time.Minute)(r)
    handler = BasicAuth(cfg.APIKey)(handler)
    handler = Logger(handler)

    return handler
}
```

**Definition of Done:**
- All endpoints routable
- Auth required on analytics routes
- Health check public
- Graceful shutdown works

**Files to Create:**
- internal/api/router.go
- internal/api/router_test.go

**Dependencies:** TASK-R009 through TASK-R013

---

#### TASK-R015: Wire Up Main Application

**Status**: NOT_STARTED
**Feature**: SIM-R002
**Estimated Hours**: 2
**Model Recommendation**: Haiku

Connect all components in main for runtime mode.

**Implementation Steps:**
1. Update `cmd/simulator/main.go`
2. Parse config flags
3. Load seed.json
4. Initialize storage
5. Load developers into storage
6. Generate commits
7. Start HTTP server
8. Handle graceful shutdown

**Main Flow:**
```go
func main() {
    cfg, err := config.ParseFlags()
    if err != nil {
        log.Fatal(err)
    }

    if cfg.Mode != "runtime" {
        log.Fatal("Replay mode not implemented yet")
    }

    seed, err := seed.LoadSeed(cfg.SeedPath)
    if err != nil {
        log.Fatal(err)
    }

    store := storage.NewMemoryStore()
    store.LoadDevelopers(seed.Developers)

    gen := generator.NewCommitGenerator(seed, store)
    if err := gen.GenerateCommits(context.Background(), cfg.Days); err != nil {
        log.Fatal(err)
    }

    router := api.NewRouter(store, cfg)
    server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: router}

    // Graceful shutdown
    go func() {
        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
        <-sigint
        server.Shutdown(context.Background())
    }()

    log.Printf("cursor-sim v2 running on :%d", cfg.Port)
    server.ListenAndServe()
}
```

**Definition of Done:**
- Server starts with seed file
- All endpoints respond correctly
- Graceful shutdown works
- Integration tests pass

**Files to Modify:**
- cmd/simulator/main.go
- Makefile (add run target)

**Dependencies:** TASK-R014

---

#### TASK-R016: End-to-End Testing

**Status**: NOT_STARTED
**Feature**: All Phase 1
**Estimated Hours**: 4
**Model Recommendation**: Sonnet

Comprehensive integration tests for all 29 endpoints.

**Implementation Steps:**
1. Create `test/integration/` directory
2. Start server with test seed
3. Test all Cursor endpoints
4. Verify response schemas match docs
5. Test pagination edge cases
6. Test date filtering
7. Test error cases (auth, rate limit)

**Test Structure:**
```go
func TestIntegration(t *testing.T) {
    // Start server with test seed
    seed := loadTestSeed(t)
    store := storage.NewMemoryStore()
    store.LoadDevelopers(seed.Developers)

    gen := generator.NewCommitGenerator(seed, store)
    gen.GenerateCommits(context.Background(), 7)

    router := api.NewRouter(store, testConfig)
    server := httptest.NewServer(router)
    defer server.Close()

    t.Run("GET /teams/members", func(t *testing.T) {
        // ...
    })

    t.Run("GET /analytics/ai-code/commits", func(t *testing.T) {
        // ...
    })
    // ... all endpoints
}
```

**Definition of Done:**
- All 29 endpoints tested
- Response schemas validated
- Pagination tested
- Error cases tested
- 80%+ coverage

**Files to Create:**
- test/integration/cursor_api_test.go
- test/testdata/integration_seed.json

**Dependencies:** TASK-R015

---

## Phase 1 Summary

| Task | Feature | Hours | Status |
|------|---------|-------|--------|
| TASK-R001 | Project Structure | 1 | NOT_STARTED |
| TASK-R002 | Seed Types | 2 | NOT_STARTED |
| TASK-R003 | Seed Loader | 3 | NOT_STARTED |
| TASK-R004 | CLI v2 | 2 | NOT_STARTED |
| TASK-R005 | Cursor Models | 3 | NOT_STARTED |
| TASK-R006 | Commit Generator | 5 | NOT_STARTED |
| TASK-R007 | Storage v2 | 4 | NOT_STARTED |
| TASK-R008 | API Infrastructure | 2 | NOT_STARTED |
| TASK-R009 | /teams/members | 1.5 | NOT_STARTED |
| TASK-R010 | /ai-code/commits | 2 | NOT_STARTED |
| TASK-R011 | /ai-code/commits.csv | 1 | NOT_STARTED |
| TASK-R012 | /team/* (11) | 6 | NOT_STARTED |
| TASK-R013 | /by-user/* (9) | 4 | NOT_STARTED |
| TASK-R014 | Router | 2 | NOT_STARTED |
| TASK-R015 | Main | 2 | NOT_STARTED |
| TASK-R016 | E2E Tests | 4 | NOT_STARTED |
| **Total** | | **44.5** | |

---

## Phase 2: PR Lifecycle + GitHub API

*Tasks to be detailed after Phase 1 completion*

| Task | Feature | Est. Hours |
|------|---------|------------|
| TASK-R017 | PR Generation Pipeline | 6 |
| TASK-R018 | Review Simulation | 6 |
| TASK-R019 | GitHub /repos endpoints | 4 |
| TASK-R020 | GitHub /pulls endpoints | 6 |
| TASK-R021 | Quality Outcomes | 4 |
| TASK-R022 | E2E Tests | 4 |
| **Total** | | **30** |

---

## Phase 3: Research Framework

*Tasks to be detailed after Phase 2 completion*

| Task | Feature | Est. Hours |
|------|---------|------------|
| TASK-R023 | Research Dataset Export | 4 |
| TASK-R024 | Code Survival Tracking | 6 |
| TASK-R025 | Replay Mode | 4 |
| TASK-R026 | E2E Tests | 4 |
| **Total** | | **18** |

---

## Task Dependency Graph

```
TASK-R001 (Project Init)
    │
    ├── TASK-R002 (Seed Types)
    │       │
    │       └── TASK-R003 (Seed Loader)
    │               │
    │               └── TASK-R006 (Commit Generator)
    │                       │
    │                       └── TASK-R007 (Storage)
    │                               │
    │                               └── TASK-R009-R013 (Endpoints)
    │                                       │
    │                                       └── TASK-R014 (Router)
    │                                               │
    │                                               └── TASK-R015 (Main)
    │                                                       │
    │                                                       └── TASK-R016 (E2E)
    │
    ├── TASK-R004 (CLI) ────────────────────────────────────┘
    │
    ├── TASK-R005 (Cursor Models) ─────────────────────────┘
    │
    └── TASK-R008 (API Infrastructure) ────────────────────┘
```
