package corerrors

import "errors"

var (
	ErrRolePermissionDenied = errors.New("role Permission Denied")
	ErrInvalidOrder         = errors.New("order values violate spot constraints")
	ErrOrderNotFound        = errors.New("order not found")
	ErrSpotNotFound         = errors.New("spot not found")
	ErrSpotInactive         = errors.New("spot is not active")
	ErrIdempotencyConflict  = errors.New("idempotency key already used for another order")
	ErrInvalidCursor        = errors.New("invalid order cursor")
)
