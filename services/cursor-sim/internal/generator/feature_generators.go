package generator

import (
	"context"
	"math/rand"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// FeatureStore defines the interface for storing feature events and querying developers.
type FeatureStore interface {
	AddMCPTool(event models.MCPToolEvent) error
	AddCommand(event models.CommandEvent) error
	AddPlan(event models.PlanEvent) error
	AddAskMode(event models.AskModeEvent) error
	ListDevelopers() []seed.Developer
}

// FeatureGenerator generates synthetic feature usage events (MCP, Commands, Plans, Ask Mode).
type FeatureGenerator struct {
	seed      *seed.SeedData
	store     FeatureStore
	velocity  *VelocityConfig
	rng       *rand.Rand
	mcpTools  []string
	commands  []string
	askModes  []string
}

// NewFeatureGenerator creates a new feature generator with a random seed.
func NewFeatureGenerator(seedData *seed.SeedData, store FeatureStore, velocity string) *FeatureGenerator {
	return NewFeatureGeneratorWithSeed(seedData, store, velocity, time.Now().UnixNano())
}

// NewFeatureGeneratorWithSeed creates a new generator with a specific seed for reproducibility.
func NewFeatureGeneratorWithSeed(seedData *seed.SeedData, store FeatureStore, velocity string, randSeed int64) *FeatureGenerator {
	return &FeatureGenerator{
		seed:     seedData,
		store:    store,
		velocity: NewVelocityConfig(velocity),
		rng:      rand.New(rand.NewSource(randSeed)),
		// Popular MCP tools
		mcpTools: []string{
			"read_file", "write_file", "list_files", "search_web",
			"execute_command", "create_directory", "delete_file", "grep",
		},
		// Cursor commands
		commands: []string{
			"explain", "refactor", "fix", "test", "generate", "review",
			"optimize", "document", "debug", "suggest", "migrate",
		},
		// Ask modes (same as plan - uses model preferences)
		askModes: []string{"chat", "think", "analyze"},
	}
}

// GenerateFeatures generates all feature events for the specified number of days.
func (g *FeatureGenerator) GenerateFeatures(ctx context.Context, days int) error {
	startTime := time.Now().AddDate(0, 0, -days)

	// Query developers from storage (includes replicated developers)
	developers := g.store.ListDevelopers()

	for _, dev := range developers {
		if err := g.generateMCPTools(ctx, dev, startTime, days); err != nil {
			return err
		}
		if err := g.generateCommands(ctx, dev, startTime, days); err != nil {
			return err
		}
		if err := g.generatePlans(ctx, dev, startTime, days); err != nil {
			return err
		}
		if err := g.generateAskMode(ctx, dev, startTime, days); err != nil {
			return err
		}
	}

	return nil
}

// generateMCPTools generates MCP tool usage events.
func (g *FeatureGenerator) generateMCPTools(ctx context.Context, dev seed.Developer, startTime time.Time, days int) error {
	for day := 0; day < days; day++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		currentDate := startTime.AddDate(0, 0, day)

		// 1-3 MCP tool usages per developer per day
		eventCount := 1 + g.rng.Intn(3)
		for i := 0; i < eventCount; i++ {
			toolIdx := g.rng.Intn(len(g.mcpTools))
			tool := g.mcpTools[toolIdx]

			// MCP server names based on tool
			serverName := "filesystem"
			if toolIdx > 2 {
				serverName = "web_search"
			}

			hour := dev.WorkingHoursBand.Start + g.rng.Intn(dev.WorkingHoursBand.End-dev.WorkingHoursBand.Start)
			timestamp := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				hour, g.rng.Intn(60), 0, 0, currentDate.Location())

			event := models.MCPToolEvent{
				UserID:        dev.UserID,
				UserEmail:     dev.Email,
				ToolName:      tool,
				MCPServerName: serverName,
				Timestamp:     timestamp,
				EventDate:     currentDate.Format("2006-01-02"),
			}

			if err := g.store.AddMCPTool(event); err != nil {
				return err
			}
		}
	}

	return nil
}

// generateCommands generates command usage events.
func (g *FeatureGenerator) generateCommands(ctx context.Context, dev seed.Developer, startTime time.Time, days int) error {
	for day := 0; day < days; day++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		currentDate := startTime.AddDate(0, 0, day)

		// 2-5 command usages per developer per day
		eventCount := 2 + g.rng.Intn(4)
		for i := 0; i < eventCount; i++ {
			command := g.commands[g.rng.Intn(len(g.commands))]

			hour := dev.WorkingHoursBand.Start + g.rng.Intn(dev.WorkingHoursBand.End-dev.WorkingHoursBand.Start)
			timestamp := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				hour, g.rng.Intn(60), 0, 0, currentDate.Location())

			event := models.CommandEvent{
				UserID:      dev.UserID,
				UserEmail:   dev.Email,
				CommandName: command,
				Timestamp:   timestamp,
				EventDate:   currentDate.Format("2006-01-02"),
			}

			if err := g.store.AddCommand(event); err != nil {
				return err
			}
		}
	}

	return nil
}

// generatePlans generates plan usage events.
func (g *FeatureGenerator) generatePlans(ctx context.Context, dev seed.Developer, startTime time.Time, days int) error {
	for day := 0; day < days; day++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		currentDate := startTime.AddDate(0, 0, day)

		// 0-2 plan usages per developer per day (less frequent)
		eventCount := g.rng.Intn(3)
		for i := 0; i < eventCount; i++ {
			// Use developer's preferred model, or default to claude-sonnet
			model := "claude-sonnet-4.5"
			if len(dev.PreferredModels) > 0 {
				model = dev.PreferredModels[g.rng.Intn(len(dev.PreferredModels))]
			}

			hour := dev.WorkingHoursBand.Start + g.rng.Intn(dev.WorkingHoursBand.End-dev.WorkingHoursBand.Start)
			timestamp := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				hour, g.rng.Intn(60), 0, 0, currentDate.Location())

			event := models.PlanEvent{
				UserID:    dev.UserID,
				UserEmail: dev.Email,
				Model:     model,
				Timestamp: timestamp,
				EventDate: currentDate.Format("2006-01-02"),
			}

			if err := g.store.AddPlan(event); err != nil {
				return err
			}
		}
	}

	return nil
}

// generateAskMode generates ask mode usage events.
func (g *FeatureGenerator) generateAskMode(ctx context.Context, dev seed.Developer, startTime time.Time, days int) error {
	for day := 0; day < days; day++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		currentDate := startTime.AddDate(0, 0, day)

		// 1-3 ask mode usages per developer per day
		eventCount := 1 + g.rng.Intn(3)
		for i := 0; i < eventCount; i++ {
			// Use developer's preferred model
			model := "claude-sonnet-4.5"
			if len(dev.PreferredModels) > 0 {
				model = dev.PreferredModels[g.rng.Intn(len(dev.PreferredModels))]
			}

			hour := dev.WorkingHoursBand.Start + g.rng.Intn(dev.WorkingHoursBand.End-dev.WorkingHoursBand.Start)
			timestamp := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				hour, g.rng.Intn(60), 0, 0, currentDate.Location())

			event := models.AskModeEvent{
				UserID:    dev.UserID,
				UserEmail: dev.Email,
				Model:     model,
				Timestamp: timestamp,
				EventDate: currentDate.Format("2006-01-02"),
			}

			if err := g.store.AddAskMode(event); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetMCPTools returns available MCP tools.
func (g *FeatureGenerator) GetMCPTools() []string {
	return g.mcpTools
}

// GetCommands returns available commands.
func (g *FeatureGenerator) GetCommands() []string {
	return g.commands
}
