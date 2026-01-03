package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// PromptConfig holds configuration for interactive prompts.
type PromptConfig struct {
	reader     io.Reader
	writer     io.Writer
	maxRetries int
}

// NewPromptConfig creates a new PromptConfig with default settings.
// Uses stdin/stdout for I/O and sets maxRetries to 3.
func NewPromptConfig() *PromptConfig {
	return &PromptConfig{
		reader:     os.Stdin,
		writer:     os.Stdout,
		maxRetries: 3,
	}
}

// PromptForInt prompts the user for an integer input with validation.
// Displays the prompt with the default value, validates the input is within [min, max],
// and retries on invalid input up to maxRetries times.
//
// Parameters:
//   - prompt: The prompt message to display
//   - defaultVal: The default value to use if user presses Enter or max retries exceeded
//   - min: Minimum valid value (inclusive)
//   - max: Maximum valid value (inclusive)
//
// Returns the validated integer or defaultVal if:
//   - User presses Enter without input
//   - Max retries exceeded
func (p *PromptConfig) PromptForInt(prompt string, defaultVal, min, max int) (int, error) {
	scanner := bufio.NewScanner(p.reader)
	attempts := 0

	for attempts <= p.maxRetries {
		// Display prompt with default value
		fmt.Fprintf(p.writer, "%s (default: %d): ", prompt, defaultVal)

		// Read user input
		if !scanner.Scan() {
			// EOF or error - use default
			if scanner.Err() != nil {
				return defaultVal, scanner.Err()
			}
			return defaultVal, nil
		}

		input := strings.TrimSpace(scanner.Text())

		// Empty input - use default
		if input == "" {
			return defaultVal, nil
		}

		// Parse integer
		value, err := strconv.Atoi(input)
		if err != nil {
			attempts++
			if attempts > p.maxRetries {
				fmt.Fprintf(p.writer, "Max retries exceeded. Using default value: %d\n", defaultVal)
				return defaultVal, nil
			}
			fmt.Fprintf(p.writer, "Invalid input: please enter a number.\n")
			continue
		}

		// Validate range
		if value < min || value > max {
			attempts++
			if attempts > p.maxRetries {
				fmt.Fprintf(p.writer, "Max retries exceeded. Using default value: %d\n", defaultVal)
				return defaultVal, nil
			}
			fmt.Fprintf(p.writer, "Value out of range: please enter a number between %d and %d.\n", min, max)
			continue
		}

		// Valid input
		return value, nil
	}

	// Should not reach here, but return default as safety
	return defaultVal, nil
}
