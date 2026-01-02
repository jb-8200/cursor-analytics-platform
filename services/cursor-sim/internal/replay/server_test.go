package replay

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplayServer_DatasetHandler(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{CommitHash: "abc123", AIRatio: 0.5, Timestamp: baseTime},
		{CommitHash: "def456", AIRatio: 0.8, Timestamp: baseTime.Add(1 * time.Hour)},
	}

	index := NewCorpusIndex(dataPoints)
	server := NewReplayServer(index)
	handler := server.DatasetHandler()

	req := httptest.NewRequest("GET", "/research/dataset?from=2026-01-14&to=2026-01-16", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "data")

	params := response["params"].(map[string]interface{})
	assert.Equal(t, "replay", params["mode"])
}

func TestReplayServer_VelocityMetricsHandler(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{CommitHash: "abc123", AIRatio: 0.1, CodingLeadTimeHours: 4.0, Timestamp: baseTime},
		{CommitHash: "def456", AIRatio: 0.9, CodingLeadTimeHours: 2.0, Timestamp: baseTime.Add(1 * time.Hour)},
	}

	index := NewCorpusIndex(dataPoints)
	server := NewReplayServer(index)
	handler := server.VelocityMetricsHandler()

	req := httptest.NewRequest("GET", "/research/metrics/velocity?from=2026-01-14&to=2026-01-16", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "data")

	params := response["params"].(map[string]interface{})
	assert.Equal(t, "replay", params["mode"])
}

func TestReplayServer_ReviewCostMetricsHandler(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{CommitHash: "abc123", AIRatio: 0.5, ReviewIterations: 2, Timestamp: baseTime},
	}

	index := NewCorpusIndex(dataPoints)
	server := NewReplayServer(index)
	handler := server.ReviewCostMetricsHandler()

	req := httptest.NewRequest("GET", "/research/metrics/review-costs?from=2026-01-14&to=2026-01-16", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "data")
}

func TestReplayServer_QualityMetricsHandler(t *testing.T) {
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	dataPoints := []models.ResearchDataPoint{
		{CommitHash: "abc123", AIRatio: 0.5, WasReverted: true, Timestamp: baseTime},
	}

	index := NewCorpusIndex(dataPoints)
	server := NewReplayServer(index)
	handler := server.QualityMetricsHandler()

	req := httptest.NewRequest("GET", "/research/metrics/quality?from=2026-01-14&to=2026-01-16", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "data")
}

func TestReplayServer_InvalidDateFormat(t *testing.T) {
	index := NewCorpusIndex(nil)
	server := NewReplayServer(index)
	handler := server.DatasetHandler()

	req := httptest.NewRequest("GET", "/research/dataset?from=invalid", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
