package mapper

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/corerrors"
	"ITK_Code/m/v2/internal/core/user"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPC(err error) error {
	switch {
	case errors.Is(err, corerrors.ErrUserIDEmpty):
		return status.Error(codes.InvalidArgument, "the user ID in the request is empty")
	case errors.Is(err, corerrors.ErrEmailEmpty):
		return status.Error(codes.InvalidArgument, "the email in the request is empty")
	case errors.Is(err, corerrors.ErrUsernameEmpty):
		return status.Error(codes.InvalidArgument, "the username in the request is empty")
	case errors.Is(err, corerrors.ErrAssetEmpty):
		return status.Error(codes.InvalidArgument, "the asset in the request is empty")
	case errors.Is(err, corerrors.ErrInvalidAsset):
		return status.Error(codes.InvalidArgument, "the asset in the request is invalid")
	case errors.Is(err, corerrors.ErrAmountEmpty):
		return status.Error(codes.InvalidArgument, "the amount in the request is empty")
	case errors.Is(err, corerrors.ErrAmountIsZero):
		return status.Error(codes.InvalidArgument, "the amount in the request is zero value")
	case errors.Is(err, corerrors.ErrAmountIsNegative):
		return status.Error(codes.InvalidArgument, "the amount in the request is negative value")
	case errors.Is(err, corerrors.ErrPasswordEmpty):
		return status.Error(codes.InvalidArgument, "the password in the request is empty")
	case errors.Is(err, corerrors.ErrPasswordsMatch):
		return status.Error(codes.InvalidArgument, "the new password matches the old password")
	case errors.Is(err, corerrors.ErrPasswordWrongUpperSymbol):
		return status.Error(codes.InvalidArgument, "the new password wrong, upper symbol not found")
	case errors.Is(err, corerrors.ErrPasswordWrongLowerSymbol):
		return status.Error(codes.InvalidArgument, "the new password wrong, lower symbol not found")
	case errors.Is(err, corerrors.ErrPasswordWrongDigitSymbol):
		return status.Error(codes.InvalidArgument, "the new password wrong, digit symbol not found")

	case errors.Is(err, user.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, corerrors.ErrComparePassword):
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case errors.Is(err, user.ErrUpdateUser):
		return status.Error(codes.Internal, "failed to update user")
	case errors.Is(err, user.ErrEmailIsExist):
		return status.Error(codes.AlreadyExists, "email is exist")

	case errors.Is(err, corerrors.ErrRefreshExpired):
		return status.Error(codes.Unauthenticated, "refresh token expired")
	case errors.Is(err, corerrors.ErrGenerateToken):
		return status.Error(codes.Internal, "failed to generate token")
	case errors.Is(err, corerrors.ErrDeviceIDEmpty):
		return status.Error(codes.InvalidArgument, "device ID is empty")
	case errors.Is(err, corerrors.ErrInvalidToken):
		return status.Error(codes.InvalidArgument, "invalid token")
	case errors.Is(err, auth.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, "session not found")
	case errors.Is(err, corerrors.ErrInvalidContext):
		return status.Error(codes.Internal, "invalid context")
	case errors.Is(err, auth.ErrIncorrectCredentials):
		return status.Error(codes.Unauthenticated, "incorrect login or password")
	case errors.Is(err, auth.ErrIncorrectPassword):
		return status.Error(codes.Aborted, "incorrect password")
	case errors.Is(err, auth.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, "not authorized")
	case errors.Is(err, auth.ErrNoAccess):
		return status.Error(codes.PermissionDenied, "no access")
	case errors.Is(err, corerrors.ErrTooManyRequests):
		return status.Error(codes.ResourceExhausted, "too many requests")

	case errors.Is(err, corerrors.Canceled):
		return status.Error(codes.Canceled, "request canceled")
	case errors.Is(err, corerrors.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request timeout")

	default:
		return status.Error(codes.Internal, "internal UserServer error")
	}
}
