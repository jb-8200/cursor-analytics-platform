// Package tui provides terminal user interface components for cursor-sim.
// Includes banner, spinners, progress bars, and interactive prompts using Charmbracelet stack.
package tui

import (
	"os"

	"github.com/muesli/termenv"
)

// SupportsColor checks if the terminal supports color output.
// Returns false if NO_COLOR environment variable is set or terminal doesn't support colors.
func SupportsColor() bool {
	// Respect NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check terminal color profile
	profile := termenv.ColorProfile()
	return profile != termenv.Ascii
}

// IsTTY checks if stdout is connected to a terminal (TTY).
// Returns false for piped output, redirects, or non-interactive environments.
func IsTTY() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// ShouldUseTUI determines if full TUI features should be enabled.
// Returns true only if both color support and TTY are available.
func ShouldUseTUI() bool {
	return SupportsColor() && IsTTY()
}
