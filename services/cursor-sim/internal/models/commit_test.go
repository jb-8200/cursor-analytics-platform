package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommit_JSONMarshaling(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)

	commit := Commit{
		CommitHash:           "abc123def456",
		UserID:               "user_001",
		UserEmail:            "alice@example.com",
		UserName:             "Alice Developer",
		RepoName:             "acme-corp/payment-service",
		BranchName:           "feature/new-feature",
		IsPrimaryBranch:      false,
		TotalLinesAdded:      150,
		TotalLinesDeleted:    20,
		TabLinesAdded:        80,
		TabLinesDeleted:      10,
		ComposerLinesAdded:   40,
		ComposerLinesDeleted: 5,
		NonAILinesAdded:      30,
		NonAILinesDeleted:    5,
		Message:              "feat: add new payment feature",
		CommitTs:             now,
		CreatedAt:            now,
	}

	data, err := json.Marshal(commit)
	require.NoError(t, err)

	var parsed Commit
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, commit.CommitHash, parsed.CommitHash)
	assert.Equal(t, commit.UserEmail, parsed.UserEmail)
	assert.Equal(t, commit.TotalLinesAdded, parsed.TotalLinesAdded)
	assert.Equal(t, commit.TabLinesAdded, parsed.TabLinesAdded)
	assert.Equal(t, commit.CommitTs.Unix(), parsed.CommitTs.Unix())
}

func TestCommit_JSONFieldNames(t *testing.T) {
	now := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)

	commit := Commit{
		CommitHash:           "abc123",
		UserID:               "user_001",
		UserEmail:            "test@example.com",
		UserName:             "Test User",
		RepoName:             "test/repo",
		BranchName:           "main",
		IsPrimaryBranch:      true,
		TotalLinesAdded:      100,
		TotalLinesDeleted:    10,
		TabLinesAdded:        60,
		TabLinesDeleted:      5,
		ComposerLinesAdded:   30,
		ComposerLinesDeleted: 3,
		NonAILinesAdded:      10,
		NonAILinesDeleted:    2,
		Message:              "test commit",
		CommitTs:             now,
		CreatedAt:            now,
	}

	data, err := json.Marshal(commit)
	require.NoError(t, err)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	// Verify camelCase field names (Cursor API format)
	assert.Contains(t, raw, "commitHash")
	assert.Contains(t, raw, "userId")
	assert.Contains(t, raw, "userEmail")
	assert.Contains(t, raw, "userName")
	assert.Contains(t, raw, "repoName")
	assert.Contains(t, raw, "branchName")
	assert.Contains(t, raw, "isPrimaryBranch")
	assert.Contains(t, raw, "totalLinesAdded")
	assert.Contains(t, raw, "totalLinesDeleted")
	assert.Contains(t, raw, "tabLinesAdded")
	assert.Contains(t, raw, "tabLinesDeleted")
	assert.Contains(t, raw, "composerLinesAdded")
	assert.Contains(t, raw, "composerLinesDeleted")
	assert.Contains(t, raw, "nonAiLinesAdded")
	assert.Contains(t, raw, "nonAiLinesDeleted")
	assert.Contains(t, raw, "message")
	assert.Contains(t, raw, "commitTs")
	assert.Contains(t, raw, "createdAt")
}

func TestCommit_AIRatio(t *testing.T) {
	tests := []struct {
		name          string
		commit        Commit
		expectedRatio float64
	}{
		{
			name: "80% AI ratio",
			commit: Commit{
				TotalLinesAdded:    100,
				TabLinesAdded:      50,
				ComposerLinesAdded: 30,
				NonAILinesAdded:    20,
			},
			expectedRatio: 0.80,
		},
		{
			name: "100% AI ratio",
			commit: Commit{
				TotalLinesAdded:    100,
				TabLinesAdded:      60,
				ComposerLinesAdded: 40,
				NonAILinesAdded:    0,
			},
			expectedRatio: 1.0,
		},
		{
			name: "0% AI ratio",
			commit: Commit{
				TotalLinesAdded:    100,
				TabLinesAdded:      0,
				ComposerLinesAdded: 0,
				NonAILinesAdded:    100,
			},
			expectedRatio: 0.0,
		},
		{
			name: "no lines added",
			commit: Commit{
				TotalLinesAdded: 0,
			},
			expectedRatio: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ratio := tt.commit.AIRatio()
			assert.InDelta(t, tt.expectedRatio, ratio, 0.01)
		})
	}
}

func TestCommit_NetLines(t *testing.T) {
	commit := Commit{
		TotalLinesAdded:   150,
		TotalLinesDeleted: 30,
	}

	assert.Equal(t, 120, commit.NetLines())
}

func TestCommit_HasAIContent(t *testing.T) {
	tests := []struct {
		name   string
		commit Commit
		hasAI  bool
	}{
		{
			name: "has tab content",
			commit: Commit{
				TabLinesAdded:      10,
				ComposerLinesAdded: 0,
			},
			hasAI: true,
		},
		{
			name: "has composer content",
			commit: Commit{
				TabLinesAdded:      0,
				ComposerLinesAdded: 5,
			},
			hasAI: true,
		},
		{
			name: "has both",
			commit: Commit{
				TabLinesAdded:      10,
				ComposerLinesAdded: 5,
			},
			hasAI: true,
		},
		{
			name: "no AI content",
			commit: Commit{
				TabLinesAdded:      0,
				ComposerLinesAdded: 0,
				NonAILinesAdded:    100,
			},
			hasAI: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.hasAI, tt.commit.HasAIContent())
		})
	}
}
