package seed

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeveloper_JSONRoundtrip(t *testing.T) {
	dev := Developer{
		UserID:         "user_abc123xyz789",
		Email:          "jane@example.com",
		Name:           "Jane Developer",
		Org:            "acme-corp",
		Division:       "Engineering",
		Team:           "Backend",
		Role:           "member",
		Region:         "US",
		Timezone:       "America/New_York",
		Locale:         "en-US",
		Seniority:      "senior",
		ActivityLevel:  "high",
		AcceptanceRate: 0.85,
		PRBehavior: PRBehavior{
			PRsPerWeek:         5.0,
			AvgPRSizeLOC:       200,
			AvgFilesPerPR:      6,
			ReviewThoroughness: 0.85,
			IterationTolerance: 2,
		},
		CodingSpeed: CodingSpeed{
			Mean: 2.0,
			Std:  1.0,
		},
		PreferredModels:  []string{"gpt-4-turbo", "claude-3-sonnet"},
		ChatVsCodeRatio:  ChatCodeRatio{Chat: 0.15, Code: 0.85},
		WorkingHoursBand: WorkingHours{Start: 9, End: 18, Peak: 14},
	}

	data, err := json.Marshal(dev)
	require.NoError(t, err)

	var parsed Developer
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, dev.UserID, parsed.UserID)
	assert.Equal(t, dev.Email, parsed.Email)
	assert.Equal(t, dev.AcceptanceRate, parsed.AcceptanceRate)
	assert.Equal(t, dev.PRBehavior.PRsPerWeek, parsed.PRBehavior.PRsPerWeek)
	assert.Equal(t, dev.PreferredModels, parsed.PreferredModels)
}

func TestDeveloper_JSONTags(t *testing.T) {
	jsonStr := `{
		"user_id": "user_001",
		"email": "dev@example.com",
		"name": "Test Dev",
		"org": "acme-corp",
		"division": "Engineering",
		"team": "Backend",
		"role": "member",
		"region": "US",
		"timezone": "America/New_York",
		"locale": "en-US",
		"seniority": "senior",
		"activity_level": "high",
		"acceptance_rate": 0.85,
		"pr_behavior": {
			"prs_per_week": 5.0,
			"avg_pr_size_loc": 200,
			"avg_files_per_pr": 6,
			"review_thoroughness": 0.85,
			"iteration_tolerance": 2
		},
		"coding_speed": {
			"mean": 2.0,
			"std": 1.0
		},
		"preferred_models": ["gpt-4-turbo"],
		"chat_vs_code_ratio": {
			"chat": 0.15,
			"code": 0.85
		},
		"working_hours_band": {
			"start": 9,
			"end": 18,
			"peak": 14
		}
	}`

	var dev Developer
	err := json.Unmarshal([]byte(jsonStr), &dev)
	require.NoError(t, err)

	assert.Equal(t, "user_001", dev.UserID)
	assert.Equal(t, "dev@example.com", dev.Email)
	assert.Equal(t, 0.85, dev.AcceptanceRate)
	assert.Equal(t, 5.0, dev.PRBehavior.PRsPerWeek)
	assert.Equal(t, 200, dev.PRBehavior.AvgPRSizeLOC)
	assert.Equal(t, "high", dev.ActivityLevel)
	assert.Equal(t, []string{"gpt-4-turbo"}, dev.PreferredModels)
}

func TestRepository_JSONRoundtrip(t *testing.T) {
	repo := Repository{
		RepoName:        "acme-corp/payment-service",
		PrimaryLanguage: "Go",
		ServiceType:     "api",
		DefaultBranch:   "main",
		Teams:           []string{"Backend", "DevOps"},
		Maturity: Maturity{
			AgeDays:           400,
			TotalCommits:      600,
			TotalPRs:          250,
			TotalContributors: 12,
			SizeBytes:         500000,
		},
		CodeQualityBaseline: CodeQualityBaseline{
			AvgFileAgeDays:      180,
			GreenfieldFileRatio: 0.15,
			RevertRateBaseline:  0.02,
			HotfixRateBaseline:  0.08,
		},
		CommonFilePatterns: []string{"**/*.go", "cmd/**/*.go"},
	}

	data, err := json.Marshal(repo)
	require.NoError(t, err)

	var parsed Repository
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, repo.RepoName, parsed.RepoName)
	assert.Equal(t, repo.Teams, parsed.Teams)
	assert.Equal(t, repo.Maturity.AgeDays, parsed.Maturity.AgeDays)
	assert.Equal(t, repo.CodeQualityBaseline.RevertRateBaseline, parsed.CodeQualityBaseline.RevertRateBaseline)
}

