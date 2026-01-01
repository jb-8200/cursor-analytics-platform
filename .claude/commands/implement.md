Implement the specified task following Test-Driven Development.

When I provide a task ID (e.g., TASK-SIM-003):

1. Read the task details from docs/TASKS.md
2. Check .claude/skills/model-selection-guide.md for recommended model
3. Read relevant SPEC.md sections
4. Ask if I want to use the recommended model or override

Then follow TDD workflow:
- RED: Write failing tests first
- GREEN: Implement just enough to pass
- REFACTOR: Clean up code
- Update documentation

Use skills:
- .claude/skills/go-best-practices.md for Go code
- .claude/skills/cursor-api-patterns.md for API endpoints
- .claude/skills/spec-driven-development.md for TDD approach
