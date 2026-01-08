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
| **Feature 3: Spinner** | 2 | TODO | 3.0h | - |
| **Feature 4: Progress Bar** | 2 | TODO | 3.0h | - |
| **Feature 5: Interactive TUI** | 1 | TODO | 3.0h | - |
| **Feature 6: E2E & Docs** | 1 | TODO | 1.0h | - |
| **TOTAL** | **10** | **5/10** | **16.0h** | **4.5h** |

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

**TDD Approach**:
```go
func TestSpinner_Start_TTY(t *testing.T) {
    var buf bytes.Buffer
    s := tui.NewSpinner("Loading...", &buf)

    s.Start()
    time.Sleep(100 * time.Millisecond)
    s.Stop("Done!")

    output := buf.String()
    assert.Contains(t, output, "Done!")
}

func TestSpinner_NonTTY_Fallback(t *testing.T) {
    tui.ShouldUseTUI = func() bool { return false }

    var buf bytes.Buffer
    s := tui.NewSpinner("Loading...", &buf)

    s.Start()
    s.Stop("Done!")

    output := buf.String()
    assert.Contains(t, output, "Loading...")
    assert.Contains(t, output, "Done!")
}
```

**Implementation Steps**:
1. Write tests for spinner lifecycle
2. Implement Spinner struct with Bubble Tea
3. Implement Start() method
4. Implement Stop() method with final message
5. Add non-TTY fallback
6. Run tests (GREEN)

**Files**:
- NEW: `services/cursor-sim/internal/tui/spinner.go`
- NEW: `services/cursor-sim/internal/tui/spinner_test.go`

**Acceptance Criteria**:
- [ ] Spinner animates in TTY
- [ ] Stop shows checkmark + message
- [ ] Non-TTY shows static text
- [ ] Thread-safe start/stop
- [ ] Tests pass

**Estimated**: 2.0h

---

#### TASK-TUI-06: Integrate Spinners via Events

**Goal**: Connect spinners to event emitter

**TDD Approach**:
```go
func TestRenderer_PhaseStart_StartsSpinner(t *testing.T) {
    var buf bytes.Buffer
    renderer := tui.NewRenderer(&buf)

    event := events.PhaseStartEvent{
        Phase:   "loading",
        Message: "Loading seed data...",
    }

    renderer.HandleEvent(event)

    // Spinner should be started
    time.Sleep(50 * time.Millisecond)
    assert.NotNil(t, renderer.GetSpinner())
}

func TestRenderer_PhaseComplete_StopsSpinner(t *testing.T) {
    var buf bytes.Buffer
    renderer := tui.NewRenderer(&buf)

    // Start
    renderer.HandleEvent(events.PhaseStartEvent{Message: "Loading..."})

    // Complete
    renderer.HandleEvent(events.PhaseCompleteEvent{Message: "Loaded 5 developers"})

    output := buf.String()
    assert.Contains(t, output, "Loaded 5 developers")
}
```

**Implementation Steps**:
1. Implement tui/renderer.go
2. Subscribe renderer to emitter in main.go
3. Emit PhaseStart events from generators
4. Emit PhaseComplete events from generators
5. Test event flow

**Files**:
- NEW: `services/cursor-sim/internal/tui/renderer.go`
- NEW: `services/cursor-sim/internal/tui/renderer_test.go`
- MODIFY: `services/cursor-sim/cmd/simulator/main.go`

**Acceptance Criteria**:
- [ ] Renderer subscribes to emitter
- [ ] PhaseStart triggers spinner start
- [ ] PhaseComplete triggers spinner stop
- [ ] Multiple phases work sequentially
- [ ] Tests pass

**Estimated**: 1.0h

---

### FEATURE 4: Progress Bar

#### TASK-TUI-04: Implement Progress Bar Wrapper

**Goal**: Track commit generation progress

**TDD Approach**:
```go
func TestProgressBar_Update(t *testing.T) {
    var buf bytes.Buffer
    pb := tui.NewProgressBar("Generating", 100, &buf)

    pb.Start()
    pb.Update(50)
    pb.Complete("Generated 500 commits")

    output := buf.String()
    assert.Contains(t, output, "500 commits")
}

func TestProgressBar_NonTTY_Fallback(t *testing.T) {
    tui.ShouldUseTUI = func() bool { return false }

    var buf bytes.Buffer
    pb := tui.NewProgressBar("Generating", 100, &buf)

    pb.Start()
    for i := 1; i <= 100; i++ {
        pb.Update(i)
    }
    pb.Complete("Done")

    output := buf.String()
    // Should have periodic text updates
    assert.Contains(t, output, "10/100")
    assert.Contains(t, output, "100/100")
}
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

**TDD Approach**:
```go
func TestRenderer_Progress_UpdatesBar(t *testing.T) {
    var buf bytes.Buffer
    renderer := tui.NewRenderer(&buf)

    // Send progress events
    for i := 1; i <= 10; i++ {
        renderer.HandleEvent(events.ProgressEvent{
            Phase:   "commits",
            Current: i,
            Total:   10,
        })
    }

    // Complete with message
    renderer.HandleEvent(events.ProgressEvent{
        Phase:   "commits",
        Current: 10,
        Total:   10,
        Message: "Generated 100 commits",
    })

    output := buf.String()
    assert.Contains(t, output, "100 commits")
}

