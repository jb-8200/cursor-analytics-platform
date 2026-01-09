// Package storage provides in-memory storage implementations for external data sources.
// This file implements storage for Harvey AI, Microsoft 365 Copilot, and Qualtrics data.
package storage

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// Harvey Types
// ============================================================================

// HarveyTask represents the type of AI assistant task.
type HarveyTask string

const (
	HarveyTaskAssist   HarveyTask = "Assist"   // General questions
	HarveyTaskDraft    HarveyTask = "Draft"    // Document drafting
	HarveyTaskReview   HarveyTask = "Review"   // Contract review
	HarveyTaskResearch HarveyTask = "Research" // Legal research
)

// HarveySource represents the data source for the task.
type HarveySource string

const (
	HarveySourceFiles     HarveySource = "Files"     // Uploaded documents
	HarveySourceWeb       HarveySource = "Web"       // Web search
	HarveySourceKnowledge HarveySource = "Knowledge" // Knowledge base
)

// HarveySentiment represents feedback sentiment.
type HarveySentiment string

const (
	HarveySentimentPositive HarveySentiment = "positive"
	HarveySentimentNegative HarveySentiment = "negative"
	HarveySentimentNeutral  HarveySentiment = "neutral"
)

// HarveyUsage represents a single Harvey AI usage event.
type HarveyUsage struct {
	EventID           int64           `json:"event_id"`
	MessageID         string          `json:"message_ID"`
	Time              time.Time       `json:"Time"`
	User              string          `json:"User"`
	Task              HarveyTask      `json:"Task"`
	ClientMatter      float64         `json:"Client Matter #"`
	Source            HarveySource    `json:"Source"`
	NumberOfDocuments int             `json:"Number of documents"`
	FeedbackComments  string          `json:"Feedback Comments"`
	FeedbackSentiment HarveySentiment `json:"Feedback Sentiment"`
}

// HarveyParams defines query parameters for Harvey usage data.
type HarveyParams struct {
	From time.Time
	To   time.Time
	User string     // Optional: filter by user email
	Task HarveyTask // Optional: filter by task type
}

// ============================================================================
// Copilot Types
// ============================================================================

// CopilotPeriod represents valid report periods.
type CopilotPeriod string

const (
	CopilotPeriodD7   CopilotPeriod = "D7"
	CopilotPeriodD30  CopilotPeriod = "D30"
	CopilotPeriodD90  CopilotPeriod = "D90"
	CopilotPeriodD180 CopilotPeriod = "D180"
	CopilotPeriodAll  CopilotPeriod = "ALL"
)

// Days returns the number of days for this period.
func (p CopilotPeriod) Days() int {
	switch p {
	case CopilotPeriodD7:
		return 7
	case CopilotPeriodD30:
		return 30
	case CopilotPeriodD90:
		return 90
	case CopilotPeriodD180:
		return 180
	case CopilotPeriodAll:
		return 180 // ALL includes all available data
	default:
		return 30
	}
}

// CopilotUsage represents user-level Copilot usage data.
// Matches Microsoft Graph API beta response schema.
type CopilotUsage struct {
	ReportRefreshDate                     string  `json:"reportRefreshDate"`
	ReportPeriod                          int     `json:"reportPeriod"`
	UserPrincipalName                     string  `json:"userPrincipalName"`
	DisplayName                           string  `json:"displayName"`
	LastActivityDate                      *string `json:"lastActivityDate"`
	MicrosoftTeamsCopilotLastActivityDate *string `json:"microsoftTeamsCopilotLastActivityDate"`
	WordCopilotLastActivityDate           *string `json:"wordCopilotLastActivityDate"`
	ExcelCopilotLastActivityDate          *string `json:"excelCopilotLastActivityDate"`
	PowerPointCopilotLastActivityDate     *string `json:"powerPointCopilotLastActivityDate"`
	OutlookCopilotLastActivityDate        *string `json:"outlookCopilotLastActivityDate"`
	OneNoteCopilotLastActivityDate        *string `json:"oneNoteCopilotLastActivityDate"`
	LoopCopilotLastActivityDate           *string `json:"loopCopilotLastActivityDate"`
	CopilotChatLastActivityDate           *string `json:"copilotChatLastActivityDate"`
}

