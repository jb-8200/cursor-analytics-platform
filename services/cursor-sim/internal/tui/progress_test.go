package tui

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgressBar_New(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("Generating commits", 100, output)

	require.NotNil(t, pb)
	assert.Equal(t, "Generating commits", pb.title)
	assert.Equal(t, 100, pb.total)
	assert.Equal(t, 0, pb.current)
}

func TestProgressBar_NewWithNilWriter(t *testing.T) {
	// Should not panic with nil writer
	pb := NewProgressBar("Test", 100, nil)
	require.NotNil(t, pb)
}

func TestProgressBar_Update(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("Generating", 10, output)

	// Update to 5/10
	pb.Update(5)
	assert.Equal(t, 5, pb.current)
	assert.Equal(t, 10, pb.total)

	// Update to 10/10
	pb.Update(10)
	assert.Equal(t, 10, pb.current)
}

func TestProgressBar_GetProgress(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("Test", 100, output)

	assert.Equal(t, 0, pb.GetProgress())

	pb.Update(50)
	assert.Equal(t, 50, pb.GetProgress())

	pb.Update(100)
	assert.Equal(t, 100, pb.GetProgress())
}

func TestProgressBar_GetPercentage(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("Test", 100, output)

	// 0/100 = 0%
	assert.Equal(t, 0, pb.GetPercentage())

	// 50/100 = 50%
	pb.Update(50)
	assert.Equal(t, 50, pb.GetPercentage())

	// 100/100 = 100%
	pb.Update(100)
	assert.Equal(t, 100, pb.GetPercentage())
}

func TestProgressBar_GetPercentage_Fractions(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("Test", 3, output)

	// 1/3 ~= 33%
	pb.Update(1)
	pct := pb.GetPercentage()
	assert.True(t, pct >= 33 && pct <= 34, "Expected ~33%, got %d", pct)

	// 2/3 ~= 67%
	pb.Update(2)
	pct = pb.GetPercentage()
	assert.True(t, pct >= 66 && pct <= 67, "Expected ~67%, got %d", pct)
}

func TestProgressBar_SetTitle(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("Initial", 100, output)

	assert.Equal(t, "Initial", pb.title)

	pb.SetTitle("Updated")
	assert.Equal(t, "Updated", pb.title)
}

func TestProgressBar_UpdateBeyondTotal(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("Test", 10, output)

	// Update to 15 (beyond total)
	pb.Update(15)
	// Should cap at total or handle gracefully
	assert.True(t, pb.current >= 10)
}

func TestProgressBar_Render(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("Test", 10, output)

	pb.Update(5)
	rendered := pb.Render()

	assert.NotEmpty(t, rendered)
	// Should contain progress bar characters or percentage
}

func TestProgressBar_RenderCompletion(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("Test", 10, output)

	pb.Update(10)
	rendered := pb.Render()

	assert.NotEmpty(t, rendered)
	// Completed progress bar should render differently
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

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should complete without panic or data race
	assert.True(t, pb.current >= 0)
}

func TestProgressBar_MultipleRenders(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("Multi", 100, output)

	// Render multiple times at different progress levels
	for i := 0; i <= 100; i += 10 {
		pb.Update(i)
		rendered := pb.Render()
		assert.NotEmpty(t, rendered)
	}
}

func TestProgressBar_ZeroTotal(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("Zero", 0, output)

	// Should handle zero total gracefully
	pct := pb.GetPercentage()
	assert.True(t, pct >= 0)
}

func TestProgressBar_WithEmptyTitle(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("", 100, output)

	assert.Equal(t, "", pb.title)

	pb.Update(50)
	rendered := pb.Render()
	assert.NotEmpty(t, rendered)
}

func TestProgressBar_AllFunctions_DontPanic(t *testing.T) {
	output := &bytes.Buffer{}
	pb := NewProgressBar("Test", 100, output)

	assert.NotPanics(t, func() { pb.Update(50) })
	assert.NotPanics(t, func() { _ = pb.GetProgress() })
	assert.NotPanics(t, func() { _ = pb.GetPercentage() })
	assert.NotPanics(t, func() { pb.SetTitle("New") })
	assert.NotPanics(t, func() { _ = pb.Render() })
}
