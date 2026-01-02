package config

import (
	"flag"
	"os"
	"testing"
)

func TestConfig_DefaultValues(t *testing.T) {
	cfg := NewConfig()

	if cfg.Port != DefaultPort {
		t.Errorf("Expected default port %d, got %d", DefaultPort, cfg.Port)
	}

	if cfg.Developers != DefaultDevelopers {
		t.Errorf("Expected default developers %d, got %d", DefaultDevelopers, cfg.Developers)
	}

	if cfg.Velocity != DefaultVelocity {
		t.Errorf("Expected default velocity %s, got %s", DefaultVelocity, cfg.Velocity)
	}
}

func TestConfig_Validate_ValidConfig(t *testing.T) {
	cfg := &Config{
		Port:        8080,
		Developers:  50,
		Velocity:    "high",
		Fluctuation: 0.2,
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}
}

func TestConfig_Validate_InvalidPort(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"port too low", 100},
		{"port too high", 70000},
		{"port zero", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Port:        tt.port,
				Developers:  50,
				Velocity:    "high",
				Fluctuation: 0.2,
			}

			err := cfg.Validate()
			if err == nil {
				t.Errorf("Expected error for invalid port %d, got nil", tt.port)
			}
		})
	}
}

func TestConfig_Validate_InvalidDevelopers(t *testing.T) {
	tests := []struct {
		name       string
		developers int
	}{
		{"negative developers", -1},
		{"zero developers", 0},
		{"too many developers", 10001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Port:        8080,
				Developers:  tt.developers,
				Velocity:    "high",
				Fluctuation: 0.2,
			}

			err := cfg.Validate()
			if err == nil {
				t.Errorf("Expected error for invalid developers %d, got nil", tt.developers)
			}
		})
	}
}

func TestConfig_Validate_InvalidVelocity(t *testing.T) {
	tests := []struct {
		name     string
		velocity string
	}{
		{"invalid velocity", "super-fast"},
		{"empty velocity", ""},
		{"uppercase", "HIGH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Port:        8080,
				Developers:  50,
				Velocity:    tt.velocity,
				Fluctuation: 0.2,
			}

			err := cfg.Validate()
			if err == nil {
				t.Errorf("Expected error for invalid velocity %q, got nil", tt.velocity)
			}
		})
	}
}

func TestConfig_Validate_InvalidFluctuation(t *testing.T) {
	tests := []struct {
		name        string
		fluctuation float64
	}{
		{"negative fluctuation", -0.1},
		{"fluctuation too high", 1.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Port:        8080,
				Developers:  50,
				Velocity:    "high",
				Fluctuation: tt.fluctuation,
			}

			err := cfg.Validate()
			if err == nil {
				t.Errorf("Expected error for invalid fluctuation %f, got nil", tt.fluctuation)
			}
		})
	}
}

func TestParseFlags_Defaults(t *testing.T) {
	// Reset flag.CommandLine for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Simulate running with no flags
	os.Args = []string{"cmd"}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("Expected no error with default flags, got: %v", err)
	}

	if cfg.Port != DefaultPort {
		t.Errorf("Expected default port %d, got %d", DefaultPort, cfg.Port)
	}

	if cfg.Developers != DefaultDevelopers {
		t.Errorf("Expected default developers %d, got %d", DefaultDevelopers, cfg.Developers)
	}

	if cfg.Velocity != DefaultVelocity {
		t.Errorf("Expected default velocity %s, got %s", DefaultVelocity, cfg.Velocity)
	}

	if cfg.Fluctuation != DefaultFluctuation {
		t.Errorf("Expected default fluctuation %f, got %f", DefaultFluctuation, cfg.Fluctuation)
	}

	if cfg.Seed != DefaultSeed {
		t.Errorf("Expected default seed %d, got %d", DefaultSeed, cfg.Seed)
	}
}

func TestParseFlags_CustomValues(t *testing.T) {
	// Reset flag.CommandLine for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Simulate custom flags
	os.Args = []string{
		"cmd",
		"--port", "9000",
		"--developers", "100",
		"--velocity", "high",
		"--fluctuation", "0.5",
		"--seed", "99999",
	}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("Expected no error with valid custom flags, got: %v", err)
	}

	if cfg.Port != 9000 {
		t.Errorf("Expected port 9000, got %d", cfg.Port)
	}

	if cfg.Developers != 100 {
		t.Errorf("Expected developers 100, got %d", cfg.Developers)
	}

	if cfg.Velocity != "high" {
		t.Errorf("Expected velocity 'high', got %s", cfg.Velocity)
	}

	if cfg.Fluctuation != 0.5 {
		t.Errorf("Expected fluctuation 0.5, got %f", cfg.Fluctuation)
	}

	if cfg.Seed != 99999 {
		t.Errorf("Expected seed 99999, got %d", cfg.Seed)
	}
}

func TestParseFlags_InvalidPort(t *testing.T) {
	tests := []struct {
		name string
		port string
	}{
		{"port too low", "100"},
		{"port too high", "70000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag.CommandLine for testing
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			os.Args = []string{"cmd", "--port", tt.port}

			cfg, err := ParseFlags()
			if err == nil {
				t.Errorf("Expected error for invalid port %s, got nil", tt.port)
			}
			if cfg != nil {
				t.Errorf("Expected nil config on error, got %+v", cfg)
			}
		})
	}
}

func TestParseFlags_InvalidVelocity(t *testing.T) {
	// Reset flag.CommandLine for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	os.Args = []string{"cmd", "--velocity", "super-fast"}

	cfg, err := ParseFlags()
	if err == nil {
		t.Errorf("Expected error for invalid velocity, got nil")
	}
	if cfg != nil {
		t.Errorf("Expected nil config on error, got %+v", cfg)
	}
}

func TestParseFlags_InvalidFluctuation(t *testing.T) {
	tests := []struct {
		name        string
		fluctuation string
	}{
		{"negative", "-0.5"},
		{"too high", "2.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag.CommandLine for testing
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			os.Args = []string{"cmd", "--fluctuation", tt.fluctuation}

			cfg, err := ParseFlags()
			if err == nil {
				t.Errorf("Expected error for invalid fluctuation %s, got nil", tt.fluctuation)
			}
			if cfg != nil {
				t.Errorf("Expected nil config on error, got %+v", cfg)
			}
		})
	}
}
