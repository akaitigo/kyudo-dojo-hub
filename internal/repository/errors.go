package repository

import "errors"

var (
	// ErrNotFound indicates the requested resource was not found.
	ErrNotFound = errors.New("not found")

	// ErrReservationConflict indicates a scheduling conflict for a reservation.
	ErrReservationConflict = errors.New("reservation conflict: same lane and overlapping time")
)