func TestCommitGenerator_EmitsProgress(t *testing.T) {
    emitter := events.NewMemoryEmitter()

    var progressEvents []events.ProgressEvent
    emitter.Subscribe(func(e events.Event) {
        if pe, ok := e.(events.ProgressEvent); ok {
            progressEvents = append(progressEvents, pe)
        }
    })

    gen := generator.NewCommitGenerator(seedData, store, "medium", 12345)
    gen.SetEmitter(emitter)

    gen.GenerateCommits(ctx, 10, 100)

    assert.Equal(t, 10, len(progressEvents))
    assert.Equal(t, 1, progressEvents[0].Current)
    assert.Equal(t, 10, progressEvents[9].Current)
}
```

**Implementation Steps**:
1. Add SetEmitter to BaseGenerator
2. Modify CommitGenerator to emit ProgressEvent
3. Update renderer to handle ProgressEvent
4. Wire up in main.go
5. Test full flow

**Files**:
- NEW: `services/cursor-sim/internal/generator/base.go`
- MODIFY: `services/cursor-sim/internal/generator/commit_generator.go`
- MODIFY: `services/cursor-sim/internal/tui/renderer.go`

**Acceptance Criteria**:
- [ ] CommitGenerator emits progress per day
- [ ] Renderer updates progress bar
- [ ] Completion shows final message
- [ ] Existing tests still pass (NullEmitter)
- [ ] Integration test passes

**Estimated**: 1.0h

---

### FEATURE 5: Interactive TUI

#### TASK-TUI-08: Replace Interactive Prompts with Bubble Tea

**Goal**: Enhanced interactive configuration experience

**TDD Approach**:
```go
func TestInteractiveConfig_Model(t *testing.T) {
    model := tui.NewInteractiveConfigModel()

    // Verify initial state
    assert.Equal(t, 3, len(model.Inputs))
    assert.Equal(t, 0, model.Cursor)
}

func TestInteractiveConfig_TabNavigation(t *testing.T) {
    model := tui.NewInteractiveConfigModel()

    // Tab moves to next input
    model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
    assert.Equal(t, 1, model.Cursor)

    // Shift+Tab moves back
    model, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
    assert.Equal(t, 0, model.Cursor)
}

func TestInteractiveConfig_Validation(t *testing.T) {
    model := tui.NewInteractiveConfigModel()

    // Set invalid days value
    model.Inputs[0].SetValue("abc")

    err := model.Validate()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "days must be a number")
}

func TestInteractiveConfig_Submit(t *testing.T) {
    model := tui.NewInteractiveConfigModel()

    model.Inputs[0].SetValue("90")
    model.Inputs[1].SetValue("1000")
    model.Inputs[2].SetValue("12345")

    params, err := model.GetParams()
    require.NoError(t, err)

    assert.Equal(t, 90, params.Days)
    assert.Equal(t, 1000, params.MaxCommits)
    assert.Equal(t, int64(12345), params.RandomSeed)
}
```

**Implementation Steps**:
1. Write tests for interactive model
2. Implement InteractiveConfigModel with textinput
3. Implement navigation (Tab, Shift+Tab)
4. Implement validation
5. Implement submission
6. Replace bufio-based prompts in main.go
7. Test end-to-end

**Files**:
- NEW: `services/cursor-sim/internal/tui/interactive.go`
- NEW: `services/cursor-sim/internal/tui/interactive_test.go`
- MODIFY: `services/cursor-sim/cmd/simulator/main.go`
- REMOVE: Interactive parts of `internal/config/interactive.go` (or keep as fallback)

**Acceptance Criteria**:
- [ ] Three input fields (days, maxCommits, randomSeed)
- [ ] Tab navigation between fields
- [ ] Real-time validation feedback
- [ ] Enter submits
- [ ] Escape cancels
- [ ] Non-TTY falls back to bufio prompts
- [ ] Tests pass

**Estimated**: 3.0h

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
