# User Story: Schema Validation Tests

**Feature ID**: P5-F03
**Epic**: P5 - cursor-analytics-core (Testing Enhancement)
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: MEDIUM

---

## User Story

**As a** backend developer working on cursor-analytics-core,
**I want** automated tests that validate the GraphQL schema structure,
**So that** schema changes are caught by tests before they break consumers.

---

## Background

The P5+P6 integration failure (January 4, 2026) revealed that P5 schema changes weren't validated against expected structure. Schema validation tests ensure:
- Required types exist
- Required fields exist on types
- Field types match expected (singular vs array, nullable vs required)
- Queries execute without errors

---

## Acceptance Criteria

### AC-01: Type Existence Tests
- **Given** GraphQL schema is defined
- **When** schema validation tests run
- **Then** all required types (Developer, TeamStats, DashboardKPI, etc.) exist

### AC-02: Field Structure Tests
- **Given** TeamStats type exists
- **When** field validation tests run
- **Then** topPerformer field is singular Developer (not array)

### AC-03: Query Validation Tests
- **Given** valid query syntax
- **When** query is validated against schema
- **Then** validation passes without errors

### AC-04: Breaking Change Detection
- **Given** developer removes a field
- **When** tests run
- **Then** test fails indicating missing field

---

## Out of Scope

- Resolver logic testing (covered by existing tests)
- Database integration tests (separate feature)
- Apollo Studio integration (P5-F02)

---

## Dependencies

| Dependency | Description | Status |
|------------|-------------|--------|
| GraphQL schema | Must be defined | Exists |
| graphql package | Schema introspection | Installed |

---

## References

- **E2E Testing Strategy**: `docs/e2e-testing-strategy.md` (Phase 1.2)
- **Current Test Coverage**: 91.49% (unit tests)
