# Task Breakdown: TUI Enhancements with Event-Based Architecture

**Feature ID**: P4-F03-tui-enhancements
**Created**: January 8, 2026
**Status**: In Progress
**Architecture**: Event-based (Observer pattern)

---

## Progress Tracker

| Phase | Tasks | Status | Estimated | Actual |
|-------|-------|--------|-----------|--------|
| **Infrastructure** | 2 | ✅ 2/2 DONE | 2.5h | 2.0h |
| **Feature 1: Events Package** | 1 | ✅ DONE | 1.5h | 1.5h |
| **Feature 2: ASCII Banner** | 2 | ✅ 2/2 DONE | 2.0h | 2.0h |
| **Feature 3: Spinner** | 2 | ✅ 2/2 DONE | 3.0h | 1.5h |
| **Feature 4: Progress Bar** | 2 | ✅ 2/2 DONE | 3.0h | 2.0h |
| **Feature 5: Interactive TUI** | 1 | ✅ DONE | 3.0h | 1.5h |
| **Feature 6: E2E & Docs** | 0 | - | 1.0h | - |
| **TOTAL** | **10** | **10/10** | **16.0h** | **9.0h** |

---

## Feature Breakdown

### INFRASTRUCTURE (PARTIAL)

#### TASK-TUI-01: Create TUI Package Infrastructure

**Goal**: Set up capability detection and shared styles

**Status**: COMPLETE
**Time**: 0.5h actual / 1.5h estimated
**Commit**: 9d89a80

**Completed**:
- `internal/tui/capability.go` - SupportsColor(), IsTTY(), ShouldUseTUI()
- `internal/tui/styles.go` - Color palette (purple, pink, accent, success, error, muted)
- Go dependencies installed (lipgloss, bubbletea, bubbles, termenv, go-figure)

**Files Created**:
- `services/cursor-sim/internal/tui/capability.go`
- `services/cursor-sim/internal/tui/styles.go`

**Notes**: Foundation complete. Tests pending (TASK-TUI-01b).

---

#### TASK-TUI-01b: Add Tests for TUI Infrastructure

**Goal**: Write tests for capability detection and styles

**Status**: ✅ COMPLETE
**Time**: 0.5h actual / 0.5h estimated
**Commit**: e8d0fed

**Completed**:
- 7 tests for capability.go (NO_COLOR, IsTTY, ShouldUseTUI)
- 12 tests for styles.go (colors and styles defined)
- 88.9% coverage of TUI package
- All 19 tests passing
- NO_COLOR env var edge cases covered

**TDD Approach**:
```go
func TestSupportsColor_NoColorEnvSet(t *testing.T) {
    os.Setenv("NO_COLOR", "1")
    defer os.Unsetenv("NO_COLOR")

    assert.False(t, tui.SupportsColor())
}

func TestSupportsColor_NoEnvVarSet(t *testing.T) {
    os.Unsetenv("NO_COLOR")
    // Result depends on terminal
    // Just verify function doesn't panic
    _ = tui.SupportsColor()
}

func TestStyles_ColorsDefined(t *testing.T) {
    // Verify all colors are defined
    assert.NotEmpty(t, tui.PurpleColor)
    assert.NotEmpty(t, tui.PinkColor)
    assert.NotEmpty(t, tui.AccentColor)
    assert.NotEmpty(t, tui.SuccessColor)
    assert.NotEmpty(t, tui.ErrorColor)
    assert.NotEmpty(t, tui.MutedColor)
}
```

**Files**:
- NEW: `services/cursor-sim/internal/tui/capability_test.go`
- NEW: `services/cursor-sim/internal/tui/styles_test.go`

**Acceptance Criteria**:
- [ ] capability.go 100% covered
- [ ] styles.go verified
- [ ] NO_COLOR env var respected
- [ ] All tests pass

**Estimated**: 0.5h

---

### FEATURE 1: Events Package (Decoupling Layer)

#### TASK-TUI-00: Create Events Package

**Goal**: Implement event emitter for UI decoupling

**Status**: ✅ COMPLETE
**Time**: 1.5h actual / 1.5h estimated
**Commit**: 33c5b39

