package generator

import (
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopilotGenerator_GenerateUsageReport(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{
				Email: "jane@company.com",
				Name:  "Jane Developer",
			},
			{
				Email: "bob@company.com",
				Name:  "Bob Manager",
			},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled:       true,
				TotalLicenses: 2,
				ActiveUsers:   2,
			},
		},
	}

	gen := NewCopilotGeneratorWithSeed(seedData, 12345)
	usage := gen.GenerateUsageReport(models.CopilotPeriodD30)

	require.Len(t, usage, 2)
	assert.Equal(t, 30, usage[0].ReportPeriod)
	assert.Equal(t, "jane@company.com", usage[0].UserPrincipalName)
	assert.Equal(t, "Jane Developer", usage[0].DisplayName)
	assert.Equal(t, "bob@company.com", usage[1].UserPrincipalName)
	assert.Equal(t, "Bob Manager", usage[1].DisplayName)

	// Verify report refresh date is set
	assert.NotEmpty(t, usage[0].ReportRefreshDate)
	assert.NotEmpty(t, usage[1].ReportRefreshDate)
}

func TestCopilotGenerator_AppAdoption(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "user@company.com", Name: "Test User"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled:       true,
				TotalLicenses: 1,
				ActiveUsers:   1,
			},
		},
	}

	// Generate many samples to verify adoption rates
	teamsSeen := 0
	wordSeen := 0
	excelSeen := 0
	oneNoteSeen := 0
	samples := 100

	for i := 0; i < samples; i++ {
		gen := NewCopilotGeneratorWithSeed(seedData, int64(i))
		usage := gen.GenerateUsageReport(models.CopilotPeriodD30)

		if usage[0].MicrosoftTeamsCopilotLastActivityDate != nil {
			teamsSeen++
		}
		if usage[0].WordCopilotLastActivityDate != nil {
			wordSeen++
		}
		if usage[0].ExcelCopilotLastActivityDate != nil {
			excelSeen++
		}
		if usage[0].OneNoteCopilotLastActivityDate != nil {
			oneNoteSeen++
		}
	}

	// Teams should have ~85% adoption
	assert.True(t, teamsSeen > 70, "Teams adoption too low: %d", teamsSeen)

	// Word should have ~70% adoption
	assert.True(t, wordSeen > 55, "Word adoption too low: %d", wordSeen)
	assert.True(t, wordSeen < 85, "Word adoption too high: %d", wordSeen)

	// Excel should have ~40% adoption
	assert.True(t, excelSeen > 25, "Excel adoption too low: %d", excelSeen)
	assert.True(t, excelSeen < 55, "Excel adoption too high: %d", excelSeen)

	// OneNote should have ~10% adoption
	assert.True(t, oneNoteSeen < 25, "OneNote adoption too high: %d", oneNoteSeen)
}

func TestCopilotGenerator_ActivityDatesWithinPeriod(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "user@company.com", Name: "Test User"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled:       true,
				TotalLicenses: 1,
				ActiveUsers:   1,
			},
		},
	}

	gen := NewCopilotGeneratorWithSeed(seedData, 12345)
	usage := gen.GenerateUsageReport(models.CopilotPeriodD30)

	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	for _, app := range models.AllCopilotApps() {
		dateStr := usage[0].GetAppLastActivityDate(app)
		if dateStr != nil {
			activityDate, err := time.Parse("2006-01-02", *dateStr)
			require.NoError(t, err, "Failed to parse date for app %s: %s", app, *dateStr)
			assert.True(t, activityDate.After(thirtyDaysAgo) || activityDate.Equal(thirtyDaysAgo),
				"Activity date for %s too old: %s", app, *dateStr)
			assert.True(t, activityDate.Before(now) || activityDate.Equal(now),
				"Activity date for %s in future: %s", app, *dateStr)
		}
	}
}

