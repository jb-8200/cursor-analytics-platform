// Package services provides business logic services for cursor-sim.
// This file implements the Qualtrics export job state machine.
package services

import (
	"encoding/hex"
	"fmt"
	mathrand "math/rand"
	"sync"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

// ExportJobManager manages Qualtrics survey export jobs and their state transitions.
// It provides thread-safe operations for starting exports, tracking progress, and retrieving files.
type ExportJobManager struct {
	generator *generator.SurveyGenerator
	mu        sync.RWMutex
	jobs      map[string]*models.ExportJob // progressID -> job
	files     map[string][]byte            // fileID -> ZIP data
	rng       *mathrand.Rand
}

// NewExportJobManager creates a new export job manager with the given survey generator.
func NewExportJobManager(gen *generator.SurveyGenerator) *ExportJobManager {
	return &ExportJobManager{
		generator: gen,
		jobs:      make(map[string]*models.ExportJob),
		files:     make(map[string][]byte),
		rng:       mathrand.New(mathrand.NewSource(time.Now().UnixNano())),
	}
}

// StartExport creates a new export job for the given survey ID.
// Returns a job with status inProgress and 0% complete.
func (m *ExportJobManager) StartExport(surveyID string) (*models.ExportJob, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate unique progress ID
	progressID := m.generateProgressID()

	// Create job
	job := &models.ExportJob{
		ProgressID:      progressID,
		SurveyID:        surveyID,
		Status:          models.ExportStatusInProgress,
		PercentComplete: 0,
		CreatedAt:       time.Now(),
	}

	// Store job
	m.jobs[progressID] = job

	// Return a copy to avoid external modification
	jobCopy := *job
	return &jobCopy, nil
}

// GetProgress retrieves the current status of an export job and advances progress.
// Each call advances progress by 20% until reaching 100%, at which point the job completes.
func (m *ExportJobManager) GetProgress(progressID string) (*models.ExportJob, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if job exists
	job, exists := m.jobs[progressID]
	if !exists {
		return nil, fmt.Errorf("export job not found: %s", progressID)
	}

	// If already complete, return current state
	if job.Status == models.ExportStatusComplete {
		jobCopy := *job
		return &jobCopy, nil
	}

	// Advance progress by 20%
	job.PercentComplete += 20

	// Check if complete
	if job.PercentComplete >= 100 {
		job.PercentComplete = 100
		job.Status = models.ExportStatusComplete

		// Generate file
		fileID, zipData, err := m.generateFile(job.SurveyID)
		if err != nil {
			job.Status = models.ExportStatusFailed
			return nil, fmt.Errorf("failed to generate export file: %w", err)
		}

		job.FileID = fileID
		m.files[fileID] = zipData
	}

	// Return a copy
	jobCopy := *job
	return &jobCopy, nil
}

// GetFile retrieves the ZIP file data for a completed export.
func (m *ExportJobManager) GetFile(fileID string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, exists := m.files[fileID]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", fileID)
	}

	// Return a copy to prevent external modification
	result := make([]byte, len(data))
	copy(result, data)
	return result, nil
}

// generateProgressID generates a unique progress ID in the format ES_xxxxx.
func (m *ExportJobManager) generateProgressID() string {
	bytes := make([]byte, 8)
	for i := range bytes {
		bytes[i] = byte(m.rng.Intn(256))
	}
	return "ES_" + hex.EncodeToString(bytes)
}

// generateFile generates survey responses and creates a ZIP file.
// Returns the file ID and ZIP data.
func (m *ExportJobManager) generateFile(surveyID string) (string, []byte, error) {
	// Generate responses using the survey generator
	responses := m.generator.GenerateResponses(surveyID)

	// Generate ZIP file
	zipData, err := models.GenerateZIPFile(responses)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create ZIP file: %w", err)
	}

	// Generate file ID
	fileBytes := make([]byte, 8)
	for i := range fileBytes {
		fileBytes[i] = byte(m.rng.Intn(256))
	}
	fileID := "FILE_" + hex.EncodeToString(fileBytes)

	return fileID, zipData, nil
}
