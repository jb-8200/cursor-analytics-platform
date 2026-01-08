# Technical Design: TUI Enhancements with Event-Based Architecture

**Feature ID**: P4-F03-tui-enhancements
**Created**: January 8, 2026
**Status**: In Progress
**Architecture Level**: Service Enhancement (cursor-sim)

---

## Overview

Implement TUI enhancements using the Charmbracelet stack with **event-based architecture** to ensure clean separation between business logic and UI. This design ensures the CLI layer is decoupled from generators, enabling future migration to web interfaces without code changes.

---

## Goals

1. **Visual Polish**: DOXAPI branding, spinners, progress bars
2. **User Feedback**: Real-time progress during long operations
3. **Decoupled Architecture**: Generators emit events, UI subscribes
4. **Web-Ready**: Same events could feed WebSocket dashboard
5. **Graceful Degradation**: Works in CI/CD (non-TTY) environments

---

## Architecture Principles

### Separation of Concerns

```
┌────────────────────────────────────────────────────────────────────┐
│                        PRESENTATION LAYER                          │
│                                                                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                │
│  │   Banner    │  │  Spinners   │  │  Progress   │                │
│  │  (lipgloss) │  │ (bubbles)   │  │  (bubbles)  │                │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘                │
│         │                │                │                        │
│         └────────────────┴────────────────┘                        │
│                          │                                         │
│                  ┌───────▼───────┐                                 │
│                  │  TUI Renderer │ ◄──── Subscribes to events      │
│                  └───────┬───────┘                                 │
└──────────────────────────┼─────────────────────────────────────────┘
                           │
                    ┌──────▼──────┐
                    │   Events    │ ◄──── Decoupling boundary
                    │  Emitter    │
                    └──────┬──────┘
                           │
┌──────────────────────────┼─────────────────────────────────────────┐
│                          │                                         │
│                  ┌───────▼───────┐                                 │
│                  │  Generators   │ ◄──── Emit events, don't know   │
│                  │ CommitGen     │       about UI                  │
│                  │ ModelGen      │                                 │
│                  └───────────────┘                                 │
│                                                                     │
│                       BUSINESS LOGIC LAYER                         │
└────────────────────────────────────────────────────────────────────┘
```

### Why Event-Based?

**Problem with Callbacks**:
```go
// BAD: Couples generator to UI
func (g *CommitGenerator) GenerateCommits(
    ctx context.Context,
    days int,
    progressFn func(current, total int),  // <-- Tight coupling!
) error {
    for day := 0; day < days; day++ {
        // generate...
        progressFn(day+1, days)  // Generator knows about UI
    }
}
```

**Solution with Events**:
```go
// GOOD: Generator emits events, UI subscribes
func (g *CommitGenerator) GenerateCommits(
    ctx context.Context,
    days int,
) error {
    for day := 0; day < days; day++ {
        // generate...
        g.emitter.Emit(events.ProgressEvent{
            Phase:   "commits",
            Current: day + 1,
            Total:   days,
        })
    }
}

// UI subscribes separately
emitter.Subscribe(func(e events.Event) {
    if pe, ok := e.(events.ProgressEvent); ok {
        progressBar.Update(pe.Current)
    }
})
```

### Benefits

1. **Testability**: Generators test without UI
2. **Extensibility**: Add logging, metrics, WebSocket without generator changes
3. **Maintainability**: UI changes don't affect business logic
4. **Web Migration**: Same events feed React dashboard via WebSocket

---

## Component Design

### 1. Events Package (NEW)

**File**: `internal/events/events.go`

```go
package events

import "time"

// EventType identifies the type of event
type EventType string

const (
    EventTypePhaseStart    EventType = "phase_start"
    EventTypePhaseComplete EventType = "phase_complete"
    EventTypeProgress      EventType = "progress"
    EventTypeWarning       EventType = "warning"
    EventTypeError         EventType = "error"
)

// Event is the base interface for all events
type Event interface {
    Type() EventType
    Timestamp() time.Time
}

// BaseEvent provides common fields
type BaseEvent struct {
    EventType EventType `json:"type"`
    Time      time.Time `json:"timestamp"`
}

func (e BaseEvent) Type() EventType    { return e.EventType }
func (e BaseEvent) Timestamp() time.Time { return e.Time }

// PhaseStartEvent signals the start of a named phase
type PhaseStartEvent struct {
    BaseEvent
    Phase   string `json:"phase"`   // "loading_seed", "generating_commits"
    Message string `json:"message"` // "Loading seed data..."
}

// PhaseCompleteEvent signals completion of a phase
type PhaseCompleteEvent struct {
    BaseEvent
    Phase   string `json:"phase"`
    Message string `json:"message"` // "Loaded 5 developers"
    Success bool   `json:"success"`
}

// ProgressEvent reports progress within a phase
type ProgressEvent struct {
    BaseEvent
    Phase   string `json:"phase"`
    Current int    `json:"current"`
    Total   int    `json:"total"`
    Message string `json:"message,omitempty"`
}

// WarningEvent reports non-fatal issues
type WarningEvent struct {
    BaseEvent
    Message string `json:"message"`
    Context string `json:"context,omitempty"`
}
```

