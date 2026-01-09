package models

import (
	"fmt"
	"time"
)

// Review represents a code review on a pull request.
// Field names use snake_case to match GitHub API format.
type Review struct {
	ID          int              `json:"id"`
	PRID        int              `json:"pr_id"`
	Reviewer    string           `json:"reviewer"`      // Developer email
	State       ReviewState      `json:"state"`
	SubmittedAt time.Time        `json:"submitted_at"`
	Body        string           `json:"body,omitempty"`
	Comments    []ReviewComment  `json:"comments,omitempty"`
}

// Validate checks that all required fields are present and valid.
func (r *Review) Validate() error {
	if r.PRID <= 0 {
		return fmt.Errorf("pr_id must be positive")
	}
	if r.Reviewer == "" {
		return fmt.Errorf("reviewer is required")
	}
	if r.State != ReviewStateApproved && r.State != ReviewStateChangesRequested && r.State != ReviewStatePending {
		return fmt.Errorf("invalid state: %s", r.State)
	}
	if r.SubmittedAt.IsZero() {
		return fmt.Errorf("submitted_at is required")
	}
	return nil
}

// IsApproval returns true if this review approves the PR.
func (r *Review) IsApproval() bool {
	return r.State == ReviewStateApproved
}

// CommentCount returns the number of comments in this review.
func (r *Review) CommentCount() int {
	return len(r.Comments)
}
