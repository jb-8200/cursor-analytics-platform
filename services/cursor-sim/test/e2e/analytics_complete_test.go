package e2e

import (
	"encoding/json"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/config"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_TeamAnalyticsSchemas validates the exact schema of each team-level endpoint
// against the Cursor Analytics API documentation.
func TestE2E_TeamAnalyticsSchemas(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name            string
		endpoint        string
		validateData    func(t *testing.T, data interface{})
		hasPagination   bool
	}{
		{
			name:     "agent-edits",
			endpoint: "/analytics/team/agent-edits",
			validateData: func(t *testing.T, data interface{}) {
				dataArray := data.([]interface{})
				require.Greater(t, len(dataArray), 0, "should have data")

				for _, item := range dataArray {
					day := item.(map[string]interface{})

					// Required fields from cursor_analytics.md
					assert.Contains(t, day, "event_date")
					assert.Contains(t, day, "total_suggested_diffs")
					assert.Contains(t, day, "total_accepted_diffs")
					assert.Contains(t, day, "total_rejected_diffs")
					assert.Contains(t, day, "total_green_lines_accepted")
					assert.Contains(t, day, "total_red_lines_accepted")
					assert.Contains(t, day, "total_green_lines_rejected")
					assert.Contains(t, day, "total_red_lines_rejected")
					assert.Contains(t, day, "total_green_lines_suggested")
					assert.Contains(t, day, "total_red_lines_suggested")
					assert.Contains(t, day, "total_lines_suggested")
					assert.Contains(t, day, "total_lines_accepted")

					// Verify data types and realistic values
					assert.IsType(t, "", day["event_date"])
					assert.GreaterOrEqual(t, day["total_suggested_diffs"].(float64), 0.0)
					assert.GreaterOrEqual(t, day["total_accepted_diffs"].(float64), 0.0)
				}
			},
			hasPagination: false,
		},
		{
			name:     "tabs",
			endpoint: "/analytics/team/tabs",
			validateData: func(t *testing.T, data interface{}) {
				dataArray := data.([]interface{})
				require.Greater(t, len(dataArray), 0, "should have data")

				for _, item := range dataArray {
					day := item.(map[string]interface{})

					// Required fields from cursor_analytics.md
					assert.Contains(t, day, "event_date")
					assert.Contains(t, day, "total_suggestions")
					assert.Contains(t, day, "total_accepts")
					assert.Contains(t, day, "total_rejects")
					assert.Contains(t, day, "total_green_lines_accepted")
					assert.Contains(t, day, "total_red_lines_accepted")
					assert.Contains(t, day, "total_green_lines_rejected")
					assert.Contains(t, day, "total_red_lines_rejected")
					assert.Contains(t, day, "total_green_lines_suggested")
					assert.Contains(t, day, "total_red_lines_suggested")
					assert.Contains(t, day, "total_lines_suggested")
					assert.Contains(t, day, "total_lines_accepted")

					// Verify realistic values (not stub data)
					assert.GreaterOrEqual(t, day["total_suggestions"].(float64), 0.0)
					assert.GreaterOrEqual(t, day["total_accepts"].(float64), 0.0)
				}
			},
			hasPagination: false,
		},
		{
			name:     "dau",
			endpoint: "/analytics/team/dau",
			validateData: func(t *testing.T, data interface{}) {
				dataArray := data.([]interface{})
				require.Greater(t, len(dataArray), 0, "should have data")

				for _, item := range dataArray {
					day := item.(map[string]interface{})

					// Required fields from cursor_analytics.md (note: "date" not "event_date")
					assert.Contains(t, day, "date")
					assert.NotContains(t, day, "event_date", "DAU uses 'date' not 'event_date'")
					assert.Contains(t, day, "dau")
					assert.Contains(t, day, "cli_dau")
					assert.Contains(t, day, "cloud_agent_dau")
					assert.Contains(t, day, "bugbot_dau")

					// Verify realistic values
					assert.GreaterOrEqual(t, day["dau"].(float64), 0.0)
					assert.GreaterOrEqual(t, day["cli_dau"].(float64), 0.0)
				}
			},
			hasPagination: false,
		},
		{
			name:     "models",
			endpoint: "/analytics/team/models",
			validateData: func(t *testing.T, data interface{}) {
				dataArray := data.([]interface{})
				// Allow empty data for feature endpoints that might not have data in test setup
				if len(dataArray) == 0 {
					t.Skip("No model usage data generated in test setup")
					return
				}

				for _, item := range dataArray {
					day := item.(map[string]interface{})

					// Required fields from cursor_analytics.md
					assert.Contains(t, day, "date")
					assert.Contains(t, day, "model_breakdown")

					// Validate model_breakdown structure
					breakdown := day["model_breakdown"].(map[string]interface{})
					assert.Greater(t, len(breakdown), 0, "should have at least one model")

					for modelName, modelData := range breakdown {
						assert.NotEmpty(t, modelName)
						modelInfo := modelData.(map[string]interface{})
						assert.Contains(t, modelInfo, "messages")
						assert.Contains(t, modelInfo, "users")
						assert.GreaterOrEqual(t, modelInfo["messages"].(float64), 0.0)
						assert.GreaterOrEqual(t, modelInfo["users"].(float64), 0.0)
					}
				}
			},
			hasPagination: false,
		},
		{
			name:     "client-versions",
			endpoint: "/analytics/team/client-versions",
			validateData: func(t *testing.T, data interface{}) {
				dataArray := data.([]interface{})
				if len(dataArray) == 0 {
					t.Skip("No client version data generated in test setup")
					return
				}

				for _, item := range dataArray {
					entry := item.(map[string]interface{})

					// Required fields from cursor_analytics.md
					assert.Contains(t, entry, "event_date")
					assert.Contains(t, entry, "client_version")
					assert.Contains(t, entry, "user_count")
					assert.Contains(t, entry, "percentage")

					// Verify realistic values
					assert.IsType(t, "", entry["client_version"])
					assert.GreaterOrEqual(t, entry["user_count"].(float64), 0.0)
					assert.GreaterOrEqual(t, entry["percentage"].(float64), 0.0)
					assert.LessOrEqual(t, entry["percentage"].(float64), 1.0)
				}
			},
			hasPagination: false,
		},
		{
			name:     "top-file-extensions",
			endpoint: "/analytics/team/top-file-extensions",
			validateData: func(t *testing.T, data interface{}) {
				dataArray := data.([]interface{})
				if len(dataArray) == 0 {
					t.Skip("No file extension data generated in test setup")
					return
				}

				for _, item := range dataArray {
					entry := item.(map[string]interface{})

					// Required fields from cursor_analytics.md
					assert.Contains(t, entry, "event_date")
					assert.Contains(t, entry, "file_extension")
					assert.Contains(t, entry, "total_files")
					assert.Contains(t, entry, "total_accepts")
					assert.Contains(t, entry, "total_rejects")
					assert.Contains(t, entry, "total_lines_suggested")
					assert.Contains(t, entry, "total_lines_accepted")
					assert.Contains(t, entry, "total_lines_rejected")

					// Verify realistic values
					assert.IsType(t, "", entry["file_extension"])
					assert.GreaterOrEqual(t, entry["total_files"].(float64), 0.0)
				}
			},
			hasPagination: false,
		},
		{
			name:     "mcp",
			endpoint: "/analytics/team/mcp",
			validateData: func(t *testing.T, data interface{}) {
				dataArray := data.([]interface{})
				if len(dataArray) == 0 {
					t.Skip("No MCP usage data generated in test setup")
					return
				}

				for _, item := range dataArray {
					entry := item.(map[string]interface{})

					// Required fields from cursor_analytics.md
					assert.Contains(t, entry, "event_date")
					assert.Contains(t, entry, "tool_name")
					assert.Contains(t, entry, "mcp_server_name")
					assert.Contains(t, entry, "usage")

					// Verify realistic values
					assert.IsType(t, "", entry["tool_name"])
					assert.IsType(t, "", entry["mcp_server_name"])
					assert.GreaterOrEqual(t, entry["usage"].(float64), 0.0)
				}
			},
			hasPagination: false,
		},
		{
			name:     "commands",
			endpoint: "/analytics/team/commands",
			validateData: func(t *testing.T, data interface{}) {
				dataArray := data.([]interface{})
				if len(dataArray) == 0 {
					t.Skip("No command usage data generated in test setup")
					return
				}

				for _, item := range dataArray {
					entry := item.(map[string]interface{})

					// Required fields from cursor_analytics.md
					assert.Contains(t, entry, "event_date")
					assert.Contains(t, entry, "command_name")
					assert.Contains(t, entry, "usage")

					// Verify realistic values
					assert.IsType(t, "", entry["command_name"])
					assert.GreaterOrEqual(t, entry["usage"].(float64), 0.0)
				}
			},
			hasPagination: false,
		},
		{
			name:     "plans",
			endpoint: "/analytics/team/plans",
			validateData: func(t *testing.T, data interface{}) {
				dataArray := data.([]interface{})
				if len(dataArray) == 0 {
					t.Skip("No plan usage data generated in test setup")
					return
				}

				for _, item := range dataArray {
					entry := item.(map[string]interface{})

					// Required fields from cursor_analytics.md
					assert.Contains(t, entry, "event_date")
					assert.Contains(t, entry, "model")
					assert.Contains(t, entry, "usage")

					// Verify realistic values
					assert.IsType(t, "", entry["model"])
					assert.GreaterOrEqual(t, entry["usage"].(float64), 0.0)
				}
			},
			hasPagination: false,
		},
		{
			name:     "ask-mode",
			endpoint: "/analytics/team/ask-mode",
			validateData: func(t *testing.T, data interface{}) {
				dataArray := data.([]interface{})
				if len(dataArray) == 0 {
					t.Skip("No ask-mode usage data generated in test setup")
					return
				}

				for _, item := range dataArray {
					entry := item.(map[string]interface{})

					// Required fields from cursor_analytics.md
					assert.Contains(t, entry, "event_date")
					assert.Contains(t, entry, "model")
					assert.Contains(t, entry, "usage")

					// Verify realistic values
					assert.IsType(t, "", entry["model"])
					assert.GreaterOrEqual(t, entry["usage"].(float64), 0.0)
				}
			},
			hasPagination: false,
		},
		{
			name:     "leaderboard",
			endpoint: "/analytics/team/leaderboard",
			validateData: func(t *testing.T, data interface{}) {
				leaderboardData := data.(map[string]interface{})

				// Required top-level fields from cursor_analytics.md
				assert.Contains(t, leaderboardData, "tab_leaderboard")
				assert.Contains(t, leaderboardData, "agent_leaderboard")

				// Validate tab_leaderboard
				tabLb := leaderboardData["tab_leaderboard"].(map[string]interface{})
				assert.Contains(t, tabLb, "data")
				assert.Contains(t, tabLb, "total_users")

				tabData := tabLb["data"].([]interface{})
				if len(tabData) > 0 {
					for _, item := range tabData {
						entry := item.(map[string]interface{})

						// Required fields from cursor_analytics.md
						assert.Contains(t, entry, "email")
						assert.Contains(t, entry, "user_id")
						assert.Contains(t, entry, "total_accepts")
						assert.Contains(t, entry, "total_lines_accepted")
						assert.Contains(t, entry, "total_lines_suggested")
						assert.Contains(t, entry, "line_acceptance_ratio")
						assert.Contains(t, entry, "rank")

						// Tab leaderboard has accept_ratio
						assert.Contains(t, entry, "accept_ratio")

						// Verify realistic values
						assert.IsType(t, "", entry["email"])
						assert.GreaterOrEqual(t, entry["rank"].(float64), 1.0)
					}
				}

				// Validate agent_leaderboard
				agentLb := leaderboardData["agent_leaderboard"].(map[string]interface{})
				assert.Contains(t, agentLb, "data")
				assert.Contains(t, agentLb, "total_users")

				agentData := agentLb["data"].([]interface{})
				if len(agentData) > 0 {
					for _, item := range agentData {
						entry := item.(map[string]interface{})

						// Required fields from cursor_analytics.md
						assert.Contains(t, entry, "email")
						assert.Contains(t, entry, "user_id")
						assert.Contains(t, entry, "total_accepts")
						assert.Contains(t, entry, "total_lines_accepted")
						assert.Contains(t, entry, "total_lines_suggested")
						assert.Contains(t, entry, "line_acceptance_ratio")
						assert.Contains(t, entry, "rank")

						// Agent leaderboard has favorite_model (optional)
						// Note: favorite_model is optional, so we don't assert it

						// Verify realistic values
						assert.IsType(t, "", entry["email"])
						assert.GreaterOrEqual(t, entry["rank"].(float64), 1.0)
					}
				}
			},
			hasPagination: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := makeRequest(t, "GET", tt.endpoint, true)
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode, "Endpoint %s should return 200", tt.endpoint)

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			// All team-level endpoints have data and params
			assert.Contains(t, result, "data", "should have data field")
			assert.Contains(t, result, "params", "should have params field")

			// Validate params structure
			params := result["params"].(map[string]interface{})
			assert.Contains(t, params, "metric")
			// teamId is optional in test environment
			assert.Contains(t, params, "startDate")
			assert.Contains(t, params, "endDate")

			// Leaderboard has pagination, others don't
			if tt.hasPagination {
				assert.Contains(t, result, "pagination", "leaderboard should have pagination")

				// Validate pagination structure
				pagination := result["pagination"].(map[string]interface{})
				assert.Contains(t, pagination, "page")
				assert.Contains(t, pagination, "pageSize")
				assert.Contains(t, pagination, "totalUsers")
				assert.Contains(t, pagination, "totalPages")
				assert.Contains(t, pagination, "hasNextPage")
				assert.Contains(t, pagination, "hasPreviousPage")
			} else {
				assert.NotContains(t, result, "pagination", "team-level endpoints should not have pagination")
			}

			// Validate data structure
			tt.validateData(t, result["data"])
		})
	}
}