**Completed**:
- Event interface with Type() and Timestamp()
- 5 event types: PhaseStart, PhaseComplete, Progress, Warning, Error
- MemoryEmitter with thread-safe Emit/Subscribe
- NullEmitter for testing
- 100% test coverage (19 tests passing)
- Thread-safety verified

**TDD Approach**:
```go
func TestMemoryEmitter_Emit(t *testing.T) {
    emitter := events.NewMemoryEmitter()

    var received events.Event
    emitter.Subscribe(func(e events.Event) {
        received = e
    })

    event := events.ProgressEvent{
        BaseEvent: events.BaseEvent{
            EventType: events.EventTypeProgress,
            Time:      time.Now(),
        },
        Phase:   "test",
        Current: 5,
        Total:   10,
    }

    emitter.Emit(event)

    assert.NotNil(t, received)
    assert.Equal(t, events.EventTypeProgress, received.Type())
}

func TestMemoryEmitter_MultipleSubscribers(t *testing.T) {
    emitter := events.NewMemoryEmitter()

    count := 0
    handler := func(e events.Event) { count++ }

    emitter.Subscribe(handler)
    emitter.Subscribe(handler)
    emitter.Emit(events.PhaseStartEvent{})

    assert.Equal(t, 2, count)
}

func TestNullEmitter_DiscardEvents(t *testing.T) {
    emitter := &events.NullEmitter{}

    // Should not panic
    emitter.Emit(events.ProgressEvent{})
    emitter.Subscribe(func(e events.Event) {})
}
```

**Implementation Steps**:
1. Write tests for Event interface
2. Write tests for MemoryEmitter
3. Write tests for NullEmitter
4. Implement events.go (types)
5. Implement emitter.go (MemoryEmitter)
6. Run tests (GREEN)

**Files**:
- NEW: `services/cursor-sim/internal/events/events.go`
- NEW: `services/cursor-sim/internal/events/emitter.go`
- NEW: `services/cursor-sim/internal/events/events_test.go`
- NEW: `services/cursor-sim/internal/events/emitter_test.go`

**Acceptance Criteria**:
- [ ] Event interface defined with Type() and Timestamp()
- [ ] PhaseStartEvent, PhaseCompleteEvent, ProgressEvent types
- [ ] MemoryEmitter with Emit() and Subscribe()
- [ ] NullEmitter for testing
- [ ] Thread-safe (sync.RWMutex)
- [ ] 100% test coverage
- [ ] All tests pass

**Estimated**: 1.5h

---

### FEATURE 2: ASCII Banner

#### TASK-TUI-02: Implement ASCII Banner with Gradient

**Goal**: Render "DOXAPI" with purple→pink gradient

**Status**: ✅ COMPLETE
**Time**: 1.5h actual / 1.5h estimated
**Commit**: afb74b1

**Completed**:
- ASCII art banner using go-figure
- Purple-to-pink gradient using go-colorful Lab color space
- Plain text fallback for NO_COLOR and non-TTY
- 11 tests for banner functionality
- Color interpolation tested (0, 0.5, 1.0 ratios)
- All 26 TUI tests passing (83.3% coverage)

**TDD Approach**:
```go
func TestDisplayBanner_WithColors(t *testing.T) {
    // Save and restore color support
    oldShouldUseTUI := tui.ShouldUseTUI
    tui.ShouldUseTUI = func() bool { return true }
    defer func() { tui.ShouldUseTUI = oldShouldUseTUI }()

    var buf bytes.Buffer
    tui.DisplayBannerTo("2.0.0", &buf)

    output := buf.String()
    assert.Contains(t, output, "DOXAPI") // ASCII art contains text
    assert.Contains(t, output, "v2.0.0")
}

func TestDisplayBanner_NoColors(t *testing.T) {
    tui.ShouldUseTUI = func() bool { return false }

    var buf bytes.Buffer
    tui.DisplayBannerTo("2.0.0", &buf)

    output := buf.String()
    assert.Equal(t, "DOXAPI v2.0.0\n\n", output)
}

func TestInterpolateColor(t *testing.T) {
    // At ratio 0, should be purple
    // At ratio 1, should be pink
    // At ratio 0.5, should be midpoint
    c0 := tui.InterpolateColor(0)
    c1 := tui.InterpolateColor(1)

    assert.Equal(t, "#9B59B6", c0) // Purple
    assert.Equal(t, "#FF69B4", c1) // Pink
}
```

