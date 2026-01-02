package seed

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSeed_ValidFile(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "valid_seed.json")
	seed, err := LoadSeed(path)

	require.NoError(t, err)
	require.NotNil(t, seed)

	assert.Equal(t, "1.0", seed.Version)
	assert.Len(t, seed.Developers, 2)
	assert.Equal(t, "user_001", seed.Developers[0].UserID)
	assert.Equal(t, "alice@example.com", seed.Developers[0].Email)
	assert.Equal(t, 0.85, seed.Developers[0].AcceptanceRate)
	assert.Len(t, seed.Repositories, 2)
}

func TestLoadSeed_FileNotFound(t *testing.T) {
	seed, err := LoadSeed("nonexistent.json")

	require.Error(t, err)
	assert.Nil(t, seed)
	assert.Contains(t, err.Error(), "failed to read seed file")
}

func TestLoadSeed_InvalidJSON(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "invalid_json.json")
	seed, err := LoadSeed(path)

	require.Error(t, err)
	assert.Nil(t, seed)
	assert.Contains(t, err.Error(), "failed to parse seed file")
}

func TestLoadSeed_InvalidUserID(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "invalid_user_id.json")
	seed, err := LoadSeed(path)

	require.Error(t, err)
	assert.Nil(t, seed)
	assert.Contains(t, err.Error(), "user_id must start with 'user_'")
}

func TestLoadSeed_InvalidAcceptanceRate(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "invalid_acceptance_rate.json")
	seed, err := LoadSeed(path)

	require.Error(t, err)
	assert.Nil(t, seed)
	assert.Contains(t, err.Error(), "acceptance_rate must be between 0 and 1")
}

func TestLoadSeed_InvalidSeniority(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "invalid_seniority.json")
	seed, err := LoadSeed(path)

	require.Error(t, err)
	assert.Nil(t, seed)
	assert.Contains(t, err.Error(), "invalid seniority")
}

func TestLoadSeed_InvalidEmail(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "invalid_email.json")
	seed, err := LoadSeed(path)

	require.Error(t, err)
	assert.Nil(t, seed)
	assert.Contains(t, err.Error(), "invalid email format")
}

func TestLoadSeed_EmptyDevelopers(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "empty_developers.json")
	seed, err := LoadSeed(path)

	require.Error(t, err)
	assert.Nil(t, seed)
	assert.Contains(t, err.Error(), "must have at least one developer")
}

func TestValidate_ValidSeed(t *testing.T) {
	seed := &SeedData{
		Version: "1.0",
		Developers: []Developer{
			{
				UserID:         "user_001",
				Email:          "test@example.com",
				Name:           "Test",
				Seniority:      "senior",
				AcceptanceRate: 0.85,
			},
		},
	}

	err := seed.Validate()
	assert.NoError(t, err)
}

func TestValidate_MissingUserID(t *testing.T) {
	seed := &SeedData{
		Version: "1.0",
		Developers: []Developer{
			{
				UserID:         "",
				Email:          "test@example.com",
				AcceptanceRate: 0.85,
			},
		},
	}

	err := seed.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user_id is required")
}

func TestValidate_UserIDFormat(t *testing.T) {
	tests := []struct {
		name   string
		userID string
		valid  bool
	}{
		{"valid user_001", "user_001", true},
		{"valid user_abc123", "user_abc123", true},
		{"invalid prefix", "usr_001", false},
		{"no prefix", "001", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seed := &SeedData{
				Version: "1.0",
				Developers: []Developer{
					{
						UserID:         tt.userID,
						Email:          "test@example.com",
						Seniority:      "senior",
						AcceptanceRate: 0.85,
					},
				},
			}

			err := seed.Validate()
			if tt.valid {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestValidate_AcceptanceRate(t *testing.T) {
	tests := []struct {
		name  string
		rate  float64
		valid bool
	}{
		{"valid 0.0", 0.0, true},
		{"valid 0.5", 0.5, true},
		{"valid 1.0", 1.0, true},
		{"invalid negative", -0.1, false},
		{"invalid > 1", 1.1, false},
		{"invalid large", 2.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seed := &SeedData{
				Version: "1.0",
				Developers: []Developer{
					{
						UserID:         "user_001",
						Email:          "test@example.com",
						Seniority:      "senior",
						AcceptanceRate: tt.rate,
					},
				},
			}

			err := seed.Validate()
			if tt.valid {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "acceptance_rate")
			}
		})
	}
}

func TestValidate_Seniority(t *testing.T) {
	tests := []struct {
		name      string
		seniority string
		valid     bool
	}{
		{"valid junior", "junior", true},
		{"valid mid", "mid", true},
		{"valid senior", "senior", true},
		{"invalid expert", "expert", false},
		{"invalid lead", "lead", false},
		{"invalid empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seed := &SeedData{
				Version: "1.0",
				Developers: []Developer{
					{
						UserID:         "user_001",
						Email:          "test@example.com",
						Seniority:      tt.seniority,
						AcceptanceRate: 0.85,
					},
				},
			}

			err := seed.Validate()
			if tt.valid {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "seniority")
			}
		})
	}
}

