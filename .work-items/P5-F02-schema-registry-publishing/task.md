# Task Breakdown: Schema Registry Publishing

**Feature ID**: P5-F02
**Epic**: P5 - cursor-analytics-core (Service Decoupling Phase 2)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Progress Tracker

| Task | Description | Status | Time Est | Time Actual |
|------|-------------|--------|----------|-------------|
| TASK01 | Set up Apollo Studio account | TODO | 0.5h | - |
| TASK02 | Install @apollo/rover in P5 | TODO | 0.25h | - |
| TASK03 | Create schema export script | TODO | 0.5h | - |
| TASK04 | Add schema publishing npm scripts | TODO | 0.25h | - |
| TASK05 | Configure CI for schema publishing | TODO | 1.0h | - |
| TASK06 | Configure CI for breaking change detection | TODO | 0.5h | - |
| TASK07 | Update P6 codegen to use registry | TODO | 0.5h | - |
| TASK08 | Document workflow | TODO | 0.5h | - |

**Total Estimated**: 4.0 hours

---

## Task Details

### TASK01: Set Up Apollo Studio Account

**Objective**: Create Apollo Studio account and configure graph.

**Steps**:
1. Go to https://studio.apollographql.com
2. Sign in with GitHub (project owner account)
3. Create new graph: `cursor-analytics`
4. Select graph type: "Self-Hosted GraphQL"
5. Select variant: `main`
6. Copy API key

**Secrets to Configure**:
- `APOLLO_KEY` in GitHub Actions secrets

**Acceptance Criteria**:
- [ ] Apollo Studio account created
- [ ] `cursor-analytics` graph created
- [ ] API key generated and stored in GitHub secrets

---

### TASK02: Install @apollo/rover in P5

**Objective**: Add Apollo Rover CLI for schema publishing.

**Commands**:
```bash
cd services/cursor-analytics-core
npm install -D @apollo/rover
```

**Acceptance Criteria**:
- [ ] @apollo/rover in devDependencies
- [ ] `npx rover --version` works

---

### TASK03: Create Schema Export Script

**Objective**: Create script to export GraphQL schema to SDL file.

**File**: `services/cursor-analytics-core/scripts/export-schema.ts`

**Script Content**:
```typescript
import { printSchema } from 'graphql';
import { schema } from '../src/graphql/schema';
import { writeFileSync, mkdirSync } from 'fs';

mkdirSync('./dist', { recursive: true });
const schemaSDL = printSchema(schema);
writeFileSync('./dist/schema.graphql', schemaSDL);
console.log('Schema exported to dist/schema.graphql');
```

**Acceptance Criteria**:
- [ ] Script created
- [ ] `npx ts-node scripts/export-schema.ts` generates dist/schema.graphql
- [ ] Generated file contains valid GraphQL SDL

---

### TASK04: Add Schema Publishing npm Scripts

**Objective**: Add npm scripts for schema management.

**Scripts to Add**:
```json
{
  "scripts": {
    "schema:export": "ts-node scripts/export-schema.ts",
    "schema:publish": "npm run schema:export && rover graph publish cursor-analytics@main --schema ./dist/schema.graphql",
    "schema:check": "npm run schema:export && rover graph check cursor-analytics@main --schema ./dist/schema.graphql"
  }
}
```

**Acceptance Criteria**:
- [ ] Scripts added to package.json
- [ ] `npm run schema:export` works
- [ ] `npm run schema:publish` works (with APOLLO_KEY set)
- [ ] `npm run schema:check` works (with APOLLO_KEY set)

---

### TASK05: Configure CI for Schema Publishing

**Objective**: Auto-publish schema on push to main.

**File**: `.github/workflows/schema-publish.yml`

**Trigger**: Push to main with changes in `services/cursor-analytics-core/src/graphql/**`

**Steps**:
1. Checkout code
2. Install dependencies
3. Export schema
4. Publish to Apollo Studio

**Acceptance Criteria**:
- [ ] Workflow file created
- [ ] Workflow triggers on P5 GraphQL changes
- [ ] Schema published successfully to Apollo Studio
- [ ] Version visible in Apollo Studio dashboard

---

### TASK06: Configure CI for Breaking Change Detection

**Objective**: Check for breaking changes on PRs.

**File**: `.github/workflows/schema-check.yml`

**Trigger**: Pull request with changes in `services/cursor-analytics-core/src/graphql/**`

**Behavior**:
- Run `rover graph check`
- Display breaking changes (warning, not blocking)
- Allow PR to merge with warning

**Acceptance Criteria**:
- [ ] Workflow file created
- [ ] Breaking changes displayed in PR
- [ ] Non-blocking (warning only)
- [ ] Clear messaging about affected operations

---

### TASK07: Update P6 Codegen to Use Registry

**Objective**: Configure P6 to fetch schema from Apollo Studio.

**File**: `services/cursor-viz-spa/codegen.yml`

**Changes**:
```yaml
schema:
  - apollo:
      graph: cursor-analytics@main
```

**Environment**:
- Add `APOLLO_KEY` to P6 CI environment

**Acceptance Criteria**:
- [ ] codegen.yml updated
- [ ] `npm run codegen` works without P5 running
- [ ] CI codegen works with Apollo Studio

---

### TASK08: Document Workflow

**Objective**: Document the schema registry workflow.

**Files to Update**:
- `services/cursor-analytics-core/README.md` - Add schema publishing section
- `docs/INTEGRATION.md` - Add Apollo Studio section
- Create `docs/schema-registry.md` - Full workflow documentation

**Content**:
- How to publish schema manually
- How to check for breaking changes
- How to view version history
- Troubleshooting guide

**Acceptance Criteria**:
- [ ] README updated
- [ ] INTEGRATION.md updated
- [ ] schema-registry.md created
- [ ] All commands documented

---

## Dependencies

| Task | Depends On |
|------|------------|
| TASK02 | TASK01 (account exists) |
| TASK03 | TASK02 (rover installed) |
| TASK04 | TASK03 (export script exists) |
| TASK05 | TASK04 (scripts exist), TASK01 (API key) |
| TASK06 | TASK05 (publishing works) |
| TASK07 | TASK06 (registry populated), P6-F02 (codegen setup) |
| TASK08 | All previous tasks |

---

## External Dependencies

| Dependency | Description |
|------------|-------------|
| Apollo Studio | Free for OSS projects |
| GitHub Secrets | Store APOLLO_KEY |
| P6-F02 | Must be complete for TASK07 |

---

## Notes

- This is Phase 2 of the data contract testing mitigation plan
- Depends on P6-F02 (GraphQL Code Generator) being complete
- Apollo Studio is free for open source projects
