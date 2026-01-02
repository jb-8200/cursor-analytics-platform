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

	// Parse date range
	params.From = r.URL.Query().Get("from")
	params.To = r.URL.Query().Get("to")

	// Set defaults if not provided
	if params.From == "" {
		params.From = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if params.To == "" {
		params.To = time.Now().Format("2006-01-02")
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", params.From); err != nil {
		return params, fmt.Errorf("invalid from date: must be YYYY-MM-DD format")
	}
	if _, err := time.Parse("2006-01-02", params.To); err != nil {
		return params, fmt.Errorf("invalid to date: must be YYYY-MM-DD format")
	}

	// Parse optional filters
	params.UserID = r.URL.Query().Get("userId")
	params.RepoName = r.URL.Query().Get("repoName")

	return params, nil
}