func TestRepository_JSONTags(t *testing.T) {
	jsonStr := `{
		"repo_name": "acme-corp/api-gateway",
		"primary_language": "Python",
		"service_type": "api",
		"default_branch": "main",
		"teams": ["Backend"],
		"maturity": {
			"age_days": 200,
			"total_commits": 300,
			"total_prs": 100,
			"total_contributors": 5,
			"size_bytes": 250000
		},
		"code_quality_baseline": {
			"avg_file_age_days": 100,
			"greenfield_file_ratio": 0.25,
			"revert_rate_baseline": 0.03,
			"hotfix_rate_baseline": 0.10
		},
		"common_file_patterns": ["src/**/*.py"]
	}`

	var repo Repository
	err := json.Unmarshal([]byte(jsonStr), &repo)
	require.NoError(t, err)

	assert.Equal(t, "acme-corp/api-gateway", repo.RepoName)
	assert.Equal(t, "Python", repo.PrimaryLanguage)
	assert.Equal(t, 200, repo.Maturity.AgeDays)
	assert.Equal(t, 0.25, repo.CodeQualityBaseline.GreenfieldFileRatio)
}

func TestTextTemplates_JSONRoundtrip(t *testing.T) {
	templates := TextTemplates{
		CommitMessages: CommitMessageTemplates{
			Feature:  []string{"Add {{ feature_name }} to {{ component }}"},
			Bugfix:   []string{"Fix {{ issue }} in {{ component }}"},
			Refactor: []string{"Refactor {{ component }}"},
			Chore:    []string{"Update dependencies"},
		},
		PRTitles:       []string{"feat: add new feature"},
		PRDescriptions: []string{"## Summary\n\nAdded feature X"},
		ReviewComments: ReviewCommentTemplates{
			Style:      []string{"Consider using const here"},
			Logic:      []string{"What happens if X is null?"},
			Suggestion: []string{"Could we extract this?"},
			Approval:   []string{"LGTM!"},
		},
		ChatPromptThemes: ChatPromptThemes{
			CodeGeneration: []string{"implement function"},
			Debugging:      []string{"fix error"},
			Refactoring:    []string{"improve code"},
			Explanation:    []string{"explain code"},
			Learning:       []string{"how to"},
		},
	}

	data, err := json.Marshal(templates)
	require.NoError(t, err)

	var parsed TextTemplates
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, templates.CommitMessages.Feature, parsed.CommitMessages.Feature)
	assert.Equal(t, templates.ReviewComments.Approval, parsed.ReviewComments.Approval)
}

func TestCorrelations_JSONRoundtrip(t *testing.T) {
	corr := Correlations{
		SeniorityToBehavior: map[string]SeniorityBehavior{
			"senior": {
				EventsPerDay:       StatParams{Mean: 150, Std: 50},
				TabToComposerRatio: 0.65,
				AcceptsPerSession:  StatParams{Mean: 15, Std: 6},
			},
		},
		RegionToActivity: map[string]RegionActivity{
			"US": {
				WeekdayWeight: 1.0,
				WeekendWeight: 0.15,
				PeakHours:     []int{10, 11, 14, 15, 16},
			},
		},
		LinesPerChange: map[string]LineChangeParams{
			"TAB": {
				Added:   StatParamsMax{Mean: 3, Std: 2, Max: 25},
				Deleted: StatParamsMax{Mean: 0.5, Std: 0.8, Max: 10},
			},
		},
		AIRatioBands: AIRatioBands{
			Low:    AIRatioBand{Max: 0.3},
			Medium: AIRatioBand{Min: 0.3, Max: 0.6},
			High:   AIRatioBand{Min: 0.6},
		},
	}

	data, err := json.Marshal(corr)
	require.NoError(t, err)

	var parsed Correlations
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, corr.SeniorityToBehavior["senior"].TabToComposerRatio, parsed.SeniorityToBehavior["senior"].TabToComposerRatio)
	assert.Equal(t, corr.RegionToActivity["US"].PeakHours, parsed.RegionToActivity["US"].PeakHours)
}

