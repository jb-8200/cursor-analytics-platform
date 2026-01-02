package generator

import (
	"testing"

	"github.com/jb-8200/cursor-analytics-platform/services/cursor-sim/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDeveloperGenerator(t *testing.T) {
	cfg := &config.Config{
		Developers: 10,
		Seed:       12345,
	}

	gen := NewDeveloperGenerator(cfg)
	require.NotNil(t, gen)
	assert.NotNil(t, gen.config)
	assert.NotNil(t, gen.rng)
}

func TestDeveloperGenerator_Generate_Count(t *testing.T) {
	tests := []struct {
		name      string
		devCount  int
		wantCount int
	}{
		{
			name:      "small team",
			devCount:  10,
			wantCount: 10,
		},
		{
			name:      "medium team",
			devCount:  50,
			wantCount: 50,
		},
		{
			name:      "large team",
			devCount:  200,
			wantCount: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Developers: tt.devCount,
				Seed:       12345,
			}

			gen := NewDeveloperGenerator(cfg)
			developers, err := gen.Generate()

			require.NoError(t, err)
			assert.Len(t, developers, tt.wantCount)
		})
	}
}

func TestDeveloperGenerator_Generate_RegionDistribution(t *testing.T) {
	cfg := &config.Config{
		Developers: 1000, // Large sample for statistical accuracy
		Seed:       12345,
	}

	gen := NewDeveloperGenerator(cfg)
	developers, err := gen.Generate()

	require.NoError(t, err)
	require.Len(t, developers, 1000)

	// Count regions
	regionCounts := make(map[string]int)
	for _, dev := range developers {
		regionCounts[dev.Region]++
	}

	// Check distribution (allow ±5% variance)
	assert.InDelta(t, 500, regionCounts["US"], 50, "US distribution should be ~50%")
	assert.InDelta(t, 350, regionCounts["EU"], 35, "EU distribution should be ~35%")
	assert.InDelta(t, 150, regionCounts["APAC"], 15, "APAC distribution should be ~15%")
}

func TestDeveloperGenerator_Generate_DivisionDistribution(t *testing.T) {
	cfg := &config.Config{
		Developers: 1000,
		Seed:       12345,
	}

	gen := NewDeveloperGenerator(cfg)
	developers, err := gen.Generate()

	require.NoError(t, err)

	divisionCounts := make(map[string]int)
	for _, dev := range developers {
		divisionCounts[dev.Division]++
	}

	assert.InDelta(t, 400, divisionCounts["AGS"], 40, "AGS distribution should be ~40%")
	assert.InDelta(t, 350, divisionCounts["AT"], 35, "AT distribution should be ~35%")
	assert.InDelta(t, 250, divisionCounts["ST"], 25, "ST distribution should be ~25%")
}

func TestDeveloperGenerator_Generate_GroupDistribution(t *testing.T) {
	cfg := &config.Config{
		Developers: 1000,
		Seed:       12345,
	}

	gen := NewDeveloperGenerator(cfg)
	developers, err := gen.Generate()

	require.NoError(t, err)

	groupCounts := make(map[string]int)
	for _, dev := range developers {
		groupCounts[dev.Group]++
	}

	assert.InDelta(t, 600, groupCounts["TMOBILE"], 60, "TMOBILE distribution should be ~60%")
	assert.InDelta(t, 400, groupCounts["ATANT"], 40, "ATANT distribution should be ~40%")
}

func TestDeveloperGenerator_Generate_TeamDistribution(t *testing.T) {
	cfg := &config.Config{
		Developers: 1000,
		Seed:       12345,
	}

	gen := NewDeveloperGenerator(cfg)
	developers, err := gen.Generate()

	require.NoError(t, err)

	teamCounts := make(map[string]int)
	for _, dev := range developers {
		teamCounts[dev.Team]++
	}

	assert.InDelta(t, 750, teamCounts["dev"], 75, "dev distribution should be ~75%")
	assert.InDelta(t, 250, teamCounts["support"], 25, "support distribution should be ~25%")
}

func TestDeveloperGenerator_Generate_SeniorityDistribution(t *testing.T) {
	cfg := &config.Config{
		Developers: 1000,
		Seed:       12345,
	}

	gen := NewDeveloperGenerator(cfg)
	developers, err := gen.Generate()

	require.NoError(t, err)

	seniorityCounts := make(map[string]int)
	for _, dev := range developers {
		seniorityCounts[dev.Seniority]++
	}

	// Spec: 20% junior, 50% mid, 30% senior
	assert.InDelta(t, 200, seniorityCounts["junior"], 20, "junior distribution should be ~20%")
	assert.InDelta(t, 500, seniorityCounts["mid"], 50, "mid distribution should be ~50%")
	assert.InDelta(t, 300, seniorityCounts["senior"], 30, "senior distribution should be ~30%")
}

