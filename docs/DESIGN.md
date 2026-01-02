# System Design Document: Cursor Usage Analytics Platform

**Version**: 2.0.0
**Last Updated**: January 2026
**Status**: Active - Major Revision

## 1. Executive Summary

The Cursor Usage Analytics Platform is a microservices-based system designed to simulate, aggregate, and visualize AI coding assistant usage metrics. Version 2.0 introduces significant architectural changes to support **SDLC research on AI impact** through correlated datasets spanning AI telemetry, PR lifecycle, and code quality outcomes.

The system consists of three decoupled services following the ETL (Extract, Transform, Load) pattern:
- **cursor-sim**: High-fidelity API simulator with seed-based generation
- **cursor-analytics-core**: GraphQL aggregator with PostgreSQL persistence
- **cursor-viz-spa**: Interactive React dashboard

### What's New in v2.0

| Area | v1.0 | v2.0 |
|------|------|------|
| Data Generation | Internal random | Seed-based from DataDesigner |
| API Surface | Generic endpoints | Exact Cursor + GitHub APIs |
| Research Support | None | Full SDLC research framework |
| Operation Modes | Single | Runtime + Replay |
| Endpoints | ~5 | 29 Cursor + 20 GitHub |

## 2. Problem Statement

Organizations adopting AI coding assistants like Cursor lack visibility into how effectively their teams utilize these tools. More critically, **research teams** need correlated datasets to study AI's impact on:

- **Velocity**: Does AI assistance speed up coding and review cycles?
- **Review Costs**: Does AI-generated code require more review iterations?
- **Quality**: Does AI code have higher revert rates or survival issues?

This platform addresses these needs by providing:
1. A simulator that exactly matches Cursor's production API
2. GitHub-style PR lifecycle simulation for SDLC metrics
3. Pre-joined research datasets with proper JOIN keys
4. Reproducible data generation via seed files

## 3. System Architecture

### 3.1 High-Level Architecture (v2.0)

