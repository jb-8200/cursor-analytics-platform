# Task Breakdown: cursor-sim Phase 2 (GitHub PR Simulation)

## Overview

**Feature**: cursor-sim v2 Phase 2
**Total Estimated Hours**: 20-25
**Number of Steps**: 8
**Current Step**: None - NOT_STARTED

## Progress Tracker

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| 01 | PR and Review Models | 2.0 | NOT_STARTED | - |
| 02 | Seed Schema Extensions | 2.0 | NOT_STARTED | - |
| 03 | PR Storage Methods | 2.5 | NOT_STARTED | - |
| 04 | PR Generator | 4.0 | NOT_STARTED | - |
| 05 | Review Generator | 3.0 | NOT_STARTED | - |
| 06 | Quality Outcomes | 2.0 | NOT_STARTED | - |
| 07 | GitHub API Handlers | 4.0 | NOT_STARTED | - |
| 08 | Integration & E2E Tests | 3.0 | NOT_STARTED | - |

## Dependency Graph

```
Step 01 (Models)
    │
    ├── Step 02 (Seed Extensions)
    │       │
    │       └── Step 04 (PR Generator)
    │               │
    │               └── Step 05 (Review Generator)
    │                       │
    │                       └── Step 06 (Quality Outcomes)
    │
    └── Step 03 (PR Storage)
            │
            └── Step 07 (GitHub API)
                    │
                    └── Step 08 (E2E Tests)
```

## Critical Path

01 → 02 → 04 → 05 → 06 → 08

## Step Details

### Step 01: PR and Review Models

**File**: `internal/models/pr.go`

**Tasks**:
- [ ] Define PullRequest struct with all fields
- [ ] Define ReviewComment struct
- [ ] Define Repository struct
- [ ] Add JSON tags matching GitHub API schema
- [ ] Write unit tests for model validation

**Acceptance Criteria**:
- Models compile and have correct JSON serialization
- Tests verify field constraints

---

### Step 02: Seed Schema Extensions

**Files**:
- `internal/seed/types.go`
- `internal/seed/loader.go`

**Tasks**:
- [ ] Add PRLifecycle to SeedData
- [ ] Add review_comments to TextTemplates
- [ ] Add ai_ratio_revert_rate to Correlations
- [ ] Update validation for new fields
- [ ] Update testdata/valid_seed.json

**Acceptance Criteria**:
- Seed loader accepts new PR-related fields
- Validation ensures required correlations exist

---

### Step 03: PR Storage Methods

**File**: `internal/storage/store.go`

**Tasks**:
- [ ] Add PR map to MemoryStore
- [ ] Implement StorePR(pr *models.PullRequest)
- [ ] Implement GetPRsByRepo(repoName string) []PullRequest
- [ ] Implement GetPRByNumber(repoName string, number int) *PullRequest
- [ ] Implement GetPRCommits(repoName string, number int) []Commit
- [ ] Add thread-safe access

**Acceptance Criteria**:
- All PR storage operations work correctly
- Thread-safety tests pass

---

### Step 04: PR Generator

**File**: `internal/generator/pr_generator.go`

**Tasks**:
- [ ] Implement commit clustering algorithm
- [ ] Generate PR metadata from commit cluster
- [ ] Calculate aggregated AI metrics
- [ ] Link commits to PR via pr_number field
- [ ] Generate PR timeline (created, updated, merged)
- [ ] Unit tests for all generation logic

**Acceptance Criteria**:
- PRs generated from commit clusters
- Commits properly linked to PRs
- Timeline respects seed parameters

---

### Step 05: Review Generator

**File**: `internal/generator/review_generator.go`

**Tasks**:
- [ ] Implement reviewer selection (same team, not author)
- [ ] Generate review comments from templates
- [ ] Simulate review iterations
- [ ] Generate approval/rejection decisions
- [ ] Unit tests

**Acceptance Criteria**:
- Reviews assigned to appropriate team members
- Comment count matches thoroughness settings
- Iteration count follows seed parameters

---

### Step 06: Quality Outcomes

**File**: `internal/generator/quality.go`

**Tasks**:
- [ ] Implement revert probability calculation
- [ ] Implement AI ratio categorization
- [ ] Apply quality signals to PRs
- [ ] Identify bug-fix commits/PRs
- [ ] Unit tests with deterministic seed

**Acceptance Criteria**:
- Revert rates match configured correlations
- Bug fixes properly identified
- Reproducible with same random seed

---

### Step 07: GitHub API Handlers

**Files**:
- `internal/api/github/repos.go`
- `internal/api/github/pulls.go`
- `internal/api/github/reviews.go`

**Tasks**:
- [ ] Implement GET /repos
- [ ] Implement GET /repos/{owner}/{repo}
- [ ] Implement GET /repos/{owner}/{repo}/pulls
- [ ] Implement GET /repos/{owner}/{repo}/pulls/{number}
- [ ] Implement GET /repos/{owner}/{repo}/pulls/{number}/commits
- [ ] Implement GET /repos/{owner}/{repo}/pulls/{number}/reviews
- [ ] Add pagination support
- [ ] Register routes in router
- [ ] Handler unit tests

**Acceptance Criteria**:
- All endpoints return GitHub-compatible JSON
- Pagination works correctly
- Auth required on all endpoints

---

### Step 08: Integration & E2E Tests

**Files**:
- `test/e2e/pr_test.go`
- Integration tests in various packages

**Tasks**:
- [ ] E2E test: full PR lifecycle
- [ ] E2E test: GitHub API endpoints
- [ ] E2E test: quality outcome correlations
- [ ] Integration test: seed → commits → PRs → reviews
- [ ] Verify schema compatibility with real GitHub API

**Acceptance Criteria**:
- All E2E tests pass
- API responses validate against GitHub schema
- Coverage > 80% for new packages

---

## Model Recommendations

| Step | Model | Rationale |
|------|-------|-----------|
| 01, 02, 03 | Haiku | Well-specified, low complexity |
| 04, 05, 06 | Sonnet | Generation logic complexity |
| 07, 08 | Sonnet | Multiple endpoints, integration |

## TDD Checklist (Per Step)

- [ ] Read step details and acceptance criteria
- [ ] Write failing test(s) for the step
- [ ] Run tests, confirm RED
- [ ] Implement minimal code to pass
- [ ] Run tests, confirm GREEN
- [ ] Refactor while green
- [ ] Run linter (golangci-lint)
- [ ] Update step status to DONE
- [ ] Commit with time tracking
