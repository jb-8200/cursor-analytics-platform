package services

import (
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResearchMetricsService_CalculateVelocityMetrics(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		// Low AI ratio (0.1)
		{AIRatio: 0.1, Additions: 100, Deletions: 20, CodingLeadTimeHours: 4.0, Timestamp: baseTime},
		{AIRatio: 0.2, Additions: 150, Deletions: 30, CodingLeadTimeHours: 6.0, Timestamp: baseTime.Add(1 * time.Hour)},
		// High AI ratio (0.8)
		{AIRatio: 0.8, Additions: 200, Deletions: 10, CodingLeadTimeHours: 2.0, Timestamp: baseTime.Add(2 * time.Hour)},
		{AIRatio: 0.9, Additions: 180, Deletions: 15, CodingLeadTimeHours: 3.0, Timestamp: baseTime.Add(3 * time.Hour)},
	}

	svc := NewResearchMetricsService(dataPoints)

	metrics := svc.CalculateVelocityMetrics("2026-01")

	// Should have metrics for low and high bands
	require.Len(t, metrics, 2)

	// Find low band metrics
	var lowMetrics, highMetrics *models.VelocityMetrics
	for i := range metrics {
		if metrics[i].AIRatioBand == models.AIRatioBandLow {
			lowMetrics = &metrics[i]
		} else if metrics[i].AIRatioBand == models.AIRatioBandHigh {
			highMetrics = &metrics[i]
		}
	}

	require.NotNil(t, lowMetrics)
	require.NotNil(t, highMetrics)

	// Verify low band aggregates
	assert.Equal(t, 2, lowMetrics.TotalCommits)
	assert.Equal(t, 250, lowMetrics.TotalAdditions)   // 100 + 150
	assert.Equal(t, 50, lowMetrics.TotalDeletions)    // 20 + 30
	assert.Equal(t, 5.0, lowMetrics.AvgLeadTimeHours) // (4 + 6) / 2

	// Verify high band aggregates
	assert.Equal(t, 2, highMetrics.TotalCommits)
	assert.Equal(t, 380, highMetrics.TotalAdditions)   // 200 + 180
	assert.Equal(t, 2.5, highMetrics.AvgLeadTimeHours) // (2 + 3) / 2
}

func TestResearchMetricsService_CalculateReviewCostMetrics(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		// Medium AI ratio with reviews
		{AIRatio: 0.4, ReviewIterations: 2, ReviewLeadTimeHours: 8.0, PRNumber: 1, Timestamp: baseTime},
		{AIRatio: 0.5, ReviewIterations: 1, ReviewLeadTimeHours: 4.0, PRNumber: 2, Timestamp: baseTime.Add(1 * time.Hour)},
		// High AI ratio with reviews
		{AIRatio: 0.8, ReviewIterations: 3, ReviewLeadTimeHours: 12.0, PRNumber: 3, Timestamp: baseTime.Add(2 * time.Hour)},
	}

	svc := NewResearchMetricsService(dataPoints)

	metrics := svc.CalculateReviewCostMetrics("2026-01")

	// Should have metrics for medium and high bands
	require.Len(t, metrics, 2)

	// Find medium band metrics
	var mediumMetrics *models.ReviewCostMetrics
	for i := range metrics {
		if metrics[i].AIRatioBand == models.AIRatioBandMedium {
			mediumMetrics = &metrics[i]
		}
	}

	require.NotNil(t, mediumMetrics)
	assert.Equal(t, 2, mediumMetrics.TotalPRsReviewed)
	assert.Equal(t, 3, mediumMetrics.TotalReviewIterations) // 2 + 1
	assert.Equal(t, 1.5, mediumMetrics.AvgIterationsPerPR)  // 3 / 2
	assert.Equal(t, 6.0, mediumMetrics.AvgReviewTimeHours)  // (8 + 4) / 2
}

