# User Story: Component Integration Tests

**Feature ID**: P6-F04
**Epic**: P6 - cursor-viz-spa (Testing Enhancement)
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: MEDIUM

---

## User Story

**As a** frontend developer working on cursor-viz-spa,
**I want** integration tests that verify components work correctly with their hooks and child components,
**So that** component integration issues are caught before deployment.

---

## Background

Current P6 testing focuses on unit tests (91.68% coverage) but lacks integration tests that verify:
- Pages correctly integrate with data hooks
- Props are passed correctly between components
- Loading/error/success states render properly
- GraphQL queries return expected data structures

The P5+P6 integration issues (January 4, 2026) revealed that unit tests alone don't catch integration problems.

---

## Acceptance Criteria

### AC-01: Dashboard Integration Test
- **Given** Dashboard page with mock GraphQL data
- **When** page renders
- **Then** KPI cards, charts, and tables display correct data

### AC-02: Teams Integration Test
- **Given** Teams page with mock GraphQL data
- **When** page renders
- **Then** team list displays with correct team data

### AC-03: Developers Integration Test
- **Given** Developers page with mock GraphQL data
- **When** page renders
- **Then** developer list displays with pagination working

### AC-04: Error State Testing
- **Given** GraphQL query returns error
- **When** page renders
- **Then** error message displays (not crash)

### AC-05: Loading State Testing
- **Given** GraphQL query is pending
- **When** page renders
- **Then** loading indicator displays

### AC-06: Data Flow Validation
- **Given** mock data matches P5 schema
- **When** tests run
- **Then** data flows correctly from hooks to components to UI

---

## Out of Scope

- E2E tests with real P5 (handled by P6-F05)
- Visual regression tests (handled by P6-F05)
- Unit tests (already exist)

---

## Dependencies

| Dependency | Description | Status |
|------------|-------------|--------|
| @apollo/client/testing | MockedProvider | Installed |
| @testing-library/react | Component testing | Installed |
| P6-F02 | Generated types for mock data | Must complete first |

---

## References

- **E2E Testing Strategy**: `docs/e2e-testing-strategy.md` (Phase 1.1)
- **Current Test Coverage**: 91.68% (unit tests)
