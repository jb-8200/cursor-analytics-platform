package cursor

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// Regenerate creates a handler for POST /admin/regenerate.
// Allows runtime reconfiguration and data regeneration without service restart.
//
// Supports two modes:
// - append: Adds new data to existing storage
// - override: Clears storage and generates fresh data
func Regenerate(store storage.Store, seedData *seed.SeedData) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only POST is allowed
		if r.Method != http.MethodPost {
			api.RespondError(w, http.StatusMethodNotAllowed, "Only POST is allowed")
			return
		}

		// Parse request body
		var req models.RegenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse request: %v", err))
			return
		}

		// Validate request
		if err := validateRegenerateRequest(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Track start time for duration
		startTime := time.Now()

		// Get stats before operation
		statsBefore := store.GetStats()

		// Handle override mode (clear data)
		dataCleaned := false
		if req.Mode == "override" {
			if err := store.ClearAllData(); err != nil {
				api.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to clear data: %v", err))
				return
			}
			dataCleaned = true
		}

		// Replicate developers if needed
		targetDevelopers := req.Developers
		if targetDevelopers == 0 {
			targetDevelopers = len(seedData.Developers)
		}

		// Load developers with replication
		developers := seedData.Developers
		if targetDevelopers > len(seedData.Developers) {
			// Replicate developers to match target count
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			var err error
			developers, err = seed.ReplicateDevelopers(seedData, targetDevelopers, rng)
			if err != nil {
				api.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to replicate developers: %v", err))
				return
			}
		}

		// Load developers into storage
		if err := store.LoadDevelopers(developers); err != nil {
			api.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to load developers: %v", err))
			return
		}

		// Create context for generation
		ctx := context.Background()

		// Generate commits
		commitGen := generator.NewCommitGenerator(seedData, store, req.Velocity)
		if err := commitGen.GenerateCommits(ctx, req.Days, req.MaxCommits); err != nil {
			api.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate commits: %v", err))
			return
		}

		// Generate PRs from commits
		prGen := generator.NewPRGeneratorWithSeed(seedData, store, time.Now().UnixNano())
		startDate := time.Now().AddDate(0, 0, -req.Days)
		endDate := time.Now().Add(24 * time.Hour)
		if err := prGen.GeneratePRsFromCommits(startDate, endDate); err != nil {
			api.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate PRs: %v", err))
			return
		}

		// Generate reviews for PRs
		reviewGen := generator.NewReviewGenerator(seedData, rand.New(rand.NewSource(time.Now().UnixNano())))
		repos := store.ListRepositories()
		for _, repoName := range repos {
			prs := store.GetPRsByRepo(repoName)
			for _, pr := range prs {
				reviews := reviewGen.GenerateReviewsForPR(pr)
				for _, review := range reviews {
					if err := store.StoreReview(review); err != nil {
						// Log warning but continue
						continue
					}
				}
			}
		}

		// Generate issues for PRs
		issueGen := generator.NewIssueGeneratorWithStore(seedData, store, time.Now().UnixNano())
		for _, repoName := range repos {
			prs := store.GetPRsByRepo(repoName)
			if _, err := issueGen.GenerateAndStoreIssuesForPRs(prs, repoName); err != nil {
				// Log warning but continue
				continue
			}
		}

		// Generate feature events
		featureGen := generator.NewFeatureGenerator(seedData, store, req.Velocity)
		if err := featureGen.GenerateFeatures(ctx, req.Days); err != nil {
			// Log warning but continue (non-critical)
		}

		// Generate file extension events
		extGen := generator.NewExtensionGenerator(seedData, store, req.Velocity)
		if err := extGen.GenerateFileExtensions(ctx, req.Days); err != nil {
			// Log warning but continue (non-critical)
		}

		// Get stats after operation
		statsAfter := store.GetStats()

		// Calculate duration
		duration := time.Since(startTime)

		// Build response
		response := models.RegenerateResponse{
			Status:          "success",
			Mode:            req.Mode,
			DataCleaned:     dataCleaned,
			CommitsAdded:    statsAfter.Commits - statsBefore.Commits,
			PRsAdded:        statsAfter.PullRequests - statsBefore.PullRequests,
			ReviewsAdded:    statsAfter.Reviews - statsBefore.Reviews,
			IssuesAdded:     statsAfter.Issues - statsBefore.Issues,
			TotalCommits:    statsAfter.Commits,
			TotalPRs:        statsAfter.PullRequests,
			TotalDevelopers: statsAfter.Developers,
			Duration:        duration.Round(time.Millisecond).String(),
			Config: models.ConfigParams{
				Days:       req.Days,
				Velocity:   req.Velocity,
				Developers: targetDevelopers,
				MaxCommits: req.MaxCommits,
			},
		}

		api.RespondJSON(w, http.StatusOK, response)
	})
}

// validateRegenerateRequest validates the regenerate request parameters.
// Returns an error if any parameter is invalid.
func validateRegenerateRequest(req *models.RegenerateRequest) error {
	// Validate mode
	if req.Mode != "append" && req.Mode != "override" {
		return fmt.Errorf("mode must be 'append' or 'override', got: %s", req.Mode)
	}

	// Validate days (1-3650, approximately 10 years)
	if req.Days < 1 || req.Days > 3650 {
		return fmt.Errorf("days must be between 1 and 3650, got: %d", req.Days)
	}

	// Validate velocity
	if req.Velocity != "low" && req.Velocity != "medium" && req.Velocity != "high" {
		return fmt.Errorf("velocity must be 'low', 'medium', or 'high', got: %s", req.Velocity)
	}

	// Validate developers (0-10000, 0 means use seed count)
	if req.Developers < 0 || req.Developers > 10000 {
		return fmt.Errorf("developers must be between 0 and 10000, got: %d", req.Developers)
	}

	// Validate max_commits (0-100000, 0 means unlimited)
	if req.MaxCommits < 0 || req.MaxCommits > 100000 {
		return fmt.Errorf("max_commits must be between 0 and 100000, got: %d", req.MaxCommits)
	}

	return nil
}
