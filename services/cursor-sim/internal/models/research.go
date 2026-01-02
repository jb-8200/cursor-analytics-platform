package models

import (
	"errors"
	"time"
)

// AIRatioBand categorizes AI ratio into low/medium/high bands for research analysis.
type AIRatioBand string

const (
	AIRatioBandLow    AIRatioBand = "low"    // 0.0 - 0.3
	AIRatioBandMedium AIRatioBand = "medium" // 0.3 - 0.7
	AIRatioBandHigh   AIRatioBand = "high"   // 0.7 - 1.0
)

// SurvivalWindow represents time windows for code survival analysis.
type SurvivalWindow int

const (
	SurvivalWindow30d SurvivalWindow = 30
	SurvivalWindow60d SurvivalWindow = 60
	SurvivalWindow90d SurvivalWindow = 90
)

// Days returns the number of days in this survival window.
func (sw SurvivalWindow) Days() int {
	return int(sw)
}

// ResearchDataPoint represents a pre-joined row for research export.
// Contains all relevant metrics from commits, PRs, and reviews for SDLC analysis.
// Uses snake_case JSON tags for compatibility with data science tools (pandas, R).
type ResearchDataPoint struct {
	// Identifiers
	CommitHash string `json:"commit_hash"`
	PRNumber   int    `json:"pr_number"`
	AuthorID   string `json:"author_id"`
	RepoName   string `json:"repo_name"`

	// AI Metrics (Independent Variables)
	AIRatio       float64 `json:"ai_ratio"`
	TabLines      int     `json:"tab_lines"`
	ComposerLines int     `json:"composer_lines"`

	// PR Metrics
	Additions    int `json:"additions"`
	Deletions    int `json:"deletions"`
	FilesChanged int `json:"files_changed"`

	// Cycle Times (Dependent Variables)
	CodingLeadTimeHours float64 `json:"coding_lead_time_hours"`
	ReviewLeadTimeHours float64 `json:"review_lead_time_hours"`
	MergeLeadTimeHours  float64 `json:"merge_lead_time_hours"`

	// Quality Outcomes (Dependent Variables)
	WasReverted      bool `json:"was_reverted"`
	RequiredHotfix   bool `json:"required_hotfix"`
	ReviewIterations int  `json:"review_iterations"`

	// Control Variables
	AuthorSeniority string `json:"author_seniority"`
	RepoMaturity    string `json:"repo_maturity"`
	IsGreenfield    bool   `json:"is_greenfield"`

	Timestamp time.Time `json:"timestamp"`
}

// GetAIRatioBand returns the AI ratio band for this data point.
func (dp *ResearchDataPoint) GetAIRatioBand() AIRatioBand {
	switch {
	case dp.AIRatio < 0.3:
		return AIRatioBandLow
	case dp.AIRatio < 0.7:
		return AIRatioBandMedium
	default:
		return AIRatioBandHigh
	}
}

// TotalLeadTimeHours returns the sum of all lead time components.
func (dp *ResearchDataPoint) TotalLeadTimeHours() float64 {
	return dp.CodingLeadTimeHours + dp.ReviewLeadTimeHours + dp.MergeLeadTimeHours
}

// Validate checks that the data point has valid values.
func (dp *ResearchDataPoint) Validate() error {
	if dp.CommitHash == "" {
		return errors.New("commit_hash is required")
	}
	if dp.AuthorID == "" {
		return errors.New("author_id is required")
	}
	if dp.RepoName == "" {
		return errors.New("repo_name is required")
	}
	if dp.AIRatio < 0 || dp.AIRatio > 1 {
		return errors.New("ai_ratio must be between 0 and 1")
	}
	if dp.Additions < 0 {
		return errors.New("additions cannot be negative")
	}
	if dp.Deletions < 0 {
		return errors.New("deletions cannot be negative")
	}
	return nil
}

// VelocityMetrics aggregates coding velocity metrics by AI ratio band.
type VelocityMetrics struct {
	Period              string      `json:"period"` // e.g., "2026-01" for monthly
	AIRatioBand         AIRatioBand `json:"ai_ratio_band"`
	TotalCommits        int         `json:"total_commits"`
	TotalPRs            int         `json:"total_prs"`
	TotalAdditions      int         `json:"total_additions"`
	TotalDeletions      int         `json:"total_deletions"`
	AvgCommitsPerDay    float64     `json:"avg_commits_per_day"`
	AvgPRsPerWeek       float64     `json:"avg_prs_per_week"`
	AvgLeadTimeHours    float64     `json:"avg_lead_time_hours"`
	MedianLeadTimeHours float64     `json:"median_lead_time_hours"`
	StdDevLeadTimeHours float64     `json:"std_dev_lead_time_hours"`
}

// ReviewCostMetrics aggregates review effort metrics by AI ratio band.
type ReviewCostMetrics struct {
	Period                string      `json:"period"`
	AIRatioBand           AIRatioBand `json:"ai_ratio_band"`
	TotalPRsReviewed      int         `json:"total_prs_reviewed"`
	TotalReviewComments   int         `json:"total_review_comments"`
	TotalReviewIterations int         `json:"total_review_iterations"`
	AvgCommentsPerPR      float64     `json:"avg_comments_per_pr"`
	AvgIterationsPerPR    float64     `json:"avg_iterations_per_pr"`
	AvgReviewTimeHours    float64     `json:"avg_review_time_hours"`
	MedianReviewTimeHours float64     `json:"median_review_time_hours"`
	StdDevReviewTimeHours float64     `json:"std_dev_review_time_hours"`
}

// QualityMetrics aggregates quality outcome metrics by AI ratio band.
type QualityMetrics struct {
	Period          string      `json:"period"`
	AIRatioBand     AIRatioBand `json:"ai_ratio_band"`
	TotalMergedPRs  int         `json:"total_merged_prs"`
	RevertedPRs     int         `json:"reverted_prs"`
	HotfixPRs       int         `json:"hotfix_prs"`
	RevertRate      float64     `json:"revert_rate"`
	HotfixRate      float64     `json:"hotfix_rate"`
	AvgTimeToRevert float64     `json:"avg_time_to_revert_hours"`
	AvgTimeToHotfix float64     `json:"avg_time_to_hotfix_hours"`
}

// CodeSurvivalRecord tracks survival of code blocks over time.
type CodeSurvivalRecord struct {
	CommitHash    string     `json:"commit_hash"`
	FilePath      string     `json:"file_path"`
	StartLine     int        `json:"start_line"`
	EndLine       int        `json:"end_line"`
	LinesAdded    int        `json:"lines_added"`
	AIRatio       float64    `json:"ai_ratio"`
	AuthorID      string     `json:"author_id"`
	AddedAt       time.Time  `json:"added_at"`
	SurvivedAt30d bool       `json:"survived_at_30d"`
	SurvivedAt60d bool       `json:"survived_at_60d"`
	SurvivedAt90d bool       `json:"survived_at_90d"`
	ModifiedCount int        `json:"modified_count"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

// SurvivalRate returns the survival rate (0.0 or 1.0) for the given window.
func (csr *CodeSurvivalRecord) SurvivalRate(window SurvivalWindow) float64 {
	switch window {
	case SurvivalWindow30d:
		if csr.SurvivedAt30d {
			return 1.0
		}
	case SurvivalWindow60d:
		if csr.SurvivedAt60d {
			return 1.0
		}
	case SurvivalWindow90d:
		if csr.SurvivedAt90d {
			return 1.0
		}
	}
	return 0.0
}
