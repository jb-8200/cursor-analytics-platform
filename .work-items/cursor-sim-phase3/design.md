# Design Document: cursor-sim Phase 3

## Overview

Phase 3 completes cursor-sim with research-focused capabilities:
1. Research dataset export (Parquet/JSON/CSV)
2. Code survival tracking
3. Replay mode for reproducible research

## Architecture Decisions

### AD-1: Parquet Library Choice

**Decision:** Use `github.com/parquet-go/parquet-go` (v4+)

**Rationale:**
- More actively maintained than xitongsys/parquet-go
- Better performance for large datasets
- Native Go implementation, no CGO dependencies
- Supports schema evolution

**Alternatives Considered:**
- `github.com/xitongsys/parquet-go` - Less active maintenance
- Apache Arrow Go - Overkill for our needs
- Manual Parquet writing - Too complex

### AD-2: Research Dataset Schema

**Decision:** Single flat table with all metrics pre-joined

**Rationale:**
- Optimized for data science workflows (pandas, R)
- No JOIN operations required by analysts
- Each row = one PR with all context
- Denormalized for query performance

**Schema:**
```go
type ResearchRow struct {
    // PR Identity
    PRNumber       int64     `parquet:"name=pr_number, type=INT64"`
    AuthorEmail    string    `parquet:"name=author_email, type=BYTE_ARRAY, convertedtype=UTF8"`
    RepoName       string    `parquet:"name=repo_name, type=BYTE_ARRAY, convertedtype=UTF8"`

    // AI Metrics
    AILinesAdded   int32     `parquet:"name=ai_lines_added, type=INT32"`
    AIRatio        float64   `parquet:"name=ai_ratio, type=DOUBLE"`
    PRVolume       int32     `parquet:"name=pr_volume, type=INT32"`
    GreenfieldIdx  float64   `parquet:"name=greenfield_index, type=DOUBLE"`

    // Timing Metrics
    CodingLeadTime float64   `parquet:"name=coding_lead_time_hours, type=DOUBLE"`
    PickupTime     float64   `parquet:"name=pickup_time_hours, type=DOUBLE"`
    ReviewLeadTime float64   `parquet:"name=review_lead_time_hours, type=DOUBLE"`

    // Review Metrics
    ReviewDensity  float64   `parquet:"name=review_density, type=DOUBLE"`
    Iterations     int32     `parquet:"name=iterations, type=INT32"`
    ReworkRatio    float64   `parquet:"name=rework_ratio, type=DOUBLE"`
    ScopeCreep     float64   `parquet:"name=scope_creep, type=DOUBLE"`

    // Quality Metrics
    IsReverted     bool      `parquet:"name=is_reverted, type=BOOLEAN"`
    Survival30d    float64   `parquet:"name=survival_rate_30d, type=DOUBLE"`
    HasHotfix      bool      `parquet:"name=has_hotfix_followup, type=BOOLEAN"`

    // Context
    Seniority      string    `parquet:"name=author_seniority, type=BYTE_ARRAY, convertedtype=UTF8"`
    RepoAgeDays    int32     `parquet:"name=repo_age_days, type=INT32"`
    Language       string    `parquet:"name=primary_language, type=BYTE_ARRAY, convertedtype=UTF8"`

    // Timestamps
    CreatedAt      int64     `parquet:"name=created_at, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
    MergedAt       int64     `parquet:"name=merged_at, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
}
```

### AD-3: Code Survival Tracking

**Decision:** File-based line tracking with hash-based matching

**Algorithm:**
1. When commit added, store: `{file_path, line_hash, commit_hash, timestamp}`
2. When later commit modifies file, recalculate hashes
3. Match by content hash (not line number - handles insertions)
4. At T+N days, count surviving hashes

