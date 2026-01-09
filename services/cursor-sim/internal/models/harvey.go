package models

import (
	"fmt"
	"time"
)

// HarveyTask represents the type of AI assistant task
type HarveyTask string

const (
	HarveyTaskAssist   HarveyTask = "Assist"   // General questions
	HarveyTaskDraft    HarveyTask = "Draft"    // Document drafting
	HarveyTaskReview   HarveyTask = "Review"   // Contract review
	HarveyTaskResearch HarveyTask = "Research" // Legal research
)

// HarveySource represents the data source for the task
type HarveySource string

const (
	HarveySourceFiles     HarveySource = "Files"     // Uploaded documents
	HarveySourceWeb       HarveySource = "Web"       // Web search
	HarveySourceKnowledge HarveySource = "Knowledge" // Knowledge base
)

// HarveySentiment represents feedback sentiment
type HarveySentiment string

const (
	HarveySentimentPositive HarveySentiment = "positive"
	HarveySentimentNegative HarveySentiment = "negative"
	HarveySentimentNeutral  HarveySentiment = "neutral"
)

// HarveyUsageEvent represents a single AI assistant usage event.
// Field names match the Harvey API specification exactly.
type HarveyUsageEvent struct {
	EventID           int64           `json:"event_id"`
	MessageID         string          `json:"message_ID"`
	Time              time.Time       `json:"Time"`
	User              string          `json:"User"`
	Task              HarveyTask      `json:"Task"`
	ClientMatter      float64         `json:"Client Matter #"`
	Source            HarveySource    `json:"Source"`
	NumberOfDocuments int             `json:"Number of documents"`
	FeedbackComments  string          `json:"Feedback Comments"`
	FeedbackSentiment HarveySentiment `json:"Feedback Sentiment"`
}

// Validate checks that all required fields are present and valid.
func (e *HarveyUsageEvent) Validate() error {
	if e.EventID == 0 {
		return fmt.Errorf("event_id is required")
	}
	if e.MessageID == "" {
		return fmt.Errorf("message_ID is required")
	}
	if e.Time.IsZero() {
		return fmt.Errorf("Time is required")
	}
	if e.User == "" {
		return fmt.Errorf("User is required")
	}
	if e.Task == "" {
		return fmt.Errorf("Task is required")
	}
	return nil
}
