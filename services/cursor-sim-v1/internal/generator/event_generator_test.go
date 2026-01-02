package generator

import (
	"testing"
	"time"

	"github.com/jb-8200/cursor-analytics-platform/services/cursor-sim/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEventGenerator(t *testing.T) {
	cfg := &config.Config{
		Developers:  10,
		Velocity:    "high",
		Fluctuation: 0.2,
		Seed:        12345,
	}

	devGen := NewDeveloperGenerator(cfg)
	developers, err := devGen.Generate()
	require.NoError(t, err)

	eventGen := NewEventGenerator(cfg, developers)
	require.NotNil(t, eventGen)
	assert.NotNil(t, eventGen.config)
	assert.NotNil(t, eventGen.developers)
	assert.NotNil(t, eventGen.commitChan)
	assert.NotNil(t, eventGen.changeChan)
}

func TestEventGenerator_Start(t *testing.T) {
	cfg := &config.Config{
		Developers:  5,
		Velocity:    "high",
		Fluctuation: 0.2,
		Seed:        12345,
	}

	devGen := NewDeveloperGenerator(cfg)
	developers, err := devGen.Generate()
	require.NoError(t, err)

	eventGen := NewEventGenerator(cfg, developers)

	// Start generation
	err = eventGen.Start()
	require.NoError(t, err)

	// Wait briefly for initial events (each dev generates one immediately)
	time.Sleep(100 * time.Millisecond)

	// Stop generation
	eventGen.Stop()

	// Should have generated some events
	commits := eventGen.GetCommits()
	changes := eventGen.GetChanges()

	assert.Greater(t, len(commits), 0, "Should have generated some commits")
	assert.Greater(t, len(changes), 0, "Should have generated some changes")
}

func TestEventGenerator_Stop(t *testing.T) {
	cfg := &config.Config{
		Developers:  3,
		Velocity:    "high",
		Fluctuation: 0.2,
		Seed:        12345,
	}

	devGen := NewDeveloperGenerator(cfg)
	developers, err := devGen.Generate()
	require.NoError(t, err)

	eventGen := NewEventGenerator(cfg, developers)
	err = eventGen.Start()
	require.NoError(t, err)

	time.Sleep(500 * time.Millisecond)

	// Stop should not panic
	eventGen.Stop()

	// Calling Stop again should not panic
	eventGen.Stop()
}

func TestEventGenerator_VelocityImpact(t *testing.T) {
	tests := []struct {
		velocity     string
		minEvents    int
		duration     time.Duration
		expectedRate float64 // events per second (approximate)
	}{
		{
			velocity:     "low",
			minEvents:    1,
			duration:     10 * time.Second,
			expectedRate: 0.014, // ~5 per hour per dev = 0.0014/sec per dev, with 10 devs ~0.014/sec
		},
		{
			velocity:     "high",
			minEvents:    2,
			duration:     3 * time.Second,
			expectedRate: 0.14, // ~50 per hour per dev = 0.014/sec per dev, with 10 devs ~0.14/sec
		},
	}

	for _, tt := range tests {
		t.Run(tt.velocity, func(t *testing.T) {
			cfg := &config.Config{
				Developers:  10,
				Velocity:    tt.velocity,
				Fluctuation: 0.1,
				Seed:        12345,
			}

			devGen := NewDeveloperGenerator(cfg)
			developers, err := devGen.Generate()
			require.NoError(t, err)

			eventGen := NewEventGenerator(cfg, developers)
			err = eventGen.Start()
			require.NoError(t, err)

			time.Sleep(tt.duration)
			eventGen.Stop()

			commits := eventGen.GetCommits()
			assert.GreaterOrEqual(t, len(commits), tt.minEvents,
				"Should have generated at least %d commits with %s velocity", tt.minEvents, tt.velocity)
		})
	}
}

