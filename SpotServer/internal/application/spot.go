package application

import (
	"ITK_Code/m/v2/internal/application/validate"
	"ITK_Code/m/v2/internal/core/dto"
	errorsCore "ITK_Code/m/v2/internal/core/errors"
	"ITK_Code/m/v2/internal/core/spot"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

func (s *Spot) CreateSpot(ctx context.Context, reqSpot dto.CreateSpot) (string, time.Time, error) {
	log := s.log.Named("Create spot")

	err := validate.CreateSpot(log, reqSpot)
	if err != nil {
		return "", time.Time{}, err
	}
	log.Info("data validation passed")

	spotID, err := s.spotRepository.Save(ctx, reqSpot)
	if err != nil {
		log.Error("spot save failed", zap.Error(err))
		return "", time.Time{}, spot.ErrSaveSpot
	}
	log.Info("spot saved", zap.String("id", spotID))

	return spotID, time.Now(), nil
}

func (s *Spot) GetSpot(ctx context.Context, spotID string) (dto.Spot, error) {
	log := s.log.Named("Get spot")

	gotSpot, err := s.spotRepository.Get(ctx, spotID)
	if err != nil {
		log.Error("spot get failed", zap.Error(err))
		return dto.Spot{}, spot.ErrGetSpot
	}
	log.Info("got spot", zap.String("id", spotID))

	return gotSpot, nil
}

func (s *Spot) EnableSpot(ctx context.Context, spotID string) error {
	log := s.log.Named("Enable spot")

	err := s.spotRepository.Enable(ctx, spotID)
	if errors.Is(err, errorsCore.ErrSpotNotFound) {
		log.Error("spot not found", zap.String("id", spotID))
		return errorsCore.ErrSpotNotFound
	}
	if err != nil {
		log.Error("spot enable failed", zap.Error(err))
		return spot.ErrEnableSpot
	}
	log.Info("spot enable", zap.String("id", spotID))

	return nil
}

func (s *Spot) DisableSpot(ctx context.Context, spotID string) error {
	log := s.log.Named("Disable spot")

	err := s.spotRepository.Disable(ctx, spotID)
	if errors.Is(err, errorsCore.ErrSpotNotFound) {
		log.Error("spot not found", zap.String("id", spotID))
		return errorsCore.ErrSpotNotFound
	}
	if err != nil {
		log.Error("spot disable failed", zap.Error(err))
		return spot.ErrDisableSpot
	}
	log.Info("spot disabled", zap.String("id", spotID))

	return nil
}