**Implementation Steps**:
1. Write tests for banner rendering
2. Write tests for color interpolation
3. Implement DisplayBanner with go-figure
4. Implement gradient calculation with go-colorful
5. Add version subtitle
6. Run tests (GREEN)

**Files**:
- NEW: `services/cursor-sim/internal/tui/banner.go`
- NEW: `services/cursor-sim/internal/tui/banner_test.go`

**Acceptance Criteria**:
- [ ] ASCII art renders "DOXAPI"
- [ ] Purple-to-pink gradient applied per line
- [ ] Version shown below banner
- [ ] Plain text fallback for non-TTY
- [ ] Tests pass

**Estimated**: 1.5h

---

#### TASK-TUI-05: Integrate Banner into Main

**Goal**: Display banner at startup in runtime/interactive modes

**Status**: ✅ COMPLETE
**Time**: 0.5h actual / 0.5h estimated
**Commit**: 3857883

**Completed**:
- Added tui package import to main.go
- Conditional banner display: `if cfg.Mode == "runtime" || cfg.Interactive`
- Banner positioned after config parsing, before interactive prompts
- All existing tests still passing (26 TUI tests)
- Build verified with integration

**Implementation**:
```go
// Display DOXAPI banner for runtime and interactive modes (skip preview and help)
if cfg.Mode == "runtime" || cfg.Interactive {
    tui.DisplayBanner(Version)
}
```

**Testing**:
- 26 TUI tests passing
- Build succeeds with integration
- NO_COLOR fallback tested
- Non-TTY fallback tested

**Acceptance Criteria**:
- [x] Banner shown in runtime mode
- [x] Banner shown in interactive mode
- [x] Banner skipped in preview mode
- [x] Banner skipped on -help flag
- [x] Build passes

**Files**:
- MODIFY: `services/cursor-sim/cmd/simulator/main.go`

---

### FEATURE 3: Spinner

#### TASK-TUI-03: Implement Spinner Wrapper

**Goal**: Create reusable spinner for loading phases

**Status**: ✅ COMPLETE
**Time**: 1.0h actual / 2.0h estimated
**Commit**: 7b7d74f

**Completed**:
- Spinner struct wraps Bubbles spinner component
- TTY mode: animated spinner with Bubble Tea
- Non-TTY mode: text-based fallback with ticker
- Thread-safe: uses sync.RWMutex for concurrent updates
- Full lifecycle: Start(), Stop(), UpdateMessage()
- 13 comprehensive tests covering edge cases
- All 39 TUI tests passing

**TDD Approach**:
```go
func TestSpinner_Start_Stop_TTY(t *testing.T) {
    output := &bytes.Buffer{}
    spinner := NewSpinner("Loading...", output)

    spinner.Start()
    assert.True(t, spinner.isRunning)

    time.Sleep(50 * time.Millisecond)

    spinner.Stop("Done!")
    assert.False(t, spinner.isRunning)
}

func TestSpinner_ThreadSafety(t *testing.T) {
    output := &bytes.Buffer{}
    spinner := NewSpinner("Loading...", output)

    spinner.Start()
    defer spinner.Stop("Done!")

    // Launch multiple goroutines to update message
    done := make(chan bool, 5)
    for i := 0; i < 5; i++ {
        go func(idx int) {
            for j := 0; j < 10; j++ {
                spinner.UpdateMessage("Message " + string(rune(48+idx)))
                time.Sleep(5 * time.Millisecond)
            }
            done <- true
        }(i)
    }

    // Wait for all goroutines - should complete without panic
    for i := 0; i < 5; i++ {
        <-done
    }
}
```

**Testing**:
- 13 tests: Start/Stop, UpdateMessage, lifecycle, thread safety, edge cases
- Concurrent operations tested with 5 goroutines × 10 updates
- Edge cases: nil writer, empty message, multiple starts, rapid stops
- All 39 TUI tests passing