func TestEventGenerator_CommitStructure(t *testing.T) {
	cfg := &config.Config{
		Developers:  3,
		Velocity:    "high",
		Fluctuation: 0.2,
		Seed:        12345,
	}

	devGen := NewDeveloperGenerator(cfg)
	developers, err := devGen.Generate()
	require.NoError(t, err)

	eventGen := NewEventGenerator(cfg, developers)
	err = eventGen.Start()
	require.NoError(t, err)

	time.Sleep(3 * time.Second)
	eventGen.Stop()

	commits := eventGen.GetCommits()
	require.Greater(t, len(commits), 0, "Should have generated commits")

	// Check first commit structure
	commit := commits[0]
	assert.NotEmpty(t, commit.Hash, "Commit hash should not be empty")
	assert.NotEmpty(t, commit.Message, "Commit message should not be empty")
	assert.NotEmpty(t, commit.UserID, "User ID should not be empty")
	assert.NotEmpty(t, commit.UserEmail, "User email should not be empty")
	assert.NotEmpty(t, commit.Repository, "Repository should not be empty")
	assert.NotEmpty(t, commit.Branch, "Branch should not be empty")
	assert.False(t, commit.Timestamp.IsZero(), "Timestamp should not be zero")
	assert.False(t, commit.IngestionTime.IsZero(), "Ingestion time should not be zero")

	// Check line counts are reasonable
	assert.GreaterOrEqual(t, commit.TotalLines, commit.LinesFromTAB+commit.LinesFromComposer,
		"Total lines should be >= AI lines")
	assert.GreaterOrEqual(t, commit.LinesFromTAB, 0, "Lines from TAB should be >= 0")
	assert.GreaterOrEqual(t, commit.LinesFromComposer, 0, "Lines from COMPOSER should be >= 0")
}

func TestEventGenerator_ChangeStructure(t *testing.T) {
	cfg := &config.Config{
		Developers:  3,
		Velocity:    "high",
		Fluctuation: 0.2,
		Seed:        12345,
	}

	devGen := NewDeveloperGenerator(cfg)
	developers, err := devGen.Generate()
	require.NoError(t, err)

	eventGen := NewEventGenerator(cfg, developers)
	err = eventGen.Start()
	require.NoError(t, err)

	time.Sleep(3 * time.Second)
	eventGen.Stop()

	changes := eventGen.GetChanges()
	require.Greater(t, len(changes), 0, "Should have generated changes")

	// Check first change structure
	change := changes[0]
	assert.NotEmpty(t, change.ChangeID, "Change ID should not be empty")
	assert.NotEmpty(t, change.CommitHash, "Commit hash should not be empty")
	assert.NotEmpty(t, change.UserID, "User ID should not be empty")
	assert.NotEmpty(t, change.Source, "Source should not be empty")
	assert.Contains(t, []string{"TAB", "COMPOSER"}, change.Source, "Source should be TAB or COMPOSER")
	assert.NotEmpty(t, change.Model, "Model should not be empty")
	assert.NotEmpty(t, change.FilePath, "File path should not be empty")
	assert.NotEmpty(t, change.FileExtension, "File extension should not be empty")
	assert.Regexp(t, `^\.\w+$`, change.FileExtension, "File extension should start with dot")
	assert.False(t, change.Timestamp.IsZero(), "Timestamp should not be zero")
	assert.GreaterOrEqual(t, change.LinesAdded, 0, "Lines added should be >= 0")
	assert.GreaterOrEqual(t, change.LinesRemoved, 0, "Lines removed should be >= 0")
}

