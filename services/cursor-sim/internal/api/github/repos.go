package github

import (
	"net/http"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// RepoInfo represents repository information returned by the API.
type RepoInfo struct {
	FullName      string `json:"full_name"`
	DefaultBranch string `json:"default_branch"`
	OpenPRs       int    `json:"open_pull_requests_count"`
	TotalPRs      int    `json:"total_pull_requests_count"`
}

// ListRepos returns an HTTP handler for GET /repos.
// It returns all repositories that have PRs.
func ListRepos(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repoNames := store.ListRepositories()

		repos := make([]RepoInfo, 0, len(repoNames))
		for _, name := range repoNames {
			prs := store.GetPRsByRepo(name)

			openCount := 0
			for _, pr := range prs {
				if pr.State == models.PRStateOpen {
					openCount++
				}
			}

			repos = append(repos, RepoInfo{
				FullName:      name,
				DefaultBranch: "main",
				OpenPRs:       openCount,
				TotalPRs:      len(prs),
			})
		}

		respondJSON(w, http.StatusOK, repos)
	})
}

// GetRepo returns an HTTP handler for GET /repos/{owner}/{repo}.
// It returns information about a single repository.
func GetRepo(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repoName := parseRepoFromPath(r.URL.Path)
		if repoName == "" {
			respondError(w, http.StatusBadRequest, "invalid repository path")
			return
		}

		prs := store.GetPRsByRepo(repoName)

		// If no PRs, repo doesn't exist in our store
		if len(prs) == 0 {
			respondError(w, http.StatusNotFound, "repository not found")
			return
		}

		openCount := 0
		for _, pr := range prs {
			if pr.State == models.PRStateOpen {
				openCount++
			}
		}

		repo := RepoInfo{
			FullName:      repoName,
			DefaultBranch: "main",
			OpenPRs:       openCount,
			TotalPRs:      len(prs),
		}

		respondJSON(w, http.StatusOK, repo)
	})
}
