package seed

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReplicateDevelopers_Downsample tests sampling when N < seed count
func TestReplicateDevelopers_Downsample(t *testing.T) {
	// Arrange: Seed with 5 developers, request 2
	seed := &SeedData{
		Developers: []Developer{
			{UserID: "dev1", Email: "dev1@example.com", Name: "Dev One"},
			{UserID: "dev2", Email: "dev2@example.com", Name: "Dev Two"},
			{UserID: "dev3", Email: "dev3@example.com", Name: "Dev Three"},
			{UserID: "dev4", Email: "dev4@example.com", Name: "Dev Four"},
			{UserID: "dev5", Email: "dev5@example.com", Name: "Dev Five"},
		},
	}
	rng := rand.New(rand.NewSource(12345)) // Deterministic

	// Act
	result, err := ReplicateDevelopers(seed, 2, rng)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	// Verify they are from the original seed
	for _, dev := range result {
		found := false
		for _, orig := range seed.Developers {
			if dev.UserID == orig.UserID && dev.Email == orig.Email {
				found = true
				break
			}
		}
		assert.True(t, found, "developer %s should be from original seed", dev.UserID)
	}
}

// TestReplicateDevelopers_ExactMatch tests when N == seed count
func TestReplicateDevelopers_ExactMatch(t *testing.T) {
	// Arrange: Seed with 3 developers, request 3
	seed := &SeedData{
		Developers: []Developer{
			{UserID: "alice", Email: "alice@example.com", Name: "Alice"},
			{UserID: "bob", Email: "bob@example.com", Name: "Bob"},
			{UserID: "charlie", Email: "charlie@example.com", Name: "Charlie"},
		},
	}
	rng := rand.New(rand.NewSource(12345))

	// Act
	result, err := ReplicateDevelopers(seed, 3, rng)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 3)
	// Should return all developers (possibly shuffled)
	userIDs := make(map[string]bool)
	for _, dev := range result {
		userIDs[dev.UserID] = true
	}
	assert.True(t, userIDs["alice"])
	assert.True(t, userIDs["bob"])
	assert.True(t, userIDs["charlie"])
}

// TestReplicateDevelopers_Replicate tests cloning when N > seed count
func TestReplicateDevelopers_Replicate(t *testing.T) {
	// Arrange: Seed with 2 developers, request 5
	seed := &SeedData{
		Developers: []Developer{
			{UserID: "user_001", Email: "alice@example.com", Name: "Alice Developer"},
			{UserID: "user_002", Email: "bob@example.com", Name: "Bob Developer"},
		},
	}
	rng := rand.New(rand.NewSource(12345))

	// Act
	result, err := ReplicateDevelopers(seed, 5, rng)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 5)

	// Check that we have 5 unique user IDs
	userIDs := make(map[string]bool)
	for _, dev := range result {
		userIDs[dev.UserID] = true
	}
	assert.Len(t, userIDs, 5, "should have 5 unique user IDs")
}

// TestReplicateDevelopers_UniqueIDs ensures all clones have unique IDs
func TestReplicateDevelopers_UniqueIDs(t *testing.T) {
	// Arrange: Seed with 1 developer, request 10
	seed := &SeedData{
		Developers: []Developer{
			{UserID: "dev_base", Email: "base@example.com", Name: "Base Dev"},
		},
	}
	rng := rand.New(rand.NewSource(12345))

	// Act
	result, err := ReplicateDevelopers(seed, 10, rng)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 10)

	// Verify all user IDs are unique
	userIDs := make(map[string]bool)
	for _, dev := range result {
		assert.False(t, userIDs[dev.UserID], "duplicate user ID: %s", dev.UserID)
		userIDs[dev.UserID] = true
	}
	assert.Len(t, userIDs, 10, "should have 10 unique user IDs")
}

// TestReplicateDevelopers_CloneNamingConvention verifies clone naming format
func TestReplicateDevelopers_CloneNamingConvention(t *testing.T) {
	// Arrange
	seed := &SeedData{
		Developers: []Developer{
			{UserID: "user_001", Email: "alice@example.com", Name: "Alice"},
		},
	}
	rng := rand.New(rand.NewSource(12345))

	// Act: Request 4 developers (1 original + 3 clones)
	result, err := ReplicateDevelopers(seed, 4, rng)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 4)

	// Check naming pattern
	assert.Equal(t, "user_001", result[0].UserID)
	assert.Equal(t, "alice@example.com", result[0].Email)
	assert.Equal(t, "Alice", result[0].Name)

	assert.Equal(t, "user_001_clone1", result[1].UserID)
	assert.Equal(t, "clone1_alice@example.com", result[1].Email)
	assert.Equal(t, "Alice (Clone 1)", result[1].Name)

	assert.Equal(t, "user_001_clone2", result[2].UserID)
	assert.Equal(t, "clone2_alice@example.com", result[2].Email)
	assert.Equal(t, "Alice (Clone 2)", result[2].Name)

	assert.Equal(t, "user_001_clone3", result[3].UserID)
	assert.Equal(t, "clone3_alice@example.com", result[3].Email)
	assert.Equal(t, "Alice (Clone 3)", result[3].Name)
}

// TestReplicateDevelopers_InvalidInputs tests error handling
func TestReplicateDevelopers_InvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		seed        *SeedData
		targetCount int
		wantErr     string
	}{
		{
			name:        "zero target count",
			seed:        &SeedData{Developers: []Developer{{UserID: "dev1"}}},
			targetCount: 0,
			wantErr:     "target count must be >= 1",
		},
		{
			name:        "negative target count",
			seed:        &SeedData{Developers: []Developer{{UserID: "dev1"}}},
			targetCount: -5,
			wantErr:     "target count must be >= 1",
		},
		{
			name:        "empty seed developers",
			seed:        &SeedData{Developers: []Developer{}},
			targetCount: 5,
			wantErr:     "seed data has no developers",
		},
		{
			name:        "nil developers",
			seed:        &SeedData{},
			targetCount: 5,
			wantErr:     "seed data has no developers",
		},
	}

	rng := rand.New(rand.NewSource(12345))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReplicateDevelopers(tt.seed, tt.targetCount, rng)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// TestReplicateDevelopers_PreservesAttributes verifies clones preserve attributes
func TestReplicateDevelopers_PreservesAttributes(t *testing.T) {
	// Arrange
	seed := &SeedData{
		Developers: []Developer{
			{
				UserID:         "dev1",
				Email:          "dev1@example.com",
				Name:           "Developer One",
				Org:            "Acme Corp",
				Division:       "Engineering",
				Team:           "Platform",
				Role:           "engineer",
				Region:         "us-west",
				Timezone:       "America/Los_Angeles",
				Seniority:      "senior",
				ActivityLevel:  "high",
				AcceptanceRate: 0.85,
			},
		},
	}
	rng := rand.New(rand.NewSource(12345))

	// Act: Request 3 developers
	result, err := ReplicateDevelopers(seed, 3, rng)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 3)

	// Verify all clones preserve attributes (except ID, email, name)
	for _, dev := range result {
		assert.Equal(t, "Acme Corp", dev.Org)
		assert.Equal(t, "Engineering", dev.Division)
		assert.Equal(t, "Platform", dev.Team)
		assert.Equal(t, "engineer", dev.Role)
		assert.Equal(t, "us-west", dev.Region)
		assert.Equal(t, "America/Los_Angeles", dev.Timezone)
		assert.Equal(t, "senior", dev.Seniority)
		assert.Equal(t, "high", dev.ActivityLevel)
		assert.Equal(t, 0.85, dev.AcceptanceRate)
	}
}
