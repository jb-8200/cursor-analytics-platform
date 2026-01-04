---
description: SDD (Spec-Driven Development) process enforcement after task completion
---

# SDD Process Rules

Mandatory workflow steps that enforce Spec-Driven Development discipline.

---

## NEVER

- **Proceed to next task** without completing all 7 steps below
- **Commit without tests passing**: All tests must pass first
- **Commit without reflections checked**: Run dependency-reflection skill
- **Commit without SPEC.md synced**: Update if spec-sync-check triggers
- **Skip git commits**: Every task = commit (no accumulated changes)
- **Update DEVELOPMENT.md** without running SPEC.md sync first
- **Leave task.md stale**: Must reflect actual progress
- **Abandon work mid-step**: Mark as BLOCKED, document blocker

---

## ALWAYS - The 7-Step Workflow

After every task completion, follow **ALL** steps in order:

### Step 1: Verify Tests Pass
```bash
go test ./...              # cursor-sim
npm test                   # analytics-core, viz-spa
```
**Requirement**: All tests passing, no failures or skips

### Step 2: Check Dependency Reflections (REFLECT)
**Ask yourself**:
- Did I modify data models? → Check generators, handlers, schemas, tests
- Did I add/modify endpoints? → Check SPEC.md, E2E tests
- Did I change interfaces? → Check all implementations using them
- Did I complete a phase step? → Check SPEC.md status, task.md

**Run skill**: `dependency-reflection`

**Actions**:
- [ ] Check documentation drift
- [ ] Verify test coverage
- [ ] Ensure all cross-file changes are propagated
- [ ] Run regression tests if high-priority changes

### Step 3: Update SPEC.md if Triggered (SYNC)
**Run skill**: `spec-sync-check`

**Update SPEC.md if ANY**:
- [ ] Completed a phase step (update Implementation Status)
- [ ] Added new endpoint (update Endpoints table)
- [ ] Added new service/package (update Package Structure)
- [ ] Modified API response format (update Response Schemas)
- [ ] Changed CLI flags (update CLI Configuration)

**Process**:
1. Open `services/{service}/SPEC.md`
2. Update relevant sections
3. Update "Last Updated" date
4. Verify accuracy

**Include SPEC.md in commit** if updated

### Step 4: Stage Changes
```bash
git add <files-related-to-task>
# If SPEC.md updated:
git add services/{service}/SPEC.md
```

**Staging rules** (from 02-repo-guardrails.md):
- Only task-related files
- Include SPEC.md if updated
- Exclude unrelated changes
- Verify with `git status` first

### Step 5: Commit with Descriptive Message
```bash
git commit -m "$(cat <<'EOF'
feat(service): complete TASK## - Task Name

Brief description of what was accomplished.

## Changes
- Key change 1
- Key change 2
- Key change 3

Files changed: N
Tests: All passing
Time: Xh actual / Yh estimated

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
)"
```

**Commit guidelines**:
- Descriptive subject line (feat, fix, docs, test, refactor)
- Include what and why
- List key changes
- Track time spent vs estimated
- Sign with generator attribution

### Step 6: Update Progress Tracking

**Update task.md**:
```markdown
### TASK##: Task Name (Date)

**Status**: COMPLETE
**Time**: Xh actual / Yh estimated

**Completed**:
- Deliverable 1
- Deliverable 2

**Changes**:
- Modified: file.go
- Added: file_test.go

**Commit**: {hash}
```

**Update DEVELOPMENT.md**:
- Current status section
- Recently completed work (last 1-2 items)
- Next steps

### Step 7: Only Then Proceed to Next Task
- Verify commit succeeded
- Review git log for correct message
- Check task.md reflects completion
- Acknowledge readiness for next task

---

## Why This Matters

Without enforcement:
- Work gets lost (no commits)
- Dependencies break (no reflections)
- Specs diverge from reality (no SPEC.md sync)
- Progress is unclear (no task.md updates)
- Collaboration becomes impossible

**Enhanced SDD Flow**:
```
Spec → Tests → Code → Tests Pass → REFLECT → SYNC → COMMIT → Next Task
                                      ^^^^    ^^^^    ^^^^^^
                                   Check    Update   Include
                                  deps     SPEC.md   all docs
```

---

## Red Flags

If you catch yourself saying **WITHOUT completing all 7 steps**:
- "Now let's move to Step B02..."
- "Ready for the next step?"
- "Step B01 complete! Would you like to continue..."

**STOP and complete the workflow first!**

---

## Correct Pattern

1. "Step A01 complete. Let me check reflections and SPEC.md sync..."
2. [Run dependency-reflection and spec-sync-check]
3. [Update SPEC.md if triggered]
4. [Stage, commit with proper message, update docs]
5. "Step A01 committed. Ready to start Step A02?"

---

## See Also

- Security rules in `01-security.md`
- Repository guardrails in `02-repo-guardrails.md`
- Coding standards in `03-coding-standards.md`
- sdd-checklist skill: Detailed guidance and examples
