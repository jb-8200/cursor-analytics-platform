# Design Document: SDD Subagent Orchestration Protocol

**Feature ID**: P0-F01
**Epic**: P0 - Project Management
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

This document defines the operational protocol for coordinating master agents and subagents during parallel development sessions using the SDD (Spec-Driven Development) workflow.

---

## Architecture

### Agent Hierarchy

```
┌───────────────────────────────────────────────────────────────────────┐
│                        Master Agent (Orchestrator)                    │
│  ───────────────────────────────────────────────────────────────────  │
│  Responsibilities:                                                    │
│  • Task delegation to subagents                                       │
│  • Code review across all changes                                     │
│  • E2E testing for cross-service issues                              │
│  • DEVELOPMENT.md and plan folder updates                            │
│  • Final commit coordination                                          │
└───────────────────────────────────────────────────────────────────────┘
          │                    │                    │
          ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ cursor-sim-cli  │  │ analytics-core  │  │   viz-spa-dev   │
│   Subagent      │  │   Subagent      │  │   Subagent      │
│  ─────────────  │  │  ─────────────  │  │  ─────────────  │
│ P4: CLI only    │  │ P5: GraphQL     │  │ P6: React/Vite  │
│ NEVER API/Gen   │  │ + PostgreSQL    │  │ + Apollo Client │
└─────────────────┘  └─────────────────┘  └─────────────────┘
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Feature-Level task.md (each subagent updates only its own)         │
│  .work-items/{feature}/task.md                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Protocol Phases

### Phase 1: Task Delegation

```
Master Agent
    │
    ├─→ Identifies parallel work opportunities
    ├─→ Creates/assigns tasks to appropriate subagents
    ├─→ Provides context (relevant files, acceptance criteria)
    └─→ Launches subagents (in parallel when possible)
```

**Master Agent Actions**:
1. Read DEVELOPMENT.md for current state
2. Identify tasks that can run in parallel
3. Use Task tool to spawn subagents
4. Provide clear, detailed prompts with context

---

### Phase 2: Subagent Execution

```
Subagent (per feature/service)
    │
    ├─→ Read SPEC.md for service specification
    ├─→ Follow TDD: RED → GREEN → REFACTOR
    ├─→ Run sdd-checklist before commit
    ├─→ Update task.md with status and notes
    └─→ Report completion to master agent
```

**Subagent Constraints**:
- **ONLY** update `.work-items/{feature}/task.md`
- **NEVER** update `.claude/DEVELOPMENT.md` (master agent only)
- **NEVER** modify plan folder symlinks
- **FOCUS** on assigned service (no cross-service changes)

**task.md Update Format**:
```markdown
### TASK##: {Task Name} (Date)

**Status**: {TODO|IN_PROGRESS|COMPLETE|BLOCKED}
**Time**: {actual}h / {estimated}h

**Completed**:
- {deliverable-1}: {brief description}
- {deliverable-2}: {brief description}

**Changes**:
- Modified: {file-path}
- Added: {file-path}
- Tests: {test-file-path}

**Issues/Blockers** (if any):
- {issue-1}: {description and impact}

**Commit**: {commit-hash} (if applicable)

**Next Steps**:
- {next-action-1}
```

**Update Triggers**:
1. When starting task: Set status to IN_PROGRESS
2. After each commit: Update Completed, Changes, Commit hash
3. When blocked: Set status to BLOCKED, document issue
4. When complete: Set status to COMPLETE, document final time

**Subagent Completion Report Format**:
```
TASK COMPLETE: {task-id}
Status: {PASSED|BLOCKED|PARTIAL}
Commit: {commit-hash} (if applicable)
Notes: {any blockers, dependencies, or follow-up items}
```

---

### Phase 3: Master Agent Code Review

```
Master Agent (after all subagents complete)
    │
    ├─→ Wait for all subagent completion reports
    ├─→ Review code changes across services
    │   ├─→ Check for code quality issues
    │   ├─→ Check for cross-service consistency
    │   └─→ Check for contract alignment (API, GraphQL, types)
    │
    ├─→ If issues found:
    │   ├─→ Create new task in appropriate task.md
    │   └─→ Delegate to subagent for resolution
    │
    └─→ If no issues: proceed to E2E testing
