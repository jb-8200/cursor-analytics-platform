package tui

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// ProgressBar represents a progress bar for tracking operations.
// It wraps the Bubbles progress component and provides a simple API.
type ProgressBar struct {
	title   string
	current int
	total   int
	writer  io.Writer
	mu      sync.RWMutex

	// Bubbles progress model
	bubblesProgress progress.Model
	program         *tea.Program
}

// NewProgressBar creates a new progress bar with the given title and total count.
// Writer can be nil; defaults to os.Stdout.
func NewProgressBar(title string, total int, writer io.Writer) *ProgressBar {
	if writer == nil {
		writer = os.Stdout
	}

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	pb := &ProgressBar{
		title:           title,
		current:         0,
		total:           total,
		writer:          writer,
		bubblesProgress: p,
	}

	return pb
}

// Update sets the current progress to the given value.
func (pb *ProgressBar) Update(current int) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	pb.current = current

	// Cap at total
	if pb.current > pb.total && pb.total > 0 {
		pb.current = pb.total
	}
}

// GetProgress returns the current progress value.
func (pb *ProgressBar) GetProgress() int {
	pb.mu.RLock()
	defer pb.mu.RUnlock()
	return pb.current
}

// GetPercentage returns the progress as a percentage (0-100).
func (pb *ProgressBar) GetPercentage() int {
	pb.mu.RLock()
	defer pb.mu.RUnlock()

	if pb.total == 0 {
		return 0
	}

	return (pb.current * 100) / pb.total
}

// SetTitle updates the progress bar title.
func (pb *ProgressBar) SetTitle(title string) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	pb.title = title
}

// Render returns a formatted string representation of the progress bar.
func (pb *ProgressBar) Render() string {
	pb.mu.RLock()
	defer pb.mu.RUnlock()

	if pb.total == 0 {
		return fmt.Sprintf("%s: N/A\n", pb.title)
	}

	// Calculate progress ratio (0.0 to 1.0)
	ratio := float64(pb.current) / float64(pb.total)

	// Clamp to [0, 1]
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	// Render progress bar using Bubbles
	barWidth := 40
	filled := int(float64(barWidth) * ratio)

	// Build progress bar string
	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "]"

	// Format output
	percentage := (pb.current * 100) / pb.total
	output := fmt.Sprintf("%s %s %d/%d (%d%%)\n",
		pb.title,
		bar,
		pb.current,
		pb.total,
		percentage,
	)

	return output
}

// progressBarModel implements the Bubble Tea Model interface.
type progressBarModel struct {
	title   string
	current int
	total   int
	prog    progress.Model
}

// Init implements tea.Model.
func (m progressBarModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m progressBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	default:
		return m, nil
	}
}

// View implements tea.Model.
func (m progressBarModel) View() string {
	ratio := 0.0
	if m.total > 0 {
		ratio = float64(m.current) / float64(m.total)
	}

	return fmt.Sprintf("%s %s %d/%d\n",
		m.title,
		m.prog.ViewAs(ratio),
		m.current,
		m.total,
	)
}
