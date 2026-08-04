package errors

import "errors"

var (
	ErrUserIDEmpty      = errors.New("user id empty")
	ErrEmailEmpty       = errors.New("email empty")
	ErrUsernameEmpty    = errors.New("username empty")
	ErrAssetEmpty       = errors.New("asset empty")
	ErrInvalidAsset     = errors.New("invalid asset")
	ErrAmountEmpty      = errors.New("amount empty")
	ErrAmountIsZero     = errors.New("amount is zero value")
	ErrAmountIsNegative = errors.New("amount is negative")
)
