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

// RevertEvent tracks a revert commit that reverses a previous PR.
type RevertEvent struct {
	RepoName      string    `json:"repo_name"`
	PRNumber      int       `json:"pr_number"`
	RevertCommit  string    `json:"revert_commit"`  // SHA of the revert commit
	OriginalPR    int       `json:"original_pr"`    // PR number that was reverted
	MergedAt      time.Time `json:"merged_at"`      // When original PR was merged
	RevertedAt    time.Time `json:"reverted_at"`    // When revert commit was created
	DaysToRevert  float64   `json:"days_to_revert"` // Time between merge and revert
	RevertMessage string    `json:"revert_message"` // Commit message of revert
}

// RevertedPR represents a PR that was reverted in the response.
type RevertedPR struct {
	PRNumber     int     `json:"pr_number"`
	MergedAt     string  `json:"merged_at"`
	RevertedAt   string  `json:"reverted_at"`
	DaysToRevert float64 `json:"days_to_revert"`
}

// RevertAnalysis represents the aggregated revert metrics for a repository.
type RevertAnalysis struct {
	WindowDays       int          `json:"window_days"`
	TotalPRsMerged   int          `json:"total_prs_merged"`
	TotalPRsReverted int          `json:"total_prs_reverted"`
	RevertRate       float64      `json:"revert_rate"`
	RevertedPRs      []RevertedPR `json:"reverted_prs"`
}

// HotfixEvent tracks a hotfix PR that follows a merged PR.
type HotfixEvent struct {
	RepoName      string    `json:"repo_name"`
	OriginalPR    int       `json:"original_pr"`
	HotfixPR      int       `json:"hotfix_pr"`
	FilesInCommon []string  `json:"files_in_common"`
	MergedAt      time.Time `json:"merged_at"`      // When original PR was merged
	HotfixAt      time.Time `json:"hotfix_at"`      // When hotfix PR was merged
	HoursBetween  float64   `json:"hours_between"`
}

// HotfixPRInfo represents a hotfix PR in the response.
type HotfixPRInfo struct {
	OriginalPR   int      `json:"original_pr"`
	HotfixPR     int      `json:"hotfix_pr"`
	HoursBetween float64  `json:"hours_between"`
	FilesInCommon []string `json:"files_in_common"`
}

// HotfixAnalysis represents the aggregated hotfix metrics for a repository.
type HotfixAnalysis struct {
	WindowHours  int            `json:"window_hours"`
	TotalPRsMerged int          `json:"total_prs_merged"`
	PRsWithHotfix  int          `json:"prs_with_hotfix"`
	HotfixRate   float64        `json:"hotfix_rate"`
	HotfixPRs    []HotfixPRInfo `json:"hotfix_prs"`
}
