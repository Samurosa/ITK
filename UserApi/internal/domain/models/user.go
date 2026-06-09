package models

import "time"

type Role string

const (
	UnspecifiedRole = "ROLE_UNSPECIFIED"

	UserRole Role = "ROLE_USER"

	AdminRole Role = "ROLE_ADMIN"
)

type Balance struct {
	Asset string

	Available string

	Locked string
}

type User struct {
	ID           string
	Name         string
	Login        string
	PasswordHash []byte
	Balances     map[string]*Balance
	Role         Role
	CreateTime   time.Time
	UpdateTime   time.Time
}
