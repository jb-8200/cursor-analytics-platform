# Show Next Task

When the user runs `/next-task [service-name]`, display the next unfinished task from the implementation plan.

## Task Sources (Priority Order)

1. **Active Work Item** (if `.claude/plans/active` symlink exists)
   - Read `.claude/plans/active/task.md`
   - Find first unchecked task `- [ ]`
   - Display with context

2. **Service Tasks** (from docs/TASKS.md)
   - Read `docs/TASKS.md`
   - Find section for specified service
   - Identify next incomplete task based on dependencies

3. **P0 Tasks** (from P0_MAKERUNNABLE.md)
   - If no service specified and P0 incomplete
   - Read `P0_MAKERUNNABLE.md`
   - Show next P0 task

## Display Format

```
Next Task: {TASK-ID} - {Task Name}
Service: {service-name}
Priority: {P0/P1/P2}

Description:
{Task description from docs}

Acceptance Criteria:
- {Criterion 1}
- {Criterion 2}

Dependencies:
- {Dependency 1} (Status: ✓ Complete)
- {Dependency 2} (Status: ✗ Incomplete - BLOCKER)

Files to Modify:
- {file1}
- {file2}

Relevant Specs:
- services/{service}/SPEC.md (lines {nn}-{mm})
- docs/USER_STORIES.md (lines {nn}-{mm})

Start with: Write tests first! See test-plan.md or TESTING_STRATEGY.md

Recommended Model: haiku ⚡
(Well-specified struct from SPEC.md - see .claude/skills/model-selection-guide.md)
```

## Dependency Resolution

If next task has incomplete dependencies, show:
```
⚠️  Next task (TASK-SIM-006) is blocked by:
   - TASK-SIM-005: Implement In-Memory Storage (Incomplete)

Suggested: Complete TASK-SIM-005 first
```

## Example Usage

```
User: /next-task cursor-sim
Assistant: [Analyzes task status and displays next actionable task]
```

## Implementation Instructions

1. Determine task source (active work item > service tasks > P0)
2. Read task list and identify status (completed vs incomplete)
3. Check dependencies before suggesting a task
4. Provide all context needed to start work immediately
5. Always remind to write tests first (TDD)
