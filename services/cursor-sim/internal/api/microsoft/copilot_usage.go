// Package microsoft provides HTTP handlers for Microsoft Graph API endpoints.
// TASK-DS-09: Create Copilot API Handler
package microsoft

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// periodPattern matches Graph API period format: (period='D30')
var periodPattern = regexp.MustCompile(`period='([^']+)'`)

// CopilotUsageHandler returns an HTTP handler for Microsoft 365 Copilot usage endpoint.
// GET /reports/getMicrosoft365CopilotUsageUserDetail(period='D30')
// Query params: $format (application/json or text/csv)
// Auth: Basic Authentication required
func CopilotUsageHandler(store storage.CopilotStore, gen *generator.CopilotGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify authentication
		_, _, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "authentication required"})
			return
		}

		// Extract period from path (e.g., period='D30')
		period, err := extractPeriod(r.URL.Path)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		// Convert to storage period type
		storagePeriod := toStoragePeriod(period)

		// Query the store
		params := storage.CopilotParams{
			Period: storagePeriod,
		}

		usageData, err := store.GetUsage(r.Context(), params)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to query usage data"})
			return
		}

		// If store is empty, use generator
		if len(usageData) == 0 {
			modelPeriod := toModelPeriod(period)
			generatedData := gen.GenerateUsageReport(modelPeriod)

			// Convert to storage format and store for future requests
			usageData = make([]*storage.CopilotUsage, len(generatedData))
			for i, detail := range generatedData {
				usageData[i] = &storage.CopilotUsage{
					ReportRefreshDate:                     detail.ReportRefreshDate,
					ReportPeriod:                          detail.ReportPeriod,
					UserPrincipalName:                     detail.UserPrincipalName,
					DisplayName:                           detail.DisplayName,
					LastActivityDate:                      detail.LastActivityDate,
					MicrosoftTeamsCopilotLastActivityDate: detail.MicrosoftTeamsCopilotLastActivityDate,
					WordCopilotLastActivityDate:           detail.WordCopilotLastActivityDate,
					ExcelCopilotLastActivityDate:          detail.ExcelCopilotLastActivityDate,
					PowerPointCopilotLastActivityDate:     detail.PowerPointCopilotLastActivityDate,
					OutlookCopilotLastActivityDate:        detail.OutlookCopilotLastActivityDate,
					OneNoteCopilotLastActivityDate:        detail.OneNoteCopilotLastActivityDate,
					LoopCopilotLastActivityDate:           detail.LoopCopilotLastActivityDate,
					CopilotChatLastActivityDate:           detail.CopilotChatLastActivityDate,
				}
			}

			// Store generated data
			if len(usageData) > 0 {
				_ = store.StoreUsage(r.Context(), usageData)
			}
		}

		// Determine response format
		format := r.URL.Query().Get("$format")
		if format == "" {
			format = "application/json" // Default to JSON
		}

		// Convert to response format
		responseData := make([]models.CopilotUsageUserDetail, len(usageData))
		for i, u := range usageData {
			responseData[i] = models.CopilotUsageUserDetail{
				ReportRefreshDate:                     u.ReportRefreshDate,
				ReportPeriod:                          u.ReportPeriod,
				UserPrincipalName:                     u.UserPrincipalName,
				DisplayName:                           u.DisplayName,
				LastActivityDate:                      u.LastActivityDate,
				MicrosoftTeamsCopilotLastActivityDate: u.MicrosoftTeamsCopilotLastActivityDate,
				WordCopilotLastActivityDate:           u.WordCopilotLastActivityDate,
				ExcelCopilotLastActivityDate:          u.ExcelCopilotLastActivityDate,
				PowerPointCopilotLastActivityDate:     u.PowerPointCopilotLastActivityDate,
				OutlookCopilotLastActivityDate:        u.OutlookCopilotLastActivityDate,
				OneNoteCopilotLastActivityDate:        u.OneNoteCopilotLastActivityDate,
				LoopCopilotLastActivityDate:           u.LoopCopilotLastActivityDate,
				CopilotChatLastActivityDate:           u.CopilotChatLastActivityDate,
			}
		}

		// Handle CSV export
		if format == "text/csv" {
			writeCopilotCSV(w, responseData, period)
			return
		}

		// Default: JSON response with OData structure
		writeJSONResponse(w, responseData)
	}
}

