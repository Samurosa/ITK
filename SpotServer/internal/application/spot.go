package application

import (
	"ITK_Code/m/v2/internal/application/ports"
	"ITK_Code/m/v2/internal/application/validate"
	"ITK_Code/m/v2/internal/core/coreErrors"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/spot"
	"ITK_Code/m/v2/internal/core/spot/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"go.uber.org/zap"
)

type Spot struct {
	spotRepository ports.SpotRepository
}

func NewSpot(spotRepository ports.SpotRepository) *Spot {
	return &Spot{
		spotRepository: spotRepository,
	}
}

func (s *Spot) CreateSpot(ctx context.Context, reqSpot models.CreateSpot) (string, time.Time, error) {
	log := logging.FromContext(ctx).Named("spot.create").With(
		zap.String("base_asset", reqSpot.BaseAsset),
		zap.String("quote_asset", reqSpot.QuoteAsset),
	)

	now := time.Now()

	err := validate.CreateSpot(reqSpot)
	if err != nil {
		log.Warn("spot validation failed", zap.Error(err))
		return "", time.Time{}, err
	}
	log.Debug("spot validation passed")

	spotID, err := s.spotRepository.Save(ctx, reqSpot, now)
	if err != nil {
		logOperationError(log, "spot save failed", err)
		return "", time.Time{}, fmt.Errorf("%w: %w", spot.ErrSaveSpot, err)
	}

	log.Debug("spot create request completed", zap.String("spot_id", spotID))

	return spotID, now, nil
}

func (s *Spot) GetSpot(ctx context.Context, spotID string) (dto.Spot, error) {
	log := logging.FromContext(ctx).Named("spot.get").With(zap.String("spot_id", spotID))

	gotSpot, err := s.spotRepository.Get(ctx, spotID)
	if err != nil {
		if errors.Is(err, coreErrors.ErrSpotNotFound) {
			log.Debug("spot not found")
			return dto.Spot{}, coreErrors.ErrSpotNotFound
		}
		logOperationError(log, "spot get failed", err)
		return dto.Spot{}, fmt.Errorf("%w: %w", spot.ErrGetSpot, err)
	}
	log.Debug("spot retrieved")

	return gotSpot, nil
}

func (s *Spot) EnableSpot(ctx context.Context, spotID string) error {
	log := logging.FromContext(ctx).Named("spot.enable").With(zap.String("spot_id", spotID))

	err := s.spotRepository.Enable(ctx, spotID)
	if errors.Is(err, coreErrors.ErrSpotNotFound) {
		log.Debug("spot not found")
		return coreErrors.ErrSpotNotFound
	}
	if err != nil {
		logOperationError(log, "spot enable failed", err)
		return fmt.Errorf("%w: %w", spot.ErrEnableSpot, err)
	}
	log.Info("spot enabled")

	return nil
}

func (s *Spot) DisableSpot(ctx context.Context, spotID string) error {
	log := logging.FromContext(ctx).Named("spot.disable").With(zap.String("spot_id", spotID))

	err := s.spotRepository.Disable(ctx, spotID)
	if errors.Is(err, coreErrors.ErrSpotNotFound) {
		log.Debug("spot not found")
		return coreErrors.ErrSpotNotFound
	}
	if err != nil {
		logOperationError(log, "spot disable failed", err)
		return fmt.Errorf("%w: %w", spot.ErrDisableSpot, err)
	}
	log.Info("spot disabled")

	return nil
}

func (s *Spot) ListSpots(ctx context.Context, request models.ListSpotsRequest) ([]dto.SpotListItem, string, bool, error) {
	log := logging.FromContext(ctx).Named("spot.list").With(zap.Int32("page_size", request.PageSize))

	spotsList, cursor, hasMore, err := s.spotRepository.List(ctx, request)
	if err != nil {
		logOperationError(log, "spot list failed", err)
		return spotsList, "", false, err
	}
	log.Debug("spots listed", zap.Int("count", len(spotsList)), zap.Bool("has_more", hasMore))

	return spotsList, cursor, hasMore, nil
}

func logOperationError(log *zap.Logger, message string, err error) {
	switch {
	case errors.Is(err, context.Canceled):
		log.Debug(message, zap.Error(err))
	case errors.Is(err, context.DeadlineExceeded), errors.Is(err, coreErrors.ErrInvalidCursor):
		log.Warn(message, zap.Error(err))
	default:
		log.Error(message, zap.Error(err))
	}
}
