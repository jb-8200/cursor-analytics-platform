package generator

import (
	"strings"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSurveyGenerator_GenerateResponses(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev1@company.com", Name: "Dev One"},
			{Email: "dev2@company.com", Name: "Dev Two"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				SurveyName:    "AI Tools Survey",
				ResponseCount: 10,
			},
		},
	}

	gen := NewSurveyGeneratorWithSeed(seedData, 12345)
	responses := gen.GenerateResponses("SV_abc123")

	// Should generate exactly the configured number of responses
	assert.Len(t, responses, 10)

	// Verify all responses have required fields
	for _, r := range responses {
		assert.NotEmpty(t, r.ResponseID, "ResponseID should not be empty")
		assert.NotEmpty(t, r.RespondentEmail, "RespondentEmail should not be empty")
		assert.True(t, r.OverallAISatisfaction >= 1 && r.OverallAISatisfaction <= 5, "OverallAISatisfaction should be 1-5")
		assert.True(t, r.CursorSatisfaction >= 1 && r.CursorSatisfaction <= 5, "CursorSatisfaction should be 1-5")
		assert.True(t, r.CopilotSatisfaction >= 1 && r.CopilotSatisfaction <= 5, "CopilotSatisfaction should be 1-5")
		assert.NotEmpty(t, r.MostUsedTool, "MostUsedTool should not be empty")
		assert.False(t, r.RecordedAt.IsZero(), "RecordedAt should not be zero")
	}
}

func TestSurveyGenerator_SatisfactionDistribution(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev1@company.com", Name: "Dev One"},
			{Email: "dev2@company.com", Name: "Dev Two"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				SurveyName:    "AI Tools Survey",
				ResponseCount: 100,
			},
		},
	}

	gen := NewSurveyGeneratorWithSeed(seedData, 12345)
	responses := gen.GenerateResponses("SV_abc123")

	require.Len(t, responses, 100)

	// Count satisfaction scores
	counts := make(map[int]int)
	for _, r := range responses {
		counts[r.OverallAISatisfaction]++
	}

	// Default distribution should have bell curve centered around 3-4
	// Allow tolerance for randomness
	total := float64(len(responses))

	// Most responses should be 3, 4, or 5 (high satisfaction)
	highSatisfaction := counts[3] + counts[4] + counts[5]
	assert.True(t, float64(highSatisfaction)/total >= 0.6, "Expected at least 60%% satisfaction in 3-5 range")

	// Should have some variation (not all the same)
	assert.True(t, len(counts) >= 3, "Expected at least 3 different satisfaction levels")
}

func TestSurveyGenerator_Reproducible(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				SurveyName:    "AI Tools Survey",
				ResponseCount: 20,
			},
		},
	}

	// Same seed should produce identical results
	gen1 := NewSurveyGeneratorWithSeed(seedData, 12345)
	responses1 := gen1.GenerateResponses("SV_abc123")

	gen2 := NewSurveyGeneratorWithSeed(seedData, 12345)
	responses2 := gen2.GenerateResponses("SV_abc123")

	require.Equal(t, len(responses1), len(responses2))
	for i := range responses1 {
		assert.Equal(t, responses1[i].ResponseID, responses2[i].ResponseID)
		assert.Equal(t, responses1[i].RespondentEmail, responses2[i].RespondentEmail)
		assert.Equal(t, responses1[i].OverallAISatisfaction, responses2[i].OverallAISatisfaction)
		assert.Equal(t, responses1[i].MostUsedTool, responses2[i].MostUsedTool)
	}
}

func TestSurveyGenerator_RespondentSelection(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				SurveyName:    "AI Tools Survey",
				ResponseCount: 10,
			},
		},
	}

	gen := NewSurveyGeneratorWithSeed(seedData, 12345)
	responses := gen.GenerateResponses("SV_abc123")

	// All respondents should be from developers
	for _, r := range responses {
		assert.True(t, strings.Contains(r.RespondentEmail, "company.com"), "Expected developer email")
	}
}

