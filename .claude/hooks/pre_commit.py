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
