package models

// SeedUploadRequest represents the request body for POST /admin/seed.
// Supports multiple formats (JSON, YAML, CSV) with optional regeneration.
type SeedUploadRequest struct {
	Data             string              `json:"data"`                        // Seed data as string (JSON/YAML/CSV content)
	Format           string              `json:"format"`                      // "json", "yaml", or "csv"
	Regenerate       bool                `json:"regenerate"`                  // Whether to regenerate data after upload
	RegenerateConfig *RegenerateRequest  `json:"regenerate_config,omitempty"` // Optional regeneration parameters
}

// SeedUploadResponse represents the response from POST /admin/seed.
// Reports the uploaded seed structure and optional regeneration results.
type SeedUploadResponse struct {
	Status         string               `json:"status"`                    // "success" or "error"
	SeedLoaded     bool                 `json:"seed_loaded"`               // Whether seed was successfully loaded
	Developers     int                  `json:"developers"`                // Number of developers in seed
	Repositories   int                  `json:"repositories"`              // Number of repositories in seed
	Teams          []string             `json:"teams"`                     // Unique teams
	Divisions      []string             `json:"divisions"`                 // Unique divisions
	Organizations  []string             `json:"organizations"`             // Unique organizations
	Regenerated    bool                 `json:"regenerated"`               // Whether data was regenerated
	GenerateStats  *RegenerateResponse  `json:"generate_stats,omitempty"`  // Stats from regeneration (if regenerated)
}

// SeedPreset represents a predefined seed configuration.
type SeedPreset struct {
	Name        string   `json:"name"`        // Preset identifier (e.g., "small-team")
	Description string   `json:"description"` // Human-readable description
	Developers  int      `json:"developers"`  // Number of developers
	Teams       int      `json:"teams"`       // Number of teams
	Regions     []string `json:"regions"`     // Geographic regions
}

// SeedPresetsResponse represents the response from GET /admin/seed/presets.
type SeedPresetsResponse struct {
	Presets []SeedPreset `json:"presets"` // List of available presets
}
