package user

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUpdateUser   = errors.New("error updating user")
	ErrEmailIsExist = errors.New("email is exist")
)
