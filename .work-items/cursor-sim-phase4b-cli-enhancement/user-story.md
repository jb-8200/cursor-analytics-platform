# User Story: Interactive CLI Configuration for Data Generation

**Feature ID**: cursor-sim-phase4b-cli-enhancement
**Created**: January 3, 2026
**Status**: Planning (Ready to Start)

---

## Executive Summary

As a **SDLC researcher** using cursor-sim to generate synthetic developer data, I want to **interactively configure data generation parameters through the CLI** so that I can **easily control dataset characteristics without manually editing seed files**, enabling rapid experimentation with different scenarios.

---

## User Stories (EARS Format)

### Story 1: Interactive Developer Count Configuration

**As a** researcher generating test datasets
**I want** to specify the number of developers interactively at CLI startup
**So that** I can quickly create datasets with varying team sizes without editing JSON files

**Acceptance Criteria**:

```gherkin
Given I start cursor-sim in runtime mode
When the CLI prompts me for the number of developers
And I see a suggested default (e.g., "2")
And I press Enter to accept the default OR type a different number
Then the system validates my input (1-100 developers)
And generates data for exactly that many developers
And rejects invalid inputs with clear error messages
```

### Story 2: Interactive Time Period Configuration (Months)

**As a** researcher studying long-term trends
**I want** to specify the generation period in months instead of days
**So that** I can think in natural time units for longitudinal studies

**Acceptance Criteria**:

```gherkin
Given I start cursor-sim in runtime mode
When the CLI prompts me for the period in months
And I see a suggested default (e.g., "3 months")
And I press Enter to accept OR type a different value (1-24 months)
Then the system converts months to days (month * 30)
And generates data spanning that time period
And rejects invalid inputs (< 1, > 24) with error messages
```

### Story 3: Interactive Commit Limit Configuration

**As a** researcher controlling dataset size
**I want** to set a maximum number of commits to generate
**So that** I can create predictable dataset sizes for performance testing

**Acceptance Criteria**:

```gherkin
Given I start cursor-sim in runtime mode
When the CLI prompts me for maximum commits
And I see a suggested default (e.g., "1000")
And I press Enter to accept OR type a different value (1-100000)
Then the system stops generation when reaching the commit limit
And the limit applies globally across all developers
And rejects invalid inputs with error messages
```

### Story 4: Input Validation with Regex

**As a** system protecting data quality
**I want** to validate all CLI inputs using regex patterns
**So that** invalid inputs are caught early with helpful feedback

**Acceptance Criteria**:

```gherkin
Given the CLI prompts me for any numeric parameter
When I enter non-numeric characters (e.g., "abc", "10x", "1.5.3")
Then the system rejects the input immediately
And displays a clear error message with the expected format
And re-prompts me to enter valid input
And allows up to 3 retry attempts before defaulting
```

### Story 5: Default Value Behavior

**As a** user wanting quick defaults
**I want** pressing Enter to accept sensible defaults for all parameters
**So that** I can quickly start generation without typing every value

**Acceptance Criteria**:

```gherkin
Given the CLI prompts me for any parameter
When I see the suggested default value (e.g., "[default: 2]")
And I press Enter without typing anything
Then the system uses the default value
And proceeds to the next prompt
And displays confirmation of the chosen value
```

---

## Problem Statement

### Current Issues

1. **Hardcoded Parameters**: Current implementation hardcodes:
   - 90 days of history (`-days` flag)
   - 2 developers (defined in seed file)
   - ~700 commits (emergent from Poisson process, uncontrolled)
2. **Poor UX**: Users must edit JSON seed files to change developer count
3. **Inflexible Time Units**: CLI uses days, but researchers think in months
4. **Unpredictable Dataset Size**: No way to cap total commits for consistent testing
5. **No Interactive Configuration**: All parameters must be specified via flags or seed files

### Root Causes

- CLI flags are passive (require explicit values)
- No interactive prompting mechanism
- Generator count controlled by seed file structure, not CLI
- Missing validation layer for inputs
- Time period hardcoded in days vs. months

---

## Goals & Non-Goals

### Goals

- ⏳ Interactive CLI prompts for: developer count, period (months), max commits
- ⏳ Press Enter to use sensible defaults
- ⏳ Regex validation for all numeric inputs
- ⏳ Clear error messages with retry logic
- ⏳ Developer replication from seed file (scale from 2 → N developers)
- ⏳ Maintain backward compatibility with existing `-days`, `-seed` flags

### Non-Goals

- ❌ GUI or web interface for configuration
- ❌ Persistent storage of user preferences
- ❌ Dynamic seed file generation (still require base seed file)
- ❌ Real-time parameter adjustment (post-startup)

---

## Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| User adoption | 90% of users prefer interactive mode | Usage telemetry |
| Input validation rate | < 5% of inputs invalid | Validation logs |
| Time to first prompt | < 2 seconds from startup | E2E tests |
| Default acceptance rate | > 70% of prompts use defaults | Behavioral testing |
| Developer scaling | Support 1-100 developers from seed | Integration tests |

---

## User Journey

### Before (Current State)

```bash
# Edit seed file manually
vim testdata/custom_seed.json  # Add more developers manually

# Run with hardcoded defaults
./bin/cursor-sim -mode runtime -seed testdata/custom_seed.json -days 90

# No control over commit count, no idea how many commits will be generated
```

### After (Enhanced CLI)

```bash
./bin/cursor-sim -mode runtime -seed testdata/valid_seed.json

cursor-sim v2.1.0
================

Configuration (press Enter for defaults):

Number of developers [default: 2]: 5
Period in months [default: 3]: 6
Maximum commits [default: 1000]: 2500

Validating inputs...
✓ Developers: 5
✓ Period: 6 months (180 days)
✓ Max commits: 2500

Loading seed data from testdata/valid_seed.json...
Generating 5 developers over 180 days (max 2500 commits)...
Generated 2500 commits across 5 developers

HTTP server listening on port 8080
```

---

## Dependencies

### Upstream (Blockers)

- None (independent feature)

### Downstream (Enabled)

- Improved research workflow velocity
- Easier onboarding for new users
- More reliable E2E test scenarios

---

## Open Questions

1. **Q**: Should we allow developers to be auto-generated from seed templates?
   **A**: Deferred to Phase 5. For now, CLI repeats/samples from seed file developers.

2. **Q**: Should max commits be per-developer or total?
   **A**: Total (global limit), simpler to reason about.

3. **Q**: Should we support fractional months (e.g., "1.5 months")?
   **A**: No, integers only (1-24 months). Use `-days` flag for fine-grained control.

4. **Q**: What happens if max commits is reached before period expires?
   **A**: Generation stops early, logs a message: "Reached max commits (N) before end of period"

---

## Related Work Items

- `.work-items/cursor-sim-v2/` - Phase 1 foundation
- `.work-items/cursor-sim-phase2/` - PR lifecycle
- `.work-items/cursor-sim-phase3/` - Research framework (current)
- `services/cursor-sim/SPEC.md` - Technical specification

---

**Next Steps**: Design technical approach → Task breakdown → Implementation
