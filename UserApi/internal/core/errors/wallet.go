package errors

import "errors"

var (
	ErrBalanceNotFound  = errors.New("balance not found")
	ErrCreateNewBalance = errors.New("create new balance")
	ErrSaveBalance      = errors.New("save new balance")
)
