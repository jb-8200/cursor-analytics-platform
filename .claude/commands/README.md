# Custom Slash Commands

This directory contains custom slash commands for Claude Code.

## Available Commands

### /next-task
**Usage:** `/next-task`

Shows what task you should work on next based on:
- Current project phase from DEVELOPMENT.md
- Task dependencies from TASKS.md
- Model recommendations from model-selection-guide.md

**Output:**
- Task ID and name
- Service it belongs to
- Recommended model
- Key files to read

---

### /spec
**Usage:** `/spec [service-name]`

Displays service specifications.

**Examples:**
- `/spec cursor-sim` - Shows cursor-sim SPEC.md
- `/spec cursor-analytics-core` - Shows core SPEC.md
- `/spec` - Lists all available services

---

### /implement
**Usage:** `/implement TASK-ID`

Implements a task following Test-Driven Development (TDD).

**Examples:**
- `/implement TASK-SIM-001` - Implements Go project structure
- `/implement TASK-SIM-003` - Implements Developer generator

**Process:**
1. Reads task from docs/TASKS.md
2. Checks model recommendation
3. Reads relevant SPEC.md sections
4. Follows TDD: RED → GREEN → REFACTOR

---

### /status
**Usage:** `/status`

Shows current project status including:
- Implementation percentage
- Active work area
- Next recommended task
- Recent commits
- Uncommitted changes

---

### /model
**Usage:** `/model TASK-ID` or `/model [description]`

Recommends which Claude model to use for a task.

**Examples:**
- `/model TASK-SIM-003` - Shows: Haiku ⚡ (well-specified struct)
- `/model TASK-SIM-004` - Shows: Sonnet ⚡⚡ (complex Poisson logic)

**Output:**
- Recommended model with reasoning
- Estimated cost
- Alternative options

---

## How Commands Work

Each `.md` file in this directory becomes a slash command. The filename (without .md) is the command name.

The file content is a prompt that tells Claude what to do when you type that command.

## Command Format

```markdown
Brief description of what to do.

When user provides X:
- Step 1
- Step 2
- Step 3

Show:
- Output format
```

## Using with Model Selection

Combine commands with explicit model choice:

```bash
/next-task
# Shows: TASK-SIM-003, Recommended: Haiku

# Then implement with recommended model:
/implement TASK-SIM-003
# Claude asks: "Use Haiku as recommended? [Y/n]"
```

Or override:

```bash
# Force Sonnet even if Haiku recommended
"Use Sonnet to /implement TASK-SIM-003"
```

## Tips

1. **Start each session:** `/status` to see where you are
2. **Find next work:** `/next-task` for recommendations
3. **Check specs:** `/spec cursor-sim` before implementing
4. **Get model advice:** `/model TASK-ID` for cost optimization
5. **Implement:** `/implement TASK-ID` to build with TDD

## Integration with Skills

Commands automatically use skills in `.claude/skills/` (categorized):

**Process** (`.claude/skills/process/`):
- `spec-process-core.md` - Core SDD principles
- `spec-process-dev.md` - TDD workflow

**Standards** (`.claude/skills/standards/`):
- `spec-user-story.md` - User story format
- `spec-design.md` - Design doc format
- `spec-tasks.md` - Task breakdown format

**Guidelines** (`.claude/skills/guidelines/`):
- `go-best-practices.md` - Go standards
- `cursor-api-patterns.md` - API patterns

**Operational** (`.claude/skills/operational/`):
- `sdd-checklist.md` - **CRITICAL**: Post-task commit enforcement
- `model-selection-guide.md` - Model recommendations

---

**Note:** If a command isn't showing up, make sure:
1. The `.md` file is in `.claude/commands/`
2. The filename matches the command (no spaces, kebab-case)
3. Claude Code has reloaded (restart if needed)
