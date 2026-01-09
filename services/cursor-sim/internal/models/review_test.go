package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestReview_Validate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		review  Review
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid review with approved state",
			review: Review{
				ID:          1,
				PRID:        42,
				Reviewer:    "alice@example.com",
				State:       ReviewStateApproved,
				SubmittedAt: now,
				Body:        "LGTM!",
			},
			wantErr: false,
		},
		{
			name: "valid review with changes requested",
			review: Review{
				ID:          2,
				PRID:        43,
				Reviewer:    "bob@example.com",
				State:       ReviewStateChangesRequested,
				SubmittedAt: now,
				Body:        "Please fix the tests",
			},
			wantErr: false,
		},
		{
			name: "valid review with pending state",
			review: Review{
				ID:          3,
				PRID:        44,
				Reviewer:    "charlie@example.com",
				State:       ReviewStatePending,
				SubmittedAt: now,
			},
			wantErr: false,
		},
		{
			name: "invalid pr_id - zero",
			review: Review{
				PRID:        0,
				Reviewer:    "alice@example.com",
				State:       ReviewStateApproved,
				SubmittedAt: now,
			},
			wantErr: true,
			errMsg:  "pr_id must be positive",
		},
		{
			name: "invalid pr_id - negative",
			review: Review{
				PRID:        -1,
				Reviewer:    "alice@example.com",
				State:       ReviewStateApproved,
				SubmittedAt: now,
			},
			wantErr: true,
			errMsg:  "pr_id must be positive",
		},
		{
			name: "missing reviewer",
			review: Review{
				PRID:        42,
				Reviewer:    "",
				State:       ReviewStateApproved,
				SubmittedAt: now,
			},
			wantErr: true,
			errMsg:  "reviewer is required",
		},
		{
			name: "invalid state",
			review: Review{
				PRID:        42,
				Reviewer:    "alice@example.com",
				State:       "invalid",
				SubmittedAt: now,
			},
			wantErr: true,
			errMsg:  "invalid state: invalid",
		},
		{
			name: "missing submitted_at",
			review: Review{
				PRID:     42,
				Reviewer: "alice@example.com",
				State:    ReviewStateApproved,
			},
			wantErr: true,
			errMsg:  "submitted_at is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.review.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestReview_IsApproval(t *testing.T) {
	tests := []struct {
		name  string
		state ReviewState
		want  bool
	}{
		{
			name:  "approved state returns true",
			state: ReviewStateApproved,
			want:  true,
		},
		{
			name:  "changes_requested returns false",
			state: ReviewStateChangesRequested,
			want:  false,
		},
		{
			name:  "pending returns false",
			state: ReviewStatePending,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			review := Review{State: tt.state}
			if got := review.IsApproval(); got != tt.want {
				t.Errorf("IsApproval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReview_CommentCount(t *testing.T) {
	tests := []struct {
		name     string
		comments []ReviewComment
		want     int
	}{
		{
			name:     "no comments",
			comments: nil,
			want:     0,
		},
		{
			name:     "empty comments array",
			comments: []ReviewComment{},
			want:     0,
		},
		{
			name: "one comment",
			comments: []ReviewComment{
				{ID: 1, Body: "Good work"},
			},
			want: 1,
		},
		{
			name: "multiple comments",
			comments: []ReviewComment{
				{ID: 1, Body: "Fix this"},
				{ID: 2, Body: "Also fix that"},
				{ID: 3, Body: "And this too"},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			review := Review{Comments: tt.comments}
			if got := review.CommentCount(); got != tt.want {
				t.Errorf("CommentCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReview_JSONMarshaling(t *testing.T) {
	now := time.Now().Truncate(time.Second) // Truncate for comparison

	tests := []struct {
		name   string
		review Review
	}{
		{
			name: "review with all fields",
			review: Review{
				ID:          1,
				PRID:        42,
				Reviewer:    "alice@example.com",
				State:       ReviewStateApproved,
				SubmittedAt: now,
				Body:        "LGTM! Great work.",
				Comments: []ReviewComment{
					{
						ID:        101,
						PRNumber:  42,
						RepoName:  "test-repo",
						AuthorID:  "alice",
						Body:      "Consider refactoring this",
						Path:      "main.go",
						Line:      42,
						State:     ReviewStateApproved,
						CreatedAt: now,
					},
				},
			},
		},
		{
			name: "review with minimal fields",
			review: Review{
				ID:          2,
				PRID:        43,
				Reviewer:    "bob@example.com",
				State:       ReviewStatePending,
				SubmittedAt: now,
			},
		},
		{
			name: "review with changes requested",
			review: Review{
				ID:          3,
				PRID:        44,
				Reviewer:    "charlie@example.com",
				State:       ReviewStateChangesRequested,
				SubmittedAt: now,
				Body:        "Please address the following issues",
				Comments: []ReviewComment{
					{ID: 201, Body: "Fix error handling"},
					{ID: 202, Body: "Add tests"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.review)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			// Unmarshal back to struct
			var decoded Review
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Compare fields
			if decoded.ID != tt.review.ID {
				t.Errorf("ID mismatch: got %v, want %v", decoded.ID, tt.review.ID)
			}
			if decoded.PRID != tt.review.PRID {
				t.Errorf("PRID mismatch: got %v, want %v", decoded.PRID, tt.review.PRID)
			}
			if decoded.Reviewer != tt.review.Reviewer {
				t.Errorf("Reviewer mismatch: got %v, want %v", decoded.Reviewer, tt.review.Reviewer)
			}
			if decoded.State != tt.review.State {
				t.Errorf("State mismatch: got %v, want %v", decoded.State, tt.review.State)
			}
			if decoded.Body != tt.review.Body {
				t.Errorf("Body mismatch: got %v, want %v", decoded.Body, tt.review.Body)
			}
			if len(decoded.Comments) != len(tt.review.Comments) {
				t.Errorf("Comments length mismatch: got %v, want %v", len(decoded.Comments), len(tt.review.Comments))
			}

			// Validate the unmarshaled review
			if err := decoded.Validate(); err != nil {
				t.Errorf("Unmarshaled review failed validation: %v", err)
			}
		})
	}
}

func TestReviewState_Constants(t *testing.T) {
	// Verify ReviewState constants match expected values
	if ReviewStateApproved != "approved" {
		t.Errorf("ReviewStateApproved = %v, want 'approved'", ReviewStateApproved)
	}
	if ReviewStateChangesRequested != "changes_requested" {
		t.Errorf("ReviewStateChangesRequested = %v, want 'changes_requested'", ReviewStateChangesRequested)
	}
	if ReviewStatePending != "pending" {
		t.Errorf("ReviewStatePending = %v, want 'pending'", ReviewStatePending)
	}
}

func TestReview_WithComments(t *testing.T) {
	now := time.Now()

	review := Review{
		ID:          1,
		PRID:        42,
		Reviewer:    "alice@example.com",
		State:       ReviewStateChangesRequested,
		SubmittedAt: now,
		Body:        "Please address the following",
		Comments: []ReviewComment{
			{
				ID:        1,
				PRNumber:  42,
				AuthorID:  "alice",
				Body:      "This needs refactoring",
				Path:      "main.go",
				Line:      15,
				State:     ReviewStateChangesRequested,
				CreatedAt: now,
			},
			{
				ID:        2,
				PRNumber:  42,
				AuthorID:  "alice",
				Body:      "Add unit tests here",
				Path:      "service.go",
				Line:      28,
				State:     ReviewStateChangesRequested,
				CreatedAt: now,
			},
		},
	}

	if err := review.Validate(); err != nil {
		t.Errorf("Valid review failed validation: %v", err)
	}

	if review.CommentCount() != 2 {
		t.Errorf("CommentCount() = %v, want 2", review.CommentCount())
	}

	if review.IsApproval() {
		t.Errorf("IsApproval() = true for changes_requested state")
	}
}
