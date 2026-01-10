# System Design Document: Cursor Usage Analytics Platform

> **📚 REFERENCE DOCUMENT**
> This is a project-level overview for orientation purposes.
> **Source of truth**: `services/{service}/SPEC.md` for technical specs, `.work-items/` for active work.

**Version**: 2.2.0
**Last Updated**: January 4, 2026
**Status**: Active - Integration Testing Complete, Phase 3C Design Complete

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

### 3.1 High-Level Architecture (v2.2 - Current Deployment)

**As Implemented (January 4, 2026)**:

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
┌─────────────────────────┐      ┌────────────────────────────────┐      ┌─────────────────────┐
│   cursor-sim (P4)       │      │  cursor-analytics-core (P5)    │      │  cursor-viz-spa     │
│   Docker Container      │      │  ┌──────────────────────────┐  │      │  (P6 - React/Vite)  │
│   Port 8080             │─────▶│  │  GraphQL (Port 4000)     │  │─────▶│  Local npm dev      │
│                         │ REST │  │  Apollo Server 4         │  │ GQL  │  Port 3000          │
│  Runtime Generation:    │      │  └──────────────────────────┘  │      │                     │
│  • Load seed.json       │      │  ┌──────────────────────────┐  │      │  Apollo Client      │
│  • Generate PRs/Commits │      │  │  PostgreSQL (Port 5432)  │  │      │  Recharts viz       │
│  • Enforce correlations │      │  │  (Internal to Docker)    │  │      │  Tailwind CSS       │
│  • Cursor API (29)      │      │  └──────────────────────────┘  │      │                     │
│  • GitHub API (15)      │      │  Docker Compose Stack         │      │  Connects to:       │
│  • Research Export (5)  │      │  (GraphQL + PostgreSQL)       │      │  http://localhost:  │
└─────────────────────────┘      └────────────────────────────────┘      │  4000/graphql       │
   Docker (standalone)                                                   └─────────────────────┘
                                                                             Host (npm run dev)
```

**Deployment Models**:
- **cursor-sim (P4)**: Docker container (`cursor-sim-local`) - Runtime generation mode
- **cursor-analytics-core (P5)**: Docker Compose stack (GraphQL + PostgreSQL) - No ingestion worker (deferred)
- **cursor-viz-spa (P6)**: Local npm dev server - Fast hot-reload development

**Why This Architecture**:
- P5 GraphQL and PostgreSQL run in same Docker network for reliable container-to-container communication
- P6 runs locally for fast hot-reload development experience (Vite HMR)
- P4 runs in Docker for consistent seed data and API simulation
- **Integration tested and verified working** (January 4, 2026)

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

**Replay Mode**: Serve pre-generated corpus for reproducible research (deferred to Phase 3D).
- Fixed dataset for statistical analysis
- No randomness between runs
- Faster startup (no generation)

For Phase 3C, reproducibility is achieved via seeded RNG + deterministic event generation.

### 6.4 PR Generation Strategy (Phase 3C)

> **Design Decision** (January 3, 2026): PRs are derived on-the-fly from commit groupings using session-based parameters.

**Session-Based Generation Model**:

PRs emerge naturally from "work sessions" with developer-specific characteristics enforced through session parameters:

```go
type Session struct {
    Developer     seed.Developer
    Repo          seed.Repository
    Branch        string
    StartTime     time.Time
    MaxCommits    int           // Seniority-based
    TargetLoC     int           // Affects commit sizes
    InactivityGap time.Duration // From working hours
    Commits       []models.Commit
}

