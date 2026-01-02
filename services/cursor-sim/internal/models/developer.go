package models

import (
	"time"

	"github.com/google/uuid"
)

// Developer represents a simulated developer
type Developer struct {
	ID             string    `json:"id"`              // UUID: "user_abc123"
	Email          string    `json:"email"`           // "jane.smith@company.com"
	Name           string    `json:"name"`            // "Jane Smith"
	Region         string    `json:"region"`          // "US", "EU", "APAC"
	Division       string    `json:"division"`        // "AGS", "AT", "ST"
	Group          string    `json:"group"`           // "TMOBILE", "ATANT"
	Team           string    `json:"team"`            // "dev", "support"
	Seniority      string    `json:"seniority"`       // "junior", "mid", "senior"
	ClientVersion  string    `json:"client_version"`  // "0.43.6"
	AcceptanceRate float64   `json:"acceptance_rate"` // 0.0-1.0
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	LastActiveAt   time.Time `json:"last_active_at"`
}

// NewDeveloper creates a new Developer with defaults
func NewDeveloper(name, email string) *Developer {
	now := time.Now().UTC()
	return &Developer{
		ID:             "user_" + uuid.New().String()[:8],
		Email:          email,
		Name:           name,
		Region:         "US",
		Division:       "AGS",
		Group:          "TMOBILE",
		Team:           "dev",
		Seniority:      "mid",
		ClientVersion:  "0.43.6",
		AcceptanceRate: 0.75,
		IsActive:       true,
		CreatedAt:      now,
		LastActiveAt:   now,
	}
}

// Validate checks if the developer data is valid
func (d *Developer) Validate() error {
	if d.ID == "" {
		return ErrInvalidDeveloperID
	}
	if d.Email == "" {
		return ErrInvalidEmail
	}
	if d.Name == "" {
		return ErrInvalidName
	}
	return nil
}
