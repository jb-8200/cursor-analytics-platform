# Available Skills

Skills are auto-triggered knowledge modules that activate when your request matches their description. They provide guidance, patterns, and how-to documentation—not enforcement.

**Difference from Rules**:
- **Rules** (in `.claude/rules/`) are always-on enforcement with NEVER/ALWAYS constraints
- **Skills** (in `.claude/skills/`) provide optional guidance that activates on-demand

---

## Skills Catalog

| Skill | Triggers | Service | Purpose |
|-------|----------|---------|---------|
| `api-contract` | "API", "endpoint", "contract", "request", "response" | P4 (cursor-sim) | Reference for cursor-sim REST API endpoints, models, and response formats |
| `cursor-api-patterns` | "API implementation", "handler", "authentication", "pagination" | P4 (cursor-sim) | Patterns for implementing HTTP handlers, auth, pagination, and CSV exports |
| `go-best-practices` | Go code, "handlers", "structs", "error", "testing" | P4 (cursor-sim) | Go naming conventions, error handling, concurrency, testing patterns |
| `typescript-graphql-patterns` | GraphQL, "resolver", "schema", "mutation", "query" | P5 (analytics-core) | Apollo Server, resolver implementations, type safety, error handling |
| `react-vite-patterns` | React, Vite, "hooks", "component", "Apollo", "Tailwind" | P6 (viz-spa) | React hooks, Apollo Client integration, Tailwind styling patterns |
| `sdd-checklist` | "task complete", "commit", "ready to commit", "next task" | All | Post-task completion checklist: REFLECT, SYNC, COMMIT workflow |
| `spec-process-core` | SDD, "Spec-Driven Development", "workflow" | All | Core principles of Spec-Driven Development methodology |
| `spec-process-dev` | TDD, "test-driven", "red-green-refactor" | All | Test-Driven Development workflow (RED → GREEN → REFACTOR) |
| `dependency-reflection` | Refactoring, "model change", "interface", "endpoint added" | All | Identifies required updates to dependent files after code changes |
| `spec-sync-check` | Phase complete, "endpoint added", "mutation", "schema change" | All | Determines if SPEC.md needs updating after implementation changes |
| `spec-user-story` | User story, "EARS format", "acceptance criteria" | All | Creates/revises user stories following EARS format (As-a/I-want/So-that) |
| `spec-design` | Design document, "architecture", "ADR", "decision" | All | Creates/revises technical design documents with rationale and alternatives |
| `spec-tasks` | Task breakdown, "decompose", "work units" | All | Creates/revises task breakdowns for feature implementation |
| `model-selection-guide` | Model choice, "which model", "cost", "speed" | All | Helps select the right Claude model (Haiku vs Sonnet vs Opus) |

---

## How Skills Activate

Skills are automatically invoked when you make a request that matches their description. Examples:

**Auto-triggers**:
- Asking about GraphQL → `typescript-graphql-patterns` activates
- Completing a task → `sdd-checklist` activates
- Asking about React components → `react-vite-patterns` activates

**Explicit invocation**:
- Use the `Skill` tool with skill name (e.g., `Skill(skill: "sdd-checklist")`)
- Or use shorthand `/commit` → invokes the commit skill

---

## Rules vs Skills vs Commands

| Type | Trigger | Format | Scope | Example |
|------|---------|--------|-------|---------|
| **Rule** | Always loaded | YAML + NEVER/ALWAYS | Path-scoped or global | `cursor-sim.md`: "NEVER modify API" |
| **Skill** | Matching request | Markdown + frontmatter | Knowledge/patterns | `react-vite-patterns`: "Hook best practices" |
| **Command** | Explicit user call | `/command` | Workflow | `/commit`: Execute commit skill |

**Decision tree**:
- **MUST always enforce?** → Create Rule
- **Provide helpful guidance?** → Create Skill
- **Multi-step workflow?** → Create Command (invokes skills)

---

## Service Mapping

### P4: cursor-sim (Go CLI + REST API)
- `api-contract` - API reference
- `cursor-api-patterns` - API implementation
- `go-best-practices` - Go patterns

### P5: cursor-analytics-core (TypeScript + GraphQL)
- `typescript-graphql-patterns` - GraphQL schema, resolvers
- `api-contract` - Consuming P4 API

### P6: cursor-viz-spa (React + Vite + Apollo)
- `react-vite-patterns` - React components, hooks
- `typescript-graphql-patterns` - Consuming P5 GraphQL

### Cross-Service
- `sdd-checklist` - Post-task workflow
- `spec-*` - Specification and documentation
- `dependency-reflection` - Checking cross-file impacts
- `spec-sync-check` - SPEC.md updates
- `model-selection-guide` - Model optimization

---

## Using Skills in Your Work

1. **Request matches skill description** → Skill auto-activates with relevant guidance
2. **Need specific guidance** → Use `/skillname` to explicitly invoke
3. **Combined workflows** → Commands may invoke multiple skills sequentially

---

## Skill Development

When creating new skills:

1. **Create skill directory**: `.claude/skills/{skill-name}/`
2. **Add SKILL.md with frontmatter**:
   ```markdown
   ---
   name: skill-name
   description: Clear description with trigger keywords
   allowed-tools: [optional, only if restricting tools]
   ---

   # Skill Name

   [Guidance content]
   ```
3. **Include**: Examples, templates, decision trees
4. **Avoid**: Enforcement language (NEVER/ALWAYS) — use Rules instead

---

## See Also

- Rules enforcement: `.claude/rules/README.md`
- SDD process: `.claude/rules/04-sdd-process.md`
- Full Claude Code guide: `.claude/README.md`
