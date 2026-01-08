package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBaseEvent_Type(t *testing.T) {
	event := BaseEvent{
		EventType: EventTypePhaseStart,
		Time:      time.Now(),
	}

	assert.Equal(t, EventTypePhaseStart, event.Type())
}

func TestBaseEvent_Timestamp(t *testing.T) {
	now := time.Now()
	event := BaseEvent{
		EventType: EventTypeProgress,
		Time:      now,
	}

	assert.Equal(t, now, event.Timestamp())
}

func TestPhaseStartEvent_ImplementsEvent(t *testing.T) {
	event := PhaseStartEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypePhaseStart,
			Time:      time.Now(),
		},
		Phase:   "loading_seed",
		Message: "Loading seed data...",
	}

	// Should implement Event interface
	var _ Event = event

	assert.Equal(t, EventTypePhaseStart, event.Type())
	assert.Equal(t, "loading_seed", event.Phase)
	assert.Equal(t, "Loading seed data...", event.Message)
}

func TestPhaseCompleteEvent_ImplementsEvent(t *testing.T) {
	event := PhaseCompleteEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypePhaseComplete,
			Time:      time.Now(),
		},
		Phase:   "loading_seed",
		Message: "Loaded 5 developers",
		Success: true,
	}

	// Should implement Event interface
	var _ Event = event

	assert.Equal(t, EventTypePhaseComplete, event.Type())
	assert.Equal(t, "loading_seed", event.Phase)
	assert.Equal(t, "Loaded 5 developers", event.Message)
	assert.True(t, event.Success)
}

func TestProgressEvent_ImplementsEvent(t *testing.T) {
	event := ProgressEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeProgress,
			Time:      time.Now(),
		},
		Phase:   "generating_commits",
		Current: 45,
		Total:   90,
		Message: "Day 45/90",
	}

	// Should implement Event interface
	var _ Event = event

	assert.Equal(t, EventTypeProgress, event.Type())
	assert.Equal(t, "generating_commits", event.Phase)
	assert.Equal(t, 45, event.Current)
	assert.Equal(t, 90, event.Total)
	assert.Equal(t, "Day 45/90", event.Message)
}

func TestWarningEvent_ImplementsEvent(t *testing.T) {
	event := WarningEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeWarning,
			Time:      time.Now(),
		},
		Message: "Invalid model name",
		Context: "developer alice",
	}

	// Should implement Event interface
	var _ Event = event

	assert.Equal(t, EventTypeWarning, event.Type())
	assert.Equal(t, "Invalid model name", event.Message)
	assert.Equal(t, "developer alice", event.Context)
}

func TestErrorEvent_ImplementsEvent(t *testing.T) {
	event := ErrorEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeError,
			Time:      time.Now(),
		},
		Message: "Database connection failed",
		Context: "storage initialization",
	}

	// Should implement Event interface
	var _ Event = event

	assert.Equal(t, EventTypeError, event.Type())
	assert.Equal(t, "Database connection failed", event.Message)
	assert.Equal(t, "storage initialization", event.Context)
}

func TestEventType_Constants(t *testing.T) {
	// Verify event type constants are defined
	assert.Equal(t, EventType("phase_start"), EventTypePhaseStart)
	assert.Equal(t, EventType("phase_complete"), EventTypePhaseComplete)
	assert.Equal(t, EventType("progress"), EventTypeProgress)
	assert.Equal(t, EventType("warning"), EventTypeWarning)
	assert.Equal(t, EventType("error"), EventTypeError)
}
