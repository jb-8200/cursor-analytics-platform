package services

import (
	"math"
	"sort"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

// ResearchMetricsService calculates research metrics from data points.
type ResearchMetricsService struct {
	dataPoints []models.ResearchDataPoint
}

// NewResearchMetricsService creates a new research metrics service.
func NewResearchMetricsService(dataPoints []models.ResearchDataPoint) *ResearchMetricsService {
	return &ResearchMetricsService{
		dataPoints: dataPoints,
	}
}

// CalculateVelocityMetrics calculates velocity metrics grouped by AI ratio band.
func (s *ResearchMetricsService) CalculateVelocityMetrics(period string) []models.VelocityMetrics {
	if len(s.dataPoints) == 0 {
		return nil
	}

	// Group by AI ratio band
	groups := s.groupByAIRatioBand()

	var metrics []models.VelocityMetrics
	for band, points := range groups {
		if len(points) == 0 {
			continue
		}

		leadTimes := make([]float64, len(points))
		var totalAdditions, totalDeletions int

		for i, dp := range points {
			leadTimes[i] = dp.CodingLeadTimeHours
			totalAdditions += dp.Additions
			totalDeletions += dp.Deletions
		}

		metrics = append(metrics, models.VelocityMetrics{
			Period:              period,
			AIRatioBand:         band,
			TotalCommits:        len(points),
			TotalPRs:            countUniquePRs(points),
			TotalAdditions:      totalAdditions,
			TotalDeletions:      totalDeletions,
			AvgLeadTimeHours:    mean(leadTimes),
			MedianLeadTimeHours: median(leadTimes),
			StdDevLeadTimeHours: stdDev(leadTimes),
		})
	}

	return metrics
}

// CalculateReviewCostMetrics calculates review cost metrics grouped by AI ratio band.
func (s *ResearchMetricsService) CalculateReviewCostMetrics(period string) []models.ReviewCostMetrics {
	if len(s.dataPoints) == 0 {
		return nil
	}

	// Group by AI ratio band
	groups := s.groupByAIRatioBand()

	var metrics []models.ReviewCostMetrics
	for band, points := range groups {
		if len(points) == 0 {
			continue
		}

		reviewTimes := make([]float64, len(points))
		var totalIterations int

		for i, dp := range points {
			reviewTimes[i] = dp.ReviewLeadTimeHours
			totalIterations += dp.ReviewIterations
		}

		prCount := countUniquePRs(points)

		metrics = append(metrics, models.ReviewCostMetrics{
			Period:                period,
			AIRatioBand:           band,
			TotalPRsReviewed:      prCount,
			TotalReviewIterations: totalIterations,
			AvgIterationsPerPR:    float64(totalIterations) / float64(prCount),
			AvgReviewTimeHours:    mean(reviewTimes),
			MedianReviewTimeHours: median(reviewTimes),
			StdDevReviewTimeHours: stdDev(reviewTimes),
		})
	}

	return metrics
}

// CalculateQualityMetrics calculates quality metrics grouped by AI ratio band.
func (s *ResearchMetricsService) CalculateQualityMetrics(period string) []models.QualityMetrics {
	if len(s.dataPoints) == 0 {
		return nil
	}

	// Group by AI ratio band
	groups := s.groupByAIRatioBand()

	var metrics []models.QualityMetrics
	for band, points := range groups {
		if len(points) == 0 {
			continue
		}

		var reverted, hotfixes int
		for _, dp := range points {
			if dp.WasReverted {
				reverted++
			}
			if dp.RequiredHotfix {
				hotfixes++
			}
		}

		prCount := countUniquePRs(points)

		metrics = append(metrics, models.QualityMetrics{
			Period:         period,
			AIRatioBand:    band,
			TotalMergedPRs: prCount,
			RevertedPRs:    reverted,
			HotfixPRs:      hotfixes,
			RevertRate:     float64(reverted) / float64(prCount),
			HotfixRate:     float64(hotfixes) / float64(prCount),
		})
	}

	return metrics
}

// groupByAIRatioBand groups data points by their AI ratio band.
func (s *ResearchMetricsService) groupByAIRatioBand() map[models.AIRatioBand][]models.ResearchDataPoint {
	groups := make(map[models.AIRatioBand][]models.ResearchDataPoint)

	for _, dp := range s.dataPoints {
		band := dp.GetAIRatioBand()
		groups[band] = append(groups[band], dp)
	}

	return groups
}

// countUniquePRs counts unique PR numbers in the data points.
func countUniquePRs(points []models.ResearchDataPoint) int {
	seen := make(map[int]bool)
	for _, dp := range points {
		if dp.PRNumber > 0 {
			seen[dp.PRNumber] = true
		}
	}
	// If no PRs, count as data points
	if len(seen) == 0 {
		return len(points)
	}
	return len(seen)
}

// mean calculates the arithmetic mean of a slice of floats.
func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// median calculates the median of a slice of floats.
func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

// stdDev calculates the standard deviation of a slice of floats.
func stdDev(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	m := mean(values)
	var sumSquares float64
	for _, v := range values {
		sumSquares += (v - m) * (v - m)
	}

	return math.Sqrt(sumSquares / float64(len(values)))
}
