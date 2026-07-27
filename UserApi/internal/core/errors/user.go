package errors

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrComparePassword = errors.New("password comparison failed")
	ErrUpdateUser      = errors.New("error updating user")
	ErrPassGenHash     = errors.New("error generating password hash")
	ErrEmailIsExist    = errors.New("email is exist")
)
