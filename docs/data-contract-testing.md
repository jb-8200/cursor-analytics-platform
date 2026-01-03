# Data Contract Testing & Schema Validation Strategy

**Created**: January 4, 2026
**Status**: PROPOSED
**Priority**: CRITICAL

## Executive Summary

During P5/P6 integration (Jan 4, 2026), GraphQL schema mismatches caused complete integration failure. P6 manually defined types that didn't match P5's actual schema, resulting in 400 Bad Request errors. This document proposes a comprehensive strategy to prevent data contract drift between services.

---

## Problem Statement

### Current State (Broken)

```
┌─────────────────────────────┐       ┌──────────────────────────────┐
│  cursor-analytics-core (P5) │       │  cursor-viz-spa (P6)         │
│  ─────────────────────────  │       │  ──────────────────────────  │
│                             │       │                              │
│  src/graphql/schema.ts      │       │  src/graphql/types.ts        │
│  ┌────────────────────────┐ │       │  ┌─────────────────────────┐ │
│  │ type TeamStats {       │ │       │  │ interface TeamStats {   │ │
│  │   topPerformer: Dev    │ │  ❌   │  │   topPerformers: Dev[]  │ │
│  │ }                      │ │       │  │ }                       │ │
│  └────────────────────────┘ │       │  └─────────────────────────┘ │
│                             │       │                              │
│  (Source of Truth)          │       │  (Manual, Drifts)            │
└─────────────────────────────┘       └──────────────────────────────┘

                                      Runtime Error! ⚠️
                                      400 Bad Request
```

**Issues:**
1. **Manual Type Definitions**: P6 types manually written, based on outdated design docs
2. **No Validation**: TypeScript validates against local types, not server schema
3. **Late Detection**: Errors only appear at runtime, in browser
4. **Breaking Changes**: P5 schema changes break P6 silently

---

## Proposed Solution: Schema-First Development

### Target State (Fixed)

```
┌─────────────────────────────────────────────────────────────────────┐
│  cursor-analytics-core (P5)                                         │
│  ─────────────────────────────────────────────────────────────────  │
│                                                                     │
│  src/graphql/schema.ts (Single Source of Truth)                    │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ type TeamStats {                                               │ │
│  │   teamName: String!                                            │ │
│  │   topPerformer: Developer    # Singular                        │ │
│  │ }                                                              │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                     │
│  Step 1: Introspection Query                                       │
│  ▼                                                                  │
│  dist/schema.graphql (Generated)                                   │
└─────────────────────────────────────────────────────────────────────┘
                           │
                           │ Step 2: Copy schema to P6
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│  cursor-viz-spa (P6)                                                │
│  ─────────────────────────────────────────────────────────────────  │
│                                                                     │
│  schema.graphql (Copied from P5)                                   │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ type TeamStats {                                               │ │
│  │   teamName: String!                                            │ │
│  │   topPerformer: Developer                                      │ │
│  │ }                                                              │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                     │
│  Step 3: GraphQL Code Generator                                    │
│  ▼                                                                  │
│  src/graphql/generated.ts (Auto-generated)                         │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ export interface TeamStats {                                   │ │
│  │   teamName: string;                                            │ │
│  │   topPerformer?: Developer;  // Auto-generated ✅               │ │
│  │ }                                                              │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                     │
│  Step 4: Use Generated Types                                       │
│  ▼                                                                  │
│  src/graphql/queries.ts                                            │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ import { TeamStats } from './generated';                       │ │
│  │                                                                │ │
│  │ const GET_TEAMS = gql`                                         │ │
│  │   query {                                                      │ │
│  │     teams { topPerformer { name } }  # TypeScript validates!  │ │
│  │   }                                                            │ │
│  │ `;                                                             │ │
│  └────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘

                         ✅ Compile-time validation
                         ✅ Auto-complete in IDE
                         ✅ No runtime errors
```

---

## Implementation Plan

### Phase 1: GraphQL Code Generator (High Priority)

**Goal**: Auto-generate P6 TypeScript types from P5 GraphQL schema.

#### Step 1.1: Install GraphQL Code Generator in P6

```bash
cd services/cursor-viz-spa
npm install -D @graphql-codegen/cli @graphql-codegen/typescript @graphql-codegen/typescript-operations @graphql-codegen/typed-document-node
```

#### Step 1.2: Create `codegen.yml` in P6

