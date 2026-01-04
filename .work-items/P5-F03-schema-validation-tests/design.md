# Design Document: Schema Validation Tests

**Feature ID**: P5-F03
**Epic**: P5 - cursor-analytics-core (Testing Enhancement)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

Implement automated tests that validate the GraphQL schema structure, catching breaking changes before they affect consumers (P6).

---

## Test Categories

### 1. Type Existence Tests

Verify all required types are defined in schema.

```typescript
describe('Schema Type Existence', () => {
  const requiredTypes = [
    'Developer',
    'DeveloperStats',
    'DailyStats',
    'TeamStats',
    'DashboardKPI',
    'Commit',
    'PageInfo',
    'DeveloperConnection',
    'CommitConnection',
  ];

  requiredTypes.forEach(typeName => {
    it(`should have type ${typeName}`, () => {
      const type = schema.getType(typeName);
      expect(type).toBeDefined();
    });
  });
});
```

### 2. Field Structure Tests

Verify fields have correct types and nullability.

```typescript
describe('TeamStats Field Structure', () => {
  let teamStatsType: GraphQLObjectType;

  beforeAll(() => {
    teamStatsType = schema.getType('TeamStats') as GraphQLObjectType;
  });

  it('should have topPerformer as Developer (singular)', () => {
    const field = teamStatsType.getFields().topPerformer;
    expect(field).toBeDefined();
    // Should be Developer, not [Developer]
    expect(field.type.toString()).toBe('Developer');
  });

  it('should NOT have topPerformers field', () => {
    const fields = teamStatsType.getFields();
    expect(fields.topPerformers).toBeUndefined();
  });
});
```

### 3. Query Validation Tests

Verify queries execute against schema without errors.

```typescript
describe('Query Validation', () => {
  const queries = [
    {
      name: 'dashboardSummary',
      query: `{ dashboardSummary { totalDevelopers teamComparison { topPerformer { name } } } }`,
    },
    {
      name: 'developers',
      query: `{ developers(first: 10) { edges { node { id name } } } }`,
    },
  ];

  queries.forEach(({ name, query }) => {
    it(`should validate ${name} query`, async () => {
      const result = await graphql({ schema, source: query });
      expect(result.errors).toBeUndefined();
    });
  });
});
```

### 4. Invalid Query Rejection Tests

Verify schema rejects invalid queries.

```typescript
describe('Invalid Query Rejection', () => {
  it('should reject query with invalid field name', async () => {
    const query = `{
      dashboardSummary {
        teamComparison {
          topPerformers { name }  # WRONG: Should be topPerformer
        }
      }
    }`;

    const result = await graphql({ schema, source: query });
    expect(result.errors).toBeDefined();
    expect(result.errors[0].message).toContain('topPerformers');
  });
});
```

---

## File Structure

```
src/graphql/__tests__/
├── schema.validation.test.ts   # Type and field validation
├── queries.validation.test.ts  # Query execution tests
└── breaking-changes.test.ts    # Known breaking change prevention
```

---

## NPM Scripts

```json
{
  "scripts": {
    "test:schema": "vitest run --include '**/schema.*.test.ts'"
  }
}
```

---

## Test Implementation

```typescript
// src/graphql/__tests__/schema.validation.test.ts
import { GraphQLObjectType, GraphQLSchema, graphql } from 'graphql';
import { schema } from '../schema';

describe('GraphQL Schema Validation', () => {
  describe('Required Types', () => {
    const requiredTypes = [
      'Query',
      'Developer',
      'DeveloperStats',
      'DailyStats',
      'TeamStats',
      'DashboardKPI',
      'Commit',
      'PageInfo',
    ];

    requiredTypes.forEach(typeName => {
      it(`should define ${typeName} type`, () => {
        expect(schema.getType(typeName)).toBeDefined();
      });
    });
  });

  describe('TeamStats Type', () => {
    let teamStatsType: GraphQLObjectType;

    beforeAll(() => {
      teamStatsType = schema.getType('TeamStats') as GraphQLObjectType;
    });

    it('should have topPerformer field', () => {
      const fields = teamStatsType.getFields();
      expect(fields.topPerformer).toBeDefined();
    });

    it('topPerformer should be singular Developer', () => {
      const field = teamStatsType.getFields().topPerformer;
      // Unwrap nullable
      const typeName = field.type.toString().replace(/!/g, '');
      expect(typeName).toBe('Developer');
      expect(typeName).not.toContain('['); // Not an array
    });

    it('should NOT have topPerformers (plural)', () => {
      const fields = teamStatsType.getFields();
      expect(fields.topPerformers).toBeUndefined();
    });
  });

  describe('DailyStats Type', () => {
    let dailyStatsType: GraphQLObjectType;

    beforeAll(() => {
      dailyStatsType = schema.getType('DailyStats') as GraphQLObjectType;
    });

    it('should have linesAdded field', () => {
      expect(dailyStatsType.getFields().linesAdded).toBeDefined();
    });

    it('should NOT have humanLinesAdded field', () => {
      expect(dailyStatsType.getFields().humanLinesAdded).toBeUndefined();
    });
  });
});
```

---

## Success Metrics

| Metric | Before | After |
|--------|--------|-------|
| Schema structure coverage | 0% | 100% |
| Breaking change detection | Manual | Automated |
| Field name typos | Caught at runtime | Caught in tests |

---

## References

- [GraphQL Schema Introspection](https://graphql.org/learn/introspection/)
- `docs/e2e-testing-strategy.md` (Phase 1.2)