func TestCopilotGenerator_Reproducible(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "user1@company.com", Name: "User One"},
			{Email: "user2@company.com", Name: "User Two"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled:       true,
				TotalLicenses: 2,
				ActiveUsers:   2,
			},
		},
	}

	// Same seed should produce identical results
	gen1 := NewCopilotGeneratorWithSeed(seedData, 12345)
	usage1 := gen1.GenerateUsageReport(models.CopilotPeriodD30)

	gen2 := NewCopilotGeneratorWithSeed(seedData, 12345)
	usage2 := gen2.GenerateUsageReport(models.CopilotPeriodD30)

	require.Equal(t, len(usage1), len(usage2))
	for i := range usage1 {
		assert.Equal(t, usage1[i].UserPrincipalName, usage2[i].UserPrincipalName)
		assert.Equal(t, usage1[i].MicrosoftTeamsCopilotLastActivityDate, usage2[i].MicrosoftTeamsCopilotLastActivityDate)
		assert.Equal(t, usage1[i].WordCopilotLastActivityDate, usage2[i].WordCopilotLastActivityDate)
		assert.Equal(t, usage1[i].LastActivityDate, usage2[i].LastActivityDate)
	}
}

func TestCopilotGenerator_AllPeriods(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "user@company.com", Name: "Test User"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled:       true,
				TotalLicenses: 1,
				ActiveUsers:   1,
			},
		},
	}

	gen := NewCopilotGeneratorWithSeed(seedData, 12345)

	periods := []models.CopilotReportPeriod{
		models.CopilotPeriodD7,
		models.CopilotPeriodD30,
		models.CopilotPeriodD90,
		models.CopilotPeriodD180,
	}

	for _, period := range periods {
		t.Run(string(period), func(t *testing.T) {
			usage := gen.GenerateUsageReport(period)
			require.Len(t, usage, 1)
			assert.Equal(t, period.Days(), usage[0].ReportPeriod)

			// Verify dates within period
			if usage[0].LastActivityDate != nil {
				activityDate, err := time.Parse("2006-01-02", *usage[0].LastActivityDate)
				require.NoError(t, err)

				now := time.Now()
				periodStart := now.AddDate(0, 0, -period.Days())

				assert.True(t, activityDate.After(periodStart) || activityDate.Equal(periodStart),
					"Activity date %s before period start %s", activityDate, periodStart)
				assert.True(t, activityDate.Before(now) || activityDate.Equal(now),
					"Activity date %s after now %s", activityDate, now)
			}
		})
	}
}

func TestCopilotGenerator_LastActivityDateComputed(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "user@company.com", Name: "Test User"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled:       true,
				TotalLicenses: 1,
				ActiveUsers:   1,
			},
		},
	}

	// Use multiple seeds to ensure we get some app activities
	for i := 0; i < 20; i++ {
		gen := NewCopilotGeneratorWithSeed(seedData, int64(i))
		usage := gen.GenerateUsageReport(models.CopilotPeriodD30)

		if usage[0].LastActivityDate != nil {
			_, err := time.Parse("2006-01-02", *usage[0].LastActivityDate)
			require.NoError(t, err)

			// Find the max date from all app-specific dates
			var maxAppDate time.Time
			for _, app := range models.AllCopilotApps() {
				appDateStr := usage[0].GetAppLastActivityDate(app)
				if appDateStr != nil {
					appDate, err := time.Parse("2006-01-02", *appDateStr)
					require.NoError(t, err)
					if appDate.After(maxAppDate) {
						maxAppDate = appDate
					}
				}
			}

			// LastActivityDate should equal the max of all app dates
			if !maxAppDate.IsZero() {
				expectedDate := maxAppDate.Format("2006-01-02")
				assert.Equal(t, expectedDate, *usage[0].LastActivityDate,
					"LastActivityDate should be max of all app dates")
			}
			break // Found a valid case, exit
		}
	}
}

func TestCopilotGenerator_NoUsersReturnsEmpty(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled:       true,
				TotalLicenses: 0,
				ActiveUsers:   0,
			},
		},
	}

	gen := NewCopilotGeneratorWithSeed(seedData, 12345)
	usage := gen.GenerateUsageReport(models.CopilotPeriodD30)

	assert.Empty(t, usage)
}

func TestCopilotGenerator_DisabledReturnsEmpty(t *testing.T) {
	seedData := &seed.SeedData{
		Developers: []seed.Developer{
			{Email: "user@company.com", Name: "Test User"},
		},
		ExternalDataSources: &seed.ExternalDataSourcesSeed{
			Copilot: &seed.CopilotSeedConfig{
				Enabled: false,
			},
		},
	}

	gen := NewCopilotGeneratorWithSeed(seedData, 12345)
	usage := gen.GenerateUsageReport(models.CopilotPeriodD30)

	assert.Empty(t, usage)
}
