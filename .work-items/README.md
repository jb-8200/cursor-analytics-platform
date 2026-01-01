# Work Items

This directory contains feature work items following Spec-Driven Development (SDD) methodology.

## Structure

Each work item is a numbered directory:

```
.work-items/
├── README.md                    # This file
├── 01-p0-scaffolding/          # Example work item
│   ├── user-story.md           # Acceptance criteria (Given-When-Then)
│   ├── design.md               # Technical design decisions
│   ├── task.md                 # Implementation checklist
│   └── test-plan.md            # Test cases (write these FIRST)
├── 02-developer-generator/
│   ├── user-story.md
│   ├── design.md
│   ├── task.md
│   └── test-plan.md
└── archive/                    # Completed work items (optional)
    └── 01-p0-scaffolding/
```

## File Templates

### user-story.md
```markdown
# User Story: {Feature Name}

## As a...
[User role or system component]

## I want to...
[Capability or feature]

## So that...
[Business value or technical benefit]

## Acceptance Criteria

### Scenario 1: {Scenario Name}
**Given** [initial context]
**When** [action occurs]
**Then** [expected outcome]

### Scenario 2: {Another Scenario}
**Given** [initial context]
**When** [action occurs]
**Then** [expected outcome]
```

### design.md
```markdown
# Design: {Feature Name}

## Overview
[1-2 paragraph summary of the feature and its purpose]

## Technical Decisions

### Decision 1: {Topic}
- **Options Considered**: A, B, C
- **Chosen**: B
- **Rationale**: [why B is best for this use case]
- **Trade-offs**: [what we're giving up by not choosing A or C]

### Decision 2: {Topic}
...

## API Changes
[New endpoints, modified endpoints, or schema changes]

## Data Model Changes
[Database schema updates, new types, modified structs]

## Dependencies
[External libraries, services, or APIs]

## Security Considerations
[Auth requirements, input validation, rate limiting, etc.]

## Performance Considerations
[Expected load, scaling concerns, optimization strategies]

## References
- services/{service}/SPEC.md (lines {nn}-{mm})
- docs/USER_STORIES.md (lines {nn}-{mm})
```

### task.md
```markdown
# Tasks: {Feature Name}

## Phase 1: Tests (Red)
- [ ] Write failing unit test for {component}
- [ ] Write failing integration test for {endpoint}
- [ ] Verify tests fail with correct error messages
- [ ] Verify test coverage baseline established

## Phase 2: Implementation (Green)
- [ ] Implement {component} to pass unit tests
- [ ] Implement {endpoint} to pass integration tests
- [ ] Verify all tests pass
- [ ] Verify no existing tests broken

## Phase 3: Refactor
- [ ] Extract reusable utilities
- [ ] Add inline documentation (comments where logic isn't obvious)
- [ ] Run linter and fix all issues
- [ ] Verify tests still pass after refactoring

## Phase 4: Documentation
- [ ] Update SPEC.md if behavior changed
- [ ] Add/update API documentation
- [ ] Update CHANGELOG.md with user-facing changes
- [ ] Update README.md if needed

## Phase 5: Review
- [ ] Self-review code changes
- [ ] Check test coverage (>= 80%)
- [ ] Verify all acceptance criteria met (user-story.md)
- [ ] Ready for PR
```

### test-plan.md
```markdown
# Test Plan: {Feature Name}

## Unit Tests

### {Component}.test.{ext}

#### Test 1: {Description}
- **Input**: {test data}
- **Expected**: {assertion}
- **Edge Cases**: {list edge cases}

#### Test 2: {Description}
...

## Integration Tests

### {Endpoint} - {Scenario}
- **Method**: GET/POST/PUT/DELETE
- **Path**: /path/to/endpoint
- **Auth**: Required/Optional/None
- **Request Headers**: {...}
- **Request Body**: {...}
- **Expected Status**: 200/201/400/etc
- **Expected Response**: {...}
- **Edge Cases**: {list edge cases}

## Edge Cases to Cover
1. Empty input
2. Null values
3. Extremely large input
4. Invalid format
5. Unauthorized access
6. Missing required fields
7. Duplicate entries

## Performance Tests (if applicable)
- **Load**: {number} requests/second
- **Latency**: < {threshold} ms
- **Memory**: < {threshold} MB

## Security Tests (if applicable)
- SQL injection attempts
- XSS attempts
- Auth bypass attempts
- Rate limit enforcement

## Coverage Target
- **Minimum**: 80%
- **Target**: 90%+
```

## Workflow

### 1. Create Work Item
```bash
# Using slash command
/start-feature basic-auth-middleware

# Or manually
mkdir -p .work-items/03-basic-auth-middleware
cd .work-items/03-basic-auth-middleware
touch user-story.md design.md task.md test-plan.md
```

### 2. Write Specifications FIRST
1. Define acceptance criteria in `user-story.md`
2. Make technical decisions in `design.md`
3. Plan tests in `test-plan.md`
4. Create task checklist in `task.md`

### 3. TDD Implementation (Red-Green-Refactor)
1. **RED**: Write failing tests based on test-plan.md
2. **GREEN**: Implement just enough to pass tests
3. **REFACTOR**: Clean up while keeping tests green
4. Check off tasks in task.md as you complete them

### 4. Complete Work Item
1. Verify all acceptance criteria met
2. Verify all tasks checked off
3. Run full test suite
4. Update documentation
5. Optionally move to `archive/`

## Numbering Convention

Work items are numbered sequentially starting from 01:

- `01-p0-scaffolding` - First work item (P0 scaffolding)
- `02-developer-generator` - Second work item
- `03-basic-auth-middleware` - Third work item
- etc.

Use leading zeros for sorting (01, 02, ..., 09, 10, 11, ...)

## Archiving Completed Work

Optional: Move completed work items to archive:

```bash
mkdir -p .work-items/archive
mv .work-items/01-p0-scaffolding .work-items/archive/
```

Or keep them in the main directory for historical reference.

## Integration with .claude/plans/

The active work item is tracked via symlink:

```bash
.claude/plans/active -> ../../.work-items/03-basic-auth-middleware/
```

This allows AI assistants to quickly find the current work context.

## Best Practices

1. **One work item = One PR** - Keep work items focused
2. **Write specs before code** - user-story.md and design.md before task.md
3. **Write tests before implementation** - test-plan.md before actual test files
4. **Small, focused work items** - Better than large, monolithic ones
5. **Check off tasks as you go** - Provides progress visibility
6. **Reference SPECs** - Always link to relevant SPEC.md sections in design.md

## Example Work Item Names

- `01-p0-scaffolding` - Initial project setup
- `02-developer-profile-generator` - Generate realistic developer profiles
- `03-basic-auth-middleware` - HTTP Basic Auth implementation
- `04-commit-event-generator` - Generate commit events
- `05-pagination-helper` - Pagination utility
- `06-rate-limiter` - Rate limiting middleware

Use descriptive, kebab-case names that clearly indicate what the work item covers.

---

**Remember**: Specifications → Tests → Implementation → Refactor → Documentation

This is the Spec-Driven Development way.
