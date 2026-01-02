package models

import "time"

// ===========================================================================
// Team-Level Analytics Response Models (matches cursor_analytics.md)
// ===========================================================================

// AgentEditsDay represents daily agent edit metrics.
// Matches Cursor Analytics API /analytics/team/agent-edits
// Reference: docs/api-reference/cursor_analytics.md
type AgentEditsDay struct {
	EventDate                string `json:"event_date"`
	TotalSuggestedDiffs      int    `json:"total_suggested_diffs"`
	TotalAcceptedDiffs       int    `json:"total_accepted_diffs"`
	TotalRejectedDiffs       int    `json:"total_rejected_diffs"`
	TotalGreenLinesAccepted  int    `json:"total_green_lines_accepted"`
	TotalRedLinesAccepted    int    `json:"total_red_lines_accepted"`
	TotalGreenLinesRejected  int    `json:"total_green_lines_rejected"`
	TotalRedLinesRejected    int    `json:"total_red_lines_rejected"`
	TotalGreenLinesSuggested int    `json:"total_green_lines_suggested"`
	TotalRedLinesSuggested   int    `json:"total_red_lines_suggested"`
	TotalLinesSuggested      int    `json:"total_lines_suggested"`
	TotalLinesAccepted       int    `json:"total_lines_accepted"`
}

// AcceptanceRate calculates the ratio of accepted diffs to suggested diffs.
// Returns a value between 0.0 and 1.0.
func (a *AgentEditsDay) AcceptanceRate() float64 {
	if a.TotalSuggestedDiffs == 0 {
		return 0.0
	}
	return float64(a.TotalAcceptedDiffs) / float64(a.TotalSuggestedDiffs)
}

// TabUsageDay represents daily tab completion metrics.
// Matches Cursor Analytics API /analytics/team/tabs
// Reference: docs/api-reference/cursor_analytics.md
type TabUsageDay struct {
	EventDate                string `json:"event_date"`
	TotalSuggestions         int    `json:"total_suggestions"`
	TotalAccepts             int    `json:"total_accepts"`
	TotalRejects             int    `json:"total_rejects"`
	TotalGreenLinesAccepted  int    `json:"total_green_lines_accepted"`
	TotalRedLinesAccepted    int    `json:"total_red_lines_accepted"`
	TotalGreenLinesRejected  int    `json:"total_green_lines_rejected"`
	TotalRedLinesRejected    int    `json:"total_red_lines_rejected"`
	TotalGreenLinesSuggested int    `json:"total_green_lines_suggested"`
	TotalRedLinesSuggested   int    `json:"total_red_lines_suggested"`
	TotalLinesSuggested      int    `json:"total_lines_suggested"`
	TotalLinesAccepted       int    `json:"total_lines_accepted"`
}

// DAUDay represents daily active users metrics.
// Matches Cursor Analytics API /analytics/team/dau
// Reference: docs/api-reference/cursor_analytics.md
// NOTE: Field is "date" not "event_date" for this endpoint!
type DAUDay struct {
	Date          string `json:"date"`
	DAU           int    `json:"dau"`
	CLIDAU        int    `json:"cli_dau"`
	CloudAgentDAU int    `json:"cloud_agent_dau"`
	BugbotDAU     int    `json:"bugbot_dau"`
}

// ModelUsageDay represents daily model usage metrics.
// Matches Cursor Analytics API /analytics/team/models
// Reference: docs/api-reference/cursor_analytics.md
type ModelUsageDay struct {
	Date           string                        `json:"date"`
	ModelBreakdown map[string]ModelBreakdownItem `json:"model_breakdown"`
}

// ModelBreakdownItem represents usage stats for a single model.
type ModelBreakdownItem struct {
	Messages int `json:"messages"`
	Users    int `json:"users"`
}

// ClientVersionDay represents daily client version distribution.
// Matches Cursor Analytics API /analytics/team/client-versions
// Reference: docs/api-reference/cursor_analytics.md
type ClientVersionDay struct {
	EventDate     string  `json:"event_date"`
	ClientVersion string  `json:"client_version"`
	UserCount     int     `json:"user_count"`
	Percentage    float64 `json:"percentage"`
}