func TestValidate_Email(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"valid simple", "test@example.com", true},
		{"valid subdomain", "user@mail.example.com", true},
		{"valid hyphen", "test-user@example.com", true},
		{"invalid no @", "testexample.com", false},
		{"invalid no domain", "test@", false},
		{"invalid no user", "@example.com", false},
		{"invalid spaces", "test @example.com", false},
		{"invalid empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seed := &SeedData{
				Version: "1.0",
				Developers: []Developer{
					{
						UserID:         "user_001",
						Email:          tt.email,
						Seniority:      "senior",
						AcceptanceRate: 0.85,
					},
				},
			}

			err := seed.Validate()
			if tt.valid {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "email")
			}
		})
	}
}

func TestValidate_DuplicateUserIDs(t *testing.T) {
	seed := &SeedData{
		Version: "1.0",
		Developers: []Developer{
			{
				UserID:         "user_001",
				Email:          "alice@example.com",
				Seniority:      "senior",
				AcceptanceRate: 0.85,
			},
			{
				UserID:         "user_001",
				Email:          "bob@example.com",
				Seniority:      "mid",
				AcceptanceRate: 0.75,
			},
		},
	}

	err := seed.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate user_id")
}

func TestValidate_DuplicateEmails(t *testing.T) {
	seed := &SeedData{
		Version: "1.0",
		Developers: []Developer{
			{
				UserID:         "user_001",
				Email:          "alice@example.com",
				Seniority:      "senior",
				AcceptanceRate: 0.85,
			},
			{
				UserID:         "user_002",
				Email:          "alice@example.com",
				Seniority:      "mid",
				AcceptanceRate: 0.75,
			},
		},
	}

	err := seed.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate email")
}

func TestLoadSeed_WithTempFile(t *testing.T) {
	// Create a temporary valid seed file
	tmpFile, err := os.CreateTemp("", "seed_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := `{
		"version": "1.0",
		"developers": [
			{
				"user_id": "user_temp",
				"email": "temp@example.com",
				"name": "Temp User",
				"seniority": "mid",
				"acceptance_rate": 0.8
			}
		],
		"repositories": [],
		"text_templates": {
			"commit_messages": {"feature": [], "bugfix": [], "refactor": [], "chore": []},
			"pr_titles": [],
			"pr_descriptions": [],
			"review_comments": {"style": [], "logic": [], "suggestion": [], "approval": []},
			"chat_prompt_themes": {"code_generation": [], "debugging": [], "refactoring": [], "explanation": [], "learning": []}
		},
		"correlations": {
			"seniority_to_behavior": {},
			"region_to_activity": {},
			"lines_per_change": {},
			"ai_ratio_bands": {"low": {}, "medium": {}, "high": {}}
		},
		"pr_lifecycle": {
			"cycle_times": {
				"coding_lead_time": {"base_distribution": "", "params": {}, "modifiers": {}},
				"pickup_time": {"base_distribution": "", "params": {}, "modifiers": {}},
				"review_lead_time": {"base_distribution": "", "params": {}, "modifiers": {}}
			},
			"review_patterns": {
				"comments_per_100_loc": {"base": 0, "modifiers": {}},
				"iterations": {"base_distribution": "", "params": {}, "modifiers": {}},
				"reviewer_count": {"base": 0, "modifiers": {}}
			},
			"quality_outcomes": {
				"revert_probability": {"base": 0, "modifiers": {}},
				"hotfix_probability": {"base": 0, "modifiers": {}},
				"code_survival_30d": {"base": 0, "modifiers": {}}
			},
			"scope_creep": {"base_ratio": 0, "modifiers": {}},
			"rework_ratio": {"base_ratio": 0, "modifiers": {}}
		}
	}`

	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)
	tmpFile.Close()

	seed, err := LoadSeed(tmpFile.Name())
	require.NoError(t, err)
	require.NotNil(t, seed)
	assert.Equal(t, "user_temp", seed.Developers[0].UserID)
}