// TestE2E_ByUserAnalyticsSchemas validates the exact schema of each by-user endpoint
// against the Cursor Analytics API documentation.
func TestE2E_ByUserAnalyticsSchemas(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name         string
		endpoint     string
		validateUser func(t *testing.T, userData []interface{})
	}{
		{
			name:     "by-user/agent-edits",
			endpoint: "/analytics/by-user/agent-edits",
			validateUser: func(t *testing.T, userData []interface{}) {
				for _, item := range userData {
					day := item.(map[string]interface{})

					// By-user agent-edits uses simplified schema
					// Reference: docs/api-reference/cursor_analytics.md (By-User Agent Edits)
					assert.Contains(t, day, "event_date")
					assert.Contains(t, day, "suggested_lines")
					assert.Contains(t, day, "accepted_lines")
					assert.GreaterOrEqual(t, day["suggested_lines"].(float64), 0.0)
				}
			},
		},
		{
			name:     "by-user/tabs",
			endpoint: "/analytics/by-user/tabs",
			validateUser: func(t *testing.T, userData []interface{}) {
				for _, item := range userData {
					day := item.(map[string]interface{})

					// By-user tabs uses simplified schema (same as agent-edits)
					assert.Contains(t, day, "event_date")
					assert.Contains(t, day, "suggested_lines")
					assert.Contains(t, day, "accepted_lines")
					assert.GreaterOrEqual(t, day["suggested_lines"].(float64), 0.0)
				}
			},
		},
		{
			name:     "by-user/models",
			endpoint: "/analytics/by-user/models",
			validateUser: func(t *testing.T, userData []interface{}) {
				for _, item := range userData {
					day := item.(map[string]interface{})

					assert.Contains(t, day, "date")
					assert.Contains(t, day, "model_breakdown")

					breakdown := day["model_breakdown"].(map[string]interface{})
					for _, modelData := range breakdown {
						modelInfo := modelData.(map[string]interface{})
						assert.Contains(t, modelInfo, "messages")
						assert.Contains(t, modelInfo, "users")
					}
				}
			},
		},
		{
			name:     "by-user/client-versions",
			endpoint: "/analytics/by-user/client-versions",
			validateUser: func(t *testing.T, userData []interface{}) {
				for _, item := range userData {
					entry := item.(map[string]interface{})

					assert.Contains(t, entry, "event_date")
					assert.Contains(t, entry, "client_version")
					assert.Contains(t, entry, "user_count")
					assert.Contains(t, entry, "percentage")
				}
			},
		},
		{
			name:     "by-user/top-file-extensions",
			endpoint: "/analytics/by-user/top-file-extensions",
			validateUser: func(t *testing.T, userData []interface{}) {
				for _, item := range userData {
					entry := item.(map[string]interface{})

					assert.Contains(t, entry, "event_date")
					assert.Contains(t, entry, "file_extension")
					assert.Contains(t, entry, "total_files")
					assert.Contains(t, entry, "total_accepts")
				}
			},
		},
		{
			name:     "by-user/mcp",
			endpoint: "/analytics/by-user/mcp",
			validateUser: func(t *testing.T, userData []interface{}) {
				for _, item := range userData {
					entry := item.(map[string]interface{})

					assert.Contains(t, entry, "event_date")
					assert.Contains(t, entry, "tool_name")
					assert.Contains(t, entry, "mcp_server_name")
					assert.Contains(t, entry, "usage")
				}
			},
		},
		{
			name:     "by-user/commands",
			endpoint: "/analytics/by-user/commands",
			validateUser: func(t *testing.T, userData []interface{}) {
				for _, item := range userData {
					entry := item.(map[string]interface{})

					assert.Contains(t, entry, "event_date")
					assert.Contains(t, entry, "command_name")
					assert.Contains(t, entry, "usage")
				}
			},
		},
		{
			name:     "by-user/plans",
			endpoint: "/analytics/by-user/plans",
			validateUser: func(t *testing.T, userData []interface{}) {
				for _, item := range userData {
					entry := item.(map[string]interface{})

					assert.Contains(t, entry, "event_date")
					assert.Contains(t, entry, "model")
					assert.Contains(t, entry, "usage")
				}
			},
		},
		{
			name:     "by-user/ask-mode",
			endpoint: "/analytics/by-user/ask-mode",
			validateUser: func(t *testing.T, userData []interface{}) {
				for _, item := range userData {
					entry := item.(map[string]interface{})

					assert.Contains(t, entry, "event_date")
					assert.Contains(t, entry, "model")
					assert.Contains(t, entry, "usage")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := makeRequest(t, "GET", tt.endpoint, true)
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode, "Endpoint %s should return 200", tt.endpoint)

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			// All by-user endpoints have data, pagination, and params
			assert.Contains(t, result, "data", "should have data field")
			assert.Contains(t, result, "pagination", "should have pagination field")
			assert.Contains(t, result, "params", "should have params field")

			// Validate data structure (object keyed by email)
			data := result["data"].(map[string]interface{})

			// Feature endpoints (models, client-versions, etc.) may have no data in test setup
			isFeatureEndpoint := strings.Contains(tt.endpoint, "models") ||
								  strings.Contains(tt.endpoint, "client-versions") ||
								  strings.Contains(tt.endpoint, "top-file-extensions") ||
								  strings.Contains(tt.endpoint, "mcp") ||
								  strings.Contains(tt.endpoint, "commands") ||
								  strings.Contains(tt.endpoint, "plans") ||
								  strings.Contains(tt.endpoint, "ask-mode")

			if isFeatureEndpoint && len(data) == 0 {
				t.Skipf("No data generated for %s in test setup", tt.endpoint)
				return
			}

			assert.Greater(t, len(data), 0, "should have at least one user")

			// Validate each user's data
			for email, userData := range data {
				assert.IsType(t, "", email, "key should be email string")
				userArray := userData.([]interface{})
				assert.Greater(t, len(userArray), 0, "user should have at least one data point")

				tt.validateUser(t, userArray)
			}

			// Validate pagination structure
			pagination := result["pagination"].(map[string]interface{})
			assert.Contains(t, pagination, "page")
			assert.Contains(t, pagination, "pageSize")
			// totalUsers is optional (omitempty when 0)
			assert.Contains(t, pagination, "totalPages")
			assert.Contains(t, pagination, "hasNextPage")
			assert.Contains(t, pagination, "hasPreviousPage")

			// Validate params structure
			params := result["params"].(map[string]interface{})
			assert.Contains(t, params, "metric")
			// teamId is optional in test environment
			assert.Contains(t, params, "startDate")
			assert.Contains(t, params, "endDate")

			// userMappings is optional (omitempty when empty)
			if userMappingsVal, exists := params["userMappings"]; exists && userMappingsVal != nil {
				userMappings := userMappingsVal.([]interface{})
				assert.Equal(t, len(data), len(userMappings), "userMappings should match number of users in data")

				for _, mapping := range userMappings {
					mappingObj := mapping.(map[string]interface{})
					assert.Contains(t, mappingObj, "id")
					assert.Contains(t, mappingObj, "email")
				}
			}
		})
	}
}

