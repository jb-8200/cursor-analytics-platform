package storage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Harvey Store Tests
// ============================================================================

func TestHarveyStore_StoreAndGetUsage(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()

	usage := &HarveyUsage{
		EventID:           12345,
		MessageID:         "uuid-abc-123",
		Time:              time.Now(),
		User:              "attorney@firm.com",
		Task:              HarveyTaskReview,
		ClientMatter:      2024.001,
		Source:            HarveySourceFiles,
		NumberOfDocuments: 3,
		FeedbackSentiment: HarveySentimentPositive,
	}

	// Store usage
	err := store.Harvey().StoreUsage(ctx, []*HarveyUsage{usage})
	require.NoError(t, err)

	// Retrieve usage
	params := HarveyParams{
		From: time.Now().Add(-1 * time.Hour),
		To:   time.Now().Add(1 * time.Hour),
	}
	results, err := store.Harvey().GetUsage(ctx, params)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, int64(12345), results[0].EventID)
	assert.Equal(t, "attorney@firm.com", results[0].User)
	assert.Equal(t, HarveyTaskReview, results[0].Task)
}

func TestHarveyStore_GetUsage_FilterByUser(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()
	now := time.Now()

	usages := []*HarveyUsage{
		{EventID: 1, User: "alice@firm.com", Task: HarveyTaskAssist, Time: now},
		{EventID: 2, User: "bob@firm.com", Task: HarveyTaskDraft, Time: now},
		{EventID: 3, User: "alice@firm.com", Task: HarveyTaskReview, Time: now},
	}

	err := store.Harvey().StoreUsage(ctx, usages)
	require.NoError(t, err)

	params := HarveyParams{
		From: now.Add(-1 * time.Hour),
		To:   now.Add(1 * time.Hour),
		User: "alice@firm.com",
	}
	results, err := store.Harvey().GetUsage(ctx, params)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	for _, r := range results {
		assert.Equal(t, "alice@firm.com", r.User)
	}
}

func TestHarveyStore_GetUsage_FilterByTask(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()
	now := time.Now()

	usages := []*HarveyUsage{
		{EventID: 1, User: "user@firm.com", Task: HarveyTaskAssist, Time: now},
		{EventID: 2, User: "user@firm.com", Task: HarveyTaskDraft, Time: now},
		{EventID: 3, User: "user@firm.com", Task: HarveyTaskAssist, Time: now},
	}

	err := store.Harvey().StoreUsage(ctx, usages)
	require.NoError(t, err)

	params := HarveyParams{
		From: now.Add(-1 * time.Hour),
		To:   now.Add(1 * time.Hour),
		Task: HarveyTaskAssist,
	}
	results, err := store.Harvey().GetUsage(ctx, params)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	for _, r := range results {
		assert.Equal(t, HarveyTaskAssist, r.Task)
	}
}

func TestHarveyStore_GetUsage_TimeRange(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()
	now := time.Now()

	usages := []*HarveyUsage{
		{EventID: 1, User: "user@firm.com", Task: HarveyTaskAssist, Time: now.Add(-2 * time.Hour)},
		{EventID: 2, User: "user@firm.com", Task: HarveyTaskDraft, Time: now},
		{EventID: 3, User: "user@firm.com", Task: HarveyTaskReview, Time: now.Add(2 * time.Hour)},
	}

	err := store.Harvey().StoreUsage(ctx, usages)
	require.NoError(t, err)

	params := HarveyParams{
		From: now.Add(-30 * time.Minute),
		To:   now.Add(30 * time.Minute),
	}
	results, err := store.Harvey().GetUsage(ctx, params)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, int64(2), results[0].EventID)
}

func TestHarveyStore_EmptyStore(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()

	params := HarveyParams{
		From: time.Now().Add(-1 * time.Hour),
		To:   time.Now().Add(1 * time.Hour),
	}
	results, err := store.Harvey().GetUsage(ctx, params)
	require.NoError(t, err)
	assert.Len(t, results, 0)
}

// ============================================================================
// Copilot Store Tests
// ============================================================================

func TestCopilotStore_StoreAndGetUsage(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()

	teamsDate := "2026-01-08"
	usage := &CopilotUsage{
		ReportRefreshDate:                     "2026-01-09",
		ReportPeriod:                          30,
		UserPrincipalName:                     "user@company.com",
		DisplayName:                           "Jane Developer",
		MicrosoftTeamsCopilotLastActivityDate: &teamsDate,
	}

	// Store usage
	err := store.Copilot().StoreUsage(ctx, []*CopilotUsage{usage})
	require.NoError(t, err)

	// Retrieve usage
	params := CopilotParams{
		Period: CopilotPeriodD30,
	}
	results, err := store.Copilot().GetUsage(ctx, params)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "user@company.com", results[0].UserPrincipalName)
	assert.Equal(t, 30, results[0].ReportPeriod)
	assert.Equal(t, &teamsDate, results[0].MicrosoftTeamsCopilotLastActivityDate)
}

