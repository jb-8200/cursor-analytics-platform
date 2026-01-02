# Claude Code Hooks

This directory contains automation hooks for the Spec-Driven Development workflow.

## Implemented Hooks

### pre_prompt.py

**Purpose**: Automatically injects context when Claude starts a coding session.

**Functionality**:
- Reads active work item from `.claude/plans/active` symlink
- Loads project rules from `rules/*.mdc`
- Includes development context from `.claude/DEVELOPMENT.md`

**Usage**:
```bash
python .claude/hooks/pre_prompt.py
```

**When to Use**: Called by Claude Code before processing prompts to ensure Claude "knows" the active context without explicit reminders.

### pre_commit.py

**Purpose**: Enforces TDD by blocking commits if tests fail.

**Functionality**:
- Detects which services have staged changes
- Runs relevant test suites (Go, TypeScript, React)
- Blocks commit if any tests fail
- Reports test results and coverage

**Installation** (as git hook):
```bash
ln -sf ../../.claude/hooks/pre_commit.py .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

**Manual Usage**:
```bash
python .claude/hooks/pre_commit.py
```

## Hook Integration with Workflow

### Feature Start
1. `/start-feature {name}` creates active symlink
2. `pre_prompt.py` detects active feature
3. Context is automatically loaded

### During Development
1. TDD: RED → GREEN → REFACTOR
2. Attempt `git commit`
3. `pre_commit.py` runs tests
4. Commit blocked if tests fail

### Feature Complete
1. `/complete-feature {name}` verifies all steps
2. Full test suite runs
3. Symlink removed
4. Completion committed

## Directory Integration

```
.claude/
├── hooks/
│   ├── README.md           # This file
│   ├── pre_prompt.py       # Context injection
│   └── pre_commit.py       # Test enforcement
├── plans/
│   ├── README.md           # Symlink mechanism docs
│   └── active -> ...       # Active feature symlink
└── commands/
    ├── start-feature.md    # Begin feature work
    └── complete-feature.md # Verify and close

.work-items/
└── {feature}/
    ├── user-story.md       # Requirements
    ├── design.md           # Technical design
    ├── task.md             # Step breakdown
    └── {NN}_step.md        # Implementation steps

rules/
├── process-core.mdc        # Core SDD principles
├── standards-tdd.mdc       # TDD standards
└── standards-user-story.mdc # User story format
```

## Future Hooks

The following hooks are planned but not yet implemented:

- **post_test.py**: Coverage analysis after test runs
- **validate_spec.py**: Ensures spec files are complete
- **lint_check.py**: Pre-commit linting enforcement
