package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullRequest_JSONMarshaling(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	mergedAt := now.Add(24 * time.Hour)

	pr := PullRequest{
		Number:        42,
		Title:         "feat: add user authentication",
		Body:          "This PR adds OAuth2 authentication",
		State:         PRStateMerged,
		AuthorID:      "user_001",
		AuthorEmail:   "alice@example.com",
		AuthorName:    "Alice Developer",
		RepoName:      "acme-corp/platform",
		BaseBranch:    "main",
		HeadBranch:    "feature/auth",
		Reviewers:     []string{"user_002", "user_003"},
		Labels:        []string{"enhancement", "security"},
		Additions:     250,
		Deletions:     30,
		ChangedFiles:  8,
		CommitCount:   5,
		AIRatio:       0.75,
		TabLines:      150,
		ComposerLines: 50,
		CreatedAt:     now,
		UpdatedAt:     now.Add(12 * time.Hour),
		MergedAt:      &mergedAt,
		ClosedAt:      nil,
		WasReverted:   false,
		IsBugFix:      false,
	}

	data, err := json.Marshal(pr)
	require.NoError(t, err)

	var parsed PullRequest
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, pr.Number, parsed.Number)
	assert.Equal(t, pr.Title, parsed.Title)
	assert.Equal(t, pr.State, parsed.State)
	assert.Equal(t, pr.AuthorID, parsed.AuthorID)
	assert.Equal(t, pr.Reviewers, parsed.Reviewers)
	assert.Equal(t, pr.Labels, parsed.Labels)
	assert.Equal(t, pr.AIRatio, parsed.AIRatio)
	assert.NotNil(t, parsed.MergedAt)
	assert.Nil(t, parsed.ClosedAt)
}

func TestPullRequest_JSONFieldNames(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)

	pr := PullRequest{
		Number:        1,
		Title:         "test",
		Body:          "body",
		State:         PRStateOpen,
		AuthorID:      "user_001",
		AuthorEmail:   "test@example.com",
		AuthorName:    "Test User",
		RepoName:      "test/repo",
		BaseBranch:    "main",
		HeadBranch:    "feature",
		Reviewers:     []string{},
		Labels:        []string{},
		Additions:     100,
		Deletions:     10,
		ChangedFiles:  5,
		CommitCount:   2,
		AIRatio:       0.5,
		TabLines:      30,
		ComposerLines: 20,
		CreatedAt:     now,
		UpdatedAt:     now,
		WasReverted:   false,
		IsBugFix:      false,
	}

	data, err := json.Marshal(pr)
	require.NoError(t, err)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	// Verify snake_case field names (GitHub API format)
	assert.Contains(t, raw, "number")
	assert.Contains(t, raw, "title")
	assert.Contains(t, raw, "body")
	assert.Contains(t, raw, "state")
	assert.Contains(t, raw, "author_id")
	assert.Contains(t, raw, "author_email")
	assert.Contains(t, raw, "author_name")
	assert.Contains(t, raw, "repo_name")
	assert.Contains(t, raw, "base_branch")
	assert.Contains(t, raw, "head_branch")
	assert.Contains(t, raw, "reviewers")
	assert.Contains(t, raw, "labels")
	assert.Contains(t, raw, "additions")
	assert.Contains(t, raw, "deletions")
	assert.Contains(t, raw, "changed_files")
	assert.Contains(t, raw, "commit_count")
	assert.Contains(t, raw, "ai_ratio")
	assert.Contains(t, raw, "tab_lines")
	assert.Contains(t, raw, "composer_lines")
	assert.Contains(t, raw, "created_at")
	assert.Contains(t, raw, "updated_at")
	assert.Contains(t, raw, "was_reverted")
	assert.Contains(t, raw, "is_bug_fix")
}

func TestPullRequest_States(t *testing.T) {
	assert.Equal(t, PRState("open"), PRStateOpen)
	assert.Equal(t, PRState("closed"), PRStateClosed)
	assert.Equal(t, PRState("merged"), PRStateMerged)
}

func TestPullRequest_IsOpen(t *testing.T) {
	tests := []struct {
		name   string
		state  PRState
		isOpen bool
	}{
		{"open PR", PRStateOpen, true},
		{"closed PR", PRStateClosed, false},
		{"merged PR", PRStateMerged, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PullRequest{State: tt.state}
			assert.Equal(t, tt.isOpen, pr.IsOpen())
		})
	}
}

func TestPullRequest_IsMerged(t *testing.T) {
	tests := []struct {
		name     string
		state    PRState
		isMerged bool
	}{
		{"open PR", PRStateOpen, false},
		{"closed PR", PRStateClosed, false},
		{"merged PR", PRStateMerged, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PullRequest{State: tt.state}
			assert.Equal(t, tt.isMerged, pr.IsMerged())
		})
	}
}

