package config

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromptForInt_ValidInput(t *testing.T) {
	input := strings.NewReader("5\n")
	output := &bytes.Buffer{}

	cfg := &PromptConfig{
		reader:     input,
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
		reader:     input,
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
		reader:     input,
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
		reader:     input,
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
		reader:     input,
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
		reader:     input,
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
		reader:     input,
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
		reader:     input,
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
		reader:     input,
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
