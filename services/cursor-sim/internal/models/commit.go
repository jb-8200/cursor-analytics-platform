package models

import "time"

// Commit represents a git commit with AI-generated code tracking
type Commit struct {
	Hash              string    `json:"commit_hash"`
	Timestamp         time.Time `json:"timestamp"`
	Message           string    `json:"message"`
	UserID            string    `json:"user_id"`
	UserEmail         string    `json:"user_email"`
	Repository        string    `json:"repository"`
	Branch            string    `json:"branch"`
	TotalLines        int       `json:"total_lines"`
	LinesFromTAB      int       `json:"lines_from_tab"`
	LinesFromComposer int       `json:"lines_from_composer"`
	LinesNonAI        int       `json:"lines_non_ai"`
	IngestionTime     time.Time `json:"ingestion_time"`
}
