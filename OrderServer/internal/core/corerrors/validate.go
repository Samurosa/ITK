package corerrors

import "errors"

var (
	ErrRolePermissionDenied = errors.New("role Permission Denied")
	ErrInvalidOrder         = errors.New("order values violate spot constraints")
)
