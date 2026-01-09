package generator

import (
	"crypto/rand"
	"encoding/hex"
	"math"
	mathrand "math/rand"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// HarveyGenerator generates Harvey AI usage events based on seed configuration.
type HarveyGenerator struct {
	seedData *seed.SeedData
	rng      *mathrand.Rand
}

// HarveyConfig defines configurable parameters for Harvey event generation.
type HarveyConfig struct {
	// BaseEventsPerDay is the average number of events per user per day
	BaseEventsPerDay float64

	// TaskDistribution defines the probability distribution for task types
	TaskDistribution map[models.HarveyTask]float64

	// SentimentRates defines the probability distribution for feedback sentiment
	SentimentRates map[models.HarveySentiment]float64

	// WorkingHours defines the time window when events occur
	WorkingHours struct {
		Start int // Hour of day (0-23)
		End   int // Hour of day (0-23)
	}
}

// NewHarveyGenerator creates a new Harvey generator with a random seed.
func NewHarveyGenerator(seedData *seed.SeedData) *HarveyGenerator {
	return NewHarveyGeneratorWithSeed(seedData, time.Now().UnixNano())
}

// NewHarveyGeneratorWithSeed creates a new Harvey generator with a specific seed for reproducibility.
func NewHarveyGeneratorWithSeed(seedData *seed.SeedData, randSeed int64) *HarveyGenerator {
	return &HarveyGenerator{
		seedData: seedData,
		rng:      mathrand.New(mathrand.NewSource(randSeed)),
	}
}

// DefaultHarveyConfig returns the default configuration for Harvey event generation.
func DefaultHarveyConfig() HarveyConfig {
	config := HarveyConfig{
		BaseEventsPerDay: 5.0,
		TaskDistribution: map[models.HarveyTask]float64{
			models.HarveyTaskAssist:   0.35, // 35% general questions
			models.HarveyTaskDraft:    0.30, // 30% document drafting
			models.HarveyTaskReview:   0.25, // 25% contract review
			models.HarveyTaskResearch: 0.10, // 10% legal research
		},
		SentimentRates: map[models.HarveySentiment]float64{
			models.HarveySentimentPositive: 0.70, // 70% positive
			models.HarveySentimentNeutral:  0.20, // 20% neutral
			models.HarveySentimentNegative: 0.10, // 10% negative
		},
	}
	config.WorkingHours.Start = 8  // 8 AM
	config.WorkingHours.End = 18   // 6 PM

	return config
}

// GenerateEvents generates Harvey usage events for the specified time range.
// Uses Poisson distribution for event counts and uniform distribution for timing within working hours.
func (g *HarveyGenerator) GenerateEvents(from, to time.Time, config HarveyConfig) []models.HarveyUsageEvent {
	var events []models.HarveyUsageEvent

	// If no Harvey config or no developers, return empty
	if g.seedData.ExternalDataSources == nil ||
		g.seedData.ExternalDataSources.Harvey == nil ||
		!g.seedData.ExternalDataSources.Harvey.Enabled ||
		len(g.seedData.Developers) == 0 {
		return events
	}

	harveyConfig := g.seedData.ExternalDataSources.Harvey

	// Calculate total days in range
	duration := to.Sub(from)
	days := int(math.Ceil(duration.Hours() / 24.0))
	if days <= 0 {
		return events
	}

	eventID := int64(1)

	// Generate events for each developer (treating them as Harvey users)
	for _, dev := range g.seedData.Developers {
		// Generate event count for this user over the entire period using Poisson distribution
		expectedEvents := config.BaseEventsPerDay * float64(days)
		eventCount := g.poisson(expectedEvents)

		// Generate events uniformly distributed across the time range
		for i := 0; i < eventCount; i++ {
			// Generate random timestamp within working hours
			eventTime := g.randomWorkingHoursTime(from, to, config.WorkingHours)

			// Select task type based on distribution
			task := g.selectTask(config.TaskDistribution)

			// Select source
			source := g.selectSource()

			// Select sentiment
			sentiment := g.selectSentiment(config.SentimentRates)

			// Generate number of documents (0-5, with bias toward 1-2)
			numDocs := g.rng.Intn(6)

			// Generate client matter number from practice areas
			clientMatter := g.generateClientMatter(harveyConfig.PracticeAreas)

			// Generate feedback comments (optional, ~30% have feedback)
			feedbackComments := ""
			if g.rng.Float64() < 0.30 {
				feedbackComments = g.generateFeedbackComment(sentiment)
			}

			event := models.HarveyUsageEvent{
				EventID:           eventID,
				MessageID:         g.generateMessageID(),
				Time:              eventTime,
				User:              dev.Email,
				Task:              task,
				ClientMatter:      clientMatter,
				Source:            source,
				NumberOfDocuments: numDocs,
				FeedbackComments:  feedbackComments,
				FeedbackSentiment: sentiment,
			}

			events = append(events, event)
			eventID++
		}
	}

	// Sort events by time
	return g.sortEventsByTime(events)
}

// poisson generates a random number from a Poisson distribution with the given lambda.
func (g *HarveyGenerator) poisson(lambda float64) int {
	if lambda <= 0 {
		return 0
	}

	// For large lambda, use normal approximation
	if lambda > 30 {
		return int(math.Max(0, g.rng.NormFloat64()*math.Sqrt(lambda)+lambda))
	}

	// Knuth's algorithm for small lambda
	L := math.Exp(-lambda)
	k := 0
	p := 1.0

	for p > L {
		k++
		p *= g.rng.Float64()
	}

	return k - 1
}

// randomWorkingHoursTime generates a random timestamp within the time range during working hours.
func (g *HarveyGenerator) randomWorkingHoursTime(from, to time.Time, workingHours struct{ Start, End int }) time.Time {
	// Calculate total working hours in the range
	duration := to.Sub(from)
	totalDays := duration.Hours() / 24.0

	// Pick a random day
	dayOffset := g.rng.Float64() * totalDays
	baseTime := from.Add(time.Duration(dayOffset*24) * time.Hour)

	// Set to a random hour within working hours
	workingHourSpan := workingHours.End - workingHours.Start
	randomHourOffset := g.rng.Intn(workingHourSpan)
	targetHour := workingHours.Start + randomHourOffset

	// Random minute and second
	randomMinute := g.rng.Intn(60)
	randomSecond := g.rng.Intn(60)

	// Construct the timestamp
	result := time.Date(
		baseTime.Year(),
		baseTime.Month(),
		baseTime.Day(),
		targetHour,
		randomMinute,
		randomSecond,
		0,
		time.UTC,
	)

	// Ensure it's within bounds
	if result.Before(from) {
		result = from
	}
	if result.After(to) {
		result = to
	}

	return result
}

// selectTask selects a task type based on the configured distribution.
func (g *HarveyGenerator) selectTask(distribution map[models.HarveyTask]float64) models.HarveyTask {
	r := g.rng.Float64()
	cumulative := 0.0

	tasks := []models.HarveyTask{
		models.HarveyTaskAssist,
		models.HarveyTaskDraft,
		models.HarveyTaskReview,
		models.HarveyTaskResearch,
	}

	for _, task := range tasks {
		cumulative += distribution[task]
		if r < cumulative {
			return task
		}
	}

	// Fallback
	return models.HarveyTaskAssist
}

// selectSource selects a random source for the event.
func (g *HarveyGenerator) selectSource() models.HarveySource {
	// Weighted: 50% Files, 30% Knowledge, 20% Web
	r := g.rng.Float64()
	if r < 0.50 {
		return models.HarveySourceFiles
	} else if r < 0.80 {
		return models.HarveySourceKnowledge
	}
	return models.HarveySourceWeb
}

// selectSentiment selects a sentiment based on the configured distribution.
func (g *HarveyGenerator) selectSentiment(distribution map[models.HarveySentiment]float64) models.HarveySentiment {
	r := g.rng.Float64()
	cumulative := 0.0

	sentiments := []models.HarveySentiment{
		models.HarveySentimentPositive,
		models.HarveySentimentNeutral,
		models.HarveySentimentNegative,
	}

	for _, sentiment := range sentiments {
		cumulative += distribution[sentiment]
		if r < cumulative {
			return sentiment
		}
	}

	// Fallback
	return models.HarveySentimentNeutral
}

// generateClientMatter generates a client matter number based on practice areas.
func (g *HarveyGenerator) generateClientMatter(practiceAreas []string) float64 {
	// Generate a matter number like 2024.001, 2025.042, etc.
	year := 2024 + g.rng.Intn(2) // 2024 or 2025
	matterNum := g.rng.Intn(100) + 1
	return float64(year) + float64(matterNum)/1000.0
}

// generateFeedbackComment generates a realistic feedback comment based on sentiment.
func (g *HarveyGenerator) generateFeedbackComment(sentiment models.HarveySentiment) string {
	positive := []string{
		"Very helpful response",
		"Exactly what I needed",
		"Saved me hours of research",
		"Great analysis",
		"Comprehensive and accurate",
	}

	neutral := []string{
		"Adequate response",
		"Needs more detail",
		"Acceptable",
		"Could be more specific",
	}

	negative := []string{
		"Not relevant to my question",
		"Inaccurate information",
		"Too generic",
		"Missing key details",
		"Did not answer my question",
	}

	var pool []string
	switch sentiment {
	case models.HarveySentimentPositive:
		pool = positive
	case models.HarveySentimentNeutral:
		pool = neutral
	case models.HarveySentimentNegative:
		pool = negative
	default:
		pool = neutral
	}

	return pool[g.rng.Intn(len(pool))]
}

// generateMessageID generates a unique message ID (UUID-like).
func (g *HarveyGenerator) generateMessageID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// sortEventsByTime sorts events by timestamp.
func (g *HarveyGenerator) sortEventsByTime(events []models.HarveyUsageEvent) []models.HarveyUsageEvent {
	// Simple bubble sort (fine for reasonable event counts)
	n := len(events)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if events[j].Time.After(events[j+1].Time) {
				events[j], events[j+1] = events[j+1], events[j]
			}
		}
	}
	return events
}