func TestResearchMetricsService_CalculateQualityMetrics(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		// Low AI ratio
		{AIRatio: 0.1, WasReverted: false, RequiredHotfix: false, PRNumber: 1, Timestamp: baseTime},
		{AIRatio: 0.2, WasReverted: true, RequiredHotfix: false, PRNumber: 2, Timestamp: baseTime.Add(1 * time.Hour)},
		// High AI ratio
		{AIRatio: 0.8, WasReverted: true, RequiredHotfix: true, PRNumber: 3, Timestamp: baseTime.Add(2 * time.Hour)},
		{AIRatio: 0.9, WasReverted: false, RequiredHotfix: false, PRNumber: 4, Timestamp: baseTime.Add(3 * time.Hour)},
	}

	svc := NewResearchMetricsService(dataPoints)

	metrics := svc.CalculateQualityMetrics("2026-01")

	// Should have metrics for low and high bands
	require.Len(t, metrics, 2)

	// Find low band metrics
	var lowMetrics, highMetrics *models.QualityMetrics
	for i := range metrics {
		if metrics[i].AIRatioBand == models.AIRatioBandLow {
			lowMetrics = &metrics[i]
		} else if metrics[i].AIRatioBand == models.AIRatioBandHigh {
			highMetrics = &metrics[i]
		}
	}

	require.NotNil(t, lowMetrics)
	require.NotNil(t, highMetrics)

	// Verify low band quality
	assert.Equal(t, 2, lowMetrics.TotalMergedPRs)
	assert.Equal(t, 1, lowMetrics.RevertedPRs)
	assert.Equal(t, 0.5, lowMetrics.RevertRate) // 1/2

	// Verify high band quality
	assert.Equal(t, 2, highMetrics.TotalMergedPRs)
	assert.Equal(t, 1, highMetrics.RevertedPRs)
	assert.Equal(t, 1, highMetrics.HotfixPRs)
	assert.Equal(t, 0.5, highMetrics.HotfixRate) // 1/2
}

func TestResearchMetricsService_EmptyData(t *testing.T) {
	svc := NewResearchMetricsService(nil)

	velocityMetrics := svc.CalculateVelocityMetrics("2026-01")
	assert.Empty(t, velocityMetrics)

	reviewMetrics := svc.CalculateReviewCostMetrics("2026-01")
	assert.Empty(t, reviewMetrics)

	qualityMetrics := svc.CalculateQualityMetrics("2026-01")
	assert.Empty(t, qualityMetrics)
}

func TestResearchMetricsService_StatisticalAggregations(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{AIRatio: 0.1, CodingLeadTimeHours: 2.0, Timestamp: baseTime},
		{AIRatio: 0.1, CodingLeadTimeHours: 4.0, Timestamp: baseTime.Add(1 * time.Hour)},
		{AIRatio: 0.1, CodingLeadTimeHours: 6.0, Timestamp: baseTime.Add(2 * time.Hour)},
		{AIRatio: 0.1, CodingLeadTimeHours: 8.0, Timestamp: baseTime.Add(3 * time.Hour)},
	}

	svc := NewResearchMetricsService(dataPoints)

	metrics := svc.CalculateVelocityMetrics("2026-01")
	require.Len(t, metrics, 1)

	lowMetrics := metrics[0]

	// Mean: (2 + 4 + 6 + 8) / 4 = 5
	assert.Equal(t, 5.0, lowMetrics.AvgLeadTimeHours)

	// Median of [2, 4, 6, 8] = (4 + 6) / 2 = 5
	assert.Equal(t, 5.0, lowMetrics.MedianLeadTimeHours)

	// Standard deviation should be calculated
	assert.Greater(t, lowMetrics.StdDevLeadTimeHours, 0.0)
}

func TestResearchMetricsService_GroupByAIRatioBand(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{AIRatio: 0.1, PRNumber: 1, Timestamp: baseTime},  // Low
		{AIRatio: 0.29, PRNumber: 2, Timestamp: baseTime}, // Low (boundary)
		{AIRatio: 0.30, PRNumber: 3, Timestamp: baseTime}, // Medium (boundary)
		{AIRatio: 0.5, PRNumber: 4, Timestamp: baseTime},  // Medium
		{AIRatio: 0.69, PRNumber: 5, Timestamp: baseTime}, // Medium (boundary)
		{AIRatio: 0.70, PRNumber: 6, Timestamp: baseTime}, // High (boundary)
		{AIRatio: 0.9, PRNumber: 7, Timestamp: baseTime},  // High
	}

	svc := NewResearchMetricsService(dataPoints)

	metrics := svc.CalculateVelocityMetrics("2026-01")
	require.Len(t, metrics, 3) // All three bands should be present

	// Count by band
	counts := make(map[models.AIRatioBand]int)
	for _, m := range metrics {
		counts[m.AIRatioBand] = m.TotalCommits
	}

	assert.Equal(t, 2, counts[models.AIRatioBandLow])    // 0.1, 0.29
	assert.Equal(t, 3, counts[models.AIRatioBandMedium]) // 0.30, 0.5, 0.69
	assert.Equal(t, 2, counts[models.AIRatioBandHigh])   // 0.70, 0.9
}
