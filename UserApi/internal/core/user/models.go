package user

import "time"

type Role string

const (
	UnspecifiedRole Role = "ROLE_UNSPECIFIED"

	UserRole Role = "ROLE_USER"

	GuestRole Role = "ROLE_GUEST"

	PremiumRole Role = "ROLE_PREMIUM"

	AdminRole Role = "ROLE_ADMIN"
)

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"password_hash"`
	Role         Role      `json:"role"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
}

type UpdateUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
}
