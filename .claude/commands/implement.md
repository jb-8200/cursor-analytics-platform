# Implement Task with Model Selection

When the user runs `/implement [TASK-ID] [--model=haiku|sonnet|opus]`, execute the task using the specified or recommended model.

## Syntax

```bash
# With explicit model
/implement TASK-SIM-003 --model=haiku
/implement TASK-SIM-004 --model=sonnet

# Auto-recommend model
/implement TASK-SIM-003
# Output: "Recommended model: haiku ⚡ (well-specified struct from SPEC.md)"
```

## Behavior

### 1. Parse Task ID
- Extract task ID (e.g., TASK-SIM-003)
- Read docs/TASKS.md to find task details
- Check if task exists in service specifications

### 2. Determine Model

**If user specifies --model flag:**
- Use the specified model (haiku, sonnet, or opus)

**If no model specified:**
- Consult .claude/skills/model-selection-guide.md
- Find recommendation for this task ID
- Display recommendation to user
- Ask for confirmation or proceed with recommendation

### 3. Load Context

Read relevant specifications:
- services/{service}/SPEC.md
- .claude/skills/go-best-practices.md (for Go tasks)
- .claude/skills/cursor-api-patterns.md (for API tasks)
- docs/USER_STORIES.md (related stories)
- .work-items/*/task.md (if active work item exists)

### 4. Execute with TDD

**Phase 1: Tests (Red)**
1. Read or create test-plan.md for this task
2. Write failing tests
3. Run tests to verify they fail
4. Display: "✗ Tests failing as expected"

**Phase 2: Implementation (Green)**
5. Implement code to pass tests
6. Run tests to verify they pass
7. Display: "✓ All tests passing"

**Phase 3: Refactor**
8. Clean up code (extract helpers, add comments)
9. Run linter
10. Run tests again to ensure still passing
11. Display: "✓ Refactored, tests still passing"

**Phase 4: Document**
12. Update relevant documentation
13. Check off task in task.md or docs/TASKS.md

### 5. Report Results

```
Task: TASK-SIM-003 - Implement Developer Profile Generator
Model: haiku ⚡
Status: ✓ Complete

Changes Made:
- services/cursor-sim/internal/generator/developer.go (new)
- services/cursor-sim/internal/generator/developer_test.go (new)

Test Results:
- Unit tests: 15/15 passing
- Coverage: 92.3%

Next Task: TASK-SIM-004 - Implement Event Generation Engine
Recommended model: sonnet ⚡⚡ (complex Poisson distribution logic)
```

## Model Recommendations

### How Recommendations Work

The command reads `.claude/skills/model-selection-guide.md` and matches:

```markdown
| TASK-SIM-003 | Developer Profile Generator | **Haiku** | Struct implementation from SPEC.md:145-250 |
```

### Recommendation Display

```
Task: TASK-SIM-003 - Implement Developer Profile Generator

Model Recommendation: haiku ⚡
Reason: Well-specified struct implementation from SPEC.md lines 145-250

Estimated Cost: ~$0.20
Estimated Time: ~5 minutes

Options:
1. Use recommended model (haiku)
2. Override with sonnet
3. Override with opus

Proceed with haiku? [Y/n]
```

## Task Tool Integration

When executing with a sub-agent:

```python
# If model specified or recommended as haiku
Task(
    subagent_type="general-purpose",
    model="haiku",
    prompt=f"""
    Implement {task_id} following TDD:

    1. Read SPEC: services/cursor-sim/SPEC.md lines 145-250
    2. Write failing tests for Developer struct
    3. Implement Developer struct with all fields
    4. Ensure tests pass
    5. Follow go-best-practices.md patterns

    Return: Summary of changes and test results
    """
)
```

## Examples

### Example 1: Well-Specified Task (Haiku)

```bash
User: /implement TASK-SIM-003

Claude: Analyzing TASK-SIM-003...

Task: TASK-SIM-003 - Implement Developer Profile Generator
Model Recommendation: haiku �
Reason: Well-specified struct implementation from SPEC.md lines 145-250

Using model: haiku

[Haiku agent spawned via Task tool]
 Tests written (developer_test.go)
 Developer struct implemented
 All tests passing (15/15)
 Coverage: 92.3%
 Linter: No issues

Task complete!
```

### Example 2: Complex Task (Sonnet)

```bash
User: /implement TASK-SIM-004

Claude: Analyzing TASK-SIM-004...

Task: TASK-SIM-004 - Implement Event Generation Engine
Model Recommendation: sonnet ��
Reason: Complex Poisson distribution logic, concurrent event generation

Using model: sonnet

[Sonnet agent spawned via Task tool]
 Tests written (event_generator_test.go)
 Poisson distribution helper implemented
 Event scheduler implemented with goroutines
 All tests passing (28/28)
 Coverage: 87.5%
 Race detector: Clean

Task complete!
```

### Example 3: Override Model

```bash
User: /implement TASK-SIM-003 --model=sonnet

Claude: Task: TASK-SIM-003 - Implement Developer Profile Generator

Note: Recommended model for this task is haiku (well-specified struct)
You requested: sonnet

Proceeding with sonnet...

[Proceeds with Sonnet instead of Haiku]
```

## Implementation Notes

When you implement this command, you should:

1. Parse the task ID and optional --model flag
2. Look up the task in model-selection-guide.md
3. Display the recommendation with cost/time estimates
4. Use the Task tool with the specified/recommended model
5. Follow TDD workflow (Red-Green-Refactor)
6. Report results with next task suggestion

This enables cost-effective development by using Haiku for well-specified tasks while reserving Sonnet/Opus for complex logic and architectural decisions.
