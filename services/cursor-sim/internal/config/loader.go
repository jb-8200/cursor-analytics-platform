package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadFromFile reads configuration from a JSON file
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration in file: %w", err)
	}

	return &cfg, nil
}

// MergeWithFlags merges configuration from a file with command-line flags
// Flags override file values when explicitly set
func MergeWithFlags(filePath string) (*Config, error) {
	var cfg *Config
	var err error

	// Load from file if path is provided
	if filePath != "" {
		cfg, err = LoadFromFile(filePath)
		if err != nil {
			return nil, err
		}
	} else {
		cfg = NewConfig()
	}

	// Apply flag overrides (this would be called from main after flag.Parse)
	// For now, just return the loaded config
	return cfg, nil
}
