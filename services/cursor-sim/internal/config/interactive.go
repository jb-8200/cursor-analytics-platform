package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// GenerationParams holds the parsed generation parameters.
type GenerationParams struct {
	Developers int
	Days       int
	MaxCommits int
}

// PromptConfig holds configuration for interactive prompts.
type PromptConfig struct {
	reader     *bufio.Reader
	writer     io.Writer
	maxRetries int
}

// NewPromptConfig creates a new PromptConfig with default settings.
// Uses stdin/stdout for I/O and sets maxRetries to 3.
func NewPromptConfig() *PromptConfig {
	return &PromptConfig{
		reader:     bufio.NewReader(os.Stdin),
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
	attempts := 0

	for attempts <= p.maxRetries {
		// Display prompt with default value
		fmt.Fprintf(p.writer, "%s (default: %d): ", prompt, defaultVal)

		// Read user input
		line, err := p.reader.ReadString('\n')
		if err != nil {
			// EOF or error - use default
			if err.Error() == "EOF" {
				return defaultVal, nil
			}
			return defaultVal, err
		}

		input := strings.TrimSpace(line)

		// Empty input - use default
		if input == "" {
			return defaultVal, nil
		}

		// Parse integer
		value, parseErr := strconv.Atoi(input)
		if parseErr != nil {
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

// InteractiveConfig prompts the user for all configuration values.
// Displays prompts for number of developers, time period (in months), and max commits.
// Converts months to days (months * 30) and returns a validated GenerationParams struct.
//
// Default values:
//   - Developers: 10 (range: 1-100)
//   - Period: 6 months (range: 1-24 months)
//   - MaxCommits: 500 (range: 100-2000)
//
// Returns the configured parameters or an error if input reading fails.
func (p *PromptConfig) InteractiveConfig() (*GenerationParams, error) {
	fmt.Fprintln(p.writer, "\nCursor Simulator - Interactive Configuration")
	fmt.Fprintln(p.writer, "Press Enter to use default values")

	// Prompt 1: Number of developers
	developers, err := p.PromptForInt("Number of developers", 10, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to read developers: %w", err)
	}

	// Prompt 2: Period in months
	months, err := p.PromptForInt("Period in months", 6, 1, 24)
	if err != nil {
		return nil, fmt.Errorf("failed to read period: %w", err)
	}

	// Prompt 3: Maximum commits per developer
	maxCommits, err := p.PromptForInt("Maximum commits per developer", 500, 100, 2000)
	if err != nil {
		return nil, fmt.Errorf("failed to read max commits: %w", err)
	}

	// Convert months to days
	days := months * 30

	// Create params
	params := &GenerationParams{
		Developers: developers,
		Days:       days,
		MaxCommits: maxCommits,
	}

	// Display configuration summary
	p.displayConfigSummary(params, months)

	return params, nil
}

// displayConfigSummary prints a formatted summary of the configuration parameters.
func (p *PromptConfig) displayConfigSummary(params *GenerationParams, months int) {
	fmt.Fprintln(p.writer, "\nConfiguration Summary:")
	fmt.Fprintf(p.writer, "  Developers: %d\n", params.Developers)
	fmt.Fprintf(p.writer, "  Period: %d months (%d days)\n", months, params.Days)
	fmt.Fprintf(p.writer, "  Max commits: %d\n\n", params.MaxCommits)
}