**File**: `internal/events/emitter.go`

```go
package events

import "sync"

// Handler processes events
type Handler func(Event)

// Emitter broadcasts events to subscribers
type Emitter interface {
    Emit(Event)
    Subscribe(Handler)
    Unsubscribe(Handler)
}

// MemoryEmitter is an in-memory event emitter
type MemoryEmitter struct {
    handlers []Handler
    mu       sync.RWMutex
}

func NewMemoryEmitter() *MemoryEmitter {
    return &MemoryEmitter{
        handlers: make([]Handler, 0),
    }
}

func (e *MemoryEmitter) Emit(event Event) {
    e.mu.RLock()
    defer e.mu.RUnlock()
    for _, h := range e.handlers {
        h(event)
    }
}

func (e *MemoryEmitter) Subscribe(h Handler) {
    e.mu.Lock()
    defer e.mu.Unlock()
    e.handlers = append(e.handlers, h)
}

// NullEmitter discards all events (for testing)
type NullEmitter struct{}

func (e *NullEmitter) Emit(Event)        {}
func (e *NullEmitter) Subscribe(Handler) {}
func (e *NullEmitter) Unsubscribe(Handler) {}
```

---

### 2. TUI Package Structure

```
internal/tui/
├── capability.go      # Terminal detection (COMPLETE ✅)
├── capability_test.go
├── styles.go          # Shared lipgloss styles (COMPLETE ✅)
├── styles_test.go
├── banner.go          # ASCII art + gradient
├── banner_test.go
├── spinner.go         # Bubbles spinner wrapper
├── spinner_test.go
├── progress.go        # Bubbles progress wrapper
├── progress_test.go
├── renderer.go        # Event subscriber that renders TUI
├── renderer_test.go
├── interactive.go     # Bubble Tea interactive prompts
└── interactive_test.go
```

---

### 3. Banner Component

**File**: `internal/tui/banner.go`

```go
package tui

import (
    "fmt"
    "strings"

    figure "github.com/common-nighthawk/go-figure"
    "github.com/charmbracelet/lipgloss"
    "github.com/lucasb-eyer/go-colorful"
)

// DisplayBanner renders the DOXAPI ASCII banner with gradient
func DisplayBanner(version string) {
    if !ShouldUseTUI() {
        // Plain text fallback
        fmt.Printf("DOXAPI v%s\n\n", version)
        return
    }

    // Generate ASCII art
    fig := figure.NewFigure("DOXAPI", "standard", true)
    asciiArt := fig.String()
    lines := strings.Split(asciiArt, "\n")

    // Apply gradient (purple → pink)
    purpleHex := "#9B59B6"
    pinkHex := "#FF69B4"
    purple, _ := colorful.Hex(purpleHex)
    pink, _ := colorful.Hex(pinkHex)

    for i, line := range lines {
        if line == "" {
            continue
        }
        // Interpolate color based on line position
        ratio := float64(i) / float64(len(lines)-1)
        color := purple.BlendLab(pink, ratio)
        hexColor := color.Hex()

        style := lipgloss.NewStyle().Foreground(lipgloss.Color(hexColor))
        fmt.Println(style.Render(line))
    }

    // Version info
    versionStyle := SubtitleStyle.Copy()
    fmt.Println(versionStyle.Render(fmt.Sprintf("v%s", version)))
    fmt.Println()
}
```

---

### 4. Spinner Component

**File**: `internal/tui/spinner.go`

