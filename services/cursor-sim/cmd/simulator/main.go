// Package main provides the cursor-sim CLI application.
// cursor-sim v2 is a seed-based Cursor API simulator that generates
// synthetic usage data matching the exact Cursor Business API.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/preview"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/server"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/tui"
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

	// Display DOXAPI banner for runtime and interactive modes (skip preview and help)
	if cfg.Mode == "runtime" || cfg.Interactive {
		tui.DisplayBanner(Version)
	}

	// TASK-CLI-10: If interactive mode is enabled, prompt for configuration
	if cfg.Interactive {
		promptConfig := config.NewPromptConfig()
		params, err := promptConfig.InteractiveConfig()
		if err != nil {
			log.Fatalf("Interactive config error: %v\n", err)
		}
		// Override config with interactive parameters
		cfg.GenParams = *params
		cfg.Days = params.Days
		log.Printf("Using interactive configuration: %d developers, %d days, max %d commits\n",
			params.Developers, params.Days, params.MaxCommits)
	}

	// Run the application
	if err := run(ctx, cfg); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

// run is the main application logic, separated for testability.
func run(ctx context.Context, cfg *config.Config) error {

	log.Printf("Starting cursor-sim v%s with config: %s\n", Version, cfg)

	// Validate mode
	switch cfg.Mode {
	case "runtime":
		return runRuntimeMode(ctx, cfg)
	case "preview":
		return runPreviewMode(ctx, cfg)
	default:
		return fmt.Errorf("invalid mode: '%s' (must be 'runtime' or 'preview')", cfg.Mode)
	}
}

// runPreviewMode executes preview mode for quick seed validation.
func runPreviewMode(ctx context.Context, cfg *config.Config) error {
	log.Printf("Preview mode: validating seed file %s\n", cfg.SeedPath)

	// Load seed data
	seedData, err := seed.LoadSeed(cfg.SeedPath)
	if err != nil {
		return fmt.Errorf("failed to load seed data: %w", err)
	}

	// Create preview config
	previewCfg := preview.Config{
		Days:       cfg.Days,
		MaxCommits: 10,  // Limit for fast preview
		MaxEvents:  100, // Limit for fast preview
	}

	// Create and run preview
	p := preview.New(seedData, previewCfg, os.Stdout)
	if err := p.Run(ctx); err != nil {
		return fmt.Errorf("preview failed: %w", err)
	}

	log.Println("Preview complete")
	return nil
}

