// Package preview implements seed file validation and preview mode.
// Preview mode allows quick validation of seed files without full data generation.
// Inspired by NVIDIA NeMo DataDesigner preview patterns.
package preview

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// Config defines preview mode configuration parameters.
type Config struct {
	Days       int // Number of days of history to preview
	MaxCommits int // Maximum number of sample commits to generate
	MaxEvents  int // Maximum number of sample events to generate
}

// Preview handles seed file validation and preview output.
type Preview struct {
	seedData *seed.SeedData
	config   Config
	writer   io.Writer
}

// New creates a new Preview instance.
// If writer is nil, defaults to os.Stdout.
func New(seedData *seed.SeedData, config Config, writer io.Writer) *Preview {
	if writer == nil {
		writer = os.Stdout
	}

	return &Preview{
		seedData: seedData,
		config:   config,
		writer:   writer,
	}
}

// Run executes the preview mode, displaying seed validation and sample data.
func (p *Preview) Run(ctx context.Context) error {
	// Display header
	p.displayHeader()

	// Display developer summary
	p.displayDeveloperSummary()

	// Generate sample commits (limited)
	store, err := p.generateSampleData(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate sample data: %w", err)
	}

	// Display sample commits
	p.displaySampleCommits(store)

	// Display statistics
	p.displayStatistics(store)

	return nil
}

// displayHeader prints the preview mode header.
func (p *Preview) displayHeader() {
	fmt.Fprintln(p.writer, strings.Repeat("=", 60))
	fmt.Fprintln(p.writer, "  PREVIEW MODE - Seed Validation")
	fmt.Fprintln(p.writer, strings.Repeat("=", 60))
	fmt.Fprintln(p.writer)
}

// displayDeveloperSummary prints summary of developers from seed.
func (p *Preview) displayDeveloperSummary() {
	fmt.Fprintf(p.writer, "Developers: %d\n", len(p.seedData.Developers))
	fmt.Fprintln(p.writer, strings.Repeat("-", 60))

	if len(p.seedData.Developers) == 0 {
		fmt.Fprintln(p.writer, "  No developers configured (0 developers)")
		fmt.Fprintln(p.writer)
		return
	}

	for i, dev := range p.seedData.Developers {
		if i >= 5 { // Limit to first 5 developers
			fmt.Fprintf(p.writer, "  ... and %d more\n", len(p.seedData.Developers)-5)
			break
		}
		fmt.Fprintf(p.writer, "  • %s (%s)\n", dev.Name, dev.Email)
		if dev.UserID != "" {
			fmt.Fprintf(p.writer, "    ID: %s", dev.UserID)
			if dev.Seniority != "" {
				fmt.Fprintf(p.writer, " | Seniority: %s", dev.Seniority)
			}
			fmt.Fprintln(p.writer)
		}
	}
	fmt.Fprintln(p.writer)
}

// generateSampleData creates limited sample commits for preview.
func (p *Preview) generateSampleData(ctx context.Context) (*storage.MemoryStore, error) {
	store := storage.NewMemoryStore()

	// Load developers into storage
	if err := store.LoadDevelopers(p.seedData.Developers); err != nil {
		return nil, fmt.Errorf("failed to load developers: %w", err)
	}

	// Generate limited commits (much faster than full generation)
	gen := generator.NewCommitGenerator(p.seedData, store, "low")
	days := p.config.Days
	if days > 1 {
		days = 1 // Only generate 1 day for preview
	}

	if err := gen.GenerateCommits(ctx, days, p.config.MaxCommits); err != nil {
		return nil, fmt.Errorf("failed to generate commits: %w", err)
	}

	return store, nil
}

// displaySampleCommits shows a few sample commits.
func (p *Preview) displaySampleCommits(store *storage.MemoryStore) {
	commits := store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour))

	fmt.Fprintln(p.writer, "Sample Commits")
	fmt.Fprintln(p.writer, strings.Repeat("-", 60))

	if len(commits) == 0 {
		fmt.Fprintln(p.writer, "  No commits generated")
		fmt.Fprintln(p.writer)
		return
	}

	// Show first 3 commits
	showCount := 3
	if len(commits) < showCount {
		showCount = len(commits)
	}

	for i := 0; i < showCount; i++ {
		commit := commits[i]
		fmt.Fprintf(p.writer, "  [%s] %s\n", commit.CommitTs.Format("15:04:05"), commit.Message)
		fmt.Fprintf(p.writer, "    Author: %s | Repo: %s\n", commit.UserName, commit.RepoName)
	}

	if len(commits) > showCount {
		fmt.Fprintf(p.writer, "  ... and %d more commits\n", len(commits)-showCount)
	}
	fmt.Fprintln(p.writer)
}

// displayStatistics shows summary statistics.
func (p *Preview) displayStatistics(store *storage.MemoryStore) {
	commits := store.GetCommitsByTimeRange(time.Time{}, time.Now().Add(24*time.Hour))

	fmt.Fprintln(p.writer, "Statistics")
	fmt.Fprintln(p.writer, strings.Repeat("-", 60))
	fmt.Fprintf(p.writer, "  Total commits generated: %d\n", len(commits))
	fmt.Fprintf(p.writer, "  Developers: %d\n", len(p.seedData.Developers))
	fmt.Fprintf(p.writer, "  Repositories: %d\n", len(p.seedData.Repositories))
	fmt.Fprintf(p.writer, "  Preview duration: %d day(s)\n", p.config.Days)
	fmt.Fprintln(p.writer)
	fmt.Fprintln(p.writer, strings.Repeat("=", 60))
}