func TestCopilotStore_GetUsage_ByPeriod(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()

	usages := []*CopilotUsage{
		{UserPrincipalName: "user1@company.com", ReportPeriod: 7},
		{UserPrincipalName: "user2@company.com", ReportPeriod: 30},
		{UserPrincipalName: "user3@company.com", ReportPeriod: 30},
	}

	err := store.Copilot().StoreUsage(ctx, usages)
	require.NoError(t, err)

	// Get D30 period
	params := CopilotParams{Period: CopilotPeriodD30}
	results, err := store.Copilot().GetUsage(ctx, params)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Get D7 period
	params = CopilotParams{Period: CopilotPeriodD7}
	results, err = store.Copilot().GetUsage(ctx, params)
	require.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestCopilotStore_EmptyStore(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()

	params := CopilotParams{Period: CopilotPeriodD30}
	results, err := store.Copilot().GetUsage(ctx, params)
	require.NoError(t, err)
	assert.Len(t, results, 0)
}

// ============================================================================
// Qualtrics Store Tests
// ============================================================================

func TestQualtricsStore_GetSurveys(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()

	surveys := []*Survey{
		{SurveyID: "SV_abc123", Name: "AI Tools Survey Q1 2026", ResponseCount: 150},
		{SurveyID: "SV_def456", Name: "Developer Experience", ResponseCount: 75},
	}

	err := store.Qualtrics().StoreSurveys(ctx, surveys)
	require.NoError(t, err)

	results, err := store.Qualtrics().GetSurveys(ctx)
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestQualtricsStore_ExportJob_Lifecycle(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()

	job := &ExportJob{
		ProgressID:      "ES_xyz789",
		SurveyID:        "SV_abc123",
		Status:          ExportStatusInProgress,
		PercentComplete: 0,
		StartedAt:       time.Now(),
	}

	// Store the job
	err := store.Qualtrics().StoreExportJob(ctx, job)
	require.NoError(t, err)

	// Retrieve the job
	retrieved, err := store.Qualtrics().GetExportJob(ctx, "ES_xyz789")
	require.NoError(t, err)
	assert.Equal(t, "SV_abc123", retrieved.SurveyID)
	assert.Equal(t, ExportStatusInProgress, retrieved.Status)

	// Update job progress
	job.PercentComplete = 50
	err = store.Qualtrics().StoreExportJob(ctx, job)
	require.NoError(t, err)

	// Verify update
	retrieved, err = store.Qualtrics().GetExportJob(ctx, "ES_xyz789")
	require.NoError(t, err)
	assert.Equal(t, 50, retrieved.PercentComplete)

	// Complete the job
	job.Status = ExportStatusComplete
	job.PercentComplete = 100
	job.FileID = "FILE_abc123"
	now := time.Now()
	job.CompletedAt = &now
	err = store.Qualtrics().StoreExportJob(ctx, job)
	require.NoError(t, err)

	// Verify completion
	retrieved, err = store.Qualtrics().GetExportJob(ctx, "ES_xyz789")
	require.NoError(t, err)
	assert.Equal(t, ExportStatusComplete, retrieved.Status)
	assert.Equal(t, "FILE_abc123", retrieved.FileID)
}

func TestQualtricsStore_GetExportJob_NotFound(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()

	job, err := store.Qualtrics().GetExportJob(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, job)
}

func TestQualtricsStore_StoreAndGetFile(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()

	fileData := []byte("ZIP file content here")

	err := store.Qualtrics().StoreFile(ctx, "FILE_abc123", fileData)
	require.NoError(t, err)

	retrieved, err := store.Qualtrics().GetFile(ctx, "FILE_abc123")
	require.NoError(t, err)
	assert.Equal(t, fileData, retrieved)
}

func TestQualtricsStore_GetFile_NotFound(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()

	data, err := store.Qualtrics().GetFile(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestQualtricsStore_EmptyStore(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()

	surveys, err := store.Qualtrics().GetSurveys(ctx)
	require.NoError(t, err)
	assert.Len(t, surveys, 0)
}

// ============================================================================
// ExternalDataStore Container Tests
// ============================================================================

func TestExternalDataStore_AccessSubstores(t *testing.T) {
	store := NewExternalMemoryStore()

	// Verify all substores are accessible
	assert.NotNil(t, store.Harvey())
	assert.NotNil(t, store.Copilot())
	assert.NotNil(t, store.Qualtrics())
}

// ============================================================================
// Concurrent Access Tests
// ============================================================================

func TestExternalStore_ConcurrentAccess(t *testing.T) {
	store := NewExternalMemoryStore()
	ctx := context.Background()
	now := time.Now()

	var wg sync.WaitGroup
	errChan := make(chan error, 300)

	// Concurrent Harvey writes
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			usage := &HarveyUsage{
				EventID: int64(idx),
				User:    "user@firm.com",
				Task:    HarveyTaskAssist,
				Time:    now,
			}
			if err := store.Harvey().StoreUsage(ctx, []*HarveyUsage{usage}); err != nil {
				errChan <- err
			}
		}(i)
	}

	// Concurrent Harvey reads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			params := HarveyParams{
				From: now.Add(-1 * time.Hour),
				To:   now.Add(1 * time.Hour),
			}
			if _, err := store.Harvey().GetUsage(ctx, params); err != nil {
				errChan <- err
			}
		}()
	}

	// Concurrent Copilot writes
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			usage := &CopilotUsage{
				UserPrincipalName: "user@company.com",
				ReportPeriod:      30,
			}
			if err := store.Copilot().StoreUsage(ctx, []*CopilotUsage{usage}); err != nil {
				errChan <- err
			}
		}(i)
	}

	// Concurrent Copilot reads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			params := CopilotParams{Period: CopilotPeriodD30}
			if _, err := store.Copilot().GetUsage(ctx, params); err != nil {
				errChan <- err
			}
		}()
	}

	// Concurrent Qualtrics job writes
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			job := &ExportJob{
				ProgressID:      "ES_" + string(rune('a'+idx)),
				SurveyID:        "SV_abc",
				Status:          ExportStatusInProgress,
				PercentComplete: idx,
			}
			if err := store.Qualtrics().StoreExportJob(ctx, job); err != nil {
				errChan <- err
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		t.Errorf("concurrent access error: %v", err)
	}
}

