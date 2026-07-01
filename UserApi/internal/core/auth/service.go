package auth

import (
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
		createdAt time.Time,
		err error,
	)

	Login(ctx context.Context,
		email string,
		password string,
		deviceId string,
	) (
		tokensPairs TokensModel,
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

	RefreshToken(ctx context.Context,
		refreshToken string,
		deviceID string,
	) (
		tokensPairs TokensModel,
		err error,
	)
}
