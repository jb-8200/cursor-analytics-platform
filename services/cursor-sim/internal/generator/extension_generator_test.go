package generator

import (
	"context"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtensionGenerator_GenerateFileExtensions(t *testing.T) {
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

	store := &mockExtensionStore{
		events: make([]models.FileExtensionEvent, 0),
	}

	gen := NewExtensionGeneratorWithSeed(seedData, store, "medium", 42)

	err := gen.GenerateFileExtensions(context.Background(), 7)
	require.NoError(t, err)

	// Should have generated multiple events (2-5 per developer per day)
	minEvents := len(seedData.Developers) * 7 * 2 // At least 2 per day
	assert.GreaterOrEqual(t, len(store.events), minEvents, "should generate multiple events")

	// All events should have valid data
	validExtensions := gen.GetExtensions()
	for _, event := range store.events {
		assert.Contains(t, validExtensions, event.FileExtension)
		assert.Greater(t, event.LinesSuggested, 0)
		assert.GreaterOrEqual(t, event.LinesAccepted, 0)
		assert.GreaterOrEqual(t, event.LinesRejected, 0)
		assert.Equal(t, event.LinesSuggested, event.LinesAccepted+event.LinesRejected)
		assert.NotEmpty(t, event.UserID)
		assert.NotEmpty(t, event.EventDate)
		assert.False(t, event.Timestamp.IsZero())
	}
}

func TestExtensionGenerator_AcceptanceRatio(t *testing.T) {
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

	store := &mockExtensionStore{
		events: make([]models.FileExtensionEvent, 0),
	}

	gen := NewExtensionGeneratorWithSeed(seedData, store, "medium", 99)

	err := gen.GenerateFileExtensions(context.Background(), 30)
	require.NoError(t, err)

	// Calculate acceptance ratio
	totalSuggested := 0
	totalAccepted := 0
	for _, event := range store.events {
		totalSuggested += event.LinesSuggested
		totalAccepted += event.LinesAccepted
	}

	acceptanceRatio := float64(totalAccepted) / float64(totalSuggested)
	// Should be around 50-70% acceptance (accounting for random variation)
	assert.InDelta(t, 0.60, acceptanceRatio, 0.2, "acceptance ratio should be around 50-70%%")
}

func TestExtensionGenerator_FavoriteExtensions(t *testing.T) {
	gen := NewExtensionGeneratorWithSeed(&seed.SeedData{}, nil, "medium", 111)

	// Test favorite selection
	favorites := gen.selectFavoriteExtensions(3)
	assert.Equal(t, 3, len(favorites))

	// Generate events and check extension distribution
	extCount := make(map[string]int)
	for i := 0; i < 1000; i++ {
		ext := gen.selectExtension(favorites)
		extCount[ext]++
	}

	// Favorites should appear more frequently
	favoriteCount := 0
	for _, fav := range favorites {
		favoriteCount += extCount[fav]
	}

	// 80% should be from favorites
	percentage := float64(favoriteCount) / 1000.0
	assert.InDelta(t, 0.8, percentage, 0.15, "favorites should be selected 80%% of the time")
}

func TestExtensionGenerator_EventStructure(t *testing.T) {
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

	store := &mockExtensionStore{events: make([]models.FileExtensionEvent, 0)}
	gen := NewExtensionGeneratorWithSeed(seedData, store, "medium", 222)

	err := gen.GenerateFileExtensions(context.Background(), 2)
	require.NoError(t, err)

	for _, event := range store.events {
		assert.Equal(t, "user_001", event.UserID)
		assert.Equal(t, "alice@example.com", event.UserEmail)
		assert.NotEmpty(t, event.FileExtension)
		assert.Greater(t, event.LinesSuggested, 0)
		assert.Greater(t, event.LinesSuggested, event.LinesAccepted+event.LinesRejected-1) // Allow rounding
		assert.NotEmpty(t, event.EventDate)

		// Check accept/reject balance
		if event.WasAccepted {
			assert.Greater(t, event.LinesAccepted, 0)
		} else {
			assert.Greater(t, event.LinesRejected, 0)
		}
	}
}

// mockExtensionStore implements the FileExtensionStore interface for testing
type mockExtensionStore struct {
	events []models.FileExtensionEvent
}

func (m *mockExtensionStore) AddFileExtension(event models.FileExtensionEvent) error {
	m.events = append(m.events, event)
	return nil
}
