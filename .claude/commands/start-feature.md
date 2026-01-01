# Start New Feature (SDD Workflow)

When the user runs `/start-feature [feature-name]`, guide them through creating a new feature following Spec-Driven Development methodology.

## Workflow Steps

### 1. Check Specifications First
- Ask which service this feature belongs to (cursor-sim, cursor-analytics-core, cursor-viz-spa)
- Read the service's SPEC.md to understand existing architecture
- Check if feature is already documented in docs/FEATURES.md or docs/USER_STORIES.md

### 2. Create Work Item Structure

Create directory: `.work-items/{nn}-{feature-name}/` where {nn} is next number

Files to create:
```
.work-items/{nn}-{feature-name}/
├── user-story.md       # Given-When-Then acceptance criteria
├── design.md           # Technical design decisions
├── task.md             # Implementation checklist
└── test-plan.md        # Test cases (write these FIRST)
```

### 3. Write Specification Documents

**user-story.md template:**
```markdown
# User Story: {Feature Name}

## As a...
[User role]

## I want to...
[Capability]

## So that...
[Business value]

## Acceptance Criteria

### Scenario 1: {Scenario Name}
**Given** [initial context]
**When** [action occurs]
**Then** [expected outcome]

### Scenario 2: ...
```

**design.md template:**
```markdown
# Design: {Feature Name}

## Overview
[1-paragraph summary]

## Technical Decisions

### Decision 1: {Topic}
- **Options Considered**: A, B, C
- **Chosen**: B
- **Rationale**: [why]
- **Trade-offs**: [what we're giving up]

## API Changes
[New endpoints, types, or schema changes]

## Data Model Changes
[Database schema updates if any]

## Dependencies
[External libraries or services]

## Security Considerations
[Auth, input validation, etc.]
```

**task.md template:**
```markdown
# Tasks: {Feature Name}

## Phase 1: Tests (Red)
- [ ] Write failing unit tests for {component}
- [ ] Write failing integration tests for {endpoint}
- [ ] Verify tests fail with correct error messages

## Phase 2: Implementation (Green)
- [ ] Implement {component} to pass unit tests
- [ ] Implement {endpoint} to pass integration tests
- [ ] Verify all tests pass

## Phase 3: Refactor
- [ ] Extract reusable utilities
- [ ] Add inline documentation
- [ ] Run linter and fix issues
- [ ] Verify tests still pass

## Phase 4: Documentation
- [ ] Update SPEC.md if behavior changed
- [ ] Add API documentation
- [ ] Update CHANGELOG.md
```

**test-plan.md template:**
```markdown
# Test Plan: {Feature Name}

## Unit Tests

### {Component}.test.{ext}
- **Test**: {description}
- **Input**: {test data}
- **Expected**: {assertion}

## Integration Tests

### {Endpoint} - {Scenario}
- **Method**: GET/POST/etc
- **Path**: /path/to/endpoint
- **Auth**: Required/Optional
- **Request Body**: {...}
- **Expected Status**: 200
- **Expected Response**: {...}

## Edge Cases
1. {Edge case 1}
2. {Edge case 2}

## Performance Tests
- **Load**: {number} requests/second
- **Latency**: < {threshold} ms
```

### 4. Link to Active Plan

If `.claude/plans/` directory exists, create symlink:
```
.claude/plans/active -> ../../.work-items/{nn}-{feature-name}/
```

### 5. Remind Developer of TDD Workflow

Display:
```
✓ Feature structure created: .work-items/{nn}-{feature-name}/

Next Steps (TDD):
1. Write tests first (test-plan.md → actual test files)
2. Run tests and verify they FAIL
3. Implement just enough code to make tests PASS
4. Refactor while keeping tests green
5. Update documentation

Remember: Red → Green → Refactor

Recommended Model for Planning: sonnet ⚡⚡
(Architectural decisions and design - then use haiku for implementation)
```

## Example Usage

```
User: /start-feature basic-auth-middleware
Assistant: [Creates .work-items/03-basic-auth-middleware/ with all templates]
```

## Implementation Instructions

Guide the user through creating proper SDD artifacts before writing any implementation code. Ensure tests are specified before implementation begins.