func StartSession(dev seed.Developer, repo seed.Repository) *Session {
    return &Session{
        Developer:     dev,
        Repo:          repo,
        MaxCommits:    sampleMaxCommits(dev.Seniority),    // seniors: 5-12, juniors: 2-5
        TargetLoC:     sampleTargetLoC(dev.Seniority),     // affects commit sizes
        InactivityGap: sampleGap(dev.WorkingHoursBand),    // 15-60 minutes
    }
}
```

**Grouping Rules**:
1. Open PR when work session starts (first commit on new branch)
2. Keep adding commits until:
   - Inactivity gap > N minutes (developer-specific)
   - Max commits per PR reached (seniority-based)
   - Random early close triggered (volatility)
3. Finalize PR metrics and store the envelope

**Correlation Enforcement**:

| Correlation | Enforcement Point | Mechanism |
|-------------|-------------------|-----------|
| Seniority → PR Size | Session.TargetLoC | Sample from seniority-specific distribution |
| Seniority → Commits/PR | Session.MaxCommits | Juniors: 2-5, Seniors: 5-12 |
| Working Hours → Gap | Session.InactivityGap | Clip to developer's work schedule |
| AI Ratio → Review Iterations | PR generation | Higher AI → more iterations (probabilistic) |

**Memory Efficiency**:
- Persist only the PR envelope (id, timestamps, author, repo, branch, commit list, summary metrics)
- Do not copy full commit data into PR storage
- Support continuous run without unbounded memory growth

### 6.5 Quality Correlation Enforcement (Phase 3C)

> **Design Decision**: Use probabilistic enforcement with sigmoid risk scoring, not deterministic.

Deterministic "high AI ⇒ revert" would look artificial and distort aggregates. Instead:

```go
func CalculateRevertRisk(pr models.PullRequest, dev seed.Developer) float64 {
    // Sigmoid: high AI + low seniority + high volatility → higher risk
    rawScore := a*pr.AIRatio + b*volatility + c*seniorityPenalty(dev.Seniority)
    return 1.0 / (1.0 + math.Exp(-rawScore))
}

func ShouldRevert(pr models.PullRequest, dev seed.Developer, rng *rand.Rand) bool {
    risk := CalculateRevertRisk(pr, dev)
    return rng.Float64() < risk  // Bernoulli sampling
}
```

**Benefits**:
- Correlations hold at population level (statistically significant over 1000+ PRs)
- Individual PRs retain realistic variability
- No artificial "high AI always reverts" patterns

### 6.6 Code Survival Tracking (Phase 3C)

> **Design Decision**: File-level survival tracking (simple, fast, sufficient for research).

**File Survival Model**:

```go
type FileSurvival struct {
    FilePath        string    `json:"file_path"`
    RepoName        string    `json:"repo_name"`
    CreatedAt       time.Time `json:"created_at"`      // First commit timestamp
    LastModifiedAt  time.Time `json:"last_modified_at"`
    AILinesAdded    int       `json:"ai_lines_added"`
    HumanLinesAdded int       `json:"human_lines_added"`
    TotalLines      int       `json:"total_lines"`
    RevertEvents    int       `json:"revert_events"`
    IsDeleted       bool      `json:"is_deleted"`
    DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}
```

**Key Decisions**:

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Greenfield threshold | First commit timestamp for file | OS file creation is meaningless in simulator |
| Survival granularity | File-level | Simple, fast, good enough for research |
| Line-level tracking | Deferred | Add only if very specific metrics needed |
| Survival windows | 30d, 60d, 90d | Standard cohort intervals |

**Calculation**:
- Track each file from first appearance (birth) to deletion (death) or observation date
- `survival_rate` = files_surviving / files_added_in_cohort
- Aggregate AI vs human lines per file for correlation analysis

### 6.7 Greenfield Index Calculation

> **Design Decision**: Greenfield = file created < 30 days before commit timestamp.

```go
func IsGreenfield(file models.CommitFile, commitTime time.Time) bool {
    fileAge := commitTime.Sub(file.CreatedAt)
    return fileAge < 30 * 24 * time.Hour
}

func CalculateGreenfieldIndex(pr models.PullRequest) float64 {
    var greenfieldLines, totalLines int
    for _, file := range pr.Files {
        if IsGreenfield(file, pr.CreatedAt) {
            greenfieldLines += file.Additions
        }
        totalLines += file.Additions
    }
    if totalLines == 0 {
        return 0
    }
    return float64(greenfieldLines) / float64(totalLines)
}
```

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

### 10.1 Glossary (v2.1 Additions)

| Term | Definition |
|------|------------|
| DataDesigner | NVIDIA NeMo tool for generating seed data |
| Greenfield Index | % of PR lines in files < 30 days old |
| Pickup Time | Hours from PR open to first review |
| Rework Ratio | LoC changed during review / initial LoC |
| Scope Creep | (final additions - initial) / final |
| Survival Rate | % of files still present after N days |
| Session | A work period that produces a PR from grouped commits |
| Risk Score | Sigmoid-based probability for quality outcomes |
| Hotfix | Fix-PR within 48 hours of original PR merge |
| Revert Chain | Sequence linking original PR to its revert commit |
| File Birth | First commit timestamp containing the file path |
| Cohort Window | Time period for grouping files in survival analysis |

### 10.2 References

- Cursor Business API: https://cursor.com/docs/account/teams/analytics-api
- Cursor AI Code Tracking: https://docs.cursor.com/business/api-reference/ai-code-tracking
- NVIDIA NeMo DataDesigner: https://github.com/NVIDIA/NeMo
- GitHub REST API: https://docs.github.com/en/rest
- **Methods Proposal**: `docs/design/External - Methods Proposal - AI on SDLC Study.md` - Scientific framework for SDLC metrics
- **GitHub Sim API**: `cursor-analytics-platform-research/packages/shared-schemas/openapi/github-sim-api.yaml`

---

## 11. Integration Testing & Lessons Learned (v2.2)

### 11.1 Integration Testing Results (January 4, 2026)

**Status**: ✅ COMPLETE - Full stack P4 → P5 → P6 integration verified

**Testing Architecture**:
```
cursor-sim (Docker)  →  cursor-analytics-core (Docker)  →  cursor-viz-spa (npm local)
   Port 8080         →       Port 4000 (GraphQL)        →        Port 3000
                     →       Port 5432 (PostgreSQL)     →
