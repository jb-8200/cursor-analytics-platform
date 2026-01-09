# Task Breakdown: GitHub Simulation

**Feature ID**: P2-F01-github-simulation
**Status**: Not Started
**Created**: January 8, 2026

---

## Progress Tracker

| Phase | Tasks | Status | Estimated | Actual |
|-------|-------|--------|-----------|--------|
| **Models** | 3 | 1/3 | 2.5h | 0.5h |
| **Generators** | 3 | 0/3 | 8.0h | - |
| **Storage** | 2 | 0/2 | 3.0h | - |
| **API Handlers** | 5 | 0/5 | 5.0h | - |
| **Testing & Docs** | 2 | 0/2 | 3.5h | - |
| **TOTAL** | **15** | **1/15** | **22.0h** | **0.5h** |

---

## Task List

### PHASE 1: Data Models

#### TASK-GH-01: Create PullRequest Model (1.0h) ✅ COMPLETE

**Status**: COMPLETE
**Time**: 0.5h actual / 1.0h estimated
**Commit**: 3677559

**Goal**: Define PullRequest struct with relationships

**Files**:
- MODIFIED: `internal/models/pr.go` (enhanced existing model)
- MODIFIED: `internal/models/pr_test.go` (added tests)
- NEW: `internal/models/pull_request.go` (placeholder)
- NEW: `internal/models/pull_request_test.go` (placeholder)

**Acceptance Criteria**:
- [x] PullRequest struct with all fields (ID, CommitIDs, ReviewIDs, IssueIDs added)
- [x] PRStatus enum (open, merged, closed) (used existing PRState)
- [x] JSON marshaling works (tests verify marshaling/unmarshaling)
- [x] Validation for required fields (Validate() method added)
- [x] Tests pass (54 tests passing)

**Implementation Notes**:
- Enhanced existing PullRequest model in pr.go instead of creating duplicate
- Added ID field for internal tracking
- Added relationship arrays: CommitIDs, ReviewIDs, IssueIDs with omitempty tags
- Validate() method checks title, author_email, branches, state, created_at
- 4 new test functions cover all new functionality

---

#### TASK-GH-02: Create Review Model (0.75h)

**Goal**: Define Review struct with comments

**Files**:
- NEW: `internal/models/review.go`
- NEW: `internal/models/review_test.go`

**Acceptance Criteria**:
- [ ] Review struct with state
- [ ] ReviewComment struct
- [ ] ReviewState enum (approved, changes_requested, commented)
- [ ] JSON marshaling works
- [ ] Tests pass

---

#### TASK-GH-03: Create Issue Model (0.75h)

**Goal**: Define Issue struct with PR linkage

**Files**:
- NEW: `internal/models/issue.go`
- NEW: `internal/models/issue_test.go`

**Acceptance Criteria**:
- [ ] Issue struct with state and labels
- [ ] IssueState enum (open, in_progress, closed)
- [ ] PR linkage field (ClosedByPRID)
- [ ] JSON marshaling works
- [ ] Tests pass

---

### PHASE 2: Generators

#### TASK-GH-04: Implement PR Generator (3.0h)

**Goal**: Generate PRs from commits

**Files**:
- NEW: `internal/generator/pr_generator.go`
- NEW: `internal/generator/pr_generator_test.go`