```go
package tui

import (
    "fmt"
    "io"
    "sync"

    "github.com/charmbracelet/bubbles/spinner"
    tea "github.com/charmbracelet/bubbletea"
)

// Spinner provides animated loading indicator
type Spinner struct {
    message string
    writer  io.Writer
    program *tea.Program
    done    chan struct{}
    mu      sync.Mutex
}

// NewSpinner creates a new spinner with message
func NewSpinner(message string, writer io.Writer) *Spinner {
    return &Spinner{
        message: message,
        writer:  writer,
        done:    make(chan struct{}),
    }
}

// Start begins the spinner animation
func (s *Spinner) Start() {
    if !ShouldUseTUI() {
        // Plain text fallback
        fmt.Fprintf(s.writer, "%s...\n", s.message)
        return
    }

    model := spinnerModel{
        spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
        message: s.message,
    }
    model.spinner.Style = lipgloss.NewStyle().Foreground(AccentColor)

    s.program = tea.NewProgram(model, tea.WithOutput(s.writer))
    go func() {
        s.program.Run()
        close(s.done)
    }()
}

// Stop ends the spinner and shows completion message
func (s *Spinner) Stop(finalMsg string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if s.program != nil {
        s.program.Quit()
        <-s.done
    }

    // Clear spinner line and show completion
    successMsg := SuccessStyle.Render("✓") + " " + finalMsg
    fmt.Fprintln(s.writer, successMsg)
}

// spinnerModel implements tea.Model for bubble tea
type spinnerModel struct {
    spinner spinner.Model
    message string
}

func (m spinnerModel) Init() tea.Cmd {
    return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m, tea.Quit
    case spinner.TickMsg:
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd
    }
    return m, nil
}

func (m spinnerModel) View() string {
    return fmt.Sprintf("\r%s %s", m.spinner.View(), m.message)
}
```

---

### 5. Progress Bar Component

**File**: `internal/tui/progress.go`

```go
package tui

import (
    "fmt"
    "io"
    "sync"

    "github.com/charmbracelet/bubbles/progress"
    tea "github.com/charmbracelet/bubbletea"
)

// ProgressBar tracks progress of multi-step operations
type ProgressBar struct {
    title   string
    total   int
    current int
    writer  io.Writer
    program *tea.Program
    mu      sync.Mutex
}

// NewProgressBar creates a progress bar with total steps
func NewProgressBar(title string, total int, writer io.Writer) *ProgressBar {
    return &ProgressBar{
        title:  title,
        total:  total,
        writer: writer,
    }
}

// Start initializes the progress bar
func (p *ProgressBar) Start() {
    if !ShouldUseTUI() {
        fmt.Fprintf(p.writer, "%s: 0/%d\n", p.title, p.total)
        return
    }

    model := progressModel{
        progress: progress.New(
            progress.WithDefaultGradient(),
            progress.WithWidth(40),
        ),
        title: p.title,
        total: p.total,
    }

    p.program = tea.NewProgram(model, tea.WithOutput(p.writer))
    go p.program.Run()
}

// Update sets current progress
func (p *ProgressBar) Update(current int) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.current = current

    if !ShouldUseTUI() {
        // Periodic text updates for non-TTY
        if current%10 == 0 || current == p.total {
            fmt.Fprintf(p.writer, "%s: %d/%d\n", p.title, current, p.total)
        }
        return
    }

    if p.program != nil {
        p.program.Send(progressUpdateMsg{current: current})
    }
}

// Complete finishes the progress bar
func (p *ProgressBar) Complete(message string) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.program != nil {
        p.program.Quit()
    }

    fmt.Fprintln(p.writer)
    successMsg := SuccessStyle.Render("✓") + " " + message
    fmt.Fprintln(p.writer, successMsg)
}

type progressUpdateMsg struct{ current int }

type progressModel struct {
    progress progress.Model
    title    string
    total    int
    current  int
}

func (m progressModel) Init() tea.Cmd { return nil }

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case progressUpdateMsg:
        m.current = msg.current
        return m, nil
    case tea.KeyMsg:
        return m, tea.Quit
    }
    return m, nil
}

func (m progressModel) View() string {
    percent := float64(m.current) / float64(m.total)
    return fmt.Sprintf("\r%s %s %d/%d (%d%%)",
        m.title,
        m.progress.ViewAs(percent),
        m.current,
        m.total,
        int(percent*100),
    )
}
```

---

### 6. TUI Renderer (Event Subscriber)

**File**: `internal/tui/renderer.go`

