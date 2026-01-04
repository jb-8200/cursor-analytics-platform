# Design Document: Schema Registry Publishing

**Feature ID**: P5-F02
**Epic**: P5 - cursor-analytics-core (Service Decoupling Phase 2)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

Set up Apollo Studio schema registry for cursor-analytics-core to enable centralized schema management, version history, and breaking change detection.

---

## Architecture

### Schema Registry Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           cursor-analytics-core (P5)                    │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  src/graphql/schema.ts (Source of Truth)                               │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │ type Query {                                                      │  │
│  │   dashboardSummary: DashboardKPI                                  │  │
│  │   developers(...): DeveloperConnection                            │  │
│  │ }                                                                 │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  CI Pipeline: rover graph publish                                      │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Publish on every main push
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        Apollo Studio Registry                          │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  Graph: cursor-analytics@main                                          │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │ Version History:                                                  │  │
│  │ • v1.2.3 (current) - Added topContributor field                  │  │
│  │ • v1.2.2 - Fixed dailyStats type                                 │  │
│  │ • v1.2.1 - Added pagination to developers                        │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  Features:                                                             │
│  • Breaking change detection                                           │
│  • Schema documentation                                                │
│  • GraphQL playground                                                  │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Fetch for codegen
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         cursor-viz-spa (P6)                             │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  codegen.yml:                                                          │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │ # Fetch from registry instead of local P5                        │  │
│  │ schema:                                                           │  │
│  │   - apollo:                                                       │  │
│  │       graph: cursor-analytics@main                                │  │
│  └───────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Apollo Studio Setup

### 1. Create Account and Graph

1. Go to https://studio.apollographql.com
2. Sign in with GitHub
3. Create new graph: `cursor-analytics`
4. Select variant: `main`
5. Copy API key for CI

### 2. Install Rover CLI

```bash
# P5
cd services/cursor-analytics-core
npm install -D @apollo/rover
```

### 3. Configure Environment

```bash
# Set in CI secrets
APOLLO_KEY=service:cursor-analytics:xxxxx
```

### 4. NPM Scripts

```json
{
  "scripts": {
    "schema:export": "ts-node scripts/export-schema.ts",
    "schema:publish": "rover graph publish cursor-analytics@main --schema ./dist/schema.graphql",
    "schema:check": "rover graph check cursor-analytics@main --schema ./dist/schema.graphql"
  }
}
```

---

## Schema Export Script

```typescript
// services/cursor-analytics-core/scripts/export-schema.ts
import { printSchema } from 'graphql';
import { schema } from '../src/graphql/schema';
import { writeFileSync } from 'fs';

const schemaSDL = printSchema(schema);
writeFileSync('./dist/schema.graphql', schemaSDL);
console.log('Schema exported to dist/schema.graphql');
```

---

## CI Integration

### GitHub Actions Workflow

```yaml
# .github/workflows/schema-publish.yml
name: Publish Schema

on:
  push:
    branches: [main]
    paths:
      - 'services/cursor-analytics-core/src/graphql/**'

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'

      - name: Install dependencies
        run: |
          cd services/cursor-analytics-core
          npm ci

      - name: Export schema
        run: |
          cd services/cursor-analytics-core
          npm run schema:export

      - name: Publish to Apollo Studio
        env:
          APOLLO_KEY: ${{ secrets.APOLLO_KEY }}
        run: |
          cd services/cursor-analytics-core
          npm run schema:publish
```

### Schema Check on PR

```yaml
# .github/workflows/schema-check.yml
name: Schema Check

on:
  pull_request:
    paths:
      - 'services/cursor-analytics-core/src/graphql/**'

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'

      - name: Install dependencies
        run: |
          cd services/cursor-analytics-core
          npm ci

      - name: Export schema
        run: |
          cd services/cursor-analytics-core
          npm run schema:export

      - name: Check for breaking changes
        env:
          APOLLO_KEY: ${{ secrets.APOLLO_KEY }}
        run: |
          cd services/cursor-analytics-core
          npm run schema:check || {
            echo "⚠️ Breaking changes detected!"
            echo "Review the changes above and ensure P6 is updated."
            exit 0  # Warning only, not blocking
          }
```

---

## P6 Configuration Update

After registry is set up, update P6's codegen.yml:

```yaml
# services/cursor-viz-spa/codegen.yml
overwrite: true

# Option 1: Fetch from Apollo Studio (production)
schema:
  - apollo:
      graph: cursor-analytics@main
      # Uses APOLLO_KEY from environment

# Option 2: Fallback to local P5 (development)
# schema: "http://localhost:4000/graphql"

documents: "src/**/*.{ts,tsx}"
generates:
  src/graphql/generated.ts:
    plugins:
      - "typescript"
      - "typescript-operations"
      - "typed-document-node"
```

---

## Breaking Change Detection

Apollo Studio automatically detects breaking changes:

### Safe Changes (No Alert)
- Adding optional fields
- Adding new types/queries
- Deprecating fields

### Breaking Changes (Alert)
- Removing fields
- Changing field types
- Making optional fields required
- Renaming fields

### Example CI Output

```
❌ BREAKING CHANGES DETECTED

Removed fields:
  - TeamStats.topPerformers (was: [Developer])

Changed types:
  - DailyStats.linesAdded: Int → Int! (now required)

Affected operations:
  - GetDashboardSummary (uses topPerformers)
  - GetTeamStats (uses linesAdded)

⚠️ These changes may break cursor-viz-spa (P6)
```

---

## Success Metrics

| Metric | Before | After |
|--------|--------|-------|
| Schema versioning | None | Full history |
| Breaking change detection | Manual review | Automated |
| P6 codegen dependency | P5 running locally | Registry (always available) |
| Documentation | Manual | Auto-generated |

---

## References

- [Apollo Studio Documentation](https://www.apollographql.com/docs/studio/)
- [Rover CLI Reference](https://www.apollographql.com/docs/rover/)
- [Schema Checks](https://www.apollographql.com/docs/studio/schema-checks/)