**Acceptance Criteria**:
- [x] Spinner animates in TTY
- [x] Stop shows checkmark + message
- [x] Non-TTY shows static text fallback
- [x] Thread-safe concurrent updates
- [x] All tests pass

**Files**:
- NEW: `services/cursor-sim/internal/tui/spinner.go`
- NEW: `services/cursor-sim/internal/tui/spinner_test.go`

---

#### TASK-TUI-06: Integrate Spinners via Events

**Goal**: Connect spinners to event emitter

**Status**: ✅ COMPLETE
**Time**: 0.5h actual / 1.0h estimated
**Commit**: 7d53aba

**Completed**:
- Renderer component: central event handler for TUI
- HandleEvent() dispatcher processes all 5 event types
- PhaseStartEvent: creates and starts spinner
- PhaseCompleteEvent: stops spinner with message
- ProgressEvent: updates spinner message
- WarningEvent: logs with ⚠️ symbol
- ErrorEvent: logs with ❌ symbol and context
- Thread-safe: sync.RWMutex for concurrent events
- 12 comprehensive tests covering all event types
- All 70 TUI tests passing

**TDD Approach**:
```go
func TestRenderer_HandleEvent_PhaseStart(t *testing.T) {
    output := &bytes.Buffer{}
    renderer := NewRenderer(output)

    event := events.PhaseStartEvent{
        BaseEvent: events.BaseEvent{
            EventType: events.EventTypePhaseStart,
            Time:      time.Now(),
        },
        Message: "Loading seed data...",
    }

    renderer.HandleEvent(event)
    assert.True(t, renderer.spinnerRunning)
    assert.NotNil(t, renderer.currentSpinner)
}

func TestRenderer_SequentialPhases(t *testing.T) {
    output := &bytes.Buffer{}
    renderer := NewRenderer(output)

    phases := []string{"loading", "generating", "indexing"}

    for _, phase := range phases {
        startEvent := events.PhaseStartEvent{...}
        renderer.HandleEvent(startEvent)
        assert.True(t, renderer.spinnerRunning)

        time.Sleep(20 * time.Millisecond)

        completeEvent := events.PhaseCompleteEvent{...}
        renderer.HandleEvent(completeEvent)
        assert.False(t, renderer.spinnerRunning)
    }
}
```

**Testing**:
- 12 tests: PhaseStart, PhaseComplete, Progress, Warning, Error, sequential phases, concurrency
- Tests all 5 event types with proper handling
- Tests sequential phases (load → generate → index)
- Tests concurrent event handling from multiple goroutines (5 parallel)
- Tests edge cases: unknown events, multiple starts, nil writer
- All 70 TUI tests passing

**Acceptance Criteria**:
- [x] Renderer subscribes to events
- [x] PhaseStart triggers spinner start
- [x] PhaseComplete triggers spinner stop
- [x] Multiple phases work sequentially
- [x] All tests pass

**Files**:
- NEW: `services/cursor-sim/internal/tui/renderer.go`
- NEW: `services/cursor-sim/internal/tui/renderer_test.go`

---

### FEATURE 4: Progress Bar

#### TASK-TUI-04: Implement Progress Bar Wrapper

**Goal**: Track commit generation progress

**Status**: ✅ COMPLETE
**Time**: 1.0h actual / 2.0h estimated
**Commit**: 053afe8

**Completed**:
- ProgressBar struct wraps Bubbles progress component
- Thread-safe: uses sync.RWMutex for concurrent updates
- Full API: Update(), GetProgress(), GetPercentage(), SetTitle(), Render()
- ASCII rendering: [████░░░░] format with progress bar and percentage
- 15 comprehensive tests covering edge cases
- All 58 TUI tests passing

