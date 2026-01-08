package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette for DOXAPI branding
var (
	// Primary gradient colors (purple → pink)
	PurpleColor = lipgloss.Color("#9B59B6")
	PinkColor   = lipgloss.Color("#FF69B4")

	// UI colors
	AccentColor  = lipgloss.Color("#00E5FF") // Cyan accent
	SuccessColor = lipgloss.Color("#00FF66") // Green checkmark
	ErrorColor   = lipgloss.Color("#FF6B6B") // Red error
	MutedColor   = lipgloss.Color("#888888") // Gray text
)

// Predefined styles for consistent UI
var (
	// TitleStyle for section headers
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PurpleColor)

	// SubtitleStyle for version info and subtitles
	SubtitleStyle = lipgloss.NewStyle().
			Faint(true).
			Foreground(MutedColor)

	// SuccessStyle for completion messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	// ErrorStyle for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	// PromptStyle for interactive prompts
	PromptStyle = lipgloss.NewStyle().
			Foreground(AccentColor)

	// HelpStyle for help text
	HelpStyle = lipgloss.NewStyle().
			Faint(true).
			Foreground(MutedColor).
			Italic(true)
)
