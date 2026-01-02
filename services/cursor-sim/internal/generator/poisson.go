package generator

import (
	"math"
	"math/rand"
	"time"
)

// PoissonTimer generates inter-event times following a Poisson process
// Events in a Poisson process have exponentially distributed inter-arrival times
type PoissonTimer struct {
	lambda float64    // Rate parameter (events per hour)
	rng    *rand.Rand // Random number generator
}

// NewPoissonTimer creates a new Poisson timer with the given rate (events/hour) and seed
func NewPoissonTimer(lambda float64, seed int64) *PoissonTimer {
	return &PoissonTimer{
		lambda: lambda,
		rng:    rand.New(rand.NewSource(seed)),
	}
}

// NextInterval returns the next inter-event time duration
// For a Poisson process with rate λ events/hour, inter-event times follow
// an exponential distribution with mean 1/λ hours
func (p *PoissonTimer) NextInterval() time.Duration {
	// Generate exponentially distributed random variable
	// Using inverse transform: -ln(U) / λ where U ~ Uniform(0,1)
	u := p.rng.Float64()

	// Avoid log(0) by ensuring u > 0
	for u == 0 {
		u = p.rng.Float64()
	}

	// Calculate interval in hours
	intervalHours := -math.Log(u) / p.lambda

	// Convert to nanoseconds and return as duration
	intervalNs := intervalHours * 3600.0 * 1e9
	return time.Duration(int64(intervalNs))
}

// VelocityToLambda converts a velocity string to events/hour rate
func VelocityToLambda(velocity string) float64 {
	switch velocity {
	case "low":
		return 5.0 // 5 events per hour
	case "medium":
		return 25.0 // 25 events per hour
	case "high":
		return 50.0 // 50 events per hour
	default:
		return 25.0 // Default to medium
	}
}

// ApplyVolatility adjusts the base lambda with per-developer variance
// Returns a new lambda value within [baseLambda*(1-volatility), baseLambda*(1+volatility)]
func ApplyVolatility(baseLambda, volatility float64, seed int64) float64 {
	rng := rand.New(rand.NewSource(seed))

	// Generate random factor in range [-volatility, +volatility]
	factor := (rng.Float64()*2.0 - 1.0) * volatility

	// Apply factor to base lambda
	adjustedLambda := baseLambda * (1.0 + factor)

	// Ensure it stays positive
	if adjustedLambda < 0.1 {
		adjustedLambda = 0.1
	}

	return adjustedLambda
}
