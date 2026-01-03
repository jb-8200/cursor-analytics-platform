package generator

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

func TestSampleMaxCommits(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	tests := []struct {
		name      string
		seniority string
		minCommits int
		maxCommits int
	}{
		{"junior developer", "junior", 2, 5},
		{"mid-level developer", "mid", 4, 8},
		{"senior developer", "senior", 5, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sample many times to verify range
			for i := 0; i < 100; i++ {
				commits := sampleMaxCommits(tt.seniority, rng)
				if commits < tt.minCommits || commits > tt.maxCommits {
					t.Errorf("sampleMaxCommits(%s) = %d, want between %d and %d",
						tt.seniority, commits, tt.minCommits, tt.maxCommits)
				}
			}
		})
	}
}

func TestSampleTargetLoC(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	tests := []struct {
		name      string
		seniority string
		minLoC    int
		maxLoC    int
	}{
		{"junior developer", "junior", 50, 150},
		{"mid-level developer", "mid", 100, 300},
		{"senior developer", "senior", 150, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sample many times to verify range
			for i := 0; i < 100; i++ {
				loc := sampleTargetLoC(tt.seniority, rng)
				if loc < tt.minLoC || loc > tt.maxLoC {
					t.Errorf("sampleTargetLoC(%s) = %d, want between %d and %d",
						tt.seniority, loc, tt.minLoC, tt.maxLoC)
				}
			}
		})
	}
}

func TestSampleInactivityGap(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	workingHours := seed.WorkingHours{
		Start: 9,
		End:   18,
	}

	// Sample many times to verify range
	for i := 0; i < 100; i++ {
		gap := sampleInactivityGap(workingHours, rng)
		minGap := 15 * time.Minute
		maxGap := 60 * time.Minute

		if gap < minGap || gap > maxGap {
			t.Errorf("sampleInactivityGap() = %v, want between %v and %v",
				gap, minGap, maxGap)
		}
	}
}

func TestStartSession(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	developer := seed.Developer{
		UserID:    "user_001",
		Email:     "alice@example.com",
		Name:      "Alice",
		Seniority: "senior",
		WorkingHoursBand: seed.WorkingHours{
			Start: 9,
			End:   18,
		},
	}

	repository := seed.Repository{
		RepoName: "acme/platform",
		Maturity: seed.Maturity{
			AgeDays: 365,
		},
	}

	branch := "feature/auth-improvement"
	startTime := time.Now()

	session := StartSession(developer, repository, branch, startTime, rng)

	// Verify session fields
	if session.Developer.UserID != developer.UserID {
		t.Errorf("session.Developer.UserID = %s, want %s", session.Developer.UserID, developer.UserID)
	}

	if session.Repo.RepoName != repository.RepoName {
		t.Errorf("session.Repo.RepoName = %s, want %s", session.Repo.RepoName, repository.RepoName)
	}

	if session.Branch != branch {
		t.Errorf("session.Branch = %s, want %s", session.Branch, branch)
	}

	if !session.StartTime.Equal(startTime) {
		t.Errorf("session.StartTime = %v, want %v", session.StartTime, startTime)
	}

	// Verify seniority-based parameters are set
	if session.MaxCommits < 5 || session.MaxCommits > 12 {
		t.Errorf("session.MaxCommits = %d, want between 5 and 12 for senior developer", session.MaxCommits)
	}

	if session.TargetLoC < 150 || session.TargetLoC > 500 {
		t.Errorf("session.TargetLoC = %d, want between 150 and 500 for senior developer", session.TargetLoC)
	}

	if session.InactivityGap < 15*time.Minute || session.InactivityGap > 60*time.Minute {
		t.Errorf("session.InactivityGap = %v, want between 15m and 60m", session.InactivityGap)
	}

	// Verify commits slice is initialized
	if session.Commits == nil {
		t.Error("session.Commits should be initialized")
	}

	if len(session.Commits) != 0 {
		t.Errorf("session.Commits length = %d, want 0", len(session.Commits))
	}
}

func TestSessionReproducibility(t *testing.T) {
	// Same seed should produce same session parameters
	developer := seed.Developer{
		UserID:    "user_001",
		Seniority: "mid",
		WorkingHoursBand: seed.WorkingHours{
			Start: 9,
			End:   18,
		},
	}

	repository := seed.Repository{
		RepoName: "acme/platform",
	}

	branch := "feature/test"
	startTime := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	// Create two sessions with same seed
	rng1 := rand.New(rand.NewSource(12345))
	session1 := StartSession(developer, repository, branch, startTime, rng1)

	rng2 := rand.New(rand.NewSource(12345))
	session2 := StartSession(developer, repository, branch, startTime, rng2)

	// Verify they are identical
	if session1.MaxCommits != session2.MaxCommits {
		t.Errorf("MaxCommits not reproducible: %d vs %d", session1.MaxCommits, session2.MaxCommits)
	}

	if session1.TargetLoC != session2.TargetLoC {
		t.Errorf("TargetLoC not reproducible: %d vs %d", session1.TargetLoC, session2.TargetLoC)
	}

	if session1.InactivityGap != session2.InactivityGap {
		t.Errorf("InactivityGap not reproducible: %v vs %v", session1.InactivityGap, session2.InactivityGap)
	}
}
