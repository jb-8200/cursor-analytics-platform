package generator

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/jb-8200/cursor-analytics-platform/services/cursor-sim/internal/config"
	"github.com/jb-8200/cursor-analytics-platform/services/cursor-sim/internal/models"
)

var (
	// AI models to rotate through
	aiModels = []string{
		"claude-3.5-sonnet",
		"claude-opus-4",
		"gpt-4-turbo",
		"gpt-4",
	}

	// File extensions with weights
	fileExtensions = []string{
		".ts", ".tsx", ".js", ".jsx", // TypeScript/JavaScript
		".py",              // Python
		".go",              // Go
		".java",            // Java
		".rb",              // Ruby
		".cpp", ".c", ".h", // C/C++
		".rs",                   // Rust
		".md", ".json", ".yaml", // Config
	}

	// Repository names
	repositories = []string{
		"backend-api",
		"frontend-app",
		"data-pipeline",
		"mobile-app",
		"analytics-service",
	}

	// Branch names
	branches = []string{
		"main",
		"develop",
		"feature/auth",
		"feature/dashboard",
		"bugfix/login",
		"hotfix/critical",
	}

	// Commit message templates
	commitMessages = []string{
		"feat: add %s functionality",
		"fix: resolve %s issue",
		"refactor: improve %s implementation",
		"chore: update %s dependencies",
		"docs: add %s documentation",
		"test: add %s tests",
		"perf: optimize %s performance",
	}

	commitFeatures = []string{
		"authentication", "user management", "data validation",
		"API endpoints", "database queries", "caching logic",
		"error handling", "logging", "monitoring",
	}
)

// EventGenerator generates commits and changes with realistic timing
type EventGenerator struct {
	config      *config.Config
	developers  []*models.Developer
	commitChan  chan *models.Commit
	changeChan  chan *models.Change
	commits     []*models.Commit
	changes     []*models.Change
	mu          sync.Mutex
	ctx         context.Context
	cancel      context.CancelFunc
	genWg       sync.WaitGroup // Waits for generator goroutines
	collectorWg sync.WaitGroup // Waits for collector goroutines
	stopped     bool
	rng         *rand.Rand
}

// NewEventGenerator creates a new event generator
func NewEventGenerator(cfg *config.Config, developers []*models.Developer) *EventGenerator {
	ctx, cancel := context.WithCancel(context.Background())
	return &EventGenerator{
		config:     cfg,
		developers: developers,
		commitChan: make(chan *models.Commit, 100),
		changeChan: make(chan *models.Change, 1000),
		commits:    make([]*models.Commit, 0),
		changes:    make([]*models.Change, 0),
		ctx:        ctx,
		cancel:     cancel,
		rng:        rand.New(rand.NewSource(cfg.Seed)),
	}
}

// Start begins event generation for all developers
func (e *EventGenerator) Start() error {
	baseLambda := VelocityToLambda(e.config.Velocity)

	// Start collector goroutines
	e.collectorWg.Add(2)
	go e.collectCommits()
	go e.collectChanges()

	// Start generator for each developer
	for i, dev := range e.developers {
		// Apply per-developer volatility
		devLambda := ApplyVolatility(baseLambda, e.config.Fluctuation, e.config.Seed+int64(i))

		e.genWg.Add(1)
		go e.generateForDeveloper(dev, devLambda, int64(i))
	}

	return nil
}

// Stop stops event generation
func (e *EventGenerator) Stop() {
	e.mu.Lock()
	if e.stopped {
		e.mu.Unlock()
		return
	}
	e.stopped = true
	e.mu.Unlock()

	// Signal all generators to stop
	e.cancel()

	// Wait for all generator goroutines to finish
	e.genWg.Wait()

	// Close channels to signal collectors
	close(e.commitChan)
	close(e.changeChan)

	// Wait for collectors to finish
	e.collectorWg.Wait()
}

// GetCommits returns all generated commits
func (e *EventGenerator) GetCommits() []*models.Commit {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]*models.Commit{}, e.commits...)
}

// GetChanges returns all generated changes
func (e *EventGenerator) GetChanges() []*models.Change {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]*models.Change{}, e.changes...)
}

// collectCommits collects commits from the channel
func (e *EventGenerator) collectCommits() {
	defer e.collectorWg.Done()
	for commit := range e.commitChan {
		e.mu.Lock()
		e.commits = append(e.commits, commit)
		e.mu.Unlock()
	}
}

