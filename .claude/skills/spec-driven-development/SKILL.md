---
name: spec-driven-development
description: Comprehensive SDD methodology for Cursor Analytics Platform. Use when learning about SDD practices, understanding procedures, or handling specification gaps. Covers spec location, test writing, and quality standards.
---

# Spec-Driven Development

This skill enables following spec-driven development practices effectively within the Cursor Analytics Platform project.

## Overview

Spec-driven development is an approach where specifications and documentation drive the implementation process. Tests are derived from specifications, and code is written to satisfy those tests.

## Core Principles

### Specifications Come First

Every feature has a specification that defines its behavior. Before implementing anything, locate and understand the relevant specification.

**Service-level specifications** in `services/{service}/SPEC.md` provide detailed technical requirements including data models, API contracts, configuration options, and performance requirements.

**User stories** in `docs/USER_STORIES.md` describe features from the user's perspective with Given-When-Then acceptance criteria that map directly to test cases.

**Task breakdown** in `docs/TASKS.md` lists specific implementation tasks with their dependencies and definitions of done.

### Tests Before Code

Following Test-Driven Development principles, tests must exist and fail before implementation code is written.

**RED phase**: Write a test based on the specification's acceptance criteria. The test should fail because the functionality doesn't exist yet.

**GREEN phase**: Write the minimum code necessary to make the test pass. Avoid additional functionality not covered by tests.

**REFACTOR phase**: Improve the code while keeping tests green. Extract functions, improve names, restructure as needed.

### Documentation Stays Current

When implementation reveals gaps or changes needed in specifications, the documentation must be updated. Specifications should remain accurate reflections of actual behavior.

## Procedures

### Starting a New Feature

1. **Locate the specification**: Check service SPEC.md, search USER_STORIES.md for relevant stories, find the task in TASKS.md
2. **Verify dependencies**: Check that prerequisite tasks are complete
3. **Derive test cases**: Each Given-When-Then scenario becomes at least one test case
4. **Write failing tests**: Encode expected behavior, confirm tests fail for expected reason
5. **Implement minimally**: Write just enough code to make each test pass
6. **Refactor**: Clean up code while tests are green
7. **Update documentation**: If implementation revealed specification gaps

### Locating Relevant Specifications

| Service | Specification | Story Prefix |
|---------|--------------|--------------|
| cursor-sim | `services/cursor-sim/SPEC.md` | US-SIM-* |
| cursor-analytics-core | `services/cursor-analytics-core/SPEC.md` | US-CORE-* |
| cursor-viz-spa | `services/cursor-viz-spa/SPEC.md` | US-VIZ-* |

For cross-cutting features, check `docs/DESIGN.md` for architectural decisions.

### Writing Tests from Specifications

Acceptance criteria in Given-When-Then format map directly to test structure:

| Criterion | Test Section | Purpose |
|-----------|--------------|---------|
| **Given** | Arrange | Test setup, preconditions |
| **When** | Act | Action being tested |
| **Then** | Assert | Expected outcome verification |

Example transformation:

**Criterion**: "Given a developer has 100 suggestions shown and 75 accepted, when I query their acceptance rate, then the result should be 75.0"

```go
func TestAcceptanceRate(t *testing.T) {
    // Arrange (Given)
    dev := NewDeveloper("dev-1")
    dev.SuggestionsShown = 100
    dev.SuggestionsAccepted = 75

    // Act (When)
    rate := dev.AcceptanceRate()

    // Assert (Then)
    assert.Equal(t, 75.0, rate)
}
```

### Handling Specification Gaps

When implementation reveals incomplete or ambiguous specifications:

1. Make a reasonable assumption based on the specification's intent
2. Document the assumption in a code comment
3. Implement based on that assumption
4. Flag the gap by suggesting a specification update

Example: "The SPEC.md should clarify whether acceptance rate returns 0 or null when a developer has no suggestions. I implemented null to avoid implying certainty about a rate that cannot be calculated."

## Quality Standards

### Test Coverage

| Category | Required | Target |
|----------|----------|--------|
| All services | 80% | 85% |
| Critical paths | 95% | 98% |
| Metric calculations | 95% | 98% |

### Code Style

- Go code: Must pass `golangci-lint`
- TypeScript: Strict mode enabled
- React: Functional components with hooks

### Documentation

Code comments should explain **why**, not **what**. The code itself should be readable enough that what it does is clear. Comments add value by explaining reasoning behind decisions.