```
                           ┌─────────────────────────────────┐
                           │     DataDesigner (NVIDIA)       │
                           │  ┌───────────────────────────┐  │
                           │  │   generate_seed.py        │  │
                           │  │   ─────────────────────   │  │
                           │  │   • Developer roster      │  │
                           │  │   • Repository catalog    │  │
                           │  │   • Correlations          │  │
                           │  │   • Text templates        │  │
                           │  └───────────────────────────┘  │
                           └───────────────┬─────────────────┘
                                           │ seed.json
                                           ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                           Docker Network (cursor-net)                         │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌────────────────────────────────────┐                                      │
│  │         cursor-sim (Go)            │                                      │
│  │         Port: 8080                 │                                      │
│  │  ┌──────────────────────────────┐  │                                      │
│  │  │ Mode A: Runtime Generation   │  │                                      │
│  │  │   Load seed.json             │  │                                      │
│  │  │   Generate PR/Commits/Reviews│  │                                      │
│  │  │   Enforce correlations       │  │                                      │
│  │  └──────────────────────────────┘  │                                      │
│  │  ┌──────────────────────────────┐  │                                      │
│  │  │ Mode B: Replay Server        │  │                                      │
│  │  │   Load events.parquet        │  │                                      │
│  │  │   Serve read-only            │  │                                      │
│  │  └──────────────────────────────┘  │                                      │
│  │                                    │                                      │
│  │  APIs:                             │                                      │
│  │  ├── Cursor Admin API (4)          │                                      │
│  │  ├── Cursor AI Code Tracking (4)   │                                      │
│  │  ├── Cursor Team Analytics (11)    │                                      │
│  │  ├── Cursor By-User Analytics (9)  │                                      │
│  │  ├── GitHub Repos/PRs/Reviews (15) │                                      │
│  │  └── Research Dataset Export (5)   │                                      │
│  └─────────────────┬──────────────────┘                                      │
│                    │ REST (JSON/CSV/Parquet)                                 │
│                    ▼                                                         │
│  ┌────────────────────────────────────┐      ┌────────────────────────────┐  │
│  │    cursor-analytics-core (TS)      │      │        PostgreSQL          │  │
│  │    Port: 4000                      │◄────►│        Port: 5432          │  │
│  │  ┌──────────────────────────────┐  │      │                            │  │
│  │  │ Data Ingestion Worker        │  │      │  • developers              │  │
│  │  │   Poll cursor-sim APIs       │  │      │  • commits                 │  │
│  │  │   Normalize & store          │  │      │  • pull_requests           │  │
│  │  │   Calculate KPIs             │  │      │  • reviews                 │  │
│  │  └──────────────────────────────┘  │      │  • daily_stats (MV)        │  │
│  │  ┌──────────────────────────────┐  │      └────────────────────────────┘  │
│  │  │ GraphQL API                  │  │                                      │
│  │  │   Developer/Team queries     │  │                                      │
│  │  │   PR lifecycle queries       │  │                                      │
│  │  │   Research metrics           │  │                                      │
│  │  └──────────────────────────────┘  │                                      │
│  └─────────────────┬──────────────────┘                                      │
│                    │ GraphQL                                                 │
│                    ▼                                                         │
│  ┌────────────────────────────────────┐                                      │
│  │    cursor-viz-spa (React)          │                                      │
│  │    Port: 3000                      │                                      │
│  │  ┌──────────────────────────────┐  │                                      │
│  │  │ Dashboard Views              │  │                                      │
│  │  │   • AI Adoption metrics      │  │                                      │
│  │  │   • PR Velocity charts       │  │                                      │
│  │  │   • Review Cost analysis     │  │                                      │
│  │  │   • Quality indicators       │  │                                      │
│  │  └──────────────────────────────┘  │                                      │
│  └────────────────────────────────────┘                                      │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Data Flow: Research Framework

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           RESEARCH DATA FLOW                                    │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  DataDesigner (seed.json)                                                       │
│  ├── developers[]: user_id, email, seniority, pr_behavior, acceptance_rate     │
│  ├── repositories[]: repo_name, primary_language, maturity, quality_baseline   │
│  ├── correlations: seniority→behavior, ai_ratio→quality, cycle_times           │
│  └── text_templates: commit_messages, pr_titles, review_comments               │
│                                                                                 │
│                              ▼                                                  │
│                                                                                 │
│  cursor-sim (runtime generation)                                                │
│  ├── PR Generation Pipeline                                                     │
│  │   ├── Select developer (weighted by prs_per_week)                           │
│  │   ├── Select repository (by team assignment)                                 │
│  │   ├── Generate PR attributes (size, scatter, greenfield)                    │
│  │   ├── Generate commits (1-8 per PR)                                         │
│  │   ├── Apply AI contribution (TAB/COMPOSER split)                            │
│  │   ├── Generate timeline (coding → review → merge)                           │
│  │   └── Determine quality outcomes (revert, hotfix)                           │
│  │                                                                              │
│  ├── Cursor API (AI Telemetry)     ─┐                                          │
│  │   • /analytics/ai-code/commits   │  JOIN KEY: commit_sha                    │
│  │   • /analytics/ai-code/changes   │  JOIN KEY: user_email                    │
│  │   • /analytics/team/*            │  JOIN KEY: repo_name                     │
│  │   • /analytics/by-user/*         │                                          │
│  │                                  │                                          │
│  └── GitHub API (PR Lifecycle)     ─┘                                          │
│      • /repos/{o}/{r}/pulls                                                     │
│      • /repos/{o}/{r}/pulls/{n}/reviews                                         │
│      • /repos/{o}/{r}/analysis/survival                                         │
│      • /research/dataset (pre-joined)                                           │
│                                                                                 │
│                              ▼                                                  │
│                                                                                 │
│  Research Output (/research/dataset)                                            │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │ pr_number | author_email | ai_lines_added | pr_volume | greenfield_idx │   │
│  │ coding_lead_time | pickup_time | review_lead_time | iterations        │   │
│  │ review_density | rework_ratio | scope_creep | is_reverted | survival   │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Service Boundaries

**cursor-sim (Data Generator)**

The simulator owns:
- Loading and validating seed.json from DataDesigner
- PR lifecycle event generation with correlation enforcement
- Exact Cursor Business API contract compliance
- GitHub-style PR/review simulation for research
- Research dataset export (CSV/Parquet)

The simulator does NOT own:
- Seed data generation (DataDesigner responsibility)
- Long-term persistence (aggregator responsibility)
- Visualization (dashboard responsibility)

**cursor-analytics-core (Data Processor)**

The aggregator owns:
- Persistent storage of all ingested data
- Historical analytics and trend calculations
- GraphQL API for complex queries
- Time-series aggregations and KPIs

**cursor-viz-spa (Data Consumer)**

The dashboard owns:
- Interactive visualization of analytics
- User interaction and filtering
- Client-side caching and state

## 4. Service Specifications

### 4.1 Service A: cursor-sim (v2.0)

**Technology Stack:**
- Go 1.21+
- Standard `net/http` server
- In-memory storage (sync.Map)
- Poisson/lognormal distributions

**Operation Modes:**

| Mode | Command | Description |
|------|---------|-------------|
| Runtime | `--mode=runtime --seed=seed.json` | Generate events from seed |
| Replay | `--mode=replay --corpus=data.parquet` | Serve pre-generated corpus |

**API Surface (49 endpoints):**

```
Cursor Admin API (4):
  GET  /teams/members
  POST /teams/daily-usage-data
  POST /teams/filtered-usage-events
  POST /teams/spend

