# Design Document: Hooks Configuration

**Feature ID**: P0-F05
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

Configure existing hook scripts in `settings.local.json` so they run automatically.

---

## Settings Configuration

Update `.claude/settings.local.json`:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [{
          "type": "command",
          "command": "python3 \"$CLAUDE_PROJECT_DIR/.claude/hooks/pre_commit.py\""
        }]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [{
          "type": "command",
          "command": "python3 \"$CLAUDE_PROJECT_DIR/.claude/hooks/markdown_formatter.py\""
        }]
      }
    ],
    "Stop": [
      {
        "matcher": "",
        "hooks": [{
          "type": "command",
          "command": "python3 \"$CLAUDE_PROJECT_DIR/.claude/hooks/sdd_reminder.py\""
        }]
      }
    ]
  }
}
```

---

## Hook Behavior

### PreToolUse: pre_commit.py

**Trigger**: Before any Bash command
**Matcher**: `Bash`
**Action**: If command contains `git commit`, show SDD reminder
**Exit Code**: 0 (allow with reminder)

### PostToolUse: markdown_formatter.py

**Trigger**: After Edit or Write
**Matcher**: `Edit|Write`
**Action**: If file is .md/.mdx, auto-format
**Exit Code**: 0 (success)

### Stop: sdd_reminder.py

**Trigger**: When Claude finishes responding
**Matcher**: "" (all stops)
**Action**: Show SDD checklist reminder
**Exit Code**: 0 (informational)

---

## Script Updates

May need to update scripts for robustness:
- Better error handling
- Cleaner output formatting
- Proper JSON parsing
