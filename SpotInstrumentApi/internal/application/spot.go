package application

import (
	"ITK_Code/m/v2/internal/core/dto"

	"context"
	"time"
)

func (s *Spot) CreateSpot(ctx context.Context, reqSpot dto.CreateSpot) (string, time.Time, error) {
	panic("implement me")
}

func (s *Spot) GetSpot(ctx context.Context, spotID string) (dto.Spot, error) {
	panic("implement me")
}

func (s *Spot) DisableSpot(ctx context.Context, spotID string) (bool, time.Time, error) {
	panic("implement me")
}

func (s *Spot) ViewMarkets(ctx context.Context, markets []string, page int32, pageSize int32) ([]dto.Market, int64, int32, error) {
	panic("implement me")
}

func (s *Spot) DescribeMarket(ctx context.Context, spotID string) (dto.DescriptionMarket, error) {
	panic("implement me")
}
