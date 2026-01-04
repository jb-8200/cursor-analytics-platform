---
name: spec-design
description: Create or revise technical design documents. Use when planning architecture, making component decisions, designing interfaces, or documenting ADRs (Architecture Decision Records). Produces design.md files with rationale and alternatives. (project)
---

# Design Document Standard

Design documents capture **how** we'll implement a feature technically. They include architecture decisions, component designs, data structures, and interfaces.

## Output Location

`.work-items/{feature-name}/design.md`

## Template

```markdown
# Design Document: {Feature Name}

## Overview

{Brief description of the technical approach. 2-3 sentences max.}

## Architecture Decisions

### AD-1: {Decision Title}

**Decision:** {What was decided}

**Rationale:** {Why this approach was chosen}

**Alternatives Considered:**
- {Option A} - {Why rejected}
- {Option B} - {Why rejected}

## Component Design

### Package Structure

```
internal/
├── {package1}/
│   ├── {file1}.go
│   └── {file2}.go
└── {package2}/
```

### Key Interfaces

```go
type {InterfaceName} interface {
    {Method1}({params}) {returns}
    {Method2}({params}) {returns}
}
```

### Data Structures

```go
type {StructName} struct {
    {Field1} {Type} `json:"{field1}"`
    {Field2} {Type} `json:"{field2}"`
}
```

## Data Flow

```
1. Request → Handler
2. Handler → Service
3. Service → Storage
4. Storage → Response
```

## API Design

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /path | Description |
| POST | /path | Description |

## Testing Strategy

### Unit Tests
- {What to test}
- {Key scenarios}

### Integration Tests
- {Component interactions}

## Performance Considerations

| Metric | Target | Measurement |
|--------|--------|-------------|
| Response time | < 100ms | p99 latency |
| Memory | < 100MB | Peak usage |

## Related Documents

- `.work-items/{feature}/user-story.md` - Requirements
- `.work-items/{feature}/task.md` - Implementation tasks
```

## Architecture Decisions (ADRs)

Each significant decision should include:

1. **Decision** - What was decided (clear, specific)
2. **Rationale** - Why this approach (evidence-based)
3. **Alternatives** - What else was considered

Bad:
```
AD-1: Use Parquet
Decision: Use Parquet
```

Good:
```
AD-1: Data Export Format

Decision: Use Parquet with SNAPPY compression

Rationale:
- Columnar format optimized for analytics workloads
- 10x smaller than JSON for tabular data
- Native pandas/R support without parsing

Alternatives Considered:
- JSON: Simpler but 10x larger, slower to parse
- CSV: Universal but no schema, type ambiguity
```

## Interface Design

Define interfaces before implementation:

```go
// Good: Clear contract
type DatasetBuilder interface {
    BuildDataset(from, to time.Time) ([]ResearchRow, error)
    FilterByRepo(rows []ResearchRow, repo string) []ResearchRow
}

// Bad: Implementation leaked into interface
type DatasetBuilder interface {
    BuildDatasetWithSQLQuery(query string) ([]ResearchRow, error)
}
```

## Process

1. **Read user story** for requirements
2. **Identify key decisions** that need documentation
3. **Draft component design** with interfaces first
4. **Document data flow** through the system
5. **Define testing strategy**
6. **Note open questions** for discussion
7. **Create file** in `.work-items/{feature}/design.md`
8. **Proceed to task.md** for implementation breakdown

## Red Flags

| Flag | Fix |
|------|-----|
| No alternatives considered | Think harder, there's always another way |
| Implementation details without "why" | Add rationale |
| Missing error handling design | Add failure scenarios |
| No testing strategy | Add before implementation |
| Vague interfaces | Make specific with types |
