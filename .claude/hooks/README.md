# Claude Code Hooks

## CRITICAL: Hooks Do Not Execute in Claude Code

**The Python hooks in this directory are documentation-only.** Claude Code does not execute custom Python scripts as workflow automation.

These files describe **intended behavior** that must be implemented through alternative means.

---

## Hook Files (Documentation Only)

| File | Intended Purpose | Status |
|------|-----------------|--------|
| `pre_prompt.py` | Inject context before each prompt | **NOT EXECUTED** |
| `pre_commit.py` | Run tests before git commit | **NOT EXECUTED** |
| `pre_patch.py` | Lint Markdown before applying patches | **NOT EXECUTED** |

---

## Alternative Implementations

Since hooks don't execute, use these alternatives:

### pre_prompt.py Intent: Context Injection

**What it was designed to do:**
- Read `.claude/plans/active` symlink
- Load rules from `rules/*.mdc`
- Include `DEVELOPMENT.md` context

**Claude Code Alternative:**

1. **Session start**: Read `.claude/DEVELOPMENT.md` manually
2. **Check active work**: `ls -la .claude/plans/active`
3. **Reference skills**: Include skill names in requests

```
"Following spec-process-core and go-best-practices, implement the handler"
```

### pre_commit.py Intent: Test Enforcement

**What it was designed to do:**
- Run test suite before commit
- Block commit if tests fail
- Report coverage

**Claude Code Alternative:**

Use the `sdd-checklist` skill which enforces:

1. Run tests: `go test ./...`
2. Verify all pass
3. Only then commit

The checklist is manual but explicit. Include test running in TodoWrite:

```javascript
[
  {"content": "Implement Step B03", "status": "completed"},
  {"content": "Run tests", "status": "in_progress"},  // <-- Explicit step
  {"content": "Commit changes", "status": "pending"}
]
```

### pre_patch.py Intent: Lint Enforcement

**What it was designed to do:**
- Run vale and markdownlint on Markdown files
- Add lint warnings as comments

**Claude Code Alternative:**

Run linters manually before committing documentation:

```bash
# If installed
vale docs/
markdownlint "**/*.md"

# Or use Go linting for code
golangci-lint run
```

---

## Why Hooks Don't Work

Claude Code is a CLI-based AI assistant that:
- Runs in a sandboxed environment
- Does not have a plugin/extension system
- Cannot execute arbitrary Python scripts as workflow hooks
- Uses a different architecture than Cursor IDE

The hook files exist because:
1. They document intended workflow automation
2. They could work if Claude Code adds hook support in the future
3. They serve as reference for manual enforcement

---

## Recommended Workflow Enforcement

### TodoWrite Pattern

Use TodoWrite to track enforcement steps explicitly:

```javascript
[
  {"content": "Implement feature code", "status": "completed"},
  {"content": "Write unit tests", "status": "completed"},
  {"content": "Run tests (go test ./...)", "status": "completed"},
  {"content": "Run linter (golangci-lint)", "status": "completed"},
  {"content": "Commit with message", "status": "completed"},
  {"content": "Update task.md progress", "status": "completed"},
  {"content": "Update DEVELOPMENT.md", "status": "completed"}
]
```

### sdd-checklist Skill

Reference the `sdd-checklist` skill after completing any task:

```
"Task complete. Following sdd-checklist, let me commit and update progress."
```

The skill contains the 5-step post-task process:
1. Tests pass
2. Git commit
3. Update task.md
4. Update DEVELOPMENT.md
5. Proceed to next task

---

## Future Possibilities

### If Claude Code Adds Hook Support

The current Python files follow a reasonable interface:

```python
def pre_prompt_hook(prompt: str, context: dict) -> str:
    # Modify prompt before processing
    return modified_prompt

def pre_commit_hook(context: dict) -> str:
    # Return empty string to allow, message to block
    return ""
```

If Claude Code implements hooks, these could be adapted.

### MCP (Model Context Protocol) Alternative

MCP servers could provide similar functionality:

| MCP Server | Equivalent Hook |
|------------|-----------------|
| `mcp-git-validator` | pre_commit.py |
| `mcp-test-runner` | pre_commit.py |
| `mcp-lint` | pre_patch.py |
| `mcp-context-loader` | pre_prompt.py |

MCP is not currently implemented but documented for future development.

---

## Summary

| Hook | Status | Alternative |
|------|--------|-------------|
| `pre_prompt.py` | Documentation only | Read DEVELOPMENT.md at session start |
| `pre_commit.py` | Documentation only | sdd-checklist skill + TodoWrite |
| `pre_patch.py` | Documentation only | Manual linter execution |

**The SDD workflow works through discipline, not automation.**

Skills + TodoWrite + explicit workflow steps replace automated hooks.