```

**Data Flow Verified**:
1. ✅ P4 generates simulated data (3 developers, 2 teams, 7 usage events)
2. ✅ P5 PostgreSQL stores developer and event data
3. ✅ P5 GraphQL serves `dashboardSummary` query with aggregated data
4. ✅ P6 Dashboard fetches and displays data via Apollo Client
5. ✅ Full stack renders: KPI cards, velocity heatmap, team radar, developer table

### 11.2 Critical Issues Discovered & Resolved

#### Issue 1: Dashboard Component Not Integrated (commit 57dc089)
**Problem**: Dashboard.tsx was still a placeholder from initial scaffolding, never integrated with hooks or chart components built in later tasks.

**Impact**:
- Dashboard showed "Chart placeholder" text instead of data
- No GraphQL POST requests visible in browser Network tab
- Components and hooks existed but were disconnected

**Root Cause**: Task completion checklist didn't verify end-to-end integration, only individual component creation.

**Fix**: Updated Dashboard.tsx to:
- Import and call `useDashboard()` hook
- Add loading/error states
- Render KPI cards with real data from GraphQL
- Pass data to VelocityHeatmap, TeamRadarChart, DeveloperTable components

**Lesson**: **Task checklists must verify end-to-end integration, not just unit creation.**

---

#### Issue 2: Import/Export Mismatches (commit 293f4fc)
**Problem**: Chart components used `export default` but Dashboard imported them as named exports `{ Component }`.

**Impact**: `Uncaught SyntaxError: The requested module does not provide an export named 'DeveloperTable'`

**Fix**: Changed all chart imports from `import { Component }` to `import Component`.

**Lesson**: **Enforce consistent export style (all default OR all named) in ESLint config.**

---

#### Issue 3: Component Prop Type Mismatches (commit 26d3567)
**Problem**: Dashboard created custom data objects (e.g., `{ date, count, level }`) that didn't match component prop interfaces (e.g., `DailyStats[]`).

**Impact**: TypeScript errors, components received wrong data shape.

**Fix**: Pass data directly without transformation, matching component interfaces exactly.

**Lesson**: **Component integration tests should validate prop contracts, not just component rendering.**

---

#### Issue 4: GraphQL Schema Mismatches (commit 2dfd06b) ⚠️ **CRITICAL**

**Problem**: P6 GraphQL queries were manually defined in `src/graphql/types.ts` based on design docs, but didn't match P5's actual implemented schema.

**Mismatches**:
| P6 Query (Manual) | P5 Schema (Actual) | Impact |
|-------------------|-------------------|---------|
| `topPerformers: Developer[]` | `topPerformer?: Developer` | Field doesn't exist error |
| `humanLinesAdded: Int` | `linesAdded: Int` | Field doesn't exist error |
| `aiLinesDeleted: Int` | Not in schema | Field doesn't exist error |
| `humanLinesAdded: Int` | `totalLinesAdded: Int` | Wrong field name |

**Impact**: **Complete integration failure**
- 400 Bad Request errors from GraphQL server
- Browser console: `[GraphQL error]: Cannot query field "topPerformers" on type "TeamStats"`
- Dashboard showed error state, no data rendered

**Root Cause**:
- P6 defined types manually based on outdated design docs
- P5 schema evolved during implementation (singular `topPerformer`, not plural)
- No automated validation between P6 queries and P5 schema
- TypeScript provided false sense of type safety (types matched local definitions, not server schema)

**Fix**: Manually compared P6 `queries.ts` with P5 `schema.ts` and aligned all field names and types.

**Lesson**: **NEVER manually define GraphQL types in client code. Always auto-generate from server schema.**

**Prevention Strategy**: See `docs/data-contract-testing.md` for comprehensive mitigation plan.

---

### 11.3 Testing Gaps Identified

Based on integration testing, the following test coverage gaps exist:

| Gap | Current | Impact | Proposed Solution |
|-----|---------|--------|-------------------|
| **Schema Contract Tests** | None | GraphQL 400 errors at runtime | GraphQL Code Generator + GraphQL Inspector |
| **Component Integration Tests** | Incomplete | Dashboard placeholders not caught | Page-level integration tests |
| **E2E Full Stack Tests** | None | Integration failures discovered manually | Playwright E2E tests |
| **Visual Regression Tests** | None | UI placeholder issues not caught | Playwright snapshots |
| **Pre-commit Schema Validation** | None | Schema drift enters codebase | Husky + codegen validation |

**See**: `docs/e2e-testing-strategy.md` for comprehensive testing enhancement plan.

---

### 11.4 Data Contract Testing Strategy

**Problem Statement**: GraphQL schema drift between services caused complete integration failure.

**Solution**: Automated schema validation and type generation.

**4-Phase Mitigation Plan**:

**Phase 1: Automated Type Generation (HIGH PRIORITY)**
- Install GraphQL Code Generator in P6
- Auto-generate TypeScript types from P5 GraphQL schema
- Delete manual `types.ts` file
- Add `predev` hook to run codegen before every dev server start
- **Impact**: Eliminates 100% of schema mismatch errors at compile-time

**Phase 2: Pre-Commit Schema Validation**
- Add Husky pre-commit hook to validate generated types are committed
- Prevent PRs with outdated schema from entering main branch
- **Impact**: 100% prevention of schema drift in version control

**Phase 3: CI/CD Schema Compatibility Check**
- GitHub Actions workflow to validate P6 queries against running P5 server
- Block PRs with incompatible GraphQL queries
- **Impact**: Automated validation, no manual review needed

**Phase 4: Schema Registry (Apollo Studio or Self-Hosted)**
- Centralized schema management with versioning
- Breaking change detection
- Team notifications on schema changes
- **Impact**: Long-term schema governance

**See**: `docs/data-contract-testing.md` for full implementation guide.
**See**: `docs/MITIGATION-PLAN.md` for executive summary and rollout timeline.

---

### 11.5 Recommended Testing Pyramid

```
                         ┌─────────────┐
                         │   E2E Tests │  ← 5% (Playwright, full stack)
                         │  (Proposed) │
                         └─────────────┘
                       ┌────────────────────┐
                       │ Integration Tests  │  ← 15% (Page-level, GraphQL)
                       │  (Need expansion)  │
                       └────────────────────┘
                  ┌──────────────────────────────┐
                  │   Component Tests            │  ← 30% (React Testing Lib)
                  │   (Current: Complete)        │
                  └──────────────────────────────┘
            ┌────────────────────────────────────────┐
            │   Unit Tests                           │  ← 50% (Vitest, Go test)
            │   (Current: 91% coverage)              │
            └────────────────────────────────────────┘
