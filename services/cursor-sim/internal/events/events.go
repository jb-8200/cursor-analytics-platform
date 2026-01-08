// Package events provides event-based communication between business logic and UI layers.
// This enables decoupling of generators from TUI components, allowing future migration to web interfaces.
package events

import "time"

// EventType identifies the type of event
type EventType string

const (
	// EventTypePhaseStart signals the start of a named phase
	EventTypePhaseStart EventType = "phase_start"

	// EventTypePhaseComplete signals completion of a phase
	EventTypePhaseComplete EventType = "phase_complete"

	// EventTypeProgress reports incremental progress within a phase
	EventTypeProgress EventType = "progress"

	// EventTypeWarning reports non-fatal issues
	EventTypeWarning EventType = "warning"

	// EventTypeError reports fatal errors
	EventTypeError EventType = "error"
)

// Event is the base interface for all events.
// All event types must implement Type() and Timestamp().
type Event interface {
	Type() EventType
	Timestamp() time.Time
}

// BaseEvent provides common fields for all event types.
// Embed this in concrete event types to satisfy the Event interface.
type BaseEvent struct {
	EventType EventType `json:"type"`
	Time      time.Time `json:"timestamp"`
}

// Type returns the event type
func (e BaseEvent) Type() EventType {
	return e.EventType
}

// Timestamp returns the event timestamp
func (e BaseEvent) Timestamp() time.Time {
	return e.Time
}

// PhaseStartEvent signals the start of a named phase.
// Subscribers (like TUI) can start spinners or display "Loading..." messages.
type PhaseStartEvent struct {
	BaseEvent
	Phase   string `json:"phase"`   // e.g., "loading_seed", "generating_commits"
	Message string `json:"message"` // e.g., "Loading seed data..."
}

// PhaseCompleteEvent signals completion of a phase.
// Subscribers can stop spinners and display completion messages.
type PhaseCompleteEvent struct {
	BaseEvent
	Phase   string `json:"phase"`   // e.g., "loading_seed"
	Message string `json:"message"` // e.g., "Loaded 5 developers"
	Success bool   `json:"success"` // true if phase completed successfully
}

// ProgressEvent reports incremental progress within a phase.
// Subscribers can update progress bars or display "50/90 days" text.
type ProgressEvent struct {
	BaseEvent
	Phase   string `json:"phase"`             // e.g., "generating_commits"
	Current int    `json:"current"`           // Current step (e.g., day 45)
	Total   int    `json:"total"`             // Total steps (e.g., 90 days)
	Message string `json:"message,omitempty"` // Optional message (e.g., "Generated 500 commits")
}

// WarningEvent reports non-fatal issues.
// Examples: unknown model names, invalid working hours
type WarningEvent struct {
	BaseEvent
	Message string `json:"message"`           // Warning message
	Context string `json:"context,omitempty"` // Additional context (e.g., "developer alice")
}

// ErrorEvent reports fatal errors.
// Examples: database connection failure, file not found
type ErrorEvent struct {
	BaseEvent
	Message string `json:"message"`           // Error message
	Context string `json:"context,omitempty"` // Additional context (e.g., "storage initialization")
}