```go
package tui

import (
    "io"
    "sync"

    "github.com/cursor-analytics-platform/services/cursor-sim/internal/events"
)

// Renderer subscribes to events and renders TUI components
type Renderer struct {
    writer   io.Writer
    spinner  *Spinner
    progress *ProgressBar
    mu       sync.Mutex
}

// NewRenderer creates a TUI renderer for the given writer
func NewRenderer(writer io.Writer) *Renderer {
    return &Renderer{writer: writer}
}

// HandleEvent processes events and updates TUI
func (r *Renderer) HandleEvent(e events.Event) {
    r.mu.Lock()
    defer r.mu.Unlock()

    switch event := e.(type) {
    case events.PhaseStartEvent:
        r.handlePhaseStart(event)
    case events.PhaseCompleteEvent:
        r.handlePhaseComplete(event)
    case events.ProgressEvent:
        r.handleProgress(event)
    }
}

func (r *Renderer) handlePhaseStart(e events.PhaseStartEvent) {
    // Stop any existing spinner
    if r.spinner != nil {
        r.spinner.Stop("")
    }

    // Start new spinner
    r.spinner = NewSpinner(e.Message, r.writer)
    r.spinner.Start()
}

func (r *Renderer) handlePhaseComplete(e events.PhaseCompleteEvent) {
    if r.spinner != nil {
        r.spinner.Stop(e.Message)
        r.spinner = nil
    }
}

func (r *Renderer) handleProgress(e events.ProgressEvent) {
    // Initialize progress bar on first progress event
    if r.progress == nil || r.progress.total != e.Total {
        if r.progress != nil {
            r.progress.Complete("")
        }
        r.progress = NewProgressBar(e.Phase, e.Total, r.writer)
        r.progress.Start()
    }
    r.progress.Update(e.Current)

    // Complete on reaching total
    if e.Current >= e.Total && e.Message != "" {
        r.progress.Complete(e.Message)
        r.progress = nil
    }
}
```

---

### 7. Generator Integration

**File**: `internal/generator/base.go` (NEW - base generator with emitter)

```go
package generator

import "github.com/cursor-analytics-platform/services/cursor-sim/internal/events"

// BaseGenerator provides common functionality for all generators
type BaseGenerator struct {
    emitter events.Emitter
}

// SetEmitter sets the event emitter (Dependency Injection)
func (g *BaseGenerator) SetEmitter(e events.Emitter) {
    g.emitter = e
}

// emit sends an event if emitter is set
func (g *BaseGenerator) emit(e events.Event) {
    if g.emitter != nil {
        g.emitter.Emit(e)
    }
}
```

**Modification**: `internal/generator/commit_generator.go`

```go
type CommitGenerator struct {
    BaseGenerator  // Embed for event support
    seedData *seed.SeedData
    store    storage.Store
    // ... existing fields
}

func (g *CommitGenerator) GenerateCommits(ctx context.Context, days int, maxCommits int) error {
    // Emit phase start
    g.emit(events.PhaseStartEvent{
        BaseEvent: events.BaseEvent{EventType: events.EventTypePhaseStart, Time: time.Now()},
        Phase:     "generating_commits",
        Message:   "Generating commits...",
    })

    for day := 0; day < days; day++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        // ... existing generation logic ...

        // Emit progress
        g.emit(events.ProgressEvent{
            BaseEvent: events.BaseEvent{EventType: events.EventTypeProgress, Time: time.Now()},
            Phase:     "generating_commits",
            Current:   day + 1,
            Total:     days,
        })
    }

    // Emit phase complete
    g.emit(events.PhaseCompleteEvent{
        BaseEvent: events.BaseEvent{EventType: events.EventTypePhaseComplete, Time: time.Now()},
        Phase:     "generating_commits",
        Message:   fmt.Sprintf("Generated %d commits", g.store.GetCommitCount()),
        Success:   true,
    })

    return nil
}
```

---

### 8. Main Integration

**File**: `cmd/simulator/main.go`

