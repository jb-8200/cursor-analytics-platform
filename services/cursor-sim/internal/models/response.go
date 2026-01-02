package models

// ===========================================================================
// AI Code Tracking API Response Types (matches OpenAPI spec)
// ===========================================================================

// CommitsResponse is the response for GET /analytics/ai-code/commits.
// Matches the CommitsResponse schema in cursor-api.yaml.
type CommitsResponse struct {
	Items      []Commit `json:"items"`
	TotalCount int      `json:"totalCount"`
	Page       int      `json:"page"`
	PageSize   int      `json:"pageSize"`
}

// ChangesResponse is the response for GET /analytics/ai-code/changes.
// Matches the ChangesResponse schema in cursor-api.yaml.
type ChangesResponse struct {
	Items      []Change `json:"items"`
	TotalCount int      `json:"totalCount"`
	Page       int      `json:"page"`
	PageSize   int      `json:"pageSize"`
}

// ===========================================================================
// Admin API Response Types (matches OpenAPI spec)
// ===========================================================================

// DailyUsageResponse is the response for POST /teams/daily-usage-data.
type DailyUsageResponse struct {
	Data []DailyUsageRecord `json:"data"`
}

// DailyUsageRecord represents a single day's usage data.
type DailyUsageRecord struct {
	Date   string       `json:"date"`
	Org    string       `json:"org"`
	Totals UsageTotals  `json:"totals"`
	ByUser []UserUsage  `json:"by_user,omitempty"`
}

// UsageTotals contains aggregated usage metrics.
type UsageTotals struct {
	TotalAccepts      int `json:"total_accepts"`
	TotalRejects      int `json:"total_rejects"`
	TotalTabsShown    int `json:"total_tabs_shown"`
	TotalTabsAccepted int `json:"total_tabs_accepted"`
	TotalLinesAdded   int `json:"total_lines_added"`
	TotalLinesDeleted int `json:"total_lines_deleted"`
}

// UserUsage contains per-user usage for a day.
type UserUsage struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Accepts      int    `json:"accepts"`
	Rejects      int    `json:"rejects"`
	TabsShown    int    `json:"tabs_shown"`
	TabsAccepted int    `json:"tabs_accepted"`
	LinesAdded   int    `json:"lines_added"`
	LinesDeleted int    `json:"lines_deleted"`
}

// FilteredUsageEventsResponse is the response for POST /teams/filtered-usage-events.
type FilteredUsageEventsResponse struct {
	UsageEvents []UsageEvent `json:"usageEvents"`
	TotalCount  int          `json:"totalCount"`
	Page        int          `json:"page"`
	PageSize    int          `json:"pageSize"`
	TotalPages  int          `json:"totalPages"`
}

// UsageEvent represents a single usage event.
type UsageEvent struct {
	ID               string `json:"id"`
	UserEmail        string `json:"userEmail"`
	Org              string `json:"org,omitempty"`
	EventType        string `json:"eventType"` // "composer", "chat", "agent"
	Model            string `json:"model,omitempty"`
	TotalAccepts     int    `json:"totalAccepts,omitempty"`
	TotalRejects     int    `json:"totalRejects,omitempty"`
	InputTokens      int    `json:"inputTokens,omitempty"`
	OutputTokens     int    `json:"outputTokens,omitempty"`
	CacheWriteTokens int    `json:"cacheWriteTokens,omitempty"`
	CacheReadTokens  int    `json:"cacheReadTokens,omitempty"`
	TotalCents       int    `json:"totalCents,omitempty"`
	CreatedAt        string `json:"createdAt"`
}

// SpendResponse is the response for POST /teams/spend.
type SpendResponse struct {
	TeamMemberSpend        []MemberSpend `json:"teamMemberSpend"`
	SubscriptionCycleStart int64         `json:"subscriptionCycleStart"`
	TotalMembers           int           `json:"totalMembers"`
	TotalPages             int           `json:"totalPages"`
}

// MemberSpend represents spending for a single team member.
type MemberSpend struct {
	SpendCents              int    `json:"spendCents"`
	FastPremiumRequests     int    `json:"fastPremiumRequests"`
	Name                    string `json:"name"`
	Email                   string `json:"email"`
	Role                    string `json:"role"`
	HardLimitOverrideDollars int   `json:"hardLimitOverrideDollars"`
}

// ===========================================================================
// Analytics API Response Types (matches cursor_analytics.md)
// ===========================================================================

// AnalyticsTeamResponse is the response format for team-level analytics endpoints.
// Used by: /analytics/team/agent-edits, /analytics/team/tabs, /analytics/team/dau,
// /analytics/team/models, /analytics/team/client-versions, /analytics/team/top-file-extensions,
// /analytics/team/mcp, /analytics/team/commands, /analytics/team/plans, /analytics/team/ask-mode
//
// Reference: docs/api-reference/cursor_analytics.md (Team-Level Endpoints)
type AnalyticsTeamResponse struct {
	Data   interface{}      `json:"data"`
	Params AnalyticsParams  `json:"params"`
}

// AnalyticsParams contains the request parameters for analytics endpoints.
type AnalyticsParams struct {
	Metric    string `json:"metric"`
	TeamID    int    `json:"teamId,omitempty"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Users     string `json:"users,omitempty"`
	Page      int    `json:"page,omitempty"`
	PageSize  int    `json:"pageSize,omitempty"`
}

// ===========================================================================
// Common Types
// ===========================================================================

// Error represents an API error response.
type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ===========================================================================
// Legacy Types (for backwards compatibility during migration)
// ===========================================================================

// PaginatedResponse wraps API responses with pagination metadata.
// Deprecated: Use endpoint-specific response types instead.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination,omitempty"`
	Params     Params      `json:"params,omitempty"`
}

// Pagination contains pagination metadata for list responses.
type Pagination struct {
	Page            int  `json:"page"`
	PageSize        int  `json:"pageSize"`
	TotalUsers      int  `json:"totalUsers,omitempty"`
	TotalPages      int  `json:"totalPages"`
	HasNextPage     bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
}

// Params contains the request parameters.
type Params struct {
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
	User      string `json:"user,omitempty"`
	Page      int    `json:"page,omitempty"`
	PageSize  int    `json:"pageSize,omitempty"`
	// Legacy fields - used internally during migration
	From     string `json:"-"`
	To       string `json:"-"`
	UserID   string `json:"-"`
	RepoName string `json:"-"`
}