Cursor AI Code Tracking (4):
  GET  /analytics/ai-code/commits
  GET  /analytics/ai-code/commits.csv
  GET  /analytics/ai-code/changes
  GET  /analytics/ai-code/changes.csv

Cursor Team Analytics (11):
  GET  /analytics/team/agent-edits
  GET  /analytics/team/tabs
  GET  /analytics/team/dau
  GET  /analytics/team/client-versions
  GET  /analytics/team/models
  GET  /analytics/team/top-file-extensions
  GET  /analytics/team/mcp
  GET  /analytics/team/commands
  GET  /analytics/team/plans
  GET  /analytics/team/ask-mode
  GET  /analytics/team/leaderboard

Cursor By-User Analytics (9):
  GET  /analytics/by-user/agent-edits
  GET  /analytics/by-user/tabs
  GET  /analytics/by-user/models
  GET  /analytics/by-user/top-file-extensions
  GET  /analytics/by-user/client-versions
  GET  /analytics/by-user/mcp
  GET  /analytics/by-user/commands
  GET  /analytics/by-user/plans
  GET  /analytics/by-user/ask-mode

GitHub Simulation (15):
  GET  /repos
  GET  /repos/{owner}/{repo}
  GET  /repos/{owner}/{repo}/pulls
  GET  /repos/{owner}/{repo}/pulls/{n}
  GET  /repos/{owner}/{repo}/pulls/{n}/commits
  GET  /repos/{owner}/{repo}/pulls/{n}/files
  GET  /repos/{owner}/{repo}/pulls/{n}/reviews
  GET  /repos/{owner}/{repo}/commits
  GET  /repos/{owner}/{repo}/commits/{sha}
  GET  /repos/{owner}/{repo}/analysis/survival
  GET  /repos/{owner}/{repo}/analysis/reverts
  GET  /repos/{owner}/{repo}/analysis/hotfixes

Research Dataset (5):
  GET  /research/dataset
  GET  /research/metrics/velocity
  GET  /research/metrics/review-costs
  GET  /research/metrics/quality
  GET  /health
```

**Data Models (v2.0):**

```go
// Loaded from seed.json
type SeedData struct {
    Developers   []Developer    `json:"developers"`
    Repositories []Repository   `json:"repositories"`
    Correlations Correlations   `json:"correlations"`
    TextTemplates TextTemplates `json:"text_templates"`
}

type Developer struct {
    UserID         string  `json:"user_id"`      // user_xxx format
    Email          string  `json:"email"`
    Name           string  `json:"name"`
    Org            string  `json:"org"`
    Division       string  `json:"division"`
    Team           string  `json:"team"`
    Seniority      string  `json:"seniority"`
    Region         string  `json:"region"`
    AcceptanceRate float64 `json:"acceptance_rate"`
    PRBehavior     PRBehavior `json:"pr_behavior"`
}

type PullRequest struct {
    Number              int       `json:"number"`
    RepoName            string    `json:"repo_name"`
    AuthorEmail         string    `json:"author_email"`
    Title               string    `json:"title"`
    State               string    `json:"state"` // open, closed, merged
    Additions           int       `json:"additions"`
    Deletions           int       `json:"deletions"`
    ChangedFiles        int       `json:"changed_files"`
    InitialAdditions    int       `json:"initial_additions"`
    FirstCommitAt       time.Time `json:"first_commit_at"`
    CreatedAt           time.Time `json:"created_at"`
    FirstReviewAt       *time.Time `json:"first_review_at"`
    MergedAt            *time.Time `json:"merged_at"`
    CodingLeadTimeHours float64   `json:"coding_lead_time_hours"`
    PickupTimeHours     float64   `json:"pickup_time_hours"`
    ReviewLeadTimeHours float64   `json:"review_lead_time_hours"`
    ReviewComments      int       `json:"review_comments"`
    Iterations          int       `json:"iterations"`
    ReviewDensity       float64   `json:"review_density"`
    ReworkRatio         float64   `json:"rework_ratio"`
    ScopeCreep          float64   `json:"scope_creep"`
    CommitSHAs          []string  `json:"commit_shas"`
    AISummary           AISummary `json:"ai_summary"`
    IsReverted          bool      `json:"is_reverted"`
    HasHotfixFollowup   bool      `json:"has_hotfix_followup"`
}

