package generator

import (
	"math"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHarveyGenerator_GenerateEvents(t *testing.T) {
	seedData := &seed.SeedData{
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Harvey: &seed.HarveySeedConfig{
				Enabled:       true,
				TotalUsage:    seed.UsageRange{Min: 100, Max: 200},
				ModelsUsed:    []string{"GPT-4", "Claude"},
				PracticeAreas: []string{"Corporate", "Litigation"},
			},
		},
		Developers: []seed.Developer{
			{
				UserID: "atty_001",
				Email:  "john@firm.com",
				Name:   "John Attorney",
			},
		},
	}

	gen := NewHarveyGeneratorWithSeed(seedData, 12345)
	config := DefaultHarveyConfig()

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)

	events := gen.GenerateEvents(from, to, config)

	// Expect ~5 events/day * 7 days * 1 user = ~35 events
	assert.True(t, len(events) >= 20, "Expected at least 20 events, got %d", len(events))
	assert.True(t, len(events) <= 50, "Expected at most 50 events, got %d", len(events))

	// Verify all events are from our user
	for _, e := range events {
		assert.Equal(t, "john@firm.com", e.User)
		assert.NotEmpty(t, e.EventID)
		assert.NotEmpty(t, e.MessageID)
		assert.False(t, e.Time.IsZero())
		assert.NotEmpty(t, e.Task)
	}
}

func TestHarveyGenerator_TaskDistribution(t *testing.T) {
	seedData := &seed.SeedData{
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Harvey: &seed.HarveySeedConfig{
				Enabled:       true,
				TotalUsage:    seed.UsageRange{Min: 100, Max: 200},
				ModelsUsed:    []string{"GPT-4"},
				PracticeAreas: []string{"Corporate"},
			},
		},
		Developers: []seed.Developer{
			{UserID: "atty_001", Email: "user@firm.com", Name: "Attorney"},
		},
	}

	gen := NewHarveyGeneratorWithSeed(seedData, 12345)
	config := DefaultHarveyConfig()

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	events := gen.GenerateEvents(from, to, config)

	// Count task types
	counts := make(map[models.HarveyTask]int)
	for _, e := range events {
		counts[e.Task]++
	}

	total := float64(len(events))
	require.True(t, total > 0, "Expected some events to be generated")

	// Allow 15% tolerance from configured rates (35%, 30%, 25%, 10%)
	assistRate := float64(counts[models.HarveyTaskAssist]) / total
	draftRate := float64(counts[models.HarveyTaskDraft]) / total
	reviewRate := float64(counts[models.HarveyTaskReview]) / total
	researchRate := float64(counts[models.HarveyTaskResearch]) / total

	assertRate(t, assistRate, 0.35, 0.15)
	assertRate(t, draftRate, 0.30, 0.15)
	assertRate(t, reviewRate, 0.25, 0.15)
	assertRate(t, researchRate, 0.10, 0.15)
}

func TestHarveyGenerator_Reproducible(t *testing.T) {
	seedData := &seed.SeedData{
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Harvey: &seed.HarveySeedConfig{
				Enabled:       true,
				TotalUsage:    seed.UsageRange{Min: 50, Max: 100},
				ModelsUsed:    []string{"GPT-4"},
				PracticeAreas: []string{"Corporate"},
			},
		},
		Developers: []seed.Developer{
			{UserID: "atty_001", Email: "user@firm.com", Name: "Attorney"},
		},
	}

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)
	config := DefaultHarveyConfig()

	// Same seed should produce identical results
	gen1 := NewHarveyGeneratorWithSeed(seedData, 12345)
	events1 := gen1.GenerateEvents(from, to, config)

	gen2 := NewHarveyGeneratorWithSeed(seedData, 12345)
	events2 := gen2.GenerateEvents(from, to, config)

	require.Equal(t, len(events1), len(events2))
	for i := range events1 {
		assert.Equal(t, events1[i].EventID, events2[i].EventID)
		assert.Equal(t, events1[i].Task, events2[i].Task)
		assert.Equal(t, events1[i].User, events2[i].User)
		assert.Equal(t, events1[i].Time.Unix(), events2[i].Time.Unix())
	}
}