**Data Structure:**
```go
type LineTracker struct {
    mu sync.RWMutex

    // Map: file_path -> []TrackedLine
    linesByFile map[string][]*TrackedLine

    // Index: commit_hash -> []TrackedLine
    linesByCommit map[string][]*TrackedLine
}

type TrackedLine struct {
    CommitHash  string
    FilePath    string
    LineHash    string    // SHA256 of trimmed line content
    AddedAt     time.Time
    RemovedAt   *time.Time // nil if still exists
    IsAICode    bool
}
```

**Optimization:** Only track lines from commits with AI attribution

**Memory Estimate:** ~200 bytes/line × 100k commits × 50 lines/commit = ~1GB (acceptable)

### AD-4: Replay Mode Architecture

**Decision:** In-memory Parquet loading with same storage interface

**Implementation:**
```go
type ReplayStore struct {
    commits    []*models.Commit
    prs        []*models.PullRequest
    reviews    []*models.ReviewComment
    developers map[string]*seed.Developer

    // Same indices as MemoryStore
    commitsByHash   map[string]*models.Commit
    commitsByUser   map[string][]*models.Commit
    prsByRepo       map[string][]*models.PullRequest
}

func NewReplayStore(corpusPath string) (*ReplayStore, error) {
    // Read Parquet file
    // Populate indices
    // Return read-only store
}
```

**Interface Compliance:** Implements same `storage.Store` interface as MemoryStore

**Benefits:**
- API handlers don't need changes
- Same query performance as runtime mode
- Deterministic results for testing

### AD-5: Research API Endpoints

**Decision:** New `/research/*` endpoint group

**Endpoints:**
```
GET /research/dataset?format={json|csv|parquet}&startDate=X&endDate=Y
GET /research/metrics/velocity
GET /research/metrics/review-costs
GET /research/metrics/quality
```

**Response Caching:** No caching in MVP (data is static per run)

### AD-6: Export Performance

**Decision:** Streaming export for large datasets

**Implementation:**
```go
func (h *ResearchHandler) ExportDataset(w http.ResponseWriter, r *http.Request) {
    format := r.URL.Query().Get("format")

    switch format {
    case "parquet":
        w.Header().Set("Content-Type", "application/octet-stream")
        w.Header().Set("Content-Disposition", "attachment; filename=research_dataset.parquet")

        pw := parquet.NewWriter(w, new(ResearchRow))
        defer pw.Close()

        for _, row := range h.buildDataset() {
            pw.Write(row)
        }

    case "csv":
        w.Header().Set("Content-Type", "text/csv")
        cw := csv.NewWriter(w)
        for _, row := range h.buildDataset() {
            cw.Write(rowToCSV(row))
        }
        cw.Flush()
    }
}
```

**Memory:** Streaming avoids loading entire dataset in memory

## Component Design

### Package Structure

```
internal/
├── export/
│   ├── parquet.go          # Parquet writer
│   ├── csv.go              # CSV writer
│   └── research.go         # Dataset builder
├── survival/
│   ├── tracker.go          # Line tracking
│   ├── calculator.go       # Survival rate computation
│   └── matcher.go          # Content-based line matching
├── replay/
│   ├── loader.go           # Parquet corpus loader
│   └── store.go            # ReplayStore implementation
└── api/
    └── research/
        ├── dataset.go      # Dataset export handlers
        └── metrics.go      # Metrics aggregation handlers
```

### Key Interfaces

```go
// export/research.go
type DatasetBuilder interface {
    BuildDataset(from, to time.Time) []ResearchRow
    FilterByDateRange(rows []ResearchRow, from, to time.Time) []ResearchRow
}

// survival/tracker.go
type SurvivalTracker interface {
    TrackCommit(commit *models.Commit) error
    CalculateSurvival(commitHash string, asOf time.Time) SurvivalRates
    GetSurvivalByPR(prNumber int64) SurvivalRates
}

type SurvivalRates struct {
    Survival7d  float64
    Survival14d float64
    Survival30d float64
}

// replay/store.go
type CorpusLoader interface {
    Load(path string) (storage.Store, error)
    Validate(store storage.Store) error
}
```

