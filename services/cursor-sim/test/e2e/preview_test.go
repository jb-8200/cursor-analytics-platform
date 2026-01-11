package e2e

import (
	"os"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TASK-PREV-10: Add E2E Test for Preview Mode

func TestE2E_PreviewMode(t *testing.T) {
	// Get absolute path to binary and service root
	wd, _ := os.Getwd()
	serviceRoot := filepath.Join(wd, "..", "..")
	binPath := filepath.Join(serviceRoot, "bin", "cursor-sim")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Skipf("binary not found at %s, run: go build -o bin/cursor-sim ./cmd/simulator", binPath)
	}

	cmd := exec.Command(
		binPath,
		"-mode", "preview",
		"-seed", "testdata/valid_seed.yaml",
	)
	cmd.Dir = serviceRoot  // Ensure command runs from service root

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Combined output:\n%s", string(output))
	}
	require.NoError(t, err, "preview mode should exit cleanly")

	outputStr := string(output)
	assert.Contains(t, outputStr, "PREVIEW MODE", "output should contain PREVIEW MODE header")
	assert.Contains(t, outputStr, "Developers", "output should contain Developers summary")
	assert.Contains(t, outputStr, "Sample Commits", "output should contain Sample Commits section")
	assert.Contains(t, outputStr, "Statistics", "output should contain Statistics section")
	assert.Contains(t, outputStr, "Validation Warnings", "output should contain Validation Warnings section")
}

func TestE2E_PreviewWithWarnings(t *testing.T) {
	// Get absolute path to binary and service root
	wd, _ := os.Getwd()
	serviceRoot := filepath.Join(wd, "..", "..")
	binPath := filepath.Join(serviceRoot, "bin", "cursor-sim")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Skipf("binary not found at %s, run: go build -o bin/cursor-sim ./cmd/simulator", binPath)
	}

	cmd := exec.Command(
		binPath,
		"-mode", "preview",
		"-seed", "testdata/invalid_models.yaml",
	)
	cmd.Dir = serviceRoot  // Ensure command runs from service root

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "preview mode should exit cleanly even with warnings")

	outputStr := string(output)
	assert.Contains(t, outputStr, "Validation Warnings", "output should contain Validation Warnings section")
	assert.Contains(t, outputStr, "⚠️", "output should contain warning emoji")
	assert.Contains(t, outputStr, "Unknown model", "output should mention unknown model")
}

func TestE2E_PreviewPerformance(t *testing.T) {
	// Get absolute path to binary and service root
	wd, _ := os.Getwd()
	serviceRoot := filepath.Join(wd, "..", "..")
	binPath := filepath.Join(serviceRoot, "bin", "cursor-sim")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Skipf("binary not found at %s, run: go build -o bin/cursor-sim ./cmd/simulator", binPath)
	}

	cmd := exec.Command(
		binPath,
		"-mode", "preview",
		"-seed", "testdata/large_seed.yaml",
	)
	cmd.Dir = serviceRoot  // Ensure command runs from service root

	// Should complete within 10 seconds
	done := make(chan error, 1)
	start := time.Now()

	go func() {
		output, err := cmd.CombinedOutput()
		if err == nil {
			outputStr := string(output)
			assert.Contains(t, outputStr, "PREVIEW MODE", "output should contain PREVIEW MODE")
		}
		done <- err
	}()

	select {
	case <-time.After(10 * time.Second):
		cmd.Process.Kill()
		t.Fatal("preview mode timed out after 10 seconds")
	case err := <-done:
		elapsed := time.Since(start)
		require.NoError(t, err, "preview should complete successfully")
		assert.Less(t, elapsed, 10*time.Second, "preview should complete within 10 seconds")
		t.Logf("preview completed in %v", elapsed)
	}
}
