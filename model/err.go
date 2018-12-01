package model

import (
	"errors"
)

var (
	// ErrNoResult - no data found
	ErrNoResult = errors.New("no data found")
	// ErrNoDataToDelete - no data found to delete"
	ErrNoDataToDelete = errors.New("no data found to delete")
	// ErrUnavailable - database is unavailable
	ErrUnavailable = errors.New("database is unavailable")
)
