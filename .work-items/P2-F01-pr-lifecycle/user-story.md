# User Stories: cursor-sim PR Lifecycle

**Feature ID**: P2-F01-pr-lifecycle
**Phase**: P2 (cursor-sim GitHub Simulation)
**Created**: January 2, 2026
**Status**: TODO

**Feature**: cursor-sim v2 Phase 2
**Priority**: P1
**Status**: NOT_STARTED

---

## Epic: GitHub PR Simulation

### US-SIM-P2-001: Generate PRs from Commits

**As a** SDLC researcher
**I want** cursor-sim to generate realistic PR data from commits
**So that** I can study the relationship between AI-assisted coding and PR outcomes

**Priority**: P1
**Story Points**: 8
**Feature**: SIM-R009

**Acceptance Criteria:**

**Scenario 1: PR generation from commit clusters**
Given commits exist in the storage
When PRs are generated
Then each PR should reference 1-10 commits based on developer patterns

**Scenario 2: PR metadata completeness**
Given a PR is generated
Then it should include: title, description, author, reviewers, labels, created_at, merged_at (nullable)

**Scenario 3: Commit-PR linkage**
Given commits are assigned to a PR
Then the commit's pr_number field should reference the PR

---

### US-SIM-P2-002: Simulate Code Reviews

**As a** analytics dashboard user
**I want** realistic review cycles simulated
**So that** I can analyze review patterns and time-to-merge metrics

**Priority**: P1
**Story Points**: 8
**Feature**: SIM-R010

**Acceptance Criteria:**

**Scenario 1: Review assignment**
Given a PR is created
Then 1-3 reviewers should be assigned based on team membership

**Scenario 2: Review comments**
Given reviewers are assigned
Then each reviewer may generate 0-5 comments based on thoroughness setting

**Scenario 3: Approval flow**
Given a PR has been reviewed
Then it should transition through: opened -> review -> changes_requested (optional) -> approved -> merged

---

### US-SIM-P2-003: GitHub API Compatibility

**As a** developer integrating with cursor-sim
**I want** GitHub-compatible API endpoints
**So that** I can use standard GitHub client libraries

**Priority**: P1
**Story Points**: 5
**Feature**: SIM-R011

**Acceptance Criteria:**

**Scenario 1: List repositories**
Given the server is running
When I call GET /repos
Then I receive a paginated list of repositories matching GitHub schema

**Scenario 2: List PRs**
Given a repository has PRs
When I call GET /repos/{owner}/{repo}/pulls
Then I receive PRs matching GitHub's PR response schema

**Scenario 3: PR details with commits**
Given a PR exists
When I call GET /repos/{owner}/{repo}/pulls/{number}/commits
Then I receive the commits linked to that PR

---

### US-SIM-P2-004: Track Quality Outcomes

**As a** researcher studying AI code quality
**I want** quality outcome signals in the data
**So that** I can correlate AI assistance levels with code stability

**Priority**: P1
**Story Points**: 5
**Feature**: SIM-R012

**Acceptance Criteria:**

**Scenario 1: Revert tracking**
Given a commit is generated
Then some commits (configurable %) should be marked as reverts

**Scenario 2: AI ratio correlation**
Given commits have varying AI ratios
Then higher AI ratios should correlate with configurable revert rates

**Scenario 3: Bug-fix identification**
Given commits are generated
Then commits with "fix:" prefix should be marked as bug fixes

---

## Dependencies

- **Requires**: cursor-sim Phase 1 (COMPLETE)
- **Enables**: cursor-analytics-core PR metrics, cursor-viz-spa PR dashboards
