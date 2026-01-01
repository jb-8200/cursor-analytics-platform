# Claude Code Slash Commands

This directory contains custom slash commands that streamline common development tasks in the Cursor Analytics Platform project. These commands provide shortcuts for spec-driven development and TDD workflows.

## Available Commands

### /spec

**Purpose**: Displays the specification for a given service or feature.

**Usage**:
```
/spec cursor-sim
/spec cursor-analytics-core  
/spec cursor-viz-spa
/spec US-CORE-003
/spec TASK-SIM-004
```

**Behavior**: 

When given a service name, this command reads and displays the contents of that service's SPEC.md file. This helps Claude and developers quickly reference technical specifications without leaving the conversation.

When given a user story ID (US-XXX-NNN format), the command searches USER_STORIES.md and displays the matching story with its acceptance criteria.

When given a task ID (TASK-XXX-NNN format), the command searches TASKS.md and displays the task details including dependencies and definition of done.

**Implementation Notes**:

Claude should implement this by using file reading tools to fetch the appropriate specification file based on the argument pattern. Service names map to `services/{name}/SPEC.md`. User story IDs require searching `docs/USER_STORIES.md`. Task IDs require searching `docs/TASKS.md`.

---

### /test

**Purpose**: Guides the creation of tests based on specifications.

**Usage**:
```
/test US-CORE-003
/test acceptance-rate
/test cursor-sim developer-generation
```

**Behavior**:

When given a user story ID, this command reads the acceptance criteria and generates skeleton test cases for each scenario. The skeletons include the test structure with comments indicating what each section should test.

When given a feature name and optional service, the command locates relevant specifications and generates appropriate test cases.

**Output Example**:

```typescript
// Generated test cases for US-CORE-003: Calculate Acceptance Rate

describe('calculateAcceptanceRate', () => {
    // Scenario 1: Basic calculation
    // Given: developer has 100 shown and 75 accepted
    // Then: rate should be 75.0
    it('should calculate rate as (accepted/shown)*100', async () => {
        // TODO: Implement test
        // Arrange: Create developer with specific event counts
        // Act: Call calculateAcceptanceRate
        // Assert: Verify result is 75.0
    });

    // Scenario 2: No suggestions shown
    // Given: developer has 0 shown
    // Then: rate should be null
    it('should return null when no suggestions shown', async () => {
        // TODO: Implement test
    });

    // Continue for remaining scenarios...
});
```

---

### /implement

**Purpose**: Guides implementation following the TDD workflow.

**Usage**:
```
/implement TASK-SIM-003
/implement developer-generation cursor-sim
```

**Behavior**:

This command orchestrates the full implementation workflow. It starts by reading the relevant task specification and checking dependencies. Then it identifies or creates the test file and helps write failing tests. Finally, it guides the implementation to make tests pass.

The command provides step-by-step guidance rather than generating all code at once, allowing for iterative development.

**Workflow Steps**:

1. Read task specification from TASKS.md
2. Check dependency status
3. Locate or create test file
4. Generate test cases from acceptance criteria
5. Confirm tests fail
6. Guide minimal implementation
7. Verify tests pass
8. Suggest refactoring opportunities

---

### /coverage

**Purpose**: Analyzes test coverage for a service and identifies gaps.

**Usage**:
```
/coverage cursor-sim
/coverage cursor-analytics-core
/coverage all
```

**Behavior**:

Runs the test suite with coverage enabled and analyzes the results. Reports the overall coverage percentage and lists files or functions below the coverage threshold. Suggests which additional tests would most improve coverage.

---

### /validate

**Purpose**: Validates that specifications, tests, and implementation are aligned.

**Usage**:
```
/validate cursor-sim
/validate US-CORE-003
/validate all
```

**Behavior**:

Checks that specifications have corresponding tests, tests have corresponding implementations, and implementations match specification requirements. Reports any misalignments or gaps.

**Checks Performed**:

For specifications: Verifies test files exist for each acceptance criterion.

For tests: Verifies implementation files exist for tested functionality.

For implementations: Verifies behavior matches specification requirements through static analysis where possible.

---

### /status

**Purpose**: Shows the current project status and progress.

**Usage**:
```
/status
/status phase1
/status cursor-sim
```

**Behavior**:

Displays an overview of completed, in-progress, and pending tasks. When given a phase or service filter, shows only relevant items.

**Output Example**:

```
Project Status: Cursor Analytics Platform

Phase 1 Progress: 35% complete (7/20 tasks)

cursor-sim:
  ✅ TASK-SIM-001: Initialize Go Project Structure
  ✅ TASK-SIM-002: Implement CLI Flag Parsing
  🔄 TASK-SIM-003: Implement Developer Profile Generator
  ⬜ TASK-SIM-004: Implement Event Generation Engine
  ⬜ TASK-SIM-005: Implement In-Memory Storage
  
cursor-analytics-core:
  ✅ TASK-CORE-001: Initialize TypeScript Project Structure
  ⬜ TASK-CORE-002: Define Database Schema
  ...

Next recommended task: TASK-SIM-004 (depends on SIM-003)
```

---

### /scaffold

**Purpose**: Creates boilerplate files for new features or components.

**Usage**:
```
/scaffold component TeamRadarChart cursor-viz-spa
/scaffold resolver developer cursor-analytics-core
/scaffold handler activity cursor-sim
```

**Behavior**:

Creates the file structure and boilerplate code for common patterns. Includes test files alongside implementation files. Uses project conventions for file naming and directory structure.

**Supported Scaffolds**:

For cursor-sim (Go):
- `handler` - HTTP handler with test file
- `generator` - Data generator with test file
- `model` - Domain type definition

For cursor-analytics-core (TypeScript):
- `resolver` - GraphQL resolver with test file
- `service` - Business logic service with test file
- `migration` - Database migration file

For cursor-viz-spa (React):
- `component` - React component with test file and story
- `hook` - Custom hook with test file
- `page` - Page component with routing setup

---

## Command Implementation Guide

These commands are conceptual specifications for Claude's behavior. When a user types a slash command, Claude should recognize the pattern and execute the described behavior using available tools like file reading, searching, and code generation.

### Recognition Pattern

Commands are recognized by the leading slash followed by the command name. Arguments follow separated by spaces. Service names, task IDs, and user story IDs have distinct patterns that help Claude route to the correct behavior.

### Error Handling

If a command cannot be executed due to missing files or invalid arguments, Claude should explain what went wrong and suggest corrections. For example, if `/spec US-FAKE-001` finds no matching story, Claude should list available stories or suggest checking the USER_STORIES.md file.

### Context Awareness

Commands should be context-aware. If the user has been working on a specific service, commands without explicit service arguments should default to that service. The conversation history provides this context.
