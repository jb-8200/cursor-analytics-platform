# Spec-Driven Development Workflow

This skill defines the SDD workflow for this project using Claude Code.

## Work Item Structure

Features are tracked in `.work-items/{feature-name}/`:

```
.work-items/{feature}/
├── user-story.md    # EARS-format requirements
├── design.md        # Technical design with decision log
├── task.md          # Step breakdown with progress
└── {NN}_step.md     # Detailed implementation steps
```

## Active Work Tracking

The currently active feature is tracked via symlink:
```
.claude/plans/active -> ../../.work-items/{feature}/task.md
```

Check what's active: `ls -la .claude/plans/active`

## Phase Gates

### Phase 1: User Story
- Persona (As a... I want... So that...)
- EARS-format acceptance criteria
- Out of scope defined
- Dependencies listed

### Phase 2: Design
- Decision log with rationale
- Architecture diagram
- Data models
- API contracts
- Risks and mitigations

### Phase 3: Task Breakdown
- Steps sized to 1-4 hours
- Dependencies mapped
- Model recommendations (Haiku/Sonnet/Opus)
- ACs mapped to steps

### Phase 4: Implementation
- TDD: RED → GREEN → REFACTOR
- One step at a time
- Commit after each green state
- Time tracking in commits

## TDD Cycle

1. **RED**: Write failing test
   - Test defines expected behavior
   - Confirm test fails for the right reason

2. **GREEN**: Minimal implementation
   - Write just enough to pass
   - Don't add extra functionality

3. **REFACTOR**: Clean up
   - Improve code while tests green
   - Extract patterns, improve names

## Commit Message Format

```
{type}: {brief description}

{details}

Time Tracking:
- Estimated: {E} hours
- Actual: {A} hours
- Delta: {D} hours

Generated with Claude Code
Co-Authored-By: Claude <noreply@anthropic.com>
```

## Commands

- `/start-feature {name}` - Begin work on a feature
- `/complete-feature {name}` - Verify and close feature
- `/implement {task-id}` - Implement specific task with TDD
- `/next-task` - Show next task to work on

## Coverage Targets

- Minimum: 80% line coverage
- Generators/Core: 90%+ coverage
- Run: `go test ./... -cover`