```yaml
# services/cursor-viz-spa/codegen.yml
overwrite: true
schema: "http://localhost:4000/graphql"  # P5 GraphQL endpoint
documents: "src/**/*.{ts,tsx}"
generates:
  src/graphql/generated.ts:
    plugins:
      - "typescript"
      - "typescript-operations"
      - "typed-document-node"
    config:
      avoidOptionals: false
      maybeValue: T | null | undefined
      inputMaybeValue: T | null | undefined
```

#### Step 1.3: Add npm scripts

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

#### Step 1.4: Update P6 queries to use generated types

**Before (Manual):**
```typescript
// src/graphql/types.ts (DELETE THIS FILE)
export interface TeamStats {
  topPerformers: Developer[];  // ❌ Wrong!
}
```

**After (Generated):**
```typescript
// src/graphql/queries.ts
import { gql } from '@apollo/client';
import type { GetDashboardSummaryQuery } from './generated';

export const GET_DASHBOARD_SUMMARY = gql`
  query GetDashboardSummary {
    dashboardSummary {
      teamComparison {
        topPerformer { name }  # ✅ TypeScript validates this field exists
      }
    }
  }
`;
```

#### Step 1.5: CI/CD Integration

Add to P6's CI pipeline:

```yaml
# .github/workflows/ci.yml
- name: Validate GraphQL Schema
  run: |
    cd services/cursor-viz-spa
    npm run codegen
    git diff --exit-code src/graphql/generated.ts
```

**Benefit**: CI fails if generated types have uncommitted changes (schema drift).

---

### Phase 2: Schema Registry (Medium Priority)

**Goal**: Centralize schema management and enable breaking change detection.

#### Option 2A: Apollo Studio (Recommended)

**Pros**:
- Free for open source
- Schema versioning
- Breaking change detection
- GraphQL playground
- Performance monitoring

**Setup**:

```bash
# P5: Publish schema on every deploy
cd services/cursor-analytics-core
npm install -D @apollo/rover

# Add to package.json
{
  "scripts": {
    "schema:publish": "rover graph publish cursor-analytics@main --schema ./src/graphql/schema.ts"
  }
}

# P6: Fetch schema before codegen
cd services/cursor-viz-spa
{
  "scripts": {
    "schema:fetch": "rover graph fetch cursor-analytics@main > schema.graphql",
    "precodegen": "npm run schema:fetch"
  }
}
```

#### Option 2B: GraphQL Inspector (Self-Hosted)

**Pros**:
- No external dependencies
- Git-based workflow
- Pre-commit validation

**Setup**:

```bash
npm install -D @graphql-inspector/cli

# Add to package.json
{
  "scripts": {
    "schema:validate": "graphql-inspector diff schema.graphql.old schema.graphql"
  }
}
```

---

### Phase 3: Contract Testing (Medium Priority)

**Goal**: Validate that P6 queries are compatible with P5 schema **before** deployment.

#### Step 3.1: Install GraphQL Inspector

```bash
cd services/cursor-viz-spa
npm install -D @graphql-inspector/cli
```

#### Step 3.2: Create validation script

```bash
# scripts/validate-queries.sh
#!/bin/bash
set -e

echo "Fetching P5 schema..."
curl -s http://localhost:4000/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ __schema { types { name } } }"}' \
  --fail > /dev/null || {
  echo "❌ P5 GraphQL server not running!"
  exit 1
}

echo "Validating P6 queries against P5 schema..."
graphql-inspector validate 'src/**/*.ts' 'http://localhost:4000/graphql'

echo "✅ All queries valid!"
```

#### Step 3.3: Add to pre-commit hook

```bash
# .husky/pre-commit
#!/bin/sh
cd services/cursor-viz-spa
npm run schema:validate || {
  echo ""
  echo "❌ GraphQL schema validation failed!"
  echo "Your queries don't match the P5 schema."
  echo ""
  echo "Fix by running: npm run codegen"
  exit 1
}
```

---

### Phase 4: Visual Regression Testing (Low Priority)

**Goal**: Catch when UI renders placeholders instead of real data.

#### Step 4.1: Install Playwright

```bash
cd services/cursor-viz-spa
npm install -D @playwright/test
npx playwright install
```

#### Step 4.2: Create E2E test with visual snapshots

