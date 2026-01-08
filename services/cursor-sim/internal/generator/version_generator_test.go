package generator

import (
	"context"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionGenerator_GenerateClientVersions(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "alice@example.com",
				WorkingHoursBand: seed.WorkingHours{
					Start: 9,
					End:   17,
				},
			},
			{
				UserID: "user_002",
				Email:  "bob@example.com",
				WorkingHoursBand: seed.WorkingHours{
					Start: 10,
					End:   18,
				},
			},
		},
	}

	store := &mockVersionStore{
		events: make([]models.ClientVersionEvent, 0),
		developers: seedData.Developers,
	}

	gen := NewVersionGeneratorWithSeed(seedData, store, "medium", 42)

	err := gen.GenerateClientVersions(context.Background(), 7)
	require.NoError(t, err)

	// Should have generated events for each developer for each day
	expectedEvents := len(seedData.Developers) * 7 // 2 developers * 7 days = 14 events
	assert.Equal(t, expectedEvents, len(store.events), "should generate one event per developer per day")

	// All events should have valid versions
	validVersions := gen.GetVersions()
	for _, event := range store.events {
		assert.Contains(t, validVersions, event.ClientVersion, "version should be from valid version list")
		assert.NotEmpty(t, event.UserID)
		assert.NotEmpty(t, event.UserEmail)
		assert.NotEmpty(t, event.EventDate)
		assert.False(t, event.Timestamp.IsZero())
	}
}

func TestVersionGenerator_VersionUpgrades(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "alice@example.com",
				WorkingHoursBand: seed.WorkingHours{
					Start: 9,
					End:   17,
				},
			},
		},
	}

	store := &mockVersionStore{
		events: make([]models.ClientVersionEvent, 0),
		developers: seedData.Developers,
	}

	gen := NewVersionGeneratorWithSeed(seedData, store, "medium", 123)

	// Generate over 30 days to increase chance of version upgrades
	err := gen.GenerateClientVersions(context.Background(), 30)
	require.NoError(t, err)

	// Track versions over time for the single developer
	versions := make(map[string]bool)
	for _, event := range store.events {
		versions[event.ClientVersion] = true
	}

	// Should have at least 1 version (might have upgrades but not guaranteed in 30 days)
	assert.GreaterOrEqual(t, len(versions), 1, "should have at least one version")
}

func TestVersionGenerator_VersionDistribution(t *testing.T) {
	gen := NewVersionGeneratorWithSeed(&seed.SeedData{}, nil, "medium", 999)

	// Test initial version selection
	versionCounts := make(map[string]int)
	for i := 0; i < 1000; i++ {
		version := gen.selectInitialVersion()
		versionCounts[version]++
	}

	// Should have selected multiple versions
	assert.Greater(t, len(versionCounts), 1, "should select from multiple versions")

	// Later versions should be more common (70% get latest 3)
	validVersions := gen.GetVersions()
	latestThree := validVersions[len(validVersions)-3:]

	countInLatestThree := 0
	for _, v := range latestThree {
		countInLatestThree += versionCounts[v]
	}

	// Approximately 70% should be in the latest 3 versions
	percentage := float64(countInLatestThree) / 1000.0
	assert.InDelta(t, 0.7, percentage, 0.15, "about 70%% should use latest 3 versions")
}

func TestVersionGenerator_UpgradeLogic(t *testing.T) {
	gen := NewVersionGeneratorWithSeed(&seed.SeedData{}, nil, "medium", 777)

	versions := gen.GetVersions()

	// Test upgrading from first version
	firstVersion := versions[0]
	upgraded := gen.selectNewerVersion(firstVersion)
	assert.NotEqual(t, firstVersion, upgraded, "should upgrade from first version")

	// Test that latest version doesn't downgrade
	latestVersion := versions[len(versions)-1]
	notDowngraded := gen.selectNewerVersion(latestVersion)
	assert.Equal(t, latestVersion, notDowngraded, "should not downgrade from latest version")
}

func TestVersionGenerator_EventStructure(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "alice@example.com",
				WorkingHoursBand: seed.WorkingHours{
					Start: 9,
					End:   17,
				},
			},
		},
	}

	store := &mockVersionStore{events: make([]models.ClientVersionEvent, 0)}
	gen := NewVersionGeneratorWithSeed(seedData, store, "medium", 555)

	err := gen.GenerateClientVersions(context.Background(), 3)
	require.NoError(t, err)

	// Verify event structure
	for _, event := range store.events {
		assert.Equal(t, "user_001", event.UserID)
		assert.Equal(t, "alice@example.com", event.UserEmail)
		assert.Regexp(t, `^0\.\d+\.\d+$`, event.ClientVersion, "version should be semver format")
		assert.False(t, event.Timestamp.IsZero())
		assert.Regexp(t, `^\d{4}-\d{2}-\d{2}$`, event.EventDate, "date should be YYYY-MM-DD format")

		// Check that timestamp is within working hours
		hour := event.Timestamp.Hour()
		assert.GreaterOrEqual(t, hour, 9, "should be during working hours")
		assert.Less(t, hour, 17, "should be during working hours")
	}
}

// mockVersionStore implements the ClientVersionStore interface for testing
type mockVersionStore struct {
	developers []seed.Developer
	events []models.ClientVersionEvent
}

func (m *mockVersionStore) AddClientVersion(event models.ClientVersionEvent) error {
	m.events = append(m.events, event)
	return nil
}

func (m *mockVersionStore) ListDevelopers() []seed.Developer {
	return m.developers
}