func TestSurveyGenerator_FreeTextFeedback(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				SurveyName:    "AI Tools Survey",
				ResponseCount: 50,
			},
		},
	}

	gen := NewSurveyGeneratorWithSeed(seedData, 12345)
	responses := gen.GenerateResponses("SV_abc123")

	// At least some responses should have feedback (not all, not none)
	hasPositive := false
	hasImprovement := false
	emptyPositive := 0
	emptyImprovement := 0

	for _, r := range responses {
		if r.PositiveFeedback != "" {
			hasPositive = true
		} else {
			emptyPositive++
		}
		if r.ImprovementAreas != "" {
			hasImprovement = true
		} else {
			emptyImprovement++
		}
	}

	assert.True(t, hasPositive, "Expected some positive feedback")
	assert.True(t, hasImprovement, "Expected some improvement areas")
	assert.True(t, emptyPositive > 0, "Expected some empty positive feedback")
	assert.True(t, emptyImprovement > 0, "Expected some empty improvement areas")
}

func TestSurveyGenerator_MostUsedToolDistribution(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				SurveyName:    "AI Tools Survey",
				ResponseCount: 100,
			},
		},
	}

	gen := NewSurveyGeneratorWithSeed(seedData, 12345)
	responses := gen.GenerateResponses("SV_abc123")

	// Count tool mentions
	toolCounts := make(map[string]int)
	for _, r := range responses {
		toolCounts[r.MostUsedTool]++
	}

	// Should have multiple tools mentioned
	assert.True(t, len(toolCounts) >= 2, "Expected at least 2 different tools mentioned")

	// All tools should be valid
	validTools := map[string]bool{
		"Cursor":         true,
		"GitHub Copilot": true,
		"Both Equally":   true,
		"Neither":        true,
	}
	for tool := range toolCounts {
		assert.True(t, validTools[tool], "Tool %s should be valid", tool)
	}
}

func TestSurveyGenerator_TimestampWithinReasonableRange(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				SurveyName:    "AI Tools Survey",
				ResponseCount: 10,
			},
		},
	}

	gen := NewSurveyGeneratorWithSeed(seedData, 12345)
	responses := gen.GenerateResponses("SV_abc123")

	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	// All timestamps should be within last 30 days
	for _, r := range responses {
		assert.True(t, r.RecordedAt.After(thirtyDaysAgo) || r.RecordedAt.Equal(thirtyDaysAgo),
			"RecordedAt should be within last 30 days")
		assert.True(t, r.RecordedAt.Before(now) || r.RecordedAt.Equal(now),
			"RecordedAt should not be in the future")
	}
}

func TestSurveyGenerator_NoQualtricsConfig(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Developer"},
		},
		// No Qualtrics config
	}

	gen := NewSurveyGeneratorWithSeed(seedData, 12345)
	responses := gen.GenerateResponses("SV_abc123")

	// Should return empty slice when no config
	assert.Empty(t, responses)
}

func TestSurveyGenerator_QualtricsDisabled(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       false, // Disabled
				SurveyID:      "SV_abc123",
				SurveyName:    "AI Tools Survey",
				ResponseCount: 10,
			},
		},
	}

	gen := NewSurveyGeneratorWithSeed(seedData, 12345)
	responses := gen.GenerateResponses("SV_abc123")

	// Should return empty slice when disabled
	assert.Empty(t, responses)
}

func TestSurveyGenerator_NoDevelopers(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{}, // No developers
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				SurveyName:    "AI Tools Survey",
				ResponseCount: 10,
			},
		},
	}

	gen := NewSurveyGeneratorWithSeed(seedData, 12345)
	responses := gen.GenerateResponses("SV_abc123")

	// Should return empty slice when no developers
	assert.Empty(t, responses)
}

func TestSurveyGenerator_WrongSurveyID(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				SurveyName:    "AI Tools Survey",
				ResponseCount: 10,
			},
		},
	}

	gen := NewSurveyGeneratorWithSeed(seedData, 12345)
	responses := gen.GenerateResponses("SV_wrong_id")

	// Should return empty slice when survey ID doesn't match
	assert.Empty(t, responses)
}

func TestSurveyGenerator_NewWithRandomSeed(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:       true,
				SurveyID:      "SV_abc123",
				SurveyName:    "AI Tools Survey",
				ResponseCount: 5,
			},
		},
	}

	// Test that NewSurveyGenerator (without explicit seed) works
	gen := NewSurveyGenerator(seedData)
	responses := gen.GenerateResponses("SV_abc123")

	// Should generate responses
	assert.Len(t, responses, 5)
}
