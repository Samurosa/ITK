package auth

import (
	"ITK_Code/m/v2/internal/core/dto"
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
		tokensPairs dto.TokensModel,
		err error,
	)

	Logout(ctx context.Context,
		jti string,
		refreshToken string,
	) (
		err error,
	)

	LogoutAllDevices(ctx context.Context,
		jti string,
		refreshToken string,
	) (
		err error,
	)

	RefreshToken(ctx context.Context,
		refreshToken string,
	) (
		tokensPairs dto.TokensModel,
		err error,
	)
}
