package models

import "time"

type Role string

const (
	UnspecifiedRole Role = "ROLE_UNSPECIFIED"

	UserRole Role = "ROLE_USER"

	AdminRole Role = "ROLE_ADMIN"
)

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
type UpdateUser struct {
	Name     string
	Login    string
	Password []byte
}
