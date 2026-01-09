# Task Breakdown: GitHub Simulation

**Feature ID**: P2-F01-github-simulation
**Status**: Not Started
**Created**: January 8, 2026

---

## Progress Tracker

| Phase | Tasks | Status | Estimated | Actual |
|-------|-------|--------|-----------|--------|
| **Models** | 3 | ✅ 3/3 | 2.5h | 1.25h |
| **Generators** | 3 | ✅ 3/3 | 8.0h | 3.0h |
| **Storage** | 2 | ✅ 2/2 | 3.0h | 2.0h |
| **API Handlers** | 5 | ✅ 5/5 | 5.0h | 5.0h |
| **Testing & Docs** | 2 | ⬜ 0/2 | 3.5h | - |
| **TOTAL** | **15** | **13/15** | **22.0h** | **11.25h** |

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

#### TASK-GH-02: Create Review Model (0.75h) ✅ COMPLETE

**Status**: COMPLETE
**Time**: 0.5h actual / 0.75h estimated
**Commit**: bfa3f64

**Goal**: Define Review struct with comments

**Files**:
- NEW: `internal/models/review.go`
- NEW: `internal/models/review_test.go`

**Acceptance Criteria**:
- [x] Review struct with state
- [x] ReviewComment struct (already in pr.go)
- [x] ReviewState enum (approved, changes_requested, pending)
- [x] JSON marshaling works
- [x] Tests pass (6 test functions, 19 subtests)

**Implementation Notes**:
- Review model with ID, PRID, Reviewer, State, SubmittedAt, Body, Comments
- Validate() method checks pr_id, reviewer, state, submitted_at
- IsApproval() and CommentCount() helper methods
- 93.9% coverage for internal/models package
- ReviewState and ReviewComment already existed in pr.go, used those definitions

---

#### TASK-GH-03: Create Issue Model (0.75h) ✅ COMPLETE

**Status**: COMPLETE
**Time**: 0.25h actual / 0.75h estimated

**Goal**: Define Issue struct with PR linkage

**Files**:
- NEW: `internal/models/issue.go`
- NEW: `internal/models/issue_test.go`

**Acceptance Criteria**:
- [x] Issue struct with state and labels
- [x] IssueState enum (open, closed)
- [x] PR linkage field (ClosedByPRID)
- [x] JSON marshaling works
- [x] Tests pass (8 test functions, 22 subtests)

**Implementation Notes**:
- Issue model with Number, Title, Body, State, AuthorID, RepoName, Labels, Assignees, timestamps
- IssueState enum: open, closed (simpler than originally planned)
- ClosedByPRID pointer for tracking PR that closed the issue
- Validate() method checks number, title, state, author_id, repo_name, created_at
- Helper methods: IsOpen(), IsClosed(), WasClosedByPR()
- Comprehensive tests for validation, JSON marshaling, helper methods

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

#### TASK-GH-05: Implement Review Generator (3.0h) ✅ COMPLETE

**Status**: COMPLETE
**Time**: 1.5h actual / 3.0h estimated

**Goal**: Generate reviews for PRs

**Files**:
- MODIFIED: `internal/generator/review_generator.go` (added GenerateReviewsForPR method)
- MODIFIED: `internal/generator/review_generator_test.go` (added 12 new test functions)

**Acceptance Criteria**:
- [x] Assigns 1-3 reviewers per PR
- [x] Reviewer ≠ PR author (enforced by filtering)
- [x] Review state distribution (70% approved, 20% changes_requested, 10% pending)
- [x] Review timestamp between PR creation and merge/close
- [x] Generates 0-5 comments per non-approved review
- [x] Tests cover all states (19 test functions total)
- [x] Edge cases handled (no reviewers, unique reviewers, reproducibility)

**Implementation Notes**:
- Added `GenerateReviewsForPR(pr PullRequest) []Review` method
- New constructors: `NewReviewGenerator(seedData, rng)` for simple use
- Review timing uses PR.MergedAt, PR.ClosedAt, or 7-day window for open PRs
- State distribution: 70% approved, 20% changes_requested, 10% pending
- Non-approved reviews get body text and 0-5 inline comments
- Approved reviews get short positive messages (LGTM!)
- Auto-incrementing review IDs for uniqueness
- Full backward compatibility with existing store-based methods

---

#### TASK-GH-06: Implement Issue Generator (2.0h) ✅ COMPLETE

**Status**: COMPLETE
**Time**: 1.0h actual / 2.0h estimated
**Completed**: 2026-01-09

**Goal**: Generate issues linked to PRs

