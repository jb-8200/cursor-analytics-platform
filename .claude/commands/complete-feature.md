Complete a feature by verifying all steps are done and removing the active plan symlink.

When I provide a feature name (e.g., `/complete-feature cursor-sim-v2`):

1. **Verify feature is active**:
   - Check `.claude/plans/active` symlink exists
   - Verify it points to the specified feature
   - Report error if different feature is active

2. **Verify all steps complete**:
   - Read `task.md` and check all steps are DONE
   - Report any incomplete steps
   - Block completion if steps remain

3. **Verify acceptance criteria**:
   - Read `user-story.md` acceptance criteria
   - For each AC, check if corresponding tests exist and pass
   - Report any uncovered ACs

4. **Run full test suite**:
   - Execute all tests for the service
   - Ensure no regressions
   - Report coverage

5. **Calculate time tracking**:
   - Sum estimated hours from task.md
   - Sum actual hours from commit messages
   - Calculate time saved/overrun

6. **Remove symlink and commit**:
   ```bash
   rm .claude/plans/active
   git add .claude/plans/
   git commit -m "feat: complete feature {feature-name}

   Summary:
   - Steps completed: {N}
   - Tests passed: {M}
   - Coverage: {X}%

   Time Tracking:
   - Estimated: {E} hours
   - Actual: {A} hours
   - Delta: {D} hours

   Generated with Claude Code"
   ```

7. **Archive or cleanup**:
   - Feature work items remain in `.work-items/` for reference
   - Consider moving to `.work-items/_completed/` if desired

## Output Format

```
=== Completing Feature: {feature-name} ===

Step Verification:
- [x] Step 01: Project structure (DONE)
- [x] Step 02: Seed types (DONE)
...

Acceptance Criteria:
- [x] AC-1: Seed Loading - VERIFIED
- [x] AC-2: Admin API - VERIFIED
...

Test Results:
- Go tests: 45 passed, 0 failed
- Coverage: 87.3%

Time Tracking:
- Estimated: 44.5 hours
- Actual: 38.2 hours
- Time saved: 6.3 hours (14%)

Feature {feature-name} completed successfully!
```

## Error Handling

- If steps incomplete: List remaining steps, do not complete
- If tests fail: Report failures, do not complete
- If wrong feature: Show which feature is active
- If no active feature: Report nothing to complete
