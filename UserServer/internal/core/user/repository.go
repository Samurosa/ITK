package user

import (
	"context"
)

type Repository interface {
	SaveUser(ctx context.Context, user User) (string, error)

	Get(ctx context.Context, uid string) (User, error)

	GetByEmail(ctx context.Context, email string) (User, error)

	Update(ctx context.Context, userID string, update UpdateUser) error

	UpdatePassword(ctx context.Context, user User, newPass string) error

	Delete(ctx context.Context, uid string) error

	IsAdmin(ctx context.Context, uid string) (bool, error)
}
