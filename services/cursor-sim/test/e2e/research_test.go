package e2e

import (
	"bytes"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"context"
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/export"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/replay"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/services"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_ResearchDatasetGeneration tests the full pipeline:
// seed -> commit generation -> PR generation -> dataset export
func TestE2E_ResearchDatasetGeneration(t *testing.T) {
	// Load test seed data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)

	// Initialize storage and generate data
	store := storage.NewMemoryStore()
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err)

	// Generate commits
	commitGen := generator.NewCommitGeneratorWithSeed(seedData, store, "medium", 42)
	err = commitGen.GenerateCommits(context.Background(), 7, 0)
	require.NoError(t, err)

	// Generate PRs
	prGen := generator.NewPRGeneratorWithSeed(seedData, store, 42)
	err = prGen.GeneratePRsFromCommits(time.Now().Add(-7*24*time.Hour), time.Now())
	require.NoError(t, err)

	// Generate research dataset
	researchGen := generator.NewResearchGeneratorWithSeed(seedData, store, 42)
	dataPoints, err := researchGen.GenerateDataset(
		time.Now().Add(-8*24*time.Hour),
		time.Now().Add(24*time.Hour),
	)
	require.NoError(t, err)
	require.Greater(t, len(dataPoints), 0, "should have research data points")

	// Verify data point structure
	for _, dp := range dataPoints {
		assert.NotEmpty(t, dp.CommitHash)
		assert.NotEmpty(t, dp.AuthorID)
		assert.GreaterOrEqual(t, dp.AIRatio, 0.0)
		assert.LessOrEqual(t, dp.AIRatio, 1.0)
		assert.NotZero(t, dp.Timestamp)
	}

	t.Logf("Generated %d research data points", len(dataPoints))
}

// TestE2E_CSVExportFormat tests that CSV exports are valid and loadable
func TestE2E_CSVExportFormat(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{
			CommitHash:          "abc123",
			PRNumber:            1,
			AuthorID:            "user_001",
			RepoName:            "test/repo",
			AIRatio:             0.65,
			TabLines:            100,
			ComposerLines:       50,
			Additions:           200,
			Deletions:           50,
			FilesChanged:        5,
			CodingLeadTimeHours: 4.5,
			ReviewLeadTimeHours: 2.0,
			MergeLeadTimeHours:  1.5,
			WasReverted:         false,
			RequiredHotfix:      false,
			ReviewIterations:    2,
			AuthorSeniority:     "senior",
			RepoMaturity:        "mature",
			IsGreenfield:        false,
			Timestamp:           baseTime,
		},
	}

	// Export to CSV
	var buf bytes.Buffer
	exporter := export.NewCSVExporter(&buf)
	err := exporter.ExportDataPoints(dataPoints)
	require.NoError(t, err)

	// Parse the CSV
	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Verify structure
	require.Len(t, records, 2) // header + 1 data row

	header := records[0]
	assert.Equal(t, "commit_hash", header[0])
	assert.Equal(t, "author_email", header[3])
	assert.Equal(t, "ai_ratio", header[5])

	// Verify data row
	row := records[1]
	assert.Equal(t, "abc123", row[0])
	assert.Contains(t, row[5], "0.65") // ai_ratio

	// Verify we have all 38 columns
	assert.Len(t, header, 38, "Expected 38 columns in CSV export")

	t.Logf("CSV export validated: %d columns, %d rows", len(header), len(records)-1)
}

// TestE2E_JSONExportFormat tests that JSON exports are valid
func TestE2E_JSONExportFormat(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{
			CommitHash: "abc123",
			AIRatio:    0.65,
			Timestamp:  baseTime,
		},
		{
			CommitHash: "def456",
			AIRatio:    0.30,
			Timestamp:  baseTime.Add(1 * time.Hour),
		},
	}

	// Export to JSON
	var buf bytes.Buffer
	exporter := export.NewJSONExporter(&buf)
	err := exporter.ExportDataPoints(dataPoints)
	require.NoError(t, err)

	// Parse the JSON
	var decoded []models.ResearchDataPoint
	err = json.Unmarshal(buf.Bytes(), &decoded)
	require.NoError(t, err)

	assert.Len(t, decoded, 2)
	assert.Equal(t, "abc123", decoded[0].CommitHash)
	assert.Equal(t, 0.65, decoded[0].AIRatio)

	t.Logf("JSON export validated: %d data points", len(decoded))
}

// TestE2E_ResearchMetricsAggregation tests metrics calculation from generated data
func TestE2E_ResearchMetricsAggregation(t *testing.T) {
	// Create test data with known distribution
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		// Low AI ratio group
		{AIRatio: 0.1, CodingLeadTimeHours: 4.0, ReviewIterations: 1, WasReverted: false, Timestamp: baseTime},
		{AIRatio: 0.2, CodingLeadTimeHours: 6.0, ReviewIterations: 2, WasReverted: true, Timestamp: baseTime},
		// High AI ratio group
		{AIRatio: 0.8, CodingLeadTimeHours: 2.0, ReviewIterations: 1, WasReverted: false, Timestamp: baseTime},
		{AIRatio: 0.9, CodingLeadTimeHours: 3.0, ReviewIterations: 3, WasReverted: true, Timestamp: baseTime},
	}

	svc := services.NewResearchMetricsService(dataPoints)

	// Test velocity metrics
	velocityMetrics := svc.CalculateVelocityMetrics("2026-01")
	require.Len(t, velocityMetrics, 2) // low and high bands

	// Test quality metrics
	qualityMetrics := svc.CalculateQualityMetrics("2026-01")
	require.Len(t, qualityMetrics, 2)

	// Both bands should have 50% revert rate (1 of 2)
	for _, m := range qualityMetrics {
		assert.Equal(t, 0.5, m.RevertRate)
	}

	t.Logf("Metrics validated: %d velocity groups, %d quality groups",
		len(velocityMetrics), len(qualityMetrics))
}

