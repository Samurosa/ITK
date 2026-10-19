package mapper

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/errorcore"
	"ITK_Code/m/v2/internal/core/user"
	"errors"

	"github.com/Samurosa/exchange-common/shared/transport/grpcerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPC(err error) error {
	if mapped, ok := grpcerrors.Map(err); ok {
		return mapped
	}
	switch {
	case errors.Is(err, errorcore.ErrPasswordsMatch):
		return status.Error(codes.InvalidArgument, "the new password matches the old password")
	case errors.Is(err, errorcore.ErrPasswordWrongUpperSymbol):
		return status.Error(codes.InvalidArgument, "the new password wrong, upper symbol not found")
	case errors.Is(err, errorcore.ErrPasswordWrongLowerSymbol):
		return status.Error(codes.InvalidArgument, "the new password wrong, lower symbol not found")
	case errors.Is(err, errorcore.ErrPasswordWrongDigitSymbol):
		return status.Error(codes.InvalidArgument, "the new password wrong, digit symbol not found")

	case errors.Is(err, user.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, user.ErrEmailIsExist):
		return status.Error(codes.AlreadyExists, "email is exist")

	case errors.Is(err, errorcore.ErrGenerateToken):
		return status.Error(codes.Internal, "failed to generate token")
	case errors.Is(err, errorcore.ErrGenerateTokenProcessing):
		return status.Error(codes.Aborted, "token refresh already in progress")
	case errors.Is(err, errorcore.ErrDeviceIDEmpty):
		return status.Error(codes.InvalidArgument, "device ID is empty")
	case errors.Is(err, errorcore.ErrInvalidToken):
		return status.Error(codes.InvalidArgument, "invalid token")
	case errors.Is(err, auth.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, "session not found")
	case errors.Is(err, auth.ErrIncorrectCredentials):
		return status.Error(codes.Unauthenticated, "incorrect login or password")
	case errors.Is(err, auth.ErrIncorrectPassword):
		return status.Error(codes.Aborted, "incorrect password")
	case errors.Is(err, auth.ErrNoAccess):
		return status.Error(codes.PermissionDenied, "no access")
	case errors.Is(err, errorcore.ErrTooManyRequests):
		return status.Error(codes.ResourceExhausted, "too many requests")

	case errors.Is(err, errorcore.ErrPassGenHash):
		return status.Error(codes.Internal, "failed to generate password hash")
	case errors.Is(err, errorcore.ErrPassVerifyHash):
		return status.Error(codes.Internal, "failed to verify password hash")
	case errors.Is(err, errorcore.ErrSyncRedis):
		return status.Error(codes.Internal, "failed to acquire session lock")

	default:
		return status.Error(codes.Internal, "internal UserServer error")
	}
}
