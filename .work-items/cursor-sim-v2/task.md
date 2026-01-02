# Task Breakdown: cursor-sim v2

## Overview

**Feature**: cursor-sim v2 - Cursor API Simulator
**Total Estimated Hours**: 44.5
**Number of Steps**: 16
**Current Step**: 13 (Complete) - Next: 14

## Progress Tracker

| Step | Task                          | Hours | Status      | Actual |
|------|-------------------------------|-------|-------------|--------|
| 01   | Project structure             | 1.0   | DONE        | 0.25   |
| 02   | Seed schema types             | 2.0   | DONE        | 0.5    |
| 03   | Seed loader + validation      | 3.0   | DONE        | 0.75   |
| 04   | CLI v2 flags                  | 2.0   | DONE        | 0.5    |
| 05   | Cursor data models            | 3.0   | DONE        | 0.5    |
| 06   | Commit generation engine      | 5.0   | DONE        | 1.25   |
| 07   | In-memory storage v2          | 4.0   | DONE        | 1.0    |
| 08   | API infrastructure            | 2.0   | DONE        | 0.75   |
| 09   | /teams/members endpoint       | 1.5   | DONE        | 0.5    |
| 10   | /ai-code/commits endpoint     | 2.0   | DONE        | 0.75   |
| 11   | /ai-code/commits.csv endpoint | 1.0   | DONE        | 0.25   |
| 12   | /team/* endpoints (11)        | 6.0   | DONE        | 1.5    |
| 13   | /by-user/* endpoints (9)      | 4.0   | DONE        | 0.5    |
| 14   | HTTP router                   | 2.0   | NOT_STARTED | -      |
| 15   | Main application              | 2.0   | NOT_STARTED | -      |
| 16   | E2E testing                   | 4.0   | NOT_STARTED | -      |

## Dependency Graph

```text
Step 01 (Structure)
    │
    ├── Step 02 (Seed Types)
    │       │
    │       └── Step 03 (Seed Loader)
    │               │
    │               └── Step 06 (Commit Generator)
    │
    ├── Step 04 (CLI) ──────────────────────┐
    │                                        │
    ├── Step 05 (Cursor Models) ────────────┤
    │                                        │
    └── Step 08 (API Infrastructure) ───────┤
                                             │
Step 06 + Step 07 (Storage) ─────────────────┤
                                             │
Steps 09-13 (Endpoints) ◄────────────────────┘
    │
    └── Step 14 (Router)
            │
            └── Step 15 (Main)
                    │
                    └── Step 16 (E2E Tests)
```

## Critical Path

01 → 02 → 03 → 06 → 07 → 09 → 14 → 15 → 16

## Model Recommendations

| Step                                   | Model  | Rationale                                     |
|----------------------------------------|--------|-----------------------------------------------|
| 01, 02, 04, 05, 08, 09, 11, 14, 15     | Haiku  | Well-specified, low complexity                |
| 03, 06, 07, 10, 12, 13, 16             | Sonnet | Validation/generation/aggregation complexity  |

## Step Files

Each step has a detailed implementation file:

- `01_project_structure.md`
- `02_seed_types.md`
- `03_seed_loader.md`
- ... etc.

## TDD Checklist (Per Step)

- [ ] Read step file and acceptance criteria
- [ ] Write failing test(s) for the step
- [ ] Run tests, confirm RED
- [ ] Implement minimal code to pass
- [ ] Run tests, confirm GREEN
- [ ] Refactor while green
- [ ] Run linter (golangci-lint)
- [ ] Update step status to DONE
- [ ] Commit with time tracking

## Time Tracking Summary

| Category            | Estimated | Actual           | Delta                 |
|---------------------|-----------|------------------|-----------------------|
| Foundation (01-05)  | 11.0h     | 2.5h             | -77% ✅                |
| Generation (06-07)  | 9.0h      | 2.25h            | -75% ✅                |
| Endpoints (08-13)   | 16.5h     | 4.25h            | -74% ✅                |
| Integration (14-16) | 8.0h      | -                | -                     |
| **Total**           | **44.5h** | **9.0h**         | **-80% (so far)** ✅   |

**Note**: Actual time is significantly lower due to:

- TDD methodology (fewer rewrites)
- Auto-approved file operations (no permission delays)
- Parallel tool execution
- Well-specified requirements in SPEC.md

## Acceptance Criteria Mapping

| AC | Steps |
| ---- | ----- |
| AC-1: Seed Loading | 02, 03 |
| AC-2: Admin API | 09 |
| AC-3: AI Code Tracking | 10, 11 |
| AC-4: Team Analytics | 12 |
| AC-5: By-User Analytics | 13 |
| AC-6: CSV Export | 11 |
| AC-7: Reproducibility | 06 |
| AC-8: Performance | 07, 16 |
