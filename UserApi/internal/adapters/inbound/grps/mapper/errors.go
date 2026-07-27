package mapper

import (
	"ITK_Code/m/v2/internal/adapters/outbound/repository/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/repository/redis"
	errors2 "ITK_Code/m/v2/internal/core/errors"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPC(err error) error {
	switch {
	case errors.Is(err, errors2.ErrBalanceNotFound):
		return status.Error(codes.NotFound, "balance not found")
	case errors.Is(err, errors2.ErrCreateNewBalance):
		return status.Error(codes.Internal, "failed to create balance")
	case errors.Is(err, errors2.ErrSaveBalance):
		return status.Error(codes.Internal, "failed to save balance")

	case errors.Is(err, errors2.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, errors2.ErrComparePassword):
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case errors.Is(err, errors2.ErrUpdateUser):
		return status.Error(codes.Internal, "failed to update user")
	case errors.Is(err, errors2.ErrEmailIsExist):
		return status.Error(codes.AlreadyExists, "email is exist")

	case errors.Is(err, errors2.ErrRefreshExpired):
		return status.Error(codes.Unauthenticated, "refresh token expired")
	case errors.Is(err, errors2.ErrGenerateToken):
		return status.Error(codes.Internal, "failed to generate token")
	case errors.Is(err, errors2.ErrInvalidToken):
		return status.Error(codes.InvalidArgument, "invalid token")
	case errors.Is(err, errors2.ErrSessionExpired):
		return status.Error(codes.Unauthenticated, "session expired")
	case errors.Is(err, errors2.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, "session not found")
	case errors.Is(err, errors2.ErrInvalidContext):
		return status.Error(codes.Internal, "invalid context")
	case errors.Is(err, errors2.ErrInvalidLoginCredentials):
		return status.Error(codes.InvalidArgument, "invalid login")
	case errors.Is(err, errors2.Unauthorized):
		return status.Error(codes.Unauthenticated, "incorrect login or password")
	case errors.Is(err, errors2.ErrNoAccess):
		return status.Error(codes.PermissionDenied, "no access")
	case errors.Is(err, errors2.ErrTooManyRequests):
		return status.Error(codes.Aborted, "too many requests")
	case errors.Is(err, postgres.ErrPingDB):
		return status.Error(codes.Internal, "failed connect to database")
	case errors.Is(err, redis.ErrPingToRedis):
		return status.Error(codes.Internal, "failed connect to redis")

	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
