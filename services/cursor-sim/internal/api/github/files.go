package github

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// FileChange represents a file changed in a pull request.
type FileChange struct {
	Filename       string  `json:"filename"`
	Additions      int     `json:"additions"`
	Deletions      int     `json:"deletions"`
	Changes        int     `json:"changes"`
	GreenfieldIdx  float64 `json:"greenfield_index"`
	IsGreenfield   bool    `json:"is_greenfield"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
}

// ListPullFiles returns an HTTP handler for GET /repos/{owner}/{repo}/pulls/{number}/files.
// It returns all files changed in a pull request with greenfield index.
func ListPullFiles(store storage.Store) http.Handler {
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

		// Get the PR to find its branch and metrics
		pr, err := store.GetPR(repoName, prNumber)
		if err != nil || pr == nil {
			respondError(w, http.StatusNotFound, "pull request not found")
			return
		}

		// Generate file list based on PR metrics
		// Each PR has ChangedFiles entries
		files := generatePRFiles(pr)

		if files == nil {
			files = []FileChange{}
		}

		respondJSON(w, http.StatusOK, files)
	})
}

// generatePRFiles generates a synthetic list of files changed in a PR.
// The number of files is based on the PR's ChangedFiles count.
// File sizes are distributed proportionally to PR additions/deletions.
func generatePRFiles(pr *models.PullRequest) []FileChange {
	if pr.ChangedFiles == 0 {
		// Default to 1 file if not specified
		pr.ChangedFiles = 1
	}

	rng := rand.New(rand.NewSource(hashPRNumber(pr.Number)))
	files := make([]FileChange, pr.ChangedFiles)

	// Distribute additions and deletions across files
	remainingAdditions := pr.Additions
	remainingDeletions := pr.Deletions

	for i := 0; i < pr.ChangedFiles; i++ {
		// Allocate portion of changes to this file
		var additions, deletions int
		if i == pr.ChangedFiles-1 {
			// Last file gets remainder
			additions = remainingAdditions
			deletions = remainingDeletions
		} else {
			// Distribute remaining across files
			additions = remainingAdditions / (pr.ChangedFiles - i)
			deletions = remainingDeletions / (pr.ChangedFiles - i)
			remainingAdditions -= additions
			remainingDeletions -= deletions
		}

		// Generate filename
		filename := generateFileName(rng, pr.RepoName, i)

		// Calculate greenfield index
		// Files created in the last 30 days are greenfield
		// Use PR creation time as file creation reference
		fileCreatedAt := pr.CreatedAt.AddDate(0, 0, -rng.Intn(30))
		isGreenfield := fileCreatedAt.After(pr.CreatedAt.AddDate(0, 0, -30))
		greenfieldIdx := 0.0
		if isGreenfield {
			greenfieldIdx = 1.0
		}

		files[i] = FileChange{
			Filename:      filename,
			Additions:     additions,
			Deletions:     deletions,
			Changes:       additions + deletions,
			GreenfieldIdx: greenfieldIdx,
			IsGreenfield:  isGreenfield,
			CreatedAt:     &fileCreatedAt,
		}
	}

	return files
}

// generateFileName generates a realistic filename based on repository and index.
func generateFileName(rng *rand.Rand, repoName string, index int) string {
	extensions := []string{".go", ".ts", ".tsx", ".py", ".js", ".jsx", ".java", ".rs", ".cpp"}
	paths := []string{"cmd/", "internal/", "src/", "lib/", "pkg/", "test/", "tests/"}

	ext := extensions[rng.Intn(len(extensions))]
	path := paths[rng.Intn(len(paths))]

	// Generate filename based on repo name and index
	parts := strings.Split(repoName, "/")
	repoShort := parts[len(parts)-1]
	if len(repoShort) > 3 {
		repoShort = repoShort[:3]
	}

	names := []string{
		repoShort + "_handler",
		repoShort + "_service",
		repoShort + "_model",
		repoShort + "_client",
		repoShort + "_config",
		repoShort + "_utils",
	}

	name := names[index%len(names)]
	return path + name + ext
}

// hashPRNumber generates a deterministic seed from PR number.
// This ensures files generated for the same PR are consistent.
func hashPRNumber(prNumber int) int64 {
	return int64(prNumber)*31 + 17
}
