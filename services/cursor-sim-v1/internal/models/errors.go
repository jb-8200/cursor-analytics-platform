package models

import "errors"

var (
	// ErrInvalidDeveloperID indicates an invalid developer ID
	ErrInvalidDeveloperID = errors.New("invalid developer ID")

	// ErrInvalidEmail indicates an invalid email address
	ErrInvalidEmail = errors.New("invalid email address")

	// ErrInvalidName indicates an invalid name
	ErrInvalidName = errors.New("invalid name")

	// ErrNotFound indicates a requested item was not found
	ErrNotFound = errors.New("not found")
)
