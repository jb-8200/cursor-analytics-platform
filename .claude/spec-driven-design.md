# Spec-Driven Design (SDD) for Claude Code

**Version**: 2.1
**Last Updated**: January 3, 2026
**Adapted from**: [genai-specs](https://github.com/betsalel-williamson/genai-specs) methodology

---

## Overview

Spec-Driven Design (SDD) is a development methodology where **specifications drive all implementation**. This document adapts the genai-specs model for Claude Code, leveraging Skills, Commands, and structured workflows instead of Cursor-specific features.

### Core Flow (Enhanced)

```
Specifications → Tests → Implementation → Refactor → Reflect → Sync → Commit
     ↑                                                            │
     └────────────────────────────────────────────────────────────┘
         (Bidirectional: code informs specs via REFLECT/SYNC)
```

**Enhanced 7-Step Workflow**:
1. **SPEC** → Read specification before coding
2. **TEST** → Write failing tests (RED)
3. **CODE** → Minimal implementation (GREEN)
4. **REFACTOR** → Clean up while tests pass
5. **REFLECT** → Check dependency reflections (NEW)
6. **SYNC** → Update SPEC.md if triggered (NEW)
7. **COMMIT** → Checkpoint with code + docs

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

| Category | Skills | Trigger Examples |
|----------|--------|------------------|
| **Process** | spec-process-core, spec-process-dev | "start feature", "implement", "TDD" |
| **Standards** | spec-user-story, spec-design, spec-tasks | "write user story", "create design doc" |
| **Guidelines** | go-best-practices, cursor-api-patterns | Language/tech-specific coding |
| **Operational** | sdd-checklist, spec-sync-check, dependency-reflection, model-selection-guide | "commit", "update SPEC", "check reflections", "which model" |

Skills activate automatically when your request matches their description. You can also reference explicitly:

```
"Help me create a user story for the login feature"
  → spec-user-story activates automatically

"Implement the HTTP handler"
  → go-best-practices activates automatically
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

Each skill has its own directory with a `SKILL.md` file containing YAML frontmatter:

```
.claude/skills/
├── go-best-practices/          # Go patterns (cursor-sim)
│   └── SKILL.md
├── cursor-api-patterns/        # Cursor API implementation
│   └── SKILL.md
├── sdd-checklist/              # Post-task commit enforcement (enhanced with REFLECT/SYNC)
│   └── SKILL.md
├── spec-sync-check/            # SPEC.md update trigger detection (NEW)
│   └── SKILL.md
├── dependency-reflection/      # Dependency reflection checking (NEW)
│   └── SKILL.md
├── model-selection-guide/      # Model optimization
│   └── SKILL.md
├── spec-process-core/          # Core SDD principles (enhanced 7-step workflow)
│   └── SKILL.md
├── spec-process-dev/           # TDD development workflow
│   └── SKILL.md
├── spec-user-story/            # User story format (EARS)
│   └── SKILL.md
├── spec-design/                # Design doc format (ADRs)
│   └── SKILL.md
├── spec-tasks/                 # Task breakdown format
│   └── SKILL.md
├── sdd-workflow/               # Full workflow reference
│   └── SKILL.md
└── spec-driven-development/    # SDD methodology details
    └── SKILL.md
```

### SKILL.md Format

Each skill requires YAML frontmatter with `name` and `description`:

```markdown
---
name: spec-user-story
description: Create or revise user stories following EARS format. Use when
             writing PRDs, feature narratives, acceptance criteria, or
             product requirements. Produces user-story.md files.
---

# User Story Standard

## Process
1. Gather requirements from stakeholder
2. Draft story using template
3. Review acceptance criteria
4. Create file in .work-items/{feature}/user-story.md

## Template

[Template content here]
```

### How Skills Are Discovered

1. **At startup**: Claude loads skill names and descriptions
2. **On request**: Claude matches your request against descriptions
3. **Activation**: You approve, Claude loads full skill content

### Skill Categories

| Category | Skills | Purpose |
|----------|--------|---------|
| **Process** | spec-process-core, spec-process-dev | Workflow guidance |
| **Standards** | spec-user-story, spec-design, spec-tasks | Artifact templates |
| **Guidelines** | go-best-practices, cursor-api-patterns | Tech-specific patterns |
| **Operational** | sdd-checklist, spec-sync-check, dependency-reflection, model-selection-guide | Day-to-day workflow, commit hygiene |

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

### Phase 4: Reflect (Dependency Checking)

**NEW**: Before committing, check if changes require updates to related files.

**Purpose**: Detect when changes in one file require updates to related files (documentation, tests, or other code).

**What to Check**:

1. **Documentation Reflections**
   - Did I modify models? → Check SPEC.md schema section
   - Did I add/modify endpoints? → Check SPEC.md endpoints table
   - Did I complete a phase step? → Check SPEC.md implementation status
   - Did I add packages? → Check SPEC.md package structure

2. **Code Synchronization Reflections**
   - Did I modify models? → Check generators that create these models
   - Did I change storage interface? → Check all handlers using storage
   - Did I add new events? → Check generators that produce them
   - Did I modify config? → Check CLI that consumes it

3. **Test Synchronization Reflections**
   - Did I add new code paths? → Write tests for them
   - Did I modify behavior? → Update test assertions
   - Did I add endpoints? → Write E2E tests
   - Did I change models? → Update model tests

**Regression Testing Protocol**:

After identifying reflections, run appropriate tests:

```bash
# After model changes
go test ./... -cover

# After handler changes
go test ./internal/api/... -v
go test ./test/e2e/... -v

# After storage interface changes
go test ./internal/... -v -race

# After major refactor
go test ./... -v -race -count=5
```

**Skill**: Use `dependency-reflection` for detailed guidance.

**Detection Matrix Example**:

| Files Changed | Check These Files | Type |
|--------------|-------------------|------|
| `internal/models/*.go` | Generators, handlers, SPEC.md, tests | All 3 types |
| `internal/api/**/*.go` | SPEC.md endpoints, E2E tests | Documentation + Test |
| `internal/storage/*.go` | All handlers using storage, tests | Code + Test |

### Phase 5: Sync (SPEC.md Updates)

**NEW**: Update service specifications when implementation changes warrant it.

**Purpose**: Keep `services/{service}/SPEC.md` synchronized with actual implementation.

**Trigger Detection**:

Run `spec-sync-check` to determine if SPEC.md needs updating:

**High-Priority Triggers (MUST update)**:

1. **Phase Completion**
   - Completed final step of any phase (e.g., C06, B08)
   - Update: Implementation Status table, Phase Features section

2. **New Endpoint Added**
   - Created/modified handler in `internal/api/`
   - Update: Endpoints table with method, path, auth, status

3. **New Service/Package Created**
   - Created directory in `internal/`
   - Update: Package Structure section

4. **CLI Configuration Changes**
   - Modified `internal/config/` or `cmd/`
   - Update: CLI Configuration, environment variables

**Medium-Priority Triggers (SHOULD update)**:

1. **Model Changes**
   - Modified struct fields in `internal/models/`
   - Update: Response Format section, schema examples

2. **Generator Changes**
   - New/modified generators in `internal/generator/`
   - Update: Generation Algorithm section

**Update Checklist**:

When updating `services/{service}/SPEC.md`:

- [ ] Line 5: **Last Updated** date is today
- [ ] Lines 17-22: **Implementation Status** reflects current phase
- [ ] **Endpoints tables** include all new/modified endpoints
- [ ] **Schema examples** match actual struct fields
- [ ] **Phase features** marked as "Implemented ✅" when complete
- [ ] **Package Structure** includes new directories
- [ ] **Decision Log** includes significant technical decisions

**Critical Rule**: If SPEC.md is updated, include it in the same commit as code changes.

**Skill**: Use `spec-sync-check` for detailed guidance.

### Phase 6: Commit (Checkpoint)

**CRITICAL**: Every completed task requires commit before proceeding.

**Enhanced Commit Sequence**:

```
Tests Pass → REFLECT → SYNC → Stage → Commit → Update Progress → Next Task
```

**Enforcement**: Use `sdd-checklist` skill.

**7-Step Checklist**:
1. ✅ **Tests pass**: `go test ./...`
2. ✅ **Check reflections**: Run `dependency-reflection` check
3. ✅ **Update SPEC.md**: Run `spec-sync-check` if triggered
4. ✅ **Stage changes**: `git add {files}` (include SPEC.md if updated)
5. ✅ **Git commit**: Descriptive message
6. ✅ **Update task.md**: Mark step as DONE
7. ✅ **Update DEVELOPMENT.md**: Current status

**Never** move to next task before completing all 7 steps.

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
├── .work-items/{P#-F##-feature}/  # Active work tracking
│   ├── user-story.md              # Feature requirements
│   ├── design.md                  # Technical design
│   ├── task.md                    # Implementation tasks (TASK01, TASK02...)
│   └── adr-{NNN}.md               # Architecture decisions
│
│   Hierarchy: Phase (P#) → Feature (F##) → Task (TASK##)
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

**Feature ID**: P#-F##-feature-name
**Phase**: P# (Phase Name)
**Status**: IN PROGRESS

## Progress Tracker

| Task ID | Task | Hours | Status | Actual |
|---------|------|-------|--------|--------|
| TASK01 | {Task name} | 2.0 | DONE | 1.5 |
| TASK02 | {Task name} | 1.5 | IN_PROGRESS | - |

## Current Task: TASK##

### TASK##: {Task Name}

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
| **Committing without checking reflections** | **STOP. Run dependency-reflection first.** |
| **Adding endpoint without updating SPEC.md** | **STOP. Run spec-sync-check first.** |
| **Completing phase without updating status** | **STOP. Update SPEC.md Implementation Status.** |

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

### Task Completion (Enhanced)

1. ✅ Tests pass
2. ✅ Check reflections (dependency-reflection)
3. ✅ Update SPEC.md if triggered (spec-sync-check)
4. ✅ Git commit (code + SPEC.md if updated)
5. ✅ Update task.md
6. ✅ Update DEVELOPMENT.md
7. ✅ Next task

### Feature Lifecycle (Enhanced)

```
/start-feature {name}
    ↓
Read user-story.md + design.md + task.md
    ↓
For each task:
    RED → GREEN → REFACTOR → REFLECT → SYNC → COMMIT
    ↓
/complete-feature {name}
```

**7-Step Per-Task Cycle**:
1. RED - Write failing test
2. GREEN - Minimal implementation
3. REFACTOR - Clean up code
4. REFLECT - Check dependency reflections
5. SYNC - Update SPEC.md if triggered
6. COMMIT - Checkpoint with code + docs
7. REPEAT - Next task

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

**Enhanced 7-Step Workflow** (v2.1):
1. SPEC → Read specification before coding
2. TEST → Write failing tests (RED)
3. CODE → Minimal implementation (GREEN)
4. REFACTOR → Clean up while tests pass
5. **REFLECT** → Check dependency reflections (NEW)
6. **SYNC** → Update SPEC.md if triggered (NEW)
7. COMMIT → Checkpoint with code + docs

The methodology works without automated hooks by encoding expectations in skills and relying on structured workflow discipline.

**Key Enhancements**:
- **REFLECT** prevents documentation drift and missing test updates
- **SYNC** keeps SPEC.md synchronized with implementation
- **Commit hygiene** now includes code + documentation updates together

**Spec first. Tests first. Reflect always. Sync when triggered. Commit with docs.**
