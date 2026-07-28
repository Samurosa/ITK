package errors

import "errors"

var (
	ErrCompareBaseQuoteAsset = errors.New("base asset cannot be equal to quote asset")

	ErrInvalidMinOrderGreaterMaxOrder = errors.New("minimal order size cannot be greater than max order size")
	ErrInvalidOrderSizePrecision      = errors.New("invalid min order size relative to size precision")
	ErrInvalidMinOrderSize            = errors.New("invalid minimal order size")
	ErrInvalidMaxOrderSize            = errors.New("invalid maximal order size")

	ErrSaveSpot    = errors.New("spot save failed")
	ErrDisableSpot = errors.New("spot disable failed")
	ErrEnableSpot  = errors.New("spot enable failed")
)
