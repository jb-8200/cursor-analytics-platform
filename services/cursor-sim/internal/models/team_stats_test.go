package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentEditsDay_JSON(t *testing.T) {
	day := AgentEditsDay{
		EventDate:               "2026-01-15",
		TotalSuggestedDiffs:     100,
		TotalAcceptedDiffs:      75,
		TotalRejectedDiffs:      25,
		TotalGreenLinesAccepted: 150,
		TotalRedLinesAccepted:   30,
	}

	data, err := json.Marshal(day)
	require.NoError(t, err)

	var parsed AgentEditsDay
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "2026-01-15", parsed.EventDate)
	assert.Equal(t, 100, parsed.TotalSuggestedDiffs)
	assert.Equal(t, 75, parsed.TotalAcceptedDiffs)
}

func TestAgentEditsDay_FieldNames(t *testing.T) {
	day := AgentEditsDay{
		EventDate:                "2026-01-15",
		TotalSuggestedDiffs:      100,
		TotalAcceptedDiffs:       75,
		TotalRejectedDiffs:       25,
		TotalGreenLinesAccepted:  150,
		TotalRedLinesAccepted:    30,
		TotalGreenLinesRejected:  40,
		TotalRedLinesRejected:    10,
		TotalGreenLinesSuggested: 190,
		TotalRedLinesSuggested:   40,
		TotalLinesSuggested:      230,
		TotalLinesAccepted:       180,
	}

	data, err := json.Marshal(day)
	require.NoError(t, err)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	// Verify snake_case field names (Cursor API format for analytics)
	assert.Contains(t, raw, "event_date")
	assert.Contains(t, raw, "total_suggested_diffs")
	assert.Contains(t, raw, "total_accepted_diffs")
	assert.Contains(t, raw, "total_rejected_diffs")
	assert.Contains(t, raw, "total_green_lines_accepted")
	assert.Contains(t, raw, "total_red_lines_accepted")
	assert.Contains(t, raw, "total_green_lines_rejected")
	assert.Contains(t, raw, "total_red_lines_rejected")
	assert.Contains(t, raw, "total_green_lines_suggested")
	assert.Contains(t, raw, "total_red_lines_suggested")
	assert.Contains(t, raw, "total_lines_suggested")
	assert.Contains(t, raw, "total_lines_accepted")
}

func TestAgentEditsDay_AcceptanceRate(t *testing.T) {
	tests := []struct {
		name     string
		day      AgentEditsDay
		expected float64
	}{
		{
			name: "75% acceptance",
			day: AgentEditsDay{
				TotalSuggestedDiffs: 100,
				TotalAcceptedDiffs:  75,
			},
			expected: 0.75,
		},
		{
			name: "100% acceptance",
			day: AgentEditsDay{
				TotalSuggestedDiffs: 50,
				TotalAcceptedDiffs:  50,
			},
			expected: 1.0,
		},
		{
			name: "0% acceptance",
			day: AgentEditsDay{
				TotalSuggestedDiffs: 50,
				TotalAcceptedDiffs:  0,
			},
			expected: 0.0,
		},
		{
			name: "no suggestions",
			day: AgentEditsDay{
				TotalSuggestedDiffs: 0,
				TotalAcceptedDiffs:  0,
			},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate := tt.day.AcceptanceRate()
			assert.InDelta(t, tt.expected, rate, 0.01)
		})
	}
}

func TestTabUsageDay_JSON(t *testing.T) {
	day := TabUsageDay{
		EventDate:        "2026-01-15",
		TotalSuggestions: 500,
		TotalAccepts:     400,
		TotalRejects:     100,
	}

	data, err := json.Marshal(day)
	require.NoError(t, err)

	var parsed TabUsageDay
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "2026-01-15", parsed.EventDate)
	assert.Equal(t, 500, parsed.TotalSuggestions)
	assert.Equal(t, 400, parsed.TotalAccepts)
}

// Test backwards compatibility alias
func TestTabCompletionDay_Alias(t *testing.T) {
	// TabCompletionDay should be the same type as TabUsageDay
	day := TabCompletionDay{
		EventDate:        "2026-01-15",
		TotalSuggestions: 500,
		TotalAccepts:     400,
		TotalRejects:     100,
	}

	data, err := json.Marshal(day)
	require.NoError(t, err)

	var parsed TabCompletionDay
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "2026-01-15", parsed.EventDate)
	assert.Equal(t, 500, parsed.TotalSuggestions)
	assert.Equal(t, 400, parsed.TotalAccepts)
}