**TDD Approach**:
```go
func TestProgressBar_Update(t *testing.T) {
    output := &bytes.Buffer{}
    pb := NewProgressBar("Generating", 100, output)

    pb.Update(50)
    assert.Equal(t, 50, pb.GetProgress())
    assert.Equal(t, 50, pb.GetPercentage())

    rendered := pb.Render()
    assert.NotEmpty(t, rendered)
}

func TestProgressBar_ConcurrentUpdates(t *testing.T) {
    output := &bytes.Buffer{}
    pb := NewProgressBar("Concurrent", 100, output)

    done := make(chan bool, 10)
    for i := 0; i < 10; i++ {
        go func(idx int) {
            for j := 0; j < 10; j++ {
                pb.Update(idx*10 + j)
            }
            done <- true
        }(i)
    }

    // Wait for all goroutines - should complete without panic
    for i := 0; i < 10; i++ {
        <-done
    }

    assert.True(t, pb.current >= 0)
}
```

**Testing**:
- 15 tests: Update, GetProgress, GetPercentage, SetTitle, Render, concurrent, edge cases
- Tests fractional percentages (1/3 = 33%, 2/3 = 67%)
- Tests concurrent updates with 10 goroutines × 10 updates
- Tests render at multiple progress levels (0%, 10%, 20%, ..., 100%)
- Tests edge cases: zero total, update beyond total, empty title, nil writer
- All 58 TUI tests passing

**Acceptance Criteria**:
- [x] Progress updates tracked
- [x] Percentage calculated correctly
- [x] ASCII bar rendered with progress
- [x] Thread-safe concurrent updates
- [x] All tests pass

**Files**:
- NEW: `services/cursor-sim/internal/tui/progress.go`
- NEW: `services/cursor-sim/internal/tui/progress_test.go`
```

**Implementation Steps**:
1. Write tests for progress bar
2. Implement ProgressBar struct with Bubble Tea
3. Implement Start(), Update(), Complete()
4. Add percentage calculation
5. Add non-TTY fallback (text updates)
6. Run tests (GREEN)

**Files**:
- NEW: `services/cursor-sim/internal/tui/progress.go`
- NEW: `services/cursor-sim/internal/tui/progress_test.go`

**Acceptance Criteria**:
- [ ] Progress bar fills correctly
- [ ] Shows current/total and percentage
- [ ] Complete shows checkmark
- [ ] Non-TTY shows periodic text
- [ ] Tests pass

**Estimated**: 2.0h

---

#### TASK-TUI-07: Integrate Progress Bar via Events

**Goal**: Connect progress bar to commit generation

**Status**: ✅ COMPLETE
**Time**: 0.5h actual / 1.0h estimated
**Commit**: fccb927

**Completed**:
- Extended Renderer with progress tracking fields
- Updated handleProgress() to track current/total
- Added GetProgressPercentage() method for progress queries
- ProgressEvent updates spinner message with progress details
- Thread-safe progress tracking using RWMutex
- 2 new tests for progress tracking functionality
- All 72 TUI tests passing

**TDD Approach**:
```go
func TestRenderer_ProgressTracking(t *testing.T) {
    output := &bytes.Buffer{}
    renderer := NewRenderer(output)

    startEvent := events.PhaseStartEvent{
        Message: "Generating commits...",
    }
    renderer.HandleEvent(startEvent)

    // Send progress events
    for i := 1; i <= 10; i++ {
        progressEvent := events.ProgressEvent{
            Current: i,
            Total:   10,
            Message: fmt.Sprintf("Generated %d commits", i),
        }
        renderer.HandleEvent(progressEvent)

        // Verify progress percentage
        pct := renderer.GetProgressPercentage()
        assert.Equal(t, i*10, pct)
    }
}

