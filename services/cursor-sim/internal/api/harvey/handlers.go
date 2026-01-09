// Package harvey provides HTTP handlers for Harvey AI usage endpoints.
// TASK-DS-05: Create Harvey API Handler
package harvey

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// UsageResponse represents the response for the usage endpoint.
type UsageResponse struct {
	Data       []UsageEvent   `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

// UsageEvent represents a single Harvey usage event in the response.
type UsageEvent struct {
	EventID           int64   `json:"event_id"`
	MessageID         string  `json:"message_id"`
	Time              string  `json:"time"`
	User              string  `json:"User"`
	Task              string  `json:"Task"`
	ClientMatter      float64 `json:"Client Matter #"`
	Source            string  `json:"Source"`
	NumberOfDocuments int     `json:"Number of documents"`
	FeedbackComments  string  `json:"Feedback Comments"`
	FeedbackSentiment string  `json:"Feedback Sentiment"`
}

// PaginationInfo represents pagination metadata.
type PaginationInfo struct {
	Page        int  `json:"page"`
	PageSize    int  `json:"pageSize"`
	TotalCount  int  `json:"totalCount"`
	TotalPages  int  `json:"totalPages"`
	HasNextPage bool `json:"hasNextPage"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// UsageHandler returns an HTTP handler for the Harvey usage endpoint.
// GET /harvey/api/v1/history/usage
// Query params: from, to, user, task, page, page_size
func UsageHandler(store storage.HarveyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Parse query parameters
		query := r.URL.Query()

		// Parse date range (required)
		fromStr := query.Get("from")
		toStr := query.Get("to")

		from, err := parseDate(fromStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid from date: " + err.Error()})
			return
		}

		to, err := parseDate(toStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid to date: " + err.Error()})
			return
		}

		// Parse optional filters
		user := query.Get("user")
		task := storage.HarveyTask(query.Get("task"))

		// Parse pagination
		page := 1
		pageSize := 50

		if pageStr := query.Get("page"); pageStr != "" {
			page, err = strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid page: must be >= 1"})
				return
			}
		}

		if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
			pageSize, err = strconv.Atoi(pageSizeStr)
			if err != nil || pageSize < 1 {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid page_size: must be >= 1"})
				return
			}
			if pageSize > 1000 {
				pageSize = 1000
			}
		}

		// Query the store
		params := storage.HarveyParams{
			From: from,
			To:   to,
			User: user,
			Task: task,
		}

		events, err := store.GetUsage(r.Context(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to query usage data"})
			return
		}

		// Calculate pagination
		totalCount := len(events)
		totalPages := 0
		if totalCount > 0 {
			totalPages = (totalCount + pageSize - 1) / pageSize
		}

		// Apply pagination
		start := (page - 1) * pageSize
		end := start + pageSize

		var pageEvents []*storage.HarveyUsage
		if start < totalCount {
			if end > totalCount {
				end = totalCount
			}
			pageEvents = events[start:end]
		}

		// Convert to response format
		data := make([]UsageEvent, 0, len(pageEvents))
		for _, e := range pageEvents {
			data = append(data, UsageEvent{
				EventID:           e.EventID,
				MessageID:         e.MessageID,
				Time:              e.Time.Format(time.RFC3339),
				User:              e.User,
				Task:              string(e.Task),
				ClientMatter:      e.ClientMatter,
				Source:            string(e.Source),
				NumberOfDocuments: e.NumberOfDocuments,
				FeedbackComments:  e.FeedbackComments,
				FeedbackSentiment: string(e.FeedbackSentiment),
			})
		}

		response := UsageResponse{
			Data: data,
			Pagination: PaginationInfo{
				Page:        page,
				PageSize:    pageSize,
				TotalCount:  totalCount,
				TotalPages:  totalPages,
				HasNextPage: page < totalPages,
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// parseDate parses a date string in various formats.
func parseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}

	// Try RFC3339 first
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	// Try date-only format
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}

	// Try relative format (7d, 30d, etc.)
	if len(s) > 1 && s[len(s)-1] == 'd' {
		days, err := strconv.Atoi(s[:len(s)-1])
		if err == nil {
			return time.Now().UTC().AddDate(0, 0, -days), nil
		}
	}

	return time.Time{}, strconv.ErrSyntax
}
