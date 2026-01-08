package tui

import (
	"bytes"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderer_New(t *testing.T) {
	output := &bytes.Buffer{}
	renderer := NewRenderer(output)

	require.NotNil(t, renderer)
	assert.NotNil(t, renderer.writer)
	assert.Nil(t, renderer.currentSpinner)
}

func TestRenderer_NewWithNilWriter(t *testing.T) {
	// Should not panic with nil writer
	renderer := NewRenderer(nil)
	require.NotNil(t, renderer)
}

func TestRenderer_HandleEvent_PhaseStart(t *testing.T) {
	output := &bytes.Buffer{}
	renderer := NewRenderer(output)

	event := events.PhaseStartEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventTypePhaseStart,
			Time:      time.Now(),
		},
		Phase:   "loading",
		Message: "Loading seed data...",
	}

	renderer.HandleEvent(event)

	// Spinner should be started
	assert.True(t, renderer.spinnerRunning)
	assert.NotNil(t, renderer.currentSpinner)
}

func TestRenderer_HandleEvent_PhaseComplete(t *testing.T) {
	output := &bytes.Buffer{}
	renderer := NewRenderer(output)

	// Start phase
	startEvent := events.PhaseStartEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventTypePhaseStart,
			Time:      time.Now(),
		},
		Phase:   "loading",
		Message: "Loading...",
	}
	renderer.HandleEvent(startEvent)
	assert.True(t, renderer.spinnerRunning)

	// Let spinner run briefly
	time.Sleep(50 * time.Millisecond)

	// Complete phase
	completeEvent := events.PhaseCompleteEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventTypePhaseComplete,
			Time:      time.Now(),
		},
		Phase:   "loading",
		Message: "Loaded successfully",
	}
	renderer.HandleEvent(completeEvent)

	// Spinner should be stopped
	assert.False(t, renderer.spinnerRunning)
}

func TestRenderer_HandleEvent_Progress(t *testing.T) {
	output := &bytes.Buffer{}
	renderer := NewRenderer(output)

	// Start phase first
	startEvent := events.PhaseStartEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventTypePhaseStart,
			Time:      time.Now(),
		},
		Phase:   "generating",
		Message: "Generating commits...",
	}
	renderer.HandleEvent(startEvent)

	// Send progress events
	for i := 1; i <= 5; i++ {
		progressEvent := events.ProgressEvent{
			BaseEvent: events.BaseEvent{
				EventType: events.EventTypeProgress,
				Time:      time.Now(),
			},
			Phase:   "generating",
			Current: i,
			Total:   10,
			Message: "Generating commits...",
		}
		renderer.HandleEvent(progressEvent)
	}

	// Spinner should still be running
	assert.True(t, renderer.spinnerRunning)
}

func TestRenderer_HandleEvent_Warning(t *testing.T) {
	output := &bytes.Buffer{}
	renderer := NewRenderer(output)

	warningEvent := events.WarningEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventTypeWarning,
			Time:      time.Now(),
		},
		Message: "This is a warning",
	}

	// Should not panic
	assert.NotPanics(t, func() {
		renderer.HandleEvent(warningEvent)
	})
}

func TestRenderer_HandleEvent_Error(t *testing.T) {
	output := &bytes.Buffer{}
	renderer := NewRenderer(output)

	errorEvent := events.ErrorEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventTypeError,
			Time:      time.Now(),
		},
		Message: "An error occurred",
		Context: "test context",
	}

	// Should not panic
	assert.NotPanics(t, func() {
		renderer.HandleEvent(errorEvent)
	})
}

func TestRenderer_SequentialPhases(t *testing.T) {
	output := &bytes.Buffer{}
	renderer := NewRenderer(output)

	phases := []string{"loading", "generating", "indexing"}

	for _, phase := range phases {
		// Start phase
		startEvent := events.PhaseStartEvent{
			BaseEvent: events.BaseEvent{
				EventType: events.EventTypePhaseStart,
				Time:      time.Now(),
			},
			Phase:   phase,
			Message: "Processing " + phase + "...",
		}
		renderer.HandleEvent(startEvent)
		assert.True(t, renderer.spinnerRunning)

		// Simulate work
		time.Sleep(20 * time.Millisecond)

		// Complete phase
		completeEvent := events.PhaseCompleteEvent{
			BaseEvent: events.BaseEvent{
				EventType: events.EventTypePhaseComplete,
				Time:      time.Now(),
			},
			Phase:   phase,
			Message: "Completed " + phase,
		}
		renderer.HandleEvent(completeEvent)
		assert.False(t, renderer.spinnerRunning)

		// Small delay before next phase
		time.Sleep(10 * time.Millisecond)
	}
}

