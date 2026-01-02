package generator

// VelocityConfig controls the rate of event generation.
type VelocityConfig struct {
	multiplier float64
	baseRate   float64
}

// NewVelocityConfig creates a velocity configuration for the specified level.
// Valid levels: "low", "medium", "high". Defaults to "medium" for invalid values.
func NewVelocityConfig(velocity string) *VelocityConfig {
	cfg := &VelocityConfig{}

	switch velocity {
	case "low":
		cfg.multiplier = 0.5
		cfg.baseRate = 80.0 // events per day
	case "high":
		cfg.multiplier = 2.0
		cfg.baseRate = 200.0
	case "medium":
		fallthrough
	default:
		cfg.multiplier = 1.0
		cfg.baseRate = 120.0
	}

	return cfg
}

// CommitsPerDay calculates the expected number of commits per day
// based on the developer's PRs per week and velocity multiplier.
func (v *VelocityConfig) CommitsPerDay(prsPerWeek float64) float64 {
	// Assume ~3 commits per PR on average
	commitsPerWeek := prsPerWeek * 3.0
	commitsPerDay := commitsPerWeek / 7.0
	return commitsPerDay * v.multiplier
}

// EventsPerDay returns the base event rate (tab completions, etc.) per day.
func (v *VelocityConfig) EventsPerDay() float64 {
	return v.baseRate
}