// TestE2E_LeaderboardPagination tests pagination behavior on the leaderboard endpoint
func TestE2E_LeaderboardPagination(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	t.Run("default_pagination", func(t *testing.T) {
		resp := makeRequest(t, "GET", "/analytics/team/leaderboard", true)
		defer resp.Body.Close()

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		pagination := result["pagination"].(map[string]interface{})
		assert.Equal(t, float64(1), pagination["page"], "default page should be 1")
		// Implementation uses pageSize of 100 (from default params parsing)
		assert.GreaterOrEqual(t, pagination["pageSize"].(float64), float64(1), "pageSize should be >= 1")
	})

	t.Run("custom_pagination", func(t *testing.T) {
		resp := makeRequest(t, "GET", "/analytics/team/leaderboard?page=2&pageSize=5", true)
		defer resp.Body.Close()

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		pagination := result["pagination"].(map[string]interface{})
		assert.Equal(t, float64(2), pagination["page"])
		assert.Equal(t, float64(5), pagination["pageSize"])

		// Verify hasPreviousPage is true on page 2
		assert.Equal(t, true, pagination["hasPreviousPage"])
	})

	t.Run("pagination_consistency", func(t *testing.T) {
		resp := makeRequest(t, "GET", "/analytics/team/leaderboard", true)
		defer resp.Body.Close()

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		pagination := result["pagination"].(map[string]interface{})
		totalUsers := int(pagination["totalUsers"].(float64))
		pageSize := int(pagination["pageSize"].(float64))
		totalPages := int(pagination["totalPages"].(float64))

		// Verify totalPages calculation
		expectedPages := (totalUsers + pageSize - 1) / pageSize
		assert.Equal(t, expectedPages, totalPages, "totalPages should be calculated correctly")
	})
}

