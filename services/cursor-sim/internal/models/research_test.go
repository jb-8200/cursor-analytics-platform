package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResearchDataPoint_JSONSerialization(t *testing.T) {
	ts := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	dp := ResearchDataPoint{
		CommitHash:          "abc123",
		PRNumber:            42,
		AuthorID:            "user_001",
		RepoName:            "test/repo",
		AIRatio:             0.65,
		TabLines:            100,
		ComposerLines:       50,
		Additions:           200,
		Deletions:           50,
		FilesChanged:        5,
		CodingLeadTimeHours: 4.5,
		ReviewLeadTimeHours: 2.0,
		MergeLeadTimeHours:  1.5,
		WasReverted:         false,
		RequiredHotfix:      false,
		ReviewIterations:    2,
		AuthorSeniority:     "senior",
		RepoMaturity:        "mature",
		IsGreenfield:        false,
		Timestamp:           ts,
	}

	// Test JSON serialization
	data, err := json.Marshal(dp)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"commit_hash":"abc123"`)
	assert.Contains(t, string(data), `"ai_ratio":0.65`)
	assert.Contains(t, string(data), `"author_seniority":"senior"`)

	// Test JSON deserialization
	var decoded ResearchDataPoint
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, dp.CommitHash, decoded.CommitHash)
	assert.Equal(t, dp.AIRatio, decoded.AIRatio)
	assert.Equal(t, dp.AuthorSeniority, decoded.AuthorSeniority)
}

func TestResearchDataPoint_AIRatioCategory(t *testing.T) {
	tests := []struct {
		name     string
		aiRatio  float64
		expected AIRatioBand
	}{
		{"Low AI ratio", 0.15, AIRatioBandLow},
		{"Medium AI ratio", 0.45, AIRatioBandMedium},
		{"High AI ratio", 0.85, AIRatioBandHigh},
		{"Zero ratio", 0.0, AIRatioBandLow},
		{"Full AI", 1.0, AIRatioBandHigh},
		{"Boundary low-medium", 0.3, AIRatioBandMedium},
		{"Boundary medium-high", 0.7, AIRatioBandHigh},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := ResearchDataPoint{AIRatio: tt.aiRatio}
			assert.Equal(t, tt.expected, dp.GetAIRatioBand())
		})
	}
}

func TestResearchDataPoint_TotalLeadTime(t *testing.T) {
	dp := ResearchDataPoint{
		CodingLeadTimeHours: 4.5,
		ReviewLeadTimeHours: 2.0,
		MergeLeadTimeHours:  1.5,
	}

	assert.Equal(t, 8.0, dp.TotalLeadTimeHours())
}

func TestResearchDataPoint_Validate(t *testing.T) {
	validDP := ResearchDataPoint{
		CommitHash:    "abc123",
		AuthorID:      "user_001",
		RepoName:      "test/repo",
		AIRatio:       0.5,
		Additions:     100,
		TabLines:      30,
		ComposerLines: 20,
		Timestamp:     time.Now(),
	}

	// Valid data point should pass
	assert.NoError(t, validDP.Validate())

	// Missing commit hash
	invalidDP := validDP
	invalidDP.CommitHash = ""
	assert.Error(t, invalidDP.Validate())

	// AI ratio out of range
	invalidDP = validDP
	invalidDP.AIRatio = 1.5
	assert.Error(t, invalidDP.Validate())

	// Negative additions
	invalidDP = validDP
	invalidDP.Additions = -10
	assert.Error(t, invalidDP.Validate())
}

func TestVelocityMetrics_JSONSerialization(t *testing.T) {
	vm := VelocityMetrics{
		Period:              "2026-01",
		AIRatioBand:         AIRatioBandHigh,
		TotalCommits:        150,
		TotalPRs:            30,
		TotalAdditions:      5000,
		TotalDeletions:      2000,
		AvgCommitsPerDay:    5.0,
		AvgPRsPerWeek:       7.5,
		AvgLeadTimeHours:    24.5,
		MedianLeadTimeHours: 18.0,
		StdDevLeadTimeHours: 8.5,
	}

	data, err := json.Marshal(vm)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"ai_ratio_band":"high"`)
	assert.Contains(t, string(data), `"total_commits":150`)

	var decoded VelocityMetrics
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, vm.AIRatioBand, decoded.AIRatioBand)
	assert.Equal(t, vm.TotalCommits, decoded.TotalCommits)
}