func TestDeveloperGenerator_Generate_AcceptanceRates(t *testing.T) {
	cfg := &config.Config{
		Developers: 300,
		Seed:       12345,
	}

	gen := NewDeveloperGenerator(cfg)
	developers, err := gen.Generate()

	require.NoError(t, err)

	juniorRates := []float64{}
	midRates := []float64{}
	seniorRates := []float64{}

	for _, dev := range developers {
		switch dev.Seniority {
		case "junior":
			juniorRates = append(juniorRates, dev.AcceptanceRate)
		case "mid":
			midRates = append(midRates, dev.AcceptanceRate)
		case "senior":
			seniorRates = append(seniorRates, dev.AcceptanceRate)
		}
	}

	// Check that all acceptance rates are in expected ranges
	// Junior: 55-65%
	for _, rate := range juniorRates {
		assert.GreaterOrEqual(t, rate, 0.55, "junior acceptance rate should be >= 55%")
		assert.LessOrEqual(t, rate, 0.65, "junior acceptance rate should be <= 65%")
	}

	// Mid: 70-80%
	for _, rate := range midRates {
		assert.GreaterOrEqual(t, rate, 0.70, "mid acceptance rate should be >= 70%")
		assert.LessOrEqual(t, rate, 0.80, "mid acceptance rate should be <= 80%")
	}

	// Senior: 85-95%
	for _, rate := range seniorRates {
		assert.GreaterOrEqual(t, rate, 0.85, "senior acceptance rate should be >= 85%")
		assert.LessOrEqual(t, rate, 0.95, "senior acceptance rate should be <= 95%")
	}
}

func TestDeveloperGenerator_Generate_DeterministicNames(t *testing.T) {
	cfg := &config.Config{
		Developers: 50,
		Seed:       99999, // Fixed seed
	}

	// Generate twice with same seed
	gen1 := NewDeveloperGenerator(cfg)
	developers1, err1 := gen1.Generate()
	require.NoError(t, err1)

	gen2 := NewDeveloperGenerator(cfg)
	developers2, err2 := gen2.Generate()
	require.NoError(t, err2)

	// Names should be identical when using same seed
	require.Len(t, developers1, len(developers2))
	for i := range developers1 {
		assert.Equal(t, developers1[i].Name, developers2[i].Name,
			"Names should be deterministic with same seed at index %d", i)
		assert.Equal(t, developers1[i].Email, developers2[i].Email,
			"Emails should be deterministic with same seed at index %d", i)
	}
}

func TestDeveloperGenerator_Generate_UniqueIDs(t *testing.T) {
	cfg := &config.Config{
		Developers: 100,
		Seed:       12345,
	}

	gen := NewDeveloperGenerator(cfg)
	developers, err := gen.Generate()

	require.NoError(t, err)

	// Check all IDs are unique
	idSet := make(map[string]bool)
	for _, dev := range developers {
		assert.False(t, idSet[dev.ID], "Developer ID should be unique: %s", dev.ID)
		idSet[dev.ID] = true
	}
}

func TestDeveloperGenerator_Generate_UniqueEmails(t *testing.T) {
	cfg := &config.Config{
		Developers: 100,
		Seed:       12345,
	}

	gen := NewDeveloperGenerator(cfg)
	developers, err := gen.Generate()

	require.NoError(t, err)

	// Check all emails are unique
	emailSet := make(map[string]bool)
	for _, dev := range developers {
		assert.False(t, emailSet[dev.Email], "Developer email should be unique: %s", dev.Email)
		emailSet[dev.Email] = true
	}
}

func TestDeveloperGenerator_Generate_ClientVersions(t *testing.T) {
	cfg := &config.Config{
		Developers: 100,
		Seed:       12345,
	}

	gen := NewDeveloperGenerator(cfg)
	developers, err := gen.Generate()

	require.NoError(t, err)

	// Check all developers have valid client versions (0.42.x or 0.43.x)
	for _, dev := range developers {
		assert.NotEmpty(t, dev.ClientVersion, "ClientVersion should not be empty")
		// Should start with "0.42." or "0.43."
		assert.Regexp(t, `^0\.(42|43)\.\d+$`, dev.ClientVersion,
			"ClientVersion should match pattern 0.42.x or 0.43.x")
	}
}

func TestDeveloperGenerator_Generate_AllFieldsPopulated(t *testing.T) {
	cfg := &config.Config{
		Developers: 10,
		Seed:       12345,
	}

	gen := NewDeveloperGenerator(cfg)
	developers, err := gen.Generate()

	require.NoError(t, err)

	for i, dev := range developers {
		assert.NotEmpty(t, dev.ID, "ID should not be empty at index %d", i)
		assert.NotEmpty(t, dev.Email, "Email should not be empty at index %d", i)
		assert.NotEmpty(t, dev.Name, "Name should not be empty at index %d", i)
		assert.NotEmpty(t, dev.Region, "Region should not be empty at index %d", i)
		assert.NotEmpty(t, dev.Division, "Division should not be empty at index %d", i)
		assert.NotEmpty(t, dev.Group, "Group should not be empty at index %d", i)
		assert.NotEmpty(t, dev.Team, "Team should not be empty at index %d", i)
		assert.NotEmpty(t, dev.Seniority, "Seniority should not be empty at index %d", i)
		assert.NotEmpty(t, dev.ClientVersion, "ClientVersion should not be empty at index %d", i)
		assert.Greater(t, dev.AcceptanceRate, 0.0, "AcceptanceRate should be > 0 at index %d", i)
		assert.True(t, dev.IsActive, "IsActive should be true at index %d", i)
		assert.False(t, dev.CreatedAt.IsZero(), "CreatedAt should not be zero at index %d", i)
		assert.False(t, dev.LastActiveAt.IsZero(), "LastActiveAt should not be zero at index %d", i)
	}
}
