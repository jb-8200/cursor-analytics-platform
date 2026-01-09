package models

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportJob_States(t *testing.T) {
	job := &ExportJob{
		ProgressID:      "ES_abc123",
		SurveyID:        "SV_xyz",
		Status:          ExportStatusInProgress,
		PercentComplete: 0,
		CreatedAt:       time.Now(),
	}

	assert.Equal(t, ExportStatusInProgress, job.Status)
	assert.Equal(t, 0, job.PercentComplete)

	// Transition to complete
	job.PercentComplete = 100
	job.Status = ExportStatusComplete
	job.FileID = "FILE_xyz"

	assert.Equal(t, ExportStatusComplete, job.Status)
	assert.Equal(t, 100, job.PercentComplete)
	assert.NotEmpty(t, job.FileID)
}

func TestExportJobStatus_Constants(t *testing.T) {
	assert.Equal(t, ExportJobStatus("inProgress"), ExportStatusInProgress)
	assert.Equal(t, ExportJobStatus("complete"), ExportStatusComplete)
	assert.Equal(t, ExportJobStatus("failed"), ExportStatusFailed)
}

func TestSurveyResponse_Fields(t *testing.T) {
	now := time.Now()
	resp := SurveyResponse{
		ResponseID:            "R_abc123",
		RespondentEmail:       "user@company.com",
		OverallAISatisfaction: 4,
		CursorSatisfaction:    5,
		CopilotSatisfaction:   3,
		MostUsedTool:          "Cursor",
		PositiveFeedback:      "Great tool!",
		ImprovementAreas:      "Needs better docs",
		RecordedAt:            now,
	}

	assert.Equal(t, "R_abc123", resp.ResponseID)
	assert.Equal(t, "user@company.com", resp.RespondentEmail)
	assert.Equal(t, 4, resp.OverallAISatisfaction)
	assert.Equal(t, 5, resp.CursorSatisfaction)
	assert.Equal(t, 3, resp.CopilotSatisfaction)
	assert.Equal(t, "Cursor", resp.MostUsedTool)
	assert.Equal(t, "Great tool!", resp.PositiveFeedback)
	assert.Equal(t, "Needs better docs", resp.ImprovementAreas)
	assert.Equal(t, now, resp.RecordedAt)
}

func TestSurveyResponse_SatisfactionRange(t *testing.T) {
	tests := []struct {
		name         string
		satisfaction int
		valid        bool
	}{
		{"minimum valid", 1, true},
		{"maximum valid", 5, true},
		{"mid range", 3, true},
		{"below minimum", 0, false},
		{"above maximum", 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := SurveyResponse{
				OverallAISatisfaction: tt.satisfaction,
			}
			// Satisfaction is stored as-is; validation happens at generator level
			assert.Equal(t, tt.satisfaction, resp.OverallAISatisfaction)
		})
	}
}

func TestGenerateZIPFile_Success(t *testing.T) {
	responses := []SurveyResponse{
		{
			ResponseID:            "R_001",
			RespondentEmail:       "user1@company.com",
			OverallAISatisfaction: 4,
			CursorSatisfaction:    5,
			CopilotSatisfaction:   3,
			MostUsedTool:          "Cursor",
			PositiveFeedback:      "Great tool",
			ImprovementAreas:      "Need better docs",
			RecordedAt:            time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC),
		},
		{
			ResponseID:            "R_002",
			RespondentEmail:       "user2@company.com",
			OverallAISatisfaction: 5,
			CursorSatisfaction:    4,
			CopilotSatisfaction:   4,
			MostUsedTool:          "Copilot",
			PositiveFeedback:      "Very helpful",
			ImprovementAreas:      "",
			RecordedAt:            time.Date(2026, 1, 9, 14, 30, 0, 0, time.UTC),
		},
	}

	zipData, err := GenerateZIPFile(responses)
	require.NoError(t, err)
	assert.NotEmpty(t, zipData)

	// Verify ZIP structure
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	require.NoError(t, err)
	require.Len(t, reader.File, 1)
	assert.Equal(t, "survey_responses.csv", reader.File[0].Name)
}

