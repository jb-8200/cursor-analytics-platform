package config

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromptForInt_ValidInput(t *testing.T) {
	input := strings.NewReader("5\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 3,
	}

	result, err := cfg.PromptForInt("Enter number:", 10, 1, 100)
	require.NoError(t, err)
	assert.Equal(t, 5, result)
	assert.Contains(t, output.String(), "Enter number:")
}

func TestPromptForInt_DefaultOnEmpty(t *testing.T) {
	input := strings.NewReader("\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 3,
	}

	result, err := cfg.PromptForInt("Enter number:", 42, 1, 100)
	require.NoError(t, err)
	assert.Equal(t, 42, result)
	assert.Contains(t, output.String(), "default: 42")
}

func TestPromptForInt_InvalidInputRetry(t *testing.T) {
	// First input is invalid (non-numeric), second is valid
	input := strings.NewReader("abc\n5\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 3,
	}

	result, err := cfg.PromptForInt("Enter number:", 10, 1, 100)
	require.NoError(t, err)
	assert.Equal(t, 5, result)
	assert.Contains(t, output.String(), "Invalid input")
}

func TestPromptForInt_OutOfRangeRetry(t *testing.T) {
	// First input is out of range (too high), second is valid
	input := strings.NewReader("150\n50\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 3,
	}

	result, err := cfg.PromptForInt("Enter number:", 10, 1, 100)
	require.NoError(t, err)
	assert.Equal(t, 50, result)
	assert.Contains(t, output.String(), "out of range")
}

func TestPromptForInt_OutOfRangeLowRetry(t *testing.T) {
	// First input is out of range (too low), second is valid
	input := strings.NewReader("0\n25\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 3,
	}

	result, err := cfg.PromptForInt("Enter number:", 10, 1, 100)
	require.NoError(t, err)
	assert.Equal(t, 25, result)
	assert.Contains(t, output.String(), "out of range")
}

func TestPromptForInt_MaxRetriesExceeded(t *testing.T) {
	// All inputs are invalid
	input := strings.NewReader("abc\nxyz\n999\nfoo\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 3,
	}

	result, err := cfg.PromptForInt("Enter number:", 42, 1, 100)
	require.NoError(t, err)
	assert.Equal(t, 42, result) // Should return default after max retries
	assert.Contains(t, output.String(), "Using default")
}

func TestPromptForInt_NegativeNumbers(t *testing.T) {
	input := strings.NewReader("-5\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 3,
	}

	// Test range that includes negative numbers
	result, err := cfg.PromptForInt("Enter number:", 0, -10, 10)
	require.NoError(t, err)
	assert.Equal(t, -5, result)
}

func TestPromptForInt_ZeroValue(t *testing.T) {
	input := strings.NewReader("0\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 3,
	}

	// Test that 0 is valid when in range
	result, err := cfg.PromptForInt("Enter number:", 5, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 0, result)
}

func TestPromptForInt_MultipleInvalidInputs(t *testing.T) {
	// Multiple invalid inputs followed by a valid one
	input := strings.NewReader("abc\n999\nxyz\n42\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 5,
	}

	result, err := cfg.PromptForInt("Enter number:", 10, 1, 100)
	require.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestNewPromptConfig_Defaults(t *testing.T) {
	cfg := NewPromptConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, 3, cfg.maxRetries)
	assert.NotNil(t, cfg.reader)
	assert.NotNil(t, cfg.writer)
}

// TASK-CLI-02: InteractiveConfig Function Tests

func TestInteractiveConfig_AllDefaults(t *testing.T) {
	// Simulate pressing Enter 3 times (use all defaults)
	input := strings.NewReader("\n\n\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 3,
	}

	params, err := cfg.InteractiveConfig()
	require.NoError(t, err)
	assert.NotNil(t, params)

	// Verify defaults: 10 developers, 6 months (180 days), 500 max commits
	assert.Equal(t, 10, params.Developers)
	assert.Equal(t, 180, params.Days) // 6 months * 30 days
	assert.Equal(t, 500, params.MaxCommits)

	// Verify output contains prompts
	outputStr := output.String()
	assert.Contains(t, outputStr, "Number of developers")
	assert.Contains(t, outputStr, "Period in months")
	assert.Contains(t, outputStr, "Maximum commits")
}

