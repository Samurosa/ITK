package ports

import (
	"ITK_Code/m/v2/internal/core/user"
	"context"
)

type PasswordHasher interface {
	GeneratePasswordHash(password string) ([]byte, error)
	VerifyPasswordHash(password string, hash []byte) error
}

type UserRepository interface {
	SaveUser(context.Context, user.User) (string, error)
	Get(context.Context, string) (user.User, error)
	GetByEmail(context.Context, string) (user.User, error)
	Update(context.Context, string, user.UpdateUser) error
	UpdatePassword(context.Context, user.User, string) error
	Delete(context.Context, string) error
	IsAdmin(context.Context, string) (bool, error)
}