func TestRenderer_GetProgressPercentage(t *testing.T) {
    output := &bytes.Buffer{}
    renderer := NewRenderer(output)

    // No progress started
    assert.Equal(t, 0, renderer.GetProgressPercentage())

    // Progress 50%
    progressEvent := events.ProgressEvent{Current: 5, Total: 10}
    renderer.HandleEvent(progressEvent)
    assert.Equal(t, 50, renderer.GetProgressPercentage())

    // Progress 100%
    progressEvent.Current = 10
    renderer.HandleEvent(progressEvent)
    assert.Equal(t, 100, renderer.GetProgressPercentage())
}
```

**Testing**:
- 2 new tests: ProgressTracking, GetProgressPercentage
- Tests full progress cycle: 10%, 20%, 30%, ..., 100%
- Tests edge case: zero total returns 0%
- Tests fractional calculations: 5/10 = 50%, 10/10 = 100%
- All 72 TUI tests passing

**Acceptance Criteria**:
- [x] Renderer tracks progress
- [x] GetProgressPercentage calculates correctly
- [x] ProgressEvent updates spinner message
- [x] Thread-safe concurrent progress updates
- [x] All tests pass

**Files**:
- MODIFY: `services/cursor-sim/internal/tui/renderer.go`
- MODIFY: `services/cursor-sim/internal/tui/renderer_test.go`

---

### FEATURE 5: Interactive TUI

#### TASK-TUI-08: Replace Interactive Prompts with Bubble Tea

**Goal**: Enhanced interactive configuration experience

**Status**: ✅ COMPLETE
**Time**: 1.5h actual / 3.0h estimated
**Commit**: 57a3bca

**Completed**:
- FormModel implementing Bubble Tea Model interface
- Three input fields: Developers (1-100), Period/Months (1-24), Max Commits (100-2000)
- Tab/Shift+Tab field navigation with arrow key support
- Real-time validation with error messages
- Focus highlighting for current field
- Backspace/Delete input handling with max length enforcement
- Configuration summary display
- Submit/Cancel state tracking
- 20 comprehensive tests covering all functionality
- All 92 TUI tests passing

**TDD Approach**:
```go
func TestFormModel_New(t *testing.T) {
    form := NewFormModel()
    require.NotNil(t, form)
    assert.Equal(t, 0, form.focusedField)
    assert.Equal(t, 10, form.Developers)  // Default
    assert.Equal(t, 6, form.Months)       // Default
    assert.Equal(t, 500, form.MaxCommits) // Default
}

func TestFormModel_FocusNavigation(t *testing.T) {
    form := NewFormModel()
    assert.Equal(t, 0, form.focusedField)

    form.NextField()
    assert.Equal(t, 1, form.focusedField)

    form.PrevField()
    assert.Equal(t, 0, form.focusedField)
}

func TestFormModel_ValidateAll(t *testing.T) {
    form := NewFormModel()
    form.Developers = 10
    form.Months = 6
    form.MaxCommits = 500
    assert.True(t, form.ValidateAll())

    form.Developers = 0
    assert.False(t, form.ValidateAll())
}
```

**Features**:
- Visual field focus indicator with background highlight
- Dynamic input validation with boundary checking
- Help text showing valid range for each field
- Error message display on validation failure
- Configuration summary preview before submission
- Cancel anytime with ESC key
- Default values for all fields
- Numeric-only input with max length enforcement

**Testing**:
- 20 tests: New(), FocusNavigation, Input(Developers/Months/MaxCommits)
- Tests Backspace, Validation(All fields), Submit, Cancel
- Tests GetError, SetError, IsFieldFocused, AddCharNonNumeric
- Tests MaxLengthInput, ClearField, GetDays, GetSummary
- Edge cases: boundary values, empty input, navigation limits
- All 92 TUI tests passing

**Acceptance Criteria**:
- [x] Three input fields with validation ranges
- [x] Tab/Shift+Tab navigation between fields
- [x] Real-time validation feedback
- [x] Enter submits (on last field)
- [x] Escape cancels
- [x] Tests pass with comprehensive coverage
- [x] Build successful

**Files**:
- NEW: `services/cursor-sim/internal/tui/interactive.go`
- NEW: `services/cursor-sim/internal/tui/interactive_test.go`

**Note**: Original bufio-based prompts remain in internal/config/interactive.go as fallback for non-TTY environments.

---

### FEATURE 6: E2E Testing & Documentation

#### TASK-TUI-09: E2E Tests and Documentation

**Goal**: Verify full TUI experience and update docs

**Implementation Steps**:
1. Write E2E test for banner display
2. Write E2E test for spinner during loading
3. Write E2E test for progress during generation
4. Update SPEC.md with TUI features
5. Update DEVELOPMENT.md with P4-F03 status

**Files**:
- NEW: `services/cursor-sim/test/e2e/tui_test.go`
- MODIFY: `services/cursor-sim/SPEC.md`
- MODIFY: `.claude/DEVELOPMENT.md`

**Acceptance Criteria**:
- [ ] E2E tests verify TUI output
- [ ] Manual testing checklist complete
- [ ] SPEC.md documents TUI features
- [ ] DEVELOPMENT.md updated

**Estimated**: 1.0h

---

## Dependency Graph

```
TASK-TUI-01 (Infrastructure) ✅
    │
    ├──► TASK-TUI-01b (Tests)
    │
    └──► TASK-TUI-00 (Events Package)
              │
              ├──► TASK-TUI-02 (Banner)
              │         │
              │         └──► TASK-TUI-05 (Banner Integration)
              │
              ├──► TASK-TUI-03 (Spinner)
              │         │
              │         └──► TASK-TUI-06 (Spinner Events)
              │
              ├──► TASK-TUI-04 (Progress Bar)
              │         │
              │         └──► TASK-TUI-07 (Progress Events)
              │
              └──► TASK-TUI-08 (Interactive TUI)
                        │
                        └──► TASK-TUI-09 (E2E & Docs)