func TestDAUDay_JSON(t *testing.T) {
	day := DAUDay{
		Date:          "2026-01-15",
		DAU:           50,
		CLIDAU:        10,
		CloudAgentDAU: 5,
		BugbotDAU:     3,
	}

	data, err := json.Marshal(day)
	require.NoError(t, err)

	var parsed DAUDay
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "2026-01-15", parsed.Date)
	assert.Equal(t, 50, parsed.DAU)
	assert.Equal(t, 10, parsed.CLIDAU)
	assert.Equal(t, 5, parsed.CloudAgentDAU)
	assert.Equal(t, 3, parsed.BugbotDAU)
}

func TestDAUDay_FieldNames(t *testing.T) {
	day := DAUDay{
		Date:          "2026-01-15",
		DAU:           50,
		CLIDAU:        10,
		CloudAgentDAU: 5,
		BugbotDAU:     3,
	}

	data, err := json.Marshal(day)
	require.NoError(t, err)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	// Verify field names match Cursor API (NOTE: "date" not "event_date"!)
	assert.Contains(t, raw, "date")
	assert.Contains(t, raw, "dau")
	assert.Contains(t, raw, "cli_dau")
	assert.Contains(t, raw, "cloud_agent_dau")
	assert.Contains(t, raw, "bugbot_dau")

	// Verify event_date is NOT present
	assert.NotContains(t, raw, "event_date")
}

func TestModelUsageDay_JSON(t *testing.T) {
	day := ModelUsageDay{
		Date: "2026-01-15",
		ModelBreakdown: map[string]ModelBreakdownItem{
			"claude-sonnet-4.5": {Messages: 150, Users: 25},
			"gpt-4o":            {Messages: 80, Users: 15},
		},
	}

	data, err := json.Marshal(day)
	require.NoError(t, err)

	var parsed ModelUsageDay
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "2026-01-15", parsed.Date)
	assert.Equal(t, 2, len(parsed.ModelBreakdown))
	assert.Equal(t, 150, parsed.ModelBreakdown["claude-sonnet-4.5"].Messages)
	assert.Equal(t, 25, parsed.ModelBreakdown["claude-sonnet-4.5"].Users)
}

func TestModelUsageDay_FieldNames(t *testing.T) {
	day := ModelUsageDay{
		Date: "2026-01-15",
		ModelBreakdown: map[string]ModelBreakdownItem{
			"claude-sonnet-4.5": {Messages: 150, Users: 25},
		},
	}

	data, err := json.Marshal(day)
	require.NoError(t, err)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	// Verify field names
	assert.Contains(t, raw, "date")
	assert.Contains(t, raw, "model_breakdown")
}

func TestLeaderboardResponse_JSON(t *testing.T) {
	response := LeaderboardResponse{
		TabLeaderboard: LeaderboardSection{
			Data: []LeaderboardEntry{
				{
					Email:               "alice@example.com",
					UserID:              "user_001",
					TotalAccepts:        500,
					TotalLinesAccepted:  2000,
					TotalLinesSuggested: 2500,
					LineAcceptanceRatio: 0.8,
					Rank:                1,
				},
			},
			TotalUsers: 10,
		},
		AgentLeaderboard: LeaderboardSection{
			Data: []LeaderboardEntry{
				{
					Email:               "bob@example.com",
					UserID:              "user_002",
					TotalAccepts:        300,
					TotalLinesAccepted:  1200,
					TotalLinesSuggested: 1500,
					LineAcceptanceRatio: 0.8,
					FavoriteModel:       "claude-sonnet-4.5",
					Rank:                1,
				},
			},
			TotalUsers: 10,
		},
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var parsed LeaderboardResponse
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, 1, len(parsed.TabLeaderboard.Data))
	assert.Equal(t, 1, len(parsed.AgentLeaderboard.Data))
	assert.Equal(t, "alice@example.com", parsed.TabLeaderboard.Data[0].Email)
	assert.Equal(t, "bob@example.com", parsed.AgentLeaderboard.Data[0].Email)
}
