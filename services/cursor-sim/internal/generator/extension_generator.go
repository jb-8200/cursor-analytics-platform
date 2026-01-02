package generator

import (
	"context"
	"math/rand"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// FileExtensionStore defines the interface for storing file extension events.
type FileExtensionStore interface {
	AddFileExtension(event models.FileExtensionEvent) error
}

// ExtensionGenerator generates synthetic file extension usage events.
type ExtensionGenerator struct {
	seed        *seed.SeedData
	store       FileExtensionStore
	velocity    *VelocityConfig
	rng         *rand.Rand
	extensions  []string // Popular file extensions
}

// NewExtensionGenerator creates a new file extension generator with a random seed.
func NewExtensionGenerator(seedData *seed.SeedData, store FileExtensionStore, velocity string) *ExtensionGenerator {
	return NewExtensionGeneratorWithSeed(seedData, store, velocity, time.Now().UnixNano())
}

// NewExtensionGeneratorWithSeed creates a new generator with a specific seed for reproducibility.
func NewExtensionGeneratorWithSeed(seedData *seed.SeedData, store FileExtensionStore, velocity string, randSeed int64) *ExtensionGenerator {
	// Popular file extensions developers work with
	extensions := []string{
		"tsx", "ts", "jsx", "js", "py", "go", "java", "rs", "cpp", "c",
		"rb", "php", "cs", "kt", "swift", "m", "scala", "r", "sql", "json",
	}

	return &ExtensionGenerator{
		seed:       seedData,
		store:      store,
		velocity:   NewVelocityConfig(velocity),
		rng:        rand.New(rand.NewSource(randSeed)),
		extensions: extensions,
	}
}

// GenerateFileExtensions generates file extension events for the specified number of days.
// Creates synthetic file edit suggestions with accept/reject outcomes.
func (g *ExtensionGenerator) GenerateFileExtensions(ctx context.Context, days int) error {
	startTime := time.Now().AddDate(0, 0, -days)

	for _, dev := range g.seed.Developers {
		if err := g.generateForDeveloper(ctx, dev, startTime, days); err != nil {
			return err
		}
	}

	return nil
}

// generateForDeveloper generates file extension events for a single developer.
func (g *ExtensionGenerator) generateForDeveloper(ctx context.Context, dev seed.Developer, startTime time.Time, days int) error {
	// Developer has 3-4 favorite extensions
	favExtensions := g.selectFavoriteExtensions(3)

	// Generate events for each day
	for day := 0; day < days; day++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		currentDate := startTime.AddDate(0, 0, day)

		// Generate 2-5 file extension events per developer per day
		eventCount := 2 + g.rng.Intn(4)
		for i := 0; i < eventCount; i++ {
			// Select extension (80% chance of favorite, 20% other)
			ext := g.selectExtension(favExtensions)

			// Generate realistic line counts
			linesSuggested := 10 + g.rng.Intn(191) // 10-200 lines
			linesAccepted := 0
			wasAccepted := g.rng.Float64() < 0.65 // 65% accept rate

			if wasAccepted {
				linesAccepted = int(float64(linesSuggested) * (0.7 + g.rng.Float64()*0.3)) // 70-100% accepted
			}
			linesRejected := linesSuggested - linesAccepted

			// Random time during working hours
			hour := dev.WorkingHoursBand.Start + g.rng.Intn(dev.WorkingHoursBand.End-dev.WorkingHoursBand.Start)
			timestamp := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				hour, g.rng.Intn(60), 0, 0, currentDate.Location())

			event := models.FileExtensionEvent{
				UserID:         dev.UserID,
				UserEmail:      dev.Email,
				FileExtension:  ext,
				LinesSuggested: linesSuggested,
				LinesAccepted:  linesAccepted,
				LinesRejected:  linesRejected,
				WasAccepted:    wasAccepted,
				Timestamp:      timestamp,
				EventDate:      currentDate.Format("2006-01-02"),
			}

			if err := g.store.AddFileExtension(event); err != nil {
				return err
			}
		}
	}

	return nil
}

// selectFavoriteExtensions picks n favorite extensions for a developer.
func (g *ExtensionGenerator) selectFavoriteExtensions(count int) []string {
	if count > len(g.extensions) {
		count = len(g.extensions)
	}

	// Randomly select favorite extensions
	favorites := make([]string, 0, count)
	indices := g.rng.Perm(len(g.extensions))
	for i := 0; i < count; i++ {
		favorites = append(favorites, g.extensions[indices[i]])
	}

	return favorites
}

// selectExtension picks an extension (biased towards favorites).
func (g *ExtensionGenerator) selectExtension(favorites []string) string {
	if len(favorites) == 0 {
		return g.extensions[g.rng.Intn(len(g.extensions))]
	}

	// 80% chance of favorite, 20% random
	if g.rng.Float64() < 0.8 {
		return favorites[g.rng.Intn(len(favorites))]
	}

	return g.extensions[g.rng.Intn(len(g.extensions))]
}

// GetExtensions returns the list of available file extensions.
func (g *ExtensionGenerator) GetExtensions() []string {
	return g.extensions
}
