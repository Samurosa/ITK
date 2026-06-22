package user

import (
	models2 "ITK_Code/m/v2/internal/core/user/models"
	"context"
	"time"
)

type Service interface {
	Registration(ctx context.Context,
		email string,
		password string,
		name string,
	) (
		id string,
		tokenPairs models2.Tokens,
		createdAt time.Time,
		err error,
	)

	GetUser(ctx context.Context,
		id string,
	) (
		name string,
		email string,
		balances map[string]models2.Balance,
		role models2.Role,
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
		amount models2.Money,
	) (
		success bool,
		balance models2.Balance,
		err error,
	)

	Login(ctx context.Context,
		email string,
		password string,
	) (
		tokensPairs models2.Tokens,
		err error,
	)

	Logout(ctx context.Context,
		refreshToken string,
		deviceID string,
	) (
		success bool,
		loggedOutAt time.Time,
		err error,
	)

	GetBalance(ctx context.Context,
		id string,
	) (
		balances map[string]models2.Balance,
		err error,
	)

	RefreshToken(ctx context.Context,
		refreshToken string,
		deviceID string,
	) (
		tokensPairs models2.Tokens,
		err error,
	)
}
