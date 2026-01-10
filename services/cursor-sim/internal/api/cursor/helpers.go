package cursor

import (
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// extractUniqueTeams extracts unique team names from seed data.
func extractUniqueTeams(sd *seed.SeedData) []string {
	teamMap := make(map[string]bool)
	for _, dev := range sd.Developers {
		if dev.Team != "" {
			teamMap[dev.Team] = true
		}
	}

	teams := make([]string, 0, len(teamMap))
	for team := range teamMap {
		teams = append(teams, team)
	}
	return teams
}

// extractUniqueDivisions extracts unique division names from seed data.
func extractUniqueDivisions(sd *seed.SeedData) []string {
	divisionMap := make(map[string]bool)
	for _, dev := range sd.Developers {
		if dev.Division != "" {
			divisionMap[dev.Division] = true
		}
	}

	divisions := make([]string, 0, len(divisionMap))
	for division := range divisionMap {
		divisions = append(divisions, division)
	}
	return divisions
}

// extractUniqueOrgs extracts unique organization names from seed data.
func extractUniqueOrgs(sd *seed.SeedData) []string {
	orgMap := make(map[string]bool)
	for _, dev := range sd.Developers {
		if dev.Org != "" {
			orgMap[dev.Org] = true
		}
	}

	orgs := make([]string, 0, len(orgMap))
	for org := range orgMap {
		orgs = append(orgs, org)
	}
	return orgs
}

// extractUniqueRegions extracts unique region names from seed data.
func extractUniqueRegions(sd *seed.SeedData) []string {
	regionMap := make(map[string]bool)
	for _, dev := range sd.Developers {
		if dev.Region != "" {
			regionMap[dev.Region] = true
		}
	}

	regions := make([]string, 0, len(regionMap))
	for region := range regionMap {
		regions = append(regions, region)
	}
	return regions
}

// extractUniqueRepos extracts unique repository names from seed data.
func extractUniqueRepos(sd *seed.SeedData) []string {
	repos := make([]string, 0, len(sd.Repositories))
	for _, repo := range sd.Repositories {
		if repo.RepoName != "" {
			repos = append(repos, repo.RepoName)
		}
	}
	return repos
}

// groupBySeniority groups developers by seniority level.
func groupBySeniority(sd *seed.SeedData) map[string]int {
	result := make(map[string]int)
	for _, dev := range sd.Developers {
		result[dev.Seniority]++
	}
	return result
}

// groupByRegion groups developers by region.
func groupByRegion(sd *seed.SeedData) map[string]int {
	result := make(map[string]int)
	for _, dev := range sd.Developers {
		result[dev.Region]++
	}
	return result
}

// groupByTeam groups developers by team.
func groupByTeam(sd *seed.SeedData) map[string]int {
	result := make(map[string]int)
	for _, dev := range sd.Developers {
		result[dev.Team]++
	}
	return result
}

// groupByActivity groups developers by activity level.
func groupByActivity(sd *seed.SeedData) map[string]int {
	result := make(map[string]int)
	for _, dev := range sd.Developers {
		result[dev.ActivityLevel]++
	}
	return result
}
