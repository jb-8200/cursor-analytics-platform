# Claude Code Hooks

## How Hooks Actually Work

Hooks are **shell commands** configured via `/hooks` command or `settings.json`. They:
- Run at specific lifecycle events (PreToolUse, PostToolUse, etc.)
- Receive JSON input via stdin
- Control behavior via exit codes (0=success, 2=block with feedback)

## Available Hook Events

| Event | When it Runs | Use Case |
|-------|--------------|----------|
| `PreToolUse` | Before tool calls | Block/validate commands |
| `PostToolUse` | After tool calls | Format files, log actions |
| `UserPromptSubmit` | When user submits prompt | Inject context |
| `Notification` | When Claude sends notifications | Custom alerts |
| `Stop` | When Claude finishes responding | Post-response actions |

## Our Project Hooks

### 1. Pre-Commit Test Enforcement

Blocks `git commit` if tests haven't been run recently.

**Configure via `/hooks`:**
- Event: `PreToolUse`
- Matcher: `Bash`
- Command: `"$CLAUDE_PROJECT_DIR"/.claude/hooks/pre_commit.py`

### 2. Markdown Formatter

Auto-formats markdown files after editing.

**Configure via `/hooks`:**
- Event: `PostToolUse`
- Matcher: `Edit|Write`
- Command: `"$CLAUDE_PROJECT_DIR"/.claude/hooks/markdown_formatter.py`

### 3. SDD Context Reminder

Reminds about sdd-checklist after completing tasks.

**Configure via `/hooks`:**
- Event: `Stop`
- Matcher: (empty for all)
- Command: `"$CLAUDE_PROJECT_DIR"/.claude/hooks/sdd_reminder.py`

---

## Setup Instructions

### Option 1: Via `/hooks` Command (Recommended)

1. Run `/hooks` in Claude Code
2. Select the event (e.g., `PreToolUse`)
3. Add matcher (e.g., `Bash`)
4. Add hook command (e.g., path to script)
5. Choose storage location (User or Project)

### Option 2: Via settings.json

Add to `~/.claude/settings.json` or `.claude/settings.local.json`:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "python3 \"$CLAUDE_PROJECT_DIR/.claude/hooks/pre_commit.py\""
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "python3 \"$CLAUDE_PROJECT_DIR/.claude/hooks/markdown_formatter.py\""
          }
        ]
      }
    ]
  }
}
```

---

## Hook Scripts

### pre_commit.py

Validates git commits follow SDD (tests must pass first):

```python
#!/usr/bin/env python3
import json
import sys

data = json.load(sys.stdin)
command = data.get('tool_input', {}).get('command', '')

# Check if this is a git commit
if 'git commit' in command:
    # Output feedback message
    print("SDD Reminder: Have you run tests before committing?")
    print("Checklist: 1) go test ./... 2) All pass 3) Then commit")
    # Exit 0 to allow, exit 2 to block with feedback
    sys.exit(0)  # Allow but remind

sys.exit(0)
```

### sdd_reminder.py

Reminds about post-task workflow:

```python
#!/usr/bin/env python3
import json
import sys

data = json.load(sys.stdin)
# Check if response mentions task completion
# Could analyze stop_reason or transcript

print("---")
print("SDD Checklist: If task complete, remember to:")
print("1. Run tests  2. Git commit  3. Update task.md  4. Update DEVELOPMENT.md")
sys.exit(0)
```

---

## Hook Input/Output

### Input (JSON via stdin)

```json
{
  "session_id": "abc123",
  "tool_name": "Bash",
  "tool_input": {
    "command": "git commit -m 'feat: add feature'",
    "description": "Commit changes"
  }
}
```

### Output (Exit Codes)

| Exit Code | Meaning |
|-----------|---------|
| 0 | Success, continue |
| 2 | Block with feedback (stdout shown to Claude) |
| Other | Error (logged, doesn't block) |

---

## Security Considerations

⚠️ Hooks run with your environment's credentials. Always:
- Review hook code before enabling
- Don't expose secrets in hook output
- Be careful with file paths and command injection

---

## Debugging

1. Test hooks manually:
```bash
echo '{"tool_input":{"command":"git commit"}}' | python3 .claude/hooks/pre_commit.py
```

2. Check exit code:
```bash
echo $?
```

3. View logs in Claude Code output

---

## Current Status

| Hook | Script | Configured? |
|------|--------|-------------|
| Pre-commit validation | `pre_commit.py` | ✅ Configured |
| Markdown formatter | `markdown_formatter.py` | ✅ Configured |
| SDD reminder | `sdd_reminder.py` | ✅ Configured |

All hooks are configured in `.claude/settings.local.json` and ready to use.