```

**Code Review Checklist**:
1. [ ] All tests passing in each service
2. [ ] No TypeScript/Go lint errors
3. [ ] API contracts match between services
4. [ ] GraphQL schema/queries aligned (P5↔P6)
5. [ ] No hardcoded values or secrets
6. [ ] Error handling appropriate
7. [ ] Documentation updated where needed

**Issue Escalation**:
- **Minor issues**: Create task, delegate to subagent
- **Cross-service issues**: Master agent resolves directly
- **Blocking issues**: Pause workflow, clarify with user

---

### Phase 4: Master Agent E2E Testing

```
Master Agent
    │
    ├─→ Start all services (Docker or local)
    ├─→ Run E2E test suite
    │   ├─→ API endpoint validation (P4)
    │   ├─→ GraphQL query execution (P5)
    │   └─→ Frontend data display (P6)
    │
    ├─→ If E2E issues found:
    │   └─→ Master agent fixes directly (avoids cross-subagent coordination)
    │
    └─→ Document any E2E fixes made
```

**Why Master Agent Fixes E2E Issues**:
- Cross-service issues require understanding of multiple services
- Delegating to subagent requires additional context transfer
- Master agent has full session context
- Faster resolution than round-trip to subagent

**E2E Test Categories**:
1. **Service Health**: All endpoints responding
2. **Data Flow**: P4 → P5 → P6 data pipeline
3. **Contract Compliance**: Request/response formats
4. **Integration**: Frontend displays correct data

---

### Phase 5: Documentation & Final Commit

```
Master Agent
    │
    ├─→ Update DEVELOPMENT.md
    │   ├─→ Session summary
    │   ├─→ Tasks completed
    │   ├─→ Issues resolved
    │   └─→ Next steps
    │
    ├─→ Update plan folder
    │   ├─→ Remove active symlink if feature complete
    │   └─→ Update progress tracking
    │
    └─→ Final commit
        ├─→ Any uncommitted changes from E2E fixes
        └─→ Documentation updates
```

**DEVELOPMENT.md Update Template**:
```markdown
### Session Summary ({date})
**Completed**:
- {task-1}: {brief description}
- {task-2}: {brief description}

**Issues Resolved**:
- {issue-1}: {how resolved}

**Next Steps**:
- {next-task-1}
```

---

## Decision Matrix

| Situation | Actor | Action |
|-----------|-------|--------|
| Task assignment | Master Agent | Use Task tool to spawn subagent |
| Task completion | Subagent | Update task.md, report to master |
| Code quality issue (single service) | Subagent | Create task, delegate to same subagent |
| Cross-service issue | Master Agent | Fix directly |
| E2E failure | Master Agent | Fix directly |
| DEVELOPMENT.md update | Master Agent | Update after E2E passes |
| Feature completion | Master Agent | Remove active symlink, final commit |

---

## File Ownership

| File/Path | Owner | Update Trigger |
|-----------|-------|----------------|
| `.work-items/{feature}/task.md` | Subagent | After each task |
| `.work-items/{feature}/design.md` | Either | When design changes |
| `.work-items/{feature}/user-story.md` | Either | When requirements change |
| `.claude/DEVELOPMENT.md` | Master Agent ONLY | After E2E passes |
| `.claude/plans/active` | Master Agent ONLY | Feature start/complete |
| `services/{service}/SPEC.md` | Subagent | After API changes |

---

## Error Handling

### Subagent Failure
```
If subagent reports BLOCKED:
    1. Master agent reviews blocker
    2. If resolvable: provide guidance, retry
    3. If not resolvable: escalate to user
```

### E2E Test Failure
```
If E2E tests fail:
    1. Identify failing test(s)
    2. Determine root cause service
    3. Fix directly (master agent)
    4. Re-run E2E tests
    5. Document fix in session summary
```

### Conflicting Changes
```
If subagents modify same file:
    1. Master agent reviews both changes
    2. Merge manually if compatible
    3. If incompatible: revert one, create clarification task
```

---

## Anti-Patterns

### DO NOT:
- Allow subagents to update DEVELOPMENT.md
- Allow subagents to modify plan folder symlinks
- Delegate E2E fixes to subagents (coordination overhead)
- Skip code review after subagent completion
- Commit without running sdd-checklist

### DO:
- Run subagents in parallel when tasks are independent
- Provide full context when delegating tasks
- Wait for all subagents before code review
- Document all E2E fixes made by master agent
- Update DEVELOPMENT.md as final step

---

## Implementation Notes

This is a **protocol document**, not a code implementation. The protocol is enforced through:

1. **Claude Code skills**: `sdd-checklist`, `spec-process-core`
2. **TodoWrite tool**: Track tasks and status
3. **Agent prompts**: Include protocol reminders
4. **CLAUDE.md**: Reference this protocol

No code changes required. This document serves as the authoritative reference for agent behavior.