```

**Current Coverage**:
- ✅ Unit Tests: 91% (P5), 91.68% (P6), 80%+ (P4)
- ✅ Component Tests: 162 tests (P6)
- ⚠️ Integration Tests: Incomplete (need Dashboard integration tests)
- ❌ E2E Tests: Missing (need Playwright)
- ❌ Visual Regression: Missing (need Playwright snapshots)
- ❌ Schema Contract Tests: Missing (need GraphQL Code Generator)

**Target State (6 months)**:
- Unit: 90%+ (maintain)
- Component: 30% of test suite
- Integration: 15% of test suite (page-level + GraphQL schema validation)
- E2E: 5% of test suite (critical user journeys)
- Visual Regression: 100% of pages
- Schema Contract: 100% automated validation

---

### 11.6 Key Takeaways

> **"Never manually define GraphQL types in client code. Always generate from server schema."**

> **"TypeScript type safety is only as good as your runtime contract validation."**

> **"Automate everything that can drift: types, schemas, contracts, tests."**

> **"Task completion means end-to-end integration, not just individual component creation."**

**Success Metrics**:

| Metric | Before (v2.1) | After (v2.2) | Target |
|--------|---------------|--------------|--------|
| Schema mismatch errors | Runtime (browser) | Compile-time (IDE) | Zero runtime errors |
| Integration test coverage | 0% | Page tests added | 15% of suite |
| Schema validation | Manual | Automated (codegen) | 100% automated |
| Time to detect integration issues | Days (manual testing) | Minutes (CI failure) | < 5 minutes |

---

## 12. Updated Documentation (v2.2)

### 12.1 New Documentation Created

- **docs/data-contract-testing.md**: Comprehensive GraphQL schema validation strategy (600+ lines)
- **docs/e2e-testing-strategy.md**: E2E and integration testing enhancement plan (400+ lines)
- **docs/MITIGATION-PLAN.md**: Executive summary of mitigation strategies (400+ lines)

### 12.2 Updated Documentation

- **docs/INTEGRATION.md**: Added current architecture, troubleshooting for all 4 integration issues, data contract testing section
- **.work-items/P6-cursor-viz-spa/task.md**: Documented integration issues with root causes, fixes, lessons learned
- **.claude/DEVELOPMENT.md**: Added P5+P6 integration testing section with status and next steps
- **docs/DESIGN.md** (this file): Updated architecture diagram, added integration testing section

### 12.3 Documentation Hierarchy

| Document | Purpose | Audience |
|----------|---------|----------|
| **DESIGN.md** (this file) | System architecture and design decisions | Architects, new developers |
| **INTEGRATION.md** | Integration testing guide and troubleshooting | QA, DevOps, developers |
| **data-contract-testing.md** | GraphQL schema validation strategy | Backend + frontend developers |
| **e2e-testing-strategy.md** | E2E testing implementation plan | QA, test engineers |
| **MITIGATION-PLAN.md** | Executive summary for leadership | Product managers, tech leads |
| **services/{service}/SPEC.md** | Technical specification per service | Service developers |
| **.work-items/{feature}/task.md** | Task-level implementation tracking | Individual contributors |

---

## 13. Data Pipeline Architecture (P8/P9)

### 13.1 Alternative Analytics Path: dbt + DuckDB/Snowflake

Beyond the primary GraphQL path (cursor-sim → analytics-core → viz-spa), the platform includes an alternative analytics implementation using dbt and DuckDB for local development, with Snowflake for production:

```
Path 1 (GraphQL):   cursor-sim → cursor-analytics-core → cursor-viz-spa
                    (P4)        (P5 TypeScript/GraphQL)   (P6 React)