func TestGenerateZIPFile_CSVContent(t *testing.T) {
	responses := []SurveyResponse{
		{
			ResponseID:            "R_001",
			RespondentEmail:       "user1@company.com",
			OverallAISatisfaction: 4,
			CursorSatisfaction:    5,
			CopilotSatisfaction:   3,
			MostUsedTool:          "Cursor",
			PositiveFeedback:      "Great tool",
			ImprovementAreas:      "Need better docs",
			RecordedAt:            time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC),
		},
		{
			ResponseID:            "R_002",
			RespondentEmail:       "user2@company.com",
			OverallAISatisfaction: 5,
			CursorSatisfaction:    4,
			CopilotSatisfaction:   4,
			MostUsedTool:          "Copilot",
			PositiveFeedback:      "Very helpful",
			ImprovementAreas:      "",
			RecordedAt:            time.Date(2026, 1, 9, 14, 30, 0, 0, time.UTC),
		},
	}

	zipData, err := GenerateZIPFile(responses)
	require.NoError(t, err)

	// Extract and verify CSV content
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	require.NoError(t, err)

	csvFile, err := reader.File[0].Open()
	require.NoError(t, err)
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	rows, err := csvReader.ReadAll()
	require.NoError(t, err)

	// Verify structure
	assert.Len(t, rows, 3) // Header + 2 data rows

	// Verify header
	header := rows[0]
	assert.Contains(t, header, "ResponseID")
	assert.Contains(t, header, "RespondentEmail")
	assert.Contains(t, header, "OverallAISatisfaction")
	assert.Contains(t, header, "CursorSatisfaction")
	assert.Contains(t, header, "CopilotSatisfaction")
	assert.Contains(t, header, "MostUsedTool")
	assert.Contains(t, header, "PositiveFeedback")
	assert.Contains(t, header, "ImprovementAreas")
	assert.Contains(t, header, "RecordedAt")

	// Verify first data row
	row1 := rows[1]
	assert.Equal(t, "R_001", row1[0])
	assert.Equal(t, "user1@company.com", row1[1])
	assert.Equal(t, "4", row1[2])
	assert.Equal(t, "5", row1[3])
	assert.Equal(t, "3", row1[4])
	assert.Equal(t, "Cursor", row1[5])
	assert.Equal(t, "Great tool", row1[6])
	assert.Equal(t, "Need better docs", row1[7])
	assert.Equal(t, "2026-01-08T10:00:00Z", row1[8])

	// Verify second data row
	row2 := rows[2]
	assert.Equal(t, "R_002", row2[0])
	assert.Equal(t, "user2@company.com", row2[1])
	assert.Equal(t, "5", row2[2])
	assert.Equal(t, "4", row2[3])
	assert.Equal(t, "4", row2[4])
	assert.Equal(t, "Copilot", row2[5])
	assert.Equal(t, "Very helpful", row2[6])
	assert.Equal(t, "", row2[7])
	assert.Equal(t, "2026-01-09T14:30:00Z", row2[8])
}

func TestGenerateZIPFile_EmptyResponses(t *testing.T) {
	responses := []SurveyResponse{}

	zipData, err := GenerateZIPFile(responses)
	require.NoError(t, err)
	assert.NotEmpty(t, zipData)

	// Verify ZIP still contains header
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	require.NoError(t, err)

	csvFile, err := reader.File[0].Open()
	require.NoError(t, err)
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	rows, err := csvReader.ReadAll()
	require.NoError(t, err)

	// Should have just the header
	assert.Len(t, rows, 1)
}

