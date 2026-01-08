package tui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	figure "github.com/common-nighthawk/go-figure"
	"github.com/lucasb-eyer/go-colorful"
)

// DisplayBanner renders the DOXAPI ASCII banner with gradient if terminal supports colors.
// Uses ShouldUseTUI() to determine if colors should be applied.
func DisplayBanner(version string) {
	displayBannerTo(version, nil, ShouldUseTUI())
}

// displayBannerTo is the internal implementation that accepts a writer and color flag.
// This allows testing with custom writers and color settings.
func displayBannerTo(version string, writer io.Writer, useColors bool) {
	if writer == nil {
		writer = os.Stdout
	}

	if !useColors {
		// Plain text fallback for non-TTY or NO_COLOR environments
		fmt.Fprintf(writer, "DOXAPI v%s\n\n", version)
		return
	}

	// Generate ASCII art using go-figure
	fig := figure.NewFigure("DOXAPI", "standard", true)
	asciiArt := fig.String()
	lines := strings.Split(asciiArt, "\n")

	// Apply purple-to-pink gradient per line
	for i, line := range lines {
		if line == "" {
			continue
		}

		// Calculate ratio (0 = purple, 1 = pink)
		ratio := 0.0
		if len(lines) > 1 {
			ratio = float64(i) / float64(len(lines)-1)
		}

		// Get interpolated color and render line
		hexColor := interpolateColor(ratio)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(hexColor))
		fmt.Fprintln(writer, style.Render(line))
	}

	// Version info below banner
	versionStyle := SubtitleStyle
	fmt.Fprintf(writer, "%s\n", versionStyle.Render(fmt.Sprintf("v%s", version)))
	fmt.Fprintln(writer)
}

// interpolateColor returns a hex color code interpolated between purple and pink.
// ratio: 0 = purple (#9B59B6), 1 = pink (#FF69B4), 0.5 = midpoint blend
func interpolateColor(ratio float64) string {
	// Parse purple and pink hex colors
	purple, err := colorful.Hex("#9B59B6")
	if err != nil {
		return "#9B59B6" // Fallback to purple on error
	}

	pink, err := colorful.Hex("#FF69B4")
	if err != nil {
		return "#FF69B4" // Fallback to pink on error
	}

	// Clamp ratio to [0, 1]
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	// Interpolate using Lab color space for better perceptual blend
	blended := purple.BlendLab(pink, ratio)
	return blended.Hex()
}
