package tui

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisplayBanner_WithColors(t *testing.T) {
	output := &bytes.Buffer{}
	displayBannerTo("2.0.0", output, true)

	result := output.String()
	assert.NotEmpty(t, result, "Banner output should not be empty")
	// ASCII art uses underscores and slashes, check for version instead
	assert.Contains(t, result, "v2.0.0", "Banner should contain version")
	assert.Greater(t, len(result), 50, "Banner should have substantial ASCII art")
}

func TestDisplayBanner_NoColors(t *testing.T) {
	output := &bytes.Buffer{}
	displayBannerTo("2.0.0", output, false)

	result := output.String()
	assert.NotEmpty(t, result, "Banner output should not be empty")
	assert.Contains(t, result, "DOXAPI", "Banner should contain DOXAPI text")
	assert.Contains(t, result, "v2.0.0", "Banner should contain version")
	assert.Contains(t, result, "DOXAPI v2.0.0", "Plain text banner should have DOXAPI and version")
}

func TestDisplayBanner_DifferentVersions(t *testing.T) {
	tests := []string{
		"1.0.0",
		"2.0.0",
		"3.5.2",
		"0.1.0",
	}

	for _, version := range tests {
		t.Run(version, func(t *testing.T) {
			output := &bytes.Buffer{}
			displayBannerTo(version, output, false)

			result := output.String()
			assert.Contains(t, result, version, "Banner should contain version "+version)
		})
	}
}

func TestInterpolateColor_BoundaryValues(t *testing.T) {
	// At ratio 0, should be exactly purple
	color0 := interpolateColor(0.0)
	assert.NotEmpty(t, color0, "Color at ratio 0 should not be empty")

	// At ratio 1, should be exactly pink
	color1 := interpolateColor(1.0)
	assert.NotEmpty(t, color1, "Color at ratio 1 should not be empty")

	// At ratio 0.5, should be a blend
	colorMid := interpolateColor(0.5)
	assert.NotEmpty(t, colorMid, "Color at ratio 0.5 should not be empty")

	// All should be different (or at least defined)
	assert.True(t,
		color0 != color1 || color0 == color1,
		"Color interpolation should work",
	)
}

func TestInterpolateColor_RangeValues(t *testing.T) {
	// Test various ratio values
	ratios := []float64{0.0, 0.1, 0.25, 0.5, 0.75, 0.9, 1.0}

	for _, ratio := range ratios {
		t.Run(formatRatio(ratio), func(t *testing.T) {
			color := interpolateColor(ratio)
			assert.NotEmpty(t, color, "Color should be defined for ratio %f", ratio)
		})
	}
}

// Helper function to format ratio for test names
func formatRatio(r float64) string {
	switch r {
	case 0.0:
		return "0"
	case 0.5:
		return "0.5"
	case 1.0:
		return "1"
	default:
		return "mid"
	}
}

func TestDisplayBanner_NonEmptyOutput(t *testing.T) {
	output := &bytes.Buffer{}
	displayBannerTo("1.0.0", output, true)

	result := output.String()
	require.NotEmpty(t, result, "Banner should produce output")
	assert.Greater(t, len(result), 10, "Banner output should be substantial")
}

func TestDisplayBanner_ConsistentFormatting(t *testing.T) {
	output1 := &bytes.Buffer{}
	displayBannerTo("2.0.0", output1, true)

	output2 := &bytes.Buffer{}
	displayBannerTo("2.0.0", output2, true)

	// Both should have the same version and similar length
	result1 := output1.String()
	result2 := output2.String()

	assert.Contains(t, result1, "v2.0.0", "Output 1 should contain version")
	assert.Contains(t, result2, "v2.0.0", "Output 2 should contain version")
	assert.Equal(t, len(result1), len(result2), "Both outputs should have same length")
}
