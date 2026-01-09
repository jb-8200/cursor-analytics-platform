package services

import (
	"archive/zip"
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportJobManager_StartExport(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Test Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:  true,
				SurveyID: "SV_abc123",
				ResponseCount: 10,
			},
		},
	}
	gen := generator.NewSurveyGenerator(seedData)
	manager := NewExportJobManager(gen)

	job, err := manager.StartExport("SV_abc123")
	require.NoError(t, err)

	assert.NotEmpty(t, job.ProgressID)
	assert.True(t, len(job.ProgressID) > 3, "ProgressID should have ES_ prefix plus content")
	assert.Equal(t, "SV_abc123", job.SurveyID)
	assert.Equal(t, models.ExportStatusInProgress, job.Status)
	assert.Equal(t, 0, job.PercentComplete)
	assert.Empty(t, job.FileID)
}

func TestExportJobManager_ProgressAdvancement(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Test Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:  true,
				SurveyID: "SV_abc123",
				ResponseCount: 10,
			},
		},
	}
	gen := generator.NewSurveyGenerator(seedData)
	manager := NewExportJobManager(gen)

	job, _ := manager.StartExport("SV_abc123")

	// Poll until complete
	callCount := 0
	for job.Status == models.ExportStatusInProgress {
		job, _ = manager.GetProgress(job.ProgressID)
		callCount++
		assert.True(t, callCount < 10, "Should complete within 10 calls (5 calls for 20% increments)")
	}

	assert.Equal(t, models.ExportStatusComplete, job.Status)
	assert.Equal(t, 100, job.PercentComplete)
	assert.NotEmpty(t, job.FileID)
	assert.Equal(t, 5, callCount, "Should take exactly 5 calls to advance from 0 to 100% at 20% per call")
}

func TestExportJobManager_FileDownload(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Test Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:  true,
				SurveyID: "SV_abc123",
				ResponseCount: 5,
			},
		},
	}
	gen := generator.NewSurveyGenerator(seedData)
	manager := NewExportJobManager(gen)

	job, _ := manager.StartExport("SV_abc123")

	// Poll until complete
	for job.Status == models.ExportStatusInProgress {
		job, _ = manager.GetProgress(job.ProgressID)
	}

	// Download file
	data, err := manager.GetFile(job.FileID)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Verify it's a valid ZIP
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)
	assert.Len(t, reader.File, 1, "ZIP should contain exactly one file")
	assert.Equal(t, "survey_responses.csv", reader.File[0].Name)

	// Verify CSV has content
	csvFile, err := reader.File[0].Open()
	require.NoError(t, err)
	defer csvFile.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(csvFile)
	require.NoError(t, err)
	assert.True(t, buf.Len() > 0, "CSV file should have content")
	assert.Contains(t, buf.String(), "ResponseID", "CSV should have ResponseID header")
}

func TestExportJobManager_NotFound(t *testing.T) {
	seedData := &seed.SeedData{}
	gen := generator.NewSurveyGenerator(seedData)
	manager := NewExportJobManager(gen)

	_, err := manager.GetProgress("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found", "Error should indicate job not found")

	_, err = manager.GetFile("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found", "Error should indicate file not found")
}

