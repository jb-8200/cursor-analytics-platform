package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVelocityConfig_CommitsPerDay(t *testing.T) {
	tests := []struct {
		name          string
		velocity      string
		prsPerWeek    float64
		expectedRange [2]float64 // min, max
	}{
		{
			name:          "low velocity, moderate PRs",
			velocity:      "low",
			prsPerWeek:    3.0,
			expectedRange: [2]float64{0.5, 1.0},
		},
		{
			name:          "medium velocity, moderate PRs",
			velocity:      "medium",
			prsPerWeek:    3.0,
			expectedRange: [2]float64{1.0, 1.5},
		},
		{
			name:          "high velocity, moderate PRs",
			velocity:      "high",
			prsPerWeek:    3.0,
			expectedRange: [2]float64{2.0, 3.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewVelocityConfig(tt.velocity)
			cpd := cfg.CommitsPerDay(tt.prsPerWeek)

			assert.GreaterOrEqual(t, cpd, tt.expectedRange[0])
			assert.LessOrEqual(t, cpd, tt.expectedRange[1])
		})
	}
}

func TestVelocityConfig_EventsPerDay(t *testing.T) {
	tests := []struct {
		name     string
		velocity string
		minRate  float64
	}{
		{"low velocity", "low", 50.0},
		{"medium velocity", "medium", 100.0},
		{"high velocity", "high", 150.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewVelocityConfig(tt.velocity)
			epd := cfg.EventsPerDay()

			assert.GreaterOrEqual(t, epd, tt.minRate)
		})
	}
}

func TestVelocityConfig_DefaultsToMedium(t *testing.T) {
	cfg := NewVelocityConfig("invalid")
	assert.NotNil(t, cfg)

	// Should behave like medium velocity
	cpd := cfg.CommitsPerDay(3.0)
	assert.Greater(t, cpd, 0.0)
}