// ============================================================================
// Type Constant Tests
// ============================================================================

func TestHarveyTask_Constants(t *testing.T) {
	assert.Equal(t, HarveyTask("Assist"), HarveyTaskAssist)
	assert.Equal(t, HarveyTask("Draft"), HarveyTaskDraft)
	assert.Equal(t, HarveyTask("Review"), HarveyTaskReview)
	assert.Equal(t, HarveyTask("Research"), HarveyTaskResearch)
}

func TestHarveySource_Constants(t *testing.T) {
	assert.Equal(t, HarveySource("Files"), HarveySourceFiles)
	assert.Equal(t, HarveySource("Web"), HarveySourceWeb)
	assert.Equal(t, HarveySource("Knowledge"), HarveySourceKnowledge)
}

func TestHarveySentiment_Constants(t *testing.T) {
	assert.Equal(t, HarveySentiment("positive"), HarveySentimentPositive)
	assert.Equal(t, HarveySentiment("negative"), HarveySentimentNegative)
	assert.Equal(t, HarveySentiment("neutral"), HarveySentimentNeutral)
}

func TestCopilotPeriod_Constants(t *testing.T) {
	assert.Equal(t, CopilotPeriod("D7"), CopilotPeriodD7)
	assert.Equal(t, CopilotPeriod("D30"), CopilotPeriodD30)
	assert.Equal(t, CopilotPeriod("D90"), CopilotPeriodD90)
	assert.Equal(t, CopilotPeriod("D180"), CopilotPeriodD180)
	assert.Equal(t, CopilotPeriod("ALL"), CopilotPeriodAll)
}

func TestExportStatus_Constants(t *testing.T) {
	assert.Equal(t, ExportStatus("inProgress"), ExportStatusInProgress)
	assert.Equal(t, ExportStatus("complete"), ExportStatusComplete)
	assert.Equal(t, ExportStatus("failed"), ExportStatusFailed)
}

func TestCopilotPeriod_Days(t *testing.T) {
	tests := []struct {
		period CopilotPeriod
		days   int
	}{
		{CopilotPeriodD7, 7},
		{CopilotPeriodD30, 30},
		{CopilotPeriodD90, 90},
		{CopilotPeriodD180, 180},
		{CopilotPeriodAll, 180},
		{CopilotPeriod("INVALID"), 30}, // Default case
	}

	for _, tt := range tests {
		t.Run(string(tt.period), func(t *testing.T) {
			assert.Equal(t, tt.days, tt.period.Days())
		})
	}
}
