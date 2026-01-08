package generator

import (
	"context"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelGenerator_GenerateModelUsage(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:          "user_001",
				Email:           "alice@example.com",
				PreferredModels: []string{"gpt-4-turbo", "claude-3-sonnet"},
				PRBehavior: seed.PRBehavior{
					PRsPerWeek: 5.0,
				},
			},
			{
				UserID:          "user_002",
				Email:           "bob@example.com",
				PreferredModels: []string{"claude-3-opus"},
				PRBehavior: seed.PRBehavior{
					PRsPerWeek: 2.0,
				},
			},
		},
	}

	store := &mockModelStore{
		developers: seedData.Developers,
		usage:      make([]models.ModelUsageEvent, 0),
	}

	gen := NewModelGeneratorWithSeed(seedData, store, "medium", 42)

	err := gen.GenerateModelUsage(context.Background(), 7)
	require.NoError(t, err)

	// Should have generated model usage events
	assert.True(t, len(store.usage) > 0, "should generate model usage events")

	// Check that models match preferred models
	modelCounts := make(map[string]int)
	for _, usage := range store.usage {
		modelCounts[usage.ModelName]++
		assert.Contains(t, []string{"gpt-4-turbo", "claude-3-sonnet", "claude-3-opus"}, usage.ModelName)
	}

	// Both developers should have events (alice has higher PRsPerWeek so tends to have more)
	aliceEvents := 0
	bobEvents := 0
	for _, usage := range store.usage {
		if usage.UserEmail == "alice@example.com" {
			aliceEvents++
		} else if usage.UserEmail == "bob@example.com" {
			bobEvents++
		}
	}

	assert.True(t, aliceEvents > 0, "alice should have model usage events")
	assert.True(t, bobEvents > 0, "bob should have model usage events")
}

func TestModelGenerator_UsageTypes(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:          "user_001",
				Email:           "alice@example.com",
				PreferredModels: []string{"gpt-4-turbo"},
				ChatVsCodeRatio: seed.ChatCodeRatio{
					Chat: 0.3,
					Code: 0.7,
				},
				PRBehavior: seed.PRBehavior{
					PRsPerWeek: 10.0,
				},
			},
		},
	}

	store := &mockModelStore{
		developers: seedData.Developers,
		usage:      make([]models.ModelUsageEvent, 0),
	}

	gen := NewModelGeneratorWithSeed(seedData, store, "high", 123)

	err := gen.GenerateModelUsage(context.Background(), 7)
	require.NoError(t, err)

	// Count chat vs code usage
	chatCount := 0
	codeCount := 0
	for _, usage := range store.usage {
		if usage.UsageType == "chat" {
			chatCount++
		} else if usage.UsageType == "code" {
			codeCount++
		}
	}

	// Ratio should be approximately 30/70
	totalUsage := chatCount + codeCount
	chatRatio := float64(chatCount) / float64(totalUsage)

	assert.InDelta(t, 0.3, chatRatio, 0.2, "chat ratio should be approximately 30%%")
}

func TestModelGenerator_EventStructure(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID:          "user_001",
				Email:           "alice@example.com",
				PreferredModels: []string{"gpt-4-turbo"},
				PRBehavior: seed.PRBehavior{
					PRsPerWeek: 3.0,
				},
			},
		},
	}

	store := &mockModelStore{
		developers: seedData.Developers,
		usage:      make([]models.ModelUsageEvent, 0),
	}
	gen := NewModelGeneratorWithSeed(seedData, store, "medium", 999)

	err := gen.GenerateModelUsage(context.Background(), 7)
	require.NoError(t, err)

	// Should generate events
	assert.True(t, len(store.usage) > 0, "should generate model usage events")

	// Verify event structure
	for _, usage := range store.usage {
		assert.Equal(t, "user_001", usage.UserID)
		assert.Equal(t, "alice@example.com", usage.UserEmail)
		assert.Equal(t, "gpt-4-turbo", usage.ModelName)
		assert.Contains(t, []string{"chat", "code"}, usage.UsageType)
		assert.False(t, usage.Timestamp.IsZero())
		assert.NotEmpty(t, usage.EventDate)
	}
}

// mockModelStore implements the ModelStore interface for testing
type mockModelStore struct {
	developers []seed.Developer
	usage []models.ModelUsageEvent
}

func (m *mockModelStore) AddModelUsage(usage models.ModelUsageEvent) error {
	m.usage = append(m.usage, usage)
	return nil
}

func (m *mockModelStore) ListDevelopers() []seed.Developer {
	return m.developers
}
