# Data Contract Mismatch Mitigation Plan

**Created**: January 4, 2026
**Owner**: Platform Team
**Status**: ACTIVE
**Priority**: CRITICAL

## Executive Summary

This mitigation plan addresses the critical issue of data contract mismatches between services (cursor-sim, cursor-analytics-core, cursor-viz-spa) discovered during integration testing on January 4, 2026.

**Problem**: GraphQL schema drift between P5 (server) and P6 (client) caused complete integration failure.

**Solution**: Implement automated schema validation, code generation, and contract testing.

---

## Mitigation Strategies

### Strategy 1: Automated Type Generation (HIGHEST PRIORITY)

**Goal**: Eliminate manual type definitions in P6, auto-generate from P5 schema.

**Implementation**: GraphQL Code Generator

```bash
# services/cursor-viz-spa/
npm install -D @graphql-codegen/cli @graphql-codegen/typescript \
  @graphql-codegen/typescript-operations @graphql-codegen/typed-document-node

# Add to package.json
{
  "scripts": {
    "codegen": "graphql-codegen --config codegen.yml",
    "predev": "npm run codegen",
    "prebuild": "npm run codegen"
  }
}
```

**Impact**:
- ✅ Compile-time type checking
- ✅ No manual schema sync required
- ✅ IDE autocomplete for GraphQL queries
- ✅ Catches schema mismatches before `npm run dev`

**Effort**: 2-4 hours
**ROI**: Eliminates 100% of schema mismatch errors

---

### Strategy 2: Pre-Commit Schema Validation

**Goal**: Prevent committing P6 code with outdated schema.

**Implementation**: Husky + GraphQL Inspector

```bash
# services/cursor-viz-spa/
npm install -D @graphql-inspector/cli husky

# .husky/pre-commit
npm run codegen
git diff --exit-code src/graphql/generated.ts || {
  echo "❌ Generated types have changed!"
  echo "Run: npm run codegen && git add src/graphql/generated.ts"
  exit 1
}
```

**Impact**:
- ✅ Prevents drift from entering version control
- ✅ Forces developers to update types before commit
- ✅ CI fails if schema is out of sync

**Effort**: 1-2 hours
**ROI**: 100% prevention of schema drift in main branch

---

### Strategy 3: CI/CD Schema Compatibility Check

**Goal**: Validate P6 queries against P5 schema in every PR.

**Implementation**: GitHub Actions + GraphQL Validator

```yaml
# .github/workflows/schema-validation.yml
name: Schema Validation

on: [pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - name: Start P5 GraphQL Server
        run: |
          cd services/cursor-analytics-core
          docker-compose up -d
          sleep 5

      - name: Validate P6 Queries
        run: |
          cd services/cursor-viz-spa
          graphql-inspector validate 'src/**/*.ts' \
            'http://localhost:4000/graphql'
```

**Impact**:
- ✅ Blocks PRs with incompatible queries
- ✅ Automated validation (no manual review needed)
- ✅ Early detection (before merge)

**Effort**: 2-3 hours
**ROI**: Prevents production incidents

---

### Strategy 4: Schema Registry (Apollo Studio or Self-Hosted)

**Goal**: Centralized schema management with versioning.

**Implementation**: Apollo Studio (Free for OSS)

```bash
# P5: Publish schema on deploy
rover graph publish cursor-analytics@main \
  --schema ./src/graphql/schema.ts

# P6: Fetch schema before codegen
rover graph fetch cursor-analytics@main > schema.graphql
```

**Impact**:
- ✅ Schema versioning and history
- ✅ Breaking change detection
- ✅ Team notifications on schema changes
- ✅ Documentation auto-generated

**Effort**: 4-6 hours
**ROI**: Long-term schema governance

---

### Strategy 5: Contract Testing with GraphQL Inspector

**Goal**: Detect breaking changes before deployment.

**Implementation**: graphql-inspector diff

```bash
# Compare current schema vs production
graphql-inspector diff \
  https://api.cursor-analytics.com/graphql \
  ./src/graphql/schema.ts

# Output:
# ✅ Field 'topPerformer' was added (safe)
# ❌ Field 'topPerformers' was removed (BREAKING)
```

**Impact**:
- ✅ Identifies breaking vs. safe changes
- ✅ Prevents accidental breaking changes
- ✅ Enables safe schema evolution

**Effort**: 2-3 hours
**ROI**: Reduces rollback risk

---

## Enforcement Workflow

```
┌──────────────────────────────────────────────────────────────┐
│ 1. Developer modifies P5 schema                              │
│    └─ Edit: services/cursor-analytics-core/src/graphql/schema.ts
└──────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│ 2. P5 Tests Run                                              │
│    └─ npm run test:schema  (validates GraphQL syntax)       │
│    └─ ✅ Schema is valid                                     │
└──────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│ 3. Commit P5 Schema                                          │
│    └─ git commit -m "feat(p5): add topContributor field"    │
└──────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│ 4. P5 CI Publishes Schema to Apollo Studio                  │
│    └─ rover graph publish cursor-analytics@main             │
│    └─ ✅ Schema v1.2.3 published                             │
└──────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│ 5. P6 Developer Runs Codegen                                 │
│    └─ cd services/cursor-viz-spa                            │
│    └─ npm run dev  (triggers predev → codegen)              │
│    └─ Types auto-generated from latest P5 schema            │
└──────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│ 6. TypeScript Compile Errors (if incompatible)              │
│    └─ ❌ Property 'topPerformers' does not exist             │
│    └─ 💡 Did you mean 'topPerformer' or 'topContributor'?   │
└──────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│ 7. Developer Fixes P6 Code                                   │
│    └─ Update queries to use new field names                 │
│    └─ npm test  (all tests pass)                            │
└──────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│ 8. Pre-Commit Hook Validates                                │
│    └─ Checks generated.ts is committed                      │
│    └─ ✅ Schema in sync                                      │
└──────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│ 9. CI Schema Validation                                      │
│    └─ graphql-inspector validate                            │
│    └─ ✅ All queries valid against P5 schema                │
└──────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│ 10. Merge to Main                                            │
│     └─ ✅ Zero schema mismatch risk                          │
└──────────────────────────────────────────────────────────────┘
```

