package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaginatedResponse_JSON(t *testing.T) {
	resp := PaginatedResponse{
		Data: []string{"item1", "item2"},
		Pagination: Pagination{
			Page:            1,
			PageSize:        10,
			TotalPages:      5,
			HasNextPage:     true,
			HasPreviousPage: false,
		},
		Params: Params{
			From: "2026-01-01",
			To:   "2026-01-15",
		},
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var parsed PaginatedResponse
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, 1, parsed.Pagination.Page)
	assert.Equal(t, 10, parsed.Pagination.PageSize)
	assert.True(t, parsed.Pagination.HasNextPage)
}

func TestPagination_FieldNames(t *testing.T) {
	p := Pagination{
		Page:            2,
		PageSize:        25,
		TotalUsers:      100,
		TotalPages:      4,
		HasNextPage:     true,
		HasPreviousPage: true,
	}

	data, err := json.Marshal(p)
	require.NoError(t, err)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	// Verify camelCase field names
	assert.Contains(t, raw, "page")
	assert.Contains(t, raw, "pageSize")
	assert.Contains(t, raw, "totalPages")
	assert.Contains(t, raw, "hasNextPage")
	assert.Contains(t, raw, "hasPreviousPage")
}

func TestParams_FieldNames(t *testing.T) {
	p := Params{
		From:     "2026-01-01",
		To:       "2026-01-31",
		Page:     1,
		PageSize: 100,
		UserID:   "user_001",
		RepoName: "acme/platform",
	}

	data, err := json.Marshal(p)
	require.NoError(t, err)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	// Verify camelCase for params (matching Cursor API)
	assert.Contains(t, raw, "from")
	assert.Contains(t, raw, "to")
	assert.Contains(t, raw, "page")
	assert.Contains(t, raw, "pageSize")
	assert.Contains(t, raw, "userId")
	assert.Contains(t, raw, "repoName")
}
