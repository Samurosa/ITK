package mapper

import (
	"ITK_Code/m/v2/internal/core/corerrors"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPC(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	case errors.Is(err, corerrors.ErrOrderNotFound):
		return status.Error(codes.NotFound, "order not found")
	case errors.Is(err, corerrors.ErrSpotNotFound):
		return status.Error(codes.NotFound, "spot not found")
	case errors.Is(err, corerrors.ErrSpotInactive):
		return status.Error(codes.FailedPrecondition, "spot is not active")
	case errors.Is(err, corerrors.ErrIdempotencyConflict):
		return status.Error(codes.AlreadyExists, "idempotency key already used for another order")
	case errors.Is(err, corerrors.ErrInvalidCursor):
		return status.Error(codes.InvalidArgument, "invalid order cursor")
	case errors.Is(err, corerrors.ErrRolePermissionDenied):
		return status.Error(codes.PermissionDenied, "role is not allowed to create orders")
	case errors.Is(err, corerrors.ErrInvalidOrder):
		return status.Error(codes.InvalidArgument, "order violates spot constraints")
	default:
		switch status.Code(err) {
		case codes.Canceled:
			return status.Error(codes.Canceled, "request canceled")
		case codes.DeadlineExceeded:
			return status.Error(codes.DeadlineExceeded, "dependency deadline exceeded")
		case codes.Unavailable:
			return status.Error(codes.Unavailable, "dependency unavailable")
		}
		return status.Error(codes.Internal, "internal Order Server error")
	}
}
