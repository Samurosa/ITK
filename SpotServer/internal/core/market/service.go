package market

import (
	"ITK_Code/m/v2/internal/core/dto"
	"context"

	"go.uber.org/zap"
)

type Service interface {
	ViewMarkets(ctx context.Context, log *zap.Logger, markets []string, page int32, pageSize int32) ([]dto.Market, int64, int32, error)
	DescribeMarket(ctx context.Context, log *zap.Logger, spotID string) (dto.DescriptionMarket, error)
}
