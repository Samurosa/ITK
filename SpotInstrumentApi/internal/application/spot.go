package application

import (
	"ITK_Code/m/v2/internal/application/validate"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errors"
	"context"
	"time"

	"go.uber.org/zap"
)

func (s *Spot) CreateSpot(ctx context.Context, log *zap.Logger, reqSpot dto.CreateSpot) (string, time.Time, error) {
	log = log.Named("CreateSpot")

	err := validate.CreateSpot(log, reqSpot)
	if err != nil {
		return "", time.Time{}, err
	}
	log.Info("data validation passed")

	reqSpot.Symbol = reqSpot.Symbol + "/" + reqSpot.QuoteAsset

	spotID, err := s.spotRepository.Save(ctx, reqSpot)
	if err != nil {
		log.Error("spot save failed", zap.Error(err))
		return "", time.Time{}, errors.ErrSaveSpot
	}
	log.Info("spot saved", zap.String("id", spotID))

	return spotID, time.Now(), nil
}

func (s *Spot) GetSpot(ctx context.Context, log *zap.Logger, spotID string) (dto.Spot, error) {
	panic("implement me")
}

func (s *Spot) EnableSpot(ctx context.Context, log *zap.Logger, spotID string) (bool, time.Time, error) {
	log.Named("EnableSpot")

	err := s.spotRepository.Enable(ctx, spotID)
	if err != nil {
		log.Error("spot enable failed", zap.Error(err))
		return false, time.Time{}, errors.ErrEnableSpot
	}
	log.Info("spot enable", zap.String("id", spotID))

	return true, time.Now(), nil
}

func (s *Spot) DisableSpot(ctx context.Context, log *zap.Logger, spotID string) (bool, time.Time, error) {
	log.Named("DisableSpot")

	err := s.spotRepository.Disable(ctx, spotID)
	if err != nil {
		log.Error("spot disable failed", zap.Error(err))
		return false, time.Time{}, errors.ErrDisableSpot
	}
	log.Info("spot disabled", zap.String("id", spotID))

	return true, time.Now(), nil
}