type Commit struct {
    Hash              string    `json:"commitHash"`
    UserID            string    `json:"userId"`
    UserEmail         string    `json:"userEmail"`
    RepoName          string    `json:"repoName"`
    BranchName        string    `json:"branchName"`
    IsPrimaryBranch   bool      `json:"isPrimaryBranch"`
    TotalLinesAdded   int       `json:"totalLinesAdded"`
    TotalLinesDeleted int       `json:"totalLinesDeleted"`
    TabLinesAdded     int       `json:"tabLinesAdded"`
    TabLinesDeleted   int       `json:"tabLinesDeleted"`
    ComposerLinesAdded int      `json:"composerLinesAdded"`
    ComposerLinesDeleted int    `json:"composerLinesDeleted"`
    NonAILinesAdded   int       `json:"nonAiLinesAdded"`
    NonAILinesDeleted int       `json:"nonAiLinesDeleted"`
    Message           string    `json:"message"`
    CommitTs          time.Time `json:"commitTs"`
    CreatedAt         time.Time `json:"createdAt"`
    PRNumber          *int      `json:"pull_request_number"`
}
```

### 4.2 Relationship to DataDesigner

DataDesigner (NVIDIA NeMo) is a **seed generator**, not an event generator.

| Component | Responsibility |
|-----------|----------------|
| DataDesigner | Generates dimension data: developers, repos, correlations, text |
| cursor-sim | Generates time-series events: commits, PRs, reviews |

**Workflow:**
```bash
# Step 1: Generate seed (one-time or periodic)
python tools/data-designer/generate_seed.py -o seed.json

# Step 2: Run simulator with seed
cursor-sim --mode=runtime --seed=seed.json --port=8080

# Step 3: Export research dataset
curl "http://localhost:8080/research/dataset?format=csv" > research.csv
```

### 4.3 Correlation Enforcement

cursor-sim enforces research-valid correlations from the seed:

| Correlation | Implementation |
|-------------|----------------|
| Seniority → Acceptance Rate | Sample from seniority band in seed |
| Seniority → PR Size | Adjust avg_pr_size_loc by seniority |
| Seniority → Revert Rate | Multiply base rate by modifier |
| AI Ratio → Review Iterations | High AI = more iterations |
| AI Ratio → Review Density | High AI = more comments/LoC |
| AI Ratio → Revert Rate | High AI = higher revert risk |
| Region → Working Hours | Clip to region's work schedule |
| PR Size → Cycle Time | Larger PRs take longer |

### 4.4 JOIN Key Consistency

All APIs share JOIN keys for unified analysis:

| Key | Cursor API Field | GitHub API Field |
|-----|------------------|------------------|
| Commit | `commitHash` | `sha` |
| User | `userEmail` | `author.email` |
| Repository | `repoName` | `full_name` |

## 5. Research Framework Support

### 5.1 Research Variables by Category

**Independent Variables (Treatment):**
- AI Usage Intensity: `tabLinesAdded + composerLinesAdded`
- AI Ratio: `(tab + composer) / total`

**Control Variables:**
- PR Volume: `additions + deletions`
- PR Scatter: `changed_files`
- Greenfield Index: % of lines in new files
- Repo Maturity: `age_days`, `primary_language`, `total_size_bytes`

**Velocity Metrics (Outcome):**
- Coding Lead Time: `created_at - first_commit_at`
- Pickup Time: `first_review_at - created_at`
- Review Lead Time: `merged_at - first_review_at`

**Review Cost Metrics (Outcome):**
- Review Density: `comments / LoC`
- Iteration Count: Review → commit cycles
- Rework Ratio: `LoC_changed_in_review / initial_LoC`
- Scope Creep: `(final - initial) / final`

**Quality Metrics (Outcome):**
- Revert Rate: PRs reverted within 7 days
- Code Survival: % of lines still present after 30 days
- Hotfix Rate: Fix-PRs within 48 hours

### 5.2 Pre-Joined Research Dataset

The `/research/dataset` endpoint provides a single table optimized for statistical analysis:

```csv
pr_number,author_email,repo_name,
ai_lines_added,ai_lines_deleted,non_ai_lines_added,ai_ratio,
pr_volume,pr_scatter,greenfield_index,
coding_lead_time_hours,pickup_time_hours,review_lead_time_hours,
review_density,iteration_count,rework_ratio,scope_creep,
is_reverted,survival_rate_30d,has_hotfix_followup,
repo_age_days,primary_language,author_seniority
```

## 6. Technical Decisions

### 6.1 Why Seed-Based Generation?

**Problem**: Random generation produces statistically independent data, unsuitable for research.

**Solution**: DataDesigner generates correlated seed data with realistic relationships:
- Senior developers have higher acceptance rates
- AI-heavy PRs take longer to review
- New repos have higher greenfield ratios

**Benefits**:
- Reproducible research datasets
- Controllable correlation strengths
- Realistic organizational structure

### 6.2 Why Exact API Matching?

**Problem**: Custom API contracts require translation logic and increase maintenance burden.

**Solution**: Exactly match Cursor Business API:
- Same endpoint paths
- Same field names (camelCase for Cursor, snake_case for GitHub)
- Same pagination patterns
- Same authentication scheme

**Benefits**:
- Drop-in replacement for Cursor API
- Aggregator works with real or simulated data
- Reduced testing matrix

### 6.3 Why Two Operation Modes?

**Runtime Mode**: Generate events continuously for development and demos.
- Fresh data each run
- Configurable velocity
- Real-time event streams

**Replay Mode**: Serve pre-generated corpus for reproducible research.
- Fixed dataset for statistical analysis
- No randomness between runs
- Faster startup (no generation)

## 7. Deployment Architecture

### 7.1 Docker Compose (v2.0)

```yaml
version: '3.8'

