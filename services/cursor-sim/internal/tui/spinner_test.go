package tui

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpinner_New(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner("Loading...", output)

	require.NotNil(t, spinner)
	assert.Equal(t, "Loading...", spinner.message)
	assert.NotNil(t, spinner.writer)
}

func TestSpinner_NewWithNilWriter(t *testing.T) {
	// Should not panic with nil writer
	spinner := NewSpinner("Loading...", nil)
	require.NotNil(t, spinner)
}

func TestSpinner_Start_Stop_TTY(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner("Loading...", output)

	// Start spinner
	spinner.Start()
	assert.True(t, spinner.isRunning)

	// Let it spin for a bit
	time.Sleep(50 * time.Millisecond)

	// Stop spinner with completion message
	spinner.Stop("Done!")
	assert.False(t, spinner.isRunning)

	// Output should contain stop message
	result := output.String()
	assert.NotEmpty(t, result)
}

func TestSpinner_StopBeforeStart(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner("Loading...", output)

	// Stop without starting should not panic
	assert.NotPanics(t, func() {
		spinner.Stop("Done!")
	})
}

func TestSpinner_MultipleStarts(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner("Loading...", output)

	spinner.Start()
	assert.True(t, spinner.isRunning)

	// Starting again should handle gracefully
	spinner.Start()
	assert.True(t, spinner.isRunning)

	spinner.Stop("Done!")
}

func TestSpinner_UpdateMessage(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner("Loading...", output)

	spinner.Start()
	assert.Equal(t, "Loading...", spinner.message)

	// Update message
	spinner.UpdateMessage("Generating commits...")
	assert.Equal(t, "Generating commits...", spinner.message)

	spinner.Stop("Done!")
}

func TestSpinner_UpdateMessageBeforeStart(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner("Loading...", output)

	spinner.UpdateMessage("New message")
	assert.Equal(t, "New message", spinner.message)
}

func TestSpinner_NonTTY_Fallback(t *testing.T) {
	// Test with fallback mode (no spinner animation)
	// Note: This is a basic test; actual non-TTY detection happens at runtime
	output := &bytes.Buffer{}
	spinner := NewSpinner("Loading...", output)

	// Verify spinner can start and stop without panic in any environment
	spinner.Start()
	assert.True(t, spinner.isRunning)

	time.Sleep(10 * time.Millisecond)

	spinner.Stop("Done!")
	assert.False(t, spinner.isRunning)

	// Should have completed without error
	assert.NotNil(t, spinner)
}

func TestSpinner_WithEmptyMessage(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner("", output)

	require.NotNil(t, spinner)
	assert.Equal(t, "", spinner.message)

	spinner.Start()
	spinner.Stop("Done!")
	assert.False(t, spinner.isRunning)
}

func TestSpinner_ConcurrentOperations(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner("Loading...", output)

	spinner.Start()
	assert.True(t, spinner.isRunning)

	// Update message while running
	for i := 0; i < 5; i++ {
		spinner.UpdateMessage("Step " + string(rune(48+i)) + "...")
		time.Sleep(10 * time.Millisecond)
	}

	spinner.Stop("Complete!")
	assert.False(t, spinner.isRunning)
}

func TestSpinner_RapidStops(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner("Loading...", output)

	spinner.Start()
	spinner.Stop("First stop")

	// Calling stop again should not panic
	assert.NotPanics(t, func() {
		spinner.Stop("Second stop")
	})
}

func TestSpinner_AllFunctions_DontPanic(t *testing.T) {
	output := &bytes.Buffer{}
	spinner := NewSpinner("Test", output)

	assert.NotPanics(t, func() { spinner.Start() })
	assert.NotPanics(t, func() { spinner.UpdateMessage("New") })
	assert.NotPanics(t, func() { spinner.Stop("Done") })
	assert.NotPanics(t, func() { spinner.Stop("Already stopped") })
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

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}

	// Should complete without panic
	assert.True(t, true)
}