// CopilotParams defines query parameters for Copilot usage data.
type CopilotParams struct {
	Period CopilotPeriod
}

// ============================================================================
// Qualtrics Types
// ============================================================================

// ExportStatus represents the status of an export job.
type ExportStatus string

const (
	ExportStatusInProgress ExportStatus = "inProgress"
	ExportStatusComplete   ExportStatus = "complete"
	ExportStatusFailed     ExportStatus = "failed"
)

// ExportJob represents an active survey export job.
type ExportJob struct {
	ProgressID      string       `json:"progressId"`
	SurveyID        string       `json:"surveyId"`
	Status          ExportStatus `json:"status"`
	PercentComplete int          `json:"percentComplete"`
	FileID          string       `json:"fileId,omitempty"`
	StartedAt       time.Time    `json:"startedAt"`
	CompletedAt     *time.Time   `json:"completedAt,omitempty"`
	Error           string       `json:"error,omitempty"`
}

// Survey represents a configured survey for simulation.
type Survey struct {
	SurveyID      string `json:"surveyId"`
	Name          string `json:"name"`
	ResponseCount int    `json:"responseCount"`
}

// ============================================================================
// Store Interfaces
// ============================================================================

// HarveyStore defines the interface for Harvey AI usage storage.
type HarveyStore interface {
	// GetUsage retrieves Harvey usage events matching the given parameters.
	GetUsage(ctx context.Context, params HarveyParams) ([]*HarveyUsage, error)
	// StoreUsage stores Harvey usage events.
	StoreUsage(ctx context.Context, usage []*HarveyUsage) error
}

// CopilotStore defines the interface for Microsoft 365 Copilot usage storage.
type CopilotStore interface {
	// GetUsage retrieves Copilot usage data matching the given parameters.
	GetUsage(ctx context.Context, params CopilotParams) ([]*CopilotUsage, error)
	// StoreUsage stores Copilot usage data.
	StoreUsage(ctx context.Context, usage []*CopilotUsage) error
}

// QualtricsStore defines the interface for Qualtrics survey export storage.
type QualtricsStore interface {
	// GetSurveys retrieves all configured surveys.
	GetSurveys(ctx context.Context) ([]*Survey, error)
	// StoreSurveys stores survey configurations.
	StoreSurveys(ctx context.Context, surveys []*Survey) error
	// GetExportJob retrieves an export job by progress ID.
	GetExportJob(ctx context.Context, progressID string) (*ExportJob, error)
	// StoreExportJob stores or updates an export job.
	StoreExportJob(ctx context.Context, job *ExportJob) error
	// GetFile retrieves a generated file by file ID.
	GetFile(ctx context.Context, fileID string) ([]byte, error)
	// StoreFile stores a generated file.
	StoreFile(ctx context.Context, fileID string, data []byte) error
}

// ExternalDataStore provides access to all external data source stores.
type ExternalDataStore interface {
	Harvey() HarveyStore
	Copilot() CopilotStore
	Qualtrics() QualtricsStore
}

// ============================================================================
// In-Memory Implementations
// ============================================================================

// harveyMemoryStore is a thread-safe in-memory implementation of HarveyStore.
type harveyMemoryStore struct {
	mu     sync.RWMutex
	events []*HarveyUsage
}

// GetUsage retrieves Harvey usage events matching the given parameters.
func (s *harveyMemoryStore) GetUsage(ctx context.Context, params HarveyParams) ([]*HarveyUsage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*HarveyUsage
	for _, event := range s.events {
		// Filter by time range
		if !event.Time.Before(params.From) && !event.Time.After(params.To) {
			// Apply optional filters
			if params.User != "" && event.User != params.User {
				continue
			}
			if params.Task != "" && event.Task != params.Task {
				continue
			}
			// Make a copy to avoid returning internal pointers
			eventCopy := *event
			results = append(results, &eventCopy)
		}
	}

	return results, nil
}

// StoreUsage stores Harvey usage events.
func (s *harveyMemoryStore) StoreUsage(ctx context.Context, usage []*HarveyUsage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, u := range usage {
		// Make a copy to avoid storing external pointers
		usageCopy := *u
		s.events = append(s.events, &usageCopy)
	}

	return nil
}

