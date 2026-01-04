---
name: model-selection-guide
description: Guide for selecting the right Claude model (Haiku, Sonnet, Opus) for different tasks. Use when starting a new task, planning implementation, or optimizing costs. Helps choose between fast/cheap Haiku and powerful Sonnet/Opus. (project)
model: haiku
---

# Model Selection Guide

This skill provides guidance on which Claude model to use for different tasks.

## Model Capabilities Overview

| Model | Cost | Speed | Best For |
|-------|------|-------|----------|
| **Haiku** | $ | Fast | Well-specified implementation, pattern following |
| **Sonnet** | $$ | Balanced | Complex logic, architectural decisions, debugging |
| **Opus** | $$$ | Powerful | Novel architectures, major refactoring |

## Decision Tree

```
Is this well-specified in SPEC.md?
├─ YES → Does it follow an existing pattern?
│   ├─ YES → Use Haiku
│   └─ NO → Is it complex logic?
│       ├─ YES → Use Sonnet
│       └─ NO → Use Haiku
│
└─ NO → Is it a major architectural decision?
    ├─ YES → Use Opus
    └─ NO → Use Sonnet
```

## Use Haiku (Fast & Cheap) for:

**Implementation from Specs**
- Implementing structs/types defined in SPEC.md
- Creating HTTP handlers following cursor-api-patterns
- Writing Go code following go-best-practices
- Building React components with clear requirements

**Writing Tests**
- Table-driven tests from test-plan.md
- Unit tests for well-specified functions
- Integration tests following patterns

**Pattern Following**
- Adding pagination to endpoints
- Implementing Basic Auth
- Error handling patterns
- JSON marshaling

**Simple Refactoring**
- Extract helper functions
- Rename variables/functions
- Format code

## Use Sonnet (Balanced) for:

**Complex Logic**
- Poisson distribution event generation
- Rate limiting algorithms
- Caching strategies
- Concurrent worker pools

**Architectural Decisions**
- Service integration patterns
- Database schema design
- API versioning strategy

**Multi-File Changes**
- Feature implementation spanning multiple files
- Refactoring across packages
- Adding middleware layers

**Debugging**
- Race condition investigation
- Performance bottleneck identification
- Integration test failures

## Use Opus (Powerful) for:

**Novel Architectures**
- Designing new major features not in specs
- Complex refactoring across all services
- Performance optimization strategy

**Cross-Service Integration**
- End-to-end feature implementation
- Contract negotiation between services
- Breaking API changes

**Strategic Planning**
- Choosing technology stack
- Migration planning
- Security architecture review

## Cost Optimization Strategy

### Session-Based Approach

1. **Planning Session (Sonnet/Opus)**: Create work items, design docs
2. **Implementation Session (Haiku)**: Write tests and code from specs
3. **Integration Session (Sonnet)**: Wire up components, debug

### Cost Comparison

**All Sonnet**: ~$4.25 per feature
**Hybrid Approach**: ~$1.55 per feature (64% savings)

## Task-to-Model Mapping

### cursor-sim Tasks

| Task | Model | Rationale |
|------|-------|-----------|
| Initialize Go Project | Sonnet | Project structure decisions |
| CLI Flag Parsing | Haiku | Standard flag/JSON parsing |
| Developer Profile Generator | Haiku | Struct from SPEC.md |
| Event Generation Engine | Sonnet | Complex timing logic |
| In-Memory Storage | Haiku | Standard sync.Map patterns |
| REST API Handlers | Haiku | Following cursor-api-patterns |
| Wire Up Main | Sonnet | Integration and error handling |

## When in Doubt

**Default to Sonnet** if:
- You're unsure of the complexity
- It's the first time implementing this pattern
- The task involves multiple services

**Switch to Haiku** once:
- The pattern is established
- You have clear acceptance criteria
- The task is well-isolated

## Summary

The comprehensive specs in this project make it ideal for Haiku usage. Use Sonnet/Opus for planning and complex logic, then leverage Haiku for well-specified implementation work.
