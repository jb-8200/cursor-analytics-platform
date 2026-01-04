---
name: spec-sync-check
description: Determines if SPEC.md needs updating after implementation changes. Triggers automatically when completing tasks that modify endpoints, models, services, or complete phases. Use when finishing any implementation task. (project)
---

# SPEC.md Sync Check

**Purpose**: Ensure service specifications stay synchronized with implementation.

## When to Use

Run this check after completing ANY task that:
- Adds or modifies API endpoints
- Changes data models or schemas
- Completes a phase or major milestone
- Adds new packages or services
- Modifies CLI configuration

## Automatic Trigger Detection

After completing a task, check if any of these conditions apply:

### High-Priority Triggers (MUST update SPEC.md)

#### 1. Phase Completion
**Detection**: Completing the final step of any phase (e.g., C06, B08, A07)

**Update Sections**:
- Implementation Status table (lines 16-23)
- Phase features section (mark as "Implemented ✅")
- Last Updated date

**Example**: Completing C06 (last step of Phase 3 Part C)
- Change "Phase 3 (P2)" status from "IN_PROGRESS" to "MOSTLY COMPLETE ✅"
- Update Phase 3 Features section from "Planned" to "Mostly Complete ✅"

#### 2. New Endpoint Added
**Detection**: New handler file created/modified in `internal/api/`

**Update Sections**:
- Endpoints table (lines 89-148)
- Add new row with method, path, auth requirement, status
- Update endpoint count in Overview if needed

**Example**: Adding `/analytics/team/models`
```markdown
| GET | `/analytics/team/models` | Yes | ✅ Implemented |
```

#### 3. New Service/Package Created
**Detection**: New directory created in `internal/`

**Update Sections**:
- Package Structure (lines 299-323)
- Add new package with description
- Architecture diagram if significant

**Example**: Adding `internal/services/`
```
├── services/         # Business logic (code survival, hotfix, revert analysis)
```

### Medium-Priority Triggers (SHOULD update SPEC.md)

#### 4. Model Changes
**Detection**: Modifications to `internal/models/` struct fields

**Update Sections**:
- Response Format section (lines 152-200)
- Update schema examples
- Update field descriptions

**Example**: Adding new field to Commit model
- Update Commit Schema JSON example
- Update invariant if calculation changes

#### 5. Generator Changes
**Detection**: New or significantly modified generators in `internal/generator/`

**Update Sections**:
- Generation Algorithm section (lines 328-346)
- Phase features if new generator type
- Update generation rate if performance changed

#### 6. CLI Configuration Changes
**Detection**: Modifications to `internal/config/` or `cmd/`

**Update Sections**:
- CLI Configuration (lines 43-65)
- Update flags table
- Update environment variables table
- Update Quick Start examples

### Low-Priority Triggers (MAY update SPEC.md)

#### 7. Test Coverage Changes (> 5% delta)
**Detection**: Running `go test ./... -cover` shows significant change

**Update Sections**:
- Test Coverage table (lines 367-378)
- Update coverage percentages by package
- Update overall percentage

#### 8. Performance Improvements (> 20% improvement)
**Detection**: Benchmarks show significant improvement

**Update Sections**:
- Performance Targets table (lines 354-362)
- Update "Actual" column
- Add decision log entry for major optimizations

## Decision Matrix

| Files Changed | Update SPEC.md? | Sections to Update |
|---------------|-----------------|-------------------|
| `internal/api/**/*.go` | **YES** | Endpoints table, API Reference |
| `internal/models/*.go` | **YES** | Response Format, Schema examples |
| `internal/services/*.go` | **YES** | Phase Features, Architecture, Package Structure |
| `internal/generator/*.go` | **MAYBE** | Generation Algorithm, Phase Features |
| `internal/storage/*.go` | **MAYBE** | Architecture diagram |
| `internal/config/*.go` | **YES** | CLI Configuration, Environment Variables |
| `cmd/**/*.go` | **YES** | CLI Configuration, Quick Start |
| `test/**/*.go` | **MAYBE** | Test Coverage table |
| Completing phase step | **YES** | Implementation Status, Phase Features |

## SPEC.md Update Checklist

When updating `services/cursor-sim/SPEC.md`:

### Always Check
- [ ] Line 5: **Last Updated** date is today
- [ ] Lines 17-22: **Implementation Status** table reflects current phase status

### If Endpoints Changed
- [ ] Lines 89-148: **Endpoints tables** include all new/modified endpoints
- [ ] Endpoint status: ✅ Implemented (not ⚡ Stub) if fully working

### If Models Changed
- [ ] Lines 175-200: **Schema examples** match actual struct fields
- [ ] Lines 199: **Invariants** still valid

### If Phase Completed
- [ ] Lines 17-22: Phase marked as **COMPLETE ✅** or **MOSTLY COMPLETE ✅**
- [ ] Phase Features section: Changed from "Planned" to "Implemented ✅"

### If Architecture Changed
- [ ] Lines 299-323: **Package Structure** includes new packages/directories
- [ ] File count approximations updated if significant changes

### If CLI Changed
- [ ] Lines 43-65: **CLI flags** table updated
- [ ] Environment variables table updated
- [ ] Quick Start examples still work

### If Decision Made
- [ ] Lines 437-449: **Decision Log** includes significant technical decisions

## Integration with SDD Workflow

This check happens at **Step 6: SYNC** in the enhanced SDD cycle:

```
4. REFACTOR → Clean up while tests pass
5. REFLECT  → Check dependency reflections
6. SYNC     → Update SPEC.md if triggered ← YOU ARE HERE
7. COMMIT   → Commit code + docs together
```

**Critical Rule**: If SPEC.md is updated, include it in the same commit as the code changes.

## Example Usage

### Scenario: Completed implementing `/analytics/team/models` endpoint

**Step 1: Check triggers**
- ✅ New endpoint added (High-Priority Trigger #2)
- ✅ Model generators modified (Medium-Priority Trigger #5)

**Step 2: Open SPEC.md and update**
```bash
# Line 5
**Last Updated**: January 3, 2026

# Lines 125 (in Team Analytics API table)
| GET | `/analytics/team/models` | Yes | ✅ Implemented |

# Lines 310 (in Package Structure)
├── generator/        # Event generation (26 files: commits, PRs, reviews, quality, models)
```

**Step 3: Stage with code changes**
```bash
git add internal/api/cursor/team.go
git add internal/generator/model_usage.go
git add services/cursor-sim/SPEC.md
git commit -m "feat(cursor-sim): implement /analytics/team/models endpoint"
```

## Quick Reference Card

**Before every commit, ask:**

1. Did I add/modify an endpoint? → Update Endpoints table
2. Did I complete a phase step? → Update Implementation Status
3. Did I change models? → Update Response Format
4. Did I add a package? → Update Package Structure
5. Did I change CLI? → Update CLI Configuration

**If answer is YES to any** → Update SPEC.md and include in commit

---

**Remember**: SPEC.md is the source of truth. Keep it current!