**Files**:
- NEW: `internal/generator/issue_generator.go`
- NEW: `internal/generator/issue_generator_test.go`

**Acceptance Criteria**:
- [x] 40% of merged PRs close an issue
- [x] 10% of issues remain open
- [x] Issue created before PR (1-7 days)
- [x] Labels assigned (bug, feature, enhancement)
- [x] Issue titles derived from PR titles
- [x] Tests pass (12 test functions)

**Implementation Notes**:
- IssueGenerator struct with seed data and RNG
- GenerateIssuesForPRs(prs, repoName) method
- 40% of merged PRs get linked issues
- 10% of generated issues remain open
- Issues created 1-7 days before PR creation
- Labels randomly assigned (1-2 from bug/feature/enhancement)
- Reproducible with seeded RNG
- Full validation passes on generated issues

---

### PHASE 3: Storage

#### TASK-GH-07: Add Storage Methods for GitHub Data (2.0h) ✅ COMPLETE

**Status**: COMPLETE
**Time**: 1.5h actual / 2.0h estimated
**Completed**: 2026-01-09

**Goal**: Extend Store interface for PRs, reviews, issues

**Files**:
- MODIFIED: `internal/storage/store.go` (added PR, Review, Issue interfaces)
- MODIFIED: `internal/storage/memory.go` (added storage fields and implementations)
- NEW: `internal/storage/github_storage_test.go` (14 comprehensive test functions)

**Acceptance Criteria**:
- [x] StorePR, GetPRByID, GetPRsByStatus, GetPRsByAuthorEmail
- [x] GetPRsByRepoWithPagination with state filtering
- [x] StoreReview, GetReviewsByPRID, GetReviewsByReviewer, GetReviewsByRepoPR
- [x] StoreIssue, GetIssueByNumber, GetIssuesByState, GetIssuesByRepo
- [x] Thread-safe operations (sync.RWMutex)
- [x] Tests for all methods (100% pass rate)

**Implementation Notes**:
- Added 5 new PR storage methods with ID-based indexing
- Added 4 new Review storage methods with PRID and reviewer indexing
- Added 4 new Issue storage methods with repo/number and state indexing
- Pagination support for PR queries with state filtering
- Multiple indexes per entity type for efficient lookups
- Thread-safety verified with concurrent access tests
- All 73 storage tests pass including 14 new GitHub storage tests

---

#### TASK-GH-08: Integrate Generators with Storage (1.0h)

**Status**: COMPLETE
**Actual**: 0.5h
**Completed**: 2026-01-09

**Goal**: Wire up generators to persist data

**Files**:
- MODIFY: `internal/generator/issue_generator.go` - Added IssueStore interface, NewIssueGeneratorWithStore(), GenerateAndStoreIssuesForPRs()

**Notes**:
- PRGenerator already has storage integration via GeneratePRsFromCommits() (implemented in earlier task)
- ReviewGenerator already has storage integration (implemented in earlier task)
- IssueGenerator now has storage capability via new methods

**Acceptance Criteria**:
- [x] PRs stored after generation (pre-existing)
- [x] Reviews stored after generation (pre-existing)
- [x] Issues stored after generation (new)
- [ ] Integration test verifies end-to-end flow (deferred to E2E phase)
- [x] Tests pass (all generator tests pass)

---

### PHASE 4: API Handlers

#### TASK-GH-09: Implement `/analytics/github/prs` Endpoint (1.0h) ✅ COMPLETE

**Status**: COMPLETE
**Time**: 0.5h actual / 1.0h estimated
**Completed**: 2026-01-09

**Goal**: PR listing with filters

**Files**:
- NEW: `internal/api/github/prs.go`
- NEW: `internal/api/github/prs_test.go`
- MODIFIED: `internal/server/router.go`

**Acceptance Criteria**:
- [x] Query params: status, author, start_date, end_date
- [x] Pagination support
- [x] Returns PR list with metrics
- [x] Handler tests pass (12 test cases)

---

#### TASK-GH-10: Implement `/analytics/github/reviews` Endpoint (1.0h)

**Status**: COMPLETE
**Time**: 1.0h actual / 1.0h estimated
**Completed**: 2026-01-09

**Goal**: Review activity listing

**Files**:
- NEW: `internal/api/github/reviews.go`
- NEW: `internal/api/github/reviews_test.go`
- MODIFIED: `internal/server/router.go` (registered endpoint)
- MODIFIED: `SPEC.md` (added GitHub Analytics API section)

