package harvey

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockHarveyStore is a mock implementation of storage.HarveyStore for testing.
type mockHarveyStore struct {
	events []*storage.HarveyUsage
	err    error
}

func (m *mockHarveyStore) GetUsage(ctx context.Context, params storage.HarveyParams) ([]*storage.HarveyUsage, error) {
	if m.err != nil {
		return nil, m.err
	}

	// Apply filters
	var results []*storage.HarveyUsage
	for _, event := range m.events {
		// Time range filter
		if !event.Time.Before(params.From) && !event.Time.After(params.To) {
			// Optional user filter
			if params.User != "" && event.User != params.User {
				continue
			}
			// Optional task filter
			if params.Task != "" && event.Task != params.Task {
				continue
			}
			results = append(results, event)
		}
	}

	return results, nil
}

func (m *mockHarveyStore) StoreUsage(ctx context.Context, usage []*storage.HarveyUsage) error {
	if m.err != nil {
		return m.err
	}
	m.events = append(m.events, usage...)
	return nil
}

// createTestEvents creates test Harvey usage events.
func createTestEvents() []*storage.HarveyUsage {
	baseTime := time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC)

	return []*storage.HarveyUsage{
		{
			EventID:           1,
			MessageID:         "msg-001",
			Time:              baseTime,
			User:              "alice@firm.com",
			Task:              storage.HarveyTaskAssist,
			ClientMatter:      2024.001,
			Source:            storage.HarveySourceFiles,
			NumberOfDocuments: 2,
			FeedbackComments:  "Very helpful",
			FeedbackSentiment: storage.HarveySentimentPositive,
		},
		{
			EventID:           2,
			MessageID:         "msg-002",
			Time:              baseTime.Add(1 * time.Hour),
			User:              "bob@firm.com",
			Task:              storage.HarveyTaskDraft,
			ClientMatter:      2024.002,
			Source:            storage.HarveySourceKnowledge,
			NumberOfDocuments: 1,
			FeedbackComments:  "",
			FeedbackSentiment: storage.HarveySentimentNeutral,
		},
		{
			EventID:           3,
			MessageID:         "msg-003",
			Time:              baseTime.Add(2 * time.Hour),
			User:              "alice@firm.com",
			Task:              storage.HarveyTaskReview,
			ClientMatter:      2024.003,
			Source:            storage.HarveySourceWeb,
			NumberOfDocuments: 3,
			FeedbackComments:  "Excellent analysis",
			FeedbackSentiment: storage.HarveySentimentPositive,
		},
	}
}

func TestUsageHandler_Success(t *testing.T) {
	events := createTestEvents()
	store := &mockHarveyStore{events: events}

	handler := UsageHandler(store)

	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=2026-01-08&to=2026-01-09", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response UsageResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	// Verify data
	assert.Equal(t, 3, len(response.Data))

	// Verify pagination
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 50, response.Pagination.PageSize)
	assert.Equal(t, 3, response.Pagination.TotalCount)
	assert.Equal(t, 1, response.Pagination.TotalPages)
	assert.Equal(t, false, response.Pagination.HasNextPage)
}

func TestUsageHandler_UserFilter(t *testing.T) {
	events := createTestEvents()
	store := &mockHarveyStore{events: events}

	handler := UsageHandler(store)

	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=2026-01-08&to=2026-01-09&user=alice@firm.com", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response UsageResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	// Verify only alice's events returned
	assert.Equal(t, 2, len(response.Data))

	for _, event := range response.Data {
		assert.Equal(t, "alice@firm.com", event.User)
	}
}

func TestUsageHandler_TaskFilter(t *testing.T) {
	events := createTestEvents()
	store := &mockHarveyStore{events: events}

	handler := UsageHandler(store)

	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=2026-01-08&to=2026-01-09&task=Assist", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response UsageResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	// Verify only Assist tasks returned
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, "Assist", response.Data[0].Task)
}

func TestUsageHandler_Pagination(t *testing.T) {
	// Create many events
	var events []*storage.HarveyUsage
	baseTime := time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC)
	for i := 1; i <= 100; i++ {
		events = append(events, &storage.HarveyUsage{
			EventID:           int64(i),
			MessageID:         fmt.Sprintf("msg-%03d", i),
			Time:              baseTime.Add(time.Duration(i) * time.Minute),
			User:              "alice@firm.com",
			Task:              storage.HarveyTaskAssist,
			ClientMatter:      2024.001,
			Source:            storage.HarveySourceFiles,
			NumberOfDocuments: 1,
			FeedbackSentiment: storage.HarveySentimentPositive,
		})
	}

	store := &mockHarveyStore{events: events}
	handler := UsageHandler(store)

	// Test first page
	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=2026-01-08&to=2026-01-09&page=1&page_size=20", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response UsageResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	// Verify first page
	assert.Equal(t, 20, len(response.Data))
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 20, response.Pagination.PageSize)
	assert.Equal(t, 100, response.Pagination.TotalCount)
	assert.Equal(t, 5, response.Pagination.TotalPages)
	assert.Equal(t, true, response.Pagination.HasNextPage)
}

func TestUsageHandler_EmptyResults(t *testing.T) {
	store := &mockHarveyStore{events: []*storage.HarveyUsage{}}
	handler := UsageHandler(store)

	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=2026-01-08&to=2026-01-09", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response UsageResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	// Verify empty data
	assert.Equal(t, 0, len(response.Data))
	assert.Equal(t, 0, response.Pagination.TotalCount)
	assert.Equal(t, 0, response.Pagination.TotalPages)
	assert.Equal(t, false, response.Pagination.HasNextPage)
}

func TestUsageHandler_InvalidDateFormat(t *testing.T) {
	store := &mockHarveyStore{}
	handler := UsageHandler(store)

	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=invalid-date&to=2026-01-09", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response.Error, "from")
}

func TestUsageHandler_InvalidPageNumber(t *testing.T) {
	store := &mockHarveyStore{}
	handler := UsageHandler(store)

	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=2026-01-08&to=2026-01-09&page=-1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response.Error, "page")
}

func TestUsageHandler_StoreError(t *testing.T) {
	store := &mockHarveyStore{err: fmt.Errorf("database connection error")}
	handler := UsageHandler(store)

	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=2026-01-08&to=2026-01-09", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response.Error, "failed")
}

func TestUsageHandler_DefaultPagination(t *testing.T) {
	events := createTestEvents()
	store := &mockHarveyStore{events: events}
	handler := UsageHandler(store)

	// Request without pagination params
	req := httptest.NewRequest("GET", "/harvey/api/v1/history/usage?from=2026-01-08&to=2026-01-09", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response UsageResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	// Verify default pagination values
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 50, response.Pagination.PageSize) // Default page size
}
