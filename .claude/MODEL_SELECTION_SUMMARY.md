# Model Selection Integration - Implementation Summary

## What Was Implemented

The Claude Code native SDD structure now includes **automatic model selection and recommendations** to optimize cost and speed while maintaining quality.

---

## New Components

### 1. Model Selection Guide Skill
**File**: `.claude/skills/model-selection-guide.md` (340 lines)

Provides task-to-model mapping for all project tasks:

| Model | Use Case | Example Tasks |
|-------|----------|---------------|
| **Haiku** ⚡ | Well-specified structs, pattern following, simple tests | TASK-SIM-003, TASK-SIM-005, TASK-VIZ-002 |
| **Sonnet** ⚡⚡ | Complex logic, architectural decisions, integration | TASK-SIM-004, TASK-CORE-004, TASK-VIZ-004 |
| **Opus** ⚡ | Novel architectures, major refactoring, cross-service | New major features, performance redesign |

**Key Sections**:
- Task-to-Model mapping for all P0 and P1 tasks
- Activity-based recommendations (implementation vs design)
- Decision tree for model selection
- Cost optimization strategies (hybrid approach saves ~64%)
- Integration with slash commands

### 2. Enhanced Slash Commands

#### `/implement [TASK-ID] [--model=haiku|sonnet|opus]`
**File**: `.claude/commands/implement.md` (200+ lines)

New command that:
1. Parses task ID and optional --model flag
2. Consults model-selection-guide.md for recommendation
3. Displays cost/time estimates
4. Spawns Task agent with appropriate model
5. Follows TDD workflow (Red-Green-Refactor)
6. Reports results with next task suggestion

**Usage**:
```bash
# Auto-recommend model
/implement TASK-SIM-003
# → Uses Haiku (⚡ cheap & fast)

# Override model
/implement TASK-SIM-003 --model=sonnet
# → Uses Sonnet (more powerful but slower)

# Complex task
/implement TASK-SIM-004
# → Auto-recommends Sonnet for Poisson distribution logic
```

#### Enhanced `/next-task [service-name]`
Now includes model recommendation in output:

```
Next Task: TASK-SIM-003 - Implement Developer Profile Generator
Recommended Model: haiku ⚡
(Well-specified struct from SPEC.md - see model-selection-guide.md)
```

#### Enhanced `/start-feature [name]`
Now recommends Sonnet for planning phase:

```
Recommended Model for Planning: sonnet ⚡⚡
(Architectural decisions and design - then use haiku for implementation)
```

---

## How It Works

### Model Selection Flow

```
User runs: /implement TASK-SIM-003
    ↓
Claude reads: .claude/skills/model-selection-guide.md
    ↓
Finds: TASK-SIM-003 → Haiku (well-specified struct)
    ↓
Displays: "Recommended model: haiku ⚡ (~$0.20, ~5 min)"
    ↓
Spawns: Task(model="haiku", prompt="Implement Developer struct...")
    ↓
Haiku executes: Read SPEC → Write tests → Implement → Refactor
    ↓
Reports: "✓ Complete. Coverage: 92.3%. Next: TASK-SIM-004 (use sonnet)"
```

### Task Tool Integration

Commands can now spawn sub-agents with specific models:

```python
# Well-specified implementation → Use Haiku
Task(
    subagent_type="general-purpose",
    model="haiku",
    prompt="Implement Developer struct from SPEC.md lines 145-250"
)

# Complex logic → Use Sonnet
Task(
    subagent_type="general-purpose",
    model="sonnet",
    prompt="Implement Poisson distribution event generator"
)
```

---

## Cost Optimization Benefits

### Example: P0 Scaffolding Tasks

**All-Sonnet Approach**:
- 8 tasks × ~$0.50/task = **$4.00 total**
- Time: ~2 hours

