Start a feature by creating the active plan symlink and loading context.

When I provide a feature name (e.g., `/start-feature cursor-sim-v2`):

1. **Validate** the feature exists:
   - Check `.work-items/{feature-name}/` directory exists
   - Verify `user-story.md`, `design.md`, and `task.md` are present
   - Report any missing files

2. **Create active plan symlink**:
   ```bash
   cd .claude/plans
   ln -sf ../../.work-items/{feature-name}/task.md active
   ```

3. **Commit the change**:
   ```bash
   git add .claude/plans/active
   git commit -m "feat: start feature {feature-name}

   Activated work item: .work-items/{feature-name}/

   Generated with Claude Code"
   ```

4. **Load context**:
   - Read `user-story.md` to understand requirements
   - Read `design.md` to understand technical approach
   - Read `task.md` to see step breakdown

5. **Display first step**:
   - Show the first step from task.md that is NOT_STARTED
   - Summarize what needs to be done

## Output Format

```
=== Starting Feature: {feature-name} ===

User Story Summary:
{brief summary of user-story.md}

Technical Approach:
{brief summary of design.md}

Current Progress: {N}/{total} steps complete

Next Step: {step-name}
{step description}

TDD Reminder:
1. Write failing test first (RED)
2. Implement minimal code (GREEN)
3. Refactor while green
4. Commit with time tracking

Ready to begin? Acknowledge to start Step {N}.
```

## Error Handling

- If feature directory doesn't exist: List available features from `.work-items/`
- If already active: Show current active feature and ask to `/complete-feature` first
- If files missing: Report which files are missing
