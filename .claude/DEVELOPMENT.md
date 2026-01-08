# Development Session Context

**Last Updated**: January 8, 2026
**Active Features**: P3-F04 (Preview Mode & YAML Support - Planning Complete)
**Primary Focus**: Add preview mode and YAML seed support inspired by NVIDIA DataDesigner

---

## Current Status

### Phase Overview

| Phase | Description | Status |
|-------|-------------|--------|
| **P0** | Project Management | **COMPLETE** ✅ (8/8 features) |
| **P1** | cursor-sim Foundation | COMPLETE ✅ |
| **P2** | cursor-sim GitHub Simulation | TODO |
| **P3** | cursor-sim Research Framework | COMPLETE ✅ |
| **P4** | cursor-sim CLI Enhancements | IN PROGRESS (P4-F02) |
| **P5** | cursor-analytics-core | COMPLETE ✅ |
| **P6** | cursor-viz-spa | COMPLETE ✅ |
| **P7** | Deployment Infrastructure | COMPLETE ✅ |

### Recently Completed: P0 Phase (8 features, 8.0 hours)

1. **P0-F01**: SDD Subagent Orchestration Protocol (1.0h)
2. **P0-F02**: Rules Layer Implementation (3.0h)
3. **P0-F03**: Skills Cleanup & Catalog (1.5h)
4. **P0-F04**: Commands/Prompts Restructure (1.0h)
5. **P0-F05**: Hooks Configuration (0.65h)
6. **P0-F06**: Agents Documentation (0.5h)
7. **P0-F07**: Selection Heuristic Guide (0.5h)
8. **P0-F08**: DEVELOPMENT.md Optimization (0.7h)

**Infrastructure Complete**: All Claude Code features configured and ready for service development.

---

## Next Steps

### In Progress
- **P3-F04**: Preview Mode & YAML Support - 0/11 tasks complete (Planning done)
  - Inspired by NVIDIA NeMo DataDesigner patterns
  - Preview mode for quick seed validation (< 5 seconds)
  - YAML seed file support with comments
  - Validation framework with warnings
  - Estimated: 10.5h

### Recently Completed
- **P4-F02**: CLI Enhancement - 14/14 tasks complete ✅ (13.0h total)
  - Interactive configuration prompts with validation
  - Developer replication (downsample/upsample from seed)
  - Commit limit enforcement
  - Generator bug fix (empty dataset issue resolved)
  - All manual tests passed, UX polished

### Ready to Start
- **P6-F02**: GraphQL Code Generator (viz-spa) - New feature
- **P2**: GitHub Simulation (PR, issues, reviews) - Deferred from earlier phases

### Parallel Development Enabled
- Subagent infrastructure: cursor-sim-cli-dev, analytics-core-dev, viz-spa-dev, cursor-sim-infra-dev
- Commands: `/subagent/cursor-sim-cli`, `/subagent/analytics-core`, `/subagent/viz-spa`
- Rules enforcement: 7 rule files with NEVER/ALWAYS constraints
- Hooks: Pre-commit reminders, markdown formatting, SDD checklists

---

## Quick Reference

### Session Start Checklist
1. [ ] Read `.claude/DEVELOPMENT.md` (this file)
2. [ ] Check active work: `readlink .claude/plans/active`
3. [ ] Review current task status
4. [ ] Follow SDD workflow: SPEC → TEST → CODE → COMMIT

### Common Commands
| Command | Purpose |
|---------|---------|
| `/start-feature {name}` | Start feature, create symlink |
| `/implement {task-id}` | TDD implementation |
| `/status` | Show current state |
| `/spec {service}` | Display specification |

### Running Services
```bash
# cursor-sim (port 8080)
cd services/cursor-sim && go build && ./bin/cursor-sim -port 8080

# cursor-analytics-core (port 4000)
cd services/cursor-analytics-core && npm run dev

# cursor-viz-spa (port 3000)
cd services/cursor-viz-spa && npm run dev
```

---

## Key Files

| File | Purpose |
|------|---------|
| `.claude/rules/` | Enforcement constraints (NEVER/ALWAYS) |
| `.claude/skills/` | Knowledge guides (14 skills) |
| `.claude/commands/` | Slash commands and workflows |
| `.claude/agents/` | Subagent definitions (4 agents) |
| `.claude/hooks/README.md` | Hook configuration |
| `.work-items/P*/` | Active feature directories |

---

**See `.claude/archive/` for historical session summaries and integration testing notes from earlier development.**