```go
func runRuntimeMode(seedData *seed.SeedData, cfg *config.Config) error {
    // Display banner (only in runtime + interactive modes)
    tui.DisplayBanner(Version)

    // Set up event emitter and TUI renderer
    emitter := events.NewMemoryEmitter()
    renderer := tui.NewRenderer(os.Stdout)
    emitter.Subscribe(renderer.HandleEvent)

    // Phase 1: Load developers
    emitter.Emit(events.PhaseStartEvent{
        BaseEvent: events.BaseEvent{EventType: events.EventTypePhaseStart, Time: time.Now()},
        Phase:     "loading_seed",
        Message:   "Loading seed data...",
    })

    developers := seedData.Developers
    devStore := storage.NewDeveloperStore()
    for _, dev := range developers {
        devStore.Add(dev)
    }

    emitter.Emit(events.PhaseCompleteEvent{
        BaseEvent: events.BaseEvent{EventType: events.EventTypePhaseComplete, Time: time.Now()},
        Phase:     "loading_seed",
        Message:   fmt.Sprintf("Loaded %d developers", len(developers)),
        Success:   true,
    })

    // Phase 2: Generate commits with progress
    commitGen := generator.NewCommitGenerator(seedData, store, cfg.Velocity, cfg.RandomSeed)
    commitGen.SetEmitter(emitter)  // Inject emitter

    if err := commitGen.GenerateCommits(ctx, cfg.Days, cfg.MaxCommits); err != nil {
        return err
    }

    // ... rest of runtime mode
}
```

---

## Data Flow

```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│                  │     │                  │     │                  │
│  CommitGenerator │────►│   MemoryEmitter  │────►│   TUI Renderer   │
│                  │     │                  │     │                  │
│  • Generate      │     │  • Emit()        │     │  • HandleEvent() │
│  • emit(event)   │     │  • Subscribe()   │     │  • Spinner       │
│                  │     │                  │     │  • ProgressBar   │
└──────────────────┘     └──────────────────┘     └──────────────────┘
                                │
                                │ (future)
                                ▼
                         ┌──────────────────┐
                         │                  │
                         │  WebSocket Relay │
                         │                  │
                         │  • Same events   │
                         │  • React UI      │
                         │                  │
                         └──────────────────┘
```

---

## Testing Strategy

### Unit Tests

| Component | Test Coverage |
|-----------|---------------|
| `events/emitter.go` | Emit, Subscribe, Unsubscribe, thread-safety |
| `tui/banner.go` | Gradient calculation, plain text fallback |
| `tui/spinner.go` | Start/Stop lifecycle, non-TTY fallback |
| `tui/progress.go` | Update tracking, percentage calculation |
| `tui/renderer.go` | Event handling, state transitions |

### Integration Tests

| Test | Scope |
|------|-------|
| Generator with emitter | CommitGenerator → Events → Assertions |
| Renderer with events | Events → Renderer → Output |
| Full pipeline | Generator → Emitter → Renderer → Output |

### Manual Testing Checklist

- [ ] TTY with colors: `./bin/cursor-sim -mode runtime ...`
- [ ] Non-TTY: `./bin/cursor-sim ... | cat`
- [ ] NO_COLOR: `NO_COLOR=1 ./bin/cursor-sim ...`
- [ ] Interactive mode: `./bin/cursor-sim -interactive ...`
- [ ] Preview mode (no banner): `./bin/cursor-sim -mode preview ...`
- [ ] Help flag (no banner): `./bin/cursor-sim -help`

---

## Backward Compatibility

### API Layer
- **Zero impact**: Events are internal, API responses unchanged
- **No generator signature changes**: `GenerateCommits(ctx, days, maxCommits)` same

### CLI Flags
- **No new required flags**: TUI is automatic
- **Optional environment**: `NO_COLOR=1` disables colors

### Testing
- **NullEmitter for tests**: Generators work without UI
- **Existing tests pass**: No callback parameters needed

---

## Future Extensions

### WebSocket Dashboard (P6+)

```go
// Future: WebSocket relay
wsRelay := websocket.NewEventRelay(conn)
emitter.Subscribe(wsRelay.HandleEvent)

// React component receives same events:
// { type: "progress", phase: "commits", current: 45, total: 90 }
```

### Metrics Collection

```go
// Future: Prometheus metrics
metricsHandler := metrics.NewEventHandler(registry)
emitter.Subscribe(metricsHandler.HandleEvent)
```

### Structured Logging

```go
// Future: JSON logging
logger := logging.NewEventLogger(os.Stderr)
emitter.Subscribe(logger.HandleEvent)
```

---

## Success Criteria

- [ ] Events package with Emitter interface
- [ ] TUI components (banner, spinner, progress)
- [ ] Renderer subscribes to events
- [ ] Generators emit events without UI knowledge
- [ ] All tests pass (unit + integration)
- [ ] Manual testing checklist complete
- [ ] Clean fallback in non-TTY environments

---

**Next Step**: Create task.md with implementation tasks following TDD approach.