// copilotMemoryStore is a thread-safe in-memory implementation of CopilotStore.
type copilotMemoryStore struct {
	mu    sync.RWMutex
	usage []*CopilotUsage
}

// GetUsage retrieves Copilot usage data matching the given parameters.
func (s *copilotMemoryStore) GetUsage(ctx context.Context, params CopilotParams) ([]*CopilotUsage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*CopilotUsage
	periodDays := params.Period.Days()

	for _, u := range s.usage {
		// Filter by period
		if u.ReportPeriod == periodDays {
			// Make a copy to avoid returning internal pointers
			usageCopy := *u
			results = append(results, &usageCopy)
		}
	}

	return results, nil
}

// StoreUsage stores Copilot usage data.
func (s *copilotMemoryStore) StoreUsage(ctx context.Context, usage []*CopilotUsage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, u := range usage {
		// Make a copy to avoid storing external pointers
		usageCopy := *u
		s.usage = append(s.usage, &usageCopy)
	}

	return nil
}

// qualtricsMemoryStore is a thread-safe in-memory implementation of QualtricsStore.
type qualtricsMemoryStore struct {
	mu      sync.RWMutex
	surveys []*Survey
	jobs    map[string]*ExportJob
	files   map[string][]byte
}

// GetSurveys retrieves all configured surveys.
func (s *qualtricsMemoryStore) GetSurveys(ctx context.Context) ([]*Survey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]*Survey, 0, len(s.surveys))
	for _, survey := range s.surveys {
		surveyCopy := *survey
		results = append(results, &surveyCopy)
	}

	return results, nil
}

// StoreSurveys stores survey configurations.
func (s *qualtricsMemoryStore) StoreSurveys(ctx context.Context, surveys []*Survey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, survey := range surveys {
		surveyCopy := *survey
		s.surveys = append(s.surveys, &surveyCopy)
	}

	return nil
}

// GetExportJob retrieves an export job by progress ID.
func (s *qualtricsMemoryStore) GetExportJob(ctx context.Context, progressID string) (*ExportJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[progressID]
	if !exists {
		return nil, fmt.Errorf("export job not found: %s", progressID)
	}

	// Return a copy
	jobCopy := *job
	return &jobCopy, nil
}

// StoreExportJob stores or updates an export job.
func (s *qualtricsMemoryStore) StoreExportJob(ctx context.Context, job *ExportJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Make a copy to avoid storing external pointers
	jobCopy := *job
	s.jobs[job.ProgressID] = &jobCopy

	return nil
}

// GetFile retrieves a generated file by file ID.
func (s *qualtricsMemoryStore) GetFile(ctx context.Context, fileID string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.files[fileID]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", fileID)
	}

	// Return a copy
	result := make([]byte, len(data))
	copy(result, data)
	return result, nil
}

// StoreFile stores a generated file.
func (s *qualtricsMemoryStore) StoreFile(ctx context.Context, fileID string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Make a copy to avoid storing external slices
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	s.files[fileID] = dataCopy

	return nil
}

// ExternalMemoryStore is a thread-safe in-memory implementation of ExternalDataStore.
type ExternalMemoryStore struct {
	harvey    *harveyMemoryStore
	copilot   *copilotMemoryStore
	qualtrics *qualtricsMemoryStore
}

// NewExternalMemoryStore creates a new thread-safe in-memory store for external data.
func NewExternalMemoryStore() *ExternalMemoryStore {
	return &ExternalMemoryStore{
		harvey: &harveyMemoryStore{
			events: make([]*HarveyUsage, 0, 1000),
		},
		copilot: &copilotMemoryStore{
			usage: make([]*CopilotUsage, 0, 500),
		},
		qualtrics: &qualtricsMemoryStore{
			surveys: make([]*Survey, 0, 10),
			jobs:    make(map[string]*ExportJob),
			files:   make(map[string][]byte),
		},
	}
}

// Harvey returns the Harvey store.
func (s *ExternalMemoryStore) Harvey() HarveyStore {
	return s.harvey
}

// Copilot returns the Copilot store.
func (s *ExternalMemoryStore) Copilot() CopilotStore {
	return s.copilot
}

// Qualtrics returns the Qualtrics store.
func (s *ExternalMemoryStore) Qualtrics() QualtricsStore {
	return s.qualtrics
}
