package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormModel_New(t *testing.T) {
	form := NewFormModel()

	require.NotNil(t, form)
	assert.Equal(t, 0, form.focusedField) // Should start at first field
	assert.Equal(t, 10, form.Developers)  // Default
	assert.Equal(t, 6, form.Months)       // Default
	assert.Equal(t, 500, form.MaxCommits) // Default
	assert.False(t, form.submitted)
	assert.False(t, form.cancelled)
}

func TestFormModel_FocusNavigation(t *testing.T) {
	form := NewFormModel()

	// Start at field 0 (developers)
	assert.Equal(t, 0, form.focusedField)

	// Move to next field
	form.NextField()
	assert.Equal(t, 1, form.focusedField)

	// Move to next field
	form.NextField()
	assert.Equal(t, 2, form.focusedField)

	// At last field, NextField should not go beyond
	form.NextField()
	assert.Equal(t, 2, form.focusedField)

	// Move back to previous field
	form.PrevField()
	assert.Equal(t, 1, form.focusedField)

	// Move to previous field
	form.PrevField()
	assert.Equal(t, 0, form.focusedField)

	// At first field, PrevField should not go below
	form.PrevField()
	assert.Equal(t, 0, form.focusedField)
}

func TestFormModel_InputDevelopers(t *testing.T) {
	form := NewFormModel()

	// Start at developers field
	form.focusedField = 0

	// Add digits
	form.AddChar('2')
	form.AddChar('5')

	// Value should be "25"
	assert.Equal(t, "25", form.developerInput)
	assert.Equal(t, 25, form.Developers)
}

func TestFormModel_InputMonths(t *testing.T) {
	form := NewFormModel()

	// Move to months field
	form.focusedField = 1

	form.AddChar('1')
	form.AddChar('2')

	assert.Equal(t, "12", form.monthInput)
	assert.Equal(t, 12, form.Months)
}

func TestFormModel_InputMaxCommits(t *testing.T) {
	form := NewFormModel()

	// Move to maxCommits field
	form.focusedField = 2

	form.AddChar('1')
	form.AddChar('5')
	form.AddChar('0')
	form.AddChar('0')

	assert.Equal(t, "1500", form.commitInput)
	assert.Equal(t, 1500, form.MaxCommits)
}

func TestFormModel_Backspace(t *testing.T) {
	form := NewFormModel()
	form.focusedField = 0

	form.AddChar('2')
	form.AddChar('5')
	assert.Equal(t, "25", form.developerInput)

	form.Backspace()
	assert.Equal(t, "2", form.developerInput)
	assert.Equal(t, 2, form.Developers)

	form.Backspace()
	assert.Equal(t, "", form.developerInput)
	assert.Equal(t, 10, form.Developers) // Back to default
}

func TestFormModel_ValidateDevelopers(t *testing.T) {
	form := NewFormModel()
	form.focusedField = 0

	// Test valid range
	form.developerInput = "50"
	form.Developers = 50
	assert.True(t, form.ValidateDevelopers())

	// Test minimum boundary
	form.developerInput = "1"
	form.Developers = 1
	assert.True(t, form.ValidateDevelopers())

	// Test maximum boundary
	form.developerInput = "100"
	form.Developers = 100
	assert.True(t, form.ValidateDevelopers())

	// Test below minimum
	form.developerInput = "0"
	form.Developers = 0
	assert.False(t, form.ValidateDevelopers())

	// Test above maximum
	form.developerInput = "101"
	form.Developers = 101
	assert.False(t, form.ValidateDevelopers())
}

func TestFormModel_ValidateMonths(t *testing.T) {
	form := NewFormModel()

	// Test valid range
	form.monthInput = "12"
	form.Months = 12
	assert.True(t, form.ValidateMonths())

	// Test boundaries
	form.monthInput = "1"
	form.Months = 1
	assert.True(t, form.ValidateMonths())

	form.monthInput = "24"
	form.Months = 24
	assert.True(t, form.ValidateMonths())

	// Test invalid
	form.monthInput = "0"
	form.Months = 0
	assert.False(t, form.ValidateMonths())

	form.monthInput = "25"
	form.Months = 25
	assert.False(t, form.ValidateMonths())
}

