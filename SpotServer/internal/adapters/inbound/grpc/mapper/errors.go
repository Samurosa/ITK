package mapper

import (
	"ITK_Code/m/v2/internal/core/coreErrors"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPC(err error) error {
	switch {

	case errors.Is(err, coreErrors.ErrSpotNotFound):
		return status.Error(codes.NotFound, "spot not found")
	case errors.Is(err, coreErrors.ErrInvalidMinOrderGreaterMaxOrder):
		return status.Error(codes.OutOfRange, "invalid min order greater than max order")
	case errors.Is(err, coreErrors.ErrInvalidMaxOrderSize):
		return status.Error(codes.OutOfRange, "invalid max order size")
	case errors.Is(err, coreErrors.ErrInvalidMinOrderSize):
		return status.Error(codes.OutOfRange, "invalid min order size")
	case errors.Is(err, coreErrors.ErrInvalidOrderSizePrecision):
		return status.Error(codes.OutOfRange, "invalid order size precision")
	case errors.Is(err, coreErrors.ErrCompareBaseQuoteAsset):
		return status.Error(codes.InvalidArgument, "invalid compare base quote asset")

	default:
		return status.Error(codes.Internal, "internal SpotServer error")
	}
}
