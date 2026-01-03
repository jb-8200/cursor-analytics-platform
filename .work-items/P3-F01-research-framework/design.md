# Technical Design: Research Framework

**Feature ID**: P3-F01-research-framework
**Phase**: P3 (cursor-sim Research Framework)
**Created**: January 3, 2026
**Status**: COMPLETE

## Overview

This feature implements the research data export framework for SDLC studies, including data models, generators, exporters, and API handlers.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Research Framework                    │
├─────────────────────────────────────────────────────────┤
│  internal/models/research.go     - ResearchDataPoint    │
│  internal/generator/research.go  - Dataset Generator    │
│  internal/export/csv.go          - CSV Exporter         │
│  internal/export/json.go         - JSON Exporter        │
│  internal/services/research.go   - Metrics Calculator   │
│  internal/api/research/          - API Handlers         │
└─────────────────────────────────────────────────────────┘
```

## Data Models

### ResearchDataPoint

Core model for research dataset with fields supporting:
- AI usage metrics (ai_ratio, ai_lines_added)
- Cycle times (coding_lead_time, pickup_time, review_lead_time)
- Review costs (review_density, iteration_count, rework_ratio)
- Quality outcomes (is_reverted, survival_rate, hotfix_followup)

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /research/dataset` | Export full dataset (CSV/JSON) |
| `GET /research/metrics/velocity` | Velocity metrics by AI band |
| `GET /research/metrics/review-costs` | Review cost metrics |
| `GET /research/metrics/quality` | Quality metrics |

---

**Status**: COMPLETE (1.75h actual / 15-20h estimated)
