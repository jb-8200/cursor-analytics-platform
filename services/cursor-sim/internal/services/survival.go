package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

// CommitStore defines the interface for accessing commits needed by survival analysis.
type CommitStore interface {
	GetCommitsByRepo(repoName string, from, to time.Time) []models.Commit
}

// SurvivalService calculates code survival metrics for files in a repository.
type SurvivalService struct {
	store CommitStore
	rng   *rand.Rand
}

// NewSurvivalService creates a new survival analysis service.
func NewSurvivalService(store CommitStore) *SurvivalService {
	return &SurvivalService{
		store: store,
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewSurvivalServiceWithSeed creates a survival service with a specific random seed for reproducibility.
func NewSurvivalServiceWithSeed(store CommitStore, seed int64) *SurvivalService {
	return &SurvivalService{
		store: store,
		rng:   rand.New(rand.NewSource(seed)),
	}
}

// CalculateSurvival calculates survival rates for files added in a cohort period.
// It tracks files from creation in the cohort window to the observation date.
func (s *SurvivalService) CalculateSurvival(
	repoName string,
	cohortStart, cohortEnd, observationDate time.Time,
) (*models.SurvivalAnalysis, error) {
	// Get all commits in the cohort period
	cohortCommits := s.store.GetCommitsByRepo(repoName, cohortStart, cohortEnd)

	if len(cohortCommits) == 0 {
		return &models.SurvivalAnalysis{
			CohortStart:     cohortStart.Format("2006-01-02"),
			CohortEnd:       cohortEnd.Format("2006-01-02"),
			ObservationDate: observationDate.Format("2006-01-02"),
			TotalLinesAdded: 0,
			LinesSurviving:  0,
			SurvivalRate:    0.0,
			ByDeveloper:     []models.DeveloperSurvival{},
		}, nil
	}

	// Track file births in cohort period
	fileMap := make(map[string]*models.FileSurvival)
	fileToEmail := make(map[string]string) // Map file to developer email
	developerStats := make(map[string]*developerSurvivalStats)

	// Process cohort commits to track file births
	for _, commit := range cohortCommits {
		// Generate files for this commit
		files := s.generateFilesForCommit(commit, cohortStart)

		for _, file := range files {
			// Track file birth (only count first occurrence)
			if _, exists := fileMap[file.FilePath]; !exists {
				fileMap[file.FilePath] = file
				fileToEmail[file.FilePath] = commit.UserEmail
			}

			// Track developer stats (only for new files)
			email := commit.UserEmail
			if _, exists := developerStats[email]; !exists {
				developerStats[email] = &developerSurvivalStats{
					Email:          email,
					LinesAdded:     0,
					LinesSurviving: 0,
				}
			}

			// Only add lines once per file
			if fileToEmail[file.FilePath] == email {
				developerStats[email].LinesAdded += file.AILinesAdded + file.HumanLinesAdded
			}
		}
	}

	// Simulate file survival based on observation date
	// Files have a probabilistic chance of being deleted over time
	daysSinceCohort := observationDate.Sub(cohortEnd).Hours() / 24
	deletionProbability := calculateDeletionProbability(daysSinceCohort)

	totalLinesAdded := 0
	linesSurviving := 0

	for _, file := range fileMap {
		totalLinesAdded += file.TotalLines

		// Determine if file survived
		if s.rng.Float64() > deletionProbability {
			// File survived
			linesSurviving += file.TotalLines

			// Update developer survival stats
			if email, exists := fileToEmail[file.FilePath]; exists {
				if stats, exists := developerStats[email]; exists {
					stats.LinesSurviving += file.TotalLines
				}
			}
		} else {
			// File was deleted
			file.IsDeleted = true
			deletedTime := cohortEnd.Add(time.Duration(s.rng.Intn(int(daysSinceCohort))) * 24 * time.Hour)
			file.DeletedAt = &deletedTime
		}
	}

	// Build developer breakdown
	byDeveloper := make([]models.DeveloperSurvival, 0, len(developerStats))
	for _, stats := range developerStats {
		survivalRate := 0.0
		if stats.LinesAdded > 0 {
			survivalRate = float64(stats.LinesSurviving) / float64(stats.LinesAdded)
		}

		byDeveloper = append(byDeveloper, models.DeveloperSurvival{
			Email:          stats.Email,
			LinesAdded:     stats.LinesAdded,
			LinesSurviving: stats.LinesSurviving,
			SurvivalRate:   survivalRate,
		})
	}

	// Calculate overall survival rate
	survivalRate := 0.0
	if totalLinesAdded > 0 {
		survivalRate = float64(linesSurviving) / float64(totalLinesAdded)
	}

	return &models.SurvivalAnalysis{
		CohortStart:     cohortStart.Format("2006-01-02"),
		CohortEnd:       cohortEnd.Format("2006-01-02"),
		ObservationDate: observationDate.Format("2006-01-02"),
		TotalLinesAdded: totalLinesAdded,
		LinesSurviving:  linesSurviving,
		SurvivalRate:    survivalRate,
		ByDeveloper:     byDeveloper,
	}, nil
}

// developerSurvivalStats tracks survival stats for a developer during calculation.
type developerSurvivalStats struct {
	Email          string
	LinesAdded     int
	LinesSurviving int
}

// generateFilesForCommit generates synthetic files for a commit.
// In a real implementation, this would parse the commit diff.
func (s *SurvivalService) generateFilesForCommit(commit models.Commit, cohortStart time.Time) []*models.FileSurvival {
	// Generate 1-3 files per commit
	numFiles := 1 + s.rng.Intn(3)
	files := make([]*models.FileSurvival, numFiles)

	for i := 0; i < numFiles; i++ {
		// Distribute lines across files
		linesPerFile := commit.TotalLinesAdded / numFiles
		aiLines := commit.TabLinesAdded + commit.ComposerLinesAdded
		aiLinesPerFile := aiLines / numFiles
		humanLinesPerFile := linesPerFile - aiLinesPerFile

		files[i] = &models.FileSurvival{
			FilePath:        fmt.Sprintf("%s/file_%d.go", commit.RepoName, s.rng.Intn(1000)),
			RepoName:        commit.RepoName,
			CreatedAt:       commit.CommitTs,
			LastModifiedAt:  commit.CommitTs,
			AILinesAdded:    aiLinesPerFile,
			HumanLinesAdded: humanLinesPerFile,
			TotalLines:      linesPerFile,
			RevertEvents:    0,
			IsDeleted:       false,
		}
	}

	return files
}

// calculateDeletionProbability returns the probability a file is deleted based on days elapsed.
// Uses a sigmoid curve: files are more likely to be deleted as time passes.
func calculateDeletionProbability(daysSinceCohort float64) float64 {
	// Sigmoid curve centered at 45 days
	// At 30 days: ~10% deletion
	// At 45 days: ~20% deletion
	// At 60 days: ~30% deletion
	// At 90 days: ~40% deletion
	if daysSinceCohort <= 0 {
		return 0.0
	}

	// Simple linear model for simulation
	probability := daysSinceCohort / 200.0
	if probability > 0.5 {
		probability = 0.5 // Cap at 50% deletion
	}

	return probability
}

// getEmailFromFilePath extracts the developer email from file path and commits.
// This is a helper for attribution in the simulation.
func getEmailFromFilePath(filePath string, commits []models.Commit) string {
	// For simulation, just use the first commit's email
	if len(commits) > 0 {
		return commits[0].UserEmail
	}
	return "unknown@example.com"
}
