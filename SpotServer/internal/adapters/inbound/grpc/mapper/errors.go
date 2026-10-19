package mapper

import (
	"ITK_Code/m/v2/internal/core/errorcore"
	"ITK_Code/m/v2/internal/core/spot"
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
	case errors.Is(err, errorcore.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, "session not found")
	case errors.Is(err, errorcore.ErrInvalidCursor):
		return status.Error(codes.InvalidArgument, "invalid spot list cursor")
	case errors.Is(err, errorcore.ErrSpotNotFound):
		return status.Error(codes.NotFound, "spot not found")
	case errors.Is(err, errorcore.ErrSpotAlreadyExists):
		return status.Error(codes.AlreadyExists, "spot already exists")
	case errors.Is(err, errorcore.ErrInvalidMinOrderGreaterMaxOrder):
		return status.Error(codes.OutOfRange, "invalid min order greater than max order")
	case errors.Is(err, errorcore.ErrInvalidMaxOrderSize):
		return status.Error(codes.OutOfRange, "invalid max order size")
	case errors.Is(err, errorcore.ErrInvalidMinOrderSize):
		return status.Error(codes.OutOfRange, "invalid min order size")
	case errors.Is(err, errorcore.ErrInvalidOrderSizePrecision):
		return status.Error(codes.OutOfRange, "invalid order size precision")
	case errors.Is(err, errorcore.ErrCompareBaseQuoteAsset):
		return status.Error(codes.InvalidArgument, "invalid compare base quote asset")
	case errors.Is(err, spot.ErrSaveSpot):
		return status.Error(codes.Internal, "failed to save spot")
	case errors.Is(err, spot.ErrGetSpot):
		return status.Error(codes.Internal, "failed to get spot")
	case errors.Is(err, spot.ErrEnableSpot):
		return status.Error(codes.Internal, "failed to enable spot")
	case errors.Is(err, spot.ErrDisableSpot):
		return status.Error(codes.Internal, "failed to disable spot")

	default:
		return status.Error(codes.Internal, "internal SpotServer error")
	}
}
