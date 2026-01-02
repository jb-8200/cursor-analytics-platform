package models

import "time"

// Commit represents a git commit with AI code attribution.
// Field names use camelCase to match the Cursor AI Code Tracking API.
type Commit struct {
	CommitHash           string    `json:"commitHash"`
	UserID               string    `json:"userId"`
	UserEmail            string    `json:"userEmail"`
	UserName             string    `json:"userName"`
	RepoName             string    `json:"repoName"`
	BranchName           string    `json:"branchName"`
	IsPrimaryBranch      bool      `json:"isPrimaryBranch"`
	TotalLinesAdded      int       `json:"totalLinesAdded"`
	TotalLinesDeleted    int       `json:"totalLinesDeleted"`
	TabLinesAdded        int       `json:"tabLinesAdded"`
	TabLinesDeleted      int       `json:"tabLinesDeleted"`
	ComposerLinesAdded   int       `json:"composerLinesAdded"`
	ComposerLinesDeleted int       `json:"composerLinesDeleted"`
	NonAILinesAdded      int       `json:"nonAiLinesAdded"`
	NonAILinesDeleted    int       `json:"nonAiLinesDeleted"`
	Message              string    `json:"message"`
	CommitTs             time.Time `json:"commitTs"`
	CreatedAt            time.Time `json:"createdAt"`
}

// AIRatio calculates the ratio of AI-generated lines to total lines added.
// Returns a value between 0.0 and 1.0.
func (c *Commit) AIRatio() float64 {
	if c.TotalLinesAdded == 0 {
		return 0.0
	}
	aiLines := c.TabLinesAdded + c.ComposerLinesAdded
	return float64(aiLines) / float64(c.TotalLinesAdded)
}

// NetLines returns the net line change (added - deleted).
func (c *Commit) NetLines() int {
	return c.TotalLinesAdded - c.TotalLinesDeleted
}

// HasAIContent returns true if the commit contains any AI-generated code.
func (c *Commit) HasAIContent() bool {
	return c.TabLinesAdded > 0 || c.ComposerLinesAdded > 0
}
