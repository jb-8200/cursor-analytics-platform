---
description: Shared coding standards for Go, TypeScript, and React
---

# Coding Standards

Language-specific patterns and conventions applied across services.

---

## Go (cursor-sim, P4)

### Formatting
- **ALWAYS** use `gofmt` (or `goimports` for imports)
- No custom formatting, trust standard tools
- Line length: convention is ~80 chars but gofmt decides

### Error Handling
- **NEVER** panic in production code (only init/main for setup)
- **ALWAYS** return errors explicitly: `return err` not `panic(err)`
- **ALWAYS** wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Check errors immediately: `if err != nil { return err }`
- Log errors with context, not just "failed"

### Testing
- **Test files**: `*_test.go` in same package
- **Table-driven tests** for multiple cases
- **100% coverage target** for handler code
- **Use `t.Run()`** for subtest organization
- **Mock external calls**, don't call real APIs

### Packages and Imports
- **Single responsibility**: One concept per package
- **Avoid circular imports**: Use interfaces to break cycles
- **Standard library first**: Only add dependencies if essential
- **Version locked**: Use go.mod, commit go.sum

### Naming
- **Exported functions**: PascalCase (GetUser, HandleRequest)
- **Unexported functions**: camelCase (getData, processEvent)
- **Interfaces**: Reader, Writer (convention: -er suffix)
- **Receiver names**: 1-2 chars or abbreviation (r for Receiver, h for Handler)

---

## TypeScript (analytics-core & viz-spa, P5 & P6)

### Strict Mode
- **ALWAYS** use `strict: true` in tsconfig.json
- **ALWAYS** specify return types: `function foo(): string {}`
- **ALWAYS** use explicit types: `const x: number = 5` (unless obvious)
- **NO `any` type**: Use `unknown` with type guards instead

### Imports and Exports
- **Named exports**: `export function foo() {}`
- **Default exports**: Only for large modules
- **Absolute imports**: Use tsconfig.json paths if available
- **Organize imports**: Standard library → packages → local

### Error Handling
- **Throw descriptive errors**: `throw new Error("context: what failed")`
- **Use custom error classes** for domain errors
- **Catch specific errors**, not bare `catch (e)`
- **Log with context**: Include stack trace and relevant data

### Testing
- **Test files**: `*.test.ts` or `*.spec.ts`
- **Unit tests**: Single responsibility per test
- **Use describe/it**: Organize into logical groups
- **Mock dependencies**: Don't call real APIs
- **80%+ coverage target**: Focus on logic, not lines

### React-Specific (viz-spa)

#### Components
- **Functional components only**: No class components
- **Props interface**: Define clear prop types
- **Hooks only**: useState, useEffect, custom hooks
- **Extract logic**: Separate hooks from components
- **Memoization**: Use memo/useMemo only when needed

#### Styling
- **Tailwind CSS only**: No custom CSS files
- **Utility-first**: Stack classes, don't create custom ones
- **Responsive**: Mobile-first, use breakpoints (sm, md, lg)
- **Dark mode**: Support through Tailwind dark: variant

#### Accessibility
- **WCAG 2.1 AA**: Minimum standard
- **Semantic HTML**: Use proper elements
- **ARIA labels**: For complex components
- **Keyboard navigation**: Test with keyboard only
- **Color contrast**: 4.5:1 ratio for text

---

## GraphQL (analytics-core, P5)

### Schema Design
- **Null safety**: Use `!` for required fields intentionally
- **Pagination**: Cursor-based (not offset) for consistency
- **Mutations**: Return fields for success confirmation
- **Error handling**: Use response types with error field
- **Documentation**: Add descriptions to fields

### Resolvers
- **Type safety**: Use generated types
- **Error handling**: Throw with message, not silent failures
- **Caching**: Use appropriate fetch policies
- **N+1 prevention**: Use batch loaders
- **Input validation**: Validate at resolver boundary

---

## General Principles

### Code Quality
- **DRY**: Don't Repeat Yourself (extract helpers)
- **KISS**: Keep It Simple, Stupid (avoid over-engineering)
- **Composition**: Small functions, composed together
- **Single Responsibility**: One job per function

### Testing Strategy
- **Test behavior, not implementation**: Focus on what, not how
- **Arrange-Act-Assert**: Clear test structure
- **Descriptive names**: Test name explains what it tests
- **Fast tests**: < 100ms each ideally
- **Isolated tests**: No dependencies between tests

### Documentation
- **Comments explain WHY**: Not what (code shows what)
- **README per service**: How to run, test, deploy
- **SPEC.md**: Keep technical spec current
- **Docstrings**: For public APIs only

---

## See Also

- SDD process in `04-sdd-process.md`
- Service-specific rules (cursor-sim.md, analytics-core.md, viz-spa.md)
- Dependency-reflection skill for impact analysis
