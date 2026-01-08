package generator

import (
	"context"
	"math/rand"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// ModelStore defines the interface for storing model usage events and querying developers.
type ModelStore interface {
	AddModelUsage(usage models.ModelUsageEvent) error
	ListDevelopers() []seed.Developer
}

// ModelGenerator generates synthetic model usage events based on seed data.
type ModelGenerator struct {
	seed     *seed.SeedData
	store    ModelStore
	velocity *VelocityConfig
	rng      *rand.Rand
}

// NewModelGenerator creates a new model usage generator with a random seed.
func NewModelGenerator(seedData *seed.SeedData, store ModelStore, velocity string) *ModelGenerator {
	return NewModelGeneratorWithSeed(seedData, store, velocity, time.Now().UnixNano())
}

// NewModelGeneratorWithSeed creates a new model usage generator with a specific seed for reproducibility.
func NewModelGeneratorWithSeed(seedData *seed.SeedData, store ModelStore, velocity string, randSeed int64) *ModelGenerator {
	return &ModelGenerator{
		seed:     seedData,
		store:    store,
		velocity: NewVelocityConfig(velocity),
		rng:      rand.New(rand.NewSource(randSeed)),
	}
}

// GenerateModelUsage generates model usage events for the specified number of days.
// Uses Poisson process for event timing based on developer activity.
func (g *ModelGenerator) GenerateModelUsage(ctx context.Context, days int) error {
	startTime := time.Now().AddDate(0, 0, -days)

	// Query developers from storage (includes replicated developers)
	developers := g.store.ListDevelopers()

	for _, dev := range developers {
		if err := g.generateForDeveloper(ctx, dev, startTime); err != nil {
			return err
		}
	}

	return nil
}

// generateForDeveloper generates model usage events for a single developer.
func (g *ModelGenerator) generateForDeveloper(ctx context.Context, dev seed.Developer, startTime time.Time) error {
	// Calculate event rate based on PRs per week and velocity
	// Assume ~20 model queries per PR (mix of chat and code)
	// Apply velocity config (similar to commits per day calculation)
	adjustedQueriesPerDay := g.velocity.CommitsPerDay(dev.PRBehavior.PRsPerWeek) * (20.0 / 3.0) // Scale by query-to-commit ratio
	queriesPerHour := adjustedQueriesPerDay / 24.0

	current := startTime
	maxIterations := 100000 // Safety limit to prevent infinite loops
	iterations := 0

	for current.Before(time.Now()) && iterations < maxIterations {
		iterations++

		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Use exponential distribution for wait time (Poisson process)
		waitHours := g.exponential(1.0 / queriesPerHour)
		current = current.Add(time.Duration(waitHours * float64(time.Hour)))

		if current.After(time.Now()) {
			break
		}

		// Skip if outside working hours (basic check)
		if !g.isWorkingHour(current, dev) {
			continue
		}

		// Generate model usage event
		usage := g.generateEvent(dev, current)
		if err := g.store.AddModelUsage(usage); err != nil {
			return err
		}
	}

	return nil
}

// generateEvent creates a single model usage event for a developer.
func (g *ModelGenerator) generateEvent(dev seed.Developer, timestamp time.Time) models.ModelUsageEvent {
	// Select model from developer's preferences
	modelName := g.selectModel(dev)

	// Determine usage type (chat vs code) based on developer's ratio
	usageType := g.selectUsageType(dev)

	return models.ModelUsageEvent{
		UserID:    dev.UserID,
		UserEmail: dev.Email,
		ModelName: modelName,
		UsageType: usageType,
		Timestamp: timestamp,
		EventDate: timestamp.Format("2006-01-02"),
	}
}

// selectModel randomly selects a model from the developer's preferred models.
func (g *ModelGenerator) selectModel(dev seed.Developer) string {
	if len(dev.PreferredModels) == 0 {
		// Default models if none specified
		return "gpt-4-turbo"
	}

	// Weighted selection - first model has higher probability
	if g.rng.Float64() < 0.7 && len(dev.PreferredModels) > 0 {
		return dev.PreferredModels[0]
	}

	// Otherwise, pick randomly from all preferred models
	idx := g.rng.Intn(len(dev.PreferredModels))
	return dev.PreferredModels[idx]
}

// selectUsageType determines if this is a chat or code usage based on developer's ratio.
func (g *ModelGenerator) selectUsageType(dev seed.Developer) string {
	// Default to code-heavy if not specified
	chatRatio := 0.15
	if dev.ChatVsCodeRatio.Chat > 0 {
		chatRatio = dev.ChatVsCodeRatio.Chat
	}

	if g.rng.Float64() < chatRatio {
		return "chat"
	}
	return "code"
}

// isWorkingHour checks if the timestamp falls within developer's working hours.
// Simplified check - only considers hour of day, not timezone.
func (g *ModelGenerator) isWorkingHour(t time.Time, dev seed.Developer) bool {
	hour := t.Hour()

	// Default working hours if not specified
	start := 9
	end := 18

	if dev.WorkingHoursBand.Start > 0 {
		start = dev.WorkingHoursBand.Start
	}
	if dev.WorkingHoursBand.End > 0 {
		end = dev.WorkingHoursBand.End
	}

	return hour >= start && hour < end
}

// exponential generates a random number from exponential distribution.
func (g *ModelGenerator) exponential(lambda float64) float64 {
	return -1.0 / lambda * (g.rng.Float64() + 1e-10) // Add small epsilon to avoid log(0)
}
