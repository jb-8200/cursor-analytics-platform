package models

// ConfigResponse represents the response structure for GET /admin/config.
// It provides comprehensive runtime configuration inspection including generation parameters,
// seed structure, external data sources, and server information.
type ConfigResponse struct {
	Generation      GenerationConfig      `json:"generation"`
	Seed            SeedConfig            `json:"seed"`
	ExternalSources ExternalSourcesConfig `json:"external_sources"`
	Server          ServerConfig          `json:"server"`
}

// GenerationConfig contains active generation parameters.
type GenerationConfig struct {
	Days       int    `json:"days"`
	Velocity   string `json:"velocity"`
	Developers int    `json:"developers"`
	MaxCommits int    `json:"max_commits"`
}

// SeedConfig contains seed data structure information.
type SeedConfig struct {
	Version       string         `json:"version"`
	Developers    int            `json:"developers"`
	Repositories  int            `json:"repositories"`
	Organizations []string       `json:"organizations"`
	Divisions     []string       `json:"divisions"`
	Teams         []string       `json:"teams"`
	Regions       []string       `json:"regions"`
	BySeniority   map[string]int `json:"by_seniority"`
	ByRegion      map[string]int `json:"by_region"`
	ByTeam        map[string]int `json:"by_team"`
}

// ExternalSourcesConfig contains external data source configurations.
type ExternalSourcesConfig struct {
	Harvey    HarveyConfig    `json:"harvey"`
	Copilot   CopilotConfig   `json:"copilot"`
	Qualtrics QualtricsConfig `json:"qualtrics"`
}

// HarveyConfig contains Harvey AI configuration.
type HarveyConfig struct {
	Enabled bool     `json:"enabled"`
	Models  []string `json:"models"`
}

// CopilotConfig contains Microsoft 365 Copilot configuration.
type CopilotConfig struct {
	Enabled       bool `json:"enabled"`
	TotalLicenses int  `json:"total_licenses"`
	ActiveUsers   int  `json:"active_users"`
}

// QualtricsConfig contains Qualtrics survey configuration.
type QualtricsConfig struct {
	Enabled       bool   `json:"enabled"`
	SurveyID      string `json:"survey_id"`
	ResponseCount int    `json:"response_count"`
}

// ServerConfig contains server runtime information.
type ServerConfig struct {
	Port    int    `json:"port"`
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
}
