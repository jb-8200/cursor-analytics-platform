# User Story: Preview Mode & YAML Seed Support

**Feature ID**: P3-F04-preview-yaml
**Created**: January 8, 2026
**Status**: Planning
**Inspired By**: NVIDIA NeMo DataDesigner patterns

---

## Story (EARS Format)

### As-a / I-want / So-that

**As a** cursor-sim user (researcher, developer, QA engineer)
**I want** to preview generated data before full generation and use YAML seed files
**So that** I can validate configurations quickly and use more readable seed formats

---

## Context

### Inspiration from DataDesigner

NVIDIA's DataDesigner library demonstrates best practices for synthetic data generation:
- **Preview Mode**: Test configurations with small samples before full-scale generation
- **YAML Configuration**: Human-readable, comment-friendly format for complex configs
- **Validation Focus**: Quick feedback loops for data quality

### Current Limitations

1. **No Preview Capability**
   - Must generate full dataset to see output quality
   - Slow iteration: 90 days × 10 developers = ~500+ commits before validation
   - Wastes time on misconfigured seeds

2. **JSON-Only Seeds**
   - No comments allowed (JSON limitation)
   - Hard to read for complex configurations
   - Less maintainable for large seed files

3. **Limited Validation Feedback**
   - Only shows "Generated X commits" at end
   - No quick sample inspection
   - Can't test seed changes rapidly

---

## Requirements

### Functional Requirements (EARS)

#### FR1: Preview Mode Flag
**WHEN** user runs `cursor-sim -mode preview -seed seed.json`
**THEN** system generates sample data (10% or 7 days max)
**AND** displays formatted output to console
**AND** exits without starting server

#### FR2: Preview Output Format
**WHEN** preview mode completes
**THEN** system displays:
- Developer summary (count, IDs, working hours)
- Sample commits (first 5 per developer)
- Sample events (model usage, PRs, reviews)
- Event distribution statistics
- Validation warnings (if any)

#### FR3: YAML Seed File Support
**WHEN** user provides `-seed config.yaml`
**THEN** system detects YAML extension
**AND** parses using YAML parser
**AND** converts to internal SeedData structure

#### FR4: Backward Compatibility
**WHEN** user provides `-seed config.json`
**THEN** system uses existing JSON parser
**AND** behaves identically to current implementation

#### FR5: Validation Warnings
**WHEN** seed file has potential issues
**THEN** preview mode displays warnings:
- Working hours overlap issues
- Invalid model names
- Missing required fields
- Range violations (e.g., velocity not in [low, medium, high])

---

### Non-Functional Requirements

#### NFR1: Performance
- Preview mode completes in < 5 seconds for typical seeds
- YAML parsing overhead < 50ms

#### NFR2: Usability
- Preview output is human-readable and well-formatted
- YAML syntax errors show clear line numbers and context
- Preview clearly indicates it's a sample, not full dataset

#### NFR3: Maintainability
- YAML and JSON share same internal data structure
- No code duplication between parsers
- Easy to add CSV or other formats later

---

## Acceptance Criteria (Given-When-Then)

### AC1: Preview Mode Basic Functionality
**GIVEN** a valid JSON seed file with 2 developers
**WHEN** I run `./cursor-sim -mode preview -seed testdata/valid_seed.json`
**THEN** I see:
- Developer summary (2 developers listed)
- 5 sample commits per developer (10 total)
- Sample model usage events (3-5 events)
- Statistics: "Preview generated 10 commits, 5 model events across 2 developers"
- Message: "This is a preview. Use -mode runtime for full generation."
- Exit code 0

### AC2: Preview Mode with YAML
**GIVEN** a YAML seed file with 3 developers
**WHEN** I run `./cursor-sim -mode preview -seed testdata/config.yaml`
**THEN** I see identical output format as JSON
**AND** system parses YAML correctly
**AND** exit code 0

### AC3: YAML Comments Supported
**GIVEN** a YAML seed with inline comments
```yaml
developers:
  - user_id: alice  # Engineering lead
    email: alice@example.com
    working_hours:
      start: 9  # PST timezone
      end: 17
```
**WHEN** I load this seed in preview or runtime mode
**THEN** comments are ignored (YAML spec)
**AND** parsing succeeds
**AND** data loads correctly

