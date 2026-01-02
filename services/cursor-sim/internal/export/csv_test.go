package export

import (
	"bytes"
	"encoding/csv"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSVExporter_ExportDataPoints(t *testing.T) {
	ts := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{
			CommitHash:          "abc123",
			PRNumber:            42,
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
			Timestamp:           ts,
		},
		{
			CommitHash:          "def456",
			PRNumber:            43,
			AuthorID:            "user_002",
			RepoName:            "test/repo",
			AIRatio:             0.30,
			TabLines:            50,
			ComposerLines:       25,
			Additions:           150,
			Deletions:           75,
			FilesChanged:        3,
			CodingLeadTimeHours: 8.0,
			ReviewLeadTimeHours: 4.0,
			MergeLeadTimeHours:  2.0,
			WasReverted:         true,
			RequiredHotfix:      true,
			ReviewIterations:    3,
			AuthorSeniority:     "junior",
			RepoMaturity:        "greenfield",
			IsGreenfield:        true,
			Timestamp:           ts.Add(1 * time.Hour),
		},
	}

	var buf bytes.Buffer
	exporter := NewCSVExporter(&buf)

	err := exporter.ExportDataPoints(dataPoints)
	require.NoError(t, err)

	// Parse the CSV output
	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Should have header + 2 data rows
	assert.Len(t, records, 3)

	// Verify header
	header := records[0]
	assert.Contains(t, header, "commit_hash")
	assert.Contains(t, header, "ai_ratio")
	assert.Contains(t, header, "author_seniority")

	// Verify first data row
	row1 := records[1]
	assert.Equal(t, "abc123", row1[0]) // commit_hash is first column
}

func TestCSVExporter_Headers(t *testing.T) {
	var buf bytes.Buffer
	exporter := NewCSVExporter(&buf)

	headers := exporter.GetHeaders()

	// Verify all expected headers are present
	expectedHeaders := []string{
		"commit_hash", "pr_number", "author_id", "repo_name",
		"ai_ratio", "tab_lines", "composer_lines",
		"additions", "deletions", "files_changed",
		"coding_lead_time_hours", "review_lead_time_hours", "merge_lead_time_hours",
		"was_reverted", "required_hotfix", "review_iterations",
		"author_seniority", "repo_maturity", "is_greenfield",
		"timestamp",
	}

	for _, expected := range expectedHeaders {
		assert.Contains(t, headers, expected, "missing header: %s", expected)
	}
}

func TestCSVExporter_EmptyDataset(t *testing.T) {
	var buf bytes.Buffer
	exporter := NewCSVExporter(&buf)

	err := exporter.ExportDataPoints(nil)
	require.NoError(t, err)

	// Should still have header
	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
	records, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, records, 1) // Just header
}

func TestCSVExporter_FilterByDateRange(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{CommitHash: "early", Timestamp: baseTime.Add(-24 * time.Hour)},
		{CommitHash: "inrange1", Timestamp: baseTime},
		{CommitHash: "inrange2", Timestamp: baseTime.Add(1 * time.Hour)},
		{CommitHash: "late", Timestamp: baseTime.Add(48 * time.Hour)},
	}

	var buf bytes.Buffer
	exporter := NewCSVExporter(&buf)

	from := baseTime.Add(-1 * time.Hour)
	to := baseTime.Add(2 * time.Hour)
	err := exporter.ExportDataPointsFiltered(dataPoints, from, to)
	require.NoError(t, err)

	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Header + 2 rows in range
	assert.Len(t, records, 3)
	assert.Equal(t, "inrange1", records[1][0])
	assert.Equal(t, "inrange2", records[2][0])
}

func TestCSVExporter_BooleanFormatting(t *testing.T) {
	ts := time.Now()
	dataPoints := []models.ResearchDataPoint{
		{
			CommitHash:     "test1",
			WasReverted:    true,
			RequiredHotfix: false,
			IsGreenfield:   true,
			Timestamp:      ts,
		},
	}

	var buf bytes.Buffer
	exporter := NewCSVExporter(&buf)
	err := exporter.ExportDataPoints(dataPoints)
	require.NoError(t, err)

	content := buf.String()
	// Booleans should be "true"/"false" not "1"/"0"
	assert.Contains(t, content, "true")
	assert.Contains(t, content, "false")
}

func TestCSVExporter_FloatPrecision(t *testing.T) {
	ts := time.Now()
	dataPoints := []models.ResearchDataPoint{
		{
			CommitHash:          "test1",
			AIRatio:             0.12345678,
			CodingLeadTimeHours: 4.567,
			Timestamp:           ts,
		},
	}

	var buf bytes.Buffer
	exporter := NewCSVExporter(&buf)
	err := exporter.ExportDataPoints(dataPoints)
	require.NoError(t, err)

	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Float should have reasonable precision (not too many decimal places)
	row := records[1]
	// Find the ai_ratio column
	headers := records[0]
	aiRatioIdx := -1
	for i, h := range headers {
		if h == "ai_ratio" {
			aiRatioIdx = i
			break
		}
	}
	require.NotEqual(t, -1, aiRatioIdx)

	// AI ratio should be formatted with reasonable precision
	assert.NotContains(t, row[aiRatioIdx], "0.12345678") // Too many decimals
}
