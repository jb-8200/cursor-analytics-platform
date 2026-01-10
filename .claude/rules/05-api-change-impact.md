---
description: Enforces thorough review when cursor-sim API changes affect data pipeline
---

# API Change Impact Review Rule

## Context

cursor-sim (P4) is the authoritative source of truth for the analytics platform. Two downstream paths depend on its API contract:

```
Path 1 (GraphQL): cursor-sim → analytics-core (P5) → viz-spa (P6)
Path 2 (dbt):     cursor-sim → api-loader (P8) → dbt → streamlit-dashboard (P9)
```

When the cursor-sim API changes, **all downstream layers must be validated** to prevent:
- Data extraction failures
- Schema mismatches
- Column mapping errors
- Query failures
- Dashboard rendering issues

This rule enforces systematic impact analysis for Path 2 (data pipeline).

---

## NEVER

- **Never** modify cursor-sim API endpoints without creating an impact review task
- **Never** change API response format (`{items:[]}` structure) without validating api-loader
- **Never** rename or remove API fields without checking:
  - dbt column mapping (camelCase → snake_case)
  - DuckDB raw schema (main_raw.* tables)
  - Streamlit dashboard queries (main_mart.* references)
- **Never** add new API endpoints without updating data extraction plan
- **Never** change field types without validating dbt type casting
- **Never** merge API changes without E2E validation (cursor-sim → dbt → dashboard)
- **Never** skip SPEC.md update when API contract changes

---

## ALWAYS

### When Modifying cursor-sim API

**Step 1: Immediate Notification**
- Report to orchestrator: "⚠️ API CHANGE: [what changed]"
- Identify downstream impact: P8 (api-loader, dbt), P9 (Streamlit)

**Step 2: Create Impact Review Task**
- Feature directory: `.work-items/P4-F##-api-change-impact-review/`
- Create user-story.md, design.md, task.md
- Scope: Validate and update P8/P9 alignment

**Step 3: Document in SPEC.md**
- Update cursor-sim SPEC.md with API changes
- Mark sections affected: Endpoints, Response Schemas, Field Names

**Step 4: Validate Downstream Layers**

**P8 Data Tier Validation**:
- [ ] api-loader: Response format handling (dual format support)
- [ ] api-loader: Field extraction (all fields captured)
- [ ] DuckDB raw schema: Table columns match API fields (camelCase preserved)
- [ ] dbt staging: Column mapping correct (camelCase → snake_case)
- [ ] dbt marts: Aggregations still valid (no missing columns)

**P9 Dashboard Validation**:
- [ ] Streamlit queries: mart columns still available
- [ ] Dashboard pages: No KeyError exceptions
- [ ] Parameterized queries: No SQL syntax errors
- [ ] Data display: Charts and tables render correctly

**Step 5: Run E2E Tests**
```bash
# Full pipeline validation
1. Start cursor-sim with new API
2. Run api-loader extraction
3. Execute dbt transformations
4. Query dashboard endpoints
5. Verify all 4 dashboard pages load
```

**Step 6: Update Documentation**
- Data contract docs: `docs/data-contract-testing.md`
- Testing strategy: `docs/TESTING_STRATEGY.md`
- Architecture: `docs/design/new_data_architecture.md`

---

## Downstream Impact Matrix

| API Change Type | P8 Impact | P9 Impact | Validation Required |
|-----------------|-----------|-----------|---------------------|
| **New endpoint** | Add extractor in api-loader | Potentially add new queries | Review data availability |
| **Response format change** | Update format handling logic | N/A (consumes marts) | CRITICAL: Test extraction |
| **Field rename** | Update column mapping in dbt | Update query column refs | CRITICAL: Breaking change |
| **Field type change** | Add type casting in dbt | Validate DataFrame types | High priority |
| **Field removal** | Remove from staging models | Remove from queries | CRITICAL: Breaking change |
| **New field** | Optionally extract | Optionally visualize | Low priority (additive) |
| **Pagination change** | Update fetch logic | N/A | Medium priority |

---

## Example: API Field Rename

**Scenario**: Rename `commitHash` → `commit_sha` in cursor-sim API

**Impact Analysis**:

1. **P8 api-loader**:
   - ❌ Extraction breaks (expects `commitHash`)
   - ✅ Fix: Update field reference in extractors/base.py

2. **P8 DuckDB raw schema**:
   - ❌ Column missing (still expects `commitHash`)
   - ✅ Fix: Update raw table schema or keep old name

3. **P8 dbt staging**:
   - ❌ Column mapping breaks (commitHash → commit_hash)
   - ✅ Fix: Update staging model to map `commit_sha` → `commit_hash`

4. **P9 Streamlit**:
   - ✅ No direct impact (queries marts, not raw API)
   - ⚠️ Verify mart columns still populated

**Required Actions**:
- Create P4-F##-api-field-rename-impact task
- Update api-loader field references
- Update dbt staging column mapping
- Run E2E test: cursor-sim → dbt → dashboard
- Document in data contract testing guide

---

## Red Flags

If you see any of these patterns, **STOP and create impact review task**:

- "Modified cursor-sim endpoint response format"
- "Renamed API field from X to Y"
- "Removed field Z from API response"
- "Changed field type from string to integer"
- "Added new endpoint /analytics/new-data"
- API change merged without P8/P9 validation

---

## See Also

- **Rule 04-sdd-process.md**: REFLECT/SYNC phases
- **services/cursor-sim/SPEC.md**: API contract source of truth
- **docs/design/new_data_architecture.md**: Data pipeline architecture
- **docs/data-contract-testing.md**: Contract validation patterns
- **.claude/agents/cursor-sim-api-dev.md**: API agent constraints
