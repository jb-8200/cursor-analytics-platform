package github

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// CommitStore defines the interface for commit operations needed by the API.
type CommitStore interface {
	GetCommitByHash(hash string) (*models.Commit, error)
	GetCommitsByRepo(repoName string, from, to time.Time) []models.Commit
}

// ListCommits returns an HTTP handler for GET /repos/{owner}/{repo}/commits.
// It returns all commits for a repository with optional filtering.
func ListCommits(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repoName := parseRepoFromPath(r.URL.Path)
		if repoName == "" {
			respondError(w, http.StatusBadRequest, "invalid repository path")
			return
		}

		// Parse query parameters for date range
		fromStr := r.URL.Query().Get("since")
		toStr := r.URL.Query().Get("until")

		// Default to last 30 days if not specified
		to := time.Now().Add(24 * time.Hour)
		from := to.AddDate(0, 0, -30)

		if fromStr != "" {
			if parsed, err := time.Parse("2006-01-02", fromStr); err == nil {
				from = parsed
			}
		}

		if toStr != "" {
			if parsed, err := time.Parse("2006-01-02", toStr); err == nil {
				to = parsed.Add(24 * time.Hour)
			}
		}

		// Parse pagination parameters
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

		// Get commits from store
		commits := store.GetCommitsByRepo(repoName, from, to)

		// Apply pagination
		start := (page - 1) * perPage
		end := start + perPage

		if start > len(commits) {
			start = len(commits)
		}
		if end > len(commits) {
			end = len(commits)
		}

		pageCommits := commits[start:end]
		if pageCommits == nil {
			pageCommits = []models.Commit{}
		}

		respondJSON(w, http.StatusOK, pageCommits)
	})
}

// GetCommit returns an HTTP handler for GET /repos/{owner}/{repo}/commits/{sha}.
// It returns a single commit by SHA.
func GetCommit(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sha := parseCommitSHAFromPath(r.URL.Path)
		if sha == "" {
			respondError(w, http.StatusBadRequest, "invalid commit SHA")
			return
		}

		commit, err := store.GetCommitByHash(sha)
		if err != nil || commit == nil {
			respondError(w, http.StatusNotFound, "commit not found")
			return
		}

		respondJSON(w, http.StatusOK, commit)
	})
}

// ListPullCommits returns an HTTP handler for GET /repos/{owner}/{repo}/pulls/{number}/commits.
// It returns all commits that are part of a pull request.
func ListPullCommits(store storage.Store) http.Handler {
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

		// Get the PR to find its branch
		pr, err := store.GetPR(repoName, prNumber)
		if err != nil || pr == nil {
			respondError(w, http.StatusNotFound, "pull request not found")
			return
		}

		// Get commits from the PR's branch
		// Use a wide time range to capture all commits on the branch
		from := time.Time{}
		to := time.Now().Add(24 * time.Hour)
		commits := store.GetCommitsByRepo(repoName, from, to)

		// Filter commits to only those on the PR's head branch
		var prCommits []models.Commit
		for _, commit := range commits {
			if commit.BranchName == pr.HeadBranch && commit.RepoName == repoName {
				prCommits = append(prCommits, commit)
			}
		}

		if prCommits == nil {
			prCommits = []models.Commit{}
		}

		respondJSON(w, http.StatusOK, prCommits)
	})
}

// parseCommitSHAFromPath extracts commit SHA from the request path.
// Expected path format: /repos/{owner}/{repo}/commits/{sha}
func parseCommitSHAFromPath(path string) string {
	parts := strings.Split(path, "/")
	// Find "commits" in path and get next part
	for i, part := range parts {
		if part == "commits" && i+1 < len(parts) {
			sha := parts[i+1]
			if sha != "" {
				return sha
			}
		}
	}
	return ""
}
