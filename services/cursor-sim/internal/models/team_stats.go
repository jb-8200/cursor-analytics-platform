package models

// AgentEditsDay represents daily agent edit metrics.
// Field names use snake_case to match the Cursor Team Analytics API.
type AgentEditsDay struct {
	EventDate               string `json:"event_date"`
	TotalSuggestedDiffs     int    `json:"total_suggested_diffs"`
	TotalAcceptedDiffs      int    `json:"total_accepted_diffs"`
	TotalRejectedDiffs      int    `json:"total_rejected_diffs"`
	TotalGreenLinesAccepted int    `json:"total_green_lines_accepted"`
	TotalRedLinesAccepted   int    `json:"total_red_lines_accepted"`
}

// AcceptanceRate calculates the ratio of accepted diffs to suggested diffs.
// Returns a value between 0.0 and 1.0.
func (a *AgentEditsDay) AcceptanceRate() float64 {
	if a.TotalSuggestedDiffs == 0 {
		return 0.0
	}
	return float64(a.TotalAcceptedDiffs) / float64(a.TotalSuggestedDiffs)
}

// TabCompletionDay represents daily tab completion metrics.
type TabCompletionDay struct {
	EventDate     string `json:"event_date"`
	TotalSuggests int    `json:"total_suggests"`
	TotalAccepts  int    `json:"total_accepts"`
	TotalRejects  int    `json:"total_rejects"`
}

// DAUDay represents daily active users metrics.
type DAUDay struct {
	EventDate        string  `json:"event_date"`
	UniqueUsers      int     `json:"unique_users"`
	TotalEvents      int     `json:"total_events"`
	AvgEventsPerUser float64 `json:"avg_events_per_user"`
}

// ModelUsage represents model usage statistics.
type ModelUsage struct {
	ModelName   string `json:"model_name"`
	TotalUsages int    `json:"total_usages"`
	UniqueUsers int    `json:"unique_users"`
}

// FileExtensionStats represents file extension usage statistics.
type FileExtensionStats struct {
	Extension   string `json:"extension"`
	TotalEdits  int    `json:"total_edits"`
	UniqueUsers int    `json:"unique_users"`
}

// ClientVersion represents client version distribution.
type ClientVersion struct {
	Version     string `json:"version"`
	TotalUsers  int    `json:"total_users"`
	UniqueUsers int    `json:"unique_users"`
}
