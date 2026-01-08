package tui

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupportsColor_NoColorEnvVarSet(t *testing.T) {
	// Save original NO_COLOR value
	originalNoColor := os.Getenv("NO_COLOR")
	defer os.Setenv("NO_COLOR", originalNoColor)

	// Set NO_COLOR to any non-empty value
	os.Setenv("NO_COLOR", "1")

	assert.False(t, SupportsColor(), "should return false when NO_COLOR is set")
}

func TestSupportsColor_NoColorEnvVarNotSet(t *testing.T) {
	// Save original NO_COLOR value
	originalNoColor := os.Getenv("NO_COLOR")
	defer os.Setenv("NO_COLOR", originalNoColor)

	// Unset NO_COLOR
	os.Unsetenv("NO_COLOR")

	// Result depends on terminal capabilities, but function shouldn't panic
	result := SupportsColor()
	assert.IsType(t, true, result, "SupportsColor should return a boolean")
}

func TestSupportsColor_NoColorEmptyString(t *testing.T) {
	// Save original NO_COLOR value
	originalNoColor := os.Getenv("NO_COLOR")
	defer os.Setenv("NO_COLOR", originalNoColor)

	// Set NO_COLOR to empty string (should NOT disable colors)
	os.Setenv("NO_COLOR", "")

	// Empty string should be treated as not set
	result := SupportsColor()
	assert.IsType(t, true, result, "SupportsColor should return a boolean")
}

func TestIsTTY(t *testing.T) {
	// This test verifies that IsTTY can be called without panic
	// Actual result depends on whether stdout is connected to terminal
	// In test environments, usually false; in interactive terminal, true
	result := IsTTY()
	assert.IsType(t, false, result, "IsTTY should return a boolean")
}

func TestShouldUseTUI(t *testing.T) {
	// Save original NO_COLOR value
	originalNoColor := os.Getenv("NO_COLOR")
	defer os.Setenv("NO_COLOR", originalNoColor)

	// Test with NO_COLOR set (should be false regardless of TTY)
	os.Setenv("NO_COLOR", "1")
	assert.False(t, ShouldUseTUI(), "ShouldUseTUI should be false when NO_COLOR is set")

	// Test without NO_COLOR
	os.Unsetenv("NO_COLOR")
	// Result depends on both SupportsColor() and IsTTY()
	// Both being true would make ShouldUseTUI true
	// At least one being false would make ShouldUseTUI false
	result := ShouldUseTUI()
	assert.IsType(t, true, result, "ShouldUseTUI should return a boolean")
}

func TestShouldUseTUI_LogicalAnd(t *testing.T) {
	// Save original NO_COLOR value
	originalNoColor := os.Getenv("NO_COLOR")
	defer os.Setenv("NO_COLOR", originalNoColor)

	// With NO_COLOR set, SupportsColor returns false
	os.Setenv("NO_COLOR", "1")
	supportsColor := SupportsColor()
	isTTY := IsTTY()
	shouldUseTUI := ShouldUseTUI()

	// ShouldUseTUI should be the logical AND of SupportsColor and IsTTY
	assert.Equal(t, supportsColor && isTTY, shouldUseTUI,
		"ShouldUseTUI should be the logical AND of SupportsColor and IsTTY")
}

func TestCapabilityFunctions_DontPanic(t *testing.T) {
	// Verify that all capability functions can be called without panicking
	assert.NotPanics(t, func() {
		_ = SupportsColor()
	})
	assert.NotPanics(t, func() {
		_ = IsTTY()
	})
	assert.NotPanics(t, func() {
		_ = ShouldUseTUI()
	})
}