## Data Flow

### Research Export Flow

```
1. HTTP Request → /research/dataset?format=parquet
2. ResearchHandler.ExportDataset()
3. DatasetBuilder.BuildDataset()
   ├─ Fetch all PRs from storage
   ├─ For each PR:
   │  ├─ Get commits
   │  ├─ Get reviews
   │  ├─ Calculate AI metrics
   │  ├─ Calculate timing metrics
   │  ├─ Get survival rates
   │  └─ Build ResearchRow
   └─ Return []ResearchRow
4. Stream to Parquet writer
5. HTTP Response
```

### Survival Calculation Flow

```
1. Commit Generation
   ├─ For each line added
   │  └─ tracker.TrackLine(commit, file, lineHash, isAI)
   └─ Store in linesByFile index

2. Later Commit Modifies File
   ├─ Read current file lines
   ├─ Hash each line
   ├─ Compare with tracked lines
   ├─ Mark missing lines as removed
   └─ Update linesByFile

3. Query Time
   ├─ calculator.GetSurvival(commitHash, asOf)
   ├─ Filter lines: addedAt <= asOf
   ├─ Count: still_exists vs removed
   └─ Return survival_rate
```

### Replay Mode Flow

```
1. Startup with --mode=replay --corpus=events.parquet
2. ReplayLoader.Load(corpusPath)
   ├─ Open Parquet file
   ├─ Read all rows
   ├─ Deserialize to Go structs
   └─ Build indices
3. Initialize ReplayStore (implements storage.Store)
4. Start HTTP server with ReplayStore
5. API requests served from static data
```

## Testing Strategy

### Unit Tests

- `export/research_test.go` - Dataset building
- `export/parquet_test.go` - Parquet writing
- `survival/tracker_test.go` - Line tracking
- `survival/calculator_test.go` - Survival math
- `replay/loader_test.go` - Corpus loading

### Integration Tests

- Export full dataset and validate schema
- Track survival across multiple commits
- Load corpus and query via API

### E2E Tests

- Generate data → Export → Load → Verify identical
- Runtime mode vs Replay mode equivalence
- Performance benchmarks

### Test Data

Create `testdata/research_corpus.parquet` with known values:
- 10 PRs
- 50 commits
- Known survival rates
- All metrics calculable by hand

## Performance Considerations

### Memory Usage

| Component | Estimate (100k commits) |
|-----------|------------------------|
| Commits | 200 MB |
| PRs | 50 MB |
| Reviews | 100 MB |
| Line Tracking | 1 GB |
| **Total** | **~1.5 GB** |

**Acceptable for research simulator**

### Export Performance

Target: Export 100k PRs in < 10 seconds

Optimization:
- Streaming write (no buffering)
- Pre-calculated metrics (no JOIN during export)
- Parquet compression (SNAPPY)

### Replay Startup

Target: < 1 second for 1M events

Optimization:
- Parquet columnar read (only needed columns)
- Concurrent index building
- Memory-mapped file reading

## Security Considerations

- No authentication changes (same Basic Auth)
- Replay mode is read-only (reject POST/PUT/DELETE)
- Parquet files validated before loading
- No arbitrary file path injection

## Migration Path

Phase 3 is additive only:
- No changes to existing APIs
- New packages don't affect Phase 1/2
- Can deploy Phase 3 independently

## Open Questions

1. **Q:** Should we support incremental survival updates?
   **A:** No, calculate at generation time (simpler, reproducible)

2. **Q:** Parquet compression codec?
   **A:** SNAPPY (balance of speed and size)

3. **Q:** Max corpus file size for replay?
   **A:** 10 GB uncompressed (~2-3 GB compressed)

## Related Documents

- `services/cursor-sim/SPEC.md` - Technical specification
- `.work-items/cursor-sim-phase3/user-story.md` - Requirements
- `.work-items/cursor-sim-phase3/task.md` - Implementation tasks
