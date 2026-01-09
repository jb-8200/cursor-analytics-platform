package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestIssueState_Constants(t *testing.T) {
	// Verify IssueState constants match expected values
	if IssueStateOpen != "open" {
		t.Errorf("IssueStateOpen = %v, want 'open'", IssueStateOpen)
	}
	if IssueStateClosed != "closed" {
		t.Errorf("IssueStateClosed = %v, want 'closed'", IssueStateClosed)
	}
}

func TestIssue_Validate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		issue   Issue
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid issue with open state",
			issue: Issue{
				Number:    1,
				Title:     "Bug in login flow",
				Body:      "Description of the bug",
				State:     IssueStateOpen,
				AuthorID:  "alice@example.com",
				RepoName:  "acme/platform",
				CreatedAt: now,
			},
			wantErr: false,
		},
		{
			name: "valid issue with closed state",
			issue: Issue{
				Number:       2,
				Title:        "Feature request",
				Body:         "Please add dark mode",
				State:        IssueStateClosed,
				AuthorID:     "bob@example.com",
				RepoName:     "acme/platform",
				CreatedAt:    now.Add(-24 * time.Hour),
				ClosedAt:     &now,
				ClosedByPRID: intPtr(42),
			},
			wantErr: false,
		},
		{
			name: "valid issue with labels and assignees",
			issue: Issue{
				Number:    3,
				Title:     "Add unit tests",
				Body:      "Need more test coverage",
				State:     IssueStateOpen,
				AuthorID:  "charlie@example.com",
				RepoName:  "acme/platform",
				Labels:    []string{"bug", "priority-high"},
				Assignees: []string{"alice@example.com", "bob@example.com"},
				CreatedAt: now,
			},
			wantErr: false,
		},
		{
			name: "invalid number - zero",
			issue: Issue{
				Number:    0,
				Title:     "Some issue",
				State:     IssueStateOpen,
				AuthorID:  "alice@example.com",
				RepoName:  "acme/platform",
				CreatedAt: now,
			},
			wantErr: true,
			errMsg:  "number must be positive",
		},
		{
			name: "invalid number - negative",
			issue: Issue{
				Number:    -1,
				Title:     "Some issue",
				State:     IssueStateOpen,
				AuthorID:  "alice@example.com",
				RepoName:  "acme/platform",
				CreatedAt: now,
			},
			wantErr: true,
			errMsg:  "number must be positive",
		},
		{
			name: "missing title",
			issue: Issue{
				Number:    1,
				Title:     "",
				State:     IssueStateOpen,
				AuthorID:  "alice@example.com",
				RepoName:  "acme/platform",
				CreatedAt: now,
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "invalid state",
			issue: Issue{
				Number:    1,
				Title:     "Some issue",
				State:     "invalid",
				AuthorID:  "alice@example.com",
				RepoName:  "acme/platform",
				CreatedAt: now,
			},
			wantErr: true,
			errMsg:  "invalid state: invalid",
		},
		{
			name: "missing author_id",
			issue: Issue{
				Number:    1,
				Title:     "Some issue",
				State:     IssueStateOpen,
				AuthorID:  "",
				RepoName:  "acme/platform",
				CreatedAt: now,
			},
			wantErr: true,
			errMsg:  "author_id is required",
		},
		{
			name: "missing repo_name",
			issue: Issue{
				Number:    1,
				Title:     "Some issue",
				State:     IssueStateOpen,
				AuthorID:  "alice@example.com",
				RepoName:  "",
				CreatedAt: now,
			},
			wantErr: true,
			errMsg:  "repo_name is required",
		},
		{
			name: "missing created_at",
			issue: Issue{
				Number:   1,
				Title:    "Some issue",
				State:    IssueStateOpen,
				AuthorID: "alice@example.com",
				RepoName: "acme/platform",
			},
			wantErr: true,
			errMsg:  "created_at is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.issue.Validate()
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

func TestIssue_IsOpen(t *testing.T) {
	tests := []struct {
		name  string
		state IssueState
		want  bool
	}{
		{
			name:  "open state returns true",
			state: IssueStateOpen,
			want:  true,
		},
		{
			name:  "closed state returns false",
			state: IssueStateClosed,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issue := Issue{State: tt.state}
			if got := issue.IsOpen(); got != tt.want {
				t.Errorf("IsOpen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssue_IsClosed(t *testing.T) {
	tests := []struct {
		name  string
		state IssueState
		want  bool
	}{
		{
			name:  "closed state returns true",
			state: IssueStateClosed,
			want:  true,
		},
		{
			name:  "open state returns false",
			state: IssueStateOpen,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issue := Issue{State: tt.state}
			if got := issue.IsClosed(); got != tt.want {
				t.Errorf("IsClosed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssue_WasClosedByPR(t *testing.T) {
	tests := []struct {
		name         string
		closedByPRID *int
		want         bool
	}{
		{
			name:         "has ClosedByPRID returns true",
			closedByPRID: intPtr(42),
			want:         true,
		},
		{
			name:         "nil ClosedByPRID returns false",
			closedByPRID: nil,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issue := Issue{ClosedByPRID: tt.closedByPRID}
			if got := issue.WasClosedByPR(); got != tt.want {
				t.Errorf("WasClosedByPR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssue_JSONMarshaling(t *testing.T) {
	now := time.Now().Truncate(time.Second) // Truncate for comparison
	closedAt := now.Add(time.Hour)

	tests := []struct {
		name  string
		issue Issue
	}{
		{
			name: "issue with all fields",
			issue: Issue{
				Number:       1,
				Title:        "Implement dark mode",
				Body:         "Add a dark mode toggle in settings",
				State:        IssueStateClosed,
				AuthorID:     "alice@example.com",
				RepoName:     "acme/platform",
				Labels:       []string{"enhancement", "ui"},
				Assignees:    []string{"bob@example.com"},
				CreatedAt:    now,
				UpdatedAt:    now.Add(30 * time.Minute),
				ClosedAt:     &closedAt,
				ClosedByPRID: intPtr(42),
			},
		},
		{
			name: "issue with minimal fields",
			issue: Issue{
				Number:    2,
				Title:     "Simple bug report",
				State:     IssueStateOpen,
				AuthorID:  "bob@example.com",
				RepoName:  "acme/platform",
				CreatedAt: now,
			},
		},
		{
			name: "issue with empty arrays",
			issue: Issue{
				Number:    3,
				Title:     "No labels or assignees",
				State:     IssueStateOpen,
				AuthorID:  "charlie@example.com",
				RepoName:  "acme/platform",
				Labels:    []string{},
				Assignees: []string{},
				CreatedAt: now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.issue)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			// Unmarshal back to struct
			var decoded Issue
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Compare fields
			if decoded.Number != tt.issue.Number {
				t.Errorf("Number mismatch: got %v, want %v", decoded.Number, tt.issue.Number)
			}
			if decoded.Title != tt.issue.Title {
				t.Errorf("Title mismatch: got %v, want %v", decoded.Title, tt.issue.Title)
			}
			if decoded.Body != tt.issue.Body {
				t.Errorf("Body mismatch: got %v, want %v", decoded.Body, tt.issue.Body)
			}
			if decoded.State != tt.issue.State {
				t.Errorf("State mismatch: got %v, want %v", decoded.State, tt.issue.State)
			}
			if decoded.AuthorID != tt.issue.AuthorID {
				t.Errorf("AuthorID mismatch: got %v, want %v", decoded.AuthorID, tt.issue.AuthorID)
			}
			if decoded.RepoName != tt.issue.RepoName {
				t.Errorf("RepoName mismatch: got %v, want %v", decoded.RepoName, tt.issue.RepoName)
			}
			if len(decoded.Labels) != len(tt.issue.Labels) {
				t.Errorf("Labels length mismatch: got %v, want %v", len(decoded.Labels), len(tt.issue.Labels))
			}
			if len(decoded.Assignees) != len(tt.issue.Assignees) {
				t.Errorf("Assignees length mismatch: got %v, want %v", len(decoded.Assignees), len(tt.issue.Assignees))
			}

			// Check ClosedByPRID pointer
			if tt.issue.ClosedByPRID != nil {
				if decoded.ClosedByPRID == nil {
					t.Errorf("ClosedByPRID is nil, want %v", *tt.issue.ClosedByPRID)
				} else if *decoded.ClosedByPRID != *tt.issue.ClosedByPRID {
					t.Errorf("ClosedByPRID mismatch: got %v, want %v", *decoded.ClosedByPRID, *tt.issue.ClosedByPRID)
				}
			} else {
				if decoded.ClosedByPRID != nil {
					t.Errorf("ClosedByPRID is %v, want nil", *decoded.ClosedByPRID)
				}
			}

			// Validate the unmarshaled issue
			if err := decoded.Validate(); err != nil {
				t.Errorf("Unmarshaled issue failed validation: %v", err)
			}
		})
	}
}

func TestIssue_WithLabels(t *testing.T) {
	now := time.Now()

	issue := Issue{
		Number:    1,
		Title:     "Add CI/CD pipeline",
		Body:      "We need automated deployments",
		State:     IssueStateOpen,
		AuthorID:  "devops@example.com",
		RepoName:  "acme/platform",
		Labels:    []string{"infrastructure", "ci-cd", "priority-high"},
		Assignees: []string{"alice@example.com", "bob@example.com"},
		CreatedAt: now,
	}

	if err := issue.Validate(); err != nil {
		t.Errorf("Valid issue failed validation: %v", err)
	}

	if len(issue.Labels) != 3 {
		t.Errorf("Labels count = %v, want 3", len(issue.Labels))
	}

	if len(issue.Assignees) != 2 {
		t.Errorf("Assignees count = %v, want 2", len(issue.Assignees))
	}

	if !issue.IsOpen() {
		t.Errorf("IsOpen() = false for open issue")
	}

	if issue.WasClosedByPR() {
		t.Errorf("WasClosedByPR() = true for issue without ClosedByPRID")
	}
}

func TestIssue_ClosedByPR(t *testing.T) {
	now := time.Now()
	closedAt := now.Add(24 * time.Hour)

	issue := Issue{
		Number:       42,
		Title:        "Fix authentication bug",
		Body:         "Login fails for users with special characters in email",
		State:        IssueStateClosed,
		AuthorID:     "tester@example.com",
		RepoName:     "acme/platform",
		Labels:       []string{"bug", "security"},
		CreatedAt:    now,
		UpdatedAt:    closedAt,
		ClosedAt:     &closedAt,
		ClosedByPRID: intPtr(123),
	}

	if err := issue.Validate(); err != nil {
		t.Errorf("Valid issue failed validation: %v", err)
	}

	if issue.IsOpen() {
		t.Errorf("IsOpen() = true for closed issue")
	}

	if !issue.IsClosed() {
		t.Errorf("IsClosed() = false for closed issue")
	}

	if !issue.WasClosedByPR() {
		t.Errorf("WasClosedByPR() = false for issue closed by PR")
	}

	if *issue.ClosedByPRID != 123 {
		t.Errorf("ClosedByPRID = %v, want 123", *issue.ClosedByPRID)
	}
}

// intPtr is a helper to create *int pointers for tests
func intPtr(i int) *int {
	return &i
}