services:
  cursor-sim:
    build: ./services/cursor-sim
    ports:
      - "8080:8080"
    volumes:
      - ./tools/data-designer/output:/data:ro
    environment:
      - CURSOR_SIM_MODE=runtime
      - CURSOR_SIM_SEED=/data/seed.json
      - CURSOR_SIM_PORT=8080
      - CURSOR_SIM_DAYS=90
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3

  postgres:
    image: postgres:15-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=cursor
      - POSTGRES_PASSWORD=cursor_dev
      - POSTGRES_DB=cursor_analytics
    volumes:
      - postgres_data:/var/lib/postgresql/data

  cursor-analytics-core:
    build: ./services/cursor-analytics-core
    ports:
      - "4000:4000"
    environment:
      - DATABASE_URL=postgresql://cursor:cursor_dev@postgres:5432/cursor_analytics
      - SIMULATOR_URL=http://cursor-sim:8080
      - POLL_INTERVAL_MS=60000
    depends_on:
      postgres:
        condition: service_healthy
      cursor-sim:
        condition: service_healthy

  cursor-viz-spa:
    build: ./services/cursor-viz-spa
    ports:
      - "3000:3000"
    environment:
      - VITE_GRAPHQL_URL=http://localhost:4000/graphql
    depends_on:
      cursor-analytics-core:
        condition: service_healthy

volumes:
  postgres_data:
```

## 8. Performance Targets

| Metric | Target |
|--------|--------|
| Startup time | < 2 seconds |
| PRs/second generation | 1,000+ |
| API response time (p99) | < 50ms |
| Memory usage (10k PRs) | < 500MB |
| CSV/Parquet export | 50MB/s |

## 9. Migration from v1.0

### 9.1 Breaking Changes

| v1.0 | v2.0 | Migration |
|------|------|-----------|
| `/v1/org/users` | `/teams/members` | Update endpoint path |
| `/v1/stats/activity` | `/analytics/ai-code/commits` | Update endpoint + schema |
| `id` field | `userId` with `user_` prefix | Update client parsing |
| Internal generation | Seed-based | Create seed.json first |

### 9.2 Preserved Functionality

- Poisson-distributed event timing
- Concurrent goroutine architecture
- In-memory storage patterns
- TDD test infrastructure

## 10. Appendix

### 10.1 Glossary (v2.0 Additions)

| Term | Definition |
|------|------------|
| DataDesigner | NVIDIA NeMo tool for generating seed data |
| Greenfield Index | % of PR lines in files < 30 days old |
| Pickup Time | Hours from PR open to first review |
| Rework Ratio | LoC changed during review / initial LoC |
| Scope Creep | (final additions - initial) / final |
| Survival Rate | % of lines still present after N days |

### 10.2 References

- Cursor Business API: https://cursor.com/docs/account/teams/analytics-api
- Cursor AI Code Tracking: https://docs.cursor.com/business/api-reference/ai-code-tracking
- NVIDIA NeMo DataDesigner: https://github.com/NVIDIA/NeMo
- GitHub REST API: https://docs.github.com/en/rest
