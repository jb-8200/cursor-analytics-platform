package generator

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPoissonTimer(t *testing.T) {
	timer := NewPoissonTimer(10.0, 12345)
	require.NotNil(t, timer)
	assert.NotNil(t, timer.rng)
	assert.Equal(t, 10.0, timer.lambda)
}

func TestPoissonTimer_NextInterval(t *testing.T) {
	timer := NewPoissonTimer(10.0, 12345)

	// Generate multiple intervals
	intervals := make([]time.Duration, 1000)
	for i := 0; i < 1000; i++ {
		intervals[i] = timer.NextInterval()
		assert.Greater(t, intervals[i], time.Duration(0), "Interval should be positive")
	}

	// Check that intervals vary (not all the same)
	uniqueIntervals := make(map[time.Duration]bool)
	for _, interval := range intervals {
		uniqueIntervals[interval] = true
	}
	assert.Greater(t, len(uniqueIntervals), 100, "Should have variety in intervals")
}

func TestPoissonTimer_MeanInterval(t *testing.T) {
	// Test that mean interval approximates expected value
	// For Poisson process with rate λ events/hour, mean interval = 1/λ hours

	tests := []struct {
		name           string
		lambda         float64 // events per hour
		expectedMeanMs float64 // milliseconds
	}{
		{
			name:           "10 events/hour",
			lambda:         10.0,
			expectedMeanMs: 360000.0, // 3600000ms / 10 = 360000ms = 6 minutes
		},
		{
			name:           "60 events/hour",
			lambda:         60.0,
			expectedMeanMs: 60000.0, // 3600000ms / 60 = 60000ms = 1 minute
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := NewPoissonTimer(tt.lambda, 12345)

			// Generate many samples
			var totalMs float64
			samples := 10000
			for i := 0; i < samples; i++ {
				interval := timer.NextInterval()
				totalMs += float64(interval.Milliseconds())
			}

			meanMs := totalMs / float64(samples)

			// Allow 10% variance from expected mean
			tolerance := tt.expectedMeanMs * 0.10
			assert.InDelta(t, tt.expectedMeanMs, meanMs, tolerance,
				"Mean interval should approximate expected value")
		})
	}
}

func TestPoissonTimer_Deterministic(t *testing.T) {
	// Same seed should produce same sequence
	timer1 := NewPoissonTimer(10.0, 99999)
	timer2 := NewPoissonTimer(10.0, 99999)

	for i := 0; i < 10; i++ {
		interval1 := timer1.NextInterval()
		interval2 := timer2.NextInterval()
		assert.Equal(t, interval1, interval2,
			"Same seed should produce same intervals at iteration %d", i)
	}
}

func TestPoissonTimer_ExponentialDistribution(t *testing.T) {
	// Poisson intervals follow exponential distribution
	// We can check if the distribution roughly matches by looking at percentiles

	lambda := 60.0 // 60 events/hour
	timer := NewPoissonTimer(lambda, 12345)

	samples := 10000
	intervals := make([]float64, samples)
	for i := 0; i < samples; i++ {
		intervals[i] = float64(timer.NextInterval().Milliseconds())
	}

	// Sort intervals for percentile calculation
	// Using simple bubble sort for small dataset
	for i := 0; i < len(intervals)-1; i++ {
		for j := 0; j < len(intervals)-i-1; j++ {
			if intervals[j] > intervals[j+1] {
				intervals[j], intervals[j+1] = intervals[j+1], intervals[j]
			}
		}
	}

	// For exponential distribution with mean μ:
	// - Median ≈ 0.693 * μ
	// - 90th percentile ≈ 2.303 * μ

	meanMs := 3600000.0 / lambda // milliseconds
	median := intervals[len(intervals)/2]
	p90 := intervals[int(float64(len(intervals))*0.9)]

	expectedMedian := 0.693 * meanMs
	expectedP90 := 2.303 * meanMs

	// Allow 20% tolerance for statistical variation
	assert.InDelta(t, expectedMedian, median, expectedMedian*0.20,
		"Median should approximate exponential distribution median")
	assert.InDelta(t, expectedP90, p90, expectedP90*0.20,
		"90th percentile should approximate exponential distribution p90")
}

func TestVelocityToLambda(t *testing.T) {
	tests := []struct {
		velocity string
		expected float64
	}{
		{"low", 5.0},
		{"medium", 25.0},
		{"high", 50.0},
	}

	for _, tt := range tests {
		t.Run(tt.velocity, func(t *testing.T) {
			lambda := VelocityToLambda(tt.velocity)
			assert.Equal(t, tt.expected, lambda)
		})
	}
}

func TestVelocityToLambda_InvalidVelocity(t *testing.T) {
	// Should default to medium for unknown velocity
	lambda := VelocityToLambda("invalid")
	assert.Equal(t, 25.0, lambda, "Unknown velocity should default to medium")
}

func TestApplyVolatility(t *testing.T) {
	baseLambda := 50.0
	volatility := 0.3
	seed := int64(12345)

	// Generate adjusted lambda values for multiple developers
	adjustedLambdas := make([]float64, 100)
	for i := 0; i < 100; i++ {
		adjustedLambdas[i] = ApplyVolatility(baseLambda, volatility, seed+int64(i))
	}

	// Check that values are within expected range
	minExpected := baseLambda * (1.0 - volatility)
	maxExpected := baseLambda * (1.0 + volatility)

	for i, lambda := range adjustedLambdas {
		assert.GreaterOrEqual(t, lambda, minExpected,
			"Adjusted lambda should be >= min at index %d", i)
		assert.LessOrEqual(t, lambda, maxExpected,
			"Adjusted lambda should be <= max at index %d", i)
	}

	// Check that there's variety (not all the same)
	uniqueValues := make(map[float64]bool)
	for _, lambda := range adjustedLambdas {
		// Round to 2 decimal places for comparison
		rounded := math.Round(lambda*100) / 100
		uniqueValues[rounded] = true
	}
	assert.Greater(t, len(uniqueValues), 50,
		"Should have variety in adjusted lambda values")
}
