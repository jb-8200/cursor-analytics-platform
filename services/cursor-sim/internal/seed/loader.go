package seed

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// LoadSeed reads and validates a seed.json file from the specified path.
// Returns an error if the file cannot be read, parsed, or if validation fails.
func LoadSeed(path string) (*SeedData, error) {
	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read seed file %q: %w", path, err)
	}

	// Parse JSON
	var seed SeedData
	if err := json.Unmarshal(data, &seed); err != nil {
		return nil, fmt.Errorf("failed to parse seed file %q: %w", path, err)
	}

	// Validate the seed data
	if err := seed.Validate(); err != nil {
		return nil, err
	}

	return &seed, nil
}

// LoadSeedWithReplication reads a seed file and optionally replicates developers.
// If developerCount is 0, returns the original developers from the seed.
// If developerCount > 0, uses ReplicateDevelopers to create the desired number.
// If rng is nil, creates a new random number generator with current time as seed.
// Returns the seed data and the (possibly replicated) developer list.
func LoadSeedWithReplication(path string, developerCount int, rng *rand.Rand) (*SeedData, []Developer, error) {
	// Load and validate the seed file
	seed, err := LoadSeed(path)
	if err != nil {
		return nil, nil, err
	}

	// If no developer count specified, return original developers
	if developerCount == 0 {
		return seed, seed.Developers, nil
	}

	// Create RNG if not provided
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	// Replicate developers to target count
	developers, err := ReplicateDevelopers(seed, developerCount, rng)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to replicate developers: %w", err)
	}

	return seed, developers, nil
}

// Validate performs comprehensive validation on SeedData.
// Returns a descriptive error if any validation rule fails.
func (s *SeedData) Validate() error {
	// Check that we have at least one developer
	if len(s.Developers) == 0 {
		return fmt.Errorf("validation failed: must have at least one developer")
	}

	// Track unique user IDs and emails
	userIDs := make(map[string]int)
	emails := make(map[string]int)

	// Validate each developer
	for i, dev := range s.Developers {
		if err := validateDeveloper(&dev, i); err != nil {
			return err
		}

		// Check for duplicate user_id
		if prevIdx, exists := userIDs[dev.UserID]; exists {
			return fmt.Errorf("validation failed: duplicate user_id %q at developers[%d] (previously seen at developers[%d])",
				dev.UserID, i, prevIdx)
		}
		userIDs[dev.UserID] = i

		// Check for duplicate email
		if prevIdx, exists := emails[dev.Email]; exists {
			return fmt.Errorf("validation failed: duplicate email %q at developers[%d] (previously seen at developers[%d])",
				dev.Email, i, prevIdx)
		}
		emails[dev.Email] = i
	}

	return nil
}

// validateDeveloper validates a single Developer struct.
func validateDeveloper(dev *Developer, index int) error {
	// Validate user_id
	if dev.UserID == "" {
		return fmt.Errorf("validation failed: developers[%d]: user_id is required", index)
	}
	if !strings.HasPrefix(dev.UserID, "user_") {
		return fmt.Errorf("validation failed: developers[%d]: user_id must start with 'user_', got %q", index, dev.UserID)
	}

	// Validate email
	if dev.Email == "" {
		return fmt.Errorf("validation failed: developers[%d]: email is required", index)
	}
	if !isValidEmail(dev.Email) {
		return fmt.Errorf("validation failed: developers[%d]: invalid email format %q", index, dev.Email)
	}

	// Validate acceptance_rate
	if dev.AcceptanceRate < 0 || dev.AcceptanceRate > 1 {
		return fmt.Errorf("validation failed: developers[%d]: acceptance_rate must be between 0 and 1, got %f",
			index, dev.AcceptanceRate)
	}

	// Validate seniority
	if dev.Seniority == "" {
		return fmt.Errorf("validation failed: developers[%d]: seniority is required", index)
	}
	if !isValidSeniority(dev.Seniority) {
		return fmt.Errorf("validation failed: developers[%d]: invalid seniority %q (must be 'junior', 'mid', or 'senior')",
			index, dev.Seniority)
	}

	return nil
}

// isValidEmail performs basic email format validation.
func isValidEmail(email string) bool {
	// Basic validation: must contain @ and have text before and after it
	if email == "" {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	// Check that both parts have content
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}

	// Domain must contain at least one dot
	if !strings.Contains(parts[1], ".") {
		return false
	}

	// No spaces allowed
	if strings.Contains(email, " ") {
		return false
	}

	return true
}

// isValidSeniority checks if the seniority level is valid.
func isValidSeniority(seniority string) bool {
	validLevels := map[string]bool{
		"junior": true,
		"mid":    true,
		"senior": true,
	}
	return validLevels[seniority]
}
