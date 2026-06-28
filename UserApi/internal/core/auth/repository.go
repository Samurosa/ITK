package auth

import (
	"ITK_Code/m/v2/internal/core/user/models"
	"context"
)

type SessionRepository interface {
	Create(ctx context.Context, SessionModel SessionModel) error

	GetByUserIdAndDeviceId(ctx context.Context, userID string, deviceID string) (SessionModel, error)

	Update(ctx context.Context, SessionModel SessionModel) error

	DeleteByUserAndDevice(ctx context.Context, userID, deviceID string) error
}

type TokenManager interface {
	Generate(user models.User) (TokensModel, error)
}
