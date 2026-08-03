package errors

import "errors"

var (
	ErrPasswordEmpty            = errors.New("password empty")
	ErrPasswordsMatch           = errors.New("new password matches the old password")
	ErrPasswordWrongUpperSymbol = errors.New("password wrong, upper symbol not found")
	ErrPasswordWrongDigitSymbol = errors.New("password wrong, digit not found")

	ErrPassGenHash     = errors.New("error generating password hash")
	ErrComparePassword = errors.New("password comparison failed")
)
