# Technical Design: cursor-sim Foundation

**Feature ID**: P1-F01-foundation
**Phase**: P1 (cursor-sim Foundation)
**Created**: January 2, 2026
**Status**: COMPLETE

## Decision Log

| Date | Decision | Rationale | Evidence |
|------|----------|-----------|----------|
| Jan 2026 | Seed-based generation | Random generation produces uncorrelated data unsuitable for research | Prior v1 analysis showed statistical independence |
| Jan 2026 | Exact Cursor API matching | Drop-in replacement reduces aggregator complexity | Cursor docs: https://cursor.com/docs/account/teams/analytics-api |
| Jan 2026 | In-memory storage | MVP simplicity, acceptable data loss on restart | Performance target: <50ms p99 |
| Jan 2026 | Go 1.21+ | Existing v1 codebase, team expertise, performance | Poisson implementation proven at 93.3% coverage |

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         cursor-sim v2                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────────┐  │
│  │ CLI/Config   │───▶│ Seed Loader  │───▶│ Commit Generator     │  │
│  │ (flag pkg)   │    │ (JSON parse) │    │ (Poisson timing)     │  │
│  └──────────────┘    └──────────────┘    └──────────────────────┘  │
│                              │                      │               │
│                              ▼                      ▼               │
│                       ┌──────────────────────────────────┐         │
│                       │      In-Memory Storage           │         │
│                       │  ┌────────────┐ ┌─────────────┐  │         │
│                       │  │ Developers │ │ Commits     │  │         │
│                       │  │ (by ID)    │ │ (by time)   │  │         │
│                       │  └────────────┘ └─────────────┘  │         │
│                       └──────────────────────────────────┘         │
│                                      │                              │
│  ┌───────────────────────────────────┴────────────────────────┐   │
│  │                    HTTP Router                              │   │
│  │  ┌─────────────────────────────────────────────────────┐   │   │
│  │  │ Middleware: BasicAuth → RateLimit → Logger          │   │   │
│  │  └─────────────────────────────────────────────────────┘   │   │
│  │                                                             │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │   │
│  │  │ Admin API   │ │ AI Code API │ │ Analytics   │          │   │
│  │  │ /teams/*    │ │ /analytics/ │ │ /analytics/ │          │   │
│  │  │ (4 eps)     │ │ ai-code/*   │ │ team/* (11) │          │   │
│  │  └─────────────┘ │ (4 eps)     │ │ by-user (9) │          │   │
│  │                  └─────────────┘ └─────────────┘          │   │
│  └────────────────────────────────────────────────────────────┘   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

## Package Structure

```
services/cursor-sim/
├── cmd/simulator/
│   └── main.go              # Entry point, wiring
├── internal/
│   ├── config/
│   │   └── config.go        # CLI flags, env vars
│   ├── seed/
│   │   ├── types.go         # SeedData, Developer, Repository
│   │   ├── loader.go        # LoadSeed(), Validate()
│   │   └── validation.go    # Field validators
│   ├── models/
│   │   ├── commit.go        # Cursor Commit schema
│   │   ├── change.go        # Cursor Change schema
│   │   ├── team_stats.go    # Team analytics types
│   │   ├── user_stats.go    # By-user analytics types
│   │   └── response.go      # Pagination, Params
│   ├── generator/
│   │   ├── commit_generator.go  # Main generation logic
│   │   ├── poisson.go           # Timing distribution
│   │   └── velocity.go          # Rate configuration
│   ├── storage/
│   │   ├── store.go         # Interface definition
│   │   └── memory.go        # In-memory implementation
│   └── api/
│       ├── middleware.go    # Auth, rate limiting, logging
│       ├── response.go      # JSON/CSV helpers
│       ├── router.go        # Route registration
│       └── cursor/
│           ├── admin.go     # /teams/* handlers
│           ├── aicode.go    # /analytics/ai-code/* handlers
│           ├── team.go      # /analytics/team/* handlers
│           └── byuser.go    # /analytics/by-user/* handlers
├── testdata/
│   ├── valid_seed.json
│   └── invalid_seed_*.json
├── go.mod
├── Makefile
└── Dockerfile
```

## Data Models

### Seed Types (from DataDesigner)

```go
type SeedData struct {
    Developers    []Developer    `json:"developers"`
    Repositories  []Repository   `json:"repositories"`
    Correlations  Correlations   `json:"correlations"`
    TextTemplates TextTemplates  `json:"text_templates"`
}

type Developer struct {
    UserID         string     `json:"user_id"`      // user_xxx format
    Email          string     `json:"email"`
    Name           string     `json:"name"`
    Org            string     `json:"org"`
    Division       string     `json:"division"`
    Team           string     `json:"team"`
    Seniority      string     `json:"seniority"`    // junior|mid|senior
    Region         string     `json:"region"`       // us-west|us-east|eu|apac
    AcceptanceRate float64    `json:"acceptance_rate"`
    PRBehavior     PRBehavior `json:"pr_behavior"`
}

type PRBehavior struct {
    PRsPerWeek      float64 `json:"prs_per_week"`
    AvgPRSizeLoc    int     `json:"avg_pr_size_loc"`
    GreenfieldRatio float64 `json:"greenfield_ratio"`
}
```

### Cursor API Types (exact match)

```go
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
```

## API Contract

### Endpoints (29 total)

| Group | Endpoints | Auth | Rate Limit |
|-------|-----------|------|------------|
| Admin | /teams/members, /teams/daily-usage-data, /teams/filtered-usage-events, /teams/spend | Basic | 100/min |
| AI Code | /analytics/ai-code/commits, /analytics/ai-code/commits.csv, /analytics/ai-code/changes, /analytics/ai-code/changes.csv | Basic | 100/min |
| Team Analytics | /analytics/team/{agent-edits,tabs,dau,client-versions,models,top-file-extensions,mcp,commands,plans,ask-mode,leaderboard} | Basic | 100/min |
| By-User | /analytics/by-user/{agent-edits,tabs,models,top-file-extensions,client-versions,mcp,commands,plans,ask-mode} | Basic | 50/min |
| Health | /health | None | None |

### Authentication

Basic Auth matching Cursor API:
- Username: API key
- Password: empty string

```go
func BasicAuth(apiKey string) func(http.Handler) http.Handler
```

### Rate Limiting

Token bucket algorithm:
- Team endpoints: 100 requests/minute
- By-user endpoints: 50 requests/minute
- Return 429 when exceeded with Retry-After header

## Generation Algorithm

### Commit Generation

1. For each developer in seed:
   - Calculate commit rate from `prs_per_week`
   - Use Poisson timing (exponential inter-arrival)
2. For each commit:
   - Generate size from lognormal distribution
   - Apply developer's `acceptance_rate` for AI ratio
   - Split AI lines: TAB (60-80%), COMPOSER (20-40%)
   - Calculate: `nonAi = max(0, total - tab - composer)`
   - Select repository by team assignment
   - Generate commit message from templates
   - Assign deterministic hash from seed + index

### Reproducibility

```go
// Seed RNG with deterministic value
rng := rand.New(rand.NewSource(hashSeed(seedData)))

// Each developer gets unique but reproducible seed
devRNG := rand.New(rand.NewSource(baseSeed + int64(devIndex)*1000))
```

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| Startup | < 2s | time.Now() around load |
| API p99 | < 50ms | httptest + percentile |
| Memory | < 500MB for 100k commits | runtime.MemStats |
| Generation | 1000+ commits/sec | benchmark test |

## Testing Strategy

- Unit tests: Each package independently
- Integration tests: Full API endpoint testing
- Property tests: AI line sum invariant
- Benchmark tests: Performance targets
- Target coverage: 80% minimum, 90%+ for generators

## Risks and Mitigations

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Schema drift from Cursor | Medium | High | Pin to specific API version, diff checks |
| Memory pressure at scale | Low | Medium | Streaming responses, chunked generation |
| Seed complexity | Medium | Low | Comprehensive validation, helpful errors |

## Alternatives Considered

1. **Database storage**: Rejected for MVP complexity; in-memory sufficient
2. **gRPC API**: Rejected; must match Cursor REST API exactly
3. **External seed service**: Rejected; file-based is simpler and reproducible
