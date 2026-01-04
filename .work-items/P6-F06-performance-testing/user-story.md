# User Story: Performance Testing with Lighthouse

**Feature ID**: P6-F06
**Epic**: P6 - cursor-viz-spa (Testing Enhancement)
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: LOW

---

## User Story

**As a** product owner,
**I want** automated performance testing that verifies the dashboard loads quickly,
**So that** performance regressions are caught before they impact users.

---

## Background

As the dashboard grows in complexity with charts and data visualizations, performance can degrade without notice. Lighthouse CI provides:
- Core Web Vitals measurement
- Performance budgets
- Regression detection
- Accessibility audits (bonus)

---

## Acceptance Criteria

### AC-01: Lighthouse CI Setup
- **Given** cursor-viz-spa project
- **When** developer runs `npm run lighthouse`
- **Then** Lighthouse audit runs and reports performance scores

### AC-02: Performance Budget
- **Given** performance budget is defined
- **When** Dashboard loads
- **Then** First Contentful Paint < 2s, Time to Interactive < 5s

### AC-03: CI Integration
- **Given** PR with P6 changes
- **When** CI pipeline runs
- **Then** Lighthouse audit runs and reports in PR

### AC-04: Performance Regression Detection
- **Given** performance baseline exists
- **When** new code causes regression
- **Then** CI warns about performance decrease

---

## Out of Scope

- Backend performance testing (P5)
- Load testing (concurrent users)
- Mobile performance (initial focus on desktop)

---

## Dependencies

| Dependency | Description | Status |
|------------|-------------|--------|
| @lhci/cli | Lighthouse CI CLI | To install |
| P6 built | Must test production build | Available |

---

## References

- **E2E Testing Strategy**: `docs/e2e-testing-strategy.md` (Phase 4)
- **Lighthouse CI**: https://github.com/GoogleChrome/lighthouse-ci
