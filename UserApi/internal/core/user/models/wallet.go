package models

type Balance struct {
	Asset string

	Available string

	Locked string
}

type Money struct {
	Currency string
	Units    int64
	Nanos    int32
}
