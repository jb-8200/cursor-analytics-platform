package models

// StatsResponse represents comprehensive statistics about the generated simulation data.
type StatsResponse struct {
	Generation   Generation   `json:"generation"`
	Developers   Developers   `json:"developers"`
	Quality      Quality      `json:"quality"`
	Variance     Variance     `json:"variance"`
	Performance  Performance  `json:"performance"`
	Organization Organization `json:"organization"`
	TimeSeries   *TimeSeries  `json:"time_series,omitempty"`
}

// Generation contains overall generation statistics.
type Generation struct {
	TotalCommits    int    `json:"total_commits"`
	TotalPRs        int    `json:"total_prs"`
	TotalReviews    int    `json:"total_reviews"`
	TotalIssues     int    `json:"total_issues"`
	TotalDevelopers int    `json:"total_developers"`
	DataSize        string `json:"data_size"`
}

// Developers contains developer breakdown statistics.
type Developers struct {
	BySeniority map[string]int `json:"by_seniority"`
	ByRegion    map[string]int `json:"by_region"`
	ByTeam      map[string]int `json:"by_team"`
	ByActivity  map[string]int `json:"by_activity"`
}

// Quality contains quality metrics.
type Quality struct {
	AvgRevertRate         float64 `json:"avg_revert_rate"`
	AvgHotfixRate         float64 `json:"avg_hotfix_rate"`
	AvgCodeSurvival       float64 `json:"avg_code_survival_30d"`
	AvgReviewThoroughness float64 `json:"avg_review_thoroughness"`
	AvgIterations         float64 `json:"avg_pr_iterations"`
}

// Variance contains variance metrics (standard deviation).
type Variance struct {
	CommitsStdDev   float64 `json:"commits_std_dev"`
	PRSizeStdDev    float64 `json:"pr_size_std_dev"`
	CycleTimeStdDev float64 `json:"cycle_time_std_dev"`
}

// Performance contains performance metrics.
type Performance struct {
	LastGenerationTime string `json:"last_generation_time"`
	MemoryUsage        string `json:"memory_usage"`
	StorageEfficiency  string `json:"storage_efficiency"`
}

// Organization contains organizational structure information.
type Organization struct {
	Teams        []string `json:"teams"`
	Divisions    []string `json:"divisions"`
	Repositories []string `json:"repositories"`
}

// TimeSeries contains optional time series data.
type TimeSeries struct {
	CommitsPerDay []int     `json:"commits_per_day,omitempty"`
	PRsPerDay     []int     `json:"prs_per_day,omitempty"`
	AvgCycleTime  []float64 `json:"avg_cycle_time,omitempty"`
}
