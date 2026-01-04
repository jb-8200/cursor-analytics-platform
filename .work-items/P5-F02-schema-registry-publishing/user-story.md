# User Story: Schema Registry Publishing

**Feature ID**: P5-F02
**Epic**: P5 - cursor-analytics-core (Service Decoupling Phase 2)
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: MEDIUM

---

## User Story

**As a** backend developer working on cursor-analytics-core,
**I want** the GraphQL schema to be automatically published to a central registry,
**So that** schema changes are versioned, breaking changes are detected, and consumers are notified.

---

## Background

After implementing GraphQL Code Generator (P6-F02), P6 can auto-generate types from P5's schema. However, this requires P5 to be running locally. A schema registry provides:
- Centralized schema storage (no local P5 required for codegen)
- Schema version history
- Breaking change detection before deployment
- Documentation auto-generation

---

## Acceptance Criteria

### AC-01: Schema Registry Account
- **Given** Apollo Studio account is set up
- **When** developer logs in
- **Then** cursor-analytics graph is visible

### AC-02: Schema Publishing
- **Given** P5 has schema changes
- **When** changes are pushed to main branch
- **Then** CI automatically publishes schema to Apollo Studio

### AC-03: Breaking Change Detection
- **Given** P5 schema has breaking change (field removal, type change)
- **When** CI runs schema check
- **Then** warning is displayed with list of breaking changes

### AC-04: P6 Schema Fetching
- **Given** Apollo Studio has latest P5 schema
- **When** P6 runs codegen
- **Then** types are generated from registry (no local P5 needed)

### AC-05: Version History
- **Given** multiple schema versions published
- **When** developer views Apollo Studio
- **Then** full version history is visible with diffs

---

## Out of Scope

- Schema governance policies
- Team notifications (Slack/email)
- Production monitoring

---

## Dependencies

| Dependency | Description | Status |
|------------|-------------|--------|
| P6-F02 | GraphQL Code Generator | Must complete first |
| Apollo Studio | Free for OSS | Account needed |
| @apollo/rover | Schema publishing CLI | To install |

---

## Technical Context

### Schema Publishing Flow

```
┌─────────────────────────────────────────┐
│  Developer pushes P5 changes            │
│  ─────────────────────────────────────  │
│  git push origin main                   │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│  CI: Publish Schema                     │
│  ─────────────────────────────────────  │
│  rover graph publish cursor-analytics   │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│  Apollo Studio Registry                 │
│  ─────────────────────────────────────  │
│  • Schema v1.2.3 stored                 │
│  • Breaking changes detected            │
│  • Documentation generated              │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│  P6: Fetch Schema for Codegen           │
│  ─────────────────────────────────────  │
│  rover graph fetch cursor-analytics     │
│  npm run codegen                        │
└─────────────────────────────────────────┘
```

---

## References

- **Mitigation Plan**: `docs/MITIGATION-PLAN.md` (Strategy 4)
- **Data Contract Testing**: `docs/data-contract-testing.md` (Phase 2)
- **Apollo Studio**: https://studio.apollographql.com
- **Rover CLI**: https://www.apollographql.com/docs/rover/
