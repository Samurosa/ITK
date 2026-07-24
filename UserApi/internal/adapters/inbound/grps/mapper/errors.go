package mapper

import (
	"ITK_Code/m/v2/internal/adapters/outbound/repository/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/repository/redis"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"ITK_Code/m/v2/internal/core/wallet"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPC(err error) error {
	switch {
	case errors.Is(err, wallet.ErrBalanceNotFound):
		return status.Error(codes.NotFound, "balance not found")
	case errors.Is(err, wallet.ErrCreateNewBalance):
		return status.Error(codes.Internal, "failed to create balance")
	case errors.Is(err, wallet.ErrSaveBalance):
		return status.Error(codes.Internal, "failed to save balance")

	case errors.Is(err, user.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, user.ErrComparePassword):
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case errors.Is(err, user.ErrUpdateUser):
		return status.Error(codes.Internal, "failed to update user")
	case errors.Is(err, user.ErrEmailIsExist):
		return status.Error(codes.AlreadyExists, "email is exist")

	case errors.Is(err, auth.ErrRefreshExpired):
		return status.Error(codes.Unauthenticated, "refresh token expired")
	case errors.Is(err, auth.ErrGenerateToken):
		return status.Error(codes.Internal, "failed to generate token")
	case errors.Is(err, auth.ErrInvalidToken):
		return status.Error(codes.InvalidArgument, "invalid token")
	case errors.Is(err, auth.ErrSessionExpired):
		return status.Error(codes.Unauthenticated, "session expired")
	case errors.Is(err, auth.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, "session not found")
	case errors.Is(err, auth.ErrInvalidContext):
		return status.Error(codes.Internal, "invalid context")
	case errors.Is(err, auth.ErrInvalidLoginCredentials):
		return status.Error(codes.InvalidArgument, "invalid login")
	case errors.Is(err, auth.Unauthorized):
		return status.Error(codes.Unauthenticated, "incorrect login or password")
	case errors.Is(err, auth.ErrNoAccess):
		return status.Error(codes.PermissionDenied, "no access")
	case errors.Is(err, auth.ErrTooManyRequests):
		return status.Error(codes.Aborted, "too many requests")
	case errors.Is(err, postgres.ErrPingDB):
		return status.Error(codes.Internal, "failed connect to database")
	case errors.Is(err, redis.ErrPingToRedis):
		return status.Error(codes.Internal, "failed connect to redis")

	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
