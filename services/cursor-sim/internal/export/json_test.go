package export

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONExporter_ExportDataPoints(t *testing.T) {
	ts := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{
			CommitHash:      "abc123",
			PRNumber:        42,
			AuthorID:        "user_001",
			AIRatio:         0.65,
			AuthorSeniority: "senior",
			Timestamp:       ts,
		},
		{
			CommitHash:      "def456",
			PRNumber:        43,
			AuthorID:        "user_002",
			AIRatio:         0.30,
			AuthorSeniority: "junior",
			Timestamp:       ts.Add(1 * time.Hour),
		},
	}

	var buf bytes.Buffer
	exporter := NewJSONExporter(&buf)

	err := exporter.ExportDataPoints(dataPoints)
	require.NoError(t, err)

	// Verify JSON is valid
	var decoded []models.ResearchDataPoint
	err = json.Unmarshal(buf.Bytes(), &decoded)
	require.NoError(t, err)

	assert.Len(t, decoded, 2)
	assert.Equal(t, "abc123", decoded[0].CommitHash)
	assert.Equal(t, "def456", decoded[1].CommitHash)
}

func TestJSONExporter_EmptyDataset(t *testing.T) {
	var buf bytes.Buffer
	exporter := NewJSONExporter(&buf)

	err := exporter.ExportDataPoints(nil)
	require.NoError(t, err)

	// Should be null or empty array in JSON
	content := strings.TrimSpace(buf.String())
	assert.Equal(t, "null", content)
}

func TestJSONExporter_FilterByDateRange(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{CommitHash: "early", Timestamp: baseTime.Add(-24 * time.Hour)},
		{CommitHash: "inrange1", Timestamp: baseTime},
		{CommitHash: "inrange2", Timestamp: baseTime.Add(1 * time.Hour)},
		{CommitHash: "late", Timestamp: baseTime.Add(48 * time.Hour)},
	}

	var buf bytes.Buffer
	exporter := NewJSONExporter(&buf)

	from := baseTime.Add(-1 * time.Hour)
	to := baseTime.Add(2 * time.Hour)
	err := exporter.ExportDataPointsFiltered(dataPoints, from, to)
	require.NoError(t, err)

	var decoded []models.ResearchDataPoint
	err = json.Unmarshal(buf.Bytes(), &decoded)
	require.NoError(t, err)

	assert.Len(t, decoded, 2)
	assert.Equal(t, "inrange1", decoded[0].CommitHash)
	assert.Equal(t, "inrange2", decoded[1].CommitHash)
}

func TestJSONExporter_StreamFormat(t *testing.T) {
	ts := time.Now()
	dataPoints := []models.ResearchDataPoint{
		{CommitHash: "abc123", Timestamp: ts},
		{CommitHash: "def456", Timestamp: ts},
	}

	var buf bytes.Buffer
	exporter := NewJSONExporter(&buf)

	err := exporter.ExportDataPointsStream(dataPoints)
	require.NoError(t, err)

	// Stream format should be NDJSON (one JSON object per line)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 2)

	// Each line should be valid JSON
	for _, line := range lines {
		var dp models.ResearchDataPoint
		err := json.Unmarshal([]byte(line), &dp)
		require.NoError(t, err)
	}
}
