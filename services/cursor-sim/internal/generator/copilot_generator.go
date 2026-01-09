package generator

import (
	mathrand "math/rand"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// CopilotGenerator generates Microsoft 365 Copilot usage data based on seed configuration.
type CopilotGenerator struct {
	seedData *seed.SeedData
	rng      *mathrand.Rand
}

// CopilotConfig defines app adoption rates for Copilot usage generation.
type CopilotConfig struct {
	// AppAdoptionRates defines the probability each app has activity (0.0-1.0)
	AppAdoptionRates map[models.CopilotApp]float64
}

// NewCopilotGenerator creates a new Copilot generator with a random seed.
func NewCopilotGenerator(seedData *seed.SeedData) *CopilotGenerator {
	return NewCopilotGeneratorWithSeed(seedData, time.Now().UnixNano())
}

// NewCopilotGeneratorWithSeed creates a new Copilot generator with a specific seed for reproducibility.
func NewCopilotGeneratorWithSeed(seedData *seed.SeedData, randSeed int64) *CopilotGenerator {
	return &CopilotGenerator{
		seedData: seedData,
		rng:      mathrand.New(mathrand.NewSource(randSeed)),
	}
}

// DefaultCopilotConfig returns the default configuration for Copilot usage generation.
// Adoption rates based on typical enterprise Microsoft 365 Copilot usage patterns.
func DefaultCopilotConfig() CopilotConfig {
	return CopilotConfig{
		AppAdoptionRates: map[models.CopilotApp]float64{
			models.CopilotAppTeams:      0.85, // 85% - Most popular (meetings, chat)
			models.CopilotAppWord:       0.70, // 70% - Document drafting
			models.CopilotAppOutlook:    0.65, // 65% - Email composition
			models.CopilotAppPowerPoint: 0.50, // 50% - Presentation creation
			models.CopilotAppExcel:      0.40, // 40% - Data analysis
			models.CopilotAppLoop:       0.20, // 20% - Collaborative workspaces (newer)
			models.CopilotAppCopilot:    0.75, // 75% - Copilot Chat standalone
			models.CopilotAppOneNote:    0.10, // 10% - Note-taking (lowest)
		},
	}
}

// GenerateUsageReport generates a Copilot usage report for the specified period.
// Uses configured adoption rates to determine which apps each user has activity in.
func (g *CopilotGenerator) GenerateUsageReport(period models.CopilotReportPeriod) []models.CopilotUsageUserDetail {
	return g.GenerateUsageReportWithConfig(period, DefaultCopilotConfig())
}

// GenerateUsageReportWithConfig generates a Copilot usage report with custom configuration.
func (g *CopilotGenerator) GenerateUsageReportWithConfig(period models.CopilotReportPeriod, config CopilotConfig) []models.CopilotUsageUserDetail {
	var report []models.CopilotUsageUserDetail

	// If Copilot not enabled or no developers, return empty
	if g.seedData.ExternalDataSources == nil ||
		g.seedData.ExternalDataSources.Copilot == nil ||
		!g.seedData.ExternalDataSources.Copilot.Enabled ||
		len(g.seedData.Developers) == 0 {
		return report
	}

	// Report refresh date is "today"
	reportRefreshDate := time.Now().Format("2006-01-02")

	// Period end is today, start is period.Days() ago
	periodEnd := time.Now()
	periodStart := periodEnd.AddDate(0, 0, -period.Days())

	// Generate usage for each developer
	for _, dev := range g.seedData.Developers {
		detail := models.CopilotUsageUserDetail{
			ReportRefreshDate: reportRefreshDate,
			ReportPeriod:      period.Days(),
			UserPrincipalName: dev.Email,
			DisplayName:       dev.Name,
		}

		// Track the latest activity date across all apps
		var latestActivityDate time.Time

		// For each app, determine if this user has activity based on adoption rate
		for _, app := range models.AllCopilotApps() {
			adoptionRate := config.AppAdoptionRates[app]
			if g.rng.Float64() < adoptionRate {
				// User has activity in this app - generate a random date within period
				activityDate := g.randomDateInRange(periodStart, periodEnd)
				activityDateStr := activityDate.Format("2006-01-02")

				// Set the app-specific activity date
				g.setAppActivityDate(&detail, app, &activityDateStr)

				// Track latest date
				if activityDate.After(latestActivityDate) {
					latestActivityDate = activityDate
				}
			}
		}

		// Set LastActivityDate to the max of all app activity dates
		if !latestActivityDate.IsZero() {
			lastActivityDateStr := latestActivityDate.Format("2006-01-02")
			detail.LastActivityDate = &lastActivityDateStr
		}

		report = append(report, detail)
	}

	return report
}

// randomDateInRange generates a random date within the specified range.
func (g *CopilotGenerator) randomDateInRange(start, end time.Time) time.Time {
	duration := end.Sub(start)
	randomDuration := time.Duration(g.rng.Int63n(int64(duration)))
	return start.Add(randomDuration)
}

// setAppActivityDate sets the activity date for the specified app on the detail record.
func (g *CopilotGenerator) setAppActivityDate(detail *models.CopilotUsageUserDetail, app models.CopilotApp, date *string) {
	switch app {
	case models.CopilotAppTeams:
		detail.MicrosoftTeamsCopilotLastActivityDate = date
	case models.CopilotAppWord:
		detail.WordCopilotLastActivityDate = date
	case models.CopilotAppExcel:
		detail.ExcelCopilotLastActivityDate = date
	case models.CopilotAppPowerPoint:
		detail.PowerPointCopilotLastActivityDate = date
	case models.CopilotAppOutlook:
		detail.OutlookCopilotLastActivityDate = date
	case models.CopilotAppOneNote:
		detail.OneNoteCopilotLastActivityDate = date
	case models.CopilotAppLoop:
		detail.LoopCopilotLastActivityDate = date
	case models.CopilotAppCopilot:
		detail.CopilotChatLastActivityDate = date
	}
}