func TestHarveyGenerator_WorkingHoursConstraint(t *testing.T) {
	seedData := &seed.SeedData{
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Harvey: &seed.HarveySeedConfig{
				Enabled:       true,
				TotalUsage:    seed.UsageRange{Min: 100, Max: 200},
				ModelsUsed:    []string{"GPT-4"},
				PracticeAreas: []string{"Corporate"},
			},
		},
		Developers: []seed.Developer{
			{UserID: "atty_001", Email: "user@firm.com", Name: "Attorney"},
		},
	}

	gen := NewHarveyGeneratorWithSeed(seedData, 12345)
	config := DefaultHarveyConfig()
	config.WorkingHours.Start = 8  // 8 AM
	config.WorkingHours.End = 18   // 6 PM

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)

	events := gen.GenerateEvents(from, to, config)

	// Verify all events are within working hours
	for _, e := range events {
		hour := e.Time.Hour()
		assert.True(t, hour >= config.WorkingHours.Start && hour < config.WorkingHours.End,
			"Event at %s (hour %d) outside working hours %d-%d",
			e.Time.Format(time.RFC3339), hour, config.WorkingHours.Start, config.WorkingHours.End)
	}
}

func TestHarveyGenerator_NoUsersReturnsEmpty(t *testing.T) {
	seedData := &seed.SeedData{
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Harvey: &seed.HarveySeedConfig{
				Enabled:       true,
				TotalUsage:    seed.UsageRange{Min: 100, Max: 200},
				ModelsUsed:    []string{"GPT-4"},
				PracticeAreas: []string{"Corporate"},
			},
		},
		Developers: []seed.Developer{},
	}

	gen := NewHarveyGeneratorWithSeed(seedData, 12345)
	config := DefaultHarveyConfig()

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)

	events := gen.GenerateEvents(from, to, config)

	assert.Empty(t, events, "Expected no events when no users configured")
}

func TestHarveyGenerator_SentimentDistribution(t *testing.T) {
	seedData := &seed.SeedData{
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Harvey: &seed.HarveySeedConfig{
				Enabled:       true,
				TotalUsage:    seed.UsageRange{Min: 100, Max: 200},
				ModelsUsed:    []string{"GPT-4"},
				PracticeAreas: []string{"Corporate"},
			},
		},
		Developers: []seed.Developer{
			{UserID: "atty_001", Email: "user@firm.com", Name: "Attorney"},
		},
	}

	gen := NewHarveyGeneratorWithSeed(seedData, 12345)
	config := DefaultHarveyConfig()

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	events := gen.GenerateEvents(from, to, config)

	// Count sentiments
	counts := make(map[models.HarveySentiment]int)
	for _, e := range events {
		counts[e.FeedbackSentiment]++
	}

	total := float64(len(events))
	require.True(t, total > 0, "Expected some events to be generated")

	// Allow 15% tolerance from configured rates (70% positive, 20% neutral, 10% negative)
	positiveRate := float64(counts[models.HarveySentimentPositive]) / total
	neutralRate := float64(counts[models.HarveySentimentNeutral]) / total
	negativeRate := float64(counts[models.HarveySentimentNegative]) / total

	assertRate(t, positiveRate, 0.70, 0.15)
	assertRate(t, neutralRate, 0.20, 0.15)
	assertRate(t, negativeRate, 0.10, 0.15)
}

func TestDefaultHarveyConfig(t *testing.T) {
	config := DefaultHarveyConfig()

	assert.Equal(t, 5.0, config.BaseEventsPerDay)
	assert.Equal(t, 8, config.WorkingHours.Start)
	assert.Equal(t, 18, config.WorkingHours.End)

	// Verify task distribution sums to 1.0
	totalTask := 0.0
	for _, rate := range config.TaskDistribution {
		totalTask += rate
	}
	assert.InDelta(t, 1.0, totalTask, 0.01)

	// Verify sentiment rates sum to 1.0
	totalSentiment := 0.0
	for _, rate := range config.SentimentRates {
		totalSentiment += rate
	}
	assert.InDelta(t, 1.0, totalSentiment, 0.01)
}

// Helper function to assert rate is within tolerance
func assertRate(t *testing.T, actual, expected, tolerance float64) {
	t.Helper()
	diff := math.Abs(actual - expected)
	if diff > tolerance {
		t.Errorf("Rate %.2f outside tolerance of %.2f±%.2f (diff=%.2f)",
			actual, expected, tolerance, diff)
	}
}
