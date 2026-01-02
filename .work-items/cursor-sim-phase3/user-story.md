# User Story: cursor-sim Phase 3 (Research Framework & Completeness)

## Story

**As a** data scientist / researcher
**I want** cursor-sim to export pre-joined research datasets and support replay mode
**So that** I can conduct reproducible SDLC research without regenerating data every time

## Background

Phase 1 delivered the core Cursor API simulation with commit generation.
Phase 2 added GitHub PR simulation with review cycles and quality outcomes.
Phase 3 completes cursor-sim with research-focused capabilities:

1. **Research Dataset Export** - Pre-joined data ready for analysis
2. **Code Survival Tracking** - Measure long-term retention of AI-generated code
3. **Replay Mode** - Serve pre-generated datasets from corpus files

## Acceptance Criteria

### AC-1: Research Dataset Export (SIM-R013)

**Given** cursor-sim has generated commits and PRs
**When** I request `/research/dataset?format=parquet`
**Then** I receive a Parquet file with pre-joined SDLC metrics

**Columns must include:**
- PR metadata: pr_number, author_email, repo_name
- AI metrics: ai_lines_added, ai_ratio, pr_volume
- Timing: coding_lead_time_hours, pickup_time_hours, review_lead_time_hours
- Review: review_density, iterations, rework_ratio, scope_creep
- Quality: is_reverted, survival_rate_30d, has_hotfix_followup
- Context: author_seniority, repo_age_days, primary_language

**Also support:**
- JSON format: `/research/dataset?format=json`
- CSV format: `/research/dataset?format=csv`
- Date filtering: `?startDate=2026-01-01&endDate=2026-01-31`

### AC-2: Research Metrics Endpoints (SIM-R013)

**Given** research data is available
**When** I request specific metric endpoints
**Then** I receive aggregated metrics

**Endpoints:**
- `GET /research/metrics/velocity` - Developer velocity by AI ratio
- `GET /research/metrics/review-costs` - Review time vs AI ratio
- `GET /research/metrics/quality` - Quality outcomes by AI ratio

### AC-3: Code Survival Tracking (SIM-R014)

**Given** cursor-sim has generated commits over 30+ days
**When** code survival is calculated
**Then** survival rates are computed at 7d, 14d, 30d intervals

**Requirements:**
- Track which lines from commit N still exist at N+7, N+14, N+30 days
- Calculate: `survival_rate = lines_remaining / lines_added`
- Support per-commit and per-PR survival
- Survival should correlate with AI ratio (based on seed correlations)

**Performance:**
- Must handle 100k commits without excessive memory
- Survival calculation should complete in < 5 seconds

### AC-4: Replay Mode (SIM-R015)

**Given** a pre-generated Parquet corpus file
**When** I start cursor-sim with `--mode=replay --corpus=events.parquet`
**Then** all APIs serve data from the corpus instead of generating

**Requirements:**
- Load Parquet file using parquet-go library
- Support all existing Cursor API endpoints
- Support all GitHub API endpoints
- Support research endpoints
- Startup time < 1 second for 1M events
- Read-only mode (no POST/PUT/DELETE)

### AC-5: Integration Testing

**Given** Phase 3 features are implemented
**When** E2E tests run
**Then** all tests pass with >80% coverage

**Test scenarios:**
- Export to all formats (JSON, CSV, Parquet)
- Load and validate Parquet structure
- Verify survival rate calculations
- Replay mode serves identical data to runtime mode
- Performance benchmarks pass

## Out of Scope

- Real-time streaming of events
- External database storage (remains in-memory)
- Advanced time-series forecasting
- Multi-corpus merging

## Dependencies

- Phase 1 (Complete): Cursor API, commit generation
- Phase 2 (Complete): GitHub PR simulation, review cycles
- Go parquet library: github.com/xitongsys/parquet-go

## Success Metrics

- Research dataset exports successfully with all columns
- Code survival tracking correlates with AI ratio
- Replay mode starts in <1s and serves data correctly
- All E2E tests pass
- Test coverage >80%

## Technical Notes

### Parquet Schema

```
message ResearchDataset {
  required int64 pr_number;
  required binary author_email (UTF8);
  required binary repo_name (UTF8);
  required int32 ai_lines_added;
  required double ai_ratio;
  required int32 pr_volume;
  required double greenfield_index;
  required double coding_lead_time_hours;
  required double pickup_time_hours;
  required double review_lead_time_hours;
  required double review_density;
  required int32 iterations;
  required double rework_ratio;
  required double scope_creep;
  required boolean is_reverted;
  required double survival_rate_30d;
  required boolean has_hotfix_followup;
  required binary author_seniority (UTF8);
  required int32 repo_age_days;
  required binary primary_language (UTF8);
  required int64 created_at (TIMESTAMP_MILLIS);
  required int64 merged_at (TIMESTAMP_MILLIS);
}
```

### Survival Tracking Algorithm

For each commit C at time T:
1. Extract all lines added: L_added
2. At T+7d, check which lines still exist: L_7d
3. At T+14d, check which lines still exist: L_14d
4. At T+30d, check which lines still exist: L_30d
5. Calculate rates: survival_7d = L_7d / L_added (etc.)

**Optimization:** Use file path + line number hashing for efficient lookup

## Related Documents

- `services/cursor-sim/SPEC.md` - Technical specification
- `docs/FEATURES.md` - Feature details
- `.work-items/cursor-sim-phase3/design.md` - Design decisions
- `.work-items/cursor-sim-phase3/task.md` - Implementation tasks
