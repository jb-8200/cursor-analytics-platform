# Spec-Driven Design (SDD) for Claude Code

**Version**: 2.0
**Last Updated**: January 2, 2026
**Adapted from**: [genai-specs](https://github.com/betsalel-williamson/genai-specs) methodology

---

## Overview

Spec-Driven Design (SDD) is a development methodology where **specifications drive all implementation**. This document adapts the genai-specs model for Claude Code, leveraging Skills, Commands, and structured workflows instead of Cursor-specific features.

### Core Flow

```
Specifications → Tests → Implementation → Documentation
     ↑                                         │
     └─────────────────────────────────────────┘
         (Bidirectional: code informs specs)
```

---

## The SDD Triad (genai-specs Model)

### 1. Always Included (Spine)

Content loaded in **every** Claude Code conversation via `CLAUDE.md`:

| Content | Purpose |
|---------|---------|
| Workflow sequence | Spec → Tests → Code → Commit |
| Documentation hierarchy | Where to find specifications |
| Active work pointer | Current feature focus |
| Commit discipline | Every task = commit |

**Keep CLAUDE.md minimal** (~100 lines). Reference Skills for detailed guidance.

### 2. Conditional Guidance (Skills)

Loaded **when relevant** based on semantic match to your request:

| Category | Location | Trigger Examples |
|----------|----------|------------------|
| **Process** | `.claude/skills/process/` | "start feature", "complete task", "implement" |
| **Standards** | `.claude/skills/standards/` | "write user story", "create design doc" |
| **Guidelines** | `.claude/skills/guidelines/` | Language/tech-specific coding |

Skills activate when your request matches their description. Reference explicitly for reliability:

```
"Following spec-user-story, create a user story for the login feature"
"Using go-best-practices, implement the HTTP handler"
```

### 3. Manual Standards (Artifact Templates)

Invoked explicitly when creating specific artifacts:

| Artifact | Skill to Invoke | Output Location |
|----------|-----------------|-----------------|
| User Story | `spec-user-story` | `.work-items/{feature}/user-story.md` |
| Design Doc | `spec-design` | `.work-items/{feature}/design.md` |
| Task Breakdown | `spec-tasks` | `.work-items/{feature}/task.md` |
| ADR | `spec-adr` | `.work-items/{feature}/adr-{NNN}.md` |

---

## Skill Architecture

### Directory Structure

```
.claude/skills/
├── process/                       # Workflow stage guidance
│   ├── spec-process-core.md       # Core SDD principles (always applicable)
│   ├── spec-process-dev.md        # Development workflow (TDD)
│   └── spec-process-commit.md     # Commit discipline
│
├── standards/                     # Artifact templates
│   ├── spec-user-story.md         # User story format (EARS)
│   ├── spec-design.md             # Design doc format (ADRs)
│   ├── spec-tasks.md              # Task breakdown format
│   └── spec-adr.md                # Architecture Decision Record
│
├── guidelines/                    # Technology-specific
│   ├── go-best-practices.md       # Go patterns (cursor-sim)
│   ├── cursor-api-patterns.md     # Cursor API implementation
│   ├── typescript-patterns.md     # TypeScript patterns (future)
│   └── graphql-patterns.md        # GraphQL patterns (future)
│
└── operational/                   # Day-to-day enforcement
    ├── sdd-checklist.md           # Post-task commit enforcement
    ├── sdd-workflow.md            # Full workflow reference
    ├── model-selection-guide.md   # Model optimization
    └── spec-driven-development.md # SDD methodology details
```

### Skill Wrapper Pattern

Skills reference canonical content without duplication:

```markdown
---
name: spec-user-story
description: Create or revise user stories. Use for PRDs, feature narratives,
             acceptance criteria, and UX-oriented specifications.
---

# User Story Standard

## Process
1. Read the canonical template below
2. Ask for missing inputs (actor, goal, context)
3. Produce artifact following the template exactly
4. Keep value-focused (no technical implementation details)

## Template

[Include template content here]

## Examples

[Include examples here]
```

---

## Workflow Phases

### Phase 1: Define (Spec First)

Before writing any code:

1. **Create work item directory**: `.work-items/{feature-name}/`
2. **Write user story** using `spec-user-story` skill
3. **Write design doc** using `spec-design` skill
4. **Create task breakdown** using `spec-tasks` skill

**Command**: `/start-feature {feature-name}`

### Phase 2: Contract (Tests First)

Tests define the contract before implementation:

1. **Derive test cases** from acceptance criteria
2. **Write failing tests** (RED phase)
3. **Verify tests fail for the right reason**

Given-When-Then → Arrange-Act-Assert:

```
AC: Given a developer with 100 suggestions, 75 accepted
    When I calculate acceptance rate
    Then result is 75.0%
```

Becomes:

```go
func TestAcceptanceRate(t *testing.T) {
    // Arrange (Given)
    dev := createDeveloperWithMetrics(100, 75)

    // Act (When)
    rate := dev.AcceptanceRate()

    // Assert (Then)
    assert.Equal(t, 75.0, rate)
}
```

### Phase 3: Implement (Minimal Code)

Write minimum code to pass tests:

1. **Implement minimally** (GREEN phase)
2. **Run tests** after each change
3. **Refactor** while keeping tests green
4. **Update specs** if behavior differs

### Phase 4: Commit (Checkpoint)

**CRITICAL**: Every completed task requires commit before proceeding.

```
Tests Pass → Stage → Commit → Update Progress → Next Task
```

**Enforcement**: Use `sdd-checklist` skill (no automated hooks).

**Command**: After each task, follow the checklist:
1. Run tests: `go test ./...`
2. Stage changes: `git add {files}`
3. Commit with message
4. Update `task.md` progress
5. Update `DEVELOPMENT.md`

---

## Claude Code Hooks

Claude Code **does** support hooks - they're shell commands configured via `/hooks` or `settings.json`.

### How Hooks Work

| Component | Description |
|-----------|-------------|
| **Events** | PreToolUse, PostToolUse, Stop, Notification, etc. |
| **Matchers** | Filter which tools trigger the hook (e.g., `Bash`, `Edit|Write`) |
| **Commands** | Shell commands that receive JSON stdin, control via exit codes |

### Our Project Hooks

| Hook | Event | Matcher | Purpose |
|------|-------|---------|---------|
| `pre_commit.py` | PreToolUse | Bash | Reminds about SDD checklist on git commit |
| `markdown_formatter.py` | PostToolUse | Edit\|Write | Auto-formats markdown files |
| `sdd_reminder.py` | Stop | (all) | Reminds about post-task workflow |

### Setup

Run `/hooks` command to configure, or add to `.claude/settings.local.json`:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [{"type": "command", "command": "python3 .claude/hooks/pre_commit.py"}]
      }
    ]
  }
}
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success, continue |
| 2 | Block with feedback (stdout shown to Claude) |

