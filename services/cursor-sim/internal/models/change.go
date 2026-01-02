package models

import "time"

// Change represents an individual AI-generated code change
type Change struct {
	ChangeID      string    `json:"change_id"` // Deterministic ID
	CommitHash    string    `json:"commit_hash"`
	UserID        string    `json:"user_id"`
	Timestamp     time.Time `json:"timestamp"`
	Source        string    `json:"source"` // "TAB" or "COMPOSER"
	Model         string    `json:"model"`  // "claude-3.5-sonnet", etc.
	FilePath      string    `json:"file_path"`
	FileExtension string    `json:"file_extension"` // ".ts", ".go", ".py"
	LinesAdded    int       `json:"lines_added"`
	LinesRemoved  int       `json:"lines_removed"`
	IngestionTime time.Time `json:"ingestion_time"`
}