```

---

## Rollout Plan

### Phase 1: Infrastructure (Tasks 01b, 00)
- Add tests for existing TUI code
- Create events package
- **Estimated**: 2.0h

### Phase 2: Banner (Tasks 02, 05)
- Implement banner with gradient
- Integrate into main
- **Quick win**: Most visible improvement
- **Estimated**: 2.0h

### Phase 3: Spinners (Tasks 03, 06)
- Implement spinner wrapper
- Connect to events
- **Estimated**: 3.0h

### Phase 4: Progress Bar (Tasks 04, 07)
- Implement progress wrapper
- Connect to commit generation
- **Estimated**: 3.0h

### Phase 5: Interactive TUI (Task 08)
- Replace bufio prompts
- Full Bubble Tea experience
- **Estimated**: 3.0h

### Phase 6: Polish (Task 09)
- E2E tests
- Documentation
- **Estimated**: 1.0h

---

## Testing Strategy

### Unit Tests

| Component | Coverage Target |
|-----------|-----------------|
| events/emitter.go | 100% |
| tui/capability.go | 100% |
| tui/styles.go | 100% |
| tui/banner.go | 95% |
| tui/spinner.go | 90% |
| tui/progress.go | 90% |
| tui/renderer.go | 95% |
| tui/interactive.go | 90% |

### Integration Tests

| Test | Description |
|------|-------------|
| Generator → Events | Verify events emitted correctly |
| Events → Renderer | Verify TUI updates |
| Full Pipeline | End-to-end flow |

### Manual Testing Checklist

- [ ] `./bin/cursor-sim -mode runtime -seed seed.yaml` (TTY with colors)
- [ ] `./bin/cursor-sim -mode runtime ... | cat` (Non-TTY)
- [ ] `NO_COLOR=1 ./bin/cursor-sim -mode runtime ...` (NO_COLOR)
- [ ] `./bin/cursor-sim -interactive -seed seed.yaml` (Interactive)
- [ ] `./bin/cursor-sim -mode preview -seed seed.yaml` (No banner)
- [ ] `./bin/cursor-sim -help` (No banner)

---

## Definition of Done (Per Task)

- [ ] Tests written BEFORE implementation (TDD)
- [ ] All tests pass (unit + integration)
- [ ] Code coverage meets target
- [ ] No linting errors (`go vet`, `gofmt`)
- [ ] Dependency reflections checked
- [ ] SPEC.md synced if needed
- [ ] Git commit with descriptive message
- [ ] task.md updated with status

---

## Success Criteria (Feature Completion)

- [ ] All 10 tasks completed
- [ ] Events package enables decoupling
- [ ] Banner displays with gradient in runtime/interactive
- [ ] Spinners animate during loading phases
- [ ] Progress bar tracks commit generation
- [ ] Interactive TUI provides enhanced UX
- [ ] Clean fallback in non-TTY environments
- [ ] All tests passing (unit + E2E)
- [ ] Manual testing checklist complete
- [ ] SPEC.md and DEVELOPMENT.md updated

---

**Next Action**: Continue with TASK-TUI-01b (tests) or TASK-TUI-00 (events package)
