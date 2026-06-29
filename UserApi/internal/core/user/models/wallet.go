package models

type Balance struct {
	Asset string

	Available Money

	Locked Money
}

type Money struct {
	Currency string
	Units    int64
	Nanos    int32
}
