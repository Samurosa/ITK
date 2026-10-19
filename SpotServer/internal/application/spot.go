package application

import (
	"ITK_Code/m/v2/internal/application/ports"
	"ITK_Code/m/v2/internal/application/validate"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errorcore"
	"ITK_Code/m/v2/internal/core/spot"
	"ITK_Code/m/v2/internal/core/spot/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var spotTracer = otel.Tracer("spot-service/application/spot")

type Spot struct {
	spotRepository ports.SpotRepository
}

func NewSpot(spotRepository ports.SpotRepository) *Spot {
	return &Spot{
		spotRepository: spotRepository,
	}
}

func (s *Spot) CreateSpot(ctx context.Context, reqSpot models.CreateSpot) (string, time.Time, error) {
	ctx, span := spotTracer.Start(ctx, "Spot.CreateSpot")
	defer span.End()

	log := logging.FromContext(ctx).Named("spot.create").With(
		zap.String("base_asset", reqSpot.BaseAsset),
		zap.String("quote_asset", reqSpot.QuoteAsset),
	)

	now := time.Now()

	err := validate.CreateSpot(reqSpot)
	if err != nil {
		recordSpanError(span, err)
		log.Warn("spot validation failed", zap.Error(err))
		return "", time.Time{}, err
	}
	log.Debug("spot validation passed")

	spotID, err := s.spotRepository.Save(ctx, reqSpot, now)
	if err != nil {
		logOperationError(log, "spot save failed", err)

		saveErr := fmt.Errorf("%w: %w", spot.ErrSaveSpot, err)
		recordSpanError(span, saveErr)

		return "", time.Time{}, saveErr
	}

	log.Debug("spot create request completed", zap.String("spot_id", spotID))

	return spotID, now, nil
}

func (s *Spot) GetSpot(ctx context.Context, spotID string) (dto.Spot, error) {
	ctx, span := spotTracer.Start(ctx, "Spot.GetSpot")
	defer span.End()

	log := logging.FromContext(ctx).Named("spot.get").With(zap.String("spot_id", spotID))

	gotSpot, err := s.spotRepository.Get(ctx, spotID)
	if err != nil {
		if errors.Is(err, errorcore.ErrSpotNotFound) {
			recordSpanError(span, errorcore.ErrSpotNotFound)
			log.Debug("spot not found")
			return dto.Spot{}, errorcore.ErrSpotNotFound
		}

		logOperationError(log, "spot get failed", err)

		getErr := fmt.Errorf("%w: %w", spot.ErrGetSpot, err)
		recordSpanError(span, getErr)

		return dto.Spot{}, getErr
	}
	log.Debug("spot retrieved")

	return gotSpot, nil
}

func (s *Spot) EnableSpot(ctx context.Context, spotID string) error {
	ctx, span := spotTracer.Start(ctx, "Spot.EnableSpot")
	defer span.End()

	log := logging.FromContext(ctx).Named("spot.enable").With(zap.String("spot_id", spotID))

	err := s.spotRepository.Enable(ctx, spotID)
	if errors.Is(err, errorcore.ErrSpotNotFound) {
		recordSpanError(span, errorcore.ErrSpotNotFound)
		log.Debug("spot not found")
		return errorcore.ErrSpotNotFound
	}
	if err != nil {
		logOperationError(log, "spot enable failed", err)

		enableErr := fmt.Errorf("%w: %w", spot.ErrEnableSpot, err)
		recordSpanError(span, enableErr)

		return enableErr
	}
	log.Info("spot enabled")

	return nil
}

func (s *Spot) DisableSpot(ctx context.Context, spotID string) error {
	ctx, span := spotTracer.Start(ctx, "Spot.DisableSpot")
	defer span.End()

	log := logging.FromContext(ctx).Named("spot.disable").With(zap.String("spot_id", spotID))

	err := s.spotRepository.Disable(ctx, spotID)
	if errors.Is(err, errorcore.ErrSpotNotFound) {
		recordSpanError(span, errorcore.ErrSpotNotFound)
		log.Debug("spot not found")
		return errorcore.ErrSpotNotFound
	}
	if err != nil {
		logOperationError(log, "spot disable failed", err)

		disableErr := fmt.Errorf("%w: %w", spot.ErrDisableSpot, err)
		recordSpanError(span, disableErr)

		return disableErr
	}
	log.Info("spot disabled")

	return nil
}

func (s *Spot) ListSpots(ctx context.Context, request models.ListSpotsRequest) ([]dto.SpotListItem, string, bool, error) {
	ctx, span := spotTracer.Start(ctx, "Spot.ListSpots")
	defer span.End()

	log := logging.FromContext(ctx).Named("spot.list").With(zap.Int32("page_size", request.PageSize))

	spotsList, cursor, hasMore, err := s.spotRepository.List(ctx, request)
	if err != nil {
		recordSpanError(span, err)
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
	case errors.Is(err, context.DeadlineExceeded), errors.Is(err, errorcore.ErrInvalidCursor):
		log.Warn(message, zap.Error(err))
	default:
		log.Error(message, zap.Error(err))
	}
}

func recordSpanError(span trace.Span, err error) {
	if err == nil ||
		errors.Is(err, context.Canceled) ||
		errors.Is(err, errorcore.ErrSpotNotFound) ||
		errors.Is(err, errorcore.ErrInvalidCursor) {
		return
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, "application operation failed")
}
