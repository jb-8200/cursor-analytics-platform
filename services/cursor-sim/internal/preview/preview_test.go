package preview

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TASK-PREV-04: Create Preview Package and Config (RED)

func TestPreview_New(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		},
	}

	var buf bytes.Buffer
	cfg := Config{Days: 7, MaxCommits: 50}
	p := New(seedData, cfg, &buf)

	assert.NotNil(t, p)
	assert.Equal(t, seedData, p.seedData)
	assert.Equal(t, cfg, p.config)
	assert.Equal(t, &buf, p.writer)
}

func TestConfig_Defaults(t *testing.T) {
	// Test that Config struct can be created with default values
	cfg := Config{
		Days:       7,
		MaxCommits: 50,
		MaxEvents:  100,
	}

	assert.Equal(t, 7, cfg.Days)
	assert.Equal(t, 50, cfg.MaxCommits)
	assert.Equal(t, 100, cfg.MaxEvents)
}

func TestPreview_NewWithNilWriter(t *testing.T) {
	// Should handle nil writer gracefully (use stdout)
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{UserID: "user_001", Email: "alice@example.com"},
		},
	}

	cfg := Config{Days: 1, MaxCommits: 10}
	p := New(seedData, cfg, nil)

	assert.NotNil(t, p)
	// Writer should default to os.Stdout if nil provided
	assert.NotNil(t, p.writer)
}

func TestPreview_NewWithEmptyDevelopers(t *testing.T) {
	// Should handle empty developers list
	seedData := &seed.SeedData{
		Developers: []seed.Developer{},
	}

	var buf bytes.Buffer
	cfg := Config{Days: 7, MaxCommits: 50}
	p := New(seedData, cfg, &buf)

	require.NotNil(t, p)
	assert.Empty(t, p.seedData.Developers)
}

// TASK-PREV-05: Implement Preview Run Method (RED)

func TestPreview_Run(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:  "user_001",
				Email:   "alice@example.com",
				Name:    "Alice Developer",
				Seniority: "senior",
				WorkingHoursBand: seed.WorkingHours{Start: 9, End: 17, Peak: 14},
				PreferredModels:  []string{"claude-sonnet-4"},
			},
		},
		Repositories: []seed.Repository{
			{
				RepoName:        "acme-corp/platform",
				PrimaryLanguage: "Go",
			},
		},
	}

	var buf bytes.Buffer
	cfg := Config{Days: 7, MaxCommits: 10, MaxEvents: 5}
	p := New(seedData, cfg, &buf)

	err := p.Run(context.Background())
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "PREVIEW MODE")
	assert.Contains(t, output, "Alice Developer")
	assert.Contains(t, output, "Sample Commits")
	assert.Contains(t, output, "Statistics")
}

func TestPreview_RunWithTimeout(t *testing.T) {
	// Test that preview respects context timeout
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "alice@example.com",
				Name:   "Alice",
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var buf bytes.Buffer
	cfg := Config{Days: 1, MaxCommits: 5, MaxEvents: 5}
	p := New(seedData, cfg, &buf)

	err := p.Run(ctx)

	// Should complete within timeout or return context error
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("unexpected error: %v", err)
	}

	// If no error, output should exist
	if err == nil {
		assert.Greater(t, len(buf.String()), 0, "Should have output")
	}
}

func TestPreview_RunWithEmptyDevelopers(t *testing.T) {
	// Test with no developers
	seedData := &seed.SeedData{
		Developers: []seed.Developer{},
	}

	var buf bytes.Buffer
	cfg := Config{Days: 7, MaxCommits: 10}
	p := New(seedData, cfg, &buf)

	err := p.Run(context.Background())
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "PREVIEW MODE")
	assert.Contains(t, output, "0 developers")
}
