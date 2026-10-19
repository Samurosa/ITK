package mapper

import (
	"ITK_Code/m/v2/internal/core/errorcore"
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
	case errors.Is(err, errorcore.ErrOrderNotFound):
		return status.Error(codes.NotFound, "order not found")
	case errors.Is(err, errorcore.ErrSpotNotFound):
		return status.Error(codes.NotFound, "spot not found")
	case errors.Is(err, errorcore.ErrSpotInactive):
		return status.Error(codes.FailedPrecondition, "spot is not active")
	case errors.Is(err, errorcore.ErrIdempotencyConflict):
		return status.Error(codes.AlreadyExists, "idempotency key already used for another order")
	case errors.Is(err, errorcore.ErrInvalidCursor):
		return status.Error(codes.InvalidArgument, "invalid order cursor")
	case errors.Is(err, errorcore.ErrRolePermissionDenied):
		return status.Error(codes.PermissionDenied, "role is not allowed to create orders")
	case errors.Is(err, errorcore.ErrInvalidOrder):
		return status.Error(codes.InvalidArgument, "order violates spot constraints")
	case errors.Is(err, errorcore.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, "session not found")
	case errors.Is(err, errorcore.ErrOrderListenerClosed):
		return status.Error(codes.Unavailable, "order subscription unavailable")
	default:
		return status.Error(codes.Internal, "internal Order Server error")
	}
}
