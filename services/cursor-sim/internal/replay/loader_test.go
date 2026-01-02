package replay

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCorpusLoader_LoadJSON(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	corpusPath := filepath.Join(tmpDir, "corpus.json")

	ts := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{
			CommitHash: "abc123",
			PRNumber:   1,
			AuthorID:   "user_001",
			AIRatio:    0.5,
			Timestamp:  ts,
		},
		{
			CommitHash: "def456",
			PRNumber:   2,
			AuthorID:   "user_002",
			AIRatio:    0.8,
			Timestamp:  ts.Add(1 * time.Hour),
		},
	}

	// Write test file
	data, err := json.Marshal(dataPoints)
	require.NoError(t, err)
	err = os.WriteFile(corpusPath, data, 0644)
	require.NoError(t, err)

	// Load corpus
	loader := NewCorpusLoader()
	loaded, err := loader.LoadJSON(corpusPath)
	require.NoError(t, err)

	assert.Len(t, loaded, 2)
	assert.Equal(t, "abc123", loaded[0].CommitHash)
	assert.Equal(t, "def456", loaded[1].CommitHash)
}

func TestCorpusLoader_LoadNDJSON(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	corpusPath := filepath.Join(tmpDir, "corpus.ndjson")

	ts := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{CommitHash: "abc123", Timestamp: ts},
		{CommitHash: "def456", Timestamp: ts.Add(1 * time.Hour)},
	}

	// Write test file (one JSON per line)
	f, err := os.Create(corpusPath)
	require.NoError(t, err)
	for _, dp := range dataPoints {
		line, _ := json.Marshal(dp)
		f.Write(line)
		f.WriteString("\n")
	}
	f.Close()

	// Load corpus
	loader := NewCorpusLoader()
	loaded, err := loader.LoadNDJSON(corpusPath)
	require.NoError(t, err)

	assert.Len(t, loaded, 2)
}

func TestCorpusLoader_FileNotFound(t *testing.T) {
	loader := NewCorpusLoader()
	_, err := loader.LoadJSON("/nonexistent/path.json")
	assert.Error(t, err)
}

func TestCorpusLoader_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	corpusPath := filepath.Join(tmpDir, "invalid.json")

	err := os.WriteFile(corpusPath, []byte("not valid json"), 0644)
	require.NoError(t, err)

	loader := NewCorpusLoader()
	_, err = loader.LoadJSON(corpusPath)
	assert.Error(t, err)
}

func TestCorpusIndex_QueryByTimeRange(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{CommitHash: "early", Timestamp: baseTime.Add(-24 * time.Hour)},
		{CommitHash: "inrange1", Timestamp: baseTime},
		{CommitHash: "inrange2", Timestamp: baseTime.Add(1 * time.Hour)},
		{CommitHash: "late", Timestamp: baseTime.Add(48 * time.Hour)},
	}

	index := NewCorpusIndex(dataPoints)

	from := baseTime.Add(-1 * time.Hour)
	to := baseTime.Add(2 * time.Hour)

	results := index.QueryByTimeRange(from, to)

	assert.Len(t, results, 2)
	assert.Equal(t, "inrange1", results[0].CommitHash)
	assert.Equal(t, "inrange2", results[1].CommitHash)
}

func TestCorpusIndex_QueryByAIRatioBand(t *testing.T) {
	dataPoints := []models.ResearchDataPoint{
		{CommitHash: "low", AIRatio: 0.1},
		{CommitHash: "medium", AIRatio: 0.5},
		{CommitHash: "high", AIRatio: 0.9},
	}

	index := NewCorpusIndex(dataPoints)

	lowResults := index.QueryByAIRatioBand(models.AIRatioBandLow)
	assert.Len(t, lowResults, 1)
	assert.Equal(t, "low", lowResults[0].CommitHash)

	highResults := index.QueryByAIRatioBand(models.AIRatioBandHigh)
	assert.Len(t, highResults, 1)
	assert.Equal(t, "high", highResults[0].CommitHash)
}

func TestCorpusIndex_GetAll(t *testing.T) {
	dataPoints := []models.ResearchDataPoint{
		{CommitHash: "a"},
		{CommitHash: "b"},
		{CommitHash: "c"},
	}

	index := NewCorpusIndex(dataPoints)

	all := index.GetAll()
	assert.Len(t, all, 3)
}

func TestCorpusIndex_Empty(t *testing.T) {
	index := NewCorpusIndex(nil)

	results := index.QueryByTimeRange(time.Now(), time.Now().Add(1*time.Hour))
	assert.Empty(t, results)

	all := index.GetAll()
	assert.Empty(t, all)
}
