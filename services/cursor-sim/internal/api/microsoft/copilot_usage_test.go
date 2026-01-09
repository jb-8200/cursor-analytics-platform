package microsoft

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCopilotStore is a mock implementation of storage.CopilotStore for testing.
type mockCopilotStore struct {
	usage []*storage.CopilotUsage
	err   error
}

func (m *mockCopilotStore) GetUsage(ctx context.Context, params storage.CopilotParams) ([]*storage.CopilotUsage, error) {
	if m.err != nil {
		return nil, m.err
	}

	// Filter by period
	var results []*storage.CopilotUsage
	periodDays := params.Period.Days()
	for _, u := range m.usage {
		if u.ReportPeriod == periodDays {
			results = append(results, u)
		}
	}

	return results, nil
}

func (m *mockCopilotStore) StoreUsage(ctx context.Context, usage []*storage.CopilotUsage) error {
	if m.err != nil {
		return m.err
	}
	m.usage = append(m.usage, usage...)
	return nil
}

// createTestCopilotUsage creates test Copilot usage data.
func createTestCopilotUsage(period int) []*storage.CopilotUsage {
	lastActivityDate := "2026-01-08"
	teamsDate := "2026-01-07"
	wordDate := "2026-01-06"

	return []*storage.CopilotUsage{
		{
			ReportRefreshDate:                     "2026-01-09",
			ReportPeriod:                          period,
			UserPrincipalName:                     "alice@company.com",
			DisplayName:                           "Alice Smith",
			LastActivityDate:                      &lastActivityDate,
			MicrosoftTeamsCopilotLastActivityDate: &teamsDate,
			WordCopilotLastActivityDate:           &wordDate,
		},
		{
			ReportRefreshDate: "2026-01-09",
			ReportPeriod:      period,
			UserPrincipalName: "bob@company.com",
			DisplayName:       "Bob Jones",
			LastActivityDate:  &lastActivityDate,
		},
	}
}

// TestCopilotUsageHandler_JSONResponse tests JSON format response.
func TestCopilotUsageHandler_JSONResponse(t *testing.T) {
	usage := createTestCopilotUsage(30)
	store := &mockCopilotStore{usage: usage}

	// Mock generator that returns empty data (store will be used)
	seedData := &seed.SeedData{
		Developers: []seed.Developer{},
	}
	gen := generator.NewCopilotGenerator(seedData)

	handler := CopilotUsageHandler(store, gen)

	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=application/json", nil)
	req.SetBasicAuth("api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response models.CopilotUsageResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	// Verify OData context is included
	assert.Contains(t, response.Context, "graph.microsoft.com")

	// Verify data
	assert.Equal(t, 2, len(response.Value))
	assert.Equal(t, "alice@company.com", response.Value[0].UserPrincipalName)
	assert.Equal(t, 30, response.Value[0].ReportPeriod)
}

// TestCopilotUsageHandler_CSVExport tests CSV format export.
func TestCopilotUsageHandler_CSVExport(t *testing.T) {
	usage := createTestCopilotUsage(30)
	store := &mockCopilotStore{usage: usage}

	seedData := &seed.SeedData{
		Developers: []seed.Developer{},
	}
	gen := generator.NewCopilotGenerator(seedData)

	handler := CopilotUsageHandler(store, gen)

	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=text/csv", nil)
	req.SetBasicAuth("api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/csv", rec.Header().Get("Content-Type"))
	assert.Equal(t, "attachment; filename=copilot-usage-D30.csv", rec.Header().Get("Content-Disposition"))

	// Verify CSV content
	reader := csv.NewReader(rec.Body)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Should have header + 2 data rows
	assert.Equal(t, 3, len(records))

	// Verify header
	assert.Contains(t, records[0], "Report Refresh Date")
	assert.Contains(t, records[0], "User Principal Name")
	assert.Contains(t, records[0], "Display Name")

	// Verify data
	assert.Equal(t, "alice@company.com", records[1][2]) // User Principal Name column
}

// TestCopilotUsageHandler_AllPeriods tests all supported periods.
func TestCopilotUsageHandler_AllPeriods(t *testing.T) {
	tests := []struct {
		name       string
		period     string
		periodDays int
	}{
		{"D7", "D7", 7},
		{"D30", "D30", 30},
		{"D90", "D90", 90},
		{"D180", "D180", 180},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage := createTestCopilotUsage(tt.periodDays)
			store := &mockCopilotStore{usage: usage}

			seedData := &seed.SeedData{
				Developers: []seed.Developer{},
			}
			gen := generator.NewCopilotGenerator(seedData)

			handler := CopilotUsageHandler(store, gen)

			url := fmt.Sprintf("/reports/getMicrosoft365CopilotUsageUserDetail(period='%s')", tt.period)
			req := httptest.NewRequest("GET", url, nil)
			req.SetBasicAuth("api-key", "")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.CopilotUsageResponse
			err := json.NewDecoder(rec.Body).Decode(&response)
			require.NoError(t, err)

			// Verify correct period
			for _, user := range response.Value {
				assert.Equal(t, tt.periodDays, user.ReportPeriod)
			}
		})
	}
}

