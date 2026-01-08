package tui

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Spinner wraps a Bubbles spinner for showing loading states.
// It handles both TTY (animated) and non-TTY (text-based) environments gracefully.
type Spinner struct {
	message   string
	isRunning bool
	writer    io.Writer
	mu        sync.RWMutex

	// TTY mode fields
	bubblesSpinner spinner.Model
	program        *tea.Program
	done           chan bool

	// Non-TTY fallback
	stopChan chan struct{}
	stopOnce sync.Once
}

// NewSpinner creates a new spinner with the given message.
// Writer can be nil; defaults to os.Stdout.
func NewSpinner(message string, writer io.Writer) *Spinner {
	if writer == nil {
		writer = os.Stdout
	}

	s := &Spinner{
		message:   message,
		isRunning: false,
		writer:    writer,
		done:      make(chan bool, 1),
		stopChan:  make(chan struct{}),
	}

	// Initialize Bubbles spinner
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	s.bubblesSpinner = spin

	return s
}

// Start begins the spinner animation (TTY mode) or status output (non-TTY mode).
func (s *Spinner) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return // Already running
	}

	s.isRunning = true

	if ShouldUseTUI() {
		// TTY mode: use animated spinner
		s.startTTYSpinner()
	} else {
		// Non-TTY mode: text-based fallback
		s.startNonTTYSpinner()
	}
}

// startTTYSpinner starts the animated Bubbles spinner (TTY environments).
func (s *Spinner) startTTYSpinner() {
	// Create a model that wraps our spinner state
	m := spinnerModel{
		spinner: s.bubblesSpinner,
		message: s.message,
		done:    s.done,
	}

	p := tea.NewProgram(m, tea.WithoutRenderer())
	s.program = p

	go func() {
		if _, err := p.Run(); err != nil {
			// Silently ignore errors in spinner
		}
	}()
}

// startNonTTYSpinner starts text-based status output (non-TTY environments).
func (s *Spinner) startNonTTYSpinner() {
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		dots := 0
		for {
			select {
			case <-ticker.C:
				// Periodically output status (disabled to reduce noise in tests)
				// Uncomment for visible output:
				// fmt.Fprintf(s.writer, "⏳ %s%s\r", s.getMessage(), strings.Repeat(".", (dots%3)+1))
				dots++
			case <-s.stopChan:
				return
			}
		}
	}()
}

// Stop stops the spinner and displays the completion message.
func (s *Spinner) Stop(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return // Not running, nothing to stop
	}

	s.isRunning = false

	if s.program != nil {
		// TTY mode: send quit message to program
		s.program.Send(stopMsg{message: message})
		s.program.Quit()
		// Wait for program to finish
		<-s.done
		s.program = nil
	} else {
		// Non-TTY mode: stop the ticker
		s.stopOnce.Do(func() {
			close(s.stopChan)
		})
		// Output completion message
		fmt.Fprintf(s.writer, "✅ %s\n", message)
	}
}

// UpdateMessage updates the spinner message while it's running.
func (s *Spinner) UpdateMessage(newMessage string) {
	s.mu.Lock()
	s.message = newMessage
	s.mu.Unlock()

	if s.program != nil {
		s.program.Send(updateMsg{message: newMessage})
	}
}

// getMessage safely retrieves the current message.
func (s *Spinner) getMessage() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.message
}

// spinnerModel implements the Bubble Tea Model interface for the spinner.
type spinnerModel struct {
	spinner spinner.Model
	message string
	done    chan bool
	quitting bool
}

// Init implements tea.Model.
func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update implements tea.Model.
func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case stopMsg:
		m.quitting = true
		return m, tea.Quit
	case updateMsg:
		um := msg.(updateMsg)
		m.message = um.message
		return m, m.spinner.Tick
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

// View implements tea.Model.
func (m spinnerModel) View() string {
	if m.quitting {
		// Signal completion
		go func() {
			m.done <- true
		}()
		return ""
	}
	return fmt.Sprintf("%s %s\n", m.spinner.View(), m.message)
}

// stopMsg signals the spinner to stop.
type stopMsg struct {
	message string
}

// updateMsg signals the spinner to update its message.
type updateMsg struct {
	message string
}
