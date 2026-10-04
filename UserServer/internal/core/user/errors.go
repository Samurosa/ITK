package user

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailIsExist = errors.New("email is exist")
)
