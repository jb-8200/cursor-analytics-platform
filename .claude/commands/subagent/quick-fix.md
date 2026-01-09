---
description: Spawn quick-fix agent for small, independent tasks using Haiku model
argument-hint: "[description of fix]"
allowed-tools: Task
---

# Spawn quick-fix Agent

Delegate a small, independent task to the quick-fix agent (Haiku model).

**Task**: $ARGUMENTS

## Objective

Make the specified quick fix with minimal overhead.

## Ideal For

- Typo fixes in code or documentation
- Simple configuration changes
- Adding missing imports
- Fixing linting errors
- Updating version numbers
- Renaming single variables
- Adding simple comments
- Fixing broken links

## NOT Ideal For

- New feature implementation
- Multi-file refactoring
- Tasks requiring new tests
- Architectural decisions
- Cross-service changes
- Tasks with dependencies

## Guidelines

- Make minimal change needed
- Don't over-engineer
- Skip exploration if target is clear
- Commit quickly and return

## Completion Report

```
FIXED: {what was fixed}
File: {path}
Commit: {hash}
```

If task is too complex:

```
ESCALATE: Task requires full agent
Reason: {why this isn't a quick fix}
Recommendation: Use {data-tier-dev|streamlit-dev|cursor-sim-cli-dev}
```

## Example Usage

```bash
/subagent/quick-fix "Fix typo in README.md line 42: 'teh' -> 'the'"
/subagent/quick-fix "Update port in docker-compose.yml from 3000 to 3001"
/subagent/quick-fix "Add missing import os to tools/api-loader/loader.py"
```

---

See also:
- `.claude/agents/quick-fix.md` - Full agent definition