func TestEventGenerator_TabComposerRatio(t *testing.T) {
	// Default ratio is 0.7 (70% TAB, 30% COMPOSER)
	cfg := &config.Config{
		Developers:  10,
		Velocity:    "high",
		Fluctuation: 0.1,
		Seed:        12345,
	}

	devGen := NewDeveloperGenerator(cfg)
	developers, err := devGen.Generate()
	require.NoError(t, err)

	eventGen := NewEventGenerator(cfg, developers)
	err = eventGen.Start()
	require.NoError(t, err)

	time.Sleep(10 * time.Second)
	eventGen.Stop()

	changes := eventGen.GetChanges()
	require.Greater(t, len(changes), 20, "Need sufficient samples for ratio test")

	// Count TAB vs COMPOSER
	tabCount := 0
	composerCount := 0
	for _, change := range changes {
		if change.Source == "TAB" {
			tabCount++
		} else if change.Source == "COMPOSER" {
			composerCount++
		}
	}

	totalCount := tabCount + composerCount
	tabRatio := float64(tabCount) / float64(totalCount)

	// Should be approximately 70% TAB (allow ±15% variance)
	assert.InDelta(t, 0.70, tabRatio, 0.15,
		"TAB ratio should be approximately 0.70, got %.2f (%d TAB, %d COMPOSER)",
		tabRatio, tabCount, composerCount)
}

func TestEventGenerator_UniqueCommitHashes(t *testing.T) {
	cfg := &config.Config{
		Developers:  5,
		Velocity:    "high",
		Fluctuation: 0.2,
		Seed:        12345,
	}

	devGen := NewDeveloperGenerator(cfg)
	developers, err := devGen.Generate()
	require.NoError(t, err)

	eventGen := NewEventGenerator(cfg, developers)
	err = eventGen.Start()
	require.NoError(t, err)

	time.Sleep(5 * time.Second)
	eventGen.Stop()

	commits := eventGen.GetCommits()
	require.Greater(t, len(commits), 3, "Need multiple commits for uniqueness test")

	// Check all hashes are unique
	hashSet := make(map[string]bool)
	for _, commit := range commits {
		assert.False(t, hashSet[commit.Hash],
			"Commit hash should be unique: %s", commit.Hash)
		hashSet[commit.Hash] = true
	}
}

func TestEventGenerator_ChangesLinkedToCommits(t *testing.T) {
	cfg := &config.Config{
		Developers:  3,
		Velocity:    "high",
		Fluctuation: 0.2,
		Seed:        12345,
	}

	devGen := NewDeveloperGenerator(cfg)
	developers, err := devGen.Generate()
	require.NoError(t, err)

	eventGen := NewEventGenerator(cfg, developers)
	err = eventGen.Start()
	require.NoError(t, err)

	time.Sleep(3 * time.Second)
	eventGen.Stop()

	commits := eventGen.GetCommits()
	changes := eventGen.GetChanges()

	require.Greater(t, len(commits), 0, "Should have commits")
	require.Greater(t, len(changes), 0, "Should have changes")

	// Build set of valid commit hashes
	validHashes := make(map[string]bool)
	for _, commit := range commits {
		validHashes[commit.Hash] = true
	}

	// All changes should reference valid commit hashes
	for i, change := range changes {
		assert.True(t, validHashes[change.CommitHash],
			"Change at index %d should reference a valid commit hash: %s", i, change.CommitHash)
	}
}

func TestEventGenerator_DeveloperActivity(t *testing.T) {
	cfg := &config.Config{
		Developers:  5,
		Velocity:    "high",
		Fluctuation: 0.2,
		Seed:        12345,
	}

	devGen := NewDeveloperGenerator(cfg)
	developers, err := devGen.Generate()
	require.NoError(t, err)

	eventGen := NewEventGenerator(cfg, developers)
	err = eventGen.Start()
	require.NoError(t, err)

	time.Sleep(5 * time.Second)
	eventGen.Stop()

	commits := eventGen.GetCommits()
	require.Greater(t, len(commits), 0, "Should have commits")

	// Check that multiple developers are active
	developerActivity := make(map[string]int)
	for _, commit := range commits {
		developerActivity[commit.UserID]++
	}

	// At least 2 developers should have generated commits
	assert.GreaterOrEqual(t, len(developerActivity), 2,
		"Multiple developers should be generating commits")
}
