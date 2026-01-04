# User Story: Contract Testing with GraphQL Inspector

**Feature ID**: P6-F03
**Epic**: P6 - cursor-viz-spa (Service Decoupling Phase 3)
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: MEDIUM

---

## User Story

**As a** frontend developer working on cursor-viz-spa,
**I want** my GraphQL queries validated against the P5 schema before commit,
**So that** incompatible queries are caught before they reach production.

---

## Background

With GraphQL Code Generator (P6-F02) generating types, TypeScript catches type mismatches at compile time. However, it doesn't validate that:
- Query field selections are valid
- Required variables are provided
- Fragments match type conditions

GraphQL Inspector provides query validation against a live schema, catching errors that TypeScript cannot detect.

---

## Acceptance Criteria

### AC-01: GraphQL Inspector Installation
- **Given** cursor-viz-spa project
- **When** developer runs `npm install`
- **Then** GraphQL Inspector CLI is available

### AC-02: Query Validation Script
- **Given** P5 GraphQL server is running
- **When** developer runs `npm run schema:validate`
- **Then** all P6 queries are validated against P5 schema

### AC-03: Pre-Commit Hook
- **Given** developer has modified a GraphQL query
- **When** developer attempts to commit
- **Then** commit is blocked if queries are invalid

### AC-04: CI Integration
- **Given** PR with P6 changes
- **When** CI pipeline runs
- **Then** query validation runs and fails on invalid queries

### AC-05: Clear Error Messages
- **Given** a query has an invalid field
- **When** validation fails
- **Then** error message shows exact field and file location

---

## Out of Scope

- Schema publishing (handled by P5-F02)
- Type generation (handled by P6-F02)
- E2E testing (handled by P6-F05)

---

## Dependencies

| Dependency | Description | Status |
|------------|-------------|--------|
| P6-F02 | GraphQL Code Generator | Must complete first |
| P5 running | For schema introspection | Available |
| @graphql-inspector/cli | Validation tool | To install |

---

## References

- **Mitigation Plan**: `docs/MITIGATION-PLAN.md` (Strategy 5)
- **Data Contract Testing**: `docs/data-contract-testing.md` (Phase 3)
- **GraphQL Inspector**: https://graphql-inspector.com
