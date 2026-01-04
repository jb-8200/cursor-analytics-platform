---
description: analytics-core specific guardrails for GraphQL schema and database safety
paths: services/cursor-analytics-core/**
---

# analytics-core Rules

Service-specific constraints for P5 (analytics-core GraphQL aggregator).

---

## NEVER

- **Manually define GraphQL types** in client code (only in schema)
- **Add breaking changes** to GraphQL schema without major version bump
- **Expose internal Prisma models** directly in GraphQL responses
- **Execute raw SQL** (use Prisma ORM always)
- **Cache GraphQL queries** without TTL and invalidation strategy
- **Query without pagination** on list fields (P6 dashboard needs bounds)
- **Implement mutations** without response types including errors

---

## ALWAYS

- **Keep GraphQL schema backward compatible** unless intentional major change
- **Document breaking changes** clearly for P6 consumers
- **Use Prisma for all database access** (type safety + migrations)
- **Include error handling** in resolver responses
- **Use cursor-based pagination** (not offset) for consistency
- **Type all resolver parameters** explicitly
- **Test queries end-to-end** with real PostgreSQL database
- **Verify P6 alignment** before schema changes

---

## GraphQL Schema Design

### Null Safety
- **Required fields**: Use `!` intentionally (not by default)
- **Optional fields**: Clearly document why nullable
- **Lists**: Decide if empty list or null (prefer empty)

### Pagination
- **Cursor-based pagination only**: No offset pagination
- **Standard connection pattern**: `{edges, pageInfo, totalCount}`
- **Page info**: `{hasNextPage, hasPreviousPage, startCursor, endCursor}`
- **Example**:
```graphql
type DeveloperConnection {
  edges: [DeveloperEdge!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}
```

### Mutations
- **Return operation type**: Include result type with success/error
- **Example**:
```graphql
type UpdateDeveloperResult {
  success: Boolean!
  developer: Developer
  error: String
}
```

### Error Handling
- **Response errors**: Include error field in result type
- **HTTP errors**: Throw with context
- **GraphQL errors**: Message must be user-friendly

---

## Database (Prisma) Rules

### Migrations
- **Always use Prisma migrations**: `prisma migrate dev`
- **Track migrations in version control**: `prisma/migrations/`
- **Test migrations**: Up and down
- **Document breaking changes**: Add migration notes

### Queries
- **Use Prisma client**: No raw SQL
- **Select only needed fields**: Performance
- **Validate input**: At resolver boundary
- **Handle null responses**: Gracefully

### Types
- **Generated from schema.prisma**: Always up-to-date
- **Use Prisma types**: Not custom interfaces
- **Type Prisma results**: Explicitly in resolvers

---

## Data Flow

### From cursor-sim (P4) to analytics-core (P5)
1. Fetch data via cursor-sim API (CursorSimClient)
2. Validate response format (matches contract)
3. Transform to Prisma models
4. Store in PostgreSQL
5. Expose via GraphQL

### To viz-spa (P6)
1. P6 queries GraphQL endpoint
2. Resolvers fetch from Prisma
3. Aggregate/filter data as needed
4. Return typed response
5. P6 displays in dashboard

---

## Testing Requirements

### Unit Tests
- Resolver logic: 80%+ coverage
- Service functions: 90%+ coverage
- Error cases: All documented errors

### Integration Tests
- GraphQL queries against real database
- Full query pipelines (P4 API → Prisma → GraphQL response)
- Data consistency validation

### Contract Tests
- P4 API contract: Verify cursor-sim responses match contract
- P6 consumption: Test actual P6 queries work

---

## Performance Considerations

### Query Optimization
- **Use Prisma `select`**: Only fields you need
- **Use `include` carefully**: Can cause N+1
- **Batch loaders**: For related data loading
- **Caching**: Use DataLoader for request-scoped caching

### Response Size
- **Pagination required**: No unlimited lists
- **Default page size**: 20-50 items
- **Max page size**: 100 items

---

## Documentation

### SPEC.md Sections
- GraphQL Schema: Complete schema definition
- Queries: All available queries with parameters
- Mutations: All available mutations
- Data Models: Explain relationships

### Commit Messages
```
feat(analytics-core): add {query/mutation}

Implements {description}.

## Changes
- Added Query.{field}
- Schema: Updated {types}

## Testing
- Tests: X new cases
- Coverage: Y%

## P6 Impact
- Breaking changes: Yes/No
- Query updates needed: Yes/No
```

---

## See Also

- API contract in `.claude/skills/api-contract/SKILL.md`
- TypeScript/GraphQL patterns in `.claude/skills/typescript-graphql-patterns/`
- Global coding standards in `03-coding-standards.md`
- SPEC.md: `services/cursor-analytics-core/SPEC.md`
