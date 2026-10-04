package coreErrors

import "errors"

var (
	ErrSpotNotFound      = errors.New("spot not found")
	ErrSpotAlreadyExists = errors.New("spot already exists")
)
