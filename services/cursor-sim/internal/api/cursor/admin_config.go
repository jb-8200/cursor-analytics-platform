package cursor

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// serverStartTime tracks when the server was started for uptime calculation.
var serverStartTime = time.Now()

// GetConfig returns a handler for GET /admin/config.
// It inspects the current runtime configuration including generation parameters,
// seed structure, external data sources, and server information.
func GetConfig(cfg *config.Config, seedData *seed.SeedData, version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET
		if r.Method != http.MethodGet {
			api.RespondError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
			return
		}

		// Build response
		response := models.ConfigResponse{
			Generation: models.GenerationConfig{
				Days:       cfg.Days,
				Velocity:   cfg.Velocity,
				Developers: len(seedData.Developers),
				MaxCommits: cfg.GenParams.MaxCommits,
			},
			Seed: buildSeedConfig(seedData),
			ExternalSources: buildExternalSourcesConfig(seedData),
			Server: models.ServerConfig{
				Port:    cfg.Port,
				Version: version,
				Uptime:  formatUptime(time.Since(serverStartTime)),
			},
		}

		api.RespondJSON(w, http.StatusOK, response)
	})
}

// buildSeedConfig extracts seed configuration from seed data.
func buildSeedConfig(seedData *seed.SeedData) models.SeedConfig {
	return models.SeedConfig{
		Version:       seedData.Version,
		Developers:    len(seedData.Developers),
		Repositories:  len(seedData.Repositories),
		Organizations: extractUniqueOrgs(seedData),
		Divisions:     extractUniqueDivisions(seedData),
		Teams:         extractUniqueTeams(seedData),
		Regions:       extractUniqueRegions(seedData),
		BySeniority:   groupBySeniority(seedData),
		ByRegion:      groupByRegion(seedData),
		ByTeam:        groupByTeam(seedData),
	}
}

// buildExternalSourcesConfig extracts external data source configuration.
func buildExternalSourcesConfig(seedData *seed.SeedData) models.ExternalSourcesConfig {
	config := models.ExternalSourcesConfig{
		Harvey:    models.HarveyConfig{Enabled: false},
		Copilot:   models.CopilotConfig{Enabled: false},
		Qualtrics: models.QualtricsConfig{Enabled: false},
	}

	if seedData.ExternalDataSources == nil {
		return config
	}

	// Harvey configuration
	if seedData.ExternalDataSources.Harvey != nil {
		config.Harvey.Enabled = seedData.ExternalDataSources.Harvey.Enabled
		config.Harvey.Models = seedData.ExternalDataSources.Harvey.ModelsUsed
	}

	// Copilot configuration
	if seedData.ExternalDataSources.Copilot != nil {
		config.Copilot.Enabled = seedData.ExternalDataSources.Copilot.Enabled
		config.Copilot.TotalLicenses = seedData.ExternalDataSources.Copilot.TotalLicenses
		config.Copilot.ActiveUsers = seedData.ExternalDataSources.Copilot.ActiveUsers
	}

	// Qualtrics configuration
	if seedData.ExternalDataSources.Qualtrics != nil {
		config.Qualtrics.Enabled = seedData.ExternalDataSources.Qualtrics.Enabled
		config.Qualtrics.SurveyID = seedData.ExternalDataSources.Qualtrics.SurveyID
		config.Qualtrics.ResponseCount = seedData.ExternalDataSources.Qualtrics.ResponseCount
	}

	return config
}

// formatUptime formats a duration as a human-readable uptime string.
func formatUptime(d time.Duration) string {
	if d < time.Minute {
		return d.Round(time.Second).String()
	}
	if d < time.Hour {
		return d.Round(time.Minute).String()
	}
	if d < 24*time.Hour {
		return d.Round(time.Minute).String()
	}

	// Format as "Xd Yh Zm" for uptime >= 1 day
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	var result string
	if days > 0 {
		result = fmt.Sprintf("%dd", days)
		if hours > 0 {
			result += fmt.Sprintf("%dh", hours)
		}
		if minutes > 0 {
			result += fmt.Sprintf("%dm", minutes)
		}
	} else if hours > 0 {
		result = fmt.Sprintf("%dh", hours)
		if minutes > 0 {
			result += fmt.Sprintf("%dm", minutes)
		}
	} else {
		result = fmt.Sprintf("%dm", minutes)
	}

	return result
}
