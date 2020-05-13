package data

import "errors"

var (
	// ErrNotFound is used when a specific resource is not found.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")
)
