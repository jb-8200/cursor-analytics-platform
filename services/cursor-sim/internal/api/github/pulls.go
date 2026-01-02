package github

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// PRStore defines the interface for PR operations needed by the API.
type PRStore interface {
	GetPR(repoName string, number int) (*models.PullRequest, error)
	GetPRsByRepo(repoName string) []models.PullRequest
	GetPRsByRepoAndState(repoName string, state models.PRState) []models.PullRequest
	GetReviewComments(repoName string, prNumber int) []models.ReviewComment
}

// parseRepoFromPath extracts owner/repo from the request path.
// Expected path format: /repos/{owner}/{repo}/pulls...
func parseRepoFromPath(path string) string {
	parts := strings.Split(path, "/")
	// Find "repos" in path and get next two parts
	for i, part := range parts {
		if part == "repos" && i+2 < len(parts) {
			return parts[i+1] + "/" + parts[i+2]
		}
	}
	return ""
}

// parsePRNumberFromPath extracts PR number from the request path.
// Expected path format: /repos/{owner}/{repo}/pulls/{number}...
func parsePRNumberFromPath(path string) int {
	parts := strings.Split(path, "/")
	// Find "pulls" in path and get next part
	for i, part := range parts {
		if part == "pulls" && i+1 < len(parts) {
			num, err := strconv.Atoi(parts[i+1])
			if err == nil {
				return num
			}
		}
	}
	return 0
}

// ListPulls returns an HTTP handler for GET /repos/{owner}/{repo}/pulls.
// It returns all PRs for a repository with optional state filtering.
func ListPulls(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repoName := parseRepoFromPath(r.URL.Path)
		if repoName == "" {
			respondError(w, http.StatusBadRequest, "invalid repository path")
			return
		}

		// Parse query parameters
		state := r.URL.Query().Get("state")
		perPage := 30 // GitHub default
		page := 1

		if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
			if p, err := strconv.Atoi(perPageStr); err == nil && p > 0 {
				perPage = p
				if perPage > 100 {
					perPage = 100 // GitHub max
				}
			}
		}

		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		// Get PRs with optional state filter
		var prs []models.PullRequest
		if state != "" && state != "all" {
			prState := models.PRState(state)
			prs = store.GetPRsByRepoAndState(repoName, prState)
		} else {
			prs = store.GetPRsByRepo(repoName)
		}

		// Apply pagination
		start := (page - 1) * perPage
		end := start + perPage

		if start > len(prs) {
			start = len(prs)
		}
		if end > len(prs) {
			end = len(prs)
		}

		pagePRs := prs[start:end]

		respondJSON(w, http.StatusOK, pagePRs)
	})
}

// GetPull returns an HTTP handler for GET /repos/{owner}/{repo}/pulls/{number}.
// It returns a single PR by number.
func GetPull(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repoName := parseRepoFromPath(r.URL.Path)
		if repoName == "" {
			respondError(w, http.StatusBadRequest, "invalid repository path")
			return
		}

		prNumber := parsePRNumberFromPath(r.URL.Path)
		if prNumber == 0 {
			respondError(w, http.StatusBadRequest, "invalid pull request number")
			return
		}

		pr, err := store.GetPR(repoName, prNumber)
		if err != nil || pr == nil {
			respondError(w, http.StatusNotFound, "pull request not found")
			return
		}

		respondJSON(w, http.StatusOK, pr)
	})
}

// ListPullReviews returns an HTTP handler for GET /repos/{owner}/{repo}/pulls/{number}/reviews.
// It returns all reviews for a PR.
func ListPullReviews(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repoName := parseRepoFromPath(r.URL.Path)
		if repoName == "" {
			respondError(w, http.StatusBadRequest, "invalid repository path")
			return
		}

		prNumber := parsePRNumberFromPath(r.URL.Path)
		if prNumber == 0 {
			respondError(w, http.StatusBadRequest, "invalid pull request number")
			return
		}

		reviews := store.GetReviewComments(repoName, prNumber)
		if reviews == nil {
			reviews = []models.ReviewComment{}
		}

		respondJSON(w, http.StatusOK, reviews)
	})
}

// respondJSON writes a JSON response.
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// respondError writes a JSON error response.
func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