func TestFormModel_ValidateMaxCommits(t *testing.T) {
	form := NewFormModel()

	// Test valid range
	form.commitInput = "1000"
	form.MaxCommits = 1000
	assert.True(t, form.ValidateMaxCommits())

	// Test boundaries
	form.commitInput = "100"
	form.MaxCommits = 100
	assert.True(t, form.ValidateMaxCommits())

	form.commitInput = "2000"
	form.MaxCommits = 2000
	assert.True(t, form.ValidateMaxCommits())

	// Test invalid
	form.commitInput = "50"
	form.MaxCommits = 50
	assert.False(t, form.ValidateMaxCommits())

	form.commitInput = "2500"
	form.MaxCommits = 2500
	assert.False(t, form.ValidateMaxCommits())
}

func TestFormModel_ValidateAll(t *testing.T) {
	form := NewFormModel()

	// Valid state
	form.Developers = 10
	form.Months = 6
	form.MaxCommits = 500
	assert.True(t, form.ValidateAll())

	// Invalid developer
	form.Developers = 0
	assert.False(t, form.ValidateAll())
	form.Developers = 10

	// Invalid months
	form.Months = 0
	assert.False(t, form.ValidateAll())
	form.Months = 6

	// Invalid commits
	form.MaxCommits = 50
	assert.False(t, form.ValidateAll())
	form.MaxCommits = 500

	assert.True(t, form.ValidateAll())
}

func TestFormModel_Submit(t *testing.T) {
	form := NewFormModel()

	assert.False(t, form.submitted)

	form.Submit()

	assert.True(t, form.submitted)
	assert.False(t, form.cancelled)
}

func TestFormModel_Cancel(t *testing.T) {
	form := NewFormModel()

	assert.False(t, form.cancelled)

	form.Cancel()

	assert.True(t, form.cancelled)
	assert.False(t, form.submitted)
}

func TestFormModel_GetError(t *testing.T) {
	form := NewFormModel()

	// No error initially
	assert.Equal(t, "", form.GetError())

	// Invalid developer
	form.Developers = 0
	form.focusedField = 0
	form.error = "Developers must be between 1 and 100"
	assert.NotEmpty(t, form.GetError())

	// Clear error
	form.error = ""
	assert.Equal(t, "", form.GetError())
}

func TestFormModel_SetError(t *testing.T) {
	form := NewFormModel()

	form.SetError("Test error")
	assert.Equal(t, "Test error", form.error)
}

func TestFormModel_IsFieldFocused(t *testing.T) {
	form := NewFormModel()

	form.focusedField = 0
	assert.True(t, form.IsFieldFocused(0))
	assert.False(t, form.IsFieldFocused(1))

	form.focusedField = 1
	assert.False(t, form.IsFieldFocused(0))
	assert.True(t, form.IsFieldFocused(1))
}

func TestFormModel_AddCharNonNumeric(t *testing.T) {
	form := NewFormModel()
	form.focusedField = 0

	// Try to add non-numeric character
	form.AddChar('a')
	assert.Equal(t, "", form.developerInput)

	// Try to add space
	form.AddChar(' ')
	assert.Equal(t, "", form.developerInput)

	// Valid numeric character
	form.AddChar('5')
	assert.Equal(t, "5", form.developerInput)
}

func TestFormModel_MaxLengthInput(t *testing.T) {
	form := NewFormModel()
	form.focusedField = 0

	// Add maximum digits for developers (up to 100, so 3 chars)
	form.AddChar('1')
	form.AddChar('0')
	form.AddChar('0')

	assert.Equal(t, "100", form.developerInput)

	// Try to add more - should be limited
	form.AddChar('0')
	assert.Equal(t, "100", form.developerInput) // Should not add more
}

func TestFormModel_ClearField(t *testing.T) {
	form := NewFormModel()

	form.developerInput = "25"
	form.monthInput = "12"
	form.commitInput = "1000"

	form.ClearCurrentField()

	assert.Equal(t, "", form.developerInput)
	assert.Equal(t, "12", form.monthInput)
	assert.Equal(t, "1000", form.commitInput)
}

func TestFormModel_GetDays(t *testing.T) {
	form := NewFormModel()

	form.Months = 6
	assert.Equal(t, 180, form.GetDays()) // 6 * 30

	form.Months = 12
	assert.Equal(t, 360, form.GetDays()) // 12 * 30

	form.Months = 1
	assert.Equal(t, 30, form.GetDays()) // 1 * 30
}

func TestFormModel_GetSummary(t *testing.T) {
	form := NewFormModel()

	form.Developers = 20
	form.Months = 12
	form.MaxCommits = 1000

	summary := form.GetSummary()
	assert.Contains(t, summary, "20")
	assert.Contains(t, summary, "12")
	assert.Contains(t, summary, "360") // 12 * 30
	assert.Contains(t, summary, "1000")
}
