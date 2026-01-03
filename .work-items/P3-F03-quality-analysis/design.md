# Technical Design: Quality Analysis

**Feature ID**: P3-F03-quality-analysis
**Phase**: P3 (cursor-sim Research Framework)
**Created**: January 3, 2026
**Status**: COMPLETE

## Overview

This feature implements GitHub simulation with quality analysis, including PR generation, code survival, revert detection, and hotfix tracking.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                   Quality Analysis                       │
├─────────────────────────────────────────────────────────┤
│  PR Generation Pipeline                                  │
│  ├── Session Model (seniority-based parameters)         │
│  ├── Commit Grouping (inactivity gaps)                  │
│  └── PR Envelope (metrics aggregation)                  │
├─────────────────────────────────────────────────────────┤
│  Quality Services                                        │
│  ├── SurvivalService (file-level tracking)              │
│  ├── RevertService (risk scoring, pattern matching)     │
│  └── HotfixService (file overlap detection)             │
├─────────────────────────────────────────────────────────┤
│  GitHub API Handlers                                     │
│  ├── /repos/{owner}/{repo}/analysis/survival            │
│  ├── /repos/{owner}/{repo}/analysis/reverts             │
│  └── /repos/{owner}/{repo}/analysis/hotfixes            │
└─────────────────────────────────────────────────────────┘
```

## Key Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| PR Generation | Session-based | Natural commit groupings |
| Greenfield | First commit timestamp | Simple, deterministic |
| Quality Correlations | Sigmoid risk score | Smooth probability curves |
| Code Survival | File-level | Simple, fast, sufficient |
| Replay Mode | Deferred to P3-D | Use seeded RNG for now |

## Models

### Session Model
```go
type Session struct {
    Developer     seed.Developer
    MaxCommits    int           // junior=2-5, mid=4-8, senior=5-12
    TargetLoC     int           // junior=50-150, mid=100-300, senior=150-500
    InactivityGap time.Duration // 15-60 minutes
}
```

### Quality Models
- FileSurvival: Track file birth/death, AI vs human lines
- RevertEvent: Link reverts to original PRs
- HotfixEvent: Track follow-up fixes

---

**Status**: COMPLETE
