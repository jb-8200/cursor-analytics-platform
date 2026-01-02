package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

// CSVExporter exports research data points to CSV format.
type CSVExporter struct {
	writer *csv.Writer
}

// NewCSVExporter creates a new CSV exporter that writes to the given writer.
func NewCSVExporter(w io.Writer) *CSVExporter {
	return &CSVExporter{
		writer: csv.NewWriter(w),
	}
}

// GetHeaders returns the CSV column headers.
func (e *CSVExporter) GetHeaders() []string {
	return []string{
		"commit_hash",
		"pr_number",
		"author_id",
		"repo_name",
		"ai_ratio",
		"tab_lines",
		"composer_lines",
		"additions",
		"deletions",
		"files_changed",
		"coding_lead_time_hours",
		"review_lead_time_hours",
		"merge_lead_time_hours",
		"was_reverted",
		"required_hotfix",
		"review_iterations",
		"author_seniority",
		"repo_maturity",
		"is_greenfield",
		"timestamp",
	}
}

// ExportDataPoints exports all data points to CSV.
func (e *CSVExporter) ExportDataPoints(dataPoints []models.ResearchDataPoint) error {
	// Write header
	if err := e.writer.Write(e.GetHeaders()); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, dp := range dataPoints {
		row := e.dataPointToRow(dp)
		if err := e.writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	e.writer.Flush()
	return e.writer.Error()
}

// ExportDataPointsFiltered exports data points within the given time range.
func (e *CSVExporter) ExportDataPointsFiltered(dataPoints []models.ResearchDataPoint, from, to time.Time) error {
	// Write header
	if err := e.writer.Write(e.GetHeaders()); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write filtered data rows
	for _, dp := range dataPoints {
		if (dp.Timestamp.Equal(from) || dp.Timestamp.After(from)) &&
			(dp.Timestamp.Equal(to) || dp.Timestamp.Before(to)) {
			row := e.dataPointToRow(dp)
			if err := e.writer.Write(row); err != nil {
				return fmt.Errorf("failed to write CSV row: %w", err)
			}
		}
	}

	e.writer.Flush()
	return e.writer.Error()
}

// dataPointToRow converts a ResearchDataPoint to a CSV row.
func (e *CSVExporter) dataPointToRow(dp models.ResearchDataPoint) []string {
	return []string{
		dp.CommitHash,
		strconv.Itoa(dp.PRNumber),
		dp.AuthorID,
		dp.RepoName,
		formatFloat(dp.AIRatio, 4),
		strconv.Itoa(dp.TabLines),
		strconv.Itoa(dp.ComposerLines),
		strconv.Itoa(dp.Additions),
		strconv.Itoa(dp.Deletions),
		strconv.Itoa(dp.FilesChanged),
		formatFloat(dp.CodingLeadTimeHours, 2),
		formatFloat(dp.ReviewLeadTimeHours, 2),
		formatFloat(dp.MergeLeadTimeHours, 2),
		strconv.FormatBool(dp.WasReverted),
		strconv.FormatBool(dp.RequiredHotfix),
		strconv.Itoa(dp.ReviewIterations),
		dp.AuthorSeniority,
		dp.RepoMaturity,
		strconv.FormatBool(dp.IsGreenfield),
		dp.Timestamp.Format(time.RFC3339),
	}
}

// formatFloat formats a float with specified precision.
func formatFloat(f float64, precision int) string {
	return strconv.FormatFloat(f, 'f', precision, 64)
}