// TestCopilotUsageHandler_InvalidPeriod tests invalid period parameter.
func TestCopilotUsageHandler_InvalidPeriod(t *testing.T) {
	store := &mockCopilotStore{}
	seedData := &seed.SeedData{}
	gen := generator.NewCopilotGenerator(seedData)

	handler := CopilotUsageHandler(store, gen)

	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='INVALID')", nil)
	req.SetBasicAuth("api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&errResp)
	require.NoError(t, err)

	assert.Contains(t, errResp.Error, "period")
}

// TestCopilotUsageHandler_MissingPeriod tests missing period parameter.
func TestCopilotUsageHandler_MissingPeriod(t *testing.T) {
	store := &mockCopilotStore{}
	seedData := &seed.SeedData{}
	gen := generator.NewCopilotGenerator(seedData)

	handler := CopilotUsageHandler(store, gen)

	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail", nil)
	req.SetBasicAuth("api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestCopilotUsageHandler_GeneratorFallback tests generator fallback when store is empty.
func TestCopilotUsageHandler_GeneratorFallback(t *testing.T) {
	// Empty store
	store := &mockCopilotStore{usage: []*storage.CopilotUsage{}}

	// Seed data with developers for generator
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				UserID: "dev1",
				Name:   "Developer One",
				Email:  "dev1@company.com",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled: true,
			},
		},
	}
	gen := generator.NewCopilotGenerator(seedData)

	handler := CopilotUsageHandler(store, gen)

	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')", nil)
	req.SetBasicAuth("api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.CopilotUsageResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	// Should have generated data for the developer
	assert.Equal(t, 1, len(response.Value))
	assert.Equal(t, "dev1@company.com", response.Value[0].UserPrincipalName)
}

// TestCopilotUsageHandler_StoreError tests storage error handling.
func TestCopilotUsageHandler_StoreError(t *testing.T) {
	store := &mockCopilotStore{err: fmt.Errorf("database connection error")}
	seedData := &seed.SeedData{}
	gen := generator.NewCopilotGenerator(seedData)

	handler := CopilotUsageHandler(store, gen)

	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')", nil)
	req.SetBasicAuth("api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var errResp ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&errResp)
	require.NoError(t, err)

	assert.Contains(t, errResp.Error, "failed")
}

// TestCopilotUsageHandler_NoAuthentication tests missing authentication.
func TestCopilotUsageHandler_NoAuthentication(t *testing.T) {
	store := &mockCopilotStore{}
	seedData := &seed.SeedData{}
	gen := generator.NewCopilotGenerator(seedData)

	handler := CopilotUsageHandler(store, gen)

	// Request without authentication
	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// TestCopilotUsageHandler_EmptyResults tests empty results scenario.
func TestCopilotUsageHandler_EmptyResults(t *testing.T) {
	// Empty store and no developers in seed data
	store := &mockCopilotStore{usage: []*storage.CopilotUsage{}}
	seedData := &seed.SeedData{
		Developers: []seed.Developer{},
	}
	gen := generator.NewCopilotGenerator(seedData)

	handler := CopilotUsageHandler(store, gen)

	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')", nil)
	req.SetBasicAuth("api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.CopilotUsageResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)

	// Should return empty value array
	assert.Equal(t, 0, len(response.Value))
}

// TestCopilotUsageHandler_DefaultFormat tests default format (JSON).
func TestCopilotUsageHandler_DefaultFormat(t *testing.T) {
	usage := createTestCopilotUsage(30)
	store := &mockCopilotStore{usage: usage}

	seedData := &seed.SeedData{
		Developers: []seed.Developer{},
	}
	gen := generator.NewCopilotGenerator(seedData)

	handler := CopilotUsageHandler(store, gen)

	// Request without $format parameter - should default to JSON
	req := httptest.NewRequest("GET", "/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')", nil)
	req.SetBasicAuth("api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}
