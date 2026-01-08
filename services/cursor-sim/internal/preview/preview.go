// Package preview implements seed file validation and preview mode.
// Preview mode allows quick validation of seed files without full data generation.
// Inspired by NVIDIA NeMo DataDesigner preview patterns.
package preview

import (
	"io"
	"os"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
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
