package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate_RuntimeMode(t *testing.T) {
	cfg := &Config{
		Mode:     "runtime",
		SeedPath: "seed.json",
		Port:     8080,
		Days:     90,
		Velocity: "medium",
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestConfig_Validate_RuntimeMode_MissingSeedPath(t *testing.T) {
	cfg := &Config{
		Mode:     "runtime",
		SeedPath: "",
		Port:     8080,
	}

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "seed path is required for runtime mode")
}

func TestConfig_Validate_ReplayMode(t *testing.T) {
	cfg := &Config{
		Mode:       "replay",
		CorpusPath: "events.parquet",
		Port:       8080,
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestConfig_Validate_ReplayMode_MissingCorpusPath(t *testing.T) {
	cfg := &Config{
		Mode:       "replay",
		CorpusPath: "",
		Port:       8080,
	}

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "corpus path is required for replay mode")
}

func TestConfig_Validate_InvalidMode(t *testing.T) {
	cfg := &Config{
		Mode: "invalid",
		Port: 8080,
	}

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mode must be 'runtime' or 'replay'")
}

func TestConfig_Validate_InvalidPort(t *testing.T) {
	tests := []struct {
		name string
		port int
		err  bool
	}{
		{"valid 8080", 8080, false},
		{"valid 3000", 3000, false},
		{"invalid 0", 0, true},
		{"invalid negative", -1, true},
		{"invalid too large", 70000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Mode:     "runtime",
				SeedPath: "seed.json",
				Port:     tt.port,
				Days:     90,
				Velocity: "medium",
			}

			err := cfg.Validate()
			if tt.err {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "port")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_Validate_InvalidDays(t *testing.T) {
	tests := []struct {
		name string
		days int
		err  bool
	}{
		{"valid 1", 1, false},
		{"valid 90", 90, false},
		{"valid 365", 365, false},
		{"invalid 0", 0, true},
		{"invalid negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Mode:     "runtime",
				SeedPath: "seed.json",
				Port:     8080,
				Days:     tt.days,
				Velocity: "medium",
			}

			err := cfg.Validate()
			if tt.err {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "days")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_Validate_InvalidVelocity(t *testing.T) {
	tests := []struct {
		name     string
		velocity string
		err      bool
	}{
		{"valid low", "low", false},
		{"valid medium", "medium", false},
		{"valid high", "high", false},
		{"invalid fast", "fast", true},
		{"invalid empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Mode:     "runtime",
				SeedPath: "seed.json",
				Port:     8080,
				Days:     90,
				Velocity: tt.velocity,
			}

			err := cfg.Validate()
			if tt.err {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "velocity")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseFlags_Defaults(t *testing.T) {
	args := []string{"-seed=test.json"}

	cfg, err := parseFlagsWithArgs(args)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "runtime", cfg.Mode)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, 90, cfg.Days)
	assert.Equal(t, "medium", cfg.Velocity)
	assert.Equal(t, "test.json", cfg.SeedPath)
}

func TestParseFlags_CustomValues(t *testing.T) {
	args := []string{
		"-mode=replay",
		"-corpus=events.parquet",
		"-port=9000",
		"-days=30",
		"-velocity=high",
	}

	cfg, err := parseFlagsWithArgs(args)
	require.NoError(t, err)

	assert.Equal(t, "replay", cfg.Mode)
	assert.Equal(t, "events.parquet", cfg.CorpusPath)
	assert.Equal(t, 9000, cfg.Port)
	assert.Equal(t, 30, cfg.Days)
	assert.Equal(t, "high", cfg.Velocity)
}

func TestParseFlags_RuntimeMode(t *testing.T) {
	args := []string{
		"-mode=runtime",
		"-seed=test_seed.json",
		"-port=8080",
	}

	cfg, err := parseFlagsWithArgs(args)
	require.NoError(t, err)

	assert.Equal(t, "runtime", cfg.Mode)
	assert.Equal(t, "test_seed.json", cfg.SeedPath)
	assert.Equal(t, 8080, cfg.Port)
}

func TestParseFlags_EnvironmentOverrides(t *testing.T) {
	// Save original env and restore after test
	oldMode := os.Getenv("CURSOR_SIM_MODE")
	oldPort := os.Getenv("CURSOR_SIM_PORT")
	oldSeed := os.Getenv("CURSOR_SIM_SEED")
	defer func() {
		os.Setenv("CURSOR_SIM_MODE", oldMode)
		os.Setenv("CURSOR_SIM_PORT", oldPort)
		os.Setenv("CURSOR_SIM_SEED", oldSeed)
	}()

	os.Setenv("CURSOR_SIM_MODE", "runtime")
	os.Setenv("CURSOR_SIM_PORT", "9999")
	os.Setenv("CURSOR_SIM_SEED", "env_seed.json")

	args := []string{}
	cfg, err := parseFlagsWithArgs(args)
	require.NoError(t, err)

	assert.Equal(t, "runtime", cfg.Mode)
	assert.Equal(t, 9999, cfg.Port)
	assert.Equal(t, "env_seed.json", cfg.SeedPath)
}

func TestParseFlags_ValidationFails(t *testing.T) {
	args := []string{
		"-mode=runtime",
		// Missing -seed flag
		"-port=8080",
	}

	cfg, err := parseFlagsWithArgs(args)
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "seed path is required")
}

func TestParseFlags_InvalidFlag(t *testing.T) {
	args := []string{
		"-invalid-flag=value",
	}

	cfg, err := parseFlagsWithArgs(args)
	require.Error(t, err)
	assert.Nil(t, cfg)
}

func TestConfig_String(t *testing.T) {
	t.Run("runtime mode", func(t *testing.T) {
		cfg := &Config{
			Mode:     "runtime",
			SeedPath: "seed.json",
			Port:     8080,
			Days:     90,
			Velocity: "medium",
		}

		str := cfg.String()
		assert.Contains(t, str, "mode=runtime")
		assert.Contains(t, str, "port=8080")
		assert.Contains(t, str, "days=90")
		assert.Contains(t, str, "velocity=medium")
		assert.Contains(t, str, "seed=seed.json")
	})

	t.Run("replay mode", func(t *testing.T) {
		cfg := &Config{
			Mode:       "replay",
			CorpusPath: "events.parquet",
			Port:       9000,
		}

		str := cfg.String()
		assert.Contains(t, str, "mode=replay")
		assert.Contains(t, str, "port=9000")
		assert.Contains(t, str, "corpus=events.parquet")
		assert.NotContains(t, str, "days")
		assert.NotContains(t, str, "velocity")
	})
}