**Hybrid Approach** (with model-selection-guide):
- Planning (Sonnet): 3 tasks × $0.50 = $1.50
- Implementation (Haiku): 5 tasks × $0.10 = $0.50
- **Total: $2.00** (50% savings)
- Time: ~1.5 hours (faster)

### Example: Full cursor-sim Implementation

**All-Sonnet**:
- 7 tasks × ~$2.00/task = **$14.00**

**Hybrid** (following model-selection-guide):
- Complex (Sonnet): TASK-SIM-001, 004, 007 = $6.00
- Simple (Haiku): TASK-SIM-002, 003, 005, 006 = $0.80
- **Total: $6.80** (51% savings)

---

## Usage Examples

### Starting a New Feature

```bash
# Step 1: Plan with Sonnet (architectural decisions)
User: "Use Sonnet to run /start-feature developer-generator"
Claude: [Creates design.md, task.md, test-plan.md with Sonnet]

# Step 2: Implement with Haiku (well-specified work)
User: "Use Haiku to implement from the task list"
Claude: [Spawns Haiku agent for implementation]

# Step 3: Integrate with Sonnet (complex wiring)
User: "Use Sonnet to wire up the generator in main.go"
Claude: [Handles integration with Sonnet]
```

### Batch Implementation

```bash
# Let Claude auto-select models
User: "/implement TASK-SIM-002"
Claude: "Using haiku ⚡ (CLI flag parsing - simple)"

User: "/implement TASK-SIM-003"
Claude: "Using haiku ⚡ (Developer struct - well-specified)"

User: "/implement TASK-SIM-004"
Claude: "Using sonnet ⚡⚡ (Poisson distribution - complex)"
```

### Override When Needed

```bash
# Force Sonnet for a Haiku-recommended task
User: "/implement TASK-SIM-003 --model=sonnet"
Claude: "Note: Recommended model is haiku, but proceeding with sonnet as requested"
```

---

## Decision Tree for Users

```
Q: Do I know which model to use?
├─ YES → Use /implement TASK-ID --model=<model>
└─ NO  → Use /implement TASK-ID (auto-recommend)
    ↓
    Claude shows: "Recommended: haiku ⚡ (~$0.20, ~5 min)"
    ↓
    Q: Accept recommendation?
    ├─ YES → Proceed with recommended model
    └─ NO  → Override with --model flag
```

---

## Files Created/Modified

### New Files
1. `.claude/skills/model-selection-guide.md` (340 lines)
2. `.claude/commands/implement.md` (200+ lines)

### Modified Files
1. `.claude/commands/next-task.md` - Added model recommendation
2. `.claude/commands/start-feature.md` - Added model recommendation

### Integration
All commands now reference `model-selection-guide.md` for consistent recommendations across the project.

---

## Best Practices

### When to Override Recommendations

**Override to Haiku** if:
- You're experimenting and OK with potential retry
- Cost is critical concern
- Task is even simpler than guide suggests

**Override to Sonnet** if:
- First time implementing this pattern
- Debugging an unexpected issue
- Task is more complex than guide suggests

**Override to Opus** if:
- Major architectural decision needed
- Complex multi-service debugging
- Novel feature not covered in specs

### Session Planning

**Planning Session** (Sonnet):
- Run `/start-feature [name]`
- Make design decisions
- Write test plans

**Implementation Session** (Haiku):
- Run `/implement` for well-specified tasks
- Write tests and code
- Follow established patterns

**Integration Session** (Sonnet):
- Wire up components
- Debug integration issues
- Verify end-to-end

---

## Next Steps

To use this system:

1. **Review model-selection-guide.md** to understand task-to-model mappings
2. **Use /next-task** to see what to work on with model recommendation
3. **Use /implement** to execute tasks with optimal model selection
4. **Override when needed** with --model flag

The comprehensive specs in this project make it ideal for Haiku usage on 60-70% of tasks, with Sonnet for complex logic and Opus for novel architectures.

---

**Result**: Cost-effective development without sacrificing quality. Use the right tool for the job!
