package mapper

import (
	"ITK_Code/m/v2/internal/core/corerrors"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPC(err error) error {
	switch {
	case errors.Is(err, corerrors.ErrRolePermissionDenied):
		return status.Error(codes.PermissionDenied, "role is not allowed to create orders")
	case errors.Is(err, corerrors.ErrInvalidOrder):
		return status.Error(codes.InvalidArgument, "order violates spot constraints")
	default:
		return status.Error(codes.Internal, "internal Order Server error")
	}
}