```typescript
// tests/e2e/dashboard.spec.ts
import { test, expect } from '@playwright/test';

test('dashboard loads with real data', async ({ page }) => {
  // Start P5 in Docker (prerequisite)
  await page.goto('http://localhost:3000/dashboard');

  // Wait for GraphQL request to complete
  await page.waitForResponse(resp =>
    resp.url().includes('/graphql') && resp.status() === 200
  );

  // Verify KPI cards show numbers (not placeholders)
  await expect(page.locator('[data-testid="total-devs"]'))
    .not.toContainText('Chart placeholder');

  // Take screenshot
  await expect(page).toHaveScreenshot('dashboard.png');
});
```

#### Step 4.3: Run in CI

```yaml
# .github/workflows/e2e.yml
- name: Start P5 in Docker
  run: |
    cd services/cursor-analytics-core
    docker-compose up -d
    sleep 5

- name: Run E2E tests
  run: |
    cd services/cursor-viz-spa
    npx playwright test
```

---

## Mitigation Plan: Avoiding Schema Drift

### Development Workflow

```
┌─────────────────────────────────────────────────────────────────────┐
│ 1. Developer modifies P5 schema                                     │
│    ├─ Edit: services/cursor-analytics-core/src/graphql/schema.ts   │
│    └─ Example: Add field `topContributor: Developer`               │
└─────────────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 2. P5 tests run (includes schema validation)                        │
│    ├─ npm test                                                      │
│    └─ ✅ Schema is valid GraphQL                                    │
└─────────────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 3. Commit P5 changes                                                │
│    ├─ git add src/graphql/schema.ts                                │
│    └─ git commit -m "feat(p5): add topContributor field"           │
└─────────────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 4. P5 CI publishes schema to registry (Apollo Studio)              │
│    ├─ rover graph publish cursor-analytics@main                    │
│    └─ ✅ Schema version v1.2.3 published                            │
└─────────────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 5. P6 developer pulls latest and runs codegen                      │
│    ├─ git pull origin main                                         │
│    ├─ cd services/cursor-viz-spa                                   │
│    └─ npm run codegen  # Auto-fetches latest P5 schema             │
└─────────────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 6. TypeScript errors appear immediately                            │
│    ├─ src/components/Dashboard.tsx:42:18                           │
│    └─ ❌ Property 'topPerformers' does not exist on type 'TeamStats'│
│        Did you mean 'topPerformer' or 'topContributor'?            │
└─────────────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 7. Developer fixes P6 code to use new field                        │
│    ├─ Update queries: topPerformer → topContributor                │
│    ├─ npm run test  # All tests pass                               │
│    └─ git commit -m "feat(p6): use new topContributor field"       │
└─────────────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 8. P6 CI validates queries against P5 schema                       │
│    ├─ graphql-inspector validate src/**/*.ts http://localhost:4000 │
│    └─ ✅ All queries valid                                          │
└─────────────────────────────────────────────────────────────────────┘
```

---

### Breaking Change Policy

**Rule**: P5 schema changes MUST be backward compatible OR coordinated with P6.

#### Safe Changes (Non-Breaking)

✅ **Adding optional fields**:
```graphql
type TeamStats {
  teamName: String!
  topPerformer: Developer
  topContributor: Developer  # NEW: Optional field (safe)
}
```

✅ **Adding new types**:
```graphql
type NewFeature {  # NEW: Entire type (safe)
  id: ID!
}
```

✅ **Deprecating fields** (with `@deprecated`):
```graphql
type TeamStats {
  topPerformer: Developer @deprecated(reason: "Use topContributor instead")
  topContributor: Developer
}
```

#### Breaking Changes (Require Coordination)

❌ **Removing fields**:
```graphql
type TeamStats {
  # topPerformer: Developer  # REMOVED: P6 queries will break!
  topContributor: Developer
}
```

❌ **Changing field types**:
```graphql
type TeamStats {
  topPerformer: Developer  # Was: Developer, Now: [Developer] - BREAKS QUERIES
}
```

❌ **Making optional fields required**:
```graphql
type TeamStats {
  topPerformer: Developer!  # Was optional, now required - BREAKS QUERIES
}
```

#### Handling Breaking Changes

**Process**:
1. **Announce** breaking change in team chat/PR
2. **Deprecate** old field first (keep both for 1 sprint)
3. **Update P6** to use new field
4. **Remove** old field after P6 deployed

