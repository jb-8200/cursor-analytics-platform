package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeedUploadRequest_JSON(t *testing.T) {
	jsonData := `{
		"data": "{\"version\": \"1.0\"}",
		"format": "json",
		"regenerate": true,
		"regenerate_config": {
			"mode": "override",
			"days": 180,
			"velocity": "high",
			"developers": 100,
			"max_commits": 500
		}
	}`

	var req SeedUploadRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	require.NoError(t, err)

	assert.Equal(t, "{\"version\": \"1.0\"}", req.Data)
	assert.Equal(t, "json", req.Format)
	assert.True(t, req.Regenerate)
	require.NotNil(t, req.RegenerateConfig)
	assert.Equal(t, "override", req.RegenerateConfig.Mode)
	assert.Equal(t, 180, req.RegenerateConfig.Days)
	assert.Equal(t, "high", req.RegenerateConfig.Velocity)
	assert.Equal(t, 100, req.RegenerateConfig.Developers)
	assert.Equal(t, 500, req.RegenerateConfig.MaxCommits)
}

func TestSeedUploadResponse_JSON(t *testing.T) {
	resp := SeedUploadResponse{
		Status:        "success",
		SeedLoaded:    true,
		Developers:    50,
		Repositories:  10,
		Teams:         []string{"Backend", "Frontend"},
		Divisions:     []string{"Engineering"},
		Organizations: []string{"acme-corp"},
		Regenerated:   true,
		GenerateStats: &RegenerateResponse{
			Status:       "success",
			Mode:         "override",
			DataCleaned:  true,
			CommitsAdded: 5000,
			TotalCommits: 5000,
		},
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded SeedUploadResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "success", decoded.Status)
	assert.True(t, decoded.SeedLoaded)
	assert.Equal(t, 50, decoded.Developers)
	assert.Equal(t, 10, decoded.Repositories)
	assert.Equal(t, []string{"Backend", "Frontend"}, decoded.Teams)
	assert.Equal(t, []string{"Engineering"}, decoded.Divisions)
	assert.Equal(t, []string{"acme-corp"}, decoded.Organizations)
	assert.True(t, decoded.Regenerated)
	require.NotNil(t, decoded.GenerateStats)
	assert.Equal(t, "success", decoded.GenerateStats.Status)
	assert.Equal(t, 5000, decoded.GenerateStats.CommitsAdded)
}

func TestSeedPreset_JSON(t *testing.T) {
	preset := SeedPreset{
		Name:        "small-team",
		Description: "Small team (2 developers, 2 repos, 1 region)",
		Developers:  2,
		Teams:       2,
		Regions:     []string{"US"},
	}

	data, err := json.Marshal(preset)
	require.NoError(t, err)

	var decoded SeedPreset
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "small-team", decoded.Name)
	assert.Equal(t, "Small team (2 developers, 2 repos, 1 region)", decoded.Description)
	assert.Equal(t, 2, decoded.Developers)
	assert.Equal(t, 2, decoded.Teams)
	assert.Equal(t, []string{"US"}, decoded.Regions)
}

func TestSeedPresetsResponse_JSON(t *testing.T) {
	resp := SeedPresetsResponse{
		Presets: []SeedPreset{
			{
				Name:        "small-team",
				Description: "Small team",
				Developers:  2,
				Teams:       2,
				Regions:     []string{"US"},
			},
			{
				Name:        "enterprise",
				Description: "Enterprise",
				Developers:  100,
				Teams:       10,
				Regions:     []string{"US", "EU", "APAC"},
			},
		},
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded SeedPresetsResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	require.Len(t, decoded.Presets, 2)
	assert.Equal(t, "small-team", decoded.Presets[0].Name)
	assert.Equal(t, "enterprise", decoded.Presets[1].Name)
}