---

## Breaking Change Policy

### Safe Changes (Non-Breaking)

✅ **Adding optional fields**
✅ **Adding new types/queries**
✅ **Deprecating fields** (with `@deprecated`)
✅ **Widening types** (String → String | Int)

### Breaking Changes (Require Coordination)

❌ **Removing fields**
❌ **Changing field types**
❌ **Making optional fields required**
❌ **Renaming fields**

### Handling Breaking Changes

**Process**:
1. Add new field with new name
2. Deprecate old field
3. Update P6 to use new field
4. Wait 1 sprint for deployment
5. Remove deprecated field

**Example**:
```graphql
# Sprint N
type TeamStats {
  topPerformers: [Developer] @deprecated(reason: "Use topPerformer")
  topPerformer: Developer  # NEW
}

# Sprint N+1
type TeamStats {
  topPerformer: Developer  # OLD field removed
}
```

---

## Monitoring & Alerts

### Daily Schema Drift Check

```yaml
# .github/workflows/schema-drift.yml
name: Daily Schema Drift Check
on:
  schedule:
    - cron: '0 9 * * *'  # 9am daily

jobs:
  check-drift:
    steps:
      - name: Compare schemas
        run: |
          graphql-inspector diff \
            https://api.cursor-analytics.com/graphql \
            services/cursor-viz-spa/schema.graphql || {
            echo "⚠️ Schema drift detected!"
            exit 1
          }

      - name: Slack Alert
        if: failure()
        run: |
          curl -X POST https://hooks.slack.com/... \
            -d '{"text":"Schema drift detected between P5 and P6"}'
```

---

## Rollout Plan

### Week 1: Quick Wins
- [x] Document integration issues (DONE)
- [x] Create mitigation plan (DONE)
- [ ] Install GraphQL Code Generator in P6
- [ ] Add predev hook to run codegen
- [ ] Delete manual types.ts file

### Week 2: Validation
- [ ] Add pre-commit hook for schema validation
- [ ] Set up CI schema compatibility check
- [ ] Add tests to verify generated types

### Week 3: Registry
- [ ] Set up Apollo Studio account
- [ ] Configure P5 to publish schema
- [ ] Configure P6 to fetch schema
- [ ] Add breaking change detection

### Week 4: Monitoring
- [ ] Set up daily drift detection
- [ ] Configure Slack alerts
- [ ] Document workflow in team wiki
- [ ] Train team on new process

---

## Success Criteria

### Before (Current State)

❌ Manual schema synchronization required
❌ Schema mismatches discovered at runtime (browser 400 errors)
❌ No validation until full stack running
❌ TypeScript provides false sense of type safety

### After (Target State)

✅ Auto-generated types from P5 schema
✅ Compile-time errors for schema mismatches
✅ Pre-commit validation prevents drift
✅ CI blocks incompatible queries
✅ Zero runtime schema errors

---

## Cost-Benefit Analysis

| Strategy | Effort | Cost | Benefit | ROI |
|----------|--------|------|---------|-----|
| Code Generator | 2-4h | Free | Eliminates 100% of type drift | ⭐⭐⭐⭐⭐ |
| Pre-Commit Hook | 1-2h | Free | Prevents drift in main branch | ⭐⭐⭐⭐⭐ |
| CI Validation | 2-3h | Free | Blocks bad PRs | ⭐⭐⭐⭐ |
| Apollo Studio | 4-6h | Free (OSS) | Long-term governance | ⭐⭐⭐⭐ |
| GraphQL Inspector | 2-3h | Free | Breaking change detection | ⭐⭐⭐ |

**Total Implementation Time**: ~12-18 hours
**Total Ongoing Cost**: $0 (all tools free for open source)
**Risk Reduction**: Eliminates #1 cause of integration failures

---

## References

- **Data Contract Testing Strategy**: `docs/data-contract-testing.md`
- **E2E Testing Strategy**: `docs/e2e-testing-strategy.md`
- **Integration Guide**: `docs/INTEGRATION.md`
- **GraphQL Code Generator**: https://the-guild.dev/graphql/codegen
- **Apollo Studio**: https://studio.apollographql.com
- **GraphQL Inspector**: https://graphql-inspector.com

---

## Appendix: Lessons Learned

### What Went Wrong

1. **Manual Type Definitions**: P6 types manually written, not generated
2. **No Validation**: No automated schema validation before runtime
3. **Late Detection**: Errors only discovered during full stack integration
4. **TypeScript False Security**: Local types matched, but not server schema

### What Went Right

1. **Comprehensive Testing**: Unit and component tests caught 90%+ of bugs
2. **Docker Architecture**: Container networking simplified integration
3. **Clear Documentation**: Easy to trace issues and implement fixes
4. **SDD Workflow**: Systematic approach enabled rapid diagnosis

### Key Takeaways

> **"Never manually define GraphQL types in client code. Always generate from server schema."**

> **"TypeScript type safety is only as good as your runtime contract validation."**

> **"Automate everything that can drift: types, schemas, contracts, tests."**

---

**Status**: Mitigation plan approved and ready for implementation.
**Next Action**: Install GraphQL Code Generator (Week 1, Priority 1)
