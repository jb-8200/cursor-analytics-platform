# Cursor Analytics Platform

## Workflow: Spec-Driven Development (SDD)

```
1. SPEC    → Read specification before coding
2. TEST    → Write failing tests (RED)
3. CODE    → Minimal implementation (GREEN)
4. REFACTOR → Clean up while tests pass
5. COMMIT  → Every task = commit (CRITICAL)
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
.work-items/{feature}/
├── user-story.md    # Requirements (what + why)
├── design.md        # Technical approach (how)
└── task.md          # Implementation tasks (steps)
```

---

## Skills (Auto-Discovered)

Skills activate automatically based on your request. Available skills:

| Skill | Purpose |
|-------|---------|
| `spec-process-core` | Core SDD principles |
| `spec-process-dev` | TDD workflow |
| `spec-user-story` | User story format |
| `spec-design` | Design doc format |
| `spec-tasks` | Task breakdown format |
| `go-best-practices` | Go patterns |
| `cursor-api-patterns` | API implementation |
| `sdd-checklist` | **CRITICAL**: Post-task commit |
| `model-selection-guide` | Model optimization |

Skills trigger when your request matches their description.

---

## Task Completion (CRITICAL)

After **every** task, follow `sdd-checklist`:

1. ✅ Tests pass
2. ✅ Git commit
3. ✅ Update task.md
4. ✅ Update DEVELOPMENT.md
5. ✅ Proceed to next task

**Never** move to the next task without committing.

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

**Spec first. Tests first. Commit always.**
