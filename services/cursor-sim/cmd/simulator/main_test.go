package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_Success(t *testing.T) {
	// Create configuration using test seed file
	cfg := &config.Config{
		Mode:     "runtime",
		SeedPath: "../../testdata/valid_seed.json",
		Port:     18080,
		Days:     1,
		Velocity: "medium",
	}

	// Run the application with a context that will cancel after 100ms
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Run in a goroutine
	done := make(chan error, 1)
	go func() {
		done <- run(ctx, cfg)
	}()

	// Wait for server to start
	time.Sleep(50 * time.Millisecond)

	// Verify server is responding
	resp, err := http.Get("http://localhost:18080/health")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	// Wait for context to cancel
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("run() did not exit after context cancellation")
	}
}

func TestRun_InvalidMode(t *testing.T) {
	ctx := context.Background()

	cfg := &config.Config{
		Mode:     "replay",
		SeedPath: "test.json",
		Port:     8080,
	}

	// Replay mode not supported
	err := run(ctx, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only runtime mode is supported")
}

func TestRun_InvalidSeedFile(t *testing.T) {
	ctx := context.Background()

	cfg := &config.Config{
		Mode:     "runtime",
		SeedPath: "/nonexistent/seed.json",
		Port:     8080,
		Days:     1,
		Velocity: "medium",
	}

	// Non-existent seed file should fail
	err := run(ctx, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load seed data")
}

func TestRun_WithMaxCommits(t *testing.T) {
	// Test that maxCommits parameter flows through correctly
	cfg := &config.Config{
		Mode:     "runtime",
		SeedPath: "../../testdata/valid_seed.json",
		Port:     18081, // Different port to avoid conflicts
		Days:     30,
		Velocity: "high",
		GenParams: config.GenerationParams{
			Developers: 0,        // Use all developers from seed
			Days:       30,        // 30 days
			MaxCommits: 20,        // Limit to 20 commits total
		},
	}

	// Run with context that cancels after 100ms (enough time to generate)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- run(ctx, cfg)
	}()

	// Wait for generation to complete
	time.Sleep(50 * time.Millisecond)

	// Verify server started
	resp, err := http.Get("http://localhost:18081/health")
	if err == nil {
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	}

	// Wait for completion
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("run() did not exit after context cancellation")
	}
}
