# Task Breakdown: GraphQL Code Generator Setup

**Feature ID**: P6-F02
**Epic**: P6 - cursor-viz-spa (Service Decoupling Phase 1)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Progress Tracker

| Task | Description | Status | Time Est | Time Actual |
|------|-------------|--------|----------|-------------|
| TASK01 | Install GraphQL Code Generator dependencies | TODO | 0.5h | - |
| TASK02 | Create codegen.yml configuration | TODO | 0.5h | - |
| TASK03 | Add npm scripts (codegen, predev, prebuild) | TODO | 0.25h | - |
| TASK04 | Generate initial types from P5 schema | TODO | 0.25h | - |
| TASK05 | Update P6 queries to import generated types | TODO | 1.5h | - |
| TASK06 | Delete manual types.ts file | TODO | 0.25h | - |
| TASK07 | Add pre-commit hook for codegen validation | TODO | 0.5h | - |
| TASK08 | Add CI schema drift detection | TODO | 0.5h | - |
| TASK09 | Update documentation | TODO | 0.5h | - |

**Total Estimated**: 4.75 hours

---

## Task Details

### TASK01: Install GraphQL Code Generator Dependencies

**Objective**: Add GraphQL Code Generator and plugins to P6 devDependencies.

**Commands**:
```bash
cd services/cursor-viz-spa
npm install -D @graphql-codegen/cli @graphql-codegen/typescript \
  @graphql-codegen/typescript-operations @graphql-codegen/typed-document-node
```

**Acceptance Criteria**:
- [ ] All packages installed
- [ ] package.json updated with devDependencies
- [ ] package-lock.json updated

---

### TASK02: Create codegen.yml Configuration

**Objective**: Create configuration file for GraphQL Code Generator.

**File**: `services/cursor-viz-spa/codegen.yml`

**Configuration Requirements**:
- Schema: `http://localhost:4000/graphql`
- Documents: `src/**/*.{ts,tsx}`
- Output: `src/graphql/generated.ts`
- Plugins: typescript, typescript-operations, typed-document-node

**Acceptance Criteria**:
- [ ] codegen.yml created
- [ ] Schema endpoint configured
- [ ] Output path configured
- [ ] All required plugins configured
- [ ] Scalar types mapped (DateTime, Date)

---

### TASK03: Add npm Scripts

**Objective**: Add codegen scripts and hooks to package.json.

**Scripts to Add**:
```json
{
  "scripts": {
    "codegen": "graphql-codegen --config codegen.yml",
    "codegen:watch": "graphql-codegen --watch --config codegen.yml",
    "predev": "npm run codegen",
    "prebuild": "npm run codegen"
  }
}
```

**Acceptance Criteria**:
- [ ] `npm run codegen` works when P5 running
- [ ] `npm run dev` runs codegen first
- [ ] `npm run build` runs codegen first

---

### TASK04: Generate Initial Types

**Objective**: Run codegen to generate types from P5 schema.

**Prerequisites**:
- P5 must be running: `cd services/cursor-analytics-core && npm run dev`

**Commands**:
```bash
npm run codegen
```

**Acceptance Criteria**:
- [ ] `src/graphql/generated.ts` created
- [ ] File contains all P5 types (Developer, TeamStats, DashboardKPI, etc.)
- [ ] File contains typed document nodes for queries
- [ ] No TypeScript errors in generated file

---

### TASK05: Update P6 Queries to Import Generated Types

**Objective**: Replace manual type imports with generated type imports.

**Files to Update**:
- `src/graphql/queries.ts`
- `src/hooks/useDashboard.ts`
- `src/hooks/useDevelopers.ts`
- `src/hooks/useTeamStats.ts`
- Components that import types

**Changes**:
```typescript
// Before
import { TeamStats, Developer } from './types';

// After
import { TeamStats, Developer, GetDashboardSummaryQuery } from './generated';
```

**Acceptance Criteria**:
- [ ] All imports updated to use generated types
- [ ] TypeScript build passes (`npm run type-check`)
- [ ] All tests pass (`npm test`)
- [ ] Any schema mismatches fixed

**Expected Fixes** (based on known issues):
- `topPerformers` → `topPerformer`
- `humanLinesAdded` → `linesAdded`
- Remove `aiLinesDeleted` (doesn't exist)

---

### TASK06: Delete Manual types.ts File

**Objective**: Remove the now-obsolete manual type definitions.

**Commands**:
```bash
rm src/graphql/types.ts
```

**Verification**:
```bash
npm run build  # Should still work
npm test       # Should still pass
```

**Acceptance Criteria**:
- [ ] types.ts deleted
- [ ] No imports reference types.ts
- [ ] Build successful
- [ ] Tests pass

---

### TASK07: Add Pre-Commit Hook for Codegen Validation

**Objective**: Prevent committing when generated types are out of date.

**Option A - Using Husky**:
```bash
npm install -D husky
npx husky init
echo 'npm run codegen && git diff --exit-code src/graphql/generated.ts' > .husky/pre-commit
```

**Option B - Using package.json scripts**:
```json
{
  "scripts": {
    "codegen:check": "npm run codegen && git diff --exit-code src/graphql/generated.ts"
  }
}
```

**Acceptance Criteria**:
- [ ] Hook installed and functional
- [ ] Commit blocked if generated.ts differs
- [ ] Clear error message shown to developer

---

### TASK08: Add CI Schema Drift Detection

**Objective**: Fail CI if generated types differ from committed version.

**File**: `.github/workflows/ci.yml` (update existing or create new)

**CI Step**:
```yaml
- name: Validate GraphQL Schema
  run: |
    cd services/cursor-viz-spa
    npm run codegen
    git diff --exit-code src/graphql/generated.ts || {
      echo "❌ Generated types have changed!"
      echo "Regenerate with: npm run codegen"
      exit 1
    }
```

**Acceptance Criteria**:
- [ ] CI workflow updated
- [ ] Drift detection works (test with intentional mismatch)
- [ ] Clear error message in CI logs

---

### TASK09: Update Documentation

**Objective**: Document the new codegen workflow.

**Files to Update**:
- `services/cursor-viz-spa/README.md` - Add codegen section
- `docs/INTEGRATION.md` - Update P6 setup instructions

**Content to Add**:
- How to run codegen
- Prerequisites (P5 must be running)
- Troubleshooting common errors
- CI validation behavior

**Acceptance Criteria**:
- [ ] README updated with codegen instructions
- [ ] INTEGRATION.md updated with new workflow
- [ ] Troubleshooting section added

---

## Dependencies

| Task | Depends On |
|------|------------|
| TASK02 | TASK01 (packages installed) |
| TASK03 | TASK02 (config exists) |
| TASK04 | TASK03 (scripts exist), P5 running |
| TASK05 | TASK04 (types generated) |
| TASK06 | TASK05 (imports updated) |
| TASK07 | TASK06 (manual types deleted) |
| TASK08 | TASK07 (local validation works) |
| TASK09 | TASK08 (full workflow working) |

---

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| P5 not available during codegen | Document requirement, add error message |
| Generated types break existing code | Fix incrementally, test after each file |
| CI slowdown from codegen step | Cache node_modules, run only on P6 changes |
| Developer forgets to run codegen | predev hook runs automatically |

---

## Rollback Plan

1. Revert package.json and codegen.yml changes
2. Restore src/graphql/types.ts from git history
3. Update imports back to types.ts
4. Remove pre-commit hook

---

## Notes

- This is Phase 1 of the data contract testing mitigation plan
- Highest priority due to severity of integration issues
- Should be completed before any new P6 development
