package cursor

import (
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
	"net/http"
)

// ByUserAgentEdits returns handler for GET /analytics/by-user/agent-edits.
// Returns agent edits broken down by user.
func ByUserAgentEdits(store storage.Store) http.Handler {
	return stubHandler("by-user-agent-edits")
}

// ByUserTabs returns handler for GET /analytics/by-user/tabs.
// Returns tab completion metrics by user.
func ByUserTabs(store storage.Store) http.Handler {
	return stubHandler("by-user-tabs")
}

// ByUserModels returns handler for GET /analytics/by-user/models.
func ByUserModels(store storage.Store) http.Handler {
	return stubHandler("by-user-models")
}

// ByUserClientVersions returns handler for GET /analytics/by-user/client-versions.
func ByUserClientVersions(store storage.Store) http.Handler {
	return stubHandler("by-user-client-versions")
}

// ByUserTopFileExtensions returns handler for GET /analytics/by-user/top-file-extensions.
func ByUserTopFileExtensions(store storage.Store) http.Handler {
	return stubHandler("by-user-top-file-extensions")
}

// ByUserMCP returns handler for GET /analytics/by-user/mcp.
func ByUserMCP(store storage.Store) http.Handler {
	return stubHandler("by-user-mcp")
}

// ByUserCommands returns handler for GET /analytics/by-user/commands.
func ByUserCommands(store storage.Store) http.Handler {
	return stubHandler("by-user-commands")
}

// ByUserPlans returns handler for GET /analytics/by-user/plans.
func ByUserPlans(store storage.Store) http.Handler {
	return stubHandler("by-user-plans")
}

// ByUserAskMode returns handler for GET /analytics/by-user/ask-mode.
func ByUserAskMode(store storage.Store) http.Handler {
	return stubHandler("by-user-ask-mode")
}