**Example**:
```graphql
# Sprint N: Add new field, deprecate old
type TeamStats {
  topPerformers: [Developer] @deprecated(reason: "Use topPerformer (singular)")
  topPerformer: Developer  # NEW
}

# Sprint N+1: Remove deprecated field
type TeamStats {
  topPerformer: Developer
}
```

---

## Monitoring & Alerts

### Schema Drift Detection

**Automated Check (Daily)**:
```bash
# .github/workflows/schema-drift.yml
name: Schema Drift Detection
on:
  schedule:
    - cron: '0 9 * * *'  # Daily at 9am

jobs:
  check-drift:
    runs-on: ubuntu-latest
    steps:
      - name: Fetch P5 schema from production
        run: curl https://api.cursor-analytics.com/graphql > p5-prod.graphql

      - name: Fetch P6 schema from repo
        run: cat services/cursor-viz-spa/schema.graphql > p6-repo.graphql

      - name: Compare schemas
        run: |
          graphql-inspector diff p6-repo.graphql p5-prod.graphql || {
            echo "❌ Schema drift detected!"
            echo "P6 schema doesn't match P5 production schema."
            exit 1
          }
```

**Slack Alert**:
```yaml
      - name: Send Slack alert on drift
        if: failure()
        uses: slackapi/slack-github-action@v1
        with:
          payload: |
            {
              "text": "⚠️ Schema drift detected between P5 and P6!",
              "blocks": [
                {
                  "type": "section",
                  "text": {
                    "type": "mrkdwn",
                    "text": "P6 schema is out of sync with P5 production.\nRun `npm run codegen` in cursor-viz-spa."
                  }
                }
              ]
            }
```

---

## Success Metrics

### Before (Current State)

❌ **Manual schema sync** required after every P5 change
❌ **Runtime errors** in browser (400 Bad Request)
❌ **Integration failures** discovered during manual testing
❌ **No validation** until full stack running

### After (Target State)

✅ **Auto-generated types** from P5 schema (no manual sync)
✅ **Compile-time errors** in IDE (before `git commit`)
✅ **CI failures** on schema drift (before merge)
✅ **Zero runtime** GraphQL schema errors

---

## Rollout Plan

### Week 1: Quick Wins (Phase 1)
- [ ] Install GraphQL Code Generator in P6
- [ ] Configure codegen.yml to fetch from P5
- [ ] Run `npm run codegen` and commit generated types
- [ ] Update P6 to import from `generated.ts` instead of `types.ts`
- [ ] Add pre-commit hook to validate codegen is up-to-date

### Week 2: Automation (Phase 2)
- [ ] Set up Apollo Studio schema registry (or GraphQL Inspector)
- [ ] Configure P5 to publish schema on deploy
- [ ] Configure P6 to fetch schema before codegen
- [ ] Add CI check for schema compatibility

### Week 3: Testing (Phase 3 & 4)
- [ ] Add contract tests with GraphQL Inspector
- [ ] Set up Playwright for E2E tests
- [ ] Add visual regression tests for Dashboard
- [ ] Document testing workflow in INTEGRATION.md

### Week 4: Monitoring (Phase 4)
- [ ] Set up daily schema drift detection
- [ ] Configure Slack alerts
- [ ] Create schema change SOP document
- [ ] Train team on new workflow

---

## Appendix: Tools Comparison

| Tool | Purpose | Pros | Cons | Recommendation |
|------|---------|------|------|----------------|
| **GraphQL Code Generator** | Auto-generate TypeScript types | Industry standard, great DX | Requires running P5 locally | ✅ Use |
| **Apollo Studio** | Schema registry + monitoring | Free for OSS, full-featured | External dependency | ✅ Use |
| **GraphQL Inspector** | Schema validation | Git-based, self-hosted | Less features than Studio | ⚠️ Alternative |
| **Playwright** | E2E + visual regression | Fast, reliable | Requires running full stack | ✅ Use |
| **Pact** | Consumer-driven contracts | Language agnostic | Steep learning curve | ❌ Skip (overkill) |

---

## References

- [GraphQL Code Generator Docs](https://the-guild.dev/graphql/codegen)
- [Apollo Studio Schema Registry](https://www.apollographql.com/docs/studio/schema-registry/)
- [GraphQL Inspector](https://graphql-inspector.com/)
- [Playwright Visual Comparisons](https://playwright.dev/docs/test-snapshots)
