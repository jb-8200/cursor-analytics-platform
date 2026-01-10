package seed

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadSeed reads and validates a seed file (JSON or YAML) from the specified path.
// The file format is detected by extension: .json, .yaml, or .yml
// Returns an error if the file cannot be read, parsed, or if validation fails.
func LoadSeed(path string) (*SeedData, error) {
	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read seed file %q: %w", path, err)
	}

	// Detect format by extension
	var seed SeedData
	ext := filepath.Ext(path)

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &seed); err != nil {
			return nil, fmt.Errorf("failed to parse JSON seed file %q: %w", path, err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &seed); err != nil {
			return nil, fmt.Errorf("failed to parse YAML seed file %q: %w", path, err)
		}
	default:
		return nil, fmt.Errorf("unsupported seed file format: %s (use .json, .yaml, or .yml)", ext)
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

// LoadFromCSV reads a CSV file and creates a minimal SeedData structure.
// CSV format: user_id, email, name
// Returns a basic seed with developers only, using default values for other fields.
func LoadFromCSV(reader io.Reader) (*SeedData, error) {
	csvReader := csv.NewReader(reader)

	// Read header row
	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate headers (must contain at least user_id, email, name)
	if len(headers) < 3 {
		return nil, fmt.Errorf("CSV must have at least 3 columns: user_id, email, name")
	}

	// Find column indices
	var userIDIdx, emailIdx, nameIdx int = -1, -1, -1
	for i, header := range headers {
		switch strings.ToLower(strings.TrimSpace(header)) {
		case "user_id":
			userIDIdx = i
		case "email":
			emailIdx = i
		case "name":
			nameIdx = i
		}
	}

	if userIDIdx == -1 || emailIdx == -1 || nameIdx == -1 {
		return nil, fmt.Errorf("CSV must have columns: user_id, email, name")
	}

	// Read all developer rows
	developers := make([]Developer, 0)
	rowNum := 1 // Start at 1 (header is row 0)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row %d: %w", rowNum+1, err)
		}

		if len(row) < 3 {
			return nil, fmt.Errorf("CSV row %d has fewer than 3 columns", rowNum+1)
		}

		// Create developer with minimal fields
		dev := Developer{
			UserID:         strings.TrimSpace(row[userIDIdx]),
			Email:          strings.TrimSpace(row[emailIdx]),
			Name:           strings.TrimSpace(row[nameIdx]),
			Org:            "default-org",
			Division:       "Engineering",
			Team:           "Development",
			Role:           "developer",
			Region:         "US",
			Timezone:       "America/New_York",
			Locale:         "en-US",
			Seniority:      "mid",
			ActivityLevel:  "medium",
			AcceptanceRate: 0.7,
			PRBehavior: PRBehavior{
				PRsPerWeek:         2.0,
				AvgPRSizeLOC:       200,
				AvgFilesPerPR:      5,
				ReviewThoroughness: 0.7,
				IterationTolerance: 2,
			},
			CodingSpeed: CodingSpeed{
				Mean: 4.0,
				Std:  1.0,
			},
			PreferredModels: []string{"gpt-4"},
			ChatVsCodeRatio: ChatCodeRatio{
				Chat: 0.3,
				Code: 0.7,
			},
			WorkingHoursBand: WorkingHours{
				Start: 9,
				End:   17,
				Peak:  13,
			},
		}

		developers = append(developers, dev)
		rowNum++
	}

	if len(developers) == 0 {
		return nil, fmt.Errorf("CSV contains no developer data")
	}

	// Load the default seed template and replace developers
	// This ensures we get valid Correlations and PRLifecycle structures
	templateSeed, err := LoadSeed("testdata/valid_seed.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load seed template: %w", err)
	}

	// Replace developers with CSV data
	templateSeed.Developers = developers

	// Validate the generated seed
	if err := templateSeed.Validate(); err != nil {
		return nil, fmt.Errorf("CSV seed validation failed: %w", err)
	}

	return templateSeed, nil
}
