Implement the specified task following Test-Driven Development.

When I provide a task ID or step (e.g., TASK-SIM-003, Step B02):

1. Read the task details from `.work-items/{feature}/task.md`
2. Read relevant `services/{service}/SPEC.md` sections
3. The `model-selection-guide` skill will help choose the right model
4. Ask if I want to use the recommended model or override

Then follow TDD workflow:
- RED: Write failing tests first
- GREEN: Implement just enough to pass
- REFACTOR: Clean up code

**CRITICAL: After tests pass, follow sdd-checklist:**
1. Run tests: `go test ./...`
2. Stage changes: `git add {files}`
3. Commit with descriptive message
4. Update task.md progress (Status → DONE, add Actual hours)
5. Update DEVELOPMENT.md
6. Only then proceed to next task

Skills activate automatically:
- `go-best-practices` for Go code
- `cursor-api-patterns` for API endpoints
- `spec-process-dev` for TDD workflow
- `sdd-checklist` for post-task commit
