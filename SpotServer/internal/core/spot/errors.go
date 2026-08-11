package spot

import "errors"

var (
	ErrSaveSpot    = errors.New("spot save failed")
	ErrGetSpot     = errors.New("spot get failed")
	ErrDisableSpot = errors.New("spot disable failed")
	ErrEnableSpot  = errors.New("spot enable failed")
)
