# User Story: TUI Enhancements with Charmbracelet Stack

**Feature ID**: P4-F03-tui-enhancements
**Created**: January 8, 2026
**Status**: In Progress
**Branding**: DOXAPI

---

## Story (EARS Format)

### As-a / I-want / So-that

**As a** cursor-sim user (researcher, developer, QA engineer)
**I want** a polished terminal user interface with animated spinners, progress bars, and branded banner
**So that** I have visual feedback during data generation and a professional CLI experience

---

## Context

### Current State

1. **No Visual Branding**
   - CLI starts with plain "cursor-sim version X.X.X" text
   - No visual identity or professional appearance
   - Hard to distinguish from other tools

2. **Silent Generation Process**
   - No feedback during 90-day data generation (30-60 seconds)
   - User doesn't know if process is working or frozen
   - No progress indication for commit generation

3. **Basic Interactive Prompts**
   - Uses plain `bufio.Reader` for input
   - No validation feedback until enter pressed
   - Poor keyboard navigation

### Desired State (DOXAPI Branding)

1. **Branded Banner**
   - ASCII art "DOXAPI" text at startup
   - Purple-to-pink gradient (#9B59B6 → #FF69B4)
   - Version and copyright information
   - Displayed in runtime + interactive modes only

2. **Visual Progress Feedback**
   - Animated spinners during loading phases
   - Progress bars for commit generation (day-by-day)
   - Success/error indicators with color coding

3. **Enhanced Interactive Experience**
   - Bubble Tea-based prompts with real-time validation
   - Keyboard navigation between fields
   - Clear visual hierarchy

---

## Requirements

### Functional Requirements (EARS)

#### FR1: ASCII Banner Display
**WHEN** user runs `cursor-sim -mode runtime` or `-interactive`
**THEN** system displays DOXAPI ASCII banner with gradient
**AND** shows version information below banner
**AND** exits cleanly without displaying banner for `-help`

#### FR2: Banner Gradient
**WHEN** terminal supports colors
**THEN** banner displays with purple-to-pink vertical gradient
**AND** each line gets progressively pinker from top to bottom

#### FR3: Spinner During Loading
**WHEN** system is loading seed data or initializing generators
**THEN** animated spinner displays with context message
**AND** spinner stops with checkmark on completion

#### FR4: Progress Bar for Generation
**WHEN** system generates commits across N days
**THEN** progress bar shows current day / total days
**AND** updates smoothly as each day completes
**AND** shows completion percentage

#### FR5: Interactive TUI
**WHEN** user runs with `-interactive` flag
**THEN** Bubble Tea prompts guide configuration
**AND** fields validate in real-time
**AND** keyboard navigation works (Tab, Enter, Escape)

#### FR6: Terminal Fallback
**WHEN** terminal doesn't support colors (NO_COLOR set)
**THEN** banner displays as plain ASCII text
**AND** spinners show static text "Loading..."
**AND** progress shows periodic text updates

#### FR7: Non-TTY Fallback
**WHEN** output is piped or redirected
**THEN** all TUI features degrade to plain text
**AND** no ANSI escape codes in output

---

### Non-Functional Requirements

#### NFR1: Performance
- Banner render < 50ms
- Spinner CPU usage < 1%
- Progress updates smooth (60fps capable)

#### NFR2: Separation of Concerns
- TUI layer completely decoupled from business logic
- Generators emit events, TUI subscribes
- Web interface could replace TUI without generator changes

#### NFR3: Usability
- Banner fits 80-column terminal
- Colors accessible (4.5:1 contrast ratio)
- Clear success/error states

#### NFR4: Maintainability
- Use Charmbracelet stack consistently
- Shared color palette in styles.go
- Event-based architecture for extensibility

---

## Acceptance Criteria (Given-When-Then)

### AC1: Banner in Runtime Mode
**GIVEN** terminal supports colors
**WHEN** I run `./cursor-sim -mode runtime -seed seed.yaml`
**THEN** I see:
- DOXAPI banner with purple-to-pink gradient
- Version number below banner
- Normal runtime execution follows
- Exit code 0

### AC2: Banner Skipped for Help
**GIVEN** any terminal
**WHEN** I run `./cursor-sim -help`
**THEN** I see:
- Usage information
- Flag descriptions
- NO banner displayed
- Exit code 0

### AC3: Banner Skipped for Preview
**GIVEN** terminal supports colors
**WHEN** I run `./cursor-sim -mode preview -seed seed.yaml`
**THEN** I see:
- Preview mode output
- NO banner (preview needs clean output)
- Exit code 0

### AC4: Spinner During Loading
**GIVEN** runtime mode with colors
**WHEN** seed data is loading
**THEN** I see:
- Animated spinner (dots or line)
- Message "Loading seed data..."
- Checkmark "Loaded 5 developers" on completion

### AC5: Progress Bar During Generation
**GIVEN** runtime mode generating 90 days
**WHEN** commit generation runs
**THEN** I see:
- Progress bar filling left to right
- Text "Day 45/90 (50%)"
- Smooth updates each day
- Completion message

### AC6: Non-TTY Fallback
**GIVEN** output piped to file
**WHEN** I run `./cursor-sim -mode runtime ... | cat`
**THEN** I see:
- Plain text "DOXAPI v2.0.0"
- Static "Loading seed data..."
- Periodic "Progress: 50/90 days"
- No ANSI escape codes

### AC7: NO_COLOR Respected
**GIVEN** NO_COLOR=1 environment variable
**WHEN** I run `NO_COLOR=1 ./cursor-sim -mode runtime ...`
**THEN** I see:
- Plain ASCII banner (no colors)
- Plain text spinners
- Plain text progress
- Professional appearance maintained

### AC8: Interactive TUI
**GIVEN** TTY with colors
**WHEN** I run `./cursor-sim -interactive -seed seed.yaml`
**THEN** I can:
- Navigate between fields (Tab/Shift+Tab)
- See real-time validation
- Submit with Enter
- Cancel with Escape

---

## Out of Scope

- **Web UI integration**: Future enhancement (but architecture supports it)
- **Custom themes**: Single DOXAPI theme for now
- **Internationalization**: English only
- **Mouse support**: Keyboard navigation only
- **Configuration file for colors**: Hardcoded palette

---

## Success Metrics

- Banner renders correctly on macOS, Linux terminals
- 90%+ of users see colored output (TTY detection)
- Generation feels responsive (no "frozen" perception)
- Clean degradation in CI/CD environments

---

## User Scenarios

### Scenario 1: Developer First Run
```bash
# Terminal with iTerm2 + colors
./bin/cursor-sim -mode runtime -seed testdata/team.yaml -days 90

# Output:
#  ____   _____  __   __    _    ____  ___
# |  _ \ / _ \ \/ /  / \  |  _ \|_ _|
# | | | | | | \  /  / _ \ | |_) || |
# | |_| | |_| /  \ / ___ \|  __/ | |
# |____/ \___/_/\_\_/   \_\_|   |___|
#                        v2.0.0
#
# ⣾ Loading seed data...
# ✓ Loaded 5 developers
#
# ⣾ Initializing generators...
# ✓ Ready to generate
#
# Generating commits...
# [████████████████████░░░░░░░░░░░░░░░░░░░░] 45/90 days (50%)
#
# ✓ Generated 1,247 commits
# ✓ Server listening on :8080
```

### Scenario 2: CI/CD Pipeline
```bash
# Non-TTY environment (GitHub Actions)
./bin/cursor-sim -mode runtime -seed seed.yaml -days 90

# Output (plain text):
# DOXAPI v2.0.0
# Loading seed data...
# Loaded 5 developers
# Generating commits...
# Progress: 45/90 days
# Progress: 90/90 days
# Generated 1,247 commits
# Server listening on :8080
```

### Scenario 3: Interactive Configuration
```bash
# User runs interactive mode
./bin/cursor-sim -interactive -seed testdata/minimal.yaml

# TUI prompts:
# ┌─ Generation Configuration ─────────────────────┐
# │                                                 │
# │   Days to Generate: [90    ] ← cursor here     │
# │   Max Commits:      [1000  ]                   │
# │   Random Seed:      [12345 ]                   │
# │                                                 │
# │   [Start Generation]   [Cancel]                │
# │                                                 │
# └─────────────────────────────────────────────────┘
```

---

## Dependencies

- **Completed**: P4-F02 (CLI Enhancement with interactive prompts)
- **Completed**: P3-F04 (Preview mode)
- **New**: Charmbracelet libraries (lipgloss, bubbletea, bubbles, go-figure)

---

## Risks & Mitigations

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Library version conflicts | Low | Medium | Pin specific versions in go.mod |
| Performance on slow terminals | Low | Low | Throttle updates to 30fps |
| Colors look bad on light themes | Medium | Low | Test with multiple themes |
| Event architecture adds complexity | Medium | Medium | Keep emitter interface simple |

---

## Timeline Estimate

| Phase | Tasks | Hours |
|-------|-------|-------|
| Infrastructure | Events + Styles | 2.5h |
| Banner | Implementation + Integration | 2.0h |
| Spinners | Wrapper + Integration | 3.0h |
| Progress | Wrapper + Integration | 3.5h |
| Interactive TUI | Bubble Tea prompts | 3.0h |
| Testing | E2E + Manual | 2.0h |
| **TOTAL** | **10** | **16.0h** |

---

**Next Step**: Review design.md for technical architecture with event-based decoupling.
