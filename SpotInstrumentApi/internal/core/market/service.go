package market

import (
	"ITK_Code/m/v2/internal/core/dto"
	"context"
)

type Service interface {
	ViewMarkets(ctx context.Context, markets []string, page int32, pageSize int32) ([]dto.Market, int64, int32, error)
	DescribeMarket(ctx context.Context, spotID string) (dto.DescriptionMarket, error)
}
