#!/usr/bin/env python3
"""
Pre-commit hook for Claude Code.

Configured as PreToolUse hook with Bash matcher.
Reminds about SDD checklist when git commit is detected.

Setup: /hooks → PreToolUse → Bash → this script
"""
import json
import sys
import os

def main():
    try:
        data = json.load(sys.stdin)
    except json.JSONDecodeError:
        sys.exit(0)  # No input, allow

    tool_input = data.get('tool_input', {})
    command = tool_input.get('command', '')

    # Check if this is a git commit command
    if 'git commit' in command:
        # Check if tests were run recently (within last 5 minutes)
        # Look for test output files or just remind
        project_dir = os.environ.get('CLAUDE_PROJECT_DIR', '.')

        print("=" * 50)
        print("SDD Pre-Commit Checklist")
        print("=" * 50)
        print()
        print("Before committing, verify:")
        print("  1. ✅ Tests pass: go test ./...")
        print("  2. ✅ Coverage adequate")
        print("  3. ✅ task.md will be updated after commit")
        print("  4. ✅ DEVELOPMENT.md will be updated")
        print()
        print("Proceeding with commit...")
        print("=" * 50)

        # Exit 0 to allow (with reminder shown)
        # Exit 2 to block (feedback shown to Claude)
        sys.exit(0)

    # Not a commit, allow
    sys.exit(0)

if __name__ == '__main__':
    main()
