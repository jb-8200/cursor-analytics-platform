#!/usr/bin/env python3
import json
import sys

data = json.load(sys.stdin)
# Check if response mentions task completion
# Could analyze stop_reason or transcript

print("---")
print("SDD Checklist: If task complete, remember to:")
print("1. Run tests  2. Check reflections  3. Git commit  4. Update task.md  5. Update DEVELOPMENT.md")
sys.exit(0)