### Backup: TodoWrite Enforcement

Use TodoWrite to track workflow steps:

```javascript
[
  {"content": "Complete Step B01 implementation", "status": "completed"},
  {"content": "Run tests", "status": "completed"},
  {"content": "Commit B01 changes", "status": "completed"},
  {"content": "Update task.md progress", "status": "completed"},
  {"content": "Start Step B02", "status": "in_progress"}
]
```

**Red Flags** (stop and commit first):
- "Now let's start Step B02..."
- "Ready for the next step?"
- "Step complete! Should I continue?"

---

## Project Structure

### Documentation Hierarchy

```
Repository Root
├── CLAUDE.md                    # Always included (minimal spine ~100 lines)
│
├── docs/
│   ├── spec-driven-design.md    # This file - methodology
│   ├── DESIGN.md                # System architecture (reference)
│   ├── USER_STORIES.md          # All stories (reference)
│   └── TASKS.md                 # All tasks (reference)
│
├── services/{service}/
│   └── SPEC.md                  # Technical specification (source of truth)
│
├── .work-items/{feature}/       # Active work tracking
│   ├── user-story.md            # Feature requirements
│   ├── design.md                # Technical design
│   ├── task.md                  # Implementation tasks
│   └── adr-{NNN}.md             # Architecture decisions
│
└── .claude/
    ├── DEVELOPMENT.md           # Session context (read first)
    ├── skills/                  # Knowledge guides (categorized)
    ├── commands/                # Custom slash commands
    ├── hooks/                   # Documentation only (not executed)
    └── plans/
        └── active -> ...        # Symlink to current work
```

### Source of Truth Priority

1. **`SPEC.md`** - Technical contract
2. **Work Item Files** - Feature-specific details
3. **`DEVELOPMENT.md`** - Current session state
4. **Reference Docs** - Background context

---

## Artifact Standards

### User Story Format (EARS)

```markdown
# User Story: {Feature Name}

## Story

**As a** {role}
**I want** {capability}
**So that** {benefit}

## Acceptance Criteria

### AC-1: {Criterion Name}

**Given** {precondition}
**When** {action}
**Then** {expected result}

## Out of Scope
- {Explicit exclusions}

## Dependencies
- {Required prerequisites}
```

