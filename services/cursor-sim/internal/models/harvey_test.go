package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHarveyUsageEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   HarveyUsageEvent
		wantErr bool
	}{
		{
			name: "valid event",
			event: HarveyUsageEvent{
				EventID:           12345,
				MessageID:         "uuid-abc",
				Time:              time.Now(),
				User:              "user@firm.com",
				Task:              HarveyTaskAssist,
				Source:            HarveySourceFiles,
				FeedbackSentiment: HarveySentimentPositive,
			},
			wantErr: false,
		},
		{
			name: "missing event_id",
			event: HarveyUsageEvent{
				MessageID:         "uuid-abc",
				Time:              time.Now(),
				User:              "user@firm.com",
				Task:              HarveyTaskAssist,
				Source:            HarveySourceFiles,
				FeedbackSentiment: HarveySentimentPositive,
			},
			wantErr: true,
		},
		{
			name: "missing user",
			event: HarveyUsageEvent{
				EventID:           12345,
				MessageID:         "uuid-abc",
				Time:              time.Now(),
				Task:              HarveyTaskAssist,
				Source:            HarveySourceFiles,
				FeedbackSentiment: HarveySentimentPositive,
			},
			wantErr: true,
		},
		{
			name: "missing task",
			event: HarveyUsageEvent{
				EventID:           12345,
				MessageID:         "uuid-abc",
				Time:              time.Now(),
				User:              "user@firm.com",
				Source:            HarveySourceFiles,
				FeedbackSentiment: HarveySentimentPositive,
			},
			wantErr: true,
		},
		{
			name: "missing message_id",
			event: HarveyUsageEvent{
				EventID:           12345,
				Time:              time.Now(),
				User:              "user@firm.com",
				Task:              HarveyTaskAssist,
				Source:            HarveySourceFiles,
				FeedbackSentiment: HarveySentimentPositive,
			},
			wantErr: true,
		},
		{
			name: "zero time",
			event: HarveyUsageEvent{
				EventID:           12345,
				MessageID:         "uuid-abc",
				Time:              time.Time{},
				User:              "user@firm.com",
				Task:              HarveyTaskAssist,
				Source:            HarveySourceFiles,
				FeedbackSentiment: HarveySentimentPositive,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHarveyTask_Constants(t *testing.T) {
	assert.Equal(t, HarveyTask("Assist"), HarveyTaskAssist)
	assert.Equal(t, HarveyTask("Draft"), HarveyTaskDraft)
	assert.Equal(t, HarveyTask("Review"), HarveyTaskReview)
	assert.Equal(t, HarveyTask("Research"), HarveyTaskResearch)
}

func TestHarveySource_Constants(t *testing.T) {
	assert.Equal(t, HarveySource("Files"), HarveySourceFiles)
	assert.Equal(t, HarveySource("Web"), HarveySourceWeb)
	assert.Equal(t, HarveySource("Knowledge"), HarveySourceKnowledge)
}

func TestHarveySentiment_Constants(t *testing.T) {
	assert.Equal(t, HarveySentiment("positive"), HarveySentimentPositive)
	assert.Equal(t, HarveySentiment("negative"), HarveySentimentNegative)
	assert.Equal(t, HarveySentiment("neutral"), HarveySentimentNeutral)
}
