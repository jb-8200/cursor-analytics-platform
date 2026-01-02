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
		EventDate:               "2026-01-15",
		TotalSuggestedDiffs:     100,
		TotalAcceptedDiffs:      75,
		TotalRejectedDiffs:      25,
		TotalGreenLinesAccepted: 150,
		TotalRedLinesAccepted:   30,
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

func TestTabCompletionDay_JSON(t *testing.T) {
	day := TabCompletionDay{
		EventDate:     "2026-01-15",
		TotalSuggests: 500,
		TotalAccepts:  400,
		TotalRejects:  100,
	}

	data, err := json.Marshal(day)
	require.NoError(t, err)

	var parsed TabCompletionDay
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "2026-01-15", parsed.EventDate)
	assert.Equal(t, 500, parsed.TotalSuggests)
	assert.Equal(t, 400, parsed.TotalAccepts)
}

func TestDAUDay_JSON(t *testing.T) {
	day := DAUDay{
		EventDate:        "2026-01-15",
		UniqueUsers:      50,
		TotalEvents:      1000,
		AvgEventsPerUser: 20.0,
	}

	data, err := json.Marshal(day)
	require.NoError(t, err)

	var parsed DAUDay
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "2026-01-15", parsed.EventDate)
	assert.Equal(t, 50, parsed.UniqueUsers)
	assert.Equal(t, 20.0, parsed.AvgEventsPerUser)
}

func TestModelUsage_JSON(t *testing.T) {
	usage := ModelUsage{
		ModelName:   "gpt-4-turbo",
		TotalUsages: 150,
		UniqueUsers: 25,
	}

	data, err := json.Marshal(usage)
	require.NoError(t, err)

	var parsed ModelUsage
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "gpt-4-turbo", parsed.ModelName)
	assert.Equal(t, 150, parsed.TotalUsages)
	assert.Equal(t, 25, parsed.UniqueUsers)
}
