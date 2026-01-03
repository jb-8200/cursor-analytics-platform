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

// NewCSVExporter creates a new CSV exporter.
func NewCSVExporter(w io.Writer) *CSVExporter {
	return &CSVExporter{
		writer: csv.NewWriter(w),
	}
}

// ExportDataPoints exports data points as CSV with all columns.
func (e *CSVExporter) ExportDataPoints(dataPoints []models.ResearchDataPoint) error {
	// Write header
	header := []string{
		"commit_hash",
		"pr_number",
		"author_id",
		"author_email",
		"repo_name",
		"ai_ratio",
		"ai_lines_added",
		"ai_lines_deleted",
		"non_ai_lines_added",
		"tab_lines",
		"composer_lines",
		"pr_volume",
		"pr_scatter",
		"additions",
		"deletions",
		"files_changed",
		"greenfield_index",
		"coding_lead_time_hours",
		"pickup_time_hours",
		"review_lead_time_hours",
		"merge_lead_time_hours",
		"review_density",
		"iteration_count",
		"rework_ratio",
		"scope_creep",
		"reviewer_count",
		"review_iterations",
		"is_reverted",
		"has_hotfix_followup",
		"survival_rate_30d",
		"was_reverted",
		"required_hotfix",
		"author_seniority",
		"repo_maturity",
		"repo_age_days",
		"primary_language",
		"is_greenfield",
		"timestamp",
	}

	if err := e.writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, dp := range dataPoints {
		row := []string{
			dp.CommitHash,
			strconv.Itoa(dp.PRNumber),
			dp.AuthorID,
			dp.AuthorEmail,
			dp.RepoName,
			fmt.Sprintf("%.4f", dp.AIRatio),
			strconv.Itoa(dp.AILinesAdded),
			strconv.Itoa(dp.AILinesDeleted),
			strconv.Itoa(dp.NonAILinesAdded),
			strconv.Itoa(dp.TabLines),
			strconv.Itoa(dp.ComposerLines),
			strconv.Itoa(dp.PRVolume),
			strconv.Itoa(dp.PRScatter),
			strconv.Itoa(dp.Additions),
			strconv.Itoa(dp.Deletions),
			strconv.Itoa(dp.FilesChanged),
			fmt.Sprintf("%.4f", dp.GreenfieldIndex),
			fmt.Sprintf("%.2f", dp.CodingLeadTimeHours),
			fmt.Sprintf("%.2f", dp.PickupTimeHours),
			fmt.Sprintf("%.2f", dp.ReviewLeadTimeHours),
			fmt.Sprintf("%.2f", dp.MergeLeadTimeHours),
			fmt.Sprintf("%.6f", dp.ReviewDensity),
			strconv.Itoa(dp.IterationCount),
			fmt.Sprintf("%.4f", dp.ReworkRatio),
			fmt.Sprintf("%.4f", dp.ScopeCreep),
			strconv.Itoa(dp.ReviewerCount),
			strconv.Itoa(dp.ReviewIterations),
			strconv.FormatBool(dp.IsReverted),
			strconv.FormatBool(dp.HasHotfixFollowup),
			fmt.Sprintf("%.4f", dp.SurvivalRate30d),
			strconv.FormatBool(dp.WasReverted),
			strconv.FormatBool(dp.RequiredHotfix),
			dp.AuthorSeniority,
			dp.RepoMaturity,
			strconv.Itoa(dp.RepoAgeDays),
			dp.PrimaryLanguage,
			strconv.FormatBool(dp.IsGreenfield),
			dp.Timestamp.Format("2006-01-02T15:04:05Z"),
		}

		if err := e.writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	e.writer.Flush()
	return e.writer.Error()
}

// GetHeaders returns the CSV header row with all column names.
func (e *CSVExporter) GetHeaders() []string {
	return []string{
		"commit_hash",
		"pr_number",
		"author_id",
		"author_email",
		"repo_name",
		"ai_ratio",
		"ai_lines_added",
		"ai_lines_deleted",
		"non_ai_lines_added",
		"tab_lines",
		"composer_lines",
		"pr_volume",
		"pr_scatter",
		"additions",
		"deletions",
		"files_changed",
		"greenfield_index",
		"coding_lead_time_hours",
		"pickup_time_hours",
		"review_lead_time_hours",
		"merge_lead_time_hours",
		"review_density",
		"iteration_count",
		"rework_ratio",
		"scope_creep",
		"reviewer_count",
		"review_iterations",
		"is_reverted",
		"has_hotfix_followup",
		"survival_rate_30d",
		"was_reverted",
		"required_hotfix",
		"author_seniority",
		"repo_maturity",
		"repo_age_days",
		"primary_language",
		"is_greenfield",
		"timestamp",
	}
}

// ExportDataPointsFiltered exports data points filtered by time range.
func (e *CSVExporter) ExportDataPointsFiltered(dataPoints []models.ResearchDataPoint, from, to time.Time) error {
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
