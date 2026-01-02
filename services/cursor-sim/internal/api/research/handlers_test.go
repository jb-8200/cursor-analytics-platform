package research

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockResearchStore implements the ResearchStore interface for testing.
type MockResearchStore struct {
	commits    []models.Commit
	prs        []models.PullRequest
	reviews    map[string][]models.ReviewComment
	developers []seed.Developer
	repos      []string
}

func NewMockResearchStore() *MockResearchStore {
	return &MockResearchStore{
		reviews: make(map[string][]models.ReviewComment),
	}
}

func (m *MockResearchStore) GetCommitsByTimeRange(from, to time.Time) []models.Commit {
	var result []models.Commit
	for _, c := range m.commits {
		if (c.CommitTs.Equal(from) || c.CommitTs.After(from)) &&
			(c.CommitTs.Equal(to) || c.CommitTs.Before(to)) {
			result = append(result, c)
		}
	}
	return result
}

func (m *MockResearchStore) GetPRsByRepo(repoName string) []models.PullRequest {
	var result []models.PullRequest
	for _, pr := range m.prs {
		if pr.RepoName == repoName {
			result = append(result, pr)
		}
	}
	return result
}

func (m *MockResearchStore) GetPRsByRepoAndState(repoName string, state models.PRState) []models.PullRequest {
	var result []models.PullRequest
	for _, pr := range m.prs {
		if pr.RepoName == repoName && pr.State == state {
			result = append(result, pr)
		}
	}
	return result
}

func (m *MockResearchStore) GetReviewComments(repoName string, prNumber int) []models.ReviewComment {
	key := repoName + ":" + string(rune(prNumber))
	return m.reviews[key]
}

func (m *MockResearchStore) GetDeveloper(userID string) (*seed.Developer, error) {
	for _, d := range m.developers {
		if d.UserID == userID {
			return &d, nil
		}
	}
	return nil, nil
}

func (m *MockResearchStore) ListRepositories() []string {
	return m.repos
}

func (m *MockResearchStore) AddCommit(c models.Commit) {
	m.commits = append(m.commits, c)
}

func (m *MockResearchStore) AddPR(pr models.PullRequest) {
	m.prs = append(m.prs, pr)
}

func (m *MockResearchStore) AddDeveloper(d seed.Developer) {
	m.developers = append(m.developers, d)
}

func (m *MockResearchStore) AddRepo(name string) {
	m.repos = append(m.repos, name)
}

func TestDatasetHandler_JSON(t *testing.T) {
	store := NewMockResearchStore()
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	// Setup test data
	store.AddDeveloper(seed.Developer{UserID: "user_001", Seniority: "senior"})
	store.AddRepo("test/repo")
	store.AddCommit(models.Commit{
		CommitHash:      "abc123",
		UserID:          "user_001",
		RepoName:        "test/repo",
		BranchName:      "feature-1",
		TotalLinesAdded: 100,
		CommitTs:        baseTime,
	})

	mergeTime := baseTime.Add(24 * time.Hour)
	store.AddPR(models.PullRequest{
		Number:     1,
		RepoName:   "test/repo",
		AuthorID:   "user_001",
		State:      models.PRStateMerged,
		HeadBranch: "feature-1",
		CreatedAt:  baseTime.Add(-1 * time.Hour),
		MergedAt:   &mergeTime,
	})

	seedData := &seed.SeedData{
		Repositories: []seed.Repository{{RepoName: "test/repo"}},
	}

	gen := generator.NewResearchGeneratorWithSeed(seedData, store, 42)
	handler := DatasetHandler(gen)

	req := httptest.NewRequest("GET", "/research/dataset?from=2026-01-14&to=2026-01-16&format=json", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "data")
}

func TestDatasetHandler_CSV(t *testing.T) {
	store := NewMockResearchStore()
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	store.AddDeveloper(seed.Developer{UserID: "user_001"})
	store.AddRepo("test/repo")
	store.AddCommit(models.Commit{
		CommitHash: "abc123",
		UserID:     "user_001",
		RepoName:   "test/repo",
		CommitTs:   baseTime,
	})

	seedData := &seed.SeedData{
		Repositories: []seed.Repository{{RepoName: "test/repo"}},
	}

	gen := generator.NewResearchGeneratorWithSeed(seedData, store, 42)
	handler := DatasetHandler(gen)

	req := httptest.NewRequest("GET", "/research/dataset?from=2026-01-14&to=2026-01-16&format=csv", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/csv", rec.Header().Get("Content-Type"))

	// Verify CSV structure
	reader := csv.NewReader(strings.NewReader(rec.Body.String()))
	records, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Greater(t, len(records), 0) // At least header
}

func TestVelocityMetricsHandler(t *testing.T) {
	store := NewMockResearchStore()
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	store.AddDeveloper(seed.Developer{UserID: "user_001"})
	store.AddRepo("test/repo")
	store.AddCommit(models.Commit{
		CommitHash:         "abc123",
		UserID:             "user_001",
		RepoName:           "test/repo",
		TotalLinesAdded:    100,
		TabLinesAdded:      30,
		ComposerLinesAdded: 20,
		CommitTs:           baseTime,
	})

	seedData := &seed.SeedData{
		Repositories: []seed.Repository{{RepoName: "test/repo"}},
	}

	gen := generator.NewResearchGeneratorWithSeed(seedData, store, 42)
	handler := VelocityMetricsHandler(gen)

	req := httptest.NewRequest("GET", "/research/metrics/velocity?from=2026-01-14&to=2026-01-16", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "data")
	assert.Contains(t, response, "params")
}

func TestReviewCostMetricsHandler(t *testing.T) {
	store := NewMockResearchStore()
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	store.AddDeveloper(seed.Developer{UserID: "user_001"})
	store.AddRepo("test/repo")
	store.AddCommit(models.Commit{
		CommitHash: "abc123",
		UserID:     "user_001",
		RepoName:   "test/repo",
		CommitTs:   baseTime,
	})

	seedData := &seed.SeedData{
		Repositories: []seed.Repository{{RepoName: "test/repo"}},
	}

	gen := generator.NewResearchGeneratorWithSeed(seedData, store, 42)
	handler := ReviewCostMetricsHandler(gen)

	req := httptest.NewRequest("GET", "/research/metrics/review-costs?from=2026-01-14&to=2026-01-16", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "data")
}

func TestQualityMetricsHandler(t *testing.T) {
	store := NewMockResearchStore()
	baseTime := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	store.AddDeveloper(seed.Developer{UserID: "user_001"})
	store.AddRepo("test/repo")
	store.AddCommit(models.Commit{
		CommitHash: "abc123",
		UserID:     "user_001",
		RepoName:   "test/repo",
		CommitTs:   baseTime,
	})

	seedData := &seed.SeedData{
		Repositories: []seed.Repository{{RepoName: "test/repo"}},
	}

	gen := generator.NewResearchGeneratorWithSeed(seedData, store, 42)
	handler := QualityMetricsHandler(gen)

	req := httptest.NewRequest("GET", "/research/metrics/quality?from=2026-01-14&to=2026-01-16", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "data")
}

func TestHandlers_InvalidDateFormat(t *testing.T) {
	store := NewMockResearchStore()
	seedData := &seed.SeedData{}
	gen := generator.NewResearchGenerator(seedData, store)

	handler := DatasetHandler(gen)

	req := httptest.NewRequest("GET", "/research/dataset?from=invalid&to=2026-01-16", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