// TestE2E_ReplayModeServesData tests that replay mode correctly serves corpus data
func TestE2E_ReplayModeServesData(t *testing.T) {
	// Create temporary corpus file
	tmpDir := t.TempDir()
	corpusPath := filepath.Join(tmpDir, "corpus.json")

	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{CommitHash: "replay1", AIRatio: 0.2, Timestamp: baseTime},  // low band (<0.3)
		{CommitHash: "replay2", AIRatio: 0.5, Timestamp: baseTime.Add(1 * time.Hour)},  // medium band
		{CommitHash: "replay3", AIRatio: 0.8, Timestamp: baseTime.Add(2 * time.Hour)},  // high band (>0.7)
	}

	// Write corpus file
	data, err := json.Marshal(dataPoints)
	require.NoError(t, err)
	err = os.WriteFile(corpusPath, data, 0644)
	require.NoError(t, err)

	// Load corpus
	loader := replay.NewCorpusLoader()
	loaded, err := loader.LoadJSON(corpusPath)
	require.NoError(t, err)
	require.Len(t, loaded, 3)

	// Create index and query
	index := replay.NewCorpusIndex(loaded)

	// Query by time range
	results := index.QueryByTimeRange(
		baseTime.Add(-1*time.Hour),
		baseTime.Add(90*time.Minute),
	)
	require.Len(t, results, 2) // replay1 and replay2

	// Query by AI ratio band
	lowResults := index.QueryByAIRatioBand(models.AIRatioBandLow)
	require.Len(t, lowResults, 1)
	assert.Equal(t, "replay1", lowResults[0].CommitHash)

	highResults := index.QueryByAIRatioBand(models.AIRatioBandHigh)
	require.Len(t, highResults, 1)
	assert.Equal(t, "replay3", highResults[0].CommitHash)

	t.Logf("Replay mode validated: corpus loaded and queryable")
}

// TestE2E_FullResearchPipeline tests the complete research workflow
func TestE2E_FullResearchPipeline(t *testing.T) {
	// 1. Load seed and generate data
	seedData, err := seed.LoadSeed("../../testdata/valid_seed.json")
	require.NoError(t, err)

	store := storage.NewMemoryStore()
	err = store.LoadDevelopers(seedData.Developers)
	require.NoError(t, err)

	commitGen := generator.NewCommitGeneratorWithSeed(seedData, store, "high", 42)
	err = commitGen.GenerateCommits(context.Background(), 14, 0) // 2 weeks
	require.NoError(t, err)

	prGen := generator.NewPRGeneratorWithSeed(seedData, store, 42)
	err = prGen.GeneratePRsFromCommits(time.Now().Add(-14*24*time.Hour), time.Now())
	require.NoError(t, err)

	// 2. Generate research dataset
	researchGen := generator.NewResearchGeneratorWithSeed(seedData, store, 42)
	dataPoints, err := researchGen.GenerateDataset(
		time.Now().Add(-15*24*time.Hour),
		time.Now().Add(24*time.Hour),
	)
	require.NoError(t, err)
	require.Greater(t, len(dataPoints), 0)

	// 3. Export to JSON (simulating corpus creation)
	tmpDir := t.TempDir()
	corpusPath := filepath.Join(tmpDir, "research_corpus.json")

	corpusData, err := json.Marshal(dataPoints)
	require.NoError(t, err)
	err = os.WriteFile(corpusPath, corpusData, 0644)
	require.NoError(t, err)

	// 4. Load corpus in replay mode
	loader := replay.NewCorpusLoader()
	replayData, err := loader.LoadJSON(corpusPath)
	require.NoError(t, err)
	assert.Equal(t, len(dataPoints), len(replayData))

	// 5. Calculate metrics from replay data
	svc := services.NewResearchMetricsService(replayData)
	velocityMetrics := svc.CalculateVelocityMetrics("2026-01")
	reviewMetrics := svc.CalculateReviewCostMetrics("2026-01")
	qualityMetrics := svc.CalculateQualityMetrics("2026-01")

	// 6. Verify metrics were calculated
	assert.Greater(t, len(velocityMetrics)+len(reviewMetrics)+len(qualityMetrics), 0,
		"should have at least one metric")

	t.Logf("Full pipeline validated:")
	t.Logf("  - Generated %d data points", len(dataPoints))
	t.Logf("  - Velocity metrics: %d groups", len(velocityMetrics))
	t.Logf("  - Review metrics: %d groups", len(reviewMetrics))
	t.Logf("  - Quality metrics: %d groups", len(qualityMetrics))
}