**Acceptance Criteria**:
- [ ] Groups commits into PRs (3-10 commits per PR)
- [ ] Generates branch names (feature/*, bugfix/*)
- [ ] Assigns PR status (85% merged, 10% closed, 5% open)
- [ ] Sets merge time (1-7 days after creation)
- [ ] Respects temporal ordering
- [ ] Tests with 20+ scenarios
- [ ] Integration with commit generator

---

#### TASK-GH-05: Implement Review Generator (3.0h)

**Goal**: Generate reviews for PRs

**Files**:
- NEW: `internal/generator/review_generator.go`
- NEW: `internal/generator/review_generator_test.go`

**Acceptance Criteria**:
- [ ] Assigns 1-3 reviewers per PR
- [ ] Reviewer ≠ PR author
- [ ] Review state distribution (70% approved, 20% changes, 10% commented)
- [ ] Review timestamp between PR creation and merge
- [ ] Generates 0-5 comments per review
- [ ] Tests cover all states
- [ ] Edge cases handled

---

#### TASK-GH-06: Implement Issue Generator (2.0h)

**Goal**: Generate issues linked to PRs

**Files**:
- NEW: `internal/generator/issue_generator.go`
- NEW: `internal/generator/issue_generator_test.go`

**Acceptance Criteria**:
- [ ] 40% of merged PRs close an issue
- [ ] 10% of issues remain open
- [ ] Issue created before PR
- [ ] Labels assigned (bug, feature, enhancement)
- [ ] Issue titles match PR titles
- [ ] Tests pass

---

### PHASE 3: Storage

#### TASK-GH-07: Add Storage Methods for GitHub Data (2.0h)

**Goal**: Extend Store interface for PRs, reviews, issues

**Files**:
- MODIFY: `internal/storage/store.go` (interface)
- MODIFY: `internal/storage/memory_store.go` (implementation)
- NEW: `internal/storage/github_storage_test.go`

**Acceptance Criteria**:
- [ ] StorePR, GetPRByID, GetPRsByStatus, GetPRsByAuthor
- [ ] StoreReview, GetReviewsByPRID, GetReviewsByReviewer
- [ ] StoreIssue, GetIssueByNumber, GetIssuesByState
- [ ] Thread-safe operations
- [ ] Tests for all methods

---

#### TASK-GH-08: Integrate Generators with Storage (1.0h)

**Goal**: Wire up generators to persist data

**Files**:
- MODIFY: `internal/generator/pr_generator.go`
- MODIFY: `internal/generator/review_generator.go`
- MODIFY: `internal/generator/issue_generator.go`
- NEW: `test/integration/github_generation_test.go`

**Acceptance Criteria**:
- [ ] PRs stored after generation
- [ ] Reviews stored after generation
- [ ] Issues stored after generation
- [ ] Integration test verifies end-to-end flow
- [ ] Tests pass

---

### PHASE 4: API Handlers

#### TASK-GH-09: Implement `/analytics/github/prs` Endpoint (1.0h)

**Goal**: PR listing with filters

**Files**:
- NEW: `internal/api/github/prs.go`
- NEW: `internal/api/github/prs_test.go`

**Acceptance Criteria**:
- [ ] Query params: status, author, start_date, end_date
- [ ] Pagination support
- [ ] Returns PR list with metrics
- [ ] Handler tests pass

---

#### TASK-GH-10: Implement `/analytics/github/reviews` Endpoint (1.0h)

**Goal**: Review activity listing

**Files**:
- NEW: `internal/api/github/reviews.go`
- NEW: `internal/api/github/reviews_test.go`

**Acceptance Criteria**:
- [ ] Query params: pr_id, reviewer
- [ ] Returns review list
- [ ] Handler tests pass

---

#### TASK-GH-11: Implement `/analytics/github/issues` Endpoint (1.0h)

**Goal**: Issue tracking data

**Files**:
- NEW: `internal/api/github/issues.go`
- NEW: `internal/api/github/issues_test.go`

**Acceptance Criteria**:
- [ ] Query params: state, labels
- [ ] Returns issue list
- [ ] Handler tests pass

---

#### TASK-GH-12: Implement `/analytics/github/pr-cycle-time` Endpoint (1.0h)

**Goal**: PR lifecycle metrics

**Files**:
- NEW: `internal/api/github/pr_cycle_time.go`
- NEW: `internal/api/github/pr_cycle_time_test.go`

**Acceptance Criteria**:
- [ ] Calculates avg time to first review
- [ ] Calculates avg/median time to merge
- [ ] Returns percentiles (p50, p75, p90)
- [ ] Handler tests pass

---

#### TASK-GH-13: Implement `/analytics/github/review-quality` Endpoint (1.0h)

**Goal**: Review quality metrics

**Files**:
- NEW: `internal/api/github/review_quality.go`
- NEW: `internal/api/github/review_quality_test.go`

**Acceptance Criteria**:
- [ ] Calculates approval rate
- [ ] Avg reviewers per PR
- [ ] Avg comments per PR
- [ ] Handler tests pass

---

### PHASE 5: Testing & Documentation

#### TASK-GH-14: E2E Tests for GitHub Simulation (2.0h)

**Goal**: Comprehensive E2E tests

**Files**:
- NEW: `test/e2e/github_test.go`

**Acceptance Criteria**:
- [ ] Test full PR lifecycle (create → review → merge)
- [ ] Test issue resolution via PR
- [ ] Test all 5 API endpoints
- [ ] Test pagination and filters
- [ ] All tests pass

---

#### TASK-GH-15: Update Documentation (1.5h)

**Goal**: Document GitHub simulation in SPEC.md

**Files**:
- MODIFY: `services/cursor-sim/SPEC.md`
- MODIFY: `.claude/DEVELOPMENT.md`

**Acceptance Criteria**:
- [ ] SPEC.md documents all 5 endpoints
- [ ] SPEC.md shows response schemas
- [ ] SPEC.md explains PR/review/issue generation
- [ ] DEVELOPMENT.md marks P2-F01 complete
- [ ] Manual testing checklist included

---

## Dependencies

- **P1 (Complete)**: Commit generation
- **P3 (Complete)**: Storage layer

---

## Success Criteria

- [ ] All 15 tasks complete
- [ ] 90%+ test coverage
- [ ] All API endpoints functional
- [ ] SPEC.md updated
- [ ] E2E tests passing
- [ ] Build successful
- [ ] Manual testing verified

---

**Estimated Total**: 22.0 hours
