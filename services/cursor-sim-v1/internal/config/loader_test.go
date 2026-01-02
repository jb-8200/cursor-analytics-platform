package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFile_ValidJSON(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configJSON := `{
		"Port": 9090,
		"Developers": 75,
		"Velocity": "high",
		"Fluctuation": 0.3,
		"Seed": 54321
	}`

	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	cfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("Expected no error loading valid config, got: %v", err)
	}

	if cfg.Port != 9090 {
		t.Errorf("Expected port 9090, got %d", cfg.Port)
	}

	if cfg.Developers != 75 {
		t.Errorf("Expected developers 75, got %d", cfg.Developers)
	}

	if cfg.Velocity != "high" {
		t.Errorf("Expected velocity 'high', got %s", cfg.Velocity)
	}

	if cfg.Fluctuation != 0.3 {
		t.Errorf("Expected fluctuation 0.3, got %f", cfg.Fluctuation)
	}

	if cfg.Seed != 54321 {
		t.Errorf("Expected seed 54321, got %d", cfg.Seed)
	}
}

func TestLoadFromFile_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "bad-config.json")

	// Write invalid JSON
	badJSON := `{ "Port": "not-a-number" }`
	if err := os.WriteFile(configPath, []byte(badJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	_, err := LoadFromFile(configPath)
	if err == nil {
		t.Error("Expected error loading invalid JSON, got nil")
	}
}

func TestLoadFromFile_InvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid-config.json")

	// Valid JSON but invalid config values
	invalidJSON := `{
		"Port": 100,
		"Developers": 50,
		"Velocity": "medium",
		"Fluctuation": 0.2,
		"Seed": 12345
	}`

	if err := os.WriteFile(configPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	_, err := LoadFromFile(configPath)
	if err == nil {
		t.Error("Expected error for invalid port value, got nil")
	}
}

func TestLoadFromFile_FileNotFound(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/config.json")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestMergeWithFlags_NoFile(t *testing.T) {
	cfg, err := MergeWithFlags("")
	if err != nil {
		t.Fatalf("Expected no error with empty file path, got: %v", err)
	}

	// Should return defaults
	if cfg.Port != DefaultPort {
		t.Errorf("Expected default port %d, got %d", DefaultPort, cfg.Port)
	}
}

func TestMergeWithFlags_WithFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configJSON := `{
		"Port": 9999,
		"Developers": 200,
		"Velocity": "low",
		"Fluctuation": 0.8,
		"Seed": 11111
	}`

	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	cfg, err := MergeWithFlags(configPath)
	if err != nil {
		t.Fatalf("Expected no error loading file, got: %v", err)
	}

	if cfg.Port != 9999 {
		t.Errorf("Expected port 9999, got %d", cfg.Port)
	}
}
