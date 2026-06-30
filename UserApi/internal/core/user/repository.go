package user

import (
	"context"
)

type Save interface {
	SaveUser(ctx context.Context,
		user User,
	) (
		string,
		error,
	)

	IsExistsUserByEmail(ctx context.Context,
		login string,
	) bool
}

type Provider interface {
	Get(ctx context.Context, uid string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, user *User, update UpdateUser) (bool, error)
	Delete(ctx context.Context, uid string) error
	IsAdmin(ctx context.Context, uid string) (bool, error)
}
