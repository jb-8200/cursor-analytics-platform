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
	writer io.Writer
}

// NewJSONExporter creates a new JSON exporter.
func NewJSONExporter(w io.Writer) *JSONExporter {
	return &JSONExporter{
		writer: w,
	}
}

// ExportDataPointsStream exports data points as NDJSON (newline-delimited JSON).
func (e *JSONExporter) ExportDataPointsStream(dataPoints []models.ResearchDataPoint) error {
	encoder := json.NewEncoder(e.writer)
	for _, dp := range dataPoints {
		if err := encoder.Encode(dp); err != nil {
			return fmt.Errorf("failed to encode data point: %w", err)
		}
	}
	return nil
}

// ExportDataPoints exports data points as a JSON array.
func (e *JSONExporter) ExportDataPoints(dataPoints []models.ResearchDataPoint) error {
	encoder := json.NewEncoder(e.writer)
	if err := encoder.Encode(dataPoints); err != nil {
		return fmt.Errorf("failed to encode data points: %w", err)
	}
	return nil
}

// ExportDataPointsFiltered exports data points filtered by time range.
func (e *JSONExporter) ExportDataPointsFiltered(dataPoints []models.ResearchDataPoint, from, to time.Time) error {
	// Filter data points by timestamp
	var filtered []models.ResearchDataPoint
	for _, dp := range dataPoints {
		if !dp.Timestamp.Before(from) && dp.Timestamp.Before(to) {
			filtered = append(filtered, dp)
		}
	}

	// Export filtered data
	return e.ExportDataPoints(filtered)
}
