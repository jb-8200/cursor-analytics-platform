package generator

import (
	"math/rand"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// Session represents a work period that produces a PR from grouped commits.
// Session parameters enforce developer-specific characteristics (seniority, working hours).
type Session struct {
	Developer     seed.Developer
	Repo          seed.Repository
	Branch        string
	StartTime     time.Time
	MaxCommits    int           // Seniority-based: juniors 2-5, seniors 5-12
	TargetLoC     int           // Affects commit sizes
	InactivityGap time.Duration // From working hours: 15-60 minutes
	Commits       []models.Commit
}

// sampleMaxCommits samples the maximum number of commits per PR based on developer seniority.
// Returns:
// - junior: 2-5 commits
// - mid: 4-8 commits
// - senior: 5-12 commits
func sampleMaxCommits(seniority string, rng *rand.Rand) int {
	switch seniority {
	case "junior":
		return 2 + rng.Intn(4) // 2-5
	case "mid":
		return 4 + rng.Intn(5) // 4-8
	case "senior":
		return 5 + rng.Intn(8) // 5-12
	default:
		// Default to mid-level if seniority is unknown
		return 4 + rng.Intn(5)
	}
}

// sampleTargetLoC samples the target lines of code per PR based on developer seniority.
// This affects the size of commits within the session.
// Returns:
// - junior: 50-150 LoC
// - mid: 100-300 LoC
// - senior: 150-500 LoC
func sampleTargetLoC(seniority string, rng *rand.Rand) int {
	switch seniority {
	case "junior":
		return 50 + rng.Intn(101) // 50-150
	case "mid":
		return 100 + rng.Intn(201) // 100-300
	case "senior":
		return 150 + rng.Intn(351) // 150-500
	default:
		// Default to mid-level if seniority is unknown
		return 100 + rng.Intn(201)
	}
}

// sampleInactivityGap samples the maximum inactivity gap before closing a session.
// Returns a duration between 15 and 60 minutes.
// This is developer-specific based on their working hours and work style.
func sampleInactivityGap(workingHours seed.WorkingHours, rng *rand.Rand) time.Duration {
	// Sample between 15 and 60 minutes
	minMinutes := 15
	maxMinutes := 60
	minutes := minMinutes + rng.Intn(maxMinutes-minMinutes+1)
	return time.Duration(minutes) * time.Minute
}

// StartSession creates a new work session with seniority-based parameters.
// This enforces correlations at the session level:
// - Senior developers have larger PRs (more commits, more LoC)
// - Working hours affect inactivity gap
// - Session parameters are sampled from distributions using the provided RNG
func StartSession(developer seed.Developer, repo seed.Repository, branch string, startTime time.Time, rng *rand.Rand) *Session {
	return &Session{
		Developer:     developer,
		Repo:          repo,
		Branch:        branch,
		StartTime:     startTime,
		MaxCommits:    sampleMaxCommits(developer.Seniority, rng),
		TargetLoC:     sampleTargetLoC(developer.Seniority, rng),
		InactivityGap: sampleInactivityGap(developer.WorkingHoursBand, rng),
		Commits:       make([]models.Commit, 0),
	}
}

// AddCommit adds a commit to the session.
func (s *Session) AddCommit(commit models.Commit) {
	s.Commits = append(s.Commits, commit)
}

// ShouldClose determines if the session should be closed based on:
// 1. Max commits reached
// 2. Inactivity gap exceeded
// 3. Random early close (volatility)
func (s *Session) ShouldClose(lastCommitTime time.Time, rng *rand.Rand) bool {
	// Rule 1: Max commits reached
	if len(s.Commits) >= s.MaxCommits {
		return true
	}

	// Rule 2: Inactivity gap exceeded
	if len(s.Commits) > 0 {
		timeSinceLastCommit := lastCommitTime.Sub(s.Commits[len(s.Commits)-1].CommitTs)
		if timeSinceLastCommit > s.InactivityGap {
			return true
		}
	}

	// Rule 3: Random early close (1% chance after 3+ commits for volatility)
	if len(s.Commits) >= 3 && rng.Float64() < 0.01 {
		return true
	}

	return false
}
