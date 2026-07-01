package user

import (
	"context"
	"time"
)

type Service interface {
	GetUser(ctx context.Context,
		id string,
	) (
		user User,
		err error,
	)

	GetUserByEmail(ctx context.Context,
		email string,
	) (
		user User,
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
		err error,
	)
}