### AC4: Preview Validation Warnings
**GIVEN** a seed file with invalid model name "gpt-5000"
**WHEN** I run preview mode
**THEN** I see warning:
```
⚠️  Validation Warnings:
  - Developer alice: Unknown model "gpt-5000" in preferred_models
    Valid models: claude-sonnet-4.5, gpt-4o, claude-opus-4
```
**AND** preview continues with other data
**AND** exit code 0 (warnings, not errors)

### AC5: Runtime Mode Unaffected
**GIVEN** I've used preview mode successfully
**WHEN** I run `./cursor-sim -mode runtime -seed seed.yaml -port 8080`
**THEN** server starts normally
**AND** generates full dataset
**AND** serves 29 API endpoints
**AND** behavior identical to JSON seed

### AC6: Invalid Mode Handling
**GIVEN** invalid mode specified
**WHEN** I run `./cursor-sim -mode invalid`
**THEN** I see error: "Invalid mode: 'invalid'. Valid modes: runtime, preview"
**AND** exit code 1

---

## Out of Scope

- **CSV seed files**: Future enhancement
- **LLM-as-judge validation**: Too heavy for preview
- **Interactive YAML editing**: User uses their own editor
- **Web UI for preview**: CLI-only for now
- **Saving preview output**: User pipes to file if needed

---

## Success Metrics

- Preview mode completes in < 5 seconds
- 90% of seed validation issues caught in preview
- 50% reduction in "generate → see problem → fix → regenerate" cycles
- YAML adoption rate: >30% of new seed files use YAML

---

## User Scenarios

### Scenario 1: Researcher Testing New Seed
```bash
# Write YAML seed with comments
vim research_team.yaml

# Quick preview to validate
./cursor-sim -mode preview -seed research_team.yaml

# Output shows:
# ✅ 5 developers loaded
# ✅ Working hours validated
# ⚠️  Developer "bob": Model "claude-haiku-3.5" not found
#
# Sample commits: [10 commits shown]
# Preview complete.

# Fix model name
vim research_team.yaml

# Preview again (5 seconds)
./cursor-sim -mode preview -seed research_team.yaml

# ✅ All validation passed
# Ready for full generation
./cursor-sim -mode runtime -seed research_team.yaml -port 8080 -days 180
```

### Scenario 2: QA Engineer Testing Edge Cases
```bash
# Test extreme working hours
cat > edge_cases.yaml <<EOF
developers:
  - user_id: night_owl
    email: owl@example.com
    working_hours:
      start: 22  # 10 PM
      end: 6     # 6 AM (crosses midnight)
    velocity: high
EOF

# Preview to see if midnight crossing works
./cursor-sim -mode preview -seed edge_cases.yaml

# Output shows commits spread across 22:00-06:00
# Validates temporal logic
```

### Scenario 3: Developer Comparing Configurations
```bash
# Preview config A
./cursor-sim -mode preview -seed config_a.yaml > preview_a.txt

# Preview config B
./cursor-sim -mode preview -seed config_b.yaml > preview_b.txt

# Compare distributions
diff preview_a.txt preview_b.txt
```

---

## Dependencies

- **Phase 1**: Seed loading infrastructure (COMPLETE ✅)
- **Phase 4**: CLI flag parsing (COMPLETE ✅)
- **New**: YAML parsing library (go-yaml/yaml)

---

## Risks & Mitigations

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| YAML library adds 5MB to binary | Medium | Low | Accept trade-off, it's standard library quality |
| Preview sample not representative | Medium | Medium | Use stratified sampling across developers |
| Users confuse preview with runtime | Low | Medium | Clear messaging, different output format |
| YAML syntax errors hard to debug | Medium | High | Show line numbers, context, helpful messages |

---

## Timeline Estimate

| Phase | Tasks | Hours |
|-------|-------|-------|
| Design & Planning | 1 | 0.5h |
| YAML Parser Integration | 2 | 2.0h |
| Preview Mode Implementation | 4 | 4.0h |
| Validation Framework | 2 | 2.0h |
| Testing & Documentation | 2 | 2.0h |
| **TOTAL** | **11** | **10.5h** |

---

**Next Step**: Create design.md with technical architecture and implementation details.