### Design Document Format

```markdown
# Design Document: {Feature Name}

## Overview
{Brief description of approach}

## Architecture Decisions

### AD-1: {Decision Title}

**Decision:** {What was decided}
**Rationale:** {Why this approach}
**Alternatives Considered:**
- {Option A} - {Why rejected}

## Component Design
{Technical details, data structures, interfaces}

## Testing Strategy
{How to verify the design}
```

### Task Breakdown Format

```markdown
# Task Breakdown: {Feature Name}

## Progress Tracker

| Step | Task | Hours | Status | Actual |
|------|------|-------|--------|--------|
| 01 | {Task name} | 2.0 | DONE | 1.5 |
| 02 | {Task name} | 1.5 | IN_PROGRESS | - |

## Current Step: {NN}

### Step {NN}: {Task Name}

**Estimated**: {hours}h
**Files**: `path/to/file.go`

**Tasks**:
- [ ] Subtask 1
- [ ] Subtask 2

**Done when**: {Testable criterion}
```

---

## Commands (Slash Commands)

Custom commands in `.claude/commands/`:

| Command | Purpose | Triggers |
|---------|---------|----------|
| `/start-feature {name}` | Initialize feature context | Creates symlink, loads context |
| `/complete-feature {name}` | Verify and close feature | Checks completion, removes symlink |
| `/implement {task-id}` | TDD implementation | Follows RED→GREEN→REFACTOR |
| `/status` | Current project state | Reads DEVELOPMENT.md |
| `/next-task` | Find next work item | Scans task.md files |
| `/spec {service}` | Display service spec | Reads SPEC.md |

---

## Model Selection

Different models for different tasks:

| Task Type | Model | Rationale |
|-----------|-------|-----------|
| Spec writing, architecture | **Opus** | Complex reasoning |
| Well-specified implementation | **Haiku** | Cost-effective |
| Complex implementation | **Sonnet** | Balanced capability |
| Bug investigation | **Sonnet** | Analytical reasoning |
| Documentation | **Haiku** | Straightforward |

See `model-selection-guide` skill for detailed guidance.

---

## Red Flags and Anti-Patterns

### Red Flags (Stop Immediately)

| Flag | Action |
|------|--------|
| Moving to next task without commit | STOP. Commit first. |
| Writing code without reading spec | STOP. Read spec first. |
| Tests written after implementation | Reorder. Tests first. |
| "We can add tests later" | No. Tests now. |

### Anti-Patterns

| Anti-Pattern | Better Approach |
|--------------|-----------------|
| Huge multi-task commits | One commit per task |
| "Quick fix" without tests | Write failing test first |
| Spec says X, code does Y | Align spec and code |
| Skipping design docs | Document decisions briefly |

---

## MCP Integration (Future)

Model Context Protocol servers could enhance SDD:

| MCP Server | Potential Use |
|------------|---------------|
| `mcp-git` | Automated commit validation |
| `mcp-test-runner` | Test execution with reporting |
| `mcp-lint` | Pre-commit lint enforcement |
| `mcp-spec-validator` | Spec completeness checking |

Not currently implemented. Requirements documented for future development.

---

## Quick Reference

### Session Start

1. Read `.claude/DEVELOPMENT.md`
2. Check `.claude/plans/active`
3. Review current task status
4. Continue TDD workflow

### Task Completion

1. Tests pass
2. Git commit
3. Update task.md
4. Update DEVELOPMENT.md
5. Next task

### Feature Lifecycle

```
/start-feature {name}
    ↓
Read user-story.md + design.md + task.md
    ↓
For each task:
    RED → GREEN → REFACTOR → COMMIT
    ↓
/complete-feature {name}
```

---

## Related Documents

| Document | Location | Purpose |
|----------|----------|---------|
| Claude integration | `.claude/README.md` | How Claude Code works here |
| Session context | `.claude/DEVELOPMENT.md` | Current state |
| Commit checklist | `.claude/skills/sdd-checklist.md` | Post-task enforcement |
| SDD methodology | `.claude/skills/spec-driven-development.md` | Detailed principles |

---

## Summary

**SDD in Claude Code = Specifications + Skills + Discipline**

1. **Specifications** define behavior before code
2. **Skills** provide context-appropriate guidance (process/standards/guidelines)
3. **Discipline** (via TodoWrite + sdd-checklist) ensures commit hygiene

The methodology works without automated hooks by encoding expectations in skills and relying on structured workflow discipline.

**Spec first. Tests first. Commit always.**
