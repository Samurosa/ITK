package user

import (
	"ITK_Code/m/v2/internal/core/user/models"
	"context"
	"time"
)

type Save interface {
	SaveUser(ctx context.Context,
		login string,
		passwordHash []byte,
		name string,
		balances map[string]*models.Balance,
		role models.Role,
		createTime time.Time,
		updateTime time.Time,
	) (
		string,
		error,
	)
	IsExistsUserByLogin(ctx context.Context,
		login string,
	) bool
}

type Provider interface {
	GetUser(ctx context.Context, uid string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetBalanceUser(ctx context.Context, uid string, asset string) (*models.Balance, error)
	UpdateUser(ctx context.Context, user *models.User, update models.UpdateUser) (bool, error)
	DeleteUser(ctx context.Context, uid string) error
	IsAdmin(ctx context.Context, uid string) (bool, error)
}

type AppProvider interface {
	GetSecret() (models.App, error)
}
