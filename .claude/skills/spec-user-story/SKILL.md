---
name: spec-user-story
description: Create or revise user stories following EARS format. Use when writing PRDs, feature narratives, acceptance criteria, or product requirements. Produces user-story.md files with As-a/I-want/So-that and Given-When-Then criteria.
---

# User Story Standard

User stories define **what** we're building and **why**, focused on end-user value. They do NOT include technical implementation details.

## Output Location

`.work-items/{feature-name}/user-story.md`

## Template

```markdown
# User Story: {Feature Name}

## Story

**As a** {role/persona}
**I want** {capability/feature}
**So that** {benefit/value}

## Background

{Context explaining why this feature matters. Include:
- Current pain points
- How this fits into the larger system
- Any relevant history or constraints}

## Acceptance Criteria

### AC-1: {Criterion Name}

**Given** {precondition/context}
**When** {action/trigger}
**Then** {expected outcome}

{Additional details:
- Specific values or thresholds
- Edge cases to consider
- Performance requirements}

### AC-2: {Criterion Name}

**Given** {precondition}
**When** {action}
**Then** {outcome}

## Out of Scope

- {Explicit exclusion 1}
- {Explicit exclusion 2}

## Dependencies

- {Required prerequisite 1}
- {Required prerequisite 2}

## Success Metrics

- {How we'll measure success}
- {Quantifiable outcomes}

## Related Documents

- `services/{service}/SPEC.md` - Technical specification
- `.work-items/{feature}/design.md` - Technical design
```

## Writing Guidelines

### The "As a... I want... So that..." Format

**As a** - The persona benefiting from this feature
- Be specific: "data scientist" not "user"

**I want** - The capability being requested
- Action-oriented: "export data to Parquet"
- One clear thing, not a list

**So that** - The value delivered
- Business benefit: "I can conduct reproducible research"
- Not technical: "the system can generate files"

### Acceptance Criteria Best Practices

1. **Each AC is independently testable**
2. **Use Given-When-Then format** (maps to Arrange-Act-Assert)
3. **Include specific values** when they matter
4. **Cover edge cases** explicitly
5. **Define "done"** clearly

Bad:
```
AC-1: The export works correctly
```

Good:
```
AC-1: Export to Parquet Format

Given 1000 commits in the system
When I request GET /research/dataset?format=parquet
Then I receive a valid Parquet file
And the file contains all 22 required columns
And file size is < 10MB for this dataset
```

### Out of Scope is Critical

Explicitly state what's NOT included:
- Prevents scope creep
- Sets clear boundaries
- Avoids assumptions

## Process

1. **Gather requirements** from stakeholder/user
2. **Draft story** using template
3. **Review acceptance criteria** - each must be testable
4. **Define scope boundaries** explicitly
5. **Identify dependencies**
6. **Create file** in `.work-items/{feature}/user-story.md`
7. **Proceed to design.md** for technical approach

## Red Flags

| Flag | Fix |
|------|-----|
| Story includes implementation details | Move to design.md |
| AC not testable | Add specific, measurable criteria |
| No "out of scope" section | Add explicit exclusions |
| Vague success metrics | Add quantifiable measures |
| Multiple unrelated features | Split into separate stories |