func TestInteractiveConfig_CustomValues(t *testing.T) {
	// Simulate custom inputs: 5 developers, 6 months, 1500 max commits
	input := strings.NewReader("5\n6\n1500\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 3,
	}

	params, err := cfg.InteractiveConfig()
	require.NoError(t, err)
	assert.NotNil(t, params)

	assert.Equal(t, 5, params.Developers)
	assert.Equal(t, 180, params.Days) // 6 months * 30 days
	assert.Equal(t, 1500, params.MaxCommits)
}

func TestInteractiveConfig_MonthsToDays(t *testing.T) {
	// Test month-to-day conversion
	testCases := []struct {
		name           string
		input          string
		expectedDays   int
		expectedMonths int
	}{
		{"1 month", "\n1\n\n", 30, 1},
		{"3 months", "\n3\n\n", 90, 3},
		{"6 months", "\n6\n\n", 180, 6},
		{"12 months", "\n12\n\n", 360, 12},
		{"24 months", "\n24\n\n", 720, 24},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			output := &bytes.Buffer{}

			cfg := &PromptConfig{
				reader:     bufio.NewReader(input),
				writer:     output,
				maxRetries: 3,
			}

			params, err := cfg.InteractiveConfig()
			require.NoError(t, err)
			assert.Equal(t, tc.expectedDays, params.Days)

			// Verify summary shows both months and days
			outputStr := output.String()
			assert.Contains(t, outputStr, fmt.Sprintf("%d months (%d days)", tc.expectedMonths, tc.expectedDays))
		})
	}
}

func TestInteractiveConfig_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name            string
		input           string
		expectedDevs    int
		expectedDays    int
		expectedCommits int
	}{
		{
			name:            "Invalid then valid developers",
			input:           "abc\n5\n\n\n", // Invalid, then 5, then defaults for rest
			expectedDevs:    5,
			expectedDays:    180, // default 6 months
			expectedCommits: 500, // default
		},
		{
			name:            "Out of range developers",
			input:           "200\n10\n\n\n", // Too high, then 10, then defaults
			expectedDevs:    10,
			expectedDays:    180,
			expectedCommits: 500,
		},
		{
			name:            "Invalid months",
			input:           "\nxyz\n3\n\n", // Default devs, invalid months, then 3, then default commits
			expectedDevs:    10,
			expectedDays:    90, // 3 months
			expectedCommits: 500,
		},
		{
			name:            "Out of range months",
			input:           "\n50\n6\n\n", // Default devs, too high months, then 6, then default commits
			expectedDevs:    10,
			expectedDays:    180, // 6 months
			expectedCommits: 500,
		},
		{
			name:            "Invalid max commits",
			input:           "\n\nabc\n1000\n", // Defaults for devs and months, invalid commits, then 1000
			expectedDevs:    10,
			expectedDays:    180,
			expectedCommits: 1000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			output := &bytes.Buffer{}

			cfg := &PromptConfig{
				reader:     bufio.NewReader(input),
				writer:     output,
				maxRetries: 3,
			}

			params, err := cfg.InteractiveConfig()
			require.NoError(t, err)
			assert.Equal(t, tc.expectedDevs, params.Developers)
			assert.Equal(t, tc.expectedDays, params.Days)
			assert.Equal(t, tc.expectedCommits, params.MaxCommits)
		})
	}
}

func TestInteractiveConfig_MixedDefaultsAndCustom(t *testing.T) {
	// Use default for developers, custom for months and commits
	input := strings.NewReader("\n12\n1500\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 3,
	}

	params, err := cfg.InteractiveConfig()
	require.NoError(t, err)

	assert.Equal(t, 10, params.Developers)   // default
	assert.Equal(t, 360, params.Days)        // 12 months * 30
	assert.Equal(t, 1500, params.MaxCommits) // custom
}

func TestInteractiveConfig_DisplaysSummary(t *testing.T) {
	input := strings.NewReader("\n\n\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     bufio.NewReader(input),
		writer:     output,
		maxRetries: 3,
	}

	_, err := cfg.InteractiveConfig()
	require.NoError(t, err)

	outputStr := output.String()

	// Verify configuration summary is displayed
	assert.Contains(t, outputStr, "Configuration Summary:")
	assert.Contains(t, outputStr, "Developers:")
	assert.Contains(t, outputStr, "Period:")
	assert.Contains(t, outputStr, "Max commits:")
}