// FileExtensionDay represents daily file extension usage.
// Matches Cursor Analytics API /analytics/team/top-file-extensions
// Reference: docs/api-reference/cursor_analytics.md
type FileExtensionDay struct {
	EventDate           string `json:"event_date"`
	FileExtension       string `json:"file_extension"`
	TotalFiles          int    `json:"total_files"`
	TotalAccepts        int    `json:"total_accepts"`
	TotalRejects        int    `json:"total_rejects"`
	TotalLinesSuggested int    `json:"total_lines_suggested"`
	TotalLinesAccepted  int    `json:"total_lines_accepted"`
	TotalLinesRejected  int    `json:"total_lines_rejected"`
}

// MCPUsageDay represents daily MCP (Model Context Protocol) tool usage.
// Matches Cursor Analytics API /analytics/team/mcp
// Reference: docs/api-reference/cursor_analytics.md
type MCPUsageDay struct {
	EventDate     string `json:"event_date"`
	ToolName      string `json:"tool_name"`
	MCPServerName string `json:"mcp_server_name"`
	Usage         int    `json:"usage"`
}

// CommandUsageDay represents daily Cursor command usage.
// Matches Cursor Analytics API /analytics/team/commands
// Reference: docs/api-reference/cursor_analytics.md
type CommandUsageDay struct {
	EventDate   string `json:"event_date"`
	CommandName string `json:"command_name"`
	Usage       int    `json:"usage"`
}

// PlanUsageDay represents daily AI planning feature usage.
// Matches Cursor Analytics API /analytics/team/plans
// Reference: docs/api-reference/cursor_analytics.md
type PlanUsageDay struct {
	EventDate string `json:"event_date"`
	Model     string `json:"model"`
	Usage     int    `json:"usage"`
}

// AskModeDay represents daily ask mode usage.
// Matches Cursor Analytics API /analytics/team/ask-mode
// Reference: docs/api-reference/cursor_analytics.md
type AskModeDay struct {
	EventDate string `json:"event_date"`
	Model     string `json:"model"`
	Usage     int    `json:"usage"`
}

// LeaderboardResponse represents the team leaderboard response.
// Matches Cursor Analytics API /analytics/team/leaderboard
// Reference: docs/api-reference/cursor_analytics.md
type LeaderboardResponse struct {
	TabLeaderboard   LeaderboardSection `json:"tab_leaderboard"`
	AgentLeaderboard LeaderboardSection `json:"agent_leaderboard"`
}

// LeaderboardSection represents one section of the leaderboard.
type LeaderboardSection struct {
	Data       []LeaderboardEntry `json:"data"`
	TotalUsers int                `json:"total_users"`
}

// LeaderboardEntry represents a single user's leaderboard entry.
type LeaderboardEntry struct {
	Email               string  `json:"email"`
	UserID              string  `json:"user_id"`
	TotalAccepts        int     `json:"total_accepts"`
	TotalLinesAccepted  int     `json:"total_lines_accepted"`
	TotalLinesSuggested int     `json:"total_lines_suggested"`
	LineAcceptanceRatio float64 `json:"line_acceptance_ratio"`
	AcceptRatio         float64 `json:"accept_ratio,omitempty"`
	FavoriteModel       string  `json:"favorite_model,omitempty"`
	Rank                int     `json:"rank"`
}

// ===========================================================================
// Internal Event Models (used for data generation)
// ===========================================================================

// ModelUsageEvent represents an individual model usage event.
type ModelUsageEvent struct {
	UserID    string    `json:"user_id"`
	UserEmail string    `json:"user_email"`
	ModelName string    `json:"model_name"`
	UsageType string    `json:"usage_type"` // "chat" or "code"
	Timestamp time.Time `json:"timestamp"`
	EventDate string    `json:"event_date"` // YYYY-MM-DD format
}

// ClientVersionEvent represents a client version usage event.
type ClientVersionEvent struct {
	UserID        string    `json:"user_id"`
	UserEmail     string    `json:"user_email"`
	ClientVersion string    `json:"client_version"` // Semver format (e.g., "0.42.3")
	Timestamp     time.Time `json:"timestamp"`
	EventDate     string    `json:"event_date"` // YYYY-MM-DD format
}

// ===========================================================================
// Legacy Type Aliases (for backwards compatibility during migration)
// ===========================================================================

// TabCompletionDay is deprecated. Use TabUsageDay instead.
// Kept for backwards compatibility with existing code.
type TabCompletionDay = TabUsageDay
