package errorcore

import "errors"

var (
	ErrPasswordsMatch           = errors.New("new password matches the old password")
	ErrPasswordWrongUpperSymbol = errors.New("password wrong, upper symbol not found")
	ErrPasswordWrongLowerSymbol = errors.New("password wrong, lower symbol not found")
	ErrPasswordWrongDigitSymbol = errors.New("password wrong, digit not found")

	ErrPassGenHash    = errors.New("error generating password hash")
	ErrPassVerifyHash = errors.New("error verifying password hash")
)
