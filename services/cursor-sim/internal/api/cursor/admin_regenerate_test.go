package cursor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRegenerateAppendMode tests append mode adds data without clearing.
func TestRegenerateAppendMode(t *testing.T) {
	// Create store with some initial data
	store := storage.NewMemoryStore()

	// Load minimal seed data
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID:         "user1",
				Email:          "user1@example.com",
				Name:           "User 1",
				Org:            "TestOrg",
				Division:       "Engineering",
				Team:           "Backend",
				Region:         "US",
				Seniority:      "mid",
				ActivityLevel:  "high",
				AcceptanceRate: 0.8,
				PRBehavior: seed.PRBehavior{
					PRsPerWeek:    2.0,
					AvgPRSizeLOC:  100,
					AvgFilesPerPR: 5,
				},
				CodingSpeed: seed.CodingSpeed{
					Mean: 4.0,
					Std:  1.0,
				},
				ChatVsCodeRatio: seed.ChatCodeRatio{
					Chat: 0.3,
					Code: 0.7,
				},
				WorkingHoursBand: seed.WorkingHours{
					Start: 9,
					End:   17,
					Peak:  13,
				},
			},
		},
		Repositories: []seed.Repository{
			{
				RepoName:        "test/repo",
				PrimaryLanguage: "go",
				ServiceType:     "api",
				DefaultBranch:   "main",
				Teams:           []string{"Backend"},
			},
		},
		Correlations: seed.Correlations{},
		PRLifecycle:  seed.PRLifecycle{},
	}

	// Get initial stats
	statsBefore := store.GetStats()

	// Create request
	reqBody := models.RegenerateRequest{
		Mode:       "append",
		Days:       7,
		Velocity:   "low",
		Developers: 0, // Use seed count
		MaxCommits: 10,
	}

	reqBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	// Create handler
	handler := Regenerate(store, seedData)

	// Make request
	req := httptest.NewRequest(http.MethodPost, "/admin/regenerate", bytes.NewReader(reqBytes))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Check response
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.RegenerateResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "append", resp.Mode)
	assert.False(t, resp.DataCleaned, "Append mode should not clean data")
	assert.Greater(t, resp.CommitsAdded, 0, "Should have added commits")

	// Verify data was added (not replaced)
	statsAfter := store.GetStats()
	assert.Equal(t, statsBefore.Commits+resp.CommitsAdded, statsAfter.Commits)
}

// TestRegenerateOverrideMode tests override mode clears data before generating.
func TestRegenerateOverrideMode(t *testing.T) {
	// Create store with some initial data
	store := storage.NewMemoryStore()

	// Load minimal seed data
	seedData := &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID:         "user1",
				Email:          "user1@example.com",
				Name:           "User 1",
				Org:            "TestOrg",
				Division:       "Engineering",
				Team:           "Backend",
				Region:         "US",
				Seniority:      "mid",
				ActivityLevel:  "high",
				AcceptanceRate: 0.8,
				PRBehavior: seed.PRBehavior{
					PRsPerWeek:    2.0,
					AvgPRSizeLOC:  100,
					AvgFilesPerPR: 5,
				},
				CodingSpeed: seed.CodingSpeed{
					Mean: 4.0,
					Std:  1.0,
				},
				ChatVsCodeRatio: seed.ChatCodeRatio{
					Chat: 0.3,
					Code: 0.7,
				},
				WorkingHoursBand: seed.WorkingHours{
					Start: 9,
					End:   17,
					Peak:  13,
				},
			},
		},
		Repositories: []seed.Repository{
			{
				RepoName:        "test/repo",
				PrimaryLanguage: "go",
				ServiceType:     "api",
				DefaultBranch:   "main",
				Teams:           []string{"Backend"},
			},
		},
		Correlations: seed.Correlations{},
		PRLifecycle:  seed.PRLifecycle{},
	}

	// Create request
	reqBody := models.RegenerateRequest{
		Mode:       "override",
		Days:       7,
		Velocity:   "low",
		Developers: 0,
		MaxCommits: 10,
	}

	reqBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	// Create handler
	handler := Regenerate(store, seedData)

	// Make request
	req := httptest.NewRequest(http.MethodPost, "/admin/regenerate", bytes.NewReader(reqBytes))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Check response
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.RegenerateResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "override", resp.Mode)
	assert.True(t, resp.DataCleaned, "Override mode should clean data")
	assert.Greater(t, resp.TotalCommits, 0, "Should have generated commits")
	assert.Equal(t, resp.CommitsAdded, resp.TotalCommits, "CommitsAdded should equal TotalCommits in override mode")
}

