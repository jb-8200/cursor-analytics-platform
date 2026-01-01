# Verify Spec-Test-Code Alignment

When the user runs `/verify [service-name]`, check that specifications, tests, and implementation are aligned for the given service.

## Verification Checklist

### 1. Spec Completeness
- [ ] SPEC.md exists for service
- [ ] All API endpoints documented
- [ ] All data models defined
- [ ] Configuration options listed
- [ ] Error scenarios covered

### 2. Test Coverage
- [ ] Unit tests exist for all components
- [ ] Integration tests exist for all endpoints
- [ ] Test files follow naming convention (*.test.*)
- [ ] Run coverage report and verify >= 80%

### 3. Implementation Alignment
- [ ] All endpoints in SPEC.md have handlers implemented
- [ ] All data models in SPEC.md have types/structs defined
- [ ] All configuration in SPEC.md is loaded
- [ ] Error handling matches SPEC.md error scenarios

### 4. Documentation Sync
- [ ] CHANGELOG.md updated for user-facing changes
- [ ] API documentation matches implementation
- [ ] README examples work with current code

## Verification Commands

### For Go Services (cursor-sim)
```bash
# Check test coverage
cd services/cursor-sim
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total

# Run linter
golangci-lint run
```

### For TypeScript Services (cursor-analytics-core, cursor-viz-spa)
```bash
# Check test coverage
npm run test:coverage

# Run linter
npm run lint

# Type check
npm run type-check
```

## Report Format

Display results as:

```
Verification Report: {service-name}
=====================================

✓ Spec Completeness:     PASS
✓ Test Coverage:         85.3% (>= 80% required)
✗ Implementation:        FAIL - Missing endpoints:
    - GET /v1/analytics/team/models
    - GET /v1/analytics/team/dau
✓ Documentation:         PASS

Summary: 3/4 checks passed

Action Items:
1. Implement missing endpoints (see SPEC.md lines 245-289)
2. Add integration tests for new endpoints
```

## Example Usage

```
User: /verify cursor-sim
Assistant: [Runs verification checks and displays report]
```

## Implementation Instructions

1. Read the service's SPEC.md to know what should exist
2. Check the actual codebase for implementation
3. Run test commands to verify coverage
4. Report gaps between spec and implementation
5. Provide actionable next steps to fix misalignment
