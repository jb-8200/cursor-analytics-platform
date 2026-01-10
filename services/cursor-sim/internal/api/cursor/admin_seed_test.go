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

// TestGetSeedPresets tests getting predefined seed presets.
func TestGetSeedPresets(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/admin/seed/presets", nil)
	w := httptest.NewRecorder()

	handler := GetSeedPresets()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.SeedPresetsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	// Should return 4 presets
	require.Len(t, resp.Presets, 4)

	// Verify preset names
	presetNames := make([]string, len(resp.Presets))
	for i, p := range resp.Presets {
		presetNames[i] = p.Name
	}
	assert.Contains(t, presetNames, "small-team")
	assert.Contains(t, presetNames, "medium-team")
	assert.Contains(t, presetNames, "enterprise")
	assert.Contains(t, presetNames, "multi-region")

	// Verify each preset has required fields
	for _, preset := range resp.Presets {
		assert.NotEmpty(t, preset.Name)
		assert.NotEmpty(t, preset.Description)
		assert.Greater(t, preset.Developers, 0)
		assert.Greater(t, preset.Teams, 0)
		assert.NotEmpty(t, preset.Regions)
	}
}

// TestSeedUpload_InvalidFormat tests rejection of invalid format.
func TestSeedUpload_InvalidFormat(t *testing.T) {
	store := storage.NewMemoryStore()
	currentSeed := &seed.SeedData{}

	reqBody := models.SeedUploadRequest{
		Data:       "some data",
		Format:     "xml", // Invalid format
		Regenerate: false,
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admin/seed", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := UploadSeed(store, &currentSeed)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Contains(t, errResp["error"], "invalid format")
}

// TestSeedUpload_MethodNotAllowed tests that only POST is allowed.
func TestSeedUpload_MethodNotAllowed(t *testing.T) {
	store := storage.NewMemoryStore()
	currentSeed := &seed.SeedData{}

	req := httptest.NewRequest(http.MethodGet, "/admin/seed", nil)
	w := httptest.NewRecorder()

	handler := UploadSeed(store, &currentSeed)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
