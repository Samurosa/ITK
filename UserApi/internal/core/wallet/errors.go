package wallet

import "errors"

var (
	ErrBalanceNotFound       = errors.New("balance not found")
	ErrBalanceAllReadyExists = errors.New("balance all ready exists")
	ErrCreateNewBalance      = errors.New("create new balance")
	ErrSaveBalance           = errors.New("save new balance")
)
