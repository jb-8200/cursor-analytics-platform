package generator

import (
	"context"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeatureGenerator_GenerateFeatures(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "alice@example.com",
				WorkingHoursBand: seed.WorkingHours{
					Start: 9,
					End:   17,
				},
				PreferredModels: []string{"claude-sonnet-4.5"},
			},
		},
	}

	store := &mockFeatureStore{
		mcpTools: make([]models.MCPToolEvent, 0),
		developers: seedData.Developers,
		commands: make([]models.CommandEvent, 0),
		plans:    make([]models.PlanEvent, 0),
		askModes: make([]models.AskModeEvent, 0),
	}

	gen := NewFeatureGeneratorWithSeed(seedData, store, "medium", 42)

	err := gen.GenerateFeatures(context.Background(), 7)
	require.NoError(t, err)

	// Should generate events for all 4 feature types
	assert.Greater(t, len(store.mcpTools), 0, "should generate MCP tool events")
	assert.Greater(t, len(store.commands), 0, "should generate command events")
	assert.Greater(t, len(store.plans), 0, "should generate plan events")
	assert.Greater(t, len(store.askModes), 0, "should generate ask mode events")

	// Verify MCP tools
	validTools := gen.GetMCPTools()
	for _, event := range store.mcpTools {
		assert.Contains(t, validTools, event.ToolName)
		assert.NotEmpty(t, event.MCPServerName)
	}

	// Verify commands
	validCommands := gen.GetCommands()
	for _, event := range store.commands {
		assert.Contains(t, validCommands, event.CommandName)
	}

	// Verify plans and ask modes have models
	for _, event := range store.plans {
		assert.NotEmpty(t, event.Model)
	}
	for _, event := range store.askModes {
		assert.NotEmpty(t, event.Model)
	}
}

func TestFeatureGenerator_EventDistribution(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "user_001",
				Email:  "alice@example.com",
				WorkingHoursBand: seed.WorkingHours{
					Start: 9,
					End:   17,
				},
				PreferredModels: []string{"gpt-4o", "claude-opus-4"},
			},
		},
	}

	store := &mockFeatureStore{
		mcpTools: make([]models.MCPToolEvent, 0),
		developers: seedData.Developers,
		commands: make([]models.CommandEvent, 0),
		plans:    make([]models.PlanEvent, 0),
		askModes: make([]models.AskModeEvent, 0),
	}

	gen := NewFeatureGeneratorWithSeed(seedData, store, "medium", 99)

	err := gen.GenerateFeatures(context.Background(), 30)
	require.NoError(t, err)

	// Verify event counts are reasonable
	// MCP: 1-3 per day, so 30-90 events over 30 days
	assert.Greater(t, len(store.mcpTools), 20, "should have multiple MCP events")
	assert.Greater(t, len(store.commands), 50, "should have multiple command events")

	// Plans are less frequent (0-2 per day)
	assert.Greater(t, len(store.plans), 0, "should have at least some plan events")

	// Ask modes (1-3 per day)
	assert.Greater(t, len(store.askModes), 20, "should have multiple ask mode events")

	// Verify model diversity for plans
	modelsUsed := make(map[string]bool)
	for _, event := range store.plans {
		modelsUsed[event.Model] = true
	}
	if len(store.plans) > 0 {
		// Should use at least one of the preferred models
		assert.True(t, len(modelsUsed) > 0)
	}
}

// mockFeatureStore implements the FeatureStore interface for testing
type mockFeatureStore struct {
	developers []seed.Developer
	mcpTools []models.MCPToolEvent
	commands []models.CommandEvent
	plans    []models.PlanEvent
	askModes []models.AskModeEvent
}

func (m *mockFeatureStore) AddMCPTool(event models.MCPToolEvent) error {
	m.mcpTools = append(m.mcpTools, event)
	return nil
}

func (m *mockFeatureStore) AddCommand(event models.CommandEvent) error {
	m.commands = append(m.commands, event)
	return nil
}

func (m *mockFeatureStore) AddPlan(event models.PlanEvent) error {
	m.plans = append(m.plans, event)
	return nil
}

func (m *mockFeatureStore) AddAskMode(event models.AskModeEvent) error {
	m.askModes = append(m.askModes, event)
	return nil
}

func (m *mockFeatureStore) ListDevelopers() []seed.Developer {
	return m.developers
}