func TestReviewCostMetrics_JSONSerialization(t *testing.T) {
	rcm := ReviewCostMetrics{
		Period:                "2026-01",
		AIRatioBand:           AIRatioBandMedium,
		TotalPRsReviewed:      50,
		TotalReviewComments:   200,
		TotalReviewIterations: 75,
		AvgCommentsPerPR:      4.0,
		AvgIterationsPerPR:    1.5,
		AvgReviewTimeHours:    3.5,
		MedianReviewTimeHours: 2.5,
		StdDevReviewTimeHours: 2.0,
	}

	data, err := json.Marshal(rcm)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"ai_ratio_band":"medium"`)
	assert.Contains(t, string(data), `"avg_comments_per_pr":4`)

	var decoded ReviewCostMetrics
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, rcm.AvgCommentsPerPR, decoded.AvgCommentsPerPR)
}

func TestQualityMetrics_JSONSerialization(t *testing.T) {
	qm := QualityMetrics{
		Period:          "2026-01",
		AIRatioBand:     AIRatioBandLow,
		TotalMergedPRs:  100,
		RevertedPRs:     5,
		HotfixPRs:       3,
		RevertRate:      0.05,
		HotfixRate:      0.03,
		AvgTimeToRevert: 48.0,
		AvgTimeToHotfix: 24.0,
	}

	data, err := json.Marshal(qm)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"revert_rate":0.05`)
	assert.Contains(t, string(data), `"hotfix_rate":0.03`)

	var decoded QualityMetrics
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, qm.RevertRate, decoded.RevertRate)
}

func TestCodeSurvivalRecord_JSONSerialization(t *testing.T) {
	ts := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	csr := CodeSurvivalRecord{
		CommitHash:    "abc123",
		FilePath:      "src/main.go",
		StartLine:     10,
		EndLine:       25,
		LinesAdded:    15,
		AIRatio:       0.75,
		AuthorID:      "user_001",
		AddedAt:       ts,
		SurvivedAt30d: true,
		SurvivedAt60d: true,
		SurvivedAt90d: false,
		ModifiedCount: 2,
		DeletedAt:     &ts,
	}

	data, err := json.Marshal(csr)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"survived_at_30d":true`)
	assert.Contains(t, string(data), `"survived_at_90d":false`)

	var decoded CodeSurvivalRecord
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, csr.SurvivedAt30d, decoded.SurvivedAt30d)
	assert.Equal(t, csr.SurvivedAt90d, decoded.SurvivedAt90d)
}

func TestCodeSurvivalRecord_SurvivalRate(t *testing.T) {
	csr := CodeSurvivalRecord{
		LinesAdded:    100,
		SurvivedAt30d: true,
		SurvivedAt60d: true,
		SurvivedAt90d: false,
	}

	// 30d window - survived
	assert.Equal(t, 1.0, csr.SurvivalRate(SurvivalWindow30d))

	// 60d window - survived
	assert.Equal(t, 1.0, csr.SurvivalRate(SurvivalWindow60d))

	// 90d window - not survived
	assert.Equal(t, 0.0, csr.SurvivalRate(SurvivalWindow90d))
}

func TestAIRatioBand_String(t *testing.T) {
	assert.Equal(t, "low", string(AIRatioBandLow))
	assert.Equal(t, "medium", string(AIRatioBandMedium))
	assert.Equal(t, "high", string(AIRatioBandHigh))
}

func TestSurvivalWindow_Days(t *testing.T) {
	assert.Equal(t, 30, SurvivalWindow30d.Days())
	assert.Equal(t, 60, SurvivalWindow60d.Days())
	assert.Equal(t, 90, SurvivalWindow90d.Days())
}
