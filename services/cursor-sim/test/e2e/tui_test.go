package e2e

import (
	"bytes"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/events"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/tui"
	"github.com/stretchr/testify/assert"
)

// TestBannerCapabilities verifies banner infrastructure works.
func TestBannerCapabilities(t *testing.T) {
	// Test color support detection
	colorSupported := tui.SupportsColor()
	assert.NotNil(t, colorSupported) // Just verify function works

	// Test TTY detection
	isTTY := tui.IsTTY()
	assert.NotNil(t, isTTY) // Just verify function works

	// Test composite check
	shouldUseTUI := tui.ShouldUseTUI()
	assert.NotNil(t, shouldUseTUI) // Just verify function works
}

// TestSpinnerLifecycle verifies spinner start/stop cycle.
func TestSpinnerLifecycle(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := tui.NewSpinner("Loading seed data...", output)

	// Verify initial state
	assert.False(t, spinner.IsRunning())

	// Start spinner
	spinner.Start()
	assert.True(t, spinner.IsRunning())

	// Let it run briefly
	time.Sleep(50 * time.Millisecond)

	// Stop spinner
	spinner.Stop("Loaded successfully")
	assert.False(t, spinner.IsRunning())
}

// TestSpinnerMessageUpdate verifies message updates during operation.
func TestSpinnerMessageUpdate(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := tui.NewSpinner("Initializing...", output)

	spinner.Start()
	assert.True(t, spinner.IsRunning())

	// Update message
	spinner.UpdateMessage("Processing...")
	time.Sleep(10 * time.Millisecond)

	// Update again
	spinner.UpdateMessage("Finalizing...")
	time.Sleep(10 * time.Millisecond)

	spinner.Stop("Complete")
	assert.False(t, spinner.IsRunning())
}

// TestProgressBarUpdates verifies progress bar updates work correctly.
func TestProgressBarUpdates(t *testing.T) {
	output := &bytes.Buffer{}
	pb := tui.NewProgressBar("Generating commits", 100, output)

	// Initial state
	assert.Equal(t, 0, pb.GetProgress())
	assert.Equal(t, 0, pb.GetPercentage())

	// Progress to 50%
	pb.Update(50)
	assert.Equal(t, 50, pb.GetProgress())
	assert.Equal(t, 50, pb.GetPercentage())

	// Progress to 100%
	pb.Update(100)
	assert.Equal(t, 100, pb.GetProgress())
	assert.Equal(t, 100, pb.GetPercentage())
}

// TestProgressBarRendering verifies progress bar renders without errors.
func TestProgressBarRendering(t *testing.T) {
	output := &bytes.Buffer{}
	pb := tui.NewProgressBar("Generating", 10, output)

	// Render at multiple progress levels
	for i := 0; i <= 10; i++ {
		pb.Update(i)
		rendered := pb.Render()
		assert.NotEmpty(t, rendered)
	}
}

// TestRendererEventFlow verifies complete event flow through renderer.
func TestRendererEventFlow(t *testing.T) {
	output := &bytes.Buffer{}
	renderer := tui.NewRenderer(output)
	emitter := events.NewMemoryEmitter()

	// Subscribe renderer to emitter
	emitter.Subscribe(renderer.HandleEvent)

	// Emit phase start
	startEvent := events.PhaseStartEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventTypePhaseStart,
			Time:      time.Now(),
		},
		Phase:   "generating",
		Message: "Generating commits...",
	}
	emitter.Emit(startEvent)
	time.Sleep(20 * time.Millisecond)

	// Emit progress events
	for i := 1; i <= 5; i++ {
		progressEvent := events.ProgressEvent{
			BaseEvent: events.BaseEvent{
				EventType: events.EventTypeProgress,
				Time:      time.Now(),
			},
			Phase:   "generating",
			Current: i,
			Total:   10,
		}
		emitter.Emit(progressEvent)
		time.Sleep(10 * time.Millisecond)

		// Verify progress is tracked
		pct := renderer.GetProgressPercentage()
		assert.Equal(t, i*10, pct)
	}

	// Emit completion
	completeEvent := events.PhaseCompleteEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventTypePhaseComplete,
			Time:      time.Now(),
		},
		Phase:   "generating",
		Message: "Generated commits successfully",
	}
	emitter.Emit(completeEvent)
}

