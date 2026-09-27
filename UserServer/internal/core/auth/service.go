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
		id string,
		deviceID string,
	) (
		string,
		time.Time,
		error,
	)

	Login(ctx context.Context,
		email string,
		password string,
		ip string,
		deviceID string,
	) (
		dto.TokensModel,
		error,
	)

	Logout(ctx context.Context,
		jti string,
		refreshToken string,
	) (
		err error,
	)

	LogoutAllDevices(ctx context.Context,
		jti string,
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
