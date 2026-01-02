package export

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

// JSONExporter exports research data points to JSON format.
type JSONExporter struct {
	writer  io.Writer
	encoder *json.Encoder
}

// NewJSONExporter creates a new JSON exporter that writes to the given writer.
func NewJSONExporter(w io.Writer) *JSONExporter {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return &JSONExporter{
		writer:  w,
		encoder: encoder,
	}
}

// ExportDataPoints exports all data points to JSON.
func (e *JSONExporter) ExportDataPoints(dataPoints []models.ResearchDataPoint) error {
	if err := e.encoder.Encode(dataPoints); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

// ExportDataPointsFiltered exports data points within the given time range.
func (e *JSONExporter) ExportDataPointsFiltered(dataPoints []models.ResearchDataPoint, from, to time.Time) error {
	var filtered []models.ResearchDataPoint
	for _, dp := range dataPoints {
		if (dp.Timestamp.Equal(from) || dp.Timestamp.After(from)) &&
			(dp.Timestamp.Equal(to) || dp.Timestamp.Before(to)) {
			filtered = append(filtered, dp)
		}
	}
	return e.ExportDataPoints(filtered)
}

// ExportDataPointsStream exports data points one at a time for large datasets.
// This uses JSON Lines format (NDJSON) - one JSON object per line.
func (e *JSONExporter) ExportDataPointsStream(dataPoints []models.ResearchDataPoint) error {
	// Create a non-indented encoder for streaming
	encoder := json.NewEncoder(e.writer)
	for _, dp := range dataPoints {
		if err := encoder.Encode(dp); err != nil {
			return fmt.Errorf("failed to encode JSON line: %w", err)
		}
	}
	return nil
}