func TestRenderer_UnknownEventType(t *testing.T) {
	output := &bytes.Buffer{}
	renderer := NewRenderer(output)

	// Create a generic event with unknown type
	unknownEvent := events.PhaseStartEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventType("unknown_type"),
			Time:      time.Now(),
		},
	}

	// Should not panic on unknown event type
	assert.NotPanics(t, func() {
		renderer.HandleEvent(unknownEvent)
	})
}

func TestRenderer_MultipleStartEventsWithoutStop(t *testing.T) {
	output := &bytes.Buffer{}
	renderer := NewRenderer(output)

	// Start first phase
	startEvent1 := events.PhaseStartEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventTypePhaseStart,
			Time:      time.Now(),
		},
		Phase:   "phase1",
		Message: "Phase 1...",
	}
	renderer.HandleEvent(startEvent1)
	assert.True(t, renderer.spinnerRunning)

	// Start second phase without completing first
	startEvent2 := events.PhaseStartEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventTypePhaseStart,
			Time:      time.Now(),
		},
		Phase:   "phase2",
		Message: "Phase 2...",
	}
	renderer.HandleEvent(startEvent2)

	// Should still be running (new spinner replaces old one)
	assert.True(t, renderer.spinnerRunning)
}

func TestRenderer_AllFunctions_DontPanic(t *testing.T) {
	output := &bytes.Buffer{}
	renderer := NewRenderer(output)

	events := []events.Event{
		events.PhaseStartEvent{
			BaseEvent: events.BaseEvent{
				EventType: events.EventTypePhaseStart,
				Time:      time.Now(),
			},
			Message: "Start",
		},
		events.ProgressEvent{
			BaseEvent: events.BaseEvent{
				EventType: events.EventTypeProgress,
				Time:      time.Now(),
			},
			Current: 5,
			Total:   10,
		},
		events.PhaseCompleteEvent{
			BaseEvent: events.BaseEvent{
				EventType: events.EventTypePhaseComplete,
				Time:      time.Now(),
			},
			Message: "Complete",
		},
		events.WarningEvent{
			BaseEvent: events.BaseEvent{
				EventType: events.EventTypeWarning,
				Time:      time.Now(),
			},
			Message: "Warning",
		},
		events.ErrorEvent{
			BaseEvent: events.BaseEvent{
				EventType: events.EventTypeError,
				Time:      time.Now(),
			},
			Message: "Error",
			Context: "test",
		},
	}

	for _, e := range events {
		assert.NotPanics(t, func() {
			renderer.HandleEvent(e)
		})
	}
}

func TestRenderer_ThreadSafety(t *testing.T) {
	output := &bytes.Buffer{}
	renderer := NewRenderer(output)

	done := make(chan bool, 5)

	// Launch multiple goroutines sending events
	for i := 0; i < 5; i++ {
		go func(idx int) {
			startEvent := events.PhaseStartEvent{
				BaseEvent: events.BaseEvent{
					EventType: events.EventTypePhaseStart,
					Time:      time.Now(),
				},
				Message: "Phase started",
			}
			renderer.HandleEvent(startEvent)

			for j := 0; j < 3; j++ {
				progressEvent := events.ProgressEvent{
					BaseEvent: events.BaseEvent{
						EventType: events.EventTypeProgress,
						Time:      time.Now(),
					},
					Current: j,
					Total:   3,
				}
				renderer.HandleEvent(progressEvent)
				time.Sleep(5 * time.Millisecond)
			}

			completeEvent := events.PhaseCompleteEvent{
				BaseEvent: events.BaseEvent{
					EventType: events.EventTypePhaseComplete,
					Time:      time.Now(),
				},
				Message: "Phase complete",
			}
			renderer.HandleEvent(completeEvent)

			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}

	// Should complete without panic or race condition
	assert.True(t, true)
}
