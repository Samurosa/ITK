package user

import (
	"ITK_Code/m/v2/internal/core/user/models"
	"context"
)

type Save interface {
	SaveUser(ctx context.Context,
		user models.User,
	) (
		string,
		error,
	)

	IsExistsUserByEmail(ctx context.Context,
		login string,
	) bool
}

type Provider interface {
	Get(ctx context.Context, uid string) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetBalance(ctx context.Context, uid string, asset string) (models.Balance, error)
	Update(ctx context.Context, user *models.User, update models.UpdateUser) (bool, error)
	Delete(ctx context.Context, uid string) error
	IsAdmin(ctx context.Context, uid string) (bool, error)
}
