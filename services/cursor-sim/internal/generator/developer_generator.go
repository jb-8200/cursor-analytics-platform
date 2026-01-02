package generator

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/jb-8200/cursor-analytics-platform/services/cursor-sim/internal/config"
	"github.com/jb-8200/cursor-analytics-platform/services/cursor-sim/internal/models"
)

// Default organizational distributions from SPEC.md
var (
	// Region distribution: US: 50%, EU: 35%, APAC: 15%
	defaultRegions = map[string]float64{
		"US":   0.50,
		"EU":   0.35,
		"APAC": 0.15,
	}

	// Division distribution: AGS: 40%, AT: 35%, ST: 25%
	defaultDivisions = map[string]float64{
		"AGS": 0.40,
		"AT":  0.35,
		"ST":  0.25,
	}

	// Group distribution: TMOBILE: 60%, ATANT: 40%
	defaultGroups = map[string]float64{
		"TMOBILE": 0.60,
		"ATANT":   0.40,
	}

	// Team distribution: dev: 75%, support: 25%
	defaultTeams = map[string]float64{
		"dev":     0.75,
		"support": 0.25,
	}

	// Seniority distribution: junior: 20%, mid: 50%, senior: 30%
	defaultSeniority = map[string]float64{
		"junior": 0.20,
		"mid":    0.50,
		"senior": 0.30,
	}

	// Acceptance rate ranges by seniority
	acceptanceRates = map[string][2]float64{
		"junior": {0.55, 0.65}, // 55-65%
		"mid":    {0.70, 0.80}, // 70-80%
		"senior": {0.85, 0.95}, // 85-95%
	}

	// Client versions to randomly assign
	clientVersions = []string{
		"0.43.6", "0.43.5", "0.43.4", "0.43.3", "0.43.2", "0.43.1", "0.43.0",
		"0.42.9", "0.42.8", "0.42.7", "0.42.6", "0.42.5", "0.42.4",
	}
)

// DeveloperGenerator generates realistic developer profiles
type DeveloperGenerator struct {
	config        *config.Config
	rng           *rand.Rand
	nameGenerator *NameGenerator
}

// NewDeveloperGenerator creates a new developer generator
func NewDeveloperGenerator(cfg *config.Config) *DeveloperGenerator {
	return &DeveloperGenerator{
		config:        cfg,
		rng:           rand.New(rand.NewSource(cfg.Seed)),
		nameGenerator: NewNameGenerator(cfg.Seed),
	}
}

// Generate creates a slice of developers according to configuration
func (g *DeveloperGenerator) Generate() ([]*models.Developer, error) {
	count := g.config.Developers
	developers := make([]*models.Developer, 0, count)

	now := time.Now().UTC()

	// Track used emails to ensure uniqueness
	usedEmails := make(map[string]bool)

	for i := 0; i < count; i++ {
		// Generate unique name and email
		var firstName, lastName, email string
		for {
			firstName, lastName = g.nameGenerator.GenerateName()
			email = g.nameGenerator.GenerateEmail(firstName, lastName)

			// Check if email is unique, if not, generate another
			if !usedEmails[email] {
				usedEmails[email] = true
				break
			}
		}

		name := fmt.Sprintf("%s %s", firstName, lastName)

		// Generate organizational attributes based on distributions
		region := g.selectFromDistribution(defaultRegions)
		division := g.selectFromDistribution(defaultDivisions)
		group := g.selectFromDistribution(defaultGroups)
		team := g.selectFromDistribution(defaultTeams)
		seniority := g.selectFromDistribution(defaultSeniority)

		// Generate acceptance rate based on seniority
		acceptanceRate := g.generateAcceptanceRate(seniority)

		// Random client version
		clientVersion := clientVersions[g.rng.Intn(len(clientVersions))]

		// Create developer
		dev := &models.Developer{
			ID:             fmt.Sprintf("user_%08x", g.rng.Uint32()),
			Email:          email,
			Name:           name,
			Region:         region,
			Division:       division,
			Group:          group,
			Team:           team,
			Seniority:      seniority,
			ClientVersion:  clientVersion,
			AcceptanceRate: acceptanceRate,
			IsActive:       true,
			CreatedAt:      now,
			LastActiveAt:   now,
		}

		developers = append(developers, dev)
	}

	return developers, nil
}

// selectFromDistribution selects a key from a distribution map based on probabilities
// Uses deterministic ordering by sorting keys to ensure reproducibility
func (g *DeveloperGenerator) selectFromDistribution(dist map[string]float64) string {
	// Generate random number between 0 and 1
	r := g.rng.Float64()

	// Create sorted list of keys for deterministic ordering
	keys := make([]string, 0, len(dist))
	for key := range dist {
		keys = append(keys, key)
	}

	// Sort keys alphabetically for consistent ordering
	sort.Strings(keys)

	// Accumulate probabilities and select when we exceed r
	cumulative := 0.0
	for _, key := range keys {
		cumulative += dist[key]
		if r < cumulative {
			return key
		}
	}

	// Fallback to last key (handles floating point precision issues)
	if len(keys) > 0 {
		return keys[len(keys)-1]
	}

	return ""
}

// generateAcceptanceRate generates a rate within the range for the given seniority
func (g *DeveloperGenerator) generateAcceptanceRate(seniority string) float64 {
	rateRange, ok := acceptanceRates[seniority]
	if !ok {
		// Default to mid-range if seniority not found
		rateRange = acceptanceRates["mid"]
	}

	// Generate random rate within range
	min := rateRange[0]
	max := rateRange[1]
	return min + g.rng.Float64()*(max-min)
}
