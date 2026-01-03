# User Story: cursor-sim Foundation

**Feature ID**: P1-F01-foundation
**Phase**: P1 (cursor-sim Foundation)
**Created**: January 2, 2026
**Status**: COMPLETE

## Persona

**As a** developer analytics researcher
**I want** a simulator that exactly matches the Cursor Business API
**So that** I can test analytics pipelines without production access and generate correlated datasets for SDLC research

## Story

### Background

Organizations adopting AI coding assistants need visibility into how teams use these tools. Research teams need correlated datasets to study AI's impact on velocity, review costs, and code quality. The production Cursor API is not accessible for testing or research purposes.

### Need

A simulator that:
1. Produces data from seed files (reproducible, correlated)
2. Exposes exact Cursor Business API endpoints (drop-in replacement)
3. Supports research use cases (PR lifecycle, quality outcomes)

## Acceptance Criteria (EARS Format)

### AC-1: Seed Loading
**When** the simulator starts with `--mode=runtime --seed=seed.json`
**Then** it SHALL load developers, repositories, correlations, and text templates from the seed file
**And** validate all required fields are present
**And** report descriptive errors for invalid seeds

### AC-2: Cursor Admin API
**When** an authenticated client requests `/teams/members`
**Then** the response SHALL match the exact Cursor API schema
**And** include pagination with page/pageSize/totalUsers/totalPages
**And** return 401 for missing/invalid credentials

### AC-3: AI Code Tracking API
**When** a client requests `/analytics/ai-code/commits` with date range
**Then** the response SHALL include commits with exact field names:
- commitHash, userId, userEmail, repoName
- totalLinesAdded, tabLinesAdded, composerLinesAdded, nonAiLinesAdded
**And** AI lines SHALL sum correctly: `total = tab + composer + nonAi`

### AC-4: Team Analytics API
**When** a client requests any of the 11 `/analytics/team/*` endpoints
**Then** the response SHALL aggregate data by event_date
**And** support startDate/endDate filtering
**And** enforce rate limiting (100 req/min)

### AC-5: By-User Analytics API
**When** a client requests any of the 9 `/analytics/by-user/*` endpoints
**Then** the response SHALL group data by user_email
**And** include userMappings in params section
**And** enforce rate limiting (50 req/min)

### AC-6: CSV Export
**When** a client requests `/analytics/ai-code/commits.csv`
**Then** the response SHALL be RFC 4180 compliant CSV
**And** include Content-Disposition header with filename

### AC-7: Reproducibility
**When** the simulator runs twice with the same seed file
**Then** the generated data SHALL be identical
**And** commit hashes SHALL be deterministic

### AC-8: Performance
**When** serving 100k commits
**Then** API response time (p99) SHALL be < 50ms
**And** memory usage SHALL be < 500MB
**And** startup time SHALL be < 2 seconds

## Out of Scope (Phase 1)

- GitHub PR simulation endpoints (Phase 2)
- Quality outcomes (reverts, hotfixes) (Phase 2)
- Research dataset export (Phase 3)
- Replay mode with Parquet (Phase 3)

## Dependencies

- DataDesigner seed generator (available in tools/data-designer/)
- OpenAPI specs (available in specs/openapi/)

## References

- Cursor Analytics API: https://cursor.com/docs/account/teams/analytics-api
- Cursor AI Code Tracking: https://docs.cursor.com/business/api-reference/ai-code-tracking
- Internal design: docs/DESIGN.md
