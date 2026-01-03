# Task Breakdown: Quality Analysis

**Feature ID**: P3-F03-quality-analysis
**Phase**: P3 (cursor-sim Research Framework)
**Created**: January 3, 2026
**Status**: COMPLETE

## Overview

**Total Estimated Hours**: 18.5h
**Total Actual Hours**: 17.5h
**Number of Tasks**: 7

## Progress Tracker

| Task ID | Task | Hours Est | Status | Actual |
|---------|------|-----------|--------|--------|
| TASK00 | PR Generation Pipeline | 4.0 | DONE | 4.0 |
| TASK01 | Wire GitHub/Research Routes | 2.0 | DONE | 2.0 |
| TASK02 | Code Survival Calculator | 3.0 | DONE | 3.0 |
| TASK03 | Revert Chain Analysis | 2.5 | DONE | 2.5 |
| TASK04 | Hotfix Tracking | 2.0 | DONE | 2.0 |
| TASK05 | Research Dataset Enhancement | 2.5 | DONE | 2.5 |
| TASK06 | Integration Tests | 2.5 | DONE | 1.5 |

---

## Task Details

### TASK00: PR Generation Pipeline

**Files**:
- `internal/generator/pr_generator.go`
- `internal/generator/session.go`

Generate PRs from commit groupings using session-based parameters.
- Session model with seniority-based MaxCommits
- Inactivity gap detection for PR boundaries
- PR envelope with aggregated metrics

### TASK01: Wire GitHub/Research Routes

**Files**: `internal/server/router.go`, `internal/api/github/*.go`

Wire all 12 GitHub routes and 5 Research routes.
Implement ListCommits, GetCommit, ListPullCommits, ListPullFiles handlers.

### TASK02: Code Survival Calculator

**Files**:
- `internal/services/survival.go`
- `internal/models/quality.go`
- `internal/api/github/analysis.go`

Track file-level survival with cohort analysis.
Calculate survival_rate = files_surviving / files_added_in_cohort.

### TASK03: Revert Chain Analysis

**Files**: `internal/services/reverts.go`

Detect reverts via message patterns and link to original PRs.
Sigmoid risk function: AI ratio + seniority + volatility.

### TASK04: Hotfix Tracking

**Files**: `internal/services/hotfixes.go`

Detect fix-PRs within 48 hours of merged PR to same files.
Pattern matching: fix, hotfix, urgent, patch.

### TASK05: Research Dataset Enhancement

**Files**: `internal/models/research.go`, `internal/generator/research_generator.go`

Add all 38 columns including:
- greenfield_index, pickup_time, review_density
- rework_ratio, scope_creep, reviewer_count
- has_hotfix_followup, survival_rate_30d

### TASK06: Integration Tests

**Files**: `test/e2e/github_api_test.go`, `test/e2e/quality_analysis_test.go`

E2E tests for GitHub and Research endpoints.
Hypothesis validation tests for AI impact correlations.

---

**Status**: COMPLETE