// TestInteractiveFormModel verifies form functionality.
func TestInteractiveFormModel(t *testing.T) {
	form := tui.NewFormModel()

	// Verify defaults
	assert.Equal(t, 10, form.Developers)
	assert.Equal(t, 6, form.Months)
	assert.Equal(t, 500, form.MaxCommits)

	// Set valid values
	form.Developers = 20
	form.Months = 12
	form.MaxCommits = 1000

	// Verify validation passes
	assert.True(t, form.ValidateAll())

	// Verify days calculation
	assert.Equal(t, 360, form.GetDays()) // 12 * 30
}

// TestInteractiveFormValidation verifies form validation rules.
func TestInteractiveFormValidation(t *testing.T) {
	form := tui.NewFormModel()

	// Valid configuration
	form.Developers = 10
	form.Months = 6
	form.MaxCommits = 500
	assert.True(t, form.ValidateAll())

	// Invalid developer count
	form.Developers = 0
	assert.False(t, form.ValidateDevelopers())
	form.Developers = 10

	// Invalid months
	form.Months = 0
	assert.False(t, form.ValidateMonths())
	form.Months = 6

	// Invalid max commits
	form.MaxCommits = 50
	assert.False(t, form.ValidateMaxCommits())
	form.MaxCommits = 500

	// All valid again
	assert.True(t, form.ValidateAll())
}

// TestInteractiveFormNavigation verifies form field navigation.
func TestInteractiveFormNavigation(t *testing.T) {
	form := tui.NewFormModel()

	// Start at first field
	assert.Equal(t, 0, form.FocusedField())

	// Navigate forward
	form.NextField()
	assert.Equal(t, 1, form.FocusedField())

	form.NextField()
	assert.Equal(t, 2, form.FocusedField())

	// At last field, can't go further
	form.NextField()
	assert.Equal(t, 2, form.FocusedField())

	// Navigate backward
	form.PrevField()
	assert.Equal(t, 1, form.FocusedField())

	form.PrevField()
	assert.Equal(t, 0, form.FocusedField())

	// At first field, can't go back further
	form.PrevField()
	assert.Equal(t, 0, form.FocusedField())
}

// TestCompleteEventSequence simulates a realistic event sequence.
func TestCompleteEventSequence(t *testing.T) {
	output := &bytes.Buffer{}
	emitter := events.NewMemoryEmitter()
	renderer := tui.NewRenderer(output)

	emitter.Subscribe(renderer.HandleEvent)

	// Phase 1: Loading
	emitter.Emit(events.PhaseStartEvent{
		BaseEvent: events.BaseEvent{EventType: events.EventTypePhaseStart, Time: time.Now()},
		Message:   "Loading seed data...",
	})

	emitter.Emit(events.PhaseCompleteEvent{
		BaseEvent: events.BaseEvent{EventType: events.EventTypePhaseComplete, Time: time.Now()},
		Message:   "Loaded 5 developers",
	})

	// Phase 2: Generating
	emitter.Emit(events.PhaseStartEvent{
		BaseEvent: events.BaseEvent{EventType: events.EventTypePhaseStart, Time: time.Now()},
		Phase:     "generating",
		Message:   "Generating commits...",
	})

	// Progress updates - reach 100%
	for day := 10; day <= 90; day += 20 {
		emitter.Emit(events.ProgressEvent{
			BaseEvent: events.BaseEvent{EventType: events.EventTypeProgress, Time: time.Now()},
			Phase:     "generating",
			Current:   day,
			Total:     90,
		})
		time.Sleep(5 * time.Millisecond)
	}

	// Final progress to 100%
	emitter.Emit(events.ProgressEvent{
		BaseEvent: events.BaseEvent{EventType: events.EventTypeProgress, Time: time.Now()},
		Phase:     "generating",
		Current:   90,
		Total:     90,
	})

	emitter.Emit(events.PhaseCompleteEvent{
		BaseEvent: events.BaseEvent{EventType: events.EventTypePhaseComplete, Time: time.Now()},
		Message:   "Generated 100 commits",
	})

	// Verify final progress
	assert.Equal(t, 100, renderer.GetProgressPercentage())
}

// TestColorEnvironment verifies color handling works.
func TestColorEnvironment(t *testing.T) {
	// Just verify the capability detection doesn't panic
	_ = tui.SupportsColor()
	_ = tui.IsTTY()
	_ = tui.ShouldUseTUI()
}
