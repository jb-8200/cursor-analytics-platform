---
name: spec-process-core
description: Core SDD (Spec-Driven Development) process. Use when starting development work, planning features, or asking about workflow. Covers the fundamental sequence of Spec, Test, Code, Refactor, Commit.
---

# Core SDD Process

Every piece of work follows this sequence. No exceptions.

## The Fundamental Sequence

```
1. SPEC    → Define what we're building (user story, design)
2. TEST    → Encode expectations as tests (RED phase)
3. CODE    → Write minimal implementation (GREEN phase)
4. REFACTOR → Clean up while tests pass
5. COMMIT  → Checkpoint with descriptive message
6. REPEAT  → Next task in the breakdown
```

## Core Principles

### 1. Specifications Before Code

**Never write code without reading the specification first.**

Check these locations in order:
1. `services/{service}/SPEC.md` - Technical specification
2. `.work-items/{feature}/user-story.md` - Requirements
3. `.work-items/{feature}/design.md` - Technical design
4. `.work-items/{feature}/task.md` - Task breakdown

If specifications are missing or incomplete, **create them first**.

### 2. Tests Define the Contract

**Tests are written BEFORE implementation code.**

- Acceptance criteria become test cases
- Given-When-Then maps to Arrange-Act-Assert
- Tests fail initially (RED) - this confirms they're testing something
- Implementation makes tests pass (GREEN)

### 3. Small Batch Sizes

**Break work into 1-4 hour tasks.**

Each task should:
- Have clear acceptance criteria
- Be independently completable
- Result in a single commit
- Update progress tracking

### 4. Commit After Every Task

**CRITICAL: Never move to the next task without committing.**

The commit sequence:
1. Tests pass
2. Stage changes
3. Commit with descriptive message
4. Update task.md progress
5. Update DEVELOPMENT.md
6. Proceed to next task

### 5. Documentation Stays Current

**If code changes, documentation changes.**

When implementation reveals:
- Spec gaps → Update spec
- Design changes → Update design.md
- New requirements → Update user-story.md

## When to Apply This Process

### Full Process (New Features)

1. Create `.work-items/{feature}/` directory
2. Write user-story.md
3. Write design.md
4. Create task.md breakdown
5. Implement via TDD loop
6. Complete feature

### Abbreviated Process (Bug Fixes)

1. Read relevant SPEC.md
2. Write failing test reproducing bug
3. Fix bug
4. Verify test passes
5. Commit

### Minimal Process (Documentation Updates)

1. Make documentation changes
2. Commit
3. No TDD required for pure docs

## Red Flags

Stop immediately if you notice:

| Flag | Required Action |
|------|-----------------|
| Writing code before reading spec | STOP. Read spec first. |
| Implementing without failing test | STOP. Write test first. |
| Moving to next task without commit | STOP. Commit first. |
| Changing behavior without updating spec | STOP. Update spec. |

## Related Skills

- `spec-process-dev` - Detailed development workflow
- `sdd-checklist` - Post-task enforcement
- `spec-user-story` - User story format
- `spec-design` - Design document format
- `spec-tasks` - Task breakdown format
