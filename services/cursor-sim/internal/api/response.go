package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

// RespondJSON writes a JSON response with the given status code and data.
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// RespondError writes a JSON error response.
func RespondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// BuildPaginatedResponse creates a paginated response wrapper.
// Deprecated: Use BuildAnalyticsTeamResponse for team-level analytics endpoints.
func BuildPaginatedResponse(data interface{}, params models.Params, totalCount int) models.PaginatedResponse {
	totalPages := 0
	if totalCount > 0 {
		totalPages = (totalCount + params.PageSize - 1) / params.PageSize
	}

	return models.PaginatedResponse{
		Data: data,
		Pagination: models.Pagination{
			Page:            params.Page,
			PageSize:        params.PageSize,
			TotalPages:      totalPages,
			HasNextPage:     params.Page < totalPages,
			HasPreviousPage: params.Page > 1,
		},
		Params: params,
	}
}

// BuildAnalyticsTeamResponse creates an analytics team-level response.
// This matches the Cursor Analytics API format for team endpoints.
//
// Reference: docs/api-reference/cursor_analytics.md (Team-Level Endpoints)
// Format: { "data": [...], "params": {...} }
func BuildAnalyticsTeamResponse(data interface{}, metric string, params models.Params) models.AnalyticsTeamResponse {
	return models.AnalyticsTeamResponse{
		Data: data,
		Params: models.AnalyticsParams{
			Metric:    metric,
			TeamID:    12345, // Fixed team ID for simulator
			StartDate: params.StartDate,
			EndDate:   params.EndDate,
			Users:     params.User,
			Page:      params.Page,
			PageSize:  params.PageSize,
		},
	}
}

// BuildAnalyticsByUserResponse creates an analytics by-user response.
// This matches the Cursor Analytics API format for by-user endpoints.
//
// Reference: docs/api-reference/cursor_analytics.md (By-User Endpoints)
// Format: { "data": { "email": [...] }, "pagination": {...}, "params": {...} }
func BuildAnalyticsByUserResponse(data map[string]interface{}, metric string, params models.Params, totalUsers int, userMappings []models.UserMapping) models.AnalyticsByUserResponse {
	totalPages := 0
	if totalUsers > 0 {
		totalPages = (totalUsers + params.PageSize - 1) / params.PageSize
	}

	return models.AnalyticsByUserResponse{
		Data: data,
		Pagination: models.Pagination{
			Page:            params.Page,
			PageSize:        params.PageSize,
			TotalUsers:      totalUsers,
			TotalPages:      totalPages,
			HasNextPage:     params.Page < totalPages,
			HasPreviousPage: params.Page > 1,
		},
		Params: models.AnalyticsParams{
			Metric:       metric,
			TeamID:       12345, // Fixed team ID for simulator
			StartDate:    params.StartDate,
			EndDate:      params.EndDate,
			Users:        params.User,
			Page:         params.Page,
			PageSize:     params.PageSize,
			UserMappings: userMappings,
		},
	}
}

// RespondCSV writes a CSV response for the given data.
// The data must be a slice of structs with json tags.
func RespondCSV(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf(
		"attachment; filename=cursor-sim-export-%s.csv",
		time.Now().Format("20060102-150405"),
	))

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Use reflection to get struct fields
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice")
	}

	if val.Len() == 0 {
		// Write header only for empty data
		// Get type from slice element type
		elemType := val.Type().Elem()
		if elemType.Kind() == reflect.Struct {
			headers := getCSVHeaders(elemType)
			return csvWriter.Write(headers)
		}
		return nil
	}

	// Get headers from first element
	firstElem := val.Index(0)
	headers := getCSVHeaders(firstElem.Type())
	if err := csvWriter.Write(headers); err != nil {
		return err
	}

	// Write data rows
	for i := 0; i < val.Len(); i++ {
		row := getCSVRow(val.Index(i))
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// getCSVHeaders extracts CSV headers from struct fields using json tags.
func getCSVHeaders(t reflect.Type) []string {
	headers := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag != "" && tag != "-" {
			// Remove options like "omitempty"
			for idx := 0; idx < len(tag); idx++ {
				if tag[idx] == ',' {
					tag = tag[:idx]
					break
				}
			}
			headers = append(headers, tag)
		}
	}
	return headers
}

