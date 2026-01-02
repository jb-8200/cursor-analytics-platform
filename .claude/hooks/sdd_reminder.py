#!/usr/bin/env python3
"""
SDD Reminder hook for Claude Code.

Configured as Stop hook to remind about post-task workflow.

Setup: /hooks → Stop → (empty matcher) → this script
"""
import json
import sys

def main():
    try:
        data = json.load(sys.stdin)
    except json.JSONDecodeError:
        sys.exit(0)

    # Could analyze stop_reason or session data
    # For now, just provide a gentle reminder

    print()
    print("─" * 40)
    print("SDD Reminder: If a task was completed:")
    print("  1. Run tests")
    print("  2. Git commit")
    print("  3. Update task.md")
    print("  4. Update DEVELOPMENT.md")
    print("─" * 40)

    sys.exit(0)

if __name__ == '__main__':
    main()
