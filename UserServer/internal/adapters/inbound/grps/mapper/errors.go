package mapper

import (
	"ITK_Code/m/v2/internal/adapters/outbound/repository/postgres"
	"ITK_Code/m/v2/internal/adapters/outbound/repository/redis"
	"ITK_Code/m/v2/internal/core/auth"
	coreErrors "ITK_Code/m/v2/internal/core/errors"
	"ITK_Code/m/v2/internal/core/user"
	"ITK_Code/m/v2/internal/core/wallet"
	"context"
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
	case errors.Is(err, coreErrors.ErrInvalidAsset):
		return status.Error(codes.InvalidArgument, "the asset in the request is invalid")
	case errors.Is(err, coreErrors.ErrAmountEmpty):
		return status.Error(codes.InvalidArgument, "the amount in the request is empty")
	case errors.Is(err, coreErrors.ErrAmountIsZero):
		return status.Error(codes.InvalidArgument, "the amount in the request is zero value")
	case errors.Is(err, coreErrors.ErrAmountIsNegative):
		return status.Error(codes.InvalidArgument, "the amount in the request is negative value")
	case errors.Is(err, coreErrors.ErrPasswordEmpty):
		return status.Error(codes.InvalidArgument, "the password in the request is empty")
	case errors.Is(err, coreErrors.ErrPasswordsMatch):
		return status.Error(codes.InvalidArgument, "the new password matches the old password")
	case errors.Is(err, coreErrors.ErrPasswordWrongUpperSymbol):
		return status.Error(codes.InvalidArgument, "the new password wrong, upper symbol not found")
	case errors.Is(err, coreErrors.ErrPasswordWrongDigitSymbol):
		return status.Error(codes.InvalidArgument, "the new password wrong, digit symbol not found")

	case errors.Is(err, wallet.ErrBalanceNotFound):
		return status.Error(codes.NotFound, "balance not found")
	case errors.Is(err, wallet.ErrCreateNewBalance):
		return status.Error(codes.Internal, "failed to create balance")
	case errors.Is(err, wallet.ErrSaveBalance):
		return status.Error(codes.Internal, "failed to save balance")

	case errors.Is(err, user.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, coreErrors.ErrComparePassword):
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case errors.Is(err, user.ErrUpdateUser):
		return status.Error(codes.Internal, "failed to update user")
	case errors.Is(err, user.ErrEmailIsExist):
		return status.Error(codes.AlreadyExists, "email is exist")

	case errors.Is(err, coreErrors.ErrRefreshExpired):
		return status.Error(codes.Unauthenticated, "refresh token expired")
	case errors.Is(err, coreErrors.ErrGenerateToken):
		return status.Error(codes.Internal, "failed to generate token")
	case errors.Is(err, coreErrors.ErrInvalidToken):
		return status.Error(codes.InvalidArgument, "invalid token")
	case errors.Is(err, auth.ErrSessionExpired):
		return status.Error(codes.Unauthenticated, "session expired")
	case errors.Is(err, auth.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, "session not found")
	case errors.Is(err, coreErrors.ErrInvalidContext):
		return status.Error(codes.Internal, "invalid context")
	case errors.Is(err, auth.ErrIncorrectCredentials):
		return status.Error(codes.Unauthenticated, "incorrect login or password")
	case errors.Is(err, auth.ErrIncorrectPassword):
		return status.Error(codes.Aborted, "incorrect password")
	case errors.Is(err, auth.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, "not authorized")
	case errors.Is(err, auth.ErrNoAccess):
		return status.Error(codes.PermissionDenied, "no access")
	case errors.Is(err, coreErrors.ErrTooManyRequests):
		return status.Error(codes.Aborted, "too many requests")
	case errors.Is(err, postgres.ErrPingDB):
		return status.Error(codes.Internal, "failed connect to database")
	case errors.Is(err, redis.ErrPingToRedis):
		return status.Error(codes.Internal, "failed connect to redis")

	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request timeout")

	default:
		return status.Error(codes.Internal, "internal UserServer error")
	}
}
