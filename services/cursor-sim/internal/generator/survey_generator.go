package generator

import (
	"encoding/hex"
	mathrand "math/rand"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// SurveyGenerator generates Qualtrics survey response data based on seed configuration.
type SurveyGenerator struct {
	seedData *seed.SeedData
	rng      *mathrand.Rand
}

// NewSurveyGenerator creates a new survey generator with a random seed.
func NewSurveyGenerator(seedData *seed.SeedData) *SurveyGenerator {
	return NewSurveyGeneratorWithSeed(seedData, time.Now().UnixNano())
}

// NewSurveyGeneratorWithSeed creates a new survey generator with a specific seed for reproducibility.
func NewSurveyGeneratorWithSeed(seedData *seed.SeedData, randSeed int64) *SurveyGenerator {
	return &SurveyGenerator{
		seedData: seedData,
		rng:      mathrand.New(mathrand.NewSource(randSeed)),
	}
}

// GenerateResponses generates survey responses for the specified survey ID.
// Returns empty slice if Qualtrics is not configured, disabled, or survey ID doesn't match.
func (g *SurveyGenerator) GenerateResponses(surveyID string) []models.SurveyResponse {
	var responses []models.SurveyResponse

	// Check if Qualtrics is configured and enabled
	if g.seedData.ExternalDataSources == nil ||
		g.seedData.ExternalDataSources.Qualtrics == nil ||
		!g.seedData.ExternalDataSources.Qualtrics.Enabled {
		return responses
	}

	qualtricsConfig := g.seedData.ExternalDataSources.Qualtrics

	// Verify survey ID matches
	if qualtricsConfig.SurveyID != surveyID {
		return responses
	}

	// Check if we have developers to select from
	if len(g.seedData.Developers) == 0 {
		return responses
	}

	// Generate configured number of responses
	responseCount := qualtricsConfig.ResponseCount
	if responseCount <= 0 {
		return responses
	}

	// Time range for responses (last 30 days)
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	// Generate responses
	for i := 0; i < responseCount; i++ {
		// Select random respondent from developers
		respondent := g.seedData.Developers[g.rng.Intn(len(g.seedData.Developers))]

		// Generate satisfaction scores (bell curve centered around 3-4)
		overallSat := g.generateSatisfactionScore()
		cursorSat := g.generateSatisfactionScore()
		copilotSat := g.generateSatisfactionScore()

		// Select most used tool based on satisfaction scores
		mostUsedTool := g.selectMostUsedTool(cursorSat, copilotSat)

		// Generate free text feedback (optional, ~40% chance)
		positiveFeedback := ""
		improvementAreas := ""
		if g.rng.Float64() < 0.4 {
			positiveFeedback = g.generatePositiveFeedback(mostUsedTool)
		}
		if g.rng.Float64() < 0.4 {
			improvementAreas = g.generateImprovementAreas(mostUsedTool)
		}

		// Generate random timestamp within last 30 days
		recordedAt := g.randomTimeInRange(thirtyDaysAgo, now)

		response := models.SurveyResponse{
			ResponseID:            g.generateResponseID(),
			RespondentEmail:       respondent.Email,
			OverallAISatisfaction: overallSat,
			CursorSatisfaction:    cursorSat,
			CopilotSatisfaction:   copilotSat,
			MostUsedTool:          mostUsedTool,
			PositiveFeedback:      positiveFeedback,
			ImprovementAreas:      improvementAreas,
			RecordedAt:            recordedAt,
		}

		responses = append(responses, response)
	}

	return responses
}

// generateSatisfactionScore generates a satisfaction score (1-5) with bell curve distribution.
// Weighted toward 3-5 (higher satisfaction) to simulate realistic enterprise AI tool adoption.
func (g *SurveyGenerator) generateSatisfactionScore() int {
	// Distribution: 1:5%, 2:10%, 3:25%, 4:40%, 5:20%
	r := g.rng.Float64()
	if r < 0.05 {
		return 1
	} else if r < 0.15 {
		return 2
	} else if r < 0.40 {
		return 3
	} else if r < 0.80 {
		return 4
	}
	return 5
}

// selectMostUsedTool selects the most used tool based on satisfaction scores.
func (g *SurveyGenerator) selectMostUsedTool(cursorSat, copilotSat int) string {
	diff := cursorSat - copilotSat

	if diff > 1 {
		return "Cursor"
	} else if diff < -1 {
		return "GitHub Copilot"
	} else if diff == 0 {
		// Tie - randomly pick or "Both Equally"
		r := g.rng.Float64()
		if r < 0.6 {
			return "Both Equally"
		} else if r < 0.8 {
			return "Cursor"
		}
		return "GitHub Copilot"
	} else {
		// Small difference - pick the higher one or both equally
		r := g.rng.Float64()
		if r < 0.3 {
			return "Both Equally"
		} else if cursorSat > copilotSat {
			return "Cursor"
		}
		return "GitHub Copilot"
	}
}

// generatePositiveFeedback generates realistic positive feedback.
func (g *SurveyGenerator) generatePositiveFeedback(tool string) string {
	feedback := []string{
		"Great productivity boost",
		"Helps me write code faster",
		"Very helpful for learning new patterns",
		"Saves time on boilerplate code",
		"Excellent code suggestions",
		"Makes debugging easier",
		"Good for exploring new libraries",
		"Helpful for refactoring",
		"Assists with documentation",
		"Improves code quality",
	}

	return feedback[g.rng.Intn(len(feedback))]
}

// generateImprovementAreas generates realistic improvement suggestions.
func (g *SurveyGenerator) generateImprovementAreas(tool string) string {
	improvements := []string{
		"Sometimes suggests incorrect code",
		"Could be better at understanding context",
		"Needs improvement for edge cases",
		"Would like better integration with our workflow",
		"Occasionally slow response times",
		"More accurate suggestions needed",
		"Better support for our tech stack",
		"Could use more customization options",
		"Sometimes suggests outdated patterns",
		"Would benefit from better error handling suggestions",
	}

	return improvements[g.rng.Intn(len(improvements))]
}

// generateResponseID generates a unique response ID (Qualtrics format: R_xxx).
// Uses seeded random generator for reproducibility.
func (g *SurveyGenerator) generateResponseID() string {
	bytes := make([]byte, 8)
	for i := range bytes {
		bytes[i] = byte(g.rng.Intn(256))
	}
	return "R_" + hex.EncodeToString(bytes)
}

// randomTimeInRange generates a random time within the specified range.
func (g *SurveyGenerator) randomTimeInRange(start, end time.Time) time.Time {
	duration := end.Sub(start)
	randomDuration := time.Duration(g.rng.Int63n(int64(duration)))
	return start.Add(randomDuration)
}
