package github

import (
	"net/http"
	"strings"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// RepoRouter returns an HTTP handler that routes GitHub API requests under /repos/.
// It handles all subroutes like /repos/{owner}/{repo}/pulls, /repos/{owner}/{repo}/commits, etc.
func RepoRouter(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Route to appropriate handler based on path pattern
		switch {
		// GET /repos/{owner}/{repo}
		case !strings.Contains(path, "/pulls") && !strings.Contains(path, "/commits") && countPathSegments(path) == 4:
			GetRepo(store).ServeHTTP(w, r)

		// GET /repos/{owner}/{repo}/pulls
		case strings.Contains(path, "/pulls") && !strings.Contains(path, "/reviews") && !strings.Contains(path, "/commits") && !strings.Contains(path, "/files") && countPathSegments(path) == 5:
			ListPulls(store).ServeHTTP(w, r)

		// GET /repos/{owner}/{repo}/pulls/{number}
		case strings.Contains(path, "/pulls") && !strings.Contains(path, "/reviews") && !strings.Contains(path, "/commits") && !strings.Contains(path, "/files") && countPathSegments(path) == 6:
			GetPull(store).ServeHTTP(w, r)

		// GET /repos/{owner}/{repo}/pulls/{number}/reviews
		case strings.Contains(path, "/reviews") && countPathSegments(path) == 7:
			ListPullReviews(store).ServeHTTP(w, r)

		// GET /repos/{owner}/{repo}/pulls/{number}/commits
		case strings.Contains(path, "/pulls") && strings.Contains(path, "/commits") && countPathSegments(path) == 7:
			ListPullCommits(store).ServeHTTP(w, r)

		// GET /repos/{owner}/{repo}/pulls/{number}/files
		case strings.Contains(path, "/pulls") && strings.Contains(path, "/files") && countPathSegments(path) == 7:
			ListPullFiles(store).ServeHTTP(w, r)

		// GET /repos/{owner}/{repo}/commits
		case strings.Contains(path, "/commits") && !strings.Contains(path, "/pulls") && countPathSegments(path) == 5:
			ListCommits(store).ServeHTTP(w, r)

		// GET /repos/{owner}/{repo}/commits/{sha}
		case strings.Contains(path, "/commits") && !strings.Contains(path, "/pulls") && countPathSegments(path) == 6:
			GetCommit(store).ServeHTTP(w, r)

		// GET /repos/{owner}/{repo}/analysis/survival
		case strings.Contains(path, "/analysis/survival"):
			SurvivalAnalysisHandler(store).ServeHTTP(w, r)

		default:
			respondError(w, http.StatusNotFound, "route not found")
		}
	})
}

// countPathSegments counts the number of non-empty path segments.
func countPathSegments(path string) int {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	count := 0
	for _, part := range parts {
		if part != "" {
			count++
		}
	}
	return count
}