// getCSVRow extracts values from a struct as CSV row.
func getCSVRow(v reflect.Value) []string {
	row := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		tag := fieldType.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		var value string
		switch field.Kind() {
		case reflect.String:
			value = field.String()
		case reflect.Int, reflect.Int64:
			value = strconv.FormatInt(field.Int(), 10)
		case reflect.Float64:
			value = strconv.FormatFloat(field.Float(), 'f', -1, 64)
		case reflect.Bool:
			value = strconv.FormatBool(field.Bool())
		case reflect.Struct:
			// Handle time.Time specially
			if field.Type() == reflect.TypeOf(time.Time{}) {
				t := field.Interface().(time.Time)
				value = t.Format(time.RFC3339)
			} else {
				value = fmt.Sprintf("%v", field.Interface())
			}
		default:
			value = fmt.Sprintf("%v", field.Interface())
		}
		row = append(row, value)
	}
	return row
}

// ParseQueryParams extracts and validates query parameters from the request.
// Uses Cursor API parameter names: startDate, endDate, user, page, pageSize.
func ParseQueryParams(r *http.Request) (models.Params, error) {
	params := models.Params{
		Page:     1,
		PageSize: 100,
	}

	// Parse page
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return params, fmt.Errorf("invalid page: must be >= 1")
		}
		params.Page = page
	}

	// Parse pageSize
	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			return params, fmt.Errorf("invalid pageSize: must be >= 1")
		}
		if pageSize > 1000 {
			return params, fmt.Errorf("invalid pageSize: must be <= 1000")
		}
		params.PageSize = pageSize
	}

	// Parse date range - support both new (startDate/endDate) and old (from/to) names
	params.StartDate = r.URL.Query().Get("startDate")
	params.EndDate = r.URL.Query().Get("endDate")

	// Fallback to legacy parameter names for backwards compatibility
	if params.StartDate == "" {
		params.StartDate = r.URL.Query().Get("from")
	}
	if params.EndDate == "" {
		params.EndDate = r.URL.Query().Get("to")
	}

	// Set defaults if not provided
	if params.StartDate == "" {
		params.StartDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if params.EndDate == "" {
		params.EndDate = time.Now().Format("2006-01-02")
	}

	// Parse date with support for relative formats (7d, 30d, now, today)
	startDate, err := parseDateParam(params.StartDate)
	if err != nil {
		return params, fmt.Errorf("invalid startDate: %v", err)
	}
	endDate, err := parseDateParam(params.EndDate)
	if err != nil {
		return params, fmt.Errorf("invalid endDate: %v", err)
	}

	// Store parsed dates in canonical format
	params.StartDate = startDate.Format("2006-01-02")
	params.EndDate = endDate.Format("2006-01-02")

	// Also populate legacy fields for internal use
	params.From = params.StartDate
	params.To = params.EndDate

	// Parse optional user filter (supports both 'user' and legacy 'userId')
	params.User = r.URL.Query().Get("user")
	if params.User == "" {
		params.User = r.URL.Query().Get("userId")
	}
	params.UserID = params.User // Legacy field

	// Parse optional repo filter
	params.RepoName = r.URL.Query().Get("repoName")

	return params, nil
}

// parseDateParam parses a date string supporting multiple formats.
// Supports: ISO 8601 (2025-01-15T10:30:00Z), date-only (2025-01-15),
// and relative shortcuts (7d, 30d, 90d, today, now).
func parseDateParam(param string) (time.Time, error) {
	// Try relative shortcuts first
	if len(param) > 0 && param[len(param)-1] == 'd' {
		days, err := strconv.Atoi(param[:len(param)-1])
		if err == nil {
			return time.Now().UTC().AddDate(0, 0, -days), nil
		}
	}

	switch param {
	case "today":
		now := time.Now().UTC()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
	case "now":
		return time.Now().UTC(), nil
	}

	// Try ISO 8601 with time
	t, err := time.Parse(time.RFC3339, param)
	if err == nil {
		return t, nil
	}

	// Try date-only format
	t, err = time.Parse("2006-01-02", param)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("must be YYYY-MM-DD, ISO 8601, or relative format (7d, 30d, now, today)")
}
