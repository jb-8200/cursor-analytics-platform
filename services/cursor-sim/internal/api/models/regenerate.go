package models

// RegenerateRequest represents the request body for POST /admin/regenerate.
// Allows runtime reconfiguration of data generation parameters.
type RegenerateRequest struct {
	Mode        string `json:"mode"`         // "append" or "override"
	Days        int    `json:"days"`         // Days of history to generate (1-3650)
	Velocity    string `json:"velocity"`     // Event generation rate: "low", "medium", "high"
	Developers  int    `json:"developers"`   // Number of developers (0-10000, 0 = use seed count)
	MaxCommits  int    `json:"max_commits"`  // Maximum commits per developer (0-100000, 0 = unlimited)
}

// RegenerateResponse represents the response from POST /admin/regenerate.
// Reports the results of data regeneration operation.
type RegenerateResponse struct {
	Status          string       `json:"status"`           // "success" or "error"
	Mode            string       `json:"mode"`             // "append" or "override"
	DataCleaned     bool         `json:"data_cleaned"`     // Whether data was cleared (override mode)
	CommitsAdded    int          `json:"commits_added"`    // Number of commits added
	PRsAdded        int          `json:"prs_added"`        // Number of PRs added
	ReviewsAdded    int          `json:"reviews_added"`    // Number of reviews added
	IssuesAdded     int          `json:"issues_added"`     // Number of issues added
	TotalCommits    int          `json:"total_commits"`    // Total commits after operation
	TotalPRs        int          `json:"total_prs"`        // Total PRs after operation
	TotalDevelopers int          `json:"total_developers"` // Total developers
	Duration        string       `json:"duration"`         // Time taken for operation
	Config          ConfigParams `json:"config"`           // Configuration used
}

// ConfigParams represents the generation configuration used.
type ConfigParams struct {
	Days       int    `json:"days"`        // Days of history generated
	Velocity   string `json:"velocity"`    // Event generation rate
	Developers int    `json:"developers"`  // Number of developers
	MaxCommits int    `json:"max_commits"` // Max commits per developer
}
