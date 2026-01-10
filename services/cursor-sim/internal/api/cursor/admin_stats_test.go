package cursor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/models"
	csmodels "github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStats_Basic(t *testing.T) {
	// Setup test store and seed data
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()

	// Load developers
	require.NoError(t, store.LoadDevelopers(seedData.Developers))

	// Add test commits
	now := time.Now()
	for i := 0; i < 100; i++ {
		commit := csmodels.Commit{
			CommitHash:           generateCommitHash(i),
			UserID:               seedData.Developers[i%len(seedData.Developers)].UserID,
			UserEmail:            seedData.Developers[i%len(seedData.Developers)].Email,
			UserName:             seedData.Developers[i%len(seedData.Developers)].Name,
			RepoName:             "acme/platform",
			TotalLinesAdded:      150,
			TotalLinesDeleted:    45,
			TabLinesAdded:        90,
			TabLinesDeleted:      20,
			ComposerLinesAdded:   35,
			ComposerLinesDeleted: 10,
			NonAILinesAdded:      25,
			NonAILinesDeleted:    15,
			CommitTs:             now.Add(time.Duration(-i) * time.Hour),
			CreatedAt:            now.Add(time.Duration(-i) * time.Hour),
		}
		require.NoError(t, store.AddCommit(commit))
	}

	// Add test PRs
	for i := 0; i < 20; i++ {
		createdAt := now.Add(time.Duration(-i*24) * time.Hour)
		mergedAt := createdAt.Add(time.Duration(48) * time.Hour)
		pr := csmodels.PullRequest{
			Number:      i + 1,
			Title:       "Test PR",
			State:       csmodels.PRStateMerged,
			AuthorID:    seedData.Developers[i%len(seedData.Developers)].UserID,
			AuthorEmail: seedData.Developers[i%len(seedData.Developers)].Email,
			AuthorName:  seedData.Developers[i%len(seedData.Developers)].Name,
			RepoName:    "acme/platform",
			Additions:   200,
			Deletions:   50,
			CommitCount: 5,
			CreatedAt:   createdAt,
			MergedAt:    &mergedAt,
		}
		require.NoError(t, store.StorePR(pr))
	}

	// Add test reviews
	for i := 0; i < 40; i++ {
		review := csmodels.Review{
			ID:          i + 1,
			PRID:        i%20 + 1,
			Reviewer:    seedData.Developers[(i+1)%len(seedData.Developers)].Email,
			State:       csmodels.ReviewStateApproved,
			Body:        "LGTM",
			SubmittedAt: now.Add(time.Duration(-i*12) * time.Hour),
			Comments:    []csmodels.ReviewComment{{Body: "Comment 1"}, {Body: "Comment 2"}},
		}
		require.NoError(t, store.StoreReview(review))
	}

	// Add test issues
	for i := 0; i < 10; i++ {
		issue := csmodels.Issue{
			Number:    i + 1,
			Title:     "Test Issue",
			State:     csmodels.IssueStateClosed,
			RepoName:  "acme/platform",
			AuthorID:  seedData.Developers[i%len(seedData.Developers)].UserID,
			CreatedAt: now.Add(time.Duration(-i*24) * time.Hour),
			ClosedAt:  timePtr(now.Add(time.Duration(-i*24+48) * time.Hour)),
		}
		require.NoError(t, store.StoreIssue(issue))
	}

	// Create handler
	handler := GetStats(store, seedData)

	// Test request without time series
	req := httptest.NewRequest(http.MethodGet, "/admin/stats", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	// Parse response
	var response models.StatsResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&response))

	// Verify generation stats
	assert.Equal(t, 100, response.Generation.TotalCommits)
	assert.Equal(t, 20, response.Generation.TotalPRs)
	assert.Equal(t, 40, response.Generation.TotalReviews)
	assert.Equal(t, 10, response.Generation.TotalIssues)
	assert.Equal(t, 3, response.Generation.TotalDevelopers)
	assert.NotEmpty(t, response.Generation.DataSize)

	// Verify developer stats
	assert.Equal(t, 3, response.Developers.BySeniority["senior"]+response.Developers.BySeniority["mid"]+response.Developers.BySeniority["junior"])
	assert.Equal(t, 3, response.Developers.ByRegion["us-west"]+response.Developers.ByRegion["us-east"])
	assert.Equal(t, 3, response.Developers.ByTeam["Platform"]+response.Developers.ByTeam["API"])
	assert.Equal(t, 3, response.Developers.ByActivity["high"]+response.Developers.ByActivity["medium"])

	// Verify quality metrics
	assert.Greater(t, response.Quality.AvgReviewThoroughness, 0.0)
	assert.Greater(t, response.Quality.AvgIterations, 0.0)

	// Verify variance metrics (should be calculated from actual data)
	assert.GreaterOrEqual(t, response.Variance.CommitsStdDev, 0.0)
	assert.GreaterOrEqual(t, response.Variance.PRSizeStdDev, 0.0)
	assert.GreaterOrEqual(t, response.Variance.CycleTimeStdDev, 0.0)

	// Verify performance metrics
	assert.NotEmpty(t, response.Performance.MemoryUsage)
	assert.NotEmpty(t, response.Performance.StorageEfficiency)

	// Verify organization
	assert.Contains(t, response.Organization.Teams, "Platform")
	assert.Contains(t, response.Organization.Divisions, "Engineering")
	assert.Contains(t, response.Organization.Repositories, "acme/platform")

	// Verify time series is nil (not requested)
	assert.Nil(t, response.TimeSeries)
}

