package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for cursor-sim v2.
type Config struct {
	// Mode is the operation mode: "runtime" or "replay"
	Mode string

	// SeedPath is the path to seed.json (required for runtime mode)
	SeedPath string

	// CorpusPath is the path to events.parquet (required for replay mode)
	CorpusPath string

	// Port is the HTTP server port
	Port int

	// Days is the number of days of history to generate (runtime mode only)
	Days int

	// Velocity controls event generation rate: "low", "medium", or "high"
	Velocity string
}

// ParseFlags parses command-line flags and environment variables to create a Config.
// Environment variables override flag defaults but are overridden by explicit flags.
// Returns an error if validation fails.
func ParseFlags() (*Config, error) {
	return parseFlagsWithArgs(os.Args[1:])
}

// parseFlagsWithArgs is an internal function that allows testing with custom arguments.
func parseFlagsWithArgs(args []string) (*Config, error) {
	cfg := &Config{}

	// Create a new FlagSet to avoid global flag conflicts in tests
	fs := flag.NewFlagSet("cursor-sim", flag.ContinueOnError)

	// Define flags with defaults
	fs.StringVar(&cfg.Mode, "mode", "runtime", "Operation mode: runtime or replay")
	fs.StringVar(&cfg.SeedPath, "seed", "", "Path to seed.json (required for runtime mode)")
	fs.StringVar(&cfg.CorpusPath, "corpus", "", "Path to events.parquet (required for replay mode)")
	fs.IntVar(&cfg.Port, "port", 8080, "HTTP server port")
	fs.IntVar(&cfg.Days, "days", 90, "Days of history to generate (runtime mode)")
	fs.StringVar(&cfg.Velocity, "velocity", "medium", "Event rate: low, medium, or high")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// Apply environment variable overrides
	if v := os.Getenv("CURSOR_SIM_MODE"); v != "" {
		cfg.Mode = v
	}
	if v := os.Getenv("CURSOR_SIM_SEED"); v != "" {
		cfg.SeedPath = v
	}
	if v := os.Getenv("CURSOR_SIM_CORPUS"); v != "" {
		cfg.CorpusPath = v
	}
	if v := os.Getenv("CURSOR_SIM_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Port = port
		}
	}
	if v := os.Getenv("CURSOR_SIM_DAYS"); v != "" {
		if days, err := strconv.Atoi(v); err == nil {
			cfg.Days = days
		}
	}
	if v := os.Getenv("CURSOR_SIM_VELOCITY"); v != "" {
		cfg.Velocity = v
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the configuration is valid.
// Returns a descriptive error if any validation rule fails.
func (c *Config) Validate() error {
	// Validate mode
	if c.Mode != "runtime" && c.Mode != "replay" {
		return fmt.Errorf("validation failed: mode must be 'runtime' or 'replay', got %q", c.Mode)
	}

	// Mode-specific validation
	if c.Mode == "runtime" {
		if c.SeedPath == "" {
			return fmt.Errorf("validation failed: seed path is required for runtime mode")
		}
	}

	if c.Mode == "replay" {
		if c.CorpusPath == "" {
			return fmt.Errorf("validation failed: corpus path is required for replay mode")
		}
	}

	// Validate port
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("validation failed: port must be between 1 and 65535, got %d", c.Port)
	}

	// Validate days (only relevant for runtime mode)
	if c.Mode == "runtime" && c.Days <= 0 {
		return fmt.Errorf("validation failed: days must be greater than 0, got %d", c.Days)
	}

	// Validate velocity (only required for runtime mode)
	if c.Mode == "runtime" {
		if c.Velocity != "low" && c.Velocity != "medium" && c.Velocity != "high" {
			return fmt.Errorf("validation failed: velocity must be 'low', 'medium', or 'high', got %q", c.Velocity)
		}
	}

	return nil
}

// String returns a human-readable string representation of the Config.
func (c *Config) String() string {
	if c.Mode == "runtime" {
		return fmt.Sprintf("mode=runtime seed=%s port=%d days=%d velocity=%s",
			c.SeedPath, c.Port, c.Days, c.Velocity)
	}
	return fmt.Sprintf("mode=replay corpus=%s port=%d",
		c.CorpusPath, c.Port)
}