func TestExportJobManager_ConcurrentExports(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Test Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:  true,
				SurveyID: "SV_1",
				ResponseCount: 5,
			},
		},
	}
	gen := generator.NewSurveyGenerator(seedData)
	manager := NewExportJobManager(gen)

	var wg sync.WaitGroup
	wg.Add(2)

	var job1, job2 *models.ExportJob
	var err1, err2 error

	// Start two exports concurrently
	go func() {
		defer wg.Done()
		job1, err1 = manager.StartExport("SV_1")
		for err1 == nil && job1.Status == models.ExportStatusInProgress {
			job1, err1 = manager.GetProgress(job1.ProgressID)
			time.Sleep(1 * time.Millisecond) // Small delay to allow interleaving
		}
	}()

	go func() {
		defer wg.Done()
		job2, err2 = manager.StartExport("SV_1")
		for err2 == nil && job2.Status == models.ExportStatusInProgress {
			job2, err2 = manager.GetProgress(job2.ProgressID)
			time.Sleep(1 * time.Millisecond) // Small delay to allow interleaving
		}
	}()

	wg.Wait()

	// Both jobs should complete successfully
	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.Equal(t, models.ExportStatusComplete, job1.Status)
	assert.Equal(t, models.ExportStatusComplete, job2.Status)
	assert.NotEqual(t, job1.ProgressID, job2.ProgressID, "Jobs should have different progress IDs")
	assert.NotEqual(t, job1.FileID, job2.FileID, "Jobs should have different file IDs")

	// Both files should be downloadable
	file1, err := manager.GetFile(job1.FileID)
	require.NoError(t, err)
	assert.NotEmpty(t, file1)

	file2, err := manager.GetFile(job2.FileID)
	require.NoError(t, err)
	assert.NotEmpty(t, file2)
}

func TestExportJobManager_ProgressIDFormat(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Test Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:  true,
				SurveyID: "SV_abc123",
				ResponseCount: 5,
			},
		},
	}
	gen := generator.NewSurveyGenerator(seedData)
	manager := NewExportJobManager(gen)

	job, err := manager.StartExport("SV_abc123")
	require.NoError(t, err)

	// Progress ID should follow ES_xxx format
	assert.True(t, len(job.ProgressID) >= 6, "ProgressID should be at least ES_ + 3 chars")
	assert.Equal(t, "ES_", job.ProgressID[:3], "ProgressID should start with ES_")
}

func TestExportJobManager_FileIDFormat(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Test Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:  true,
				SurveyID: "SV_abc123",
				ResponseCount: 5,
			},
		},
	}
	gen := generator.NewSurveyGenerator(seedData)
	manager := NewExportJobManager(gen)

	job, _ := manager.StartExport("SV_abc123")

	// Poll until complete
	for job.Status == models.ExportStatusInProgress {
		job, _ = manager.GetProgress(job.ProgressID)
	}

	// FileID should be assigned and have reasonable format
	assert.NotEmpty(t, job.FileID)
	assert.True(t, len(job.FileID) >= 6, "FileID should be at least 6 characters")
}

func TestExportJobManager_CreatedAtSet(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Test Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:  true,
				SurveyID: "SV_abc123",
				ResponseCount: 5,
			},
		},
	}
	gen := generator.NewSurveyGenerator(seedData)
	manager := NewExportJobManager(gen)

	before := time.Now()
	job, err := manager.StartExport("SV_abc123")
	after := time.Now()

	require.NoError(t, err)
	assert.True(t, job.CreatedAt.After(before.Add(-time.Second)) && job.CreatedAt.Before(after.Add(time.Second)),
		"CreatedAt should be set to time near job creation")
}

func TestExportJobManager_GetProgressOnCompleteJob(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "dev@company.com", Name: "Test Developer"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Qualtrics: &seed.QualtricsSeedConfig{
				Enabled:  true,
				SurveyID: "SV_abc123",
				ResponseCount: 5,
			},
		},
	}
	gen := generator.NewSurveyGenerator(seedData)
	manager := NewExportJobManager(gen)

	job, _ := manager.StartExport("SV_abc123")

	// Poll until complete
	for job.Status == models.ExportStatusInProgress {
		job, _ = manager.GetProgress(job.ProgressID)
	}

	firstFileID := job.FileID
	firstPercentComplete := job.PercentComplete

	// Call GetProgress again on already complete job
	job2, err := manager.GetProgress(job.ProgressID)
	require.NoError(t, err)
	assert.Equal(t, models.ExportStatusComplete, job2.Status)
	assert.Equal(t, firstPercentComplete, job2.PercentComplete, "Progress should not change for complete job")
	assert.Equal(t, firstFileID, job2.FileID, "FileID should not change for complete job")
}
