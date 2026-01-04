---
description: cursor-sim specific guardrails protecting API contracts and CLI isolation
paths: services/cursor-sim/**
---

# cursor-sim Rules

Service-specific constraints for P4 (cursor-sim foundation).

---

## NEVER

- **Modify `internal/api/`** without updating SPEC.md and E2E tests
- **Change API response format** without downstream testing (P5/P6)
- **Modify `internal/generator/`** from CLI subagent tasks
- **Break backward compatibility** in API responses
- **Add public API endpoints** that aren't documented in SPEC.md
- **Touch `internal/api/`** if assigned to CLI-only work
- **Expose internal data structures** in API responses

---

## ALWAYS

- **Update SPEC.md** whenever endpoints or response schemas change
- **Write E2E tests** for all endpoint modifications
- **Keep response schemas backward compatible** for P5/P6
- **Document all endpoint changes** in commit messages
- **Include test coverage** for API handlers (80%+ target)
- **Use SPEC.md as source of truth** for contract verification
- **Test with actual P5 consumer** before considering API change complete

---

## API Contract Protection

### What P5 (analytics-core) Depends On
- REST endpoints in `/teams/*`, `/admin/*`, `/analytics/*`
- Response format: JSON with `data`, `pagination`, `params` fields
- Authentication: Basic Auth with API key
- Pagination: Offset-based with `page`, `page_size`, `total`
- Error format: Status code + error message

### What P6 (viz-spa) Depends On
- P5 GraphQL schema (not direct P4 dependency)
- P5 must transform P4 data correctly

### Breaking Changes
**NEVER make these changes**:
- Remove endpoints without deprecation notice
- Change response field names
- Change field types (string → int)
- Alter pagination behavior
- Change authentication method

---

## CLI Isolation (cursor-sim-cli-dev subagent)

### Allowed Scope
- `services/cursor-sim/internal/cli/` - CLI implementation
- `services/cursor-sim/cmd/simulator/` - Main entry point
- `.work-items/P4-*` - Task documentation

### Prohibited Scope (NEVER TOUCH)
- `internal/api/` - API handlers
- `internal/generator/` - Data generation
- `internal/models/` - Core data models
- `internal/service/` - Business logic

**Reason**: Protects API contracts for downstream P5/P6 services

---

## Testing Requirements

### Unit Tests
- `*_test.go` files for all new code
- Table-driven tests for multiple cases
- 80%+ coverage for handler code

### E2E Tests
- Test endpoints end-to-end before deploying
- Use actual request/response validation
- Cover error cases and edge conditions

### Integration Tests
- Test with Prisma database if applicable
- Test with cursor-sim API client (P5 perspective)

---

## Documentation

### SPEC.md Sections to Update
- Endpoints table: Add/modify endpoint
- Response Schemas: Document response format
- Authentication: If changed
- Error Responses: Document error codes
- Data Models: If new types added

### Commit Message Template
```
feat(cursor-sim): add {endpoint} endpoint

Implements {description} for {feature}.

## Changes
- Added POST /endpoint
- Updated response schema

## Testing
- E2E tests: X new cases
- Coverage: Y%

## API Impact
- Backward compatible: Yes/No
- Requires P5 update: Yes/No
- Requires P6 update: Yes/No
```

---

## See Also

- API contract in `.claude/skills/api-contract/SKILL.md`
- Global coding standards in `03-coding-standards.md`
- SPEC.md: `services/cursor-sim/SPEC.md`
