package models

import "time"

// FileSurvival tracks the lifecycle of a file in a repository.
// It records when a file was created, modified, and optionally deleted,
// along with metrics about AI vs human contributions.
type FileSurvival struct {
	FilePath        string     `json:"file_path"`
	RepoName        string     `json:"repo_name"`
	CreatedAt       time.Time  `json:"created_at"`        // First commit timestamp
	LastModifiedAt  time.Time  `json:"last_modified_at"`
	AILinesAdded    int        `json:"ai_lines_added"`
	HumanLinesAdded int        `json:"human_lines_added"`
	TotalLines      int        `json:"total_lines"`
	RevertEvents    int        `json:"revert_events"`
	IsDeleted       bool       `json:"is_deleted"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// SurvivalAnalysis represents the aggregated survival metrics for a repository.
type SurvivalAnalysis struct {
	CohortStart     string               `json:"cohort_start"`
	CohortEnd       string               `json:"cohort_end"`
	ObservationDate string               `json:"observation_date"`
	TotalLinesAdded int                  `json:"total_lines_added"`
	LinesSurviving  int                  `json:"lines_surviving"`
	SurvivalRate    float64              `json:"survival_rate"`
	ByDeveloper     []DeveloperSurvival  `json:"by_developer"`
}

// DeveloperSurvival represents survival metrics for a single developer.
type DeveloperSurvival struct {
	Email          string  `json:"email"`
	LinesAdded     int     `json:"lines_added"`
	LinesSurviving int     `json:"lines_surviving"`
	SurvivalRate   float64 `json:"survival_rate"`
}
