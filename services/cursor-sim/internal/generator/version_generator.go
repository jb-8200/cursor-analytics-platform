package generator

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// ClientVersionStore defines the interface for storing client version events and querying developers.
type ClientVersionStore interface {
	AddClientVersion(event models.ClientVersionEvent) error
	ListDevelopers() []seed.Developer
}

// VersionGenerator generates synthetic client version usage events.
type VersionGenerator struct {
	seed     *seed.SeedData
	store    ClientVersionStore
	velocity *VelocityConfig
	rng      *rand.Rand
	versions []string // Available client versions
}

// NewVersionGenerator creates a new client version generator with a random seed.
func NewVersionGenerator(seedData *seed.SeedData, store ClientVersionStore, velocity string) *VersionGenerator {
	return NewVersionGeneratorWithSeed(seedData, store, velocity, time.Now().UnixNano())
}

// NewVersionGeneratorWithSeed creates a new client version generator with a specific seed for reproducibility.
func NewVersionGeneratorWithSeed(seedData *seed.SeedData, store ClientVersionStore, velocity string, randSeed int64) *VersionGenerator {
	// Generate realistic Cursor client versions (0.42.x and 0.43.x series)
	versions := []string{
		"0.42.1", "0.42.2", "0.42.3", "0.42.4", "0.42.5",
		"0.43.0", "0.43.1", "0.43.2", "0.43.3",
	}

	return &VersionGenerator{
		seed:     seedData,
		store:    store,
		velocity: NewVelocityConfig(velocity),
		rng:      rand.New(rand.NewSource(randSeed)),
		versions: versions,
	}
}

// GenerateClientVersions generates client version events for the specified number of days.
// Each developer is assigned a client version at the start and may upgrade over time.
func (g *VersionGenerator) GenerateClientVersions(ctx context.Context, days int) error {
	startTime := time.Now().AddDate(0, 0, -days)

	// Query developers from storage (includes replicated developers)
	developers := g.store.ListDevelopers()

	// Assign each developer an initial version
	developerVersions := make(map[string]string)
	for _, dev := range developers {
		developerVersions[dev.UserID] = g.selectInitialVersion()
	}

	// Generate version events for each day
	for day := 0; day < days; day++ {
		currentDate := startTime.AddDate(0, 0, day)

		for _, dev := range developers {
			// Check context cancellation
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// Check if developer upgrades their client (5% chance per day)
			if g.rng.Float64() < 0.05 {
				developerVersions[dev.UserID] = g.selectNewerVersion(developerVersions[dev.UserID])
			}

			// Generate a version event for this developer on this day
			// Simulate the developer using the client once per day during working hours
			hour := dev.WorkingHoursBand.Start + g.rng.Intn(dev.WorkingHoursBand.End-dev.WorkingHoursBand.Start)
			timestamp := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				hour, g.rng.Intn(60), 0, 0, currentDate.Location())

			event := models.ClientVersionEvent{
				UserID:        dev.UserID,
				UserEmail:     dev.Email,
				ClientVersion: developerVersions[dev.UserID],
				Timestamp:     timestamp,
				EventDate:     currentDate.Format("2006-01-02"),
			}

			if err := g.store.AddClientVersion(event); err != nil {
				return err
			}
		}
	}

	return nil
}

// selectInitialVersion selects an initial client version for a developer.
// Weighted towards recent versions (70% get latest or near-latest).
func (g *VersionGenerator) selectInitialVersion() string {
	if g.rng.Float64() < 0.7 {
		// 70% get one of the latest 3 versions
		idx := len(g.versions) - 1 - g.rng.Intn(3)
		if idx < 0 {
			idx = 0
		}
		return g.versions[idx]
	}

	// 30% get a random older version
	return g.versions[g.rng.Intn(len(g.versions))]
}

// selectNewerVersion selects a newer version than the current one.
// Returns a version that's at least one version newer.
func (g *VersionGenerator) selectNewerVersion(current string) string {
	currentIdx := -1
	for i, v := range g.versions {
		if v == current {
			currentIdx = i
			break
		}
	}

	if currentIdx == -1 || currentIdx >= len(g.versions)-1 {
		// Already on latest or version not found
		return current
	}

	// Upgrade to a newer version (randomly between current+1 and latest)
	upgradeRange := len(g.versions) - currentIdx - 1
	if upgradeRange <= 0 {
		return current
	}

	newIdx := currentIdx + 1 + g.rng.Intn(upgradeRange)
	if newIdx >= len(g.versions) {
		newIdx = len(g.versions) - 1
	}

	return g.versions[newIdx]
}

// GetVersions returns the list of available client versions.
func (g *VersionGenerator) GetVersions() []string {
	return g.versions
}

// VersionInfo returns formatted version information for debugging.
func (g *VersionGenerator) VersionInfo() string {
	return fmt.Sprintf("Versions: %v (total: %d)", g.versions, len(g.versions))
}