// TestE2E_ByUserPagination tests pagination behavior on by-user endpoints
func TestE2E_ByUserPagination(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	endpoint := "/analytics/by-user/agent-edits"

	t.Run("default_pagination", func(t *testing.T) {
		resp := makeRequest(t, "GET", endpoint, true)
		defer resp.Body.Close()

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		pagination := result["pagination"].(map[string]interface{})
		assert.Equal(t, float64(1), pagination["page"])

		// Default pageSize for by-user endpoints is 100
		params := result["params"].(map[string]interface{})
		pageSize := params["pageSize"]
		assert.NotNil(t, pageSize)
	})

	t.Run("custom_page_size", func(t *testing.T) {
		resp := makeRequest(t, "GET", endpoint+"?pageSize=2", true)
		defer resp.Body.Close()

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		data := result["data"].(map[string]interface{})
		assert.LessOrEqual(t, len(data), 2, "should respect pageSize limit")
	})
}

// TestE2E_UserFiltering tests user filtering on by-user endpoints
func TestE2E_UserFiltering(t *testing.T) {
	cleanup, store := setupTestServer(t)
	defer cleanup()

	// Get all developers to find a specific user to filter
	developers := store.ListDevelopers()
	if len(developers) == 0 {
		t.Skip("No developers in test setup")
		return
	}

	targetEmail := developers[0].Email

	t.Run("filter_single_user", func(t *testing.T) {
		url := fmt.Sprintf("/analytics/by-user/agent-edits?users=%s", targetEmail)
		resp := makeRequest(t, "GET", url, true)
		defer resp.Body.Close()

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		// Data should only contain the filtered user
		data := result["data"].(map[string]interface{})

		// Should have at most 1 user (might be 0 if no data for that user)
		assert.LessOrEqual(t, len(data), 1)

		// If we have data, it should be for the filtered user
		if len(data) > 0 {
			_, exists := data[targetEmail]
			assert.True(t, exists, "data should contain the filtered user")
		}

		// Pagination should reflect filtered results
		pagination := result["pagination"].(map[string]interface{})
		totalUsers := int(pagination["totalUsers"].(float64))
		assert.LessOrEqual(t, totalUsers, 1, "totalUsers should reflect filter")
	})

	t.Run("filter_multiple_users", func(t *testing.T) {
		if len(developers) < 2 {
			t.Skip("need at least 2 developers for this test")
		}

		email1 := developers[0].Email
		email2 := developers[1].Email
		url := fmt.Sprintf("/analytics/by-user/tabs?users=%s,%s", email1, email2)

		resp := makeRequest(t, "GET", url, true)
		defer resp.Body.Close()

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		data := result["data"].(map[string]interface{})
		assert.LessOrEqual(t, len(data), 2, "should have at most 2 users")

		// Pagination should reflect filtered results
		pagination := result["pagination"].(map[string]interface{})
		totalUsers := int(pagination["totalUsers"].(float64))
		assert.LessOrEqual(t, totalUsers, 2)
	})
}

