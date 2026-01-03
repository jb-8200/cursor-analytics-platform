# Task Breakdown: Research Framework

**Feature ID**: P3-F01-research-framework
**Phase**: P3 (cursor-sim Research Framework)
**Created**: January 3, 2026
**Status**: COMPLETE

## Overview

**Total Estimated Hours**: 15-20h
**Total Actual Hours**: 1.75h
**Number of Tasks**: 7

## Progress Tracker

| Task ID | Task | Hours Est | Status | Actual |
|---------|------|-----------|--------|--------|
| TASK01 | Research Data Models | 2.0 | DONE | 0.25 |
| TASK02 | Research Dataset Generator | 3.0 | DONE | 0.25 |
| TASK03 | Parquet/CSV Export | 2.5 | DONE | 0.25 |
| TASK04 | Research Metrics Service | 3.0 | DONE | 0.25 |
| TASK05 | Research API Handlers | 2.5 | DONE | 0.25 |
| TASK06 | Replay Mode Infrastructure | 3.0 | DONE | 0.25 |
| TASK07 | Integration Tests | 2.0 | DONE | 0.25 |

---

## Task Details

### TASK01: Research Data Models

**File**: `internal/models/research.go`

Create `ResearchDataPoint` struct with all fields for SDLC research metrics.

### TASK02: Research Dataset Generator

**File**: `internal/generator/research_generator.go`

Generate correlated research data points from commits and PRs.

### TASK03: Parquet/CSV Export

**Files**: `internal/export/csv.go`, `internal/export/json.go`

Implement exporters for research dataset in CSV and JSON formats.

### TASK04: Research Metrics Service

**File**: `internal/services/research_metrics.go`

Calculate aggregated metrics by AI ratio bands.

### TASK05: Research API Handlers

**Files**: `internal/api/research/*.go`

Implement handlers for /research/* endpoints.

### TASK06: Replay Mode Infrastructure

Seeded RNG for reproducible generation.

### TASK07: Integration Tests

E2E tests for all research endpoints.

---

**Status**: COMPLETE
