package errors

import (
	"errors"
)

var (
	ErrUserIDEmpty   = errors.New("user id empty")
	ErrEmailEmpty    = errors.New("email empty")
	ErrUsernameEmpty = errors.New("username empty")
	ErrAssetEmpty    = errors.New("asset empty")
	ErrAmountEmpty   = errors.New("amount empty")

	ErrPasswordEmpty            = errors.New("password empty")
	ErrPasswordsMatch           = errors.New("new password matches the old password")
	ErrPasswordWrongUpperSymbol = errors.New("password wrong, upper symbol not found")
	ErrPasswordWrongDigitSymbol = errors.New("password wrong, digit not found")
)
