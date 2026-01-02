# Design Document Standard

**Trigger**: Creating technical designs, architecture decisions, component designs, system designs

---

## Purpose

Design documents capture **how** we'll implement a feature technically. They include architecture decisions, component designs, data structures, and interfaces.

## Output Location

`.work-items/{feature-name}/design.md`

---

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

### AD-2: {Decision Title}

{Continue for each significant decision...}

## Component Design

### Package Structure

```
internal/
├── {package1}/
│   ├── {file1}.go
│   └── {file2}.go
└── {package2}/
    └── ...
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

{Describe how data moves through the system}

```
1. Request → Handler
2. Handler → Service
3. Service → Storage
4. Storage → Response
```

{Or use ASCII diagrams:}

```
┌─────────┐    ┌─────────┐    ┌─────────┐
│ Client  │───▶│ Handler │───▶│ Storage │
└─────────┘    └─────────┘    └─────────┘
```

## API Design

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /path | Description |
| POST | /path | Description |

### Request/Response Format

```json
{
  "data": [...],
  "params": {...}
}
```

## Testing Strategy

### Unit Tests
- {What to test}
- {Key scenarios}

### Integration Tests
- {Component interactions}
- {External dependencies}

### E2E Tests
- {Full user flows}

## Performance Considerations

| Metric | Target | Measurement |
|--------|--------|-------------|
| Response time | < 100ms | p99 latency |
| Memory | < 100MB | Peak usage |
| Throughput | 1000 req/s | Sustained load |

## Security Considerations

- {Authentication approach}
- {Authorization rules}
- {Data protection}

## Migration Path

{How to deploy without breaking existing functionality}

## Open Questions

1. **Q:** {Unresolved question}
   **A:** {Current thinking or "TBD"}

## Related Documents

- `.work-items/{feature}/user-story.md` - Requirements
- `.work-items/{feature}/task.md` - Implementation tasks
- `services/{service}/SPEC.md` - Service specification
```

---

## Writing Guidelines

### Architecture Decisions (ADRs)

Each significant decision should include:

1. **Decision** - What was decided (clear, specific)
2. **Rationale** - Why this approach (evidence-based)
3. **Alternatives** - What else was considered (shows thorough analysis)

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
- SNAPPY provides good compression/speed tradeoff

Alternatives Considered:
- JSON: Simpler but 10x larger, slower to parse
- CSV: Universal but no schema, type ambiguity
- Arrow: Overkill for our batch export use case
```

### Interface Design

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

### Data Flow Documentation

Show the happy path clearly:

```
1. HTTP Request: GET /research/dataset?format=parquet
2. Router: Matches /research/* pattern
3. Handler: Parses query params, validates
4. Service: Fetches data from storage
5. Builder: Transforms to ResearchRow structs
6. Writer: Streams Parquet to response
7. HTTP Response: 200 OK with file
```

---

## Example

```markdown
# Design Document: Research Dataset Export

## Overview

Add `/research/dataset` endpoint that exports pre-joined SDLC metrics in
Parquet, JSON, or CSV format for data science workflows.

## Architecture Decisions

### AD-1: Parquet Library

**Decision:** Use `github.com/parquet-go/parquet-go` (v4+)

**Rationale:**
- Actively maintained (weekly commits)
- No CGO dependencies
- Supports schema evolution
- Good documentation

**Alternatives Considered:**
- `xitongsys/parquet-go` - Less active, older API
- Apache Arrow Go - More complex than needed

### AD-2: Streaming Export

**Decision:** Stream directly to HTTP response, no intermediate buffer

**Rationale:**
- Handles large datasets without OOM
- Lower latency (first bytes sent immediately)
- Simpler error handling

**Alternatives Considered:**
- Buffer in memory - OOM risk for large datasets
- Write to temp file - Adds disk I/O, cleanup complexity

## Component Design

### Package Structure

```
internal/
├── export/
│   ├── parquet.go      # Parquet writer
│   ├── csv.go          # CSV writer
│   └── research.go     # Dataset builder
└── api/
    └── research/
        └── handlers.go # HTTP handlers
```

### Key Interfaces

```go
type DatasetBuilder interface {
    BuildDataset(from, to time.Time) []ResearchRow
}

type Exporter interface {
    Export(w io.Writer, rows []ResearchRow) error
}
```

## Testing Strategy

### Unit Tests
- DatasetBuilder with mock storage
- Each Exporter format independently

### Integration Tests
- Full export → parse → verify cycle
- Large dataset performance

### E2E Tests
- HTTP request → file download → pandas load
```

---

## Process

1. **Read user story** for requirements
2. **Identify key decisions** that need documentation
3. **Draft component design** with interfaces first
4. **Document data flow** through the system
5. **Define testing strategy**
6. **Note open questions** for discussion
7. **Create file** in `.work-items/{feature}/design.md`
8. **Proceed to task.md** for implementation breakdown

---

## Red Flags

| Flag | Fix |
|------|-----|
| No alternatives considered | Think harder, there's always another way |
| Implementation details without "why" | Add rationale |
| Missing error handling design | Add failure scenarios |
| No testing strategy | Add before implementation |
| Vague interfaces | Make specific with types |
