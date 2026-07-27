package user

import (
	"ITK_Code/m/v2/internal/core/dto"
	"context"
)

type Save interface {
	SaveUser(ctx context.Context,
		user dto.User,
	) (
		string,
		error,
	)
}

type Provider interface {
	Get(ctx context.Context, uid string) (dto.User, error)

	GetByEmail(ctx context.Context, email string) (dto.User, error)

	Update(ctx context.Context, userID string, update dto.UpdateUser) (bool, error)

	UpdatePassword(ctx context.Context, user dto.User, newPass string) (bool, error)

	Delete(ctx context.Context, uid string) error

	IsAdmin(ctx context.Context, uid string) (bool, error)
}
