package models

import (
	"fmt"
	"time"
)

// IssueState represents the state of an issue.
type IssueState string

const (
	IssueStateOpen   IssueState = "open"
	IssueStateClosed IssueState = "closed"
)

// Issue represents a GitHub issue with PR linkage.
// Field names use snake_case to match GitHub API format.
type Issue struct {
	Number       int        `json:"number"`
	Title        string     `json:"title"`
	Body         string     `json:"body,omitempty"`
	State        IssueState `json:"state"`
	AuthorID     string     `json:"author_id"`
	RepoName     string     `json:"repo_name"`
	Labels       []string   `json:"labels,omitempty"`
	Assignees    []string   `json:"assignees,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at,omitempty"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
	ClosedByPRID *int       `json:"closed_by_pr_id,omitempty"` // PR that closed this issue
}

// Validate checks that all required fields are present and valid.
func (i *Issue) Validate() error {
	if i.Number <= 0 {
		return fmt.Errorf("number must be positive")
	}
	if i.Title == "" {
		return fmt.Errorf("title is required")
	}
	if i.State != IssueStateOpen && i.State != IssueStateClosed {
		return fmt.Errorf("invalid state: %s", i.State)
	}
	if i.AuthorID == "" {
		return fmt.Errorf("author_id is required")
	}
	if i.RepoName == "" {
		return fmt.Errorf("repo_name is required")
	}
	if i.CreatedAt.IsZero() {
		return fmt.Errorf("created_at is required")
	}
	return nil
}

// IsOpen returns true if the issue is in open state.
func (i *Issue) IsOpen() bool {
	return i.State == IssueStateOpen
}

// IsClosed returns true if the issue is in closed state.
func (i *Issue) IsClosed() bool {
	return i.State == IssueStateClosed
}

// WasClosedByPR returns true if the issue was closed by a pull request.
func (i *Issue) WasClosedByPR() bool {
	return i.ClosedByPRID != nil
}