// runRuntimeMode executes full runtime mode with API server.
func runRuntimeMode(ctx context.Context, cfg *config.Config) error {
	// Load seed data
	log.Printf("Loading seed data from %s...\n", cfg.SeedPath)

	// TASK-CLI-06: Use LoadSeedWithReplication to support developer scaling
	// If cfg.GenParams.Developers > 0, replicate developers to that count
	// Otherwise, use all developers from seed file
	seedData, developers, err := seed.LoadSeedWithReplication(cfg.SeedPath, cfg.GenParams.Developers, nil)
	if err != nil {
		return fmt.Errorf("failed to load seed data: %w", err)
	}

	if cfg.GenParams.Developers > 0 {
		log.Printf("Loaded %d developers from seed file, replicated to %d developers\n",
			len(seedData.Developers), len(developers))
	} else {
		log.Printf("Loaded %d developers from seed file\n", len(developers))
	}

	// Initialize storage
	store := storage.NewMemoryStore()

	// Load developers into storage (use replicated developers, not seed.Developers)
	if err := store.LoadDevelopers(developers); err != nil {
		return fmt.Errorf("failed to load developers into storage: %w", err)
	}
	log.Printf("Loaded %d developers into storage\n", len(developers))

	// Generate commits
	log.Printf("Generating %d days of commit history...\n", cfg.Days)
	gen := generator.NewCommitGenerator(seedData, store, cfg.Velocity)
	maxCommits := cfg.GenParams.MaxCommits
	if maxCommits > 0 {
		log.Printf("Max commits limit: %d\n", maxCommits)
	}
	if err := gen.GenerateCommits(ctx, cfg.Days, maxCommits); err != nil {
		return fmt.Errorf("failed to generate commits: %w", err)
	}

	// Count generated commits for logging
	allCommits := store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour))
	actualDevelopers := store.ListDevelopers()
	log.Printf("Generated %d commits across %d developers\n", len(allCommits), len(actualDevelopers))

	// Generate PRs from commits
	log.Printf("Generating PRs from commits...\n")
	prGen := generator.NewPRGeneratorWithSeed(seedData, store, time.Now().UnixNano())
	startDate := time.Now().AddDate(0, 0, -cfg.Days)
	endDate := time.Now().Add(24 * time.Hour)
	if err := prGen.GeneratePRsFromCommits(startDate, endDate); err != nil {
		return fmt.Errorf("failed to generate PRs: %w", err)
	}
	repos := store.ListRepositories()
	log.Printf("Generated PRs across %d repositories\n", len(repos))

	// Generate reviews for PRs
	log.Printf("Generating reviews for PRs...\n")
	reviewGen := generator.NewReviewGenerator(seedData, rand.New(rand.NewSource(time.Now().UnixNano())))
	totalReviews := 0
	for _, repoName := range repos {
		prs := store.GetPRsByRepo(repoName)
		for _, pr := range prs {
			reviews := reviewGen.GenerateReviewsForPR(pr)
			for _, review := range reviews {
				if err := store.StoreReview(review); err != nil {
					log.Printf("Warning: failed to store review: %v", err)
				}
				totalReviews++
			}
		}
	}
	log.Printf("Generated %d reviews\n", totalReviews)

	// Generate issues for PRs
	log.Printf("Generating issues for PRs...\n")
	issueGen := generator.NewIssueGeneratorWithStore(seedData, store, time.Now().UnixNano())
	totalIssues := 0
	for _, repoName := range repos {
		prs := store.GetPRsByRepo(repoName)
		count, err := issueGen.GenerateAndStoreIssuesForPRs(prs, repoName)
		if err != nil {
			log.Printf("Warning: failed to generate issues for %s: %v", repoName, err)
		}
		totalIssues += count
	}
	log.Printf("Generated %d issues\n", totalIssues)

	// Generate model usage events
	log.Printf("Generating model usage events...\n")
	modelGen := generator.NewModelGenerator(seedData, store, cfg.Velocity)
	if err := modelGen.GenerateModelUsage(ctx, cfg.Days); err != nil {
		return fmt.Errorf("failed to generate model usage: %w", err)
	}
	log.Printf("Generated model usage events\n")

	// Generate client version events
	log.Printf("Generating client version events...\n")
	versionGen := generator.NewVersionGenerator(seedData, store, cfg.Velocity)
	if err := versionGen.GenerateClientVersions(ctx, cfg.Days); err != nil {
		return fmt.Errorf("failed to generate client versions: %w", err)
	}
	log.Printf("Generated client version events\n")

	// Generate file extension events
	log.Printf("Generating file extension events...\n")
	extensionGen := generator.NewExtensionGenerator(seedData, store, cfg.Velocity)
	if err := extensionGen.GenerateFileExtensions(ctx, cfg.Days); err != nil {
		return fmt.Errorf("failed to generate file extensions: %w", err)
	}
	log.Printf("Generated file extension events\n")

	// Generate feature events (MCP, Commands, Plans, AskMode)
	log.Printf("Generating feature events...\n")
	featureGen := generator.NewFeatureGenerator(seedData, store, cfg.Velocity)
	if err := featureGen.GenerateFeatures(ctx, cfg.Days); err != nil {
		return fmt.Errorf("failed to generate features: %w", err)
	}
	log.Printf("Generated feature events\n")

	// Create HTTP router
	router := server.NewRouter(store, seedData, DefaultAPIKey, cfg, Version)

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
