// Package main provides the cursor-sim CLI application.
// cursor-sim v2 is a seed-based Cursor API simulator that generates
// synthetic usage data matching the exact Cursor Business API.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/server"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// Version indicates the current release of cursor-sim.
const Version = "2.0.0"

// DefaultAPIKey is the default API key for authentication.
// In production, this would be loaded from secure configuration.
const DefaultAPIKey = "cursor-sim-dev-key"

func main() {
	fmt.Printf("cursor-sim v%s\n", Version)

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, shutting down gracefully...")
		cancel()
	}()

	// Parse configuration
	cfg, err := config.ParseFlags()
	if err != nil {
		log.Fatalf("Config error: %v\n", err)
	}

	// Run the application
	if err := run(ctx, cfg); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

// run is the main application logic, separated for testability.
func run(ctx context.Context, cfg *config.Config) error {

	log.Printf("Starting cursor-sim v%s with config: %s\n", Version, cfg)

	// Only runtime mode is supported in v2.0.0
	if cfg.Mode != "runtime" {
		return fmt.Errorf("only runtime mode is supported in v2.0.0")
	}

	// Load seed data
	log.Printf("Loading seed data from %s...\n", cfg.SeedPath)
	seedData, err := seed.LoadSeed(cfg.SeedPath)
	if err != nil {
		return fmt.Errorf("failed to load seed data: %w", err)
	}
	log.Printf("Loaded %d developers from seed file\n", len(seedData.Developers))

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers into storage
	if err := store.LoadDevelopers(seedData.Developers); err != nil {
		return fmt.Errorf("failed to load developers into storage: %w", err)
	}
	log.Printf("Loaded %d developers into storage\n", len(seedData.Developers))

	// Generate commits
	log.Printf("Generating %d days of commit history...\n", cfg.Days)
	gen := generator.NewCommitGenerator(seedData, store, cfg.Velocity)
	if err := gen.GenerateCommits(ctx, cfg.Days); err != nil {
		return fmt.Errorf("failed to generate commits: %w", err)
	}

	// Count generated commits for logging
	allCommits := store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour))
	log.Printf("Generated %d commits across %d developers\n", len(allCommits), len(seedData.Developers))

	// Generate model usage events
	log.Printf("Generating model usage events...\n")
	modelGen := generator.NewModelGenerator(seedData, store, cfg.Velocity)
	if err := modelGen.GenerateModelUsage(ctx, cfg.Days); err != nil {
		return fmt.Errorf("failed to generate model usage: %w", err)
	}
	log.Printf("Generated model usage events\n")

	// Create HTTP router
	router := server.NewRouter(store, seedData, DefaultAPIKey)

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	// Run server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("HTTP server listening on port %d\n", cfg.Port)
		log.Printf("API Key: %s\n", DefaultAPIKey)
		log.Printf("Health check: http://localhost:%d/health\n", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		log.Println("Shutting down HTTP server...")
		// Give the server 5 seconds to shut down gracefully
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}
		log.Println("Server stopped gracefully")
	}

	return nil
}