// TestRegenerateInvalidMode tests validation of mode parameter.
func TestRegenerateInvalidMode(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := &seed.SeedData{Version: "1.0"}

	reqBody := models.RegenerateRequest{
		Mode:       "invalid",
		Days:       7,
		Velocity:   "medium",
		Developers: 0,
		MaxCommits: 0,
	}

	reqBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	handler := Regenerate(store, seedData)
	req := httptest.NewRequest(http.MethodPost, "/admin/regenerate", bytes.NewReader(reqBytes))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)

	assert.Contains(t, errResp["error"], "mode must be")
}

// TestRegenerateInvalidVelocity tests validation of velocity parameter.
func TestRegenerateInvalidVelocity(t *testing.T) {
	store := storage.NewMemoryStore()
	seedData := &seed.SeedData{Version: "1.0"}

	reqBody := models.RegenerateRequest{
		Mode:       "append",
		Days:       7,
		Velocity:   "super-fast", // Invalid
		Developers: 0,
		MaxCommits: 0,
	}

	reqBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	handler := Regenerate(store, seedData)
	req := httptest.NewRequest(http.MethodPost, "/admin/regenerate", bytes.NewReader(reqBytes))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)

	assert.Contains(t, errResp["error"], "velocity must be")
}

// TestValidateRegenerateRequest tests the validation function.
func TestValidateRegenerateRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     models.RegenerateRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid append request",
			req: models.RegenerateRequest{
				Mode:       "append",
				Days:       30,
				Velocity:   "medium",
				Developers: 10,
				MaxCommits: 100,
			},
			wantErr: false,
		},
		{
			name: "valid override request",
			req: models.RegenerateRequest{
				Mode:       "override",
				Days:       90,
				Velocity:   "high",
				Developers: 0,
				MaxCommits: 0,
			},
			wantErr: false,
		},
		{
			name: "invalid mode",
			req: models.RegenerateRequest{
				Mode:       "replace",
				Days:       30,
				Velocity:   "medium",
				Developers: 10,
				MaxCommits: 100,
			},
			wantErr: true,
			errMsg:  "mode must be 'append' or 'override'",
		},
		{
			name: "days too low",
			req: models.RegenerateRequest{
				Mode:       "append",
				Days:       0,
				Velocity:   "medium",
				Developers: 10,
				MaxCommits: 100,
			},
			wantErr: true,
			errMsg:  "days must be between 1 and 3650",
		},
		{
			name: "days too high",
			req: models.RegenerateRequest{
				Mode:       "append",
				Days:       4000,
				Velocity:   "medium",
				Developers: 10,
				MaxCommits: 100,
			},
			wantErr: true,
			errMsg:  "days must be between 1 and 3650",
		},
		{
			name: "invalid velocity",
			req: models.RegenerateRequest{
				Mode:       "append",
				Days:       30,
				Velocity:   "blazing",
				Developers: 10,
				MaxCommits: 100,
			},
			wantErr: true,
			errMsg:  "velocity must be 'low', 'medium', or 'high'",
		},
		{
			name: "developers too high",
			req: models.RegenerateRequest{
				Mode:       "append",
				Days:       30,
				Velocity:   "medium",
				Developers: 20000,
				MaxCommits: 100,
			},
			wantErr: true,
			errMsg:  "developers must be between 0 and 10000",
		},
		{
			name: "max_commits too high",
			req: models.RegenerateRequest{
				Mode:       "append",
				Days:       30,
				Velocity:   "medium",
				Developers: 10,
				MaxCommits: 200000,
			},
			wantErr: true,
			errMsg:  "max_commits must be between 0 and 100000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegenerateRequest(&tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