**Acceptance Criteria**:
- [x] Query params: pr_id, reviewer
- [x] Pagination support (page, page_size)
- [x] Returns review list with filtering
- [x] Handler tests pass (11 test cases)

**Implementation Notes**:
- Follows same pattern as PRs endpoint
- Uses existing storage methods: GetReviewsByPRID, GetReviewsByReviewer
- Response format: `{ data: [...], pagination: {...}, params: {...} }`
- Comprehensive test coverage: all filters, pagination, combined filters, error cases

---

#### TASK-GH-11: Implement `/analytics/github/issues` Endpoint (1.0h) ✅ COMPLETE

**Status**: COMPLETE
**Time**: 1.0h actual / 1.0h estimated
**Commit**: d246bef

**Goal**: Issue tracking data

**Files**:
- NEW: `internal/api/github/issues.go` (167 lines)
- NEW: `internal/api/github/issues_test.go` (346 lines, 11 test cases)
- MODIFIED: `internal/server/router.go` (registered endpoint)
- MODIFIED: `services/cursor-sim/SPEC.md` (updated API documentation)

**Acceptance Criteria**:
- [x] Query params: state, labels
- [x] Returns issue list
- [x] Handler tests pass (11 tests, 97% coverage)

**Implementation Notes**:
- Supports filtering by state (open, closed)
- Supports filtering by labels (comma-separated, AND logic)
- Pagination with page and page_size parameters
- Standard response format: data, pagination, params
- Comprehensive test coverage including error cases
- All tests passing (TestListIssuesAnalytics_*)

---

#### TASK-GH-12: Implement `/analytics/github/pr-cycle-time` Endpoint (1.0h) ✅ COMPLETE

**Status**: COMPLETE
**Time**: 1.0h actual / 1.0h estimated
**Completed**: 2026-01-09

**Goal**: PR lifecycle metrics

**Files**:
- NEW: `internal/api/github/pr_cycle_time.go` (205 lines)
- NEW: `internal/api/github/pr_cycle_time_test.go` (336 lines, 8 test cases)
- MODIFIED: `internal/server/router.go` (registered endpoint)
- MODIFIED: `services/cursor-sim/SPEC.md` (updated API documentation)

**Acceptance Criteria**:
- [x] Calculates avg time to first review
- [x] Calculates avg/median time to merge
- [x] Returns percentiles (p50, p75, p90)
- [x] Handler tests pass (8 tests, 100% pass rate)

**Implementation Notes**:
- Endpoint registered at `/analytics/github/pr-cycle-time`
- Calculates metrics from merged PRs only
- Supports date range filtering with `from` and `to` query params (YYYY-MM-DD format)
- Returns metrics in seconds for all time values
- Handles edge cases: no data, no reviews, no merged PRs, single PR
- Percentile calculation uses linear interpolation for accuracy
- Response format: `{ data: {...}, params: {...} }`
- All 8 tests passing with comprehensive coverage

---

#### TASK-GH-13: Implement `/analytics/github/review-quality` Endpoint (1.0h) ✅ COMPLETE

**Status**: COMPLETE
**Time**: 1.5h actual / 1.0h estimated
**Completed**: 2026-01-09
**Commit**: 868bfea (accidentally included with TASK-DS-12)

**Goal**: Review quality metrics

**Files**:
- NEW: `internal/api/github/review_quality.go` (205 lines)
- NEW: `internal/api/github/review_quality_test.go` (296 lines, 5 test cases)
- MODIFIED: `internal/server/router.go` (registered endpoint)
- MODIFIED: `services/cursor-sim/SPEC.md` (documented endpoint, query params, response format)

**Acceptance Criteria**:
- [x] Calculates approval rate (% of reviews that are approvals)
- [x] Avg reviewers per PR (unique reviewers per merged PR)
- [x] Avg comments per review (comment count per review)
- [x] Calculates changes requested rate and pending rate
- [x] Supports date range filtering (from/to query parameters)
- [x] Only analyzes merged PRs within date range
- [x] Handler tests pass (5 tests, 100% pass rate)

**Implementation Notes**:
- Endpoint registered at `/analytics/github/review-quality`
- Calculates metrics from reviews on merged PRs only
- Supports date range filtering with `from` and `to` query params (YYYY-MM-DD format)
- Metrics include: approval_rate, avg_reviewers_per_pr, avg_comments_per_review, changes_requested_rate, pending_rate
- Tracks unique reviewers per PR using map data structure
- Response format: `{ data: {...}, params: {...} }`
- All 5 tests passing: basic metrics, date filtering, no data, invalid dates, merged PRs only
- Comprehensive test coverage including edge cases

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
