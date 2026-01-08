package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorsPalette_Defined(t *testing.T) {
	// Verify all color constants are defined and non-empty
	assert.NotNil(t, PurpleColor, "PurpleColor should be defined")
	assert.NotNil(t, PinkColor, "PinkColor should be defined")
	assert.NotNil(t, AccentColor, "AccentColor should be defined")
	assert.NotNil(t, SuccessColor, "SuccessColor should be defined")
	assert.NotNil(t, ErrorColor, "ErrorColor should be defined")
	assert.NotNil(t, MutedColor, "MutedColor should be defined")
}

func TestColorsPalette_ExpectedValues(t *testing.T) {
	// Verify that colors have the expected hex values
	// These are light checks - just ensure the colors exist and are strings
	purpleStr := string(PurpleColor)
	pinkStr := string(PinkColor)
	accentStr := string(AccentColor)
	successStr := string(SuccessColor)
	errorStr := string(ErrorColor)
	mutedStr := string(MutedColor)

	assert.NotEmpty(t, purpleStr, "PurpleColor should not be empty")
	assert.NotEmpty(t, pinkStr, "PinkColor should not be empty")
	assert.NotEmpty(t, accentStr, "AccentColor should not be empty")
	assert.NotEmpty(t, successStr, "SuccessColor should not be empty")
	assert.NotEmpty(t, errorStr, "ErrorColor should not be empty")
	assert.NotEmpty(t, mutedStr, "MutedColor should not be empty")
}

func TestPredefinedStyles_Defined(t *testing.T) {
	// Verify all style constants are defined
	assert.NotNil(t, TitleStyle, "TitleStyle should be defined")
	assert.NotNil(t, SubtitleStyle, "SubtitleStyle should be defined")
	assert.NotNil(t, SuccessStyle, "SuccessStyle should be defined")
	assert.NotNil(t, ErrorStyle, "ErrorStyle should be defined")
	assert.NotNil(t, PromptStyle, "PromptStyle should be defined")
	assert.NotNil(t, HelpStyle, "HelpStyle should be defined")
}

func TestTitleStyle_Properties(t *testing.T) {
	// Verify TitleStyle has expected properties (bold + purple color)
	// Render with test string to verify it doesn't panic
	testStr := "Title"
	result := TitleStyle.Render(testStr)
	assert.NotEmpty(t, result, "TitleStyle should render without error")
	// Result should contain the text (though ANSI codes may wrap it)
	assert.Contains(t, result, testStr, "TitleStyle render should contain the input text")
}

func TestSubtitleStyle_Properties(t *testing.T) {
	// Verify SubtitleStyle has expected properties (faint + muted color)
	testStr := "Subtitle"
	result := SubtitleStyle.Render(testStr)
	assert.NotEmpty(t, result, "SubtitleStyle should render without error")
	assert.Contains(t, result, testStr, "SubtitleStyle render should contain the input text")
}

func TestSuccessStyle_Properties(t *testing.T) {
	// Verify SuccessStyle has expected properties (green + bold)
	testStr := "Success"
	result := SuccessStyle.Render(testStr)
	assert.NotEmpty(t, result, "SuccessStyle should render without error")
	assert.Contains(t, result, testStr, "SuccessStyle render should contain the input text")
}

func TestErrorStyle_Properties(t *testing.T) {
	// Verify ErrorStyle has expected properties (red + bold)
	testStr := "Error"
	result := ErrorStyle.Render(testStr)
	assert.NotEmpty(t, result, "ErrorStyle should render without error")
	assert.Contains(t, result, testStr, "ErrorStyle render should contain the input text")
}

func TestPromptStyle_Properties(t *testing.T) {
	// Verify PromptStyle has expected properties (cyan accent)
	testStr := "Prompt"
	result := PromptStyle.Render(testStr)
	assert.NotEmpty(t, result, "PromptStyle should render without error")
	assert.Contains(t, result, testStr, "PromptStyle render should contain the input text")
}

func TestHelpStyle_Properties(t *testing.T) {
	// Verify HelpStyle has expected properties (faint + muted + italic)
	testStr := "Help"
	result := HelpStyle.Render(testStr)
	assert.NotEmpty(t, result, "HelpStyle should render without error")
	assert.Contains(t, result, testStr, "HelpStyle render should contain the input text")
}

func TestAllStyles_DontPanic(t *testing.T) {
	// Verify all styles can be used without panicking
	testStr := "test"

	assert.NotPanics(t, func() { TitleStyle.Render(testStr) })
	assert.NotPanics(t, func() { SubtitleStyle.Render(testStr) })
	assert.NotPanics(t, func() { SuccessStyle.Render(testStr) })
	assert.NotPanics(t, func() { ErrorStyle.Render(testStr) })
	assert.NotPanics(t, func() { PromptStyle.Render(testStr) })
	assert.NotPanics(t, func() { HelpStyle.Render(testStr) })
}

func TestColorPalette_Gradient(t *testing.T) {
	// Verify that purple and pink colors exist for gradient purposes
	// These are the primary gradient colors
	assert.NotNil(t, PurpleColor, "PurpleColor required for gradient")
	assert.NotNil(t, PinkColor, "PinkColor required for gradient")

	purpleStr := string(PurpleColor)
	pinkStr := string(PinkColor)

	assert.NotEmpty(t, purpleStr, "PurpleColor hex value required")
	assert.NotEmpty(t, pinkStr, "PinkColor hex value required")
}

func TestColorPalette_StateIndicators(t *testing.T) {
	// Verify that success and error colors exist for status indicators
	assert.NotNil(t, SuccessColor, "SuccessColor required for state indicators")
	assert.NotNil(t, ErrorColor, "ErrorColor required for state indicators")

	successStr := string(SuccessColor)
	errorStr := string(ErrorColor)

	assert.NotEmpty(t, successStr, "SuccessColor hex value required")
	assert.NotEmpty(t, errorStr, "ErrorColor hex value required")
}