func TestGetStats_WithTimeSeries(t *testing.T) {
	// Setup test store and seed data
	store := storage.NewMemoryStore()
	seedData := createTestSeedData()

	// Load developers
	require.NoError(t, store.LoadDevelopers(seedData.Developers))

	// Add test commits spread over 30 days
	now := time.Now()
	for i := 0; i < 300; i++ {
		commit := csmodels.Commit{
			CommitHash:           generateCommitHash(i),
			UserID:               seedData.Developers[i%len(seedData.Developers)].UserID,
			UserEmail:            seedData.Developers[i%len(seedData.Developers)].Email,
			UserName:             seedData.Developers[i%len(seedData.Developers)].Name,
			RepoName:             "acme/platform",
			TotalLinesAdded:      150,
			TotalLinesDeleted:    45,
			TabLinesAdded:        90,
			TabLinesDeleted:      20,
			ComposerLinesAdded:   35,
			ComposerLinesDeleted: 10,
			NonAILinesAdded:      25,
			NonAILinesDeleted:    15,
			CommitTs:             now.Add(time.Duration(-i) * time.Hour),
			CreatedAt:            now.Add(time.Duration(-i) * time.Hour),
		}
		require.NoError(t, store.AddCommit(commit))
	}

	// Add test PRs spread over 30 days
	for i := 0; i < 60; i++ {
		createdAt := now.Add(time.Duration(-i*12) * time.Hour)
		mergedAt := createdAt.Add(time.Duration(48) * time.Hour)
		pr := csmodels.PullRequest{
			Number:      i + 1,
			Title:       "Test PR",
			State:       csmodels.PRStateMerged,
			AuthorID:    seedData.Developers[i%len(seedData.Developers)].UserID,
			AuthorEmail: seedData.Developers[i%len(seedData.Developers)].Email,
			AuthorName:  seedData.Developers[i%len(seedData.Developers)].Name,
			RepoName:    "acme/platform",
			Additions:   200,
			Deletions:   50,
			CommitCount: 5,
			CreatedAt:   createdAt,
			MergedAt:    &mergedAt,
		}
		require.NoError(t, store.StorePR(pr))
	}

	// Create handler
	handler := GetStats(store, seedData)

	// Test request WITH time series
	req := httptest.NewRequest(http.MethodGet, "/admin/stats?include_timeseries=true", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response models.StatsResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&response))

	// Verify time series is included
	require.NotNil(t, response.TimeSeries)
	assert.NotEmpty(t, response.TimeSeries.CommitsPerDay)
	assert.NotEmpty(t, response.TimeSeries.PRsPerDay)
	assert.NotEmpty(t, response.TimeSeries.AvgCycleTime)

	// Verify time series data length matches (should be same length)
	assert.Equal(t, len(response.TimeSeries.CommitsPerDay), len(response.TimeSeries.PRsPerDay))
	assert.Equal(t, len(response.TimeSeries.CommitsPerDay), len(response.TimeSeries.AvgCycleTime))

	// Verify time series data has non-zero values
	totalCommits := 0
	for _, count := range response.TimeSeries.CommitsPerDay {
		totalCommits += count
	}
	assert.Greater(t, totalCommits, 0, "Time series should have commits")

	totalPRs := 0
	for _, count := range response.TimeSeries.PRsPerDay {
		totalPRs += count
	}
	assert.Greater(t, totalPRs, 0, "Time series should have PRs")
}

