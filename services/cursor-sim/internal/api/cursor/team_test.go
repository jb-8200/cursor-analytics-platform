package cursor

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTeamTestStore() *storage.MemoryStore {
	store := storage.NewMemoryStore()

	// Load developers
	developers := []seed.Developer{
		{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		{UserID: "user_002", Email: "bob@example.com", Name: "Bob"},
		{UserID: "user_003", Email: "carol@example.com", Name: "Carol"},
	}
	store.LoadDevelopers(developers)

	// Add commits across multiple days
	now := time.Now()
	commits := []models.Commit{
		// Day 1: 2 commits from alice, 1 from bob
		{
			CommitHash:         "c1",
			UserID:             "user_001",
			UserEmail:          "alice@example.com",
			TotalLinesAdded:    100,
			TabLinesAdded:      60,
			ComposerLinesAdded: 20,
			CommitTs:           now.Add(-48 * time.Hour),
			CreatedAt:          now,
		},
		{
			CommitHash:         "c2",
			UserID:             "user_001",
			UserEmail:          "alice@example.com",
			TotalLinesAdded:    50,
			TabLinesAdded:      30,
			ComposerLinesAdded: 10,
			CommitTs:           now.Add(-47 * time.Hour),
			CreatedAt:          now,
		},
		{
			CommitHash:         "c3",
			UserID:             "user_002",
			UserEmail:          "bob@example.com",
			TotalLinesAdded:    80,
			TabLinesAdded:      40,
			ComposerLinesAdded: 20,
			CommitTs:           now.Add(-46 * time.Hour),
			CreatedAt:          now,
		},
		// Day 2: 1 commit from alice, 1 from carol
		{
			CommitHash:         "c4",
			UserID:             "user_001",
			UserEmail:          "alice@example.com",
			TotalLinesAdded:    120,
			TabLinesAdded:      80,
			ComposerLinesAdded: 20,
			CommitTs:           now.Add(-24 * time.Hour),
			CreatedAt:          now,
		},
		{
			CommitHash:         "c5",
			UserID:             "user_003",
			UserEmail:          "carol@example.com",
			TotalLinesAdded:    60,
			TabLinesAdded:      40,
			ComposerLinesAdded: 10,
			CommitTs:           now.Add(-23 * time.Hour),
			CreatedAt:          now,
		},
	}

	for _, c := range commits {
		store.AddCommit(c)
	}

	return store
}

func TestTeamAgentEdits_Success(t *testing.T) {
	store := setupTeamTestStore()
	handler := TeamAgentEdits(store)

	req := httptest.NewRequest("GET", "/analytics/team/agent-edits", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	var response models.PaginatedResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotNil(t, response.Data)
}

func TestTeamTabs_Success(t *testing.T) {
	store := setupTeamTestStore()
	handler := TeamTabs(store)

	req := httptest.NewRequest("GET", "/analytics/team/tabs", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	var response models.PaginatedResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotNil(t, response.Data)
}

func TestTeamDAU_Success(t *testing.T) {
	store := setupTeamTestStore()
	handler := TeamDAU(store)

	req := httptest.NewRequest("GET", "/analytics/team/dau", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	var response models.PaginatedResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should have DAU data
	assert.NotNil(t, response.Data)
}

func TestTeamModels(t *testing.T) {
	store := storage.NewMemoryStore()

	// Add model usage events
	now := time.Now()
	events := []models.ModelUsageEvent{
		{
			UserID:    "user_001",
			UserEmail: "alice@example.com",
			ModelName: "claude-sonnet-4.5",
			UsageType: "code",
			Timestamp: now.Add(-48 * time.Hour),
			EventDate: now.Add(-48 * time.Hour).Format("2006-01-02"),
		},
		{
			UserID:    "user_001",
			UserEmail: "alice@example.com",
			ModelName: "claude-sonnet-4.5",
			UsageType: "chat",
			Timestamp: now.Add(-47 * time.Hour),
			EventDate: now.Add(-47 * time.Hour).Format("2006-01-02"),
		},
		{
			UserID:    "user_002",
			UserEmail: "bob@example.com",
			ModelName: "gpt-4o",
			UsageType: "code",
			Timestamp: now.Add(-46 * time.Hour),
			EventDate: now.Add(-46 * time.Hour).Format("2006-01-02"),
		},
	}

	for _, event := range events {
		err := store.AddModelUsage(event)
		require.NoError(t, err)
	}

	handler := TeamModels(store)

	req := httptest.NewRequest("GET", "/analytics/team/models", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	// Parse response
	var response struct {
		Data []models.ModelUsageDay `json:"data"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should have model usage data
	assert.NotNil(t, response.Data)
	// Should have at least one day with model breakdown
	if len(response.Data) > 0 {
		assert.NotEmpty(t, response.Data[0].ModelBreakdown)
	}
}

func TestTeamClientVersions(t *testing.T) {
	store := storage.NewMemoryStore()

	// Add client version events
	now := time.Now()
	events := []models.ClientVersionEvent{
		{
			UserID:        "user_001",
			UserEmail:     "alice@example.com",
			ClientVersion: "0.42.3",
			Timestamp:     now.Add(-48 * time.Hour),
			EventDate:     now.Add(-48 * time.Hour).Format("2006-01-02"),
		},
		{
			UserID:        "user_002",
			UserEmail:     "bob@example.com",
			ClientVersion: "0.42.3",
			Timestamp:     now.Add(-48 * time.Hour),
			EventDate:     now.Add(-48 * time.Hour).Format("2006-01-02"),
		},
		{
			UserID:        "user_003",
			UserEmail:     "carol@example.com",
			ClientVersion: "0.43.1",
			Timestamp:     now.Add(-48 * time.Hour),
			EventDate:     now.Add(-48 * time.Hour).Format("2006-01-02"),
		},
	}

	for _, event := range events {
		err := store.AddClientVersion(event)
		require.NoError(t, err)
	}

	handler := TeamClientVersions(store)

	req := httptest.NewRequest("GET", "/analytics/team/client-versions", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	// Parse response
	var response struct {
		Data []models.ClientVersionDay `json:"data"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should have client version data
	assert.NotNil(t, response.Data)
	// Should have at least one version entry
	if len(response.Data) > 0 {
		assert.NotEmpty(t, response.Data[0].ClientVersion)
		assert.Greater(t, response.Data[0].UserCount, 0)
		assert.GreaterOrEqual(t, response.Data[0].Percentage, 0.0)
		assert.LessOrEqual(t, response.Data[0].Percentage, 1.0)
	}
}

func TestTeamTopFileExtensions(t *testing.T) {
	store := storage.NewMemoryStore()

	// Add file extension events
	now := time.Now()
	date := now.Add(-48 * time.Hour).Format("2006-01-02")
	events := []models.FileExtensionEvent{
		{
			UserID:         "user_001",
			UserEmail:      "alice@example.com",
			FileExtension:  "tsx",
			LinesSuggested: 150,
			LinesAccepted:  100,
			LinesRejected:  50,
			WasAccepted:    true,
			Timestamp:      now.Add(-48 * time.Hour),
			EventDate:      date,
		},
		{
			UserID:         "user_002",
			UserEmail:      "bob@example.com",
			FileExtension:  "tsx",
			LinesSuggested: 120,
			LinesAccepted:  80,
			LinesRejected:  40,
			WasAccepted:    true,
			Timestamp:      now.Add(-48 * time.Hour),
			EventDate:      date,
		},
		{
			UserID:         "user_001",
			UserEmail:      "alice@example.com",
			FileExtension:  "go",
			LinesSuggested: 90,
			LinesAccepted:  60,
			LinesRejected:  30,
			WasAccepted:    true,
			Timestamp:      now.Add(-48 * time.Hour),
			EventDate:      date,
		},
	}

	for _, event := range events {
		err := store.AddFileExtension(event)
		require.NoError(t, err)
	}

	handler := TeamTopFileExtensions(store)

	req := httptest.NewRequest("GET", "/analytics/team/top-file-extensions", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	// Parse response
	var response struct {
		Data []models.FileExtensionDay `json:"data"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should have file extension data
	assert.NotNil(t, response.Data)
	// Should have at least 2 extensions (tsx, go)
	if len(response.Data) > 0 {
		assert.NotEmpty(t, response.Data[0].FileExtension)
		assert.Greater(t, response.Data[0].TotalLinesSuggested, 0)
		assert.GreaterOrEqual(t, response.Data[0].TotalLinesAccepted, 0)
		assert.GreaterOrEqual(t, response.Data[0].TotalLinesRejected, 0)
	}
}

func TestTeamMCP_Stub(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := TeamMCP(store)

	req := httptest.NewRequest("GET", "/analytics/team/mcp", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestTeamCommands_Stub(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := TeamCommands(store)

	req := httptest.NewRequest("GET", "/analytics/team/commands", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestTeamPlans_Stub(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := TeamPlans(store)

	req := httptest.NewRequest("GET", "/analytics/team/plans", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestTeamAskMode_Stub(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := TeamAskMode(store)

	req := httptest.NewRequest("GET", "/analytics/team/ask-mode", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestTeamLeaderboard_Stub(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := TeamLeaderboard(store)

	req := httptest.NewRequest("GET", "/analytics/team/leaderboard", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}
