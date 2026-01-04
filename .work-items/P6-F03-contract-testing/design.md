# Design Document: Contract Testing with GraphQL Inspector

**Feature ID**: P6-F03
**Epic**: P6 - cursor-viz-spa (Service Decoupling Phase 3)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

Implement GraphQL Inspector to validate P6 queries against P5's GraphQL schema, providing an additional layer of contract validation beyond TypeScript type checking.

---

## Architecture

### Validation Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│  Developer modifies P6 query                                            │
│  src/graphql/queries.ts                                                │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  Pre-Commit Hook: npm run schema:validate                              │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  graphql-inspector validate 'src/**/*.ts' 'http://localhost:4000'      │
│                                                                         │
│  ✅ All queries valid                                                   │
│  OR                                                                     │
│  ❌ Query error:                                                        │
│     src/graphql/queries.ts:42                                          │
│     Cannot query field "topPerformers" on type "TeamStats"             │
│     Did you mean "topPerformer"?                                       │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼ (if valid)
┌─────────────────────────────────────────────────────────────────────────┐
│  Commit proceeds                                                        │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Implementation

### 1. Dependencies

```json
{
  "devDependencies": {
    "@graphql-inspector/cli": "^4.0.0"
  }
}
```

### 2. Validation Script

**File**: `services/cursor-viz-spa/scripts/validate-queries.sh`

```bash
#!/bin/bash
set -e

echo "Checking P5 GraphQL server..."
curl -s http://localhost:4000/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ __typename }"}' \
  --fail > /dev/null 2>&1 || {
  echo "❌ P5 GraphQL server not running at localhost:4000!"
  echo "Start P5: cd services/cursor-analytics-core && npm run dev"
  exit 1
}

echo "Validating P6 queries against P5 schema..."
npx graphql-inspector validate \
  'src/**/*.ts' \
  'http://localhost:4000/graphql'

echo "✅ All queries valid!"
```

### 3. NPM Scripts

```json
{
  "scripts": {
    "schema:validate": "bash scripts/validate-queries.sh",
    "schema:diff": "graphql-inspector diff http://localhost:4000/graphql schema.graphql"
  }
}
```

### 4. Pre-Commit Hook (Husky)

```bash
# .husky/pre-commit
#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

# Only validate if GraphQL files changed
if git diff --cached --name-only | grep -q 'graphql\|queries\|mutations'; then
  echo "GraphQL changes detected, validating..."
  npm run schema:validate || {
    echo ""
    echo "❌ GraphQL query validation failed!"
    echo "Fix the errors above before committing."
    exit 1
  }
fi
```

---

## CI Integration

### GitHub Actions Step

```yaml
# Add to .github/workflows/ci.yml
- name: Start P5 for validation
  run: |
    cd services/cursor-analytics-core
    npm ci
    npm run dev &
    sleep 5  # Wait for server

- name: Validate GraphQL queries
  run: |
    cd services/cursor-viz-spa
    npm run schema:validate
```

---

## Error Examples

### Invalid Field Name
```
❌ Cannot query field "topPerformers" on type "TeamStats"
   Did you mean "topPerformer"?

   Location: src/graphql/queries.ts:42:8

   Query:
   41 |   teamComparison {
   42 |     topPerformers {  # ← ERROR
   43 |       name
```

### Missing Required Variable
```
❌ Variable "$teamId" is never used in operation "GetTeamStats"

   Location: src/graphql/queries.ts:15:1
```

### Type Mismatch
```
❌ Variable "$limit" of type "String!" used in position expecting type "Int!"

   Location: src/graphql/queries.ts:22:5
```

---

## Comparison: TypeScript vs GraphQL Inspector

| Validation | TypeScript | GraphQL Inspector |
|------------|------------|-------------------|
| Field names | ✅ (with codegen) | ✅ |
| Field types | ✅ | ✅ |
| Query syntax | ❌ | ✅ |
| Variable usage | ❌ | ✅ |
| Fragment matching | ❌ | ✅ |
| Real-time validation | ✅ (IDE) | ❌ (CLI only) |

**Recommendation**: Use both - TypeScript for IDE feedback, Inspector for pre-commit/CI.

---

## References

- [GraphQL Inspector Documentation](https://graphql-inspector.com/)
- [GraphQL Inspector CLI](https://graphql-inspector.com/docs/essentials/validate)
