# Model Selection Guide

This skill provides guidance on which Claude model to use for different development tasks in the Cursor Analytics Platform.

## Model Capabilities Overview

| Model | Cost | Speed | Best For |
|-------|------|-------|----------|
| **Haiku** | $ | ⚡⚡⚡ | Well-specified implementation, pattern following, simple refactoring |
| **Sonnet** | $$ | ⚡⚡ | Complex logic, architectural decisions, multi-file changes, debugging |
| **Opus** | $$$ | ⚡ | Novel architectures, major refactoring, cross-service integration |

---

## Task-to-Model Mapping

### P0: Make Runnable (Scaffolding)

| Task | Recommended Model | Rationale |
|------|-------------------|-----------|
| P0.1: Go scaffolding | **Sonnet** | Establish patterns for first time |
| P0.2: TypeScript scaffolding | **Sonnet** | Establish patterns for first time |
| P0.3: React scaffolding | **Sonnet** | Vite config requires decisions |
| P0.4: .env.example | **Haiku** | Simple file creation from SPEC |
| P0.5: Update docker-compose.yml | **Haiku** | Following existing structure |
| P0.6: Update Makefile | **Haiku** | Pattern-based additions |
| P0.7: Verify docker-compose up | **Sonnet** | Debugging integration issues |
| P0.8: Happy path implementation | **Sonnet** | Cross-service integration |

---

### cursor-sim Tasks

| Task ID | Task Name | Model | Rationale |
|---------|-----------|-------|-----------|
| TASK-SIM-001 | Initialize Go Project | **Sonnet** | Project structure decisions |
| TASK-SIM-002 | CLI Flag Parsing | **Haiku** | Standard flag/JSON parsing |
| TASK-SIM-003 | Developer Profile Generator | **Haiku** | Struct implementation from SPEC.md:145-250 |
| TASK-SIM-004 | Event Generation Engine | **Sonnet** | Poisson distribution, complex timing logic |
| TASK-SIM-005 | In-Memory Storage | **Haiku** | Standard sync.Map patterns |
| TASK-SIM-006 | REST API Handlers | **Haiku** | Following cursor-api-patterns.md |
| TASK-SIM-007 | Wire Up Main | **Sonnet** | Integration and error handling |

---

### cursor-analytics-core Tasks

| Task ID | Task Name | Model | Rationale |
|---------|-----------|-------|-----------|
| TASK-CORE-001 | Initialize TS Project | **Sonnet** | Apollo setup, project decisions |
| TASK-CORE-002 | Database Schema | **Sonnet** | Schema design, migration strategy |
| TASK-CORE-003 | GraphQL Schema | **Haiku** | Implementing from specs/api/graphql-schema.graphql |
| TASK-CORE-004 | Data Ingestion Worker | **Sonnet** | Polling logic, error handling, retries |
| TASK-CORE-005 | Metric Calculation | **Haiku** | SQL queries from spec |
| TASK-CORE-006 | Developer Resolvers | **Haiku** | GraphQL resolver patterns |

---

### cursor-viz-spa Tasks

| Task ID | Task Name | Model | Rationale |
|---------|-----------|-------|-----------|
| TASK-VIZ-001 | Initialize React Project | **Sonnet** | Vite + TanStack Query setup |
| TASK-VIZ-002 | Dashboard Layout | **Haiku** | Component structure from SPEC |
| TASK-VIZ-003 | GraphQL Client Setup | **Haiku** | Standard Apollo Client config |
| TASK-VIZ-004 | Velocity Heatmap | **Sonnet** | Complex Recharts visualization |
| TASK-VIZ-005 | Developer Efficiency Table | **Haiku** | Table component from spec |

---

## Activity-Based Recommendations

### Use Haiku (⚡ Fast & Cheap) for:

**Implementation from Specs**
- Implementing structs/types defined in SPEC.md
- Creating HTTP handlers following cursor-api-patterns.md
- Writing Go code following go-best-practices.md
- Building React components with clear requirements

**Writing Tests**
- Table-driven tests from test-plan.md
- Unit tests for well-specified functions
- Integration tests following patterns
- Test fixtures and mocks

**Pattern Following**
- Adding pagination to endpoints (cursor-api-patterns.md)
- Implementing Basic Auth (cursor-api-patterns.md:76-94)
- Error handling patterns (go-best-practices.md:69-133)
- JSON marshaling (go-best-practices.md:366-402)

**Simple Refactoring**
- Extract helper functions
- Rename variables/functions
- Format code (gofmt, prettier)
- Fix linter warnings

**Documentation**
- Adding code comments
- Updating CHANGELOG.md
- Writing inline documentation
- Creating examples

---

### Use Sonnet (⚡⚡ Balanced) for:

**Complex Logic**
- Poisson distribution event generation
- Rate limiting algorithms
- Caching strategies
- Concurrent worker pools

**Architectural Decisions**
- Service integration patterns
- Database schema design
- API versioning strategy
- Error recovery mechanisms

**Multi-File Changes**
- Feature implementation spanning multiple files
- Refactoring across packages
- Adding middleware layers
- Updating contracts across services

**Debugging**
- Race condition investigation
- Memory leak analysis
- Performance bottleneck identification
- Integration test failures

**Planning & Design**
- /start-feature command (creates design.md)
- Writing technical design documents
- Planning migration strategies
- Evaluating trade-offs

---

### Use Opus (⚡ Powerful) for:

**Novel Architectures**
- Designing new major features not in specs
- Evaluating alternative approaches
- Complex refactoring across all services
- Performance optimization strategy

**Cross-Service Integration**
- End-to-end feature implementation
- Contract negotiation between services
- Breaking API changes
- Data flow redesign

**Major Debugging**
- Production incident investigation
- Distributed tracing analysis
- Complex race conditions
- Subtle correctness bugs

**Strategic Planning**
- Choosing technology stack
- Migration planning (e.g., database change)
- Scalability architecture
- Security architecture review

---

## Decision Tree

```
Start with question: "Is this well-specified in SPEC.md?"
├─ YES → "Does it follow an existing pattern?"
│   ├─ YES → Use Haiku ⚡
│   └─ NO → "Is it complex logic?"
│       ├─ YES → Use Sonnet ⚡⚡
│       └─ NO → Use Haiku ⚡
│
└─ NO → "Is it a major architectural decision?"
    ├─ YES → Use Opus ⚡
    └─ NO → Use Sonnet ⚡⚡
```

---

## Cost Optimization Strategy

### Session-Based Approach

**Planning Session (Sonnet/Opus)**
1. Run /start-feature to create work item
2. Make design decisions in design.md
3. Write test-plan.md
4. Create task.md checklist

**Implementation Session (Haiku)**
5. Write tests from test-plan.md
6. Implement code from SPEC.md
7. Follow patterns from skills
8. Check off tasks in task.md

**Integration Session (Sonnet)**
9. Wire up components
10. Debug integration issues
11. Verify end-to-end flow

### Batch Similar Tasks

**Haiku batch** (single session):
- Implement 5 similar HTTP handlers
- Write 10 unit tests
- Create 3 React components

**Sonnet batch** (single session):
- Design 2 related features
- Debug 3 integration issues
- Refactor 4 interconnected modules

---

## Slash Command Model Hints

When you run slash commands, they will suggest appropriate models:

```bash
/next-task cursor-sim
# Output includes: "Recommended: haiku ⚡ (well-specified struct)"

/start-feature complex-caching
# Output includes: "Recommended: sonnet ⚡⚡ (architectural decisions needed)"

/verify cursor-sim
# Output includes: "Recommended: sonnet ⚡⚡ (cross-file analysis)"
```

---

## Example Usage

### Efficient Workflow

```bash
# Use Sonnet to plan
User: "Use Sonnet to run /start-feature developer-generator"
# Creates design.md, task.md, test-plan.md

# Switch to Haiku for implementation
User: "Use Haiku to implement the Developer struct from design.md"
# Fast, cheap implementation

# Back to Sonnet for integration
User: "Use Sonnet to wire up the generator in main.go"
# Handles complex integration

# Haiku for tests
User: "Use Haiku to write tests from test-plan.md"
# Pattern-based test writing
```

### Cost Comparison

**All Sonnet Approach:**
- Planning: $0.50
- Implementation: $2.00
- Integration: $0.75
- Testing: $1.00
- **Total: $4.25**

**Hybrid Approach:**
- Planning (Sonnet): $0.50
- Implementation (Haiku): $0.20
- Integration (Sonnet): $0.75
- Testing (Haiku): $0.10
- **Total: $1.55** (64% savings)

---

## When in Doubt

**Default to Sonnet** if:
- You're unsure of the complexity
- It's the first time implementing this pattern
- The task involves multiple services
- You're debugging an unexpected issue

**Switch to Haiku** once:
- The pattern is established
- You have clear acceptance criteria
- The task is well-isolated
- You're following documented examples

---

## Integration with Commands

Enhanced slash commands support model parameters:

```bash
/implement TASK-SIM-003 --model=haiku
/implement TASK-SIM-004 --model=sonnet
```

Commands automatically recommend models based on this guide.

---

**Remember**: The comprehensive specs in this project make it ideal for Haiku usage. Use Sonnet/Opus for planning and complex logic, then leverage Haiku for the well-specified implementation work.
