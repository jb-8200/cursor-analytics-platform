package models

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"
)

// ExportJobStatus represents export job states
type ExportJobStatus string

const (
	ExportStatusInProgress ExportJobStatus = "inProgress"
	ExportStatusComplete   ExportJobStatus = "complete"
	ExportStatusFailed     ExportJobStatus = "failed"
)

// ExportJob tracks a survey export job
type ExportJob struct {
	ProgressID      string          `json:"progressId"`
	SurveyID        string          `json:"surveyId"`
	Status          ExportJobStatus `json:"status"`
	PercentComplete int             `json:"percentComplete"`
	FileID          string          `json:"fileId,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
}

// SurveyResponse represents a single survey response
type SurveyResponse struct {
	ResponseID            string    `json:"responseId"`
	RespondentEmail       string    `json:"respondentEmail"`
	OverallAISatisfaction int       `json:"overallAISatisfaction"` // 1-5 scale
	CursorSatisfaction    int       `json:"cursorSatisfaction"`    // 1-5 scale
	CopilotSatisfaction   int       `json:"copilotSatisfaction"`   // 1-5 scale
	MostUsedTool          string    `json:"mostUsedTool"`
	PositiveFeedback      string    `json:"positiveFeedback,omitempty"`
	ImprovementAreas      string    `json:"improvementAreas,omitempty"`
	RecordedAt            time.Time `json:"recordedAt"`
}

// GenerateZIPFile creates a ZIP containing survey_responses.csv
func GenerateZIPFile(responses []SurveyResponse) ([]byte, error) {
	// Create a buffer to write the ZIP to
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Create the CSV file inside the ZIP
	csvFile, err := zipWriter.Create("survey_responses.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to create CSV file in ZIP: %w", err)
	}

	// Write CSV data
	csvWriter := csv.NewWriter(csvFile)

	// Write header
	header := []string{
		"ResponseID",
		"RespondentEmail",
		"OverallAISatisfaction",
		"CursorSatisfaction",
		"CopilotSatisfaction",
		"MostUsedTool",
		"PositiveFeedback",
		"ImprovementAreas",
		"RecordedAt",
	}
	if err := csvWriter.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, resp := range responses {
		row := []string{
			resp.ResponseID,
			resp.RespondentEmail,
			strconv.Itoa(resp.OverallAISatisfaction),
			strconv.Itoa(resp.CursorSatisfaction),
			strconv.Itoa(resp.CopilotSatisfaction),
			resp.MostUsedTool,
			resp.PositiveFeedback,
			resp.ImprovementAreas,
			resp.RecordedAt.Format(time.RFC3339),
		}
		if err := csvWriter.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	// Flush the CSV writer
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	// Close the ZIP writer
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close ZIP writer: %w", err)
	}

	return buf.Bytes(), nil
}

// ExportStartResponse is the API response for starting export
type ExportStartResponse struct {
	Result struct {
		ProgressID      string `json:"progressId"`
		Status          string `json:"status"`
		PercentComplete int    `json:"percentComplete"`
	} `json:"result"`
	Meta struct {
		HTTPStatus string `json:"httpStatus"`
		RequestID  string `json:"requestId"`
	} `json:"meta"`
}

// ExportProgressResponse is the API response for checking progress
type ExportProgressResponse struct {
	Result struct {
		Status          string `json:"status"`
		PercentComplete int    `json:"percentComplete"`
		FileID          string `json:"fileId,omitempty"`
	} `json:"result"`
	Meta struct {
		HTTPStatus string `json:"httpStatus"`
		RequestID  string `json:"requestId"`
	} `json:"meta"`
}
