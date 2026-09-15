package coreErrors

import "errors"

var (
	ErrFailedGetUserID = errors.New("failed get user id")
	ErrInvalidContext  = errors.New("invalid context")
)
