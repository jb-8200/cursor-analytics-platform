# User Story: Research Framework

**Feature ID**: P3-F01-research-framework
**Phase**: P3 (cursor-sim Research Framework)
**Created**: January 3, 2026
**Status**: COMPLETE

## Overview

As a **SDLC researcher**, I want a **research data export framework** so that I can **extract correlated datasets for statistical analysis** of AI's impact on development velocity, review costs, and code quality.

## User Stories

### US-SIM-R013: Export Research Dataset

**As a** researcher analyzing AI impact on SDLC metrics
**I want** to export the complete dataset in CSV/Parquet format
**So that** I can load it into pandas/R for statistical analysis

**Acceptance Criteria**:
```gherkin
Given cursor-sim has generated PR and commit data
When I call GET /research/dataset?format=csv
Then I receive a CSV file with all research metrics
And the dataset includes AI ratio, cycle times, and quality outcomes
And all records have complete JOIN keys for cross-referencing
```

### US-SIM-R014: Velocity Metrics API

**As a** researcher studying development velocity
**I want** aggregated velocity metrics by AI usage intensity
**So that** I can test hypotheses about AI's impact on lead times

**Acceptance Criteria**:
```gherkin
Given the research dataset is populated
When I call GET /research/metrics/velocity
Then I receive velocity metrics grouped by AI ratio bands
And metrics include coding lead time, pickup time, and review lead time
```

### US-SIM-R015: Quality Metrics API

**As a** researcher studying code quality
**I want** aggregated quality metrics by AI usage intensity
**So that** I can test hypotheses about AI's impact on code survival and reverts

**Acceptance Criteria**:
```gherkin
Given the research dataset is populated
When I call GET /research/metrics/quality
Then I receive quality metrics grouped by AI ratio bands
And metrics include revert rate, survival rate, and hotfix rate
```

## Dependencies

- Phase 1 (P1-F01) foundation must be complete
- Seed data must include developer preferences and repo characteristics

---

**Status**: COMPLETE (1.75h actual / 15-20h estimated)