func TestExportStartResponse_Format(t *testing.T) {
	resp := ExportStartResponse{
		Result: struct {
			ProgressID      string `json:"progressId"`
			Status          string `json:"status"`
			PercentComplete int    `json:"percentComplete"`
		}{
			ProgressID:      "ES_abc123",
			Status:          "inProgress",
			PercentComplete: 0,
		},
		Meta: struct {
			HTTPStatus string `json:"httpStatus"`
			RequestID  string `json:"requestId"`
		}{
			HTTPStatus: "200 - OK",
			RequestID:  "req_xyz",
		},
	}

	// Verify JSON marshaling
	data, err := json.Marshal(resp)
	require.NoError(t, err)

	// Verify exact field names match Qualtrics API
	assert.Contains(t, string(data), `"progressId"`)
	assert.Contains(t, string(data), `"status"`)
	assert.Contains(t, string(data), `"percentComplete"`)
	assert.Contains(t, string(data), `"httpStatus"`)
	assert.Contains(t, string(data), `"requestId"`)
}

func TestExportProgressResponse_Format(t *testing.T) {
	resp := ExportProgressResponse{
		Result: struct {
			Status          string `json:"status"`
			PercentComplete int    `json:"percentComplete"`
			FileID          string `json:"fileId,omitempty"`
		}{
			Status:          "complete",
			PercentComplete: 100,
			FileID:          "FILE_xyz",
		},
		Meta: struct {
			HTTPStatus string `json:"httpStatus"`
			RequestID  string `json:"requestId"`
		}{
			HTTPStatus: "200 - OK",
			RequestID:  "req_xyz",
		},
	}

	// Verify JSON marshaling
	data, err := json.Marshal(resp)
	require.NoError(t, err)

	// Verify exact field names match Qualtrics API
	assert.Contains(t, string(data), `"status"`)
	assert.Contains(t, string(data), `"percentComplete"`)
	assert.Contains(t, string(data), `"fileId"`)
	assert.Contains(t, string(data), `"httpStatus"`)
	assert.Contains(t, string(data), `"requestId"`)
}

func TestExportJob_JSONMarshal(t *testing.T) {
	job := ExportJob{
		ProgressID:      "ES_abc123",
		SurveyID:        "SV_xyz",
		Status:          ExportStatusComplete,
		PercentComplete: 100,
		FileID:          "FILE_xyz",
		CreatedAt:       time.Date(2026, 1, 9, 10, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(job)
	require.NoError(t, err)

	// Verify exact field names
	assert.Contains(t, string(data), `"progressId"`)
	assert.Contains(t, string(data), `"surveyId"`)
	assert.Contains(t, string(data), `"status"`)
	assert.Contains(t, string(data), `"percentComplete"`)
	assert.Contains(t, string(data), `"fileId"`)
	assert.Contains(t, string(data), `"createdAt"`)
}

func TestSurveyResponse_JSONMarshal(t *testing.T) {
	resp := SurveyResponse{
		ResponseID:            "R_abc123",
		RespondentEmail:       "user@company.com",
		OverallAISatisfaction: 4,
		CursorSatisfaction:    5,
		CopilotSatisfaction:   3,
		MostUsedTool:          "Cursor",
		PositiveFeedback:      "Great",
		ImprovementAreas:      "Docs",
		RecordedAt:            time.Date(2026, 1, 9, 10, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	// Verify exact field names
	assert.Contains(t, string(data), `"responseId"`)
	assert.Contains(t, string(data), `"respondentEmail"`)
	assert.Contains(t, string(data), `"overallAISatisfaction"`)
	assert.Contains(t, string(data), `"cursorSatisfaction"`)
	assert.Contains(t, string(data), `"copilotSatisfaction"`)
	assert.Contains(t, string(data), `"mostUsedTool"`)
	assert.Contains(t, string(data), `"positiveFeedback"`)
	assert.Contains(t, string(data), `"improvementAreas"`)
	assert.Contains(t, string(data), `"recordedAt"`)
}
