package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Default configuration values
const (
	DefaultPort        = 8080
	DefaultDevelopers  = 50
	DefaultVelocity    = "medium"
	DefaultFluctuation = 0.2
	DefaultSeed        = int64(12345)

	MinPort        = 1024
	MaxPort        = 65535
	MinDevelopers  = 1
	MaxDevelopers  = 10000
	MinFluctuation = 0.0
	MaxFluctuation = 1.0
)

// ValidVelocities contains allowed velocity values
var ValidVelocities = []string{"low", "medium", "high"}

// Config represents the simulator configuration
type Config struct {
	// Port is the HTTP server port
	Port int

	// Developers is the number of developers to simulate
	Developers int

	// Velocity controls event generation rate (low, medium, high)
	Velocity string

	// Fluctuation controls per-developer rate variance (0.0-1.0)
	Fluctuation float64

	// Seed for reproducible random generation
	Seed int64
}

// NewConfig returns a Config with default values
func NewConfig() *Config {
	return &Config{
		Port:        DefaultPort,
		Developers:  DefaultDevelopers,
		Velocity:    DefaultVelocity,
		Fluctuation: DefaultFluctuation,
		Seed:        DefaultSeed,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate port
	if c.Port < MinPort || c.Port > MaxPort {
		return fmt.Errorf("port must be between %d and %d, got %d", MinPort, MaxPort, c.Port)
	}

	// Validate developers
	if c.Developers < MinDevelopers || c.Developers > MaxDevelopers {
		return fmt.Errorf("developers must be between %d and %d, got %d", MinDevelopers, MaxDevelopers, c.Developers)
	}

	// Validate velocity
	if !isValidVelocity(c.Velocity) {
		return fmt.Errorf("velocity must be one of %v, got %q", ValidVelocities, c.Velocity)
	}

	// Validate fluctuation
	if c.Fluctuation < MinFluctuation || c.Fluctuation > MaxFluctuation {
		return fmt.Errorf("fluctuation must be between %.1f and %.1f, got %.2f", MinFluctuation, MaxFluctuation, c.Fluctuation)
	}

	return nil
}

// isValidVelocity checks if the given velocity is valid
func isValidVelocity(v string) bool {
	for _, valid := range ValidVelocities {
		if v == valid {
			return true
		}
	}
	return false
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ErrInvalidConfig is returned when configuration is invalid
var ErrInvalidConfig = errors.New("invalid configuration")

// ParseFlags parses command-line flags and returns a validated Config
func ParseFlags() (*Config, error) {
	cfg := &Config{}

	// Define flags
	flag.IntVar(&cfg.Port, "port", DefaultPort, "HTTP server port (1024-65535)")
	flag.IntVar(&cfg.Developers, "developers", DefaultDevelopers, "Number of developers to simulate (1-10000)")
	flag.StringVar(&cfg.Velocity, "velocity", DefaultVelocity, "Event generation rate: low|medium|high")
	flag.Float64Var(&cfg.Fluctuation, "fluctuation", DefaultFluctuation, "Per-developer rate variance (0.0-1.0)")
	flag.Int64Var(&cfg.Seed, "seed", DefaultSeed, "Random seed for reproducibility")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Cursor API Simulator - Generate synthetic Cursor usage data\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s --port 8080 --developers 100 --velocity high\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --fluctuation 0.3 --seed 42\n", os.Args[0])
	}

	// Parse flags
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}