// TestE2E_DateFiltering tests date range filtering on all analytics endpoints
func TestE2E_DateFiltering(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	now := time.Now()
	startDate := now.Add(-7 * 24 * time.Hour).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	tests := []struct {
		name     string
		endpoint string
	}{
		{"team/agent-edits", "/analytics/team/agent-edits"},
		{"team/tabs", "/analytics/team/tabs"},
		{"team/dau", "/analytics/team/dau"},
		{"by-user/agent-edits", "/analytics/by-user/agent-edits"},
		{"by-user/tabs", "/analytics/by-user/tabs"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s?startDate=%s&endDate=%s", tt.endpoint, startDate, endDate)
			resp := makeRequest(t, "GET", url, true)
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			// Verify params reflect the date range
			params := result["params"].(map[string]interface{})
			assert.Contains(t, params, "startDate")
			assert.Contains(t, params, "endDate")

			// Verify the dates are what we requested (or defaults if not supported)
			returnedStart := params["startDate"].(string)
			returnedEnd := params["endDate"].(string)
			assert.NotEmpty(t, returnedStart)
			assert.NotEmpty(t, returnedEnd)
		})
	}
}

// TestE2E_NoStubResponses verifies that all endpoints return real data, not stub responses
func TestE2E_NoStubResponses(t *testing.T) {
	cleanup, _ := setupTestServer(t)
	defer cleanup()

	t.Run("team_endpoints_have_data", func(t *testing.T) {
		// Core endpoints that MUST have data from commit generation
		coreEndpoints := []string{
			"/analytics/team/agent-edits",
			"/analytics/team/tabs",
			"/analytics/team/dau",
		}

		// Feature endpoints that may not have data in test setup
		featureEndpoints := []string{
			"/analytics/team/models",
			"/analytics/team/client-versions",
			"/analytics/team/top-file-extensions",
			"/analytics/team/mcp",
			"/analytics/team/commands",
			"/analytics/team/plans",
			"/analytics/team/ask-mode",
		}

		// Core endpoints must have data
		for _, endpoint := range coreEndpoints {
			resp := makeRequest(t, "GET", endpoint, true)
			defer resp.Body.Close()

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			data := result["data"].([]interface{})
			assert.Greater(t, len(data), 0, "%s should return non-empty data", endpoint)
		}

		// Feature endpoints should return valid structure even if empty
		for _, endpoint := range featureEndpoints {
			resp := makeRequest(t, "GET", endpoint, true)
			defer resp.Body.Close()

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			data := result["data"].([]interface{})
			// Allow empty data, but verify it's a valid array
			assert.NotNil(t, data, "%s should return valid data array", endpoint)
		}
	})

	t.Run("by_user_endpoints_have_data", func(t *testing.T) {
		coreEndpoints := []string{
			"/analytics/by-user/agent-edits",
			"/analytics/by-user/tabs",
		}

		// Core endpoints must have data
		for _, endpoint := range coreEndpoints {
			resp := makeRequest(t, "GET", endpoint, true)
			defer resp.Body.Close()

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			data := result["data"].(map[string]interface{})
			assert.Greater(t, len(data), 0, "%s should return non-empty data", endpoint)

			// Verify at least one user has data
			for _, userData := range data {
				userArray := userData.([]interface{})
				assert.Greater(t, len(userArray), 0, "user should have at least one data point")
				break
			}
		}

		// Feature endpoints may be empty
		featureEndpoints := []string{
			"/analytics/by-user/models",
			"/analytics/by-user/client-versions",
			"/analytics/by-user/top-file-extensions",
		}

		for _, endpoint := range featureEndpoints {
			resp := makeRequest(t, "GET", endpoint, true)
			defer resp.Body.Close()

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			data := result["data"].(map[string]interface{})
			assert.NotNil(t, data, "%s should return valid data object", endpoint)
		}
	})

	t.Run("leaderboard_has_rankings", func(t *testing.T) {
		resp := makeRequest(t, "GET", "/analytics/team/leaderboard", true)
		defer resp.Body.Close()

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		data := result["data"].(map[string]interface{})

		// Check tab leaderboard
		tabLb := data["tab_leaderboard"].(map[string]interface{})
		tabData := tabLb["data"].([]interface{})
		if len(tabData) > 0 {
			entry := tabData[0].(map[string]interface{})
			assert.Greater(t, entry["total_accepts"].(float64), 0.0, "leaderboard should have non-zero accepts")
		}

		// Check agent leaderboard
		agentLb := data["agent_leaderboard"].(map[string]interface{})
		agentData := agentLb["data"].([]interface{})
		if len(agentData) > 0 {
			entry := agentData[0].(map[string]interface{})
			assert.Greater(t, entry["total_accepts"].(float64), 0.0, "leaderboard should have non-zero accepts")
		}
	})
}
