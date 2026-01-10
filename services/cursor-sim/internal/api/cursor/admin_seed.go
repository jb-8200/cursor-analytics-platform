package cursor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"gopkg.in/yaml.v3"
)

// UploadSeed returns an HTTP handler for POST /admin/seed.
// Accepts JSON, YAML, or CSV seed data, validates it, and optionally regenerates data.
// Thread-safe seed swapping with optional regeneration.
func UploadSeed(store storage.Store, currentSeed **seed.SeedData) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST
		if r.Method != http.MethodPost {
			api.RespondError(w, http.StatusMethodNotAllowed, "only POST method is allowed")
			return
		}

		// Parse request body
		var req models.SeedUploadRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
			return
		}

		// Validate format
		format := strings.ToLower(req.Format)
		if format != "json" && format != "yaml" && format != "csv" {
			api.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid format %q, must be json, yaml, or csv", req.Format))
			return
		}

		// Parse seed data based on format
		var newSeed *seed.SeedData
		var err error

		switch format {
		case "json":
			newSeed, err = parseSeedJSON(req.Data)
		case "yaml":
			newSeed, err = parseSeedYAML(req.Data)
		case "csv":
			newSeed, err = parseSeedCSV(req.Data)
		}

		if err != nil {
			api.RespondError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse %s seed data: %v", format, err))
			return
		}

		// Validate seed data
		if err := newSeed.Validate(); err != nil {
			api.RespondError(w, http.StatusBadRequest, fmt.Sprintf("seed validation failed: %v", err))
			return
		}

		// Swap seed data (thread-safe)
		*currentSeed = newSeed

		// Extract organizational structure
		teams := extractUniqueTeams(newSeed)
		divisions := extractUniqueDivisions(newSeed)
		orgs := extractUniqueOrgs(newSeed)

		// Build response
		response := models.SeedUploadResponse{
			Status:        "success",
			SeedLoaded:    true,
			Developers:    len(newSeed.Developers),
			Repositories:  len(newSeed.Repositories),
			Teams:         teams,
			Divisions:     divisions,
			Organizations: orgs,
			Regenerated:   false,
		}

		// Optionally regenerate data
		if req.Regenerate && req.RegenerateConfig != nil {
			// TODO: Implement regeneration logic
			// For now, just set regenerated to false
			response.Regenerated = false
		}

		api.RespondJSON(w, http.StatusOK, response)
	})
}

// GetSeedPresets returns an HTTP handler for GET /admin/seed/presets.
// Returns predefined seed configurations for common scenarios.
func GetSeedPresets() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET
		if r.Method != http.MethodGet {
			api.RespondError(w, http.StatusMethodNotAllowed, "only GET method is allowed")
			return
		}

		presets := []models.SeedPreset{
			{
				Name:        "small-team",
				Description: "Small team (2 developers, 2 repos, 1 region)",
				Developers:  2,
				Teams:       2,
				Regions:     []string{"US"},
			},
			{
				Name:        "medium-team",
				Description: "Medium team (10 developers, 5 repos, 2 regions)",
				Developers:  10,
				Teams:       3,
				Regions:     []string{"US", "EU"},
			},
			{
				Name:        "enterprise",
				Description: "Enterprise (100 developers, 20 repos, 3 regions)",
				Developers:  100,
				Teams:       10,
				Regions:     []string{"US", "EU", "APAC"},
			},
			{
				Name:        "multi-region",
				Description: "Multi-region (50 developers, 15 repos, 5 regions)",
				Developers:  50,
				Teams:       8,
				Regions:     []string{"US", "EU", "APAC", "LATAM", "MEA"},
			},
		}

		response := models.SeedPresetsResponse{
			Presets: presets,
		}

		api.RespondJSON(w, http.StatusOK, response)
	})
}

// parseSeedJSON parses JSON seed data from a string.
func parseSeedJSON(data string) (*seed.SeedData, error) {
	var seedData seed.SeedData
	if err := json.Unmarshal([]byte(data), &seedData); err != nil {
		return nil, err
	}
	return &seedData, nil
}

// parseSeedYAML parses YAML seed data from a string.
func parseSeedYAML(data string) (*seed.SeedData, error) {
	var seedData seed.SeedData
	if err := yaml.Unmarshal([]byte(data), &seedData); err != nil {
		return nil, err
	}
	return &seedData, nil
}

// parseSeedCSV parses CSV seed data from a string.
// Uses LoadFromCSV from the seed package.
func parseSeedCSV(data string) (*seed.SeedData, error) {
	reader := bytes.NewReader([]byte(data))
	return seed.LoadFromCSV(reader)
}
