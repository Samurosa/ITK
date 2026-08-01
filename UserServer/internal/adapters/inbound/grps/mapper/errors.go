package mapper

import (
	"ITK_Code/m/v2/internal/adapters/outbound/repository/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/repository/redis"
	coreErrors "ITK_Code/m/v2/internal/core/errors"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPC(err error) error {
	switch {
	case errors.Is(err, coreErrors.ErrUserIDEmpty):
		return status.Error(codes.InvalidArgument, "the user ID in the request is empty")
	case errors.Is(err, coreErrors.ErrEmailEmpty):
		return status.Error(codes.InvalidArgument, "the email in the request is empty")
	case errors.Is(err, coreErrors.ErrUsernameEmpty):
		return status.Error(codes.InvalidArgument, "the username in the request is empty")
	case errors.Is(err, coreErrors.ErrAssetEmpty):
		return status.Error(codes.InvalidArgument, "the asset in the request is empty")
	case errors.Is(err, coreErrors.ErrAmountEmpty):
		return status.Error(codes.InvalidArgument, "the amount in the request is empty")
	case errors.Is(err, coreErrors.ErrPasswordEmpty):
		return status.Error(codes.InvalidArgument, "the password in the request is empty")
	case errors.Is(err, coreErrors.ErrPasswordsMatch):
		return status.Error(codes.InvalidArgument, "the new password matches the old password")
	case errors.Is(err, coreErrors.ErrPasswordWrongUpperSymbol):
		return status.Error(codes.InvalidArgument, "the new password wrong, upper symbol not found")
	case errors.Is(err, coreErrors.ErrPasswordWrongDigitSymbol):
		return status.Error(codes.InvalidArgument, "the new password wrong, digit symbol not found")

	case errors.Is(err, coreErrors.ErrBalanceNotFound):
		return status.Error(codes.NotFound, "balance not found")
	case errors.Is(err, coreErrors.ErrCreateNewBalance):
		return status.Error(codes.Internal, "failed to create balance")
	case errors.Is(err, coreErrors.ErrSaveBalance):
		return status.Error(codes.Internal, "failed to save balance")

	case errors.Is(err, coreErrors.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, coreErrors.ErrComparePassword):
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case errors.Is(err, coreErrors.ErrUpdateUser):
		return status.Error(codes.Internal, "failed to update user")
	case errors.Is(err, coreErrors.ErrEmailIsExist):
		return status.Error(codes.AlreadyExists, "email is exist")

	case errors.Is(err, coreErrors.ErrRefreshExpired):
		return status.Error(codes.Unauthenticated, "refresh token expired")
	case errors.Is(err, coreErrors.ErrGenerateToken):
		return status.Error(codes.Internal, "failed to generate token")
	case errors.Is(err, coreErrors.ErrInvalidToken):
		return status.Error(codes.InvalidArgument, "invalid token")
	case errors.Is(err, coreErrors.ErrSessionExpired):
		return status.Error(codes.Unauthenticated, "session expired")
	case errors.Is(err, coreErrors.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, "session not found")
	case errors.Is(err, coreErrors.ErrInvalidContext):
		return status.Error(codes.Internal, "invalid context")
	case errors.Is(err, coreErrors.ErrInvalidLoginCredentials):
		return status.Error(codes.InvalidArgument, "invalid login")
	case errors.Is(err, coreErrors.Unauthorized):
		return status.Error(codes.Unauthenticated, "incorrect login or password")
	case errors.Is(err, coreErrors.ErrNoAccess):
		return status.Error(codes.PermissionDenied, "no access")
	case errors.Is(err, coreErrors.ErrTooManyRequests):
		return status.Error(codes.Aborted, "too many requests")
	case errors.Is(err, postgres.ErrPingDB):
		return status.Error(codes.Internal, "failed connect to database")
	case errors.Is(err, redis.ErrPingToRedis):
		return status.Error(codes.Internal, "failed connect to redis")

	default:
		return status.Error(codes.Internal, "internal UserServer error")
	}
}