// collectChanges collects changes from the channel
func (e *EventGenerator) collectChanges() {
	defer e.collectorWg.Done()
	for change := range e.changeChan {
		e.mu.Lock()
		e.changes = append(e.changes, change)
		e.mu.Unlock()
	}
}

// generateForDeveloper generates events for a single developer
func (e *EventGenerator) generateForDeveloper(dev *models.Developer, lambda float64, devIndex int64) {
	defer e.genWg.Done()

	// Create Poisson timer for this developer with unique seed
	timer := NewPoissonTimer(lambda, e.config.Seed+devIndex*1000)
	rng := rand.New(rand.NewSource(e.config.Seed + devIndex*1000 + 1))

	// Generate first event immediately (no delay)
	commit := e.generateCommit(dev, rng)
	changes := e.generateChanges(commit, dev, rng)

	// Send first event
	select {
	case e.commitChan <- commit:
	case <-e.ctx.Done():
		return
	}

	for _, change := range changes {
		select {
		case e.changeChan <- change:
		case <-e.ctx.Done():
			return
		}
	}

	// Continue with Poisson-timed events
	for {
		// Wait for next event time
		interval := timer.NextInterval()
		select {
		case <-time.After(interval):
			// Generate commit with changes
			commit := e.generateCommit(dev, rng)
			changes := e.generateChanges(commit, dev, rng)

			// Send to channels
			select {
			case e.commitChan <- commit:
			case <-e.ctx.Done():
				return
			}

			for _, change := range changes {
				select {
				case e.changeChan <- change:
				case <-e.ctx.Done():
					return
				}
			}

		case <-e.ctx.Done():
			return
		}
	}
}

// generateCommit creates a single commit
func (e *EventGenerator) generateCommit(dev *models.Developer, rng *rand.Rand) *models.Commit {
	now := time.Now().UTC()

	// Generate commit hash
	hash := fmt.Sprintf("%016x", rng.Uint64())

	// Generate commit message
	msgTemplate := commitMessages[rng.Intn(len(commitMessages))]
	feature := commitFeatures[rng.Intn(len(commitFeatures))]
	message := fmt.Sprintf(msgTemplate, feature)

	// Generate line counts
	linesFromTAB := rng.Intn(100) + 10
	linesFromComposer := rng.Intn(50) + 5
	linesNonAI := rng.Intn(30)
	totalLines := linesFromTAB + linesFromComposer + linesNonAI

	return &models.Commit{
		Hash:              hash,
		Timestamp:         now,
		Message:           message,
		UserID:            dev.ID,
		UserEmail:         dev.Email,
		Repository:        repositories[rng.Intn(len(repositories))],
		Branch:            branches[rng.Intn(len(branches))],
		TotalLines:        totalLines,
		LinesFromTAB:      linesFromTAB,
		LinesFromComposer: linesFromComposer,
		LinesNonAI:        linesNonAI,
		IngestionTime:     now,
	}
}

// generateChanges creates changes for a commit
func (e *EventGenerator) generateChanges(commit *models.Commit, dev *models.Developer, rng *rand.Rand) []*models.Change {
	// Generate 1-5 changes per commit
	numChanges := rng.Intn(5) + 1
	changes := make([]*models.Change, numChanges)

	// Default TAB vs COMPOSER ratio is 0.7
	tabRatio := 0.7

	for i := 0; i < numChanges; i++ {
		// Determine source (TAB or COMPOSER) based on ratio
		source := "TAB"
		if rng.Float64() > tabRatio {
			source = "COMPOSER"
		}

		// Generate change ID
		changeID := fmt.Sprintf("%s_%d", commit.Hash[:8], i)

		// Random file
		ext := fileExtensions[rng.Intn(len(fileExtensions))]
		filePath := fmt.Sprintf("src/components/Component%d%s", rng.Intn(100), ext)

		// Random model
		model := aiModels[rng.Intn(len(aiModels))]

		// Random line counts
		linesAdded := rng.Intn(50) + 1
		linesRemoved := rng.Intn(20)

		changes[i] = &models.Change{
			ChangeID:      changeID,
			CommitHash:    commit.Hash,
			UserID:        dev.ID,
			Timestamp:     commit.Timestamp,
			Source:        source,
			Model:         model,
			FilePath:      filePath,
			FileExtension: ext,
			LinesAdded:    linesAdded,
			LinesRemoved:  linesRemoved,
			IngestionTime: commit.IngestionTime,
		}
	}

	return changes
}
