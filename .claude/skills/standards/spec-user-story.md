# User Story Standard

**Trigger**: Creating user stories, PRDs, feature narratives, product requirements, UX specifications

---

## Purpose

User stories define **what** we're building and **why**, focused on end-user value. They do NOT include technical implementation details.

## Output Location

`.work-items/{feature-name}/user-story.md`

---

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

{Continue for each criterion...}

## Out of Scope

- {Explicit exclusion 1}
- {Explicit exclusion 2}
- {Things that might be assumed but are NOT included}

## Dependencies

- {Required prerequisite 1}
- {Required prerequisite 2}
- {External systems or features needed}

## Success Metrics

- {How we'll measure success}
- {Quantifiable outcomes}
- {User satisfaction indicators}

## Technical Notes

{Optional section for constraints that affect the story:
- Performance requirements
- Security considerations
- Compatibility requirements}

## Related Documents

- `services/{service}/SPEC.md` - Technical specification
- `.work-items/{feature}/design.md` - Technical design
- `docs/USER_STORIES.md` - All user stories
```

---

## Writing Guidelines

### The "As a... I want... So that..." Format

**As a** - The persona benefiting from this feature
- Be specific: "data scientist" not "user"
- Include context: "data scientist researching AI adoption"

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

---

## Example

```markdown
# User Story: Research Dataset Export

## Story

**As a** data scientist researching AI coding assistants
**I want** to export cursor-sim data in analysis-ready formats
**So that** I can conduct reproducible SDLC research without regenerating data

## Background

Cursor-sim generates synthetic developer activity data. Currently, researchers
must query the API repeatedly and manually join datasets. This story adds
direct export of pre-joined, analysis-ready datasets.

## Acceptance Criteria

### AC-1: Parquet Export

**Given** cursor-sim has generated commits and PRs
**When** I request GET /research/dataset?format=parquet
**Then** I receive a valid Parquet file with pre-joined metrics

**Required columns**: pr_number, author_email, ai_ratio, survival_rate_30d...
**Performance**: Export of 100k PRs completes in < 10 seconds

### AC-2: Multiple Format Support

**Given** the research dataset is available
**When** I request with format=json or format=csv
**Then** I receive the same data in the requested format

### AC-3: Date Filtering

**Given** data spanning January 2026
**When** I request ?startDate=2026-01-15&endDate=2026-01-20
**Then** only PRs created in that range are included

## Out of Scope

- Real-time streaming exports
- Custom column selection
- Data transformation/aggregation in export
- Multi-tenant data isolation

## Dependencies

- Phase 1: Cursor API simulation (complete)
- Phase 2: GitHub PR simulation (complete)
- parquet-go library

## Success Metrics

- Researchers can load export directly into pandas/R
- No manual data joining required
- Export performance meets targets
```

---

## Process

1. **Gather requirements** from stakeholder/user
2. **Draft story** using template
3. **Review acceptance criteria** - each must be testable
4. **Define scope boundaries** explicitly
5. **Identify dependencies**
6. **Create file** in `.work-items/{feature}/user-story.md`
7. **Proceed to design.md** for technical approach

---

## Red Flags

| Flag | Fix |
|------|-----|
| Story includes implementation details | Move to design.md |
| AC not testable | Add specific, measurable criteria |
| No "out of scope" section | Add explicit exclusions |
| Vague success metrics | Add quantifiable measures |
| Multiple unrelated features | Split into separate stories |
