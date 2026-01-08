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
	warnings []string // Validation warnings collected during seed validation
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

	// Validate seed data
	if err := p.validateSeed(); err != nil {
		return err
	}

	// Display developer summary
	p.displayDeveloperSummary()

	// Display validation warnings
	p.displayWarnings()

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
		// Truncate name and email for 80-column display
		displayName := truncate(dev.Name, 30)
		displayEmail := truncate(dev.Email, 35)
		fmt.Fprintf(p.writer, "  • %s (%s)\n", displayName, displayEmail)

		// Display ID and seniority
		if dev.UserID != "" {
			fmt.Fprintf(p.writer, "    ID: %s", dev.UserID)
			if dev.Seniority != "" {
				fmt.Fprintf(p.writer, " | Seniority: %s", dev.Seniority)
			}
			fmt.Fprintln(p.writer)
		}

		// Display working hours if configured
		if dev.WorkingHoursBand.Start > 0 || dev.WorkingHoursBand.End > 0 {
			fmt.Fprintf(p.writer, "    Hours: %s\n", formatWorkingHours(dev.WorkingHoursBand))
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
		// Truncate commit message to fit 80-column display
		displayMsg := truncate(commit.Message, 55)
		fmt.Fprintf(p.writer, "  [%s] %s\n", commit.CommitTs.Format("15:04:05"), displayMsg)

		// Truncate author and repo names
		displayAuthor := truncate(commit.UserName, 20)
		displayRepo := truncate(commit.RepoName, 25)
		fmt.Fprintf(p.writer, "    Author: %s | Repo: %s\n", displayAuthor, displayRepo)
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

	// Calculate average commits per developer
	if len(p.seedData.Developers) > 0 {
		avgCommits := float64(len(commits)) / float64(len(p.seedData.Developers))
		fmt.Fprintf(p.writer, "  Avg commits/dev: %.1f\n", avgCommits)
	}

	// Calculate total lines changed
	totalAdded, totalDeleted := 0, 0
	for _, commit := range commits {
		totalAdded += commit.TotalLinesAdded
		totalDeleted += commit.TotalLinesDeleted
	}
	fmt.Fprintf(p.writer, "  Total lines: +%d / -%d\n", totalAdded, totalDeleted)

	fmt.Fprintf(p.writer, "  Preview duration: %d day(s)\n", p.config.Days)
	fmt.Fprintln(p.writer)
	fmt.Fprintln(p.writer, strings.Repeat("=", 60))
}

// Helper functions for formatting

// truncate truncates a string to maxLen, adding "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return s[:maxLen]
	}
	// Trim trailing spaces before adding ellipsis
	truncated := strings.TrimRight(s[:maxLen-3], " ")
	return truncated + "..."
}

// formatWorkingHours formats working hours as "HH:MM - HH:MM".
func formatWorkingHours(hours seed.WorkingHours) string {
	return fmt.Sprintf("%02d:00 - %02d:00", hours.Start, hours.End)
}

// Validation functions

// validModels contains the list of supported Claude models.
var validModels = map[string]bool{
	"claude-haiku-4": true,
	"claude-haiku-4-5-20251001": true,
	"claude-sonnet-4": true,
	"claude-sonnet-4.5": true,
	"claude-sonnet-4-5-20241022": true,
	"claude-opus-4": true,
	"claude-opus-4.5": true,
	"claude-opus-4-5-20251101": true,
}

// validateSeed validates the seed data and collects warnings.
// Returns error if seed has fatal issues (e.g., no developers), nil otherwise with warnings in p.warnings.
func (p *Preview) validateSeed() error {
	p.warnings = []string{}

	// Fatal validation: must have at least one developer
	if len(p.seedData.Developers) == 0 {
		return fmt.Errorf("validation failed: no developers defined")
	}

	// Validate each developer
	for _, dev := range p.seedData.Developers {
		p.validateDeveloper(dev)
	}

	return nil
}

// validateDeveloper validates a single developer and adds warnings.
func (p *Preview) validateDeveloper(dev seed.Developer) {
	// Validate working hours
	if dev.WorkingHoursBand.Start < 0 || dev.WorkingHoursBand.Start > 23 {
		p.warnings = append(p.warnings,
			fmt.Sprintf("Developer %s: Invalid start hour %d (must be 0-23)",
				dev.UserID, dev.WorkingHoursBand.Start))
	}
	if dev.WorkingHoursBand.End < 0 || dev.WorkingHoursBand.End > 23 {
		p.warnings = append(p.warnings,
			fmt.Sprintf("Developer %s: Invalid end hour %d (must be 0-23)",
				dev.UserID, dev.WorkingHoursBand.End))
	}

	// Validate models
	for _, model := range dev.PreferredModels {
		if !validModels[model] {
			p.warnings = append(p.warnings,
				fmt.Sprintf("Developer %s: Unknown model '%s'", dev.UserID, model))
		}
	}
}

// displayWarnings displays validation warnings in a formatted section.
func (p *Preview) displayWarnings() {
	fmt.Fprintln(p.writer, "Validation Warnings")
	fmt.Fprintln(p.writer, strings.Repeat("-", 60))

	if len(p.warnings) == 0 {
		fmt.Fprintln(p.writer, "  ✅ No validation warnings")
	} else {
		for _, warning := range p.warnings {
			fmt.Fprintf(p.writer, "  ⚠️  %s\n", warning)
		}
	}
	fmt.Fprintln(p.writer)
}