func TestPRLifecycle_JSONRoundtrip(t *testing.T) {
	lifecycle := PRLifecycle{
		CycleTimes: CycleTimes{
			CodingLeadTime: TimeDistribution{
				BaseDistribution: "lognormal",
				Params:           StatParams{Mean: 4, Std: 3},
				Modifiers: TimeModifiers{
					BySeniority: map[string]float64{"senior": 0.7, "junior": 1.5},
					ByPRSize:    map[string]float64{"small": 0.5, "large": 1.8},
				},
			},
			PickupTime: TimeDistribution{
				BaseDistribution: "lognormal",
				Params:           StatParams{Mean: 6, Std: 4},
			},
			ReviewLeadTime: TimeDistribution{
				BaseDistribution: "lognormal",
				Params:           StatParams{Mean: 8, Std: 6},
			},
		},
		ReviewPatterns: ReviewPatterns{
			CommentsPer100LOC: CommentDensity{
				Base: 2.5,
				Modifiers: CommentModifiers{
					ByReviewerSeniority: map[string]float64{"senior": 3.0},
					ByAIRatio:           map[string]float64{"high": 1.4},
				},
			},
			Iterations: IterationParams{
				BaseDistribution: "poisson",
				Params:           map[string]float64{"lambda": 1.5},
			},
			ReviewerCount: ReviewerCountParams{
				Base: 1.8,
			},
		},
		QualityOutcomes: QualityOutcomes{
			RevertProbability: OutcomeParams{
				Base: 0.02,
				Modifiers: OutcomeModifiers{
					ByAuthorSeniority:  map[string]float64{"junior": 2.0},
					ByAIRatio:          map[string]float64{"high": 1.5},
					ByReviewIterations: map[string]float64{"0": 3.0},
				},
			},
			HotfixProbability: OutcomeParams{
				Base: 0.08,
			},
			CodeSurvival30D: OutcomeParams{
				Base: 0.85,
			},
		},
		ScopeCreep: ScopeCreepParams{
			BaseRatio: 0.12,
		},
		ReworkRatio: ReworkRatioParams{
			BaseRatio: 0.15,
		},
	}

	data, err := json.Marshal(lifecycle)
	require.NoError(t, err)

	var parsed PRLifecycle
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, lifecycle.CycleTimes.CodingLeadTime.Modifiers.BySeniority["senior"],
		parsed.CycleTimes.CodingLeadTime.Modifiers.BySeniority["senior"])
	assert.Equal(t, lifecycle.QualityOutcomes.RevertProbability.Base,
		parsed.QualityOutcomes.RevertProbability.Base)
}

func TestSeedData_FullStructure(t *testing.T) {
	seed := SeedData{
		Version: "1.0",
		Developers: []Developer{
			{
				UserID:         "user_001",
				Email:          "dev@example.com",
				Name:           "Test Developer",
				Org:            "acme-corp",
				Seniority:      "senior",
				AcceptanceRate: 0.85,
			},
		},
		Repositories: []Repository{
			{
				RepoName:        "acme-corp/api",
				PrimaryLanguage: "Go",
				DefaultBranch:   "main",
			},
		},
		TextTemplates: TextTemplates{
			CommitMessages: CommitMessageTemplates{
				Feature: []string{"Add feature"},
			},
		},
		Correlations: Correlations{
			SeniorityToBehavior: map[string]SeniorityBehavior{},
		},
		PRLifecycle: PRLifecycle{
			CycleTimes: CycleTimes{},
		},
	}

	data, err := json.Marshal(seed)
	require.NoError(t, err)

	var parsed SeedData
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "1.0", parsed.Version)
	assert.Len(t, parsed.Developers, 1)
	assert.Equal(t, "user_001", parsed.Developers[0].UserID)
	assert.Len(t, parsed.Repositories, 1)
}

func TestPRBehavior_Defaults(t *testing.T) {
	// Test that zero values are handled correctly
	var pr PRBehavior
	data, err := json.Marshal(pr)
	require.NoError(t, err)

	var parsed PRBehavior
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, 0.0, parsed.PRsPerWeek)
	assert.Equal(t, 0, parsed.AvgPRSizeLOC)
}

func TestMaturity_MapFields(t *testing.T) {
	jsonStr := `{
		"age_days": 365,
		"total_commits": 500,
		"total_prs": 200,
		"total_contributors": 10,
		"size_bytes": 1000000
	}`

	var m Maturity
	err := json.Unmarshal([]byte(jsonStr), &m)
	require.NoError(t, err)

	assert.Equal(t, 365, m.AgeDays)
	assert.Equal(t, 500, m.TotalCommits)
	assert.Equal(t, 200, m.TotalPRs)
	assert.Equal(t, 10, m.TotalContributors)
	assert.Equal(t, 1000000, m.SizeBytes)
}
