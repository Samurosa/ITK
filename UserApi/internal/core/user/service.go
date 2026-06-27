package user

import (
	"ITK_Code/m/v2/internal/core/user/models"
	"context"
	"time"
)

type Service interface {
	GetUser(ctx context.Context,
		id string,
	) (
		name string,
		email string,
		balances map[string]models.Balance,
		role models.Role,
		createdAt time.Time,
		updatedAt time.Time,
		err error,
	)

	UpdateUserInfo(ctx context.Context,
		id string,
		name string,
		email string,
	) (
		updated bool,
		updatedUserInfoAt time.Time,
		err error,
	)

	DeleteUser(ctx context.Context,
		id string,
	) (
		success bool,
		deletedUserAt time.Time,
		err error,
	)

	ChangePassword(ctx context.Context,
		id string,
		oldPassword string,
		newPassword string,
	) (
		success bool,
		userPasswordChangedAt time.Time,
	)

	Deposit(ctx context.Context,
		id string,
		asset string,
		amount models.Money,
	) (
		success bool,
		balance models.Balance,
		err error,
	)

	GetBalance(ctx context.Context,
		id string,
	) (
		balances map[string]models.Balance,
		err error,
	)
}
