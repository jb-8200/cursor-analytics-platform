# Cursor Analytics Platform

## Workflow: Spec-Driven Development (SDD)

```
1. SPEC     → Read specification before coding
2. TEST     → Write failing tests (RED)
3. CODE     → Minimal implementation (GREEN)
4. REFACTOR → Clean up while tests pass
5. REFLECT  → Check dependency reflections
6. SYNC     → Update SPEC.md if triggered
7. COMMIT   → Every task = commit (code + docs)
```

**Full methodology**: `docs/spec-driven-design.md`

---

## Session Start

1. Read `.claude/DEVELOPMENT.md` for current state
2. Check `.claude/plans/active` for active work
3. Review task status in `.work-items/{feature}/task.md`

---

## Documentation Hierarchy

### Source of Truth (in priority order)

| Location | Purpose |
|----------|---------|
| `services/{service}/SPEC.md` | Technical specification |
| `.work-items/{feature}/` | Active work tracking |
| `.claude/DEVELOPMENT.md` | Session context |

### Work Items Structure

```
Phase (P#) = Epic level
  └── Feature (F##) = Work item directory
       └── Task (TASK##) = Implementation step

.work-items/{P#-F##-feature-name}/
├── user-story.md    # Requirements (what + why)
├── design.md        # Technical approach (how)
└── task.md          # Implementation tasks (TASK01, TASK02...)
```

---

## Skills (Auto-Discovered)

Skills activate automatically based on your request. Available skills:

| Skill | Purpose |
|-------|---------|
| `spec-process-core` | Core SDD principles |
| `spec-process-dev` | TDD workflow |
| `sdd-checklist` | **CRITICAL**: Post-task commit |
| `spec-sync-check` | **NEW**: SPEC.md update triggers |
| `dependency-reflection` | **NEW**: Dependency checking |
| `spec-user-story` | User story format |
| `spec-design` | Design doc format |
| `spec-tasks` | Task breakdown format |
| `go-best-practices` | Go patterns |
| `cursor-api-patterns` | API implementation |
| `model-selection-guide` | Model optimization |

Skills trigger when your request matches their description.

---

## Task Completion (CRITICAL)

After **every** task, follow `sdd-checklist`:

1. ✅ Tests pass
2. ✅ **Check reflections** (dependency-reflection)
3. ✅ **Update SPEC.md** if needed (spec-sync-check)
4. ✅ Git commit (code + SPEC.md if updated)
5. ✅ Update task.md
6. ✅ Update DEVELOPMENT.md
7. ✅ Proceed to next task

**Never** move to the next task without:
- Checking reflections
- Syncing SPEC.md
- Committing all changes

---

## Subagent Orchestration Protocol

When using subagents for parallel development:

### Master Agent Responsibilities

1. **Task Delegation**: Spawn subagents with clear context and scope
2. **Code Review**: Review all subagent changes for quality and consistency
3. **E2E Testing**: Run full-stack tests, fix cross-service issues directly
4. **Documentation**: Update DEVELOPMENT.md and plan folder after E2E passes
5. **Final Commit**: Commit all changes including E2E fixes and docs

### Subagent Constraints

- **ONLY** update `.work-items/{feature}/task.md`
- **NEVER** update `.claude/DEVELOPMENT.md` (master agent only)
- **NEVER** modify plan folder symlinks
- **FOCUS** on assigned service (no cross-service changes)

### Completion Flow

```
Subagent → Update task.md → Report completion
    ↓
Master Agent → Code review → E2E testing → Fix issues → Update docs → Commit
```

**Full Protocol**: `.work-items/P0-F01-sdd-subagent-orchestration/design.md`

---

## Commands

| Command | Purpose |
|---------|---------|
| `/start-feature {name}` | Initialize feature |
| `/complete-feature {name}` | Verify and close |
| `/implement {task-id}` | TDD implementation |
| `/status` | Current state |
| `/next-task` | Find next work |

---

## Services

| Service | Tech | Port | Status |
|---------|------|------|--------|
| cursor-sim | Go | 8080 | Active |
| cursor-analytics-core | TypeScript/GraphQL | 4000 | Pending |
| cursor-viz-spa | React/Vite | 3000 | Pending |

---

## Quick Commands

```bash
# cursor-sim
cd services/cursor-sim
go test ./...                    # Run tests
go build -o bin/cursor-sim ./cmd/simulator  # Build

# Run simulator
./bin/cursor-sim -mode runtime -seed testdata/valid_seed.json -port 8080

# Test endpoint
curl http://localhost:8080/health
```

---

## Key Files

| File | Purpose |
|------|---------|
| `.claude/DEVELOPMENT.md` | **START HERE** - Session context |
| `.claude/README.md` | Claude Code integration guide |
| `services/cursor-sim/SPEC.md` | cursor-sim specification |
| `docs/spec-driven-design.md` | Full SDD methodology |

---

## Hooks: Documentation Only

The `.claude/hooks/` Python files do **NOT execute** in Claude Code.

**Alternative**: Use `sdd-checklist` skill + TodoWrite for workflow enforcement.

---

**Terminology**: Phase (P#) → Feature (F##) → Task (TASK##)

**Spec first. Tests first. Reflect before commit. SPEC.md stays current. Commit always.**
