package user

import (
	"context"
)

type Service interface {
	GetUser(ctx context.Context, id string) (user User, err error)

	GetUserByEmail(ctx context.Context, email string) (user User, err error)

	UpdateUserInfo(ctx context.Context, id string, name string, email string) (err error)

	DeleteUser(ctx context.Context, id string) (err error)

	ChangePassword(ctx context.Context, id string, oldPassword string, newPassword string) (err error)
}