Path 2 (dbt):       cursor-sim → api-loader → dbt → streamlit-dashboard
                    (P4)        (P8 Python)   (SQL)  (P9 Python)
```

**Key Distinction**:
- **Path 1 (GraphQL)**: Type-safe, real-time aggregations, complex business logic in TypeScript
- **Path 2 (dbt)**: SQL-first transformations, reproducible data pipeline, analytics-optimized

### 13.2 Data Contract Hierarchy

**cursor-sim is the authoritative source of truth**. All downstream layers validate against the API contract:

```
LEVEL 1: API CONTRACT (cursor-sim SPEC.md) ← SOURCE OF TRUTH
  • Endpoints: /analytics/ai-code/commits, /repos/*/pulls, /research/dataset
  • Response format: {items: [...], totalCount, page, pageSize}
  • Field names: camelCase (commitHash, userEmail, tabLinesAdded, composerLinesAdded, commitTs)
  • Pagination: Cursor-based with configurable page size

LEVEL 2: DATA TIER CONTRACT (api-loader → dbt → DuckDB/Snowflake)
  • Raw schema (raw_*): Preserves API fields exactly (camelCase)
  • Staging schema (stg_*): Transforms camelCase → snake_case, validates types
  • Mart schema (mart_*): Aggregations for analytics (velocity, ai_impact, quality, review_costs)

LEVEL 3: DASHBOARD CONTRACT (Streamlit)
  • Queries: SELECT from mart_* only, never raw or staging
  • Parameters: Parameterized queries ($param syntax)
  • Security: SQL injection prevention via parameter binding
```

### 13.3 Development vs Production Parity

| Layer | Development | Production | Parity |
|-------|-------------|-----------|--------|
| **API Source** | cursor-sim | Cursor + GitHub APIs | ✅ Same contract |
| **Extraction** | api-loader (Python) | SnapLogic | ✅ Same logic, tested in CI |
| **Landing** | Parquet files | Snowflake Stage | ✅ Same format |
| **Raw Tables** | DuckDB | Snowflake | ✅ Same schema |
| **Transforms** | dbt (DuckDB dialect) | dbt (Snowflake dialect) | ✅ Identical SQL |
| **Marts** | DuckDB | Snowflake | ✅ Same schema |
| **Dashboard** | Streamlit (local) | Streamlit (Cloud) | ✅ Identical |

### 13.4 DuckDB Schema Naming (Critical)

DuckDB requires the `main_` prefix for schema-qualified table names:

```sql
-- ✅ CORRECT
SELECT * FROM main_raw.commits
SELECT * FROM main_staging.stg_commits
SELECT * FROM main_mart.mart_velocity

-- ❌ INCORRECT (fails with "Catalog Error")
SELECT * FROM raw.commits
SELECT * FROM staging.stg_commits
SELECT * FROM mart.mart_velocity
```

This is a DuckDB-specific requirement where `main` is the default catalog. Queries without the `main_` prefix attempt to reference non-existent catalogs.

### 13.5 Lessons Learned: Data Pipeline Implementation

**Issue #1: API Response Format Duality**
- **Problem**: cursor-sim supports both `{items:[...]}` and raw array formats; api-loader initially didn't handle both
- **Root Cause**: API contract wasn't explicitly documented; assumption that only one format was supported
- **Resolution**: Implemented dual-format detection in BaseAPIExtractor
- **Lesson**: Document all supported response formats in API contract

**Issue #2: Column Mapping (camelCase vs snake_case)**
- **Problem**: API returns camelCase (commitHash); dashboard queries expected snake_case (commit_hash)
- **Root Cause**: Confusion about where transformation should happen (extraction vs staging)
- **Resolution**: Preserve API format in raw layer, transform in dbt staging models
- **Lesson**: Establish clear responsibility boundaries: API fields as-is → staging transforms → marts aggregate

**Issue #3: DuckDB Schema Naming**
- **Problem**: Dashboard queries failed with "Catalog Error: Table not found"; SQL used `mart.table` instead of `main_mart.table`
- **Root Cause**: DuckDB catalog/schema concept differs from standard SQL; documentation didn't specify DuckDB requirements
- **Resolution**: Added `main_` prefix requirement to all schema-qualified table names
- **Lesson**: Document database-specific quirks (e.g., DuckDB catalog naming) in early design docs

**Issue #4: Column Availability Mismatch**
- **Problem**: Dashboard tried to access columns that don't exist in dbt marts (p50_cycle_time, avg_coding_lead_time, avg_review_iterations)
- **Root Cause**: Assumed columns existed based on Snowflake documentation; dbt mart didn't define them
- **Resolution**: Removed references to non-existent columns; verified actual marts via `SELECT * FROM main_mart.mart_*`
- **Lesson**: Test column availability early; don't assume columns exist without verification

**Issue #5: INTERVAL Syntax in DuckDB**
- **Problem**: Parameterized INTERVAL syntax failed: `WHERE week >= CURRENT_DATE - INTERVAL $days DAY`
- **Root Cause**: DuckDB doesn't support parameterized INTERVAL expressions
- **Resolution**: Use f-string interpolation for days (integer is validated): `f"CURRENT_DATE - INTERVAL '{days}' DAY"`
- **Lesson**: Document database-specific SQL syntax limitations

**Issue #6: SQL Injection Prevention**
- **Problem**: Early dashboard code used f-strings for user input: `f"WHERE repo_name = '{repo_name}'"`
- **Root Cause**: Code reuse pattern from pandas/Python without SQL context
- **Resolution**: Refactored to parameterized queries with DuckDB's `$param` placeholders
- **Lesson**: Always parameterize user input; never use f-strings for SQL user input

---

## 14. Recommended Reading Order

For new developers joining the project, read in this order:

1. **This file (DESIGN.md)**: Overall architecture and design decisions
2. **services/cursor-sim/SPEC.md**: API contract and data formats (source of truth)
3. **docs/design/new_data_architecture.md**: Data pipeline details (P8/P9)
4. **docs/TESTING_STRATEGY.md**: Testing approaches and data contract validation
5. **services/{service}/README.md**: Service-specific setup and development
6. **.work-items/{feature}/design.md**: Active feature design details
7. **.work-items/{feature}/task.md**: Implementation task breakdown

