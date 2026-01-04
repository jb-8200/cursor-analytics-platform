# User Story: GraphQL Code Generator Setup

**Feature ID**: P6-F02
**Epic**: P6 - cursor-viz-spa (Service Decoupling Phase 1)
**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: HIGHEST

---

## User Story

**As a** frontend developer working on cursor-viz-spa,
**I want** TypeScript types to be auto-generated from the P5 GraphQL schema,
**So that** schema mismatches are caught at compile-time instead of runtime.

---

## Background

During P5+P6 integration testing (January 4, 2026), GraphQL schema mismatches caused complete integration failure:
- P6 manually defined `topPerformers: Developer[]` but P5 had `topPerformer: Developer`
- P6 used `humanLinesAdded` but P5 had `linesAdded`
- TypeScript passed locally because it validated against local (wrong) types

This resulted in 400 Bad Request errors only visible at runtime.

**Root Cause**: Manual type definitions in P6 drifted from P5's actual schema.

**Solution**: Auto-generate TypeScript types from P5 GraphQL schema using GraphQL Code Generator.

---

## Acceptance Criteria

### AC-01: Code Generator Installation
- **Given** cursor-viz-spa project
- **When** developer runs `npm install`
- **Then** GraphQL Code Generator and plugins are installed

### AC-02: Schema Introspection
- **Given** P5 GraphQL server is running at localhost:4000
- **When** developer runs `npm run codegen`
- **Then** types are generated from P5 schema to `src/graphql/generated.ts`

### AC-03: Type Usage in Queries
- **Given** generated types exist in `src/graphql/generated.ts`
- **When** developer imports types in query files
- **Then** TypeScript validates queries against actual P5 schema

### AC-04: Pre-Dev Hook
- **Given** developer runs `npm run dev`
- **When** development server starts
- **Then** codegen runs automatically before dev server (predev hook)

### AC-05: Pre-Build Hook
- **Given** developer runs `npm run build`
- **When** build process starts
- **Then** codegen runs automatically before build (prebuild hook)

### AC-06: CI Validation
- **Given** PR is opened with P6 changes
- **When** CI pipeline runs
- **Then** build fails if generated types differ from committed version (schema drift detection)

---

## Out of Scope

- P5 schema changes (handled by P5 team)
- Apollo Studio setup (Phase 2 feature)
- Contract testing (Phase 3 feature)

---

## Dependencies

| Dependency | Description | Status |
|------------|-------------|--------|
| P5 GraphQL Server | Must be running for introspection | Available |
| @graphql-codegen/cli | Code generator tool | To install |

---

## Technical Context

### Current State (Broken)
```
P6/src/graphql/types.ts (manual, drifts)
├── interface TeamStats { topPerformers: Developer[] }  ← WRONG
└── Used by: queries.ts, components/*
```

### Target State (Fixed)
```
P6/src/graphql/generated.ts (auto-generated)
├── interface TeamStats { topPerformer: Developer }  ← CORRECT
└── Used by: queries.ts, components/*
```

---

## References

- **Mitigation Plan**: `docs/MITIGATION-PLAN.md` (Strategy 1)
- **Data Contract Testing**: `docs/data-contract-testing.md` (Phase 1)
- **GraphQL Code Generator**: https://the-guild.dev/graphql/codegen
