package user

import "errors"

var (
	ErrUserExists      = errors.New("user exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrComparePassword = errors.New("password comparison failed")
	ErrUpdateUser      = errors.New("error updating user")
)
