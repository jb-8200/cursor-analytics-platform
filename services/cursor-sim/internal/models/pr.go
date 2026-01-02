package models

import "time"

// PRState represents the state of a pull request.
type PRState string

const (
	PRStateOpen   PRState = "open"
	PRStateClosed PRState = "closed"
	PRStateMerged PRState = "merged"
)

// PullRequest represents a GitHub pull request with AI metrics.
// Field names use snake_case to match GitHub API format.
type PullRequest struct {
	Number      int      `json:"number"`
	Title       string   `json:"title"`
	Body        string   `json:"body"`
	State       PRState  `json:"state"`
	AuthorID    string   `json:"author_id"`
	AuthorEmail string   `json:"author_email"`
	AuthorName  string   `json:"author_name"`
	RepoName    string   `json:"repo_name"`
	BaseBranch  string   `json:"base_branch"`
	HeadBranch  string   `json:"head_branch"`
	Reviewers   []string `json:"reviewers"`
	Labels      []string `json:"labels"`

	// Line metrics
	Additions    int `json:"additions"`
	Deletions    int `json:"deletions"`
	ChangedFiles int `json:"changed_files"`
	CommitCount  int `json:"commit_count"`

	// AI metrics (aggregated from commits)
	AIRatio       float64 `json:"ai_ratio"`
	TabLines      int     `json:"tab_lines"`
	ComposerLines int     `json:"composer_lines"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	MergedAt  *time.Time `json:"merged_at,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`

	// Quality signals
	WasReverted bool `json:"was_reverted"`
	IsBugFix    bool `json:"is_bug_fix"`
}

// IsOpen returns true if the PR is in open state.
func (pr *PullRequest) IsOpen() bool {
	return pr.State == PRStateOpen
}

// IsMerged returns true if the PR has been merged.
func (pr *PullRequest) IsMerged() bool {
	return pr.State == PRStateMerged
}

// NetLines returns the net line change (additions - deletions).
func (pr *PullRequest) NetLines() int {
	return pr.Additions - pr.Deletions
}

// AILines returns the total AI-generated lines (tab + composer).
func (pr *PullRequest) AILines() int {
	return pr.TabLines + pr.ComposerLines
}

// ReviewState represents the state of a review.
type ReviewState string

const (
	ReviewStatePending          ReviewState = "pending"
	ReviewStateApproved         ReviewState = "approved"
	ReviewStateChangesRequested ReviewState = "changes_requested"
)

// ReviewComment represents a review comment on a pull request.
// Field names use snake_case to match GitHub API format.
type ReviewComment struct {
	ID        int         `json:"id"`
	PRNumber  int         `json:"pr_number"`
	RepoName  string      `json:"repo_name"`
	AuthorID  string      `json:"author_id"`
	Body      string      `json:"body"`
	Path      string      `json:"path,omitempty"`
	Line      int         `json:"line,omitempty"`
	State     ReviewState `json:"state"`
	CreatedAt time.Time   `json:"created_at"`
}

// IsApproval returns true if this review approves the PR.
func (rc *ReviewComment) IsApproval() bool {
	return rc.State == ReviewStateApproved
}

// Repository represents a GitHub repository.
// Field names use snake_case to match GitHub API format.
type Repository struct {
	Name            string    `json:"name"`
	Owner           string    `json:"owner"`
	Description     string    `json:"description"`
	PrimaryLanguage string    `json:"primary_language"`
	DefaultBranch   string    `json:"default_branch"`
	Teams           []string  `json:"teams"`
	CreatedAt       time.Time `json:"created_at"`
}

// FullName returns the full repository name (owner/repo).
func (r *Repository) FullName() string {
	return r.Name
}
