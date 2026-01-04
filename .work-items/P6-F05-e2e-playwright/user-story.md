# User Story: E2E Testing with Playwright

**Feature ID**: P6-F05
**Epic**: P6 - cursor-viz-spa (Testing Enhancement)
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: MEDIUM

---

## User Story

**As a** QA engineer or developer,
**I want** automated E2E tests that verify the full P4→P5→P6 data flow,
**So that** integration issues are caught before production deployment.

---

## Background

Current testing covers:
- Unit tests (91.68% coverage)
- Component tests (with mocks)
- Integration tests (with MockedProvider)

Missing:
- Full stack E2E tests (real P5, real data)
- Visual regression tests
- Browser-based testing

E2E tests verify the actual user experience by testing through a real browser.

---

## Acceptance Criteria

### AC-01: Playwright Setup
- **Given** cursor-viz-spa project
- **When** developer runs `npm install`
- **Then** Playwright is installed with required browsers

### AC-02: Dashboard E2E Test
- **Given** P4 and P5 are running
- **When** E2E tests execute
- **Then** Dashboard loads with real data from P5

### AC-03: Navigation E2E Test
- **Given** Dashboard is loaded
- **When** user navigates to Teams and Developers pages
- **Then** navigation works and pages load correctly

### AC-04: Visual Regression Test
- **Given** baseline screenshots exist
- **When** tests run
- **Then** current screenshots are compared to baselines

### AC-05: Error Handling E2E Test
- **Given** P5 is not available
- **When** Dashboard loads
- **Then** error message is displayed (not crash)

### AC-06: CI Integration
- **Given** PR with P6 changes
- **When** CI pipeline runs
- **Then** E2E tests run with Playwright

---

## Out of Scope

- Performance testing (handled by P6-F06)
- API testing (covered by P5 tests)
- Unit/component tests (already exist)

---

## Dependencies

| Dependency | Description | Status |
|------------|-------------|--------|
| @playwright/test | E2E testing framework | To install |
| P4 cursor-sim | Must be running for E2E | Available (Docker) |
| P5 cursor-analytics-core | Must be running for E2E | Available (Docker) |

---

## References

- **E2E Testing Strategy**: `docs/e2e-testing-strategy.md` (Phase 2)
- **Playwright Docs**: https://playwright.dev