// extractPeriod extracts the period parameter from the Graph API path.
// Example: /reports/getMicrosoft365CopilotUsageUserDetail(period='D30') -> D30
func extractPeriod(path string) (models.CopilotReportPeriod, error) {
	matches := periodPattern.FindStringSubmatch(path)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid or missing period parameter")
	}

	period := models.CopilotReportPeriod(matches[1])

	// Validate period
	switch period {
	case models.CopilotPeriodD7, models.CopilotPeriodD30, models.CopilotPeriodD90, models.CopilotPeriodD180:
		return period, nil
	default:
		return "", fmt.Errorf("invalid period: %s (must be D7, D30, D90, or D180)", period)
	}
}

// toStoragePeriod converts model period to storage period.
func toStoragePeriod(period models.CopilotReportPeriod) storage.CopilotPeriod {
	switch period {
	case models.CopilotPeriodD7:
		return storage.CopilotPeriodD7
	case models.CopilotPeriodD30:
		return storage.CopilotPeriodD30
	case models.CopilotPeriodD90:
		return storage.CopilotPeriodD90
	case models.CopilotPeriodD180:
		return storage.CopilotPeriodD180
	default:
		return storage.CopilotPeriodD30
	}
}

// toModelPeriod converts model period constant to generator-compatible period.
func toModelPeriod(period models.CopilotReportPeriod) models.CopilotReportPeriod {
	return period
}

// writeJSONResponse writes the OData-formatted JSON response.
func writeJSONResponse(w http.ResponseWriter, data []models.CopilotUsageUserDetail) {
	response := models.CopilotUsageResponse{
		Context: "https://graph.microsoft.com/beta/$metadata#reports/getMicrosoft365CopilotUsageUserDetail(period='D30')",
		Value:   data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// writeCopilotCSV writes the CSV export response.
func writeCopilotCSV(w http.ResponseWriter, data []models.CopilotUsageUserDetail, period models.CopilotReportPeriod) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=copilot-usage-%s.csv", period))
	w.WriteHeader(http.StatusOK)

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	header := []string{
		"Report Refresh Date",
		"Report Period",
		"User Principal Name",
		"Display Name",
		"Last Activity Date",
		"Microsoft Teams Copilot Last Activity Date",
		"Word Copilot Last Activity Date",
		"Excel Copilot Last Activity Date",
		"PowerPoint Copilot Last Activity Date",
		"Outlook Copilot Last Activity Date",
		"OneNote Copilot Last Activity Date",
		"Loop Copilot Last Activity Date",
		"Copilot Chat Last Activity Date",
	}
	writer.Write(header)

	// Write data rows
	for _, u := range data {
		row := []string{
			u.ReportRefreshDate,
			fmt.Sprintf("%d", u.ReportPeriod),
			u.UserPrincipalName,
			u.DisplayName,
			ptrStringOrEmpty(u.LastActivityDate),
			ptrStringOrEmpty(u.MicrosoftTeamsCopilotLastActivityDate),
			ptrStringOrEmpty(u.WordCopilotLastActivityDate),
			ptrStringOrEmpty(u.ExcelCopilotLastActivityDate),
			ptrStringOrEmpty(u.PowerPointCopilotLastActivityDate),
			ptrStringOrEmpty(u.OutlookCopilotLastActivityDate),
			ptrStringOrEmpty(u.OneNoteCopilotLastActivityDate),
			ptrStringOrEmpty(u.LoopCopilotLastActivityDate),
			ptrStringOrEmpty(u.CopilotChatLastActivityDate),
		}
		writer.Write(row)
	}
}

// ptrStringOrEmpty returns the value of a string pointer or empty string if nil.
func ptrStringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return strings.TrimSpace(*s)
}