func TestPullRequest_NetLines(t *testing.T) {
	pr := PullRequest{
		Additions: 150,
		Deletions: 30,
	}
	assert.Equal(t, 120, pr.NetLines())
}

func TestPullRequest_AILines(t *testing.T) {
	pr := PullRequest{
		TabLines:      100,
		ComposerLines: 50,
	}
	assert.Equal(t, 150, pr.AILines())
}

func TestReviewComment_JSONMarshaling(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)

	comment := ReviewComment{
		ID:        123,
		PRNumber:  42,
		RepoName:  "acme-corp/platform",
		AuthorID:  "user_002",
		Body:      "LGTM! Nice work on the error handling.",
		Path:      "src/auth/handler.go",
		Line:      42,
		State:     ReviewStateApproved,
		CreatedAt: now,
	}

	data, err := json.Marshal(comment)
	require.NoError(t, err)

	var parsed ReviewComment
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, comment.ID, parsed.ID)
	assert.Equal(t, comment.PRNumber, parsed.PRNumber)
	assert.Equal(t, comment.AuthorID, parsed.AuthorID)
	assert.Equal(t, comment.Body, parsed.Body)
	assert.Equal(t, comment.Path, parsed.Path)
	assert.Equal(t, comment.Line, parsed.Line)
	assert.Equal(t, comment.State, parsed.State)
}

func TestReviewComment_JSONFieldNames(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)

	comment := ReviewComment{
		ID:        1,
		PRNumber:  1,
		RepoName:  "test/repo",
		AuthorID:  "user_001",
		Body:      "test",
		Path:      "test.go",
		Line:      10,
		State:     ReviewStatePending,
		CreatedAt: now,
	}

	data, err := json.Marshal(comment)
	require.NoError(t, err)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	// Verify snake_case field names (GitHub API format)
	assert.Contains(t, raw, "id")
	assert.Contains(t, raw, "pr_number")
	assert.Contains(t, raw, "repo_name")
	assert.Contains(t, raw, "author_id")
	assert.Contains(t, raw, "body")
	assert.Contains(t, raw, "path")
	assert.Contains(t, raw, "line")
	assert.Contains(t, raw, "state")
	assert.Contains(t, raw, "created_at")
}

func TestReviewComment_States(t *testing.T) {
	assert.Equal(t, ReviewState("pending"), ReviewStatePending)
	assert.Equal(t, ReviewState("approved"), ReviewStateApproved)
	assert.Equal(t, ReviewState("changes_requested"), ReviewStateChangesRequested)
}

func TestReviewComment_IsApproval(t *testing.T) {
	tests := []struct {
		name       string
		state      ReviewState
		isApproval bool
	}{
		{"pending", ReviewStatePending, false},
		{"approved", ReviewStateApproved, true},
		{"changes requested", ReviewStateChangesRequested, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment := ReviewComment{State: tt.state}
			assert.Equal(t, tt.isApproval, comment.IsApproval())
		})
	}
}

func TestRepository_JSONMarshaling(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)

	repo := Repository{
		Name:            "acme-corp/platform",
		Owner:           "acme-corp",
		Description:     "Main platform repository",
		PrimaryLanguage: "go",
		DefaultBranch:   "main",
		Teams:           []string{"Platform", "API"},
		CreatedAt:       now,
	}

	data, err := json.Marshal(repo)
	require.NoError(t, err)

	var parsed Repository
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, repo.Name, parsed.Name)
	assert.Equal(t, repo.Owner, parsed.Owner)
	assert.Equal(t, repo.Description, parsed.Description)
	assert.Equal(t, repo.PrimaryLanguage, parsed.PrimaryLanguage)
	assert.Equal(t, repo.Teams, parsed.Teams)
}

func TestRepository_JSONFieldNames(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)

	repo := Repository{
		Name:            "test/repo",
		Owner:           "test",
		Description:     "test repo",
		PrimaryLanguage: "go",
		DefaultBranch:   "main",
		Teams:           []string{},
		CreatedAt:       now,
	}

	data, err := json.Marshal(repo)
	require.NoError(t, err)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	// Verify snake_case field names (GitHub API format)
	assert.Contains(t, raw, "name")
	assert.Contains(t, raw, "owner")
	assert.Contains(t, raw, "description")
	assert.Contains(t, raw, "primary_language")
	assert.Contains(t, raw, "default_branch")
	assert.Contains(t, raw, "teams")
	assert.Contains(t, raw, "created_at")
}

func TestRepository_FullName(t *testing.T) {
	repo := Repository{
		Name:  "platform",
		Owner: "acme-corp",
	}
	// Name already includes owner/repo format from seed
	repo.Name = "acme-corp/platform"
	assert.Equal(t, "acme-corp/platform", repo.FullName())
}
