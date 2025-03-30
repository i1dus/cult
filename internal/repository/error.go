package repository

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")

	ErrBookingExists    = errors.New("booking already exists")
	ErrBookingNotFound  = errors.New("booking not found")
	ErrInvalidTimeRange = errors.New("invalid time range")
	ErrBookingConflict  = errors.New("booking conflict")
)
