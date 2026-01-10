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

### Primary Path (GraphQL)

| Service | Tech | Port | Status |
|---------|------|------|--------|
| cursor-sim | Go | 8080 | Active |
| cursor-analytics-core | TypeScript/GraphQL | 4000 | Pending |
| cursor-viz-spa | React/Vite | 3000 | Pending |

### Alternative Path (dbt + Streamlit)

| Service | Tech | Port | Status |
|---------|------|------|--------|
| api-loader | Python | N/A | Active |
| streamlit-dashboard | Python/Streamlit | 8501 | Active |

**Two Analytics Paths**:
- **Path 1 (GraphQL)**: cursor-sim → analytics-core (GraphQL) → viz-spa (React)
- **Path 2 (dbt)**: cursor-sim → api-loader (ETL) → dbt (transforms) → streamlit-dashboard

---

## Data Contract Hierarchy

**cursor-sim is the authoritative source of truth** for the entire platform. All downstream services and layers must validate against the API contract defined in `services/cursor-sim/SPEC.md`.

### Contract Levels

```
LEVEL 1: API CONTRACT (cursor-sim SPEC.md) ← SOURCE OF TRUTH
├─ Endpoints: /analytics/ai-code/commits, /repos/*/pulls, /research/dataset
├─ Response format: {items: [...], totalCount, page, pageSize}
├─ Field names: camelCase (commitHash, userEmail, tabLinesAdded, ...)
├─ Data types: strings, numbers, dates in ISO format
└─ Responsibility: cursor-sim (P4)

LEVEL 2: GRAPHQL PATH (cursor-analytics-core) ← PATH 1
├─ GraphQL schema reflects API contracts
├─ Resolvers fetch from cursor-sim API
├─ Aggregations and joins performed in TypeScript
└─ Responsibility: analytics-core (P5)

LEVEL 3: DATA TIER PATH (api-loader → dbt → DuckDB) ← PATH 2
├─ Raw schema: Preserves API fields exactly
├─ Staging schema: Transforms camelCase → snake_case
├─ Mart schema: Aggregations for analytics
└─ Responsibility: api-loader (P8) + dbt (P8) + streamlit-dashboard (P9)

LEVEL 4: VISUALIZATION (cursor-viz-spa + streamlit-dashboard)
├─ Frontend consumes contracted data
├─ Path 1: Queries GraphQL API (analytics-core)
├─ Path 2: Queries mart tables (DuckDB/Snowflake)
└─ Responsibility: viz-spa (P6) + streamlit-dashboard (P9)
```

### Key Principles

1. **API as Fact**: cursor-sim SPEC.md is the single source of truth
2. **Data Fidelity**: Each layer preserves data from previous layer
3. **Explicit Contracts**: All transformations must be documented
4. **Two Paths, Same Contract**: Both GraphQL and dbt paths consume same API
5. **Test at Boundaries**: Validate contract at each layer transition

### Reference Documents (in priority order)

| Priority | Document | Purpose |
|----------|----------|---------|
| **1** | `services/cursor-sim/SPEC.md` | API contract (source of truth) |
| **2** | `docs/design/new_data_architecture.md` | Data pipeline architecture (P8/P9) |
| **3** | `docs/DESIGN.md` | System architecture overview |
| **4** | `docs/TESTING_STRATEGY.md` | Testing approaches for contract validation |
| **5** | `docs/data-contract-testing.md` | Data contract validation patterns |

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