func TestGetStats_Calculations(t *testing.T) {
	// Test specific calculation functions
	t.Run("groupBySeniority", func(t *testing.T) {
		seedData := &seed.SeedData{
			Developers: []seed.Developer{
				{Seniority: "senior"},
				{Seniority: "senior"},
				{Seniority: "mid"},
				{Seniority: "junior"},
			},
		}
		result := groupBySeniority(seedData)
		assert.Equal(t, 2, result["senior"])
		assert.Equal(t, 1, result["mid"])
		assert.Equal(t, 1, result["junior"])
	})

	t.Run("groupByRegion", func(t *testing.T) {
		seedData := &seed.SeedData{
			Developers: []seed.Developer{
				{Region: "us-west"},
				{Region: "us-west"},
				{Region: "us-east"},
			},
		}
		result := groupByRegion(seedData)
		assert.Equal(t, 2, result["us-west"])
		assert.Equal(t, 1, result["us-east"])
	})

	t.Run("groupByTeam", func(t *testing.T) {
		seedData := &seed.SeedData{
			Developers: []seed.Developer{
				{Team: "Platform"},
				{Team: "Platform"},
				{Team: "API"},
			},
		}
		result := groupByTeam(seedData)
		assert.Equal(t, 2, result["Platform"])
		assert.Equal(t, 1, result["API"])
	})

	t.Run("groupByActivity", func(t *testing.T) {
		seedData := &seed.SeedData{
			Developers: []seed.Developer{
				{ActivityLevel: "high"},
				{ActivityLevel: "high"},
				{ActivityLevel: "medium"},
			},
		}
		result := groupByActivity(seedData)
		assert.Equal(t, 2, result["high"])
		assert.Equal(t, 1, result["medium"])
	})

	t.Run("calculateStdDev", func(t *testing.T) {
		values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
		stdDev := calculateStdDev(values)
		assert.InDelta(t, 2.0, stdDev, 0.1) // Approximately 2.0
	})

	t.Run("calculateStdDev_empty", func(t *testing.T) {
		values := []float64{}
		stdDev := calculateStdDev(values)
		assert.Equal(t, 0.0, stdDev)
	})

	t.Run("formatBytes", func(t *testing.T) {
		assert.Equal(t, "500 B", formatBytes(500))
		assert.Equal(t, "1.0 KB", formatBytes(1024))
		assert.Equal(t, "1.5 MB", formatBytes(1536*1024))
		assert.Equal(t, "2.0 GB", formatBytes(2048*1024*1024))
	})

	t.Run("estimateDataSize", func(t *testing.T) {
		size := estimateDataSize(1000, 100, 200, 50)
		expectedSize := 1000*500 + 100*1024 + 200*300 + 50*500 // commits + prs + reviews + issues
		assert.Equal(t, expectedSize, size)
	})
}

// Helper functions

func createTestSeedData() *seed.SeedData {
	return &seed.SeedData{
		Version: "1.0",
		Developers: []seed.Developer{
			{
				UserID:        "user_001",
				Email:         "alice@example.com",
				Name:          "Alice Developer",
				Org:           "acme-corp",
				Division:      "Engineering",
				Team:          "Platform",
				Region:        "us-west",
				Seniority:     "senior",
				ActivityLevel: "high",
			},
			{
				UserID:        "user_002",
				Email:         "bob@example.com",
				Name:          "Bob Developer",
				Org:           "acme-corp",
				Division:      "Engineering",
				Team:          "API",
				Region:        "us-east",
				Seniority:     "mid",
				ActivityLevel: "medium",
			},
			{
				UserID:        "user_003",
				Email:         "carol@example.com",
				Name:          "Carol Developer",
				Org:           "acme-corp",
				Division:      "Infrastructure",
				Team:          "Platform",
				Region:        "us-west",
				Seniority:     "junior",
				ActivityLevel: "medium",
			},
		},
		Repositories: []seed.Repository{
			{
				RepoName:        "acme/platform",
				PrimaryLanguage: "go",
				ServiceType:     "backend",
				DefaultBranch:   "main",
				Teams:           []string{"Platform", "API"},
			},
		},
	}
}

func generateCommitHash(i int) string {
	return fmt.Sprintf("%040x", i)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
