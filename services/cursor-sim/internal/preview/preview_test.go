package preview

import (
	"bytes"
	"testing"

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
