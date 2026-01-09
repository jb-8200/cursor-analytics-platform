// Package qualtrics provides HTTP handlers for Qualtrics survey export endpoints.
// TASK-DS-14: Create Qualtrics API Handlers
package qualtrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/services"
)

// ExportHandlers provides HTTP handlers for Qualtrics export operations.
type ExportHandlers struct {
	manager *services.ExportJobManager
}

// NewExportHandlers creates a new instance of export handlers.
func NewExportHandlers(manager *services.ExportJobManager) *ExportHandlers {
	return &ExportHandlers{
		manager: manager,
	}
}

// StartExportHandler handles POST /API/v3/surveys/{surveyId}/export-responses
// Starts a new export job and returns the progressId.
func (h *ExportHandlers) StartExportHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract surveyId from path
		surveyID := extractSurveyID(r.URL.Path)
		if surveyID == "" {
			writeError(w, http.StatusBadRequest, "invalid survey ID")
			return
		}

		// Start export
		job, err := h.manager.StartExport(surveyID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to start export")
			return
		}

		// Build response matching Qualtrics API format
		resp := models.ExportStartResponse{
			Result: struct {
				ProgressID      string `json:"progressId"`
				Status          string `json:"status"`
				PercentComplete int    `json:"percentComplete"`
			}{
				ProgressID:      job.ProgressID,
				Status:          string(job.Status),
				PercentComplete: job.PercentComplete,
			},
			Meta: struct {
				HTTPStatus string `json:"httpStatus"`
				RequestID  string `json:"requestId"`
			}{
				HTTPStatus: "200 - OK",
				RequestID:  job.ProgressID, // Use progressID as requestID
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// ProgressHandler handles GET /API/v3/surveys/{surveyId}/export-responses/{progressId}
// Returns the current status and progress of an export job.
func (h *ExportHandlers) ProgressHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract progressId from path
		progressID := extractProgressID(r.URL.Path)
		if progressID == "" {
			writeError(w, http.StatusBadRequest, "invalid progress ID")
			return
		}

		// Get progress
		job, err := h.manager.GetProgress(progressID)
		if err != nil {
			writeError(w, http.StatusNotFound, "export job not found")
			return
		}

		// Build response matching Qualtrics API format
		resp := models.ExportProgressResponse{
			Result: struct {
				Status          string `json:"status"`
				PercentComplete int    `json:"percentComplete"`
				FileID          string `json:"fileId,omitempty"`
			}{
				Status:          string(job.Status),
				PercentComplete: job.PercentComplete,
				FileID:          job.FileID,
			},
			Meta: struct {
				HTTPStatus string `json:"httpStatus"`
				RequestID  string `json:"requestId"`
			}{
				HTTPStatus: "200 - OK",
				RequestID:  progressID,
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// FileDownloadHandler handles GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file
// Returns the ZIP file containing survey responses.
func (h *ExportHandlers) FileDownloadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract fileId from path
		fileID := extractFileID(r.URL.Path)
		if fileID == "" {
			writeError(w, http.StatusBadRequest, "invalid file ID")
			return
		}

		// Get file
		data, err := h.manager.GetFile(fileID)
		if err != nil {
			writeError(w, http.StatusNotFound, "file not found")
			return
		}

		// Set headers for ZIP download
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=\"survey_responses.zip\"")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// extractSurveyID extracts the survey ID from the URL path.
// Example: /API/v3/surveys/SV_abc123/export-responses -> SV_abc123
func extractSurveyID(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "surveys" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// extractProgressID extracts the progress ID from the URL path.
// Example: /API/v3/surveys/SV_abc123/export-responses/ES_xyz -> ES_xyz
func extractProgressID(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "export-responses" && i+1 < len(parts) {
			// Return the ID after export-responses (not the /file suffix)
			nextPart := parts[i+1]
			if !strings.HasSuffix(nextPart, "/file") {
				return nextPart
			}
		}
	}
	return ""
}

// extractFileID extracts the file ID from the URL path.
// Example: /API/v3/surveys/SV_abc123/export-responses/FILE_xyz/file -> FILE_xyz
func extractFileID(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "export-responses" && i+1 < len(parts) && i+2 < len(parts) {
			if parts[i+2] == "file" {
				return parts[i+1]
			}
		}
	}
	return ""
}

// writeError writes an error response in JSON format.
func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResp := struct {
		Meta struct {
			Error struct {
				ErrorMessage string `json:"errorMessage"`
			} `json:"error"`
			HTTPStatus string `json:"httpStatus"`
		} `json:"meta"`
	}{}

	errResp.Meta.Error.ErrorMessage = message
	errResp.Meta.HTTPStatus = fmt.Sprintf("%d - %s", statusCode, http.StatusText(statusCode))

	json.NewEncoder(w).Encode(errResp)
}
