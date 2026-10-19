package errorcore

import "errors"

var (
	ErrOrderListenerClosed = errors.New("order listener closed unexpectedly")
)
