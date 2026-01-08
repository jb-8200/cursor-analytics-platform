package preview

import (
	"bytes"
	"context"
	"errors"
	"strings"
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
	// Test with no developers - should fail validation
	seedData := &seed.SeedData{
		Developers: []seed.Developer{},
	}

	var buf bytes.Buffer
	cfg := Config{Days: 7, MaxCommits: 10}
	p := New(seedData, cfg, &buf)

	err := p.Run(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no developers defined")
}

// TASK-PREV-06: Implement Preview Output Formatters (REFACTOR)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "short string unchanged",
			input:    "short",
			maxLen:   50,
			expected: "short",
		},
		{
			name:     "exact length unchanged",
			input:    "exactly fifty characters in this string here now!",
			maxLen:   50,
			expected: "exactly fifty characters in this string here now!",
		},
		{
			name:     "long string truncated",
			input:    "This is a very long commit message that should be truncated to fit nicely",
			maxLen:   50,
			expected: "This is a very long commit message that should...",
		},
		{
			name:     "very long truncated",
			input:    strings.Repeat("a", 100),
			maxLen:   20,
			expected: strings.Repeat("a", 17) + "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncate(tt.input, tt.maxLen)
			assert.Equal(t, tt.expected, result)
			assert.LessOrEqual(t, len(result), tt.maxLen, "truncated string should not exceed maxLen")
		})
	}
}

func TestFormatWorkingHours(t *testing.T) {
	tests := []struct {
		name     string
		hours    seed.WorkingHours
		expected string
	}{
		{
			name:     "standard hours",
			hours:    seed.WorkingHours{Start: 9, End: 17, Peak: 14},
			expected: "09:00 - 17:00",
		},
		{
			name:     "early hours",
			hours:    seed.WorkingHours{Start: 6, End: 14, Peak: 10},
			expected: "06:00 - 14:00",
		},
		{
			name:     "late hours",
			hours:    seed.WorkingHours{Start: 13, End: 21, Peak: 17},
			expected: "13:00 - 21:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatWorkingHours(tt.hours)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TASK-PREV-08: Implement Seed Validators (RED)

func TestPreview_ValidateSeed_ValidData(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "alice",
				Email:  "alice@example.com",
				WorkingHoursBand: seed.WorkingHours{Start: 9, End: 17},
				PreferredModels: []string{"claude-sonnet-4.5"},
			},
		},
	}

	var buf bytes.Buffer
	cfg := Config{Days: 7, MaxCommits: 10}
	p := New(seedData, cfg, &buf)

	err := p.validateSeed()
	assert.NoError(t, err)
	assert.Empty(t, p.warnings)
}

func TestPreview_ValidateSeed_InvalidModel(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:          "alice",
				PreferredModels: []string{"gpt-5000"},
			},
		},
	}

	var buf bytes.Buffer
	cfg := Config{Days: 7, MaxCommits: 10}
	p := New(seedData, cfg, &buf)

	err := p.validateSeed()
	assert.NoError(t, err)
	assert.NotEmpty(t, p.warnings)
	assert.Contains(t, p.warnings[0], "Unknown model 'gpt-5000'")
}

func TestPreview_ValidateSeed_InvalidWorkingHours(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:           "alice",
				WorkingHoursBand: seed.WorkingHours{Start: 25, End: 30},
			},
		},
	}

	var buf bytes.Buffer
	cfg := Config{Days: 7, MaxCommits: 10}
	p := New(seedData, cfg, &buf)

	err := p.validateSeed()
	assert.NoError(t, err)
	assert.NotEmpty(t, p.warnings)
	assert.Contains(t, p.warnings[0], "Invalid start hour 25")
}

func TestPreview_ValidateSeed_NoDevelopers(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{},
	}

	var buf bytes.Buffer
	cfg := Config{Days: 7, MaxCommits: 10}
	p := New(seedData, cfg, &buf)

	err := p.validateSeed()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no developers defined")
}

// TASK-PREV-09: Display Validation Warnings (REFACTOR)

func TestPreview_DisplayWarnings_NoWarnings(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{UserID: "alice", Email: "alice@example.com"},
		},
	}

	var buf bytes.Buffer
	cfg := Config{Days: 7, MaxCommits: 10}
	p := New(seedData, cfg, &buf)
	p.warnings = []string{}

	p.displayWarnings()

	output := buf.String()
	assert.Contains(t, output, "Validation Warnings")
	assert.Contains(t, output, "No validation warnings")
}

func TestPreview_DisplayWarnings_MultipleWarnings(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:           "alice",
				WorkingHoursBand: seed.WorkingHours{Start: 25, End: 30},
				PreferredModels: []string{"gpt-5000"},
			},
		},
	}

	var buf bytes.Buffer
	cfg := Config{Days: 7, MaxCommits: 10}
	p := New(seedData, cfg, &buf)
	p.warnings = []string{
		"Developer alice: Invalid start hour 25",
		"Developer alice: Unknown model 'gpt-5000'",
	}

	p.displayWarnings()

	output := buf.String()
	assert.Contains(t, output, "Validation Warnings")
	assert.Contains(t, output, "Invalid start hour 25")
	assert.Contains(t, output, "Unknown model 'gpt-5000'")
}
