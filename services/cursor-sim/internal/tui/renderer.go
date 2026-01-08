package tui

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/events"
)

// Renderer listens to TUI events and manages spinners and progress displays.
// It acts as the central event handler for all TUI updates.
type Renderer struct {
	writer          io.Writer
	currentSpinner  *Spinner
	spinnerRunning  bool
	currentProgress *ProgressBar
	progressCurrent int
	progressTotal   int
	mu              sync.RWMutex
}

// NewRenderer creates a new Renderer with the given writer.
// Writer can be nil; defaults to os.Stdout.
func NewRenderer(writer io.Writer) *Renderer {
	if writer == nil {
		writer = os.Stdout
	}

	return &Renderer{
		writer:         writer,
		spinnerRunning: false,
	}
}

// HandleEvent processes an event and updates the UI accordingly.
// It's thread-safe and can be called from multiple goroutines.
func (r *Renderer) HandleEvent(event events.Event) {
	if event == nil {
		return
	}

	switch e := event.(type) {
	case events.PhaseStartEvent:
		r.handlePhaseStart(e)
	case events.PhaseCompleteEvent:
		r.handlePhaseComplete(e)
	case events.ProgressEvent:
		r.handleProgress(e)
	case events.WarningEvent:
		r.handleWarning(e)
	case events.ErrorEvent:
		r.handleError(e)
	default:
		// Unknown event type - silently ignore
	}
}

// handlePhaseStart starts a new spinner for the phase.
func (r *Renderer) handlePhaseStart(event events.PhaseStartEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Stop existing spinner if running
	if r.currentSpinner != nil && r.spinnerRunning {
		r.currentSpinner.Stop("Stopped")
	}

	// Start new spinner
	message := event.Message
	if message == "" {
		message = event.Phase
	}

	r.currentSpinner = NewSpinner(message, r.writer)
	r.currentSpinner.Start()
	r.spinnerRunning = true
}

// handlePhaseComplete stops the current spinner with a completion message.
func (r *Renderer) handlePhaseComplete(event events.PhaseCompleteEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentSpinner != nil && r.spinnerRunning {
		message := event.Message
		if message == "" {
			message = "Complete"
		}

		r.currentSpinner.Stop(message)
		r.spinnerRunning = false
		r.currentSpinner = nil
	}
}

// handleProgress updates the current spinner's message and tracks progress.
func (r *Renderer) handleProgress(event events.ProgressEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Track progress for percentage calculation
	r.progressCurrent = event.Current
	r.progressTotal = event.Total

	if r.currentSpinner != nil && r.spinnerRunning {
		message := event.Message
		if message == "" && event.Total > 0 {
			// Generate progress message if not provided
			message = fmt.Sprintf("%d/%d", event.Current, event.Total)
		}

		if message != "" {
			r.currentSpinner.UpdateMessage(message)
		}
	}
}

// handleWarning logs a warning message.
func (r *Renderer) handleWarning(event events.WarningEvent) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if event.Message != "" {
		fmt.Fprintf(r.writer, "⚠️  %s\n", event.Message)
	}
}

// handleError logs an error message.
func (r *Renderer) handleError(event events.ErrorEvent) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if event.Message != "" {
		fmt.Fprintf(r.writer, "❌ %s", event.Message)
		if event.Context != "" {
			fmt.Fprintf(r.writer, " (%s)", event.Context)
		}
		fmt.Fprintf(r.writer, "\n")
	}
}

// GetProgressPercentage returns the current progress as a percentage (0-100).
func (r *Renderer) GetProgressPercentage() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.progressTotal == 0 {
		return 0
	}

	return (r.progressCurrent * 100) / r.progressTotal
}
